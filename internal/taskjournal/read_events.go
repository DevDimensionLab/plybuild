package taskjournal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// EventBatchInput binds one already-read registry snapshot. This event reader
// projects history, not current candidate validity, and never probes Git.
type EventBatchInput struct {
	Workspace workspace.WorkspaceObservation
	Projects  workspace.ProjectSnapshot
	Registry  workspace.WorkItemRegistry
	Lifecycle workspace.WorkItemLifecycleSnapshot
}

type EventBatch struct {
	Tasks       map[workspace.TaskID]Snapshot
	EpicEvents  map[workspace.EpicID][]Event
	EpicSources map[workspace.EpicID][]Source
	Reasons     []Reason
}

type sourceRead struct {
	hash string
	err  error
}
type bytesRead struct {
	bytes []byte
	err   error
}
type runRead struct {
	journal taskrun.JournalReadback
	err     error
}
type eventBatchReader struct {
	lifecycle             workspace.WorkItemLifecycleSnapshot
	sources               map[string]sourceRead
	requests              map[string]bytesRead
	runs                  map[string]runRead
	entries               []os.DirEntry
	entriesErr            error
	entriesRead           bool
	preparationsValidated bool
	preparationsErr       error
}

func cachedSourceHash(cache map[string]sourceRead, path string) (string, error) {
	if cache != nil {
		if r, ok := cache[path]; ok {
			return r.hash, r.err
		}
	}
	h, e := sourceHash(path)
	if cache != nil {
		cache[path] = sourceRead{h, e}
	}
	return h, e
}

func ReadEventBatch(d workspace.Dependencies, in EventBatchInput) (EventBatch, error) {
	out := EventBatch{Tasks: map[workspace.TaskID]Snapshot{}, EpicEvents: map[workspace.EpicID][]Event{}, EpicSources: map[workspace.EpicID][]Source{}, Reasons: []Reason{}}
	reader := &eventBatchReader{lifecycle: in.Lifecycle, sources: map[string]sourceRead{}, requests: map[string]bytesRead{}, runs: map[string]runRead{}}
	s := New(d)
	s.readBatch = reader
	registryPath := filepath.Join(in.Workspace.Root, workspace.MarkerDirectory, workspace.WorkItemsFile)
	for _, t := range in.Registry.Tasks {
		basis := workspace.TaskJournalBasis{Workspace: in.Workspace, Task: t, Registry: in.Registry, RegistryLocator: registryPath}
		records, readErr := readStore(in.Workspace.Root, string(t.ID), workspaceID(in.Workspace.Root, in.Workspace.MarkerSHA256))
		snap, _, err := s.build(basis, records)
		if err != nil {
			return out, err
		}
		if readErr != nil {
			snap.reason("source_invalid", "Journal is damaged; showing the validated prefix and independent native facts: "+readErr.Error())
		}
		if !batchTaskIdentity(in.Projects, t) {
			snap.reason("task_identity_unavailable", "Task repository does not match the captured project registration")
		}
		out.Tasks[t.ID] = snap
	}
	for _, epic := range in.Registry.Epics {
		snap := Snapshot{Events: []Event{}, Sources: []Source{}, Coverage: Coverage{State: "complete", Reasons: []Reason{}}, sourceCache: reader.sources}
		appendLifecycleEvents(&snap, in.Lifecycle, "epic", string(epic.ID))
		out.EpicEvents[epic.ID] = snap.Events
		out.EpicSources[epic.ID] = snap.Sources
	}
	// A single recheck per distinct source detects concurrent changes. This is
	// explicitly not a cross-register transaction or a filesystem snapshot.
	changed := map[string]bool{}
	for path, before := range reader.sources {
		if before.err != nil {
			continue
		}
		after, err := sourceHash(path)
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

func batchTaskIdentity(p workspace.ProjectSnapshot, t workspace.TaskRecord) bool {
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

func (s Service) lifecycle(root string) (workspace.WorkItemLifecycleSnapshot, error) {
	if s.readBatch != nil {
		return s.readBatch.lifecycle, nil
	}
	return workspace.ReadWorkItemLifecycle(s.Workspace, root)
}

func appendLifecycleEvents(out *Snapshot, l workspace.WorkItemLifecycleSnapshot, kind, id string) {
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

func (s Service) validatePreparation(b workspace.TaskJournalBasis, id string) (workspace.TaskPreparation, error) {
	if s.readBatch == nil {
		return workspace.ValidateTaskJournalPreparation(s.Workspace, b, id)
	}
	var preparation *workspace.TaskPreparation
	for i := range b.Registry.TaskPreparations {
		p := &b.Registry.TaskPreparations[i]
		if p.ID == id && p.Plan.TaskID == b.Task.ID {
			preparation = p
			break
		}
	}
	// Membership is request-specific. Only validate/cache the shared registry
	// closure after it matches, so a damaged run cannot poison sibling Tasks.
	if preparation == nil {
		return workspace.TaskPreparation{}, fmt.Errorf("preparation does not belong to Task")
	}
	r := s.readBatch
	if !r.preparationsValidated {
		_, r.preparationsErr = workspace.ValidateTaskJournalPreparation(s.Workspace, b, id)
		r.preparationsValidated = true
	}
	if r.preparationsErr != nil {
		return workspace.TaskPreparation{}, r.preparationsErr
	}
	return *preparation, nil
}

func (s Service) nativeRunEntries(dir string) ([]os.DirEntry, error) {
	if s.readBatch != nil && s.readBatch.entriesRead {
		return s.readBatch.entries, s.readBatch.entriesErr
	}
	var entries []os.DirEntry
	d, err := openDir(dir, false)
	if err == nil {
		entries, err = d.ReadDir(-1)
		_ = d.Close()
		sortEntries(entries)
	}
	if s.readBatch != nil {
		s.readBatch.entries, s.readBatch.entriesErr, s.readBatch.entriesRead = entries, err, true
	}
	return entries, err
}
func (s Service) nativeRunRequest(path string) ([]byte, error) {
	if s.readBatch != nil {
		if r, ok := s.readBatch.requests[path]; ok {
			return r.bytes, r.err
		}
	}
	b, e := readFile(path, 1<<20)
	if s.readBatch != nil {
		s.readBatch.requests[path] = bytesRead{b, e}
	}
	return b, e
}
func (s Service) nativeRunJournal(root, id string) (taskrun.JournalReadback, error) {
	if s.readBatch != nil {
		if r, ok := s.readBatch.runs[id]; ok {
			return r.journal, r.err
		}
	}
	j, e := taskrun.ReadProcessJournal(taskrun.SystemDependencies(s.Workspace), root, id)
	if s.readBatch != nil {
		s.readBatch.runs[id] = runRead{j, e}
	}
	return j, e
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
