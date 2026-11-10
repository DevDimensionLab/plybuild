package workflowtrace

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
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
	for i, outcome := range []string{"exact_effect", "conflict"} {
		e := testEvent(string(rune('a'+i)), "integration")
		e.CandidateID = cid
		e.Outcome = ptr(outcome)
		b.event(e)
	}
	bad := testEvent("missing", "integration")
	bad.CandidateID = b.candidate("d", "tree")
	bad.SourceIDs = []string{"missing-source"}
	b.source(Source{ID: "missing-source", Kind: "fixture", Locator: "/missing", Status: "missing"})
	b.event(bad)
	b.diagnostic("missing", "Unknown sibling", bad.SourceIDs, nil)
	r := b.finish()
	expectCount(t, r, "distinct_observed_candidates", 1)
	if r.Candidates[candidateIndex(r, *cid)].Axes["integration"].State != "conflicting" {
		t.Fatal("conflicting integration picked winner")
	}
	if a := analysisByID(t, r, "distinct_observed_candidates"); a.Coverage != "lower_bound" {
		t.Fatalf("partial count %+v", a)
	}
}

func candidateIndex(r Result, id string) int {
	for i, c := range r.Candidates {
		if c.ID == id {
			return i
		}
	}
	return -1
}

func TestTraceNativeRecorderIsNotAnActor(t *testing.T) {
	b := testBuilder()
	j := taskjournal.Snapshot{Sources: []taskjournal.Source{{SourceID: "source", Kind: "registry", Locator: "/fixture", Status: "valid"}}, Events: []taskjournal.Event{{EventID: "lifecycle", Type: "lifecycle_changed", Origin: "native", Actor: &taskjournal.Actor{ID: "Recorder claim", Role: "tool"}, SourceIDs: []string{"source"}, Data: raw(map[string]string{})}}}
	b.journal(j, workspace.WorkItemRegistry{})
	r := b.finish()
	e := eventByID(t, r, "lifecycle")
	if e.Role != "unknown" || !strings.Contains(e.Recorder, "Recorder claim") {
		t.Fatalf("recorder became actor or disappeared %+v", e)
	}
}

func TestTraceSourceSequenceAndExactNativeDedupPreserveProjections(t *testing.T) {
	b := testBuilder()
	e := testEvent("same-native-id", "integration")
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
