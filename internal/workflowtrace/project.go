package workflowtrace

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type builder struct {
	out        Result
	events     map[string]int
	sources    map[string]int
	candidates map[string]int
	verifiers  map[string]taskrun.TraceVerification
	runAliases map[string][]string
}

func newBuilder(root string, task workspace.TaskRecord, at string) *builder {
	return &builder{out: Result{Kind: Kind, SchemaVersion: 1, Workspace: root,
		Task: Task{string(task.ID), task.Title, string(task.ProjectID), string(task.RepoID), string(task.ParentEpicID)}, ObservedAtUTC: at,
		Freshness: "fresh", LiveState: "unknown", HistoryState: "empty",
		Ordering:     "Source chains preserve their own sequence. Independent chains have no asserted global chronology. Sorted IDs and tied positions are presentation only; spacing is not time.",
		CurrentGoals: []Goal{}, Runs: []Run{}, Sources: []Source{}, Events: []Event{}, Chains: []Chain{}, Relations: []Relation{}, Candidates: []Candidate{}, Analysis: []Analysis{}, Diagnostics: []Diagnostic{},
		Coverage: Coverage{State: "complete", Read: []string{"registered_task", "registered_goals", "task_journal", "legacy_task_runs", "workflow_runs", "native_results_qa_integration"}, Unknowns: []string{
			"Recorded-source coverage is not complete work history; missing reports do not prove no unrecorded work or waiting.",
			"Active work, attention, cost, waiting cause, live activity and time savings are unknown.",
			"Actor claims and human attestations are preserved local provenance, not authenticated human identity.",
			"Reads recheck sources but do not form a transactional filesystem snapshot.",
		}}}, events: map[string]int{}, sources: map[string]int{}, candidates: map[string]int{}, verifiers: map[string]taskrun.TraceVerification{}, runAliases: map[string][]string{}}
}

