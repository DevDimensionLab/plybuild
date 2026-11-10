package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workspaceProjectMetadataServices struct {
	show   func(workspace.ProjectID) (workspace.ProjectMetadataReadback, error)
	update func(workspace.ProjectMetadataInput) (workspace.ProjectMetadataMutation, error)
}

func addWorkspaceProjectMetadataCommands(command *cobra.Command, d workspace.Dependencies) {
	services := workspaceProjectMetadataServices{
		show: func(id workspace.ProjectID) (workspace.ProjectMetadataReadback, error) {
			return workspace.ShowProjectMetadata(d, id)
		},
		update: func(in workspace.ProjectMetadataInput) (workspace.ProjectMetadataMutation, error) {
			return workspace.UpdateProjectMetadata(d, in)
		},
	}
	command.AddCommand(newWorkspaceProjectMetadataCommandsWithServices(services)...)
}

func newWorkspaceProjectMetadataCommandsWithServices(services workspaceProjectMetadataServices) []*cobra.Command {
	metadata := &cobra.Command{Use: "metadata", Short: "Inspect explicit project companions and repository wrappers", Long: "Inspect registered project navigation metadata. Reading does not scan or require the registered directories to exist.", Example: "  ply workspace project metadata show ply --format json"}
	var format string
	show := &cobra.Command{Use: "show <project-id>", Short: "Show project navigation metadata and its history", Long: "Read explicit companion directories and repository wrappers with their recorded change history. Missing metadata gives empty companions and null repository wrappers. Revision is shared across project metadata in this workspace.", Example: "  ply workspace project metadata show ply --format json", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.ProjectInvalidArguments("expected exactly one Project ID")
		}
		if _, err := workspace.ParseProjectID(args[0]); err != nil {
			return err
		}
		return validateProjectFormat(format)
	}, RunE: func(cmd *cobra.Command, args []string) error {
		out, err := services.show(workspace.ProjectID(args[0]))
		if err != nil {
			return err
		}
		if format == "json" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Project %s metadata (workspace revision %d).\n", out.ProjectID, out.Revision); err != nil {
			return err
		}
		for _, companion := range out.Companions {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  Companion %s: %s\n", companion.Role, companion.Locator); err != nil {
				return err
			}
		}
		for _, repo := range out.Repositories {
			if repo.Wrapper != nil {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  Repository %s wrapper: %s\n", repo.RepoID, *repo.Wrapper); err != nil {
					return err
				}
			}
		}
		return nil
	}}
	show.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setProjectFlagErrors(show)
	metadata.AddCommand(show)
	companion := &cobra.Command{Use: "companion", Short: "Register planning, documentation, and other companion directories", Long: "Manage explicit companion directory metadata without creating directories or registering repositories. Registration validates a physical directory and survives later directory removal.", Example: "  ply workspace project companion add ply --role planning --path ./ply/planning --actor codex"}
	companion.AddCommand(newProjectMetadataMutationCommand(workspace.ProjectCompanionAdd, services), newProjectMetadataMutationCommand(workspace.ProjectCompanionRemove, services))
	wrapper := &cobra.Command{Use: "repo-wrapper", Short: "Register a repository's enclosing directory", Long: "Manage explicit repository wrapper metadata. A registered wrapper must contain the repository locator. It does not change project or Git identity and remains registered if the directory later disappears.", Example: "  ply workspace project repo-wrapper set ply --repo ply --wrapper ./ply --actor codex"}
	wrapper.AddCommand(newProjectMetadataMutationCommand(workspace.ProjectRepoWrapperSet, services), newProjectMetadataMutationCommand(workspace.ProjectRepoWrapperClear, services))
	return []*cobra.Command{metadata, companion, wrapper}
}

func newProjectMetadataMutationCommand(operation string, services workspaceProjectMetadataServices) *cobra.Command {
	var name, short, example string
	companion := operation == workspace.ProjectCompanionAdd || operation == workspace.ProjectCompanionRemove
	switch operation {
	case workspace.ProjectCompanionAdd:
		name, short, example = "add", "Register an explicit companion directory", "  ply workspace project companion add ply --role planning --path ./ply/planning --actor codex"
	case workspace.ProjectCompanionRemove:
		name, short, example = "remove", "Remove a companion registration without deleting its directory", "  ply workspace project companion remove ply --role planning --path ./ply/planning --actor codex"
	case workspace.ProjectRepoWrapperSet:
		name, short, example = "set", "Register the physical directory enclosing a repository", "  ply workspace project repo-wrapper set ply --repo ply --wrapper ./ply --actor codex"
	case workspace.ProjectRepoWrapperClear:
		name, short, example = "clear", "Clear a repository wrapper registration without deleting files", "  ply workspace project repo-wrapper clear ply --repo ply --actor codex"
	}
	var actor, role, repo, path, format string
	var expected int
	input := func(cmd *cobra.Command, args []string) workspace.ProjectMetadataInput {
		in := workspace.ProjectMetadataInput{Operation: operation, ProjectID: workspace.ProjectID(args[0]), RepoID: workspace.RepoID(repo), Role: workspace.CompanionRole(role), Locator: path, ActorClaim: actor}
		if cmd.Flags().Changed("expected-revision") {
			in.ExpectedRevision = &expected
		}
		return in
	}
	command := &cobra.Command{Use: name + " <project-id>", Short: short, Long: short + ". The actor is a recorded claim. An unchanged registration creates no history. --expected-revision checks the shared workspace project metadata revision before a change. No directories or Git state are changed.", Example: example, Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.ProjectInvalidArguments("expected exactly one Project ID")
		}
		if err := validateProjectFormat(format); err != nil {
			return err
		}
		return input(cmd, args).Validate()
	}, RunE: func(cmd *cobra.Command, args []string) error {
		out, err := services.update(input(cmd, args))
		if err != nil {
			return err
		}
		if format == "json" {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
		}
		verb := "Updated"
		if !out.Changed {
			verb = "Kept"
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s Project %s metadata (workspace revision %d).\n", verb, out.ProjectID, out.Revision)
		return err
	}}
	command.Flags().StringVar(&actor, "actor", "", "recorder identity or actor claim")
	command.Flags().IntVar(&expected, "expected-revision", 0, "require this workspace project metadata revision before writing")
	command.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	_ = command.MarkFlagRequired("actor")
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
