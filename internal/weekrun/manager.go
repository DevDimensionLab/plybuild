package weekrun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func (m Manager) now() time.Time {
	if m.Now != nil {
		return m.Now().UTC()
	}
	return time.Now().UTC()
}
func stamp(t time.Time) string { return t.UTC().Format(time.RFC3339Nano) }

func validateQueue(c Contract, q workspace.WorkspaceTaskQueueReadback) error {
	if q.Workspace.Root != c.Workspace || q.Target.ProjectID != c.ProjectID || q.Target.RepoID != c.RepoID || q.Target.EpicID != c.EpicID || q.Target.ParentLocator != c.ReturnWorktree || q.Target.ParentRef != c.ReturnRef {
		return fmt.Errorf("native queue target differs from the frozen week return")
	}
	return nil
}
func selectedRow(c Contract, s Selection, q workspace.WorkspaceTaskQueueReadback) (workspace.QueuePending, error) {
	if err := validateQueue(c, q); err != nil {
		return workspace.QueuePending{}, err
	}
	for _, r := range q.Pending {
		if r.TaskID != s.TaskID {
			continue
		}
		if r.Goal == nil || *r.Goal != s.Goal {
			return r, fmt.Errorf("selected Goal revision/hash changed for %s", s.TaskID)
		}
		if r.Delivery == nil || r.Delivery.Mode != workspace.DeliveryLocalEpic || !r.Delivery.AutomaticAcceptance() || r.Delivery.HumanOwnedIntegration() || r.Delivery.Acceptance.RequireHumanQA {
			return r, fmt.Errorf("selected Goal %s lacks automatic local Epic delivery", s.TaskID)
		}
		if r.Executor == nil {
			return r, fmt.Errorf("selected Goal %s has no native executor", s.TaskID)
		}
		val := func(v *string, fallback string) string {
			if v == nil {
				return fallback
			}
			return *v
		}
		if val(r.Executor.Provider, "codex") != "codex" || val(r.Executor.Model, "gpt-6-astra") != c.Model || val(r.Executor.Effort, "max") != c.Effort {
			return r, fmt.Errorf("selected Goal %s executor differs from frozen model/effort", s.TaskID)
		}
		return r, nil
	}
	return workspace.QueuePending{}, fmt.Errorf("selected Goal %s is absent from its native queue", s.TaskID)
}

// Start binds authority and the real start/deadline before any product effect.
// check validates only and does not create a directory, lock or provider process.
func (m Manager) Start(c Contract, check bool) (State, error) {
	if err := Validate(c); err != nil {
		return State{}, err
	}
	id := contractID(c)
	dir, _ := runDir(c.Workspace, id)
	if _, err := os.Lstat(filepath.Join(dir, "binding.json")); err == nil {
		return m.existingStart(c, check)
	} else if !os.IsNotExist(err) {
		return State{}, err
	}
	q, err := m.Native.Queue(c)
	if err != nil {
		return State{}, err
	}
	if err = validateQueue(c, q); err != nil {
		return State{}, err
	}
	if q.Current != nil {
		return State{}, fmt.Errorf("native queue already has preparation %s; start cannot adopt unrelated work", q.Current.PreparationID)
	}
	for _, sel := range c.Selected {
		if _, err = selectedRow(c, sel, q); err != nil {
			return State{}, err
		}
	}
	now := m.now()
	b := Binding{Contract: c, SHA256: hash(canonical(c)), ID: id, StartedAt: stamp(now), Deadline: stamp(now.Add(time.Duration(c.DurationSeconds) * time.Second))}
	s := initialState(b)
	if check {
		return s, nil
	}
	if err = directory(dir); err != nil {
		return s, err
	}
	registration, err := lock(filepath.Join(c.Workspace, ".ply", "week-runs", "registration.lock"))
	if err != nil {
		return s, err
	}
	defer registration.Close()
	// Another week may not take the same Epic while this selected contract is live.
	dirs, err := os.ReadDir(filepath.Dir(dir))
	if err != nil {
		return s, err
	}
	for _, e := range dirs {
		if !e.IsDir() || e.Name() == id {
			continue
		}
		other, od, e2 := readBinding(c.Workspace, e.Name())
		if os.IsNotExist(e2) {
			continue
		}
		if e2 != nil {
			return s, e2
		}
		if other.Contract.Target() != c.Target() {
			continue
		}
		st, e2 := readState(other, od)
		if e2 != nil {
			return s, e2
		}
		if st.Phase != "completed" && st.Phase != "completed_with_blocks" && st.Phase != "expired" && st.Phase != "stopped" {
			return s, fmt.Errorf("Epic already belongs to week %s (%s)", other.ID, st.Phase)
		}
	}
	// The registration lock also serializes identical first starts.
	if _, err = os.Lstat(filepath.Join(dir, "binding.json")); err == nil {
		return m.Show(c.Workspace, id)
	} else if !os.IsNotExist(err) {
		return s, err
	}
	controlPath, err := m.Executable()
	if err != nil {
		return s, err
	}
	source, err := executable(controlPath)
	if err != nil {
		return s, err
	}
	b.Control, err = taskexecute.PreserveControl(source, dir)
	if err != nil {
		return s, err
	}
	b.Environment = map[string]string{"HOME": os.Getenv("HOME"), "PATH": os.Getenv("PATH"), "TMPDIR": os.Getenv("TMPDIR")}
	if v := os.Getenv("CODEX_HOME"); v != "" {
		b.Environment["CODEX_HOME"] = v
	}
	if err = publish(filepath.Join(dir, "binding.json"), b, false); err != nil {
		return s, err
	}
	if err = publish(filepath.Join(dir, "state.json"), s, false); err != nil {
		return s, err
	}
	return s, nil
}

