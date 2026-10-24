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

	"github.com/sirupsen/logrus"
)

// AL04: a launched report owner completes its round after ordinary installs.
// Fixture Ply/provider bytes are synthetic; the CLI journey runs real builds.

// Replace an installed file the way an ordinary install does: new bytes at the
// same path through rename, so an already running process keeps its own image.
func al04Install(t *testing.T, path string, b []byte) {
	t.Helper()
	next := path + ".install"
	if err := os.WriteFile(next, b, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(next, path); err != nil {
		t.Fatal(err)
	}
}

func al04CaptureLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var log bytes.Buffer
	previous := logrus.StandardLogger().Out
	logrus.SetOutput(&log)
	t.Cleanup(func() { logrus.SetOutput(previous) })
	return &log
}

func al04ReceiptOwner(t *testing.T) (workflowFixture, WorkflowRun, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ply := filepath.Join(root, "installed-ply")
	if err = os.WriteFile(ply, []byte("launch ply"), 0700); err != nil {
		t.Fatal(err)
	}
	f := al01ReceiptFixture(t, ply)
	o := workflowTestStart(t, f)
	o, err = WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if err != nil {
		t.Fatal(err)
	}
	return f, o, ply
}

func al04AssertSameOwner(t *testing.T, f workflowFixture, before workflowState) workflowState {
	t.Helper()
	s, err := workflowRead(f.R.WorkspaceRoot, before.Result.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if !equal(s.Request, before.Request) || s.Result.RequestSHA256 != before.Result.RequestSHA256 || s.Result.SessionID != before.Result.SessionID || s.Result.Transport.AgentSessionID != before.Result.Transport.AgentSessionID || s.Result.Paths.Context != before.Result.Paths.Context || s.ContextSHA256 != before.ContextSHA256 || !equal(s.Acceptance, before.Acceptance) || !equal(s.StartSHA256, before.StartSHA256) || s.Result.Handoff != before.Result.Handoff {
		t.Fatal("drift rewrote the original request, mandate, receipt, context or session")
	}
	for _, c := range []string{"tab create", "agent start"} {
		if workflowTestCalls(t, f, c) != 1 {
			t.Fatal("drift started another provider or tab: " + c)
		}
	}
	return s
}

func al04AssertWarn(t *testing.T, log *bytes.Buffer, id string, facts ...string) {
	t.Helper()
	for _, s := range append([]string{"level=warning", "run_id=" + id, "post-launch"}, facts...) {
		if !strings.Contains(log.String(), s) {
			t.Fatalf("warn evidence lacks %q: %s", s, log.String())
		}
	}
}

func TestAL04ReportOwnerCompletesRoundAfterPlyUpdate(t *testing.T) {
	f, o, ply := al04ReceiptOwner(t)
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	al04Install(t, ply, []byte("updated ply"))
	log := al04CaptureLog(t)
	file := writeAny(t, f.R.WorkspaceRoot, "al04-round.json", func() WorkflowReport {
		r := f.Report
		r.Round, r.ControlID, r.PreviousReportSHA256 = o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256
		return r
	}())
	reported, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file)
	if err != nil {
		t.Fatalf("ordinary CLI update locked the working report owner: %v", err)
	}
	al04AssertWarn(t, log, o.RunID, "ply_executable")
	// A lost reply repeats the same report: same record, no second dispatch.
	again, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file)
	if err != nil || again.Round.ReportSHA256 == nil || *again.Round.ReportSHA256 != *reported.Round.ReportSHA256 {
		t.Fatalf("repeated report after lost reply differs: %v", err)
	}
	s := al04AssertSameOwner(t, f, before)
	if len(s.Records) != 1 || workflowTestCalls(t, f, "agent prompt") != 1 {
		t.Fatal("lost reply produced a second report or dispatch")
	}
	changed := f.Report
	changed.Summary = "A different report"
	if _, err = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "al04-changed.json", changed)); err == nil {
		t.Fatal("changed report replaced the recorded round")
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range shown.Reasons {
		if r.Code == "workflow_run_drift" {
			t.Fatalf("historical launch fact still reported as drift: %+v", r)
		}
	}
	o = workflowTestSettle(t, f, reported)
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if err != nil || o.FinalReturn.TerminalSHA256 == nil || o.TaskResultState != "not_published" || o.Delivery != nil {
		t.Fatalf("the same round did not complete with report-only authority: %+v %v", o, err)
	}
}

