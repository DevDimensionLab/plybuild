package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestWorkspaceEpicCreateRendersEffectsOfEarlyRejection(t *testing.T) {
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	old, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(old) })
	if err = os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	for format, want := range map[string]string{"text": "Epic creation new-epic: rejected\nEffects: worktree not_attempted, config not_attempted\n", "json": `"effects":{"worktree":"not_attempted","config":"not_attempted"}`} {
		command := newWorkspaceEpicCreateCommand(workspace.SystemDependencies())
		output := &bytes.Buffer{}
		command.SetOut(output)
		command.SetErr(&bytes.Buffer{})
		command.SetArgs([]string{"new-epic", "--title", "New", "--project", "ply", "--repo", "ply", "--parent-ref", "refs/heads/main", "--expected-parent-oid", strings.Repeat("a", 40), "--branch", "new-epic", "--path", filepath.Join(root, "new-epic"), "--format", format})
		command.SilenceErrors = true
		command.SilenceUsage = true
		if err = command.Execute(); err == nil || !strings.Contains(output.String(), want) {
			t.Fatalf("%s stdout=%q error=%v", format, output.String(), err)
		}
	}
}
