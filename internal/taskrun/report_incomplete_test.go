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

// al10Observation is an actual process run of the declared verifier.
type al10Observation struct {
	Exit           int
	Stdout, Stderr string
}

// al10Observe runs the frozen declared verifier argv in its declared cwd.
func al10Observe(t *testing.T, f workflowFixture) al10Observation {
	t.Helper()
	var verifiers []struct {
		Argv []string `json:"argv"`
		CWD  string   `json:"cwd"`
	}
	json.Unmarshal(f.Report.VerifierResults, &verifiers)
	var stdout, stderr bytes.Buffer
	c := exec.Command(verifiers[0].Argv[0], verifiers[0].Argv[1:]...)
	c.Dir, c.Stdout, c.Stderr = verifiers[0].CWD, &stdout, &stderr
	err := c.Run()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() == 0 {
		t.Fatalf("declared verifier did not fail as an actual process: %v", err)
	}
	return al10Observation{exit.ExitCode(), stdout.String(), stderr.String()}
}

// al10Report returns a worked round report whose single frozen requirement has
// the given outcome. Without an observation the verifier exit is synthetic
// unit data; with one, it is the actual process result and output.
func al10Report(t *testing.T, f workflowFixture, o WorkflowRun, name, outcome, requirement string, mutate func(map[string]any), observed *al10Observation) string {
	t.Helper()
	r := f.Report
	r.Round, r.ControlID, r.PreviousReportSHA256 = o.Round.Number, o.Round.ControlID, o.Round.PreviousReportSHA256
	r.Outcome = outcome
	usage := *r.BudgetUsage
	usage.CorrectionRounds = o.Budget.Floor.CorrectionRounds
	r.BudgetUsage = &usage
	var artifacts []map[string]any
	json.Unmarshal(r.Artifacts, &artifacts)
	old := artifacts[0]["locator"].(string)
	var evidence map[string]any
	b, _ := os.ReadFile(old)
	json.Unmarshal(b, &evidence)
	entry := evidence["requirements"].([]any)[0].(map[string]any)
	entry["outcome"] = requirement
	if requirement != "passed" {
		entry["reason"] = "The frozen verifier failed on the reported candidate; behavior is not finished."
	}
	if outcome != "complete" {
		var verifiers []map[string]any
		json.Unmarshal(r.VerifierResults, &verifiers)
		verifiers[0]["exit"] = 1
		if observed != nil {
			verifiers[0]["exit"] = observed.Exit
		}
		r.VerifierResults, _ = Canonical(verifiers)
		r.TechnicalAssessment.Gate = "failed"
		r.StopReasons = json.RawMessage(`[{"type":"unknown_or_partial_effect","primary":true,"detail":"Requirement f-01 is not finished; the frozen verifier failed."}]`)
	}
	if observed != nil {
		ids, outputs := []any{}, []map[string]any{}
		for _, out := range []struct{ id, text string }{{"observed-verifier-stderr", observed.Stderr}, {"observed-verifier-stdout", observed.Stdout}} {
			path := filepath.Join(filepath.Dir(old), name+"-"+out.id+".txt")
			if err := os.WriteFile(path, []byte(out.text), 0600); err != nil {
				t.Fatal(err)
			}
			outputs = append(outputs, map[string]any{"artifact_id": out.id, "kind": "managed", "description": "Actual output of the declared verifier process", "media_type": "text/plain", "classification": "workspace_internal", "size_bytes": len(out.text), "sha256": hash([]byte(out.text)), "locator": path})
			ids = append(ids, out.id)
		}
		artifacts = append(outputs, artifacts...)
		entry["artifact_ids"] = ids
		entry["reason"] = fmt.Sprintf("The declared verifier process exited %d on the reported candidate; behavior is not finished.", observed.Exit)
	}
	if mutate != nil {
		mutate(evidence)
	}
	path := filepath.Join(filepath.Dir(old), name+"-requirements.json")
	b, _ = Canonical(evidence)
	if err := os.WriteFile(path, b, 0600); err != nil {
		t.Fatal(err)
	}
	requirements := artifacts[len(artifacts)-1]
	requirements["locator"], requirements["sha256"], requirements["size_bytes"] = path, hash(b), len(b)
	r.Artifacts, _ = Canonical(artifacts)
	return writeAny(t, f.R.WorkspaceRoot, name+".json", r)
}

