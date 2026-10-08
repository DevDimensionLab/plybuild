package taskrun

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func positiveDeliveryClaim(a Acceptance, permission DeliveryPermissionAcceptance, r WorkflowRequest) error {
	c := a.RuntimeClaim
	// The request's profile is a launch selector. The actual runtime may expose
	// another label (or none), particularly after native permission onboarding.
	// Authority comes from confirmed, bound actual policy and validated scope,
	// never from equating the requested selector with an observed profile name.
	if permission.LaunchContractSHA256 != r.Runtime.PermissionBinding.EffectivePolicySHA256 || !permission.PermissionConfirmed || c.RuntimeID == nil || *c.RuntimeID != r.Runtime.Provider || c.EffectivePolicySHA256 == nil || c.ModelID != nil && !deliveryModelMatches(r.Runtime.Provider, r.Runtime.Model, *c.ModelID) || len(permission.ActualPolicyEvidence) == 0 {
		return workflowError(4, "delivery runtime and actual permissions are not confirmed against the launch contract")
	}
	if e := deliveryEvidence(permission.ActualPolicyEvidence); e != nil {
		return e
	}
	bound := false
	for _, evidence := range permission.ActualPolicyEvidence {
		if evidence.SHA256 == *c.EffectivePolicySHA256 {
			bound = true
		}
	}
	if !bound {
		return workflowError(4, "actual policy digest is not bound by recipient evidence")
	}
	if r.Delivery != nil && r.Delivery.Agreement != nil {
		agreement := *r.Delivery.Agreement
		if permission.DeliveryAgreementSHA256 != workspace.DeliveryAgreementDigest(agreement) || !slices.Equal(permission.AllowedEffects, workflowhandoff.DeliveryAllowedEffects(agreement)) {
			return workflowError(4, "actual permission acceptance does not cover the exact frozen delivery mode, target and effects")
		}
	} else if permission.DeliveryAgreementSHA256 != "" || len(permission.AllowedEffects) != 0 {
		return workflowError(4, "historical delivery cannot acquire new delivery authority")
	}
	return nil
}

func deliveryModelMatches(provider, requested, actual string) bool {
	if requested == actual {
		return true
	}
	// Claude --help documents family aliases separately from full model names.
	// Resolve only these bounded families, never fuzzy substrings or cross-family
	// fallbacks. Both requested and observed names remain preserved in readback.
	if provider != "claude" || requested != "opus" && requested != "sonnet" && requested != "haiku" {
		return false
	}
	prefix := "claude-" + requested + "-"
	if !strings.HasPrefix(actual, prefix) {
		return false
	}
	version := strings.TrimPrefix(actual, prefix)
	if len(version) == 0 || version[0] < '0' || version[0] > '9' {
		return false
	}
	for _, c := range version {
		if c != '-' && c != '.' && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

func deliveryCallback(d Dependencies, s workflowState, contextPath string, started bool) error {
	if !deliveryRun(s.Request) || s.Result.Delivery == nil {
		return workflowError(4, "run is not a delivery-owner v2 execution")
	}
	if s.Result.Delivery.OwnershipRelease != nil {
		return workflowError(4, "owner_released: developer callbacks are closed; continue human integration or return correction ownership through a negative human decision")
	}
	if e := workflowCallback(d, s, contextPath); e != nil {
		return e
	}
	if s.Result.Delivery.Phase == "completed" {
		return workflowError(4, "this delivery mandate is complete; further work requires a new selected goal")
	}
	if started && (s.StartSHA256 == nil || s.Result.Delivery.PermissionState != "recipient_confirmed_contract") {
		return workflowError(4, "delivery requires actual positive runtime acceptance before effects")
	}
	return deliveryLiveOwner(d, s)
}

func deliveryTarget(d Dependencies, s workflowState) (workspace.IntegrationWorktreeObservation, error) {
	x, e := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref)
	if e != nil {
		return x, e
	}
	if !x.Clean || x.Ref != s.Observed.Target.Ref || x.GitCommonDir != s.Observed.Target.GitCommonDir || len(x.InProgress) != 0 {
		return x, workflowError(4, "delivery requires the exact clean Task worktree and no in-progress Git operation")
	}
	return x, nil
}

func deliveryAppendEvent(d Dependencies, s *workflowState, id, kind string, value any) (FileBinding, error) {
	raw, e := Canonical(value)
	if e != nil {
		return FileBinding{}, e
	}
	for _, event := range s.Result.Delivery.Events {
		if event.ID == id {
			if event.Kind != kind || event.Binding.SHA256 != hash(raw) {
				return FileBinding{}, workflowError(4, "delivery event ID already has different immutable bytes")
			}
			return event.Binding, nil
		}
	}
	path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "events", fmt.Sprintf("%08d.json", len(s.Result.Delivery.Events)+1))
	binding, e := workflowKeep(d, path, value)
	if e != nil {
		return binding, e
	}
	s.Result.Delivery.Events = append(s.Result.Delivery.Events, DeliveryEventRecord{id, kind, binding})
	s.Result.Delivery.LastEventSHA256 = &binding.SHA256
	return binding, nil
}

