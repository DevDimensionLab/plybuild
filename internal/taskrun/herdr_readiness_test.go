package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func workflowPendingObservation() map[string]any {
	return map[string]any{"agent_status": "blocked", "launch_pending": true, "agent_session": nil}
}

func workflowReadinessStart(t *testing.T, f workflowFixture) (WorkflowRun, error) {
	t.Helper()
	p, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	return WorkflowStart(f.D, f.File, *p.(WorkflowPreview).Confirmation)
}

func workflowAssertNoReplay(t *testing.T, f workflowFixture, o WorkflowRun) {
	t.Helper()
	readCalls := func() ([]byte, bool) {
		data, err := os.ReadFile(f.Calls)
		if os.IsNotExist(err) {
			return nil, false
		}
		if err != nil {
			t.Fatal(err)
		}
		return data, true
	}
	before, existed := readCalls()
	repeated, err := WorkflowStart(f.D, f.File, "already reserved")
	if err != nil {
		t.Fatal(err)
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if repeated.RunID != o.RunID || shown.RunID != o.RunID || repeated.RequestSHA256 != o.RequestSHA256 || shown.RequestSHA256 != o.RequestSHA256 {
		t.Fatal("replay or show changed the reserved run identity")
	}
	after, exists := readCalls()
	if existed != exists || string(before) != string(after) {
		t.Fatal("replay or show contacted the transport")
	}
}

func TestWorkflowReadinessWaitsWithinOneStart(t *testing.T) {
	for _, earlyExit := range []bool{false, true} {
		t.Run(map[bool]string{false: "successful_start", true: "early_nonzero_start"}[earlyExit], func(t *testing.T) {
			f := workflowTestFixture(t)
			f.D.HerdrTimeout = 2 * time.Second
			model := map[string]any{"observations": []any{workflowPendingObservation(), workflowPendingObservation(), map[string]any{}}}
			if earlyExit {
				model["lost_command"] = "agent start"
			}
			workflowTestModel(t, f, model)
			o, err := workflowReadinessStart(t, f)
			if err != nil {
				t.Fatalf("same-attempt readiness did not complete: %v", err)
			}
			if o.Transport.AgentSessionID != "fixture-session" || o.Transport.State != "working" || o.NextAction.Actor != "recipient" {
				t.Fatalf("ready session not bound: %+v", o)
			}
			for _, command := range []string{"tab create", "agent start", "agent prompt"} {
				if n := workflowTestCalls(t, f, command); n != 1 {
					t.Fatalf("%s count = %d", command, n)
				}
			}
			calls, _ := os.ReadFile(f.Calls)
			gets := 0
			for _, line := range strings.Split(strings.TrimSpace(string(calls)), "\n") {
				var call struct{ Argv []string }
				if err := json.Unmarshal([]byte(line), &call); err != nil {
					t.Fatal(err)
				}
				if call.Argv[0] == "agent" && call.Argv[1] == "get" {
					gets++
				}
				if call.Argv[0] == "agent" && call.Argv[1] == "prompt" && gets < 3 {
					t.Fatal("prompt preceded fresh readiness")
				}
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestWorkflowReadinessStopsPreserveDiagnosis(t *testing.T) {
	cases := []struct {
		name  string
		model map[string]any
		want  string
	}{
		{"blocked", map[string]any{"observations": []any{workflowPendingObservation()}}, "deadline"},
		{"blocked_without_pending", map[string]any{"observations": []any{map[string]any{"agent_status": "blocked", "launch_pending": false, "agent_session": nil}}}, "deadline"},
		{"absent", map[string]any{"failure": map[string]any{"command": "agent get", "exit": 4, "stdout": `{"error":{"code":"agent_not_found","message":"synthetic-secret"}}`}}, "agent_not_found"},
		{"invalid_reply", map[string]any{"failure": map[string]any{"command": "agent get", "stdout": "synthetic-secret raw-terminal-transcript"}}, "invalid_response"},
		{"oversized_reply", map[string]any{"failure": map[string]any{"command": "agent get", "stderr": strings.Repeat("synthetic-capability", 60000)}}, "output_limit"},
		{"slow_start", map[string]any{"delay_command": "agent start", "delay_seconds": 3}, "timeout"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := workflowTestFixture(t)
			// Diagnostic cases need time to reach the intended agent-get failure,
			// including when the full suite also compiles subprocess fixtures.
			// Only deadline cases deliberately exhaust the shared startup window.
			f.D.HerdrTimeout = 5 * time.Second
			if tc.name == "blocked" || tc.name == "blocked_without_pending" || tc.name == "slow_start" {
				f.D.HerdrTimeout = 2 * time.Second
			}
			workflowTestModel(t, f, tc.model)
			o, err := workflowReadinessStart(t, f)
			if err == nil || o.RunID == "" {
				t.Fatalf("missing preserved stop: %v", err)
			}
			if o.NextAction.Actor != "coordinator" || strings.Contains(o.NextAction.Message, "Accept") {
				t.Errorf("unbound recipient was told to act: %+v", o.NextAction)
			}
			b, _ := json.Marshal(o.Reasons)
			// The shared deadline can expire between observations or inside the
			// last subprocess; both must retain the actual timeout diagnosis.
			deadlineInChild := (tc.name == "blocked" || tc.name == "blocked_without_pending") && strings.Contains(string(b), "class=timeout")
			if !strings.Contains(string(b), tc.want) && !deadlineInChild {
				t.Errorf("missing preserved %s diagnosis: %s", tc.want, b)
			}
			if len(b) > 2048 {
				t.Error("preserved diagnosis is unbounded")
			}
			rawState, err := os.ReadFile(filepath.Join(o.Paths.RunRoot, "state.json"))
			if err != nil {
				t.Fatal(err)
			}
			for _, forbidden := range []string{"synthetic-secret", "synthetic-capability", "raw-terminal-transcript"} {
				if strings.Contains(string(rawState), forbidden) {
					t.Error("child output leaked into preserved state")
				}
			}
			if o.Transport.AgentSessionID != "" || workflowTestCalls(t, f, "agent prompt") != 0 || o.TaskResultState != "not_published" {
				t.Fatal("stopped launch gained readiness or effects")
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if s.Acceptance != nil || s.StartSHA256 != nil || len(s.Records) != 0 || s.Result.FinalReturn.TerminalSHA256 != nil {
				t.Fatal("stop fabricated native return")
			}
			if _, err := os.Stat(filepathForRound(s, "prompt-attempt.json")); !os.IsNotExist(err) {
				t.Fatal("stop reserved a prompt")
			}
			if tc.name == "blocked" {
				raw, _ := os.ReadFile(filepath.Join(o.Paths.RunRoot, "state.json"))
				var state map[string]any
				json.Unmarshal(raw, &state)
				transport := state["result"].(map[string]any)["transport"].(map[string]any)
				if transport["state"] != "blocked" || transport["launch_pending"] != true {
					t.Errorf("pending observation lost: %v", transport)
				}
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestWorkflowShowCorrectsUnboundLegacyActionWithoutEffects(t *testing.T) {
	f := workflowTestFixture(t)
	f.D.HerdrTimeout = 500 * time.Millisecond
	workflowTestModel(t, f, map[string]any{"observations": []any{workflowPendingObservation()}})
	o, err := workflowReadinessStart(t, f)
	if err == nil {
		t.Fatal("expected pending startup stop")
	}
	// Model an old stored run. Readback must correct the advice without migrating
	// its saved bytes, binding a cached session, observing Herdr or restarting.
	err = workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		s.Result.NextAction = WorkflowAction{"recipient", "Accept the bound runtime before target writes, then report this round."}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(o.Paths.RunRoot, "state.json")
	before, _ := os.ReadFile(p)
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if shown.NextAction.Actor != "coordinator" || shown.Transport.AgentSessionID != "" {
		t.Errorf("unbound cached state gave recipient authority: %+v", shown.NextAction)
	}
	workflowAssertNoReplay(t, f, o)
	after, _ := os.ReadFile(p)
	if string(before) != string(after) {
		t.Fatal("readback rewrote legacy state")
	}
}

func TestWorkflowReadinessRejectsIncompleteOrChangedIdentity(t *testing.T) {
	cases := map[string]map[string]any{
		"workspace": {"workspace_id": "other"}, "tab": {"tab_id": "other"},
		"pane": {"pane_id": "other"}, "terminal": {"terminal_id": "other"},
		"name": {"name": "other"}, "agent": {"agent": "other"},
		"session_agent":        {"agent_session": map[string]any{"agent": "other", "kind": "id", "value": "fixture-session"}},
		"session_kind":         {"agent_session": map[string]any{"agent": "codex", "kind": "other", "value": "fixture-session"}},
		"session_empty":        {"agent_session": map[string]any{"agent": "codex", "kind": "id", "value": ""}},
		"pending_with_session": {"launch_pending": true},
	}
	for _, field := range []string{"workspace_id", "tab_id", "pane_id", "terminal_id", "name", "agent", "agent_session", "agent_status"} {
		cases["missing_"+field] = map[string]any{"drop_fields": []string{field}}
	}
	for name, observation := range cases {
		t.Run(name, func(t *testing.T) {
			f := workflowTestFixture(t)
			f.D.HerdrTimeout = 5 * time.Second
			model := map[string]any{"observations": []any{observation}}
			if name == "tab" {
				model["before_call_delay_command"] = "tab create"
				model["before_call_delay_millis"] = 750
			}
			workflowTestModel(t, f, model)
			o, err := workflowReadinessStart(t, f)
			if workflowTestCalls(t, f, "agent get") == 0 {
				t.Fatal("fixture stopped before observing the selected identity")
			}
			data, readErr := os.ReadFile(f.Model)
			var observed struct {
				Index int `json:"observation_index"`
			}
			if readErr != nil || json.Unmarshal(data, &observed) != nil || observed.Index == 0 {
				t.Fatal("fixture never consumed the selected identity observation")
			}
			if err == nil || workflowTestCalls(t, f, "agent prompt") != 0 {
				t.Fatal("invalid identity permitted prompt")
			}
			switch name {
			case "pending_with_session", "missing_agent_session", "missing_agent_status":
				if !strings.Contains(err.Error(), "deadline") && !strings.Contains(err.Error(), "class=timeout") {
					t.Fatalf("incomplete readiness did not retain its deadline stop: %v", err)
				}
			default:
				if !strings.Contains(err.Error(), "differs from the reserved attempt") && !strings.Contains(err.Error(), "differs from the bound run") {
					t.Fatalf("identity observation did not cause the rejection: %v", err)
				}
			}
			state, readErr := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if readErr != nil || state.Phase != "agent_start_attempted" || o.Transport.AgentSessionID != "" || workflowTestCalls(t, f, "tab create") != 1 || workflowTestCalls(t, f, "agent start") != 1 {
				t.Fatalf("identity fixture did not preserve its unbound startup: phase=%s, %v", state.Phase, readErr)
			}
			if o.NextAction.Actor != "coordinator" {
				t.Error("stop did not name coordinator")
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
	t.Run("session_changes_before_prompt", func(t *testing.T) {
		f := workflowTestFixture(t)
		workflowTestModel(t, f, map[string]any{"observations": []any{map[string]any{}, map[string]any{"agent_session": map[string]any{"agent": "codex", "kind": "id", "value": "replacement"}}}})
		o, err := workflowReadinessStart(t, f)
		if err == nil || workflowTestCalls(t, f, "agent prompt") != 0 {
			t.Fatal("replacement session permitted prompt")
		}
		if workflowTestCalls(t, f, "agent get") != 2 || !strings.Contains(err.Error(), "differs from the bound run") {
			t.Fatalf("fixture did not reject its second changed-session observation: %v", err)
		}
		workflowAssertNoReplay(t, f, o)
	})
}

func TestWorkflowReadinessUsesOneDeadline(t *testing.T) {
	f := workflowTestFixture(t)
	const window = 3 * time.Second
	f.D.HerdrTimeout = window
	workflowTestModel(t, f, map[string]any{"observations": []any{workflowPendingObservation()}, "delay_command": "agent get", "delay_millis": 200, "before_call_delay_command": "agent start", "before_call_delay_millis": 1250})
	p, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	var started time.Time
	// Mark the Go send boundary, before the child can be scheduled. For this
	// fixture no trust/environment work intervenes before the shared deadline.
	// The delayed start consumes that same window as subsequent observations.
	f.D.Fault = func(point string) error {
		if point == "workflow_before_agent_send" {
			started = time.Now()
		}
		return nil
	}
	o, err := WorkflowStart(f.D, f.File, *p.(WorkflowPreview).Confirmation)
	if started.IsZero() {
		t.Fatal("fixture never reached the Go startup boundary")
	}
	elapsed := time.Since(started)
	if err == nil || elapsed > window+900*time.Millisecond || elapsed < window-50*time.Millisecond {
		t.Fatalf("startup deadline not shared: %v, %v", elapsed, err)
	}
	if workflowTestCalls(t, f, "agent get") < 2 || workflowTestCalls(t, f, "agent prompt") != 0 {
		t.Fatal("bounded polling not observed")
	}
	workflowAssertNoReplay(t, f, o)
}

func TestWorkflowTransportDiagnosticsAreBoundedAndSafe(t *testing.T) {
	secret := "synthetic-secret-CAPABILITY-DO-NOT-EXPOSE"
	for _, tc := range []struct {
		name    string
		failure map[string]any
		want    []string
	}{
		{"exit", map[string]any{"exit": 7, "stderr": `{"detail":"` + secret + `"}`, "stdout": `{"error":{"code":"agent_not_found","message":"` + secret + `"},"terminal":"` + secret + `"}`}, []string{"agent get", "exit=7", "agent_not_found"}},
		{"mixed_invalid", map[string]any{"exit": 7, "stderr": secret, "stdout": `{"error":{"code":"agent_not_found"}}`}, []string{"agent get", "exit=7", "provider_code=unknown"}},
		{"timeout", map[string]any{"delay_seconds": 2, "stderr": secret}, []string{"agent get", "timeout"}},
		{"invalid", map[string]any{"stdout": secret}, []string{"agent get", "invalid_response"}},
		{"large", map[string]any{"exit": 7, "stderr": strings.Repeat(secret, 40000)}, []string{"agent get", "exit=7", "truncated=true", fmt.Sprintf("stderr_bytes=%d", len(secret)*40000)}},
		{"unknown_code", map[string]any{"exit": 7, "stdout": `{"error":{"code":"` + secret + `"}}`}, []string{"agent get", "unknown"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := workflowTestFixture(t)
			// Only the timeout case should depend on a short deadline. Other
			// cases must reach their actual subprocess exit/output condition.
			f.D.HerdrTimeout = 5 * time.Second
			if tc.name == "timeout" {
				f.D.HerdrTimeout = 500 * time.Millisecond
			}
			tc.failure["command"] = "agent get"
			model := map[string]any{"failure": tc.failure}
			if tc.name == "large" {
				// Reproduce a process scheduled after the old 500 ms budget.
				model["before_call_delay_command"] = "agent get"
				model["before_call_delay_millis"] = 750
			}
			workflowTestModel(t, f, model)
			_, err := workflowCall(f.D, f.R, "agent", "get", "w-fixture:p1")
			if err == nil {
				t.Fatal("failure was accepted")
			}
			message := err.Error()
			for _, want := range tc.want {
				if !strings.Contains(message, want) {
					t.Errorf("missing %s in %s", want, message)
				}
			}
			if len(message) > 1024 || strings.Contains(message, secret) {
				t.Fatalf("unsafe diagnostic (%d bytes)", len(message))
			}
		})
	}
}

func TestWorkflowStoppedBeforePromptCannotAccept(t *testing.T) {
	f := workflowTestFixture(t)
	workflowTestModel(t, f, map[string]any{"observations": []any{map[string]any{}, map[string]any{"agent_session": map[string]any{"agent": "codex", "kind": "id", "value": "replacement"}}}})
	o, err := workflowReadinessStart(t, f)
	if err == nil || workflowTestCalls(t, f, "agent prompt") != 0 {
		t.Fatal("expected pre-prompt identity stop")
	}
	p := writeAny(t, f.R.WorkspaceRoot, "acceptance.json", f.Acceptance)
	if _, err = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, p); err == nil {
		t.Fatal("a stopped pre-prompt attempt published runtime acceptance")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Acceptance != nil || s.StartSHA256 != nil {
		t.Fatal("pre-prompt stop created native acceptance")
	}
}

func TestWorkflowCorrectionKeepsNextAction(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestReady(t, f)
	o, err := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "changes_requested"))
	if err != nil {
		t.Fatal(err)
	}
	if o.NextAction.Actor != "recipient" || !strings.Contains(o.NextAction.Message, "findings") || strings.Contains(o.NextAction.Message, "Accept") {
		t.Fatalf("correction was told to accept again: %+v", o.NextAction)
	}
}

func TestWorkflowPromptCompletionPreservesCallbackAction(t *testing.T) {
	for _, reported := range []bool{false, true} {
		t.Run(map[bool]string{false: "negative_acceptance", true: "report_before_lost_reply"}[reported], func(t *testing.T) {
			f := workflowTestFixture(t)
			model := map[string]any{"delay_command": "agent prompt", "delay_millis": 3000}
			if reported {
				model["lost_command"] = "agent prompt"
			}
			workflowTestModel(t, f, model)
			p, err := WorkflowPreviewStart(f.D, f.File)
			if err != nil {
				t.Fatal(err)
			}
			type result struct {
				run WorkflowRun
				err error
			}
			finished := make(chan result, 1)
			go func() {
				o, e := WorkflowStart(f.D, f.File, *p.(WorkflowPreview).Confirmation)
				finished <- result{o, e}
			}()
			defer func() {
				if finished != nil {
					<-finished
				}
			}()
			deadline := time.Now().Add(10 * time.Second)
			for {
				b, _ := os.ReadFile(f.Model)
				var model map[string]any
				if json.Unmarshal(b, &model) == nil && model["agent_status"] == "working" {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("prompt subprocess did not begin")
				}
				time.Sleep(10 * time.Millisecond)
			}
			o, err := WorkflowShow(f.D, f.R.WorkspaceRoot, workflowID(f.R))
			if err != nil {
				t.Fatal(err)
			}
			if reported {
				workflowTestReport(t, f, workflowTestAccept(t, f, o))
			} else {
				a := f.Acceptance
				a.Acceptance = "unknown"
				a.RuntimeClaim.RuntimeID, a.RuntimeClaim.ProfileID, a.RuntimeClaim.EffectivePolicySHA256 = nil, nil, nil
				a.Sandbox = json.RawMessage("null")
				a.Issues = json.RawMessage(`[{"type":"unknown","detail":"Necessary authority is unknown."}]`)
				if _, err = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "negative.json", a)); err != nil {
					t.Fatal(err)
				}
			}
			r := <-finished
			finished = nil
			if reported {
				if r.err == nil || r.run.Round.State != "report_received" || r.run.NextAction.Actor != "coordinator" || !strings.Contains(r.run.NextAction.Message, "immutable report") {
					t.Fatalf("late lost reply overwrote a received report: %s, %+v, %v", r.run.Round.State, r.run.NextAction, r.err)
				}
				return
			}
			if r.err != nil {
				t.Fatal(r.err)
			}
			if r.run.FinalReturn.State != "blocked" || r.run.NextAction.Actor != "coordinator" {
				t.Fatalf("prompt completion overwrote the runtime stop: %+v", r.run.NextAction)
			}
		})
	}
}
