package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type sharedTargetFixture struct {
	root, wrapper, main string
	dependencies        Dependencies
	projects            *systemProjectStore
}

// The repository directory itself is the shared clean main checkout, exactly
// like plyapp's main worktree. No Epic worktree is created by the fixture.
func newSharedTargetFixture(t *testing.T) sharedTargetFixture {
	t.Helper()
	root := createProjectWorkspace(t)
	wrapper := filepath.Join(root, "app")
	main := filepath.Join(wrapper, "main")
	if err := os.MkdirAll(main, 0o755); err != nil {
		t.Fatal(err)
	}
	main = physicalPath(t, main)
	runLocalGit(t, main, "init", "-b", "main")
	runLocalGit(t, main, "config", "user.email", "tests@example.invalid")
	runLocalGit(t, main, "config", "user.name", "Ply tests")
	sharedTargetCommit(t, main, "README.md", "fixture\n")
	projects := newSystemProjectStore()
	if _, _, _, err := projects.Add(root, ProjectRecord{ID: "app", Name: "App", Wrapper: physicalPath(t, wrapper), RepoIDs: []RepoID{"app"}}, []RepoRecord{{ID: "app", Locator: main, GitCommonDir: physicalPath(t, filepath.Join(main, ".git"))}}); err != nil {
		t.Fatal(err)
	}
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	d.ProjectLocks, d.Projects = projects, projects
	return sharedTargetFixture{root: root, wrapper: physicalPath(t, wrapper), main: main, dependencies: d, projects: projects}
}

func sharedTargetCommit(t *testing.T, dir, name, text string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, dir, "add", name)
	runLocalGit(t, dir, "commit", "-m", name)
	return runLocalGit(t, dir, "rev-parse", "HEAD")
}

func (f sharedTargetFixture) adopt(t *testing.T, id EpicID, ref, oid string) (EpicMutationResult, error) {
	t.Helper()
	input, err := ParseEpicAdoptInput(string(id), "Epic "+string(id), "app", "app", f.main, ref, oid)
	if err != nil {
		t.Fatal(err)
	}
	return AdoptEpic(f.dependencies, input)
}

func (f sharedTargetFixture) task(t *testing.T, id TaskID, epic EpicID, oid string) TaskWorktreeMutationResult {
	t.Helper()
	input, _ := ParseTaskCreateInput(string(id), "Task "+string(id), "Shared target Task", string(epic), "app", "app")
	if _, err := CreateTask(f.dependencies, input); err != nil {
		t.Fatal(err)
	}
	worktree, _ := ParseTaskWorktreeCreateInput(string(id), string(id), filepath.Join(f.wrapper, string(id)), oid)
	ready, err := CreateTaskWorktree(f.dependencies, worktree)
	if err != nil || ready.Task.Worktree == nil {
		t.Fatalf("Task %s worktree = %#v, %v", id, ready, err)
	}
	return ready
}

func (f sharedTargetFixture) resolveFrom(t *testing.T, cwd string, in QueueTargetInput) (QueueTarget, error) {
	t.Helper()
	d := f.dependencies
	d.Files = projectCwdFileSystem{FileSystem: SystemDependencies().Files, cwd: cwd}
	r, err := d.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	projects, repos, err := f.projects.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	target, _, err := resolveQueueTarget(d, f.root, r, ProjectSnapshot{Projects: projects, Repos: repos}, in)
	return target, err
}

func sharedTargetEpicRecords(t *testing.T, f sharedTargetFixture, epic EpicID) (EpicRecord, []TaskRecord, []WorktreeOperationRecord, []EpicBaseVersion) {
	t.Helper()
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	record, _ := findEpic(r, epic)
	if record == nil {
		t.Fatalf("Epic %s missing", epic)
	}
	tasks, operations, bases := []TaskRecord{}, []WorktreeOperationRecord{}, []EpicBaseVersion{}
	for _, task := range r.Tasks {
		if task.ParentEpicID == epic {
			tasks = append(tasks, task)
		}
	}
	for _, operation := range r.WorktreeOperations {
		if operation.ParentEpicID == epic {
			operations = append(operations, operation)
		}
	}
	for _, base := range r.EpicBaseVersions {
		if base.EpicID == epic {
			bases = append(bases, base)
		}
	}
	return *record, tasks, operations, bases
}

