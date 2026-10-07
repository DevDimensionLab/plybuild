package taskrun

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func BuildDeliveryWorkflowRequest(in DeliveryRequestInput) (WorkflowRequest, error) {
	r := WorkflowRequest{Envelope: deliveryEnv("herdr-run-request"), RequestKey: in.RequestKey, WorkspaceRoot: in.WorkspaceRoot, PreparationID: in.PreparationID, PreparationSHA256: in.PreparationSHA256, HandoffDraft: in.HandoffDraft, Runtime: in.Runtime, HumanAuthority: in.HumanAuthority, CodexProjectTrust: in.CodexProjectTrust, ReturnMode: "delivery_owner", Delivery: &in.Delivery}
	r.Herdr.Executable, r.Herdr.WorkspaceID, r.Herdr.TabLabel = in.HerdrExecutable, in.HerdrWorkspaceID, in.TabLabel
	r.ClaudeProjectTrust = in.ClaudeProjectTrust
	return r, validateDeliveryWorkflowRequest(r)
}

func validateDeliveryWorkflowRequest(r WorkflowRequest) error {
	if r.Envelope != deliveryEnv("herdr-run-request") || r.Delivery == nil || r.ReturnMode != "delivery_owner" || r.Agreement != (Agreement{}) || r.Coordinator.ActorClaim != "" || r.Coordinator.MayRequestChanges {
		return workflowError(2, "delivery-owner v2 requires its own contract and cannot inherit Agreement A or coordinator control")
	}
	if !key(r.RequestKey) || !strings.HasPrefix(r.PreparationID, "pre_") || len(r.PreparationID) != 68 || !digestPattern.MatchString("sha256:"+strings.TrimPrefix(r.PreparationID, "pre_")) || !digestPattern.MatchString(r.PreparationSHA256) || !r.HumanAuthority.Authorized || r.HumanAuthority.StartSurface != "human_authorized_herdr" || !plain(r.HumanAuthority.ActorClaim, 1, 256) || !plain(r.Herdr.WorkspaceID, 1, 128) || !plain(r.Herdr.TabLabel, 1, 80) {
		return workflowError(2, "invalid delivery authority, preparation or Herdr binding")
	}
	if e := physical(r.WorkspaceRoot, false); e != nil {
		return e
	}
	rt, p := r.Runtime, r.Runtime.PermissionBinding
	if rt.Mode != "interactive" || !plain(rt.Model, 1, 256) || (rt.Provider != "codex" && rt.Provider != "claude") || rt.ConfigProfile != nil && !plain(*rt.ConfigProfile, 1, 256) || rt.Provider == "claude" && (rt.ConfigProfile != nil || p.ProfileID != "manual" && p.ProfileID != "auto") {
		return workflowError(2, "delivery requires the explicitly selected interactive Codex or manual/auto Claude runtime")
	}
	if p.AuthorityKind != "reported_contract_with_effective_policy" && p.AuthorityKind != "launch_contract_pending_runtime_acceptance" || !plain(p.ProfileID, 1, 256) || !digestPattern.MatchString(p.EffectivePolicySHA256) || len(p.Evidence) == 0 {
		return workflowError(2, "delivery requires an explicit launch contract and runtime evidence; actual permission acceptance remains pending")
	}
	if e := deliveryEvidence(p.Evidence); e != nil {
		return e
	}
	for _, x := range []Executable{rt.Executable, rt.PlyExecutable, r.Herdr.Executable} {
		if !digestPattern.MatchString(x.SHA256) {
			return workflowError(2, "invalid delivery executable digest")
		}
		if e := verifyExecutable(x); e != nil {
			return e
		}
	}
	c := r.Delivery
	boundary := "after_human_pass"
	if c.Agreement != nil && c.Agreement.Mode == workspace.DeliveryPullRequest {
		boundary = "none"
	}
	if !plain(c.OwnerClaim, 1, 256) || c.LocalIntegration != boundary || !filepath.IsAbs(c.AcceptancePath) || filepath.Clean(c.AcceptancePath) != c.AcceptancePath || c.ReasoningEffort != "" && !plain(c.ReasoningEffort, 1, 64) {
		return workflowError(2, "invalid delivery owner, acceptance path or local integration boundary")
	}
	if e := physical(c.AcceptancePath, true); e != nil {
		return e
	}
	if _, e := workflowBound(c.Goal, 4<<20); e != nil {
		return e
	}
	if c.NotificationContext != nil {
		if _, e := workflowBound(*c.NotificationContext, 1<<20); e != nil {
			return e
		}
	}
	if e := workflowValidateTrust(r); e != nil {
		return e
	}
	_, e := workflowhandoff.ValidateDeliveryTaskRunDraft(r.HandoffDraft)
	if e == nil {
		e = workflowhandoff.ValidateDeliveryMandate(r.HandoffDraft, c.Agreement)
	}
	return e
}

