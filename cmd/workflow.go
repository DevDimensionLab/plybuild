package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/planningrepo"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workflownotification"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowCommand(dependencies workflowhandoff.Dependencies) *cobra.Command {
	command := &cobra.Command{
		Use:               "workflow",
		Short:             "Manage agent workflows",
		Long:              "Read registered workspace workflow status, trace one Task's recorded process across its runs, create local planning repositories, execute queued goals with an assigned agent, and manage workflow transitions and notifications.",
		Example:           "  ply workflow status\n  ply workflow status --all --json\n  ply workflow trace task-id --format json\n  ply workflow execute --check\n  ply workflow epic plan --project ply --language nb\n  ply workflow handoff show hnd_0123456789abcdef0123456789abcdef",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
	command.AddCommand(newWorkflowHandoffCommand(dependencies))
	command.AddCommand(newWorkflowRunCommand(taskrun.SystemDependencies(workspace.SystemDependencies())))
	command.AddCommand(newWorkflowExecuteCommand(taskrun.SystemDependencies(workspace.SystemDependencies())))
	command.AddCommand(newWorkflowEpicCommand(planningrepo.SystemDependencies()))
	setHandoffFlagErrors(command)
	command.AddCommand(newWorkflowStatusCommand(workspace.SystemDependencies()))
	command.AddCommand(newWorkflowTraceCommand(workspace.SystemDependencies()))
	command.AddCommand(NewWorkflowNotificationCommand(workflownotification.SystemDependencies()))
	return command
}

func init() {
	RootCmd.AddCommand(newWorkflowCommand(workflowhandoff.SystemDependencies()))
	RootCmd.Example += "  ply workflow status --json\n  ply workflow trace task-id --format json\n  ply workflow epic plan --repo /work/example/main --root /work/example --check"
}
