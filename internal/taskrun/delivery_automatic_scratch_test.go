//go:build darwin

package taskrun

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func autoScratchFixture(t *testing.T, script string) (deliveryFixture, WorkflowRun, string, string) {
	t.Helper()
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement})
	ownerPaneStandin(t, &f)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "scratch fixture candidate")
	binary := filepath.Join(f.R.WorkspaceRoot, "candidate-executable")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nprintf 'fixture candidate\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	return f, o, binary, deliveryTestReview(t, f, "scratch")
}

func autoScratchBases(t *testing.T, bases ...string) {
	t.Helper()
	previous := automaticScratchBases
	automaticScratchBases = func() []string { return bases }
	t.Cleanup(func() { automaticScratchBases = previous })
}

func TestAutomaticScratchAllocationFailureBlocksBeforeTest(t *testing.T) {
	started := filepath.Join(t.TempDir(), "started")
	f, o, binary, review := autoScratchFixture(t, "#!/bin/sh\nprintf started > "+ShellQuote(started)+"\n")
	nested := filepath.Join(f.R.WorkspaceRoot, "nested-tmp")
	if err := os.Mkdir(nested, 0700); err != nil {
		t.Fatal(err)
	}
	autoScratchBases(t, nested, filepath.Join(t.TempDir(), "missing"))
	_, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err == nil || !strings.Contains(err.Error(), "temporary directory outside a Ply workspace") {
		t.Fatalf("allocation failure was not an understandable stop: %v", err)
	}
	if _, err := os.Stat(started); !os.IsNotExist(err) {
		t.Fatal("verification started with a wrong temporary directory")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || s.Result.Delivery.Attempt != nil {
		t.Fatalf("blocked allocation reserved a verification attempt: %v %+v", err, s.Result.Delivery.Attempt)
	}
	if entries, _ := os.ReadDir(nested); len(entries) != 0 {
		t.Fatalf("scratch was allocated inside the workspace: %v", entries)
	}
}

func TestAutomaticScratchFailedTestRetainsChildEvidence(t *testing.T) {
	script := "#!/bin/sh\nset -eu\nprintf data > \"$TMPDIR/written\"\nprintf 'scratch used\\n'\nprintf 'failing on purpose\\n' >&2\nexit 3\n"
	f, o, binary, review := autoScratchFixture(t, script)
	nested := filepath.Join(f.R.WorkspaceRoot, "nested-tmp")
	if err := os.Mkdir(nested, 0700); err != nil {
		t.Fatal(err)
	}
	outside, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	autoScratchBases(t, nested, outside)
	failed, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err == nil || failed.DeliveryStatus.Verification.Outcome != "failed" || len(failed.Delivery.Candidates) != 0 {
		t.Fatalf("real failure not preserved: %v %+v", err, failed.DeliveryStatus.Verification)
	}
	r := automaticTestReceipt(t, failed)
	scratch := automaticTestScratch(t, r)
	if filepath.Dir(scratch) != outside || !strings.Contains(filepath.Base(scratch), r.AttemptID) {
		t.Fatalf("scratch %s is not the attempt-bound directory under the first base outside the workspace", scratch)
	}
	if raw, err := os.ReadFile(filepath.Join(scratch, "written")); err != nil || string(raw) != "data" {
		t.Fatalf("failed attempt child evidence disappeared: %q %v", raw, err)
	}
	if info, err := os.Stat(scratch); err != nil || info.Mode().Perm() != 0700 {
		t.Fatalf("retained scratch is not private: %v %v", info, err)
	}
	for path, want := range map[string]string{r.Stdout.Locator: "scratch used", r.Stderr.Locator: "failing on purpose"} {
		if raw, err := os.ReadFile(path); err != nil || !strings.Contains(string(raw), want) {
			t.Fatalf("evidence %s lost with scratch: %q %v", path, raw, err)
		}
	}
	if r.Exit == nil || *r.Exit != 3 || r.Outcome != "fail" {
		t.Fatalf("failure receipt differs: %+v", r)
	}
}

// Headless verification keeps its sandbox: the test can write its own external
// attempt scratch, but not the shared temporary directory around it.
func TestWorkerAutomaticScratchSandboxIsBoundedToAttempt(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	script := "#!/bin/sh\nset -eu\nprintf data > \"$TMPDIR/written\"\nif printf escape > \"$TMPDIR/../ply-scratch-escape-$$\" 2>/dev/null; then exit 9; fi\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	workerTestReserveNative(t, f, o.RunID, "verify")
	result, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, o.Worker.Review.Locator, o.Worker.CandidateBinary)
	if err != nil || result.DeliveryStatus.Verification.Outcome != "passed" {
		var stderr []byte
		if result.Delivery.Attempt != nil {
			stderr, _ = os.ReadFile(filepath.Join(result.Delivery.Attempt.Path, "stderr.txt"))
		}
		t.Fatalf("sandboxed external scratch failed: %v %+v\n%s", err, result.DeliveryStatus.Verification, stderr)
	}
	r := automaticTestReceipt(t, result)
	scratch := automaticTestScratch(t, r)
	root, err := filepath.EvalSymlinks(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasPrefix(scratch+"/", root+"/") || r.Automatic.SandboxPolicy == nil {
		t.Fatalf("headless scratch %s is inside workspace or unsandboxed", scratch)
	}
	policy, err := os.ReadFile(r.Automatic.SandboxPolicy.Locator)
	if err != nil || !strings.Contains(string(policy), "(subpath \""+scratch+"\")") || strings.Contains(string(policy), "(subpath \""+filepath.Dir(scratch)+"\")") {
		t.Fatalf("sandbox policy is not bounded to the attempt scratch: %v\n%s", err, policy)
	}
	if _, err := os.Lstat(scratch); !os.IsNotExist(err) {
		t.Fatalf("headless scratch remains: %v", err)
	}
}

// TM01-F01: a preparation failure after scratch allocation hands no invocation
// onward, so it must not leave the external scratch behind or start the test.
func TestAutomaticScratchPreparationFailureRemovesScratch(t *testing.T) {
	for _, occupied := range []string{"candidate-executable", "acceptance.sh"} {
		t.Run(occupied, func(t *testing.T) {
			started := filepath.Join(t.TempDir(), "started")
			f, o, binary, review := autoScratchFixture(t, "#!/bin/sh\nprintf started > "+ShellQuote(started)+"\n")
			base, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			autoScratchBases(t, base)
			// The next attempt's snapshot already holds different bytes, so its
			// write-once preservation fails after scratch allocation.
			attempt := filepath.Join(o.Paths.RunRoot, "delivery", "attempts", "verify-00000001")
			if err := os.MkdirAll(attempt, 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(attempt, occupied), []byte("conflicting snapshot\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary); err == nil {
				t.Fatal("conflicting snapshot did not stop preparation")
			}
			if _, err := os.Stat(started); !os.IsNotExist(err) {
				t.Fatal("verification started after a failed preparation")
			}
			if entries, err := os.ReadDir(base); err != nil || len(entries) != 0 {
				t.Fatalf("failed preparation leaked scratch: %v %v", entries, err)
			}
		})
	}
}

func TestAutomaticScratchCompletionOnlyCleansOwnPass(t *testing.T) {
	for _, outcome := range []string{"pass", "fail", "blocked", ""} {
		t.Run(outcome, func(t *testing.T) {
			base := t.TempDir()
			own := filepath.Join(base, "own")
			sibling := filepath.Join(base, "other-active-attempt")
			for _, path := range []string{own, sibling} {
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(path, "child.log"), []byte("evidence"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			r := deliveryVerificationReceipt{RunID: "fixture-run", AttemptID: "fixture-attempt", Outcome: outcome, Automatic: &workflowhandoff.AutomaticVerificationInstructions{Environment: []string{"TMPDIR=" + own}}}
			finishAutomaticScratch(r)
			if _, err := os.Stat(own); (outcome == "pass" && !os.IsNotExist(err)) || (outcome != "pass" && err != nil) {
				t.Fatalf("wrong cleanup for outcome %q: %v", outcome, err)
			}
			if raw, err := os.ReadFile(filepath.Join(sibling, "child.log")); err != nil || string(raw) != "evidence" {
				t.Fatalf("another attempt was changed: %q %v", raw, err)
			}
		})
	}
}

// Policy construction only: this does not invoke the OS sandbox or establish
// live enforcement. The refused sandbox test remains unchanged and excluded.
func TestAutomaticScratchHeadlessPolicyConstruction(t *testing.T) {
	f, o, _, _ := autoScratchFixture(t, "#!/bin/sh\nexit 0\n")
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	autoScratchBases(t, base)
	path := filepath.Join(o.Paths.RunRoot, "delivery", "policy-only")
	scratch, err := automaticScratch("policy-only")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(scratch) })
	policyPath, err := workerVerifierPolicy(f.D, workflowState{}, path, scratch)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(policyPath)
	if err != nil {
		t.Fatal(err)
	}
	policy := string(raw)
	if !strings.Contains(policy, "(deny default)") || strings.Count(policy, "(allow file-write*") != 3 {
		t.Fatalf("unexpected write policy: %s", policy)
	}
	for _, permitted := range []string{filepath.Join(path, "home"), scratch} {
		if !strings.Contains(policy, "(allow file-write* (subpath "+strconv.Quote(permitted)+"))") {
			t.Fatalf("missing bounded write path %s: %s", permitted, policy)
		}
	}
	if strings.Contains(policy, "(subpath "+strconv.Quote(base)+")") || insidePlyWorkspace(scratch) {
		t.Fatalf("policy scratch escapes its physical attempt boundary: %s", policy)
	}
}
