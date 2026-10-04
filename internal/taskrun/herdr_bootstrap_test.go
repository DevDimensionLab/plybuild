package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Keep an independent, exact protocol expectation: this input grants no Task work.
const workflowExpectedBootstrapPrompt = "Transport readiness check. Reply exactly PLY_READY. Do not use tools, read files, or change anything. Wait for the next message."

type workflowBootstrapCall struct {
	Argv                       []string `json:"argv"`
	Reserved                   bool     `json:"bootstrap_reserved"`
	Phase                      string   `json:"phase_before_call"`
	Session                    string   `json:"session_before_call"`
	BootstrapSettledBeforeCall bool     `json:"bootstrap_settled_before_call"`
}

func workflowBootstrapFixture(t *testing.T, mode string) workflowFixture {
	t.Helper()
	f := workflowTestFixture(t)
	f.D.HerdrTimeout = 3 * time.Second
	workflowTestModel(t, f, map[string]any{
		"bootstrap_mode": mode, "interactive_ready": true,
		"bootstrap_expected_prompt": workflowExpectedBootstrapPrompt,
		"bootstrap_run_root":        workflowPaths(f.R, 0).RunRoot,
	})
	return f
}

func workflowBootstrapCalls(t *testing.T, f workflowFixture) []workflowBootstrapCall {
	t.Helper()
	b, err := os.ReadFile(f.Calls)
	if err != nil {
		t.Fatal(err)
	}
	var out []workflowBootstrapCall
	for _, row := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		var call workflowBootstrapCall
		if err := json.Unmarshal([]byte(row), &call); err != nil {
			t.Fatal(err)
		}
		out = append(out, call)
	}
	return out
}

func workflowBootstrapPrompts(t *testing.T, f workflowFixture) []workflowBootstrapCall {
	t.Helper()
	var out []workflowBootstrapCall
	for _, call := range workflowBootstrapCalls(t, f) {
		if len(call.Argv) >= 2 && call.Argv[0] == "agent" && call.Argv[1] == "prompt" {
			out = append(out, call)
		}
	}
	return out
}

func TestWorkflowBootstrapDiscoversSessionBeforeTask(t *testing.T) {
	f := workflowBootstrapFixture(t, "normal")
	o, err := workflowReadinessStart(t, f)
	if err != nil {
		t.Fatalf("interactive ready with no session needs one bounded bootstrap: %v", err)
	}
	if o.Transport.AgentSessionID != "fixture-session" || o.Transport.State != "working" || o.NextAction.Actor != "recipient" {
		t.Fatalf("real session and Task prompt were not bound: %+v", o)
	}
	prompts := workflowBootstrapPrompts(t, f)
	if len(prompts) != 2 {
		t.Fatalf("want exactly bootstrap then Task; got %d prompts", len(prompts))
	}
	first, task := prompts[0], prompts[1]
	want := []string{"agent", "prompt", "w-fixture:p1", workflowExpectedBootstrapPrompt, "--wait", "--until", "working", "--timeout", "10000"}
	if !equal(first.Argv, want) || !first.Reserved || first.Phase != "bootstrap_attempted" || first.Session != "" {
		t.Fatalf("bootstrap was not durably reserved before exactly one fixed input: %+v", first)
	}
	if !task.BootstrapSettledBeforeCall || task.Session != "fixture-session" || task.Argv[3] == workflowExpectedBootstrapPrompt || !strings.Contains(task.Argv[3], o.Paths.Context) {
		t.Fatalf("Task input preceded a fresh, settled native session: %+v", task)
	}
	if workflowTestCalls(t, f, "tab create") != 1 || workflowTestCalls(t, f, "agent start") != 1 {
		t.Fatal("bootstrap created a second tab or agent")
	}
	if _, err := os.Stat(filepath.Join(o.Paths.RunRoot, "startup-bootstrap-attempt.json")); err != nil {
		t.Fatal("missing durable bootstrap attempt", err)
	}
	workflowAssertNoReplay(t, f, o)
}

func TestWorkflowBootstrapLostReplyUsesFreshSessionWithoutResend(t *testing.T) {
	f := workflowBootstrapFixture(t, "lost_reply")
	o, err := workflowReadinessStart(t, f)
	if err != nil {
		t.Fatalf("fresh same-session completion should resolve a lost bootstrap reply: %v", err)
	}
	prompts := workflowBootstrapPrompts(t, f)
	if len(prompts) != 2 || prompts[0].Argv[3] != workflowExpectedBootstrapPrompt || prompts[1].Argv[3] == workflowExpectedBootstrapPrompt {
		t.Fatalf("lost reply was resent or Task input was omitted: %+v", prompts)
	}
	if !prompts[1].BootstrapSettledBeforeCall || prompts[1].Session != "fixture-session" || o.Transport.AgentSessionID != "fixture-session" {
		t.Fatal("bootstrap exit was treated as completion without fresh settled identity")
	}
	workflowAssertNoReplay(t, f, o)
}