func deliveryEvidence(evidence []Evidence) error {
	last := ""
	for _, x := range evidence {
		if x.Locator <= last || !digestPattern.MatchString(x.SHA256) || (x.Role != "effective_policy" && x.Role != "permission_proof" && x.Role != "runtime_contract") {
			return workflowError(2, "delivery policy evidence must be sorted, unique and typed")
		}
		if _, e := workflowBound(FileBinding{x.Locator, x.SHA256}, 1<<20); e != nil {
			return e
		}
		last = x.Locator
	}
	return nil
}

func deliveryReadback(d Dependencies, root, id string, operationErr error) (WorkflowRun, error) {
	out, readErr := WorkflowShow(d, root, id)
	if operationErr != nil {
		return out, operationErr
	}
	return out, readErr
}

func deliveryInstructions(s workflowState) string {
	r, o := s.Request, s.Result
	// The prompt path validates this binding before generating instructions.
	if runtime, err := workflowEffectiveRuntime(s); err == nil {
		r.Runtime = runtime
	}
	c := r.Delivery
	// The per-run guide supplies the exact callback schemas with unobserved
	// templates, so native recipients need not guess protocol fields.
	notify := "At start, use the installed ply-agent-notify skill to bind an existing authorized standing notification route, preserving the actual provider and Task cwd. Never invent a route or change global configuration. If this is an internal delegation whose goal assigns notifications to its parent, stay quiet and return to that parent."
	if c.NotificationContext != nil {
		notify = "You own the one notification context at " + c.NotificationContext.Locator + ". Use the installed ply-agent-notify skill for a necessary human answer, a real stop, or the agreed completion. Internal delegates stay quiet; preserve and reuse event identity."
	}
	notify += " Read the callback guide before accepting or reporting: " + filepath.Join(workflowDeliveryGuideDirectory(s), "callback-guide.md") + ". Copy its adjacent unobserved templates into new private files and fill actual observations."
	instructions := fmt.Sprintf("You own this selected delivery through local completion. Read the frozen goal %s, the native mandate %s, the immutable request %s and private context %s. The goal constrains the outcome; you own design, implementation, meaningful tests, review and fixes. Subagents allowed: %t. Local installation allowed: %t. These choices grant no new runtime permissions, remote effects or unrelated goals. Before target writes submit ply.workflow.run-acceptance schema_version 2 using %s workflow run accept %s --context %s --file <private-acceptance.json>. Include delivery_permission with the bound launch_contract_sha256, actual permission_confirmed and actual_policy_evidence; report actual runtime/session/policy, not requested facts as observations. Unknown authority stops dependent work. There is no inherited Agreement A or fixed correction budget. Implement the acceptance entrypoint at %s; do not treat its placeholder as a passed test. Preserve an actual review record and run %s workflow execute verify %s --context %s --review <review.json>. The verifier preserves execution evidence and qualifies a technical candidate; it does not claim human QA. Report work, necessary questions or incomplete outcomes with ply.workflow.delivery-report schema_version 2 via %s workflow execute report %s --context %s --file <report.json>. Never fabricate unrun verifier exits. After candidate qualification, prepare the actual installed human journey. Preserve the human's exact answer through execute qa; only an exact candidate pass permits execute integrate. Keep the same interactive session and own corrections after fail. Continue local integration and base update when the actual pass and original authority cover them. Technical candidate, human judgment and final delivery are separate events. Do not publish a second terminal into an existing candidate handoff. The bound control executable %s is immutable and separate from the installed candidate. Do not restart, resend uncertain input, or overwrite this control executable. %s", c.Goal.Locator, filepath.Join(o.Paths.RunRoot, "mandate.json"), workflowIndex(r.WorkspaceRoot, o.RunID), o.Paths.Context, c.AllowSubagents, c.AllowLocalInstall, ShellQuote(r.Runtime.PlyExecutable.Path), o.RunID, ShellQuote(o.Paths.Context), ShellQuote(c.AcceptancePath), ShellQuote(r.Runtime.PlyExecutable.Path), o.RunID, ShellQuote(o.Paths.Context), ShellQuote(r.Runtime.PlyExecutable.Path), o.RunID, ShellQuote(o.Paths.Context), ShellQuote(r.Runtime.PlyExecutable.Path), notify)
	if c.Agreement != nil {
		instructions = strings.ReplaceAll(instructions, "These choices grant no new runtime permissions, remote effects or unrelated goals.", "The frozen structured delivery agreement grants only its exact selected effects after actual runtime acceptance; no unrelated goals or implicit permissions.")
		instructions = strings.ReplaceAll(instructions, "You own this selected delivery through local completion.", "You own this selected delivery through its frozen completion boundary.")
		instructions += fmt.Sprintf(" Frozen delivery mode: %s. Source: %s. Target: %s. Required allowed_effects for actual permission acceptance: %s. delivery_permission.delivery_agreement_sha256 must bind %s. Keep mode, target and stop boundary unchanged.", c.Agreement.Mode, c.Agreement.SourceRef, c.Agreement.TargetRef, workflowJSON(workflowhandoff.DeliveryAllowedEffects(*c.Agreement)), workspace.DeliveryAgreementDigest(*c.Agreement))
		if c.Agreement.Mode == workspace.DeliveryPullRequest {
			instructions = strings.ReplaceAll(instructions, "only an exact candidate pass permits execute integrate.", "only an exact candidate pass permits the agreed PR publication through the delivery CLI; execute integrate is forbidden for this mode.")
			instructions = strings.ReplaceAll(instructions, "Continue local integration and base update when the actual pass and original authority cover them.", "Continue the agreed PR delivery and truthful queue closure after actual pass. Stop before merge; preserve the local target and Epic base.")
		}
	}
	return instructions
}

