package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func epicCreateFixture(t *testing.T) (workItemJourneyFixture, EpicCreateInput) {
	t.Helper()
	f := newWorkItemJourneyFixture(t)
	queueSetFixture(t, f, "create-fixture", []QueueEntry{})
	root := filepath.Dir(f.wrapper)
	r, err := f.dependencies.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	// An isolated format-6 fixture, never a migration of a real registry.
	r.FormatVersion = 6
	boundary := func(from int) *RegistryFormatBoundary {
		return &RegistryFormatBoundary{FromVersion: from, FromSHA256: "sha256:" + strings.Repeat("a", 64), ManifestSHA256: "sha256:" + strings.Repeat("b", 64), ReaderVersion: "v1.0.0", ReaderCommit: strings.Repeat("c", 40), ReaderSHA256: "sha256:" + strings.Repeat("d", 64)}
	}
	r.FormatBoundary = boundary(4)
	r.EpicCreateBoundary = boundary(5)
	r.EpicCreateOperations = []EpicCreateOperation{}
	b, err := encodeWorkItemRegistry(r)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(workItemsPath(root), b, 0600); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(f.wrapper, "main")
	in := EpicCreateInput{EpicID: "new-epic", Title: "New Epic", ProjectID: "ply", RepoID: "ply", ParentRef: "refs/heads/main", ExpectedParentOID: runLocalGit(t, repo, "rev-parse", "HEAD"), Branch: "new-epic", Path: filepath.Join(f.wrapper, "new-epic")}
	return f, in
}

func TestEpicCreateHistoricalRetryAndTaskAnchor(t *testing.T) {
	f, in := epicCreateFixture(t)
	got, err := CreateEpic(f.dependencies, in)
	if err != nil {
		t.Fatal(err)
	}
	if got.Outcome != "created" || got.Epic == nil || got.Epic.RepoBindings[0].Worktree.Origin != "ply_created" {
		t.Fatalf("%+v", got)
	}
	if _, err = ShowEpic(f.dependencies, in.EpicID); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(f.wrapper, "main")
	runLocalGit(t, repo, "commit", "--allow-empty", "-m", "parent moved")
	if err = os.WriteFile(filepath.Join(in.Path, "dirty"), []byte("dirty"), 0600); err != nil {
		t.Fatal(err)
	}
	again, err := CreateEpic(f.dependencies, in)
	if err != nil || again.Outcome != "already_created" || again.Operation.ID != got.Operation.ID {
		t.Fatalf("%+v %v", again, err)
	}
	in.Title = "Changed"
	if _, err = CreateEpic(f.dependencies, in); err == nil {
		t.Fatal("changed retry accepted")
	}
}

type epicCreateFaultGit struct {
	WorkItemGit
	mode  string
	calls int
}

