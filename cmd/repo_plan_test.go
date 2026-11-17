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

// planWorkspace has the repository struct-diagram in wrapper a/struct_diagram,
// whose directory name differs from the repository ID, two repositories
// sharing wrapper a/shared, and a planning template in a temporary home.
func planWorkspace(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	if _, err := workspace.Init(workspace.SystemDependencies()); err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	template := filepath.Join(home, ".agents", "templates", "plan-repo")
	if err := os.MkdirAll(template, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(template, "AGENTS.md"), []byte("# {{repo}}-plan\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{"a/struct_diagram/main", "a/shared"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	definitions := `{"format_version": 1, "projects": [
  {"id": "a", "name": "A", "wrapper": "a", "repos": [
    {"id": "struct-diagram", "path": "a/struct_diagram/main", "wrapper": "a/struct_diagram"},
    {"id": "one", "path": "a/shared/one", "wrapper": "a/shared"},
    {"id": "two", "path": "a/shared/two", "wrapper": "a/shared"}]}
]}
`
	if err := os.WriteFile(filepath.Join(root, workspacedef.File), []byte(definitions), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestRepoPlanCreatesPlanningRepositoryNamedByRepositoryID(t *testing.T) {
	root := planWorkspace(t)
	wrapper := filepath.Join(root, "a", "struct_diagram")
	t.Chdir(wrapper)
	stdout, _, err := executeRepo(t, process.SystemRunner().Runner, "plan")
	plan := filepath.Join(wrapper, "struct-diagram-plan")
	if err != nil || stdout != "Planning repository "+plan+": created\n" {
		t.Fatalf("stdout %q, error %v", stdout, err)
	}
	if data, err := os.ReadFile(filepath.Join(plan, "AGENTS.md")); err != nil || string(data) != "# struct-diagram-plan\n" {
		t.Fatalf("AGENTS.md %q, error %v", data, err)
	}
	if out, err := exec.Command("git", "-C", plan, "rev-list", "--count", "HEAD").Output(); err != nil || string(out) != "1\n" {
		t.Fatalf("commits %q, error %v", out, err)
	}
	stdout, _, err = executeRepo(t, process.SystemRunner().Runner, "plan")
	if err != nil || stdout != "Planning repository "+plan+": unchanged\n" {
		t.Fatalf("rerun %q, error %v", stdout, err)
	}
}

func TestRepoPlanRefusesWrapperWithoutExactlyOneRepository(t *testing.T) {
	root := planWorkspace(t)
	for dir, want := range map[string]string{
		filepath.Join(root, "a", "shared"): "2 repositories (one, two) have wrapper",
		filepath.Join(root, "a"):           "0 repositories () have wrapper",
	} {
		t.Chdir(dir)
		_, stderr, err := executeRepo(t, process.SystemRunner().Runner, "plan")
		if err == nil || !strings.Contains(stderr, want) {
			t.Fatalf("%s: stderr %q, error %v", dir, stderr, err)
		}
		entries, _ := os.ReadDir(dir)
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), "-plan") {
				t.Fatalf("%s: created %s before refusing", dir, entry.Name())
			}
		}
	}
}
