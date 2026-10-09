//go:build !windows

package taskrun

import (
	"fmt"
	"path/filepath"
)

// Only native receipts can resolve an interrupted native effect. This never
// re-launches a provider, pretest or an unobserved verification command.
func workerRecoverNative(d Dependencies, s workflowState) error {
	w := s.Result.Worker
	if w.Job == nil {
		return fmt.Errorf("worker recovery requires its exact reserved job")
	}
	in, err := workerIntent(s, *w.Job)
	if err != nil {
		return err
	}
	if err = workerNativeIntentMatches(s, in); err != nil {
		return err
	}
	a := s.Result.Delivery.Attempt
	if a == nil {
		return fmt.Errorf("native attempt has no reservation; unknown job remains blocked")
	}
	if a.CandidateOID != in.CandidateOID || a.CandidateTree != in.CandidateTree {
		return fmt.Errorf("native receipt candidate differs from the reserved worker attempt; preserve unknown work")
	}
	if s.Result.Worker.Stage == "verify" {
		if a.Kind != "verification" {
			return fmt.Errorf("native verification reservation differs")
		}
		if a.State != "attempted" {
			if _, err := workerRecordedVerification(s); err != nil {
				return err
			}
			if s.Result.Delivery.Phase == "automatic_acceptance_passed" {
				return workerStage(d, s, "integrate")
			}
			if feedback, ok := workerNativeTestFailure(s); ok {
				return workerCorrection(d, s, feedback)
			}
			return fmt.Errorf("native verification did not qualify this candidate")
		}
		// The native callback sees the reserved attempt and can only recover its
		// preserved completion spool; absent execution evidence remains blocked.
		o, err := WorkflowDeliveryVerifyWithBinary(d, s.Request.WorkspaceRoot, s.Result.RunID, s.Result.Paths.Context, s.Result.Worker.Review.Locator, s.Result.Worker.CandidateBinary)
		fresh, readErr := workflowRead(s.Request.WorkspaceRoot, s.Result.RunID)
		if readErr != nil {
			return readErr
		}
		if feedback, ok := workerNativeTestFailure(fresh); ok {
			return workerCorrection(d, fresh, feedback)
		}
		if err != nil {
			return err
		}
		if o.Delivery.Phase != "automatic_acceptance_passed" {
			return fmt.Errorf("recovered native test did not qualify")
		}
		return workerStage(d, fresh, "integrate")
	}
	if a.Kind != "integration" {
		return fmt.Errorf("native integration reservation differs; no unobserved effect will be repeated")
	}
	o, err := WorkflowDeliveryIntegrate(d, s.Request.WorkspaceRoot, s.Result.RunID, s.Result.Paths.Context)
	if err != nil {
		return err
	}
	if o.Delivery.Phase != "completed" {
		return fmt.Errorf("native integration remains incomplete")
	}
	return workerStage(d, s, "completed")
}

func workerNativeIntentMatches(s workflowState, in workerAttemptIntent) error {
	w := s.Result.Worker
	if (in.Kind != "verify" && in.Kind != "integrate") || in.Kind != w.Stage || in.CandidateOID != w.CandidateOID || in.CandidateTree != w.CandidateTree || in.Review == nil || w.Review == nil || *in.Review != *w.Review {
		return fmt.Errorf("native worker intent differs from its exact candidate and review")
	}
	if in.CandidateBinary == nil || in.CandidateBinary.Locator != w.CandidateBinary {
		return fmt.Errorf("native worker candidate executable binding differs")
	}
	if _, err := workflowBound(*in.CandidateBinary, 64<<20); err != nil {
		return err
	}
	return workerValidateReview(s, *in.Review, in.CandidateOID, in.CandidateTree)
}

func workerRecordedVerification(s workflowState) (deliveryVerificationReceipt, error) {
	w := s.Result.Worker
	if w.Review == nil || w.Job == nil {
		return deliveryVerificationReceipt{}, fmt.Errorf("worker has no current review/attempt")
	}
	in, err := workerIntent(s, *w.Job)
	if err != nil {
		return deliveryVerificationReceipt{}, err
	}
	if err = workerNativeIntentMatches(s, in); err != nil {
		return deliveryVerificationReceipt{}, err
	}
	r, err := deliveryRecoverVerification(s, w.Review.Locator)
	if err != nil {
		return r, err
	}
	if r.Review != *w.Review || r.Automatic == nil || r.Automatic.CandidateBinary.Locator != in.CandidateBinary.Locator || r.Automatic.CandidateBinary.SHA256 != in.CandidateBinary.SHA256 {
		return r, fmt.Errorf("native receipt belongs to another review or private executable")
	}
	return r, nil
}

func workerNativeTestFailure(s workflowState) (string, bool) {
	a := s.Result.Delivery.Attempt
	if a == nil || a.Kind != "verification" || a.State != "recorded" {
		return "", false
	}
	r, err := workerRecordedVerification(s)
	if err != nil || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.CandidateOID != s.Result.Worker.CandidateOID || r.CandidateTree != s.Result.Worker.CandidateTree || r.AttemptID != a.ID || r.Outcome != "fail" || r.Exit == nil || *r.Exit == 0 {
		return "", false
	}
	return "Native acceptance failed. Inspect " + filepath.Join(a.Path, "verification.json") + " and referenced logs. Demonstrate a failing regression, fix the concrete failure, and preserve the test contract.", true
}
