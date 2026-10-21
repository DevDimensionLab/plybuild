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
)

// Only bootstrap/provider observations are modeled. Real candidate commands end
// at report/readback: no verification, integration, follow or closeout fixture.
func TestDeliveryReportCLICorrectionReadback(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("private AL26 acceptance supplies an actual make-built PLY_DELIVERY_TEST_BINARY")
	}
	evidence := os.Getenv("PLY_AL26_EVIDENCE_DIR")
	if evidence == "" {
		evidence = t.TempDir()
	}
	dir, err := os.MkdirTemp(evidence, "report-cli-")
	if err != nil {
		t.Fatal(err)
	}
	capture := func(name, cwd, program string, args ...string) ([]byte, []byte, int) {
		t.Helper()
		c := exec.Command(program, args...)
		c.Dir = cwd
		var out, stderr bytes.Buffer
		c.Stdout = &out
		c.Stderr = &stderr
		err := c.Run()
		exit := 0
		if err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				exit = e.ExitCode()
			} else {
				t.Fatal(err)
			}
		}
		for suffix, b := range map[string][]byte{"stdout": out.Bytes(), "stderr": stderr.Bytes()} {
			if err := os.WriteFile(filepath.Join(dir, name+"."+suffix), b, 0600); err != nil {
				t.Fatal(err)
			}
		}
		meta, _ := json.Marshal(map[string]any{"argv": append([]string{program}, args...), "cwd": cwd, "exit": exit})
		if err := os.WriteFile(filepath.Join(dir, name+".json"), meta, 0600); err != nil {
			t.Fatal(err)
		}
		return out.Bytes(), stderr.Bytes(), exit
	}
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	_, _, exit := capture("focused-test", repo, "go", "test", "./internal/taskrun", "-run", "^TestDeliveryReportVerifierFieldDiagnostics$", "-count=1")
	if exit != 0 {
		t.Fatal("focused command failed")
	}
	f := receiptFixture(t, "codex")
	f.R.Runtime.PlyExecutable = Executable{Path: binary, SHA256: hashFileTest(t, binary)}
	f.D.Executable = func() (string, error) { return binary, nil }
	ownerPaneStandin(t, &f)
	o := workflowTestStart(t, f.workflowFixture)
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	cli := func(name string, args ...string) (WorkflowRun, string, int) {
		t.Helper()
		raw, stderr, code := capture(name, cwd, binary, append(append([]string{"workflow"}, args...), "--format", "json")...)
		var got WorkflowRun
		if code == 0 {
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("invalid CLI readback: %v", err)
			}
		}
		return got, string(raw) + string(stderr), code
	}
	o, _, exit = cli("receipt", "execute", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID)
	if exit != 0 {
		t.Fatal("receipt failed")
	}
	r := al26Report()
	r.RunID, r.RequestSHA256, r.SessionID = o.RunID, o.RequestSHA256, o.SessionID
	r.PreviousEventSHA256 = o.Delivery.LastEventSHA256
	r.VerifierResults = []DeliveryVerifierResult{{ID: "TestDeliveryReportVerifierFieldDiagnostics", Outcome: "passed", Exit: ptr(0), Argv: []string{"go", "test", "./internal/taskrun", "-run", "^TestDeliveryReportVerifierFieldDiagnostics$", "-count=1"}, CWD: repo, Reason: "Actual named command exited 0; raw stdout and stderr are bound.", Evidence: []FileBinding{{filepath.Join(dir, "focused-test.stdout"), hashFileTest(t, filepath.Join(dir, "focused-test.stdout"))}, {filepath.Join(dir, "focused-test.stderr"), hashFileTest(t, filepath.Join(dir, "focused-test.stderr"))}}}}
	submit := func(name string, r DeliveryReport) (WorkflowRun, string, int) {
		return cli(name, "execute", "report", o.RunID, "--context", o.Paths.Context, "--file", writeAny(t, dir, name+"-report.json", r))
	}
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	for i, kind := range []string{"id", "reason", "duplicate", "session", "predecessor", "evidence"} {
		bad := r
		bad.VerifierResults = append([]DeliveryVerifierResult{}, r.VerifierResults...)
		want := ""
		switch kind {
		case "id":
			want = "verifier_results[0].id"
		case "reason":
			bad.VerifierResults[0].ID = "focused-test"
			bad.VerifierResults[0].Reason = ""
			want = "verifier_results[0].reason"
		case "duplicate":
			bad.VerifierResults[0].ID = "focused-test"
			bad.VerifierResults = append(bad.VerifierResults, bad.VerifierResults[0])
			want = "verifier_results[1].id"
		case "session":
			bad.VerifierResults[0].ID = "focused-test"
			bad.SessionID = "wrong-session"
		case "predecessor":
			bad.VerifierResults[0].ID = "focused-test"
			bad.PreviousEventSHA256 = ptr(hash([]byte("stale")))
		case "evidence":
			bad.VerifierResults[0].ID = "focused-test"
			bad.VerifierResults[0].Evidence = []FileBinding{{filepath.Join(dir, "focused-test.stdout"), hash([]byte("wrong"))}}
		}
		_, message, code := submit(fmt.Sprintf("invalid-%d-%s", i, kind), bad)
		if code == 0 || !strings.Contains(message, want) {
			t.Fatalf("%s unexpectedly accepted or opaque: %d %s", kind, code, message)
		}
		after, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
		if err != nil {
			t.Fatal(err)
		}
		if !equal(before.Result.Delivery, after.Result.Delivery) {
			t.Fatal("rejection changed delivery state")
		}
	}
	r.VerifierResults[0].ID = "focused-test"
	got, message, exit := submit("corrected", r)
	if exit != 0 {
		t.Fatal(message)
	}
	repeated, message, exit := submit("repeated", r)
	if exit != 0 || !equal(got.Delivery.LastEventSHA256, repeated.Delivery.LastEventSHA256) {
		t.Fatal(message)
	}
	changed := r
	changed.Summary = "Changed immutable bytes."
	if _, _, exit = submit("changed-event", changed); exit == 0 {
		t.Fatal("event bytes replaced")
	}
	empty := r
	empty.EventID = "empty-progress"
	empty.PreviousEventSHA256 = got.Delivery.LastEventSHA256
	empty.VerifierResults = []DeliveryVerifierResult{}
	got, message, exit = submit("empty", empty)
	if exit != 0 {
		t.Fatal(message)
	}
	shown, message, exit := cli("readback", "execute", "show", o.RunID)
	if exit != 0 {
		t.Fatal(message)
	}
	if len(shown.Delivery.Events) != len(before.Result.Delivery.Events)+2 || len(shown.Delivery.Candidates) != 0 || shown.Delivery.Attempt != nil || shown.Delivery.PermissionState != RuntimePerAction || !equal(got.Delivery.LastEventSHA256, shown.Delivery.LastEventSHA256) || shown.RequestSHA256 != o.RequestSHA256 || shown.SessionID != o.SessionID {
		t.Fatalf("report changed qualification or identity: %+v", shown.Delivery)
	}
	t.Logf("Raw command evidence: %s", dir)
}
