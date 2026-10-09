package taskrun

import (
	"testing"
)

func reportCandidateQuestion(t *testing.T, f deliveryFixture, o WorkflowRun, id string) WorkflowRun {
	t.Helper()
	question := "Does this exact qualified candidate pass product QA?"
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: id, PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "needs_input", Summary: "The verified candidate awaits a human answer.", Meaning: "No candidate changes or additional effects.", Question: &question, Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	o, err := WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, id+".json", r))
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func TestQualifiedDeliveryRemainsUsableWhileAwaitingHumanInput(t *testing.T) {
	f, o := humanIntegrationCandidate(t, true)
	c := o.Delivery.Candidates[0]
	o = reportCandidateQuestion(t, f, o, "await-qa")
	if o.NextAction.Actor != "user" {
		t.Fatal("QA question was lost")
	}
	if _, err := PreviewHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID); err != nil {
		t.Fatalf("question invalidated qualified candidate: %v", err)
	}
	if _, err := ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic owner handover after QA question."); err != nil {
		t.Fatal(err)
	}
	if _, err := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, humanIntegrationTestInput(t, f, o, "pass", "question")); err != nil {
		t.Fatalf("human integration cannot record actual answer: %v", err)
	}
}

func TestQualifiedDeliveryQuestionDoesNotRestoreSupersededQualification(t *testing.T) {
	f, o := humanIntegrationCandidate(t, true)
	c := o.Delivery.Candidates[0]
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "correction", PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "working", Summary: "A discovered issue needs correction.", Meaning: "Candidate is being reconsidered.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	var err error
	o, err = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "correction.json", r))
	if err != nil {
		t.Fatal(err)
	}
	o = reportCandidateQuestion(t, f, o, "correction-question")
	if _, err = PreviewHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID); err == nil {
		t.Fatal("question restored qualification after correction started")
	}
}

func TestQualifiedDeliveryQuestionAllowsNativeQAAndPreservesItsOutcome(t *testing.T) {
	for _, outcome := range []string{"pass", "blocked", "fail"} {
		t.Run(outcome, func(t *testing.T) {
			f, o := humanIntegrationCandidate(t, true)
			o = reportCandidateQuestion(t, f, o, "before-native-qa")
			var err error
			o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, outcome, deliveryTestHuman(t, f, o, outcome))
			if err != nil {
				t.Fatalf("QA question blocked native answer: %v", err)
			}
			o = reportCandidateQuestion(t, f, o, "after-native-qa")
			_, err = PreviewHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, o.Delivery.Candidates[0].TaskResult.ID)
			if outcome == "fail" && err == nil || outcome != "fail" && err != nil {
				t.Fatalf("question changed %s QA qualification: %v", outcome, err)
			}
		})
	}
}
