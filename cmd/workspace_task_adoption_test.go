package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"strings"
	"testing"
)

func TestTaskAdoptionCommandsRemoved(t *testing.T) {
	for _, name := range []string{"adopt", "release-adoption"} {
		t.Run(name, func(t *testing.T) {
			calls := 0
			s := workspaceTaskServices{createWorktree: func(workspace.TaskWorktreeCreateInput) (workspace.TaskWorktreeMutationResult, error) {
				calls++
				return workspace.TaskWorktreeMutationResult{}, nil
			}}
			c := newWorkspaceTaskCommandWithServices(s)
			found, _, err := c.Find([]string{"worktree", name})
			if err == nil && found.Name() == name {
				t.Fatalf("removed command %s remains reachable", name)
			}
			out, err := executeTaskCommand(t, s, "worktree", name, "task", "--apply", "--confirmation", "sha256:obsolete")
			if err == nil || calls != 0 || strings.Contains(out, "Adoption receipt") {
				t.Fatalf("obsolete command executed: calls=%d err=%v out=%s", calls, err, out)
			}
		})
	}
}