func TestWorkflowBootstrapRequiresExplicitInteractiveReadiness(t *testing.T) {
	cases := []struct {
		name  string
		model map[string]any
	}{
		{"blocked", map[string]any{"agent_status": "blocked"}},
		{"working", map[string]any{"agent_status": "working"}},
		{"interactive_false", map[string]any{"interactive_ready": false}},
		{"interactive_null", map[string]any{"interactive_ready": nil}},
		{"interactive_missing", map[string]any{"observations": []any{map[string]any{"drop_fields": []string{"interactive_ready"}}}}},
		{"launch_pending", map[string]any{"launch_pending": true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := workflowBootstrapFixture(t, "normal")
			f.D.HerdrTimeout = 900 * time.Millisecond
			workflowTestModel(t, f, tc.model)
			o, err := workflowReadinessStart(t, f)
			if err == nil || o.RunID == "" {
				t.Fatalf("unready startup did not preserve its stop: %v", err)
			}
			if len(workflowBootstrapPrompts(t, f)) != 0 || o.Transport.AgentSessionID != "" {
				t.Fatal("unready metadata authorized bootstrap or Task input")
			}
			if _, err := os.Stat(filepath.Join(o.Paths.RunRoot, "startup-bootstrap-attempt.json")); !os.IsNotExist(err) {
				t.Fatal("unready startup reserved a bootstrap attempt")
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestWorkflowBootstrapUnresolvedNeverSendsTaskOrReplays(t *testing.T) {
	cases := []struct {
		name  string
		mode  string
		model map[string]any
	}{
		{"ack_without_session", "no_session", nil},
		{"lost_reply_without_session", "no_session", map[string]any{"lost_command": "agent prompt"}},
		{"working_without_settlement", "working_forever", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := workflowBootstrapFixture(t, tc.mode)
			f.D.HerdrTimeout = 1200 * time.Millisecond
			workflowTestModel(t, f, tc.model)
			o, err := workflowReadinessStart(t, f)
			if err == nil || o.RunID == "" {
				t.Fatalf("unknown bootstrap was accepted: %v", err)
			}
			prompts := workflowBootstrapPrompts(t, f)
			if len(prompts) != 1 || prompts[0].Argv[3] != workflowExpectedBootstrapPrompt {
				t.Fatalf("unresolved bootstrap allowed Task input or bootstrap replay: %+v", prompts)
			}
			if !prompts[0].Reserved || prompts[0].Phase != "bootstrap_attempted" {
				t.Fatal("uncertain bootstrap input was not durably reserved before sending")
			}
			if o.NextAction.Actor != "coordinator" || o.TaskResultState != "not_published" {
				t.Fatalf("uncertain bootstrap claimed Task readiness: %+v", o.NextAction)
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = os.Stat(filepathForRound(s, "prompt-attempt.json")); !os.IsNotExist(err) {
				t.Fatal("unknown bootstrap reserved a Task prompt")
			}
			if s.Acceptance != nil || s.StartSHA256 != nil || len(s.Records) != 0 || o.FinalReturn.TerminalSHA256 != nil {
				t.Fatal("unknown bootstrap fabricated an accepted Task or return")
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}

func TestWorkflowBootstrapRejectsSessionChangesBeforeTask(t *testing.T) {
	for _, mode := range []string{"session_change", "session_change_before_task"} {
		t.Run(mode, func(t *testing.T) {
			f := workflowBootstrapFixture(t, mode)
			o, err := workflowReadinessStart(t, f)
			if err == nil || o.RunID == "" {
				t.Fatalf("changed session was accepted: %v", err)
			}
			prompts := workflowBootstrapPrompts(t, f)
			if len(prompts) != 1 || prompts[0].Argv[3] != workflowExpectedBootstrapPrompt {
				t.Fatalf("changed session received Task input or another bootstrap: %+v", prompts)
			}
			if o.NextAction.Actor != "coordinator" {
				t.Fatal("replacement session was authorized for recipient work")
			}
			workflowAssertNoReplay(t, f, o)
		})
	}
}
