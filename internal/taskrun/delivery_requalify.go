package taskrun

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

// WorkflowDeliveryRequalify controls preserved command evidence for one exact
// attempt. It never grants a new acceptance invocation.
func WorkflowDeliveryRequalify(d Dependencies, root, id, contextPath, reviewPath, attemptID string) (WorkflowRun, error) {
	if !key(attemptID) {
		return WorkflowRun{}, workflowError(2, "--reuse requires one exact verification attempt ID")
	}
	return workflowDeliveryVerify(d, root, id, contextPath, reviewPath, attemptID)
}

func deliveryReusableVerification(d Dependencies, s workflowState, receipt deliveryVerificationReceipt) error {
	requireNew := func(reason string) error {
		return deliveryContinuityError("receipt_mismatch", reason+"; preserve this receipt and use ordinary verify to execute acceptance for the changed inputs")
	}
	if receipt.Exit == nil || *receipt.Exit != 0 || receipt.Error != "" {
		return requireNew("the preserved attempt has no successful observed command completion")
	}
	if receipt.Acceptance.Locator != s.Request.Delivery.AcceptancePath || receipt.CWD != s.Observed.Target.WorktreeLocator || receipt.Acceptance.SHA256 != receipt.AcceptanceSnapshot.SHA256 {
		return requireNew("the preserved acceptance entrypoint differs from the frozen delivery")
	}
	if err := deliveryContinuationAuthority(s); err != nil {
		return err
	}
	if err := workflowhandoff.ValidateDeliveryOwnerContinuity(d.Workflow, s.Result.Handoff.Locator); err != nil {
		return fmt.Errorf("requalify required Task publications: %w", err)
	}
	for _, b := range []FileBinding{receipt.Acceptance, receipt.AcceptanceSnapshot, receipt.Review, receipt.Stdout, receipt.Stderr} {
		if _, err := workflowBound(b, 4<<20); err != nil {
			return requireNew(fmt.Sprintf("receipt input %s no longer matches: %s", b.Locator, err))
		}
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
