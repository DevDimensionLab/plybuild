package taskrun

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

// WorkflowDeliveryVerifyWithBinary runs the existing native candidate verifier,
// additionally binding a private executable for explicitly automatic agreements.
func WorkflowDeliveryVerifyWithBinary(d Dependencies, root, id, contextPath, reviewPath, binary string) (WorkflowRun, error) {
	return workflowDeliveryVerify(d, root, id, contextPath, reviewPath, "", binary)
}

func prepareAutomaticVerification(d Dependencies, s workflowState, receipt *deliveryVerificationReceipt, path, binary string) error {
	a := deliveryEffectiveAgreement(s)
	if a == nil || !a.AutomaticAcceptance() {
		if binary != "" {
			return workflowError(2, "--candidate-binary requires an explicitly automatic acceptance agreement")
		}
		return nil
	}
	if !filepath.IsAbs(binary) {
		return workflowError(4, "automatic acceptance is blocked: --candidate-binary must name the built private candidate executable")
	}
	if err := physical(binary, true); err != nil {
		return err
	}
	info, err := os.Stat(binary)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		return workflowError(4, "automatic acceptance is blocked: candidate binary is missing or not executable")
	}
	if info.Size() > workflowhandoff.AutomaticCandidateBinaryLimit {
		return workflowError(4, "automatic acceptance is blocked: candidate binary exceeds the 64 MiB native evidence limit")
	}
	binaryBytes, err := readFile(binary, workflowhandoff.AutomaticCandidateBinaryLimit, false)
	if err != nil {
		return err
	}
	binding := FileBinding{Locator: binary, SHA256: hash(binaryBytes)}
	for _, name := range []string{"home", "tmp"} {
		if err := privateDir(filepath.Join(path, name)); err != nil {
			return err
		}
	}
	binarySnapshot := filepath.Join(path, "candidate-executable")
	if err := d.writeOnce(binarySnapshot, binaryBytes); err != nil {
		return err
	}
	if err := os.Chmod(binarySnapshot, 0700); err != nil {
		return err
	}
	v := &workflowhandoff.AutomaticVerificationInstructions{
		Policy:                  *a.Acceptance,
		CandidateBinary:         workflowhandoff.AutomaticBinaryBinding{Locator: binding.Locator, SHA256: binding.SHA256},
		CandidateBinarySnapshot: workflowhandoff.AutomaticBinaryBinding{Locator: binarySnapshot, SHA256: binding.SHA256},
		ExecutedArgv:            []string{"/bin/sh", receipt.AcceptanceSnapshot.Locator},
		Environment:             []string{"HOME=" + filepath.Join(path, "home"), "TMPDIR=" + filepath.Join(path, "tmp"), "PATH=" + os.Getenv("PATH"), "LC_ALL=C", "PLY_CANDIDATE_BINARY=" + binarySnapshot, "PLY_CANDIDATE_OID=" + receipt.CandidateOID, "PLY_CANDIDATE_TREE=" + receipt.CandidateTree},
		ExecutorClaim:           s.Request.Runtime.Provider + " / " + s.Request.Delivery.OwnerClaim,
		ExecutorSessionID:       s.Result.Transport.AgentSessionID, ExpectedExit: 0,
	}
	sort.Strings(v.Environment)
	v.InstructionsSHA256 = workflowhandoff.AutomaticInstructionsDigest(*v, receipt.Argv, receipt.CWD, receipt.Acceptance.SHA256, receipt.CandidateOID, receipt.CandidateTree)
	if err := workflowhandoff.ValidateAutomaticInstructions(*v, a, receipt.Argv, receipt.CWD, receipt.Acceptance.SHA256, receipt.CandidateOID, receipt.CandidateTree); err != nil {
		return err
	}
	receipt.Kind, receipt.SchemaVersion, receipt.Automatic = "PlyDeliveryVerification@2", 2, v
	return nil
}

func deliveryVerificationVersion(r deliveryVerificationReceipt) bool {
	return r.Kind == "PlyDeliveryVerification@1" && r.SchemaVersion == 1 && r.Automatic == nil && r.Outcome == "" || r.Kind == "PlyDeliveryVerification@2" && r.SchemaVersion == 2 && r.Automatic != nil
}

