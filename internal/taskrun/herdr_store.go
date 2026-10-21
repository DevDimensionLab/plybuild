package taskrun

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
)

func workflowRoundName(n int) string { return fmt.Sprintf("%03d", n) }
func workflowIndex(root, id string) string {
	return filepath.Join(workflowRoot(root), "requests", id+".json")
}
func workflowRead(root, id string) (workflowState, error) {
	var s workflowState
	if !regexp.MustCompile(`^wfr_[0-9a-f]{64}$`).MatchString(id) {
		return s, workflowError(2, "invalid run ID")
	}
	e := readPrivateValue(workflowIndex(root, id), privateRunState, 4<<20, &s)
	if e != nil {
		return s, e
	}
	if workflowID(s.Request) != id || s.Request.WorkspaceRoot != root || s.Result.RequestSHA256 != digest(s.Request) {
		return s, workflowError(4, "reservation binding changed")
	}
	var saved workflowState
	if headlessRun(s.Request) {
		if e = workerValidateIdentity(s); e != nil {
			return s, e
		}
	} else if s.Request.Worker != nil || s.Result.Worker != nil {
		return s, workflowError(4, "interactive reservation cannot acquire a worker")
	}
	e = readPrivateValue(filepath.Join(s.Result.Paths.RunRoot, "state.json"), privateRunState, 8<<20, &saved)
	if os.IsNotExist(e) {
		s.Result.Round.State = "unknown"
		return s, nil
	}
	if e != nil {
		return s, e
	}
	if !equal(s.Request, saved.Request) || !equal(s.CodexTrust, saved.CodexTrust) || !equal(s.ClaudeTrust, saved.ClaudeTrust) || saved.Result.RunID != id || !equal(s.Observed, saved.Observed) || saved.Result.Paths.RunRoot != s.Result.Paths.RunRoot || saved.Result.RequestSHA256 != digest(s.Request) || saved.Result.SessionID != "ply:"+id {
		return s, workflowError(4, "preserved state differs from reservation")
	}
	if _, e = workflowRecoveryRecord(saved); e != nil {
		return saved, e
	}
	if _, _, e = workflowOwnerResumeRecord(saved); e != nil {
		return saved, e
	}
	if e = workflowLoadAcceptanceSelection(&saved); e != nil {
		return saved, e
	}
	if e = deliveryLoadObservedStop(&saved); e != nil {
		return saved, e
	}
	if headlessRun(saved.Request) {
		if e = workerValidateIdentity(saved); e != nil {
			return saved, e
		}
	} else if saved.Result.Worker != nil {
		return saved, workflowError(4, "interactive state cannot acquire a worker")
	}
	if !deliveryRun(saved.Request) {
		if e = workflowLoadReportObservations(&saved); e != nil {
			return saved, e
		}
		if r, e := workflowReportRetirement(saved); e != nil {
			return saved, e
		} else if r != nil {
			binding := FileBinding{filepath.Join(saved.Result.Paths.RunRoot, "report-retirement.json"), digest(r)}
			saved.ReportRetirement = &binding
			saved.Result.Round.State, saved.Result.FinalReturn.State = "stopped", "retired"
		}
	}
	return saved, nil
}
func workflowSave(d Dependencies, s workflowState) error {
	// ActiveStop is a read projection of immutable events. Keeping it out of
	// stored state leaves historical frozen controls able to read that state.
	if s.Result.Delivery != nil {
		copy := *s.Result.Delivery
		copy.ActiveStop = nil
		s.Result.Delivery = &copy
	}
	return d.replaceValue(filepath.Join(s.Result.Paths.RunRoot, "state.json"), s)
}
func workflowUpdate(d Dependencies, root, id string, fn func(*workflowState) error) error {
	return withStore(root, func() error {
		s, e := workflowRead(root, id)
		if e != nil {
			return e
		}
		before := controlSemanticState(s)
		if e = fn(&s); e != nil {
			return e
		}
		if s.Request.ControlContract != nil && before != controlSemanticState(s) {
			s.ControlRevision++
		}
		return workflowSave(d, s)
	})
}
func workflowBound(b FileBinding, max int) ([]byte, error) {
	raw, e := readFile(b.Locator, max, true)
	if e != nil {
		return nil, e
	}
	if hash(raw) != b.SHA256 {
		return nil, workflowError(4, "artifact bytes changed: "+b.Locator)
	}
	return raw, nil
}
func workflowKeep(d Dependencies, path string, v any) (FileBinding, error) {
	b, e := Canonical(v)
	if e != nil {
		return FileBinding{}, e
	}
	if e = d.writeOnce(path, b); e != nil {
		return FileBinding{}, e
	}
	return FileBinding{path, hash(b)}, nil
}

