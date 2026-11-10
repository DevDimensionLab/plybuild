package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/repoestablish"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
	"github.com/spf13/cobra"
)

type cwdFiles struct {
	workspace.FileSystem
	cwd string
}

func (f cwdFiles) Getwd() (string, error) { return f.cwd, nil }

type fakeEstablish struct {
	requests []repoestablish.Request
	result   repoestablish.Result
	err      error
}

func (f *fakeEstablish) run(ctx context.Context, req repoestablish.Request) (repoestablish.Result, error) {
	f.requests = append(f.requests, req)
	return f.result, f.err
}

// createWorkspace builds a workspace with two Projects whose roots differ
// from the workspace root and from each other.
func createWorkspace(t *testing.T) string {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.Init(workspace.Dependencies{Files: cwdFiles{workspace.SystemDependencies().Files, root}}); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{"ply", "app"} {
		if err := os.Mkdir(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	definitions := `{"format_version": 1, "projects": [
  {"id": "ply", "name": "Ply", "wrapper": "ply", "repos": [], "epics": []},
  {"id": "app", "name": "App", "wrapper": "app", "repos": [], "epics": []}]}`
	if err := os.WriteFile(filepath.Join(root, workspacedef.File), []byte(definitions), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func executeCreate(t *testing.T, cwd string, env map[string]string, fake *fakeEstablish, args ...string) (string, string, error) {
	t.Helper()
	original, _, _ := RootCmd.Find([]string{"repo"})
	replacement := &cobra.Command{Use: "repo"}
	replacement.AddCommand(newRepoCreateCommand(repoCreateEnv{
		files:  cwdFiles{workspace.SystemDependencies().Files, cwd},
		getenv: func(name string) string { return env[name] },
		run:    fake.run,
	}))
	RootCmd.RemoveCommand(original)
	RootCmd.AddCommand(replacement)
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	RootCmd.SetOut(stdout)
	RootCmd.SetErr(stderr)
	RootCmd.SetArgs(append([]string{"repo", "create"}, args...))
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

func TestRepoCreateActorPrecedence(t *testing.T) {
	root := createWorkspace(t)
	cwd := filepath.Join(root, "ply")
	for _, tc := range []struct {
		env   map[string]string
		args  []string
		actor string
	}{
		{map[string]string{"PLY_ACTOR": "env"}, []string{"--actor", "flag", "code-map"}, "flag"},
		{map[string]string{"PLY_ACTOR": "env"}, []string{"code-map"}, "env"},
		{map[string]string{}, []string{"code-map"}, ""},
	} {
		fake := &fakeEstablish{}
		if _, _, err := executeCreate(t, cwd, tc.env, fake, tc.args...); err != nil {
			t.Fatal(err)
		}
		if len(fake.requests) != 1 || fake.requests[0].Actor != tc.actor {
			t.Fatalf("requests %+v, want actor %q", fake.requests, tc.actor)
		}
	}
}

func TestRepoCreateWorkspaceFlagWinsOverCwd(t *testing.T) {
	root, other := createWorkspace(t), createWorkspace(t)
	fake := &fakeEstablish{}
	if _, _, err := executeCreate(t, filepath.Join(other, "ply"), nil, fake, "--workspace", root, "--project", "app", "code-map"); err != nil {
		t.Fatal(err)
	}
	want := repoestablish.Request{Workspace: root, Project: "app", RepoName: "code-map"}
	if len(fake.requests) != 1 || fake.requests[0] != want {
		t.Fatalf("requests %+v", fake.requests)
	}

	outside := filepath.Dir(root)
	stdout, stderr, err := executeCreate(t, outside, nil, fake, "--project", "ply", "code-map")
	if ExitCode(err) != 1 || len(fake.requests) != 1 || stdout != "" || !strings.Contains(stderr, string(workspace.ErrorWorkspaceNotFound)) {
		t.Fatalf("outside: stdout %q, stderr %q, error %v", stdout, stderr, err)
	}
	_, stderr, err = executeCreate(t, root, nil, fake, "--workspace", "relative", "--project", "ply", "code-map")
	if ExitCode(err) != 1 || len(fake.requests) != 1 || stderr != "Error: invalid: --workspace must be an absolute path\n" {
		t.Fatalf("relative: stderr %q, error %v", stderr, err)
	}
}

func TestRepoCreateDiscoversProjectOnlyAtExactRoot(t *testing.T) {
	root := createWorkspace(t)
	fake := &fakeEstablish{}
	if _, _, err := executeCreate(t, filepath.Join(root, "app"), nil, fake, "code-map"); err != nil {
		t.Fatal(err)
	}
	if want := (repoestablish.Request{Workspace: root, Project: "app", RepoName: "code-map"}); len(fake.requests) != 1 || fake.requests[0] != want {
		t.Fatalf("requests %+v", fake.requests)
	}
	sub := filepath.Join(root, "app", "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, cwd := range []string{root, sub} {
		_, stderr, err := executeCreate(t, cwd, nil, fake, "code-map")
		if ExitCode(err) != 1 || len(fake.requests) != 1 || !strings.Contains(stderr, "--project") {
			t.Fatalf("%s: stderr %q, error %v", cwd, stderr, err)
		}
	}
	if _, _, err := executeCreate(t, sub, nil, fake, "--project", "ply", "code-map"); err != nil || fake.requests[1].Project != "ply" {
		t.Fatalf("override: %+v %v", fake.requests, err)
	}
}

// workspaceBytes reads the definitions file and every file under .ply.
func workspaceBytes(t *testing.T, root string) map[string]string {
	t.Helper()
	state := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		if rel, _ := filepath.Rel(root, path); rel != workspacedef.File && !strings.HasPrefix(rel, workspace.MarkerDirectory+string(filepath.Separator)) {
			return nil
		}
		raw, err := os.ReadFile(path)
		state[path] = string(raw)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func TestRepoCreateFromProjectRootLeavesEpicRuntimeUntouched(t *testing.T) {
	root := createWorkspace(t)
	repo := filepath.Join(root, "ply", "main")
	if out, err := exec.Command("git", "init", "-q", "-b", "main", repo).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, out)
	}
	definitions := `{"format_version": 1, "projects": [
  {"id": "ply", "name": "Ply", "wrapper": "ply", "repos": [{"id": "ply", "path": "ply/main"}],
   "epics": [{"id": "ply-main", "title": "Ply main", "repo": "ply", "path": "ply/main", "ref": "refs/heads/main"}]}]}`
	if err := os.WriteFile(filepath.Join(root, workspacedef.File), []byte(definitions), 0o644); err != nil {
		t.Fatal(err)
	}
	before := workspaceBytes(t, root)
	fake := &fakeEstablish{}
	if _, stderr, err := executeCreate(t, filepath.Join(root, "ply"), nil, fake, "code-map"); err != nil {
		t.Fatalf("stderr %q, error %v", stderr, err)
	}
	if want := (repoestablish.Request{Workspace: root, Project: "ply", RepoName: "code-map"}); len(fake.requests) != 1 || fake.requests[0] != want {
		t.Fatalf("requests %+v", fake.requests)
	}
	if after := workspaceBytes(t, root); !reflect.DeepEqual(after, before) {
		t.Fatalf("workspace changed:\nbefore %q\nafter  %q", before, after)
	}
}

func TestRepoCreatePassesOverridesAndRejectsJSON(t *testing.T) {
	root := createWorkspace(t)
	fake := &fakeEstablish{}
	_, _, err := executeCreate(t, filepath.Join(root, "ply"), nil, fake, "code-map", "--github-owner", "Org", "--leaf", "l", "--branch", "b", "--remote", "r")
	if want := (repoestablish.Overrides{GitHubOwner: "Org", Leaf: "l", Branch: "b", Remote: "r"}); err != nil || fake.requests[0].Overrides != want {
		t.Fatalf("requests %+v, error %v", fake.requests, err)
	}
	_, stderr, err := executeCreate(t, filepath.Join(root, "ply"), nil, fake, "--json", "code-map")
	if ExitCode(err) != 1 || len(fake.requests) != 1 || stderr != "Error: invalid: JSON output is not supported\n" {
		t.Fatalf("json: stderr %q, error %v", stderr, err)
	}
}

func TestRepoCreateReportsResolvedStepsAndNote(t *testing.T) {
	root := createWorkspace(t)
	resolved := repoestablish.Resolved{Workspace: root, Project: "ply", ProjectRoot: "/p", RepoName: "code-map", Wrapper: "/p/code-map", RepoDir: "/p/code-map/main",
		Leaf: "main", Branch: "main", Remote: "origin", GitHubOwner: "Org", Actor: "ply-epic"}
	header := "Resolved Project ply: root /p, wrapper /p/code-map, repository /p/code-map/main, owner Org, branch main, remote origin, actor ply-epic\n"
	fake := &fakeEstablish{result: repoestablish.Result{Resolved: resolved, Steps: []repoestablish.StepResult{{Step: repoestablish.StepWrapper, Changed: true}, {Step: repoestablish.StepInit}}},
		err: &repoestablish.Error{Step: repoestablish.StepGitHub, Err: fmt.Errorf("%w: GitHub repository Org/code-map is PUBLIC", repoestablish.ErrConflict)}}
	stdout, stderr, err := executeCreate(t, filepath.Join(root, "ply"), nil, fake, "code-map")
	if want := header + "wrapper: created\ninit: unchanged\n" + repoCreateNote + "\n"; ExitCode(err) != 1 || stdout != want {
		t.Fatalf("stdout %q, want %q, error %v", stdout, want, err)
	}
	if stderr != "Error: github: conflict: GitHub repository Org/code-map is PUBLIC\n" {
		t.Fatalf("stderr %q", stderr)
	}
	fake.err = &repoestablish.Error{Step: repoestablish.StepPrepare, Err: repoestablish.ErrInvalid}
	fake.result = repoestablish.Result{}
	if stdout, _, _ = executeCreate(t, filepath.Join(root, "ply"), nil, fake, "code-map"); stdout != repoCreateNote+"\n" {
		t.Fatalf("prepare failure stdout %q", stdout)
	}
}

func TestRepoCreateHelpSaysItCreatesAPrivateGitHubRepository(t *testing.T) {
	command := newRepoCreateCommand(repoCreateEnv{})
	if !strings.Contains(command.Long, "creates a private GitHub repository through gh") || command.Use != "create <repo-name>" {
		t.Fatalf("help %q", command.Long)
	}
	if found, _, err := newRepoCommand(nil).Find([]string{"create"}); err != nil || found.Name() != "create" {
		t.Fatalf("create is not registered under repo: %v", err)
	}
}

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

func TestRepoCreateReturnsOutputErrorAndKeepsOperationError(t *testing.T) {
	root := createWorkspace(t)
	writeErr := errors.New("stdout closed")
	stopped := &repoestablish.Error{Step: repoestablish.StepGitHub, Err: repoestablish.ErrConflict}
	for _, operationErr := range []error{nil, stopped} {
		fake := &fakeEstablish{result: repoestablish.Result{Steps: []repoestablish.StepResult{{Step: repoestablish.StepWrapper, Changed: true}}}, err: operationErr}
		command := newRepoCreateCommand(repoCreateEnv{files: cwdFiles{workspace.SystemDependencies().Files, filepath.Join(root, "ply")}, getenv: func(string) string { return "" }, run: fake.run})
		command.SetOut(failingWriter{writeErr})
		command.SetErr(&bytes.Buffer{})
		command.SetArgs([]string{"code-map"})
		command.SilenceUsage, command.SilenceErrors = true, true
		err := command.Execute()
		if !errors.Is(err, writeErr) || operationErr != nil && !errors.Is(err, operationErr) {
			t.Fatalf("operation error %v: got %v", operationErr, err)
		}
	}
}
