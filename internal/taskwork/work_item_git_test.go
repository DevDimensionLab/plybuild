package taskwork

import (
	"errors"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestSystemWorkItemGitObservesAndCreatesOnlyAdditiveWorktrees(t *testing.T) {
	t.Parallel()
	repository := initWorkItemGitRepository(t)
	git := newSystemWorkItemGit(workspace.SystemDependencies().Files)
	repo := workspace.RepoRecord{ID: "repo", Locator: repository, GitCommonDir: physicalPath(t, filepath.Join(repository, ".git"))}

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
	t.Parallel()
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
		if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, workspace.ErrorWorkGitEffect) {
			t.Fatalf("first error = %v", err)
		}
		registry, err := fixture.dependencies.WorkItems.Snapshot(testWork(t, filepath.Dir(fixture.wrapper)))
		if err != nil || registry.Tasks[0].WorktreeState != WorkItemCreating || registry.WorktreeOperations[0].LastObservation.Classification != "no_effect" {
			t.Fatalf("creating registry = %#v, %v", registry, err)
		}
		worktreeID := registry.WorktreeOperations[0].WorktreeID
		operationID := registry.WorktreeOperations[0].ID
		fixture.dependencies.WorkGit = base
		result, err := CreateTaskWorktree(fixture.dependencies, input)
		if err != nil || result.Outcome != "recovered" || result.Task.Worktree.ID != worktreeID || result.Task.Worktree.OperationID != operationID {
			t.Fatalf("retry = %#v, %v", result, err)
		}
	})

	t.Run("branch-only retry adds worktree without moving branch", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		base := fixture.dependencies.WorkGit
		fixture.dependencies.WorkGit = branchOnlyGit{WorkItemGit: base}
		input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
		if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, workspace.ErrorWorkGitEffect) {
			t.Fatalf("first error = %v", err)
		}
		fixture.dependencies.WorkGit = base
		result, err := CreateTaskWorktree(fixture.dependencies, input)
		if err != nil || result.Outcome != "recovered" || result.Task.Worktree.NativeReceiptVersion != nativeReceiptVersion {
			t.Fatalf("retry = %#v, %v", result, err)
		}
	})

	t.Run("conflicting partial effect requires reconciliation", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		repository := filepath.Join(fixture.wrapper, "main")
		runLocalGit(t, repository, "commit", "--allow-empty", "-m", "different base")
		fixture.dependencies.WorkGit = wrongBranchGit{WorkItemGit: fixture.dependencies.WorkGit}
		input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
		if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, workspace.ErrorWorkReconciliation) {
			t.Fatalf("error = %v", err)
		}
		registry, err := fixture.dependencies.WorkItems.Snapshot(testWork(t, filepath.Dir(fixture.wrapper)))
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
	t.Parallel()
	var calls [][]string
	git := &systemWorkItemGit{files: workspace.SystemDependencies().Files, run: func(arguments []string, environment []string) GitCommandOutcome {
		calls = append(calls, append([]string(nil), arguments...))
		return GitCommandOutcome{}
	}}
	input := GitCreateInput{Repository: workspace.RepoRecord{Locator: "/repo"}, Branch: "task/one", TargetLocator: "/target", ParentOID: strings.Repeat("a", 40)}
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

