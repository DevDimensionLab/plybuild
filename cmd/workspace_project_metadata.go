package cmd

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
	"github.com/spf13/cobra"
)

type workspaceProjectMetadataServices struct {
	update func(workspace.ProjectMetadataInput) (workspace.ProjectMetadataMutation, error)
}

func addWorkspaceProjectMetadataCommands(command *cobra.Command, d workspace.Dependencies) {
	services := workspaceProjectMetadataServices{
		update: func(in workspace.ProjectMetadataInput) (workspace.ProjectMetadataMutation, error) {
			return workspace.UpdateProjectMetadata(d, in)
		},
	}
	command.AddCommand(newWorkspaceProjectMetadataCommandsWithServices(services)...)
}

func newWorkspaceProjectMetadataCommandsWithServices(services workspaceProjectMetadataServices) []*cobra.Command {
	companion := &cobra.Command{Use: "companion", Short: "Register planning, documentation, and other companion directories", Long: "Manage explicit companion directories of a Project in ply-workspace.json without creating directories or registering repositories. Registration validates a physical directory and survives later directory removal.", Example: "  ply workspace project companion add ply --role planning --path ./ply/planning"}
	companion.AddCommand(newProjectMetadataMutationCommand(workspace.ProjectCompanionAdd, services), newProjectMetadataMutationCommand(workspace.ProjectCompanionRemove, services))
	wrapper := &cobra.Command{Use: "repo-wrapper", Short: "Register a repository's enclosing directory", Long: "Manage the wrapper of a repository in ply-workspace.json. A registered wrapper must contain the repository locator. It does not change project or Git identity and remains registered if the directory later disappears.", Example: "  ply workspace project repo-wrapper set ply --repo ply --wrapper ./ply"}
	wrapper.AddCommand(newProjectMetadataMutationCommand(workspace.ProjectRepoWrapperSet, services), newProjectMetadataMutationCommand(workspace.ProjectRepoWrapperClear, services))
	return []*cobra.Command{companion, wrapper}
}

func newProjectMetadataMutationCommand(operation string, services workspaceProjectMetadataServices) *cobra.Command {
	var name, short, example string
	companion := operation == workspace.ProjectCompanionAdd || operation == workspace.ProjectCompanionRemove
	switch operation {
	case workspace.ProjectCompanionAdd:
		name, short, example = "add", "Register an explicit companion directory", "  ply workspace project companion add ply --role planning --path ./ply/planning"
	case workspace.ProjectCompanionRemove:
		name, short, example = "remove", "Remove a companion registration without deleting its directory", "  ply workspace project companion remove ply --role planning --path ./ply/planning"
	case workspace.ProjectRepoWrapperSet:
		name, short, example = "set", "Register the physical directory enclosing a repository", "  ply workspace project repo-wrapper set ply --repo ply --wrapper ./ply"
	case workspace.ProjectRepoWrapperClear:
		name, short, example = "clear", "Clear a repository wrapper registration without deleting files", "  ply workspace project repo-wrapper clear ply --repo ply"
	}
	var role, repo, path, format string
	input := func(args []string) workspace.ProjectMetadataInput {
		return workspace.ProjectMetadataInput{Operation: operation, ProjectID: workspace.ProjectID(args[0]), RepoID: workspace.RepoID(repo), Role: workspace.CompanionRole(role), Locator: path}
	}
	command := &cobra.Command{Use: name + " <project-id>", Short: short, Long: short + " in ply-workspace.json. An unchanged registration writes nothing. JSON output is the Project as 'ply workspace project show' gives it, with changed. No directories or Git state are changed.", Example: example, Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.ProjectInvalidArguments("expected exactly one Project ID")
		}
		if err := validateProjectFormat(format); err != nil {
			return err
		}
		return input(args).Validate()
	}, RunE: func(cmd *cobra.Command, args []string) error {
		out, err := services.update(input(args))
		if err != nil {
			return err
		}
		if format == "json" {
			raw, err := workspace.MarshalProjectMetadataMutation(out)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(append(raw, '\n'))
			return err
		}
		verb := "Updated"
		if !out.Changed {
			verb = "Kept"
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s Project %s metadata in %s.\n", verb, out.Project.ID, workspacedef.File)
		return err
	}}
	command.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	if companion {
		command.Flags().StringVar(&role, "role", "", "companion role (planning, docs, or other)")
		command.Flags().StringVar(&path, "path", "", "companion directory path")
		_ = command.MarkFlagRequired("role")
		_ = command.MarkFlagRequired("path")
	} else {
		command.Flags().StringVar(&repo, "repo", "", "registered repository ID belonging to the project")
		_ = command.MarkFlagRequired("repo")
		if operation == workspace.ProjectRepoWrapperSet {
			command.Flags().StringVar(&path, "wrapper", "", "existing directory containing the registered repository locator")
			_ = command.MarkFlagRequired("wrapper")
		}
	}
	setProjectFlagErrors(command)
	return command
}
