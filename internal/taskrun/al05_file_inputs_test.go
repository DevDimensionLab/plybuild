//go:build darwin

package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// AL05: owner-submitted reports, evidence, acceptance entrypoints and reviews
// are hash-bound ordinary inputs. Readability by other local users grants no
// authority, so 0644/0755 inputs continue. Native private files keep their
// owner-only requirement, and every retained requirement names its input.

func al05Write(t *testing.T, path string, b []byte, mode os.FileMode) FileBinding {
	t.Helper()
	if err := os.WriteFile(path, b, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
	return FileBinding{path, hash(b)}
}

func al05Report(o WorkflowRun, event string, evidence []FileBinding, verifierLog FileBinding) DeliveryReport {
	return DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: event, PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: "working", Summary: "Focused regression passed.", Meaning: "Progress only; no candidate is qualified by this report.", Evidence: evidence, VerifierResults: []DeliveryVerifierResult{{ID: "focused", Outcome: "passed", Argv: []string{"go", "test"}, CWD: "/", Exit: ptr(0), Evidence: []FileBinding{verifierLog}, Reason: "Shell-redirected log of the actual run."}}}
}

func al05WriteReport(t *testing.T, path string, r DeliveryReport, mode os.FileMode) string {
	t.Helper()
	b, err := Canonical(r)
	if err != nil {
		t.Fatal(err)
	}
	al05Write(t, path, b, mode)
	return path
}

func al05AssertMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil || info.Mode().Perm() != want {
		t.Fatalf("%s mode differs from %04o: %v %v", path, want, info, err)
	}
}

func al05AssertError(t *testing.T, err error, facts ...string) {
	t.Helper()
	if err == nil {
		t.Fatalf("retained requirement passed; want %q", facts)
	}
	for _, fact := range facts {
		if !strings.Contains(err.Error(), fact) {
			t.Fatalf("diagnostic lacks %q: %v", fact, err)
		}
	}
	if strings.Contains(err.Error(), "AL05-CONTENT") {
		t.Fatalf("diagnostic leaked input contents: %v", err)
	}
}

func al05AssertWarn(t *testing.T, log *bytes.Buffer, facts ...string) {
	t.Helper()
	for _, fact := range append([]string{"level=warning"}, facts...) {
		if !strings.Contains(log.String(), fact) {
			t.Fatalf("warn lacks %q: %s", fact, log.String())
		}
	}
}

