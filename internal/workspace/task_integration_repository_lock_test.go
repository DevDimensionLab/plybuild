package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// lockFixture lets independent Task integrations share one workspace.
type lockFixture struct {
	integrationJourneyFixture
}

func newLockFixture(t *testing.T) lockFixture {
	t.Helper()
	return lockFixture{integrationJourneyFixture: newIntegrationJourneyFixture(t)}
}

func (f lockFixture) input() TaskIntegrationInput {
	return TaskIntegrationInput{TaskID: "task", ExpectedResultOID: f.resultOID, ExpectedParentOID: f.oid}
}

// leaveUnresolved persists the fixture Task's integration authority and then
// loses the attempt publication, the production state of an active integration.
// It returns the confirmed input that resumes that durable authority.
func (f lockFixture) leaveUnresolved(t *testing.T) TaskIntegrationInput {
	t.Helper()
	input := f.input()
	checked, err := CheckTaskIntegration(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, checked.Readback)
	publishes := 0
	d := f.dependencies
	d.WorkItems = &systemWorkItemStore{faults: &workItemStoreFaults{fail: func(operation string) error {
		if operation == "temp-open" {
			publishes++
			if publishes == 2 {
				return errInjected
			}
		}
		return nil
	}}}
	if _, err := ApplyTaskIntegration(d, input); err == nil || !hasWorkErrorClass(err, ErrorWorkIO) {
		t.Fatalf("interrupted apply error = %v", err)
	}
	registry, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(registry.IntegrationAuthorities) != 1 || len(registry.IntegrationResults) != 0 {
		t.Fatalf("fixture integration is not unresolved: %#v %v", registry.IntegrationAuthorities, err)
	}
	return input
}

