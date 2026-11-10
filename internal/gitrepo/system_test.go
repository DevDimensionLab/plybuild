package gitrepo

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

func TestEnsureInitWithSystemGitIgnoresInheritedRepository(t *testing.T) {
	base := tempBase(t)
	emptyConfig := filepath.Join(base, "gitconfig")
	if err := os.WriteFile(emptyConfig, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", base)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", emptyConfig)
	runner := process.SystemRunner().Runner
	other := filepath.Join(base, "other")
	if _, _, err := EnsureInit(context.Background(), runner, other, "trunk"); err != nil {
		t.Fatal(err)
	}
	before := listing(t, filepath.Join(other, ".git"))
	t.Setenv("GIT_DIR", filepath.Join(other, ".git"))
	t.Setenv("GIT_WORK_TREE", other)
	t.Setenv("GIT_INDEX_FILE", filepath.Join(other, ".git", "index"))
	dir := filepath.Join(base, "repo")

	repo, changed, err := EnsureInit(context.Background(), runner, dir, "main")

	if err != nil || !changed || repo.Toplevel != dir || repo.Branch != "main" {
		t.Fatalf("first EnsureInit = %#v, %v, %v", repo, changed, err)
	}
	if headRef, err := os.ReadFile(filepath.Join(dir, ".git", "HEAD")); err != nil || string(headRef) != "ref: refs/heads/main\n" {
		t.Fatalf("HEAD = %q, %v", headRef, err)
	}
	if after := listing(t, filepath.Join(other, ".git")); !reflect.DeepEqual(before, after) {
		t.Fatalf("other repository changed: %v -> %v", before, after)
	}
	if _, changed, err = EnsureInit(context.Background(), runner, dir, "main"); err != nil || changed {
		t.Fatalf("second EnsureInit = %v, %v; want unchanged", changed, err)
	}
}

func listing(t *testing.T, gitDir string) map[string]string {
	entries, err := os.ReadDir(gitDir)
	headRef, headErr := os.ReadFile(filepath.Join(gitDir, "HEAD"))
	if err != nil || headErr != nil {
		t.Fatal(err, headErr)
	}
	files := map[string]string{"HEAD": string(headRef)}
	for _, entry := range entries {
		files[entry.Name()] = entry.Type().String()
	}
	return files
}

func TestEnsureInitWithSystemGitRejectsInvalidBranchName(t *testing.T) {
	dir := filepath.Join(tempBase(t), "repo")

	_, _, err := EnsureInit(context.Background(), process.SystemRunner().Runner, dir, "bad..branch")

	if _, statErr := os.Lstat(dir); !errors.Is(err, ErrInvalid) || !os.IsNotExist(statErr) {
		t.Fatalf("EnsureInit returned %v (dir stat %v), want ErrInvalid without effect", err, statErr)
	}
}

func TestEnsureRemoteWithSystemGitWritesSharedLocalConfigOnly(t *testing.T) {
	base := tempBase(t)
	emptyConfig := filepath.Join(base, "gitconfig")
	if err := os.WriteFile(emptyConfig, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", base)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", emptyConfig)
	runner := process.SystemRunner().Runner
	a, b, w := filepath.Join(base, "a"), filepath.Join(base, "b"), filepath.Join(base, "w")
	for _, dir := range []string{a, b} {
		if _, _, err := EnsureInit(context.Background(), runner, dir, "main"); err != nil {
			t.Fatal(err)
		}
	}
	gitT(t, a, "-c", "user.name=t", "-c", "user.email=t@example.invalid", "commit", "-q", "--allow-empty", "-m", "first")
	gitT(t, a, "remote", "add", "up", "https://example.invalid/up.git")
	gitT(t, a, "worktree", "add", "-q", "--detach", w)
	aConfig, bConfig := filepath.Join(a, ".git", "config"), filepath.Join(b, ".git", "config")
	bBefore := readFile(t, bConfig)
	t.Setenv("GIT_DIR", filepath.Join(b, ".git"))
	url := "https://x:SECRET@example.invalid/origin.git"

	if changed, err := EnsureRemote(context.Background(), runner, w, "origin", url); err != nil || !changed {
		t.Fatalf("EnsureRemote from worktree = %v, %v; want added", changed, err)
	}
	aAfter := readFile(t, aConfig)
	if changed, err := EnsureRemote(context.Background(), runner, a, "origin", url); err != nil || changed {
		t.Fatalf("EnsureRemote from main = %v, %v; want unchanged", changed, err)
	}
	if _, err := EnsureRemote(context.Background(), runner, a, "origin", "https://x:OTHER@example.invalid/o.git"); !errors.Is(err, ErrConflict) ||
		strings.Contains(err.Error(), "example.invalid") {
		t.Fatalf("EnsureRemote with other URL = %v, want ErrConflict without URL", err)
	}
	if got := readFile(t, aConfig); got != aAfter || !strings.Contains(got, `[remote "up"]`) || !strings.Contains(got, "url = "+url+"\n") {
		t.Fatalf("shared config after unchanged/conflict = %q", got)
	}
	if got := readFile(t, bConfig); got != bBefore {
		t.Fatalf("GIT_DIR repository config changed: %q -> %q", bBefore, got)
	}
}

func gitT(t *testing.T, dir string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", dir}, args...)...)
	command.Env = gitEnv(os.Environ())
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", args, err, output)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
