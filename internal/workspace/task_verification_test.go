package workspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDirectVerificationRecordsFailedCommandWithoutProcessClaims(t *testing.T) {
	f := newWorkItemJourneyFixtureVersion(t, true)
	source := filepath.Join(f.wrapper, "task")
	ready, err := CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: "task", Path: source, ExpectedParentOID: f.oid})
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	stdout, stderr, receipt := filepath.Join(dir, "stdout"), filepath.Join(dir, "stderr"), filepath.Join(dir, "receipt.json")
	for _, path := range []string{stdout, stderr} {
		if err := os.WriteFile(path, []byte("observed output\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	exit := 7
	observed := taskVerificationReceipt{Kind: "PlyDeliveryVerification@1", SchemaVersion: 1, RunID: "wfr_fixture", AttemptID: "explicit-check", CandidateOID: f.oid, CandidateTree: ready.Task.Worktree.Tree, Argv: []string{"/bin/sh", "-c", "exit 7"}, CWD: source, Exit: &exit, Stdout: taskVerificationBinding{Locator: stdout, SHA256: digestTaskBytes([]byte("observed output\n"))}, Stderr: taskVerificationBinding{Locator: stderr, SHA256: digestTaskBytes([]byte("observed output\n"))}, StartedAt: "2026-10-30T10:00:00Z", FinishedAt: "2026-10-30T10:00:01Z"}
	raw, err := json.Marshal(observed)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(receipt, raw, 0600); err != nil {
		t.Fatal(err)
	}
	in := TaskVerificationInput{TaskID: "task", ReceiptPath: receipt}
	result, err := RecordTaskVerification(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	if result.TechnicalGate != "" || result.ReportedOutcome != "observed" || result.HandoffID != "" || result.StartReceiptID != "" || result.TerminalResultID != "" || len(result.VerifierResults) != 1 || result.VerifierResults[0].Outcome != "failed" || result.VerifierResults[0].Exit != 7 {
		t.Fatalf("command observation acquired process claims: %+v", result)
	}
	registry, err := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.HumanQARecords) != 0 || len(registry.TaskResultSpecBindings) != 0 || len(registry.TaskResults) != 1 {
		t.Fatal("verification generated QA, Spec binding or duplicate result")
	}
	retry, err := RecordTaskVerification(f.dependencies, in)
	if err != nil || retry.ID != result.ID {
		t.Fatalf("same observed command was not idempotent: %v", err)
	}
	altered := result
	altered.VerifierResults = append([]TaskVerifierResultRecord{}, result.VerifierResults...)
	altered.VerifierResults[0].Outcome = "passed"
	if ValidateTaskVerificationReceipt(f.dependencies, altered) == nil {
		t.Fatal("failed command was relabelled passed")
	}
	if err := os.WriteFile(stdout, []byte("changed\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = RecordTaskVerification(f.dependencies, in); err == nil {
		t.Fatal("changed command output was accepted")
	}
	if err := os.WriteFile(stdout, []byte("observed output\n"), 0600); err != nil {
		t.Fatal(err)
	}
	observed.Exit = new(int)
	raw, err = json.Marshal(observed)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(receipt, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = RecordTaskVerification(f.dependencies, in); err == nil {
		t.Fatal("same operation replaced its exact receipt")
	}
}
