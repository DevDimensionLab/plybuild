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