func TestAL04ProviderAndTrustDriftAfterLaunch(t *testing.T) {
	for _, removed := range []bool{false, true} {
		name := "provider_changed"
		if removed {
			name = "provider_removed"
		}
		t.Run(name, func(t *testing.T) {
			f := workflowTestFixture(t)
			home := workflowGrant(t, &f)
			o := workflowTestStart(t, f)
			// The grant is part of the request, so its digest binds the claims.
			f.Acceptance.RequestSHA256, f.Report.RequestSHA256 = o.RequestSHA256, o.RequestSHA256
			o = workflowTestAccept(t, f, o)
			if argv := workflowTrustAgentStart(t, f); !strings.Contains(argv, "trust_level") {
				t.Fatalf("fixture grant was not applied at launch: %s", argv)
			}
			before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if removed {
				if err = os.Remove(f.R.Runtime.Executable.Path); err != nil {
					t.Fatal(err)
				}
			} else {
				al04Install(t, f.R.Runtime.Executable.Path, []byte("#!/bin/sh\necho updated provider\n"))
			}
			workflowTrustAppendConfig(t, home, "model = \"edited-after-launch\"")
			git := exec.Command("git", "config", "ply.al04", "harmless")
			git.Dir = before.CodexTrust.RepositoryRoot
			if out, err := git.CombinedOutput(); err != nil {
				t.Fatalf("%v %s", err, out)
			}
			log := al04CaptureLog(t)
			o = workflowTestReport(t, f, o)
			if o.Round.ReportSHA256 == nil {
				t.Fatal("report was not recorded")
			}
			al04AssertWarn(t, log, o.RunID, "provider_executable", "codex_project_trust")
			al04AssertSameOwner(t, f, before)
			shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range shown.Reasons {
				if r.Code == "workflow_run_drift" {
					t.Fatalf("historical launch fact still reported as drift: %+v", r)
				}
			}
		})
	}
}

// Facts that a launch still consumes stay effective until the native session
// is bound. Delivery keeps its own continuation (TestDeliveryProviderUpgrade*).
func TestAL04PreLaunchFactsStayEffective(t *testing.T) {
	f := workflowTestFixture(t)
	home := workflowGrant(t, &f)
	o := workflowTestStart(t, f)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	provider := mustReadAL01(t, f.R.Runtime.Executable.Path)
	al04Install(t, f.R.Runtime.Executable.Path, []byte("#!/bin/sh\necho updated provider\n"))
	unbound := s
	unbound.Result.Transport.AgentSessionID = ""
	if err = workflowFresh(f.D, unbound, false); err == nil || !strings.Contains(err.Error(), "executable bytes changed") {
		t.Fatalf("an unlaunched report accepted a changed provider: %v", err)
	}
	al04Install(t, f.R.Runtime.Executable.Path, provider)
	workflowTrustAppendConfig(t, home, "model = \"edited-before-launch\"")
	if err = workflowFresh(f.D, unbound, false); err == nil || !strings.Contains(err.Error(), "Codex trust") {
		t.Fatalf("an unlaunched report accepted changed trust facts: %v", err)
	}
}