func (g *epicCreateFaultGit) CreateBranchAndWorktree(in GitCreateInput) GitCommandOutcome {
	g.calls++
	if g.mode == "branch" {
		out := runSafeIntegrationGit([]string{"-C", in.Repository.Locator, "branch", in.Branch, in.ParentOID}, nil)
		if out.Err != nil {
			return out
		}
		return GitCommandOutcome{Exit: 1, Err: errors.New("lost response")}
	}
	if g.mode == "unknown" {
		_ = os.Mkdir(in.TargetLocator, 0700)
		return GitCommandOutcome{Exit: 1, Err: errors.New("partial directory")}
	}
	if g.mode == "none" {
		return GitCommandOutcome{Exit: 1, Err: errors.New("known no effect")}
	}
	out := g.WorkItemGit.CreateBranchAndWorktree(in)
	if out.Err == nil {
		out.Exit = 1
		out.Err = errors.New("lost response")
	}
	return out
}
func TestEpicCreateRecoveryAndReservation(t *testing.T) {
	for _, mode := range []string{"branch", "exact", "none", "unknown"} {
		t.Run(mode, func(t *testing.T) {
			f, in := epicCreateFixture(t)
			g := &epicCreateFaultGit{WorkItemGit: f.dependencies.WorkGit, mode: mode}
			f.dependencies.WorkGit = g
			got, err := CreateEpic(f.dependencies, in)
			if mode == "exact" {
				if err != nil || got.Outcome != "created" {
					t.Fatalf("%+v %v", got, err)
				}
				return
			}
			if err == nil {
				t.Fatal("fault should remain pending")
			}
			r, e := f.dependencies.WorkItems.Snapshot(filepath.Dir(f.wrapper))
			if e != nil {
				t.Fatal(e)
			}
			if len(r.EpicCreateOperations) != 1 {
				t.Fatal("intent missing")
			}
			if epic, _ := findEpic(r, in.EpicID); epic != nil {
				t.Fatal("pending Epic fabricated")
			}
			if _, e = AdoptEpic(f.dependencies, EpicAdoptInput{EpicID: in.EpicID, Title: in.Title, ProjectID: in.ProjectID, RepoID: in.RepoID, Worktree: filepath.Join(f.wrapper, "main"), Ref: in.ParentRef, ExpectedOID: in.ExpectedParentOID}); e == nil {
				t.Fatal("pending identity not reserved")
			}
			if mode == "unknown" {
				_, e = CreateEpic(f.dependencies, in)
				if e == nil || g.calls != 1 {
					t.Fatal("unknown retried Git")
				}
				return
			}
			g.mode = "exact"
			got, e = CreateEpic(f.dependencies, in)
			if e != nil || got.Outcome != "created" {
				t.Fatalf("retry: %+v %v", got, e)
			}
			if mode == "branch" && g.calls != 1 {
				t.Fatal("branch recreated")
			}
		})
	}
}
func TestEpicCreateRejectsUnsafeParentTargetAndFilter(t *testing.T) {
	for _, mode := range []string{"dirty", "nested", "symlink", "owned", "filter"} {
		t.Run(mode, func(t *testing.T) {
			f, in := epicCreateFixture(t)
			repo := filepath.Join(f.wrapper, "main")
			switch mode {
			case "dirty":
				_ = os.WriteFile(filepath.Join(repo, "dirty"), []byte("x"), 0600)
			case "nested":
				in.Path = filepath.Join(repo, "nested")
			case "symlink":
				if e := os.Symlink("missing", in.Path); e != nil {
					t.Fatal(e)
				}
			case "owned":
				in.ParentRef = "refs/heads/epic"
			case "filter":
				_ = os.WriteFile(filepath.Join(repo, ".gitattributes"), []byte("README.md filter=evil\n"), 0600)
				runLocalGit(t, repo, "add", ".gitattributes")
				runLocalGit(t, repo, "commit", "-m", "attributes")
				in.ExpectedParentOID = runLocalGit(t, repo, "rev-parse", "HEAD")
				runLocalGit(t, repo, "config", "filter.evil.smudge", "touch "+filepath.Join(f.wrapper, "forbidden"))
			}
			if _, e := CreateEpic(f.dependencies, in); e == nil {
				t.Fatal("unsafe request accepted")
			}
			if _, e := os.Stat(filepath.Join(f.wrapper, "forbidden")); !os.IsNotExist(e) {
				t.Fatal("filter executed")
			}
			if out := runSafeIntegrationGit([]string{"-C", repo, "show-ref", "--verify", "refs/heads/new-epic"}, nil); out.Exit == 0 {
				t.Fatal("ref created")
			}
		})
	}
}

type epicTerminalFaultStore struct {
	WorkItemStore
	fail bool
}
type epicTerminalFaultSession struct {
	WorkItemStoreSession
	owner *epicTerminalFaultStore
}

