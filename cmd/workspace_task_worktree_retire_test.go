package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestTaskWorktreeRetirePassesApplyAndRendersJSON(t *testing.T) {
	var got []workspace.TaskWorktreeRetireInput
	c := newWorkspaceTaskCommandWithServices(workspaceTaskServices{retireWorktree: func(input workspace.TaskWorktreeRetireInput) (workspace.TaskWorktreeRetireResult, error) {
		got = append(got, input)
		return workspace.TaskWorktreeRetireResult{TaskID: input.TaskID, Applied: input.Apply, Worktree: "removed", Branch: "deleted", State: workspace.WorkItemRetired}, nil
	}})
	for _, args := range [][]string{{"worktree", "retire", "task"}, {"worktree", "retire", "task", "--apply", "--format", "json"}} {
		var out bytes.Buffer
		c.SetOut(&out)
		c.SetErr(&bytes.Buffer{})
		c.SetArgs(args)
		if err := c.Execute(); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		if len(args) > 3 && !strings.Contains(out.String(), `"kind": "WorkspaceTaskWorktreeRetire@1"`) {
			t.Fatalf("json %s", out.String())
		}
	}
	if len(got) != 2 || got[0].Apply || !got[1].Apply || got[1].TaskID != "task" {
		t.Fatalf("inputs %+v", got)
	}
}
