package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestDeliveryRecordedVerificationRequalificationDoesNotRepeatAcceptance(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "preserved correction")
	counter := filepath.Join(f.R.WorkspaceRoot, "executions.txt")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nprintf 'executed\\n' >> "+ShellQuote(counter)+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "preserved-review")
	broken := f.D
	broken.Workflow.TaskWorkspace = nil
	failed, err := WorkflowDeliveryVerify(broken, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if err == nil || failed.Delivery.Attempt.State != "recorded" || len(failed.Delivery.Candidates) != 0 {
		t.Fatalf("expected recorded qualification failure: %+v %v", failed.Delivery, err)
	}
	attempt := failed.Delivery.Attempt.ID
	qualified, err := WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, attempt)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(counter)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "executed\n" || len(qualified.Delivery.Candidates) != 1 || qualified.Delivery.Candidates[0].Key != attempt {
		t.Fatalf("requalification repeated acceptance or replaced attempt: count=%q candidates=%+v", raw, qualified.Delivery.Candidates)
	}
	qualified, err = WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, attempt)
	if err != nil || len(qualified.Delivery.Candidates) != 1 {
		t.Fatalf("retry duplicated candidate: %+v %v", qualified.Delivery, err)
	}
	if raw, _ = os.ReadFile(counter); strings.Count(string(raw), "executed") != 1 {
		t.Fatalf("retry repeated acceptance: %q", raw)
	}
}

func TestDeliveryReceiptReuseRejectsChangedInputsWithoutRunningAcceptance(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "receipt mismatch")
	counter := filepath.Join(f.R.WorkspaceRoot, "executions.txt")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nprintf 'executed\\n' >> "+ShellQuote(counter)+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "original")
	broken := f.D
	broken.Workflow.TaskWorkspace = nil
	failed, err := WorkflowDeliveryVerify(broken, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if err == nil {
		t.Fatal("expected qualification failure")
	}
	attempt := failed.Delivery.Attempt.ID
	for _, path := range []string{f.AcceptancePath, review, filepath.Join(failed.Delivery.Attempt.Path, "acceptance.sh"), filepath.Join(failed.Delivery.Attempt.Path, "stdout.txt")} {
		original, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, append(append([]byte{}, original...), []byte("changed\n")...), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err = WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, attempt); err == nil {
			t.Fatalf("changed input accepted: %s", path)
		}
		if err = os.WriteFile(path, original, 0600); err != nil {
			t.Fatal(err)
		}
	}
	deliveryTestCandidate(t, f, "changed source")
	if _, err = WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, attempt); err == nil || !strings.Contains(err.Error(), "receipt_mismatch") {
		t.Fatalf("changed candidate reused receipt: %v", err)
	}
	if raw, err := os.ReadFile(counter); err != nil || string(raw) != "executed\n" {
		t.Fatalf("mismatch executed acceptance: %q %v", raw, err)
	}
}

func TestDeliveryRequalificationAfterNativePublicationIsConcurrentAndIdempotent(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "interrupted qualification")
	counter := filepath.Join(f.R.WorkspaceRoot, "executions.txt")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nprintf 'executed\\n' >> "+ShellQuote(counter)+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "exact")
	fault := f.D
	fault.Fault = func(point string) error {
		if point == "delivery_after_candidate_qualification" {
			return errors.New("lost reply after native candidate publication")
		}
		return nil
	}
	interrupted, err := WorkflowDeliveryVerify(fault, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
	if err == nil || interrupted.Delivery.Attempt.State != "attempted" {
		t.Fatalf("expected interrupted qualification: %+v %v", interrupted.Delivery, err)
	}
	attempt := interrupted.Delivery.Attempt.ID
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, attempt)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || len(shown.Delivery.Candidates) != 1 || len(shown.Delivery.Events) != 1 {
		t.Fatalf("duplicate candidate or event: %+v %v", shown.Delivery, err)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(registry.TaskResults) != 1 {
		t.Fatalf("duplicate native result: %d %v", len(registry.TaskResults), err)
	}
	if raw, err := os.ReadFile(counter); err != nil || string(raw) != "executed\n" {
		t.Fatalf("recovery repeated acceptance: %q %v", raw, err)
	}
}