func al10CheckedIn(t *testing.T, f workflowFixture) WorkflowRun {
	t.Helper()
	o := workflowTestStart(t, f)
	o, err := WorkflowDeliveryReceipt(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, DeliveryBeginInput{Provider: "codex", NativeSession: o.Transport.AgentSessionID}, true)
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func TestAL10IncompleteReportAfterReceiptIsPreservedWithoutQualification(t *testing.T) {
	f := al01ReceiptFixture(t, "")
	o := al10CheckedIn(t, f)
	// The same unfinished requirement can never back a complete report.
	if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, al10Report(t, f, o, "complete-not-run", "complete", "not_run", nil, nil)); err == nil || !strings.Contains(err.Error(), "requirement f-01 has outcome not_run") {
		t.Fatalf("complete report with unfinished requirement: %v", err)
	}
	// Exact Spec/candidate binding, ID set and verifier references still bind.
	for name, mutate := range map[string]func(map[string]any){
		"wrong-oid": func(e map[string]any) { e["result_oid"] = strings.Repeat("0", 40) },
		"wrong-ids": func(e map[string]any) { e["requirements"].([]any)[0].(map[string]any)["id"] = "f-99" },
		"wrong-verifier": func(e map[string]any) {
			e["requirements"].([]any)[0].(map[string]any)["verifier_ids"] = []string{"absent"}
		},
	} {
		if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, al10Report(t, f, o, name, "blocked", "not_run", mutate, nil)); err == nil {
			t.Fatalf("%s: incomplete report escaped exact binding", name)
		}
	}
	// A small dirty, partially completed authorized change is reported as is.
	worktree := mustReportWorktree(t, f, o)
	if err := os.WriteFile(filepath.Join(worktree, "partial.txt"), []byte("unfinished\n"), 0600); err != nil {
		t.Fatal(err)
	}
	report := al10Report(t, f, o, "blocked-not-run", "blocked", "not_run", nil, nil)
	got, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, report)
	if err != nil {
		t.Fatalf("honest incomplete report after receipt was refused: %v", err)
	}
	if got.Round.State != "report_received" || got.FinalReturn.TerminalSHA256 != nil || got.TaskResultState != "not_published" {
		t.Fatalf("incomplete report gained qualification: %+v", got)
	}
	if _, err = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, report); err != nil {
		t.Fatalf("identical retry: %v", err)
	}
	if _, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, al10Report(t, f, o, "blocked-unknown", "blocked", "unknown", nil, nil)); err == nil {
		t.Fatal("preserved report was replaced")
	}
	o = workflowTestSettle(t, f, o)
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "blocked"))
	if err != nil {
		t.Fatal(err)
	}
	if o.FinalReturn.State != "blocked" || o.FinalReturn.TerminalSHA256 != nil || o.TaskResultState != "not_published" || o.Delivery != nil {
		t.Fatalf("blocked review gained qualification: %+v", o)
	}
	if strings.Contains(o.NextAction.Message, "human QA and integration") {
		t.Fatalf("report-only guidance invents unselected gates: %q", o.NextAction.Message)
	}
}

func TestAL10IncompleteReportCorrectsToCompleteAcceptedReturn(t *testing.T) {
	f := al01ReceiptFixture(t, "")
	o := al10CheckedIn(t, f)
	o, err := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, al10Report(t, f, o, "round-0", "blocked", "failed", nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	o = workflowTestSettle(t, f, o)
	if o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "changes_requested")); err != nil {
		t.Fatal(err)
	}
	// The documented example: protocol-permitted unknown facts stay unknown in
	// process_observation; evidence_gaps holds only actually missing evidence.
	var complete WorkflowReport
	b, _ := os.ReadFile(al10Report(t, f, o, "round-1", "complete", "passed", nil, nil))
	json.Unmarshal(b, &complete)
	note := "Model and effective policy were not observable and remain unknown from check-in; active_seconds is wall-clock time measured by the recipient."
	complete.ProcessObservation = &note
	if o, err = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "round-1-observed.json", complete)); err != nil {
		t.Fatal(err)
	}
	o = workflowTestSettle(t, f, o)
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if err != nil || o.FinalReturn.State != "accepted" || o.FinalReturn.TerminalSHA256 == nil || o.TaskResultState != "not_published" {
		t.Fatalf("corrected complete report: %+v %v", o, err)
	}
}

