package weekrun

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func (m Manager) controlDependencies() taskrun.Dependencies {
	return taskrun.Dependencies{Now: m.now, Fault: m.Fault}
}
func (m Manager) Controls() taskrun.ControlService {
	return taskrun.ControlService{Dependencies: m.controlDependencies(), Backend: m}
}
func (m Manager) ControlGuard(root, id string, fn func() error) error {
	return taskrun.WithParentControlGuard(root, id, fn)
}
func (m Manager) ControlWake(root, id string) error {
	b, dir, e := readBinding(root, id)
	if e != nil {
		return e
	}
	if m.Launch == nil {
		return fmt.Errorf("native week wake unavailable")
	}
	wake := false
	e = m.ControlGuard(root, id, func() error {
		s, err := readState(b, dir)
		if err != nil {
			return err
		}
		_, pending, err := taskrun.ControlHistory(root, taskrun.ControlTarget{Kind: "week", ID: id})
		if err != nil {
			return err
		}
		if pending == nil {
			return nil
		}
		if pending.Basis.OwnerGeneration == nil || *pending.Basis.OwnerGeneration != s.Generation {
			return taskrun.ControlError("control_outcome_unknown", "Parent ownership changed before the accepted wake.")
		}
		g, err := taskrun.ReadParentControl(root, b.ControlContract)
		if err != nil {
			return err
		}
		s.ControlWake = &ControlWake{Sequence: pending.Sequence, Predecessor: s.Generation, Successor: s.Generation + 1}
		if err = m.saveBound(b, dir, &s, &g); err != nil {
			return err
		}
		wake = true
		return nil
	})
	if e != nil || !wake {
		return e
	}
	return m.Launch(b)
}
func initialGate(b Binding) taskrun.ParentControlGate {
	return taskrun.ParentControlGate{Contract: b.ControlContract, Deadline: b.Deadline, DesiredAction: "continue"}
}

// Complete registrations must agree with their gate. Only the exact initial
// state, with no activity or control history, can repair a missing publication.
func (m Manager) existingBoundStart(b Binding, dir string, check bool) (State, error) {
	var out State
	inspect := func() error {
		s, se := readState(b, dir)
		if se != nil && !os.IsNotExist(se) {
			return se
		}
		var g taskrun.ParentControlGate
		ge := readJSON(taskrun.ParentControlPath(b.Contract.Workspace, b.ID), &g)
		if ge != nil && !os.IsNotExist(ge) {
			return ge
		}
		if se == nil && ge == nil {
			active, dispatch := "", ""
			if ix := entryIndex(&s, s.ActiveTaskID); s.ActiveTaskID != "" && ix >= 0 {
				active = s.Entries[ix].RunID
			}
			if s.Dispatch != nil {
				dispatch = s.Dispatch.Sequence
			}
			resume := s.ResumeRequested && s.DesiredAction == "continue"
			stop := taskrun.ControlRevision([]string{b.ID, b.Deadline, "deadline_expired"})
			if taskrun.ControlRevision(g.Contract) != taskrun.ControlRevision(b.ControlContract) || g.Deadline != b.Deadline || g.DesiredAction != s.DesiredAction || g.Revision != s.ControlRevision || g.ActiveRun != active || g.Dispatch != dispatch || g.ResumeAllowed != resume || g.StopSequence != "" && (g.StopSequence != stop || s.DesiredAction != "stop") || s.Phase == "expired" && g.StopSequence != stop {
				return fmt.Errorf("week state and parent authority disagree; inspect locally")
			}
			out = s
			return nil
		}
		if check {
			return fmt.Errorf("initial week registration is incomplete; repeat the same start to recover its original binding")
		}
		if se == nil && !isInitialState(b, s) {
			return fmt.Errorf("missing parent authority after state changed")
		}
		if ge == nil && hash(canonical(g)) != hash(canonical(initialGate(b))) {
			return fmt.Errorf("missing week state after parent authority changed")
		}
		if _, e := os.Lstat(filepath.Join(dir, "execution-started.json")); !os.IsNotExist(e) {
			return fmt.Errorf("incomplete registration after execution or unknown activity")
		}
		history, pending, e := taskrun.ControlHistory(b.Contract.Workspace, taskrun.ControlTarget{Kind: "week", ID: b.ID})
		if e != nil {
			return e
		}
		if pending != nil || history != taskrun.ControlRevision([]string{}) {
			return fmt.Errorf("incomplete registration after a bound control")
		}
		if e = verifyControl(b); e != nil {
			return e
		}
		q, e := m.Native.Queue(b.Contract)
		if e != nil {
			return e
		}
		if e = validateQueue(b.Contract, q); e != nil {
			return e
		}
		if q.Current != nil {
			return fmt.Errorf("initial registration recovery found native activity")
		}
		for _, sel := range b.Contract.Selected {
			if _, e = selectedRow(b.Contract, sel, q); e != nil {
				return e
			}
		}
		out = initialState(b)
		if ge != nil {
			if e = publish(taskrun.ParentControlPath(b.Contract.Workspace, b.ID), initialGate(b), false); e != nil {
				return e
			}
		}
		if se != nil {
			return publish(filepath.Join(dir, "state.json"), out, false)
		}
		return nil
	}
	var err error
	if check {
		err = inspect()
	} else {
		err = m.ControlGuard(b.Contract.Workspace, b.ID, inspect)
	}
	return out, err
}

