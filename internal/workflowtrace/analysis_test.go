package workflowtrace

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func testBuilder() *builder {
	b := newBuilder("/fixture", workspace.TaskRecord{ID: "task", RepoID: "repo"}, "2026-10-07T12:00:00Z")
	b.source(Source{ID: "source", Kind: "fixture", Locator: "/fixture/evidence", SHA256: ptr("sha256:fixture"), Status: "valid"})
	return b
}
func testEvent(id, kind string) Event {
	return Event{ID: id, Type: kind, Role: "unknown", EvidenceClass: "registered", SourceIDs: []string{"source"}, Data: raw(map[string]any{}), TimeBasis: taskjournal.TimeBasis{Kind: "unknown", Precision: "unknown"}}
}
func analysisByID(t *testing.T, r Result, id string) Analysis {
	t.Helper()
	for _, a := range r.Analysis {
		if a.ID == id {
			return a
		}
	}
	t.Fatalf("analysis %s absent", id)
	return Analysis{}
}
func expectCount(t *testing.T, r Result, id string, want int) {
	t.Helper()
	a := analysisByID(t, r, id)
	if a.Value == nil || *a.Value != float64(want) {
		t.Fatalf("%s: %+v want %d", id, a, want)
	}
}
func eventByID(t *testing.T, r Result, id string) Event {
	t.Helper()
	for _, e := range r.Events {
		if e.ID == id {
			return e
		}
	}
	t.Fatalf("event %s absent", id)
	return Event{}
}
func fixtureVerifier(oid string) taskrun.TraceVerification {
	return taskrun.TraceVerification{AttemptID: "verify-00000002", CandidateOID: oid, CandidateTree: oid + "tree", Argv: []string{"/bin/sh", "/fixture/acceptance.sh"}, CWD: "/fixture", Acceptance: taskrun.FileBinding{Locator: "/fixture/acceptance.sh", SHA256: "same"}, AcceptanceSnapshot: taskrun.FileBinding{Locator: "/fixture/snapshot", SHA256: "same"}, Review: taskrun.FileBinding{Locator: "/fixture/review", SHA256: "review"}, Exit: ptr(0), StartedAtUTC: ptr("2026-10-07T10:00:00Z"), FinishedAtUTC: ptr("2026-10-07T10:01:27.293Z"), InputsBound: true}
}

