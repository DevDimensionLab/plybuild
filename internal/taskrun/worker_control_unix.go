//go:build !windows

package taskrun

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"
)

type workerControlBackend struct{ d Dependencies }

func WorkerControls(d Dependencies) ControlService { return ControlService{d, workerControlBackend{d}} }
func ShowControl(d Dependencies, root, id string) (ControlBasis, error) {
	return WorkerControls(d).Show(root, id)
}
func ControlTargetForRun(root, id string) (ControlTarget, error) {
	s, e := workflowRead(root, id)
	if e != nil {
		return ControlTarget{}, controlSafeError("control_not_found")
	}
	c := s.Request.ControlContract
	if c != nil && c.ParentKind == "week" {
		if ValidateControlContract(c) != nil {
			return ControlTarget{}, controlSafeError("control_owner_unknown")
		}
		return ControlTarget{"week", *c.WeekID}, nil
	}
	return ControlTarget{"run", id}, nil
}
func controlLockObservation(path string) string {
	f, e := workerLock(path, false)
	if e == nil {
		f.Close()
		return "unowned"
	}
	if errors.Is(e, unix.EWOULDBLOCK) {
		return "lock_held"
	}
	if os.IsNotExist(e) {
		return "not_started"
	}
	return "unknown"
}
func (b workerControlBackend) ControlGuard(root, id string, fn func() error) error { return fn() }
func (b workerControlBackend) ControlPublicationGuard(root, expected string, fn func() error) error {
	return WithNativeControlSnapshot(b.d, root, func(current string) error {
		if current != expected {
			return controlSafeError("control_stale")
		}
		return fn()
	})
}
func (b workerControlBackend) ControlWake(root, id string) error {
	wake := false
	if err := workflowUpdate(b.d, root, id, func(s *workflowState) error {
		_, pending, err := ControlHistory(root, ControlTarget{"run", id})
		if err != nil || pending == nil {
			return err
		}
		if pending.Basis.OwnerGeneration == nil || s.Result.Worker.Generation != *pending.Basis.OwnerGeneration {
			return controlSafeError("control_outcome_unknown")
		}
		generation := s.Result.Worker.Generation
		s.ControlWake = &workerControlWake{pending.Sequence, generation, generation + 1}
		wake = true
		return nil
	}); err != nil {
		return err
	}
	if !wake {
		return nil
	}
	if b.d.ControlWake != nil {
		return b.d.ControlWake(root, id)
	}
	_, e := StartWorker(b.d, root, id, false)
	return e
}

// Called with the native owner lock, after observing state under that lock.
// Publication shares the TaskRun guard with apply. An owner already waiting to
// publish cannot acquire authority from a wake it had not observed.
func workerClaimOwner(d *Dependencies, observed workflowState, lock *os.File) error {
	root, id := observed.Request.WorkspaceRoot, observed.Result.RunID
	generation := 0
	err := workflowUpdate(*d, root, id, func(s *workflowState) error {
		w := s.Result.Worker
		s.ControlSuccessor = nil
		if s.Request.ControlContract != nil && s.Request.ControlContract.ParentKind == "standalone" {
			_, pending, err := ControlHistory(root, ControlTarget{"run", id})
			if err != nil {
				return err
			}
			link := s.ControlWake
			if pending != nil && link != nil && equal(link, observed.ControlWake) && link.Sequence == pending.Sequence && pending.Basis.OwnerGeneration != nil && link.Predecessor == *pending.Basis.OwnerGeneration && link.Predecessor == w.Generation && link.Successor == w.Generation+1 {
				copy := *link
				s.ControlSuccessor = &copy
			}
		}
		w.Generation++
		w.OwnerPID, w.OwnerStartedAt = os.Getpid(), ControlTimestamp(controlNow(*d))
		generation = w.Generation
		return nil
	})
	if err == nil {
		d.WorkerOwner = &workerOwner{root: root, runID: id, generation: generation, lock: lock}
	}
	return err
}

