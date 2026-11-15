package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

type epicRetireFixture struct {
	f          workItemJourneyFixture
	repository string
	epic       string
}

// newEpicRetireFixture has Epic work committed on refs/heads/epic, not yet in main.
func newEpicRetireFixture(t *testing.T) epicRetireFixture {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	x := epicRetireFixture{f: f, repository: physicalPath(t, filepath.Join(f.wrapper, "main")), epic: physicalPath(t, filepath.Join(f.wrapper, "epic"))}
	if err := os.WriteFile(filepath.Join(x.epic, "epic.txt"), []byte("epic work\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, x.epic, "add", "epic.txt")
	runLocalGit(t, x.epic, "commit", "-m", "epic work")
	return x
}

func (x epicRetireFixture) retire(apply bool) (EpicRetireResult, error) {
	return RetireEpic(x.f.dependencies, EpicRetireInput{EpicID: "epic", RepoID: "ply", TargetRef: "refs/heads/main", Actor: "planner-a", Apply: apply})
}

func (x epicRetireFixture) lifecycle(t *testing.T) LifecycleState {
	t.Helper()
	l, err := ShowWorkItemLifecycle(x.f.dependencies, "epic", "epic")
	if err != nil {
		t.Fatal(err)
	}
	return l.Lifecycle
}

func TestRetireEpicRemovesIntegratedWorktreeAndBranchAndArchives(t *testing.T) {
	x := newEpicRetireFixture(t)
	_, err := x.retire(true)
	requireWorkClass(t, err, ErrorWorkRefConflict)
	runLocalGit(t, x.repository, "merge", "--ff-only", "epic")
	check, err := x.retire(false)
	if err != nil || check.Applied || check.Worktree != "would_remove" || check.Branch != "would_delete" {
		t.Fatalf("check %+v %v", check, err)
	}
	if _, err = os.Stat(x.epic); err != nil || x.lifecycle(t) != LifecycleActive {
		t.Fatal("check changed something")
	}
	applied, err := x.retire(true)
	if err != nil || !applied.Applied || applied.Worktree != "removed" || applied.Branch != "deleted" || applied.Lifecycle != LifecycleArchived {
		t.Fatalf("apply %+v %v", applied, err)
	}
	if _, err = os.Stat(x.epic); !os.IsNotExist(err) || x.lifecycle(t) != LifecycleArchived {
		t.Fatal("worktree or lifecycle not retired")
	}
	if runSafeIntegrationGit([]string{"-C", x.repository, "show-ref", "--verify", "--quiet", "refs/heads/epic"}, nil).Exit == 0 {
		t.Fatal("Epic branch remains")
	}
	again, err := x.retire(true)
	if err != nil || again.Worktree != "absent" || again.Branch != "absent" || again.Lifecycle != LifecycleArchived {
		t.Fatalf("repeat %+v %v", again, err)
	}
}

func TestRetireEpicRefusesDirtyWorktree(t *testing.T) {
	x := newEpicRetireFixture(t)
	runLocalGit(t, x.repository, "merge", "--ff-only", "epic")
	if err := os.WriteFile(filepath.Join(x.epic, "unsaved.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := x.retire(true)
	requireWorkClass(t, err, ErrorWorkParentDirty)
	if _, err = os.Stat(x.epic); err != nil || x.lifecycle(t) != LifecycleActive {
		t.Fatal("refused retirement changed something")
	}
}

// An Epic whose Task still has a worktree is not retired; its base equals
// main here, so only the open Task stops it.
func TestRetireEpicRefusesOpenTasks(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	x := epicRetireFixture{f: f, repository: physicalPath(t, filepath.Join(f.wrapper, "main")), epic: physicalPath(t, filepath.Join(f.wrapper, "epic"))}
	in, err := ParseTaskCreateInput("open-task", "Open", "Open task", "epic", "ply", "ply")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTask(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "open-task", Branch: "open-task", Path: filepath.Join(f.wrapper, "open-task"), ExpectedParentOID: runLocalGit(t, x.epic, "rev-parse", "HEAD")}); err != nil {
		t.Fatal(err)
	}
	_, err = x.retire(true)
	requireWorkClass(t, err, ErrorWorkIdentityConflict)
	if _, err = os.Stat(x.epic); err != nil || x.lifecycle(t) != LifecycleActive {
		t.Fatal("refused retirement changed something")
	}
}

func TestRetireEpicRefusesTheTargetItself(t *testing.T) {
	x := newEpicRetireFixture(t)
	_, err := RetireEpic(x.f.dependencies, EpicRetireInput{EpicID: "epic", RepoID: "ply", TargetRef: "refs/heads/epic", Actor: "a", Apply: true})
	requireWorkClass(t, err, ErrorWorkInvalidArguments)
}

// definedEpicRetireFixture has the Epics epic (with work) and second, both
// defined in ply-workspace.json as in real use, and main containing both.
func definedEpicRetireFixture(t *testing.T) epicRetireFixture {
	t.Helper()
	x := newEpicRetireFixture(t)
	base := runLocalGit(t, x.repository, "rev-parse", "HEAD")
	second := filepath.Join(x.f.wrapper, "second")
	runLocalGit(t, x.repository, "worktree", "add", "-b", "second", second, base)
	if _, err := AdoptEpic(x.f.dependencies, EpicAdoptInput{EpicID: "second", Title: "Second", ProjectID: "ply", RepoID: "ply", Worktree: physicalPath(t, second), Ref: "refs/heads/second", ExpectedOID: base}); err != nil {
		t.Fatal(err)
	}
	finalReviewInstallDefinitions(t, x.f)
	runLocalGit(t, x.repository, "merge", "--ff-only", "epic")
	return x
}

// The defect seen in a real workspace: after one Epic was retired, every
// Project path lookup failed. Retiring two defined Epics in turn, and other
// writes in between, must work.
func TestRetireDefinedEpicsInTurnKeepsTheWorkspaceUsable(t *testing.T) {
	x := definedEpicRetireFixture(t)
	if r, err := x.retire(true); err != nil || r.Worktree != "removed" || r.Lifecycle != LifecycleArchived {
		t.Fatalf("first retire %+v %v", r, err)
	}
	in, err := ParseTaskCreateInput("after", "After", "After the first retire", "second", "ply", "ply")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTask(x.f.dependencies, in); err != nil {
		t.Fatalf("a write after the first retire failed: %v", err)
	}
	r, err := RetireEpic(x.f.dependencies, EpicRetireInput{EpicID: "second", RepoID: "ply", TargetRef: "refs/heads/main", Actor: "planner-a", Apply: true})
	if err != nil || r.Worktree != "removed" || r.Lifecycle != LifecycleArchived {
		t.Fatalf("second retire %+v %v", r, err)
	}
	if again, err := x.retire(true); err != nil || again.Worktree != "absent" || again.Branch != "absent" {
		t.Fatalf("rerun %+v %v", again, err)
	}
}

// The state an earlier Ply left: worktree and branch gone, Epic not archived.
// Retire names the recovery; with the worktree back on its path it finishes.
func TestRetireRecoversAWorktreeRemovedBeforeArchiving(t *testing.T) {
	x := definedEpicRetireFixture(t)
	runLocalGit(t, x.repository, "worktree", "remove", x.epic)
	runLocalGit(t, x.repository, "branch", "-D", "epic")
	_, err := x.retire(true)
	requireWorkClass(t, err, ErrorWorkPathConflict)
	runLocalGit(t, x.repository, "worktree", "add", x.epic, "-b", "epic", "refs/heads/main")
	if r, err := x.retire(true); err != nil || r.Worktree != "removed" || r.Branch != "deleted" || r.Lifecycle != LifecycleArchived {
		t.Fatalf("rerun %+v %v", r, err)
	}
}
