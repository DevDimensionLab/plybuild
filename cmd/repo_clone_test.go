package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// cloneWorkspace has a source repository with one commit, reachable as the
// remote of repository src, and members without remote or already present.
func cloneWorkspace(t *testing.T) (root, source string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	if _, err := workspace.Init(workspace.SystemDependencies()); err != nil {
		t.Fatal(err)
	}
	source = filepath.Join(t.TempDir(), "source")
	for _, args := range [][]string{{"init", "-q", "-b", "main", source}, {"-C", source, "-c", "user.name=t", "-c", "user.email=t@e.invalid", "commit", "-q", "--allow-empty", "-m", "first"}} {
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "a", "app"), 0o755); err != nil {
		t.Fatal(err)
	}
	definitions := `{"format_version": 1, "projects": [
  {"id": "a", "name": "A", "wrapper": "a", "repos": [
    {"id": "src", "path": "a/src/main", "remote": "` + source + `"},
    {"id": "local", "path": "a/local/main"},
    {"id": "app", "path": "a/app", "remote": "` + source + `"}], "epics": []},
  {"id": "b", "name": "B", "wrapper": "b", "repos": [], "epics": []}
]}
`
	if err := os.WriteFile(filepath.Join(root, workspacedef.File), []byte(definitions), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, source
}

func TestRepoCloneClonesMissingRepositoriesWithRemote(t *testing.T) {
	root, _ := cloneWorkspace(t)
	stdout, _, err := executeRepo(t, process.SystemRunner().Runner, "clone", "--project", "a")
	if err != nil {
		t.Fatalf("clone: %v", err)
	}
	if want := "src: cloned into a/src/main\nlocal: skipped, no remote\napp: present\n"; stdout != want {
		t.Fatalf("stdout %q, want %q", stdout, want)
	}
	if out, err := exec.Command("git", "-C", filepath.Join(root, "a", "src", "main"), "log", "--format=%s").Output(); err != nil || string(out) != "first\n" {
		t.Fatalf("cloned history %q, error %v", out, err)
	}
	if _, err := os.Lstat(filepath.Join(root, "a", "local")); !os.IsNotExist(err) {
		t.Fatalf("repository without remote created: %v", err)
	}
	if stdout, _, err = executeRepo(t, process.SystemRunner().Runner, "clone", "--project", "a", "src"); err != nil || stdout != "src: present\n" {
		t.Fatalf("rerun %q, error %v", stdout, err)
	}
}

func TestRepoCloneRefusesUnknownRepositoryBeforeCloning(t *testing.T) {
	root, _ := cloneWorkspace(t)
	_, stderr, err := executeRepo(t, process.SystemRunner().Runner, "clone", "--project", "a", "src", "missing")
	if err == nil || !strings.Contains(stderr, "missing is not a repository of Project a") {
		t.Fatalf("stderr %q, error %v", stderr, err)
	}
	if _, err := os.Lstat(filepath.Join(root, "a", "src")); !os.IsNotExist(err) {
		t.Fatalf("cloned before refusing: %v", err)
	}
}

func TestRepoCloneReportsFailureAndContinues(t *testing.T) {
	root, source := cloneWorkspace(t)
	if err := os.RemoveAll(source); err != nil {
		t.Fatal(err)
	}
	stdout, _, err := executeRepo(t, process.SystemRunner().Runner, "clone", "--project", "a")
	if err == nil || !strings.HasPrefix(stdout, "src: failed: ") || !strings.HasSuffix(stdout, "local: skipped, no remote\napp: present\n") {
		t.Fatalf("stdout %q, error %v", stdout, err)
	}
	if _, err := os.Lstat(filepath.Join(root, "a", "src", "main")); !os.IsNotExist(err) {
		t.Fatalf("failed clone left %v", err)
	}
}