func initialState(b Binding) State {
	s := State{Kind: "ply.workflow.week", SchemaVersion: 1, ID: b.ID, ContractSHA256: b.SHA256, Phase: "ready", StartedAt: b.StartedAt, Deadline: b.Deadline, LastProgressAt: b.StartedAt, UpdatedAt: b.StartedAt, NextAction: "Start only the selected native Goals", DesiredAction: "continue", Entries: []Entry{}}
	for _, sel := range b.Contract.Selected {
		s.Entries = append(s.Entries, Entry{TaskID: sel.TaskID, State: "pending"})
	}
	return s
}

// A first registration can die between its two publications. No scheduler may
// act until state exists; its first Tick publishes an activity marker before
// native effects. Only that provably unstarted gap can reconstruct initial state.
func (m Manager) existingStart(c Contract, check bool) (State, error) {
	b, dir, err := readBinding(c.Workspace, contractID(c))
	if err != nil {
		return State{}, err
	}
	if b.SHA256 != hash(canonical(c)) {
		return State{}, fmt.Errorf("week name already binds a different contract")
	}
	if _, err = readState(b, dir); err == nil {
		return m.Show(c.Workspace, b.ID)
	} else if !os.IsNotExist(err) {
		return State{}, err
	}
	if _, err = os.Lstat(filepath.Join(dir, "execution-started.json")); err == nil {
		return State{}, fmt.Errorf("week state is missing after execution began; inspect native effects before recovery")
	} else if !os.IsNotExist(err) {
		return State{}, err
	}
	if check {
		return State{}, fmt.Errorf("initial week state is incomplete; repeat the same start to recover its original binding")
	}
	guard, err := lock(filepath.Join(filepath.Dir(dir), "registration.lock"))
	if err != nil {
		return State{}, err
	}
	defer guard.Close()
	if _, err = readState(b, dir); err == nil {
		return m.Show(c.Workspace, b.ID)
	} else if !os.IsNotExist(err) {
		return State{}, err
	}
	if err = verifyControl(b); err != nil {
		return State{}, err
	}
	q, err := m.Native.Queue(c)
	if err != nil {
		return State{}, err
	}
	if err = validateQueue(c, q); err != nil {
		return State{}, err
	}
	if q.Current != nil {
		return State{}, fmt.Errorf("initial state recovery found native activity")
	}
	for _, sel := range c.Selected {
		if _, err = selectedRow(c, sel, q); err != nil {
			return State{}, err
		}
	}
	s := initialState(b)
	if err = publish(filepath.Join(dir, "state.json"), s, false); err != nil {
		return s, err
	}
	return s, nil
}

