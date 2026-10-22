package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func al51StartupFixture(t *testing.T, receipt bool) deliveryFixture {
	t.Helper()
	options := deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{
		SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic,
		Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic AL51 fixture", TimeoutSeconds: 30},
	}}
	if receipt {
		return receiptFixture(t, "codex", options)
	}
	return newDeliveryFixture(t, "codex", options)
}

func al51FailedStartup(t *testing.T, f deliveryFixture, bootstrap bool) WorkflowRun {
	t.Helper()
	if bootstrap {
		// Reuse the injected post-readiness failure, without calling recovery or
		// replacing a provider. All Herdr calls reach only the Python stand-in.
		herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
		if err != nil {
			t.Fatal(err)
		}
		herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+deliveryRecoveryInspectStandin), 1)
		if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
			t.Fatal(err)
		}
		f.R.Herdr.Executable.SHA256 = hash(herdr)
		workflowTestModel(t, f.workflowFixture, map[string]any{
			"bootstrap_mode": "no_session", "interactive_ready": true,
			"bootstrap_expected_prompt": workflowExpectedBootstrapPrompt,
			"bootstrap_run_root":        workflowPaths(f.R, 0).RunRoot,
		})
	} else {
		workflowTestModel(t, f.workflowFixture, map[string]any{
			"lost_command": "agent start",
			"failure":      map[string]any{"command": "agent get", "exit": 1, "stderr": `{"error":{"code":"agent_not_found"}}`},
		})
	}
	f.File = writeAny(t, f.R.WorkspaceRoot, "al51-request.json", f.R)
	o, err := workflowReadinessStart(t, f.workflowFixture)
	if err == nil || !strings.Contains(err.Error(), "provider_code=agent_not_found") {
		t.Fatalf("AL51 did not reach the injected failed observation: %+v %v", o, err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	wantPhase, wantPrompts := "agent_start_attempted", 0
	if bootstrap {
		wantPhase, wantPrompts = "bootstrap_attempted", 1
	}
	if err != nil || s.Phase != wantPhase || o.Transport.AgentSessionID != "" || s.Acceptance != nil || s.StartSHA256 != nil || s.StartDraft != nil {
		t.Fatalf("AL51 failed-start boundary differs: phase=%s, %v", s.Phase, err)
	}
	if _, err = os.Lstat(filepathForRound(s, "prompt-attempt.json")); !os.IsNotExist(err) {
		t.Fatalf("AL51 attempted Task input: %v", err)
	}
	for cmd, want := range map[string]int{"tab create": 1, "agent start": 1, "agent prompt": wantPrompts} {
		if got := workflowTestCalls(t, f.workflowFixture, cmd); got != want {
			t.Fatalf("AL51 %s count=%d, want %d", cmd, got, want)
		}
	}
	return o
}

func al51AssertInspection(t *testing.T, action WorkflowAction) {
	t.Helper()
	if action.Actor != "coordinator" || !strings.Contains(action.Message, "ply workflow execute recover-start wfr_") || !strings.Contains(action.Message, " --check --format json") {
		t.Errorf("AL51 failed unbound startup must advise coordinator inspection: %+v", action)
	}
	for _, falseClaim := range []string{"Complete any native onboarding", "shell-only", "writer absent", "restart-ready"} {
		if strings.Contains(action.Message, falseClaim) {
			t.Errorf("AL51 misleading failed-start claim %q: %+v", falseClaim, action)
		}
	}
}

func TestAL51FailedUnboundStartupGuidance(t *testing.T) {
	for _, receipt := range []bool{false, true} {
		for _, bootstrap := range []bool{false, true} {
			name := map[bool]string{false: "legacy", true: "receipt"}[receipt] + "/" + map[bool]string{false: "failed_start", true: "after_readiness_prompt"}[bootstrap]
			t.Run(name, func(t *testing.T) {
				f := al51StartupFixture(t, receipt)
				o := al51FailedStartup(t, f, bootstrap)
				al51AssertInspection(t, o.NextAction)
				// Old controls saved this misleading advice. Readback must project
				// the correction without rewriting history or observing processes.
				if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
					s.Result.NextAction = WorkflowAction{"user", "Complete any native onboarding, then run: ply workflow execute resume " + o.RunID}
					return nil
				}); err != nil {
					t.Fatal(err)
				}
				f.D.StartupProcessTree = func(int) ([]StartupProcess, error) { t.Fatal("AL51 guidance probed processes"); return nil, nil }
				before := treeState(t, f.R.WorkspaceRoot)
				shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
				if err != nil {
					t.Fatal(err)
				}
				al51AssertInspection(t, shown.NextAction)
				if !equal(shown.Reasons, o.Reasons) || shown.Transport.Observation != "cached" || shown.Transport.AgentSessionID != "" || shown.RunID != o.RunID || shown.RequestSHA256 != o.RequestSHA256 {
					t.Fatal("AL51 readback changed startup diagnosis or identity")
				}
				row := deliveryInventoryRows(t, f)[o.RunID]
				if row.Delivery == nil || row.Freshness != "fresh" {
					t.Fatalf("AL51 inventory lost the preserved startup: %+v", row)
				}
				al51AssertInspection(t, WorkflowAction{row.Delivery.NextAction.RecordedActor, row.Delivery.NextAction.Reason})
				if row.Delivery.NextAction.Actor != "agent" || row.StateFreshness != "unknown" || !hasInventoryReason(row, "herdr_start_stopped") {
					t.Fatalf("AL51 inventory invented liveness or lost the error: %+v", row)
				}
				if before != treeState(t, f.R.WorkspaceRoot) {
					t.Fatal("AL51 startup reads changed fixture files or sent input")
				}
			})
		}
	}
}

