package taskjournal

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
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
	h, e := cachedSourceHash(s.sourceCache, path)
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

func (s Service) native(b workspace.TaskJournalBasis) (Snapshot, map[string]RunBinding, error) {
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
	for _, p := range b.Registry.TaskPreparations {
		if p.Plan.TaskID != t.ID {
			continue
		}
		if _, e := s.validatePreparation(b, p.ID); e != nil {
			out.reason("source_invalid", e.Error(), registry)
			continue
		}
		ev := nativeEvent("preparation", p.ID, "Task preparation", p, registry, &p.CreatedAtUTC)
		ev.Candidate = candidate(t, p.Plan.ParentOID, p.Plan.ParentTree)
		ev.Data = raw(map[string]any{"preparation_id": p.ID, "preparation_sha256": "sha256:" + digest(p)})
		out.Events = append(out.Events, ev)
	}
	results := map[workspace.TaskResultID]workspace.TaskResultRecord{}
	for _, r := range b.Registry.TaskResults {
		if r.TaskID != t.ID {
			continue
		}
		results[r.ID] = r
		c := candidate(t, r.ResultOID, r.ResultTree)
		sources := []string{registry}
		for _, a := range []SourceRef{{r.HandoffLocator, r.HandoffSHA256}, {r.StartReceiptLocator, r.StartReceiptSHA256}, {r.TerminalResultLocator, r.TerminalResultSHA256}} {
			sources = append(sources, out.source("result_evidence", a.Locator, &a.SHA256, nil, nil))
		}
		for _, a := range r.Artifacts {
			sources = append(sources, out.source("result_artifact", a.Locator, &a.SHA256, nil, nil))
		}
		for _, axis := range []string{"reported", "technical"} {
			value := r.ReportedOutcome
			if axis == "technical" {
				value = r.TechnicalGate
			}
			e := nativeEvent("task_result_"+axis, string(r.ID), "Task result "+axis+": "+value, r, registry, &r.Recorder.RecordedAtUTC)
			e.SourceIDs = sources
			e.Actor = &Actor{r.Recorder.ActorClaim, "tool", nil}
			e.ActivityID = ptr(r.ActivityID)
			e.RunID = ptr(r.RunID)
			e = axisEvent(e, axis, value, c)
			out.Events = append(out.Events, e)
		}
	}
	for _, q := range b.Registry.HumanQARecords {
		if q.TaskID != t.ID {
			continue
		}
		r, ok := results[q.TaskResultID]
		if !ok {
			out.reason("source_invalid", "QA result binding is unavailable", registry)
			continue
		}
		e := nativeEvent("human_qa", string(q.ID), "Human QA: "+q.Outcome, q, registry, nil)
		e.Actor = &Actor{q.Actor.ActorClaim, "human", nil}
		e.OccurredAt = normalize(&q.Actor.CompletedAtUTC)
		e.OriginalOccurredAt = ptr(q.Actor.CompletedAtUTC)
		e.TimeBasis = TimeBasis{"reported", ptr("native human QA attestation"), "nanosecond", nil}
		e = axisEvent(e, "human", q.Outcome, candidate(t, q.ResultOID, q.ResultTree))
		e.RunID = ptr(r.RunID)
		for _, a := range q.Evidence {
			e.SourceIDs = append(e.SourceIDs, out.source("qa_evidence", a.Locator, &a.SHA256, nil, nil))
		}
		out.Events = append(out.Events, e)
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
		e = axisEvent(e, "integration", ir.Outcome, candidate(t, a.Plan.Task.ResultOID, a.Plan.Task.ResultTree))
		out.Events = append(out.Events, e)
	}
	runs := map[string]RunBinding{}
	dir := filepath.Join(root, ".ply", "task-runs", "v1", "runs")
	entries, e := s.nativeRunEntries(dir)
	if e != nil && !os.IsNotExist(e) {
		out.reason("source_invalid", e.Error())
	}
	if e == nil {
		for _, entry := range entries {
			id := entry.Name()
			path := filepath.Join(dir, id, "request.json")
			rb, re := s.nativeRunRequest(path)
			if re != nil {
				out.reason("source_invalid", "Run request unavailable: "+path)
				continue
			}
			var request taskrun.Request
			if json.Unmarshal(rb, &request) != nil {
				out.reason("source_invalid", "Invalid run request: "+path)
				continue
			}
			projection, pe := workflowhandoff.ValidateTaskRunDraft(request.HandoffDraft)
			if pe == nil && projection.Basis != nil && projection.Basis.TaskID != t.ID {
				continue
			}
			rs := out.source("run_request", path, ptr("sha256:"+hash(rb)), nil, nil)
			prep, ve := s.validatePreparation(b, request.PreparationID)
			if pe != nil || projection.Basis == nil || ve != nil || request.WorkspaceRoot != root || request.PreparationSHA256 != "sha256:"+digest(prep) || prep.Plan.Workspace.Root != root || prep.Plan.Workspace.MarkerSHA256 != b.Workspace.MarkerSHA256 || projection.Basis.TaskID != t.ID || prep.Outcome == nil || projection.Basis.TaskWorktreeID != prep.Outcome.WorktreeID || projection.Worktree != prep.Plan.WorktreePath || projection.OID != prep.Plan.ParentOID || projection.Ref != "refs/heads/"+prep.Plan.Branch || projection.Basis.Spec != prep.Plan.Spec || projection.Basis.Selection != prep.Plan.Selection || projection.Basis.Problem != prep.Plan.Problem || projection.Basis.Assessment != prep.Plan.Assessment {
				out.markUnbound(rs, "Run request does not match the exact historical Task preparation")
				continue
			}
			j, je := s.nativeRunJournal(root, id)
			if j.Request.RequestKey == "" || taskrun.RunID(request.RequestKey) != id || !equal(request, j.Request) || j.Binding != nil && !equal(j.Binding.Preparation, prep) {
				out.markUnbound(rs, "Run identity or preserved preparation differs")
				continue
			}
			if je != nil {
				out.reason("source_invalid", je.Error(), rs)
			} else {
				runs[id] = RunBinding{id, "sha256:" + digest(request), prep.ID, request.PreparationSHA256}
			}
			var previousTime *time.Time
			for i, event := range j.Events {
				h := j.Hashes[i]
				src := out.source("run_event", filepath.Join(dir, id, "events", fmt.Sprintf("%06d.json", event.Sequence)), &h, &event.Sequence, event.PreviousSHA256)
				ev := nativeEvent("run_"+event.Type, id+":"+fmt.Sprint(event.Sequence), "Run "+event.Type, event, src, &event.RecordedAtUTC)
				ev.SourceIDs = append(ev.SourceIDs, rs)
				ev.RunID = ptr(id)
				ev.Actor = &Actor{"ply", "tool", nil}
				ev.SourceSequence = ptr(event.Sequence)
				ev.Candidate = candidate(t, prep.Plan.ParentOID, prep.Plan.ParentTree)
				ev.Data = raw(map[string]any{"run_event_type": event.Type, "event_sha256": h, "previous_sha256": event.PreviousSHA256})
				out.Events = append(out.Events, ev)
				tm, _ := parseTime(event.RecordedAtUTC)
				if previousTime != nil && tm.Before(*previousTime) {
					out.Coverage.TimeConflictEventIDs = append(out.Coverage.TimeConflictEventIDs, ev.EventID)
					out.reason("time_conflict", "Native sequence contradicts its registration clock", src)
				}
				previousTime = &tm
			}
			if je == nil && len(j.Events) > 0 {
				last := j.Events[len(j.Events)-1]
				ev := nativeEvent("run_state", id, "Run transport: "+j.Result.Process.State, j.Result, rs, &last.RecordedAtUTC)
				ev.RunID = ptr(id)
				ev.Actor = &Actor{"ply", "tool", nil}
				ev = axisEvent(ev, "transport", j.Result.Process.State, candidate(t, prep.Plan.ParentOID, prep.Plan.ParentTree))
				out.Events = append(out.Events, ev)
			}
		}
	}
	return out, runs, nil
}
func (s *Snapshot) markUnbound(id, detail string) {
	for i := range s.Sources {
		if s.Sources[i].SourceID == id {
			s.Sources[i].Status = "unbound"
		}
	}
	s.reason("run_unbound", detail, id)
}
