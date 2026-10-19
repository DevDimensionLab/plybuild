//go:build !windows

package taskrun

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func workerLock(path string, create bool) (*os.File, error) {
	flags := os.O_RDWR | unix.O_NOFOLLOW
	if create {
		flags |= os.O_CREATE
	}
	f, err := physicalOpenFile(path, flags, 0600)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0077 != 0 {
		f.Close()
		return nil, fmt.Errorf("invalid worker lock")
	}
	if err = unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB); err != nil {
		f.Close()
		return nil, err
	}
	// Close only: explicit LOCK_UN would revoke the same lock inherited by an
	// attempt runner while its parent/controller exits.
	return f, nil
}

func ShowWorker(d Dependencies, root, id string) (WorkflowRun, error) {
	s, err := workflowRead(root, id)
	if err != nil {
		return WorkflowRun{}, err
	}
	if err = workerValidateIdentity(s); err != nil {
		return s.Result, err
	}
	o, err := WorkflowShow(d, root, id)
	if err != nil {
		return o, err
	}
	lock, probe := workerLock(filepath.Join(workerDirectory(s), "owner.lock"), false)
	if probe == nil {
		lock.Close()
		o.Worker.OwnerState = "unowned"
	} else if errors.Is(probe, unix.EWOULDBLOCK) {
		o.Worker.OwnerState = "lock_held"
	} else if os.IsNotExist(probe) {
		o.Worker.OwnerState = "not_started"
	} else {
		o.Worker.OwnerState = "unknown"
	}
	return o, nil
}

