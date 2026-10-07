package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type traceNoWrites struct{ t *testing.T }

func (f traceNoWrites) WriteOnce(string, []byte) error { f.t.Fatal("trace wrote a file"); return nil }
func (f traceNoWrites) Replace(string, []byte) error   { f.t.Fatal("trace replaced a file"); return nil }

func traceRead(t *testing.T, d Dependencies, root string) TraceHistory {
	t.Helper()
	d.Executable = func() (string, error) { t.Fatal("trace inspected current executable"); return "", nil }
	d.Runner, d.ExecRunner = nil, nil
	d.Files = traceNoWrites{t}
	d.Workspace.WorkGit, d.Workspace.IntegrationGit, d.Workflow.Git = nil, nil, nil
	d.StartupProcessTree = func(int) ([]StartupProcess, error) { t.Fatal("trace probed processes"); return nil, nil }
	before := treeState(t, root)
	result, err := ReadTraceHistory(d, root, "task")
	if err != nil {
		t.Fatal(err)
	}
	if treeState(t, root) != before {
		t.Fatal("trace changed its workspace")
	}
	return result
}

func traceRun(t *testing.T, h TraceHistory, id string) TraceRun {
	t.Helper()
	for _, run := range h.Runs {
		if run.RunID == id {
			return run
		}
	}
	t.Fatalf("run %s absent: %+v", id, h)
	return TraceRun{}
}

func traceKind(run TraceRun, kind string) []TraceEntry {
	out := []TraceEntry{}
	for _, e := range run.Entries {
		if e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}

func traceReason(run TraceRun, code string) bool {
	for _, r := range run.Reasons {
		if r.Code == code {
			return true
		}
	}
	return false
}

func TestReadTraceEmptyAndUnknownTask(t *testing.T) {
	d, request, _ := fixture(t)
	out := traceRead(t, d, request.WorkspaceRoot)
	if len(out.Runs) != 0 || len(out.Reasons) != 0 {
		t.Fatalf("empty history: %+v", out)
	}
	if _, err := ReadTraceHistory(d, request.WorkspaceRoot, "unknown-task"); err == nil {
		t.Fatal("unknown Task became empty history")
	}
	for _, path := range []string{storeRoot(request.WorkspaceRoot), workflowRoot(request.WorkspaceRoot)} {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("empty read created %s", path)
		}
	}
}

func TestReadTraceLegacyRoundsAndIndependentDamagedSibling(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestReady(t, f)
	var err error
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "changes_requested"))
	if err != nil {
		t.Fatal(err)
	}
	o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if err != nil {
		t.Fatal(err)
	}
	// A safely bound legacy native reservation is another run, not an alias of
	// the workflow report/handoff. Its absent execution evidence stays explicit.
	native := workflowNativeRequest(f.R)
	native.RequestKey = "trace/native"
	if err = os.MkdirAll(runPaths(native).RunRoot, 0700); err != nil {
		t.Fatal(err)
	}
	if err = writeValue(filepath.Join(runPaths(native).RunRoot, "request.json"), native); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(workflowRoot(f.R.WorkspaceRoot), "requests", "wfr_"+strings.Repeat("b", 64)+".json")
	if err = os.WriteFile(bad, []byte("broken"), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	history := traceRead(t, f.D, f.R.WorkspaceRoot)
	if len(history.Runs) != 2 || len(history.Reasons) == 0 {
		t.Fatalf("sibling coverage: %+v", history)
	}
	run := traceRun(t, history, o.RunID)
	if run.Family != "workflow_run_v1" || len(traceKind(run, "round_report")) != 2 || len(traceKind(run, "round_review")) != 2 || len(run.DeclaredProcess) != 1 || !strings.Contains(string(run.Contract), `"agreement"`) {
		t.Fatalf("lost rounds or actual contract: %+v", run)
	}
	for i, e := range append(traceKind(run, "round_report"), traceKind(run, "round_review")...) {
		if e.OccurredAtUTC != nil || e.ReportedAtUTC != nil || e.RegisteredAtUTC != nil {
			t.Fatalf("invented time at %d: %+v", i, e)
		}
	}
	reviews := traceKind(run, "round_review")
	if *reviews[0].Outcome != "changes_requested" || *reviews[1].Outcome != "accepted" || reviews[0].Role != "unknown" {
		t.Fatalf("changed review semantics: %+v", reviews)
	}
	legacy := traceRun(t, history, RunID(native.RequestKey))
	if legacy.Family != "legacy_task_run" || len(traceKind(legacy, "run_request")) != 1 {
		t.Fatalf("legacy run: %+v", legacy)
	}
	// mtime cannot manufacture round time or change the projection.
	s, _ := workflowRead(f.R.WorkspaceRoot, o.RunID)
	for _, record := range s.Records {
		if err = os.Chtimes(record.Report.Locator, time.Now().Add(24*time.Hour), time.Now().Add(24*time.Hour)); err != nil {
			t.Fatal(err)
		}
	}
	if after := traceRead(t, f.D, f.R.WorkspaceRoot); !reflect.DeepEqual(history, after) {
		t.Fatal("mtime changed trace facts")
	}
	if err = os.WriteFile(s.Records[0].Report.Locator, []byte("damaged earlier round"), 0600); err != nil {
		t.Fatal(err)
	}
	after := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	if len(traceKind(after, "round_report")) != 1 || len(traceKind(after, "round_review")) != 2 || !traceReason(after, "trace_round_report_unavailable") {
		t.Fatalf("damaged first round hid valid later records: %+v", after)
	}
}

