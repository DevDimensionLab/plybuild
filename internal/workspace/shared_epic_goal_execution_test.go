package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type sharedGoalSnapshotHook struct {
	ProjectSnapshotLocker
	after  func() error
	before bool
}

func (h sharedGoalSnapshotHook) WithSnapshotLock(root string, fn func(ProjectSnapshot) error) error {
	if h.before {
		if err := h.after(); err != nil {
			return err
		}
		return h.ProjectSnapshotLocker.WithSnapshotLock(root, fn)
	}
	if err := h.ProjectSnapshotLocker.WithSnapshotLock(root, fn); err != nil {
		return err
	}
	return h.after()
}

func sharedGoalExecutionFixture(t *testing.T) (workItemJourneyFixture, TaskGoalExecutePreview) {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	goal := goalFixture(t, f, "task", "shared-goal")
	goalQueueFixture(t, f, []TaskGoalQueueEntry{{TaskID: "task", Goal: &goal}})
	epic := filepath.Join(f.wrapper, "epic")
	if _, err := AdoptEpic(f.dependencies, EpicAdoptInput{EpicID: "sibling", Title: "Sibling", ProjectID: "ply", RepoID: "ply", Worktree: epic, Ref: "refs/heads/epic", ExpectedOID: runLocalGit(t, epic, "rev-parse", "HEAD")}); err != nil {
		t.Fatal(err)
	}
	f.dependencies.Files = projectCwdFileSystem{FileSystem: SystemDependencies().Files, cwd: epic}
	p, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{SpecID: goal.SpecID})
	if err != nil || p.Target.EpicID != "epic" || p.Input.Target != (QueueTargetInput{}) {
		t.Fatalf("Spec-only preview = %+v, %v", p, err)
	}
	return f, p
}

func TestSharedEpicUniqueSpecPreparesAndRecoversExactTarget(t *testing.T) {
	for _, failAt := range []int{0, 1, 2, 3, 4} {
		t.Run(fmt.Sprint(failAt), func(t *testing.T) {
			f, p := sharedGoalExecutionFixture(t)
			base := runLocalGit(t, p.Target.ParentLocator, "rev-parse", "HEAD")
			if failAt == 4 {
				// Lose the outer response after the goal has become a selected
				// solution, but before any preparation exists.
				locks := f.dependencies.ProjectLocks
				f.dependencies.ProjectLocks = sharedGoalSnapshotHook{ProjectSnapshotLocker: locks, after: func() error {
					return errors.New("synthetic lost queue-binding response")
				}}
				if _, err := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); err == nil {
					t.Fatal("lost queue response was not reported")
				}
				f.dependencies.ProjectLocks = locks
			} else if failAt != 0 {
				count := 0
				f.dependencies.TaskContent.fault = func(stage string) error {
					if stage == "postread" {
						count++
						if count == failAt {
							return errors.New("synthetic lost publication response")
						}
					}
					return nil
				}
				if _, err := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); err == nil {
					t.Fatal("lost response was not reported")
				}
				f.dependencies.TaskContent.fault = nil
			}
			if failAt != 0 {
				// The CLI repeats only --spec, not the derived acceptance path.
				resume, err := PreviewTaskGoalExecution(f.dependencies, TaskGoalExecuteInput{SpecID: p.Input.SpecID})
				if err != nil || resume.Confirmation != p.Confirmation || resume.Target != p.Target {
					t.Fatalf("retry lost preserved request: %+v, %v", resume, err)
				}
			}
			out, err := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture())
			if err != nil {
				t.Fatal(err)
			}
			plan := out.Preparation.Plan
			if plan == nil || plan.Target != p.Target || plan.TaskID != p.Goal.TaskID || plan.ParentOID != base || plan.WorktreePath != p.WorktreePath || plan.Branch != p.Branch {
				t.Fatalf("wrong preparation: %+v", plan)
			}
			if plan.WorktreePath == p.Target.ParentLocator || runLocalGit(t, plan.WorktreePath, "symbolic-ref", "HEAD") != "refs/heads/"+p.Branch {
				t.Fatal("Spec-only execution did not create its isolated Task copy")
			}
			before := queueStateBytes(t, f)
			again, err := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture())
			if err != nil || again.Preparation.Preparation.ID != out.Preparation.Preparation.ID || queueStateBytes(t, f) != before {
				t.Fatalf("retry changed history or preparation: %+v, %v", again, err)
			}
			r := mustQueueRegistry(t, f)
			if len(r.TaskSpecRevisions) != 2 || len(r.TaskSpecAssessments) != 1 || len(r.TaskSolutionSelections) != 1 || len(r.TaskPreparations) != 1 || len(r.WorktreeOperations) != 1 {
				t.Fatal("execution duplicated content or worktree history")
			}
			if runLocalGit(t, p.Target.ParentLocator, "rev-parse", "HEAD") != base {
				t.Fatal("preparation changed the shared base")
			}
		})
	}
}

func TestSharedEpicUniqueSpecRechecksTargetBeforePreparationEffects(t *testing.T) {
	for _, moveAfter := range []int{1, 2} {
		t.Run(fmt.Sprint(moveAfter), func(t *testing.T) {
			f, p := sharedGoalExecutionFixture(t)
			calls := 0
			f.dependencies.ProjectLocks = sharedGoalSnapshotHook{ProjectSnapshotLocker: f.dependencies.ProjectLocks, before: moveAfter == 2, after: func() error {
				calls++
				if calls == moveAfter {
					runLocalGit(t, p.Target.ParentLocator, "commit", "--allow-empty", "-m", "target moved between preparation phases")
				}
				return nil
			}}
			if _, err := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); err == nil {
				t.Fatal("changed target was accepted")
			}
			r := mustQueueRegistry(t, f)
			if calls < moveAfter || len(r.TaskSpecRevisions) != 2 || len(r.TaskSolutionSelections) != 1 || len(r.TaskPreparations) != 0 || len(r.WorktreeOperations) != 0 {
				t.Fatal("drift was not exercised or preparation created effects")
			}
			if _, err := os.Stat(p.WorktreePath); !os.IsNotExist(err) {
				t.Fatalf("refused preparation created a worktree: %v", err)
			}
			if runLocalGit(t, p.Target.ParentLocator, "rev-parse", "HEAD") == p.ParentOID {
				t.Fatal("refusal reset the moved target")
			}
		})
	}
}

func TestSharedEpicUniqueSpecRejectsMovedTargetOnPreservedRetry(t *testing.T) {
	f, p := sharedGoalExecutionFixture(t)
	f.dependencies.TaskContent.fault = func(stage string) error {
		if stage == "postread" {
			return errors.New("synthetic lost materialization response")
		}
		return nil
	}
	if _, err := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); err == nil {
		t.Fatal("lost response was not reported")
	}
	f.dependencies.TaskContent.fault = nil
	runLocalGit(t, p.Target.ParentLocator, "commit", "--allow-empty", "-m", "shared target moved")
	moved := runLocalGit(t, p.Target.ParentLocator, "rev-parse", "HEAD")
	before := queueStateBytes(t, f)
	if _, err := PrepareTaskGoalExecution(f.dependencies, p.Input, p.Confirmation, goalHumanFixture()); err == nil {
		t.Fatal("preserved request silently changed base")
	}
	if queueStateBytes(t, f) != before || runLocalGit(t, p.Target.ParentLocator, "rev-parse", "HEAD") != moved {
		t.Fatal("refusal changed history or shared target")
	}
	if _, err := os.Stat(p.WorktreePath); !os.IsNotExist(err) {
		t.Fatalf("refusal created a worktree: %v", err)
	}
}