// Both transports call this under the existing TaskRun -> workspace lock order.
// A reviewed report is not release evidence. Only a native TaskResult bound to
// the same handoff terminal can qualify a subsequent, separately authorized run.
func workflowReservations(d Dependencies, root string, target workspace.PlanWorktreeObservation) error {
	dir := filepath.Join(workflowRoot(root), "requests")
	if e := physical(dir, true); e != nil {
		return e
	}
	entries, e := os.ReadDir(dir)
	if os.IsNotExist(e) {
		return nil
	}
	if e != nil {
		return e
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			return workflowError(4, "unrecognized workflow reservation")
		}
		s, e := workflowRead(root, strings.TrimSuffix(entry.Name(), ".json"))
		if e != nil {
			return e
		}
		x := s.Observed.Target
		if x.WorktreeLocator != target.WorktreeLocator && !(x.GitCommonDir == target.GitCommonDir && x.Ref == target.Ref) {
			continue
		}
		if s.ReportRetirement != nil {
			if _, e := workflowReportRetirement(s); e != nil {
				return e
			}
			released, e := workflowReportQueueReleased(d, root, s.Request.PreparationID)
			if e != nil {
				return e
			}
			if !released {
				return workflowError(4, "retired report retains its reservation until exact queue release; retry the same retire input")
			}
			continue
		}
		released := false
		if deliveryRun(s.Request) {
			// A technical candidate exists while its owner is still active through
			// human QA and integration. It never releases the goal by itself.
			if s.Result.Delivery == nil || s.Result.Delivery.Phase != "completed" {
				return workflowError(4, "A delivery owner still owns this target through human QA and local integration")
			}
			if h := s.Result.Delivery.HumanIntegration; h != nil && h.Accepted {
				closed, err := workflowCompletedCloseout(root, s)
				if err != nil {
					return err
				}
				if closed == nil {
					return workflowError(4, "The accepted human integration retains this Task source until its exact closeout completes. Unrelated Task sources remain available.")
				}
			}
			// Completion ends this mandate, not the provider session. A later
			// start still needs the native qualified result below.
			if len(s.Result.Delivery.Candidates) > 0 {
				c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
				released = c.TaskResult.ID != "" && (c.Integration != nil || c.PullRequest != nil)
			}
		}
		if s.Result.FinalReturn.TerminalSHA256 != nil {
			registry, err := d.Workspace.WorkItems.Snapshot(root)
			if err != nil {
				return err
			}
			for _, r := range registry.TaskResults {
				if r.HandoffSHA256 == s.Result.Handoff.SHA256 && r.TerminalResultSHA256 == *s.Result.FinalReturn.TerminalSHA256 && s.StartSHA256 != nil && r.StartReceiptSHA256 == *s.StartSHA256 {
					released = true
				}
			}
		}
		if !released {
			return workflowError(4, "An active or unknown workflow run owns this target. A reviewed report does not attest provider inactivity; native TaskResult qualification is required before another start.")
		}
	}
	return nil
}

func workflowFresh(d Dependencies, s workflowState, target bool) error {
	runtime, err := workflowEffectiveRuntime(s)
	if err != nil {
		return err
	}
	return workflowFreshRuntime(d, s, target, runtime)
}
func workflowFreshRuntime(d Dependencies, s workflowState, target bool, runtime Runtime) error {
	continued, err := workflowContinuationRecord(s)
	if err != nil {
		return err
	}
	return workflowFreshRuntimeOperation(d, s, target, runtime, continued != nil && continued.RuntimeObservation != nil)
}

// An interactive report owner's native session is bound, so its launch has
// consumed the provider executable and any process-local trust grant. Report
// callbacks and observations never launch again; Delivery keeps its own
// continuation and headless workers still execute the bound control.
func workflowReportLaunched(s workflowState) bool {
	return !deliveryRun(s.Request) && !headlessRun(s.Request) && s.Result.Transport.AgentSessionID != ""
}

