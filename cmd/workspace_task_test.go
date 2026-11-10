package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
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
		list: func(filters workspace.TaskListFilters) (workspace.TaskListResult, error) {
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

func TestWorkspaceTaskIntegrateCallsOneServiceAndPreservesStreams(t *testing.T) {
	integrationCalls := 0
	resultOID := strings.Repeat("b", 40)
	readback := workspace.WorkspaceTaskIntegrationReadback{Value: canonicaljson.Object{{Name: "kind", Value: "WorkspaceTaskIntegrationReadback@6"}}, Classification: "ready", GitChanged: boolCommandPointer(false), RecoveryStatus: "safe-no-effect", NextAction: workspace.IntegrationNextAction{Kind: "apply_confirmed_plan", Reason: "Apply the exact plan.", Argv: []string{"ply", "workspace", "task", "integrate"}}, TaskID: "task", ResultOID: resultOID, EpicID: "epic", ParentRef: "refs/heads/epic", ParentOID: strings.Repeat("a", 40)}
	services := workspaceTaskServices{
		integrate: func(input workspace.TaskIntegrationInput) (workspace.TaskIntegrationResult, error) {
			integrationCalls++
			return workspace.TaskIntegrationResult{Readback: readback}, nil
		},
	}
	stdout, err := executeTaskCommand(t, services, "integrate", "task", "--expected-result-oid", resultOID, "--expected-parent-oid", strings.Repeat("a", 40), "--check", "--format", "json")
	if err != nil || integrationCalls != 1 || stdout != "{\"kind\":\"WorkspaceTaskIntegrationReadback@6\"}\n" {
		t.Fatalf("integrate stdout=%q calls=%d error=%v", stdout, integrationCalls, err)
	}
	if stdout, err = executeTaskCommand(t, services, "integrate", "task", "--expected-result-oid", resultOID, "--expected-parent-oid", strings.Repeat("a", 40), "--check", "--apply"); err == nil || stdout != "" || integrationCalls != 1 {
		t.Fatalf("invalid mode stdout=%q calls=%d error=%v", stdout, integrationCalls, err)
	}
}

func TestWorkspaceTaskIntegratePassesExplicitRefreshParent(t *testing.T) {
	var inputs []workspace.TaskIntegrationInput
	services := workspaceTaskServices{
		integrate: func(input workspace.TaskIntegrationInput) (workspace.TaskIntegrationResult, error) {
			inputs = append(inputs, input)
			return workspace.TaskIntegrationResult{Readback: workspace.WorkspaceTaskIntegrationReadback{Value: canonicaljson.Object{{Name: "kind", Value: "WorkspaceTaskIntegrationReadback@7"}}}}, nil
		},
	}
	args := []string{"integrate", "task", "--expected-result-oid", strings.Repeat("b", 40), "--expected-parent-oid", strings.Repeat("a", 40), "--check", "--format", "json"}
	if _, err := executeTaskCommand(t, services, args...); err != nil {
		t.Fatal(err)
	}
	if _, err := executeTaskCommand(t, services, append(args, "--refresh-parent")...); err != nil {
		t.Fatal(err)
	}
	if len(inputs) != 2 || inputs[0].RefreshParent || !inputs[1].RefreshParent {
		t.Fatalf("refresh inputs = %#v", inputs)
	}
}

func boolCommandPointer(value bool) *bool { return &value }

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

func TestTaskWorktreeReceiptReplayDoesNotClaimFreshHandoff(t *testing.T) {
	for _, version := range []int{0, 1} {
		t.Run(fmt.Sprint(version), func(t *testing.T) {
			services := workspaceTaskServices{createWorktree: func(workspace.TaskWorktreeCreateInput) (workspace.TaskWorktreeMutationResult, error) {
				return workspace.TaskWorktreeMutationResult{Outcome: "already_ready", Task: workspace.TaskRecord{ID: "task", WorktreeState: workspace.WorkItemReady, Worktree: &workspace.TaskWorktreeBinding{NativeReceiptVersion: version, Locator: "/missing/checkout", Ref: "refs/heads/task"}}}, nil
			}}
			out, err := executeTaskCommand(t, services, "worktree", "create", "task", "--branch", "task", "--path", "/missing/checkout", "--expected-parent-oid", strings.Repeat("a", 40))
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(out, "Fresh observation: clean") || strings.Contains(out, "Ready for handoff: yes") || !strings.Contains(out, "Fresh observation: unknown") {
				t.Fatalf("historical receipt claimed fresh readiness: %s", out)
			}
		})
	}
}
