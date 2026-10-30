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
