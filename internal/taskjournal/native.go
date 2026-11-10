package taskjournal

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func sortEntries(e []os.DirEntry) {
	sort.Slice(e, func(i, j int) bool { return e[i].Name() < e[j].Name() })
}
func (s *Snapshot) reason(code, detail string, ids ...string) {
	s.Coverage.State = "partial"
	s.Coverage.Reasons = append(s.Coverage.Reasons, Reason{code, append([]string{}, ids...), detail})
}
func (s *Snapshot) source(kind, path string, expected *string, seq *int, prev *string) string {
	id := "src_" + digest([]any{kind, path, expected, seq})
	for _, v := range s.Sources {
		if v.SourceID == id {
			return id
		}
	}
	status := "valid"
	var observed *string
	h, e := cachedSourceHash(s.sourceCache, s.Workspace.Root, kind, path)
	if e != nil {
		status = "invalid"
		if os.IsNotExist(e) {
			status = "missing"
		}
	} else {
		observed = &h
		if expected != nil && h != strings.TrimPrefix(*expected, "sha256:") {
			status = "changed"
		}
	}
	s.Sources = append(s.Sources, Source{id, kind, path, observed, seq, prev, status})
	if status != "valid" {
		s.reason("source_"+status, "Source could not be read at its recorded version: "+path, id)
	}
	return id
}
func nativeID(kind, id, recordHash string) string {
	return "native_" + digest([]string{kind, id, recordHash})
}
func nativeEvent(kind, id, title string, value any, source string, recorded *string) Event {
	return Event{EventID: nativeID(kind, id, digest(value)), Origin: "native", Type: kind, Title: title, RecordedAt: normalize(recorded), TimeBasis: unknownTime(), SourceIDs: []string{source}, Relations: []Relation{}, SupersededBy: []string{}, Data: raw(map[string]any{"basis": "native state; no occurred time is asserted"})}
}
func axisEvent(e Event, axis, value string, c *Candidate) Event {
	e.Outcome = ptr(value)
	e.Candidate = c
	e.Data = raw(map[string]any{"axis": axis, "value": value})
	return e
}
func candidate(t workspace.TaskRecord, oid, tree string) *Candidate {
	if t.Worktree == nil {
		return nil
	}
	return &Candidate{string(t.RepoID), string(t.Worktree.ID), oid, tree}
}

func (s Service) native(b workspace.TaskJournalBasis) (Snapshot, error) {
	root := b.Workspace.Root
	t := b.Task
	out := Snapshot{Kind: "ply.workspace.task-journal-snapshot", SchemaVersion: 1, AsOf: s.Now().UTC().Format(time.RFC3339Nano), Workspace: Workspace{workspaceID(root, b.Workspace.MarkerSHA256), root, b.Workspace.MarkerSHA256}, Task: Task{string(t.ID), t.Title, string(t.ProjectID), string(t.RepoID), string(t.ParentEpicID), t.Worktree, digest(t)}, Sources: []Source{}, Events: []Event{}, Steps: []Step{}, Lanes: []Lane{}, Observations: []Observation{}, Coverage: Coverage{"complete", []Reason{}, []string{}, []string{}}, Selection: Selection{Order: "occurred", EventIDs: []string{}}, Current: Current{Axes: map[string]Axis{}, NextAction: NextAction{"Inspect the evidence and clarify the next action with the Task owner.", nil, []string{}}}}
	if s.readBatch != nil {
		out.sourceCache = s.readBatch.sources
	}
	marker := out.source("workspace", filepath.Join(root, ".ply", "workspace.yaml"), &b.Workspace.MarkerSHA256, nil, nil)
	registry := out.source("task_registry", b.RegistryLocator, &b.Registry.RawSHA256, nil, nil)
	state := nativeEvent("task_state", string(t.ID), "Task state: "+t.Title, t, registry, nil)
	state.SourceIDs = append(state.SourceIDs, marker)
	out.Events = append(out.Events, state)
	if t.Worktree != nil && s.readBatch == nil {
		g, e := workspace.ObserveTaskJournalWorktree(s.Workspace, t.Worktree.Locator)
		if e == nil && g.Clean && g.InventoryMatch && g.Locator == t.Worktree.Locator && g.Ref == t.Worktree.Ref && g.GitCommonDir == t.GitCommonDir {
			out.Current.Candidate = candidate(t, g.OID, g.Tree)
		} else {
			out.reason("candidate_unknown", "The Task worktree is missing, dirty, changed in identity, or could not be observed.", registry)
		}
	}
	lifecycle, lifecycleErr := s.lifecycle(root)
	if lifecycleErr != nil {
		out.reason("lifecycle_unavailable", lifecycleErr.Error())
	} else {
		appendLifecycleEvents(&out, lifecycle, "task", string(t.ID))
	}
	for _, p := range b.Registry.TaskContentPublications {
		if p.TaskID != t.ID {
			continue
		}
		versions, e := workspace.ReadTaskJournalContent(s.Workspace, b, p)
		sources := []string{}
		for _, v := range versions {
			sources = append(sources, out.source(v.Kind, v.Locator, &v.SHA256, nil, nil))
		}
		src := sources[0]
		if e != nil {
			out.reason("source_invalid", e.Error(), src)
			continue
		}
		ev := nativeEvent(p.Operation, p.OutcomeRef.ManifestSHA256, "Task "+p.Operation, p, src, &p.RecordedAtUTC)
		ev.SourceIDs = sources
		ev.Data = raw(map[string]any{"operation": p.Operation, "manifest_sha256": p.OutcomeRef.ManifestSHA256})
		out.Events = append(out.Events, ev)
	}
	integrationSources := []string{}
	for _, source := range workspace.IntegrationEvidenceSources(root, b.Registry) {
		integrationSources = append(integrationSources, out.source(source.Kind, source.Locator, &source.SHA256, nil, nil))
	}
	for _, ir := range b.Registry.IntegrationResults {
		var a *workspace.IntegrationAuthority
		for i := range b.Registry.IntegrationAuthorities {
			p := &b.Registry.IntegrationAuthorities[i]
			if p.ID == ir.AuthorityID && p.TaskID == t.ID {
				a = p
				break
			}
		}
		if a == nil {
			continue
		}
		e := nativeEvent("integration", string(ir.ID), "Integration: "+ir.Outcome, ir, registry, &ir.RecordedAtUTC)
		e.SourceIDs = append(e.SourceIDs, integrationSources...)
		e = axisEvent(e, "integration", ir.Outcome, candidate(t, a.Plan.Task.ResultOID, a.Plan.Task.ResultTree))
		out.Events = append(out.Events, e)
	}
	return out, nil
}
