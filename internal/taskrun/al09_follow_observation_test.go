package taskrun

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// A request that would succeed in 400ms is killed by the observer's 200ms
// window, despite the independent 2s request budget. This measures the causal
// distinction before introducing any read continuation.
func TestAL09FollowWindowClipsHealthyRequest(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	f.D.HerdrTimeout = 2 * time.Second
	workflowTestModel(t, f.workflowFixture, map[string]any{"delay_command": "agent get", "delay_millis": 400})
	before := workflowTestCalls(t, f.workflowFixture, "agent get")
	start := time.Now()
	got, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 200*time.Millisecond)
	elapsed := time.Since(start)
	reads := workflowTestCalls(t, f.workflowFixture, "agent get") - before
	t.Logf("observer_budget=200ms request_budget=2s response_delay=400ms elapsed=%s reads=%d error=%v", elapsed, reads, err)
	if err == nil || !strings.Contains(err.Error(), "observation timed out without stopping or restarting") || strings.Contains(err.Error(), "Herdr agent get failed") {
		t.Fatalf("observer deadline misclassified as transport failure: %v", err)
	}
	if reads != 1 || elapsed > time.Second || got.Transport.AgentSessionID != o.Transport.AgentSessionID || got.Delivery.Phase != o.Delivery.Phase {
		t.Fatalf("window escaped or identity/delivery changed: elapsed=%s reads=%d got=%+v", elapsed, reads, got)
	}
	if !strings.Contains(got.NextAction.Message, "workflow execute follow "+o.RunID) {
		t.Fatalf("missing supported observation continuation: %+v", got.NextAction)
	}
}

func TestAL09FollowRequestTimeoutContinuation(t *testing.T) {
	for _, kind := range []string{"transient", "persistent", "changed_session", "actual_stop", "normal_completion"} {
		t.Run(kind, func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			o := workflowTestStart(t, f.workflowFixture)
			// Only the synthetic transport is delayed; no provider action is retried.
			first := map[string]any{"delay_millis": 700}
			last := map[string]any{"agent_status": "idle"}
			if kind == "changed_session" {
				last["agent_session"] = map[string]any{"agent": "codex", "kind": "id", "value": "replacement"}
			}
			if kind == "persistent" {
				last["delay_millis"] = 700
			}
			phase := "stopped"
			if kind == "normal_completion" {
				phase = "completed"
				first = last
			}
			// Model an already preserved owner stop/completion, without creating
			// any acceptance or delivery qualification in this isolated test.
			if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error { s.Result.Delivery.Phase = phase; return nil }); err != nil {
				t.Fatal(err)
			}
			workflowTestModel(t, f.workflowFixture, map[string]any{"observations": []any{first, last}})
			f.D.HerdrTimeout = 300 * time.Millisecond
			var warnings bytes.Buffer
			old := logrus.StandardLogger().Out
			logrus.SetOutput(&warnings)
			defer logrus.SetOutput(old)
			before := workflowTestCalls(t, f.workflowFixture, "agent get")
			start := time.Now()
			got, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 3*time.Second)
			reads := workflowTestCalls(t, f.workflowFixture, "agent get") - before
			t.Logf("case=%s observer_budget=3s request_budget=300ms first_delay=700ms elapsed=%s reads=%d error=%v", kind, time.Since(start), reads, err)
			switch kind {
			case "persistent":
				if err == nil || got.Transport.State != "unknown" || reads != 2 || !strings.Contains(got.NextAction.Message, "workflow execute follow "+o.RunID) {
					t.Fatalf("persistent read failure claimed liveness or exceeded bound: %+v %v reads=%d", got, err, reads)
				}
			case "changed_session":
				if err == nil || !strings.Contains(err.Error(), "differs from the bound run") || got.Transport.State != "unknown" || reads != 2 {
					t.Fatalf("replacement session accepted: %+v %v", got, err)
				}
			default:
				if err != nil || got.Delivery.Phase != phase || got.Transport.State != "idle" || got.Transport.Observation != "fresh" {
					t.Fatalf("same-session result lost: %+v %v", got, err)
				}
			}
			if kind != "normal_completion" && (!strings.Contains(warnings.String(), o.RunID) || !strings.Contains(warnings.String(), "agent get")) {
				t.Fatalf("missing sourced warning: %s", warnings.String())
			}
			if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
				t.Fatal("read failure replayed provider action")
			}
		})
	}
}

