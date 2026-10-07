package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
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
	e := readValue(workflowIndex(root, id), 4<<20, &s)
	if e != nil {
		return s, e
	}
	if workflowID(s.Request) != id || s.Request.WorkspaceRoot != root || s.Result.RequestSHA256 != digest(s.Request) {
		return s, workflowError(4, "reservation binding changed")
	}
	var saved workflowState
	e = readValue(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, &saved)
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
	return saved, nil
}
func workflowSave(d Dependencies, s workflowState) error {
	return d.replaceValue(filepath.Join(s.Result.Paths.RunRoot, "state.json"), s)
}
func workflowUpdate(d Dependencies, root, id string, fn func(*workflowState) error) error {
	return withStore(root, func() error {
		s, e := workflowRead(root, id)
		if e != nil {
			return e
		}
		if e = fn(&s); e != nil {
			return e
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
		released := false
		if deliveryRun(s.Request) {
			// A technical candidate exists while its owner is still active through
			// human QA and integration. It never releases the goal by itself.
			if s.Result.Delivery == nil || s.Result.Delivery.Phase != "completed" {
				return workflowError(4, "A delivery owner still owns this target through human QA and local integration")
			}
			// Completion ends this mandate, not the provider session. A later
			// start still needs the native qualified result below.
			if len(s.Result.Delivery.Candidates) > 0 {
				c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
				released = c.TaskResult.ID != "" && c.Integration != nil
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
	if e := workflowTrustFresh(s); e != nil {
		return e
	}
	if _, e := runtimeBindings(runtime); e != nil {
		return e
	}
	if e := verifyExecutable(s.Request.Herdr.Executable); e != nil {
		return e
	}
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
	for _, b := range []*FileBinding{s.Acceptance, s.StartDraft} {
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
		return workflowDeliveryFresh(d, s, target)
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
	// Project the choice from the bound request, including historical state
	// written before the optional readback field existed. No state rewrite.
	o.Provider = s.Request.Runtime.Provider
	o.Transport.Observation = "cached"
	if deliveryRun(s.Request) && o.Delivery != nil && o.Transport.AgentSessionID != "" && (s.Phase == "following" || s.StartSHA256 != nil && o.Delivery.PermissionState == "recipient_confirmed_contract") {
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
	if o.Transport.AgentSessionID == "" && o.NextAction.Actor == "recipient" {
		// Older saved attempts may still advise acceptance before session binding.
		// Correct the readback without rewriting those historical artifacts.
		o.NextAction = WorkflowAction{"coordinator", "Inspect the preserved startup and any native onboarding in the same tab; no native session is bound. Do not restart or resend input."}
	}
	if e = workflowFresh(d, s, true); e != nil {
		o.Reasons = append(o.Reasons, Reason{"workflow_run_drift", e.Error()})
		o.Round.State = "unknown"
		if o.FinalReturn.State == "accepted" {
			o.FinalReturn.State = "blocked"
		}
		o.NextAction = WorkflowAction{"coordinator", "Inspect the preserved binding drift; do not restart or delete run state."}
	}
	return o, nil
}

func deliveryCurrentCandidateQualified(s workflowState) bool {
	d := s.Result.Delivery
	if !deliveryRun(s.Request) || d == nil || d.Attempt == nil || d.Attempt.State != "recorded" || len(d.Candidates) == 0 {
		return false
	}
	switch d.Phase {
	case "awaiting_human_qa", "human_qa_passed", "integrating", "completed":
	default:
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
	case "integration":
		return c.Integration != nil && d.Phase == "completed"
	}
	return false
}