func (s *epicTerminalFaultSession) Publish(r WorkItemRegistry) error {
	if s.owner.fail && len(r.EpicCreateOperations) > 0 && r.EpicCreateOperations[0].State == "completed" {
		s.owner.fail = false
		return errors.New("terminal publication interrupted")
	}
	return s.WorkItemStoreSession.Publish(r)
}
func (s *epicTerminalFaultStore) WithLock(root string, f func(WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(inner WorkItemStoreSession) error { return f(&epicTerminalFaultSession{inner, s}) })
}
func TestEpicCreateLostTerminalPublication(t *testing.T) {
	f, in := epicCreateFixture(t)
	g := &epicCreateFaultGit{WorkItemGit: f.dependencies.WorkGit, mode: "exact"}
	f.dependencies.WorkGit = g
	f.dependencies.WorkItems = &epicTerminalFaultStore{WorkItemStore: f.dependencies.WorkItems, fail: true}
	if _, e := CreateEpic(f.dependencies, in); e == nil {
		t.Fatal("lost publication should report error")
	}
	read, e := ShowEpic(f.dependencies, in.EpicID)
	if e != nil || read.Creation == nil || read.Epic.ID != "" {
		t.Fatalf("pending read: %+v %v", read, e)
	}
	if _, e = MarshalEpicReadback(read); e != nil {
		t.Fatal(e)
	}
	if _, e = CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: "task", Branch: in.Branch, Path: in.Path, ExpectedParentOID: in.ExpectedParentOID}); e == nil {
		t.Fatal("Task consumed reserved target")
	}
	got, e := CreateEpic(f.dependencies, in)
	if e != nil || got.Outcome != "created" || g.calls != 1 {
		t.Fatalf("exact recovery replayed Git: %+v %v calls=%d", got, e, g.calls)
	}
}
func TestEpicCreateEffectBoundaryAndNewTask(t *testing.T) {
	f, in := epicCreateFixture(t)
	repo := filepath.Join(f.wrapper, "main")
	marker := filepath.Join(f.wrapper, "forbidden")
	hook := filepath.Join(repo, ".git", "hooks", "post-checkout")
	if e := os.WriteFile(hook, []byte("#!/bin/sh\ntouch '"+marker+"'\n"), 0700); e != nil {
		t.Fatal(e)
	}
	runLocalGit(t, repo, "config", "core.fsmonitor", "touch '"+marker+"'")
	t.Setenv("GIT_DIR", filepath.Join(f.wrapper, "wrong-git"))
	t.Setenv("GIT_WORK_TREE", f.wrapper)
	got, e := CreateEpic(f.dependencies, in)
	if e != nil {
		t.Fatal(e)
	}
	if got.Outcome != "created" {
		t.Fatal(got)
	}
	if _, e = os.Stat(marker); !os.IsNotExist(e) {
		t.Fatal("hook/fsmonitor executed")
	}
	// Task worktree creation remains usable with a newly created anchor.
	t.Setenv("GIT_DIR", "")
	t.Setenv("GIT_WORK_TREE", "")
	os.Unsetenv("GIT_DIR")
	os.Unsetenv("GIT_WORK_TREE")
	runLocalGit(t, repo, "config", "--unset", "core.fsmonitor")
	task, e := ParseTaskCreateInput("new-task", "New task", "A fixture task", string(in.EpicID), "ply", "ply")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = CreateTask(f.dependencies, task); e != nil {
		t.Fatal(e)
	}
	if _, e = CreateTaskWorktree(f.dependencies, TaskWorktreeCreateInput{TaskID: task.TaskID, Branch: "new-task", Path: filepath.Join(f.wrapper, "new-task"), ExpectedParentOID: in.ExpectedParentOID}); e != nil {
		t.Fatal(e)
	}
}
func TestEpicCreateStrictFormatAndPermit(t *testing.T) {
	f, in := epicCreateFixture(t)
	root := filepath.Dir(f.wrapper)
	r, e := f.dependencies.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	old := r
	old.FormatVersion = 5
	old.EpicCreateBoundary = nil
	old.EpicCreateOperations = nil
	b, e := encodeWorkItemRegistry(old)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = decodeWorkItemRegistry(append(b, []byte("epic_create_operations: []\n")...)); e == nil {
		t.Fatal("format 5 silently extended")
	}
	if e = validateRegistryFormatTransition(old, r); RegistryError(e) == nil {
		t.Fatalf("missing permit accepted: %v", e)
	}
	if e = os.WriteFile(workItemsPath(root), b, 0600); e != nil {
		t.Fatal(e)
	}
	before := runLocalGit(t, filepath.Join(f.wrapper, "main"), "for-each-ref", "--format=%(refname) %(objectname)")
	if _, e = CreateEpic(f.dependencies, in); RegistryError(e) == nil {
		t.Fatalf("legacy create: %v", e)
	}
	if _, e = ShowEpic(f.dependencies, "epic"); e != nil {
		t.Fatal(e)
	}
	if _, e = CreateTask(f.dependencies, TaskCreateInput{TaskID: "ordinary", Title: "Ordinary", Description: "Existing flow", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); e != nil {
		t.Fatal(e)
	}
	if after := runLocalGit(t, filepath.Join(f.wrapper, "main"), "for-each-ref", "--format=%(refname) %(objectname)"); after != before {
		t.Fatal("legacy rejection changed Git")
	}
}

func TestEpicCreateUnusedFilterDoesNotRunDuringStatus(t *testing.T) {
	f, in := epicCreateFixture(t)
	repo := filepath.Join(f.wrapper, "main")
	marker := filepath.Join(f.wrapper, "filter-executed")
	for _, key := range []string{"filter.unused.clean", "filter.unused.smudge", "filter.unused.process"} {
		runLocalGit(t, repo, "config", key, "touch '"+marker+"'")
	}
	runLocalGit(t, repo, "config", "filter.unused.required", "true")
	// An index-clean parent with no filter attribute in the chosen tree must work.
	got, e := CreateEpic(f.dependencies, in)
	if e != nil || got.Outcome != "created" {
		t.Fatalf("unused driver blocked ordinary create: %+v %v", got, e)
	}
	if _, e = os.Stat(marker); !os.IsNotExist(e) {
		t.Fatal("unused driver ran")
	}
	// Working attributes can differ from the selected commit. Status must inspect
	// and reject the dirty parent without running its newly selected clean driver.
	in.EpicID = "next-epic"
	in.Branch = "next-epic"
	in.Path = filepath.Join(f.wrapper, "next-epic")
	if e = os.WriteFile(filepath.Join(repo, ".gitattributes"), []byte("README.md filter=unused\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(repo, "README.md"), []byte("changed\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = CreateEpic(f.dependencies, in); e == nil {
		t.Fatal("dirty parent accepted")
	}
	if _, e = os.Stat(marker); !os.IsNotExist(e) {
		t.Fatal("working attributes ran a clean driver")
	}
}
