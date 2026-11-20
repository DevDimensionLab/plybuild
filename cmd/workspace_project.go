package cmd

import (
	"fmt"
	"io"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workspaceProjectServices struct {
	addMember    func(workspace.ProjectID, workspace.RepoInput) (workspace.ProjectResult, error)
	moveMember   func(workspace.ProjectID, workspace.RepoID, workspace.RepoInput) (workspace.ProjectResult, error)
	removeMember func(workspace.ProjectID, workspace.RepoID) (workspace.DefinitionRepoResult, error)
	add          func(workspace.ProjectAddInput) (workspace.ProjectResult, error)
	show         func(workspace.ProjectID) (workspace.ProjectResult, error)
	list         func() (workspace.ProjectListResult, error)
}

func newWorkspaceProjectCommand(dependencies taskwork.Dependencies) *cobra.Command {
	command := newWorkspaceProjectCommandWithServices(workspaceProjectServices{
		addMember: func(id workspace.ProjectID, repo workspace.RepoInput) (workspace.ProjectResult, error) {
			return workspace.AddProjectMember(dependencies.Dependencies, id, repo)
		},
		moveMember: func(id workspace.ProjectID, repo workspace.RepoID, to workspace.RepoInput) (workspace.ProjectResult, error) {
			return workspace.MoveProjectMember(dependencies.Dependencies, id, repo, to)
		},
		removeMember: func(id workspace.ProjectID, repo workspace.RepoID) (workspace.DefinitionRepoResult, error) {
			return workspace.RemoveProjectMember(dependencies.Dependencies, id, repo)
		},
		add: func(input workspace.ProjectAddInput) (workspace.ProjectResult, error) {
			return workspace.AddProject(dependencies.Dependencies, input)
		},
		show: func(id workspace.ProjectID) (workspace.ProjectResult, error) {
			return workspace.ShowProject(dependencies.Dependencies, id)
		},
		list: func() (workspace.ProjectListResult, error) {
			return workspace.ListProjects(dependencies.Dependencies)
		},
	})
	addWorkspaceProjectMetadataCommands(command, dependencies)
	return command
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

	var showFormat string
	showCommand := &cobra.Command{
		Use:     "show <project-id>",
		Short:   "Show a registered project",
		Long:    "Show the stored wrapper and repository members for one project in the containing Ply workspace.",
		Example: "  ply workspace project show ply\n  ply workspace project show --format json -- ply",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.ProjectInvalidArguments(fmt.Sprintf("expected exactly one project ID, got %d arguments", len(args)))
			}
			if _, err := workspace.ParseProjectID(args[0]); err != nil {
				return err
			}
			return validateProjectFormat(showFormat)
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
			if showFormat == "json" {
				b, err := workspace.MarshalProject(result)
				if err != nil {
					return err
				}
				_, err = cmd.OutOrStdout().Write(append(b, '\n'))
				return err
			}
			return renderProject(cmd, result, false)
		},
	}
	showCommand.Flags().StringVar(&showFormat, "format", "text", "output format (text or json)")
	setProjectFlagErrors(showCommand)

	var listFormat string
	listCommand := &cobra.Command{
		Use:     "list",
		Short:   "List registered projects",
		Long:    "List projects registered in the containing Ply workspace.",
		Example: "  ply workspace project list\n  ply workspace project list --format json",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return workspace.ProjectInvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
			}
			return validateProjectFormat(listFormat)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := services.list()
			if err != nil {
				return err
			}
			if listFormat == "json" {
				b, err := workspace.MarshalProjectList(result)
				if err != nil {
					return err
				}
				_, err = cmd.OutOrStdout().Write(append(b, '\n'))
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
	listCommand.Flags().StringVar(&listFormat, "format", "text", "output format (text or json)")
	setProjectFlagErrors(listCommand)

	var memberValue string
	memberCommand := &cobra.Command{
		Use:     "add-member <project-id>",
		Short:   "Add one explicit repository member to an existing project",
		Long:    "Add one repository member as metadata. Dirty worktrees and repositories outside the wrapper are accepted. Identical Git identity retries preserve the stored locator. This requires no Epic and changes no Task, Run, history, Git refs or source.",
		Example: "  ply workspace project add-member plyapp --repo plyapp-pwa=../plyapp-pwa/main",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.ProjectInvalidArguments("expected exactly one project ID")
			}
			_, _, err := workspace.ParseProjectMemberInput(args[0], memberValue)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, repo, err := workspace.ParseProjectMemberInput(args[0], memberValue)
			if err != nil {
				return err
			}
			result, err := services.addMember(id, repo)
			if err != nil {
				return err
			}
			action := "Added repository member"
			if !result.Created {
				action = "Already registered repository member"
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s %s in project %s.\n", action, repo.RepoID, id); err != nil {
				return err
			}
			return renderProject(cmd, result, false)
		},
	}
	memberCommand.Flags().StringVarP(&memberValue, "repo", "r", "", "repository member as <repo-id>=<path>")
	_ = memberCommand.MarkFlagRequired("repo")
	setProjectFlagErrors(memberCommand)
	var moveFrom, moveTo string
	parseMove := func(args []string) (workspace.ProjectID, workspace.RepoID, workspace.RepoInput, error) {
		if len(args) != 1 {
			return "", "", workspace.RepoInput{}, workspace.ProjectInvalidArguments("expected exactly one project ID")
		}
		repo, err := workspace.ParseRepoID(moveFrom)
		if err != nil {
			return "", "", workspace.RepoInput{}, err
		}
		id, to, err := workspace.ParseProjectMemberInput(args[0], moveTo)
		return id, repo, to, err
	}
	moveCommand := &cobra.Command{
		Use:     "move-member <project-id>",
		Short:   "Give a repository member a new ID and path",
		Long:    "Give one registered repository member a new ID, path or both after its directory has moved. The origin remote is read again at the new path and the wrapper is kept; it must contain the new path. The move is refused while a Task, Epic or operation binds the member. This moves no files and changes no Git.",
		Example: "  ply workspace project move-member ply --repo pflyt-cli --to ply-planner=ply_projects/ply/ply-planner",
		Args: func(cmd *cobra.Command, args []string) error {
			_, _, _, err := parseMove(args)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, repo, to, err := parseMove(args)
			if err != nil {
				return err
			}
			result, err := services.moveMember(id, repo, to)
			if err != nil {
				return err
			}
			line := fmt.Sprintf("Moved repository member %s to %s in project %s.\n", repo, to.RepoID, id)
			if !result.Created {
				line = fmt.Sprintf("Repository member %s is already at %s in project %s.\n", to.RepoID, to.Path, id)
			}
			if _, err := io.WriteString(cmd.OutOrStdout(), line); err != nil {
				return err
			}
			return renderProject(cmd, result, false)
		},
	}
	moveCommand.Flags().StringVarP(&moveFrom, "repo", "r", "", "repository member ID to move")
	moveCommand.Flags().StringVar(&moveTo, "to", "", "new member as <repo-id>=<path>")
	_ = moveCommand.MarkFlagRequired("repo")
	_ = moveCommand.MarkFlagRequired("to")
	setProjectFlagErrors(moveCommand)
	var removeRepo string
	parseRemove := func(args []string) (workspace.ProjectID, workspace.RepoID, error) {
		if len(args) != 1 {
			return "", "", workspace.ProjectInvalidArguments("expected exactly one project ID")
		}
		id, err := workspace.ParseProjectID(args[0])
		if err != nil {
			return "", "", err
		}
		repo, err := workspace.ParseRepoID(removeRepo)
		return id, repo, err
	}
	removeCommand := &cobra.Command{
		Use:     "remove-member <project-id>",
		Short:   "Remove a repository member from a project",
		Long:    "Remove one registered repository member from the project definition. Its directory may have moved or be gone. The removal is refused while a Task, Epic or operation binds the member. This deletes no files and changes no Git; data kept under the repository ID stays.",
		Example: "  ply workspace project remove-member ply --repo ply-observer",
		Args: func(cmd *cobra.Command, args []string) error {
			_, _, err := parseRemove(args)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id, repo, err := parseRemove(args)
			if err != nil {
				return err
			}
			result, err := services.removeMember(id, repo)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Removed repository member %s (%s) from project %s in Ply workspace %s.\n", result.RepoID, result.Path, result.ProjectID, result.Workspace)
			return err
		},
	}
	removeCommand.Flags().StringVarP(&removeRepo, "repo", "r", "", "repository member ID to remove")
	_ = removeCommand.MarkFlagRequired("repo")
	setProjectFlagErrors(removeCommand)
	projectCommand.AddCommand(addCommand, memberCommand, moveCommand, removeCommand, showCommand, listCommand)
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

func validateProjectFormat(format string) error {
	if format != "text" && format != "json" {
		return workspace.ProjectInvalidArguments("--format must be text or json")
	}
	return nil
}