type childControlReader interface {
	BoundControlBasis(Contract, string) (taskrun.ControlBasis, error)
}
type completedChildControlReader interface {
	VerifyCompletedControlAttempt(Contract, string, taskrun.ControlBasis) error
}
type nativeControlSnapshot interface {
	WithControlSnapshot(string, func(string) error) error
}

func (n NativeAdapter) WithControlSnapshot(root string, fn func(string) error) error {
	return taskrun.WithNativeControlSnapshot(n.Dependencies, root, fn)
}
func (m Manager) ControlPublicationGuard(root, expected string, fn func() error) error {
	n, ok := m.Native.(nativeControlSnapshot)
	if !ok {
		return taskrun.ControlError("control_unavailable", "Native control ownership is unavailable.")
	}
	return n.WithControlSnapshot(root, func(current string) error {
		if current != expected {
			return taskrun.ControlError("control_stale", "Native queue or policy changed.")
		}
		return fn()
	})
}

func (n NativeAdapter) BoundControlBasis(c Contract, id string) (out taskrun.ControlBasis, err error) {
	err = at(c, func() error { var e error; out, e = taskrun.ShowControl(n.Dependencies, c.Workspace, id); return e })
	return
}
func (n NativeAdapter) VerifyCompletedControlAttempt(c Contract, id string, basis taskrun.ControlBasis) error {
	return at(c, func() error { return taskrun.VerifyCompletedControlAttempt(n.Dependencies, c.Workspace, id, basis) })
}
func nullable(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func (m Manager) ControlBasis(root, id string) (taskrun.ControlBasis, error) {
	o := taskrun.NewControlBasis(root, taskrun.ControlTarget{Kind: "week", ID: id}, m.now())
	b, dir, e := readBinding(root, id)
	if e != nil {
		return o, taskrun.ControlError("control_not_found", "The exact preserved week is unavailable.")
	}
	s, e := readState(b, dir)
	if e != nil {
		return o, taskrun.ControlError("control_outcome_unknown", "The week state is inconsistent; inspect locally.")
	}
	o.ProjectID = nullable(string(b.Contract.ProjectID))
	o.RepoID = nullable(string(b.Contract.RepoID))
	o.EpicID = nullable(string(b.Contract.EpicID))
	o.Contract = b.ControlContract
	o.FrozenControl = nullable(b.Control.SHA256)
	o.Phase = s.Phase
	o.OwnerGeneration = &s.Generation
	o.RuntimeMode = nullable("codex_headless_week_v1")
	o.Provider = nullable("codex")
	o.OwnerLock = s.OwnerState
	o.Observation.State = taskrun.PublicControlState(s.Phase)
	o.Week = &taskrun.ControlWeek{ID: b.ID, Binding: b.SHA256, Selected: canonical(b.Contract.Selected), Deadline: b.Deadline, DesiredAction: s.DesiredAction, Revision: s.ControlRevision, Generation: s.Generation, FrozenControl: b.Control.SHA256}
	if s.Dispatch != nil {
		o.Week.Dispatch = &s.Dispatch.Sequence
	}
	copy := s
	copy.UpdatedAt = ""
	copy.LastProgressAt = ""
	copy.OwnerState = ""
	copy.ControllerPID = 0
	o.StateRevision = hash(canonical(copy))
	if b.ControlContract == nil {
		o.Unavailable("parent_control_unsupported")
		return o, nil
	}
	g, e := taskrun.ReadParentControl(root, b.ControlContract)
	if e != nil {
		o.Unavailable("control_owner_unknown")
		return o, nil
	}
	o.StateRevision = taskrun.ControlRevision([]string{o.StateRevision, hash(canonical(g))})
	if verifyControl(b) != nil {
		o.Unavailable("frozen_control_unavailable")
		return o, nil
	}
	if s.ActiveTaskID != "" {
		ix := entryIndex(&s, s.ActiveTaskID)
		if ix < 0 {
			o.Unavailable("control_owner_unknown")
			return o, nil
		}
		o.TaskID = nullable(string(s.ActiveTaskID))
		o.Goal = nullable(b.Contract.Selected[ix].Goal.Spec.ManifestSHA256)
		o.RunID = nullable(s.Entries[ix].RunID)
		if o.RunID != nil {
			r, er := m.Native.Observe(b.Contract, *o.RunID)
			if er != nil {
				o.Unavailable("control_owner_unknown")
				return o, nil
			}
			o.Stage = nullable(r.Stage)
			o.RetryNotBefore = nullable(r.RetryNotBefore)
			o.OwnerLock = r.OwnerState
			o.SupervisorLock = r.SupervisorState
			if reader, ok := m.Native.(childControlReader); ok {
				child, er := reader.BoundControlBasis(b.Contract, *o.RunID)
				if er != nil || taskrun.ControlRevision(child.Contract) != taskrun.ControlRevision(b.ControlContract) {
					o.Unavailable("parent_control_unsupported")
					return o, nil
				}
				if child.UnavailableReason == nil || (*child.UnavailableReason != "parent_control_unsupported" && !(r.Completed && *child.UnavailableReason == "completed")) {
					reason := "control_owner_unknown"
					if child.UnavailableReason != nil {
						reason = *child.UnavailableReason
					}
					o.Unavailable(reason)
					return o, nil
				}
				o.RunRequestHash = child.RunRequestHash
				o.WorkerSession = child.WorkerSession
				o.JobID = child.JobID
				o.JobHash = child.JobHash
				o.StateRevision = taskrun.ControlRevision([]string{o.StateRevision, child.StateRevision})
			}
			if r.Phase == "blocked" {
				o.Unavailable("control_unavailable")
				return o, nil
			}
		}
	}
	q, e := m.Native.Queue(b.Contract)
	if e != nil || validateQueue(b.Contract, q) != nil {
		o.Unavailable("queue_ownership_unavailable")
		return o, nil
	}
	// Queue snapshots include native revisions, disposition and exact selected
	// Goal identities; their hash binds an independently mutated queue as well.
	queue := taskrun.ControlRevision(q)
	o.QueueRevision = &queue
	if n, ok := m.Native.(nativeControlSnapshot); ok {
		if e = n.WithControlSnapshot(root, func(revision string) error { o.NativeRevision = &revision; return nil }); e != nil {
			o.Unavailable("queue_ownership_unavailable")
			return o, nil
		}
	}
	deadline, e := time.Parse(time.RFC3339Nano, b.Deadline)
	if e != nil || !m.now().Before(deadline) {
		o.Unavailable("deadline_expired")
		return o, nil
	}
	if s.Phase == "completed" || s.Phase == "completed_with_blocks" || s.Phase == "expired" {
		o.Unavailable("completed")
		return o, nil
	}
	if s.Phase == "blocked" {
		o.Unavailable("control_unavailable")
		return o, nil
	}
	switch s.Phase {
	case "paused":
		o.Offer("stop", "resume")
	case "stopped":
		o.Offer("resume")
	default:
		o.Offer("pause", "stop")
	}
	for i := range o.Actions {
		o.Actions[i].Caveats = append(o.Actions[i].Caveats, "This action affects the entire preserved week and selected workset, including dependent Goals.")
	}
	return o, nil
}

func (m Manager) saveBound(b Binding, dir string, s *State, g *taskrun.ParentControlGate) error {
	previous, e := readState(b, dir)
	if e != nil {
		return e
	}
	oldGate, e := taskrun.ReadParentControl(b.Contract.Workspace, b.ControlContract)
	if e != nil {
		return e
	}
	normalize := func(v State) State {
		v.ControlRevision = 0
		v.UpdatedAt = ""
		v.LastProgressAt = ""
		v.OwnerState = ""
		return v
	}
	newGate := *g
	newGate.Revision = 0
	oldGate.Revision = 0
	s.ControlRevision = previous.ControlRevision
	if hash(canonical(normalize(previous))) != hash(canonical(normalize(*s))) || hash(canonical(oldGate)) != hash(canonical(newGate)) {
		s.ControlRevision++
	}
	g.Revision = s.ControlRevision
	g.DesiredAction = s.DesiredAction
	if e := publish(taskrun.ParentControlPath(b.Contract.Workspace, b.ID), g, true); e != nil {
		return e
	}
	return m.save(b, dir, s)
}
func (m Manager) receiptBound(b Binding, dir string, s *State, g *taskrun.ParentControlGate, r taskrun.ControlRecord, effect string) error {
	d := m.controlDependencies()
	if s.ControlReceipt != nil && s.ControlReceipt.Sequence == r.Sequence {
		return taskrun.FinishControl(d, r, *s.ControlReceipt)
	}
	return taskrun.WithControlStore(b.Contract.Workspace, func() error {
		if verifyControl(b) != nil {
			return taskrun.FailControlConsumption(d, r)
		}
		if effect != "superseded" {
			if e := m.verifyChildBoundary(b, *s, r); e != nil {
				return taskrun.FailControlConsumption(d, r)
			}
		}
		if e := taskrun.BeginControlConsumption(d, r); e != nil {
			return e
		}
		basis := r.Basis
		basis.OwnerGeneration = &s.Generation
		p := taskrun.ReceiptForControl(d, r, effect, basis)
		if r.Basis.OwnerGeneration != nil && *r.Basis.OwnerGeneration != s.Generation && parentGenerationAllowed(*s, r) {
			previous := *r.Basis.OwnerGeneration
			p.PredecessorGeneration = &previous
		}
		if effect == "superseded" {
			reason := "deadline_expired"
			p.Reason = &reason
			p.NativeStopSequence = &g.StopSequence
		}
		s.ControlReceipt = &p
		if e := m.saveBound(b, dir, s, g); e != nil {
			return e
		}
		return taskrun.FinishControl(d, r, p)
	})
}

// Called under parent and TaskRun guards immediately before publishing a
// positive boundary. It permits completed integration, not arbitrary revision
// equality or a phase string standing in for positive attempt evidence.
func (m Manager) verifyChildBoundary(b Binding, s State, r taskrun.ControlRecord) error {
	reader, ok := m.Native.(childControlReader)
	if !ok {
		return nil
	}
	if r.Basis.RunID == nil {
		if s.ActiveTaskID != "" {
			return taskrun.ControlError("control_outcome_unknown", "An unapproved child appeared.")
		}
		return nil
	}
	id := *r.Basis.RunID
	child, e := reader.BoundControlBasis(b.Contract, id)
	if e != nil {
		return e
	}
	observed, e := m.Native.Observe(b.Contract, id)
	if e != nil {
		return e
	}
	equal := func(a, b any) bool { return taskrun.ControlRevision(a) == taskrun.ControlRevision(b) }
	if !equal(child.Contract, b.ControlContract) || !equal(child.RunID, r.Basis.RunID) || !equal(child.RunRequestHash, r.Basis.RunRequestHash) || !equal(child.WorkerSession, r.Basis.WorkerSession) {
		return taskrun.ControlError("control_outcome_unknown", "The approved child identity changed.")
	}
	if child.UnavailableReason == nil || (*child.UnavailableReason != "parent_control_unsupported" && !(observed.Completed && *child.UnavailableReason == "completed")) {
		return taskrun.ControlError("control_outcome_unknown", "Positive child evidence is unavailable.")
	}
	if observed.Completed {
		if !equal(child.JobID, r.Basis.JobID) || !equal(child.JobHash, r.Basis.JobHash) {
			reader, ok := m.Native.(completedChildControlReader)
			if child.JobID != nil || child.JobHash != nil || !ok {
				return taskrun.ControlError("control_outcome_unknown", "The completed child lost its approved attempt identity.")
			}
			if e := reader.VerifyCompletedControlAttempt(b.Contract, id, r.Basis); e != nil {
				return e
			}
		}
	} else {
		ix := entryIndex(&s, s.ActiveTaskID)
		if ix < 0 || s.Entries[ix].RunID != id || !equal(child.JobID, r.Basis.JobID) || !equal(child.JobHash, r.Basis.JobHash) || observed.Phase != "paused" && observed.Phase != "stopped" {
			return taskrun.ControlError("control_outcome_unknown", "The approved child boundary changed.")
		}
	}
	return nil
}

func parentGenerationAllowed(s State, r taskrun.ControlRecord) bool {
	if r.Basis.OwnerGeneration == nil {
		return false
	}
	previous := *r.Basis.OwnerGeneration
	if s.Generation == previous {
		return true
	}
	link := s.ControlSuccessor
	return link != nil && link.Sequence == r.Sequence && link.Predecessor == previous && link.Successor == previous+1 && link.Successor == s.Generation
}

func (m Manager) claimBoundOwner(b Binding, dir string, observed State) (State, error) {
	var claimed State
	err := m.ControlGuard(b.Contract.Workspace, b.ID, func() error {
		s, e := readState(b, dir)
		if e != nil {
			return e
		}
		if _, _, e = taskrun.ControlHistory(b.Contract.Workspace, taskrun.ControlTarget{Kind: "week", ID: b.ID}); e != nil {
			return e
		}
		g, e := taskrun.ReadParentControl(b.Contract.Workspace, b.ControlContract)
		if e != nil {
			return e
		}
		if e = publish(filepath.Join(dir, "execution-started.json"), map[string]string{"id": b.ID, "contract_sha256": b.SHA256}, false); e != nil {
			return e
		}
		s.ControlSuccessor = nil
		if observed.ControlWake != nil && hash(canonical(observed.ControlWake)) == hash(canonical(s.ControlWake)) && observed.ControlWake.Predecessor == s.Generation && observed.ControlWake.Successor == s.Generation+1 {
			_, pending, e := taskrun.ControlHistory(b.Contract.Workspace, taskrun.ControlTarget{Kind: "week", ID: b.ID})
			if e != nil {
				return e
			}
			if pending != nil && pending.Sequence == observed.ControlWake.Sequence && pending.Basis.OwnerGeneration != nil && *pending.Basis.OwnerGeneration == s.Generation {
				link := *observed.ControlWake
				s.ControlSuccessor = &link
			}
		}
		s.Generation++
		s.ControllerPID = os.Getpid()
		if e = m.saveBound(b, dir, &s, &g); e != nil {
			return e
		}
		claimed = s
		return nil
	})
	return claimed, err
}

func (m Manager) saveBoundFailure(b Binding, dir string, cause error) (State, error) {
	var failed State
	err := m.ControlGuard(b.Contract.Workspace, b.ID, func() error {
		s, e := readState(b, dir)
		if e != nil {
			return e
		}
		g, e := taskrun.ReadParentControl(b.Contract.Workspace, b.ControlContract)
		if e != nil {
			return e
		}
		s.Phase = "blocked"
		s.Reason = cause.Error()
		s.NextAction = "Inspect native effects before resuming the same week"
		if e = m.saveBound(b, dir, &s, &g); e != nil {
			return e
		}
		failed = s
		return nil
	})
	return failed, err
}

// The guard protects planning and reservation only. Native launch/resume and
// launchd wait happen after release. Reconciliation reloads parent state, so a
// pause published during dispatch cannot be overwritten by a stale tick.
func (m Manager) tickBound(b Binding, dir string, out *State) error {
	var dispatch *Dispatch
	err := m.ControlGuard(b.Contract.Workspace, b.ID, func() (failure error) {
		s, e := readState(b, dir)
		if e != nil {
			return e
		}
		*out = s
		g, e := taskrun.ReadParentControl(b.Contract.Workspace, b.ControlContract)
		if e != nil {
			return e
		}
		if e = publish(filepath.Join(dir, "execution-started.json"), map[string]string{"id": b.ID, "contract_sha256": b.SHA256}, false); e != nil {
			return e
		}
		_, pending, e := taskrun.ControlHistory(b.Contract.Workspace, taskrun.ControlTarget{Kind: "week", ID: b.ID})
		if e != nil {
			return e
		}
		defer func() {
			if failure != nil && pending != nil && (s.ControlReceipt == nil || s.ControlReceipt.Sequence != pending.Sequence) {
				failure = taskrun.FailControlConsumption(m.controlDependencies(), *pending)
			}
		}()
		deadline, e := time.Parse(time.RFC3339Nano, b.Deadline)
		if e != nil {
			return e
		}
		expired := !m.now().Before(deadline)
		if expired {
			s.DesiredAction = "stop"
			g.ResumeAllowed = false
			if g.StopSequence == "" {
				g.StopSequence = taskrun.ControlRevision([]string{b.ID, b.Deadline, "deadline_expired"})
			}
			if pending != nil {
				if e = m.receiptBound(b, dir, &s, &g, *pending, "superseded"); e != nil {
					return e
				}
				pending = nil
			}
		}
		if pending != nil {
			if s.ControlReceipt != nil && s.ControlReceipt.Sequence == pending.Sequence {
				if e = taskrun.FinishControl(m.controlDependencies(), *pending, *s.ControlReceipt); e != nil {
					return e
				}
				pending = nil
			} else if !parentGenerationAllowed(s, *pending) {
				return taskrun.ControlError("control_outcome_unknown", "The current parent owner is not the approved generation or its bound wake successor.")
			} else if pending.Request.Input.Action == "resume" {
				if s.Phase != "paused" && s.Phase != "stopped" {
					return taskrun.ControlError("control_outcome_unknown", "The preserved week is no longer at its approved parking boundary.")
				}
				s.DesiredAction = "continue"
				s.ResumeRequested = true
				g.ResumeAllowed = true
				s.Phase = "running"
				e = m.receiptBound(b, dir, &s, &g, *pending, "resumed")
				*out = s
				return e
			} else {
				s.DesiredAction = pending.Request.Input.Action
				g.ResumeAllowed = false
			}
		}
		if s.Dispatch != nil {
			// Reconstruct an interrupted native start only through its read-only
			// same-request lookup. Missing evidence never authorizes another start.
			if s.Dispatch.Kind == "start" {
				ix := entryIndex(&s, s.Dispatch.TaskID)
				r, er := m.Native.Start(b.Contract, b.Contract.Selected[ix], b.Control, true)
				if er != nil || r.ID == "" {
					return taskrun.ControlError("control_outcome_unknown", "A reserved Goal start has no confirmed outcome; inspect that same request locally.")
				}
				s.Entries[ix].RunID = r.ID
				s.Entries[ix].State = "running"
				g.ActiveRun = r.ID
				s.Dispatch = nil
				g.Dispatch = ""
			} else {
				return taskrun.ControlError("control_outcome_unknown", "A native dispatch has an unknown outcome; do not start another writer.")
			}
		}
		if s.ActiveTaskID != "" {
			ix := entryIndex(&s, s.ActiveTaskID)
			entry := &s.Entries[ix]
			if entry.RunID != "" {
				r, er := m.Native.Observe(b.Contract, entry.RunID)
				if er != nil {
					return er
				}
				if r.Completed {
					projectCurrentRun(&s, entry, r, true)
					s.ActiveTaskID = ""
					g.ActiveRun = ""
					s.ResumeRequested = false
					g.ResumeAllowed = false
				} else if s.DesiredAction != "continue" {
					if r.Phase == "paused" || r.Phase == "stopped" {
						s.Phase = parkPhase(s.DesiredAction, expired)
						if pending != nil {
							e = m.receiptBound(b, dir, &s, &g, *pending, s.Phase)
							*out = s
							return e
						}
					}
					if s.Phase != "paused" && s.Phase != "stopped" && s.Phase != "expired" {
						s.Phase = "pausing"
						if s.DesiredAction == "stop" {
							s.Phase = "stopping"
						}
						dispatch = &Dispatch{Kind: "control", TaskID: entry.TaskID, RunID: entry.RunID}
					}
				} else if r.Phase == "awaiting_product_answer" {
					projectCurrentRun(&s, entry, r, true)
					s.NextCheckAt = ""
					s.ResumeRequested = false
					g.ResumeAllowed = false
					if r.OwnerState != "lock_held" && r.SupervisorState != "lock_held" {
						if r.OwnerState == "unowned" && r.SupervisorState == "unowned" {
							dispatch = &Dispatch{Kind: "resume", TaskID: entry.TaskID, RunID: entry.RunID, Explicit: false}
						} else {
							return taskrun.ControlError("control_owner_unknown", "The waiting worker's ownership is unknown; do not replace it.")
						}
					}
				} else if r.Phase == "blocked" {
					s.Reason = "The native Run requires local recovery."
					projectCurrentRun(&s, entry, r, true)
				} else if r.Phase == "paused" || r.Phase == "stopped" {
					if s.ResumeRequested && g.ResumeAllowed {
						dispatch = &Dispatch{Kind: "resume", TaskID: entry.TaskID, RunID: entry.RunID, Explicit: true}
					}
					projectCurrentRun(&s, entry, r, true)
				} else {
					projectCurrentRun(&s, entry, r, true)
					s.NextCheckAt = r.RetryNotBefore
					s.ResumeRequested = false
					g.ResumeAllowed = false
				}
			} else if s.DesiredAction == "continue" {
				dispatch = &Dispatch{Kind: "start", TaskID: entry.TaskID}
			}
		}
		if s.ActiveTaskID == "" && s.DesiredAction == "continue" {
			q, er := m.Native.Queue(b.Contract)
			if er != nil {
				return er
			}
			if er = validateQueue(b.Contract, q); er != nil {
				return er
			}
			if q.Current != nil {
				return fmt.Errorf("week queue has an unrelated owner")
			}
			for i, sel := range b.Contract.Selected {
				entry := &s.Entries[i]
				if entry.State != "pending" {
					continue
				}
				ready := true
				for _, dep := range sel.DependsOn {
					ix := entryIndex(&s, dep)
					if s.Entries[ix].State != "completed" {
						ready = false
						break
					}
					r, er := m.Native.Observe(b.Contract, s.Entries[ix].RunID)
					if er != nil || !r.Completed {
						ready = false
						break
					}
				}
				if !ready {
					continue
				}
				row, er := selectedRow(b.Contract, sel, q)
				if er != nil {
					return er
				}
				if row.State != "ready" {
					continue
				}
				if er = m.Native.Preview(b.Contract, sel, b.Control); er != nil {
					return er
				}
				s.ActiveTaskID = sel.TaskID
				entry.State = "starting"
				s.Phase = "running"
				dispatch = &Dispatch{Kind: "start", TaskID: sel.TaskID}
				break
			}
			if s.ActiveTaskID == "" {
				complete := true
				for _, entry := range s.Entries {
					complete = complete && entry.State == "completed"
				}
				if complete {
					s.Phase = "completed"
				}
			}
			projectQueueGuidance(&s)
		}
		if s.DesiredAction != "continue" && s.ActiveTaskID == "" {
			s.Phase = parkPhase(s.DesiredAction, expired)
			dispatch = nil
			if pending != nil {
				e = m.receiptBound(b, dir, &s, &g, *pending, s.Phase)
				*out = s
				return e
			}
		}
		if dispatch != nil {
			dispatch.Sequence = taskrun.ControlRevision(struct {
				ID       string
				Revision uint64
				Dispatch *Dispatch
			}{b.ID, s.ControlRevision, dispatch})
			s.Dispatch = dispatch
			g.Dispatch = dispatch.Sequence
		}
		e = m.saveBound(b, dir, &s, &g)
		*out = s
		return e
	})
	if err != nil || dispatch == nil {
		return err
	}
	if m.Fault != nil {
		if err = m.Fault("week_after_dispatch_reservation"); err != nil {
			return err
		}
	}
	var run Run
	switch dispatch.Kind {
	case "start":
		ix := entryIndex(out, dispatch.TaskID)
		run, err = m.Native.Start(b.Contract, b.Contract.Selected[ix], b.Control, false)
	case "resume":
		run, err = m.Native.Resume(b.Contract, dispatch.RunID, dispatch.Explicit)
	case "control":
		run, err = m.Native.Control(b.Contract, dispatch.RunID, out.DesiredAction)
	}
	if err != nil {
		return err
	}
	return m.ControlGuard(b.Contract.Workspace, b.ID, func() error {
		s, e := readState(b, dir)
		if e != nil {
			return e
		}
		g, e := taskrun.ReadParentControl(b.Contract.Workspace, b.ControlContract)
		if e != nil {
			return e
		}
		if s.Dispatch == nil || s.Dispatch.Sequence != dispatch.Sequence {
			return taskrun.ControlError("control_outcome_unknown", "Native dispatch ownership changed.")
		}
		if dispatch.Kind == "start" {
			if run.ID == "" {
				return taskrun.ControlError("control_outcome_unknown", "Native start returned no identity.")
			}
			ix := entryIndex(&s, dispatch.TaskID)
			s.Entries[ix].RunID = run.ID
			s.Entries[ix].BaseOID = run.BaseOID
			s.Entries[ix].BaseRevision = run.BaseRevision
			s.Entries[ix].State = "running"
			g.ActiveRun = run.ID
		}
		s.Dispatch = nil
		g.Dispatch = ""
		e = m.saveBound(b, dir, &s, &g)
		*out = s
		return e
	})
}
