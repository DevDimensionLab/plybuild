package taskjournal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// work is the work of the fixture's Project "project". Tests reach journal
// paths only through it or testProcessDir.
func (f testFixture) work(t *testing.T) workspace.ProjectWork {
	t.Helper()
	works, e := workspace.DefinedWork(f.root)
	if e != nil {
		t.Fatal(e)
	}
	for _, work := range works {
		if work.Project == "project" {
			return work
		}
	}
	t.Fatalf("no Project project in %s", f.root)
	return workspace.ProjectWork{}
}

// testProcessDir is the task-process directory of Project id in the
// workspace at root. Tests reach journal paths only through it.
func testProcessDir(t *testing.T, root string, id workspace.ProjectID) string {
	t.Helper()
	works, e := workspace.DefinedWork(root)
	if e != nil {
		t.Fatal(e)
	}
	for _, work := range works {
		if work.Project == id {
			return work.TaskProcess()
		}
	}
	t.Fatalf("no Project %s in %s", id, root)
	return ""
}

// addOtherProject defines Project "other" with its own wrapper, repository,
// Epic "other-epic" and Task "other-task". The fixture Project's wrapper is
// the workspace root, so the other wrapper lies outside it.
func addOtherProject(t *testing.T, f testFixture) {
	t.Helper()
	outside, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	wrapper := filepath.Join(outside, "other")
	repo := filepath.Join(wrapper, "repo")
	if e := os.MkdirAll(repo, 0o700); e != nil {
		t.Fatal(e)
	}
	git := func(args ...string) string {
		t.Helper()
		prefix := []string{"-c", "core.hooksPath=/dev/null", "-c", "commit.gpgSign=false", "-c", "user.name=Synthetic Journal", "-c", "user.email=journal@example.invalid", "-C", repo}
		b, e := exec.Command("git", append(prefix, args...)...).CombinedOutput()
		if e != nil {
			t.Fatalf("fixture git: %v: %s", e, b)
		}
		return strings.TrimSpace(string(b))
	}
	git("init", "-b", "other-epic")
	write(t, filepath.Join(repo, "file"), []byte("other\n"))
	git("add", "file")
	git("commit", "-m", "other fixture")
	d := f.s.Workspace
	if _, e := workspace.AddProject(d.Dependencies, workspace.ProjectAddInput{ProjectID: "other", Name: "Other", Wrapper: wrapper, Repos: []workspace.RepoInput{{RepoID: "other-repo", Path: repo}}}); e != nil {
		t.Fatal(e)
	}
	if _, e := taskwork.AdoptEpic(d, taskwork.EpicAdoptInput{EpicID: "other-epic", Title: "Other Epic", ProjectID: "other", RepoID: "other-repo", Worktree: repo, Ref: "refs/heads/other-epic", ExpectedOID: git("rev-parse", "HEAD")}); e != nil {
		t.Fatal(e)
	}
	if _, e := taskwork.CreateTask(d, taskwork.TaskCreateInput{TaskID: "other-task", Title: "Other Task", Description: "Second Project", ParentEpicID: "other-epic", ProjectID: "other", RepoID: "other-repo"}); e != nil {
		t.Fatal(e)
	}
}

// Each Project keeps its Task journals in its own work; .ply holds none, and
// a journal left in .ply by an earlier version is never read.
func TestJournalsLiveInTheirProjectsWork(t *testing.T) {
	f := fixture(t)
	addOtherProject(t, f)
	old := filepath.Join(f.root, workspace.MarkerDirectory, "task-process", "v1", "task")
	if e := os.MkdirAll(old, 0o700); e != nil {
		t.Fatal(e)
	}
	write(t, filepath.Join(old, "000000001.json"), []byte("not a record"))
	f.append(t, f.input)
	other := f.input
	other.TaskID, other.PublicationKey = "other-task", "other:1"
	if _, e := f.s.Append("other-task", f.file(t, other), "event"); e != nil {
		t.Fatal(e)
	}
	for project, task := range map[workspace.ProjectID]string{"project": "task", "other": "other-task"} {
		entries, e := os.ReadDir(filepath.Join(testProcessDir(t, f.root, project), "v1", task))
		if e != nil || len(entries) == 0 {
			t.Fatalf("journal of %s is not in its Project's work: %v", task, e)
		}
	}
	if _, e := os.Stat(filepath.Join(testProcessDir(t, f.root, "project"), "v1", "other-task")); !os.IsNotExist(e) {
		t.Fatalf("other Project's journal is in project's work: %v", e)
	}
	if _, e := os.Stat(filepath.Join(testProcessDir(t, f.root, "project"), "sources")); e != nil {
		t.Fatalf("source copies are not in the Project's work: %v", e)
	}
	if s := f.show(t); len(s.Events) == 0 || s.Coverage.State == "" {
		t.Fatalf("show = %+v", s)
	}
	journals, e := Journals(f.root)
	if e != nil || len(journals) != 2 {
		t.Fatalf("journals = %+v %v", journals, e)
	}
	w, e := workspace.ObserveContaining(f.s.Workspace.Dependencies)
	if e != nil {
		t.Fatal(e)
	}
	for _, j := range journals {
		if j.Task == "task" && j.Work.Project != "project" || j.Task == "other-task" && j.Work.Project != "other" {
			t.Fatalf("journal %s found in %s", j.Task, j.Work.Project)
		}
		if records, e := ReadContributions(w, j.Work, j.Task); e != nil || len(records) != 1 {
			t.Fatalf("contributions of %s = %d, %v", j.Task, len(records), e)
		}
	}
	if b, _ := os.ReadFile(filepath.Join(old, "000000001.json")); string(b) != "not a record" {
		t.Fatal("old .ply journal was changed")
	}
}
