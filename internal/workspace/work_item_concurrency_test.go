package workspace

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConcurrentTaskWorktreeCreateSerializesIdenticalAndConflictingIntents(t *testing.T) {
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
		if got[0].Task.Worktree == nil || got[1].Task.Worktree == nil || got[0].Task.Worktree.ID != got[1].Task.Worktree.ID || got[0].Operation.ID != got[1].Operation.ID {
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
				} else if hasWorkErrorClass(err, ErrorWorkIdentityConflict) {
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
	return newWorkItemJourneyFixtureVersion(t, false)
}

func newWorkItemJourneyFixtureVersion(t *testing.T, legacy bool) workItemJourneyFixture {
	return newWorkItemJourneyFixtureHistory(t, legacy, false)
}

func newWorkItemJourneyFixtureHistory(t *testing.T, legacy, priorEpicWork bool) workItemJourneyFixture {
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
	runLocalGit(t, repository, "worktree", "add", "-b", "epic", epicPath, oid)
	epicPath = physicalPath(t, epicPath)
	if priorEpicWork {
		if err := os.WriteFile(filepath.Join(epicPath, "prior.txt"), []byte("prior Epic work\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		runLocalGit(t, epicPath, "add", "prior.txt")
		runLocalGit(t, epicPath, "commit", "-m", "prior Epic work")
		oid = runLocalGit(t, epicPath, "rev-parse", "HEAD")
	}
	projects := newSystemProjectStore()
	project := ProjectRecord{ID: "ply", Name: "Ply", Wrapper: physicalPath(t, wrapper), RepoIDs: []RepoID{"ply"}}
	repo := RepoRecord{ID: "ply", Locator: repository, GitCommonDir: common}
	if _, _, _, err := projects.Add(root, project, []RepoRecord{repo}); err != nil {
		t.Fatal(err)
	}
	dependencies := SystemDependencies()
	dependencies.Files = projectCwdFileSystem{FileSystem: dependencies.Files, cwd: root}
	dependencies.ProjectLocks = projects
	dependencies.Projects = projects
	adopt, _ := ParseEpicAdoptInput("epic", "Epic", "ply", "ply", epicPath, "refs/heads/epic", oid)
	if _, err := AdoptEpic(dependencies, adopt); err != nil {
		t.Fatal(err)
	}
	if legacy {
		// Seed an explicit pre-WS05 registry; do not weaken the policy of new Tasks.
		err := dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, err := session.Snapshot()
			if err != nil {
				return err
			}
			r.FormatVersion = 1
			r.TaskSpecPolicies = nil
			r.TaskProblemRevisions = nil
			r.TaskSpecRevisions = nil
			r.TaskSpecAssessments = nil
			r.TaskSolutionSelections = nil
			r.TaskResultSpecBindings = nil
			r.TaskContentPublications = nil
			r.TaskResults = nil
			r.HumanQARecords = nil
			r.IntegrationAuthorities = nil
			r.IntegrationIntents = nil
			r.IntegrationAttempts = nil
			r.IntegrationResults = nil
			r.Tasks = append(r.Tasks, TaskRecord{ID: "task", Title: "Task", Description: "Description", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply", GitCommonDir: common, WorktreeState: WorkItemUnbound})
			return session.Publish(r)
		})
		if err != nil {
			t.Fatal(err)
		}
	} else {
		task, _ := ParseTaskCreateInput("task", "Task", "Description", "epic", "ply", "ply")
		if _, err := CreateTask(dependencies, task); err != nil {
			t.Fatal(err)
		}
	}
	return workItemJourneyFixture{dependencies: dependencies, wrapper: wrapper, oid: oid}
}

func hasWorkErrorClass(err error, class WorkItemErrorClass) bool {
	typed, ok := err.(*WorkItemError)
	return ok && typed.Class == class
}
