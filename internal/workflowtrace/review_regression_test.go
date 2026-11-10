package workflowtrace

import (
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestReviewChangedWaitingEvidenceDoesNotKeepNumericInterval(t *testing.T) {
	b := newBuilder("/fixture", workspace.TaskRecord{ID: "task", RepoID: "repo"}, "2026-10-07T12:00:00Z")
	b.source(Source{ID: "changed", Kind: "contribution", Locator: "/fixture/changed", Status: "changed"})
	for i, id := range []string{"start", "finish"} {
		kind, at := "step_started", "2026-10-07T10:00:00Z"
		if i == 1 {
			kind, at = "step_finished", "2026-10-07T10:20:00Z"
		}
		in := taskjournal.EventInput{Type: kind, Step: &taskjournal.StepRef{ID: "waiting", Kind: "waiting"}}
		if i == 0 {
			in.Waiting = &taskjournal.Waiting{Reason: "review", Dependency: "d", NextActor: "human"}
		}
		b.event(Event{ID: id, Type: kind, Role: "agent", EvidenceClass: "reported", SourceIDs: []string{"changed"}, ReportedAtUTC: ptr(at), TimeBasis: taskjournal.TimeBasis{Kind: "reported", Clock: ptr("fixture"), Precision: "second"}, Data: raw(in)})
	}
	b.diagnostic("source_changed", "Source changed after capture", []string{"changed"}, nil)
	result := b.finish()
	for _, analysis := range result.Analysis {
		if analysis.ID == "reported_wait/start" {
			if analysis.Coverage != "unknown" || analysis.Value != nil {
				t.Fatalf("changed evidence must leave the interval unknown and null: %+v", analysis)
			}
			return
		}
	}
	t.Fatal("waiting interval disappeared")
}

func TestReviewLegacyReportDedupKeepsPerformingActorAndQuestion(t *testing.T) {
	b := newBuilder("/fixture", workspace.TaskRecord{ID: "task", RepoID: "repo"}, "2026-10-07T12:00:00Z")
	b.source(Source{ID: "valid", Kind: "run_event", Locator: "/fixture/event", Status: "valid"})
	native := Event{ID: "same-native-report", Type: "run_report_received", Role: "ply", ActorClaim: ptr("ply"), Recorder: "task journal native projection", EvidenceClass: "registered", SourceIDs: []string{"valid"}, Data: raw(map[string]any{"run_event_type": "report_received"})}
	b.event(native)
	richer := native
	richer.Role, richer.ActorClaim, richer.EvidenceClass, richer.Outcome = "agent", ptr("recipient"), "reported", ptr("needs_input")
	richer.Data = raw(map[string]any{"question": "Which input is correct?"})
	b.event(richer)
	result := b.finish()
	event := result.Events[0]
	if event.Role != "agent" || event.EvidenceClass != "reported" || event.Outcome == nil || *event.Outcome != "needs_input" {
		t.Fatalf("complementary projection lost reported actor/outcome: %+v", event)
	}
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Code == "event_identity_conflict" {
			t.Fatalf("complementary projection was called a conflict: %+v", diagnostic)
		}
	}
	for _, analysis := range result.Analysis {
		if analysis.ID == "question_records" {
			if analysis.Value == nil || *analysis.Value != 1 {
				t.Fatalf("native report question was lost: %+v", analysis)
			}
			return
		}
	}
	t.Fatal("question count absent")
}
