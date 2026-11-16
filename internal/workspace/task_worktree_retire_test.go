package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type taskRetireFixture struct {
	f        workItemJourneyFixture
	epic     string
	worktree string
}

func newTaskRetireFixture(t *testing.T) taskRetireFixture {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	epic := filepath.Join(f.wrapper, "epic")
	in, err := ParseTaskCreateInput("retire-task", "Retire", "Retire task", "epic", "ply", "ply")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = CreateTask(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	worktree := filepath.Join(f.wrapper, "retire-task")
	if _, err = CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "retire-task", Branch: "retire-task", Path: worktree, ExpectedParentOID: runLocalGit(t, epic, "rev-parse", "HEAD")}); err != nil {
		t.Fatal(err)
	}
	worktree = physicalPath(t, worktree)
	if err = os.WriteFile(filepath.Join(worktree, "work.txt"), []byte("task work\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, worktree, "add", "work.txt")
	runLocalGit(t, worktree, "commit", "-m", "task work")
	return taskRetireFixture{f: f, epic: physicalPath(t, epic), worktree: worktree}
}

func (x taskRetireFixture) integrate(t *testing.T) {
	t.Helper()
	runLocalGit(t, x.epic, "merge", "--ff-only", "retire-task")
}

func (x taskRetireFixture) state(t *testing.T) WorkItemState {
	t.Helper()
	r, err := x.f.dependencies.WorkItems.Snapshot(testWork(t, filepath.Dir(x.f.wrapper)))
	if err != nil {
		t.Fatal(err)
	}
	task, _ := findTask(r, "retire-task")
	return task.WorktreeState
}

func (x taskRetireFixture) branchExists(t *testing.T) bool {
	t.Helper()
	return runSafeIntegrationGit([]string{"-C", x.epic, "show-ref", "--verify", "--quiet", "refs/heads/retire-task"}, nil).Exit == 0
}

func requireWorkClass(t *testing.T, err error, class WorkItemErrorClass) {
	t.Helper()
	var work *WorkItemError
	if !errors.As(err, &work) || work.Class != class {
		t.Fatalf("want %s, got %v", class, err)
	}
}

func TestRetireTaskWorktreeRemovesIntegratedWorktreeAndBranch(t *testing.T) {
	x := newTaskRetireFixture(t)
	x.integrate(t)
	check, err := RetireTaskWorktree(x.f.dependencies, TaskWorktreeRetireInput{TaskID: "retire-task"})
	if err != nil || check.Applied || check.Worktree != "would_remove" || check.Branch != "would_delete" || check.State != WorkItemReady {
		t.Fatalf("check %+v %v", check, err)
	}
	if _, err = os.Stat(x.worktree); err != nil || !x.branchExists(t) {
		t.Fatal("check changed Git")
	}
	applied, err := RetireTaskWorktree(x.f.dependencies, TaskWorktreeRetireInput{TaskID: "retire-task", Apply: true})
	if err != nil || !applied.Applied || applied.Worktree != "removed" || applied.Branch != "deleted" || applied.State != WorkItemRetired {
		t.Fatalf("apply %+v %v", applied, err)
	}
	if _, err = os.Stat(x.worktree); !os.IsNotExist(err) || x.branchExists(t) || x.state(t) != WorkItemRetired {
		t.Fatal("worktree, branch or state not retired")
	}
	again, err := RetireTaskWorktree(x.f.dependencies, TaskWorktreeRetireInput{TaskID: "retire-task", Apply: true})
	if err != nil || again.Worktree != "absent" || again.Branch != "absent" || again.State != WorkItemRetired {
		t.Fatalf("repeat %+v %v", again, err)
	}
}

func TestRetireTaskWorktreeRefusesUnintegratedOrDirtyWork(t *testing.T) {
	x := newTaskRetireFixture(t)
	_, err := RetireTaskWorktree(x.f.dependencies, TaskWorktreeRetireInput{TaskID: "retire-task", Apply: true})
	requireWorkClass(t, err, ErrorWorkRefConflict)
	x.integrate(t)
	if err = os.WriteFile(filepath.Join(x.worktree, "unsaved.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = RetireTaskWorktree(x.f.dependencies, TaskWorktreeRetireInput{TaskID: "retire-task", Apply: true})
	requireWorkClass(t, err, ErrorWorkParentDirty)
	if _, err = os.Stat(x.worktree); err != nil || !x.branchExists(t) || x.state(t) != WorkItemReady {
		t.Fatal("refused retirement changed something")
	}
}

func TestRetireTaskWorktreeFinishesAfterManualWorktreeRemoval(t *testing.T) {
	x := newTaskRetireFixture(t)
	x.integrate(t)
	runLocalGit(t, x.epic, "worktree", "remove", x.worktree)
	applied, err := RetireTaskWorktree(x.f.dependencies, TaskWorktreeRetireInput{TaskID: "retire-task", Apply: true})
	if err != nil || applied.Worktree != "absent" || applied.Branch != "deleted" || x.branchExists(t) || x.state(t) != WorkItemRetired {
		t.Fatalf("recovery %+v %v", applied, err)
	}
}
