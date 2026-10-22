package workspace

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const ExternalIssuesFile = "external-issues.json"
const externalIssuesStoreKind = "WorkspaceExternalIssuesStore@1"

type ExternalIssueEvent struct {
	ID            string         `json:"event_id"`
	Sequence      int            `json:"sequence"`
	SubjectKind   string         `json:"subject_kind"`
	SubjectID     string         `json:"subject_id"`
	Operation     string         `json:"operation"`
	Issue         *ExternalIssue `json:"issue"`
	RefID         string         `json:"ref_id"`
	Primary       bool           `json:"primary"`
	ActorClaim    string         `json:"actor_claim"`
	RecordedAtUTC string         `json:"recorded_at_utc"`
}
type externalIssuesStore struct {
	Kind          string               `json:"kind"`
	SchemaVersion int                  `json:"schema_version"`
	Events        []ExternalIssueEvent `json:"events"`
}
type externalIssueState struct {
	refs    map[string]ExternalIssue
	primary string
}
type externalIssuesSnapshot struct {
	store  externalIssuesStore
	items  map[string]*externalIssueState
	owners map[string]string
	source MetadataSource
}

func emptyExternalIssues(root string) externalIssuesSnapshot {
	return externalIssuesSnapshot{store: externalIssuesStore{externalIssuesStoreKind, 1, []ExternalIssueEvent{}}, items: map[string]*externalIssueState{}, owners: map[string]string{}, source: metadataSource(filepath.Join(root, MarkerDirectory, ExternalIssuesFile), "", false)}
}
func externalEventID(e ExternalIssueEvent) string {
	e.ID = ""
	b, _ := json.Marshal(e)
	return "exe_" + strings.TrimPrefix(digestTaskBytes(b), "sha256:")
}
func (s *externalIssuesSnapshot) apply(e ExternalIssueEvent) (bool, error) {
	if err := validateLifecycleSubject(e.SubjectKind, e.SubjectID); err != nil {
		return false, err
	}
	key := e.SubjectKind + "/" + e.SubjectID
	item := s.items[key]
	if item == nil {
		item = &externalIssueState{refs: map[string]ExternalIssue{}}
		s.items[key] = item
	}
	conflict := func(detail string) (bool, error) {
		return false, workError(ErrorWorkIdentityConflict, detail+"; read current external-issue metadata and retry the intended request", nil)
	}
	switch e.Operation {
	case ExternalIssueAdd:
		if e.Issue == nil || e.RefID != "" {
			return false, WorkInvalidArguments("add requires issue fields and no ref ID")
		}
		issue, err := NormalizeExternalIssue(*e.Issue)
		if err != nil {
			return false, err
		}
		if issue != *e.Issue {
			return false, WorkInvalidArguments("stored issue is not normalized")
		}
		if owner, ok := s.owners[issue.ID]; ok && owner != key {
			return conflict(fmt.Sprintf("external issue %s is already linked to %s; explicitly unlink it there before moving it", issue.ID, owner))
		}
		previous, exists := item.refs[issue.ID]
		if exists && previous != issue {
			return conflict("same external identity has conflicting supplied metadata")
		}
		changed := !exists || e.Primary && item.primary != issue.ID
		item.refs[issue.ID] = issue
		s.owners[issue.ID] = key
		if e.Primary {
			item.primary = issue.ID
		}
		return changed, nil
	case ExternalIssueUnlink:
		if e.Issue != nil || e.Primary || e.RefID == "" {
			return false, WorkInvalidArguments("unlink requires only the selected ref ID")
		}
		if _, exists := item.refs[e.RefID]; !exists {
			return conflict("selected external ref is not linked to this work item")
		}
		delete(item.refs, e.RefID)
		delete(s.owners, e.RefID)
		if item.primary == e.RefID {
			item.primary = ""
		}
		return true, nil
	case ExternalIssuePrimary:
		if e.Issue != nil || e.Primary {
			return false, WorkInvalidArguments("primary takes only a ref ID, or an empty ID to clear")
		}
		if e.RefID != "" {
			if _, exists := item.refs[e.RefID]; !exists {
				return conflict("primary must select a ref linked to this work item")
			}
		}
		changed := item.primary != e.RefID
		item.primary = e.RefID
		return changed, nil
	default:
		return false, WorkInvalidArguments("unknown external issue operation")
	}
}
func readExternalIssues(root string) (externalIssuesSnapshot, error) {
	s := emptyExternalIssues(root)
	b, exists, err := readWorkspaceMetadata(root, ExternalIssuesFile)
	if err != nil {
		return s, err
	}
	if !exists {
		return s, nil
	}
	var store externalIssuesStore
	if err := decodeWorkspaceMetadata(b, &store); err != nil {
		return s, metadataStoreError(s.source.Locator, err)
	}
	if store.Kind != externalIssuesStoreKind || store.SchemaVersion != 1 || store.Events == nil {
		return s, metadataStoreError(s.source.Locator, fmt.Errorf("expected external issue store version 1 with events array"))
	}
	for i, e := range store.Events {
		if e.Sequence != i+1 || e.ID != externalEventID(e) || !validTaskText(e.ActorClaim, 1, 256) || !validTaskUTC(e.RecordedAtUTC) {
			return s, metadataStoreError(s.source.Locator, fmt.Errorf("invalid external issue event %d", i+1))
		}
		changed, err := s.apply(e)
		if err != nil || !changed {
			return s, metadataStoreError(s.source.Locator, fmt.Errorf("invalid external issue event %d: %v", i+1, err))
		}
	}
	s.store = store
	s.source = metadataSource(s.source.Locator, digestTaskBytes(b), true)
	return s, nil
}
func (s externalIssuesSnapshot) reference(kind, id string) *ExternalCommitReference {
	item := s.items[kind+"/"+id]
	if item == nil || item.primary == "" {
		return nil
	}
	return &ExternalCommitReference{item.primary, item.refs[item.primary].CommitKey(), kind, id}
}
func (s externalIssuesSnapshot) readback(root, kind, id string, r WorkItemRegistry) ExternalIssuesReadback {
	out := ExternalIssuesReadback{Kind: "WorkspaceExternalIssuesReadback@1", SchemaVersion: 1, Workspace: MetadataWorkspace{root}, SubjectKind: kind, SubjectID: id, ProviderStatus: "offline_unverified", Refs: []ExternalIssue{}, Events: []ExternalIssueEvent{}, Revision: len(s.store.Events), Source: s.source}
	if item := s.items[kind+"/"+id]; item != nil {
		for _, ref := range item.refs {
			out.Refs = append(out.Refs, ref)
		}
	}
	sort.Slice(out.Refs, func(i, j int) bool { return out.Refs[i].ID < out.Refs[j].ID })
	out.Primary = s.reference(kind, id)
	out.Effective = out.Primary
	if kind == "task" && out.Effective == nil {
		if task, _ := findTask(r, TaskID(id)); task != nil {
			out.Effective = s.reference("epic", string(task.ParentEpicID))
		}
	}
	for _, e := range s.store.Events {
		if e.SubjectKind == kind && e.SubjectID == id {
			out.Events = append(out.Events, e)
		}
	}
	return out
}
func externalRegistered(r WorkItemRegistry, kind, id string) error {
	if kind == "task" {
		if item, _ := findTask(r, TaskID(id)); item != nil {
			return nil
		}
	}
	if kind == "epic" {
		if item, _ := findEpic(r, EpicID(id)); item != nil {
			return nil
		}
	}
	return workError(ErrorWorkNotFound, fmt.Sprintf("%s %s is not registered", kind, id), nil)
}
func ShowExternalIssues(d Dependencies, kind, id string) (ExternalIssuesReadback, error) {
	var out ExternalIssuesReadback
	if err := validateLifecycleSubject(kind, id); err != nil {
		return out, err
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	if d.WorkItems == nil {
		return out, workError(ErrorWorkIO, "work-item store is required", nil)
	}
	r, err := d.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		return out, err
	}
	if err := externalRegistered(r, kind, id); err != nil {
		return out, err
	}
	s, err := readExternalIssues(root)
	if err != nil {
		return out, err
	}
	return s.readback(root, kind, id, r), nil
}
func UpdateExternalIssues(d Dependencies, in ExternalIssueInput) (ExternalIssuesMutation, error) {
	return updateExternalIssues(d, in, workspaceMetadataPublisher{})
}
func updateExternalIssues(d Dependencies, in ExternalIssueInput, publisher workspaceMetadataPublisher) (ExternalIssuesMutation, error) {
	var out ExternalIssuesMutation
	if err := validateLifecycleSubject(in.SubjectKind, in.SubjectID); err != nil {
		return out, err
	}
	if !validTaskText(in.ActorClaim, 1, 256) {
		return out, WorkInvalidArguments("--actor must contain 1..256 characters without surrounding whitespace or controls")
	}
	if in.ExpectedRevision != nil && *in.ExpectedRevision < 0 {
		return out, WorkInvalidArguments("expected revision must be zero or greater")
	}
	if (in.Operation == ExternalIssuePrimary || in.Operation == ExternalIssueUnlink || in.Primary) && in.ExpectedRevision == nil {
		return out, WorkInvalidArguments("primary changes and unlink require --expected-revision from current external-issue readback")
	}
	if in.Issue != nil {
		normalized, err := NormalizeExternalIssue(*in.Issue)
		if err != nil {
			return out, err
		}
		in.Issue = &normalized
	}
	if d.WorkItems == nil || d.WorkClock == nil {
		return out, workError(ErrorWorkIO, "work-item store and clock are required", nil)
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	err = d.WorkItems.WithLock(root, func(_ WorkItemStoreSession) error {
		// Registration-only read avoids source/Git/content observation for metadata.
		r, err := d.WorkItems.SnapshotRegistrations(root)
		if err != nil {
			return err
		}
		if err := externalRegistered(r, in.SubjectKind, in.SubjectID); err != nil {
			return err
		}
		s, err := readExternalIssues(root)
		if err != nil {
			return err
		}
		out.ExternalIssuesReadback = s.readback(root, in.SubjectKind, in.SubjectID, r)
		e := ExternalIssueEvent{Sequence: out.Revision + 1, SubjectKind: in.SubjectKind, SubjectID: in.SubjectID, Operation: in.Operation, Issue: in.Issue, RefID: in.RefID, Primary: in.Primary, ActorClaim: in.ActorClaim, RecordedAtUTC: d.WorkClock.Now().UTC().Format(time.RFC3339Nano)}
		e.ID = externalEventID(e)
		changed, err := s.apply(e)
		if err != nil {
			return err
		}
		// A lost-response retry may carry the preceding revision. Only an actual
		// no-effect request can bypass stale CAS; a newer primary stays protected.
		if !changed {
			return nil
		}
		if in.ExpectedRevision != nil && *in.ExpectedRevision != out.Revision {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("external issue revision is %d, expected %d; read current external-issue metadata before retrying", out.Revision, *in.ExpectedRevision), nil)
		}
		s.store.Events = append(s.store.Events, e)
		b, err := json.MarshalIndent(s.store, "", "  ")
		if err != nil {
			return err
		}
		b = append(b, '\n')
		if err := publisher.publish(root, ExternalIssuesFile, b); err != nil {
			return err
		}
		s, err = readExternalIssues(root)
		if err != nil {
			return err
		}
		out.ExternalIssuesReadback = s.readback(root, in.SubjectKind, in.SubjectID, r)
		out.Changed = true
		return nil
	})
	return out, err
}