func workerControlGenerationAllowed(s workflowState, r ControlRecord) bool {
	if r.Basis.OwnerGeneration == nil {
		return false
	}
	generation := s.Result.Worker.Generation
	if generation == *r.Basis.OwnerGeneration {
		return true
	}
	link := s.ControlSuccessor
	return link != nil && link.Sequence == r.Sequence && link.Predecessor == *r.Basis.OwnerGeneration && link.Successor == link.Predecessor+1 && link.Successor == generation
}
func (b workerControlBackend) ControlBasis(root, id string) (ControlBasis, error) {
	o := NewControlBasis(root, ControlTarget{"run", id}, controlNow(b.d))
	if !ValidateControlTarget(o.Scope) || physical(root, false) != nil {
		return o, controlSafeError("control_invalid_arguments")
	}
	s, e := workflowRead(root, id)
	if e != nil {
		if os.IsNotExist(e) {
			return o, controlSafeError("control_not_found")
		}
		return o, controlSafeError("control_outcome_unknown")
	}
	o.RunID = &s.Result.RunID
	o.RunRequestHash = &s.Result.RequestSHA256
	o.RuntimeMode = controlString(s.Request.Runtime.Mode)
	o.Provider = controlString(s.Request.Runtime.Provider)
	o.Contract = s.Request.ControlContract
	o.FrozenControl = controlString(s.Request.Runtime.PlyExecutable.SHA256)
	var draft struct {
		Task struct {
			TaskID string `json:"task_id"`
			Spec   struct {
				Hash string `json:"manifest_sha256"`
			} `json:"spec"`
		} `json:"task_spec_binding"`
		Project struct {
			ProjectID string `json:"project_id"`
			RepoID    string `json:"repo_id"`
		} `json:"project_binding"`
	}
	_ = json.Unmarshal(s.Request.HandoffDraft, &draft)
	o.ProjectID = controlString(draft.Project.ProjectID)
	o.RepoID = controlString(draft.Project.RepoID)
	o.TaskID = controlString(draft.Task.TaskID)
	o.Goal = controlString(draft.Task.Spec.Hash)
	if s.Request.Delivery != nil && s.Request.Delivery.Agreement != nil {
		o.EpicID = controlString(string(s.Request.Delivery.Agreement.EpicID))
	}
	// Full private identities/policy feed the opaque revision; public projection
	// never exposes native paths, provider streams or callback material.
	w := s.Result.Worker
	if w == nil {
		o.StateRevision = digest(s.Request)
		o.Unavailable("runtime_unsupported")
		return o, nil
	}
	o.WorkerSession = &w.SessionID
	o.OwnerGeneration = &w.Generation
	o.Phase = w.Phase
	o.Stage = controlString(w.Stage)
	o.RetryNotBefore = controlString(w.RetryNotBefore)
	if w.Job != nil {
		o.JobID = &w.Job.ID
		o.JobHash = &w.Job.IntentHash
	}
	o.OwnerLock = controlLockObservation(filepath.Join(workerDirectory(s), "owner.lock"))
	o.SupervisorLock = controlLockObservation(filepath.Join(workerDirectory(s), "supervisor.lock"))
	o.Observation.State = w.Phase
	o.StateRevision = digest(struct {
		Revision uint64
		Semantic string
	}{s.ControlRevision, controlSemanticState(s)})
	if o.Contract == nil {
		o.Unavailable("control_unsupported")
		return o, nil
	}
	if e = ValidateControlContract(o.Contract); e != nil {
		o.Unavailable(ControlFailureOf(e).Code)
		return o, nil
	}
	if s.Result.Delivery.Phase == "completed" || w.Phase == "completed" {
		o.Observation.State = "completed"
		o.Unavailable("completed")
		return o, nil
	}
	if w.QueueReleased {
		o.Unavailable("queue_released")
		return o, nil
	}
	queue := ""
	if b.d.ControlQueue != nil {
		queue, e = b.d.ControlQueue(root, id)
	} else {
		var p interface{}
		p, e = workerPreparation(b.d, s)
		if e == nil {
			queue = digest(p)
			e = workerPreparationCurrent(b.d, s)
		}
	}
	if e != nil {
		o.Unavailable("queue_ownership_unavailable")
		return o, nil
	}
	o.QueueRevision = controlString(queue)
	if b.d.ControlQueue == nil {
		if e = WithNativeControlSnapshot(b.d, root, func(revision string) error { o.NativeRevision = &revision; return nil }); e != nil {
			o.Unavailable("queue_ownership_unavailable")
			return o, nil
		}
	}
	if verifyExecutable(s.Request.Runtime.PlyExecutable) != nil || verifyExecutable(s.Request.Runtime.Executable) != nil {
		o.Unavailable("frozen_control_unavailable")
		return o, nil
	}
	if w.Phase == "blocked" || o.OwnerLock == "unknown" || o.SupervisorLock == "unknown" {
		o.Unavailable("control_unavailable")
		return o, nil
	}
	if w.Job != nil && (w.Phase == "paused" || w.Phase == "stopped") {
		result, er := workerRunReserved(b.d, s)
		if er != nil || (workerNativeKind(w.Job.Kind) && !result.NativeCompleted) {
			o.Unavailable("attempt_outcome_unknown")
			return o, nil
		}
	}
	// Parent routing is not an exemption from validating the child's runtime,
	// ownership and positive parked-attempt evidence. The parent may project a
	// healthy child, but direct child control still has no independent authority.
	if o.Contract.ParentKind != "standalone" {
		o.Unavailable("parent_control_unsupported")
		return o, nil
	}
	switch w.Phase {
	case "running", "waiting":
		o.Offer("pause", "stop")
	case "paused":
		o.Offer("stop", "resume")
	case "stopped":
		o.Offer("resume")
	default:
		o.Unavailable("control_unavailable")
	}
	return o, nil
}