// StartWorker launches a detached native controller. The supervisor can also
// invoke ServeWorker directly. No state written by an agent grants this entry.
func StartWorker(d Dependencies, root, id string, resume bool) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	s, err := workflowRead(root, id)
	if err != nil {
		return WorkflowRun{}, err
	}
	if err = workerValidateIdentity(s); err != nil {
		return s.Result, err
	}
	if s.Result.Delivery.Phase == "completed" {
		return ShowWorker(d, root, id)
	}
	if err = workerPreparationCurrent(d, s); err != nil {
		return s.Result, err
	}
	if err = verifyExecutable(s.Request.Runtime.PlyExecutable); err != nil {
		return s.Result, err
	}
	dir := workerDirectory(s)
	if err = privateDir(dir); err != nil {
		return s.Result, err
	}
	if resume {
		if err = workerControlWrite(d, s, "resume"); err != nil {
			return s.Result, err
		}
		if s.Request.ControlContract != nil && s.Request.ControlContract.ParentKind == "standalone" {
			return ShowWorker(d, root, id)
		}
	}
	// Multiple clients may race here. Only the controller that acquires the
	// owner lock may reserve a provider attempt or mutate the Task.
	log, err := physicalOpenFile(filepath.Join(dir, "controller.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND|unix.O_NOFOLLOW, 0600)
	if err != nil {
		return s.Result, err
	}
	defer log.Close()
	// Controllers stay outside the Task source they finally remove. Serve
	// selects its own cwd; only callbacks and attempts work in the source.
	cmd := exec.Command(s.Request.Runtime.PlyExecutable.Path, "workflow", "worker", "_supervise", id, "--workspace", root, "--format", "json")
	cmd.Dir, cmd.Stdout, cmd.Stderr = s.Request.WorkspaceRoot, log, log
	cmd.Stdin = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err = cmd.Start(); err != nil {
		return s.Result, err
	}
	if err = cmd.Process.Release(); err != nil {
		return s.Result, err
	}
	return ShowWorker(d, root, id)
}

func workerControlWrite(d Dependencies, s workflowState, action string) error {
	if s.Request.ControlContract != nil {
		if s.Request.ControlContract.ParentKind == "standalone" {
			return workerLegacyBoundControl(d, s, action)
		}
		return WithWorkerParentGuard(s.Request.WorkspaceRoot, s.Request.ControlContract, func() error {
			if err := CheckParentAction(s.Request.WorkspaceRoot, s.Request.ControlContract, action, controlNow(d)); err != nil {
				return err
			}
			return withStore(s.Request.WorkspaceRoot, func() error {
				return d.replaceValue(filepath.Join(workerDirectory(s), "control.json"), workerControl{s.Result.RunID, action, ControlTimestamp(controlNow(d))})
			})
		})
	}
	return withStore(s.Request.WorkspaceRoot, func() error {
		return d.replaceValue(filepath.Join(workerDirectory(s), "control.json"), workerControl{s.Result.RunID, action, d.Now().UTC().Format(time.RFC3339Nano)})
	})
}

func ControlWorker(d Dependencies, root, id, action string) (WorkflowRun, error) {
	if action != "pause" && action != "stop" {
		return WorkflowRun{}, workflowError(2, "worker action must be pause or stop")
	}
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	s, err := workflowRead(root, id)
	if err != nil {
		return WorkflowRun{}, err
	}
	if err = workerValidateIdentity(s); err != nil {
		return s.Result, err
	}
	if s.Result.Delivery.Phase == "completed" {
		return ShowWorker(d, root, id)
	}
	if err = d.fault("worker_control_after_read"); err != nil {
		return s.Result, err
	}
	if s.Request.ControlContract != nil && s.Request.ControlContract.ParentKind == "standalone" {
		if err = workerLegacyBoundControl(d, s, action); err != nil {
			return s.Result, err
		}
		return ShowWorker(d, root, id)
	}
	if err = privateDir(workerDirectory(s)); err == nil {
		err = workerControlWrite(d, s, action)
	}
	if err != nil {
		return s.Result, err
	}
	// The phase read above may have become parked before publication. Always
	// wake the supervisor; its lock makes the already-active case harmless.
	if _, err = StartWorker(d, root, id, false); err != nil {
		return s.Result, err
	}
	o, err := ShowWorker(d, root, id)
	o.NextAction = WorkflowAction{"worker", action + " requested; an active attempt retains ownership until it has actually ended. Source files are kept."}
	return o, err
}

func workerStateUpdate(d Dependencies, root, id string, fn func(*workflowState) error) error {
	return workflowUpdate(d, root, id, func(s *workflowState) error {
		if err := workerLiveOwner(d, *s); err != nil {
			return err
		}
		if err := fn(s); err != nil {
			return err
		}
		s.Result.Worker.UpdatedAt = d.Now().UTC().Format(time.RFC3339Nano)
		return nil
	})
}

func workerBlocked(d Dependencies, root, id string, cause error) error {
	return workerStateUpdate(d, root, id, func(s *workflowState) error {
		s.Result.Worker.Reason = cause.Error()
		s.Result.Worker.Phase = "blocked"
		s.Result.NextAction = WorkflowAction{"human", "Worker blocked with preserved work: " + cause.Error()}
		return nil
	})
}

func ServeWorker(d Dependencies, root, id string) error {
	s, err := workflowRead(root, id)
	if err != nil {
		return err
	}
	if err = workerValidateIdentity(s); err != nil {
		return err
	}
	dir := workerDirectory(s)
	if err = privateDir(dir); err != nil {
		return err
	}
	lock, err := workerLock(filepath.Join(dir, "owner.lock"), true)
	if errors.Is(err, unix.EWOULDBLOCK) {
		return nil
	}
	if err != nil {
		return err
	}
	defer lock.Close()
	if err = os.Chdir(workerServeDirectory(s)); err != nil {
		return err
	}
	if s, err = workflowRead(root, id); err != nil {
		return err
	}
	if err = workerClaimOwner(&d, s, lock); err != nil {
		return err
	}
	for {
		s, err = workflowRead(root, id)
		if err != nil {
			return err
		}
		if consumed, e := workerConsumeBoundControl(d, s); e != nil {
			return workerBlocked(d, root, id, e)
		} else if consumed {
			continue
		}
		if handled, e := workerProductWaiting(d, s); handled {
			if e != nil {
				return workerBlocked(d, root, id, e)
			}
			if s.Result.Worker.Phase == "stopped" {
				return nil
			}
			time.Sleep(250 * time.Millisecond)
			continue
		} else if e != nil {
			return workerBlocked(d, root, id, e)
		}
		if s.Request.ControlContract != nil && s.Request.ControlContract.ParentKind == "week" {
			if err = workerSyncParentControl(d, s); err != nil {
				return workerBlocked(d, root, id, err)
			}
		}
		if s.Result.Delivery.Phase == "completed" {
			return workerStateUpdate(d, root, id, func(s *workflowState) error {
				s.Result.Worker.Phase = "completed"
				s.Result.Worker.Stage = "completed"
				s.Result.Worker.Job = nil
				s.Result.Worker.Reason = ""
				return nil
			})
		}
		var control workerControl
		controlErr := readValue(filepath.Join(dir, "control.json"), 64<<10, &control)
		if controlErr != nil && !os.IsNotExist(controlErr) {
			return workerBlocked(d, root, id, controlErr)
		}
		if controlErr == nil && control.RunID != id {
			return workerBlocked(d, root, id, fmt.Errorf("worker control belongs to another Run"))
		}
		if control.Action == "pause" || control.Action == "stop" {
			if s.Result.Worker.Job != nil {
				result, e := workerRunReserved(d, s)
				if e == nil && workerNativeKind(s.Result.Worker.Job.Kind) && !result.NativeCompleted {
					e = fmt.Errorf("native effect did not confirm completion; inspect the preserved receipt")
				}
				if e != nil {
					return workerBlocked(d, root, id, fmt.Errorf("%s remains requested; current attempt completion is unknown: %w", control.Action, e))
				}
			}
			return workerStateUpdate(d, root, id, func(s *workflowState) error {
				if control.Action == "pause" {
					s.Result.Worker.Phase = "paused"
				} else {
					s.Result.Worker.Phase = "stopped"
				}
				s.Result.NextAction = WorkflowAction{"human", "Work is preserved; use worker resume for this same Run."}
				return nil
			})
		}
		if s.Result.Worker.Phase == "blocked" && control.Action != "resume" {
			return nil
		}
		if s.Request.ControlContract != nil && (s.Result.Worker.Phase == "paused" || s.Result.Worker.Phase == "stopped") && control.Action != "resume" {
			return nil
		}
		if control.Action == "resume" {
			if err = workerConsumeResume(d, s, control); err != nil {
				return err
			}
			continue
		}
		if err = workerPreparationCurrent(d, s); err != nil {
			return workerBlocked(d, root, id, err)
		}
		if s.Result.Worker.RetryNotBefore != "" {
			due, e := workerRetryDue(d, s)
			if e != nil {
				return workerBlocked(d, root, id, e)
			}
			if !due {
				return nil
			}
			if e = workerContinueAfterWait(d, s); e != nil {
				return workerBlocked(d, root, id, e)
			}
			continue
		}
		if err = workerProbeCurrentGeneration(d, s); err != nil {
			return workerBlocked(d, root, id, err)
		}
		if err = workerStep(d, s); err != nil {
			if s.Request.ControlContract != nil {
				if e, ok := err.(*ControlFailure); ok && (e.Code == "control_in_progress" || e.Code == "control_unavailable") {
					continue
				}
			}
			return workerBlocked(d, root, id, err)
		}
	}
}

func workerConsumeResume(d Dependencies, s workflowState, expected workerControl) error {
	return workerStateUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
		path := filepath.Join(workerDirectory(*current), "control.json")
		var actual workerControl
		if err := readValue(path, 64<<10, &actual); err != nil {
			return err
		}
		if actual != expected {
			return nil
		}
		actual.Action = "continue"
		if err := d.replaceValue(path, actual); err != nil {
			return err
		}
		current.Result.Worker.Phase = "running"
		if current.Result.Worker.RetryNotBefore != "" {
			current.Result.Worker.Phase = "waiting"
		}
		current.Result.Worker.Reason = ""
		return nil
	})
}

