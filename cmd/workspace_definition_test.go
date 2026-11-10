package cmd

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type fakeDefinitionRepos struct {
	calls  [][]string
	result workspace.DefinitionRepoResult
	err    error
}

func (f *fakeDefinitionRepos) add(project, repo, path string) (workspace.DefinitionRepoResult, error) {
	f.calls = append(f.calls, []string{project, repo, path})
	return f.result, f.err
}

func executeDefinition(t *testing.T, fake *fakeDefinitionRepos, args ...string) (string, string, error) {
	t.Helper()
	workspaceCommand, _, err := RootCmd.Find([]string{"workspace"})
	if err != nil {
		t.Fatal(err)
	}
	original, _, err := RootCmd.Find([]string{"workspace", "definition"})
	if err != nil || original == workspaceCommand {
		t.Fatalf("definition command is not registered: %v", err)
	}
	replacement := newWorkspaceDefinitionCommandWithServices(workspaceDefinitionServices{addRepo: fake.add})
	workspaceCommand.RemoveCommand(original)
	workspaceCommand.AddCommand(replacement)
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	RootCmd.SetOut(stdout)
	RootCmd.SetErr(stderr)
	RootCmd.SetArgs(append([]string{"workspace", "definition"}, args...))
	defer func() {
		workspaceCommand.RemoveCommand(replacement)
		workspaceCommand.AddCommand(original)
		RootCmd.SetOut(nil)
		RootCmd.SetErr(nil)
		RootCmd.SetArgs(nil)
		_ = RootCmd.PersistentFlags().Set("json", "false")
		_ = RootCmd.PersistentFlags().Set("force", "false")
	}()
	err = ExecuteE()
	return stdout.String(), stderr.String(), err
}

func TestWorkspaceDefinitionRepoAddReportsOutcomes(t *testing.T) {
	for _, tc := range []struct {
		changed bool
		want    string
	}{
		{true, "Added repository ply-code-map (ply-code-map) to Project ply.\n"},
		{false, "Repository ply-code-map (ply-code-map) already in Project ply.\n"},
	} {
		fake := &fakeDefinitionRepos{result: workspace.DefinitionRepoResult{Workspace: "/w", ProjectID: "ply", RepoID: "ply-code-map", Path: "ply-code-map", Changed: tc.changed}}
		stdout, stderr, err := executeDefinition(t, fake, "repo", "add", "ply", "ply-code-map", "../ply-code-map")
		if err != nil || stdout != tc.want || stderr != "" {
			t.Fatalf("stdout %q, stderr %q, error %v", stdout, stderr, err)
		}
		if want := [][]string{{"ply", "ply-code-map", "../ply-code-map"}}; !reflect.DeepEqual(fake.calls, want) {
			t.Fatalf("calls %v", fake.calls)
		}
	}
}

func TestWorkspaceDefinitionRepoAddReportsErrors(t *testing.T) {
	fake := &fakeDefinitionRepos{err: workspace.ProjectInvalidArguments("Project other is not defined")}
	stdout, stderr, err := executeDefinition(t, fake, "repo", "add", "other", "x", "x")
	if ExitCode(err) != 1 || stdout != "" || stderr != "Error: workspace_project_invalid_arguments: Project other is not defined\n" {
		t.Fatalf("stdout %q, stderr %q, error %v", stdout, stderr, err)
	}
}

func TestWorkspaceDefinitionRepoAddRejectsBeforeCall(t *testing.T) {
	for _, tc := range []struct {
		args   []string
		stderr string
	}{
		{[]string{"repo", "add", "--json", "ply", "x", "x"}, "Error: invalid: JSON output is not supported\n"},
		{[]string{"repo", "add", "--force", "ply", "x", "x"}, "Error: invalid: --force is not supported\n"},
		{[]string{"repo", "add", "ply", "x"}, "Error: accepts 3 arg(s), received 2\n"},
	} {
		fake := &fakeDefinitionRepos{err: errors.New("must not be called")}
		stdout, stderr, err := executeDefinition(t, fake, tc.args...)
		if ExitCode(err) != 1 || stdout != "" || stderr != tc.stderr || len(fake.calls) != 0 {
			t.Fatalf("%v: stdout %q, stderr %q, error %v, calls %v", tc.args, stdout, stderr, err, fake.calls)
		}
	}
}
