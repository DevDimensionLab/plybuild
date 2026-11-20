package taskwork

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"strings"
	"time"
)

type LifecycleState string

const (
	LifecycleActive               LifecycleState = "active"
	LifecycleParked               LifecycleState = "parked"
	LifecycleFrozen               LifecycleState = "frozen"
	LifecycleArchived             LifecycleState = "archived"
	LifecycleCompleted            LifecycleState = "completed"
	WorkItemLifecycleReadbackKind                = "WorkspaceWorkItemLifecycleReadback@1"
	WorkItemLifecycleMutationKind                = "WorkspaceWorkItemLifecycleMutation@1"
)

func (state LifecycleState) Validate() error {
	switch state {
	case LifecycleActive, LifecycleParked, LifecycleFrozen, LifecycleArchived, LifecycleCompleted:
		return nil
	default:
		return workspace.WorkInvalidArguments("lifecycle must be active, parked, frozen, archived, or completed")
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

// WorkItemLifecycleSnapshot is an immutable, request-local read of one
// Project's organizing metadata. A missing file is valid: Revision is zero and
// SHA256 is empty.
type WorkItemLifecycleSnapshot struct {
	Project  workspace.ProjectID
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

// WorkItemLifecycles is the lifecycle of every defined Project, sorted by
// Project ID. Epic and Task IDs are unique in the workspace, so a subject's
// state comes from whichever Project records it.
type WorkItemLifecycles []WorkItemLifecycleSnapshot

func (l WorkItemLifecycles) State(kind, id string) LifecycleState {
	for _, s := range l {
		if state, ok := s.states[kind+"/"+id]; ok {
			return state
		}
	}
	return LifecycleActive
}
func (l WorkItemLifecycles) Task(id TaskID) LifecycleState { return l.State("task", string(id)) }
func (l WorkItemLifecycles) Epic(id EpicID) LifecycleState { return l.State("epic", string(id)) }

// Events are the lifecycle events of all Projects, Project by Project.
func (l WorkItemLifecycles) Events() []WorkItemLifecycleEvent {
	out := []WorkItemLifecycleEvent{}
	for _, s := range l {
		out = append(out, s.Events...)
	}
	return out
}

// Of is the lifecycle of project; a Project without one reads as empty.
func (l WorkItemLifecycles) Of(project workspace.ProjectID) WorkItemLifecycleSnapshot {
	for _, s := range l {
		if s.Project == project {
			return s
		}
	}
	return WorkItemLifecycleSnapshot{Project: project, Events: []WorkItemLifecycleEvent{}, states: map[string]LifecycleState{}}
}

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
		return workspace.WorkInvalidArguments("--actor must contain 1..256 characters without surrounding whitespace or control characters")
	}
	if input.ExpectedRevision != nil && *input.ExpectedRevision < 0 {
		return workspace.WorkInvalidArguments("--expected-revision must be zero or greater")
	}
	return nil
}

func validateLifecycleSubject(kind, id string) error {
	if kind != "task" && kind != "epic" {
		return workspace.WorkInvalidArguments("lifecycle subject must be task or epic")
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

// ReadWorkItemLifecycle reads the lifecycle of every defined Project. It does
// not inspect Git, acquire locks, or create files. Subject membership belongs
// to the caller's registry snapshot.
func ReadWorkItemLifecycle(d Dependencies, root string) (WorkItemLifecycles, error) {
	observed, err := workspace.ObserveRoot(d.Dependencies, root)
	if err != nil {
		return nil, err
	}
	works, err := workspace.DefinedWork(observed.Root)
	if err != nil {
		return nil, err
	}
	return readLifecycles(works)
}

func readLifecycles(works []workspace.ProjectWork) (WorkItemLifecycles, error) {
	out := make(WorkItemLifecycles, 0, len(works))
	for _, work := range works {
		snapshot, err := readWorkItemLifecycle(work)
		if err != nil {
			return nil, err
		}
		out = append(out, snapshot)
	}
	return out, nil
}

// ReadProjectLifecycle reads the lifecycle of the Project of work only.
func ReadProjectLifecycle(work workspace.ProjectWork) (WorkItemLifecycleSnapshot, error) {
	return readWorkItemLifecycle(work)
}

func ShowWorkItemLifecycle(d Dependencies, kind, id string) (WorkItemLifecycleReadback, error) {
	if err := validateLifecycleSubject(kind, id); err != nil {
		return WorkItemLifecycleReadback{}, err
	}
	root, err := containingRegistrationWorkspace(d)
	if err != nil {
		return WorkItemLifecycleReadback{}, err
	}
	work, err := registeredLifecycleSubject(d, root, kind, id)
	if err != nil {
		return WorkItemLifecycleReadback{}, err
	}
	snapshot, err := readWorkItemLifecycle(work)
	if err != nil {
		return WorkItemLifecycleReadback{}, err
	}
	return lifecycleReadback(root, kind, id, snapshot), nil
}

// registeredLifecycleSubject is the work of the Project whose registry holds
// the Epic or Task.
func registeredLifecycleSubject(d Dependencies, root, kind, id string) (workspace.ProjectWork, error) {
	if d.WorkItems == nil {
		return workspace.ProjectWork{}, workspace.WorkError(workspace.ErrorWorkIO, "work-item store dependency is required", nil)
	}
	switch kind {
	case "task":
		return LocateTask(d, root, TaskID(id))
	case "epic":
		return LocateEpic(d, root, EpicID(id))
	}
	return workspace.ProjectWork{}, workspace.WorkError(workspace.ErrorWorkNotFound, fmt.Sprintf("%s %s is not registered", kind, id), nil)
}

func SetWorkItemLifecycle(d Dependencies, input WorkItemLifecycleInput) (WorkItemLifecycleMutation, error) {
	return setWorkItemLifecycle(d, input, workspaceMetadataPublisher{})
}

func setWorkItemLifecycle(d Dependencies, input WorkItemLifecycleInput, publisher workspaceMetadataPublisher) (WorkItemLifecycleMutation, error) {
	var result WorkItemLifecycleMutation
	if err := input.Validate(); err != nil {
		return result, err
	}
	if input.State == LifecycleCompleted {
		return result, workspace.WorkInvalidArguments("completed is recorded only by native integration closeout")
	}
	if d.WorkItems == nil || d.WorkClock == nil {
		return result, workspace.WorkError(workspace.ErrorWorkIO, "work-item store and clock dependencies are required", nil)
	}
	root, err := containingRegistrationWorkspace(d)
	if err != nil {
		return result, err
	}
	// Each Project's lifecycle file lives in its work, written under that
	// Project's work lock. A subject never moves to another Project.
	work, err := registeredLifecycleSubject(d, root, input.SubjectKind, input.SubjectID)
	if err != nil {
		return result, err
	}
	err = d.WorkItems.WithLock(work, func(WorkItemStoreSession) error { return publishLifecycle(d, work, input, publisher, &result) })
	return result, err
}

func publishLifecycle(d Dependencies, work workspace.ProjectWork, input WorkItemLifecycleInput, publisher workspaceMetadataPublisher, result *WorkItemLifecycleMutation) error {
	root := work.Workspace
	snapshot, err := readWorkItemLifecycle(work)
	if err != nil {
		return err
	}
	result.WorkItemLifecycleReadback = lifecycleReadback(root, input.SubjectKind, input.SubjectID, snapshot)
	result.Kind = WorkItemLifecycleMutationKind
	if input.ExpectedRevision != nil && *input.ExpectedRevision != snapshot.Revision {
		return workspace.WorkError(workspace.ErrorWorkIdentityConflict, fmt.Sprintf("lifecycle revision is %d, expected %d; read the current lifecycle before retrying", snapshot.Revision, *input.ExpectedRevision), nil)
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
	if err := work.Create(); err != nil {
		return err
	}
	if err := publisher.publish(work.Dir, workspace.WorkItemLifecycleFile, bytes); err != nil {
		return err
	}
	snapshot, err = readWorkItemLifecycle(work)
	if err != nil {
		return err
	}
	result.WorkItemLifecycleReadback = lifecycleReadback(root, input.SubjectKind, input.SubjectID, snapshot)
	result.Kind = WorkItemLifecycleMutationKind
	result.Changed, result.Event = true, &event
	return nil
}

func lifecycleEventID(event WorkItemLifecycleEvent) string {
	event.ID = ""
	bytes, _ := json.Marshal(event)
	return "lce_" + strings.TrimPrefix(digestTaskBytes(bytes), "sha256:")
}