func TestReadTraceDeliveryLoopsExactBindingsAndNoProbes(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	publishOlder := traceLegacyFixture(t, f)
	// Use the production frozen goal binding instead of the old fixture's
	// intentionally permissive human design-document shortcut.
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "trace-delivery-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	var err error
	report := func(id, phase string) {
		t.Helper()
		r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: id, PreviousEventSHA256: o.Delivery.LastEventSHA256, Phase: phase, Summary: "Explicit fixture progress", Meaning: "A report, not a native result.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
		if phase == "needs_input" {
			r.Question = ptr("Which exact behavior should the fixture implement?")
		}
		o, err = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, id+".json", r))
		if err != nil {
			t.Fatal(err)
		}
	}
	deliveryTestCandidate(t, f, "first rejected behavior")
	if err = os.WriteFile(f.AcceptancePath, []byte("printf 'actual failed check\\n'\nexit 7\n"), 0600); err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "failed"))
	if err == nil {
		t.Fatal("fixture verifier should actually fail")
	}
	report("question", "needs_input")
	report("correction", "working")
	deliveryTestCandidate(t, f, "second behavior")
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "second"))
	if err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "fail", deliveryTestHuman(t, f, o, "fail"))
	if err != nil {
		t.Fatal(err)
	}
	deliveryTestCandidate(t, f, "third behavior")
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "third"))
	if err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil {
		t.Fatal(err)
	}
	oldWorkflowID, oldNativeID := publishOlder()
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	state.Result.Reasons = append(state.Result.Reasons, Reason{"herdr_bootstrap_response", "A historical startup response was unavailable."}, Reason{"herdr_start_stopped", "Historical startup stopped before a later successful observation."})
	if err = workflowSave(f.D, state); err != nil {
		t.Fatal(err)
	}
	// A later declared goal is independent of the execution's frozen text.
	goalDraftRaw, err := os.ReadFile(filepath.Join(f.R.WorkspaceRoot, "goal-draft.json"))
	if err != nil {
		t.Fatal(err)
	}
	var goalDraft map[string]any
	if err = json.Unmarshal(goalDraftRaw, &goalDraft); err != nil {
		t.Fatal(err)
	}
	goalDraft["publication_key"], goalDraft["expected_previous"] = "fixture/goal-revision-two", f.Prepared.Preparation.Preparation.Plan.Spec
	goalDraft["objective"], goalDraft["change_reason"] = "New current goal objective after execution.", "Exercise current versus frozen declarations."
	goalDraft["requirements"] = []any{map[string]any{"id": "behavior", "acceptance": "The new goal has different acceptance text."}}
	if _, err = workspace.RecordTaskSpec(f.D.Workspace, workspace.TaskContentInput{TaskID: "task", File: writeAny(t, f.R.WorkspaceRoot, "goal-revision-two.json", goalDraft)}); err != nil {
		t.Fatal(err)
	}
	// Deliberately remove binaries. Historical reading cannot call or qualify
	// Git, provider/Herdr, an acceptance program, or the installed CLI.
	for _, path := range []string{f.R.Runtime.Executable.Path, f.R.Runtime.PlyExecutable.Path, f.R.Herdr.Executable.Path} {
		if err = os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	h := traceRead(t, f.D, f.R.WorkspaceRoot)
	if len(h.Runs) != 3 {
		t.Fatalf("mixed-family history has %d runs: %+v", len(h.Runs), h.Reasons)
	}
	older := traceRun(t, h, oldWorkflowID)
	if older.Family != "workflow_run_v1" || len(traceKind(older, "round_report")) != 2 || len(traceKind(older, "round_review")) != 2 || !strings.Contains(string(older.Contract), `"name":"A"`) {
		t.Fatalf("older semantics lost: %+v", older)
	}
	if native := traceRun(t, h, oldNativeID); native.Family != "legacy_task_run" || len(traceKind(native, "run_request")) != 1 {
		t.Fatalf("native legacy history lost: %+v", native)
	}
	run := traceRun(t, h, o.RunID)
	if run.Coverage != "complete" {
		t.Fatalf("complete run became partial: %+v", run.Reasons)
	}
	if len(run.HistoricalReasons) != 2 || len(run.NativeRunIDs) != 3 {
		t.Fatalf("historical diagnostics/owner+candidate aliases lost: %+v / %+v", run.HistoricalReasons, run.NativeRunIDs)
	}
	if run.Family != "goal_execute_v2" || run.FrozenGoal == nil || *run.FrozenGoal != f.Prepared.Goal {
		t.Fatalf("frozen goal lost: %+v", run)
	}
	var frozen struct {
		Goal struct {
			Objective    string `json:"objective"`
			Requirements []struct {
				Acceptance string `json:"acceptance"`
			} `json:"requirements"`
		} `json:"frozen_goal_declaration"`
		Execution struct {
			ContractKind string `json:"contract_kind"`
		} `json:"frozen_execution_declaration"`
	}
	if json.Unmarshal(run.Contract, &frozen) != nil || frozen.Goal.Objective != "Observe a complete delivery without attributing synthetic evidence to a real human." || len(frozen.Goal.Requirements) != 1 || frozen.Goal.Requirements[0].Acceptance != "The actual acceptance entrypoint observes the fixture behavior." || frozen.Execution.ContractKind != "execution" {
		t.Fatalf("current goal replaced frozen declarations: %s", run.Contract)
	}
	verifications, results, qa, integration := traceKind(run, "verification"), traceKind(run, "task_result_technical"), traceKind(run, "human_qa"), traceKind(run, "integration")
	if len(verifications) != 3 || len(results) != 2 || len(qa) != 2 || len(integration) != 1 {
		t.Fatalf("attempt/candidate/QA counts=%d/%d/%d/%d; reasons=%+v", len(verifications), len(results), len(qa), len(integration), run.Reasons)
	}
	if *verifications[0].Verification.Exit != 7 || verifications[0].Candidate.TaskResultID != nil || *qa[0].Outcome != "fail" || *qa[1].Outcome != "pass" {
		t.Fatal("failed attempts or exact verdict history lost")
	}
	for _, e := range verifications {
		if !e.Verification.InputsBound || e.Verification.StartedAtUTC == nil || e.Verification.FinishedAtUTC == nil {
			t.Fatalf("exact verification inputs/time lost: %+v", e)
		}
	}
	for _, e := range append(append(results, qa...), integration...) {
		if e.NativeID == nil || e.NativeEventID == nil || e.Candidate == nil || e.Candidate.TaskResultID == nil {
			t.Fatalf("native identity missing: %+v", e)
		}
	}
	if *qa[0].Candidate.TaskResultID == *qa[1].Candidate.TaskResultID || *integration[0].Candidate.TaskResultID != *qa[1].Candidate.TaskResultID {
		t.Fatal("QA transferred between candidates")
	}
	reports := traceKind(run, "delivery_report")
	if len(reports) != 2 || reports[0].RegisteredAtUTC != nil || reports[0].OccurredAtUTC != nil || reports[0].ReportedAtUTC != nil || reports[1].RegisteredAtUTC != nil {
		t.Fatal("untimed reports gained timestamps")
	}
	var question DeliveryReport
	if json.Unmarshal(reports[0].Data, &question) != nil || question.Question == nil || *question.Question != "Which exact behavior should the fixture implement?" {
		t.Fatal("exact question lost")
	}
	if strings.Contains(string(traceJSON(h)), `"reply_capability"`) || strings.Contains(string(traceJSON(h)), `"secret"`) {
		t.Fatal("private reply capability escaped")
	}
	traceMixedCLI(t, f.R.WorkspaceRoot)
	// Corrupt one verifier source; native sibling QA/integration records and
	// the other verifier attempts remain visible with bounded degradation.
	if err = os.WriteFile(verifications[0].Source.Locator, []byte("broken verification"), 0600); err != nil {
		t.Fatal(err)
	}
	degraded := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	if len(traceKind(degraded, "verification")) != 2 || len(traceKind(degraded, "human_qa")) != 2 || len(traceKind(degraded, "integration")) != 1 || degraded.Coverage != "partial" {
		t.Fatalf("source damage hid valid siblings: %+v", degraded)
	}
	if err = os.Remove(f.R.Delivery.Goal.Locator); err != nil {
		t.Fatal(err)
	}
	missing := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	if missing.FrozenGoal != nil || missing.FrozenBasis == nil || len(traceKind(missing, "human_qa")) != 2 || !traceReason(missing, "trace_frozen_goal_unavailable") {
		t.Fatalf("missing declaration hid independent native facts: %+v", missing)
	}
}

