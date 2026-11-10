package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestWorkflowCommandShape(t *testing.T) {
	command := newWorkflowCommand()
	root := &cobra.Command{Use: "ply"}
	root.AddCommand(command)
	if command.Use != "workflow" || command.Short != "Manage agent workflows" || !strings.Contains(command.Long, "planning repositories") || !strings.Contains(command.Example, "ply workflow epic plan") || !command.Runnable() || command.PersistentPreRunE == nil {
		t.Fatalf("workflow metadata = %#v", command)
	}
}
