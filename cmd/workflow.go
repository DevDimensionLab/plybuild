package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/spf13/cobra"
)

func newWorkflowCommand(dependencies workflowhandoff.Dependencies) *cobra.Command {
	command := &cobra.Command{
		Use:               "workflow",
		Short:             "Manage agent workflows",
		Long:              "Manage explicit local agent workflow transitions.",
		Example:           "  ply workflow handoff show hnd_0123456789abcdef0123456789abcdef",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
	command.AddCommand(newWorkflowHandoffCommand(dependencies))
	setHandoffFlagErrors(command)
	return command
}

func init() { RootCmd.AddCommand(newWorkflowCommand(workflowhandoff.SystemDependencies())) }