func TestAL04IdentityHistoryAndArtifactGuardsAfterDrift(t *testing.T) {
	f, o, ply := al04ReceiptOwner(t)
	al04Install(t, ply, []byte("updated ply"))
	report := f.Report
	report.Round, report.ControlID, report.PreviousReportSHA256 = o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256
	file := writeAny(t, f.R.WorkspaceRoot, "al04-guard.json", report)
	other := filepath.Join(f.R.WorkspaceRoot, "other-ply")
	if err := os.WriteFile(other, []byte("updated ply"), 0700); err != nil {
		t.Fatal(err)
	}
	executable := f.D.Executable
	f.D.Executable = func() (string, error) { return other, nil }
	if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err != nil {
		t.Fatalf("ordinary upgraded Ply could not report for the same owner: %v", err)
	}
	f.D.Executable = executable
	wrong := report
	wrong.SessionID = "ply:wrong"
	if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "al04-wrong.json", wrong)); err == nil {
		t.Fatal("wrong session reported after drift")
	}
	if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, filepath.Join(f.R.WorkspaceRoot, "context.json"), file); err == nil {
		t.Fatal("a different context reported after drift")
	}
	if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: "another-session"}, true); err == nil {
		t.Fatal("wrong native session acknowledged after drift")
	}
	var draft struct {
		Inputs []struct {
			Locator string `json:"locator"`
		} `json:"inputs"`
	}
	if err := json.Unmarshal(f.R.HandoffDraft, &draft); err != nil || len(draft.Inputs) == 0 {
		t.Fatalf("fixture has no frozen input: %v", err)
	}
	input := draft.Inputs[0].Locator
	original := mustReadAL01(t, input)
	info, err := os.Stat(input)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(input, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(input, append(original, []byte("\nchanged")...), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err == nil || !strings.Contains(err.Error(), "frozen Task input changed") {
		t.Fatalf("changed mandate input was accepted after drift: %v", err)
	}
	if err = os.WriteFile(input, original, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Chmod(input, info.Mode().Perm()); err != nil {
		t.Fatal(err)
	}
	if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err != nil {
		t.Fatalf("restored inputs did not let the same owner report: %v", err)
	}
}

func TestAL04NegativeAndStoppedHistoryIsNotRevived(t *testing.T) {
	t.Run("negative_acceptance", func(t *testing.T) {
		f := workflowTestFixture(t)
		o := workflowTestStart(t, f)
		a := f.Acceptance
		a.Acceptance = "unknown"
		a.RuntimeClaim = RuntimeClaim{NativeSessionID: a.RuntimeClaim.NativeSessionID}
		a.Sandbox = json.RawMessage(`null`)
		a.Issues = json.RawMessage(`[{"type":"unknown","detail":"Policy unknown"}]`)
		if _, err := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "negative.json", a)); err != nil {
			t.Fatal(err)
		}
		al04Install(t, f.R.Runtime.PlyExecutable.Path, []byte("updated ply"))
		al04Install(t, f.R.Runtime.Executable.Path, []byte("#!/bin/sh\necho updated provider\n"))
		if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "after-negative.json", f.Report)); err == nil {
			t.Fatal("drift revived negative acceptance")
		}
	})
	t.Run("pre_start_stop", func(t *testing.T) {
		f, o, ply := al04ReceiptOwnerBeforeCheckIn(t)
		stop := al01StopReport(f)
		if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "stop.json", stop)); err != nil {
			t.Fatal(err)
		}
		al04Install(t, ply, []byte("updated ply"))
		if _, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true); err == nil {
			t.Fatal("drift cleared a preserved pre-start stop")
		}
		shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
		if err != nil || shown.Round.State != "stopped" || shown.FinalReturn.TerminalSHA256 != nil {
			t.Fatalf("stopped history is not truthfully observable: %+v %v", shown.Round, err)
		}
	})
}

func al04ReceiptOwnerBeforeCheckIn(t *testing.T) (workflowFixture, WorkflowRun, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ply := filepath.Join(root, "installed-ply")
	if err = os.WriteFile(ply, []byte("launch ply"), 0700); err != nil {
		t.Fatal(err)
	}
	f := al01ReceiptFixture(t, ply)
	return f, workflowTestStart(t, f), ply
}

