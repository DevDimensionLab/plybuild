package workspace

import (
	"strings"
	"testing"
)

func TestTaskAdoptionWireRejectedExplicitly(t *testing.T) {
	for _, wire := range []string{
		"worktree_operations: [{adoption: {kind: WorkspaceTaskAdoptionIntent@1}}]",
		"worktree_operations: [{release: {kind: WorkspaceTaskAdoptionRelease@1}}]",
		"worktree_operations: [{adoption_observations: []}]",
		"tasks: [{worktree: {origin: adopted}}]",
		"task_preparations: [{plan: {source_mode: adopted}}]",
		"task_preparations: [{plan: {adoption_id: wop_old}}]",
	} {
		t.Run(wire, func(t *testing.T) {
			_, err := decodeWorkItemRegistry([]byte("format_version: 7\n" + wire + "\n"))
			if err == nil || !strings.Contains(err.Error(), "Task adoption is no longer supported") {
				t.Fatalf("expected explicit unsupported Task adoption, got %v", err)
			}
		})
	}
}
