package weekrun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Synthetic native adapter: records only effects at the public native boundary.
// It never claims provider execution or native product qualification.
type fakeNative struct {
	c                  Contract
	runs               map[workspace.TaskID]Run
	starts             []workspace.TaskID
	controls           []string
	resumes            []bool
	released           map[workspace.TaskID]bool
	current            workspace.TaskID
	base               int
	extra              bool
	block              map[workspace.TaskID]string
	startAfterEffect   error
	releaseAfterEffect error
}

func (f *fakeNative) Queue(c Contract) (workspace.WorkspaceTaskQueueReadback, error) {
	q := workspace.WorkspaceTaskQueueReadback{Workspace: workspace.IntegrationPlanWorkspace{Root: c.Workspace}, Target: workspace.QueueTarget{ProjectID: c.ProjectID, RepoID: c.RepoID, EpicID: c.EpicID, ParentLocator: c.ReturnWorktree, ParentRef: c.ReturnRef}, Pending: []workspace.QueuePending{}}
	if f.current != "" {
		q.Current = &workspace.QueueCurrent{TaskID: f.current, PreparationID: "pre_" + string(f.current)}
	}
	model, effort, provider := c.Model, c.Effort, "codex"
	for _, s := range c.Selected {
		if _, started := f.runs[s.TaskID]; started {
			continue
		}
		g := s.Goal
		r := workspace.QueuePending{TaskID: s.TaskID, Goal: &g, State: "ready", Executor: &workspace.TaskExecutorAssignment{Model: &model, Effort: &effort, Provider: &provider}, Delivery: automaticAgreement()}
		if reason := f.block[s.TaskID]; reason != "" {
			r.State = "blocked"
			r.Reasons = []workspace.QueueReason{{Code: "synthetic_block", Message: reason}}
		}
		q.Pending = append(q.Pending, r)
	}
	if f.extra {
		q.Pending = append([]workspace.QueuePending{{TaskID: "unselected", State: "ready"}}, q.Pending...)
	}
	return q, nil
}
func automaticAgreement() *workspace.DeliveryAgreement {
	return &workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{Mode: "automatic"}}
}
func (f *fakeNative) Preview(c Contract, s Selection, _ taskrun.Executable) error {
	if f.current != "" {
		return errors.New("current conflict")
	}
	return nil
}
func (f *fakeNative) Start(c Contract, s Selection, _ taskrun.Executable, check bool) (Run, error) {
	if r, ok := f.runs[s.TaskID]; ok {
		return r, nil
	}
	if check {
		return Run{}, nil
	}
	if f.current != "" {
		return Run{}, fmt.Errorf("duplicate writer while %s active", f.current)
	}
	r := Run{ID: "run-" + string(s.TaskID), Phase: "running", OwnerState: "lock_held", SupervisorState: "lock_held", BaseOID: fmt.Sprintf("base-%d", f.base), BaseRevision: f.base}
	f.current = s.TaskID
	f.runs[s.TaskID] = r
	f.starts = append(f.starts, s.TaskID)
	if f.startAfterEffect != nil {
		e := f.startAfterEffect
		f.startAfterEffect = nil
		return Run{}, e
	}
	return r, nil
}
func (f *fakeNative) Observe(c Contract, id string) (Run, error) {
	for _, r := range f.runs {
		if r.ID == id {
			return r, nil
		}
	}
	return Run{}, errors.New("native run missing")
}
func (f *fakeNative) Resume(c Contract, id string, explicit bool) (Run, error) {
	f.resumes = append(f.resumes, explicit)
	for task, r := range f.runs {
		if r.ID == id {
			r.OwnerState = "lock_held"
			r.SupervisorState = "lock_held"
			if explicit {
				r.Phase = "running"
				r.ControlAction = ""
			}
			f.runs[task] = r
			return r, nil
		}
	}
	return Run{}, errors.New("missing")
}
func (f *fakeNative) Control(c Contract, id, action string) (Run, error) {
	f.controls = append(f.controls, action)
	for task, r := range f.runs {
		if r.ID == id {
			r.ControlAction = action
			f.runs[task] = r
			return r, nil
		}
	}
	return Run{}, errors.New("missing")
}
func (f *fakeNative) Release(c Contract, id string, task workspace.TaskID) error {
	if f.released[task] {
		return nil
	}
	r := f.runs[task]
	if !r.Releasable || r.OwnerState != "unowned" {
		return errors.New("not positively parked")
	}
	f.released[task] = true
	f.current = ""
	if f.releaseAfterEffect != nil {
		e := f.releaseAfterEffect
		f.releaseAfterEffect = nil
		return e
	}
	return nil
}
func (f *fakeNative) complete(task workspace.TaskID) {
	r := f.runs[task]
	r.Completed = true
	r.Phase = "completed"
	r.OwnerState = "unowned"
	r.SupervisorState = "unowned"
	f.runs[task] = r
	f.current = ""
	f.base++
}

