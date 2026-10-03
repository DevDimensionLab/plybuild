package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

func workflowTestFixture(t *testing.T) workflowFixture {
	t.Helper()
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	f := newWorkflowFixture(t, root, "")
	bin := filepath.Join(root, "bin")
	if e = os.Mkdir(bin, 0700); e != nil {
		t.Fatal(e)
	}
	if e = os.Symlink(f.R.Runtime.Executable.Path, filepath.Join(bin, "codex")); e != nil {
		t.Fatal(e)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("HERDR_ENV", "1")
	return f
}
func workflowTestModel(t *testing.T, f workflowFixture, changes map[string]any) {
	t.Helper()
	b, e := os.ReadFile(f.Model)
	if e != nil {
		t.Fatal(e)
	}
	var m map[string]any
	if e = json.Unmarshal(b, &m); e != nil {
		t.Fatal(e)
	}
	for k, v := range changes {
		m[k] = v
	}
	writeAny(t, f.R.WorkspaceRoot, filepath.Base(f.Model), m)
}
func workflowTestCalls(t *testing.T, f workflowFixture, command string) int {
	t.Helper()
	b, e := os.ReadFile(f.Calls)
	if os.IsNotExist(e) {
		return 0
	}
	if e != nil {
		t.Fatal(e)
	}
	n := 0
	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var x struct {
			Argv []string `json:"argv"`
		}
		if e = json.Unmarshal(line, &x); e != nil {
			t.Fatal(e)
		}
		if strings.Join(x.Argv[:2], " ") == command {
			n++
		}
	}
	return n
}
func workflowTestStart(t *testing.T, f workflowFixture) WorkflowRun {
	t.Helper()
	p, e := WorkflowPreviewStart(f.D, f.File)
	if e != nil {
		t.Fatal(e)
	}
	preview := p.(WorkflowPreview)
	o, e := WorkflowStart(f.D, f.File, *preview.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	return o
}
func workflowTestAccept(t *testing.T, f workflowFixture, o WorkflowRun) WorkflowRun {
	t.Helper()
	p := writeAny(t, f.R.WorkspaceRoot, "acceptance.json", f.Acceptance)
	o, e := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, p)
	if e != nil {
		t.Fatal(e)
	}
	return o
}
func workflowTestReport(t *testing.T, f workflowFixture, o WorkflowRun) WorkflowRun {
	t.Helper()
	r := f.Report
	r.Round = o.Round.Number
	r.ControlID = o.Round.ControlID
	r.PreviousReportSHA256 = o.Round.PreviousReportSHA256
	if r.BudgetUsage.CorrectionRounds < o.Budget.Floor.CorrectionRounds {
		u := *r.BudgetUsage
		u.CorrectionRounds = o.Budget.Floor.CorrectionRounds
		r.BudgetUsage = &u
	}
	p := writeAny(t, f.R.WorkspaceRoot, fmt.Sprintf("round-%d.json", r.Round), r)
	o, e := WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, p)
	if e != nil {
		t.Fatal(e)
	}
	return o
}
func workflowTestSettle(t *testing.T, f workflowFixture, o WorkflowRun) WorkflowRun {
	t.Helper()
	workflowTestModel(t, f, map[string]any{"agent_status": "idle"})
	o, e := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 3*time.Second)
	if e != nil {
		t.Fatal(e)
	}
	return o
}
func workflowTestReady(t *testing.T, f workflowFixture) WorkflowRun {
	t.Helper()
	o := workflowTestAccept(t, f, workflowTestStart(t, f))
	return workflowTestSettle(t, f, workflowTestReport(t, f, o))
}
func workflowTestReview(t *testing.T, f workflowFixture, o WorkflowRun, decision string) string {
	t.Helper()
	findings := []WorkflowFinding{}
	if decision != "accepted" {
		findings = append(findings, WorkflowFinding{"A fixture check needs correction.", "The bound behavior passes.", "Run the frozen verifier."})
	}
	r := WorkflowReview{Envelope: workflowEnv("run-review"), ReviewID: fmt.Sprintf("review-%d", o.Round.Number), RunID: o.RunID, RequestSHA256: o.RequestSHA256, HandoffSHA256: o.Handoff.SHA256, Round: o.Round.Number, ReportSHA256: *o.Round.ReportSHA256, Reviewer: f.R.Coordinator.ActorClaim, Decision: decision, Findings: findings}
	return writeAny(t, f.R.WorkspaceRoot, fmt.Sprintf("review-%d.json", o.Round.Number), r)
}

