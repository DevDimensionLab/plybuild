package workflowhandoff

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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

func TestStatusPorcelainV2NormalizationMatrix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink mode assertions are Unix-specific")
	}
	worktree := physicalTempDir(t)
	files := map[string][]byte{
		"spaced æ.txt": []byte("modified\n"),
		"type.txt":     []byte("type changed\n"),
		"new name.txt": []byte("renamed\n"),
		"copied.txt":   []byte("copied\n"),
		"z odd.txt":    []byte("untracked\n"),
	}
	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(worktree, name), contents, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink("target with spaces", filepath.Join(worktree, "link")); err != nil {
		t.Fatal(err)
	}
	oid := strings.Repeat("1", 40)
	records := []string{
		fmt.Sprintf("1 .M N... 100644 100644 100644 %s %s spaced æ.txt", oid, oid),
		fmt.Sprintf("1 .D N... 100644 100644 000000 %s %s deleted.txt", oid, strings.Repeat("0", 40)),
		fmt.Sprintf("1 .T N... 100644 100644 120000 %s %s type.txt", oid, oid),
		fmt.Sprintf("2 R. N... 100644 100644 100644 %s %s R100 new name.txt", oid, oid),
		"old name.txt",
		fmt.Sprintf("2 C. N... 100644 100644 100644 %s %s C100 copied.txt", oid, oid),
		"source.txt",
		"? z odd.txt",
		fmt.Sprintf("1 .M N... 120000 120000 120000 %s %s link", oid, oid),
	}
	status, err := (systemGitObserver{files: systemFileSystem{}}).parseStatus(worktree, append([]byte(strings.Join(records, "\x00")), 0))
	if err != nil {
		t.Fatal(err)
	}
	if status.Mode != "exact_manifest" || len(status.Entries) != 9 {
		t.Fatalf("status = %#v", status)
	}
	wantPaths := []string{"copied.txt", "deleted.txt", "link", "new name.txt", "old name.txt", "source.txt", "spaced æ.txt", "type.txt", "z odd.txt"}
	gotPaths := make([]string, 0, len(status.Entries))
	byPath := map[string]StatusEntry{}
	for _, entry := range status.Entries {
		gotPaths = append(gotPaths, entry.Path)
		byPath[entry.Path] = entry
	}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Fatalf("sorted paths = %#v, want %#v", gotPaths, wantPaths)
	}
	if byPath["old name.txt"].IndexState != "deleted" || byPath["new name.txt"].IndexState != "added" || byPath["source.txt"].IndexState != "deleted" || byPath["copied.txt"].IndexState != "added" {
		t.Fatalf("rename/copy normalization = %#v", byPath)
	}
	if byPath["deleted.txt"].SHA256 != "" || byPath["type.txt"].WorktreeState != "type_changed" {
		t.Fatalf("delete/typechange normalization = %#v", byPath)
	}
	if byPath["link"].Kind != "symlink" || byPath["link"].SHA256 != digestBytes([]byte("target with spaces")) {
		t.Fatalf("symlink observation = %#v", byPath["link"])
	}
	if byPath["spaced æ.txt"].SHA256 != digestBytes(files["spaced æ.txt"]) || byPath["z odd.txt"].WorktreeState != "untracked" {
		t.Fatalf("odd path observations = %#v", byPath)
	}
}

func TestGitObserverUsesOnlyDeclaredLocalArgv(t *testing.T) {
	worktree := physicalTempDir(t)
	common := filepath.Join(worktree, ".git")
	if err := os.Mkdir(common, 0o700); err != nil {
		t.Fatal(err)
	}
	oid := strings.Repeat("1", 40)
	tree := strings.Repeat("2", 40)
	var calls [][]string
	runner := func(directory string, arguments ...string) gitResult {
		calls = append(calls, append([]string{directory}, arguments...))
		key := strings.Join(arguments, " ")
		outputs := map[string]string{
			"rev-parse --is-inside-work-tree":                   "true\n",
			"rev-parse --path-format=absolute --show-toplevel":  worktree + "\n",
			"rev-parse --path-format=absolute --git-common-dir": common + "\n",
			"symbolic-ref --quiet HEAD":                         "refs/heads/main\n",
			"rev-parse --verify refs/heads/main":                oid + "\n",
			"rev-parse --verify " + oid + "^{tree}":             tree + "\n",
			"status --porcelain=v2 --untracked-files=all -z":    "",
		}
		output, ok := outputs[key]
		if !ok {
			return gitResult{exit: 99, err: fmt.Errorf("unexpected git argv: %s", key)}
		}
		return gitResult{stdout: []byte(output)}
	}
	observed, err := (systemGitObserver{files: systemFileSystem{}, run: runner}).ObserveTarget(TargetRequest{Worktree: worktree, Ref: "refs/heads/main"})
	if err != nil {
		t.Fatal(err)
	}
	if observed.OID != oid || observed.Tree != tree || observed.GitCommonDir != common || observed.Status.Mode != "clean" {
		t.Fatalf("observation = %#v", observed)
	}
	want := [][]string{
		{worktree, "rev-parse", "--is-inside-work-tree"},
		{worktree, "rev-parse", "--path-format=absolute", "--show-toplevel"},
		{worktree, "rev-parse", "--path-format=absolute", "--git-common-dir"},
		{worktree, "symbolic-ref", "--quiet", "HEAD"},
		{worktree, "rev-parse", "--verify", "refs/heads/main"},
		{worktree, "rev-parse", "--verify", oid + "^{tree}"},
		{worktree, "status", "--porcelain=v2", "--untracked-files=all", "-z"},
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("git calls = %#v, want %#v", calls, want)
	}
	for _, call := range calls {
		joined := strings.Join(call[1:], " ")
		for _, forbidden := range []string{"fetch", "remote", "submodule", "clean", "reset"} {
			if strings.Contains(joined, forbidden) {
				t.Fatalf("forbidden git argv %q in %q", forbidden, joined)
			}
		}
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