func TestSystemTaskIntegrationGitUsesTheOnlyAllowedWriteArgvAndSafeEnvironment(t *testing.T) {
	var gotArgv, gotExtra []string
	git := &systemTaskIntegrationGit{files: workspace.SystemDependencies().Files, run: func(arguments []string, environment []string) GitCommandOutcome {
		gotArgv = append([]string(nil), arguments...)
		gotExtra = append([]string(nil), environment...)
		return GitCommandOutcome{}
	}}
	input := IntegrationMergeInput{Repository: workspace.RepoRecord{Locator: "/repo"}, ParentWorktree: "/repo/epic", ResultOID: strings.Repeat("a", 40), AttemptID: "iat_11111111111111111111111111111111"}
	git.MergeFastForward(input)
	want := []string{"-c", "core.hooksPath=" + os.DevNull, "-c", "merge.autoStash=false", "-c", "gc.auto=0", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "submodule.recurse=false", "-C", input.ParentWorktree, "merge", "--ff-only", "--no-stat", "--no-autostash", input.ResultOID}
	if strings.Join(gotArgv, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("write argv = %#v, want %#v", gotArgv, want)
	}
	wantAction := "GIT_REFLOG_ACTION=ply-workspace-task-integrate:" + string(input.AttemptID)
	if len(gotExtra) != 1 || gotExtra[0] != wantAction {
		t.Fatalf("write extra environment = %#v", gotExtra)
	}

	t.Setenv("HOME", "/secret-home")
	t.Setenv("GIT_SSH_COMMAND", "exfiltrate")
	t.Setenv("GIT_CONFIG_GLOBAL", "/unsafe")
	environment := safeIntegrationEnvironment(gotExtra)
	joined := strings.Join(environment, "\n")
	for _, forbidden := range []string{"HOME=/secret-home", "GIT_SSH_COMMAND=", "GIT_CONFIG_GLOBAL=/unsafe"} {
		if strings.Contains(joined, forbidden) {
			t.Fatalf("unsafe environment leaked %q in %#v", forbidden, environment)
		}
	}
	for _, required := range []string{"GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1", "GIT_ALLOW_PROTOCOL=file", "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=" + os.DevNull, wantAction} {
		if !containsExact(environment, required) {
			t.Fatalf("safe environment missing %q: %#v", required, environment)
		}
	}
}

func TestTaskIntegrationGitParsersFailClosedAndBoundStreams(t *testing.T) {
	t.Parallel()
	for version, want := range map[string]bool{"git version 2.44.9": false, "git version 2.45.0": true, "git version 3.0.0": true, "localized": false} {
		if got := supportedGitVersion(version); got != want {
			t.Fatalf("supportedGitVersion(%q)=%t, want %t", version, got, want)
		}
	}
	status, err := parseIntegrationStatus([]byte("? untracked file\x001 M. N... 100644 100644 100644 a b tracked\x00"))
	if err != nil || len(status) != 2 || status[0].RecordKind != "ordinary" || status[1].RecordKind != "untracked" {
		t.Fatalf("status=%#v error=%v", status, err)
	}
	for _, malformed := range [][]byte{[]byte("x unknown\x00"), []byte("1 malformed\x00"), []byte("2 R. N... 100644 100644 100644 a b R100 new\x00")} {
		if _, err := parseIntegrationStatus(malformed); err == nil {
			t.Fatalf("malformed status accepted: %q", malformed)
		}
	}
	stream := commandStream([]byte(strings.Repeat("x", (64<<10)+17)))
	if !stream.Truncated || stream.SizeBytes != (64<<10)+17 || stream.CapturedBytes != 64<<10 || len(stream.Base64) == 0 || !digestPattern.MatchString(stream.SHA256) {
		t.Fatalf("bounded stream = %#v", stream)
	}
}

func TestTaskIntegrationChangedPathsUseNULTerminatedCheckAttrStdin(t *testing.T) {
	t.Parallel()
	var calls [][]string
	var inputs [][]byte
	respond := func(arguments []string) GitCommandOutcome {
		calls = append(calls, append([]string(nil), arguments...))
		joined := strings.Join(arguments, " ")
		switch {
		case strings.Contains(joined, "merge-base --is-ancestor"):
			return GitCommandOutcome{Exit: 0}
		case strings.Contains(joined, "diff-tree"):
			return GitCommandOutcome{Exit: 0, Stdout: []byte("path with spaces\x00")}
		case strings.Contains(joined, "check-attr"):
			return GitCommandOutcome{Exit: 0, Stdout: []byte("path with spaces\x00filter\x00unspecified\x00")}
		default:
			return GitCommandOutcome{Exit: 1, Err: errors.New("unexpected call")}
		}
	}
	git := &systemTaskIntegrationGit{files: workspace.SystemDependencies().Files, run: func(arguments []string, environment []string) GitCommandOutcome {
		return respond(arguments)
	}, runInput: func(arguments []string, environment []string, stdin []byte) GitCommandOutcome {
		inputs = append(inputs, append([]byte(nil), stdin...))
		return respond(arguments)
	}}
	repo := workspace.RepoRecord{Locator: "/repo"}
	if ok, err := git.CheckAncestor(repo, strings.Repeat("a", 40), strings.Repeat("b", 40)); err != nil || !ok {
		t.Fatalf("ancestor=%t error=%v calls=%#v", ok, err, calls)
	}
	if len(calls) != 4 {
		t.Fatalf("calls=%#v", calls)
	}
	for _, call := range calls[2:] {
		joined := strings.Join(call, "\x00")
		if !strings.Contains(joined, "\x00-z\x00--stdin\x00filter") || strings.Contains(joined, "path with spaces") {
			t.Fatalf("check-attr did not use exact --stdin argv: %#v", call)
		}
	}
	if len(inputs) != 2 || string(inputs[0]) != "path with spaces\x00" || string(inputs[1]) != "path with spaces\x00" {
		t.Fatalf("check-attr stdin = %#v", inputs)
	}
}

