package taskrun

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// Model the pane disappearing after pane get, before the follow's agent get.
// The later inventory/process evidence is synthetic, never real OS coverage.
const al21CloseoutRaceStandin = `
if cmd == ["agent", "get"] and model.get("al21_disappear"):
    model["al21_disappear"] = False
    model["closed_transport"] = True
    path.write_text(json.dumps(model))
`

func TestAL21CloseoutObservation(t *testing.T) {
	for _, kind := range []string{"disappears", "transient", "persistent", "unknown_inventory", "conflicting_inventory", "local_writer", "unknown_processes", "changed_session", "changed_owner", "live_owner", "observer_window", "source_drift"} {
		t.Run(kind, func(t *testing.T) {
			f, o, _, _ := autoGateFixture(t, false, "")
			pending, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
			cleanupOwnerHandoff(t, f, pending, err)
			reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
			model := map[string]any{"agent_status": "idle", "al21_disappear": true}
			d := cleanupAt(f.D, f.Parent)
			d.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
			d.HerdrTimeout = 300 * time.Millisecond
			window := 5 * time.Second
			switch kind {
			case "live_owner":
				model["al21_disappear"] = false
				model["agent_status"] = "working"
				window = time.Second
			case "observer_window":
				model["al21_disappear"] = false
				model["delay_command"] = "agent get"
				model["delay_millis"] = 700
				d.HerdrTimeout = 2 * time.Second
				window = 200 * time.Millisecond
			case "source_drift":
				if err := os.WriteFile(f.Prepared.Preparation.Preparation.Plan.WorktreePath+"/late-output", []byte("preserve"), 0600); err != nil {
					t.Fatal(err)
				}
			case "transient", "persistent":
				model["al21_disappear"] = false
				last := map[string]any{"agent_status": "idle"}
				if kind == "persistent" {
					last["delay_millis"] = 700
				}
				model["observations"] = []any{map[string]any{"delay_millis": 700}, last}
			case "unknown_inventory":
				model["inventory"] = []any{map[string]any{"pane_id": "incomplete"}}
			case "conflicting_inventory":
				model["recovery_agents"] = []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "other", "pane_id": "other", "terminal_id": "other", "agent": "codex", "agent_session": map[string]any{"agent": "codex", "kind": "id", "value": o.Transport.AgentSessionID}}}
			case "local_writer":
				d.LocalProcesses = func(string, string) ([]LocalProcess, error) {
					return []LocalProcess{{PID: os.Getpid()}, {PID: 4242, CWD: f.Prepared.Preparation.Preparation.Plan.WorktreePath}}, nil
				}
			case "unknown_processes":
				d.LocalProcesses = func(string, string) ([]LocalProcess, error) { return nil, os.ErrPermission }
			case "changed_owner":
				model["al21_disappear"] = false
				model["observations"] = []any{map[string]any{"name": "different-owner"}}
			case "changed_session":
				model["al21_disappear"] = false
				model["observations"] = []any{map[string]any{"agent_session": map[string]any{"agent": "codex", "kind": "id", "value": "different"}}}
			}
			workflowTestModel(t, f.workflowFixture, model)
			var warnings bytes.Buffer
			old := logrus.StandardLogger().Out
			logrus.SetOutput(&warnings)
			defer logrus.SetOutput(old)
			before := workflowTestCalls(t, f.workflowFixture, "agent get")
			got, err := WorkflowFollow(d, f.R.WorkspaceRoot, o.RunID, window)
			reads := workflowTestCalls(t, f.workflowFixture, "agent get") - before
			if kind == "disappears" || kind == "transient" {
				if err != nil {
					t.Fatalf("same closeout did not continue: %v", err)
				}
				assertCleanupRemoved(t, f, got, reflog)
				if len(got.Delivery.Events) != len(pending.Delivery.Events)+1 || got.Delivery.Events[len(got.Delivery.Events)-1].Kind != "closeout" || !equal(got.Delivery.Candidates, pending.Delivery.Candidates) {
					t.Fatal("observation changed candidate/history beyond the same closeout")
				}
				if !strings.Contains(warnings.String(), o.RunID) || !strings.Contains(warnings.String(), "agent get") {
					t.Fatalf("missing read warning: %s", warnings.String())
				}
			} else {
				if err == nil || got.Delivery.Phase != "closing" || got.DeliveryStatus.FinalDelivery.State != "delivered" {
					t.Errorf("unsafe cleanup: %v %+v", err, got.DeliveryStatus)
				}
				if _, err := os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); err != nil {
					t.Error("uncertain source removed", err)
				}
				if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
					t.Error("closed an uncertain pane")
				}
				if kind == "observer_window" && (reads != 1 || err == nil || !strings.Contains(err.Error(), "observation timed out without stopping")) {
					t.Errorf("window escaped: reads=%d error=%v", reads, err)
				}
				if kind == "persistent" && (reads != 2 || got.Transport.State != "unknown") {
					t.Fatalf("unbounded or live timeout: reads=%d state=%s", reads, got.Transport.State)
				}
			}
			if len(got.Delivery.Candidates) != len(o.Delivery.Candidates) || got.Transport.AgentSessionID != o.Transport.AgentSessionID || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 || gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
				t.Fatal("observation replayed effects or changed identity")
			}
		})
	}
}

func TestAL21CloseoutCandidateCLI(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := cleanupCLIFixture(t, binary)
	pending, reflog := cleanupOwnerIntegrates(t, binary, f, o)
	workflowTestModel(t, f.workflowFixture, map[string]any{"al21_disappear": true, "inventory": []any{map[string]any{"pane_id": "incomplete"}}})
	c := exec.Command(binary, "workflow", "execute", "follow", o.RunID, "--timeout", "15", "--format", "json")
	c.Dir = f.Parent
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	err := c.Run()
	t.Logf("modeled Herdr/process race; candidate=%s exit=%d stdout=%s stderr=%s", binary, c.ProcessState.ExitCode(), stdout.String(), stderr.String())
	if err == nil || c.ProcessState.ExitCode() != 4 || !json.Valid(stdout.Bytes()) {
		t.Fatalf("follow failed or polluted JSON: %v", err)
	}
	var got WorkflowRun
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Delivery.Phase != "closing" || got.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("lost delivered closeout: %+v", got.DeliveryStatus)
	}
	if _, err := os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); err != nil {
		t.Fatal("uncertain source removed", err)
	}
	if gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic") != reflog {
		t.Fatal("integration repeated")
	}
	if !strings.Contains(stderr.String(), o.RunID) || !strings.Contains(stderr.String(), "agent_not_found") || !strings.Contains(stderr.String(), "continuing writer inspection") {
		t.Fatalf("missing sourced stderr warning: %s", stderr.String())
	}
	if len(got.Delivery.Events) != len(pending.Delivery.Events) || !equal(got.Delivery.Candidates, pending.Delivery.Candidates) || got.RunID != o.RunID || got.Transport.AgentSessionID != o.Transport.AgentSessionID || workflowTestCalls(t, f.workflowFixture, "pane close") != 0 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("race replayed an effect or changed the delivered candidate/history")
	}
	shown, output, err := cleanupCLIAt(binary, f.Parent, "execute", "show", o.RunID)
	if err != nil || shown.Delivery.Phase != "closing" {
		t.Fatalf("public readback: %v %s", err, output)
	}
}
