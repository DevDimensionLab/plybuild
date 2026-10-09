//go:build darwin

package taskrun

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workerskills"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestWorkerPromptsContainExactPinnedRoleBeforeTaskInstruction(t *testing.T) {
	s := workflowState{Request: WorkflowRequest{Delivery: &DeliveryContract{Goal: FileBinding{Locator: "/goal.json"}}}, Result: WorkflowRun{Worker: &WorkerStatus{Iteration: 1, Pretest: &FileBinding{Locator: "/receipt.json"}}}}
	for _, row := range []struct {
		role      string
		review    bool
		iteration int
	}{{workerskills.RoleDeveloper, false, 1}, {workerskills.RoleCorrection, false, 2}, {workerskills.RoleReviewer, true, 2}} {
		s.Result.Worker.Iteration = row.iteration
		prompt, err := workerPromptWithSkills(s, "/output", row.review)
		if err != nil {
			t.Fatal(err)
		}
		bundle, err := workerskills.Bundle(row.role)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(prompt, bundle) || !strings.Contains(prompt[len(bundle):], "Read the immutable Goal at /goal.json") {
			t.Fatalf("role %s not bound before actual task instruction", row.role)
		}
	}
}

func TestWorkerProviderFailureOnlyUsesProviderErrorEnvelopes(t *testing.T) {
	now := time.Date(2026, 10, 9, 18, 0, 0, 0, time.UTC)
	for _, tt := range []struct{ name, raw, kind string }{
		{"limit", `{"type":"turn.failed","error":{"code":"usage_limit_reached","resets_at":"2026-10-09T19:00:00Z"}}` + "\n", "usage_limit"},
		{"duplicate limit", `{"type":"error","code":"usage_limit_reached","resets_at":"2026-10-09T19:00:00Z"}` + "\n" + `{"type":"turn.failed","error":{"message":"usage_limit_reached"}}` + "\n", "usage_limit"},
		{"actual message", `{"type":"error","message":"You've hit your usage limit. Try again later."}` + "\n", "usage_limit"},
		{"auth", `{"type":"turn.failed","error":{"code":"authentication_error"}}` + "\n", "authentication"},
		{"tool text", `{"type":"item.completed","item":{"type":"command_execution","aggregated_output":"usage_limit_reached"}}` + "\n", ""},
		{"answer text", `{"type":"item.completed","item":{"type":"agent_message","text":"You've hit your usage limit."}}` + "\n", ""},
		{"unknown", `{"type":"error","message":"server overloaded"}` + "\n", ""},
		{"truncated", `{"type":"error","code":"usage_limit_reached"}`, ""},
		{"recovered", `{"type":"error","code":"usage_limit_reached"}` + "\n" + `{"type":"turn.completed"}` + "\n", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := workerProviderFailure([]byte(tt.raw), now)
			if tt.kind == "" {
				if got != nil {
					t.Fatalf("false retry: %+v", got)
				}
				return
			}
			if got == nil || got.Kind != tt.kind {
				t.Fatalf("classification: %+v", got)
			}
			if strings.Contains(tt.name, "limit") && got.ResetAt != "2026-10-09T19:00:00Z" {
				t.Fatalf("reset lost: %+v", got)
			}
		})
	}
}

