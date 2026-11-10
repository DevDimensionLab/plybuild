package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/agentwork"
	"github.com/devdimensionlab/plybuild/internal/humaninput"
	"github.com/devdimensionlab/plybuild/internal/planningrepo"
	"github.com/devdimensionlab/plybuild/internal/workflownotification"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowCommand() *cobra.Command {
	command := &cobra.Command{
		Use:               "workflow",
		Args:              cobra.NoArgs,
		RunE:              func(c *cobra.Command, args []string) error { return c.Help() },
		Short:             "Manage agent workflows",
		Long:              "Read registered workspace workflow status, trace one Task's recorded history, create local planning repositories, start an explicitly assigned agent, and manage notifications.",
		Example:           "  ply workflow status\n  ply workflow status --all --json\n  ply workflow trace task-id --format json\n  ply workflow epic plan --project ply --language nb",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
	command.AddCommand(newWorkflowGoalStartCommand(workspace.SystemDependencies()))
	command.AddCommand(newWorkflowEpicCommand(planningrepo.SystemDependencies()))
	command.AddCommand(newAgentWorkCommand(agentwork.System))
	command.AddCommand(newFollowupCommand(agentwork.System))
	command.AddCommand(newWorkflowStatusCommand(workspace.SystemDependencies()))
	command.AddCommand(newWorkflowTraceCommand(workspace.SystemDependencies()))
	command.AddCommand(NewWorkflowNotificationCommand(workflownotification.SystemDependencies()))
	command.AddCommand(newWorkflowInputCommand(humaninput.SystemDependencies()))
	return command
}

func init() {
	RootCmd.AddCommand(newWorkflowCommand())
	RootCmd.Example += "  ply workflow status --json\n  ply workflow trace task-id --format json\n  ply workflow epic plan --repo /work/example/main --root /work/example --check"
}
