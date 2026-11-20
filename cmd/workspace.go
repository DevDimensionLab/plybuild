package cmd

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceCommand(dependencies taskwork.Dependencies) *cobra.Command {
	workspaceCommand := &cobra.Command{
		Use:     "workspace",
		Short:   "Manage Ply workspaces",
		Long:    "Manage explicit local Ply workspaces.",
		Example: "  ply workspace init\n  ply workspace project list\n  ply workspace epic list\n  ply workspace task list",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	var discover bool
	var project string
	initCommand := &cobra.Command{
		Use:     "init",
		Short:   "Initialize a Ply workspace in the current directory",
		Long:    "Initialize a Ply workspace in the current directory by creating .ply/workspace.yaml and, when missing, ply-workspace.json without Projects.\nThe current directory does not need to be a Git repository.\nWith --discover it also defines one Project in ply-workspace.json (wrapper .) whose repositories are the Git worktree roots in <name> or <name>/main, with the URL of origin as remote.",
		Example: "  ply workspace init\n  ply workspace init --discover",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return workspace.InvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if project != "" && !discover {
				return workspace.InvalidArguments("--project requires --discover")
			}
			if discover {
				return runInitDiscover(cmd, dependencies, project)
			}
			result, err := workspace.Init(dependencies.Dependencies)
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
	initCommand.Flags().BoolVar(&discover, "discover", false, "also define one Project in ply-workspace.json with the Git repositories in <name> or <name>/main")
	initCommand.Flags().StringVar(&project, "project", "", "Project ID for --discover (default the directory name)")
	initCommand.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return workspace.InvalidArguments(err.Error())
	})
	workspaceCommand.AddCommand(initCommand, newWorkspaceProjectCommand(dependencies), newWorkspaceEpicCommand(dependencies), newWorkspaceTaskCommand(dependencies))
	workspaceCommand.AddCommand(newWorkspaceStatusCommand(dependencies), newWorkspaceAttentionCommand(dependencies), newWorkspaceActivityCommand(dependencies), newWorkspaceWorktreesCommand(dependencies))
	workspaceCommand.AddCommand(newWorkspaceOverviewCommand(dependencies), newWorkspaceDefinitionCommand(dependencies.Dependencies))
	return workspaceCommand
}

func runInitDiscover(cmd *cobra.Command, dependencies taskwork.Dependencies, project string) error {
	result, err := workspace.InitDiscover(dependencies.Dependencies, project)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	state := "already defined"
	if result.Changed {
		state = "defined"
	}
	fmt.Fprintf(out, "Ply workspace %s: Project %s %s in ply-workspace.json with %d repositories.\n", result.Root, result.ProjectID, state, len(result.Repos))
	for _, repo := range result.Repos {
		remote := repo.Remote
		if remote == "" {
			remote = "no origin"
		}
		fmt.Fprintf(out, "  %s  %s  %s\n", repo.ID, repo.Path, remote)
	}
	for _, omission := range result.Omissions {
		fmt.Fprintf(out, "  left out %s: %s\n", omission.Path, omission.Reason)
	}
	return nil
}

func init() {
	dependencies := taskwork.SystemDependencies()
	RootCmd.AddCommand(newWorkspaceCommand(dependencies))
}
