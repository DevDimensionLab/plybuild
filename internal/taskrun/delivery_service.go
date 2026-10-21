package taskrun

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
)

func BuildDeliveryWorkflowRequest(in DeliveryRequestInput) (WorkflowRequest, error) {
	if err := validateNewDeliveryInstall(in); err != nil {
		return WorkflowRequest{}, err
	}
	r := WorkflowRequest{Envelope: deliveryEnv("herdr-run-request"), RequestKey: in.RequestKey, WorkspaceRoot: in.WorkspaceRoot, PreparationID: in.PreparationID, PreparationSHA256: in.PreparationSHA256, HandoffDraft: in.HandoffDraft, Runtime: in.Runtime, HumanAuthority: in.HumanAuthority, CodexProjectTrust: in.CodexProjectTrust, ReturnMode: "delivery_owner", Delivery: &in.Delivery}
	r.Herdr.Executable, r.Herdr.WorkspaceID, r.Herdr.TabLabel = in.HerdrExecutable, in.HerdrWorkspaceID, in.TabLabel
	r.ClaudeProjectTrust = in.ClaudeProjectTrust
	r.GoalStartEnvironment = in.GoalStartEnvironment
	if in.Worker != nil {
		r.Envelope = deliveryEnv("headless-run-request")
		r.Worker = in.Worker
		r.ControlContract = in.ControlContract
	}
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
	if r.ControlContract != nil && !headlessRun(r) {
		return workflowError(2, "bound control requires headless runtime")
	}
	if env := r.GoalStartEnvironment; env != nil {
		if r.Worker != nil || r.Runtime.Provider != "codex" || !filepath.IsAbs(env.Home) || !filepath.IsAbs(env.CodexHome) || filepath.Clean(env.Home) != env.Home || filepath.Clean(env.CodexHome) != env.CodexHome || env.Path == "" || strings.ContainsAny(env.Path, "\x00\r\n") {
			return workflowError(2, "bound start requires the exact Codex HOME, CODEX_HOME and PATH")
		}
	}
	if headlessRun(r) {
		return validateWorkerRequest(r)
	}
	if r.Worker != nil {
		return workflowError(2, "worker authority requires an explicit headless request")
	}
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
	if c.StartupProtocol != "" && (c.StartupProtocol != MandateReceiptV1 || headlessRun(r)) {
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
	if e := workflowValidateTrust(r); e != nil {
		return e
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

func deliveryInstructions(s workflowState) string {
	runtime, err := workflowEffectiveRuntime(s)
	if err != nil {
		return "Startup binding is unavailable; preserve the stop and do not write the Task."
	}
	c := s.Request.Delivery
	if mandateReceipt(s.Request) {
		return deliveryReceiptInstructions(s, runtime)
	}
	notification := "Use the authorized ply-agent-notify route unless the goal assigns notifications to a parent; internal delegates stay quiet."
	if c.NotificationContext != nil {
		notification = "Reuse the existing notification identity at " + c.NotificationContext.Locator + "; do not create a second binding. Internal delegates stay quiet."
	}
	return fmt.Sprintf("You own the selected delivery through its frozen completion boundary. Work in %s. Read the frozen goal %s and native mandate %s for scope and existing human authority. Next, read %s for the bound provider/session and run %s workflow execute begin %s --context %s --help; submit only your actual runtime observations through begin before Task writes. Ply handles identities, private evidence and hashes. Unknown permission stops dependent work; never change provider policy to make acceptance pass. The installed ply-delivery skill describes the normal workflow. Use the callback guide for verification, reporting and the exact delivery gates. Own meaningful tests, review and corrections; implement the acceptance entrypoint %s. Keep this session and the immutable control executable; do not restart or resend uncertain input. Local installation allowed: %t. Delegation allowed: %t, within the goal's constraints. %s A launch, candidate, human judgment and final delivery are separate facts.", s.Observed.Target.WorktreeLocator, c.Goal.Locator, filepath.Join(s.Result.Paths.RunRoot, "mandate.json"), filepath.Join(workflowDeliveryGuideDirectory(s), "callback-guide.md"), ShellQuote(runtime.PlyExecutable.Path), s.Result.RunID, ShellQuote(s.Result.Paths.Context), c.AcceptancePath, c.AllowLocalInstall, c.AllowSubagents, notification)
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
		// The owner's integrate may release the source while this follow is
		// already observing; the same follow then finishes source cleanup.
		if o, handled, e := deliveryFollowCleanup(d, root, id, deadline, &observation); handled {
			return o, e
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
				if _, e = workflowCall(d, s.Request, "tab", "rename", s.Result.Transport.TabID, "finished "+s.Request.Herdr.TabLabel); e != nil {
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

func deliveryReceiptInstructions(s workflowState, runtime Runtime) string {
	goal, err := readDeliveryReceiptGoal(s.Request.Delivery.Goal)
	if err != nil {
		return "Frozen goal unavailable; preserve the stop."
	}
	scope := ""
	if a := s.Request.Delivery.Agreement; a != nil {
		scope = fmt.Sprintf(" Authorized delivery: %s from %s to %s at %s. No other delivery target is authorized.", a.Mode, a.SourceRef, a.TargetRef, a.TargetWorktree)
	}
	return fmt.Sprintf("Your authorized assignment is %s: %s Work in %s. Receive this frozen mandate with %s workflow execute check-in %s --context %s --provider %s --native-session %s. This idempotent receipt records identity and mandate receipt; runtime reviews each action when it happens. A known refusal or unknown previous effect stops dependent work. Preserve provider controls and the same session. Edit and test this Task, make controlled commits on its source branch, and own separate review and the exact delivery.%s Local installation allowed: %t. Delegation allowed: %t. Read %s for complete scope and %s for callbacks, verification and gates. Implement acceptance at %s. The frozen goal is %s. Parent-owned notifications stay with the parent; otherwise use the existing authorized route. Receipt, actual work, technical acceptance, human judgment and delivery are separate facts.", goal.Title, goal.Objective, s.Observed.Target.WorktreeLocator, ShellQuote(runtime.PlyExecutable.Path), s.Result.RunID, ShellQuote(s.Result.Paths.Context), s.Request.Runtime.Provider, ShellQuote(s.Result.Transport.AgentSessionID), scope, s.Request.Delivery.AllowLocalInstall, s.Request.Delivery.AllowSubagents, filepath.Join(s.Result.Paths.RunRoot, "mandate.json"), filepath.Join(workflowDeliveryGuideDirectory(s), "callback-guide.md"), s.Request.Delivery.AcceptancePath, s.Request.Delivery.Goal.Locator)
}

type deliveryReceiptGoal struct {
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

// deliveryFollowCleanup finishes a released automatic delivery's source
// cleanup from the coordinator's checkout. After two settled observations the
// released owner's exact pane is closed, every Task writer must be absent, and
// the same recorded closeout continues. A working owner is only observed.
func deliveryFollowCleanup(d Dependencies, root, id string, deadline time.Time, observation *followObservation) (WorkflowRun, bool, error) {
	s, e := workflowRead(root, id)
	if e != nil {
		return WorkflowRun{}, false, nil
	}
	if s.Result.Delivery != nil && s.Result.Delivery.Phase == "completed" && s.Result.Delivery.OwnershipRelease != nil {
		if closeout, err := workflowTaskCloseout(root, s); err == nil && closeout != nil && closeout.WorktreeRemoved {
			o, err := WorkflowShow(d, root, id)
			return o, true, err
		}
	}
	quiet := 0
	for automaticCleanupPending(s) {
		if !time.Now().Before(deadline) {
			o, err := workflowFollowReadback(d, root, id, workflowObservationWindowExpired)
			return o, true, err
		}
		finish := false
		observed := (taskWriters{d, s}).herdr("")
		if observed.state == taskWritersUnknown {
			return deliveryReadbackHandled(d, root, id, observed.err)
		}
		if observed.state == taskWritersAbsent {
			finish = true
		} else {
			a, e := workflowAgentRead(d, s, deadline)
			if e == nil {
				if e = workflowPlace(s, a); e == nil {
					e = workflowSession(s, a, false)
				}
			}
			if e != nil {
				// A read failure cannot establish liveness or absence. In
				// particular, agent_not_found may race the preceding pane get.
				if err := workflowUpdate(d, root, id, func(current *workflowState) error {
					current.Result.Transport.State = "unknown"
					current.Result.Transport.Observation = "cached"
					return nil
				}); err != nil {
					return deliveryReadbackHandled(d, root, id, err)
				}
				quiet = 0
				var call *HerdrCallError
				if errors.As(e, &call) && call.ReadOnly && call.Phase == "agent get" && call.ProviderCode == "agent_not_found" {
					logrus.WithFields(logrus.Fields{"run_id": id, "source": "Herdr agent get", "read_error": e.Error(), "handling": "reconcile fresh writer evidence for the same closeout"}).Warn("Released owner observation unavailable; continuing writer inspection")
					fresh := (taskWriters{d, s}).herdr("")
					if fresh.state != taskWritersAbsent {
						return deliveryReadbackHandled(d, root, id, fresh.err)
					}
					// Resume rechecks local writers, authority and exact source
					// before any removal; this inventory alone grants no cleanup.
					finish = true
				} else {
					failure, resume := observation.failure(d, id, deadline, e)
					if !resume {
						o, err := workflowFollowReadback(d, root, id, failure)
						return o, true, err
					}
				}
			} else if workflowSettled(a) {
				quiet++
			} else {
				quiet = 0
			}
			if quiet >= 2 {
				if _, e = deliveryRetireClosePane(d, s); e == nil {
					_, _, e = deliveryRetireAbsence(d, s)
				}
				if e != nil {
					return deliveryReadbackHandled(d, root, id, e)
				}
				finish = true
			}
		}
		if finish {
			c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
			_, _, e = ResumeAutomaticLocalDelivery(d, root, id, c.TaskResult.ID)
			return deliveryReadbackHandled(d, root, id, e)
		}
		if !time.Now().Add(250 * time.Millisecond).Before(deadline) {
			return deliveryReadbackHandled(d, root, id, workflowError(5, "The released owner is still active; source cleanup waits until it settles. Observation timed out without stopping it"))
		}
		time.Sleep(250 * time.Millisecond)
		if s, e = workflowRead(root, id); e != nil {
			return deliveryReadbackHandled(d, root, id, e)
		}
	}
	return WorkflowRun{}, false, nil
}

func deliveryReadbackHandled(d Dependencies, root, id string, err error) (WorkflowRun, bool, error) {
	o, e := deliveryReadback(d, root, id, err)
	return o, true, e
}