// This small supervisor has no Task database. It restarts only the same native
// Run, waiting while an orphaned attempt still holds the ownership lock.
func SuperviseWorker(d Dependencies, root, id string) error {
	s, err := workflowRead(root, id)
	if err != nil {
		return err
	}
	if err = workerValidateIdentity(s); err != nil {
		return err
	}
	lock, err := workerLock(filepath.Join(workerDirectory(s), "supervisor.lock"), true)
	if errors.Is(err, unix.EWOULDBLOCK) {
		return nil
	}
	if err != nil {
		return err
	}
	defer lock.Close()
	for {
		s, err = workflowRead(root, id)
		if err != nil {
			return err
		}
		parked, e := workerSupervisorExit(d, root, id, lock)
		if e != nil {
			return e
		}
		if parked {
			return nil
		}
		if s.Result.Worker.Phase == "waiting" {
			if e = d.fault("worker_supervisor_waiting"); e != nil {
				return e
			}
			due, e := workerRetryDue(d, s)
			if e != nil {
				return e
			}
			var control workerControl
			if e = readValue(filepath.Join(workerDirectory(s), "control.json"), 64<<10, &control); e != nil && !os.IsNotExist(e) {
				return e
			}
			boundPending := false
			if s.Request.ControlContract != nil && s.Request.ControlContract.ParentKind == "standalone" {
				_, pending, e := ControlHistory(root, ControlTarget{"run", id})
				if e != nil {
					return e
				}
				boundPending = pending != nil
			}
			if !due && !boundPending && control.Action != "pause" && control.Action != "stop" && control.Action != "resume" {
				time.Sleep(time.Second)
				continue
			}
		}
		probe, e := workerLock(filepath.Join(workerDirectory(s), "owner.lock"), true)
		if errors.Is(e, unix.EWOULDBLOCK) {
			time.Sleep(time.Second)
			continue
		}
		if e != nil {
			return e
		}
		probe.Close()
		if err = verifyExecutable(s.Request.Runtime.PlyExecutable); err != nil {
			return err
		}
		if err = d.fault("worker_supervisor_before_serve"); err != nil {
			return err
		}
		cmd := exec.Command(s.Request.Runtime.PlyExecutable.Path, "workflow", "worker", "_serve", id, "--workspace", root, "--format", "json")
		cmd.Dir, cmd.Stdout, cmd.Stderr = s.Request.WorkspaceRoot, os.Stdout, os.Stderr
		if e = cmd.Run(); e != nil {
			fmt.Fprintln(os.Stderr, "worker controller ended:", e)
		}
		time.Sleep(time.Second)
	}
}