func (m Manager) Show(root, id string) (State, error) {
	b, dir, err := readBinding(root, id)
	if err != nil {
		return State{}, err
	}
	s, err := readState(b, dir)
	if err != nil {
		return s, err
	}
	// Observation is read-only: do not create an ownership lock merely to show.
	if _, e := os.Lstat(filepath.Join(dir, "owner.lock")); os.IsNotExist(e) {
		s.OwnerState = "not_started"
	} else {
		s.OwnerState = ownerState(dir)
	}
	if (s.OwnerState == "unowned" || s.OwnerState == "not_started") && s.ActiveTaskID != "" && s.DesiredAction != "continue" {
		ix := entryIndex(&s, s.ActiveTaskID)
		if ix >= 0 && s.Entries[ix].RunID != "" {
			run, e := m.Native.Observe(b.Contract, s.Entries[ix].RunID)
			if e != nil {
				return s, e
			}
			if run.Phase == "paused" || run.Phase == "stopped" || run.Completed {
				deadline, _ := time.Parse(time.RFC3339Nano, b.Deadline)
				s.Phase = parkPhase(s.DesiredAction, !m.now().Before(deadline))
				s.NextAction = "Native worker reached its confirmed boundary; the controller will reconcile preserved work on resume"
			}
		}
	}
	var c Control
	if e := readJSON(filepath.Join(dir, "control.json"), &c); e == nil && c.ID != s.SeenControl {
		s.NextAction = c.Action + " requested at " + c.RequestedAt
	} else if e != nil && !os.IsNotExist(e) {
		return s, e
	}
	return s, nil
}

func (m Manager) Control(root, id, action string) (State, error) {
	if action != "pause" && action != "stop" && action != "resume" {
		return State{}, fmt.Errorf("week action must be pause, stop or resume")
	}
	b, dir, err := readBinding(root, id)
	if err != nil {
		return State{}, err
	}
	deadline, _ := time.Parse(time.RFC3339Nano, b.Deadline)
	if action == "resume" && !m.now().Before(deadline) {
		return m.Show(root, id)
	}
	guard, err := lock(filepath.Join(dir, "control.lock"))
	if err != nil {
		return State{}, err
	}
	defer guard.Close()
	now := stamp(m.now())
	ctrl := Control{ID: hash(canonical([]string{action, now})), Action: action, RequestedAt: now}
	if err = publish(filepath.Join(dir, "control.json"), ctrl, true); err != nil {
		return State{}, err
	}
	if m.Launch != nil {
		if err = m.Launch(b); err != nil {
			return State{}, err
		}
	}
	return m.Show(root, id)
}

func (m Manager) LaunchRun(root, id string) error {
	b, _, err := readBinding(root, id)
	if err != nil {
		return err
	}
	if err = verifyControl(b); err != nil {
		return err
	}
	if m.Launch == nil {
		return fmt.Errorf("week launch adapter unavailable")
	}
	return m.Launch(b)
}

func (m Manager) save(b Binding, dir string, s *State) error {
	s.UpdatedAt = stamp(m.now())
	return publish(filepath.Join(dir, "state.json"), s, true)
}

func (m Manager) recordControl(dir string, s *State) error {
	var c Control
	if err := readJSON(filepath.Join(dir, "control.json"), &c); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if c.ID == s.SeenControl {
		return nil
	}
	if c.Action != "resume" && c.Action != "pause" && c.Action != "stop" {
		return fmt.Errorf("unknown week control")
	}
	s.SeenControl = c.ID
	s.DesiredAction = c.Action
	s.ResumeRequested = c.Action == "resume"
	if s.ResumeRequested {
		s.DesiredAction = "continue"
	}
	return nil
}

