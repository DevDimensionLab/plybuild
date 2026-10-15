package weekrun

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func newMR05Week(t *testing.T) *fixture {
	t.Helper()
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	if e = os.Chmod(root, 0700); e != nil {
		t.Fatal(e)
	}
	c := Contract{Kind: "ply.workflow.week-contract", SchemaVersion: 2, Name: "mobile-week", Workspace: root, ReturnWorktree: root, ReturnRef: "refs/heads/epic", ProjectID: "test", RepoID: "test", EpicID: "epic", Selected: []Selection{selection("first"), selection("second", "first")}, Model: "gpt-6-astra", Effort: "max", DurationSeconds: 3600, MaxAttempts: 3, AgentTimeoutSeconds: 300, MaxProviderRetries: 20, RuntimeControl: &RuntimeControl{SchemaVersion: 1, Mode: taskrun.BoundControlMode}}
	n := &fakeNative{c: c, runs: map[workspace.TaskID]Run{}, released: map[workspace.TaskID]bool{}, block: map[workspace.TaskID]string{}, base: 1}
	exe := filepath.Join(root, "synthetic-control")
	if e = os.WriteFile(exe, []byte("#!/bin/sh\nexit 99\n"), 0700); e != nil {
		t.Fatal(e)
	}
	now := time.Date(2026, 10, 10, 8, 0, 0, 0, time.UTC)
	m := Manager{Native: n, Now: func() time.Time { return now }, Executable: func() (string, error) { return exe, nil }, Launch: func(Binding) error { return nil }}
	s, e := m.Start(c, false)
	if e != nil {
		t.Fatal(e)
	}
	b, dir, e := readBinding(root, s.ID)
	if e != nil {
		t.Fatal(e)
	}
	return &fixture{m: m, c: c, b: b, s: s, dir: dir, n: n, now: &now}
}

func TestMR05WeekControlScopeAndOptIn(t *testing.T) {
	f := newMR05Week(t)
	p, e := f.m.Controls().Preview(f.c.Workspace, f.b.ID, "pause", "week-pause-0001")
	if e != nil {
		t.Fatal(e)
	}
	if p.Request.Input.Target.Kind != "week" || p.Basis.RunID != nil || p.Basis.TaskID != nil || p.Basis.Week.Deadline != f.b.Deadline || len(p.Basis.Actions) != 2 {
		t.Fatalf("week preview lost scope: %+v", p)
	}
	if f.b.ControlContract == nil || f.b.ControlContract.WeekID == nil || *f.b.ControlContract.WeekID != f.b.ID {
		t.Fatal("parent contract not preserved")
	}
	var selected []Selection
	if e = json.Unmarshal(p.Basis.Week.Selected, &selected); e != nil || !reflect.DeepEqual(selected, f.c.Selected) {
		t.Fatal("preview changed workset/dependencies", e)
	}
	bad := f.c
	bad.SchemaVersion = 1
	if Validate(bad) == nil {
		t.Fatal("v1 silently upgraded")
	}
	bad = f.c
	bad.RuntimeControl = nil
	if Validate(bad) == nil {
		t.Fatal("v2 missing opt-in accepted")
	}
	old := newFixture(t, selection("old"))
	b, e := old.m.Controls().Show(old.c.Workspace, old.b.ID)
	if e != nil || len(b.Actions) != 0 || b.UnavailableReason == nil || *b.UnavailableReason != "parent_control_unsupported" {
		t.Fatal("old week supported", b, e)
	}
	f.tick(t)
	p, e = f.m.Controls().Preview(f.c.Workspace, f.b.ID, "pause", "week-active-0001")
	if e != nil || p.Basis.RunID == nil || *p.Basis.RunID != "run-first" || p.Request.Input.Target.ID != f.b.ID {
		t.Fatal("active run obscured parent scope", p, e)
	}
}

func weekPreview(t *testing.T, f *fixture, action, id string) taskrun.ControlRequest {
	t.Helper()
	p, e := f.m.Controls().Preview(f.c.Workspace, f.b.ID, action, id)
	if e != nil {
		t.Fatal(e)
	}
	return *p.Request
}
func weekApply(t *testing.T, f *fixture, r taskrun.ControlRequest) {
	t.Helper()
	s, e := f.m.Controls().Apply(r)
	if e != nil || s.RequestState != "accepted" || s.ControlEffect != "pending_boundary" {
		t.Fatalf("apply %+v %v", s, e)
	}
}
func weekEffect(t *testing.T, f *fixture, id, effect string) taskrun.ControlStatus {
	t.Helper()
	s, e := f.m.Controls().Status(f.c.Workspace, id)
	if e != nil || s.ControlEffect != effect {
		t.Fatalf("effect=%s want=%s: %+v %v", s.ControlEffect, effect, s, e)
	}
	return s
}
func nativePark(f *fixture, phase string) {
	r := f.n.runs["first"]
	r.Phase = phase
	r.OwnerState = "unowned"
	r.SupervisorState = "unowned"
	f.n.runs["first"] = r
}

