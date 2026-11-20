package taskwork

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// twoProjectWorkspace defines Projects "alpha" and "beta", each with one
// repository and an adopted Epic of the same name. Each Project's Epic has a
// Task "<project>-task" with a recorded problem.
func twoProjectWorkspace(t *testing.T) (Dependencies, string) {
	t.Helper()
	root := createProjectWorkspace(t)
	d := SystemDependencies()
	d.ProjectDataRoot = func(string) (string, error) { return "", nil }
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	for _, name := range []string{"alpha", "beta"} {
		wrapper := filepath.Join(root, name)
		repository := filepath.Join(wrapper, "main")
		if err := os.MkdirAll(repository, 0o755); err != nil {
			t.Fatal(err)
		}
		runLocalGit(t, repository, "init", "-b", "main")
		runLocalGit(t, repository, "config", "user.email", "tests@example.invalid")
		runLocalGit(t, repository, "config", "user.name", "Ply tests")
		mustWriteFile(t, filepath.Join(repository, "README.md"), name+"\n")
		runLocalGit(t, repository, "add", "README.md")
		runLocalGit(t, repository, "commit", "-m", "fixture")
		oid := runLocalGit(t, repository, "rev-parse", "HEAD")
		epicPath := filepath.Join(wrapper, "epic")
		runLocalGit(t, repository, "worktree", "add", "-b", name, epicPath, oid)
		project := workspace.ProjectRecord{ID: workspace.ProjectID(name), Name: name, Wrapper: wrapper, RepoIDs: []workspace.RepoID{workspace.RepoID(name)}}
		defineTestProject(t, root, project, []workspace.RepoRecord{{ID: workspace.RepoID(name), Locator: repository, GitCommonDir: filepath.Join(repository, ".git")}})
		adopt, _ := ParseEpicAdoptInput(name, name, name, name, epicPath, "refs/heads/"+name, oid)
		if _, err := AdoptEpic(d, adopt); err != nil {
			t.Fatal(err)
		}
		task, _ := ParseTaskCreateInput(name+"-task", "Task", "Description", name, name, name)
		if err := fixtureCreateTaskWithProblem(d, root, task); err != nil {
			t.Fatal(err)
		}
	}
	return d, root
}

// Each Project keeps its own registry under its wrapper; .ply keeps only the
// marker and locks.
func TestEachProjectHasItsOwnRegistryUnderItsWrapper(t *testing.T) {
	t.Parallel()
	d, root := twoProjectWorkspace(t)
	for _, name := range []workspace.ProjectID{"alpha", "beta"} {
		work := testWork(t, root, name)
		if work.Dir != filepath.Join(root, string(name), "data", "work") {
			t.Fatalf("work of %s is %s", name, work.Dir)
		}
		registry, err := d.WorkItems.Snapshot(work)
		if err != nil {
			t.Fatal(err)
		}
		if len(registry.Epics) != 1 || registry.Epics[0].ProjectID != name || len(registry.Tasks) != 1 || registry.Tasks[0].ProjectID != name {
			t.Fatalf("registry of %s holds %+v %+v", name, registry.Epics, registry.Tasks)
		}
		if _, err := os.Stat(filepath.Join(work.Dir, "task-content", "v1")); err != nil {
			t.Fatalf("task content of %s is not in its work: %v", name, err)
		}
	}
	entries, err := os.ReadDir(filepath.Join(root, workspace.MarkerDirectory))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
		if entry.IsDir() || entry.Name() != workspace.MarkerFile && !strings.HasSuffix(entry.Name(), ".lock") {
			t.Fatalf(".ply holds more than marker and locks: %v", names)
		}
	}
	sort.Strings(names)
	if !strings.Contains(strings.Join(names, " "), workIDsLock) {
		t.Fatalf("creation took no workspace ID lock: %v", names)
	}
}

// Commands with only an ID find the item in whichever Project owns it.
func TestTaskShowAndListFindItemsInEitherProject(t *testing.T) {
	t.Parallel()
	d, _ := twoProjectWorkspace(t)
	for _, name := range []string{"alpha", "beta"} {
		shown, err := ShowTask(d, TaskID(name+"-task"))
		if err != nil || shown.Task.ProjectID != workspace.ProjectID(name) || shown.Epic.ID != EpicID(name) {
			t.Fatalf("show %s: %+v %v", name, shown.Task, err)
		}
		epic, err := ShowEpic(d, EpicID(name))
		if err != nil || epic.Epic.ProjectID != workspace.ProjectID(name) {
			t.Fatalf("show Epic %s: %+v %v", name, epic.Epic, err)
		}
	}
	tasks, err := ListTasks(d, nil)
	if err != nil || len(tasks.Tasks) != 2 || tasks.Tasks[0].ID != "alpha-task" || tasks.Tasks[1].ID != "beta-task" {
		t.Fatalf("list: %+v %v", tasks.Tasks, err)
	}
	epics, err := ListEpics(d)
	if err != nil || len(epics.Epics) != 2 {
		t.Fatalf("Epic list: %+v %v", epics.Epics, err)
	}
	if _, err := ShowTask(d, "missing"); err == nil {
		t.Fatal("unknown Task was found")
	} else {
		assertWorkErrorClass(t, err, workspace.ErrorWorkNotFound)
	}
}