func TestTraceNativeDedupRetainsAttemptsCandidatesAndSeparateJudgments(t *testing.T) {
	b := testBuilder()
	request := testEvent("request", "run_request")
	b.event(request)
	ids := []string{"request"}
	for n, oid := range []string{"c1", "c2", "c3"} {
		id := oid + "verify"
		e := testEvent(id, "verification")
		e.Role = "ply"
		e.EvidenceClass = "controlled"
		e.CandidateID = b.candidate(oid, oid+"tree")
		e.Outcome = ptr("passed")
		e.Positions = []Position{{"run/delivery", n*3 + 1}}
		v := fixtureVerifier(oid)
		if n == 0 {
			v.Exit = ptr(1)
			e.Outcome = ptr("failed")
		}
		b.verifiers[id] = v
		b.event(e)
		ids = append(ids, id)
		if n == 0 {
			correction := testEvent("correction", "delivery_report")
			correction.EvidenceClass = "reported"
			correction.Role = "agent"
			correction.Data = raw(map[string]any{"summary": "Corrected verifier failure"})
			correction.Positions = []Position{{"run/delivery", 2}}
			b.event(correction)
			ids = append(ids, correction.ID)
			continue
		}
		tech := testEvent(oid+"result", "task_result_technical")
		tech.CandidateID = e.CandidateID
		tech.ResultID = ptr(oid + "result")
		tech.Outcome = ptr("passed")
		tech.Positions = []Position{{"run/candidates", n}}
		b.event(tech)
		ids = append(ids, tech.ID)
		q := workspace.TaskHumanQARecord{ID: workspace.HumanQARecordID(oid + "qa"), TaskID: "task", TaskResultID: workspace.TaskResultID(oid + "result"), ResultOID: oid, ResultTree: oid + "tree", Outcome: "pass", Actor: workspace.TaskHumanActorRecord{ActorClaim: "Synthetic human", StartedAtUTC: "2026-10-07T10:00:00Z", CompletedAtUTC: "2026-10-07T16:12:47Z"}}
		if n == 1 {
			q.Outcome = "fail"
		}
		qa := testEvent(nativeID("human_qa", string(q.ID), q), "human_qa")
		qa.NativeID = ptr(string(q.ID))
		qa.Data = raw(q)
		qa.Role = "human"
		qa.ActorClaim = ptr(q.Actor.ActorClaim)
		qa.CandidateID = e.CandidateID
		qa.ResultID = tech.ResultID
		qa.Outcome = ptr(q.Outcome)
		qa.EvidenceClass = "human_attestation"
		qa.Positions = []Position{{"run/delivery", n*3 + 2}}
		b.event(qa)
		// The exact QA reached through the second native read is the same event.
		b.source(Source{ID: "execute-source", Kind: "execute", Locator: "/fixture/qa", SHA256: ptr("sha256:qa"), Status: "valid"})
		qa.SourceIDs = []string{"execute-source"}
		b.event(qa)
		ids = append(ids, qa.ID)
		if n == 2 {
			integrated := testEvent("integration", "integration")
			integrated.CandidateID = e.CandidateID
			integrated.ResultID = tech.ResultID
			integrated.Outcome = ptr("exact_effect")
			integrated.Positions = []Position{{"run/delivery", n*3 + 3}}
			b.event(integrated)
			ids = append(ids, integrated.ID)
		}
	}
	b.out.Runs = append(b.out.Runs, Run{ID: "run", Family: "goal_execute_v2", EventIDs: ids, Coverage: "complete", SourceIDs: []string{"source"}})
	r := b.finish()
	if err := ValidateReferences(r); err != nil {
		t.Fatal(err)
	}
	for id, n := range map[string]int{"recorded_runs": 1, "verifier_attempts": 3, "qualified_generations": 2, "distinct_observed_candidates": 3, "human_qa_records": 2} {
		expectCount(t, r, id, n)
	}
	for _, c := range r.Candidates {
		if c.OID == "c2" && (strings.Join(c.Axes["human_qa"].Outcomes, ",") != "fail" || c.Axes["integration"].State != "not_recorded") {
			t.Fatalf("C3 pass leaked to C2: %+v", c)
		}
	}
	analysisByID(t, r, "candidate_qa_loop/run")
	for _, a := range r.Analysis {
		if strings.HasPrefix(a.ID, "reported_qa/") && (a.Value == nil || *a.Value != 22367 || !strings.Contains(a.Definition, "not active")) {
			t.Fatalf("QA interval lost provenance: %+v", a)
		}
	}
}

func TestTraceSuffixAndReportClaimsAreNotVerifierAttemptsOrLoops(t *testing.T) {
	b := testBuilder()
	v := fixtureVerifier("c")
	e := testEvent("verify-00000002", "verification")
	e.CandidateID = b.candidate("c", "ctree")
	b.event(e)
	b.verifiers[e.ID] = v
	for _, id := range []string{"question-3", "question-4"} {
		e := testEvent(id, "delivery_report")
		e.EvidenceClass = "reported"
		e.Role = "agent"
		e.Data = raw(map[string]any{"question": "Can the installed candidate pass?", "meaning": "8 findings fixed; filemode rejection happened before verifier execution"})
		e.Positions = []Position{{"run/delivery", len(b.out.Events) + 1}}
		b.event(e)
	}
	working := testEvent("working", "delivery_report")
	working.Outcome = ptr("working")
	working.Positions = []Position{{"run/delivery", 4}}
	b.event(working)
	r := b.finish()
	expectCount(t, r, "verifier_attempts", 1)
	expectCount(t, r, "question_records", 2)
	expectCount(t, r, "explicit_handoff_links", 0)
	a := analysisByID(t, r, "verifier_elapsed/verify-00000002")
	if a.Value == nil || *a.Value != 87.293 {
		t.Fatalf("elapsed %+v", a)
	}
	for _, id := range []string{"question-3", "question-4", "working"} {
		e := eventByID(t, r, id)
		if e.OccurredAtUTC != nil || e.ReportedAtUTC != nil || e.RegisteredAtUTC != nil {
			t.Fatalf("untimed event acquired clock %+v", e)
		}
	}
	for _, a := range r.Analysis {
		if strings.HasPrefix(a.ID, "reported_response") {
			t.Fatal("working manufactured an answer")
		}
	}
}

