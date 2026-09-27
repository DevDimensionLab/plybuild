package workflowhandoff

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestSystemGitObserverTargetAndInputBinding(t *testing.T) {
	repo := initObserverRepository(t)
	ref := gitOutput(t, repo, "symbolic-ref", "HEAD")
	oid := gitOutput(t, repo, "rev-parse", "HEAD")
	tree := gitOutput(t, repo, "rev-parse", "HEAD^{tree}")
	common := gitOutput(t, repo, "rev-parse", "--path-format=absolute", "--git-common-dir")
	blob := gitOutput(t, repo, "rev-parse", "HEAD:input.txt")

	observer := SystemDependencies().Git
	clean, err := observer.ObserveTarget(TargetRequest{Worktree: repo, Ref: ref})
	if err != nil {
		t.Fatal(err)
	}
	if clean.Worktree != repo || clean.GitCommonDir != common || clean.Ref != ref || clean.OID != oid || clean.Tree != tree || clean.Status.Mode != "clean" {
		t.Fatalf("clean observation = %#v", clean)
	}

	if err := os.WriteFile(filepath.Join(repo, "new file.txt"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dirty, err := observer.ObserveTarget(TargetRequest{Worktree: repo, Ref: ref})
	if err != nil {
		t.Fatal(err)
	}
	if dirty.Status.Mode != "exact_manifest" || len(dirty.Status.Entries) != 1 || dirty.Status.Entries[0].Path != "new file.txt" || dirty.Status.Entries[0].WorktreeState != "untracked" {
		t.Fatalf("dirty status = %#v", dirty.Status)
	}

	input, err := observer.ObserveInputGit(InputGitRequest{Locator: filepath.Join(repo, "input.txt"), Binding: GitBinding{GitCommonDir: common, Ref: ref, OID: oid, Blob: blob}})
	if err != nil {
		t.Fatal(err)
	}
	if !input.MatchesExpected || input.Binding.Blob != blob {
		t.Fatalf("input observation = %#v", input)
	}
}

func TestSystemGitObserverRejectsSubdirectoryAndDetachedHead(t *testing.T) {
	repo := initObserverRepository(t)
	ref := gitOutput(t, repo, "symbolic-ref", "HEAD")
	nested := filepath.Join(repo, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	observer := SystemDependencies().Git
	if _, err := observer.ObserveTarget(TargetRequest{Worktree: nested, Ref: ref}); !IsClass(err, ErrorTargetConflict) {
		t.Fatalf("subdirectory error = %v", err)
	}
	runGit(t, repo, "checkout", "--detach")
	if _, err := observer.ObserveTarget(TargetRequest{Worktree: repo, Ref: ref}); !IsClass(err, ErrorTargetConflict) {
		t.Fatalf("detached error = %v", err)
	}
}

func initObserverRepository(t *testing.T) string {
	t.Helper()
	repo, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "init", "-b", "main")
	runGit(t, repo, "config", "user.name", "WF01 Test")
	runGit(t, repo, "config", "user.email", "wf01@example.invalid")
	if err := os.WriteFile(filepath.Join(repo, "input.txt"), []byte("input\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "input.txt")
	runGit(t, repo, "commit", "-m", "base")
	return repo
}

func runGit(t *testing.T, directory string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, arguments...)...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v: %s", arguments, err, output)
	}
}

func gitOutput(t *testing.T, directory string, arguments ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, arguments...)...)
	output, err := command.Output()
	if err != nil {
		t.Fatalf("git %v: %v", arguments, err)
	}
	if len(output) == 0 || output[len(output)-1] != '\n' {
		t.Fatalf("git %v output is malformed", arguments)
	}
	return string(output[:len(output)-1])
}
