package weekrun

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Native APIs historically resolve the current physical workspace. A dedicated
// CLI/controller performs these short operations serially from the frozen return
// worktree, then restores cwd; providers themselves run in separate processes.
var cwdMu sync.Mutex

func at(c Contract, fn func() error) error {
	cwdMu.Lock()
	defer cwdMu.Unlock()
	old, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.Chdir(c.ReturnWorktree); err != nil {
		return err
	}
	defer os.Chdir(old)
	return fn()
}

type NativeAdapter struct{ Dependencies taskrun.Dependencies }

func New(d taskrun.Dependencies) Manager {
	return Manager{Native: NativeAdapter{d}, Now: time.Now, Executable: d.Executable, Launch: launchDetached}
}
func (n NativeAdapter) Queue(c Contract) (q workspace.WorkspaceTaskQueueReadback, err error) {
	err = at(c, func() error {
		var e error
		q, e = workspace.ListTaskQueue(n.Dependencies.Workspace, c.Target(), false)
		return e
	})
	return
}
func workerContract(c Contract) *taskrun.WorkerContract {
	w := &taskrun.WorkerContract{SchemaVersion: 1, MaxAttempts: c.MaxAttempts, AgentTimeoutSeconds: c.AgentTimeoutSeconds, MaxProviderRetries: c.MaxProviderRetries}
	if c.ProductQuestions != nil {
		w.SchemaVersion = 2
		v := *c.ProductQuestions
		w.ProductQuestions = &v
	}
	return w
}
func (n NativeAdapter) execute(c Contract, s Selection, control taskrun.Executable, check bool) (out taskexecute.Result, err error) {
	err = at(c, func() error {
		d := n.Dependencies
		d.Executable = func() (string, error) { return control.Path, nil }
		var e error
		var contract *taskrun.ControlContract
		if c.RuntimeControl != nil {
			contract = taskrun.NewControlContract("week", contractID(c), hash(canonical(c)))
		}
		out, e = taskexecute.Execute(d, taskexecute.Input{Target: c.Target(), SpecID: s.Goal.SpecID, Check: check, Worker: workerContract(c), ControlContract: contract, ExpectedTaskID: s.TaskID, ExpectedGoal: &s.Goal})
		return e
	})
	return
}
func (n NativeAdapter) Preview(c Contract, s Selection, control taskrun.Executable) error {
	p, err := n.execute(c, s, control, true)
	if err != nil {
		return err
	}
	if p.Goal == nil || p.Goal.Workspace != c.Workspace || p.Goal.Target.ParentLocator != c.ReturnWorktree || p.Goal.Target.ParentRef != c.ReturnRef || p.Goal.Goal.TaskID != s.TaskID || p.Goal.Goal.Goal != s.Goal {
		return fmt.Errorf("native preview differs from frozen selection")
	}
	if p.Runtime != nil && (p.Runtime.Runtime.Provider != "codex" || p.Runtime.Runtime.Model != c.Model || p.Runtime.ReasoningEffort != c.Effort) {
		return fmt.Errorf("native preview differs from frozen model/effort")
	}
	return nil
}
func (n NativeAdapter) Start(c Contract, s Selection, control taskrun.Executable, check bool) (Run, error) {
	out, err := n.execute(c, s, control, check)
	r := Run{}
	if out.Run != nil && out.Run.RunID != "" {
		r = fromNative(*out.Run)
	}
	if out.Goal != nil {
		r.BaseOID = out.Goal.ParentOID
		r.BaseRevision = out.Goal.BaseRevision
	}
	return r, err
}
func fromNative(r taskrun.WorkflowRun) Run {
	out := Run{ID: r.RunID, Completed: r.Delivery != nil && r.Delivery.Phase == "completed"}
	if r.Worker != nil {
		w := r.Worker
		out.Phase = w.Phase
		out.Stage = w.Stage
		out.OwnerState = w.OwnerState
		out.ControlAction = w.ControlAction
		out.Reason = w.Reason
		out.RetryNotBefore = w.RetryNotBefore
		out.LastProviderFailure = w.LastProviderFailure
		out.Releasable = w.Phase == "blocked" && (strings.HasPrefix(w.Reason, "worker reached its explicit limit") || strings.HasPrefix(w.Reason, "developer blocked:") || strings.HasPrefix(w.Reason, "reviewer blocked:") || strings.HasPrefix(w.Reason, "review did not give an unambiguous pass: blocked;"))
	}
	return out
}
func (n NativeAdapter) Observe(c Contract, id string) (out Run, err error) {
	err = at(c, func() error {
		r, e := taskrun.ShowWorker(n.Dependencies, c.Workspace, id)
		out = fromNative(r)
		if e == nil {
			path := filepath.Join(r.Paths.RunRoot, "worker", "supervisor.lock")
			if _, se := os.Lstat(path); os.IsNotExist(se) {
				out.SupervisorState = "not_started"
			} else if se != nil {
				out.SupervisorState = "unknown"
			} else {
				f, le := lock(path)
				if le == nil {
					f.Close()
					out.SupervisorState = "unowned"
				} else if errors.Is(le, errOwned) {
					out.SupervisorState = "lock_held"
				} else {
					out.SupervisorState = "unknown"
				}
			}
		}
		return e
	})
	return
}
func (n NativeAdapter) Resume(c Contract, id string, explicit bool) (out Run, err error) {
	err = at(c, func() error {
		r, e := taskrun.StartWorker(n.Dependencies, c.Workspace, id, explicit)
		out = fromNative(r)
		return e
	})
	return
}
func (n NativeAdapter) Control(c Contract, id, action string) (out Run, err error) {
	err = at(c, func() error {
		r, e := taskrun.ControlWorker(n.Dependencies, c.Workspace, id, action)
		out = fromNative(r)
		return e
	})
	return
}
func (n NativeAdapter) Release(c Contract, id string, task workspace.TaskID) error {
	return at(c, func() error {
		d := n.Dependencies
		// A controller may die after the native queue effect but before our local
		// observation. Read the registered preparation disposition before retrying.
		reg, err := d.Workspace.WorkItems.Snapshot(c.Workspace)
		if err != nil {
			return err
		}
		preparation := ""
		for _, p := range reg.TaskPreparations {
			if p.Plan.TaskID == task {
				if preparation != "" {
					return fmt.Errorf("ambiguous native preparation")
				}
				preparation = p.ID
			}
		}
		if preparation == "" {
			return fmt.Errorf("native Task has no registered preparation")
		}
		readback, err := workspace.ShowTaskPreparation(workspace.WithTaskContentScope(d.Workspace, task), preparation)
		if err != nil {
			return err
		}
		if readback.Target.ProjectID != c.ProjectID || readback.Target.RepoID != c.RepoID || readback.Target.EpicID != c.EpicID {
			return fmt.Errorf("native preparation target differs")
		}
		if readback.Disposition == "released" {
			return nil
		}
		return taskrun.WithParkedWorker(d, c.Workspace, id, func() error {
			q, err := workspace.ListTaskQueue(d.Workspace, c.Target(), false)
			if err != nil {
				return err
			}
			if err = validateQueue(c, q); err != nil {
				return err
			}
			if q.Current == nil || q.Current.TaskID != task || q.Current.PreparationID != preparation {
				return fmt.Errorf("native current preparation does not belong to the blocked Task")
			}
			_, err = workspace.CloseTaskQueue(d.Workspace, c.Target(), preparation, q.Revision, "Selected week: preserve positively completed blocked work and continue independent authorized Goals", "release")
			return err
		})
	})
}