func validateDeliveryReport(r DeliveryReport) error {
	if err := validateDeliveryReportContent(r); err != nil {
		return err
	}
	for _, b := range r.Evidence {
		if _, e := workflowBound(b, 64<<20); e != nil {
			return e
		}
	}
	for _, v := range r.VerifierResults {
		for _, b := range v.Evidence {
			if _, e := workflowBound(b, 64<<20); e != nil {
				return e
			}
		}
	}
	return nil
}

// The historical inventory validates the preserved report itself. Referenced
// working artifacts are checked at publication, not reread as live inputs by a
// bulk status query. They never qualify native results in that projection.
func validateDeliveryReportContent(r DeliveryReport) error {
	if r.Envelope != deliveryEnv("delivery-report") || !key(r.EventID) || !plain(r.Summary, 1, 240) || !plain(r.Meaning, 1, 2000) || r.Phase != "working" && r.Phase != "needs_input" && r.Phase != "stopped" {
		return workflowError(2, "invalid delivery progress report")
	}
	if r.Phase == "needs_input" && (r.Question == nil || !plain(*r.Question, 1, 2000) || !strings.Contains(*r.Question, "?")) || r.Phase != "needs_input" && r.Question != nil {
		return workflowError(2, "needs_input requires an actual necessary question; other phases omit it")
	}
	if r.PreviousEventSHA256 != nil && !digestPattern.MatchString(*r.PreviousEventSHA256) {
		return workflowError(2, "invalid delivery predecessor digest")
	}
	seen := map[string]bool{}
	for _, v := range r.VerifierResults {
		if !key(v.ID) || seen[v.ID] || !plain(v.Reason, 1, 2000) {
			return workflowError(2, "invalid or duplicate verifier declaration")
		}
		seen[v.ID] = true
		switch v.Outcome {
		case "passed", "failed":
			if v.Exit == nil || *v.Exit < 0 || *v.Exit > 255 || len(v.Argv) == 0 || !filepath.IsAbs(v.CWD) || (v.Outcome == "passed") != (*v.Exit == 0) {
				return workflowError(2, "executed verifier requires a consistent actual exit and command")
			}
		case "not_run", "unknown":
			if v.Exit != nil {
				return workflowError(2, "an unrun or unknown verifier cannot claim an exit code")
			}
		default:
			return workflowError(2, "unknown verifier outcome")
		}
	}
	return nil
}

func WorkflowDeliveryReport(d Dependencies, root, id, contextPath, file string) (WorkflowRun, error) {
	return workflowDeliveryReport(d, root, id, contextPath, file, false)
}

