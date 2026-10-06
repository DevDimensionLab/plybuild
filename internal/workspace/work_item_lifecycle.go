package workspace

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type LifecycleState string

const (
	LifecycleActive               LifecycleState = "active"
	LifecycleParked               LifecycleState = "parked"
	LifecycleFrozen               LifecycleState = "frozen"
	LifecycleArchived             LifecycleState = "archived"
	WorkItemLifecycleFile                        = "work-item-lifecycle.json"
	WorkItemLifecycleReadbackKind                = "WorkspaceWorkItemLifecycleReadback@1"
	WorkItemLifecycleMutationKind                = "WorkspaceWorkItemLifecycleMutation@1"
)

func (state LifecycleState) Validate() error {
	switch state {
	case LifecycleActive, LifecycleParked, LifecycleFrozen, LifecycleArchived:
		return nil
	default:
		return WorkInvalidArguments("lifecycle must be active, parked, frozen, or archived")
	}
}

type WorkItemLifecycleEvent struct {
	ID            string         `json:"event_id"`
	Sequence      int            `json:"sequence"`
	SubjectKind   string         `json:"subject_kind"`
	SubjectID     string         `json:"subject_id"`
	From          LifecycleState `json:"from"`
	To            LifecycleState `json:"to"`
	ActorClaim    string         `json:"actor_claim"`
	RecordedAtUTC string         `json:"recorded_at_utc"`
}

// WorkItemLifecycleSnapshot is an immutable, request-local read of the organizing
// metadata. A missing file is valid: Revision is zero and SHA256 is empty.
type WorkItemLifecycleSnapshot struct {
	Revision int
	Events   []WorkItemLifecycleEvent
	Locator  string
	SHA256   string
	Exists   bool
	states   map[string]LifecycleState
}

func (s WorkItemLifecycleSnapshot) State(kind, id string) LifecycleState {
	if state, ok := s.states[kind+"/"+id]; ok {
		return state
	}
	return LifecycleActive
}
func (s WorkItemLifecycleSnapshot) Task(id TaskID) LifecycleState { return s.State("task", string(id)) }
func (s WorkItemLifecycleSnapshot) Epic(id EpicID) LifecycleState { return s.State("epic", string(id)) }

type WorkItemLifecycleInput struct {
	SubjectKind      string
	SubjectID        string
	State            LifecycleState
	ActorClaim       string
	ExpectedRevision *int
}

func (input WorkItemLifecycleInput) Validate() error {
	if err := validateLifecycleSubject(input.SubjectKind, input.SubjectID); err != nil {
		return err
	}
	if err := input.State.Validate(); err != nil {
		return err
	}
	if !validTaskText(input.ActorClaim, 1, 256) {
		return WorkInvalidArguments("--actor must contain 1..256 characters without surrounding whitespace or control characters")
	}
	if input.ExpectedRevision != nil && *input.ExpectedRevision < 0 {
		return WorkInvalidArguments("--expected-revision must be zero or greater")
	}
	return nil
}

func validateLifecycleSubject(kind, id string) error {
	if kind != "task" && kind != "epic" {
		return WorkInvalidArguments("lifecycle subject must be task or epic")
	}
	return validateWorkIdentifier(kind, id)
}

type WorkItemLifecycleReadback struct {
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	Workspace     MetadataWorkspace        `json:"workspace"`
	SubjectKind   string                   `json:"subject_kind"`
	SubjectID     string                   `json:"subject_id"`
	Lifecycle     LifecycleState           `json:"lifecycle"`
	Revision      int                      `json:"revision"`
	Events        []WorkItemLifecycleEvent `json:"events"`
	Source        MetadataSource           `json:"source"`
}

type WorkItemLifecycleMutation struct {
	WorkItemLifecycleReadback
	Changed bool                    `json:"changed"`
	Event   *WorkItemLifecycleEvent `json:"event"`
}

func lifecycleReadback(root, kind, id string, snapshot WorkItemLifecycleSnapshot) WorkItemLifecycleReadback {
	events := []WorkItemLifecycleEvent{}
	for _, event := range snapshot.Events {
		if event.SubjectKind == kind && event.SubjectID == id {
			events = append(events, event)
		}
	}
	return WorkItemLifecycleReadback{WorkItemLifecycleReadbackKind, 1, MetadataWorkspace{root}, kind, id, snapshot.State(kind, id), snapshot.Revision, events, metadataSource(snapshot.Locator, snapshot.SHA256, snapshot.Exists)}
}

