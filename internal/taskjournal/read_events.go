package taskjournal

import (
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// EventBatchInput binds already-read registry snapshots, one per Project. This
// event reader projects history, not current candidate validity, and never
// probes Git.
type EventBatchInput struct {
	Workspace  workspace.WorkspaceObservation
	Projects   workspace.ProjectSnapshot
	Registries []taskwork.ProjectRegistry
	Lifecycles taskwork.WorkItemLifecycles
}

type EventBatch struct {
	Tasks       map[taskwork.TaskID]Snapshot
	EpicEvents  map[taskwork.EpicID][]Event
	EpicSources map[taskwork.EpicID][]Source
	Reasons     []Reason
}

type sourceRead struct {
	kind string
	hash string
	err  error
}
type bytesRead struct {
	bytes []byte
	err   error
}
type eventBatchReader struct {
	lifecycles taskwork.WorkItemLifecycles
	sources    map[string]sourceRead
}

func journalSourceHash(root, kind, path string) (string, error) {
	if kind == taskwork.IntegrationEvidenceSourceKind {
		digest, err := taskwork.IntegrationEvidenceSourceDigest(root, path)
		return strings.TrimPrefix(digest, "sha256:"), err
	}
	return sourceHash(path)
}

// Generic observations cannot satisfy or weaken managed evidence checks.
func cachedSourceHash(cache map[string]sourceRead, root, kind, path string) (string, error) {
	if cache != nil {
		if r, ok := cache[path]; ok && (r.kind == kind || r.kind == taskwork.IntegrationEvidenceSourceKind) {
			return r.hash, r.err
		}
	}
	h, e := journalSourceHash(root, kind, path)
	if cache != nil {
		cache[path] = sourceRead{kind: kind, hash: h, err: e}
	}
	return h, e
}

func ReadEventBatch(d taskwork.Dependencies, in EventBatchInput) (EventBatch, error) {
	return New(d).readEventBatch(in)
}

func (s Service) readEventBatch(in EventBatchInput) (EventBatch, error) {
	out := EventBatch{Tasks: map[taskwork.TaskID]Snapshot{}, EpicEvents: map[taskwork.EpicID][]Event{}, EpicSources: map[taskwork.EpicID][]Source{}, Reasons: []Reason{}}
	reader := &eventBatchReader{lifecycles: in.Lifecycles, sources: map[string]sourceRead{}}
	s.readBatch = reader
	for _, owner := range in.Registries {
		if err := s.readProjectEvents(in, reader, owner, &out); err != nil {
			return out, err
		}
	}
	// A single recheck per distinct source detects concurrent changes. This is
	// explicitly not a cross-register transaction or a filesystem snapshot.
	if err := s.check("before_freshness"); err != nil {
		return out, err
	}
	changed := map[string]bool{}
	for path, before := range reader.sources {
		if before.err != nil {
			continue
		}
		after, err := journalSourceHash(in.Workspace.Root, before.kind, path)
		if err != nil || after != before.hash {
			changed[path] = true
		}
	}
	for id, snap := range out.Tasks {
		for i := range snap.Sources {
			if changed[snap.Sources[i].Locator] {
				snap.Sources[i].Status = "changed"
				snap.reason("source_changed", "Source changed during batch reading", snap.Sources[i].SourceID)
			}
		}
		snap.deriveAxes()
		snap.SnapshotID = snapshotDigest(snap)
		out.Tasks[id] = snap
	}
	if len(changed) > 0 {
		out.Reasons = append(out.Reasons, Reason{"source_changed", []string{}, "One or more history sources changed during batch reading"})
	}
	return out, nil
}

func batchTaskIdentity(p workspace.ProjectSnapshot, t taskwork.TaskRecord) bool {
	member := false
	for _, project := range p.Projects {
		if project.ID != t.ProjectID {
			continue
		}
		for _, id := range project.RepoIDs {
			if id == t.RepoID {
				member = true
			}
		}
	}
	if !member {
		return false
	}
	for _, repo := range p.Repos {
		if repo.ID == t.RepoID {
			return repo.GitCommonDir == t.GitCommonDir
		}
	}
	return false
}

// lifecycle is the lifecycle of the Project of work.
func (s Service) lifecycle(work workspace.ProjectWork) (taskwork.WorkItemLifecycleSnapshot, error) {
	if s.readBatch != nil {
		return s.readBatch.lifecycles.Of(work.Project), nil
	}
	return taskwork.ReadProjectLifecycle(work)
}

func appendLifecycleEvents(out *Snapshot, l taskwork.WorkItemLifecycleSnapshot, kind, id string) {
	for _, change := range l.Events {
		if change.SubjectKind != kind || change.SubjectID != id {
			continue
		}
		src := out.source("lifecycle", l.Locator, &l.SHA256, nil, nil)
		e := nativeEvent("lifecycle_changed", change.ID, "Lifecycle: "+string(change.From)+" -> "+string(change.To), change, src, &change.RecordedAtUTC)
		// The sidecar already owns immutable event identity. Neither the full
		// sidecar hash nor a later lifecycle transition changes that identity.
		e.EventID = change.ID
		e.Actor = &Actor{change.ActorClaim, "tool", nil}
		e.Data = raw(map[string]any{"basis": "Explicit lifecycle transition: " + string(change.From) + " -> " + string(change.To)})
		out.Events = append(out.Events, e)
	}
}

// StableBatchReasons exposes deterministic short codes for a client's per-Task
// completeness marker. It does not turn a process report into native QA.
func StableBatchReasons(s Snapshot) []string {
	set := map[string]bool{}
	for _, r := range s.Coverage.Reasons {
		if r.Code != "no_step_reporting" {
			set[r.Code] = true
		}
	}
	out := make([]string, 0, len(set))
	for code := range set {
		out = append(out, strings.TrimSpace(code))
	}
	sort.Strings(out)
	return out
}

func (s Service) readProjectEvents(in EventBatchInput, reader *eventBatchReader, owner taskwork.ProjectRegistry, out *EventBatch) error {
	for _, t := range owner.Registry.Tasks {
		basis := taskwork.TaskJournalBasis{Workspace: in.Workspace, Task: t, Registry: owner.Registry, Work: owner.Work}
		records, readErr := readStore(owner.Work, string(t.ID), workspaceID(in.Workspace.Root, in.Workspace.MarkerSHA256))
		snap, err := s.build(basis, records)
		if err != nil {
			return err
		}
		if readErr != nil {
			snap.reason("source_invalid", "Journal is damaged; showing the validated prefix and independent native facts: "+readErr.Error())
		}
		if !batchTaskIdentity(in.Projects, t) {
			snap.reason("task_identity_unavailable", "Task repository does not match the captured project registration")
		}
		out.Tasks[t.ID] = snap
	}
	for _, epic := range owner.Registry.Epics {
		snap := Snapshot{Events: []Event{}, Sources: []Source{}, Coverage: Coverage{State: "complete", Reasons: []Reason{}}, sourceCache: reader.sources, work: owner.Work}
		appendLifecycleEvents(&snap, in.Lifecycles.Of(owner.Work.Project), "epic", string(epic.ID))
		out.EpicEvents[epic.ID] = snap.Events
		out.EpicSources[epic.ID] = snap.Sources
	}
	return nil
}