// Tick is called under the week and Epic ownership locks by Serve. It writes a
// launch intent before invoking native Execute; a restart reuses that same Goal.
func (m Manager) Tick(b Binding, dir string, s *State) error {
	if err := publish(filepath.Join(dir, "execution-started.json"), map[string]string{"id": b.ID, "contract_sha256": b.SHA256}, false); err != nil {
		return err
	}
	if err := m.recordControl(dir, s); err != nil {
		return err
	}
	deadline, err := time.Parse(time.RFC3339Nano, b.Deadline)
	if err != nil {
		return err
	}
	expired := !m.now().Before(deadline)
	if expired {
		s.DesiredAction = "stop"
		s.ResumeRequested = false
	}
	if s.Phase == "blocked" && !s.ResumeRequested && s.DesiredAction == "continue" {
		return nil
	}
	if s.ActiveTaskID != "" {
		ix := entryIndex(s, s.ActiveTaskID)
		if ix < 0 {
			return fmt.Errorf("active Task is outside week selection")
		}
		e := &s.Entries[ix]
		if e.RunID == "" {
			// Read native recovery first even for pause: an unacknowledged start may
			// already have created a Run. Check never starts a new provider.
			run, er := m.Native.Start(b.Contract, b.Contract.Selected[ix], b.Control, true)
			if run.ID != "" {
				e.RunID = run.ID
				e.State = "running"
				e.BaseOID, e.BaseRevision = run.BaseOID, run.BaseRevision
				if er = m.save(b, dir, s); er != nil {
					return er
				}
			}
			if er != nil {
				s.Phase = "blocked"
				s.Reason = "Inspect preserved native start: " + er.Error()
				s.NextAction = "Resolve this same selected Goal; no other writer is started"
				return m.save(b, dir, s)
			}
			if e.RunID == "" && s.DesiredAction != "continue" {
				s.Phase = parkPhase(s.DesiredAction, expired)
				s.NextAction = "Start intent preserved; resume retains the same selected Goal"
				return m.save(b, dir, s)
			}
			if e.RunID == "" {
				run, er = m.Native.Start(b.Contract, b.Contract.Selected[ix], b.Control, false)
				if run.ID != "" {
					e.RunID = run.ID
					e.State = "running"
					e.BaseOID, e.BaseRevision = run.BaseOID, run.BaseRevision
					s.LastProgressAt = stamp(m.now())
				}
				if er != nil {
					s.Phase = "blocked"
					s.Reason = "Inspect preserved native start: " + er.Error()
					s.NextAction = "Recover the same Goal and native request before continuing"
					return m.save(b, dir, s)
				}
				if e.RunID == "" {
					return fmt.Errorf("native start returned no Run identity")
				}
				if err = m.save(b, dir, s); err != nil {
					return err
				}
			}
		}
		run, er := m.Native.Observe(b.Contract, e.RunID)
		if er != nil {
			return er
		}
		if run.Completed {
			e.State = "completed"
			e.Reason = ""
			s.ActiveTaskID = ""
			s.Phase = "running"
			s.LastProgressAt = stamp(m.now())
			s.Reason = ""
			s.ResumeRequested = false
			// Save the native completion observation before selecting another Goal.
			return m.save(b, dir, s)
		}
		if s.DesiredAction != "continue" {
			action := s.DesiredAction
			if run.Phase == "paused" || run.Phase == "stopped" {
				s.Phase = parkPhase(action, expired)
				s.NextAction = "Work preserved; resume retains native Run and original deadline"
				s.ResumeRequested = false
				return m.save(b, dir, s)
			}
			if run.ControlAction != action {
				if _, er = m.Native.Control(b.Contract, e.RunID, action); er != nil {
					return er
				}
			} else if (run.OwnerState == "unowned" || run.OwnerState == "not_started") && (run.SupervisorState == "unowned" || run.SupervisorState == "not_started") {
				if _, er = m.Native.Resume(b.Contract, e.RunID, false); er != nil {
					return er
				}
			} else if run.OwnerState != "lock_held" && run.SupervisorState != "lock_held" {
				return fmt.Errorf("native worker control ownership is unknown")
			}

			s.Phase = "stopping"
			if action == "pause" {
				s.Phase = "pausing"
			}
			s.NextAction = "Waiting for the active native attempt's confirmed boundary"
			return m.save(b, dir, s)
		}
		if run.Phase == "blocked" {
			if run.Releasable {
				if run.OwnerState == "lock_held" {
					s.Phase = "running"
					s.Reason = run.Reason
					s.NextAction = "Wait for the concretely blocked native controller to release its process lock"
					return m.save(b, dir, s)
				}
				// The native adapter reacquires the worker lock and proves result/process
				// quiescence before this irreversible native queue release.
				e.State = "blocked"
				e.Reason = run.Reason
				if er = m.save(b, dir, s); er != nil {
					return er
				}
				if er = m.Native.Release(b.Contract, e.RunID, e.TaskID); er != nil {
					s.Phase = "blocked"
					s.Reason = "Native blocked work cannot safely be released: " + er.Error()
					return m.save(b, dir, s)
				}
				s.ActiveTaskID = ""
				s.Phase = "running"
				s.LastProgressAt = stamp(m.now())
				s.ResumeRequested = false
				return m.save(b, dir, s)
			}
			if s.ResumeRequested {
				if _, er = m.Native.Resume(b.Contract, e.RunID, true); er != nil {
					return er
				}
				s.ResumeRequested = false
				s.Phase = "running"
				s.Reason = ""
				return m.save(b, dir, s)
			}
			s.Phase = "blocked"
			s.Reason = run.Reason
			s.NextAction = "Resolve native authentication or unknown effects; the same Run remains active"
			return m.save(b, dir, s)
		}
		if run.Phase == "paused" || run.Phase == "stopped" {
			if !s.ResumeRequested {
				s.Phase = run.Phase
				s.NextAction = "Use week resume for the preserved native Run"
				return m.save(b, dir, s)
			}
			if _, er = m.Native.Resume(b.Contract, e.RunID, true); er != nil {
				return er
			}
			s.ResumeRequested = false
		} else if (run.OwnerState == "unowned" || run.OwnerState == "not_started") && (run.SupervisorState == "unowned" || run.SupervisorState == "not_started") {
			// false recovers supervision without clearing a native blocker or control.
			if _, er = m.Native.Resume(b.Contract, e.RunID, false); er != nil {
				return er
			}
		} else if run.OwnerState != "lock_held" && run.SupervisorState != "lock_held" {
			return fmt.Errorf("native worker ownership is unknown")
		}
		s.Phase = "running"
		s.ResumeRequested = false
		e.State = "running"
		s.Reason = ""
		s.NextAction = "Observe the same native Run"
		s.NextCheckAt = ""
		if run.Phase == "waiting" {
			s.Phase = "waiting"
			e.State = "waiting"
			s.NextCheckAt = run.RetryNotBefore
			s.Reason = run.Reason
			s.NextAction = "Wait for the persisted provider reset; keep the selected model and Run"
		}
		return m.save(b, dir, s)
	}
	if s.DesiredAction != "continue" {
		s.Phase = parkPhase(s.DesiredAction, expired)
		s.NextAction = "No active Run; selected work and original deadline preserved"
		return m.save(b, dir, s)
	}
	if s.Phase == "blocked" && !s.ResumeRequested {
		return nil
	}
	s.ResumeRequested = false
	// Native delivery receipts, rather than cached entry labels, authorize
	// dependent work. The operational record only remembers native Run IDs.
	for i := range s.Entries {
		if s.Entries[i].State == "completed" {
			r, e := m.Native.Observe(b.Contract, s.Entries[i].RunID)
			if e != nil {
				return e
			}
			if !r.Completed {
				return fmt.Errorf("native completion no longer confirms %s", s.Entries[i].TaskID)
			}
		}
	}
	q, err := m.Native.Queue(b.Contract)
	if err != nil {
		return err
	}
	if err = validateQueue(b.Contract, q); err != nil {
		return err
	}
	if q.Current != nil {
		s.Phase = "blocked"
		s.Reason = "Native queue has an unowned current preparation: " + q.Current.PreparationID
		s.NextAction = "Inspect the existing native preparation"
		return m.save(b, dir, s)
	}
	changed := true
	for changed {
		changed = false
		for i, sel := range b.Contract.Selected {
			e := &s.Entries[i]
			if e.State != "pending" {
				continue
			}
			for _, dep := range sel.DependsOn {
				d := s.Entries[entryIndex(s, dep)]
				if d.State == "blocked" || d.State == "dependency_blocked" {
					e.State = "dependency_blocked"
					e.Reason = "Depends on blocked selected Task " + string(dep)
					changed = true
					break
				}
			}
		}
	}
	for i, sel := range b.Contract.Selected {
		e := &s.Entries[i]
		if e.State != "pending" {
			continue
		}
		ready := true
		for _, dep := range sel.DependsOn {
			if s.Entries[entryIndex(s, dep)].State != "completed" {
				ready = false
			}
		}
		if !ready {
			continue
		}
		row, er := selectedRow(b.Contract, sel, q)
		if er != nil {
			e.State = "blocked"
			e.Reason = er.Error()
			continue
		}
		if row.State != "ready" {
			e.State = "blocked"
			e.Reason = "Native Goal is not buildable: " + queueReason(row.Reasons)
			continue
		}
		if er = m.Native.Preview(b.Contract, sel, b.Control); er != nil {
			e.State = "blocked"
			e.Reason = "Native preflight: " + er.Error()
			continue
		}
		e.State = "starting"
		s.ActiveTaskID = e.TaskID
		s.Phase = "running"
		s.Reason = ""
		s.NextAction = "Start this exact selected Goal from the freshly observed native base"
		s.LastProgressAt = stamp(m.now())
		return m.save(b, dir, s)
	}
	pending, blocked := false, false
	for _, e := range s.Entries {
		pending = pending || e.State == "pending"
		blocked = blocked || e.State == "blocked" || e.State == "dependency_blocked"
	}
	if pending {
		s.Phase = "running"
		s.NextAction = "Reevaluate dependencies after concrete native blockers"
	} else if blocked {
		s.Phase = "completed_with_blocks"
		s.NextAction = "Independent selected work is finished; blocked Tasks and dependent Goals are preserved"
	} else {
		s.Phase = "completed"
		s.NextAction = "Selected queue complete; no additional work is authorized"
	}
	s.LastProgressAt = stamp(m.now())
	return m.save(b, dir, s)
}
func entryIndex(s *State, id workspace.TaskID) int {
	for i, e := range s.Entries {
		if e.TaskID == id {
			return i
		}
	}
	return -1
}
func queueReason(rr []workspace.QueueReason) string {
	v := []string{}
	for _, r := range rr {
		v = append(v, r.Code+": "+r.Message)
	}
	return strings.Join(v, "; ")
}
func parkPhase(action string, expired bool) string {
	if expired {
		return "expired"
	}
	if action == "pause" {
		return "paused"
	}
	return "stopped"
}
func terminal(s State) bool {
	switch s.Phase {
	case "completed", "completed_with_blocks", "expired", "paused", "stopped", "blocked":
		return true
	}
	return false
}