func workflowDeliveryFresh(d Dependencies, s workflowState, target bool) error {
	if s.Result.Delivery == nil {
		return workflowError(4, "delivery state is missing")
	}
	if _, e := workflowBound(s.Request.Delivery.Goal, 4<<20); e != nil {
		return e
	}
	if s.Request.Delivery.NotificationContext != nil {
		if _, e := workflowBound(*s.Request.Delivery.NotificationContext, 1<<20); e != nil {
			return e
		}
	}
	if e := deliveryEvidence(s.Result.Delivery.ActualPolicyEvidence); e != nil {
		return e
	}
	for _, event := range s.Result.Delivery.Events {
		if _, e := workflowBound(event.Binding, 4<<20); e != nil {
			return e
		}
	}
	for _, c := range s.Result.Delivery.Candidates {
		if _, e := workflowBound(c.Verification, 4<<20); e != nil {
			return e
		}
		if c.Integration != nil {
			if _, e := workflowBound(*c.Integration, 4<<20); e != nil {
				return e
			}
		}
		if c.PullRequest != nil {
			if _, e := workflowBound(*c.PullRequest, 4<<20); e != nil {
				return e
			}
		}
	}
	phase := s.Result.Delivery.Phase
	if target && len(s.Result.Delivery.Candidates) > 0 && (phase == "awaiting_human_qa" || phase == "human_qa_passed" || phase == "integrating" || phase == "completed") {
		c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		x, e := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref)
		if e != nil {
			return e
		}
		if !x.Clean || x.OID != c.OID || x.Tree != c.Tree || x.Ref != s.Observed.Target.Ref || x.GitCommonDir != s.Observed.Target.GitCommonDir {
			return workflowError(4, "delivery candidate changed; a human verdict cannot transfer to changed bytes")
		}
	}
	return nil
}

func workflowDeliveryFollow(d Dependencies, root, id string, timeout time.Duration) (WorkflowRun, error) {
	if timeout <= 0 {
		return WorkflowRun{}, workflowError(2, "timeout must be positive")
	}
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	deadline, quiet, previous := time.Now().Add(timeout), 0, ""
	observationTimeout := func() error {
		return workflowError(5, "Delivery continues in the same session; observation timed out without stopping or restarting it")
	}
	for {
		if !time.Now().Before(deadline) {
			return deliveryReadback(d, root, id, observationTimeout())
		}
		ready := false
		e := workflowUpdate(d, root, id, func(s *workflowState) error {
			if e := workflowFresh(d, *s, true); e != nil {
				return e
			}
			// Freshness checks or the interval between observations may consume
			// the observer's budget. That does not make the delivery unknown and
			// must not start a transport call with an already expired deadline.
			if !time.Now().Before(deadline) {
				return observationTimeout()
			}
			a, e := workflowAgentGetUntil(d, *s, false, deadline)
			if e != nil {
				return e
			}
			workflowObserve(d, s, a)
			phase := s.Result.Delivery.Phase
			binding := phase + workflowJSON(s.Result.Delivery.LastEventSHA256)
			if binding != previous || !workflowSettled(a) {
				quiet = 0
			}
			previous = binding
			if phase == "needs_input" || phase == "stopped" || phase == "awaiting_human_qa" || phase == "completed" {
				if workflowSettled(a) {
					quiet++
				}
			}
			ready = quiet >= 2
			return nil
		})
		if e != nil {
			o, _ := WorkflowShow(d, root, id)
			return o, e
		}
		if ready {
			s, e := workflowRead(root, id)
			if e != nil {
				return WorkflowRun{}, e
			}
			if s.Result.Delivery.Phase == "completed" {
				// A settled completed delivery can mark its preserved tab finished.
				// Waiting for a human retains the same active owner and tab label.
				if _, e = workflowCall(d, s.Request, "tab", "rename", s.Result.Transport.TabID, "finished "+s.Request.Herdr.TabLabel); e != nil {
					return deliveryReadback(d, root, id, e)
				}
			}
			o, e := WorkflowShow(d, root, id)
			o.Transport.Observation = "fresh"
			return o, e
		}
		pause := 250 * time.Millisecond
		if left := time.Until(deadline); left < pause {
			pause = left
		}
		if pause > 0 {
			time.Sleep(pause)
		}
	}
}