func TestAL51FailedStartupCLIRead(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("AL51 public reads require an explicitly supplied candidate")
	}
	f := al51StartupFixture(t, true)
	o := al51FailedStartup(t, f, false)
	before := treeState(t, f.R.WorkspaceRoot)
	for _, format := range []string{"json", "text"} {
		c := exec.Command(binary, "workflow", "execute", "show", o.RunID, "--format", format)
		c.Dir = f.Prepared.Preparation.Preparation.Plan.WorktreePath
		output, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("AL51 CLI read: %v\n%s", err, output)
		}
		if format == "json" {
			var shown WorkflowRun
			if err = json.Unmarshal(output, &shown); err != nil {
				t.Fatal(err)
			}
			al51AssertInspection(t, shown.NextAction)
		} else if !strings.Contains(string(output), "Next: coordinator") || !strings.Contains(string(output), "recover-start "+o.RunID+" --check --format json") || !strings.Contains(string(output), "provider_code=agent_not_found") {
			t.Fatalf("AL51 text read lost inspection or diagnosis: %s", output)
		}
	}
	if before != treeState(t, f.R.WorkspaceRoot) {
		t.Fatal("AL51 public read changed fixture files or sent input")
	}
	if _, err := os.Stat(filepath.Join(o.Paths.RunRoot, "rounds", "000", "prompt-attempt.json")); !os.IsNotExist(err) {
		t.Fatal("AL51 public read attempted Task input")
	}
}

func TestAL51StartupGuidanceBoundaries(t *testing.T) {
	f, o := deliveryPausedStartup(t)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	// An in-progress onboarding observation has no completed startup stop.
	// Cached fields alone must not select failed-start or restart guidance.
	live := s
	live.Result.Reasons = nil
	live.Result.Transport.State, live.Result.Transport.LaunchPending = "blocked", ptr(true)
	action := deliveryStartupAction(live)
	if action.Actor != "user" || !strings.Contains(action.Message, "workflow execute resume "+o.RunID) || strings.Contains(action.Message, "recover-start") {
		t.Fatalf("AL51 changed in-progress onboarding guidance: %+v", action)
	}
	if !strings.Contains(o.NextAction.Message, "If native onboarding is still active") || !strings.Contains(o.NextAction.Message, "workflow execute resume "+o.RunID) {
		t.Fatalf("AL51 stopped observation lost live onboarding's conditional continuation: %+v", o.NextAction)
	}
	// A later preserved stop or necessary question must retain its actor and
	// message even when the original pre-prompt startup diagnosis remains.
	for _, phase := range []string{"stopped", "needs_input"} {
		later := s
		delivery := *s.Result.Delivery
		later.Result.Delivery = &delivery
		later.Result.Delivery.Phase = phase
		later.Result.NextAction = WorkflowAction{"user", "Inspect the preserved necessary decision?"}
		projected := inventoryDeliveryAction(later)
		if projected.RecordedActor != "user" || projected.Reason != later.Result.NextAction.Message {
			t.Errorf("AL51 startup advice replaced a later %s: %+v", phase, projected)
		}
	}
	// A possible Task send takes priority over pre-prompt recovery advice, even
	// with the same failed observation and no bound session in cached state.
	writeAny(t, filepath.Dir(filepathForRound(s, "prompt-attempt.json")), "prompt-attempt.json", map[string]any{"argv": []string{"agent", "prompt", "possibly-sent"}})
	action = deliveryStartupAction(s)
	if action.Actor != "coordinator" || !strings.Contains(action.Message, "may already have been attempted") || strings.Contains(action.Message, "recover-start") || strings.Contains(action.Message, "execute resume") {
		t.Fatalf("AL51 uncertain Task input lost its inspection boundary: %+v", action)
	}
	before := treeState(t, f.R.WorkspaceRoot)
	if _, err = WorkflowResumeDeliveryStart(f.D, f.R.WorkspaceRoot, o.RunID, time.Second); err == nil {
		t.Fatal("AL51 uncertain Task input allowed resume")
	}
	if before != treeState(t, f.R.WorkspaceRoot) {
		t.Fatal("AL51 rejected uncertain resume changed history or sent input")
	}
}