// Serve owns one controller and one Epic lock. A signal/process death releases
// only the scheduler; native worker ownership and receipts survive separately.
func (m Manager) Serve(ctx context.Context, root, id string) error {
	err := m.serve(ctx, root, id, false)
	if errors.Is(err, errOwned) {
		return nil
	}
	return err
}
func (m Manager) serve(ctx context.Context, root, id string, service bool) error {
	b, dir, err := readBinding(root, id)
	if err != nil {
		return err
	}
	if err = verifyControl(b); err != nil {
		return err
	}
	owner, err := lock(filepath.Join(dir, "owner.lock"))
	if err != nil {
		return err
	}
	defer owner.Close()
	if service {
		if e := os.Remove(filepath.Join(dir, "handover.json")); e != nil && !os.IsNotExist(e) {
			return e
		}
	}
	targetKey := strings.TrimPrefix(hash(canonical(b.Contract.Target())), "sha256:")
	target, err := lock(filepath.Join(filepath.Dir(dir), "epic-"+targetKey+".lock"))
	if err != nil {
		return fmt.Errorf("another week controls this Epic: %w", err)
	}
	defer target.Close()
	s, err := readState(b, dir)
	if err != nil {
		return err
	}
	s.Generation++
	s.ControllerPID = os.Getpid()
	if err = m.save(b, dir, &s); err != nil {
		return err
	}
	for {
		if err = ctx.Err(); err != nil {
			return err
		}
		if !service {
			if _, e := os.Lstat(filepath.Join(dir, "handover.json")); e == nil {
				return nil
			} else if !os.IsNotExist(e) {
				return e
			}
		}
		if err = m.Tick(b, dir, &s); err != nil {
			s.Phase = "blocked"
			s.Reason = err.Error()
			s.NextAction = "Inspect native effects before resuming the same week"
			_ = m.save(b, dir, &s)
			exited, exitErr := exitAtBoundary(dir, s, owner, target)
			if exitErr != nil {
				return exitErr
			}
			if exited {
				return err
			}
			continue
		}
		if terminal(s) {
			exited, e := exitAtBoundary(dir, s, owner, target)
			if e != nil {
				return e
			}
			if exited {
				return nil
			}
			continue
		}
		delay := time.Second * 3
		if s.ActiveTaskID != "" {
			delay = 5 * time.Second
		}
		// Backoff is bounded by control polling, so pause/stop remain available.
		if reset, e := time.Parse(time.RFC3339Nano, s.NextCheckAt); e == nil && reset.After(m.now()) {
			delay = 15 * time.Second
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}

// Serialize release with control publication: a new request either remains for
// this owner to consume or observes both locks released and launches its successor.
func exitAtBoundary(dir string, s State, owner, target *os.File) (bool, error) {
	guard, err := lock(filepath.Join(dir, "control.lock"))
	if errors.Is(err, errOwned) {
		time.Sleep(10 * time.Millisecond)
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer guard.Close()
	var c Control
	if err = readJSON(filepath.Join(dir, "control.json"), &c); err != nil && !os.IsNotExist(err) {
		return false, err
	}
	if err == nil && c.ID != s.SeenControl {
		return false, nil
	}
	if err = target.Close(); err != nil {
		return false, err
	}
	if err = owner.Close(); err != nil {
		return false, err
	}
	return true, nil
}
