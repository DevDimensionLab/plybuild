package workspace

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSystemWorkItemGitObservesAndCreatesOnlyAdditiveWorktrees(t *testing.T) {
	repository := initWorkItemGitRepository(t)
	git := newSystemWorkItemGit(systemFileSystem{})
	repo := RepoRecord{ID: "repo", Locator: repository, GitCommonDir: physicalPath(t, filepath.Join(repository, ".git"))}

	repoObservation, err := git.ObserveRepo(repo)
	if err != nil || repoObservation.Locator != repository || repoObservation.GitCommonDir != repo.GitCommonDir {
		t.Fatalf("repo observation = %#v, %v", repoObservation, err)
	}
	parent, err := git.ObserveWorktree(repository)
	if err != nil || !parent.Clean || !parent.InventoryMatch || parent.Ref != "refs/heads/main" {
		t.Fatalf("parent observation = %#v, %v", parent, err)
	}
	parentRef, err := git.ObserveRef(repo, parent.Ref)
	if err != nil || !parentRef.Exists || parentRef.OID != parent.OID || parentRef.CheckedOutAt != repository {
		t.Fatalf("parent ref = %#v, %v", parentRef, err)
	}

	target := filepath.Join(t.TempDir(), "task-worktree")
	if err := git.ValidateBranch("task"); err != nil {
		t.Fatal(err)
	}
	outcome := git.CreateBranchAndWorktree(GitCreateInput{Repository: repo, Branch: "task", SourceRef: "refs/heads/task", TargetLocator: target, ParentOID: parent.OID})
	if outcome.Err != nil {
		t.Fatalf("create outcome = %#v; stderr=%s", outcome, outcome.Stderr)
	}
	created, err := git.ObserveWorktree(physicalPath(t, target))
	if err != nil || created.Ref != "refs/heads/task" || created.OID != parent.OID || !created.Clean || !created.InventoryMatch {
		t.Fatalf("created observation = %#v, %v", created, err)
	}
}

func TestTaskWorktreeCreateClassifiesGitOutcomesFromObservation(t *testing.T) {
	t.Run("non-zero with exact effect", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		base := fixture.dependencies.WorkGit
		fixture.dependencies.WorkGit = effectThenErrorGit{WorkItemGit: base}
		input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
		result, err := CreateTaskWorktree(fixture.dependencies, input)
		if err != nil || result.Outcome != "recovered" || result.Task.WorktreeState != WorkItemReady {
			t.Fatalf("result = %#v, %v", result, err)
		}
	})

	t.Run("safe no effect retries with persisted IDs", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		base := fixture.dependencies.WorkGit
		fixture.dependencies.WorkGit = noEffectGit{WorkItemGit: base}
		input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
		if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkGitEffect) {
			t.Fatalf("first error = %v", err)
		}
		registry, err := fixture.dependencies.WorkItems.Snapshot(filepath.Dir(fixture.wrapper))
		if err != nil || registry.Tasks[0].WorktreeState != WorkItemCreating || registry.WorktreeOperations[0].LastObservation.Classification != "no_effect" {
			t.Fatalf("creating registry = %#v, %v", registry, err)
		}
		worktreeID := registry.WorktreeOperations[0].WorktreeID
		operationID := registry.WorktreeOperations[0].ID
		fixture.dependencies.WorkGit = base
		result, err := CreateTaskWorktree(fixture.dependencies, input)
		if err != nil || result.Outcome != "recovered" || result.Task.Worktree.ID != worktreeID || result.Operation.ID != operationID {
			t.Fatalf("retry = %#v, %v", result, err)
		}
	})

	t.Run("branch-only retry adds worktree without moving branch", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		base := fixture.dependencies.WorkGit
		fixture.dependencies.WorkGit = branchOnlyGit{WorkItemGit: base}
		input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
		if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkGitEffect) {
			t.Fatalf("first error = %v", err)
		}
		fixture.dependencies.WorkGit = base
		result, err := CreateTaskWorktree(fixture.dependencies, input)
		if err != nil || result.Outcome != "recovered" || result.Operation.LastObservation.Classification != "exact_effect" {
			t.Fatalf("retry = %#v, %v", result, err)
		}
	})

	t.Run("conflicting partial effect requires reconciliation", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		repository := filepath.Join(fixture.wrapper, "main")
		runLocalGit(t, repository, "commit", "--allow-empty", "-m", "different base")
		fixture.dependencies.WorkGit = wrongBranchGit{WorkItemGit: fixture.dependencies.WorkGit}
		input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
		if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkReconciliation) {
			t.Fatalf("error = %v", err)
		}
		registry, err := fixture.dependencies.WorkItems.Snapshot(filepath.Dir(fixture.wrapper))
		if err != nil || registry.Tasks[0].WorktreeState != WorkItemReconciliationRequired || registry.WorktreeOperations[0].State != "reconciliation_required" {
			t.Fatalf("registry = %#v, %v", registry, err)
		}
		if _, err := os.Lstat(filepath.Join(fixture.wrapper, "task")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("partial branch unexpectedly created worktree: %v", err)
		}
	})
}

