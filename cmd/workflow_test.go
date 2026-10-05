package cmd

import (
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/spf13/cobra"
)

func TestWorkflowCommandShape(t *testing.T) {
	command := newWorkflowCommand(workflowhandoff.Dependencies{})
	root := &cobra.Command{Use: "ply"}
	root.AddCommand(command)
	if command.Use != "workflow" || command.Short != "Manage agent workflows" || !strings.Contains(command.Long, "planning repositories") || !strings.Contains(command.Example, "ply workflow epic plan") || command.Runnable() || command.PersistentPreRunE == nil {
		t.Fatalf("workflow metadata = %#v", command)
	}
	handoff, _, err := command.Find([]string{"handoff"})
	if err != nil || handoff == command || handoff.Runnable() {
		t.Fatalf("handoff command = %#v, %v", handoff, err)
	}
}