func TestAL05OrdinaryReportAndEvidenceAreConsumed(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	log := al04CaptureLog(t)
	root := f.R.WorkspaceRoot
	evidence := al05Write(t, filepath.Join(root, "acceptance.log"), []byte("15 PASS\n"), 0644)
	verifierLog := al05Write(t, filepath.Join(root, "verifier.log"), []byte("ok\n"), 0644)
	file := al05WriteReport(t, filepath.Join(root, "report.json"), al05Report(o, "progress-001", []FileBinding{evidence}, verifierLog), 0644)
	o, err := WorkflowDeliveryReport(f.D, root, o.RunID, o.Paths.Context, file)
	if err != nil {
		t.Fatalf("ordinary 0644 report and evidence did not continue: %v", err)
	}
	if o.Delivery.Phase != "working" || len(o.Delivery.Events) != 1 || len(o.Delivery.Candidates) != 0 {
		t.Fatalf("report not preserved as progress only: %+v", o.Delivery)
	}
	// Ply's own copy stays private; the submitted files are not rewritten.
	al05AssertMode(t, o.Delivery.Events[0].Binding.Locator, 0600)
	for _, path := range []string{file, evidence.Locator, verifierLog.Locator} {
		al05AssertMode(t, path, 0644)
	}
	al05AssertWarn(t, log, "delivery report", file, "report evidence", evidence.Locator, verifierLog.Locator, "mode=0644")
	// A lost reply is retried with the same bytes: no second event.
	o, err = WorkflowDeliveryReport(f.D, root, o.RunID, o.Paths.Context, file)
	if err != nil || len(o.Delivery.Events) != 1 {
		t.Fatalf("lost-reply retry duplicated or failed: %v %+v", err, o.Delivery.Events)
	}
	// Historical qualification later re-reads the same evidence by digest.
	s, err := workflowRead(root, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := workflowBound(s.Result.Delivery.Events[0].Binding, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	var report DeliveryReport
	if err = decode(raw, 1<<20, &report); err != nil || validateDeliveryReport(report) != nil {
		t.Fatalf("preserved report no longer validates with its ordinary evidence: %v", err)
	}
}

func TestAL05RetainedReportRequirementsNameInputAndCause(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	root := f.R.WorkspaceRoot
	good := al05Write(t, filepath.Join(root, "good.log"), []byte("ok\n"), 0644)
	report := func(name string, evidence FileBinding, mode os.FileMode) string {
		return al05WriteReport(t, filepath.Join(root, name), al05Report(o, "progress-"+strings.TrimSuffix(name, ".json"), []FileBinding{evidence}, good), mode)
	}
	writable := report("writable.json", good, 0664)
	dir := filepath.Join(root, "directory.json")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	oversized := filepath.Join(root, "oversized.json")
	al05Write(t, oversized, append([]byte("AL05-CONTENT"), bytes.Repeat([]byte("x"), 1<<20)...), 0644)
	link := filepath.Join(root, "link.json")
	if err := os.Symlink(report("target.json", good, 0644), link); err != nil {
		t.Fatal(err)
	}
	changed := al05Write(t, filepath.Join(root, "changed.log"), []byte("AL05-CONTENT before\n"), 0644)
	changedReport := report("changed.json", changed, 0644)
	al05Write(t, changed.Locator, []byte("AL05-CONTENT after\n"), 0644)
	sharedWritable := al05Write(t, filepath.Join(root, "shared.log"), []byte("AL05-CONTENT\n"), 0666)
	for _, tc := range []struct {
		name, file string
		facts      []string
	}{
		{"writable report", writable, []string{"delivery report", writable, "writable by group or other users", "mode 0664", "chmod go-w"}},
		{"directory", dir, []string{"delivery report", dir, "regular file", "directory"}},
		{"oversized", oversized, []string{"delivery report", oversized, "at most 1048576 bytes", "1048588 bytes"}},
		{"symlink", link, []string{"delivery report", link, "symlink"}},
		{"evidence changed after binding", changedReport, []string{"report evidence", changed.Locator, "changed"}},
		{"writable evidence", report("shared.json", sharedWritable, 0644), []string{"report evidence", sharedWritable.Locator, "mode 0666"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := WorkflowDeliveryReport(f.D, root, o.RunID, o.Paths.Context, tc.file)
			al05AssertError(t, err, tc.facts...)
			s, readErr := workflowRead(root, o.RunID)
			if readErr != nil || len(s.Result.Delivery.Events) != 0 || s.Result.Delivery.LastEventSHA256 != nil {
				t.Fatalf("rejected input appended an event: %v %+v", readErr, s.Result.Delivery.Events)
			}
		})
	}
}

// The candidate binary was already read with private=false and snapshotted
// 0700 before AL05; this is a baseline, not a new behavior.
func TestAL05BaselineOrdinaryCandidateBinary(t *testing.T) {
	f, o, binary, review := autoScratchFixture(t, "#!/bin/sh\nset -eu\n\"$PLY_CANDIDATE_BINARY\"\n")
	if err := os.Chmod(binary, 0755); err != nil {
		t.Fatal(err)
	}
	o, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err != nil || len(o.Delivery.Candidates) != 1 {
		t.Fatalf("0755 candidate binary did not qualify: %v %+v", err, o.DeliveryStatus.Verification)
	}
	r := automaticTestReceipt(t, o)
	al05AssertMode(t, r.Automatic.CandidateBinarySnapshot.Locator, 0700)
	al05AssertMode(t, binary, 0755)
}

func TestAL05OrdinaryAcceptanceAndReviewQualify(t *testing.T) {
	f, o, binary, review := autoScratchFixture(t, "#!/bin/sh\nset -eu\n\"$PLY_CANDIDATE_BINARY\"\n")
	log := al04CaptureLog(t)
	for path, mode := range map[string]os.FileMode{f.AcceptancePath: 0644, review: 0644, binary: 0755} {
		if err := os.Chmod(path, mode); err != nil {
			t.Fatal(err)
		}
	}
	o, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err != nil || len(o.Delivery.Candidates) != 1 || o.Delivery.Phase != "automatic_acceptance_passed" {
		t.Fatalf("ordinary 0644 acceptance and review did not qualify: %v %+v", err, o.DeliveryStatus.Verification)
	}
	if v := o.DeliveryStatus.Verification; v == nil || !v.InputsMatch {
		t.Fatalf("status no longer matches ordinary inputs: %+v", v)
	}
	r := automaticTestReceipt(t, o)
	al05AssertMode(t, r.AcceptanceSnapshot.Locator, 0600)
	if err = deliveryVerificationInputs(r); err != nil {
		t.Fatalf("receipt reuse cannot re-read ordinary inputs: %v", err)
	}
	al05AssertWarn(t, log, "acceptance entrypoint", f.AcceptancePath, "candidate review", review, "mode=0644")
}

func TestAL05ChangedOrdinaryInputsDoNotQualify(t *testing.T) {
	for _, target := range []string{"acceptance entrypoint", "candidate review"} {
		t.Run(target, func(t *testing.T) {
			f, o, binary, review := autoScratchFixture(t, "")
			changed := f.AcceptancePath
			if target == "candidate review" {
				changed = review
			}
			script := "#!/bin/sh\nset -eu\nprintf ' ' >> " + ShellQuote(changed) + "\n"
			al05Write(t, f.AcceptancePath, []byte(script), 0644)
			if err := os.Chmod(review, 0644); err != nil {
				t.Fatal(err)
			}
			o, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
			al05AssertError(t, err, target, changed, "changed")
			if len(o.Delivery.Candidates) != 0 {
				t.Fatalf("changed %s qualified a candidate: %+v", target, o.Delivery.Candidates)
			}
		})
	}
}

// The existing fault boundary sits after the inputs were read and bound and
// before the test runs. It proves the bound-byte protection, not detection
// inside readFile: no read seam exists, and an in-place same-inode, same-size
// rewrite during the read itself is not detected there (see the AL05 audit).
func TestAL05InPlaceOverwriteAfterBindingDoesNotQualify(t *testing.T) {
	for _, target := range []string{"acceptance entrypoint", "candidate review"} {
		t.Run(target, func(t *testing.T) {
			f, o, binary, review := autoScratchFixture(t, "")
			marker := filepath.Join(t.TempDir(), "executed")
			script := func(word string) []byte {
				return []byte("#!/bin/sh\nprintf " + word + " > " + ShellQuote(marker) + "\n")
			}
			al05Write(t, f.AcceptancePath, script("original"), 0644)
			if err := os.Chmod(review, 0644); err != nil {
				t.Fatal(err)
			}
			changed := f.AcceptancePath
			replacement := script("replaced")
			if target == "candidate review" {
				changed = review
				raw, err := os.ReadFile(review)
				if err != nil {
					t.Fatal(err)
				}
				replacement = bytes.Replace(raw, []byte("fixture-review-session"), []byte("fixture-review-SESSION"), 1)
			}
			f.D.Fault = func(point string) error {
				if point != "delivery_after_verification_reservation" {
					return nil
				}
				before, err := os.Stat(changed)
				if err != nil {
					return err
				}
				if err = os.WriteFile(changed, replacement, 0644); err != nil {
					return err
				}
				after, err := os.Stat(changed)
				if err != nil || !os.SameFile(before, after) || before.Size() != after.Size() {
					t.Errorf("overwrite was not in place with the same size: %v", err)
				}
				return nil
			}
			o, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
			al05AssertError(t, err, target, changed, "changed since they were bound")
			if len(o.Delivery.Candidates) != 0 {
				t.Fatalf("overwritten %s qualified a candidate", target)
			}
			if raw, readErr := os.ReadFile(marker); readErr != nil || string(raw) != "original" {
				t.Fatalf("executed bytes were not the bound snapshot: %q %v", raw, readErr)
			}
		})
	}
}

func TestAL05NativePrivateRolesStayProtected(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	root := f.R.WorkspaceRoot
	shared := al05Write(t, filepath.Join(root, "shared-private.json"), []byte("{\"AL05-CONTENT\":1}"), 0644)
	_, err := PrivateRead(shared.Locator, 1<<20)
	al05AssertError(t, err, shared.Locator, "owner-only", "mode 0644")
	// Runtime acceptance carries the actual authority claim; it stays owner-only.
	acceptance := writeAny(t, root, "runtime-acceptance.json", deliveryTestAcceptance(t, f, o))
	if err = os.Chmod(acceptance, 0644); err != nil {
		t.Fatal(err)
	}
	_, err = WorkflowAccept(f.D, root, o.RunID, o.Paths.Context, acceptance)
	al05AssertError(t, err, "runtime acceptance", acceptance, "owner-only", "mode 0644", "chmod go-rwx")
	s, err := workflowRead(root, o.RunID)
	if err != nil || s.Acceptance != nil {
		t.Fatalf("shared runtime acceptance changed state: %v %+v", err, s.Acceptance)
	}
	// A native store file whose mode changed outside Ply is not consumed.
	state := filepath.Join(o.Paths.RunRoot, "state.json")
	if err = os.Chmod(state, 0644); err != nil {
		t.Fatal(err)
	}
	_, err = workflowRead(root, o.RunID)
	al05AssertError(t, err, "native workflow run state", state, "owner-only", "mode 0644")
	if err = os.Chmod(state, 0600); err != nil {
		t.Fatal(err)
	}
}

// The exact built CLI consumes ordinary inputs, keeps one JSON value on stdout
// and warns on stderr, then delivers the same candidate without a mode question.
func TestAL05CandidateCLIJourneyWithOrdinaryModes(t *testing.T) {
	supplied := automaticCLIBinary(t)
	suppliedInfo, err := os.Stat(supplied)
	if err != nil {
		t.Fatal(err)
	}
	// Native verify supplies Ply's private candidate snapshot; never chmod it.
	t.Cleanup(func() {
		if info, err := os.Stat(supplied); err != nil || info.Mode() != suppliedInfo.Mode() {
			t.Errorf("supplied candidate binary %s changed mode: %v %v", supplied, info, err)
		}
	})
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(dir, "ply")
	b, err := os.ReadFile(supplied)
	if err != nil {
		t.Fatal(err)
	}
	al05Write(t, binary, b, 0755)
	f, o := automaticCLIFixture(t, binary, true, 30)
	root := f.R.WorkspaceRoot
	deliveryTestCandidate(t, f, "expected")
	al05Write(t, f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = expected\n\"$PLY_CANDIDATE_BINARY\" capabilities --format json >/dev/null\n"), 0644)
	evidence := al05Write(t, filepath.Join(root, "progress.log"), []byte("focused PASS\n"), 0644)
	file := al05WriteReport(t, filepath.Join(root, "progress.json"), al05Report(o, "cli-progress", []FileBinding{evidence}, evidence), 0644)
	run := func(args ...string) (WorkflowRun, string, error) {
		c := exec.Command(binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
		c.Dir = f.Prepared.Preparation.Preparation.Plan.WorktreePath
		var stdout, stderr bytes.Buffer
		c.Stdout, c.Stderr = &stdout, &stderr
		err := c.Run()
		var out WorkflowRun
		if decodeErr := json.Unmarshal(stdout.Bytes(), &out); decodeErr != nil {
			t.Fatalf("stdout is not one JSON value: %v\n%s", decodeErr, stdout.Bytes())
		}
		return out, stderr.String(), err
	}
	reported, stderr, err := run("execute", "report", o.RunID, "--context", o.Paths.Context, "--file", file)
	if err != nil || len(reported.Delivery.Events) != 1 {
		t.Fatalf("CLI report with ordinary modes failed: %v\n%s", err, stderr)
	}
	for _, fact := range []string{"delivery report", file, "report evidence", evidence.Locator} {
		if !strings.Contains(stderr, fact) {
			t.Fatalf("stderr warn lacks %q: %s", fact, stderr)
		}
	}
	if again, _, err := run("execute", "report", o.RunID, "--context", o.Paths.Context, "--file", file); err != nil || len(again.Delivery.Events) != 1 {
		t.Fatalf("lost-reply CLI retry duplicated or failed: %v %+v", err, again.Delivery.Events)
	}
	writable := al05WriteReport(t, filepath.Join(root, "writable.json"), al05Report(reported, "cli-writable", []FileBinding{evidence}, evidence), 0666)
	if _, _, err := run("execute", "report", o.RunID, "--context", o.Paths.Context, "--file", writable); err == nil {
		t.Fatal("group/other-writable report was consumed")
	} else if exit := new(exec.ExitError); !errors.As(err, &exit) {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "cli")
	if err := os.Chmod(review, 0644); err != nil {
		t.Fatal(err)
	}
	verified, stderr, err := run(automaticCLIVerifyArgs(o, review, binary)...)
	if err != nil || verified.Delivery.Phase != "automatic_acceptance_passed" {
		t.Fatalf("CLI verify with ordinary modes failed: %v %+v\n%s", err, verified.DeliveryStatus.Verification, stderr)
	}
	for _, fact := range []string{"acceptance entrypoint", f.AcceptancePath, "candidate review", review} {
		if !strings.Contains(stderr, fact) {
			t.Fatalf("stderr warn lacks %q: %s", fact, stderr)
		}
	}
	assertAutomaticCLIClosure(t, binary, f, verified)
}

// r001-f01: the submitted round report is an ordinary input like the delivery
// report. The coordinator review stays owner-only as a control decision.
func TestAL05OrdinaryRoundReportIsConsumed(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestAccept(t, f, workflowTestStart(t, f))
	log := al04CaptureLog(t)
	r := f.Report
	r.Round, r.ControlID, r.PreviousReportSHA256 = o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256
	b, err := Canonical(r)
	if err != nil {
		t.Fatal(err)
	}
	writable := filepath.Join(f.R.WorkspaceRoot, "round-writable.json")
	al05Write(t, writable, b, 0664)
	_, err = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writable)
	al05AssertError(t, err, "round report", writable, "mode 0664", "chmod go-w")
	file := filepath.Join(f.R.WorkspaceRoot, "round.json")
	al05Write(t, file, b, 0644)
	for i := 0; i < 2; i++ {
		o, err = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file)
		if err != nil || o.Round.ReportSHA256 == nil {
			t.Fatalf("ordinary 0644 round report did not continue: %v", err)
		}
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || len(s.Records) != 1 {
		t.Fatalf("lost-reply retry duplicated the round record: %v %d", err, len(s.Records))
	}
	// Later reads use Ply's private copy of the report, not the submitted path.
	al05AssertMode(t, s.Records[0].Report.Locator, 0600)
	if _, err = workflowReadReport(s.Records[0]); err != nil {
		t.Fatal(err)
	}
	al05AssertWarn(t, log, "round report", file, "mode=0644")
	review := workflowTestReview(t, f, o, "accepted")
	if err = os.Chmod(review, 0644); err != nil {
		t.Fatal(err)
	}
	_, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, review)
	al05AssertError(t, err, "coordinator review", review, "owner-only", "mode 0644")
}