func workerLegacyBoundControl(d Dependencies, s workflowState, action string) error {
	if s.Request.ControlContract.ParentKind != "standalone" {
		return controlSafeError("parent_control_unsupported")
	}
	p, e := WorkerControls(d).Preview(s.Request.WorkspaceRoot, s.Result.RunID, action, "local-"+newControlNonce())
	if e != nil {
		return e
	}
	_, e = WorkerControls(d).Apply(*p.Request)
	return e
}

func workerSyncParentControl(d Dependencies, s workflowState) error {
	return WithWorkerParentGuard(s.Request.WorkspaceRoot, s.Request.ControlContract, func() error {
		g, e := ReadParentControl(s.Request.WorkspaceRoot, s.Request.ControlContract)
		if e != nil {
			return e
		}
		if g.ActiveRun != "" && g.ActiveRun != s.Result.RunID {
			return controlSafeError("control_owner_unknown")
		}
		action := g.DesiredAction
		_, pending, e := ControlHistory(s.Request.WorkspaceRoot, ControlTarget{"week", *s.Request.ControlContract.WeekID})
		if e != nil {
			return e
		}
		if pending != nil && (pending.Request.Input.Action == "pause" || pending.Request.Input.Action == "stop") {
			action = pending.Request.Input.Action
		}
		deadline, e := time.Parse(time.RFC3339Nano, g.Deadline)
		if e != nil {
			return e
		}
		if !controlNow(d).Before(deadline) {
			action = "stop"
		}
		if action != "pause" && action != "stop" {
			return nil
		}
		return withStore(s.Request.WorkspaceRoot, func() error {
			return d.replaceValue(filepath.Join(workerDirectory(s), "control.json"), workerControl{s.Result.RunID, action, ControlTimestamp(controlNow(d))})
		})
	})
}

