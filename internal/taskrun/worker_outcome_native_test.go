//go:build darwin

package taskrun_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/weekrun"
)

func outcomeAdapter(f taskrun.WorkerOutcomeFixture) (weekrun.NativeAdapter, weekrun.Contract) {
	return weekrun.NativeAdapter{Dependencies: f.Dependencies}, weekrun.Contract{SchemaVersion: 1, Workspace: f.Root, ReturnWorktree: f.Parent, ReturnRef: "refs/heads/epic", ProjectID: "ply", RepoID: "ply", EpicID: "epic"}
}

func TestAR02NativeOutcomeWordingAndRelease(t *testing.T) {
	for _, kind := range []string{"exhausted", "review_blocked", "developer_blocked"} {
		t.Run(kind, func(t *testing.T) {
			f := taskrun.NewWorkerOutcomeFixture(t, kind)
			n, c := outcomeAdapter(f)
			before, e := n.Observe(c, f.RunID)
			if e != nil || !before.Releasable {
				t.Fatalf("concrete outcome has no advice: %+v %v", before, e)
			}
			f.EditWorker(func(w *taskrun.WorkerStatus) { w.Reason = "No further corrections authorized; preserve the candidate." })
			after, e := n.Observe(c, f.RunID)
			if e != nil || !after.Releasable || before.Phase != after.Phase {
				t.Fatalf("wording changed current advice: %+v %v", after, e)
			}
			q, e := n.Queue(c)
			if e != nil {
				t.Fatal(e)
			}
			for i := 0; i < 2; i++ {
				if e = n.Release(c, f.RunID, "task"); e != nil {
					t.Fatal(e)
				}
			}
			released, e := n.Queue(c)
			if e != nil || released.Current != nil || released.Revision != q.Revision+1 {
				t.Fatalf("queue effect repeated or absent: %+v %v", released, e)
			}
			cmd := exec.Command("git", "rev-parse", "HEAD")
			cmd.Dir = f.Source
			raw, e := cmd.Output()
			if e != nil || strings.TrimSpace(string(raw)) != f.OID {
				t.Fatal("release changed source", e)
			}
			if _, e = taskrun.StartWorker(f.Dependencies, f.Root, f.RunID, true); e == nil {
				t.Fatal("released worker resumed")
			}
		})
	}
}

