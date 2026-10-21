package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Bootstrap is synthetic, but every acceptance, readback and delivery action
// below crosses the public built CLI boundary in a disposable native workspace.
func automaticCLIBinary(t *testing.T) string {
	t.Helper()
	path, err := filepath.EvalSymlinks(deliveryCLIBinary(t))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func automaticCLIFixture(t *testing.T, binary string, automatic bool, timeout int) (deliveryFixture, WorkflowRun) {
	t.Helper()
	a := &workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic}
	if automatic {
		a.SchemaVersion = 3
		a.Acceptance = &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic CLI orchestrator", TimeoutSeconds: timeout}
	}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: a})
	// Public native closeout reads the actual registered goal manifest for both
	// historical and automatic Runs, including an explicit later policy choice.
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	ownerPaneStandin(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	acceptance := writeAny(t, f.R.WorkspaceRoot, "cli-runtime-acceptance.json", deliveryTestAcceptance(t, f, o))
	o = automaticCLIOK(t, binary, f, "run", "accept", o.RunID, "--context", o.Paths.Context, "--file", acceptance)
	return f, o
}

func automaticCLI(binary string, f deliveryFixture, args ...string) (WorkflowRun, []byte, error) {
	c := exec.Command(binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir = f.Prepared.Preparation.Preparation.Plan.WorktreePath
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	var o WorkflowRun
	if decodeErr := json.Unmarshal(stdout.Bytes(), &o); decodeErr != nil && err == nil {
		err = decodeErr
	}
	return o, append(stdout.Bytes(), stderr.Bytes()...), err
}

func automaticCLIOK(t *testing.T, binary string, f deliveryFixture, args ...string) WorkflowRun {
	t.Helper()
	o, output, err := automaticCLI(binary, f, args...)
	if err != nil {
		t.Fatalf("CLI %v: %v\n%s", args, err, output)
	}
	return o
}

func automaticCLIVerifyArgs(o WorkflowRun, review, binary string) []string {
	return []string{"execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review, "--candidate-binary", binary}
}

func assertAutomaticCLIClosure(t *testing.T, binary string, f deliveryFixture, o WorkflowRun) {
	t.Helper()
	c := o.Delivery.Candidates[len(o.Delivery.Candidates)-1]
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	registered := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "automatic"))
	id := deliveryCLIString(t, registered, "id")
	ready := deliveryCLIOK(t, binary, cwd, "check", id)
	if deliveryCLIString(t, ready, "state") != "ready" {
		t.Fatalf("automatic gate not ready: %s", ready["reasons"])
	}
	cleanup, err := workflowhandoff.DeliveryCleanupAllowed(f.R.HandoffDraft)
	if err != nil {
		t.Fatal(err)
	}
	if cleanup {
		// CL02 default: the owner's own Delivery execute integrates and leaves
		// cleanup pending; the coordinator's follow removes the released source
		// after the owner leaves; Delivery then records the native completion.
		pending, output, err := deliveryCLIRun(binary, cwd, "execute", id)
		if err == nil || deliveryCLIString(t, pending, "state") != "pending_closeout" || string(pending["native_closed"]) != "false" || string(pending["local_integration"]) == "null" {
			t.Fatalf("owner Delivery execute did not leave observed integration with pending cleanup: %v\n%s", err, output)
		}
		if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
			t.Fatal("automatic Delivery did not integrate the exact candidate")
		}
		if _, err = os.Stat(cwd); err != nil {
			t.Fatal("active owner's worktree was removed")
		}
		workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "idle"})
		if _, output, err = cleanupCLIAt(binary, f.Parent, "execute", "follow", o.RunID, "--timeout", "30"); err != nil {
			t.Fatalf("coordinator follow did not finish cleanup: %v\n%s", err, output)
		}
		if _, err = os.Lstat(cwd); !os.IsNotExist(err) {
			t.Fatalf("released Task worktree remains: %v", err)
		}
		cwd = f.Parent
	}
	done := deliveryCLIOK(t, binary, cwd, "execute", id)
	if deliveryCLIString(t, done, "state") != "delivered" || string(done["native_closed"]) != "true" {
		t.Fatalf("delivery incomplete: %+v", done)
	}
	before := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	deliveryCLIOK(t, binary, cwd, "execute", id)
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != before {
		t.Fatal("wrong or repeated Git effect")
	}
	r, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(r.HumanQARecords) != 0 {
		t.Fatalf("machine result fabricated human QA: %v %+v", err, r.HumanQARecords)
	}
	base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || deliveryCLILatestBase(base).OID != c.OID {
		t.Fatalf("Epic base not advanced: %v %+v", err, base)
	}
	status, output, err := cleanupCLIAt(binary, f.Parent, "execute", "show", o.RunID)
	if err != nil || status.Delivery.Phase != "completed" || status.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("native closure invisible: %v %+v\n%s", err, status.DeliveryStatus, output)
	}
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || closeout.State != "complete" || closeout.Plan.Keep == cleanup {
		t.Fatalf("native closeout differs from the frozen mandate's cleanup grant %t: %+v %v", cleanup, closeout, err)
	}
}