// r001-f02: a missing or unopenable ordinary input names its role, path and
// next step and keeps the actual errno; native callers keep os.IsNotExist.
func TestAL05MissingOrUnopenableInputKeepsErrnoAndNamesRole(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	root := f.R.WorkspaceRoot
	missing := filepath.Join(root, "missing-report.json")
	_, err := WorkflowDeliveryReport(f.D, root, o.RunID, o.Paths.Context, missing)
	al05AssertError(t, err, "delivery report", missing, "no such file or directory", "existing regular file")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("missing input lost its errno: %#v", err)
	}
	unreadable := al05WriteReport(t, filepath.Join(root, "unreadable.json"), al05Report(o, "unreadable", []FileBinding{}, al05Write(t, filepath.Join(root, "ok.log"), []byte("ok\n"), 0644)), 0200)
	_, err = WorkflowDeliveryReport(f.D, root, o.RunID, o.Paths.Context, unreadable)
	al05AssertError(t, err, "delivery report", unreadable, "permission denied")
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("unopenable input lost its errno: %#v", err)
	}
	gone := filepath.Join(root, "gone.log")
	report := al05WriteReport(t, filepath.Join(root, "gone.json"), al05Report(o, "gone", []FileBinding{{gone, hash([]byte("x"))}}, FileBinding{gone, hash([]byte("x"))}), 0644)
	_, err = WorkflowDeliveryReport(f.D, root, o.RunID, o.Paths.Context, report)
	al05AssertError(t, err, "report evidence", gone, "no such file or directory")
	// Native recovery still depends on the unwrapped not-exist error.
	var v map[string]any
	if err = readValue(filepath.Join(o.Paths.RunRoot, "absent.json"), 1<<10, &v); !os.IsNotExist(err) {
		t.Fatalf("native missing-state recovery changed: %v", err)
	}
	s, err := workflowRead(root, o.RunID)
	if err != nil || len(s.Result.Delivery.Events) != 0 {
		t.Fatalf("refused input appended an event: %v", err)
	}
}

