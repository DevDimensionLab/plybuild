package taskrun

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func BuildDeliveryWorkflowRequest(in DeliveryRequestInput) (WorkflowRequest, error) {
	if err := validateNewDeliveryInstall(in); err != nil {
		return WorkflowRequest{}, err
	}
	r := WorkflowRequest{Envelope: deliveryEnv("herdr-run-request"), RequestKey: in.RequestKey, WorkspaceRoot: in.WorkspaceRoot, PreparationID: in.PreparationID, PreparationSHA256: in.PreparationSHA256, HandoffDraft: in.HandoffDraft, Runtime: in.Runtime, HumanAuthority: in.HumanAuthority, ReturnMode: "delivery_owner", Delivery: &in.Delivery}
	if in.Delivery.AdmissionMode == workspace.AdmissionParallel {
		r.SchemaVersion = 3
	}
	r.Herdr.Executable, r.Herdr.WorkspaceID, r.Herdr.TabLabel = in.HerdrExecutable, in.HerdrWorkspaceID, in.TabLabel
	r.GoalStartEnvironment = in.GoalStartEnvironment
	return r, validateDeliveryWorkflowRequest(r)
}

// These checks apply only to creation. Historical frozen request reading keeps
// its original validation and authority, including omitted Goal fields.
func validateNewDeliveryInstall(in DeliveryRequestInput) error {
	install, err := workflowhandoff.DeliveryLocalInstallAllowed(in.HandoffDraft)
	if err != nil {
		return err
	}
	if install != in.Delivery.AllowLocalInstall {
		return workflowError(4, "local installation authority differs between request and handoff")
	}
	draft, err := workflowhandoff.ValidateDeliveryTaskRunDraft(in.HandoffDraft)
	if err != nil {
		return err
	}
	var handoff struct {
		Inputs []FileBinding `json:"inputs"`
	}
	if err := json.Unmarshal(in.HandoffDraft, &handoff); err != nil {
		return err
	}
	for _, input := range handoff.Inputs {
		if input.SHA256 != draft.Basis.Spec.ManifestSHA256 {
			continue
		}
		raw, err := workflowBound(input, 4<<20)
		if err != nil {
			return err
		}
		var spec struct {
			ContractKind     string                `json:"contract_kind"`
			GoalOrigin       workspace.TaskGoalRef `json:"goal_origin"`
			ExecutionRequest struct {
				Goal workspace.TaskGoalContract `json:"goal"`
			} `json:"execution_request"`
		}
		if err := json.Unmarshal(raw, &spec); err != nil {
			return err
		}
		if spec.ContractKind != "execution" {
			return nil
		}
		goal := spec.ExecutionRequest.Goal
		if goal.TaskID == "" || goal.Goal != spec.GoalOrigin || !digestPattern.MatchString(goal.Goal.Spec.ManifestSHA256) {
			return workflowError(4, "frozen execution Spec has an invalid structured Goal basis")
		}
		if install && !goal.Executor.LocalInstallAllowed() {
			return workflowError(4, "local installation was not selected by the frozen Goal executor")
		}
		// The ceiling always comes from the exact execution Spec, including
		// manual human-gated deliveries. Only legacy manual Goal locators may
		// still name a design document instead of that Goal's manifest.
		goalManifestRequired := in.Delivery.StartupProtocol == MandateReceiptV1 || in.Delivery.Agreement != nil && in.Delivery.Agreement.AutomaticAcceptance()
		if goalManifestRequired && in.Delivery.Goal.SHA256 != goal.Goal.Spec.ManifestSHA256 {
			return workflowError(4, "delivery Goal differs from the frozen execution Spec origin")
		}
		return nil
	}
	return workflowError(4, "frozen execution Spec input is missing")
}

