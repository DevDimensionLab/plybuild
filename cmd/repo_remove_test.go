package cmd

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/reporemove"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type fakeRemove struct {
	requests []reporemove.Request
	result   reporemove.Result
	err      error
}

func (f *fakeRemove) run(ctx context.Context, req reporemove.Request) (reporemove.Result, error) {
	f.requests = append(f.requests, req)
	return f.result, f.err
}

func executeRemove(t *testing.T, cwd string, env map[string]string, fake *fakeRemove, args ...string) (string, string, error) {
	t.Helper()
	original, _, _ := RootCmd.Find([]string{"repo"})
	replacement := &cobra.Command{Use: "repo"}
	replacement.AddCommand(newRepoRemoveCommand(repoRemoveEnv{
		files:  cwdFiles{workspace.SystemDependencies().Files, cwd},
		getenv: func(name string) string { return env[name] },
		run:    fake.run,
	}))
	RootCmd.RemoveCommand(original)
	RootCmd.AddCommand(replacement)
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	RootCmd.SetOut(stdout)
	RootCmd.SetErr(stderr)
	RootCmd.SetArgs(append([]string{"repo", "remove"}, args...))
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

func TestRepoRemoveResolvesRequest(t *testing.T) {
	root, other := createWorkspace(t), createWorkspace(t)
	for name, tc := range map[string]struct {
		cwd  string
		env  map[string]string
		args []string
		want reporemove.Request
	}{
		"preview from Project root": {cwd: filepath.Join(root, "ply"), args: []string{"code-map"},
			want: reporemove.Request{Workspace: root, Project: "ply", RepoName: "code-map"}},
		"apply with flags": {cwd: filepath.Join(other, "ply"), env: map[string]string{"PLY_ACTOR": "env"},
			args: []string{"--workspace", root, "--project", "app", "--remote", "upstream", "--actor", "flag", "--apply", "--expect", "sha256:x", "code-map"},
			want: reporemove.Request{Workspace: root, Project: "app", RepoName: "code-map", Remote: "upstream", Actor: "flag", Apply: true, Expect: "sha256:x"}},
		"actor from environment": {cwd: filepath.Join(root, "app"), env: map[string]string{"PLY_ACTOR": "env"}, args: []string{"--apply", "--expect", "sha256:x", "code-map"},
			want: reporemove.Request{Workspace: root, Project: "app", RepoName: "code-map", Actor: "env", Apply: true, Expect: "sha256:x"}},
	} {
		t.Run(name, func(t *testing.T) {
			fake := &fakeRemove{}
			if _, stderr, err := executeRemove(t, tc.cwd, tc.env, fake, tc.args...); err != nil {
				t.Fatalf("stderr %q, error %v", stderr, err)
			}
			if len(fake.requests) != 1 || fake.requests[0] != tc.want {
				t.Fatalf("requests %+v", fake.requests)
			}
		})
	}
}

func TestRepoRemoveRejectsInvalidFlagsWithoutRunning(t *testing.T) {
	root := createWorkspace(t)
	for name, args := range map[string][]string{
		"apply without expect": {"--apply", "code-map"},
		"expect without apply": {"--expect", "sha256:x", "code-map"},
		"json":                 {"--json", "code-map"},
		"relative workspace":   {"--workspace", "relative", "--project", "ply", "code-map"},
	} {
		t.Run(name, func(t *testing.T) {
			fake := &fakeRemove{}
			stdout, stderr, err := executeRemove(t, filepath.Join(root, "ply"), nil, fake, args...)
			if ExitCode(err) != 1 || len(fake.requests) != 0 || stdout != "" || !strings.Contains(stderr, "invalid") {
				t.Fatalf("stdout %q, stderr %q, error %v", stdout, stderr, err)
			}
		})
	}
}

var removeResolved = reporemove.Resolved{Workspace: "/ws", Project: "ply", ProjectRoot: "/ws", RepoName: "code-map",
	Wrapper: "/ws/code-map", RepoDir: "/ws/code-map/main", Leaf: "main", Path: "code-map/main", Remote: "origin"}

func TestRepoRemoveOutput(t *testing.T) {
	plan := reporemove.Plan{Checks: []reporemove.Check{{Name: "provenance"}, {Name: "history", Detail: "unborn, no objects; remote not contacted"}},
		RemoteURL: "git@github.com:o/code-map.git", MetadataRevision: 14, Verdict: reporemove.VerdictUnborn, GitFiles: []string{"config"}, Digest: "sha256:abc"}
	resolvedLine := "Resolved Project ply: root /ws, wrapper /ws/code-map, repository /ws/code-map/main, repository ID code-map, remote origin\n"
	for name, tc := range map[string]struct {
		args   []string
		result reporemove.Result
		err    error
		stdout string
	}{
		"preview": {result: reporemove.Result{Resolved: removeResolved, Plan: plan}, stdout: resolvedLine +
			"check provenance: ok\ncheck history: ok (unborn, no objects; remote not contacted)\n" +
			"plan: clear wrapper metadata (revision 14), remove definitions member code-map, delete wrapper /ws/code-map\n" +
			"local Git files deleted with the repository: config\nplan digest: sha256:abc\n" +
			"Remote repository git@github.com:o/code-map.git is not changed.\n"},
		"apply": {args: []string{"--apply", "--expect", "sha256:abc", "--actor", "a"},
			result: reporemove.Result{Resolved: removeResolved, Plan: plan, Steps: []reporemove.StepResult{{Step: reporemove.StepMetadata}, {Step: reporemove.StepDefine, Changed: true}, {Step: reporemove.StepWrapper, Changed: true}}},
			stdout: resolvedLine + "metadata: unchanged\ndefine: removed\nwrapper: removed\nRemote repository git@github.com:o/code-map.git is not changed.\n"},
		"apply failure": {args: []string{"--apply", "--expect", "sha256:abc", "--actor", "a"},
			result: reporemove.Result{Resolved: removeResolved, Plan: plan, Steps: []reporemove.StepResult{{Step: reporemove.StepMetadata, Changed: true}}},
			err:    &reporemove.Error{Step: reporemove.StepDefine, Err: fmt.Errorf("%w: injected", reporemove.ErrUnknown)},
			stdout: resolvedLine + "metadata: removed\nRemote repository git@github.com:o/code-map.git is not changed.\n"},
		"check failure": {result: reporemove.Result{Resolved: removeResolved}, err: &reporemove.Error{Step: reporemove.StepCheck, Err: fmt.Errorf("%w: history unproven", reporemove.ErrConflict)},
			stdout: resolvedLine + "The remote repository is not changed.\n"},
		"nothing": {result: reporemove.Result{Resolved: removeResolved, Nothing: true},
			stdout: resolvedLine + "Nothing to remove: repository code-map is not defined in Project ply and wrapper /ws/code-map is absent.\nThe remote repository is not changed.\n"},
	} {
		t.Run(name, func(t *testing.T) {
			root := createWorkspace(t)
			fake := &fakeRemove{result: tc.result, err: tc.err}

			stdout, stderr, err := executeRemove(t, filepath.Join(root, "ply"), nil, fake, append(tc.args, "code-map")...)

			if stdout != tc.stdout {
				t.Fatalf("stdout:\n%s\nwant:\n%s", stdout, tc.stdout)
			}
			if tc.err == nil && err != nil || tc.err != nil && (ExitCode(err) != 1 || !reflect.DeepEqual(err, tc.err) || !strings.Contains(stderr, tc.err.Error())) {
				t.Fatalf("stderr %q, error %v", stderr, err)
			}
		})
	}
}
