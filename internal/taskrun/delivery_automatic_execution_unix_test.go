//go:build !windows

package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestAutomaticVerificationInterruptionPreservesBlockedAndStopsChildren(t *testing.T) {
	f, o, binary, review := autoGateFixture(t, false, "")
	started := filepath.Join(f.R.WorkspaceRoot, "running-test")
	survived := filepath.Join(f.R.WorkspaceRoot, "surviving-child")
	script := "#!/bin/sh\nprintf 'actual output before interruption\\n'\n(sleep 2; printf alive > " + ShellQuote(survived) + ") &\nprintf ready > " + ShellQuote(started) + "\nwait\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	parent := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	type completion struct {
		run WorkflowRun
		err error
	}
	finished := make(chan completion, 1)
	go func() {
		run, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
		finished <- completion{run, err}
	}()
	deadline := time.Now().Add(15 * time.Second)
	for {
		if _, err := os.Stat(started); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("automatic test did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	var result completion
	select {
	case result = <-finished:
	case <-time.After(10 * time.Second):
		t.Fatal("interrupted automatic test did not finish")
	}
	if result.err == nil || result.run.DeliveryStatus.Verification.Outcome != "blocked" || result.run.DeliveryStatus.Verification.Exit != nil || len(result.run.Delivery.Candidates) != 1 || gitOutput(t, f.Parent, "rev-parse", "HEAD") != parent {
		t.Fatalf("interruption invented an exit, qualification or integration: err=%v status=%+v", result.err, result.run.DeliveryStatus.Verification)
	}
	raw, err := os.ReadFile(filepath.Join(result.run.Delivery.Attempt.Path, "stdout.txt"))
	if err != nil || !strings.Contains(string(raw), "actual output before interruption") {
		t.Fatalf("interrupted output was lost: %q %v", raw, err)
	}
	time.Sleep(2500 * time.Millisecond)
	if _, err := os.Stat(survived); !os.IsNotExist(err) {
		t.Fatalf("interrupted child survived its process group: %v", err)
	}
}

func TestAutomaticVerificationPreservesFinishedResultDuringTermination(t *testing.T) {
	f, o, binary, review := autoGateFixture(t, false, "")
	f.D.Fault = func(stage string) error {
		if stage == "delivery_before_verification_receipt" {
			if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
				return err
			}
			time.Sleep(10 * time.Millisecond)
		}
		return nil
	}
	result, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err != nil || result.DeliveryStatus.Verification.Outcome != "passed" {
		t.Fatalf("finished execution was lost during result publication: %v %+v", err, result.DeliveryStatus)
	}
	f.D.Fault = nil
	reused, err := WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, result.Delivery.Attempt.ID)
	if err != nil || reused.Delivery.Candidates[len(reused.Delivery.Candidates)-1].TaskResult.ID != result.Delivery.Candidates[len(result.Delivery.Candidates)-1].TaskResult.ID {
		t.Fatalf("completed result not recoverable once: %v", err)
	}
}