// ReadWorkItemLifecycle does not inspect Git, acquire locks, or create files.
// Subject membership belongs to the caller's registry snapshot.
func ReadWorkItemLifecycle(d Dependencies, root string) (WorkItemLifecycleSnapshot, error) {
	observed, err := ObserveRoot(d, root)
	if err != nil {
		return WorkItemLifecycleSnapshot{}, err
	}
	return readWorkItemLifecycle(observed.Root)
}

func ShowWorkItemLifecycle(d Dependencies, kind, id string) (WorkItemLifecycleReadback, error) {
	if err := validateLifecycleSubject(kind, id); err != nil {
		return WorkItemLifecycleReadback{}, err
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return WorkItemLifecycleReadback{}, err
	}
	if err := registeredLifecycleSubject(d, root, kind, id); err != nil {
		return WorkItemLifecycleReadback{}, err
	}
	snapshot, err := readWorkItemLifecycle(root)
	if err != nil {
		return WorkItemLifecycleReadback{}, err
	}
	return lifecycleReadback(root, kind, id, snapshot), nil
}

func registeredLifecycleSubject(d Dependencies, root, kind, id string) error {
	if d.WorkItems == nil {
		return workError(ErrorWorkIO, "work-item store dependency is required", nil)
	}
	r, err := d.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		return err
	}
	if kind == "task" {
		if found, _ := findTask(r, TaskID(id)); found != nil {
			return nil
		}
	}
	if kind == "epic" {
		if found, _ := findEpic(r, EpicID(id)); found != nil {
			return nil
		}
	}
	return workError(ErrorWorkNotFound, fmt.Sprintf("%s %s is not registered", kind, id), nil)
}

func SetWorkItemLifecycle(d Dependencies, input WorkItemLifecycleInput) (WorkItemLifecycleMutation, error) {
	return setWorkItemLifecycle(d, input, workspaceMetadataPublisher{})
}

func setWorkItemLifecycle(d Dependencies, input WorkItemLifecycleInput, publisher workspaceMetadataPublisher) (WorkItemLifecycleMutation, error) {
	var result WorkItemLifecycleMutation
	if err := input.Validate(); err != nil {
		return result, err
	}
	if d.WorkItems == nil || d.WorkClock == nil {
		return result, workError(ErrorWorkIO, "work-item store and clock dependencies are required", nil)
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return result, err
	}
	err = d.WorkItems.WithLock(root, func(_ WorkItemStoreSession) error {
		if err := registeredLifecycleSubject(d, root, input.SubjectKind, input.SubjectID); err != nil {
			return err
		}
		snapshot, err := readWorkItemLifecycle(root)
		if err != nil {
			return err
		}
		result.WorkItemLifecycleReadback = lifecycleReadback(root, input.SubjectKind, input.SubjectID, snapshot)
		result.Kind = WorkItemLifecycleMutationKind
		if input.ExpectedRevision != nil && *input.ExpectedRevision != snapshot.Revision {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("lifecycle revision is %d, expected %d; read the current lifecycle before retrying", snapshot.Revision, *input.ExpectedRevision), nil)
		}
		if result.Lifecycle == input.State {
			return nil
		}
		event := WorkItemLifecycleEvent{Sequence: snapshot.Revision + 1, SubjectKind: input.SubjectKind, SubjectID: input.SubjectID, From: result.Lifecycle, To: input.State, ActorClaim: input.ActorClaim, RecordedAtUTC: d.WorkClock.Now().UTC().Format(time.RFC3339Nano)}
		event.ID = lifecycleEventID(event)
		events := append(append([]WorkItemLifecycleEvent{}, snapshot.Events...), event)
		bytes, err := encodeWorkItemLifecycle(events)
		if err != nil {
			return err
		}
		if err := publisher.publish(root, WorkItemLifecycleFile, bytes); err != nil {
			return err
		}
		snapshot, err = readWorkItemLifecycle(root)
		if err != nil {
			return err
		}
		result.WorkItemLifecycleReadback = lifecycleReadback(root, input.SubjectKind, input.SubjectID, snapshot)
		result.Kind = WorkItemLifecycleMutationKind
		result.Changed, result.Event = true, &event
		return nil
	})
	return result, err
}

func lifecycleEventID(event WorkItemLifecycleEvent) string {
	event.ID = ""
	bytes, _ := json.Marshal(event)
	return "lce_" + strings.TrimPrefix(digestTaskBytes(bytes), "sha256:")
}

func lifecycleLocator(root string) string {
	return filepath.Join(root, MarkerDirectory, WorkItemLifecycleFile)
}
