package taskrun

import (
	"errors"
	"testing"
)

func TestAutomaticDeliveryReadbackKeepsResolvedQualificationInactive(t *testing.T) {
	for _, phase := range []string{"needs_input", "closing"} {
		t.Run(phase, func(t *testing.T) {
			f, o, binary, _ := autoGateFixture(t, false, "")
			hasRejection := func(run WorkflowRun) bool {
				for _, reason := range run.Reasons {
					if reason.Code == "delivery_candidate_not_qualified" {
						return true
					}
				}
				return false
			}
			badReview := workflowProviderDocument(t, deliveryTestReview(t, f, "unobserved"))
			badReview["decision"] = "unobserved"
			var err error
			o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "unobserved-review.json", badReview), binary)
			if err == nil || !hasRejection(o) {
				t.Fatalf("fixture did not preserve qualification rejection: %+v %v", o.Reasons, err)
			}
			o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "corrected"), binary)
			if err != nil || hasRejection(o) {
				t.Fatalf("corrected review did not resolve qualification: %+v %v", o.Reasons, err)
			}
			if phase == "needs_input" {
				o = reportCandidateQuestion(t, f, o, "qualified-question")
			} else {
				d := f.D
				d.Fault = func(point string) error {
					if point == "automatic_closeout_after_release" {
						return errors.New("fixture interruption before retained closeout")
					}
					return nil
				}
				o, err = WorkflowDeliveryIntegrate(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
				if err == nil || o.DeliveryStatus.FinalDelivery.State != "delivered" {
					t.Fatalf("fixture did not interrupt after integration: %+v %v", o.DeliveryStatus, err)
				}
			}
			if o.Delivery.Phase != phase {
				t.Fatalf("expected phase %s, got %s", phase, o.Delivery.Phase)
			}
			c := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
			if _, err = nativeDeliveryCandidate(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID); err != nil {
				t.Fatalf("candidate actually unqualified: %v", err)
			}
			if o.DeliveryStatus.Acceptance == nil || o.DeliveryStatus.Acceptance.Outcome != "pass" || !o.DeliveryStatus.Acceptance.Current {
				t.Fatalf("fixture lost current acceptance: %+v", o.DeliveryStatus.Acceptance)
			}
			if hasRejection(o) {
				t.Errorf("%s restored a resolved qualification rejection: %+v", phase, o.Reasons)
			}
			shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil || hasRejection(shown) {
				t.Errorf("readback restored a resolved qualification rejection: %+v %v", shown.Reasons, err)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil || !hasRejection(s.Result) {
				t.Fatalf("historical qualification rejection was removed: %v", err)
			}
		})
	}
}
