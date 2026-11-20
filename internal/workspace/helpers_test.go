package workspace

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Test helpers of the register tests. internal/taskwork keeps a copy of those
// its tests share with them.

type projectCwdFileSystem struct {
	FileSystem
	cwd string
}

func (filesystem projectCwdFileSystem) Getwd() (string, error) { return filesystem.cwd, nil }

var errInjected = errors.New("injected failure")

func physicalPath(t *testing.T, path string) string {
	t.Helper()
	physical, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(physical)
}

func mustWriteFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func createProjectWorkspace(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, MarkerDirectory), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(root, MarkerDirectory, MarkerFile), fmt.Sprintf("format_version: 1\nroot: %s\n", root))
	mustWriteFile(t, filepath.Join(root, DefinitionsFile), "{\"format_version\": 1, \"projects\": []}\n")
	return root
}

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

func runGitForTest(t *testing.T, directory string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", arguments...)
	if directory != "" {
		command.Dir = directory
	}
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
}

func runLocalGit(t *testing.T, directory string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, arguments...)...)
	command.Env = append(os.Environ(), "LC_ALL=C", "LANG=C")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", arguments, err, output)
	}
	return strings.TrimSpace(string(output))
}

func assertWorkErrorClass(t *testing.T, err error, class WorkItemErrorClass) {
	t.Helper()
	var typed *WorkItemError
	if !errors.As(err, &typed) || typed.Class != class {
		t.Fatalf("error = %v, want class %s", err, class)
	}
}