// A prompt whose reply was lost keeps its unknown state and is never resent;
// the owner that did receive it still checks in and reports after an update.
func TestAL04PartialPromptEffectAfterDrift(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ply := filepath.Join(root, "installed-ply")
	if err = os.WriteFile(ply, []byte("launch ply"), 0700); err != nil {
		t.Fatal(err)
	}
	f := al01ReceiptFixture(t, ply)
	workflowTestModel(t, f, map[string]any{"lost_command": "agent prompt"})
	p, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = WorkflowStart(f.D, f.File, *p.(WorkflowPreview).Confirmation); err == nil {
		t.Fatal("lost prompt reply was reported as delivered")
	}
	workflowTestModel(t, f, map[string]any{"lost_command": ""})
	id := workflowID(f.R)
	al04Install(t, ply, []byte("updated ply"))
	al04Install(t, f.R.Runtime.Executable.Path, []byte("#!/bin/sh\necho updated provider\n"))
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, id)
	if err != nil {
		t.Fatal(err)
	}
	if shown.Transport.State != "unknown" || shown.NextAction.Actor != "coordinator" || !strings.Contains(shown.NextAction.Message, "resend") {
		t.Fatalf("partial prompt effect is not truthfully observable: %+v %+v", shown.Transport, shown.NextAction)
	}
	o, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, id, shown.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: shown.Transport.AgentSessionID}, true)
	if err != nil {
		t.Fatalf("the owner that received the prompt could not check in after drift: %v", err)
	}
	o = workflowTestReport(t, f, o)
	if o.Round.ReportSHA256 == nil || workflowTestCalls(t, f, "agent prompt") != 1 || workflowTestCalls(t, f, "agent start") != 1 {
		t.Fatal("partial prompt effect was replayed or the report was lost")
	}
}

// Rewrite the saved trust binding in both the reservation index and state, so
// workflowRead accepts it and the callback reaches the trust check itself.
func al04TamperTrust(t *testing.T, f workflowFixture, id string, facts *workflowTrustFacts) {
	t.Helper()
	s, err := workflowRead(f.R.WorkspaceRoot, id)
	if err != nil {
		t.Fatal(err)
	}
	s.CodexTrust = facts
	for _, path := range []string{workflowIndex(f.R.WorkspaceRoot, id), filepath.Join(s.Result.Paths.RunRoot, "state.json")} {
		b, err := Canonical(s)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.Chmod(path, 0600); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, b, 0600); err != nil {
			t.Fatal(err)
		}
	}
}

// Historical drift is only a changed observation of a confirmed grant. A saved
// binding that contradicts the request is an integrity error, never a warning.
func TestAL04InvalidTrustBindingIsNotHistoricalDrift(t *testing.T) {
	t.Run("binding_without_grant", func(t *testing.T) {
		f, o, _ := al04ReceiptOwner(t)
		al04TamperTrust(t, f, o.RunID, &workflowTrustFacts{RepositoryRoot: f.R.WorkspaceRoot, Paths: []workflowTrustPath{}})
		report := f.Report
		report.Round, report.ControlID, report.PreviousReportSHA256 = o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("contradictory trust binding panicked: %v", r)
				}
			}()
			_, err = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "al04-trust.json", report))
		}()
		if err == nil || !strings.Contains(err.Error(), "unexpected saved Codex trust binding") {
			t.Fatalf("contradictory trust binding was treated as drift: %v", err)
		}
		shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
		if err != nil || shown.Round.ReportSHA256 != nil || shown.NextAction.Actor != "coordinator" || !strings.Contains(fmt.Sprint(shown.Reasons), "unexpected saved Codex trust binding") {
			t.Fatalf("integrity stop is not truthfully observable: %+v %v", shown, err)
		}
	})
	t.Run("grant_without_binding", func(t *testing.T) {
		f := workflowTestFixture(t)
		workflowGrant(t, &f)
		o := workflowTestStart(t, f)
		f.Acceptance.RequestSHA256, f.Report.RequestSHA256 = o.RequestSHA256, o.RequestSHA256
		o = workflowTestAccept(t, f, o)
		al04TamperTrust(t, f, o.RunID, nil)
		r := f.Report
		r.Round, r.ControlID, r.PreviousReportSHA256 = o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256
		if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "al04-missing.json", r)); err == nil {
			t.Fatal("a grant without its saved confirmation was treated as historical drift")
		}
	})
}