func TestTraceRepeatBasisRequiresBoundIdenticalInputs(t *testing.T) {
	for _, different := range []string{"none", "unbound", "review", "candidate", "command"} {
		t.Run(different, func(t *testing.T) {
			b := testBuilder()
			for _, id := range []string{"v1", "v2"} {
				e := testEvent(id, "verification")
				b.event(e)
				v := fixtureVerifier("c")
				if id == "v2" {
					switch different {
					case "unbound":
						v.InputsBound = false
					case "review":
						v.Review.SHA256 = "different"
					case "candidate":
						v.CandidateOID = "new"
					case "command":
						v.Argv = []string{"other"}
					}
				}
				b.verifiers[id] = v
			}
			r := b.finish()
			same := 0
			for _, a := range r.Analysis {
				if strings.HasPrefix(a.ID, "same_verification_basis/") {
					same++
				}
			}
			want := 0
			if different == "none" {
				want = 1
			}
			if same != want {
				t.Fatalf("same basis patterns %d want %d", same, want)
			}
		})
	}
}

func TestTraceWaitResponseClocksAndExplicitRelations(t *testing.T) {
	for _, scenario := range []string{"known", "open", "reversed", "uncertain", "different_clock", "different_kind"} {
		t.Run(scenario, func(t *testing.T) {
			b := testBuilder()
			start := testEvent("request", "step_started")
			start.TimeBasis = taskjournal.TimeBasis{Kind: "reported", Clock: ptr("user clock"), Precision: "second"}
			start.ReportedAtUTC = ptr("2026-10-07T10:00:00Z")
			start.EvidenceClass = "reported"
			start.Data = raw(taskjournal.EventInput{Type: "step_started", Step: &taskjournal.StepRef{ID: "waiting", Kind: "waiting"}, Waiting: &taskjournal.Waiting{Reason: "review requested", Dependency: "d1", NextActor: "coordinator"}})
			b.event(start)
			if scenario != "open" {
				end := testEvent("answer", "step_finished")
				end.TimeBasis = start.TimeBasis
				end.ReportedAtUTC = ptr("2026-10-07T10:20:00Z")
				end.EvidenceClass = "reported"
				end.Data = raw(taskjournal.EventInput{Type: "step_finished", Step: &taskjournal.StepRef{ID: "waiting", Kind: "waiting"}})
				switch scenario {
				case "reversed":
					end.ReportedAtUTC = ptr("2026-10-07T09:00:00Z")
				case "uncertain":
					end.TimeBasis.Uncertainty = &taskjournal.Uncertainty{Earliest: "2026-10-07T10:00:00Z", Latest: "2026-10-07T11:00:00Z"}
				case "different_clock":
					end.TimeBasis.Clock = ptr("other clock")
				case "different_kind":
					end.TimeBasis.Kind = "observed"
					end.OccurredAtUTC = end.ReportedAtUTC
					end.ReportedAtUTC = nil
				}
				b.event(end)
				b.out.Relations = append(b.out.Relations, Relation{Type: "responds_to", From: end.ID, To: start.ID, Basis: "explicit journal relation", SourceIDs: []string{"source"}})
			}
			r := b.finish()
			a := analysisByID(t, r, "reported_wait/request")
			if scenario == "known" {
				if a.Value == nil || *a.Value != 1200 {
					t.Fatalf("interval %+v", a)
				}
			} else if a.Value != nil {
				t.Fatalf("uncertain time became interval %+v", a)
			}
			if !strings.Contains(a.Definition, "dependency: d1; next actor: coordinator") {
				t.Fatal("wait detail lost")
			}
		})
	}
}