func validateAutomaticReceipt(s workflowState, r deliveryVerificationReceipt, live bool) error {
	a := deliveryEffectiveAgreement(s)
	if a == nil || !a.AutomaticAcceptance() {
		if r.Automatic != nil {
			return workflowError(4, "historical agreement cannot acquire automatic acceptance")
		}
		return nil
	}
	if !deliveryVerificationVersion(r) || r.Automatic == nil || r.Outcome != "pass" || r.Exit == nil || *r.Exit != 0 || r.Error != "" {
		return workflowError(4, "automatic acceptance has no successful preserved execution")
	}
	if err := workflowhandoff.ValidateAutomaticInstructions(*r.Automatic, a, r.Argv, r.CWD, r.Acceptance.SHA256, r.CandidateOID, r.CandidateTree); err != nil {
		return err
	}
	if r.Automatic.ExecutorSessionID != s.Result.Transport.AgentSessionID || r.Acceptance.Locator != s.Request.Delivery.AcceptancePath || r.CWD != s.Observed.Target.WorktreeLocator || r.Acceptance.SHA256 != r.AcceptanceSnapshot.SHA256 {
		return workflowError(4, "automatic acceptance differs from its actual owner or frozen entrypoint")
	}
	if r.Automatic.ExecutedArgv[1] != r.AcceptanceSnapshot.Locator || r.Automatic.CandidateBinarySnapshot.Locator != filepath.Join(filepath.Dir(r.AcceptanceSnapshot.Locator), "candidate-executable") {
		return workflowError(4, "automatic execution snapshots differ from the reserved attempt")
	}
	if _, err := workflowBound(FileBinding{Locator: r.Automatic.CandidateBinarySnapshot.Locator, SHA256: r.Automatic.CandidateBinarySnapshot.SHA256}, workflowhandoff.AutomaticCandidateBinaryLimit); err != nil {
		return err
	}
	if live {
		for _, b := range []FileBinding{r.Acceptance, r.Review} {
			if _, err := workflowBound(b, 256<<20); err != nil {
				return fmt.Errorf("automatic acceptance input changed: %w", err)
			}
		}
		raw, err := readFile(r.Automatic.CandidateBinary.Locator, workflowhandoff.AutomaticCandidateBinaryLimit, false)
		if err != nil || hash(raw) != r.Automatic.CandidateBinary.SHA256 {
			return workflowError(4, "automatic acceptance candidate binary changed or is unavailable")
		}
	}
	return nil
}

// Called both at delivery preview and the native workspace's final authority
// read before effects, including qualification before Candidates is appended.
func validateAutomaticVerification(d Dependencies, s workflowState, c DeliveryCandidate) error {
	if deliveryEffectiveAgreement(s) == nil || !deliveryEffectiveAgreement(s).AutomaticAcceptance() {
		return nil
	}
	raw, err := workflowBound(c.Verification, 256<<10)
	if err != nil {
		return err
	}
	var r deliveryVerificationReceipt
	if err = decode(raw, 256<<10, &r); err != nil {
		return err
	}
	if r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.AttemptID != c.Key || r.CandidateOID != c.OID || r.CandidateTree != c.Tree {
		return workflowError(4, "automatic acceptance is bound to another candidate or run")
	}
	live := true
	// Once this exact reserved integration is observed at its target, recovery
	// controls the recorded effect. Disappearing working test inputs cannot
	// turn that effect into an unperformed merge or force a second execution.
	if attempt := s.Result.Delivery.Attempt; attempt != nil && attempt.Kind == "integration" && attempt.CandidateOID == c.OID && attempt.CandidateTree == c.Tree {
		a := deliveryEffectiveAgreement(s)
		target, err := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(a.TargetWorktree, a.TargetRef)
		if err == nil && target.Clean && len(target.InProgress) == 0 && target.Ref == a.TargetRef && target.GitCommonDir == s.Observed.Target.GitCommonDir && target.OID == c.OID && target.Tree == c.Tree {
			live = false
		}
	}
	return validateAutomaticReceipt(s, r, live)
}

func runDeliveryVerification(d Dependencies, r *deliveryVerificationReceipt) (workflowLimitedOutput, workflowLimitedOutput, func()) {
	var stdout, stderr workflowLimitedOutput
	ctx := context.Background()
	cancel := func() {}
	if r.Automatic != nil {
		interrupted, stop := signal.NotifyContext(ctx, automaticVerificationSignals()...)
		timed, stopTimeout := context.WithTimeout(interrupted, time.Duration(r.Automatic.Policy.TimeoutSeconds)*time.Second)
		ctx, cancel = timed, func() { stopTimeout(); stop() }
	}
	argv := r.Argv
	if r.Automatic != nil {
		argv = r.Automatic.ExecutedArgv
	}
	command := exec.CommandContext(ctx, argv[0], argv[1:]...)
	command.Dir, command.Stdout, command.Stderr = r.CWD, &stdout, &stderr
	if r.Automatic != nil {
		command.Env = r.Automatic.Environment
		configureAutomaticProcess(command)
		command.WaitDelay = time.Second
	}
	err := command.Run()
	if command.ProcessState != nil && command.ProcessState.Exited() && command.ProcessState.ExitCode() >= 0 {
		r.Exit = ptr(command.ProcessState.ExitCode())
	}
	if err != nil {
		r.Error = err.Error()
	}
	if ctx.Err() != nil {
		r.Error = "automatic acceptance blocked: " + ctx.Err().Error()
	}
	if stdout.overflow || stderr.overflow {
		r.Error = "verifier output exceeded the preserved evidence limit"
	}
	r.FinishedAt = d.Now().UTC().Format(time.RFC3339Nano)
	if r.Automatic != nil {
		r.Outcome = "blocked"
		if r.Exit != nil && ctx.Err() == nil && !stdout.overflow && !stderr.overflow {
			r.Outcome = "fail"
			if *r.Exit == 0 && r.Error == "" {
				r.Outcome = "pass"
			}
		}
	}
	// Keep signal handling installed until the caller durably preserves the
	// completed result; a termination request cannot discard a finished test
	// merely because its registration has not completed yet.
	return stdout, stderr, cancel
}