func gitWorktreeCount(t *testing.T, dir string) int {
	t.Helper()
	return strings.Count(runLocalGit(t, dir, "worktree", "list", "--porcelain"), "worktree ")
}

// AC01: plyapp's journey. A historical Epic adopted M0 in the main directory;
// M0 moved to its own worktree and the directory became clean main. A new
// logical Epic adopts the actual clean main without rewriting history.
func TestSharedEpicTargetAdoptsCleanMainAfterHistoricalBindingMoved(t *testing.T) {
	f := newSharedTargetFixture(t)
	runLocalGit(t, f.main, "checkout", "-b", "codex/m0-bootstrap")
	m0 := sharedTargetCommit(t, f.main, "m0.txt", "M0\n")
	if _, err := f.adopt(t, "first-cut", "refs/heads/codex/m0-bootstrap", m0); err != nil {
		t.Fatal(err)
	}
	f.task(t, "old-task", "first-cut", m0)

	runLocalGit(t, f.main, "checkout", "main")
	runLocalGit(t, f.main, "merge", "--ff-only", "codex/m0-bootstrap")
	mainOID := sharedTargetCommit(t, f.main, "main.txt", "main continues\n")
	runLocalGit(t, f.main, "worktree", "add", filepath.Join(f.wrapper, "m0"), "codex/m0-bootstrap")
	oldEpic, oldTasks, oldOperations, oldBases := sharedTargetEpicRecords(t, f, "first-cut")
	worktrees := gitWorktreeCount(t, f.main)

	adopted, err := f.adopt(t, "live-access-log", "refs/heads/main", mainOID)
	if err != nil || !adopted.Created {
		t.Fatalf("new Epic on clean main = %#v, %v", adopted, err)
	}
	epic, tasks, operations, bases := sharedTargetEpicRecords(t, f, "first-cut")
	if !reflect.DeepEqual(epic, oldEpic) || !reflect.DeepEqual(tasks, oldTasks) || !reflect.DeepEqual(operations, oldOperations) || !reflect.DeepEqual(bases, oldBases) {
		t.Fatal("new adoption rewrote historical Epic, Task, operation or base")
	}
	if gitWorktreeCount(t, f.main) != worktrees {
		t.Fatal("adoption created a Git worktree")
	}
	old, err := ShowEpic(f.dependencies, "first-cut")
	if err != nil || old.Ready || !sharedTargetHas(old.Reasons, "epic_ref_stale") {
		t.Fatalf("historical Epic must remain visibly stale: %#v, %v", old, err)
	}
	fresh, err := ShowEpic(f.dependencies, "live-access-log")
	if err != nil || !fresh.Ready || len(fresh.Reasons) != 0 {
		t.Fatalf("new Epic readback = %#v, %v", fresh, err)
	}
	// The stale historical hit is never chosen over the fresh one, from the
	// shared directory or by explicit identity.
	for _, in := range []QueueTargetInput{{}, {ProjectID: "app", RepoID: "app", EpicID: "live-access-log"}} {
		target, err := f.resolveFrom(t, f.main, in)
		if err != nil || target.EpicID != "live-access-log" || target.ParentRef != "refs/heads/main" {
			t.Fatalf("resolve %+v = %+v, %v", in, target, err)
		}
	}
	ready := f.task(t, "new-task", "live-access-log", mainOID)
	if ready.Task.Worktree.ParentOID != mainOID || ready.Task.Worktree.Ref == oldTasks[0].Worktree.Ref || ready.Task.Worktree.Locator == oldTasks[0].Worktree.Locator {
		t.Fatalf("new Task must have its own branch and worktree from fresh main: %#v", ready.Task.Worktree)
	}
}