// Launch facts are history after workflowReportLaunched. Observed drift is
// logged and the report-only operation continues under unchanged authority.
func workflowReportLaunchDrift(s workflowState, fact, bound string, err error) {
	logrus.WithFields(logrus.Fields{"run_id": s.Result.RunID, "fact": fact, "bound": bound, "observed": err.Error(), "handling": "post-launch report operation continues under the unchanged mandate, session and report-only authority; a new launch still requires current facts"}).Warn("Report launch fact changed after launch")
}

func workflowFreshRuntimeOperation(d Dependencies, s workflowState, target bool, runtime Runtime, historicalLauncher bool) error {
	if workflowReportLaunched(s) {
		// Only a confirmed grant's observation can become history. A missing or
		// contradictory saved binding stays an integrity error.
		if s.Request.CodexProjectTrust == nil || s.CodexTrust == nil {
			if e := workflowTrustFresh(s); e != nil {
				return e
			}
		} else if e := workflowTrustFresh(s); e != nil {
			workflowReportLaunchDrift(s, "codex_project_trust", s.Request.CodexProjectTrust.RepositoryRoot, e)
		}
		for i, x := range []Executable{runtime.Executable, runtime.PlyExecutable} {
			if e := verifyExecutable(x); e != nil {
				workflowReportLaunchDrift(s, []string{"provider_executable", "ply_executable"}[i], x.Path+" "+x.SHA256, e)
			}
		}
		// Frozen policy evidence is an artifact, not an installed launcher.
		if _, e := runtimeOperationBindings(runtime, nil); e != nil {
			return e
		}
		if e := verifyExecutable(s.Request.Herdr.Executable); e != nil {
			return e
		}
		return workflowFreshArtifacts(d, s, target, true)
	}
	if e := workflowTrustFresh(s); e != nil {
		return e
	}
	if historicalLauncher {
		if e := deliveryAcceptedAuthority(s); e != nil {
			return e
		}
		if _, e := runtimeOperationBindings(runtime, []Executable{runtime.PlyExecutable}); e != nil {
			return deliveryContinuityError("runtime_dependency", "accepted delivery control or launch-contract evidence: "+e.Error())
		}
	} else if _, e := runtimeBindings(runtime); e != nil {
		if !deliveryRun(s.Request) {
			return e
		}
		if providerErr := verifyExecutable(runtime.Executable); providerErr != nil {
			if deliveryRun(s.Request) && deliveryStarted(s) {
				return deliveryContinuityError("runtime_observation_required", "accepted delivery's historical provider launcher is missing or changed at "+runtime.Executable.Path+"; use installed workflow execute continue with the original owner's current --runtime-evidence; the new installed provider is not authority")
			}
			return deliveryContinuityError("launch_dependency", "unaccepted runtime requires its original provider launcher: "+providerErr.Error())
		}
		return e
	}
	if headlessRun(s.Request) {
		return workflowFreshArtifacts(d, s, target, true)
	}
	if e := verifyExecutable(s.Request.Herdr.Executable); e != nil {
		if !deliveryRun(s.Request) {
			return e
		}
		return deliveryContinuityError("session_dependency", "live-session observation requires the bound Herdr executable: "+e.Error())
	}
	return workflowFreshArtifacts(d, s, target, true)
}