func raw(v any) json.RawMessage { b, _ := json.Marshal(v); return b }
func ptr[T any](v T) *T         { return &v }
func nonempty(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
func value(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
func unique(a []string) []string {
	m := map[string]bool{}
	out := []string{}
	for _, s := range a {
		if s != "" && !m[s] {
			m[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}
func hash(v any) string { b, _ := taskjournal.Canonical(v); return fmt.Sprintf("%x", sha256.Sum256(b)) }
func nativeID(kind, id string, record any) string {
	return "native_" + hash([]string{kind, id, hash(record)})
}

func (b *builder) diagnostic(code, detail string, sources, events []string) {
	b.out.Diagnostics = append(b.out.Diagnostics, Diagnostic{code, detail, unique(sources), unique(events)})
	b.out.Coverage.State = "partial"
}
func (b *builder) source(s Source) string {
	if s.SHA256 != nil && len(*s.SHA256) == 64 {
		s.SHA256 = ptr("sha256:" + *s.SHA256)
	}
	if s.ID == "" {
		s.ID = "source_" + hash([]any{s.Kind, s.Locator, s.SHA256})
	}
	if i, ok := b.sources[s.ID]; ok {
		if s.Status != "valid" {
			b.out.Sources[i].Status = s.Status
		}
		return s.ID
	}
	b.sources[s.ID] = len(b.out.Sources)
	b.out.Sources = append(b.out.Sources, s)
	return s.ID
}
func (b *builder) candidate(oid, tree string) *string {
	if oid == "" || tree == "" {
		return nil
	}
	id := "candidate_" + hash([]string{b.out.Task.RepoID, oid, tree})
	if _, ok := b.candidates[id]; !ok {
		b.candidates[id] = len(b.out.Candidates)
		b.out.Candidates = append(b.out.Candidates, Candidate{ID: id, OID: oid, Tree: tree, RunIDs: []string{}, ResultIDs: []string{}, EventIDs: []string{}, Axes: map[string]Axis{}})
	}
	return &id
}
func (b *builder) event(e Event) {
	if e.Data == nil {
		e.Data = raw(nil)
	}
	if e.RunIDs == nil {
		e.RunIDs = []string{}
	}
	if e.SourceIDs == nil {
		e.SourceIDs = []string{}
	}
	if e.Positions == nil {
		e.Positions = []Position{}
	}
	if i, ok := b.events[e.ID]; ok {
		old := &b.out.Events[i]
		if old.Type != e.Type || old.Outcome != nil && e.Outcome != nil && *old.Outcome != *e.Outcome || old.CandidateID != nil && e.CandidateID != nil && *old.CandidateID != *e.CandidateID {
			b.diagnostic("event_identity_conflict", "The same exact event identity has conflicting projections; both source bindings remain visible.", append(old.SourceIDs, e.SourceIDs...), []string{e.ID})
		}
		old.SourceIDs = unique(append(old.SourceIDs, e.SourceIDs...))
		old.RunIDs = unique(append(old.RunIDs, e.RunIDs...))
		old.Positions = append(old.Positions, e.Positions...)
		if old.NativeID == nil {
			old.NativeID = e.NativeID
		}
		if old.ResultID == nil {
			old.ResultID = e.ResultID
		}
		if old.Outcome == nil {
			old.Outcome = e.Outcome
		}
		if old.CandidateID == nil {
			old.CandidateID = e.CandidateID
		}
		if old.Role == "unknown" && e.Role != "unknown" {
			old.Role, old.ActorClaim = e.Role, e.ActorClaim
		}
		if old.EvidenceClass == "registered" && e.EvidenceClass != "registered" {
			old.EvidenceClass = e.EvidenceClass
			old.Role, old.ActorClaim = e.Role, e.ActorClaim
		}
		var data map[string]json.RawMessage
		if json.Unmarshal(old.Data, &data) == nil && data != nil {
			var additional map[string]json.RawMessage
			if json.Unmarshal(e.Data, &additional) == nil {
				for key, v := range additional {
					if _, exists := data[key]; !exists {
						data[key] = v
					}
				}
			}
			var projections []json.RawMessage
			_ = json.Unmarshal(data["projection_records"], &projections)
			data["projection_records"] = raw(append(projections, raw(map[string]any{"actor_claim": e.ActorClaim, "recorder": e.Recorder, "evidence_class": e.EvidenceClass, "data": e.Data})))
			old.Data = raw(data)
		}
		return
	}
	b.events[e.ID] = len(b.out.Events)
	b.out.Events = append(b.out.Events, e)
}

func (b *builder) journal(s taskjournal.Snapshot, r workspace.WorkItemRegistry) {
	for _, src := range s.Sources {
		b.source(Source{src.SourceID, src.Kind, src.Locator, src.SHA256, src.Status})
	}
	for _, d := range s.Coverage.Reasons {
		// This describes the unobserved world, not failure to read recorded facts.
		if d.Code == "no_step_reporting" {
			b.out.Coverage.Unknowns = append(b.out.Coverage.Unknowns, d.Detail)
			continue
		}
		b.diagnostic(d.Code, d.Detail, d.SourceIDs, nil)
	}
	for _, in := range s.Events {
		e := Event{ID: in.EventID, Type: in.Type, Title: in.Title, Role: "unknown", Recorder: "task journal native projection", EvidenceClass: "registered", RunIDs: []string{}, Outcome: in.Outcome, OccurredAtUTC: in.OccurredAt, RegisteredAtUTC: in.RecordedAt, TimeBasis: in.TimeBasis, SourceIDs: in.SourceIDs, Data: in.Data, Positions: []Position{}}
		if in.RunID != nil {
			e.RunIDs = append(e.RunIDs, *in.RunID)
		}
		if in.Origin == "native" && in.Actor != nil {
			e.Recorder = in.Actor.ID + " (native record claim)"
		}
		if in.Candidate != nil && in.Type != "preparation" && !strings.HasPrefix(in.Type, "run_") {
			e.CandidateID = b.candidate(in.Candidate.OID, in.Candidate.Tree)
		}
		if in.Origin == "contribution" {
			e.EvidenceClass = "reported"
			e.Recorder = "task journal contribution"
			if in.Actor != nil {
				e.ActorClaim = nonempty(in.Actor.ID)
				e.Role = journalRole(in.Actor.Role)
			}
			if in.SourceSequence != nil {
				e.Positions = append(e.Positions, Position{"journal/" + b.out.Task.ID, *in.SourceSequence})
			}
		} else if in.SourceSequence != nil && in.RunID != nil {
			e.Positions = append(e.Positions, Position{*in.RunID + "/events", *in.SourceSequence})
			e.Role = "ply"
			e.ActorClaim = ptr("ply")
		}
		if in.TimeBasis.Kind == "reported" {
			e.ReportedAtUTC, e.OccurredAtUTC = in.OccurredAt, nil
		}
		if in.Type == "unavailable_reference" {
			e.EvidenceClass = "unknown"
		}
		b.enrichNative(&e, r)
		b.event(e)
		for _, rel := range in.Relations {
			b.out.Relations = append(b.out.Relations, Relation{rel.Type, e.ID, rel.EventID, "explicit journal relation", unique(in.SourceIDs)})
		}
	}
}

func journalRole(role string) string {
	switch role {
	case "human":
		return "human"
	case "developer", "coordinator", "reviewer", "agent":
		return "agent"
	default:
		return "unknown"
	}
}

func (b *builder) enrichNative(e *Event, registry workspace.WorkItemRegistry) {
	for _, r := range registry.TaskResults {
		if string(r.TaskID) != b.out.Task.ID {
			continue
		}
		for _, axis := range []string{"reported", "technical"} {
			if e.ID != nativeID("task_result_"+axis, string(r.ID), r) {
				continue
			}
			e.NativeID, e.ResultID = ptr(string(r.ID)), ptr(string(r.ID))
			e.CandidateID = b.candidate(r.ResultOID, r.ResultTree)
			e.Recorder = r.Recorder.ActorClaim + " (" + r.Recorder.ControlSurface + ")"
			e.Role, e.ActorClaim = "unknown", nil
			e.Data = raw(r)
			if axis == "technical" {
				e.EvidenceClass = "controlled"
				e.Role, e.ActorClaim = "ply", ptr("ply")
			} else {
				e.EvidenceClass = "reported"
			}
		}
	}
	for _, q := range registry.HumanQARecords {
		if string(q.TaskID) != b.out.Task.ID || e.ID != nativeID("human_qa", string(q.ID), q) {
			continue
		}
		e.NativeID, e.ResultID = ptr(string(q.ID)), ptr(string(q.TaskResultID))
		e.CandidateID = b.candidate(q.ResultOID, q.ResultTree)
		e.Role, e.ActorClaim, e.EvidenceClass = "human", nonempty(q.Actor.ActorClaim), "human_attestation"
		e.ReportedAtUTC = nonempty(q.Actor.CompletedAtUTC)
		e.OccurredAtUTC = nil
		e.Data = raw(q)
	}
	for _, ir := range registry.IntegrationResults {
		if e.ID != nativeID("integration", string(ir.ID), ir) {
			continue
		}
		for _, a := range registry.IntegrationAuthorities {
			if a.ID == ir.AuthorityID && string(a.TaskID) == b.out.Task.ID {
				e.NativeID, e.ResultID = ptr(string(ir.ID)), ptr(string(a.TaskResultID))
				e.CandidateID = b.candidate(a.Plan.Task.ResultOID, a.Plan.Task.ResultTree)
				e.Role, e.ActorClaim, e.EvidenceClass = "ply", ptr("ply"), "controlled"
				e.Data = raw(ir)
			}
		}
	}
}

func (b *builder) finish() Result {
	positions := map[string][]struct {
		id  string
		seq int
	}{}
	for i := range b.out.Events {
		e := &b.out.Events[i]
		e.SourceIDs, e.RunIDs = unique(e.SourceIDs), unique(e.RunIDs)
		sort.Slice(e.Positions, func(i, j int) bool {
			if e.Positions[i].ChainID != e.Positions[j].ChainID {
				return e.Positions[i].ChainID < e.Positions[j].ChainID
			}
			return e.Positions[i].Sequence < e.Positions[j].Sequence
		})
		seen := map[Position]bool{}
		clean := []Position{}
		for _, p := range e.Positions {
			if !seen[p] {
				seen[p] = true
				clean = append(clean, p)
				positions[p.ChainID] = append(positions[p.ChainID], struct {
					id  string
					seq int
				}{e.ID, p.Sequence})
			}
		}
		e.Positions = clean
		if e.CandidateID != nil {
			c := &b.out.Candidates[b.candidates[*e.CandidateID]]
			c.RunIDs = append(c.RunIDs, e.RunIDs...)
			c.ResultIDs = append(c.ResultIDs, value(e.ResultID))
			c.EventIDs = append(c.EventIDs, e.ID)
			axis := map[string]string{"task_result_reported": "reported", "task_result_technical": "technical", "human_qa": "human_qa", "integration": "integration"}[e.Type]
			if axis != "" && e.Outcome != nil && b.validEvent(*e) {
				a := c.Axes[axis]
				a.Outcomes = unique(append(a.Outcomes, *e.Outcome))
				a.EvidenceIDs = unique(append(a.EvidenceIDs, e.ID))
				a.State = "recorded"
				if len(a.Outcomes) > 1 {
					a.State = "conflicting"
				}
				c.Axes[axis] = a
			}
		}
	}
	for id, p := range positions {
		sort.Slice(p, func(i, j int) bool {
			if p[i].seq != p[j].seq {
				return p[i].seq < p[j].seq
			}
			return p[i].id < p[j].id
		})
		c := Chain{ID: id, Ordering: "source sequence; equal positions have no asserted order", EventIDs: []string{}}
		for i, x := range p {
			c.EventIDs = append(c.EventIDs, x.id)
			if i > 0 && p[i-1].seq < x.seq {
				b.out.Relations = append(b.out.Relations, Relation{"sequence", p[i-1].id, x.id, "source chain " + id, []string{}})
			}
		}
		b.out.Chains = append(b.out.Chains, c)
	}
	for i := range b.out.Candidates {
		c := &b.out.Candidates[i]
		c.RunIDs = unique(c.RunIDs)
		c.ResultIDs = unique(c.ResultIDs)
		c.EventIDs = unique(c.EventIDs)
		for _, k := range []string{"reported", "technical", "human_qa", "integration"} {
			if _, ok := c.Axes[k]; !ok {
				c.Axes[k] = Axis{"not_recorded", []string{}, []string{}}
			}
		}
	}
	b.analyze()
	if len(b.out.Runs) > 0 {
		b.out.HistoryState = "recorded"
	}
	if b.out.Coverage.State == "partial" {
		b.out.HistoryState = "partial"
	}
	for _, s := range b.out.Sources {
		if s.Status != "valid" {
			b.out.Freshness = "unknown"
		}
	}
	sort.Slice(b.out.Events, func(i, j int) bool { return b.out.Events[i].ID < b.out.Events[j].ID })
	sort.Slice(b.out.Sources, func(i, j int) bool { return b.out.Sources[i].ID < b.out.Sources[j].ID })
	sort.Slice(b.out.Chains, func(i, j int) bool { return b.out.Chains[i].ID < b.out.Chains[j].ID })
	sort.Slice(b.out.Candidates, func(i, j int) bool { return b.out.Candidates[i].ID < b.out.Candidates[j].ID })
	sort.Slice(b.out.Runs, func(i, j int) bool { return b.out.Runs[i].ID < b.out.Runs[j].ID })
	sort.Slice(b.out.Relations, func(i, j int) bool {
		a, z := b.out.Relations[i], b.out.Relations[j]
		return strings.Join([]string{a.From, a.Type, a.To}, "\x00") < strings.Join([]string{z.From, z.Type, z.To}, "\x00")
	})
	sort.Slice(b.out.Analysis, func(i, j int) bool { return b.out.Analysis[i].ID < b.out.Analysis[j].ID })
	sort.Slice(b.out.Diagnostics, func(i, j int) bool {
		return b.out.Diagnostics[i].Code+b.out.Diagnostics[i].Detail < b.out.Diagnostics[j].Code+b.out.Diagnostics[j].Detail
	})
	b.out.Coverage.Unknowns = unique(b.out.Coverage.Unknowns)
	return b.out
}

func (b *builder) validEvent(e Event) bool {
	if e.EvidenceClass == "unknown" || len(e.SourceIDs) == 0 {
		return false
	}
	for _, id := range e.SourceIDs {
		i, ok := b.sources[id]
		if !ok || b.out.Sources[i].Status != "valid" {
			return false
		}
	}
	return true
}
