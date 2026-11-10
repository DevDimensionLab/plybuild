package cmd

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/gitrepo"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workspaceDefinitionServices struct {
	addRepo func(project, repo, path string) (workspace.DefinitionRepoResult, error)
}

func newWorkspaceDefinitionCommand(dependencies workspace.Dependencies) *cobra.Command {
	return newWorkspaceDefinitionCommandWithServices(workspaceDefinitionServices{
		addRepo: func(project, repo, path string) (workspace.DefinitionRepoResult, error) {
			return workspace.AddDefinitionRepo(dependencies, project, repo, path)
		},
	})
}

func newWorkspaceDefinitionCommandWithServices(services workspaceDefinitionServices) *cobra.Command {
	definitionCommand := &cobra.Command{
		Use:   "definition",
		Short: "Manage the workspace definitions in ply-workspace.json",
	}
	repoCommand := &cobra.Command{Use: "repo", Short: "Manage repository members of defined Projects"}
	repoCommand.AddCommand(&cobra.Command{
		Use:     "add <project-id> <repo-id> <path>",
		Short:   "Add an existing Git worktree root to a defined Project, or verify that it is already there",
		Long:    "Add the Git worktree root at path as the last repository member of the Project in ply-workspace.json.\nThe path is stored relative to the workspace. An identical member is left unchanged; anything else is a conflict and is never overwritten.",
		Example: "  ply workspace definition repo add ply ply-code-map ./ply-code-map",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return fmt.Errorf("%w: JSON output is not supported", gitrepo.ErrInvalid)
			}
			if force, _ := cmd.Flags().GetBool("force"); force {
				return fmt.Errorf("%w: --force is not supported", gitrepo.ErrInvalid)
			}
			result, err := services.addRepo(args[0], args[1], args[2])
			if err != nil {
				return err
			}
			format := "Repository %s (%s) already in Project %s.\n"
			if result.Changed {
				format = "Added repository %s (%s) to Project %s.\n"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), format, result.RepoID, result.Path, result.ProjectID)
			return err
		},
	})
	definitionCommand.AddCommand(repoCommand)
	return definitionCommand
}