// Report-only validation keeps frozen inputs and receipts intact without
// claiming that unavailable runtime policy currently authorizes effects.
func workflowFreshArtifacts(d Dependencies, s workflowState, target, policy bool) error {
	if s.Result.Handoff.ID == "" {
		return workflowError(4, "reservation has no completed handoff binding; no automatic restart")
	}
	if _, e := workflowBound(FileBinding{s.Result.Handoff.Locator, s.Result.Handoff.SHA256}, 2<<20); e != nil {
		return e
	}
	facts, e := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, s.Result.Handoff.Locator)
	if e != nil {
		return e
	}
	if s.StartSHA256 != nil && (facts.Start == nil || facts.Start.SHA256 != *s.StartSHA256) {
		return workflowError(4, "native start receipt binding changed")
	}
	if s.Result.FinalReturn.TerminalSHA256 != nil && (facts.Terminal == nil || facts.Terminal.SHA256 != *s.Result.FinalReturn.TerminalSHA256) {
		return workflowError(4, "native terminal binding changed")
	}
	var draft struct {
		Inputs []struct {
			Locator string `json:"locator"`
			SHA256  string `json:"sha256"`
		} `json:"inputs"`
	}
	// The native validator already checked the complete input schema.
	if e := jsonUnmarshal(s.Request.HandoffDraft, &draft); e != nil {
		return e
	}
	for _, in := range draft.Inputs {
		b, e := readFile(in.Locator, 64<<20, false)
		if e != nil {
			return e
		}
		if hash(b) != in.SHA256 {
			return workflowError(4, "frozen Task input changed")
		}
	}
	if s.ContextSHA256 != "" {
		if _, e := workflowBound(FileBinding{s.Result.Paths.Context, s.ContextSHA256}, 1<<20); e != nil {
			return e
		}
	}
	for _, b := range []*FileBinding{s.Acceptance, s.StartDraft, s.ReportStop, s.ReportRetirement, s.ReportProviderStop} {
		if b != nil {
			if _, e := workflowBound(*b, 1<<20); e != nil {
				return e
			}
		}
	}
	for _, r := range s.Records {
		b, e := workflowBound(r.Report, 1<<20)
		if e != nil {
			return e
		}
		var report WorkflowReport
		if e = decode(b, 1<<20, &report); e != nil {
			return e
		}
		terminal, e := workflowBound(r.Terminal, 2<<20)
		if e != nil {
			return e
		}
		slotPath := filepath.Join(filepath.Dir(r.Report.Locator), "report-slot.json")
		if _, e = workflowBound(FileBinding{slotPath, digest(workflowReportSlot{report, terminal, r.TargetSHA256})}, 4<<20); e != nil {
			return e
		}
		if r.Review != nil {
			if _, e = workflowBound(*r.Review, 256<<10); e != nil {
				return e
			}
		}
		// Preserve freshness of every immutable report, including previous rounds.
		var artifacts []struct {
			Kind    string `json:"kind"`
			Locator string `json:"locator"`
			SHA256  string `json:"sha256"`
		}
		if e = jsonUnmarshal(report.Artifacts, &artifacts); e != nil {
			return e
		}
		for _, a := range artifacts {
			if a.Kind == "managed" {
				b, e := readFile(a.Locator, 64<<20, false)
				if e != nil {
					return e
				}
				if hash(b) != a.SHA256 {
					return workflowError(4, "reported artifact changed")
				}
			}
		}
	}
	if target && len(s.Records) > 0 {
		r := s.Records[len(s.Records)-1]
		if r.Number == s.Result.Round.Number {
			h, e := workflowTarget(d, s)
			if e != nil {
				return e
			}
			if h != r.TargetSHA256 {
				return workflowError(4, "candidate HEAD, identity, status or file contents changed since reporting")
			}
		}
	}
	if deliveryRun(s.Request) {
		return workflowDeliveryFreshArtifacts(d, s, target, policy)
	}
	return nil
}

