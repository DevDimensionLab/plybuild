package cmd

import (
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/spf13/cobra"
)

func TestWorkflowCommandShape(t *testing.T) {
	command := newWorkflowCommand(workflowhandoff.Dependencies{})
	root := &cobra.Command{Use: "ply"}
	root.AddCommand(command)
	if command.Use != "workflow" || command.Short != "Manage agent workflows" || command.Long != "Manage explicit local agent workflow transitions and report-ready notifications." || command.Example != "  ply workflow handoff show hnd_0123456789abcdef0123456789abcdef" || command.Runnable() || command.PersistentPreRunE == nil {
		t.Fatalf("workflow metadata = %#v", command)
	}
	handoff, _, err := command.Find([]string{"handoff"})
	if err != nil || handoff == command || handoff.Runnable() {
		t.Fatalf("handoff command = %#v, %v", handoff, err)
	}
}