// workerConsumeBoundControl executes only with the long-lived native owner
// capability. Positive attempt proof precedes parking; the atomic state write
// retains a receipt so a crash before immutable publication is recoverable.
func workerConsumeBoundControl(d Dependencies, s workflowState) (bool, error) {
	if s.Request.ControlContract == nil || s.Request.ControlContract.ParentKind != "standalone" {
		return false, nil
	}
	_, r, e := ControlHistory(s.Request.WorkspaceRoot, ControlTarget{"run", s.Result.RunID})
	if e != nil {
		return true, e
	}
	if r == nil {
		return false, nil
	}
	if s.ControlReceipt != nil && s.ControlReceipt.Sequence == r.Sequence {
		return true, FinishControl(d, *r, *s.ControlReceipt)
	}
	if verifyExecutable(s.Request.Runtime.PlyExecutable) != nil || verifyExecutable(s.Request.Runtime.Executable) != nil {
		_ = BeginControlConsumption(d, *r)
		return true, controlSafeError("control_outcome_unknown")
	}
	w := s.Result.Worker
	// Only the original generation or its recorded native wake successor can
	// consume approval. The logical worker and exact active job stay bound.
	if r.Basis.WorkerSession == nil || *r.Basis.WorkerSession != w.SessionID || r.Basis.RunRequestHash == nil || *r.Basis.RunRequestHash != s.Result.RequestSHA256 || !workerControlGenerationAllowed(s, *r) {
		return true, controlSafeError("control_outcome_unknown")
	}
	if w.Job != nil && (r.Basis.JobID == nil || *r.Basis.JobID != w.Job.ID || r.Basis.JobHash == nil || *r.Basis.JobHash != w.Job.IntentHash) {
		return true, controlSafeError("control_outcome_unknown")
	}
	a := r.Request.Input.Action
	if w.Job != nil && a != "resume" {
		result, er := workerRunReserved(d, s)
		if er != nil {
			return true, er
		}
		if workerNativeKind(w.Job.Kind) && !result.NativeCompleted {
			return true, controlSafeError("control_outcome_unknown")
		}
	}
	if e = BeginControlConsumption(d, *r); e != nil {
		return true, e
	}
	var receipt ControlReceipt
	e = workerStateUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
		w := current.Result.Worker
		effect := map[string]string{"pause": "paused", "stop": "stopped", "resume": "resumed"}[a]
		if current.Result.Delivery.Phase == "completed" {
			effect = "completed_before_boundary"
			w.Phase = "completed"
		} else {
			if w.QueueReleased {
				return controlSafeError("control_unavailable")
			}
			if a == "resume" {
				if w.Phase != "paused" && w.Phase != "stopped" {
					return controlSafeError("control_outcome_unknown")
				}
				w.Phase = "running"
				if w.RetryNotBefore != "" {
					w.Phase = "waiting"
				}
				w.Reason = ""
			} else {
				w.Phase = effect
			}
		}
		basis := r.Basis
		basis.OwnerGeneration = &w.Generation
		receipt = ReceiptForControl(d, *r, effect, basis)
		if w.Generation != *r.Basis.OwnerGeneration {
			receipt.PredecessorGeneration = r.Basis.OwnerGeneration
		}
		current.ControlReceipt = &receipt
		return nil
	})
	if e != nil {
		return true, e
	}
	return true, FinishControl(d, *r, receipt)
}

func workerControlDispatch(d Dependencies, s workflowState, fn func(*workflowState) error) error {
	return WithWorkerParentGuard(s.Request.WorkspaceRoot, s.Request.ControlContract, func() error {
		return workflowUpdate(d, s.Request.WorkspaceRoot, s.Result.RunID, func(current *workflowState) error {
			if e := workerLiveOwner(d, *current); e != nil {
				return e
			}
			if current.Request.ControlContract != nil {
				if e := CheckParentDispatch(current.Request.WorkspaceRoot, current.Request.ControlContract, controlNow(d)); e != nil {
					return e
				}
				_, pending, e := ControlHistory(current.Request.WorkspaceRoot, ControlTarget{"run", current.Result.RunID})
				if e != nil {
					return e
				}
				if pending != nil {
					return controlSafeError("control_in_progress")
				}
				if current.Result.Worker.Phase == "paused" || current.Result.Worker.Phase == "stopped" {
					return controlSafeError("control_unavailable")
				}
			}
			return fn(current)
		})
	})
}
