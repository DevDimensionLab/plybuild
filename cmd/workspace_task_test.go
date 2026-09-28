package cmd

import (
	"bytes"
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

func TestWorkspaceTaskLifecycleCommandsCallOneServiceAndPreserveStreams(t *testing.T) {
	resultCalls, qaCalls, integrationCalls := 0, 0, 0
	resultRecord := workspace.TaskResultRecord{ID: "trs_11111111111111111111111111111111", TaskID: "task", ResultOID: strings.Repeat("b", 40), TechnicalGate: "passed", HandoffID: "hnd_11111111111111111111111111111111", TerminalResultID: "res_11111111111111111111111111111111", StoreTransition: "none"}
	qaRecord := workspace.TaskHumanQARecord{ID: "hqa_11111111111111111111111111111111", TaskID: "task", TaskResultID: resultRecord.ID, ResultOID: resultRecord.ResultOID, Outcome: "pass"}
	readback := workspace.WorkspaceTaskIntegrationReadback{Value: canonicaljson.Object{{Name: "kind", Value: "WorkspaceTaskIntegrationReadback@1"}}, Classification: "ready", GitChanged: boolCommandPointer(false), RecoveryStatus: "safe-no-effect", NextAction: workspace.IntegrationNextAction{Kind: "apply_confirmed_plan", Reason: "Apply the exact plan.", Argv: []string{"ply", "workspace", "task", "integrate"}}, TaskID: "task", ResultOID: resultRecord.ResultOID, TechnicalGate: "passed", HumanQAOutcome: "pass", EpicID: "epic", ParentRef: "refs/heads/epic", ParentOID: strings.Repeat("a", 40)}
	services := workspaceTaskServices{
		recordResult: func(input workspace.TaskResultRecordInput) (workspace.TaskResultMutationResult, error) {
			resultCalls++
			return workspace.TaskResultMutationResult{Workspace: "/workspace", Record: resultRecord, Created: true}, nil
		},
		recordQA: func(input workspace.TaskHumanQARecordInput) (workspace.TaskHumanQAMutationResult, error) {
			qaCalls++
			return workspace.TaskHumanQAMutationResult{Workspace: "/workspace", Record: qaRecord, Created: true}, nil
		},
		integrate: func(input workspace.TaskIntegrationInput) (workspace.TaskIntegrationResult, error) {
			integrationCalls++
			return workspace.TaskIntegrationResult{Readback: readback}, nil
		},
	}
	stdout, err := executeTaskCommand(t, services, "result", "record", "task", "--file", "/tmp/result.json")
	if err != nil || resultCalls != 1 || !strings.Contains(stdout, "Recorded Task result trs_") || strings.Count(stdout, "\n") != 6 {
		t.Fatalf("result stdout=%q calls=%d error=%v", stdout, resultCalls, err)
	}
	stdout, err = executeTaskCommand(t, services, "qa", "record", "task", "--file", "/tmp/qa.json")
	if err != nil || qaCalls != 1 || !strings.Contains(stdout, "Recorded human QA hqa_") || strings.Count(stdout, "\n") != 6 {
		t.Fatalf("qa stdout=%q calls=%d error=%v", stdout, qaCalls, err)
	}
	stdout, err = executeTaskCommand(t, services, "integrate", "task", "--result", string(resultRecord.ID), "--qa", string(qaRecord.ID), "--expected-result-oid", resultRecord.ResultOID, "--expected-parent-oid", strings.Repeat("a", 40), "--check", "--format", "json")
	if err != nil || integrationCalls != 1 || stdout != "{\"kind\":\"WorkspaceTaskIntegrationReadback@1\"}\n" {
		t.Fatalf("integrate stdout=%q calls=%d error=%v", stdout, integrationCalls, err)
	}
	if stdout, err = executeTaskCommand(t, services, "integrate", "task", "--result", string(resultRecord.ID), "--qa", string(qaRecord.ID), "--expected-result-oid", resultRecord.ResultOID, "--expected-parent-oid", strings.Repeat("a", 40), "--check", "--apply"); err == nil || stdout != "" || integrationCalls != 1 {
		t.Fatalf("invalid mode stdout=%q calls=%d error=%v", stdout, integrationCalls, err)
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