func validateDeliveryWorkflowRequest(r WorkflowRequest) error {
	if env := r.GoalStartEnvironment; env != nil {
		if r.Runtime.Provider != "codex" || !filepath.IsAbs(env.Home) || !filepath.IsAbs(env.CodexHome) || filepath.Clean(env.Home) != env.Home || filepath.Clean(env.CodexHome) != env.CodexHome || env.Path == "" || strings.ContainsAny(env.Path, "\x00\r\n") {
			return workflowError(2, "bound start requires the exact Codex HOME, CODEX_HOME and PATH")
		}
	}
	expected := deliveryEnv("herdr-run-request")
	if r.SchemaVersion == 3 {
		expected.SchemaVersion = 3
	}
	if err := validateParallelRequest(r); err != nil {
		return err
	}
	if r.Envelope != expected || r.Delivery == nil || r.ReturnMode != "delivery_owner" || r.Agreement != (Agreement{}) || r.Coordinator.ActorClaim != "" || r.Coordinator.MayRequestChanges {
		return workflowError(2, "delivery-owner v2 requires its own contract and cannot inherit Agreement A or coordinator control")
	}
	if !key(r.RequestKey) || !strings.HasPrefix(r.PreparationID, "pre_") || len(r.PreparationID) != 68 || !digestPattern.MatchString("sha256:"+strings.TrimPrefix(r.PreparationID, "pre_")) || !digestPattern.MatchString(r.PreparationSHA256) || !r.HumanAuthority.Authorized || r.HumanAuthority.StartSurface != "human_authorized_herdr" || !plain(r.HumanAuthority.ActorClaim, 1, 256) || !plain(r.Herdr.WorkspaceID, 1, 128) || !plain(r.Herdr.TabLabel, 1, 80) {
		return workflowError(2, "invalid delivery authority, preparation or Herdr binding")
	}
	if e := physical(r.WorkspaceRoot, false); e != nil {
		return e
	}
	rt, p := r.Runtime, r.Runtime.PermissionBinding
	// The serialized observation remains well-formed history even if that
	// installation has since disappeared. It supplies no runtime authority.
	if x := rt.PlyExecutable; !digestPattern.MatchString(x.SHA256) || !filepath.IsAbs(x.Path) || filepath.Clean(x.Path) != x.Path {
		return workflowError(2, "invalid historical Ply observation")
	}
	if rt.Mode != "interactive" || !plain(rt.Model, 1, 256) || (rt.Provider != "codex" && rt.Provider != "claude") || rt.ConfigProfile != nil && !plain(*rt.ConfigProfile, 1, 256) || rt.Provider == "claude" && (rt.ConfigProfile != nil || p.ProfileID != "manual" && p.ProfileID != "auto") {
		return workflowError(2, "delivery requires the explicitly selected interactive Codex or manual/auto Claude runtime")
	}
	if p.AuthorityKind != "reported_contract_with_effective_policy" && p.AuthorityKind != "launch_contract_pending_runtime_acceptance" || !plain(p.ProfileID, 1, 256) || !digestPattern.MatchString(p.EffectivePolicySHA256) || len(p.Evidence) == 0 {
		return workflowError(2, "delivery requires an explicit launch contract and runtime evidence; actual permission acceptance remains pending")
	}
	if e := deliveryEvidence(p.Evidence); e != nil {
		return e
	}
	for _, x := range []Executable{rt.Executable, r.Herdr.Executable} {
		if !digestPattern.MatchString(x.SHA256) {
			return workflowError(2, "invalid delivery executable digest")
		}
		if e := verifyExecutable(x); e != nil {
			return e
		}
	}
	c := r.Delivery
	if c.StartupProtocol != "" && (c.StartupProtocol != MandateReceiptV1) {
		return workflowError(2, "unsupported startup protocol")
	}
	protocol, protocolErr := workflowhandoff.DeliveryStartupProtocol(r.HandoffDraft)
	if protocolErr != nil {
		return protocolErr
	}
	if mandateReceipt(r) {
		if _, err := readDeliveryReceiptGoal(c.Goal); err != nil {
			return err
		}
	}
	if protocol != c.StartupProtocol {
		return workflowError(4, "startup protocol differs from frozen mandate")
	}
	boundary := DeliveryLocalIntegrationBoundary(c.Agreement)
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
	_, e := workflowhandoff.ValidateDeliveryTaskRunDraft(r.HandoffDraft)
	if e == nil {
		e = workflowhandoff.ValidateDeliveryMandate(r.HandoffDraft, c.Agreement)
	}
	return e
}

