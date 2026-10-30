package taskrun

import (
	"fmt"
)

// WorkflowDeliveryRequalify controls preserved command evidence for one exact
// attempt. It never grants a new acceptance invocation.
func WorkflowDeliveryRequalify(d Dependencies, root, id, contextPath, attemptID string) (WorkflowRun, error) {
	if !key(attemptID) {
		return WorkflowRun{}, workflowError(2, "--reuse requires one exact verification attempt ID")
	}
	return workflowDeliveryVerify(d, root, id, contextPath, attemptID)
}

func deliveryReusableVerification(d Dependencies, s workflowState, receipt deliveryVerificationReceipt) error {
	if err := validateAutomaticReceipt(s, receipt, true); err != nil {
		return err
	}
	requireNew := func(reason string) error {
		return deliveryContinuityError("receipt_mismatch", reason+"; preserve this receipt and use ordinary verify to execute acceptance for the changed inputs")
	}
	if receipt.Exit == nil || receipt.Error != "" {
		return requireNew("the preserved attempt has no observed command completion")
	}
	if receipt.Acceptance.Locator != s.Request.Delivery.AcceptancePath || receipt.CWD != s.Observed.Target.WorktreeLocator || receipt.Acceptance.SHA256 != receipt.AcceptanceSnapshot.SHA256 {
		return requireNew("the preserved acceptance entrypoint differs from the frozen delivery")
	}
	if err := deliveryContinuationAuthority(s); err != nil {
		return err
	}
	if err := deliveryVerificationInputs(receipt); err != nil {
		return requireNew(err.Error())
	}
	x, err := deliveryTarget(d, s)
	if err != nil {
		return requireNew("source worktree no longer matches the completed check: " + err.Error())
	}
	if x.OID != receipt.CandidateOID || x.Tree != receipt.CandidateTree {
		return requireNew("source commit or tree differs from the completed check")
	}
	return nil
}

// These are the exact inputs checked by receipt reuse and by historical report
// qualification. Reading them grants no new verification, candidate or phase.
func deliveryVerificationInputs(receipt deliveryVerificationReceipt) error {
	if _, err := inputBound(receipt.Acceptance, inputAcceptance, 4<<20); err != nil {
		return fmt.Errorf("receipt input %s no longer matches: %w", receipt.Acceptance.Locator, err)
	}
	return deliveryPreservedVerificationInputs(receipt)
}

func deliveryPreservedVerificationInputs(receipt deliveryVerificationReceipt) error {
	for _, b := range []FileBinding{receipt.AcceptanceSnapshot, receipt.Stdout, receipt.Stderr} {
		if _, err := workflowBound(b, 4<<20); err != nil {
			return fmt.Errorf("receipt input %s no longer matches: %w", b.Locator, err)
		}
	}
	return nil
}
