package taskrun

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func deliveryPausedStartup(t *testing.T) (deliveryFixture, WorkflowRun) {
	t.Helper()
	f := newDeliveryFixture(t, "codex")
	// Stop at an explicit failed observation after the tab and agent exist.
	// Scheduler delay must not turn this into an earlier tab-creation timeout.
	workflowTestModel(t, f.workflowFixture, map[string]any{"observations": []any{workflowPendingObservation()}, "before_call_delay_command": "tab create", "before_call_delay_millis": 750, "failure": map[string]any{"command": "agent get", "exit": 7}})
	o, err := workflowReadinessStart(t, f.workflowFixture)
	if err == nil || o.RunID == "" || o.Transport.AgentSessionID != "" || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 0 {
		t.Fatalf("fixture must stop before initial Task prompt: %+v %v", o, err)
	}
	if !strings.Contains(err.Error(), "Herdr agent get failed: class=nonzero_exit exit=7") {
		t.Fatalf("fixture did not reach its explicit paused observation: %v", err)
	}
	state, readErr := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if readErr != nil || state.Phase != "agent_start_attempted" || state.Acceptance != nil || state.StartDraft != nil || state.StartSHA256 != nil || workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 {
		t.Fatalf("fixture did not establish its pre-Task tab/start: phase=%s, %v", state.Phase, readErr)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"failure": map[string]any{}})
	return f, o
}

func TestDeliveryResumeFirstPromptInExistingTrustedTab(t *testing.T) {
	f, o := deliveryPausedStartup(t)
	if !strings.Contains(o.NextAction.Message, "workflow run show "+o.RunID) {
		t.Fatalf("missing actionable recovery: %+v", o.NextAction)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"observations": nil, "agent_status": "idle"})
	var err error
	o, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, 3*time.Second)
	if err != nil || o.Transport.AgentSessionID != "fixture-session" {
		t.Fatalf("same-tab resume failed: %+v %v", o, err)
	}
	for _, cmd := range []string{"tab create", "agent start", "agent prompt"} {
		if n := workflowTestCalls(t, f.workflowFixture, cmd); n != 1 {
			t.Fatalf("%s count=%d", cmd, n)
		}
	}
	if _, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, time.Second); err == nil {
		t.Fatal("attempted Task prompt could be resent")
	}
	if n := workflowTestCalls(t, f.workflowFixture, "agent prompt"); n != 1 {
		t.Fatal("prompt replayed")
	}
}

func TestDeliveryResumeRefusesUncertainPromptAndIdentityDrift(t *testing.T) {
	t.Run("prompt_reserved", func(t *testing.T) {
		f := newDeliveryFixture(t, "codex")
		f.D.Fault = func(p string) error {
			if p == "workflow_before_prompt_send" {
				return errors.New("synthetic lost send outcome")
			}
			return nil
		}
		o, err := workflowReadinessStart(t, f.workflowFixture)
		if err == nil {
			t.Fatal("fixture did not reserve uncertain prompt")
		}
		f.D.Fault = nil
		if _, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, time.Second); err == nil {
			t.Fatal("uncertain prompt replay allowed")
		}
		if n := workflowTestCalls(t, f.workflowFixture, "agent prompt"); n != 0 {
			t.Fatal("uncertain prompt sent")
		}
	})
	t.Run("identity_drift", func(t *testing.T) {
		f, o := deliveryPausedStartup(t)
		workflowTestModel(t, f.workflowFixture, map[string]any{"observations": []any{map[string]any{"pane_id": "wrong-pane"}}, "agent_status": "idle"})
		if _, err := WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, time.Second); err == nil {
			t.Fatal("changed pane accepted")
		}
		if n := workflowTestCalls(t, f.workflowFixture, "agent prompt"); n != 0 {
			t.Fatal("Task prompt sent to changed pane")
		}
	})
}

func TestDeliveryResumeUsesOnlyOneSupportedReadinessExchange(t *testing.T) {
	f, o := deliveryPausedStartup(t)
	workflowTestModel(t, f.workflowFixture, map[string]any{"observations": nil, "agent_status": "idle", "bootstrap_mode": "normal", "interactive_ready": true, "bootstrap_expected_prompt": workflowExpectedBootstrapPrompt, "bootstrap_run_root": o.Paths.RunRoot})
	var err error
	o, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	ps := workflowBootstrapPrompts(t, f.workflowFixture)
	if len(ps) != 2 || ps[0].Argv[3] != workflowExpectedBootstrapPrompt || ps[1].Argv[3] == workflowExpectedBootstrapPrompt {
		t.Fatalf("expected one readiness exchange then first Task: %+v", ps)
	}
	if _, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, time.Second); err == nil {
		t.Fatal("prompt attempted should stop resume")
	}
	if len(workflowBootstrapPrompts(t, f.workflowFixture)) != 2 {
		t.Fatal("resume replayed a readiness exchange")
	}
}