// Epic and Task IDs are unique across the workspace.
func TestCreatingAnIDOfAnotherProjectIsRejected(t *testing.T) {
	t.Parallel()
	d, root := twoProjectWorkspace(t)
	before := testWorkRegistryBytes(t, root, "beta")
	task, _ := ParseTaskCreateInput("alpha-task", "Task", "Description", "beta", "beta", "beta")
	_, err := CreateTask(d, task)
	assertWorkErrorClass(t, err, workspace.ErrorWorkIdentityConflict)
	if !strings.Contains(err.Error(), "Project alpha") {
		t.Fatalf("error does not name the owner: %v", err)
	}
	epicPath := filepath.Join(root, "beta", "epic")
	oid := runLocalGit(t, epicPath, "rev-parse", "HEAD")
	adopt, _ := ParseEpicAdoptInput("alpha", "Alpha again", "beta", "beta", epicPath, "refs/heads/beta", oid)
	_, err = AdoptEpic(d, adopt)
	assertWorkErrorClass(t, err, workspace.ErrorWorkIdentityConflict)
	if after := testWorkRegistryBytes(t, root, "beta"); after != before {
		t.Fatal("rejected creation changed the registry")
	}
}

// A registry left in .ply by an earlier version is never read.
func TestOldWorkspaceRegistryIsIgnored(t *testing.T) {
	t.Parallel()
	d, root := twoProjectWorkspace(t)
	mustWriteFile(t, filepath.Join(root, workspace.MarkerDirectory, workspace.WorkItemsFile), "not: [a registry")
	if _, err := ShowTask(d, "alpha-task"); err != nil {
		t.Fatal(err)
	}
	if _, err := ListTasks(d, nil); err != nil {
		t.Fatal(err)
	}
}

func testWorkRegistryBytes(t *testing.T, root string, project workspace.ProjectID) string {
	t.Helper()
	b, err := os.ReadFile(testWork(t, root, project).Registry())
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// Each Project keeps its own lifecycle history in its work; .ply holds none
// and an old .ply lifecycle file is never read.
func TestEachProjectHasItsOwnLifecycleUnderItsWrapper(t *testing.T) {
	t.Parallel()
	d, root := twoProjectWorkspace(t)
	mustWriteFile(t, filepath.Join(root, workspace.MarkerDirectory, workspace.WorkItemLifecycleFile), "not json")
	for _, name := range []string{"alpha", "beta"} {
		result, err := SetWorkItemLifecycle(d, WorkItemLifecycleInput{SubjectKind: "task", SubjectID: name + "-task", State: LifecycleParked, ActorClaim: "actor"})
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join(testWork(t, root, workspace.ProjectID(name)).Dir, workspace.WorkItemLifecycleFile); result.Source.Locator != want || result.Revision != 1 {
			t.Fatalf("lifecycle of %s at %s revision %d", name, result.Source.Locator, result.Revision)
		}
	}
	if _, err := SetWorkItemLifecycle(d, WorkItemLifecycleInput{SubjectKind: "epic", SubjectID: "beta", State: LifecycleFrozen, ActorClaim: "actor"}); err != nil {
		t.Fatal(err)
	}
	all, err := ReadWorkItemLifecycle(d, root)
	if err != nil {
		t.Fatal(err)
	}
	if all.Task("alpha-task") != LifecycleParked || all.Task("beta-task") != LifecycleParked || all.Epic("beta") != LifecycleFrozen || all.Epic("alpha") != LifecycleActive || len(all.Events()) != 3 {
		t.Fatalf("workspace lifecycle = %+v", all)
	}
	if alpha := all.Of("alpha"); alpha.Revision != 1 || alpha.Epic("beta") != LifecycleActive {
		t.Fatalf("alpha lifecycle holds beta's history: %+v", alpha)
	}
	shown, err := ShowWorkItemLifecycle(d, "epic", "beta")
	if err != nil || shown.Lifecycle != LifecycleFrozen || shown.Revision != 2 {
		t.Fatalf("show = %+v %v", shown, err)
	}
	if b, _ := os.ReadFile(filepath.Join(root, workspace.MarkerDirectory, workspace.WorkItemLifecycleFile)); string(b) != "not json" {
		t.Fatalf(".ply lifecycle was rewritten: %q", b)
	}
}
