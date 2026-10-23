package workspaceview

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"testing"
)

func TestWorkflowStatusParallelMembershipIsExplicitAndPure(t *testing.T) {
	s := workflowFixture()
	b := s.Registry.Tasks[0]
	b.ID = "task-b"
	s.Registry.Tasks = append(s.Registry.Tasks, b)
	s.Titles[b.ID] = workspace.TaskListTitle{Title: stringValue("B"), Source: "registration", Status: "available"}
	q := workspace.RegisteredTaskQueue{QueueID: "queue", Freshness: "fresh", Target: workspace.QueueTarget{ProjectID: "project", RepoID: "repo", EpicID: "epic"}, Active: []workspace.RegisteredQueueCurrent{
		{QueueCurrent: workspace.QueueCurrent{PreparationID: "A", TaskID: "task", Disposition: "active", State: "prepared"}, Freshness: "fresh"},
		{QueueCurrent: workspace.QueueCurrent{PreparationID: "B", TaskID: "task-b", Disposition: "active", State: "prepared"}, Freshness: "fresh"},
	}}
	queues := workspace.RegisteredTaskQueues{Freshness: "fresh", Queues: []workspace.RegisteredTaskQueue{q}}
	if _, e := BuildWorkflowStatus(s, WorkflowStatusOptions{SchemaVersion: 1, IncludeAll: true}, nil, taskrun.Inventory{}, queues); e == nil {
		t.Fatal("schema1 silently hid multiple active members")
	}
	out, e := BuildWorkflowStatus(s, WorkflowStatusOptions{SchemaVersion: 2, IncludeAll: true}, nil, taskrun.Inventory{}, queues)
	if e != nil {
		t.Fatal(e)
	}
	if out.Kind != "WorkflowStatusReadback@2" || len(out.Items) != 2 {
		t.Fatalf("wrong parallel view: %+v", out)
	}
	for _, row := range out.Items {
		if len(row.Queues) != 1 || row.Queues[0].State != "active" || len(row.Queues[0].Active) != 1 || row.Queues[0].Active[0].TaskID != row.TaskID {
			t.Fatalf("member hidden or misbound: %+v", row)
		}
	}
}
