package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Real public CLI verification/integration/closeout. The acceptance prints a
// native-looking JSON citation of removed test scratch to stdout and stderr.
func verifierLogCLIFixture(t *testing.T, binary string) (deliveryFixture, WorkflowRun, []byte) {
	t.Helper()
	root := ""
	if parent := os.Getenv("PLY_INFORMATIVE_JOURNEY_ROOT"); parent != "" {
		var err error
		root, err = os.MkdirTemp(parent, "verifier-log-")
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Preserved synthetic consumer workspace:", root)
	}
	agreement := &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic verifier log consumer", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: agreement, Root: root})
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	ownerPaneStandin(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	acceptance := writeAny(t, f.R.WorkspaceRoot, "cli-runtime-acceptance.json", deliveryTestAcceptance(t, f, o))
	o = automaticCLIOK(t, binary, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", acceptance)
	deliveryTestCandidate(t, f, "verifier log consumer")
	citation := filepath.Join(f.R.WorkspaceRoot, "removed-test-scratch", "historical.json")
	line, err := json.Marshal(map[string]any{"kind": "PlyDeliveryVerification@2", "diagnostic": FileBinding{citation, hash([]byte("historical scratch"))}})
	if err != nil {
		t.Fatal(err)
	}
	log := append(line, '\n')
	script := "#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'verifier log consumer'\nprintf '%s\\n' " + ShellQuote(string(line)) + "\nprintf '%s\\n' " + ShellQuote(string(line)) + " >&2\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "verifier-log")
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, review, binary)...)
	return f, o, log
}

func assertVerifierLogsRetained(t *testing.T, f deliveryFixture, o WorkflowRun, log []byte) {
	t.Helper()
	c := o.Delivery.Candidates[0]
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || closeout.State != "complete" || closeout.ResourceState != "retired" || !closeout.WorktreeRemoved || !closeout.BranchRemoved || closeout.Retention == nil {
		t.Fatalf("native closeout with retained evidence incomplete: %+v %v", closeout, err)
	}
	logs := 0
	for _, e := range closeout.Retention.Entries {
		if strings.HasPrefix(e.OriginLocator, filepath.Join(f.R.WorkspaceRoot, "removed-test-scratch")) {
			t.Fatal("historical log citation became a required dependency")
		}
		if e.SHA256 != hash(log) {
			continue
		}
		logs++
		got, err := workspace.ReadTaskRetainedEvidence(f.R.WorkspaceRoot, c.TaskResult.TaskID, e.OriginLocator, e.SHA256)
		if err != nil || !bytes.Equal(got, log) || e.SizeBytes != int64(len(log)) {
			t.Fatal("exact log binding did not resolve retained bytes", err)
		}
		if original, err := os.ReadFile(e.OriginLocator); err == nil && !bytes.Equal(original, log) {
			t.Fatal("frozen log rewritten")
		}
	}
	// stdout and stderr: the original verifier logs and their candidate copies.
	if logs < 2 {
		t.Fatalf("verifier logs missing from retained evidence: %d", logs)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("wrong integrated target")
	}
	t.Logf("Native closeout: run=%s result=%s candidate=%s manifest=%s; %d exact log entries retained", o.RunID, c.TaskResult.ID, c.OID, closeout.Retention.ManifestSHA256, logs)
}

func TestVerifierLogCLIRetainedCloseout(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o, log := verifierLogCLIFixture(t, binary)
	handoff, output, err := automaticCLI(binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if err == nil || !bytes.Contains(output, []byte("owner_handoff:")) || handoff.Delivery == nil || handoff.Delivery.Phase != "closing" {
		t.Fatalf("owner integrate did not end in the pending-cleanup handoff: %v\n%s", err, output)
	}
	cleanupOwnerLeaves(t, f)
	if o, output, err = cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30"); err != nil || o.Delivery.Phase != "completed" {
		t.Fatalf("native closeout did not finish: %v\n%s", err, output)
	}
	assertVerifierLogsRetained(t, f, o, log)
}

// Set PLY_INFORMATIVE_OLD_READER to the unmodified base CLI. It integrates;
// its coordinator follow then reaches the actual log-citation retention
// failure. The candidate's follow resumes the same Run's closeout without a
// repeated integration. No private state is edited to construct this boundary.
func TestVerifierLogCLIResumesOldReaderCloseout(t *testing.T) {
	old := os.Getenv("PLY_INFORMATIVE_OLD_READER")
	if old == "" {
		t.Skip("requires a preserved pre-fix CLI")
	}
	binary := automaticCLIBinary(t)
	f, o, log := verifierLogCLIFixture(t, old)
	handoff, output, err := automaticCLI(old, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context)
	if err == nil || !bytes.Contains(output, []byte("owner_handoff:")) || handoff.Delivery == nil || handoff.Delivery.Phase != "closing" {
		t.Fatalf("old reader integrate did not end in the pending-cleanup handoff: %v\n%s", err, output)
	}
	cleanupOwnerLeaves(t, f)
	failed, output, err := cleanupCLIAt(old, f.Parent, "execute", "follow", o.RunID, "--timeout", "30")
	if err == nil || !strings.Contains(string(output), "evidence_retention_failed") || !strings.Contains(string(output), "removed-test-scratch") || failed.Delivery == nil || failed.Delivery.Phase != "closing" || failed.Delivery.OwnershipRelease == nil || failed.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("old reader did not reproduce integrated/released log retention failure: %v %s", err, output)
	}
	c := failed.Delivery.Candidates[0]
	t.Logf("Old reader reached actual log retention failure after integration/release: run=%s result=%s candidate=%s", failed.RunID, c.TaskResult.ID, c.OID)
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("old reader had not actually integrated")
	}
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	bindings := []FileBinding{f.R.Delivery.Goal, {o.Paths.Context, hashFileTest(t, o.Paths.Context)}, {o.Handoff.Locator, o.Handoff.SHA256}, *failed.Delivery.OwnershipRelease, *c.Integration}
	for _, e := range failed.Delivery.Events {
		bindings = append(bindings, e.Binding)
	}
	done, output, err := cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30")
	if err != nil {
		t.Fatalf("candidate follow did not resume closeout: %v\n%s", err, output)
	}
	if done.RunID != o.RunID || done.Delivery.Phase != "completed" || len(done.Delivery.Candidates) != 1 || done.Delivery.Candidates[0].TaskResult.ID != c.TaskResult.ID || !equal(done.Delivery.OwnershipRelease, failed.Delivery.OwnershipRelease) {
		t.Fatal("resume changed result/owner or omitted completion")
	}
	assertVerifierLogsRetained(t, f, done, log)
	if actual := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic"); actual != reflog {
		t.Fatal("closeout resume repeated Git integration")
	}
	after, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(registry.IntegrationResults, after.IntegrationResults) || !equal(registry.IntegrationAttempts, after.IntegrationAttempts) || len(after.HumanQARecords) != 0 {
		t.Fatal("resume changed integration or manufactured human QA")
	}
	for _, b := range bindings {
		if hashFileTest(t, b.Locator) != b.SHA256 {
			t.Fatalf("resume rewrote immutable input %s", b.Locator)
		}
	}
	t.Log("Same-Run/result closeout resumed; prior Goal/context/handoff/events/release unchanged; integration attempts/results and target reflog unchanged")
}