func al04CLI(t *testing.T, binary, cwd string, args ...string) (json.RawMessage, string) {
	t.Helper()
	c := exec.Command(binary, append([]string{"workflow"}, append(args, "--format", "json")...)...)
	c.Dir = cwd
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	if err := c.Run(); err != nil {
		t.Fatalf("CLI %v: %v\nstdout: %s\nstderr: %s", args, err, stdout.String(), stderr.String())
	}
	// Exactly one JSON value; a log line before or after it corrupts the result.
	var value json.RawMessage
	dec := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	if err := dec.Decode(&value); err != nil {
		t.Fatalf("CLI %v stdout is not JSON: %v\n%s", args, err, stdout.String())
	}
	if dec.More() {
		t.Fatalf("CLI %v stdout has more than one value:\n%s", args, stdout.String())
	}
	return value, stderr.String()
}

// The installed Ply is the actual AL03 base build; the owner checks in with it,
// an ordinary install puts the candidate at the same path, and the same owner
// completes its round through the candidate CLI. Provider bytes also change.
func TestAL04ReportCLIJourneyAcrossInstall(t *testing.T) {
	candidate := automaticCLIBinary(t)
	base := continuitySourceBinary(t, "f5535ef88b8bf2e2b75720285234c2fa71070a34", "PLY_AL04_BASE_BINARY")
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	installed := filepath.Join(root, "ply")
	al04Install(t, installed, mustReadAL01(t, base))
	f := al01ReceiptFixture(t, installed)
	o := workflowTestStart(t, f)
	cwd, _ := f.D.CWD()
	al04CLI(t, installed, cwd, "run", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID)
	before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	al04Install(t, installed, mustReadAL01(t, candidate))
	al04Install(t, f.R.Runtime.Executable.Path, []byte("#!/bin/sh\necho updated provider\n"))
	r := f.Report
	r.Round, r.ControlID, r.PreviousReportSHA256 = o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256
	file := writeAny(t, f.R.WorkspaceRoot, "al04-cli-round.json", r)
	for i := 0; i < 2; i++ {
		raw, stderr := al04CLI(t, installed, cwd, "run", "report", o.RunID, "--context", o.Paths.Context, "--file", file)
		var got WorkflowRun
		if err = json.Unmarshal(raw, &got); err != nil || got.Round.ReportSHA256 == nil {
			t.Fatalf("report was not recorded: %v %s", err, raw)
		}
		for _, s := range []string{"level=warning", "ply_executable", "provider_executable", "run_id=" + o.RunID} {
			if !strings.Contains(stderr, s) {
				t.Fatalf("report warn lacks %q: %s", s, stderr)
			}
		}
	}
	s := al04AssertSameOwner(t, f, before)
	if len(s.Records) != 1 || workflowTestCalls(t, f, "agent prompt") != 1 {
		t.Fatal("repeated CLI report produced a second record or dispatch")
	}
	for _, surface := range [][]string{{"run", "show", o.RunID}, {"execute", "show", o.RunID}} {
		raw, stderr := al04CLI(t, installed, cwd, surface...)
		var shown WorkflowRun
		if err = json.Unmarshal(raw, &shown); err != nil || shown.RunID != o.RunID {
			t.Fatalf("%v readback differs: %v %s", surface, err, raw)
		}
		if strings.Contains(string(raw), "workflow_run_drift") || !strings.Contains(stderr, "level=warning") {
			t.Fatalf("%v did not keep drift as a warning on stderr: %s\n%s", surface, raw, stderr)
		}
	}
}
