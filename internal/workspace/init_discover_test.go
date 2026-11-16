package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

func discoverFixture(t *testing.T) (Dependencies, string) {
	t.Helper()
	parent, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(parent, "trip")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: root}
	return d, root
}

func discoverRepo(t *testing.T, path, remote string) {
	t.Helper()
	runGitForTest(t, "", "init", "-q", "-b", "main", path)
	runGitForTest(t, path, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "-q", "--allow-empty", "-m", "fixture")
	if remote != "" {
		runGitForTest(t, path, "remote", "add", "origin", remote)
	}
}

// A directory of repositories becomes a workspace whose ply-workspace.json
// defines one Project named after the directory, with direct-child and
// <name>/main repositories and their origin remotes.
func TestInitDiscoverDefinesProjectFromRepositories(t *testing.T) {
	d, root := discoverFixture(t)
	discoverRepo(t, filepath.Join(root, "api"), "git@example.invalid:api.git")
	discoverRepo(t, filepath.Join(root, "web", "main"), "git@example.invalid:web.git")
	discoverRepo(t, filepath.Join(root, "notes"), "")
	if err := os.Mkdir(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}

	result, err := InitDiscover(d, "")
	if err != nil {
		t.Fatal(err)
	}
	if result.Root != root || result.ProjectID != "trip" || !result.Changed {
		t.Fatalf("result %#v", result)
	}
	defs, _, err := workspacedef.Read(root)
	if err != nil || defs == nil {
		t.Fatalf("definitions %v %v", defs, err)
	}
	want := []workspacedef.Project{{ID: "trip", Name: "trip", Wrapper: ".", Repos: []workspacedef.Repo{
		{ID: "api", Path: "api", Remote: "git@example.invalid:api.git"},
		{ID: "notes", Path: "notes"},
		{ID: "web", Path: "web/main", Remote: "git@example.invalid:web.git"},
	}}}
	if !reflect.DeepEqual(defs.Projects, want) {
		t.Fatalf("projects %#v", defs.Projects)
	}
	if len(result.Omissions) != 1 || result.Omissions[0].Path != filepath.Join(root, "docs") {
		t.Fatalf("omissions %#v", result.Omissions)
	}
	shown, err := ShowProject(d, "trip")
	if err != nil || len(shown.Repos) != 3 {
		t.Fatalf("show %#v %v", shown, err)
	}

	before, err := os.ReadFile(filepath.Join(root, workspacedef.File))
	if err != nil {
		t.Fatal(err)
	}
	again, err := InitDiscover(d, "")
	if err != nil || again.Changed {
		t.Fatalf("retry %#v %v", again, err)
	}
	if after, _ := os.ReadFile(filepath.Join(root, workspacedef.File)); string(after) != string(before) {
		t.Fatal("retry rewrote the definitions")
	}

	discoverRepo(t, filepath.Join(root, "extra"), "")
	if _, err := InitDiscover(d, ""); err == nil || !strings.Contains(err.Error(), "already defined") {
		t.Fatalf("changed repositories: %v", err)
	}
}

// Without any repository nothing is initialized.
func TestInitDiscoverWithoutRepositoriesWritesNothing(t *testing.T) {
	d, root := discoverFixture(t)
	if _, err := InitDiscover(d, ""); err == nil || !strings.Contains(err.Error(), "no Git repositories") {
		t.Fatalf("empty directory: %v", err)
	}
	if entries, _ := os.ReadDir(root); len(entries) != 0 {
		t.Fatalf("wrote %v", entries)
	}
}
