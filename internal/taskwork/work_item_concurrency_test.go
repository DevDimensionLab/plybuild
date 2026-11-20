package taskwork

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConcurrentTaskWorktreeCreateSerializesIdenticalAndConflictingIntents(t *testing.T) {
	t.Parallel()
	t.Run("identical", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
		start := make(chan struct{})
		results := make(chan TaskWorktreeMutationResult, 2)
		errors := make(chan error, 2)
		for index := 0; index < 2; index++ {
			go func() {
				<-start
				result, err := CreateTaskWorktree(fixture.dependencies, input)
				results <- result
				errors <- err
			}()
		}
		close(start)
		var got []TaskWorktreeMutationResult
		for index := 0; index < 2; index++ {
			select {
			case err := <-errors:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("concurrent create timed out")
			}
			got = append(got, <-results)
		}
		if got[0].Task.Worktree == nil || got[1].Task.Worktree == nil || got[0].Task.Worktree.ID != got[1].Task.Worktree.ID || got[0].Task.Worktree.OperationID != got[1].Task.Worktree.OperationID {
			t.Fatalf("results = %#v", got)
		}
		if !((got[0].Outcome == "created" && got[1].Outcome == "already_ready") || (got[1].Outcome == "created" && got[0].Outcome == "already_ready")) {
			t.Fatalf("outcomes = %q, %q", got[0].Outcome, got[1].Outcome)
		}
	})

	t.Run("different", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		inputs := []TaskWorktreeCreateInput{}
		first, _ := ParseTaskWorktreeCreateInput("task", "task-one", filepath.Join(fixture.wrapper, "task-one"), fixture.oid)
		second, _ := ParseTaskWorktreeCreateInput("task", "task-two", filepath.Join(fixture.wrapper, "task-two"), fixture.oid)
		inputs = append(inputs, first, second)
		start := make(chan struct{})
		errors := make(chan error, 2)
		for _, input := range inputs {
			input := input
			go func() { <-start; _, err := CreateTaskWorktree(fixture.dependencies, input); errors <- err }()
		}
		close(start)
		successes, conflicts := 0, 0
		for index := 0; index < 2; index++ {
			select {
			case err := <-errors:
				if err == nil {
					successes++
				} else if hasWorkErrorClass(err, workspace.ErrorWorkIdentityConflict) {
					conflicts++
				} else {
					t.Fatalf("unexpected error: %v", err)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("concurrent create timed out")
			}
		}
		if successes != 1 || conflicts != 1 {
			t.Fatalf("successes=%d conflicts=%d", successes, conflicts)
		}
	})
}

type workItemJourneyFixture struct {
	dependencies Dependencies
	wrapper      string
	oid          string
}

func newWorkItemJourneyFixture(t *testing.T) workItemJourneyFixture {
	return newWorkItemJourneyFixtureHistory(t, false)
}

func newWorkItemJourneyFixtureHistory(t *testing.T, priorEpicWork bool, parentBranch ...string) workItemJourneyFixture {
	t.Helper()
	root := createProjectWorkspace(t)
	wrapper := filepath.Join(root, "ply")
	repository := filepath.Join(wrapper, "main")
	if err := os.MkdirAll(repository, 0o755); err != nil {
		t.Fatal(err)
	}
	repository = physicalPath(t, repository)
	runLocalGit(t, repository, "init", "-b", "main")
	runLocalGit(t, repository, "config", "user.email", "tests@example.invalid")
	runLocalGit(t, repository, "config", "user.name", "Ply tests")
	if err := os.WriteFile(filepath.Join(repository, "README.md"), []byte("fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, repository, "add", "README.md")
	runLocalGit(t, repository, "commit", "-m", "fixture")
	oid := runLocalGit(t, repository, "rev-parse", "HEAD")
	common := physicalPath(t, filepath.Join(repository, ".git"))
	epicPath := filepath.Join(wrapper, "epic")
	branch := "epic"
	if len(parentBranch) > 0 {
		branch = parentBranch[0]
	}
	if branch == "main" || branch == "master" {
		epicPath = repository
		if branch == "master" {
			runLocalGit(t, repository, "branch", "-m", "master")
		}
	} else {
		runLocalGit(t, repository, "worktree", "add", "-b", branch, epicPath, oid)
	}
	epicPath = physicalPath(t, epicPath)
	if priorEpicWork {
		if err := os.WriteFile(filepath.Join(epicPath, "prior.txt"), []byte("prior Epic work\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		runLocalGit(t, epicPath, "add", "prior.txt")
		runLocalGit(t, epicPath, "commit", "-m", "prior Epic work")
		oid = runLocalGit(t, epicPath, "rev-parse", "HEAD")
	}
	project := workspace.ProjectRecord{ID: "ply", Name: "Ply", Wrapper: physicalPath(t, wrapper), RepoIDs: []workspace.RepoID{"ply"}}
	repo := workspace.RepoRecord{ID: "ply", Locator: repository, GitCommonDir: common}
	defineTestProject(t, root, project, []workspace.RepoRecord{repo})
	dependencies := SystemDependencies()
	dependencies.ProjectDataRoot = func(string) (string, error) { return "", nil }
	dependencies.Files = projectCwdFileSystem{FileSystem: dependencies.Files, cwd: root}
	adopt, _ := ParseEpicAdoptInput("epic", "Epic", "ply", "ply", epicPath, "refs/heads/"+branch, oid)
	if _, err := AdoptEpic(dependencies, adopt); err != nil {
		t.Fatal(err)
	}
	task, _ := ParseTaskCreateInput("task", "Task", "Description", "epic", "ply", "ply")
	if err := fixtureCreateTaskWithProblem(dependencies, root, task); err != nil {
		t.Fatal(err)
	}
	return workItemJourneyFixture{dependencies: dependencies, wrapper: wrapper, oid: oid}
}

func hasWorkErrorClass(err error, class workspace.WorkItemErrorClass) bool {
	typed, ok := err.(*workspace.WorkItemError)
	return ok && typed.Class == class
}

func fixtureCreateTaskWithProblem(d Dependencies, root string, in TaskCreateInput) error {
	work, err := workspace.DefinedProjectWork(root, in.ProjectID)
	if err != nil {
		return err
	}
	return d.ProjectLocks.WithSnapshotLock(root, func(projects workspace.ProjectSnapshot) error {
		_, repo, err := projectAndRepo(projects, in.ProjectID, in.RepoID)
		if err != nil {
			return err
		}
		return d.WorkItems.WithLock(work, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			task := TaskRecord{ID: in.TaskID, Title: in.Title, Description: in.Description, ParentEpicID: in.ParentEpicID, ProjectID: in.ProjectID, RepoID: in.RepoID, GitCommonDir: repo.GitCommonDir, WorktreeState: WorkItemUnbound}
			_, err = createTaskProblem(d, work, session, registry, projects, task)
			return err
		})
	})
}
