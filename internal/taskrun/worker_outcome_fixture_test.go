//go:build darwin

package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Exported only into the external test binary so tests can cross both the
// actual Week adapter and Worker seam using a real disposable native queue.
type WorkerOutcomeFixture struct {
	Dependencies                     Dependencies
	Root, RunID, Parent, Source, OID string
	EditWorker                       func(func(*WorkerStatus))
	EditResult                       func(func(map[string]any))
	HoldOwner                        func() func()
	LoseReleaseReply                 func() error
}

func NewWorkerOutcomeFixture(t *testing.T, kind string) WorkerOutcomeFixture {
	t.Helper()
	f, o := newWorkerDeliveryFixture(t)
	if err := workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		w := s.Result.Worker
		w.Stage = "reviewer"
		w.Iteration = s.Request.Worker.MaxAttempts
		dir := filepath.Join(workerDirectory(*s), "attempts", "reviewer-0001")
		in, _, err := deliveryReadBinding(filepath.Join(dir, "intent.json"), 1<<20)
		if err != nil {
			return err
		}
		w.Job = &WorkerJob{ID: "reviewer-0001", Kind: "reviewer", Directory: dir, IntentHash: in.SHA256}
		var result workerAttemptResult
		if err = readValue(filepath.Join(dir, "result.json"), 1<<20, &result); err != nil {
			return err
		}
		decision := "changes_requested"
		if kind == "review_blocked" {
			decision = "blocked"
		}
		answer, err := workflowKeep(f.D, filepath.Join(dir, "concrete-answer.json"), workerReviewAnswer{Decision: decision, Summary: "Modeled concrete review failure", Findings: []string{"invalid input accepted"}})
		if err != nil {
			return err
		}
		result.Answer = &answer
		if kind == "unknown" {
			result.Process.Quiescence = ptr(false)
		}
		if kind == "developer_blocked" {
			w.Stage = "developer"
			dir = filepath.Join(workerDirectory(*s), "attempts", "developer-0001")
			intent, e := workflowKeep(f.D, filepath.Join(dir, "intent.json"), workerAttemptIntent{RunID: o.RunID, JobID: "developer-0001", Kind: "developer", Generation: w.Generation, Argv: []string{}, Environment: []string{}})
			if e != nil {
				return e
			}
			w.Job = &WorkerJob{ID: "developer-0001", Kind: "developer", Directory: dir, IntentHash: intent.SHA256}
			result.JobID, result.IntentSHA256 = w.Job.ID, intent.SHA256
			answer, e = workflowKeep(f.D, filepath.Join(dir, "answer.json"), map[string]string{"status": "blocked", "summary": "Modeled developer block", "block_reason": "No authorized solution for this requirement"})
			if e != nil {
				return e
			}
			result.Answer = &answer
		}
		return f.D.replaceValue(filepath.Join(dir, "result.json"), result)
	}); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	cause := workerStep(f.D, s)
	if cause == nil {
		t.Fatal("fixture did not produce a concrete failed review")
	}
	if err = workerBlocked(f.D, f.R.WorkspaceRoot, o.RunID, cause); err != nil {
		t.Fatal(err)
	}
	f.D.WorkerOwner.lock.Close()
	out := WorkerOutcomeFixture{Dependencies: f.D, Root: f.R.WorkspaceRoot, RunID: o.RunID, Parent: f.Parent, Source: s.Observed.Target.WorktreeLocator, OID: gitOutput(t, s.Observed.Target.WorktreeLocator, "rev-parse", "HEAD")}
	out.EditWorker = func(edit func(*WorkerStatus)) {
		t.Helper()
		if e := workflowUpdate(f.D, out.Root, out.RunID, func(s *workflowState) error { edit(s.Result.Worker); return nil }); e != nil {
			t.Fatal(e)
		}
	}
	out.EditResult = func(edit func(map[string]any)) {
		t.Helper()
		path := filepath.Join(s.Result.Worker.Job.Directory, "result.json")
		var result map[string]any
		if e := readValue(path, 1<<20, &result); e != nil {
			t.Fatal(e)
		}
		edit(result)
		if e := f.D.replaceValue(path, result); e != nil {
			t.Fatal(e)
		}
	}
	out.HoldOwner = func() func() {
		t.Helper()
		lock, e := workerLock(filepath.Join(workerDirectory(s), "owner.lock"), false)
		if e != nil {
			t.Fatal(e)
		}
		return func() { lock.Close() }
	}
	out.LoseReleaseReply = func() error {
		return WithParkedWorker(f.D, out.Root, out.RunID, func() error {
			target := workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}
			q, e := workspace.ListTaskQueue(f.D.Workspace, target, false)
			if e != nil {
				return e
			}
			if _, e = workspace.CloseTaskQueue(f.D.Workspace, target, s.Request.PreparationID, q.Revision, "Disposable lost release reply fixture", "release"); e != nil {
				return e
			}
			return fmt.Errorf("release response lost after the actual native queue effect")
		})
	}
	// Keep original answer/result evidence available after queue release.
	t.Cleanup(func() {
		if _, e := os.Stat(filepath.Join(s.Result.Worker.Job.Directory, "result.json")); e != nil {
			t.Error(e)
		}
	})
	return out
}
