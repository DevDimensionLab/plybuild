package cmd

import (
	"fmt"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workspaceProjectServices struct {
	add  func(workspace.ProjectAddInput) (workspace.ProjectResult, error)
	show func(workspace.ProjectID) (workspace.ProjectResult, error)
	list func() (workspace.ProjectListResult, error)
}

func newWorkspaceProjectCommand(dependencies workspace.Dependencies) *cobra.Command {
	return newWorkspaceProjectCommandWithServices(workspaceProjectServices{
		add: func(input workspace.ProjectAddInput) (workspace.ProjectResult, error) {
			return workspace.AddProject(dependencies, input)
		},
		show: func(id workspace.ProjectID) (workspace.ProjectResult, error) {
			return workspace.ShowProject(dependencies, id)
		},
		list: func() (workspace.ProjectListResult, error) {
			return workspace.ListProjects(dependencies)
		},
	})
}

func newWorkspaceProjectCommandWithServices(services workspaceProjectServices) *cobra.Command {
	projectCommand := &cobra.Command{
		Use:     "project",
		Short:   "Manage projects in a Ply workspace",
		Long:    "Register and inspect projects owned by the containing Ply workspace.",
		Example: "  ply workspace project add ply --name Ply --wrapper ../ply --repo ply=../ply/main\n  ply workspace project list",
	}

	var name string
	var wrapper string
	var repositoryValues []string
	addCommand := &cobra.Command{
		Use:     "add <project-id>",
		Short:   "Register a project and its explicit repository members",
		Long:    "Register one project, its wrapper, and one or more explicit Git repository members in the containing Ply workspace.",
		Example: "  ply workspace project add ply --name Ply --wrapper ../ply --repo ply=../ply/main",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.ProjectInvalidArguments(fmt.Sprintf("expected exactly one project ID, got %d arguments", len(args)))
			}
			_, err := workspace.ParseProjectAddInput(args[0], name, wrapper, repositoryValues)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := workspace.ParseProjectAddInput(args[0], name, wrapper, repositoryValues)
			if err != nil {
				return err
			}
			result, err := services.add(input)
			if err != nil {
				return err
			}
			return renderProject(cmd, result, true)
		},
	}
	addCommand.Flags().StringVarP(&name, "name", "n", "", "project display name")
	addCommand.Flags().StringVarP(&wrapper, "wrapper", "w", "", "project wrapper directory")
	addCommand.Flags().StringArrayVarP(&repositoryValues, "repo", "r", nil, "repository member as <repo-id>=<path>")
	_ = addCommand.MarkFlagRequired("name")
	_ = addCommand.MarkFlagRequired("wrapper")
	_ = addCommand.MarkFlagRequired("repo")
	setProjectFlagErrors(addCommand)

	showCommand := &cobra.Command{
		Use:     "show <project-id>",
		Short:   "Show a registered project",
		Long:    "Show the stored wrapper and repository members for one project in the containing Ply workspace.",
		Example: "  ply workspace project show ply",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.ProjectInvalidArguments(fmt.Sprintf("expected exactly one project ID, got %d arguments", len(args)))
			}
			_, err := workspace.ParseProjectID(args[0])
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := workspace.ParseProjectID(args[0])
			if err != nil {
				return err
			}
			result, err := services.show(id)
			if err != nil {
				return err
			}
			return renderProject(cmd, result, false)
		},
	}
	setProjectFlagErrors(showCommand)

	listCommand := &cobra.Command{
		Use:     "list",
		Short:   "List registered projects",
		Long:    "List projects registered in the containing Ply workspace.",
		Example: "  ply workspace project list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return workspace.ProjectInvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := services.list()
			if err != nil {
				return err
			}
			if len(result.Projects) == 0 {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "No projects are registered in Ply workspace %s.\n", result.Workspace)
				return err
			}
			if _, err = fmt.Fprintf(cmd.OutOrStdout(), "Projects in Ply workspace %s:\n", result.Workspace); err != nil {
				return err
			}
			projects := append([]workspace.ProjectRecord(nil), result.Projects...)
			sort.Slice(projects, func(i, j int) bool { return projects[i].ID < projects[j].ID })
			for _, project := range projects {
				label := "repositories"
				if len(project.RepoIDs) == 1 {
					label = "repository"
				}
				if _, err = fmt.Fprintf(cmd.OutOrStdout(), "  %s: %s (%d %s) %s\n", project.ID, project.Name, len(project.RepoIDs), label, project.Wrapper); err != nil {
					return err
				}
			}
			return nil
		},
	}
	setProjectFlagErrors(listCommand)

	projectCommand.AddCommand(addCommand, showCommand, listCommand)
	return projectCommand
}

func setProjectFlagErrors(command *cobra.Command) {
	command.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return workspace.ProjectInvalidArguments(err.Error())
	})
}

func renderProject(command *cobra.Command, result workspace.ProjectResult, add bool) error {
	firstLine := "Project %s (%s) in Ply workspace %s.\n"
	if add && result.Created {
		firstLine = "Added project %s (%s) to Ply workspace %s.\n"
	} else if add {
		firstLine = "Project %s (%s) is already registered in Ply workspace %s.\n"
	}
	if _, err := fmt.Fprintf(command.OutOrStdout(), firstLine, result.Project.ID, result.Project.Name, result.Workspace); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(command.OutOrStdout(), "Wrapper: %s\nRepositories:\n", result.Project.Wrapper); err != nil {
		return err
	}
	repositories := append([]workspace.RepoRecord(nil), result.Repos...)
	sort.Slice(repositories, func(i, j int) bool { return repositories[i].ID < repositories[j].ID })
	for _, repository := range repositories {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "  %s:\n    locator: %s\n    git common directory: %s\n", repository.ID, repository.Locator, repository.GitCommonDir); err != nil {
			return err
		}
	}
	return nil
}
