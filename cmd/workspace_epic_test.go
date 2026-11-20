package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
)

func TestWorkspaceEpicCommandsRenderMutationAndCanonicalJSON(t *testing.T) {
	oid := strings.Repeat("a", 40)
	tree := strings.Repeat("b", 40)
	epic := taskwork.EpicRecord{ID: "epic", Title: "Epic", ProjectID: "ply", RepoBindings: []taskwork.EpicRepoBinding{{RepoID: "ply", GitCommonDir: "/repo/.git", Worktree: taskwork.EpicWorktreeBinding{ID: "wt_11111111111111111111111111111111", OwnerKind: "epic", OwnerID: "epic", Origin: "adopted", Locator: "/repo/epic", Ref: "refs/heads/epic", OID: oid, Tree: tree, GitCommonDir: "/repo/.git"}}}}
	services := workspaceEpicServices{
		adopt: func(input taskwork.EpicAdoptInput) (taskwork.EpicMutationResult, error) {
			return taskwork.EpicMutationResult{Workspace: "/workspace", Epic: epic, Created: true}, nil
		},
		show: func(id taskwork.EpicID) (taskwork.EpicReadbackResult, error) {
			return taskwork.EpicReadbackResult{Workspace: "/workspace", Epic: epic, Observed: []taskwork.EpicObservedRepo{{RepoID: "ply"}}, Freshness: []taskwork.EpicFreshnessRepo{{RepoID: "ply", Locator: "unknown", Ref: "unknown", OID: "unknown", Tree: "unknown", GitCommonDir: "unknown", Clean: "unknown", InventoryMatch: "unknown"}}, Reasons: []string{"epic_observation_unknown"}}, nil
		},
		list: func() (taskwork.EpicListResult, error) {
			return taskwork.EpicListResult{Workspace: "/workspace", Epics: []taskwork.EpicRecord{epic}}, nil
		},
	}
	stdout, err := executeEpicCommand(t, services, "adopt", "epic", "--title", "Epic", "--project", "ply", "--repo", "ply", "--worktree", "/repo/epic", "--ref", "refs/heads/epic", "--expected-oid", oid)
	if err != nil || !strings.HasPrefix(stdout, "Adopted Epic epic: Epic\nProject / repository: ply / ply\n") || !strings.HasSuffix(stdout, "Git changed by this command: no\n") {
		t.Fatalf("adopt stdout=%q error=%v", stdout, err)
	}
	stdout, err = executeEpicCommand(t, services, "show", "epic", "--format", "json")
	if err != nil || !strings.Contains(stdout, `"kind":"WorkspaceEpicReadback@1"`) || strings.Count(stdout, "\n") != 1 {
		t.Fatalf("json stdout=%q error=%v", stdout, err)
	}
	stdout, err = executeEpicCommand(t, services, "list")
	if err != nil || stdout != "Epics in Ply workspace /workspace:\n  epic: Epic (Project ply; repositories: ply)\n" {
		t.Fatalf("list stdout=%q error=%v", stdout, err)
	}
}

func executeEpicCommand(t *testing.T, services workspaceEpicServices, args ...string) (string, error) {
	t.Helper()
	command := newWorkspaceEpicCommandWithServices(services)
	output := &bytes.Buffer{}
	command.SetOut(output)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(args)
	command.SilenceErrors = true
	command.SilenceUsage = true
	err := command.Execute()
	return output.String(), err
}
