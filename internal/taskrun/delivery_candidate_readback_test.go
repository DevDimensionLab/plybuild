package taskrun

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeliveryQualifiedCandidateReadbackKeepsHistoryAndCurrentFailures(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "review correction")
	badReview := deliveryTestReview(t, f, "not-yet-reviewed")
	review := workflowProviderDocument(t, badReview)
	review["decision"] = "unobserved"
	badReview = writeAny(t, f.R.WorkspaceRoot, "unobserved-review.json", review)
	hasRejection := func(o WorkflowRun) bool {
		for _, reason := range o.Reasons {
			if reason.Code == "delivery_candidate_not_qualified" {
				return true
			}
		}
		return false
	}
	var err error
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, badReview)
	if err == nil || !hasRejection(o) || len(o.Delivery.Candidates) != 0 {
		t.Fatalf("unresolved review rejection missing: %+v %v", o, err)
	}
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "actually-reviewed"))
	if err != nil || len(o.Delivery.Candidates) != 1 {
		t.Fatalf("corrected review did not qualify candidate: %+v %v", o, err)
	}
	if hasRejection(o) {
		t.Error("resolved qualification rejection remains active after candidate qualification")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || !hasRejection(s.Result) {
		t.Fatalf("historical rejection was removed from persisted state: %v", err)
	}
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	raw, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	after, readErr := os.ReadFile(statePath)
	if err != nil || readErr != nil || !bytes.Equal(raw, after) || hasRejection(shown) {
		t.Errorf("readback changed history or presented a resolved rejection: %v %v", err, readErr)
	}
	for _, condition := range []string{"unknown_attempt", "different_candidate", "different_verification", "progress_only", "changed_inputs"} {
		t.Run(condition, func(t *testing.T) {
			changed := s
			delivery := *s.Result.Delivery
			attempt := *delivery.Attempt
			delivery.Attempt, changed.Result.Delivery = &attempt, &delivery
			switch condition {
			case "unknown_attempt":
				attempt.State = "attempted"
			case "different_candidate":
				attempt.CandidateOID = strings.Repeat("0", 40)
			case "different_verification":
				attempt.ID = "another-verification"
			case "progress_only":
				// Conversational progress cannot restore an older, resolved
				// review failure when the qualified inputs are unchanged.
				delivery.Phase = "working"
			case "changed_inputs":
				delivery.Phase = "working"
				original, err := os.ReadFile(f.AcceptancePath)
				if err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(f.AcceptancePath, append(append([]byte{}, original...), []byte("# unverified correction\n")...), 0600); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.WriteFile(f.AcceptancePath, original, 0600); err != nil {
						t.Error(err)
					}
				})
			}
			if err := workflowSave(f.D, changed); err != nil {
				t.Fatal(err)
			}
			shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if condition == "progress_only" {
				if err != nil || hasRejection(shown) || shown.DeliveryStatus.QualifiedCandidate == nil || !shown.DeliveryStatus.QualifiedCandidate.Current {
					t.Fatalf("progress restored an older rejection for unchanged qualification: %+v %v", shown, err)
				}
			} else if err != nil || !hasRejection(shown) {
				t.Fatalf("unresolved %s was hidden: %+v %v", condition, shown, err)
			}
			if condition == "changed_inputs" && (shown.DeliveryStatus.QualifiedCandidate == nil || shown.DeliveryStatus.QualifiedCandidate.Current || shown.DeliveryStatus.Verification == nil || shown.DeliveryStatus.Verification.InputsMatch) {
				t.Fatal("changed acceptance input retained current qualification")
			}
		})
	}
	if err = workflowSave(f.D, s); err != nil {
		t.Fatal(err)
	}
	// A later failed verification of the same bytes is a new current failure,
	// even though the previous qualified candidate remains in the history.
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, badReview)
	if err == nil || !hasRejection(o) || o.Delivery.Phase != "working" || len(o.Delivery.Candidates) != 1 {
		t.Fatalf("new qualification failure was hidden: %+v %v", o, err)
	}
}