func TestAutomaticAcceptanceCLIGreenAndFailureCorrection(t *testing.T) {
	binary := automaticCLIBinary(t)
	for _, failFirst := range []bool{false, true} {
		t.Run(fmt.Sprint(failFirst), func(t *testing.T) {
			f, o := automaticCLIFixture(t, binary, true, 30)
			value := "expected"
			if failFirst {
				value = "broken"
			}
			deliveryTestCandidate(t, f, value)
			script := "#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = expected\n\"$PLY_CANDIDATE_BINARY\" capabilities --format json\n"
			if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
				t.Fatal(err)
			}
			review := deliveryTestReview(t, f, "first")
			if failFirst {
				failed, output, err := automaticCLI(binary, f, automaticCLIVerifyArgs(o, review, binary)...)
				if err == nil || failed.DeliveryStatus.Verification.Outcome != "failed" || len(failed.Delivery.Candidates) != 0 {
					t.Fatalf("real failure not preserved: %v %s", err, output)
				}
				if _, _, err := automaticCLI(binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context); err == nil {
					t.Fatal("failed test integrated")
				}
				deliveryTestCandidate(t, f, "expected")
				if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
					t.Fatal(err)
				}
				review = deliveryTestReview(t, f, "corrected")
			}
			o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, review, binary)...)
			if o.Delivery.Phase != "automatic_acceptance_passed" || o.DeliveryStatus.Acceptance == nil || o.DeliveryStatus.Acceptance.Outcome != "pass" {
				t.Fatalf("automatic acceptance missing: %+v", o)
			}
			if failFirst && len(o.Delivery.Events) != 2 {
				t.Fatalf("failure history lost: %+v", o.Delivery.Events)
			}
			assertAutomaticCLIClosure(t, binary, f, o)
		})
	}
}

func TestAutomaticAcceptanceCLIBlocksTimeoutAndCompetingOwner(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 2)
	deliveryTestCandidate(t, f, "timeout")
	marker := filepath.Join(f.R.WorkspaceRoot, "test-started")
	script := "#!/bin/sh\nprintf 'once\\n' >> " + ShellQuote(marker) + "\nsleep 30\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "timeout")
	args := automaticCLIVerifyArgs(o, review, binary)
	type answer struct {
		run    WorkflowRun
		output []byte
		err    error
	}
	completed := make(chan answer, 1)
	go func() { run, output, err := automaticCLI(binary, f, args...); completed <- answer{run, output, err} }()
	deadline := time.Now().Add(15 * time.Second)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("verifier did not start")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if _, output, err := automaticCLI(binary, f, args...); err == nil || !bytes.Contains(output, []byte("unknown")) {
		t.Fatalf("competing owner not blocked: %v %s", err, output)
	}
	finished := <-completed
	if finished.err == nil || finished.run.DeliveryStatus.Verification.Outcome != "blocked" {
		t.Fatalf("timeout not blocked: %v %s", finished.err, finished.output)
	}
	raw, _ := os.ReadFile(marker)
	if string(raw) != "once\n" {
		t.Fatalf("duplicate verifier execution: %q", raw)
	}
	if _, _, err := automaticCLI(binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context); err == nil {
		t.Fatal("timeout integrated")
	}
}

