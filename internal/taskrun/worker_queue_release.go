//go:build !windows

package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func workerPreparation(d Dependencies, s workflowState) (workspace.TaskPreparationReadback, error) {
	var out workspace.TaskPreparationReadback
	var draft struct {
		Binding struct {
			TaskID workspace.TaskID `json:"task_id"`
		} `json:"task_spec_binding"`
	}
	if err := json.Unmarshal(s.Request.HandoffDraft, &draft); err != nil {
		return out, err
	}
	if draft.Binding.TaskID == "" {
		return out, fmt.Errorf("worker task binding missing")
	}
	err := workspace.WithTaskSpecSnapshot(workspace.WithTaskContentScope(d.Workspace, draft.Binding.TaskID), s.Request.WorkspaceRoot, func(session *workspace.TaskSpecSession) error {
		var err error
		out, err = session.TaskRunPreparation(s.Request.PreparationID)
		return err
	})
	if err != nil {
		return out, err
	}
	p := out.Preparation
	if p == nil || p.ID != s.Request.PreparationID || digest(*p) != s.Request.PreparationSHA256 || p.Plan.TaskID != draft.Binding.TaskID || p.Plan.WorktreePath != s.Observed.Target.WorktreeLocator || p.Plan.Target.ParentLocator != s.Observed.Epic.WorktreeLocator || p.Plan.Target.ParentRef != s.Observed.Epic.Ref {
		return out, fmt.Errorf("worker native preparation identity differs")
	}
	return out, nil
}

func workerPreparationCurrent(d Dependencies, s workflowState) error {
	p, err := workerPreparation(d, s)
	if err != nil {
		return err
	}
	w := s.Result.Worker
	if !w.QueueReleased && p.Disposition == "current" {
		return nil
	}
	// Integration advances the native queue before the Run closeout receipt is
	// recorded. Only its exact reserved native recovery may cross that boundary;
	// never permit a new provider, commit or verification on advanced ownership.
	if !w.QueueReleased && p.Disposition == "advanced" && w.Stage == "integrate" && w.Job != nil && w.Job.Kind == "integrate" {
		a := s.Result.Delivery.Attempt
		if a != nil && a.Kind == "integration" && a.CandidateOID == w.CandidateOID && a.CandidateTree == w.CandidateTree {
			in, e := workerIntent(s, *w.Job)
			if e != nil {
				return e
			}
			return workerNativeIntentMatches(s, in)
		}
	}
	return fmt.Errorf("worker_preparation_released: native queue ownership is %s; no new writer is authorized", p.Disposition)
}

// WithParkedWorker serializes the permanent release of a concretely failed,
// quiescent Run with its worker. Usage limits, authentication, paused work and
// unknown effects are not permission to abandon its native queue ownership.
// The caller releases only this preparation through the native queue API.
func WithParkedWorker(d Dependencies, root, id string, release func() error) error {
	s, err := workflowRead(root, id)
	if err != nil {
		return err
	}
	if err = workerValidateIdentity(s); err != nil {
		return err
	}
	lock, err := workerLock(filepath.Join(workerDirectory(s), "owner.lock"), false)
	if err != nil {
		return fmt.Errorf("parked worker still owned or unknown: %w", err)
	}
	defer lock.Close()
	s, err = workflowRead(root, id)
	if err != nil {
		return err
	}
	d.WorkerOwner = &workerOwner{root: root, runID: id, generation: s.Result.Worker.Generation, lock: lock}
	prep, err := workerPreparation(d, s)
	if err != nil {
		return err
	}
	if prep.Disposition != "released" {
		w := s.Result.Worker
		if prep.Disposition != "current" || !w.QueueReleaseAdvice() {
			return fmt.Errorf("worker is not a releasable concrete failed attempt")
		}
		var control workerControl
		if e := readValue(filepath.Join(workerDirectory(s), "control.json"), 64<<10, &control); e != nil && !os.IsNotExist(e) {
			return e
		}
		if control.Action == "resume" {
			return fmt.Errorf("worker resume is pending")
		}
		result, e := workerRunReserved(d, s)
		if e != nil {
			return e
		}
		if result.ProviderFailure != nil {
			return fmt.Errorf("provider failures retain their native Run for recovery")
		}
		if err = release(); err != nil {
			return err
		}
		prep, err = workerPreparation(d, s)
		if err != nil {
			return err
		}
		if prep.Disposition != "released" {
			return fmt.Errorf("native queue did not release the exact preparation")
		}
	}
	// A crash before this marker is safe: every worker generation independently
	// checks the native preparation disposition before any new effect.
	return workerStateUpdate(d, root, id, func(current *workflowState) error {
		current.Result.Worker.QueueReleased = true
		current.Result.Worker.Phase = "blocked"
		return nil
	})
}