func containsExact(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestSystemWorkItemGitInventoryToleratesUnrelatedDetachedWorktree(t *testing.T) {
	t.Parallel()
	repository := initWorkItemGitRepository(t)
	detached := filepath.Join(t.TempDir(), "detached")
	runLocalGit(t, repository, "worktree", "add", "--detach", detached, "HEAD")
	git := newSystemWorkItemGit(workspace.SystemDependencies().Files)
	observation, err := git.ObserveWorktree(repository)
	if err != nil || !observation.InventoryMatch {
		t.Fatalf("observation = %#v, %v", observation, err)
	}
}

func TestTaskWorktreeCreateRejectsInvalidGitBranchBeforeIntent(t *testing.T) {
	t.Parallel()
	fixture := newWorkItemJourneyFixture(t)
	input := TaskWorktreeCreateInput{TaskID: "task", Branch: "-invalid", Path: filepath.Join(fixture.wrapper, "task"), ExpectedParentOID: fixture.oid}
	if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, workspace.ErrorWorkInvalidArguments) {
		t.Fatalf("error = %v", err)
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(testWork(t, filepath.Dir(fixture.wrapper)))
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

// ASSIGNMENT-300: a Task whose worktree creation has not finished observes
// its branch twice; one failed Git observation is reported once, as
// source_observation_unknown with Git's cause, not also as a stale start base.
func TestShowTaskInCreationReportsAFailedSourceObservationOnlyAsUnknown(t *testing.T) {
	t.Parallel()
	fixture := newWorkItemJourneyFixture(t)
	base := fixture.dependencies.WorkGit
	fixture.dependencies.WorkGit = noEffectGit{WorkItemGit: base}
	input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
	if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil {
		t.Fatal("fixture: worktree creation unexpectedly succeeded")
	}
	fixture.dependencies.WorkGit = base
	repository := filepath.Join(fixture.wrapper, "main")
	runLocalGit(t, repository, "worktree", "add", "-q", "-b", "task", filepath.Join(fixture.wrapper, "first"), fixture.oid)
	runLocalGit(t, repository, "worktree", "add", "-q", "--force", filepath.Join(fixture.wrapper, "second"), "task")
	// A file at the target path adds a creation reason that must be kept.
	if err := os.WriteFile(filepath.Join(fixture.wrapper, "task"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	shown, err := ShowTask(fixture.dependencies, "task")
	if err != nil {
		t.Fatal(err)
	}
	if shown.Task.WorktreeState != WorkItemCreating {
		t.Fatalf("fixture: worktree state = %s, want creating", shown.Task.WorktreeState)
	}
	if !slices.Contains(shown.Reasons, "source_observation_unknown") || slices.Contains(shown.Reasons, "source_ref_stale") || !slices.Contains(shown.Reasons, "target_worktree_stale") {
		t.Errorf("reasons = %q, want source_observation_unknown and target_worktree_stale and no source_ref_stale", shown.Reasons)
	}
	if shown.Source.Diagnostic == nil || !strings.Contains(*shown.Source.Diagnostic, "checked out more than once") {
		t.Errorf("source diagnostic = %v, want Git's cause", shown.Source.Diagnostic)
	}
}
