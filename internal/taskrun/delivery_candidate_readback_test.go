package taskrun

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestDeliveryCandidateReadbackKeepsActualCommandHistory(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "observed command result")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil || len(o.Delivery.Candidates) != 1 {
		t.Fatalf("command result missing: %v", err)
	}
	first := o.Delivery.Candidates[0]
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	before, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	after, readErr := os.ReadFile(statePath)
	if err != nil || readErr != nil || !bytes.Equal(before, after) || shown.DeliveryStatus.QualifiedCandidate == nil || !shown.DeliveryStatus.QualifiedCandidate.Current {
		t.Fatalf("readback changed history or lost actual evidence: %v %v", err, readErr)
	}
	counter := filepath.Join(f.R.WorkspaceRoot, "failed-checks.txt")
	if err = os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nprintf 'run\\n' >> "+ShellQuote(counter)+"\nexit 7\n"), 0600); err != nil {
		t.Fatal(err)
	}
	shown, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || shown.DeliveryStatus.QualifiedCandidate.Current || shown.DeliveryStatus.Verification.InputsMatch {
		t.Fatalf("changed input retained current result: %v", err)
	}
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil || len(o.Delivery.Candidates) != 2 || o.Delivery.Candidates[0].TaskResult.ID != first.TaskResult.ID {
		t.Fatalf("failed command lost history: %v", err)
	}
	result := o.Delivery.Candidates[1].TaskResult
	if len(result.VerifierResults) != 1 || result.VerifierResults[0].Exit != 7 || result.VerifierResults[0].Outcome != "failed" || result.TechnicalGate != "" || o.DeliveryStatus.Verification.Outcome != "failed" {
		t.Fatal("failed command was relabelled as success")
	}
	o, err = WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Delivery.Attempt.ID)
	count, readErr := os.ReadFile(counter)
	if err != nil || readErr != nil || string(count) != "run\n" || len(o.Delivery.Candidates) != 2 {
		t.Fatalf("failed command receipt was replayed or lost: %v %v", err, readErr)
	}
}