func TestMR05WeekPauseDispatchReservationRace(t *testing.T) {
	t.Run("owner_generation_publication_serializes_with_apply", mr05OwnerPublicationRace)
	t.Run("failure_publication_serializes_with_apply", mr05FailurePublicationRace)
	t.Run("accepted_wake_links_only_its_successor", mr05BoundWakeSuccessor)
	t.Run("pause_at_completed_native_child_parks_between_goals", mr05CompletedChildPause)
	t.Run("pause_wins", func(t *testing.T) {
		f := newMR05Week(t)
		r := weekPreview(t, f, "pause", "pause-first-0001")
		weekApply(t, f, r)
		f.tick(t)
		weekEffect(t, f, r.RequestID, "paused")
		if len(f.n.starts) != 0 || f.s.ActiveTaskID != "" {
			t.Fatal("launched past pause")
		}
	})
	t.Run("dispatch_wins", func(t *testing.T) {
		f := newMR05Week(t)
		reserved, release := make(chan struct{}), make(chan struct{})
		f.m.Fault = func(p string) error {
			if p == "week_after_dispatch_reservation" {
				close(reserved)
				<-release
			}
			return nil
		}
		done := make(chan error, 1)
		go func() { done <- f.m.Tick(f.b, f.dir, &f.s) }()
		<-reserved
		r := weekPreview(t, f, "pause", "dispatch-first-0001")
		p, e := f.m.Controls().Show(f.c.Workspace, f.b.ID)
		if e != nil || p.Week.Dispatch == nil || len(f.n.starts) != 0 {
			t.Fatal("dispatch gap unbound", p, e)
		}
		weekApply(t, f, r)
		weekEffect(t, f, r.RequestID, "pending_boundary")
		close(release)
		if e := <-done; e != nil {
			t.Fatal(e)
		}
		f.m.Fault = nil
		f.tick(t)
		weekEffect(t, f, r.RequestID, "pending_boundary")
		if len(f.n.starts) != 1 {
			t.Fatal("dispatch count", f.n.starts)
		}
		nativePark(f, "paused")
		f.tick(t)
		weekEffect(t, f, r.RequestID, "paused")
		f.tick(t)
		if len(f.n.starts) != 1 {
			t.Fatal("started next Goal while paused")
		}
	})
	t.Run("final_return_race", func(t *testing.T) {
		f := newMR05Week(t)
		f.tick(t)
		r := weekPreview(t, f, "stop", "final-return-0001")
		weekApply(t, f, r)
		f.n.complete("first")
		f.tick(t)
		weekEffect(t, f, r.RequestID, "stopped")
		if len(f.n.starts) != 1 || f.s.ActiveTaskID != "" {
			t.Fatal("next Goal raced past parent stop")
		}
	})
}