func traceMixedCLI(t *testing.T, root string) {
	t.Helper()
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("test source location unavailable")
	}
	repo := filepath.Clean(filepath.Join(filepath.Dir(source), "..", ".."))
	binary := filepath.Join(t.TempDir(), "ply-trace-fixture")
	build := exec.Command("go", "build", "-o", binary, "./cmd/ply")
	build.Dir = repo
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build actual trace CLI: %s %v", output, err)
	}
	command := exec.Command(binary, "workflow", "trace", "task", "--format", "json")
	command.Dir = root
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	before := treeState(t, root)
	if err := command.Run(); err != nil {
		t.Fatalf("actual mixed trace: %v %s", err, stderr.String())
	}
	if treeState(t, root) != before {
		t.Fatal("actual CLI trace changed fixture source bytes")
	}
	var result struct {
		Runs []struct {
			ID     string `json:"id"`
			Family string `json:"family"`
		} `json:"runs"`
		Events []struct {
			ID         string   `json:"id"`
			Type       string   `json:"type"`
			RunIDs     []string `json:"run_ids"`
			Occurred   *string  `json:"occurred_at_utc"`
			Reported   *string  `json:"reported_at_utc"`
			Registered *string  `json:"registered_at_utc"`
		} `json:"events"`
		Analysis []struct {
			ID    string   `json:"id"`
			Value *float64 `json:"value"`
		} `json:"analysis"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("actual CLI did not emit one JSON object: %v", err)
	}
	runs, families := map[string]bool{}, map[string]bool{}
	for _, run := range result.Runs {
		runs[run.ID], families[run.Family] = true, true
	}
	if len(runs) != 3 || !families["legacy_task_run"] || !families["workflow_run_v1"] || !families["goal_execute_v2"] {
		t.Fatalf("actual projection lost mixed run families: %+v", result.Runs)
	}
	for _, event := range result.Events {
		for _, id := range event.RunIDs {
			if !runs[id] {
				t.Fatalf("handoff alias became unresolvable run: event=%s run=%s", event.ID, id)
			}
		}
		if event.Type == "delivery_report" || event.Type == "round_report" {
			if event.Occurred != nil || event.Reported != nil || event.Registered != nil {
				t.Fatalf("actual projection invented report clock: %+v", event)
			}
		}
	}
	want := map[string]float64{"recorded_runs": 3, "verifier_attempts": 3, "qualified_generations": 2, "distinct_observed_candidates": 3, "human_qa_records": 2}
	for _, a := range result.Analysis {
		if value, ok := want[a.ID]; ok {
			if a.Value == nil || *a.Value != value {
				t.Fatalf("actual projection %s=%v, expected %v", a.ID, a.Value, value)
			}
			delete(want, a.ID)
		}
	}
	if len(want) > 0 {
		t.Fatalf("actual projection omitted counts: %+v", want)
	}
}

// Build separate historical Handoff@2 contracts on the same registered Task.
// The legacy reports explicitly have unknown technical outcomes and no verifier
// exits. Only the delivery part of the mixed fixture executes verifiers and
// native QA/integration callbacks. Publishing historical bytes after that run
// avoids granting these synthetic old reservations present runtime ownership.
func traceLegacyFixture(t *testing.T, f deliveryFixture) func() (string, string) {
	t.Helper()
	p := f.Prepared.Preparation.Preparation.Plan
	var delivery map[string]json.RawMessage
	if err := json.Unmarshal(f.R.HandoffDraft, &delivery); err != nil {
		t.Fatal(err)
	}
	draft := minimalHandoffDraftValue(p.WorktreePath, "refs/heads/"+p.Branch, p.ParentOID)
	draft = replaceObjectMember(draft, "schema_version", int64(2))
	for _, name := range []string{"task_spec_binding", "inputs", "verifiers"} {
		v, err := canonicaljson.DecodeStrict(delivery[name])
		if err != nil {
			t.Fatal(err)
		}
		draft = replaceObjectMember(draft, name, v)
	}
	budget, _ := objectMember(draft, "budget")
	draft = replaceObjectMember(draft, "budget", replaceObjectMember(budget.(canonicaljson.Object), "max_rounds", int64(4)))
	raw, err := canonicaljson.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	raw = []byte(strings.ReplaceAll(string(raw), `"test"`, `"goal-acceptance"`))
	if _, err = workflowhandoff.ValidateTaskRunDraft(raw); err != nil {
		t.Fatal(err)
	}
	draftPath := writeAny(t, f.R.WorkspaceRoot, "trace-legacy-draft.json", json.RawMessage(raw))
	handoff, err := workflowhandoff.Create(f.D.Workflow, workflowhandoff.CreateInput{DraftPath: draftPath})
	if err != nil {
		t.Fatal(err)
	}
	link, err := workflowhandoff.ReadTaskRunLink(f.D.Workflow, handoff.Locator)
	if err != nil {
		t.Fatal(err)
	}
	old := f.R
	old.Envelope, old.RequestKey, old.HandoffDraft, old.Delivery = workflowEnv("herdr-run-request"), "trace/older-workflow", raw, nil
	old.ReturnMode, old.Agreement = "reviewed_report_only", Agreement{"A", 3, 5400, 2}
	old.Coordinator.ActorClaim, old.Coordinator.MayRequestChanges = "Historical coordinator claim", true
	old.Runtime.PermissionBinding.AuthorityKind = "reported_contract_with_effective_policy"
	native := workflowNativeRequest(old)
	native.RequestKey = "trace/older-native"
	if _, err = parseRequest(traceJSON(native)); err != nil {
		t.Fatal(err)
	}
	return func() (string, string) {
		t.Helper()
		x := workspace.PlanWorktreeObservation{StatusEntries: []workspace.IntegrationStatusEntry{}, InProgress: []string{}}
		s := workflowInitial(old, Observed{Target: x, Epic: x, RuntimeBindings: []FileBinding{}})
		s.Result.Handoff = WorkflowHandoff{link.HandoffID, link.Locator, link.SHA256}
		for _, path := range []string{filepath.Dir(workflowIndex(old.WorkspaceRoot, s.Result.RunID)), s.Result.Paths.RunRoot, runPaths(native).RunRoot} {
			if err = os.MkdirAll(path, 0700); err != nil {
				t.Fatal(err)
			}
		}
		if err = writeValue(workflowIndex(old.WorkspaceRoot, s.Result.RunID), s); err != nil {
			t.Fatal(err)
		}
		if err = writeValue(filepath.Join(runPaths(native).RunRoot, "request.json"), native); err != nil {
			t.Fatal(err)
		}
		var previous *string
		var control *string
		for round := 0; round < 2; round++ {
			dir := filepath.Join(s.Result.Paths.RunRoot, "rounds", workflowRoundName(round))
			report := WorkflowReport{Report: Report{Envelope: workflowEnv("round-report"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, Outcome: "unknown", Summary: "Preserved incomplete older round.", Meaning: "Unknown technical outcome; no verifier execution is asserted.", BudgetUsage: &BudgetUsage{InitialExecutionStarted: true, CorrectionRounds: round}, StopReasons: traceJSON([]any{map[string]any{"type": "unknown_or_partial_effect", "detail": "Preserved fixture has no verifier execution evidence.", "primary": true}}), ObservedEffects: json.RawMessage(`[]`), VerifierResults: json.RawMessage(`[]`), Review: json.RawMessage(`{"findings":[],"fixes":[],"open_actionable_findings":[]}`), Artifacts: json.RawMessage(`[]`), EvidenceGaps: json.RawMessage(`[]`), ForbiddenEffectsObserved: json.RawMessage(`[]`), TechnicalAssessment: TechnicalAssessment{Gate: "unknown", RequiredVerifierIDs: []string{}, AcceptedDebt: []workspace.TaskAcceptedDebtRecord{}}}, Round: round, ControlID: control, PreviousReportSHA256: previous}
			if _, err = workflowParseReport(traceJSON(report)); err != nil {
				t.Fatal(err)
			}
			reportBinding, err := workflowKeep(f.D, filepath.Join(dir, "report.json"), report)
			if err != nil {
				t.Fatal(err)
			}
			terminal := json.RawMessage(`{"historical_terminal":"unavailable native envelope; no qualification asserted"}`)
			terminalBinding, err := workflowKeep(f.D, filepath.Join(dir, "terminal-draft.json"), terminal)
			if err != nil {
				t.Fatal(err)
			}
			targetSHA := hash([]byte("explicitly unknown historical target"))
			if _, err = workflowKeep(f.D, filepath.Join(dir, "report-slot.json"), workflowReportSlot{report, terminal, targetSHA}); err != nil {
				t.Fatal(err)
			}
			review := WorkflowReview{Envelope: workflowEnv("run-review"), ReviewID: fmt.Sprintf("older-review-%d", round), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, HandoffSHA256: s.Result.Handoff.SHA256, Round: round, ReportSHA256: reportBinding.SHA256, Reviewer: old.Coordinator.ActorClaim, Decision: "changes_requested", Findings: []WorkflowFinding{{"No technical evidence.", "Preserve actual check results.", "Verify the frozen behavior."}}}
			if round == 1 {
				review.Decision = "blocked"
			}
			reviewBinding, err := workflowKeep(f.D, filepath.Join(dir, "review.json"), review)
			if err != nil {
				t.Fatal(err)
			}
			s.Records = append(s.Records, workflowRecord{Number: round, Report: reportBinding, Terminal: terminalBinding, TargetSHA256: targetSHA, Review: &reviewBinding})
			previous, control = ptr(reportBinding.SHA256), ptr(review.ReviewID)
		}
		s.Result.Round = WorkflowRound{Number: 1, State: "sealed", ReportSHA256: previous}
		s.Result.FinalReturn = WorkflowFinal{State: "blocked", ReportSHA256: previous}
		if err = workflowSave(f.D, s); err != nil {
			t.Fatal(err)
		}
		return s.Result.RunID, RunID(native.RequestKey)
	}
}

func TestReadTraceRecheckDetectsChangedSource(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "captured.json")
	if err = os.WriteFile(path, []byte(`{"version":1}`), 0600); err != nil {
		t.Fatal(err)
	}
	r := traceReader{run: TraceRun{Coverage: "complete", Sources: []FileBinding{}, Reasons: []Reason{}}, captured: map[string]inventoryCapturedSource{}}
	if _, err = r.read(path, 1024); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte(`{"version":2}`), 0600); err != nil {
		t.Fatal(err)
	}
	r.recheck()
	if r.run.Coverage != "stale" || !traceReason(r.run, "trace_source_changed") {
		t.Fatalf("changed source invisible: %+v", r.run)
	}
}

func TestReadTraceRejectsForeignGenerationResultForSameCommit(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "same candidate bytes")
	review := deliveryTestReview(t, f, "same-candidate")
	var err error
	for i := 0; i < 2; i++ {
		o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
		if err != nil {
			t.Fatal(err)
		}
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Result.Delivery.Candidates) != 2 {
		t.Fatal("fixture needs two native same-commit generations")
	}
	first, second := s.Result.Delivery.Candidates[0], s.Result.Delivery.Candidates[1]
	if first.OID != second.OID || first.TaskResult.ID == second.TaskResult.ID {
		t.Fatal("fixture lacks distinct same-commit native results")
	}
	// Every substituted record is real and registry-equal. It belongs to a
	// different verifier receipt/generation, which must not be transferable.
	s.Result.Delivery.Candidates[1].TaskResult = first.TaskResult
	s.Result.Delivery.Candidates[1].Handoff = first.Handoff
	if err = workflowSave(f.D, s); err != nil {
		t.Fatal(err)
	}
	run := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	if len(traceKind(run, "verification")) != 2 {
		t.Fatal("independent native verifier receipts disappeared")
	}
	if got := len(traceKind(run, "task_result_technical")); got != 1 || !traceReason(run, "trace_candidate_unbound") {
		t.Fatalf("foreign generation TaskResult accepted: qualified=%d reasons=%+v", got, run.Reasons)
	}
}

func TestReadTraceOmitsRunWhoseRequestDisappearsAfterInventory(t *testing.T) {
	d, request, _ := fixture(t)
	basis, err := workspace.ReadTaskJournalBasis(d.Workspace, "task")
	if err != nil {
		t.Fatal(err)
	}
	for _, family := range []string{"native", "workflow"} {
		t.Run(family, func(t *testing.T) {
			id, path := RunID(request.RequestKey), filepath.Join(runPaths(request).RunRoot, "request.json")
			var content any = request
			if family == "workflow" {
				wr := WorkflowRequest{Envelope: workflowEnv("herdr-run-request"), RequestKey: "trace/disappearing", WorkspaceRoot: request.WorkspaceRoot, PreparationID: request.PreparationID, PreparationSHA256: request.PreparationSHA256, HandoffDraft: request.HandoffDraft, Runtime: request.Runtime}
				x := workspace.PlanWorktreeObservation{StatusEntries: []workspace.IntegrationStatusEntry{}, InProgress: []string{}}
				s := workflowInitial(wr, Observed{Target: x, Epic: x, RuntimeBindings: []FileBinding{}})
				id, path, content = s.Result.RunID, workflowIndex(request.WorkspaceRoot, s.Result.RunID), s
			}
			if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				t.Fatal(err)
			}
			if err = writeValue(path, content); err != nil {
				t.Fatal(err)
			}
			inventory, err := ReadInventory(d, request.WorkspaceRoot)
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, row := range inventory.Runs {
				found = found || row.RunID == id && row.TaskID != nil && *row.TaskID == "task"
			}
			if !found {
				t.Fatal("request was not safely bound before disappearance")
			}
			if err = os.Remove(path); err != nil {
				t.Fatal(err)
			}
			r := traceReader{d: d, basis: basis, run: TraceRun{RunID: id, TaskID: "task", Coverage: "complete", Sources: []FileBinding{}, Entries: []TraceEntry{}, Reasons: []Reason{}}, captured: map[string]inventoryCapturedSource{}}
			if family == "native" {
				r.native(request.WorkspaceRoot, id)
			} else {
				r.workflow(request.WorkspaceRoot, id)
			}
			r.recheck()
			out := TraceHistory{Runs: []TraceRun{}, Reasons: []Reason{}}
			appendTraceRun(&out, r.run)
			if len(out.Runs) != 0 || len(out.Reasons) == 0 || !strings.Contains(out.Reasons[0].Detail, id) {
				t.Fatalf("disappeared request yielded incomplete run identity: %+v", out)
			}
		})
	}
}

func TestReadTraceNativeEventIDMatchesJournalFormula(t *testing.T) {
	e := TraceEntry{}
	v := workspace.TaskRecorderRecord{ActorClaim: "recorder", ControlSurface: "fixture", RecordedAtUTC: "2026-10-07T00:00:00Z"}
	traceNative(&e, "human_qa", "hqa_fixture", v)
	canonical, _ := Canonical([]string{"human_qa", "hqa_fixture", strings.TrimPrefix(digest(v), "sha256:")})
	want := "native_" + strings.TrimPrefix(hash(canonical), "sha256:")
	if e.NativeEventID == nil || *e.NativeEventID != want {
		t.Fatal(fmt.Sprint(e.NativeEventID))
	}
}
