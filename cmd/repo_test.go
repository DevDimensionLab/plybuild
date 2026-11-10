package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type fakeGit struct {
	calls       int
	initialized bool
	branch      string
	initErr     error
}

func (f *fakeGit) Run(c process.Command) error {
	f.calls++
	switch c.Args[2] {
	case "check-ref-format":
		fmt.Fprintln(c.Stdout, c.Args[4])
	case "init":
		f.initialized = f.initErr == nil
		return f.initErr
	case "rev-parse":
		if !f.initialized {
			fmt.Fprintln(c.Stderr, "fatal: not a git repository")
			return exec.Command("/bin/sh", "-c", "exit 128").Run()
		}
		if c.Args[3] == "--is-inside-work-tree" {
			fmt.Fprintln(c.Stdout, "true")
		} else {
			fmt.Fprintln(c.Stdout, c.Args[1])
		}
	case "symbolic-ref":
		fmt.Fprintln(c.Stdout, f.branch)
	}
	return nil
}

func executeRepo(t *testing.T, runner process.Runner, args ...string) (string, string, error) {
	t.Helper()
	original, _, _ := RootCmd.Find([]string{"repo"})
	replacement := newRepoCommand(runner)
	RootCmd.RemoveCommand(original)
	RootCmd.AddCommand(replacement)
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	RootCmd.SetOut(stdout)
	RootCmd.SetErr(stderr)
	RootCmd.SetArgs(append([]string{"repo"}, args...))
	defer func() {
		RootCmd.RemoveCommand(replacement)
		RootCmd.AddCommand(original)
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
		RootCmd.SetArgs(nil)
		_ = RootCmd.PersistentFlags().Set("json", "false")
	}()
	err := ExecuteE()
	return stdout.String(), stderr.String(), err
}

func TestRepoInitReportsOutcomes(t *testing.T) {
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	created := filepath.Join(base, "new")
	for _, test := range []struct {
		dir            string
		git            *fakeGit
		stdout, stderr string
		exit           int
	}{
		{created, &fakeGit{branch: "main"}, "Initialized Git repository at " + created + " on branch main.\n", "", 0},
		{base, &fakeGit{branch: "main", initialized: true}, "Git repository already initialized at " + base + " on branch main.\n", "", 0},
		{base, &fakeGit{branch: "dev", initialized: true}, "", "Error: conflict: ", 1},
		{filepath.Join(base, "failed"), &fakeGit{branch: "main", initErr: errors.New("init failed")}, "", "Error: unknown: ", 1},
	} {
		stdout, stderr, err := executeRepo(t, test.git, "init", test.dir, "--branch", "main")

		if ExitCode(err) != test.exit || stdout != test.stdout || !strings.HasPrefix(stderr, test.stderr) {
			t.Fatalf("repo init stdout %q, stderr %q, exit %d; want %q, %q, %d", stdout, stderr, ExitCode(err), test.stdout, test.stderr, test.exit)
		}
	}
}

func TestRepoInitRejectsInvalidUsageWithoutGitCalls(t *testing.T) {
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init", dir},
		{"init", "--branch", "main"},
		{"init", dir, dir, "--branch", "main"},
		{"init", dir, "--branch", "main", "--json"},
	} {
		git := &fakeGit{}
		_, stderr, err := executeRepo(t, git, args...)

		if err == nil || ExitCode(err) != 1 || git.calls != 0 {
			t.Fatalf("%v: error %v, %d git calls", args, err, git.calls)
		}
		if args[len(args)-1] == "--json" && !strings.Contains(stderr, "Error: invalid: JSON output is not supported") {
			t.Fatalf("--json stderr %q", stderr)
		}
	}
}

func TestRepoCommandSkipsLegacyInitialization(t *testing.T) {
	repo, _, err := RootCmd.Find([]string{"repo"})
	if err != nil || repo.Name() != "repo" || repo.PersistentPreRunE == nil || repo.PersistentPreRunE(repo, nil) != nil {
		t.Fatalf("repo is not registered with its own no-op PersistentPreRunE: %v", err)
	}
}

type fakeRemote struct {
	url         string
	calls, adds int
}

func (f *fakeRemote) Run(c process.Command) error {
	f.calls++
	switch strings.Join(c.Args[2:4], " ") {
	case "rev-parse --is-inside-work-tree":
		fmt.Fprintln(c.Stdout, "true")
	case "rev-parse --path-format=absolute":
		fmt.Fprintln(c.Stdout, c.Args[1])
	case "symbolic-ref -q":
		fmt.Fprintln(c.Stdout, "main")
	case "config --show-scope":
		if f.url == "" {
			return exec.Command("/bin/sh", "-c", "exit 1").Run()
		}
		fmt.Fprintf(c.Stdout, "local\x00remote.origin.url\n%s\x00", f.url)
	case "remote get-url":
		if f.url == "" {
			fmt.Fprintln(c.Stderr, "error: No such remote 'origin'")
			return exec.Command("/bin/sh", "-c", "exit 2").Run()
		}
		fmt.Fprintln(c.Stdout, f.url)
	case "remote add":
		f.adds++
		f.url = c.Args[6]
	}
	return nil
}

