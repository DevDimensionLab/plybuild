package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseWorkItemInputsEnforceTypedIdentifiersAndText(t *testing.T) {
	adopt, err := ParseEpicAdoptInput("epic-one", "Epic One", "ply", "ply", "/tmp/epic", "refs/heads/epic", strings.Repeat("a", 40))
	if err != nil {
		t.Fatal(err)
	}
	if adopt.EpicID != EpicID("epic-one") || adopt.ProjectID != ProjectID("ply") || adopt.RepoID != RepoID("ply") {
		t.Fatalf("adopt input = %#v", adopt)
	}

	for _, value := range []string{"", "Epic", "1-epic", "epic_1", strings.Repeat("a", 64)} {
		if _, err := ParseEpicID(value); err == nil {
			t.Fatalf("ParseEpicID(%q) succeeded", value)
		}
	}
	if _, err := ParseTaskCreateInput("task-one", " Task", "description", "epic-one", "ply", "ply"); err == nil {
		t.Fatal("surrounding title whitespace accepted")
	}
	if _, err := ParseTaskCreateInput("task-one", "Task", "line one\nline two", "epic-one", "ply", "ply"); err == nil {
		t.Fatal("multiline description accepted")
	}
}

func TestWorkspaceWorkItemsHappyPathAndIdenticalRetries(t *testing.T) {
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
	tree := runLocalGit(t, repository, "rev-parse", "HEAD^{tree}")
	common := physicalPath(t, filepath.Join(repository, ".git"))
	epicPath := filepath.Join(wrapper, "epic")
	runLocalGit(t, repository, "worktree", "add", "-b", "epic", epicPath, oid)
	epicPath = physicalPath(t, epicPath)

	project := ProjectRecord{ID: "ply", Name: "Ply", Wrapper: physicalPath(t, wrapper), RepoIDs: []RepoID{"ply"}}
	repo := RepoRecord{ID: "ply", Locator: repository, GitCommonDir: common}
	projects := newSystemProjectStore()
	if _, _, _, err := projects.Add(root, project, []RepoRecord{repo}); err != nil {
		t.Fatal(err)
	}

	dependencies := SystemDependencies()
	dependencies.Files = projectCwdFileSystem{FileSystem: dependencies.Files, cwd: root}
	dependencies.ProjectLocks = projects
	dependencies.Projects = projects
	dependencies.WorkIDs = cryptoWorkItemIDSource{Reader: strings.NewReader(strings.Repeat("a", 16) + strings.Repeat("b", 16) + strings.Repeat("c", 16))}

	adoptInput, _ := ParseEpicAdoptInput("epic", "Epic", "ply", "ply", epicPath, "refs/heads/epic", oid)
	adopted, err := AdoptEpic(dependencies, adoptInput)
	if err != nil || !adopted.Created || adopted.Epic.RepoBindings[0].Worktree.Tree != tree {
		t.Fatalf("adopt = %#v, %v", adopted, err)
	}
	taskInput, _ := ParseTaskCreateInput("task", "Task", "Task description", "epic", "ply", "ply")
	created, err := CreateTask(dependencies, taskInput)
	if err != nil || !created.Created || created.Task.WorktreeState != WorkItemUnbound {
		t.Fatalf("task create = %#v, %v", created, err)
	}
	target := filepath.Join(wrapper, "task")
	worktreeInput, _ := ParseTaskWorktreeCreateInput("task", "task", target, oid)
	ready, err := CreateTaskWorktree(dependencies, worktreeInput)
	if err != nil || ready.Outcome != "created" || ready.Task.WorktreeState != WorkItemReady || ready.Task.Worktree == nil {
		t.Fatalf("worktree create = %#v, %v", ready, err)
	}
	if ready.Task.Worktree.OID != oid || ready.Task.Worktree.Tree != tree || ready.Operation.State != "ready" {
		t.Fatalf("ready binding = %#v, %#v", ready.Task.Worktree, ready.Operation)
	}
	epicReadback, err := ShowEpic(dependencies, "epic")
	if err != nil || !epicReadback.Ready || len(epicReadback.Reasons) != 0 {
		t.Fatalf("Epic readback = %#v, %v", epicReadback, err)
	}
	taskReadback, err := ShowTask(dependencies, "task")
	if err != nil || !taskReadback.WorktreeReady || !taskReadback.ReadyForHandoff || len(taskReadback.Reasons) != 0 {
		t.Fatalf("Task readback = %#v, %v", taskReadback, err)
	}

	registryPath := workItemsPath(root)
	before, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	oldTime := time.Unix(946684800, 0)
	if err := os.Chtimes(registryPath, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	retry, err := CreateTaskWorktree(dependencies, worktreeInput)
	if err != nil || retry.Outcome != "already_ready" || retry.Task.Worktree.ID != ready.Task.Worktree.ID || retry.Operation.ID != ready.Operation.ID {
		t.Fatalf("worktree retry = %#v, %v", retry, err)
	}
	after, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) || info.ModTime().Unix() != oldTime.Unix() {
		t.Fatalf("retry rewrote registry: bytes equal=%t mtime=%v", string(before) == string(after), info.ModTime())
	}

	if err := os.WriteFile(filepath.Join(epicPath, "new.txt"), []byte("new\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, epicPath, "add", "new.txt")
	runLocalGit(t, epicPath, "commit", "-m", "move Epic")
	if _, err := AdoptEpic(dependencies, adoptInput); err == nil || !strings.HasPrefix(err.Error(), string(ErrorWorkGitObservation)+":") {
		t.Fatalf("stale adopted Epic retry error = %v", err)
	}
}

func TestTaskWorktreeCreateFailsClosedForPersistedReconciliation(t *testing.T) {
	t.Run("no effect does not authorize another Git attempt", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		input, before, oldTime := persistReconciliationFixture(t, &fixture)

		result, err := CreateTaskWorktree(fixture.dependencies, input)
		if err == nil || !hasWorkErrorClass(err, ErrorWorkReconciliation) {
			t.Fatalf("result = %#v, error = %v", result, err)
		}
		assertReconciliationRetryUnchanged(t, fixture, before, oldTime)
		if output, err := execGitOutput(filepath.Join(fixture.wrapper, "main"), "show-ref", "--verify", "refs/heads/task"); err == nil {
			t.Fatalf("retry created source ref: %s", output)
		}
		if _, err := os.Lstat(input.Path); !os.IsNotExist(err) {
			t.Fatalf("retry created target path: %v", err)
		}
	})

	t.Run("exact external effect does not authorize store recovery", func(t *testing.T) {
		fixture := newWorkItemJourneyFixture(t)
		input, _, _ := persistReconciliationFixture(t, &fixture)
		registry, err := fixture.dependencies.WorkItems.Snapshot(filepath.Dir(fixture.wrapper))
		if err != nil {
			t.Fatal(err)
		}
		operation := registry.WorktreeOperations[0]
		repository := RepoRecord{ID: operation.RepoID, Locator: filepath.Join(fixture.wrapper, "main"), GitCommonDir: operation.GitCommonDir}
		outcome := fixture.dependencies.WorkGit.CreateBranchAndWorktree(GitCreateInput{Repository: repository, Branch: input.Branch, SourceRef: operation.SourceRef, TargetLocator: operation.TargetLocator, ParentOID: operation.ParentOID})
		if outcome.Err != nil {
			t.Fatalf("external exact effect = %#v", outcome)
		}
		before, err := os.ReadFile(workItemsPath(filepath.Dir(fixture.wrapper)))
		if err != nil {
			t.Fatal(err)
		}
		oldTime := time.Unix(946684800, 0)
		if err := os.Chtimes(workItemsPath(filepath.Dir(fixture.wrapper)), oldTime, oldTime); err != nil {
			t.Fatal(err)
		}

		result, err := CreateTaskWorktree(fixture.dependencies, input)
		if err == nil || !hasWorkErrorClass(err, ErrorWorkReconciliation) {
			t.Fatalf("result = %#v, error = %v", result, err)
		}
		assertReconciliationRetryUnchanged(t, fixture, before, oldTime)
	})
}

func persistReconciliationFixture(t *testing.T, fixture *workItemJourneyFixture) (TaskWorktreeCreateInput, []byte, time.Time) {
	t.Helper()
	base := fixture.dependencies.WorkGit
	fixture.dependencies.WorkGit = noEffectGit{WorkItemGit: base}
	input, _ := ParseTaskWorktreeCreateInput("task", "task", filepath.Join(fixture.wrapper, "task"), fixture.oid)
	if _, err := CreateTaskWorktree(fixture.dependencies, input); err == nil || !hasWorkErrorClass(err, ErrorWorkGitEffect) {
		t.Fatalf("creating error = %v", err)
	}
	fixture.dependencies.WorkGit = base
	root := filepath.Dir(fixture.wrapper)
	if err := fixture.dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
		registry, err := session.Snapshot()
		if err != nil {
			return err
		}
		registry.Tasks[0].WorktreeState = WorkItemReconciliationRequired
		registry.WorktreeOperations[0].State = "reconciliation_required"
		registry.WorktreeOperations[0].LastObservation = emptyCreateObservation("partial_or_unknown")
		registry.WorktreeOperations[0].LastObservation.Reasons = []string{"task_worktree_reconciliation_required"}
		return session.Publish(registry)
	}); err != nil {
		t.Fatal(err)
	}
	path := workItemsPath(root)
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	oldTime := time.Unix(946684800, 0)
	if err := os.Chtimes(path, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	return input, before, oldTime
}