func TestMR05WeekControlResumePreservesSelectionAndDeadline(t *testing.T) {
	t.Run("missing_state_cannot_reset_confirmed_pause", func(t *testing.T) {
		f := newMR05Week(t)
		wakes := 0
		f.m.Launch = func(Binding) error { wakes++; return nil }
		r := weekPreview(t, f, "pause", "missing-state-pause-0001")
		weekApply(t, f, r)
		f.tick(t)
		weekEffect(t, f, r.RequestID, "paused")
		paths := []string{filepath.Join(f.dir, "binding.json"), taskrun.ParentControlPath(f.c.Workspace, f.b.ID), filepath.Join(f.c.Workspace, ".ply", "runtime-controls", "v1", r.RequestID, "request.json"), filepath.Join(f.c.Workspace, ".ply", "runtime-controls", "v1", r.RequestID, "receipt.json")}
		before := map[string]string{}
		for _, p := range paths {
			b, e := os.ReadFile(p)
			if e != nil {
				t.Fatal(e)
			}
			before[p] = string(b)
		}
		if e := os.Remove(filepath.Join(f.dir, "state.json")); e != nil {
			t.Fatal(e)
		}
		*f.now = f.now.Add(time.Minute)
		s, err := f.m.Start(f.c, false)
		if err == nil && (s.Phase != "paused" || s.DesiredAction != "pause" || s.Deadline != f.b.Deadline) {
			t.Errorf("missing state reset confirmed pause: phase=%s desired=%s deadline=%s", s.Phase, s.DesiredAction, s.Deadline)
		}
		if err == nil {
			if e := f.m.Tick(f.b, f.dir, &s); e != nil {
				t.Fatal(e)
			}
		}
		if len(f.n.starts) != 0 || wakes != 1 {
			t.Errorf("missing state dispatched or woke new work: starts=%v wakes=%d", f.n.starts, wakes)
		}
		for _, p := range paths {
			b, e := os.ReadFile(p)
			if e != nil || string(b) != before[p] {
				t.Errorf("retained authority/history changed: %s (%v)", filepath.Base(p), e)
			}
		}
	})
	t.Run("unstarted_registration_gap_keeps_original_deadline", func(t *testing.T) {
		f := newMR05Week(t)
		if e := os.Remove(filepath.Join(f.dir, "state.json")); e != nil {
			t.Fatal(e)
		}
		if e := os.Remove(taskrun.ParentControlPath(f.c.Workspace, f.b.ID)); e != nil {
			t.Fatal(e)
		}
		*f.now = f.now.Add(time.Minute)
		s, err := f.m.Start(f.c, false)
		if err != nil || s.Deadline != f.b.Deadline || s.StartedAt != f.b.StartedAt || len(f.n.starts) != 0 {
			t.Fatal("legitimate unstarted registration gap was not recovered", s, err)
		}
		if _, err := taskrun.ReadParentControl(f.c.Workspace, f.b.ControlContract); err != nil {
			t.Fatal("recovered registration is missing its parent authority", err)
		}
	})
	for _, action := range []string{"pause", "stop"} {
		t.Run(action, func(t *testing.T) {
			f := newMR05Week(t)
			original := canonical(f.b)
			r := weekPreview(t, f, action, "park-week-0001")
			weekApply(t, f, r)
			f.tick(t)
			effect := map[string]string{"pause": "paused", "stop": "stopped"}[action]
			weekEffect(t, f, r.RequestID, effect)
			// Reconstruct the scheduler using the retained binding and operational state.
			f.m = Manager{Native: f.n, Now: func() time.Time { return *f.now }, Executable: f.m.Executable, Launch: func(Binding) error { return nil }}
			f.reload(t)
			f.tick(t)
			if len(f.n.starts) != 0 {
				t.Fatal("restart cleared park")
			}
			r = weekPreview(t, f, "resume", "resume-week-0001")
			weekApply(t, f, r)
			f.tick(t)
			weekEffect(t, f, r.RequestID, "resumed")
			if len(f.n.starts) != 0 {
				t.Fatal("receipt fabricated provider start")
			}
			f.tick(t)
			if len(f.n.starts) != 1 || f.s.Entries[0].BaseRevision != 1 {
				t.Fatal("first bound Goal not selected")
			}
			f.n.complete("first")
			f.tick(t)
			if len(f.n.starts) != 2 || f.s.Entries[1].BaseRevision != 2 {
				t.Fatal("dependent Goal did not use fresh native base", f.s)
			}
			if string(canonical(f.b)) != string(original) {
				t.Fatal("resume rewrote binding/deadline")
			}
		})
	}
	t.Run("same_active_run", func(t *testing.T) {
		f := newMR05Week(t)
		f.tick(t)
		r := weekPreview(t, f, "pause", "active-pause-0001")
		weekApply(t, f, r)
		f.tick(t)
		nativePark(f, "paused")
		f.tick(t)
		r = weekPreview(t, f, "resume", "active-resume-0001")
		weekApply(t, f, r)
		f.tick(t)
		f.tick(t)
		if len(f.n.starts) != 1 || f.s.Entries[0].RunID != "run-first" || len(f.n.resumes) != 1 {
			t.Fatal("resume replaced active Run")
		}
	})
}

func TestMR05WeekControlDeadlineSupersedesPendingResume(t *testing.T) {
	f := newMR05Week(t)
	r := weekPreview(t, f, "pause", "deadline-pause-0001")
	weekApply(t, f, r)
	f.tick(t)
	r = weekPreview(t, f, "resume", "deadline-resume-0001")
	weekApply(t, f, r)
	deadline, _ := time.Parse(time.RFC3339Nano, f.b.Deadline)
	*f.now = deadline.Add(time.Second)
	f.tick(t)
	s := weekEffect(t, f, r.RequestID, "superseded")
	if s.Receipt.Reason == nil || *s.Receipt.Reason != "deadline_expired" || s.Receipt.NativeStopSequence == nil || *s.Receipt.NativeStopSequence == "" {
		t.Fatal("deadline lost stop reference", s)
	}
	if len(f.n.starts) != 0 || f.s.Phase != "expired" {
		t.Fatal("deadline dispatched")
	}
	if _, e := f.m.Controls().Apply(r); e != nil {
		t.Fatal("accepted identity lost after deadline", e)
	}
	active := newMR05Week(t)
	active.tick(t)
	r = weekPreview(t, active, "pause", "active-deadline-0001")
	weekApply(t, active, r)
	end, _ := time.Parse(time.RFC3339Nano, active.b.Deadline)
	*active.now = end
	active.tick(t)
	weekEffect(t, active, r.RequestID, "superseded")
	if active.s.Phase == "expired" {
		t.Fatal("active attempt declared quiescent on deadline")
	}
	nativePark(active, "stopped")
	active.tick(t)
	if active.s.Phase != "expired" || len(active.n.starts) != 1 {
		t.Fatal("deadline lost positive boundary")
	}
}