func TestAR02NativeUnknownCannotBecomeConcreteByWording(t *testing.T) {
	f := taskrun.NewWorkerOutcomeFixture(t, "unknown")
	n, c := outcomeAdapter(f)
	f.EditWorker(func(w *taskrun.WorkerStatus) { w.Reason = "worker reached its explicit limit of 3 developer attempts" })
	o, e := n.Observe(c, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if o.Releasable {
		t.Error("unknown effect became releasable by display text")
	}
	if e = n.Release(c, f.RunID, "task"); e == nil {
		t.Error("unknown effect released the real native queue")
	}
	q, e := n.Queue(c)
	if e != nil || q.Current == nil {
		t.Fatal("unknown effect lost queue ownership", e)
	}
}

func TestAR02NativeAdviceIsNotReleaseAuthority(t *testing.T) {
	for _, mutation := range []string{"owner", "resume", "process", "provider", "retry", "paused", "stage"} {
		t.Run(mutation, func(t *testing.T) {
			f := taskrun.NewWorkerOutcomeFixture(t, "exhausted")
			n, c := outcomeAdapter(f)
			o, e := n.Observe(c, f.RunID)
			if e != nil || !o.Releasable {
				t.Fatal(o, e)
			}
			switch mutation {
			case "owner":
				unlock := f.HoldOwner()
				defer unlock()
			case "resume":
				if e = os.WriteFile(oPath(f)+"/worker/control.json", []byte(`{"run_id":"`+f.RunID+`","action":"resume","requested_at":"2026-10-19T00:00:00Z"}`), 0600); e != nil {
					t.Fatal(e)
				}
			case "process":
				f.EditResult(func(r map[string]any) { r["process"].(map[string]any)["quiescence"] = false })
			case "provider":
				f.EditResult(func(r map[string]any) { r["provider_failure"] = map[string]any{"kind": "authentication"} })
			case "retry":
				f.EditWorker(func(w *taskrun.WorkerStatus) { w.RetryNotBefore = "2026-10-20T00:00:00Z" })
			case "paused":
				f.EditWorker(func(w *taskrun.WorkerStatus) { w.Phase = "paused" })
			case "stage":
				f.EditWorker(func(w *taskrun.WorkerStatus) { w.Stage = "integrate" })
			}
			if e = n.Release(c, f.RunID, "task"); e == nil {
				t.Fatal("cached advice authorized release after", mutation)
			}
			q, e := n.Queue(c)
			if e != nil || q.Current == nil {
				t.Fatal("lost queue", e)
			}
		})
	}
}

func oPath(f taskrun.WorkerOutcomeFixture) string {
	return f.Root + "/.ply/workflow-runs/runs/" + f.RunID
}

func TestAR02NativeReleaseLostReply(t *testing.T) {
	f := taskrun.NewWorkerOutcomeFixture(t, "exhausted")
	n, c := outcomeAdapter(f)
	q, e := n.Queue(c)
	if e != nil {
		t.Fatal(e)
	}
	if e = f.LoseReleaseReply(); e == nil || !strings.Contains(e.Error(), "response lost") {
		t.Fatal(e)
	}
	for i := 0; i < 2; i++ {
		if e = n.Release(c, f.RunID, "task"); e != nil {
			t.Fatal(e)
		}
	}
	after, e := n.Queue(c)
	if e != nil || after.Current != nil || after.Revision != q.Revision+1 {
		t.Fatal("lost reply replayed queue effect", after, e)
	}
	if _, e = taskrun.StartWorker(f.Dependencies, f.Root, f.RunID, true); e == nil {
		t.Fatal("lost marker allowed new writer")
	}
}

func TestAR02HistoricalOutcomeReadPreservesBytes(t *testing.T) {
	f := taskrun.NewWorkerOutcomeFixture(t, "exhausted")
	n, c := outcomeAdapter(f)
	f.EditWorker(func(w *taskrun.WorkerStatus) { w.Outcome = "" })
	paths := []string{oPath(f) + "/state.json", f.Root + "/.ply/workflow-runs/requests/" + f.RunID + ".json"}
	saved := make([][]byte, len(paths))
	for i, p := range paths {
		var e error
		saved[i], e = os.ReadFile(p)
		if e != nil {
			t.Fatal(e)
		}
	}
	for i := 0; i < 2; i++ {
		o, e := n.Observe(c, f.RunID)
		if e != nil || !o.Releasable {
			t.Fatal("historical native failure lost meaning", o, e)
		}
	}
	for i, p := range paths {
		raw, e := os.ReadFile(p)
		if e != nil || !bytes.Equal(saved[i], raw) {
			t.Fatal("historical read rewrote bytes", p, e)
		}
	}
	if e := n.Release(c, f.RunID, "task"); e != nil {
		t.Fatal("historical failure cannot use original release policy", e)
	}
}

// This probe requires the exact private candidate in acceptance. Providers are
// modeled fixture facts; all reads and resume rejection use the public CLI.
func TestAR02OutcomePublicCLI(t *testing.T) {
	binary := os.Getenv("PLY_BIN")
	if binary == "" {
		t.Skip("set PLY_BIN to the exact private candidate")
	}
	for _, kind := range []string{"exhausted", "unknown"} {
		t.Run(kind, func(t *testing.T) {
			f := taskrun.NewWorkerOutcomeFixture(t, kind)
			n, c := outcomeAdapter(f)
			f.EditWorker(func(w *taskrun.WorkerStatus) { w.Reason = "developer blocked: display-only wording" })
			cmd := exec.Command(binary, "workflow", "worker", "show", f.RunID, "--workspace", f.Root, "--format", "json")
			cmd.Dir = f.Source
			raw, e := cmd.Output()
			if e != nil {
				t.Fatal("public worker show", e)
			}
			var run taskrun.WorkflowRun
			if e = json.Unmarshal(raw, &run); e != nil {
				t.Fatal(e)
			}
			expected := kind == "exhausted"
			if run.Worker == nil || run.Worker.QueueReleaseAdvice() != expected || run.Worker.Outcome == "" {
				t.Fatalf("CLI lost preserved outcome: %s", raw)
			}
			observed, e := n.Observe(c, f.RunID)
			if e != nil || observed.Releasable != expected {
				t.Fatal("CLI/adapter disagreement", observed, e)
			}
			if !expected {
				return
			}
			if e = n.Release(c, f.RunID, "task"); e != nil {
				t.Fatal(e)
			}
			cmd = exec.Command(binary, "workflow", "worker", "resume", f.RunID, "--workspace", f.Root, "--format", "json")
			cmd.Dir = f.Source
			raw, e = cmd.CombinedOutput()
			if e == nil || !strings.Contains(string(raw), "worker_preparation_released") {
				t.Fatalf("public CLI resumed released failure: %s %v", raw, e)
			}
		})
	}
}

func TestAR02NativeStoppedFailureKeepsExistingRecoveryPolicy(t *testing.T) {
	f := taskrun.NewWorkerOutcomeFixture(t, "exhausted")
	n, c := outcomeAdapter(f)
	f.EditWorker(func(w *taskrun.WorkerStatus) { w.Phase = "stopped" })
	o, e := n.Observe(c, f.RunID)
	if e != nil || o.Releasable {
		t.Fatal("stopped recovery became automatic Week advice", o, e)
	}
	if e = n.Release(c, f.RunID, "task"); e != nil {
		t.Fatal("stopped concrete failure lost native recovery", e)
	}
}

func TestAR02NativeReleasedFailureRetainsReconciliationAdvice(t *testing.T) {
	f := taskrun.NewWorkerOutcomeFixture(t, "exhausted")
	n, c := outcomeAdapter(f)
	// The effect and Worker marker completed, but the Week caller lost the reply.
	if e := n.Release(c, f.RunID, "task"); e != nil {
		t.Fatal(e)
	}
	o, e := n.Observe(c, f.RunID)
	if e != nil || !o.Releasable {
		t.Fatal("completed release can no longer be reconciled by the same Week", o, e)
	}
	if e = n.Release(c, f.RunID, "task"); e != nil {
		t.Fatal("reconciliation repeated a completed effect", e)
	}
}