// Serialize the exit decision and lock release with control publication. A
// new control request either keeps this supervisor alive or sees its lock free;
// it cannot lose its replacement to a supervisor that already decided to exit.
func workerSupervisorExit(d Dependencies, root, id string, lock *os.File) (bool, error) {
	parked := false
	err := withStore(root, func() error {
		s, err := workflowRead(root, id)
		if err != nil {
			return err
		}
		var control workerControl
		if e := readValue(filepath.Join(workerDirectory(s), "control.json"), 64<<10, &control); e != nil && !os.IsNotExist(e) {
			return e
		}
		phase := s.Result.Worker.Phase
		if s.Request.ControlContract != nil && s.Request.ControlContract.ParentKind == "standalone" {
			_, pending, e := ControlHistory(root, ControlTarget{"run", id})
			if e != nil {
				return e
			}
			if pending != nil {
				return nil
			}
		}
		changeParkedControl := phase == "paused" && control.Action == "stop" || phase == "stopped" && control.Action == "pause"
		shouldExit := phase == "completed" || (phase == "blocked" || phase == "paused" || phase == "stopped") && control.Action != "resume" && !changeParkedControl
		if !shouldExit {
			return nil
		}
		if err = d.fault("worker_supervisor_before_exit"); err != nil {
			return err
		}
		if err = lock.Close(); err != nil {
			return err
		}
		parked = true
		return nil
	})
	return parked, err
}

// workerServeDirectory is the Task worktree while owner callbacks are open.
// After the automatic release only closeout remains, which removes the source
// and therefore runs from the workspace root.
func workerServeDirectory(s workflowState) string {
	if s.Result.Delivery != nil && s.Result.Delivery.OwnershipRelease != nil {
		return s.Request.WorkspaceRoot
	}
	return s.Observed.Target.WorktreeLocator
}
