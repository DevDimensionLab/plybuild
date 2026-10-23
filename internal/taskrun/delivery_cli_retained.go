package taskrun

import (
	"encoding/json"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Read the native integration ledger, release and target without resuming
// closeout. Writer absence is irrelevant to copying an immutable snapshot.
func qualifiedRetainedCLIState(d Dependencies, s workflowState, c DeliveryCandidate, integrated workflowhandoff.DeliveryIntegrationResult, registry workspace.WorkItemRegistry) error {
	a := deliveryEffectiveAgreement(s)
	if a.HumanOwnedIntegration() || s.Result.Delivery.HumanIntegration != nil {
		return workflowError(4, "retained-source installation requires automatic local delivery")
	}
	if attempt := s.Result.Delivery.Attempt; attempt != nil && attempt.State != "recorded" {
		return workflowError(4, "retained-source installation has an unresolved delivery effect")
	}
	for _, attempt := range registry.IntegrationAttempts {
		for _, authority := range registry.IntegrationAuthorities {
			if attempt.AuthorityID != authority.ID || authority.Plan.Epic.ParentLocator != a.TargetWorktree {
				continue
			}
			resolved := false
			for _, result := range registry.IntegrationResults {
				if attempt.ResultID != nil && result.ID == *attempt.ResultID && result.AttemptID != nil && *result.AttemptID == attempt.ID && (result.RecoveryStatus == "complete" || result.Outcome == "no_effect" && result.GitChanged != nil && !*result.GitChanged) {
					resolved = true
				}
			}
			if !resolved {
				return workflowError(4, "retained-source installation has an unresolved target integration attempt")
			}
		}
	}
	release, err := deliveryRelease(s, c)
	if err != nil {
		return err
	}
	if release == nil {
		return workflowError(4, "retained-source installation requires automatic ownership release")
	}
	raw, err := workflowBound(*release, 1<<20)
	if err != nil {
		return err
	}
	var r deliveryOwnershipRelease
	if err = decode(raw, 1<<20, &r); err != nil {
		return err
	}
	if r.Origin != "automatic_delivery_completed" {
		return workflowError(4, "retained-source installation rejects manual ownership release")
	}
	if integrated.Closeout != nil && integrated.Closeout.State == "complete" && s.Result.Delivery.Phase != "completed" {
		return workflowError(4, "native closeout and Run completion differ")
	}
	target, err := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(a.TargetWorktree, a.TargetRef)
	if err != nil {
		return err
	}
	if target.Locator != a.TargetWorktree || target.Ref != a.TargetRef || target.OID != c.OID || target.Tree != c.Tree || target.GitCommonDir != s.Observed.Target.GitCommonDir || !target.Clean || !target.Symbolic || len(target.InProgress) != 0 {
		return workflowError(4, "retained-source installation requires the current clean exact delivery target")
	}
	base := workspace.EpicBaseVersion{}
	for _, b := range registry.EpicBaseVersions {
		if b.ProjectID == a.ProjectID && b.RepoID == a.RepoID && b.EpicID == a.EpicID && b.Revision > base.Revision {
			base = b
		}
	}
	if base.OID != c.OID || base.Tree != c.Tree || base.Locator != a.TargetWorktree || base.Ref != a.TargetRef || base.GitCommonDir != target.GitCommonDir {
		return workflowError(4, "retained-source installation requires the exact current Epic base")
	}
	closed := false
	for _, event := range registry.TaskQueueEvents {
		if event.Kind == "advance" && event.Request["preparation_id"] == s.Request.PreparationID {
			closed = true
		}
	}
	if !closed {
		return workflowError(4, "retained-source installation requires observed Task queue closure")
	}
	// Match the recorded completion to the durable native ledger. Reusing
	// live integration authorization here would incorrectly reopen released
	// ownership; installation only observes the already completed effect.
	raw, err = Canonical(integrated.Integration.Readback.Value)
	if err != nil {
		return err
	}
	var ledger struct {
		Authority *workspace.IntegrationAuthority `json:"authority"`
		Intent    *workspace.IntegrationIntent    `json:"intent"`
		Result    *workspace.IntegrationResult    `json:"integration_result"`
	}
	if err = json.Unmarshal(traceWireJSON(raw), &ledger); err != nil {
		return err
	}
	if ledger.Authority == nil || ledger.Intent == nil || ledger.Result == nil || ledger.Authority.TaskResultID != c.TaskResult.ID || ledger.Authority.TaskID != c.TaskResult.TaskID || ledger.Authority.DeliveryOwner == nil || ledger.Authority.DeliveryOwner.RunID != s.Result.RunID || ledger.Authority.DeliveryOwner.RequestSHA256 != s.Result.RequestSHA256 || ledger.Authority.Plan.Task.ResultOID != c.OID || ledger.Authority.Plan.Task.ResultTree != c.Tree || ledger.Authority.Plan.Epic.ParentLocator != a.TargetWorktree || ledger.Authority.Plan.Epic.ParentRef != a.TargetRef || ledger.Intent.AuthorityID != ledger.Authority.ID || ledger.Result.AuthorityID != ledger.Authority.ID || ledger.Result.IntentID != ledger.Intent.ID || ledger.Result.RecoveryStatus != "complete" {
		return workflowError(4, "retained-source installation requires the exact native integration ledger")
	}
	authorityFound, intentFound, resultFound := false, false, false
	for _, value := range registry.IntegrationAuthorities {
		if equal(value, *ledger.Authority) {
			authorityFound = true
		}
	}
	for _, value := range registry.IntegrationIntents {
		if equal(value, *ledger.Intent) {
			intentFound = true
		}
	}
	for _, value := range registry.IntegrationResults {
		if equal(value, *ledger.Result) {
			resultFound = true
		}
	}
	if !authorityFound || !intentFound || !resultFound {
		return workflowError(4, "recorded delivery differs from the native integration ledger")
	}
	return nil
}
