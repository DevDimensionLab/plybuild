package cmd

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceCommand(dependencies workspace.Dependencies) *cobra.Command {
	workspaceCommand := &cobra.Command{
		Use:     "workspace",
		Short:   "Manage Ply workspaces",
		Long:    "Manage explicit local Ply workspaces.",
		Example: "  ply workspace init\n  ply workspace project list\n  ply workspace epic list\n  ply workspace task list",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	initCommand := &cobra.Command{
		Use:     "init",
		Short:   "Initialize a Ply workspace in the current directory",
		Long:    "Initialize a Ply workspace in the current directory by creating .ply/workspace.yaml.\nThe current directory does not need to be a Git repository.",
		Example: "  ply workspace init",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return workspace.InvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := workspace.Init(dependencies)
			if err != nil {
				return err
			}
			if result.Created {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "Initialized Ply workspace at %s.\n", result.Root)
			} else {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "Ply workspace already initialized at %s.\n", result.Root)
			}
			return err
		},
	}
	initCommand.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return workspace.InvalidArguments(err.Error())
	})
	workspaceCommand.AddCommand(initCommand, newWorkspaceProjectCommand(dependencies), newWorkspaceEpicCommand(dependencies), newWorkspaceTaskCommand(dependencies))
	workspaceCommand.AddCommand(newWorkspaceStatusCommand(dependencies), newWorkspaceAttentionCommand(dependencies), newWorkspaceActivityCommand(dependencies), newWorkspaceRunsCommand(dependencies), newWorkspaceWorktreesCommand(dependencies))
	workspaceCommand.AddCommand(newWorkspaceOverviewCommand(dependencies))
	return workspaceCommand
}

func init() {
	dependencies := taskrun.SystemDependencies(workspace.SystemDependencies()).Workspace
	RootCmd.AddCommand(newWorkspaceCommand(dependencies))
}
