package workspaceview

import (
	"slices"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestCloseoutProgressDistinguishesRetiredKeptAndPending(t *testing.T) {
	for _, resource := range []string{"retired", "kept", "retained"} {
		t.Run(resource, func(t *testing.T) {
			s := progressFixture()
			state := "complete"
			want := "completed"
			if resource == "retained" {
				state, want = "cleanup_pending", "cleanup_pending"
			}
			s.Closeouts = map[workspace.TaskID]workspace.TaskCloseoutReceipt{"task": {OperationID: "integration-one", Plan: workspace.TaskCloseoutPlan{TaskID: "task", TaskResultID: "result"}, State: state, ResourceState: resource, WorktreeRemoved: resource == "retired", UpdatedAtUTC: "2026-10-08T12:00:00Z"}}
			list, err := BuildTaskList(s, workspace.TaskListFilters{}, nil)
			if err != nil {
				t.Fatal(err)
			}
			row := list.Tasks[0]
			if row.Progress.State != want {
				t.Fatalf("progress %+v", row.Progress)
			}
			if resource == "retired" && row.WorktreeState != workspace.WorkItemRetired {
				t.Fatal("removed source still ready")
			}
			if resource == "retained" && len(row.Progress.NextActions) != 1 {
				t.Fatal("pending cleanup hidden")
			}
			if resource != "retained" && len(row.Progress.NextActions) != 0 {
				t.Fatal("completed Task has active instructions")
			}
		})
	}
}

func TestRetiredTaskDoesNotClaimAReusedWorktreePath(t *testing.T) {
	s := progressFixture()
	s.Registry.Tasks[0].WorktreeState = workspace.WorkItemRetired
	s.Registry.Tasks[0].Worktree = &workspace.TaskWorktreeBinding{Locator: "/work/reused", Ref: "refs/heads/task", GitCommonDir: "/work/repo/.git"}
	if owners := inventoryOwners(s, "project", workspace.RepoRecord{ID: "repo", GitCommonDir: "/work/repo/.git"}); owners["/work/reused"].kind != "" {
		t.Fatalf("historical Task claims active resource: %+v", owners)
	}
}

func TestRequiredInstallationFailureRemainsVisibleAcrossWorkspaceViews(t *testing.T) {
	for _, phase := range []string{"installation_pending", "install_failed", "install_unknown"} {
		t.Run(phase, func(t *testing.T) {
			s := progressFixture()
			s.Closeouts = map[workspace.TaskID]workspace.TaskCloseoutReceipt{"task": {OperationID: "integration-one", Plan: workspace.TaskCloseoutPlan{TaskID: "task", TaskResultID: "result"}, State: "cleanup_pending", ResourceState: "retained", Phase: phase, UpdatedAtUTC: "2026-10-08T12:00:00Z"}}
			list, err := BuildTaskList(s, workspace.TaskListFilters{}, nil)
			if err != nil || len(list.Tasks) != 1 || !slices.Contains(list.Tasks[0].Progress.Reasons, phase) || list.Tasks[0].Progress.State != "cleanup_pending" {
				t.Fatalf("Task installation hidden: %+v %v", list, err)
			}
			epic, err := BuildEpics(s, EpicFilters{}, nil)
			if err != nil || len(epic.Epics) != 1 || epic.Epics[0].ProgressCounts.Attention != 1 {
				t.Fatalf("Epic installation hidden: %+v %v", epic, err)
			}
			attention, err := BuildAttention(s, workspace.TaskListFilters{}, nil)
			if err != nil || len(attention.Items) != 1 || attention.Items[0].Kind == "resume_closeout" {
				t.Fatalf("Workspace installation misreported as cleanup: %+v %v", attention, err)
			}
			status, err := BuildWorkflowStatus(s, WorkflowStatusOptions{}, nil, taskrun.Inventory{}, workspace.RegisteredTaskQueues{})
			if err != nil || len(status.Items) != 1 || status.Items[0].Category != "follow_up" || !slices.Contains(status.Items[0].Progress.Reasons, phase) {
				t.Fatalf("Workflow installation hidden: %+v %v", status, err)
			}
		})
	}
}