func assertReconciliationRetryUnchanged(t *testing.T, fixture workItemJourneyFixture, before []byte, oldTime time.Time) {
	t.Helper()
	path := workItemsPath(filepath.Dir(fixture.wrapper))
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) || !info.ModTime().Equal(oldTime) {
		t.Fatalf("reconciliation retry rewrote registry: bytes equal=%t mtime=%v", string(after) == string(before), info.ModTime())
	}
	registry, err := fixture.dependencies.WorkItems.Snapshot(filepath.Dir(fixture.wrapper))
	if err != nil {
		t.Fatal(err)
	}
	if registry.Tasks[0].WorktreeState != WorkItemReconciliationRequired || registry.WorktreeOperations[0].State != "reconciliation_required" {
		t.Fatalf("retry changed reconciliation state: %#v", registry)
	}
}

func TestWorkItemIDSourceUsesStableTypedPrefixes(t *testing.T) {
	source := cryptoWorkItemIDSource{Reader: strings.NewReader(strings.Repeat("a", 32))}
	worktreeID, err := source.NewWorktreeID()
	if err != nil {
		t.Fatal(err)
	}
	operationID, err := source.NewOperationID()
	if err != nil {
		t.Fatal(err)
	}
	if string(worktreeID) != "wt_61616161616161616161616161616161" {
		t.Fatalf("worktree ID = %s", worktreeID)
	}
	if string(operationID) != "wop_61616161616161616161616161616161" {
		t.Fatalf("operation ID = %s", operationID)
	}
}