func TestMR05WeekControlChildAndLegacyAuthority(t *testing.T) {
	t.Run("native_registry_publication_contention", mr05NativeRegistryContention)
	t.Run("native_child_evidence_revalidated", mr05NativeChildEvidence)
	t.Run("shutdown_sees_reserved_resume", func(t *testing.T) {
		f := newMR05Week(t)
		r := weekPreview(t, f, "pause", "exit-pause-0001")
		weekApply(t, f, r)
		f.tick(t)
		r = weekPreview(t, f, "resume", "exit-resume-0001")
		weekApply(t, f, r)
		owner, e := lock(filepath.Join(f.dir, "owner.lock"))
		if e != nil {
			t.Fatal(e)
		}
		defer owner.Close()
		target, e := lock(filepath.Join(f.dir, "test-target.lock"))
		if e != nil {
			t.Fatal(e)
		}
		defer target.Close()
		exited, e := exitAtBoundary(f.dir, f.s, owner, target)
		if e != nil || exited {
			t.Fatal("shutdown lost durable resume wake", exited, e)
		}
	})
	f := newMR05Week(t)
	r := weekPreview(t, f, "pause", "authority-0001")
	weekApply(t, f, r)
	if e := taskrun.WithWorkerParentGuard(f.c.Workspace, f.b.ControlContract, func() error {
		return taskrun.WithControlStore(f.c.Workspace, func() error { return taskrun.CheckParentDispatch(f.c.Workspace, f.b.ControlContract, *f.now) })
	}); e == nil {
		t.Fatal("R4-like continuation bypassed parent reservation")
	}
	f.tick(t)
	if e := taskrun.CheckParentAction(f.c.Workspace, f.b.ControlContract, "resume", *f.now); e == nil {
		t.Fatal("child resume bypassed paused parent")
	}
	if _, e := f.m.Control(f.c.Workspace, f.b.ID, "resume"); e != nil {
		t.Fatal(e)
	}
	f.tick(t)
	if e := taskrun.CheckParentAction(f.c.Workspace, f.b.ControlContract, "resume", *f.now); e != nil {
		t.Fatal("parent resume did not grant child permit", e)
	}
	g, e := taskrun.ReadParentControl(f.c.Workspace, f.b.ControlContract)
	if e != nil || g.Deadline != f.b.Deadline {
		t.Fatal("permit changed authority", g, e)
	}
	for _, first := range []string{"public", "legacy"} {
		t.Run(first, func(t *testing.T) {
			f := newMR05Week(t)
			r := weekPreview(t, f, "pause", "legacy-race-0001")
			if first == "public" {
				weekApply(t, f, r)
				if _, e := f.m.Control(f.c.Workspace, f.b.ID, "stop"); e == nil {
					t.Fatal("legacy overwrote public request")
				}
			} else {
				if _, e := f.m.Control(f.c.Workspace, f.b.ID, "stop"); e != nil {
					t.Fatal(e)
				}
				if _, e := f.m.Controls().Apply(r); e == nil {
					t.Fatal("public overwrote legacy request")
				}
			}
		})
	}
	t.Run("unknown_dispatch_survives_restart", func(t *testing.T) {
		f := newMR05Week(t)
		f.m.Fault = func(p string) error {
			if p == "week_after_dispatch_reservation" {
				return errors.New("crash before spawn")
			}
			return nil
		}
		if e := f.m.Tick(f.b, f.dir, &f.s); e == nil {
			t.Fatal("barrier missed")
		}
		f.m.Fault = nil
		f.reload(t)
		if e := f.m.Tick(f.b, f.dir, &f.s); e == nil || !strings.Contains(e.Error(), "control_outcome_unknown") {
			t.Fatal("unobserved start replayed", e)
		}
		if len(f.n.starts) != 0 {
			t.Fatal("second writer")
		}
	})
}