type fixture struct {
	m   Manager
	c   Contract
	b   Binding
	s   State
	dir string
	n   *fakeNative
	now *time.Time
}

func newFixture(t *testing.T, selections ...Selection) *fixture {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ret := filepath.Join(root, "epic")
	if err = os.Mkdir(ret, 0700); err != nil {
		t.Fatal(err)
	}
	c := Contract{Kind: "ply.workflow.week-contract", SchemaVersion: 1, Name: "test-week", Workspace: root, ReturnWorktree: ret, ReturnRef: "refs/heads/epic", ProjectID: "test", RepoID: "test", EpicID: "epic", Selected: selections, Model: "gpt-6-astra", Effort: "max", DurationSeconds: 3600, MaxAttempts: 3, AgentTimeoutSeconds: 300, MaxProviderRetries: 20}
	n := &fakeNative{c: c, runs: map[workspace.TaskID]Run{}, released: map[workspace.TaskID]bool{}, block: map[workspace.TaskID]string{}, base: 1}
	executable := filepath.Join(root, "synthetic-control")
	if err = os.WriteFile(executable, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	m := Manager{Native: n, Now: func() time.Time { return now }, Executable: func() (string, error) { return executable, nil }}
	s, err := m.Start(c, false)
	if err != nil {
		t.Fatal(err)
	}
	b, dir, err := readBinding(root, s.ID)
	if err != nil {
		t.Fatal(err)
	}
	return &fixture{m: m, c: c, b: b, dir: dir, s: s, n: n, now: &now}
}
func selection(id string, deps ...string) Selection {
	s := Selection{TaskID: workspace.TaskID(id), Goal: workspace.TaskGoalRef{SpecID: id + "-v1", Spec: workspace.TaskRevisionRef{Revision: 1, ManifestSHA256: "sha256:" + strings.Repeat("a", 64)}}, DependsOn: []workspace.TaskID{}}
	for _, d := range deps {
		s.DependsOn = append(s.DependsOn, workspace.TaskID(d))
	}
	return s
}
func (f *fixture) tick(t *testing.T) {
	t.Helper()
	*f.now = f.now.Add(time.Second)
	if err := f.m.Tick(f.b, f.dir, &f.s); err != nil {
		t.Fatal(err)
	}
}
func (f *fixture) reload(t *testing.T) {
	t.Helper()
	s, err := readState(f.b, f.dir)
	if err != nil {
		t.Fatal(err)
	}
	f.s = s
}
func (f *fixture) control(t *testing.T, a string) {
	t.Helper()
	*f.now = f.now.Add(time.Second)
	if _, err := f.m.Control(f.c.Workspace, f.s.ID, a); err != nil {
		t.Fatal(err)
	}
}

func TestWeekSelectedDependenciesUseFreshNativeBase(t *testing.T) {
	f := newFixture(t, selection("first"), selection("second", "first"))
	f.n.extra = true
	f.tick(t)
	if len(f.n.starts) != 0 || f.s.ActiveTaskID != "first" {
		t.Fatal("intent must precede launch", f.s)
	}
	f.tick(t)
	for i := 0; i < 4; i++ {
		f.tick(t)
	}
	if !reflect.DeepEqual(f.n.starts, []workspace.TaskID{"first"}) {
		t.Fatal(f.n.starts)
	}
	f.n.complete("first")
	f.tick(t)
	f.reload(t)
	f.tick(t)
	f.tick(t)
	if !reflect.DeepEqual(f.n.starts, []workspace.TaskID{"first", "second"}) {
		t.Fatal(f.n.starts)
	}
	if f.s.Entries[1].BaseOID != "base-2" || f.s.Entries[1].BaseRevision != 2 {
		t.Fatal("second Goal did not bind the updated native base", f.s.Entries)
	}
	f.n.complete("second")
	f.tick(t)
	f.tick(t)
	if f.s.Phase != "completed" || len(f.n.starts) != 2 {
		t.Fatal(f.s)
	}
	f.tick(t)
	if len(f.n.starts) != 2 {
		t.Fatal("finished queue invented new work")
	}
}

func TestWeekConcreteBlockSkipsDependentsAndContinuesIndependent(t *testing.T) {
	f := newFixture(t, selection("blocked"), selection("dependent", "blocked"), selection("independent"))
	f.tick(t)
	f.tick(t)
	r := f.n.runs["blocked"]
	r.Phase = "blocked"
	r.OwnerState = "unowned"
	r.Reason = "developer blocked: missing authorized requirement"
	r.Releasable = true
	f.n.runs["blocked"] = r
	f.tick(t)
	if !f.n.released["blocked"] {
		t.Fatal("native preparation was not released")
	}
	f.tick(t)
	f.tick(t)
	if !reflect.DeepEqual(f.n.starts, []workspace.TaskID{"blocked", "independent"}) {
		t.Fatal(f.n.starts)
	}
	if f.s.Entries[1].State != "dependency_blocked" {
		t.Fatal(f.s.Entries)
	}
	f.n.complete("independent")
	f.tick(t)
	f.tick(t)
	if f.s.Phase != "completed_with_blocks" {
		t.Fatal(f.s)
	}
}
func TestWeekNativePreflightBlockDoesNotReserveRun(t *testing.T) {
	f := newFixture(t, selection("blocked"), selection("dependent", "blocked"), selection("independent"))
	f.n.block["blocked"] = "problem changed"
	f.tick(t)
	f.tick(t)
	f.tick(t)
	if !reflect.DeepEqual(f.n.starts, []workspace.TaskID{"independent"}) {
		t.Fatal(f.n.starts)
	}
}
func TestWeekLostStartReplyRecoversSameRunWithoutDuplicate(t *testing.T) {
	f := newFixture(t, selection("first"))
	f.n.startAfterEffect = errors.New("synthetic controller lost native reply")
	f.tick(t)
	f.tick(t)
	if f.s.Phase != "blocked" || f.s.Entries[0].RunID != "" || len(f.n.starts) != 1 {
		t.Fatal(f.s)
	}
	f.reload(t)
	f.control(t, "resume")
	f.tick(t)
	if len(f.n.starts) != 1 || f.s.Entries[0].RunID != "run-first" {
		t.Fatal("native reservation was duplicated", f.s, f.n.starts)
	}
}
func TestWeekCrashBeforeRunObservationReusesIntent(t *testing.T) {
	f := newFixture(t, selection("first"))
	f.tick(t)
	if _, err := f.n.Start(f.c, f.c.Selected[0], f.b.Control, false); err != nil {
		t.Fatal(err)
	}
	f.reload(t)
	f.tick(t)
	if len(f.n.starts) != 1 || f.s.Entries[0].RunID != "run-first" {
		t.Fatal(f.s)
	}
}
func TestWeekReleaseReplyLossUsesNativeDisposition(t *testing.T) {
	f := newFixture(t, selection("first"), selection("next"))
	f.tick(t)
	f.tick(t)
	r := f.n.runs["first"]
	r.Phase = "blocked"
	r.OwnerState = "unowned"
	r.Releasable = true
	r.Reason = "worker reached its explicit limit"
	f.n.runs["first"] = r
	f.n.releaseAfterEffect = errors.New("synthetic release reply lost")
	f.tick(t)
	if !f.n.released["first"] || f.s.Phase != "blocked" {
		t.Fatal(f.s)
	}
	f.reload(t)
	f.control(t, "resume")
	f.tick(t)
	f.tick(t)
	f.tick(t)
	if !reflect.DeepEqual(f.n.starts, []workspace.TaskID{"first", "next"}) {
		t.Fatal(f.n.starts)
	}
}
func TestWeekBackoffSurvivesRestartWithoutSupervisorSpin(t *testing.T) {
	f := newFixture(t, selection("first"))
	f.tick(t)
	f.tick(t)
	r := f.n.runs["first"]
	r.Phase = "waiting"
	r.OwnerState = "unowned"
	r.SupervisorState = "lock_held"
	r.RetryNotBefore = stamp(f.now.Add(time.Hour))
	r.Reason = "provider usage reset"
	f.n.runs["first"] = r
	for i := 0; i < 10; i++ {
		f.reload(t)
		f.tick(t)
	}
	if f.s.Phase != "waiting" || f.s.NextCheckAt != r.RetryNotBefore || len(f.n.resumes) != 0 || len(f.n.starts) != 1 {
		t.Fatal(f.s, f.n.resumes, f.n.starts)
	}
	r.SupervisorState = "unowned"
	f.n.runs["first"] = r
	f.tick(t)
	if !reflect.DeepEqual(f.n.resumes, []bool{false}) || len(f.n.starts) != 1 {
		t.Fatal("must recover only the same native supervisor", f.n.resumes)
	}
}
func TestWeekAuthAndUnknownEffectsDoNotReleaseOrContinue(t *testing.T) {
	for _, reason := range []string{"authentication required", "attempt has no durable result"} {
		t.Run(reason, func(t *testing.T) {
			f := newFixture(t, selection("first"), selection("independent"))
			f.tick(t)
			f.tick(t)
			r := f.n.runs["first"]
			r.Phase = "blocked"
			r.Reason = reason
			r.OwnerState = "unowned"
			f.n.runs["first"] = r
			f.tick(t)
			f.reload(t)
			f.tick(t)
			if f.s.Phase != "blocked" || len(f.n.released) != 0 || len(f.n.starts) != 1 || len(f.n.resumes) != 0 {
				t.Fatal(f.s)
			}
			f.control(t, "resume")
			f.tick(t)
			if !reflect.DeepEqual(f.n.resumes, []bool{true}) {
				t.Fatal("explicit resume was not forwarded")
			}
		})
	}
}
func TestWeekPauseStopResumeKeepActiveIdentityAndDeadline(t *testing.T) {
	f := newFixture(t, selection("first"), selection("second"))
	f.tick(t)
	f.tick(t)
	deadline := f.s.Deadline
	for _, action := range []string{"pause", "stop"} {
		f.control(t, action)
		f.tick(t)
		want := "pausing"
		phase := "paused"
		if action == "stop" {
			want = "stopping"
			phase = "stopped"
		}
		if f.s.Phase != want {
			t.Fatal(f.s)
		}
		r := f.n.runs["first"]
		r.Phase = phase
		r.OwnerState = "unowned"
		r.SupervisorState = "unowned"
		f.n.runs["first"] = r
		f.tick(t)
		f.reload(t)
		if f.s.Phase != phase || f.s.Entries[0].RunID != "run-first" || f.s.Deadline != deadline || len(f.n.starts) != 1 {
			t.Fatal(f.s)
		}
		f.control(t, "resume")
		f.tick(t)
		if f.s.Phase != "running" || f.s.Entries[0].RunID != "run-first" || f.s.Deadline != deadline {
			t.Fatal(f.s)
		}
	}
	if !reflect.DeepEqual(f.n.controls, []string{"pause", "stop"}) || !reflect.DeepEqual(f.n.resumes, []bool{true, true}) {
		t.Fatal(f.n.controls, f.n.resumes)
	}
}
func TestWeekDeadlineStopsActiveWorkAndNeverExtends(t *testing.T) {
	f := newFixture(t, selection("first"), selection("second"))
	f.tick(t)
	f.tick(t)
	*f.now = f.now.Add(2 * time.Hour)
	f.tick(t)
	if f.s.Phase != "stopping" || !reflect.DeepEqual(f.n.controls, []string{"stop"}) {
		t.Fatal(f.s)
	}
	r := f.n.runs["first"]
	r.Phase = "stopped"
	r.OwnerState = "unowned"
	f.n.runs["first"] = r
	f.tick(t)
	if f.s.Phase != "expired" {
		t.Fatal(f.s)
	}
	f.control(t, "resume")
	f.tick(t)
	if len(f.n.starts) != 1 || f.s.Phase != "expired" {
		t.Fatal(f.s)
	}
}
func TestWeekEmptySelectionCompletesWithoutProductEffects(t *testing.T) {
	f := newFixture(t)
	f.tick(t)
	if f.s.Phase != "completed" || len(f.n.starts) != 0 {
		t.Fatal(f.s)
	}
}
func TestWeekRejectsChangedImmutableBindingAndUnselectedDependencies(t *testing.T) {
	f := newFixture(t, selection("first"))
	c := f.c
	c.DurationSeconds = MaxDurationSeconds + 1
	if err := Validate(c); err == nil {
		t.Fatal("allowed beyond seven days")
	}
	c = f.c
	c.Selected = []Selection{selection("first", "unknown")}
	if err := Validate(c); err == nil {
		t.Fatal("allowed unselected dependency")
	}
	c.Selected = []Selection{selection("first", "second"), selection("second", "first")}
	if err := Validate(c); err == nil {
		t.Fatal("allowed dependency cycle")
	}
	c = f.c
	c.Model = "different"
	if _, err := f.m.Start(c, false); err == nil {
		t.Fatal("rebound existing week")
	}
}
func TestWeekCheckAndShowAreReadOnly(t *testing.T) {
	f := newFixture(t, selection("first"))
	if _, err := f.m.Show(f.c.Workspace, f.s.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(filepath.Join(f.dir, "owner.lock")); !os.IsNotExist(err) {
		t.Fatal("show created lock")
	}
	c := f.c
	c.Name = "preview"
	if _, err := f.m.Start(c, true); err != nil {
		t.Fatal(err)
	}
	dir, _ := runDir(c.Workspace, contractID(c))
	if _, err := os.Lstat(dir); !os.IsNotExist(err) {
		t.Fatal("preview wrote state")
	}
}
func TestWeekFrozenExecutableChangeBlocksRecovery(t *testing.T) {
	f := newFixture(t)
	if err := os.Chmod(f.b.Control.Path, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f.b.Control.Path, []byte("changed"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := f.m.Serve(context.Background(), f.c.Workspace, f.s.ID); err == nil {
		t.Fatal("executed changed control")
	}
}
func TestWeekOwnedControllerDoesNotDuplicateEffects(t *testing.T) {
	f := newFixture(t, selection("first"))
	held, err := lock(filepath.Join(f.dir, "owner.lock"))
	if err != nil {
		t.Fatal(err)
	}
	defer held.Close()
	if err = f.m.Serve(context.Background(), f.c.Workspace, f.s.ID); err != nil {
		t.Fatal(err)
	}
	if len(f.n.starts) != 0 {
		t.Fatal("duplicate writer")
	}
}

func TestWeekStartupRecoversBindingBeforeInitialStateWithoutExtendingDeadline(t *testing.T) {
	f := newFixture(t, selection("first"))
	deadline := f.b.Deadline
	started := f.b.StartedAt
	if err := os.Remove(filepath.Join(f.dir, "state.json")); err != nil {
		t.Fatal(err)
	}
	*f.now = f.now.Add(10 * time.Minute)
	s, err := f.m.Start(f.c, false)
	if err != nil {
		t.Fatal(err)
	}
	if s.Deadline != deadline || s.StartedAt != started || s.Entries[0].State != "pending" || len(f.n.starts) != 0 {
		t.Fatal("startup recovery changed authority", s)
	}
}
func TestWeekUnboundStartupDirectoryDoesNotPoisonOtherNames(t *testing.T) {
	f := newFixture(t)
	f.tick(t)
	if err := directory(filepath.Join(filepath.Dir(f.dir), "week_deadbeef")); err != nil {
		t.Fatal(err)
	}
	c := f.c
	c.Name = "another-selected-week"
	if _, err := f.m.Start(c, false); err != nil {
		t.Fatal("unbound directory poisoned native selection", err)
	}
}

type queueControlNative struct {
	Native
	once func()
}

func (n *queueControlNative) Queue(c Contract) (workspace.WorkspaceTaskQueueReadback, error) {
	q, err := n.Native.Queue(c)
	if n.once != nil {
		fn := n.once
		n.once = nil
		fn()
	}
	return q, err
}
func TestWeekControlAtControllerExitIsNotLost(t *testing.T) {
	f := newFixture(t)
	n := &queueControlNative{Native: f.n}
	f.m.Native = n
	f.m.Launch = func(b Binding) error {
		if ownerState(f.dir) != "lock_held" {
			t.Fatal("repro requires existing controller at launch check")
		}
		return nil
	}
	n.once = func() { f.control(t, "stop") }
	if err := f.m.Serve(context.Background(), f.c.Workspace, f.s.ID); err != nil {
		t.Fatal(err)
	}
	f.reload(t)
	if f.s.Phase != "stopped" || f.s.DesiredAction != "stop" {
		t.Fatal("control was lost between terminal decision and owner release", f.s)
	}
}
func TestWeekShowProjectsNativeStopAfterServiceUnloads(t *testing.T) {
	f := newFixture(t, selection("first"))
	f.tick(t)
	f.tick(t)
	f.control(t, "stop")
	f.tick(t)
	before, err := os.ReadFile(filepath.Join(f.dir, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	r := f.n.runs["first"]
	r.Phase = "stopped"
	r.OwnerState = "unowned"
	f.n.runs["first"] = r
	s, err := f.m.Show(f.c.Workspace, f.s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if s.Phase != "stopped" || s.Entries[0].RunID != "run-first" {
		t.Fatal("readback hides native confirmed stop", s)
	}
	after, _ := os.ReadFile(filepath.Join(f.dir, "state.json"))
	if string(before) != string(after) {
		t.Fatal("show wrote operational state")
	}
}

func TestWeekConcreteBlockWaitsForNativeOwnerBeforeRelease(t *testing.T) {
	f := newFixture(t, selection("first"), selection("independent"))
	f.tick(t)
	f.tick(t)
	r := f.n.runs["first"]
	r.Phase = "blocked"
	r.Releasable = true
	r.Reason = "worker reached its explicit limit"
	r.OwnerState = "lock_held"
	f.n.runs["first"] = r
	f.tick(t)
	if f.s.Phase == "blocked" || len(f.n.released) != 0 {
		t.Fatal("ordinary native lock release boundary became a terminal blocker", f.s)
	}
	r.OwnerState = "unowned"
	f.n.runs["first"] = r
	f.tick(t)
	f.tick(t)
	f.tick(t)
	if !reflect.DeepEqual(f.n.starts, []workspace.TaskID{"first", "independent"}) {
		t.Fatal("required mechanical resume after ordinary native exit", f.n.starts)
	}
}

func TestWeekRestoresLostSupervisorForPreservedStopRequest(t *testing.T) {
	f := newFixture(t, selection("first"))
	f.tick(t)
	f.tick(t)
	f.control(t, "stop")
	f.tick(t)
	r := f.n.runs["first"]
	r.ControlAction = "stop"
	r.OwnerState = "unowned"
	r.SupervisorState = "unowned"
	f.n.runs["first"] = r
	f.reload(t)
	f.tick(t)
	f.tick(t)
	if !reflect.DeepEqual(f.n.resumes, []bool{false}) || len(f.n.starts) != 1 || f.s.Entries[0].RunID != "run-first" {
		t.Fatal("persisted stop lost native supervision", f.n.resumes, f.s)
	}
}

func TestWeekMissingStateAfterExecutionCannotRecreateSelection(t *testing.T) {
	f := newFixture(t, selection("first"))
	f.tick(t)
	f.tick(t)
	if err := os.Remove(filepath.Join(f.dir, "state.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := f.m.Start(f.c, false); err == nil {
		t.Fatal("lost operational state silently recreated an active selection")
	}
	if len(f.n.starts) != 1 {
		t.Fatal("recovery duplicated native work")
	}
}