// WorkflowDeliveryIncompleteReport allows the live owner to preserve a problem
// from the installed reader without switching controls or granting any effects.
func WorkflowDeliveryIncompleteReport(d Dependencies, root, id, contextPath, file string) (WorkflowRun, error) {
	return workflowDeliveryReport(d, root, id, contextPath, file, true)
}

func workflowDeliveryReport(d Dependencies, root, id, contextPath, file string, incomplete bool) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	raw, e := readFile(file, 1<<20, true)
	if e != nil {
		return WorkflowRun{}, e
	}
	var report DeliveryReport
	if e = decode(raw, 1<<20, &report); e != nil {
		return WorkflowRun{}, e
	}
	if e = validateDeliveryReport(report); e != nil {
		return WorkflowRun{}, e
	}
	if incomplete && report.Phase != "stopped" && report.Phase != "needs_input" {
		return WorkflowRun{}, workflowError(2, "--incomplete permits only stopped or needs_input reports; it cannot grant working or delivery authority")
	}
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if !deliveryRun(s.Request) || s.Result.Delivery == nil {
			return workflowError(4, "run is not a delivery owner")
		}
		if e := workflowClaim(*s, report.RunID, report.RequestSHA256, report.SessionID); e != nil {
			return e
		}
		if incomplete {
			if e := deliveryIncompleteReportCallback(d, *s, contextPath); e != nil {
				return e
			}
		}
		canonical, _ := Canonical(report)
		for _, old := range s.Result.Delivery.Events {
			if old.ID == report.EventID {
				if old.Kind != "report" || old.Binding.SHA256 != hash(canonical) {
					return workflowError(4, "delivery event already has different bytes")
				}
				return nil
			}
		}
		// An owner must be able to report unresolved runtime authority without
		// first claiming that authority. Incomplete reports grant no Task effects.
		if !incomplete {
			if e := deliveryCallback(d, *s, contextPath, report.Phase == "working"); e != nil {
				return e
			}
		}
		if !incomplete && s.Result.Delivery.Attempt != nil && s.Result.Delivery.Attempt.State == "attempted" {
			return workflowError(4, "a reserved effect is unresolved; inspect that attempt before another transition")
		}
		if !equal(report.PreviousEventSHA256, s.Result.Delivery.LastEventSHA256) {
			return workflowError(4, "delivery report predecessor differs from current state")
		}
		if _, e := deliveryAppendEvent(d, s, report.EventID, "report", report); e != nil {
			return e
		}
		if incomplete {
			if e := d.fault("delivery_after_incomplete_report_event"); e != nil {
				return e
			}
		}
		s.Result.Delivery.Phase = report.Phase
		s.Result.Round.State = report.Phase
		s.Result.NextAction = WorkflowAction{"recipient", "Continue the selected delivery within the bound authority."}
		if report.Phase == "needs_input" {
			s.Result.NextAction = WorkflowAction{"user", *report.Question}
		}
		if report.Phase == "stopped" {
			s.Result.NextAction = WorkflowAction{"user", report.Meaning}
		}
		return nil
	})
	return deliveryReadback(d, root, id, e)
}

type deliveryVerificationReceipt struct {
	Kind               string      `json:"kind"`
	SchemaVersion      int         `json:"schema_version"`
	RunID              string      `json:"run_id"`
	RequestSHA256      string      `json:"request_sha256"`
	AttemptID          string      `json:"attempt_id"`
	CandidateOID       string      `json:"candidate_oid"`
	CandidateTree      string      `json:"candidate_tree"`
	Argv               []string    `json:"argv"`
	CWD                string      `json:"cwd"`
	Acceptance         FileBinding `json:"acceptance"`
	AcceptanceSnapshot FileBinding `json:"acceptance_snapshot"`
	Review             FileBinding `json:"review"`
	Exit               *int        `json:"exit"`
	Stdout             FileBinding `json:"stdout"`
	Stderr             FileBinding `json:"stderr"`
	StartedAt          string      `json:"started_at"`
	FinishedAt         string      `json:"finished_at"`
	Error              string      `json:"error,omitempty"`
}