func WorkflowShow(d Dependencies, root, id string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	s, e := workflowRead(root, id)
	if e != nil {
		return WorkflowRun{}, e
	}
	o := s.Result
	if !deliveryRun(s.Request) {
		protocol := s.Request.StartupProtocol
		if protocol == "" {
			protocol = ReportLegacyAcceptanceV1
		}
		o.ReportStartup = &ReportStartupStatus{Protocol: protocol, Stop: s.ReportStop, ProviderStop: s.ReportProviderStop, Retirement: s.ReportRetirement, SandboxCoverage: "unknown", RuntimeDecision: "per_action_not_observed"}
		o.ReportStartup.SourceSHA256, _ = workflowTarget(d, s)
		o.ReportStartup.EffectResolutions = []FileBinding{}
		for _, stop := range []*FileBinding{s.ReportStop, s.ReportProviderStop} {
			if stop == nil {
				continue
			}
			_, binding, err := readReportStopResolution(s, *stop)
			if err != nil {
				return o, err
			}
			if binding != nil {
				o.ReportStartup.EffectResolutions = append(o.ReportStartup.EffectResolutions, *binding)
			}
		}
		if reportReceipt(s.Request) && s.StartSHA256 != nil {
			o.ReportStartup.Receipt = s.Acceptance
		}
		if s.ReportRetirement != nil {
			o.NextAction = WorkflowAction{"coordinator", "Report owner retired with source/history and unresolved effects retained, without successful Task delivery. A replacement requires a fresh prepared Task mandate."}
		}
	}
	o.DeliveryStatus = deliveryStatus(d, s)
	if o.OwnerResume, e = ownerResumeReadback(s); e != nil {
		return o, e
	}
	// Project the choice from the bound request, including historical state
	// written before the optional readback field existed. No state rewrite.
	o.Provider = s.Request.Runtime.Provider
	o.Transport.Observation = "cached"
	if deliveryRun(s.Request) && o.Delivery != nil && o.Transport.AgentSessionID != "" && (s.Phase == "following" || deliveryStarted(s)) {
		// Keep the immutable startup evidence in saved state. Once the first
		// Task prompt is known to have arrived, these are history, not current
		// instructions to repair startup. Unrelated/current problems remain.
		o.Reasons = make([]Reason, 0, len(s.Result.Reasons))
		for _, reason := range s.Result.Reasons {
			switch reason.Code {
			case "herdr_start_response", "herdr_start_stopped", "herdr_bootstrap_response", "delivery_start_resume_stopped", "delivery_recovery_start_response":
				continue
			}
			o.Reasons = append(o.Reasons, reason)
		}
		if o.Round.State == "unknown" && o.Delivery.Phase == "awaiting_acceptance" && s.Phase == "following" {
			o.Round.State = o.Delivery.Phase
		}
	}
	if deliveryCurrentCandidateQualified(s) {
		// A later qualified candidate resolves earlier qualification failures.
		// Keep those attempts in raw state/events and keep every other reason.
		reasons := make([]Reason, 0, len(o.Reasons))
		for _, reason := range o.Reasons {
			if reason.Code != "delivery_candidate_not_qualified" {
				reasons = append(reasons, reason)
			}
		}
		o.Reasons = reasons
	}
	if deliveryStartupPending(s) {
		// Derive a useful recovery action for pre-prompt v2 starts saved by an
		// older control binary, without rewriting their frozen startup history.
		o.NextAction = deliveryStartupAction(s)
	}
	if !headlessRun(s.Request) && o.Transport.AgentSessionID == "" && o.NextAction.Actor == "recipient" {
		// Older saved attempts may still advise acceptance before session binding.
		// Correct the readback without rewriting those historical artifacts.
		o.NextAction = WorkflowAction{"coordinator", "Inspect the preserved startup and any native onboarding in the same tab; no native session is bound. Do not restart or resend input."}
	}
	if headlessRun(s.Request) {
		if e = workerReadControl(s, &o); e != nil {
			return o, e
		}
	}
	if closeout, closeoutErr := workflowTaskCloseout(root, s); closeoutErr != nil {
		o.Reasons = append(o.Reasons, Reason{"task_closeout_unavailable", closeoutErr.Error()})
	} else if closeout != nil {
		o.Closeout = closeout
		if closeout.State == "complete" {
			o.Round.State, o.FinalReturn.State = "completed", "completed"
			o.NextAction = workflowCloseoutAction(s, *closeout)
			return o, nil
		}
		reason, _, detail := workspace.TaskCloseoutPendingAction(*closeout)
		o.Reasons = append(o.Reasons, Reason{reason, detail})
		o.NextAction = workflowCloseoutAction(s, *closeout)
		if !deliveryObservedStop(s) && (s.Result.Delivery.Attempt == nil || s.Result.Delivery.Attempt.State != "attempted") && (closeout.WorktreeRemoved || workspace.TaskCloseoutSourceRemovalPending(*closeout)) {
			return o, nil
		}
	}
	if workflowReportStopped(s) {
		o.Round.State = "stopped"
		o.NextAction = workflowReportStopAction(s)
		if e = workflowFreshArtifacts(d, s, false, false); e != nil {
			o.Reasons = append(o.Reasons, Reason{"report_stop_artifact_drift", e.Error()})
		}
		return o, nil
	}
	if e = workflowFresh(d, s, true); e != nil {
		o.Reasons = append(o.Reasons, Reason{"workflow_run_drift", e.Error()})
		if status := o.DeliveryStatus; status != nil {
			if status.QualifiedCandidate != nil {
				status.QualifiedCandidate.Current = false
			}
			if status.Verification != nil {
				status.Verification.InputsMatch = false
			}
			if status.Acceptance != nil {
				status.Acceptance.Current = false
				if status.Acceptance.Outcome == "pass" {
					status.Acceptance.Outcome, status.Acceptance.Reason = "blocked", "The preserved candidate acceptance is not current: "+e.Error()
				}
			}
		}
		o.Round.State = "unknown"
		if o.FinalReturn.State == "accepted" {
			o.FinalReturn.State = "blocked"
		}
		o.NextAction = WorkflowAction{"coordinator", "Inspect the preserved binding drift; do not restart or delete run state."}
		if automaticCleanupPending(s) {
			// Integration is observed; the recorded cleanup cause stays current.
			o.NextAction = s.Result.NextAction
			if o.Closeout != nil {
				o.NextAction = workflowCloseoutAction(s, *o.Closeout)
			}
		}
		var dependency *Error
		if errors.As(e, &dependency) && dependency.Code == "delivery_runtime_observation_required" {
			o.NextAction = WorkflowAction{"recipient", "Use the installed Ply's workflow execute continue with this run, its unchanged --context and the original owner's actual current --runtime-evidence. Keep the same session and original controls."}
		}
		if action, ok := deliveryIncompleteReadback(s); ok {
			// The report was durably recorded even when runtime authority is
			// unavailable. Keep its actual question distinct from eligibility.
			o.Round.State, o.NextAction = s.Result.Delivery.Phase, action
		}
	} else if status := o.DeliveryStatus; status != nil && status.Verification != nil && status.Verification.Current && status.Verification.Outcome == "passed" && status.Verification.Qualification != "qualified" && o.NextAction.Actor == "recipient" && (o.Delivery.Phase == "working" || o.Delivery.Phase == "verifying") {
		// Older controls preserved successful checks but advised rerunning them
		// after qualification failed. Explain the supported continuation without
		// replacing a human question or recommending replay of an unknown check.
		o.NextAction = WorkflowAction{"recipient", "Acceptance succeeded; inspect the separate qualification failure. For a control compatibility update use the installed Ply's workflow execute continue with this run and its original context. Reuse unchanged successful evidence with verify --reuse " + status.Verification.AttemptID + "; changed source, review, script or authority requires new verification."}
	}
	if headlessRun(s.Request) {
		// Resolve current evidence first, so unavailable closeout or binding
		// diagnostics cannot be replaced with ordinary worker observation advice.
		o.NextAction = currentRunAction(d, s, o)
	}
	deliveryProjectStop(s, &o)
	if r, e := deliveryRetirement(s); e != nil {
		return o, e
	} else if r != nil {
		o.Round.State = "stopped"
		o.NextAction = WorkflowAction{"coordinator", "This owner is retired as history; its callbacks are closed and its queue slot is released. Start a fresh Task copy from the current Epic base under its own new mandate; reuse code only as reviewed input."}
	}
	return o, nil
}

func deliveryCurrentCandidateQualified(s workflowState) bool {
	d := s.Result.Delivery
	if !deliveryRun(s.Request) || d == nil || d.Attempt == nil || d.Attempt.State != "recorded" || len(d.Candidates) == 0 {
		return false
	}
	if !deliveryCandidateQualified(d) {
		return false
	}
	c := d.Candidates[len(d.Candidates)-1]
	if c.TaskResult.ID == "" || d.Attempt.CandidateOID != c.OID || d.Attempt.CandidateTree != c.Tree {
		return false
	}
	switch d.Attempt.Kind {
	case "verification":
		return d.Attempt.ID == c.Key
	case "human_qa":
		return c.HumanQA != nil && c.HumanQA.TaskResultID == c.TaskResult.ID
	case "human_integration":
		return c.HumanQA != nil && c.HumanQA.TaskResultID == c.TaskResult.ID
	case "integration":
		return c.Integration != nil && (d.Phase == "closing" || d.Phase == "completed")
	}
	return false
}