func TestAutomaticAcceptanceCLIHistoricalAgreementRetainsHumanGate(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, false, 30)
	deliveryTestCandidate(t, f, "historical")
	review := deliveryTestReview(t, f, "historical")
	o = automaticCLIOK(t, binary, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review)
	if o.Delivery.Phase != "awaiting_human_qa" {
		t.Fatal("historical agreement gained automatic acceptance")
	}
	if _, output, err := automaticCLI(binary, f, "execute", "integrate", o.RunID, "--context", o.Paths.Context); err == nil || !strings.Contains(string(output), "human") {
		t.Fatalf("human gate bypassed: %v %s", err, output)
	}
}

const automaticReceiptCrashGit = `import json, os, signal, sys
from pathlib import Path
p = Path(__file__).resolve().parent.parent / "receipt-crash.json"
m = json.loads(p.read_text())
if list(Path(m["run_root"]).glob("delivery/attempts/verify-*/verification.json")) and not m.get("crashed"):
    m["crashed"] = True
    p.write_text(json.dumps(m))
    os.kill(os.getppid(), signal.SIGKILL)
    raise SystemExit(9)
if any(a in sys.argv[1:] for a in ["push", "fetch", "pull", "clone", "ls-remote"]):
    raise SystemExit("unexpected network command")
os.execv(m["git"], [m["git"], *sys.argv[1:]])
`

func automaticCLIAdapter(t *testing.T, f deliveryFixture, name, source string) {
	t.Helper()
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(f.R.WorkspaceRoot, name)
	if err = os.Mkdir(bin, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(bin, "git"), []byte("#!"+python+"\n"+source), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestAutomaticAcceptanceCLIReceiptRecoveryDoesNotRepeatTest(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	deliveryTestCandidate(t, f, "receipt recovery")
	counter := filepath.Join(f.R.WorkspaceRoot, "test-count")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\nprintf 'once\\n' >> "+ShellQuote(counter)+"\ntest \"$(cat README.md)\" = 'receipt recovery'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "receipt-recovery")
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	model := writeAny(t, f.R.WorkspaceRoot, "receipt-crash.json", map[string]any{"git": git, "run_root": o.Paths.RunRoot})
	automaticCLIAdapter(t, f, "receipt-adapter", automaticReceiptCrashGit)
	args := automaticCLIVerifyArgs(o, review, binary)
	if _, output, err := automaticCLI(binary, f, args...); err == nil {
		t.Fatalf("crashed verifier claimed success: %s", output)
	}
	if workflowProviderDocument(t, model)["crashed"] != true {
		t.Fatal("did not interrupt after completed receipt")
	}
	o = automaticCLIOK(t, binary, f, args...)
	if len(o.Delivery.Candidates) != 1 {
		t.Fatal("recovery failed to preserve one candidate")
	}
	// Explicit receipt reuse also returns the exact original candidate.
	reused := automaticCLIOK(t, binary, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review, "--reuse", o.Delivery.Candidates[0].Key)
	if reused.Delivery.Candidates[0].TaskResult.ID != o.Delivery.Candidates[0].TaskResult.ID {
		t.Fatal("receipt retry published another candidate")
	}
	raw, _ := os.ReadFile(counter)
	if string(raw) != "once\n" {
		t.Fatalf("recovery repeated acceptance: %q", raw)
	}
	assertAutomaticCLIClosure(t, binary, f, o)
}

func TestAutomaticAcceptanceCLIAfterMergeRecovery(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	deliveryTestCandidate(t, f, "merge recovery")
	privateBinary := filepath.Join(f.R.WorkspaceRoot, "private-candidate")
	binaryBytes, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(privateBinary, binaryBytes, 0700); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "merge")
	o = automaticCLIOK(t, binary, f, automaticCLIVerifyArgs(o, review, privateBinary)...)
	c := o.Delivery.Candidates[0]
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	model := writeAny(t, f.R.WorkspaceRoot, "local-crash-fixture.json", map[string]any{"git": git, "parent": f.Parent, "oid": c.OID, "crash": true})
	automaticCLIAdapter(t, f, "merge-adapter", deliveryCLILocalCrashGit)
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	r := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "merge-crash"))
	id := deliveryCLIString(t, r, "id")
	if _, _, err := deliveryCLIRun(binary, cwd, "execute", id); err == nil {
		t.Fatal("crashed merge claimed completion")
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("crash did not follow real merge")
	}
	base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || deliveryCLILatestBase(base).OID == c.OID {
		t.Fatalf("crash did not precede base receipt: %v", err)
	}
	// Registration recovery must use preserved passing evidence once Git has
	// already integrated, even if the original test inputs are gone. Their
	// exact snapshots and managed native candidate artifacts remain available.
	for _, path := range []string{privateBinary, f.AcceptancePath, review} {
		if err = os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	// The owner's own recovery observes the integration; CL02 cleanup then
	// waits for the owner to leave and finishes from the return checkout.
	pending, output, err := deliveryCLIRun(binary, cwd, "execute", id)
	if err == nil || deliveryCLIString(t, pending, "state") != "pending_closeout" || string(pending["local_integration"]) == "null" {
		t.Fatalf("merge recovery did not observe integration with pending cleanup: %v\n%s", err, output)
	}
	cleanupOwnerLeaves(t, f)
	done := deliveryCLIOK(t, binary, f.Parent, "execute", id)
	if deliveryCLIString(t, done, "state") != "delivered" {
		t.Fatalf("merge recovery incomplete: %+v", done)
	}
	if _, err = os.Lstat(cwd); !os.IsNotExist(err) {
		t.Fatalf("released Task worktree remains: %v", err)
	}
	deliveryCLIOK(t, binary, f.Parent, "execute", id)
	if workflowProviderDocument(t, model)["merge_count"] != float64(1) {
		t.Fatal("merge repeated")
	}
	run, output, err := cleanupCLIAt(binary, f.Parent, "execute", "show", o.RunID)
	if err != nil || run.Delivery.Phase != "completed" || run.DeliveryStatus.HumanJudgment != nil {
		t.Fatalf("wrong native closure or fabricated human QA: %v\n%s", err, output)
	}
}