func TestWorkflowConcurrentStartsAndLegacyReservation(t *testing.T) {
	f := workflowTestFixture(t)
	p, e := WorkflowPreviewStart(f.D, f.File)
	if e != nil {
		t.Fatal(e)
	}
	confirmation := *p.(WorkflowPreview).Confirmation
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, e := WorkflowStart(f.D, f.File, confirmation); errs <- e }()
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	for _, cmd := range []string{"tab create", "agent start", "agent prompt"} {
		if n := workflowTestCalls(t, f, cmd); n != 1 {
			t.Fatalf("%s: %d", cmd, n)
		}
	}
	r := workflowNativeRequest(f.R)
	r.RequestKey = "another/native"
	if _, e = PreviewStart(f.D, writeAny(t, f.R.WorkspaceRoot, "native-other.json", r)); e == nil {
		t.Fatal("legacy start ignored workflow target reservation")
	}
	other := f.R
	other.RequestKey = "another/workflow"
	if _, e = WorkflowPreviewStart(f.D, writeAny(t, f.R.WorkspaceRoot, "other.json", other)); e == nil {
		t.Fatal("second workflow ignored target reservation")
	}
}
func TestWorkflowLegacyReservationBlocksHerdr(t *testing.T) {
	f := workflowTestFixture(t)
	r := workflowNativeRequest(f.R)
	file := writeAny(t, r.WorkspaceRoot, "native.json", r)
	p, e := PreviewStart(f.D, file)
	if e != nil {
		t.Fatal(e)
	}
	d := f.D
	d.Fault = func(point string) error {
		if point == "after_request_reservation" {
			return errors.New("crash")
		}
		return nil
	}
	_, _ = Start(d, file, *p.(Preview).Confirmation)
	if _, e = WorkflowPreviewStart(f.D, f.File); e == nil {
		t.Fatal("native unknown reservation did not block Herdr")
	}
	if workflowTestCalls(t, f, "tab create") != 0 {
		t.Fatal("unexpected transport effect")
	}
}
func TestWorkflowConcurrentReviewsReserveAndSendOnce(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestReady(t, f)
	file := workflowTestReview(t, f, o, "changes_requested")
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, _ = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file) }()
	}
	wg.Wait()
	after, e := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if after.Round.Number != 1 || after.Budget.Floor.CorrectionRounds != 1 || workflowTestCalls(t, f, "agent prompt") != 2 {
		t.Fatalf("duplicate correction: %+v", after)
	}
}
func TestWorkflowSubprocessTimeoutAndLostStartDoNotReplay(t *testing.T) {
	for _, command := range []string{"tab create", "agent start", "agent prompt"} {
		t.Run(command, func(t *testing.T) {
			f := workflowTestFixture(t)
			workflowTestModel(t, f, map[string]any{"delay_command": command, "delay_seconds": 2})
			d := f.D
			d.HerdrTimeout = 500 * time.Millisecond
			p, e := WorkflowPreviewStart(d, f.File)
			if e != nil {
				t.Fatal(e)
			}
			o, e := WorkflowStart(d, f.File, *p.(WorkflowPreview).Confirmation)
			if e == nil || o.RunID == "" {
				t.Fatal("timeout lost its run binding")
			}
			if _, e = WorkflowStart(d, f.File, *p.(WorkflowPreview).Confirmation); e != nil {
				t.Fatal(e)
			}
			if n := workflowTestCalls(t, f, command); n != 1 {
				t.Fatalf("replayed %s: %d", command, n)
			}
		})
	}
}
func TestWorkflowCrashReservationsCannotRestart(t *testing.T) {
	for _, point := range []string{"workflow_after_reservation", "workflow_before_tab_send", "workflow_before_agent_send", "workflow_before_prompt_send"} {
		t.Run(point, func(t *testing.T) {
			f := workflowTestFixture(t)
			p, e := WorkflowPreviewStart(f.D, f.File)
			if e != nil {
				t.Fatal(e)
			}
			d := f.D
			d.Fault = func(p string) error {
				if p == point {
					return errors.New("simulated crash")
				}
				return nil
			}
			_, _ = WorkflowStart(d, f.File, *p.(WorkflowPreview).Confirmation)
			before := workflowTestCalls(t, f, "tab create") + workflowTestCalls(t, f, "agent start") + workflowTestCalls(t, f, "agent prompt")
			if _, e = WorkflowStart(f.D, f.File, *p.(WorkflowPreview).Confirmation); e != nil {
				t.Fatal(e)
			}
			after := workflowTestCalls(t, f, "tab create") + workflowTestCalls(t, f, "agent start") + workflowTestCalls(t, f, "agent prompt")
			if before != after {
				t.Fatal("crash reservation replayed an effect")
			}
		})
	}
}
func TestWorkflowCorrectionCrashKeepsBudgetAndNoResend(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestReady(t, f)
	file := workflowTestReview(t, f, o, "changes_requested")
	d := f.D
	d.Fault = func(point string) error {
		if point == "workflow_after_correction_reservation" {
			return errors.New("crash before send")
		}
		return nil
	}
	_, _ = WorkflowReviewRun(d, f.R.WorkspaceRoot, o.RunID, file)
	after, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file)
	if e == nil || after.Round.Number != 1 || after.Budget.Floor.CorrectionRounds != 1 {
		t.Fatalf("reservation lost: %+v %v", after, e)
	}
	if workflowTestCalls(t, f, "agent prompt") != 1 {
		t.Fatal("pre-send crash was retried")
	}
}
func TestWorkflowTerminalRecoveryAndReservationAfterAcceptance(t *testing.T) {
	for _, point := range []string{"workflow_before_terminal_publish", "workflow_after_terminal_publish"} {
		t.Run(point, func(t *testing.T) {
			f := workflowTestFixture(t)
			o := workflowTestReady(t, f)
			file := workflowTestReview(t, f, o, "accepted")
			d := f.D
			d.Fault = func(p string) error {
				if p == point {
					return errors.New("publication interrupted")
				}
				return nil
			}
			if _, e := WorkflowReviewRun(d, f.R.WorkspaceRoot, o.RunID, file); e == nil {
				t.Fatal("fault did not interrupt")
			}
			after, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file)
			if e != nil {
				t.Fatal(e)
			}
			if after.FinalReturn.State != "accepted" || after.FinalReturn.TerminalSHA256 == nil {
				t.Fatal("terminal recovery failed")
			}
			calls := workflowTestCalls(t, f, "agent get")
			again, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file)
			if e != nil || !equal(after.FinalReturn, again.FinalReturn) || calls != workflowTestCalls(t, f, "agent get") {
				t.Fatal("sealed duplicate contacted Herdr")
			}
			r := f.R
			r.RequestKey = "later/without-task-result"
			if _, e = WorkflowPreviewStart(f.D, writeAny(t, r.WorkspaceRoot, "later.json", r)); e == nil {
				t.Fatal("accepted report released target without native TaskResult")
			}
		})
	}
}
func TestWorkflowStrictJSONAndBoundedCounters(t *testing.T) {
	f := workflowTestFixture(t)
	good, _ := Canonical(f.R)
	for name, b := range map[string][]byte{"duplicate": append([]byte(`{"kind":"x",`), good[1:]...), "unknown": append([]byte(`{"unknown":true,`), good[1:]...), "utf8": []byte("{\"kind\":\"\xff\"}"), "oversize": bytes.Repeat([]byte(" "), (1<<20)+1), "float": bytes.Replace(good, []byte(`"schema_version":1`), []byte(`"schema_version":1.0`), 1), "bool": bytes.Replace(good, []byte(`"schema_version":1`), []byte(`"schema_version":true`), 1), "huge": bytes.Replace(good, []byte(`"schema_version":1`), []byte(`"schema_version":9007199254740993`), 1)} {
		t.Run(name, func(t *testing.T) {
			p := filepath.Join(f.R.WorkspaceRoot, "invalid.json")
			os.WriteFile(p, b, 0600)
			if _, e := ReadWorkflowRequest(p); e == nil {
				t.Fatal("invalid request accepted")
			}
		})
	}
	r := f.Report
	r.BudgetUsage = &BudgetUsage{true, 100, 1, 0}
	b, _ := Canonical(r)
	if _, e := workflowParseReport(b); e == nil {
		t.Fatal("large native counter accepted")
	}
	if workflowTestCalls(t, f, "tab create") != 0 {
		t.Fatal("invalid JSON had effects")
	}
}
func TestWorkflowAllBudgetLimitsAndMonotonicFloor(t *testing.T) {
	for _, u := range []BudgetUsage{{true, 3, 10, 0}, {true, 0, 5400, 0}, {true, 0, 10, 2}} {
		t.Run(fmt.Sprint(u), func(t *testing.T) {
			f := workflowTestFixture(t)
			f.Report.BudgetUsage = &u
			o := workflowTestReady(t, f)
			file := workflowTestReview(t, f, o, "changes_requested")
			if _, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file); e == nil {
				t.Fatal("exhausted A allowed correction")
			}
			if workflowTestCalls(t, f, "agent prompt") != 1 {
				t.Fatal("exhausted A sent input")
			}
		})
	}
	f := workflowTestFixture(t)
	f.Report.BudgetUsage = &BudgetUsage{true, 1, 100, 1}
	o := workflowTestReady(t, f)
	file := workflowTestReview(t, f, o, "changes_requested")
	o, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, file)
	if e != nil {
		t.Fatal(e)
	}
	r := f.Report
	r.Round = 1
	r.ControlID = o.Round.ControlID
	r.PreviousReportSHA256 = o.Round.PreviousReportSHA256
	r.BudgetUsage = &BudgetUsage{true, 1, 101, 1}
	if _, e = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "decreasing.json", r)); e == nil {
		t.Fatal("usage below reserved floor accepted")
	}
}
func TestWorkflowRuntimeDriftAndNegativeAuthority(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestStart(t, f)
	a := f.Acceptance
	a.RuntimeClaim.ModelID = ptr("another-model")
	file := writeAny(t, f.R.WorkspaceRoot, "model-mismatch.json", a)
	if _, e := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); e == nil {
		t.Fatal("known model mismatch accepted")
	}
	a = f.Acceptance
	a.Acceptance = "unknown"
	a.RuntimeClaim.RuntimeID = nil
	a.RuntimeClaim.ProfileID = nil
	a.RuntimeClaim.EffectivePolicySHA256 = nil
	a.Sandbox = json.RawMessage("null")
	a.Issues = json.RawMessage(`[{"type":"unknown","detail":"Necessary authority is unknown."}]`)
	o, e := WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "negative.json", a))
	if e != nil {
		t.Fatal(e)
	}
	facts, e := workflowhandoff.ReadTaskRunReturnFacts(f.D.Workflow, f.Handoff)
	if e != nil {
		t.Fatal(e)
	}
	if facts.Start != nil || o.FinalReturn.State != "blocked" {
		t.Fatal("negative authority manufactured a start")
	}
	if _, e = WorkflowReportRound(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "negative-report.json", f.Report)); e == nil {
		t.Fatal("negative acceptance reported success")
	}
}