func deliveryReadBinding(path string, max int) (FileBinding, []byte, error) {
	raw, e := readFile(path, max, true)
	return FileBinding{path, hash(raw)}, raw, e
}

func WorkflowDeliveryVerify(d Dependencies, root, id, contextPath, reviewPath string) (WorkflowRun, error) {
	return workflowDeliveryVerify(d, root, id, contextPath, reviewPath, "")
}

func workflowDeliveryVerify(d Dependencies, root, id, contextPath, reviewPath, reuseAttempt string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	var s workflowState
	var receipt deliveryVerificationReceipt
	var verifierID string
	recoverReceipt := false
	alreadyQualified := false
	e := workflowUpdate(d, root, id, func(current *workflowState) error {
		if e := deliveryCallback(d, *current, contextPath, true); e != nil {
			return e
		}
		if reuseAttempt != "" && (current.Result.Delivery.Attempt == nil || current.Result.Delivery.Attempt.ID != reuseAttempt || current.Result.Delivery.Attempt.Kind != "verification") {
			return deliveryContinuityError("receipt_mismatch", "--reuse must name the latest verification attempt; inspect its exact receipt and use ordinary verify only when a new command execution is intended")
		}
		if attempt := current.Result.Delivery.Attempt; attempt != nil && (attempt.State == "attempted" || reuseAttempt != "") {
			if attempt.Kind != "verification" {
				return workflowError(4, "previous delivery effect is unresolved; it will not be replayed")
			}
			var e error
			receipt, e = deliveryRecoverVerification(*current, reviewPath)
			if e != nil {
				return e
			}
			if reuseAttempt != "" {
				if e = deliveryReusableVerification(d, *current, receipt); e != nil {
					return e
				}
				for _, candidate := range current.Result.Delivery.Candidates {
					if candidate.Key == reuseAttempt {
						if candidate.OID != receipt.CandidateOID || candidate.Tree != receipt.CandidateTree || candidate.Verification.SHA256 != digest(receipt) {
							return deliveryContinuityError("receipt_mismatch", "qualified candidate differs from the preserved verification attempt")
						}
						alreadyQualified = true
						return nil
					}
				}
				if current.Result.Delivery.Phase == "human_qa_passed" || current.Result.Delivery.Phase == "integrating" {
					return workflowError(4, "withdraw the current candidate explicitly before requalifying another attempt")
				}
				// Retain the original key and input digest while reserving its
				// bookkeeping continuation. No command is sent again.
				current.Result.Delivery.Attempt.State = "attempted"
				current.Result.Delivery.Phase = "verifying"
			}
			verifierID, e = deliveryVerifier(d, *current, receipt.CWD, receipt.Acceptance.Locator)
			if e != nil {
				return e
			}
			s, recoverReceipt = *current, true
			return nil
		}
		if current.Result.Delivery.Phase == "human_qa_passed" || current.Result.Delivery.Phase == "integrating" {
			return workflowError(4, "withdraw the current candidate explicitly before another verification")
		}
		x, e := deliveryTarget(d, *current)
		if e != nil {
			return e
		}
		acceptance, acceptanceRaw, e := deliveryReadBinding(current.Request.Delivery.AcceptancePath, 4<<20)
		if e != nil {
			return e
		}
		review, _, e := deliveryReadBinding(reviewPath, 1<<20)
		if e != nil {
			return e
		}
		attemptID := fmt.Sprintf("verify-%08d", len(current.Result.Delivery.Events)+1)
		path := filepath.Join(current.Result.Paths.RunRoot, "delivery", "attempts", attemptID)
		verifierID, e = deliveryVerifier(d, *current, x.Locator, acceptance.Locator)
		if e != nil {
			return e
		}
		snapshot := FileBinding{filepath.Join(path, "acceptance.sh"), acceptance.SHA256}
		if e = d.writeOnce(snapshot.Locator, acceptanceRaw); e != nil {
			return e
		}
		receipt = deliveryVerificationReceipt{Kind: "PlyDeliveryVerification@1", SchemaVersion: 1, RunID: id, RequestSHA256: current.Result.RequestSHA256, AttemptID: attemptID, CandidateOID: x.OID, CandidateTree: x.Tree, Argv: []string{"/bin/sh", acceptance.Locator}, CWD: x.Locator, Acceptance: acceptance, Review: review, StartedAt: d.Now().UTC().Format(time.RFC3339Nano)}
		receipt.AcceptanceSnapshot = snapshot
		if e = d.writeValue(filepath.Join(path, "attempt.json"), receipt); e != nil {
			return e
		}
		current.Result.Delivery.Attempt = &DeliveryAttempt{ID: attemptID, Kind: "verification", State: "attempted", Path: path, CandidateOID: x.OID, CandidateTree: x.Tree, InputSHA256: digest(receipt)}
		current.Result.Delivery.Phase = "verifying"
		s = *current
		return nil
	})
	if e != nil {
		o, _ := WorkflowShow(d, root, id)
		return o, e
	}
	if alreadyQualified {
		return deliveryReadback(d, root, id, nil)
	}
	if !recoverReceipt {
		if e = d.fault("delivery_after_verification_reservation"); e != nil {
			o, _ := WorkflowShow(d, root, id)
			return o, e
		}
		var stdout, stderr workflowLimitedOutput
		command := exec.Command(receipt.Argv[0], receipt.Argv[1:]...)
		command.Dir, command.Stdout, command.Stderr = receipt.CWD, &stdout, &stderr
		commandErr := command.Run()
		if command.ProcessState != nil && command.ProcessState.Exited() {
			receipt.Exit = ptr(command.ProcessState.ExitCode())
		}
		if commandErr != nil {
			receipt.Error = commandErr.Error()
		}
		if stdout.overflow || stderr.overflow {
			receipt.Error = "verifier output exceeded the preserved evidence limit"
		}
		receipt.FinishedAt = d.Now().UTC().Format(time.RFC3339Nano)
		path := s.Result.Delivery.Attempt.Path
		receipt.Stdout = FileBinding{filepath.Join(path, "stdout.txt"), hash(stdout.data)}
		receipt.Stderr = FileBinding{filepath.Join(path, "stderr.txt"), hash(stderr.data)}
		if e = d.writeOnce(receipt.Stdout.Locator, stdout.data); e == nil {
			e = d.writeOnce(receipt.Stderr.Locator, stderr.data)
		}
		if e == nil {
			e = d.writeValue(filepath.Join(path, "verification.json"), receipt)
		}
		if e != nil {
			o, _ := WorkflowShow(d, root, id)
			return o, e
		}
	}
	verification := FileBinding{filepath.Join(s.Result.Delivery.Attempt.Path, "verification.json"), digest(receipt)}
	if e = d.fault("delivery_after_verification_receipt"); e != nil {
		o, _ := WorkflowShow(d, root, id)
		return o, e
	}
	var candidate *workflowhandoff.DeliveryCandidateResult
	if receipt.Exit != nil && *receipt.Exit == 0 && receipt.Error == "" {
		if _, e = workflowBound(receipt.Acceptance, 4<<20); e == nil {
			_, e = workflowBound(receipt.Review, 1<<20)
		}
		if e == nil {
			x, err := deliveryTarget(d, s)
			e = err
			if e == nil && (x.OID != receipt.CandidateOID || x.Tree != receipt.CandidateTree) {
				e = workflowError(4, "candidate changed during verification")
			}
		}
		if e == nil {
			// Candidate publication follows the existing TaskRun -> workspace
			// lock order. Concurrent receipt observers may recover the same
			// native publication, but cannot race its staging/start/terminal.
			e = withStore(root, func() error {
				current, err := workflowRead(root, id)
				if err != nil {
					return err
				}
				if current.Result.Delivery.Attempt == nil || current.Result.Delivery.Attempt.ID != receipt.AttemptID || current.Result.Delivery.Attempt.InputSHA256 != s.Result.Delivery.Attempt.InputSHA256 {
					return workflowError(4, "verification reservation changed before qualification")
				}
				if err = deliveryCallback(d, current, contextPath, true); err != nil {
					return err
				}
				value, err := workflowhandoff.QualifyDeliveryCandidate(d.Workflow, workflowhandoff.DeliveryCandidateInput{ParentHandoffLocator: s.Result.Handoff.Locator, RunID: id, RequestSHA256: s.Result.RequestSHA256, CandidateKey: receipt.AttemptID, Summary: "Acceptance and the preserved review qualify this exact technical candidate.", CandidateOID: receipt.CandidateOID, CandidateTree: receipt.CandidateTree, VerifierID: verifierID, CWD: receipt.CWD, Argv: receipt.Argv, Exit: *receipt.Exit, StdoutPath: receipt.Stdout.Locator, StderrPath: receipt.Stderr.Locator, ReviewPath: reviewPath, VerificationPath: verification.Locator})
				if err == nil {
					candidate = &value
				}
				return err
			})
		}
	}
	qualificationErr := e
	if e = d.fault("delivery_after_candidate_qualification"); e != nil {
		o, _ := WorkflowShow(d, root, id)
		return o, e
	}
	e = workflowUpdate(d, root, id, func(current *workflowState) error {
		if current.Result.Delivery.Attempt == nil || current.Result.Delivery.Attempt.ID != receipt.AttemptID {
			return workflowError(4, "verification reservation changed")
		}
		for _, old := range current.Result.Delivery.Candidates {
			if old.Key == receipt.AttemptID {
				if candidate == nil || old.OID != receipt.CandidateOID || old.Tree != receipt.CandidateTree || old.Verification != verification || old.TaskResult.ID != candidate.TaskResult.ID {
					return workflowError(4, "candidate publication differs from the completed verification reservation")
				}
				return nil
			}
		}
		if _, err := deliveryAppendEvent(d, current, receipt.AttemptID, "verification", receipt); err != nil {
			return err
		}
		current.Result.Delivery.Attempt.State = "recorded"
		current.Result.Delivery.Phase = "working"
		current.Result.NextAction = WorkflowAction{"recipient", "Inspect the preserved verification and review evidence, fix within the selected goal, then verify again."}
		if qualificationErr != nil {
			current.Result.Reasons = append(current.Result.Reasons, Reason{"delivery_candidate_not_qualified", qualificationErr.Error()})
		}
		if candidate != nil {
			c := DeliveryCandidate{Number: len(current.Result.Delivery.Candidates) + 1, Key: receipt.AttemptID, OID: receipt.CandidateOID, Tree: receipt.CandidateTree, Verification: verification, Handoff: WorkflowHandoff{candidate.Handoff.HandoffID, candidate.Handoff.Locator, candidate.Handoff.SHA256}, TaskResult: candidate.TaskResult}
			current.Result.Delivery.Candidates = append(current.Result.Delivery.Candidates, c)
			current.Result.Delivery.Phase = "awaiting_human_qa"
			current.Result.TaskResultState = "candidate_qualified"
			current.Result.NextAction = WorkflowAction{"recipient", "Prepare the exact installed candidate journey and request the human's actual product judgment; retain ownership of the same session."}
			if a := current.Request.Delivery.Agreement; a != nil && a.HumanOwnedIntegration() && a.Mode != workspace.DeliveryPullRequest {
				current.Result.NextAction = WorkflowAction{"recipient", "Register the exact Delivery and explicitly release source ownership with ply workflow execute release. The human starts ply integration; stop before integration."}
			}
		}
		current.Result.Round.State = current.Result.Delivery.Phase
		return nil
	})
	if e == nil {
		e = qualificationErr
	}
	if e == nil && candidate == nil {
		e = workflowError(4, "verification did not qualify a candidate; actual output and exit are preserved")
	}
	return deliveryReadback(d, root, id, e)
}