func TestAL09FollowCandidateCLI(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := workflowTestStart(t, f.workflowFixture)
	binary := deliveryCLIBinary(t)
	workflowTestModel(t, f.workflowFixture, map[string]any{"delay_command": "agent get", "delay_millis": 1500})
	c := exec.Command(binary, "workflow", "execute", "follow", o.RunID, "--timeout", "1", "--format", "json")
	c.Dir = f.Parent
	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr
	start := time.Now()
	err := c.Run()
	exit := c.ProcessState.ExitCode()
	t.Logf("candidate=%s observer_budget=1s delay=1500ms elapsed=%s exit=%d stdout=%s stderr=%s", binary, time.Since(start), exit, stdout.String(), stderr.String())
	if err == nil || exit != 5 || !strings.Contains(stderr.String(), "observation timed out without stopping or restarting") {
		t.Fatalf("CLI deadline classification: %v %s", err, stderr.String())
	}
	if !json.Valid(stdout.Bytes()) {
		t.Fatalf("stdout is not JSON: %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), o.Transport.AgentSessionID) || !strings.Contains(stdout.String(), "workflow execute follow "+o.RunID) {
		t.Fatal("CLI lost identity or continuation")
	}
	if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("CLI replayed provider action")
	}
}

func TestAL09ReportFollowObservation(t *testing.T) {
	for _, kind := range []string{"short_remaining", "window", "transient", "persistent", "changed_session"} {
		t.Run(kind, func(t *testing.T) {
			f := workflowTestFixture(t)
			o := workflowTestAccept(t, f, workflowTestStart(t, f))
			o = workflowTestReport(t, f, o)
			first := map[string]any{"delay_millis": 700, "agent_status": "idle"}
			last := map[string]any{"agent_status": "idle"}
			window := 3 * time.Second
			f.D.HerdrTimeout = 300 * time.Millisecond
			if kind == "window" {
				window = 200 * time.Millisecond
				f.D.HerdrTimeout = 2 * time.Second
			}
			if kind == "short_remaining" {
				window = time.Nanosecond
			}
			if kind == "persistent" {
				last["delay_millis"] = 700
			}
			if kind == "changed_session" {
				last["agent_session"] = map[string]any{"agent": "codex", "kind": "id", "value": "replacement"}
			}
			workflowTestModel(t, f, map[string]any{"observations": []any{first, last}})
			before := workflowTestCalls(t, f, "agent get")
			start := time.Now()
			got, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, window)
			reads := workflowTestCalls(t, f, "agent get") - before
			t.Logf("case=%s window=%s request_budget=%s elapsed=%s reads=%d error=%v", kind, window, f.D.HerdrTimeout, time.Since(start), reads, err)
			switch kind {
			case "short_remaining":
				if err == nil || !strings.Contains(err.Error(), "observation timed out without stopping or restarting") || reads != 0 || got.Round.State != o.Round.State || got.Transport.ObservedAt != o.Transport.ObservedAt {
					t.Fatalf("expired report window probed transport or changed round: round=%s error=%v reads=%d", got.Round.State, err, reads)
				}
			case "window":
				if err == nil || !strings.Contains(err.Error(), "observation timed out without stopping or restarting") || strings.Contains(err.Error(), "Transport observation is unknown") || reads != 1 || got.Round.State != o.Round.State {
					t.Fatalf("report observer window misclassified: round=%s error=%v reads=%d", got.Round.State, err, reads)
				}
			case "transient":
				if err != nil || got.Round.State != "ready_for_review" || got.Transport.State != "idle" {
					t.Fatalf("report continuation failed: %+v %v", got, err)
				}
			case "persistent":
				if err == nil || got.Transport.State != "unknown" || reads != 2 {
					t.Fatalf("persistent report timeout claimed liveness: %+v %v", got, err)
				}
			case "changed_session":
				if err == nil || !strings.Contains(err.Error(), "differs from the bound run") || reads != 2 {
					t.Fatalf("report session drift hidden: %+v %v", got, err)
				}
			}
			if got.Transport.AgentSessionID != o.Transport.AgentSessionID || workflowTestCalls(t, f, "agent start") != 1 || workflowTestCalls(t, f, "agent prompt") != 1 {
				t.Fatal("report observation changed identity or replayed input")
			}
		})
	}
}
