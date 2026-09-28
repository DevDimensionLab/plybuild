package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestWorkspaceTaskCommandsRenderHumanAndMachineReadback(t *testing.T) {
	epic := workspace.EpicRecord{ID: "epic", Title: "Epic", ProjectID: "ply"}
	task := workspace.TaskRecord{ID: "task", Title: "Task", Description: "Description", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply", GitCommonDir: "/repo/.git", WorktreeState: workspace.WorkItemUnbound}
	services := workspaceTaskServices{
		create: func(input workspace.TaskCreateInput) (workspace.TaskMutationResult, error) {
			return workspace.TaskMutationResult{Workspace: "/workspace", Task: task, Epic: epic, Created: true}, nil
		},
		show: func(id workspace.TaskID) (workspace.TaskReadbackResult, error) {
			return workspace.TaskReadbackResult{Workspace: "/workspace", Task: task, Epic: epic, ProjectFreshness: "fresh", ParentFreshness: workspace.TaskFreshnessWorktree{Locator: "unknown", Ref: "unknown", OID: "unknown", Tree: "unknown", GitCommonDir: "unknown", Clean: "unknown", InventoryMatch: "unknown"}, SourceFreshness: workspace.TaskFreshnessSource{Ref: "unknown", OID: "unknown", Tree: "unknown", CheckedOutAt: "unknown"}, TargetFreshness: workspace.TaskFreshnessTarget{Kind: "fresh", Locator: "unknown", Ref: "unknown", OID: "unknown", Tree: "unknown", GitCommonDir: "unknown", Clean: "unknown", InventoryMatch: "unknown"}, Target: workspace.TaskObservedTarget{Kind: "absent"}, Reasons: []string{"task_worktree_unbound"}}, nil
		},
		list: func(id *workspace.EpicID) (workspace.TaskListResult, error) {
			return workspace.TaskListResult{Workspace: "/workspace", Tasks: []workspace.TaskRecord{task}}, nil
		},
		createWorktree: func(input workspace.TaskWorktreeCreateInput) (workspace.TaskWorktreeMutationResult, error) {
			return workspace.TaskWorktreeMutationResult{}, nil
		},
	}
	stdout, err := executeTaskCommand(t, services, "create", "task", "--title", "Task", "--description", "Description", "--epic", "epic", "--project", "ply", "--repo", "ply")
	if err != nil || !strings.Contains(stdout, "Created Task task: Task\n") || !strings.HasSuffix(stdout, "Agent started by this command: no\n") {
		t.Fatalf("create stdout=%q error=%v", stdout, err)
	}
	stdout, err = executeTaskCommand(t, services, "show", "task")
	if err != nil || !strings.Contains(stdout, "Worktree: not created\nState: unbound\n") || !strings.Contains(stdout, "Ready for handoff: no\n") {
		t.Fatalf("show stdout=%q error=%v", stdout, err)
	}
	stdout, err = executeTaskCommand(t, services, "show", "task", "--format", "json")
	if err != nil || !strings.Contains(stdout, `"kind":"WorkspaceTaskReadback@1"`) || strings.Count(stdout, "\n") != 1 {
		t.Fatalf("json stdout=%q error=%v", stdout, err)
	}
}

func executeTaskCommand(t *testing.T, services workspaceTaskServices, args ...string) (string, error) {
	t.Helper()
	command := newWorkspaceTaskCommandWithServices(services)
	output := &bytes.Buffer{}
	command.SetOut(output)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(args)
	command.SilenceErrors = true
	command.SilenceUsage = true
	err := command.Execute()
	return output.String(), err
}