// The candidate CLI consumes a 0644 round report: one record, warn on stderr.
func TestAL05RoundReportCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	installed := filepath.Join(root, "ply")
	al04Install(t, installed, mustReadAL01(t, binary))
	f := al01ReceiptFixture(t, installed)
	o := workflowTestStart(t, f)
	cwd, _ := f.D.CWD()
	al04CLI(t, installed, cwd, "run", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID)
	r := f.Report
	r.Round, r.ControlID, r.PreviousReportSHA256 = o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256
	b, err := Canonical(r)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(f.R.WorkspaceRoot, "cli-round.json")
	al05Write(t, file, b, 0644)
	for i := 0; i < 2; i++ {
		raw, stderr := al04CLI(t, installed, cwd, "run", "report", o.RunID, "--context", o.Paths.Context, "--file", file)
		var got WorkflowRun
		if err = json.Unmarshal(raw, &got); err != nil || got.Round.ReportSHA256 == nil {
			t.Fatalf("CLI round report was not recorded: %v %s", err, raw)
		}
		for _, fact := range []string{"level=warning", "round report", file} {
			if !strings.Contains(stderr, fact) {
				t.Fatalf("stderr warn lacks %q: %s", fact, stderr)
			}
		}
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || len(s.Records) != 1 {
		t.Fatalf("repeated CLI round report produced another record: %v", err)
	}
}