func TestWorkerUsageLimitPersistsAcrossRestartWithoutNewRun(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	now := time.Date(2026, 10, 9, 18, 0, 0, 0, time.UTC)
	f.D.Now = func() time.Time { return now }
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		w := s.Result.Worker
		w.Stage = "reviewer"
		w.Job = &WorkerJob{ID: "reviewer-0001", Kind: "reviewer", Directory: filepath.Join(workerDirectory(*s), "attempts", "reviewer-0001")}
		b, _, e := deliveryReadBinding(filepath.Join(w.Job.Directory, "intent.json"), 1<<20)
		w.Job.IntentHash = b.SHA256
		return e
	}); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	job := s.Result.Worker.Job
	raw := []byte(`{"type":"error","code":"usage_limit_reached","resets_at":"2026-10-09T19:00:00Z"}` + "\n")
	stream := filepath.Join(job.Directory, "injected-limit.jsonl")
	if err = f.D.writeOnce(stream, raw); err != nil {
		t.Fatal(err)
	}
	r := workerAttemptResult{RunID: o.RunID, JobID: job.ID, IntentSHA256: job.IntentHash, FinishedAt: now.Format(time.RFC3339Nano), Exit: ptr(1), Error: "provider limit", Stdout: FileBinding{stream, hash(raw)}, Stderr: FileBinding{stream, hash(raw)}, Process: &Process{State: "exited", ExitCode: ptr(1), Quiescence: ptr(true)}, ProviderFailure: workerProviderFailure(raw, now)}
	if err = f.D.replaceValue(filepath.Join(job.Directory, "result.json"), r); err != nil {
		t.Fatal(err)
	}
	if err = workerStep(f.D, s); err != nil {
		t.Fatal(err)
	}
	// A new controller reads only persisted Run state and the original result.
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	w := s.Result.Worker
	if w.Phase != "waiting" || w.ProviderRetries != 1 || w.RetryNotBefore != "2026-10-09T19:00:00Z" || w.Job.ID != job.ID || w.Iteration != 1 || s.Result.RunID != o.RunID {
		t.Fatalf("wait lost native identity: %+v", w)
	}
	if err = workerContinueAfterWait(f.D, s); err == nil {
		t.Fatal("retry before reset")
	}
	now = now.Add(time.Hour)
	if err = workerContinueAfterWait(f.D, s); err != nil {
		t.Fatal(err)
	}
	s, err = workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Result.Worker.Job != nil || s.Result.Worker.Stage != "reviewer" || s.Result.Worker.ProviderRetries != 1 || s.Request.Runtime.Model != f.R.Runtime.Model || s.Result.Worker.Iteration != 1 {
		t.Fatalf("retry changed work or model: %+v", s.Result.Worker)
	}
	if _, err = workflowBound(r.Stdout, providerStreamLimit); err != nil {
		t.Fatal("retry lost old evidence:", err)
	}
	// Bounds and positive completion cannot be bypassed by classifying an error.
	s.Result.Worker.Job = job
	s.Result.Worker.ProviderRetries = 20
	if err = workerWaitForProvider(f.D, s, r); err == nil || !strings.Contains(err.Error(), "bound") {
		t.Fatalf("missing retry bound: %v", err)
	}
	r.Process.Quiescence = ptr(false)
	if err = workerWaitForProvider(f.D, s, r); err == nil {
		t.Fatal("unknown process retried")
	}
}

func TestWorkerNativeReleasedPreparationBlocksResumeBeforeMarker(t *testing.T) {
	f, o := newWorkerDeliveryFixture(t)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = workerPreparationCurrent(f.D, s); err != nil {
		t.Fatal(err)
	}
	prep, err := workerPreparation(f.D, s)
	if err != nil {
		t.Fatal(err)
	}
	target := workspace.QueueTargetInput{ProjectID: prep.Target.ProjectID, RepoID: prep.Target.RepoID, EpicID: prep.Target.EpicID}
	q, err := workspace.ListTaskQueue(f.D.Workspace, target, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.CloseTaskQueue(f.D.Workspace, target, s.Request.PreparationID, q.Revision, "disposable AW queue release crash fixture", "release"); err != nil {
		t.Fatal(err)
	}
	// Simulated crash after the real native effect but before the worker marker.
	if s.Result.Worker.QueueReleased {
		t.Fatal("fixture wrote marker before crash")
	}
	if err = workerPreparationCurrent(f.D, s); err == nil || !strings.Contains(err.Error(), "worker_preparation_released") {
		t.Fatal(fmt.Sprintf("released native preparation resumed: %v", err))
	}
}