func TestRepoRemoteSetReportsOutcomes(t *testing.T) {
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	url := "https://x:SECRET@example.invalid/a.git"
	git := &fakeRemote{}
	for _, test := range []struct {
		url, stdout, stderr string
		exit, adds          int
	}{
		{url, "Added Git remote origin to " + dir + ".\n", "", 0, 1},
		{url, "Git remote origin already set in " + dir + ".\n", "", 0, 1},
		{"https://example.invalid/other.git", "", "Error: conflict: ", 1, 1},
	} {
		stdout, stderr, err := executeRepo(t, git, "remote", "set", dir, "origin", test.url)

		if ExitCode(err) != test.exit || stdout != test.stdout || !strings.HasPrefix(stderr, test.stderr) || git.adds != test.adds ||
			strings.Contains(stderr, "example.invalid") {
			t.Fatalf("remote set stdout %q, stderr %q, exit %d, adds %d", stdout, stderr, ExitCode(err), git.adds)
		}
	}
}

func TestRepoRemoteSetRejectsInvalidUsageWithoutAdding(t *testing.T) {
	dir := t.TempDir()
	for _, args := range [][]string{
		{"remote", "set", dir, "origin"},
		{"remote", "set", dir, "origin", "u", "extra"},
		{"remote", "set", "--json", dir, "origin", "u"},
		{"remote", "set", dir, "origin", "u", "--json"},
	} {
		git := &fakeRemote{}
		_, stderr, err := executeRepo(t, git, args...)

		if ExitCode(err) != 1 || git.calls != 0 {
			t.Fatalf("%v: error %v, %d git calls", args, err, git.calls)
		}
		if strings.Contains(strings.Join(args, " "), "--json") && !strings.Contains(stderr, "Error: invalid: JSON output is not supported") {
			t.Fatalf("%v: stderr %q", args, stderr)
		}
	}
	t.Cleanup(func() { _ = RootCmd.PersistentFlags().Set("force", "false") })
	git := &fakeRemote{url: "https://example.invalid/old.git"}
	_, stderr, err := executeRepo(t, git, "remote", "set", "--force", dir, "origin", "https://example.invalid/new.git")
	if ExitCode(err) != 1 || git.adds != 0 || !strings.HasPrefix(stderr, "Error: conflict: ") {
		t.Fatalf("--force: error %v, stderr %q, %d adds", err, stderr, git.adds)
	}
}

// fakeGitHub answers gh reads for Org/svc with visibility, or NOT_FOUND when it is empty.
type fakeGitHub struct {
	visibility     string
	calls, creates int
}

func (f *fakeGitHub) Run(c process.Command) error {
	f.calls++
	if c.Args[0] == "repo" {
		f.creates++
		f.visibility = "PRIVATE"
		return nil
	}
	if f.visibility == "" {
		fmt.Fprint(c.Stdout, `{"data":{"repository":null},"errors":[{"type":"NOT_FOUND","path":["repository"]}]}`)
		return exec.Command("/bin/sh", "-c", "exit 1").Run()
	}
	fmt.Fprintf(c.Stdout, `{"data":{"repository":{"nameWithOwner":"Org/svc","visibility":%q,"isEmpty":true,"sshUrl":"git@github.com:Org/svc.git"}}}`, f.visibility)
	return nil
}

func TestRepoGitHubEnsureReportsOutcomes(t *testing.T) {
	gh := &fakeGitHub{}
	for _, test := range []struct {
		visibility, stdout, stderr string
		exit, creates              int
	}{
		{"", "git@github.com:Org/svc.git\n", "Created private GitHub repository org/svc.\n", 0, 1},
		{"PRIVATE", "git@github.com:Org/svc.git\n", "GitHub repository org/svc is already private and empty.\n", 0, 1},
		{"PUBLIC", "", "Error: conflict: GitHub repository org/svc is PUBLIC\n", 1, 1},
	} {
		gh.visibility = test.visibility
		stdout, stderr, err := executeRepo(t, gh, "github", "ensure", "org/svc")

		if ExitCode(err) != test.exit || stdout != test.stdout || stderr != test.stderr || gh.creates != test.creates {
			t.Fatalf("github ensure stdout %q, stderr %q, exit %d, creates %d", stdout, stderr, ExitCode(err), gh.creates)
		}
	}
}

func TestRepoGitHubEnsureRejectsInvalidUsageWithoutCalls(t *testing.T) {
	for _, args := range [][]string{
		{"github", "ensure"},
		{"github", "ensure", "a/b", "c/d"},
		{"github", "ensure", "ab"},
		{"github", "ensure", "a/b/c"},
		{"github", "ensure", "--json", "a/b"},
		{"github", "ensure", "a/b", "--json"},
	} {
		gh := &fakeGitHub{}
		_, stderr, err := executeRepo(t, gh, args...)

		if ExitCode(err) != 1 || gh.calls != 0 {
			t.Fatalf("%v: error %v, %d gh calls", args, err, gh.calls)
		}
		if strings.Contains(strings.Join(args, " "), "--json") && !strings.Contains(stderr, "Error: invalid: JSON output is not supported") {
			t.Fatalf("%v: stderr %q", args, stderr)
		}
	}
	t.Cleanup(func() { _ = RootCmd.PersistentFlags().Set("force", "false") })
	gh := &fakeGitHub{visibility: "PUBLIC"}
	_, stderr, err := executeRepo(t, gh, "github", "ensure", "--force", "Org/svc")
	if ExitCode(err) != 1 || gh.calls != 1 || gh.creates != 0 || !strings.HasPrefix(stderr, "Error: conflict: ") {
		t.Fatalf("--force: error %v, stderr %q, %d calls, %d creates", err, stderr, gh.calls, gh.creates)
	}
}