func TestTraceConflictsAndPartialSourcesDoNotChooseWinners(t *testing.T) {
	b := testBuilder()
	cid := b.candidate("c", "tree")
	for i, outcome := range []string{"pass", "fail"} {
		e := testEvent(string(rune('a'+i)), "human_qa")
		e.CandidateID = cid
		e.Outcome = ptr(outcome)
		b.event(e)
	}
	bad := testEvent("missing", "human_qa")
	bad.SourceIDs = []string{"missing-source"}
	b.source(Source{ID: "missing-source", Kind: "fixture", Locator: "/missing", Status: "missing"})
	b.event(bad)
	b.diagnostic("missing", "Unknown sibling", bad.SourceIDs, nil)
	r := b.finish()
	expectCount(t, r, "human_qa_records", 2)
	if r.Candidates[0].Axes["human_qa"].State != "conflicting" {
		t.Fatal("conflicting QA picked winner")
	}
	if a := analysisByID(t, r, "human_qa_records"); a.Coverage != "lower_bound" {
		t.Fatalf("partial count %+v", a)
	}
}

func TestTraceJournalPreparationIsDeclarationNotCandidate(t *testing.T) {
	b := testBuilder()
	j := taskjournal.Snapshot{Sources: []taskjournal.Source{{SourceID: "source", Kind: "registry", Locator: "/fixture", Status: "valid"}}, Events: []taskjournal.Event{{EventID: "prep", Origin: "native", Type: "preparation", Candidate: &taskjournal.Candidate{OID: "base", Tree: "base-tree"}, SourceIDs: []string{"source"}, Data: raw(map[string]string{})}, {EventID: "lifecycle", Type: "lifecycle_changed", Origin: "native", Actor: &taskjournal.Actor{ID: "Recorder claim", Role: "tool"}, SourceIDs: []string{"source"}, Data: raw(map[string]string{})}}}
	b.journal(j, workspace.WorkItemRegistry{})
	r := b.finish()
	expectCount(t, r, "distinct_observed_candidates", 0)
	e := eventByID(t, r, "lifecycle")
	if e.Role != "unknown" || !strings.Contains(e.Recorder, "Recorder claim") {
		t.Fatalf("recorder became actor or disappeared %+v", e)
	}
}

func TestTraceSourceSequenceAndExactNativeDedupPreserveProjections(t *testing.T) {
	b := testBuilder()
	e := testEvent("same-native-id", "task_result_reported")
	e.Role = "unknown"
	e.Positions = []Position{{"journal", 1}}
	e.Data = raw(map[string]string{"native": "record"})
	b.event(e)
	e.Role = "agent"
	e.ActorClaim = ptr("exact owner")
	e.Positions = []Position{{"execute", 3}}
	e.Data = raw(map[string]any{"generation": 2})
	b.event(e)
	for id, pos := range map[string]int{"late": 2, "early": 1} {
		v := testEvent(id, "note")
		v.Positions = []Position{{"independent", pos}}
		b.event(v)
	}
	r := b.finish()
	got := eventByID(t, r, e.ID)
	if got.Role != "agent" || len(got.Positions) != 2 || !strings.Contains(string(got.Data), "generation") {
		t.Fatalf("dedup provenance %+v", got)
	}
	if len(r.Events) != 3 {
		t.Fatal("same identity duplicated")
	}
	for _, rel := range r.Relations {
		if rel.Type == "sequence" && (rel.From != "early" || rel.To != "late") {
			t.Fatalf("global chronology invented %+v", rel)
		}
	}
	if err := ValidateReferences(r); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(r)
	if err != nil || !json.Valid(encoded) {
		t.Fatal(err)
	}
}