// TM01: ordinary automatic acceptance creates a disposable workspace fixture
// in its default temp directory, as Go os.TempDir does. The script sets no
// TMPDIR/GOTMPDIR itself.
func TestAutomaticAcceptanceCLIDefaultTempHostsWorkspaceFixture(t *testing.T) {
	binary := automaticCLIBinary(t)
	for _, nestedAmbient := range []bool{false, true} {
		t.Run(fmt.Sprint(nestedAmbient), func(t *testing.T) {
			f, o := automaticCLIFixture(t, binary, true, 60)
			if nestedAmbient {
				// The verifier's own inherited temp choice lies inside the workspace.
				ambient := filepath.Join(f.R.WorkspaceRoot, "ambient-tmp")
				if err := os.Mkdir(ambient, 0700); err != nil {
					t.Fatal(err)
				}
				t.Setenv("TMPDIR", ambient)
			}
			deliveryTestCandidate(t, f, "temp fixture")
			script := "#!/bin/sh\nset -eu\nfixture=$(mktemp -d \"$TMPDIR/fixture.XXXXXX\")\ncd \"$fixture\"\n\"$PLY_CANDIDATE_BINARY\" workspace init\nprintf 'fixture created\\n'\n"
			if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
				t.Fatal(err)
			}
			review := deliveryTestReview(t, f, "temp-fixture")
			o, output, err := automaticCLI(binary, f, automaticCLIVerifyArgs(o, review, binary)...)
			if err != nil || o.DeliveryStatus.Acceptance == nil || o.DeliveryStatus.Acceptance.Outcome != "pass" {
				var stderr []byte
				if o.Delivery.Attempt != nil {
					stderr, _ = os.ReadFile(filepath.Join(o.Delivery.Attempt.Path, "stderr.txt"))
				}
				t.Fatalf("default temp could not host a workspace fixture: %v\nverifier stderr: %s\n%s", err, stderr, output)
			}
			r := automaticTestReceipt(t, o)
			scratch := automaticTestScratch(t, r)
			workspaceRoot, err := filepath.EvalSymlinks(f.R.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			if resolved, err := filepath.EvalSymlinks(filepath.Dir(scratch)); err != nil || resolved != filepath.Dir(scratch) || strings.HasPrefix(scratch+"/", workspaceRoot+"/") {
				t.Fatalf("scratch %s is not physical or lies inside workspace %s: %v", scratch, workspaceRoot, err)
			}
			if _, err := os.Lstat(scratch); !os.IsNotExist(err) {
				t.Fatalf("attempt scratch was not removed after its preserved result: %v", err)
			}
			if raw, err := os.ReadFile(r.Stdout.Locator); err != nil || !strings.Contains(string(raw), "fixture created") {
				t.Fatalf("stdout evidence lost with scratch: %q %v", raw, err)
			}
			// Completed receipt recovery never needs the removed scratch.
			reused := automaticCLIOK(t, binary, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review, "--reuse", o.Delivery.Candidates[0].Key)
			if reused.Delivery.Candidates[0].TaskResult.ID != o.Delivery.Candidates[0].TaskResult.ID {
				t.Fatal("receipt recovery published another candidate")
			}
			assertAutomaticCLIClosure(t, binary, f, o)
		})
	}
}