// addRepository registers an independent repository in its own Project with an
// Epic worktree on the given ref.
func (f lockFixture) addRepository(t *testing.T, id string, ref string) (EpicID, string, string) {
	t.Helper()
	wrapper := filepath.Join(f.root, id)
	repository := filepath.Join(wrapper, "main")
	if err := os.MkdirAll(repository, 0o755); err != nil {
		t.Fatal(err)
	}
	repository = physicalPath(t, repository)
	runLocalGit(t, repository, "init", "-b", "main")
	runLocalGit(t, repository, "config", "user.email", "tests@example.invalid")
	runLocalGit(t, repository, "config", "user.name", "Ply tests")
	oid := sharedTargetCommit(t, repository, "README.md", id+"\n")
	epicPath := filepath.Join(wrapper, "epic")
	runLocalGit(t, repository, "worktree", "add", "-b", strings.TrimPrefix(ref, "refs/heads/"), epicPath, oid)
	project := ProjectRecord{ID: ProjectID(id), Name: id, Wrapper: physicalPath(t, wrapper), RepoIDs: []RepoID{RepoID(id)}}
	if _, _, _, err := f.dependencies.Projects.(*systemProjectStore).Add(f.root, project, []RepoRecord{{ID: RepoID(id), Locator: repository, GitCommonDir: physicalPath(t, filepath.Join(repository, ".git"))}}); err != nil {
		t.Fatal(err)
	}
	epic := EpicID(id + "-epic")
	adopt, err := ParseEpicAdoptInput(string(epic), "Epic "+id, id, id, physicalPath(t, epicPath), ref, oid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := AdoptEpic(f.dependencies, adopt); err != nil {
		t.Fatal(err)
	}
	return epic, physicalPath(t, wrapper), oid
}

// addTask registers a Task with a selected Spec, a committed result and a
// passing QA record, returning the input for its own integration.
func (f lockFixture) addTask(t *testing.T, id TaskID, epic EpicID, project ProjectID, repo RepoID, wrapper, parentOID string, n byte) TaskIntegrationInput {
	t.Helper()
	create, err := ParseTaskCreateInput(string(id), "Task "+string(id), "Independent integration", string(epic), string(project), string(repo))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := CreateTask(f.dependencies, create); err != nil {
		t.Fatal(err)
	}
	worktree, _ := ParseTaskWorktreeCreateInput(string(id), string(id), filepath.Join(wrapper, string(id)), parentOID)
	ready, err := CreateTaskWorktree(f.dependencies, worktree)
	if err != nil {
		t.Fatal(err)
	}
	taskPath := ready.Task.Worktree.Locator
	resultOID := sharedTargetCommit(t, taskPath, string(id)+".txt", "result\n")
	return TaskIntegrationInput{TaskID: id, ExpectedResultOID: resultOID, ExpectedParentOID: parentOID}
}

func checkedIntegration(t *testing.T, d Dependencies, input TaskIntegrationInput) (WorkspaceTaskIntegrationReadback, string) {
	t.Helper()
	checked, err := CheckTaskIntegration(d, input)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := MarshalTaskIntegrationReadback(checked.Readback)
	if err != nil {
		t.Fatal(err)
	}
	return checked.Readback, string(encoded)
}

// IL01 ac-01: an unresolved integration of repository A to refs/heads/epic must
// not block repository B from integrating to its own refs/heads/epic.
func TestTaskIntegrationUnresolvedRepositoryDoesNotBlockOtherRepositorySameRef(t *testing.T) {
	f := newLockFixture(t)
	original := f.leaveUnresolved(t)
	epic, wrapper, parentOID := f.addRepository(t, "other", "refs/heads/epic")
	input := f.addTask(t, "other-task", epic, "other", "other", wrapper, parentOID, '2')

	readback, encoded := checkedIntegration(t, f.dependencies, input)
	if readback.Classification != "ready" || strings.Contains(encoded, "parent_integration_already_active") {
		t.Fatalf("independent repository was blocked: %s", encoded)
	}
	input.Apply, input.Confirmation = true, integrationDigest(t, readback)
	applied, err := ApplyTaskIntegration(f.dependencies, input)
	if err != nil || applied.Readback.Classification != "exact_effect" {
		t.Fatalf("independent apply = %#v %v", applied.Readback, err)
	}
	if got := runLocalGit(t, filepath.Join(wrapper, "epic"), "rev-parse", "HEAD"); got != input.ExpectedResultOID {
		t.Fatalf("independent parent HEAD = %s", got)
	}
	if got := runLocalGit(t, filepath.Join(f.wrapper, "epic"), "rev-parse", "HEAD"); got != f.oid {
		t.Fatalf("unresolved repository parent moved: %s", got)
	}

	// The unresolved integration keeps its own authority and completes once.
	git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
	f.dependencies.IntegrationGit = git
	recovered, err := ApplyTaskIntegration(f.dependencies, original)
	if err != nil || recovered.Readback.Classification != "exact_effect" || git.count() != 1 {
		t.Fatalf("recovered = %#v %v merges=%d", recovered.Readback, err, git.count())
	}
	registry, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(registry.IntegrationAuthorities) != 2 || len(registry.IntegrationResults) != 2 {
		t.Fatalf("registry authorities=%d results=%d %v", len(registry.IntegrationAuthorities), len(registry.IntegrationResults), err)
	}
}

// IL01 ac-02: the same repository and target stays serialized, while another
// target in the same repository is independent.
func TestTaskIntegrationSameRepositoryLockFollowsTargetRef(t *testing.T) {
	t.Run("same repository and ref", func(t *testing.T) {
		f := newLockFixture(t)
		f.leaveUnresolved(t)
		input := f.addTask(t, "sibling", "epic", "ply", "ply", f.wrapper, f.oid, '2')
		readback, encoded := checkedIntegration(t, f.dependencies, input)
		if readback.Classification == "ready" || !strings.Contains(encoded, "parent_integration_already_active") {
			t.Fatalf("same repository and ref was not serialized: %s", encoded)
		}
	})
	t.Run("worktree alias of the same repository and ref", func(t *testing.T) {
		f := newLockFixture(t)
		f.leaveUnresolved(t)
		alias := filepath.Join(f.wrapper, "alias")
		runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "move", filepath.Join(f.wrapper, "epic"), alias)
		adopt, err := ParseEpicAdoptInput("alias", "Alias Epic", "ply", "ply", physicalPath(t, alias), "refs/heads/epic", f.oid)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := AdoptEpic(f.dependencies, adopt); err != nil {
			t.Fatal(err)
		}
		input := f.addTask(t, "alias-task", "alias", "ply", "ply", f.wrapper, f.oid, '2')
		readback, encoded := checkedIntegration(t, f.dependencies, input)
		if readback.Classification == "ready" || !strings.Contains(encoded, "parent_integration_already_active") {
			t.Fatalf("worktree alias of the same repository and ref was not serialized: %s", encoded)
		}
	})
	t.Run("same repository other ref", func(t *testing.T) {
		f := newLockFixture(t)
		f.leaveUnresolved(t)
		path := filepath.Join(f.wrapper, "second")
		runLocalGit(t, filepath.Join(f.wrapper, "main"), "worktree", "add", "-b", "second", path, f.oid)
		adopt, err := ParseEpicAdoptInput("second", "Second Epic", "ply", "ply", physicalPath(t, path), "refs/heads/second", f.oid)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := AdoptEpic(f.dependencies, adopt); err != nil {
			t.Fatal(err)
		}
		input := f.addTask(t, "second-task", "second", "ply", "ply", f.wrapper, f.oid, '2')
		readback, encoded := checkedIntegration(t, f.dependencies, input)
		if readback.Classification != "ready" {
			t.Fatalf("other ref in the same repository was blocked: %s", encoded)
		}
	})
}