func deliveryEvidence(evidence []Evidence) error { return deliveryEvidenceRoles(evidence, false) }
func deliveryEvidenceRoles(evidence []Evidence, observation bool) error {
	last := ""
	for _, x := range evidence {
		if x.Locator <= last || !digestPattern.MatchString(x.SHA256) || (x.Role != "effective_policy" && x.Role != "permission_proof" && x.Role != "runtime_contract" && !(observation && x.Role == "runtime_observation")) {
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

func workflowDeliveryFresh(d Dependencies, s workflowState, target bool) error {
	return workflowDeliveryFreshArtifacts(d, s, target, true)
}

func workflowDeliveryFreshArtifacts(d Dependencies, s workflowState, target, policy bool) error {
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
	if policy {
		if e := deliveryEvidenceRoles(s.Result.Delivery.ActualPolicyEvidence, mandateReceipt(s.Request)); e != nil {
			return e
		}
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
	if target && len(s.Result.Delivery.Candidates) > 0 && (phase == "awaiting_human_qa" || phase == "human_qa_passed" || phase == "automatic_acceptance_passed" || phase == "integrating" || phase == "closing" || phase == "completed") {
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
	observation := followObservation{}
	observationReadback := func(err error) (WorkflowRun, error) {
		return workflowFollowReadback(d, root, id, err)
	}
	for {
		if !time.Now().Before(deadline) {
			return observationReadback(workflowObservationWindowExpired)
		}
		ready := false
		var observationErr error
		continueRead := false
		e := workflowUpdate(d, root, id, func(s *workflowState) error {
			if e := workflowFresh(d, *s, true); e != nil {
				return e
			}
			// Freshness checks or the interval between observations may consume
			// the observer's budget. That does not make the delivery unknown and
			// must not start a transport call with an already expired deadline.
			if !time.Now().Before(deadline) {
				return workflowObservationWindowExpired
			}
			// Follow only observes the bound session. A transient read cwd is
			// not a Task mutation; owner callbacks recheck the live cwd.
			a, e := workflowAgentRead(d, *s, deadline)
			if e == nil {
				if e = workflowPlace(*s, a); e == nil {
					e = workflowSession(*s, a, false)
				}
			}
			if e != nil {
				// A failed read supplies no current liveness, including when our
				// own shorter window killed it. Preserve delivery and owner facts.
				s.Result.Transport.State = "unknown"
				s.Result.Transport.Observation = "cached"
				quiet = 0
				observationErr, continueRead = observation.failure(d, id, deadline, e)
				return nil
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
			if errors.Is(e, workflowObservationWindowExpired) {
				return observationReadback(e)
			}
			return deliveryReadback(d, root, id, e)
		}
		if observationErr != nil && !continueRead {
			return observationReadback(observationErr)
		}
		if ready {
			s, e := workflowRead(root, id)
			if e != nil {
				return WorkflowRun{}, e
			}
			if s.Result.Delivery.Phase == "completed" {
				// A settled completed delivery can mark its preserved tab finished.
				// Waiting for a human retains the same active owner and tab label.
				if _, e = workflowCall(d, s.Request, "tab", "rename", s.Result.Transport.TabID, workflowTaskTabLabel(s, "finished")); e != nil {
					return deliveryReadback(d, root, id, e)
				}
			}
			o, e := WorkflowShow(d, root, id)
			o.Transport.Observation = "fresh"
			return o, e
		}
		if time.Until(deadline) <= workflowFollowInterval {
			return observationReadback(workflowObservationWindowExpired)
		}
		time.Sleep(workflowFollowInterval)
	}
}

type deliveryReceiptGoal struct {
	TaskID    string `json:"task_id"`
	Title     string `json:"title"`
	Objective string `json:"objective"`
}

func readDeliveryReceiptGoal(binding FileBinding) (deliveryReceiptGoal, error) {
	var goal deliveryReceiptGoal
	raw, err := workflowBound(binding, 4<<20)
	if err != nil {
		return goal, err
	}
	if err = json.Unmarshal(raw, &goal); err != nil {
		return goal, err
	}
	if strings.TrimSpace(goal.Title) == "" || strings.TrimSpace(goal.Objective) == "" {
		return goal, workflowError(2, "mandate receipt requires the frozen goal title and objective")
	}
	return goal, nil
}
