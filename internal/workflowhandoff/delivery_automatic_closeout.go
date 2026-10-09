package workflowhandoff

import (
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// CloseAutomaticDeliveryCandidate closes only the lifecycle of an observed
// automatic local return after explicit native owner release. The source and
// branch are kept; this adapter has no merge or resource removal operation.
func CloseAutomaticDeliveryCandidate(d Dependencies, taskID string, resultID workspace.TaskResultID, authorization *workspace.DeliveryAuthorization, ownership workspace.TaskCloseoutOwnership, guard func() error) (workspace.TaskCloseoutReceipt, error) {
	var out workspace.TaskCloseoutReceipt
	if d.TaskWorkspace == nil || authorization == nil || !authorization.Agreement.AutomaticAcceptance() || workspace.ValidateDeliveryAuthorization(*authorization) != nil || guard == nil {
		return out, fmt.Errorf("automatic closeout requires the exact native automatic agreement and current ownership guard")
	}
	wd := workspace.WithTaskContentScope(*d.TaskWorkspace, workspace.TaskID(taskID))
	a := authorization.Agreement
	// The delivery already names an exact candidate and observed integration.
	// Older TaskResults remain evidence; they are not candidates for this closeout.
	in := workspace.TaskCloseoutInput{OperationID: "automatic-closeout/" + authorization.RunID + "/" + string(resultID), TaskID: workspace.TaskID(taskID), TaskResultID: resultID, ExpectedResultOID: ownership.ResultOID, ExpectedResultTree: ownership.ResultTree, TargetRef: a.TargetRef, Keep: true, Ownership: ownership, RequireIntegration: true}
	preview, err := workspace.PreviewTaskCloseout(wd, in)
	if err != nil {
		return out, err
	}
	p := preview.Plan
	if p.TaskID != workspace.TaskID(taskID) || p.TaskResultID != resultID || p.ResultOID != ownership.ResultOID || p.ResultTree != ownership.ResultTree || p.Source.Ref != a.SourceRef || p.TargetRef != a.TargetRef || p.ReturnLocator != a.TargetWorktree || !p.Keep {
		return out, fmt.Errorf("automatic closeout differs from the exact accepted candidate and local Epic target")
	}
	registry, err := wd.WorkItems.Snapshot(p.WorkspaceRoot)
	if err != nil {
		return out, err
	}
	matched := false
	for _, effect := range registry.IntegrationResults {
		if effect.RecoveryStatus != "complete" || (effect.Outcome != "exact_effect" && effect.Outcome != "already_integrated") {
			continue
		}
		for _, prior := range registry.IntegrationAuthorities {
			bound := prior.Plan.DeliveryAuthorization
			if prior.ID == effect.AuthorityID && prior.TaskID == p.TaskID && prior.TaskResultID == resultID && prior.Plan.Task.ResultOID == p.ResultOID && prior.Plan.Task.ResultTree == p.ResultTree && prior.Plan.Epic.ParentRef == a.TargetRef && bound != nil && bound.Agreement.AutomaticAcceptance() && workspace.DeliveryAgreementDigest(bound.Agreement) == workspace.DeliveryAgreementDigest(a) && bound.RunID == authorization.RunID && bound.CandidateRunID == authorization.CandidateRunID && bound.RequestSHA256 == authorization.RequestSHA256 && bound.MandateSHA256 == authorization.MandateSHA256 && canonicalEqual(bridgeValue(bound.AcceptanceAmendment), bridgeValue(authorization.AcceptanceAmendment)) {
				matched = true
			}
		}
	}
	if !matched {
		return out, fmt.Errorf("automatic closeout has no observed integration under the same frozen native agreement")
	}
	if !preview.Ready {
		return out, fmt.Errorf("automatic closeout pending: %s; %s", strings.Join(preview.Reasons, "; "), preview.NextAction)
	}
	return workspace.ApplyTaskCloseoutWithOwnershipGuard(wd, in, preview.PlanSHA256, guard)
}