type effectThenErrorGit struct{ WorkItemGit }

func (git effectThenErrorGit) CreateBranchAndWorktree(input GitCreateInput) GitCommandOutcome {
	outcome := git.WorkItemGit.CreateBranchAndWorktree(input)
	if outcome.Err == nil {
		outcome.Err = errors.New("injected lost response")
		outcome.Exit = 1
	}
	return outcome
}

type noEffectGit struct{ WorkItemGit }

func (git noEffectGit) CreateBranchAndWorktree(GitCreateInput) GitCommandOutcome {
	return GitCommandOutcome{Exit: 1, Err: errors.New("injected no effect")}
}

type branchOnlyGit struct{ WorkItemGit }

func (git branchOnlyGit) CreateBranchAndWorktree(input GitCreateInput) GitCommandOutcome {
	command := exec.Command("git", "-C", input.Repository.Locator, "branch", input.Branch, input.ParentOID)
	command.Env = append(os.Environ(), "LC_ALL=C", "LANG=C")
	output, err := command.CombinedOutput()
	if err != nil {
		return GitCommandOutcome{Exit: 1, Stderr: output, Err: err}
	}
	return GitCommandOutcome{Exit: 1, Err: errors.New("injected branch-only response")}
}

type wrongBranchGit struct{ WorkItemGit }

func (git wrongBranchGit) CreateBranchAndWorktree(input GitCreateInput) GitCommandOutcome {
	command := exec.Command("git", "-C", input.Repository.Locator, "branch", input.Branch, "HEAD")
	command.Env = append(os.Environ(), "LC_ALL=C", "LANG=C")
	output, err := command.CombinedOutput()
	if err != nil {
		return GitCommandOutcome{Exit: 1, Stderr: output, Err: err}
	}
	return GitCommandOutcome{}
}

func TestSystemWorkItemGitWriteArgvDisablesHooksAndNeverForces(t *testing.T) {
	var calls [][]string
	git := &systemWorkItemGit{files: systemFileSystem{}, run: func(arguments []string, environment []string) GitCommandOutcome {
		calls = append(calls, append([]string(nil), arguments...))
		return GitCommandOutcome{}
	}}
	input := GitCreateInput{Repository: RepoRecord{Locator: "/repo"}, Branch: "task/one", TargetLocator: "/target", ParentOID: strings.Repeat("a", 40)}
	if err := git.ValidateBranch(input.Branch); err != nil {
		t.Fatal(err)
	}
	git.CreateBranchAndWorktree(input)
	git.AddWorktreeForExistingBranch(input)
	if len(calls) != 3 {
		t.Fatalf("calls = %#v", calls)
	}
	writes := calls[1:]
	for _, arguments := range writes {
		joined := strings.Join(arguments, " ")
		if !strings.Contains(joined, "core.hooksPath="+os.DevNull) || strings.Contains(joined, "--force") || strings.Contains(joined, " -B ") {
			t.Fatalf("unsafe Git argv: %q", joined)
		}
	}
}

func TestSystemWorkItemGitInventoryToleratesUnrelatedDetachedWorktree(t *testing.T) {
	repository := initWorkItemGitRepository(t)
	detached := filepath.Join(t.TempDir(), "detached")
	runLocalGit(t, repository, "worktree", "add", "--detach", detached, "HEAD")
	git := newSystemWorkItemGit(systemFileSystem{})
	observation, err := git.ObserveWorktree(repository)
	if err != nil || !observation.InventoryMatch {
		t.Fatalf("observation = %#v, %v", observation, err)
	}
}

func TestTaskWorktreeCreateRejectsInvalidGitBranchBeforeIntent(t *testing.T) {
	fixture := newWorkItemJourneyFixture(t)
	input := TaskWorktreeCreateInput{TaskID: "task", Branch: "-invalid", Path: filepath.Join(fixture.wrapper, "task"), ExpectedParentOID: fixture.oid}
	if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkInvalidArguments) {
		t.Fatalf("error = %v", err)
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(filepath.Dir(fixture.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	if registry.Tasks[0].WorktreeState != WorkItemUnbound || len(registry.WorktreeOperations) != 0 {
		t.Fatalf("invalid branch persisted intent: %#v", registry)
	}
}

func initWorkItemGitRepository(t *testing.T) string {
	t.Helper()
	root := physicalPath(t, t.TempDir())
	runLocalGit(t, root, "init", "-b", "main")
	runLocalGit(t, root, "config", "user.email", "tests@example.invalid")
	runLocalGit(t, root, "config", "user.name", "Ply tests")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("fixture\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, root, "add", "README.md")
	runLocalGit(t, root, "commit", "-m", "fixture")
	return root
}

func runLocalGit(t *testing.T, directory string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, arguments...)...)
	command.Env = append(os.Environ(), "LC_ALL=C", "LANG=C")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
	return strings.TrimSpace(string(output))
}