func automaticTestReceipt(t *testing.T, o WorkflowRun) deliveryVerificationReceipt {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(o.Delivery.Attempt.Path, "verification.json"))
	if err != nil {
		t.Fatal(err)
	}
	var r deliveryVerificationReceipt
	if err = json.Unmarshal(raw, &r); err != nil || r.Automatic == nil {
		t.Fatalf("automatic receipt unreadable: %v", err)
	}
	return r
}

func automaticTestScratch(t *testing.T, r deliveryVerificationReceipt) string {
	t.Helper()
	for _, entry := range r.Automatic.Environment {
		if value, ok := strings.CutPrefix(entry, "TMPDIR="); ok {
			return value
		}
	}
	t.Fatal("receipt has no TMPDIR")
	return ""
}

// AL25: bootstrap/provider ownership is modeled by the existing fixture. Verify
// and subsequent readback are real, separate candidate CLI processes.
func TestAutomaticAcceptanceCLIFailedChildEvidenceAfterExit(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	deliveryTestCandidate(t, f, "failed child evidence")
	script := "#!/bin/sh\nset -eu\nmkdir -m 700 \"$TMPDIR/child\"\nprintf 'useful child failure\\n' > \"$TMPDIR/child/failure.log\"\nprintf 'child evidence: %s/child/failure.log\\n' \"$TMPDIR\"\nprintf 'synthetic verifier failed\\n' >&2\nexit 3\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "child-evidence")
	invoke := func(name string, args ...string) (WorkflowRun, []byte, error) {
		t.Helper()
		argv := append([]string{"workflow"}, append(args, "--format", "json")...)
		c := exec.Command(binary, argv...)
		c.Dir = f.Prepared.Preparation.Preparation.Plan.WorktreePath
		var stdout, stderr bytes.Buffer
		c.Stdout, c.Stderr = &stdout, &stderr
		err := c.Run()
		var run WorkflowRun
		if e := json.Unmarshal(stdout.Bytes(), &run); e != nil {
			t.Fatalf("%s did not return one valid JSON result: %v\n%s\n%s", name, e, stdout.Bytes(), stderr.Bytes())
		}
		if dir := os.Getenv("PLY_AL25_EVIDENCE_DIR"); dir != "" {
			rawArgv, _ := json.Marshal(append([]string{binary}, argv...))
			for suffix, raw := range map[string][]byte{"stdout": stdout.Bytes(), "stderr": stderr.Bytes(), "argv.json": rawArgv, "exit": []byte(fmt.Sprint(c.ProcessState.ExitCode())), "cwd": []byte(c.Dir)} {
				if e := os.WriteFile(filepath.Join(dir, name+"."+suffix), raw, 0600); e != nil {
					t.Fatal(e)
				}
			}
		}
		return run, stderr.Bytes(), err
	}
	failed, stderr, err := invoke("failed-verify", automaticCLIVerifyArgs(o, review, binary)...)
	if err == nil || failed.DeliveryStatus.Verification.Outcome != "failed" || len(failed.Delivery.Candidates) != 0 {
		t.Fatalf("failure was not preserved: %v %+v", err, failed.DeliveryStatus.Verification)
	}
	r := automaticTestReceipt(t, failed)
	scratch := automaticTestScratch(t, r)
	if os.Getenv("PLY_AL25_EVIDENCE_DIR") == "" {
		t.Cleanup(func() { _ = os.RemoveAll(scratch) })
	}
	for _, fact := range []string{"level=warning", "Retained automatic verifier scratch", r.RunID, r.AttemptID, scratch} {
		if !bytes.Contains(stderr, []byte(fact)) {
			t.Fatalf("missing diagnostic %q: %s", fact, stderr)
		}
	}
	readback, _, err := invoke("failed-readback", "execute", "show", o.RunID)
	if err != nil || readback.DeliveryStatus.Verification.Outcome != "failed" || len(readback.Delivery.Candidates) != 0 || readback.Delivery.Attempt.ID != r.AttemptID || r.RunID != o.RunID || r.Exit == nil || *r.Exit != 3 || r.Outcome != "fail" {
		t.Fatalf("public readback changed failed attempt: %v %+v", err, readback.DeliveryStatus.Verification)
	}
	if raw, err := os.ReadFile(filepath.Join(scratch, "child", "failure.log")); err != nil || string(raw) != "useful child failure\n" {
		t.Fatalf("child evidence missing after both CLI exits: %q %v", raw, err)
	}
	if info, err := os.Stat(scratch); err != nil || info.Mode().Perm() != 0700 {
		t.Fatalf("retained evidence is not private: %v %v", info, err)
	}
	if dir := os.Getenv("PLY_AL25_EVIDENCE_DIR"); dir != "" {
		raw, _ := json.MarshalIndent(r, "", "  ")
		if err := os.WriteFile(filepath.Join(dir, "failed-receipt.json"), raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("Failed Run %s attempt %s retained child evidence at %s", r.RunID, r.AttemptID, scratch)
}

// A bounded green verification/readback control. Stop before any integration,
// follow or closeout, which are outside this scratch-lifetime regression.
func TestAutomaticAcceptanceCLIGreenScratchCleanupAfterExit(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	deliveryTestCandidate(t, f, "green scratch cleanup")
	script := "#!/bin/sh\nset -eu\nprintf 'disposable child data\\n' > \"$TMPDIR/child.log\"\nprintf 'green scratch used\\n'\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "green-scratch")
	invoke := func(name string, args ...string) WorkflowRun {
		t.Helper()
		argv := append([]string{"workflow"}, append(args, "--format", "json")...)
		c := exec.Command(binary, argv...)
		c.Dir = f.Prepared.Preparation.Preparation.Plan.WorktreePath
		var stdout, stderr bytes.Buffer
		c.Stdout, c.Stderr = &stdout, &stderr
		err := c.Run()
		if dir := os.Getenv("PLY_AL25_EVIDENCE_DIR"); dir != "" {
			rawArgv, _ := json.Marshal(append([]string{binary}, argv...))
			for suffix, raw := range map[string][]byte{"stdout": stdout.Bytes(), "stderr": stderr.Bytes(), "argv.json": rawArgv, "exit": []byte(fmt.Sprint(c.ProcessState.ExitCode())), "cwd": []byte(c.Dir)} {
				if e := os.WriteFile(filepath.Join(dir, name+"."+suffix), raw, 0600); e != nil {
					t.Fatal(e)
				}
			}
		}
		var run WorkflowRun
		if e := json.Unmarshal(stdout.Bytes(), &run); e != nil || err != nil {
			t.Fatalf("%s failed or returned invalid JSON: %v %v\n%s\n%s", name, err, e, stdout.Bytes(), stderr.Bytes())
		}
		return run
	}
	passed := invoke("green-verify", automaticCLIVerifyArgs(o, review, binary)...)
	r := automaticTestReceipt(t, passed)
	scratch := automaticTestScratch(t, r)
	readback := invoke("green-readback", "execute", "show", o.RunID)
	if readback.DeliveryStatus.Verification.Outcome != "passed" || len(readback.Delivery.Candidates) != 1 || readback.Delivery.Attempt.ID != r.AttemptID || r.RunID != o.RunID || r.Exit == nil || *r.Exit != 0 || r.Outcome != "pass" {
		t.Fatalf("green readback differs from successful attempt: %+v %+v", readback.DeliveryStatus.Verification, r)
	}
	if _, err := os.Lstat(scratch); !os.IsNotExist(err) {
		t.Fatalf("successful attempt scratch remains after CLI exits: %v", err)
	}
	if raw, err := os.ReadFile(r.Stdout.Locator); err != nil || string(raw) != "green scratch used\n" {
		t.Fatalf("green stdio evidence did not survive cleanup: %q %v", raw, err)
	}
	t.Logf("Passed Run %s attempt %s cleaned scratch %s; stdio and readback survive", r.RunID, r.AttemptID, scratch)
}
