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
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
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
	status := automaticCLIOK(t, binary, f, "execute", "show", o.RunID)
	if status.Delivery.Phase != "completed" || status.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("native closure invisible: %+v", status.DeliveryStatus)
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
	done := deliveryCLIOK(t, binary, cwd, "execute", id)
	if deliveryCLIString(t, done, "state") != "delivered" {
		t.Fatalf("merge recovery incomplete: %+v", done)
	}
	deliveryCLIOK(t, binary, cwd, "execute", id)
	if workflowProviderDocument(t, model)["merge_count"] != float64(1) {
		t.Fatal("merge repeated")
	}
	run := automaticCLIOK(t, binary, f, "execute", "show", o.RunID)
	if run.Delivery.Phase != "completed" || run.DeliveryStatus.HumanJudgment != nil {
		t.Fatal("wrong native closure or fabricated human QA")
	}
}