func deliveryVerifier(d Dependencies, s workflowState, cwd, acceptance string) (string, error) {
	view, e := workflowhandoff.TaskRunHandoffView(d.Workflow, s.Result.Handoff.Locator)
	if e != nil {
		return "", e
	}
	var mandate struct {
		Verifiers []struct {
			ID           string   `json:"id"`
			Argv         []string `json:"argv"`
			CWD          string   `json:"cwd"`
			ExpectedExit int      `json:"expected_exit"`
		} `json:"verifiers"`
	}
	if e = json.Unmarshal(view, &mandate); e != nil {
		return "", e
	}
	if len(mandate.Verifiers) != 1 || !equal(mandate.Verifiers[0].Argv, []string{"/bin/sh", acceptance}) || mandate.Verifiers[0].CWD != cwd || mandate.Verifiers[0].ExpectedExit != 0 {
		return "", workflowError(4, "delivery verifier differs from the frozen acceptance entrypoint")
	}
	return mandate.Verifiers[0].ID, nil
}

// A completed command receipt can be reconciled after a crash; an attempt with
// no receipt is unknown and must never cause the command to run a second time.
func deliveryRecoverVerification(s workflowState, reviewPath string) (deliveryVerificationReceipt, error) {
	var receipt, initial deliveryVerificationReceipt
	a := s.Result.Delivery.Attempt
	raw, e := readFile(filepath.Join(a.Path, "attempt.json"), 256<<10, true)
	if e != nil || hash(raw) != a.InputSHA256 {
		return receipt, workflowError(4, "verification attempt evidence is missing or changed; no replay is safe")
	}
	if e = decode(raw, 256<<10, &initial); e != nil {
		return receipt, e
	}
	raw, e = readFile(filepath.Join(a.Path, "verification.json"), 256<<10, true)
	if e != nil {
		return receipt, deliveryContinuityError("effect_unknown", "verification outcome is unknown; inspect the reserved attempt and running process without replaying it")
	}
	if e = decode(raw, 256<<10, &receipt); e != nil {
		return receipt, e
	}
	copy := receipt
	copy.Exit, copy.Stdout, copy.Stderr, copy.FinishedAt, copy.Error = nil, FileBinding{}, FileBinding{}, "", ""
	started, startErr := time.Parse(time.RFC3339Nano, receipt.StartedAt)
	finished, finishErr := time.Parse(time.RFC3339Nano, receipt.FinishedAt)
	if !equal(initial, copy) || receipt.Kind != "PlyDeliveryVerification@1" || receipt.SchemaVersion != 1 || a.Path != filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", a.ID) || receipt.AttemptID != a.ID || receipt.CandidateOID != a.CandidateOID || receipt.CandidateTree != a.CandidateTree || receipt.RunID != s.Result.RunID || receipt.RequestSHA256 != s.Result.RequestSHA256 || receipt.Review.Locator != reviewPath || startErr != nil || finishErr != nil || finished.Before(started) {
		return receipt, deliveryContinuityError("receipt_mismatch", "completed verification receipt differs from its reserved attempt, candidate or original review; preserve it and reverify changed inputs explicitly")
	}
	for _, event := range s.Result.Delivery.Events {
		if event.ID == a.ID && (event.Kind != "verification" || event.Binding.SHA256 != hash(raw)) {
			return receipt, deliveryContinuityError("receipt_mismatch", "completed verification receipt differs from its immutable recorded event")
		}
	}
	for _, b := range []FileBinding{receipt.AcceptanceSnapshot, receipt.Stdout, receipt.Stderr} {
		if _, e = workflowBound(b, 4<<20); e != nil {
			return receipt, e
		}
	}
	return receipt, nil
}