func TestAL10IncompleteReportCLIJourney(t *testing.T) {
	binary := automaticCLIBinary(t)
	f := al01ReceiptFixture(t, binary)
	o := workflowTestStart(t, f)
	if !strings.Contains(o.NextAction.Message, "check-in") || strings.Contains(o.NextAction.Message, "Accept the bound runtime") {
		t.Fatalf("report receipt advised legacy acceptance before check-in: %q", o.NextAction.Message)
	}
	o = al01CLIOK(t, binary, f, "run", "check-in", o.RunID, "--context", o.Paths.Context, "--provider", "codex", "--native-session", o.Transport.AgentSessionID)
	if o.NextAction.Message == "" || strings.Contains(o.NextAction.Message, "Accept the bound runtime") {
		t.Fatalf("received mandate kept legacy acceptance guidance: %q", o.NextAction.Message)
	}
	observed := al10Observe(t, f)
	report := al10Report(t, f, o, "cli-blocked", "blocked", "failed", nil, &observed)
	complete := al10Report(t, f, o, "cli-complete", "complete", "not_run", nil, nil)
	c := exec.Command(binary, "workflow", "run", "report", o.RunID, "--context", o.Paths.Context, "--file", complete, "--format", "json")
	c.Dir, _ = f.D.CWD()
	if out, err := c.CombinedOutput(); err == nil || !strings.Contains(string(out), "requirement f-01 has outcome not_run") {
		t.Fatalf("candidate CLI accepted a complete report with an unfinished requirement: %v %s", err, out)
	}
	for i := 0; i < 2; i++ {
		var stdout, stderr bytes.Buffer
		c := exec.Command(binary, "workflow", "run", "report", o.RunID, "--context", o.Paths.Context, "--file", report, "--format", "json")
		c.Dir, _ = f.D.CWD()
		c.Stdout, c.Stderr = &stdout, &stderr
		if err := c.Run(); err != nil {
			t.Fatalf("candidate CLI refused an honest blocked report: %v %s", err, stderr.Bytes())
		}
		if err := json.Unmarshal(stdout.Bytes(), &o); err != nil {
			t.Fatalf("stdout is not unchanged JSON: %v %s", err, stdout.Bytes())
		}
		warned := strings.Contains(stderr.String(), "Non-complete round report preserved") && strings.Contains(stderr.String(), o.RunID) && strings.Contains(stderr.String(), report)
		if warned != (i == 0) {
			t.Fatalf("attempt %d: sourced warning %v: %s", i, warned, stderr.Bytes())
		}
	}
	if o.Round.State != "report_received" || o.FinalReturn.TerminalSHA256 != nil || o.TaskResultState != "not_published" {
		t.Fatalf("candidate CLI qualified an incomplete report: %+v", o)
	}
	t.Logf("Candidate CLI received the mandate and preserved an honest blocked report: requirement failed, declared verifier process exited %d", observed.Exit)
}

func TestAL10ProviderStopAfterReceiptIsNotCalledPreStart(t *testing.T) {
	f := al01ReceiptFixture(t, "")
	o := al10CheckedIn(t, f)
	evidence := writeAny(t, f.R.WorkspaceRoot, "denied-report.json", map[string]any{"source": "synthetic runtime", "action": "report callback", "result": "denied after receipt"})
	in := DeliveryStopInput{EventID: "denied-report", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "report", Observer: "fixture coordinator", Reason: "Synthetic provider refused report transport", Action: "report callback", EffectStatus: DeliveryEffectUnknown, EvidencePath: evidence}
	o, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil {
		t.Fatal(err)
	}
	if shown, _ := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID); strings.Contains(shown.NextAction.Message, "pre-start") || strings.Contains(shown.NextAction.Message, "without acceptance") || !strings.Contains(shown.NextAction.Message, "Post-start provider stop preserved") {
		t.Fatalf("stop after receipt is mislabeled: %q", shown.NextAction.Message)
	}
	if _, err = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, al10Report(t, f, o, "after-stop", "blocked", "not_run", nil, nil)); err == nil || strings.Contains(err.Error(), "pre-start") {
		t.Fatalf("preserved stop after receipt: %v", err)
	}
}

func TestAL10ProviderStopAfterLegacyAcceptanceNamesRecordedStart(t *testing.T) {
	f := workflowTestFixture(t)
	if f.R.StartupProtocol == MandateReceiptV1 {
		t.Fatal("fixture is not a frozen legacy acceptance contract")
	}
	o := workflowTestAccept(t, f, workflowTestStart(t, f))
	evidence := writeAny(t, f.R.WorkspaceRoot, "denied-legacy-report.json", map[string]any{"source": "synthetic runtime", "action": "report callback", "result": "denied after acceptance"})
	in := DeliveryStopInput{EventID: "denied-legacy-report", Provider: "codex", NativeSession: o.Transport.AgentSessionID, Stage: "report", Observer: "fixture coordinator", Reason: "Synthetic provider refused report transport", Action: "report callback", EffectStatus: DeliveryEffectUnknown, EvidencePath: evidence}
	if _, err := WorkflowObserveDeliveryStop(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
		t.Fatal(err)
	}
	shown, _ := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	for _, wrong := range []string{"mandate receipt", "pre-start", "without acceptance"} {
		if strings.Contains(shown.NextAction.Message, wrong) {
			t.Fatalf("legacy stop after positive acceptance says %q: %q", wrong, shown.NextAction.Message)
		}
	}
	if !strings.Contains(shown.NextAction.Message, "Post-start provider stop preserved") {
		t.Fatalf("stop does not name the recorded native start: %q", shown.NextAction.Message)
	}
}