// AC02: two new Epics share one clean main path/ref; each keeps its own Task
// copy. The ambiguous cwd is refused with usable explicit commands.
func TestSharedEpicTargetTwoEpicsShareMainAndRequireExplicitSelection(t *testing.T) {
	f := newSharedTargetFixture(t)
	oid := runLocalGit(t, f.main, "rev-parse", "HEAD")
	for _, id := range []EpicID{"alpha", "beta"} {
		if result, err := f.adopt(t, id, "refs/heads/main", oid); err != nil || !result.Created {
			t.Fatalf("adopt %s = %#v, %v", id, result, err)
		}
	}
	if retry, err := f.adopt(t, "beta", "refs/heads/main", oid); err != nil || retry.Created {
		t.Fatalf("identical retry = %#v, %v", retry, err)
	}
	if _, err := f.dependencies.WorkItems.Snapshot(f.root); err != nil {
		t.Fatalf("strict registry read rejected shared target: %v", err)
	}
	alpha := f.task(t, "alpha-task", "alpha", oid)
	beta := f.task(t, "beta-task", "beta", oid)
	if alpha.Task.Worktree.Ref == beta.Task.Worktree.Ref || alpha.Task.Worktree.Locator == beta.Task.Worktree.Locator {
		t.Fatal("Tasks of shared Epics must keep separate branches and worktrees")
	}
	_, err := f.resolveFrom(t, f.main, QueueTargetInput{})
	if err == nil {
		t.Fatal("ambiguous shared cwd selected an Epic")
	}
	for _, want := range []string{"task_queue_context_required", "--project app --repo app --epic alpha", "--project app --repo app --epic beta"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("ambiguity error %q lacks %q", err, want)
		}
	}
	for _, id := range []EpicID{"alpha", "beta"} {
		target, err := f.resolveFrom(t, f.main, QueueTargetInput{ProjectID: "app", RepoID: "app", EpicID: id})
		if err != nil || target.EpicID != id || target.ParentLocator != f.main || target.ParentRef != "refs/heads/main" {
			t.Fatalf("explicit %s = %+v, %v", id, target, err)
		}
	}
	// An exclusive Task copy is never adopted as an Epic base.
	copyInput, _ := ParseEpicAdoptInput("gamma", "Epic gamma", "app", "app", alpha.Task.Worktree.Locator, alpha.Task.Worktree.Ref, oid)
	if _, err := AdoptEpic(f.dependencies, copyInput); err == nil || !strings.Contains(err.Error(), "reserved by Task alpha-task") {
		t.Fatalf("Epic adoption of a Task copy = %v", err)
	}
	// A Task copy still may not reuse a shared Epic path or ref.
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	r.WorktreeOperations[0].SourceRef = "refs/heads/main"
	if r.WorktreeOperations[0].IntentDigest, err = intentDigest(r.WorktreeOperations[0]); err != nil {
		t.Fatal(err)
	}
	if err := validateWorkItemRegistry(r); err == nil || !strings.Contains(err.Error(), "Git resource shared") {
		t.Fatalf("Task operation on the Epic ref must be rejected: %v", err)
	}
}

// Exclusivity belongs to the concrete mutation: an unfinished integration
// into the shared ref blocks every Epic using that ref, not only its own.
func TestSharedEpicTargetUnresolvedIntegrationBlocksSiblingEpic(t *testing.T) {
	target := QueueTarget{ProjectID: "app", RepoID: "app", EpicID: "beta", GitCommonDir: "/repo/.git", ParentWorktreeID: "wt_beta", ParentLocator: "/repo", ParentRef: "refs/heads/main"}
	authority := IntegrationAuthority{ID: "ia_alpha"}
	authority.Plan.Epic.EpicID = "alpha"
	authority.Plan.Repository.RepoID = "app"
	authority.Plan.Repository.GitCommonDir = "/repo/.git"
	authority.Plan.Epic.ParentRef = "refs/heads/main"
	r := WorkItemRegistry{IntegrationAuthorities: []IntegrationAuthority{authority}}
	if got := queueUnresolved(r, target); !reflect.DeepEqual(got, []string{"ia_alpha"}) {
		t.Fatalf("sibling integration on shared ref = %v", got)
	}
	r.IntegrationAuthorities[0].Plan.Epic.ParentRef = "refs/heads/other"
	if got := queueUnresolved(r, target); len(got) != 0 {
		t.Fatalf("unrelated ref blocked: %v", got)
	}
}

func sharedTargetHas(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
