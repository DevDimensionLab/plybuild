package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workspaceEpicServices struct {
	adopt func(workspace.EpicAdoptInput) (workspace.EpicMutationResult, error)
	show  func(workspace.EpicID) (workspace.EpicReadbackResult, error)
	list  func() (workspace.EpicListResult, error)
}

func newWorkspaceEpicCommand(dependencies workspace.Dependencies) *cobra.Command {
	command := newWorkspaceEpicCommandWithServices(workspaceEpicServices{adopt: func(input workspace.EpicAdoptInput) (workspace.EpicMutationResult, error) {
		return workspace.AdoptEpic(dependencies, input)
	}, show: func(id workspace.EpicID) (workspace.EpicReadbackResult, error) {
		return workspace.ShowEpic(dependencies, id)
	}, list: func() (workspace.EpicListResult, error) { return workspace.ListEpics(dependencies) }})
	addWorkspaceEpicBaseCommands(command, dependencies)
	addWorkspaceEpicReadCommands(command, dependencies)
	command.AddCommand(newWorkspaceLifecycleCommand(dependencies, "epic"))
	return command
}

func newWorkspaceEpicCommandWithServices(services workspaceEpicServices) *cobra.Command {
	command := &cobra.Command{Use: "epic", Short: "Manage Epics in a Ply workspace", Long: "Adopt and inspect Epics owned by the containing Ply workspace.", Example: "  ply workspace epic adopt ply-agentic-workflow-support --title \"Ply agentic workflow support\" --project ply --repo ply --worktree ../ply/ply_agentic_workflow_support --ref refs/heads/ply_agentic_workflow_support --expected-oid 54f3631cbea789f25a4134945c7ca16d343139df\n  ply workspace epic list"}
	var title, projectID, repoID, worktree, ref, expectedOID string
	adopt := &cobra.Command{Use: "adopt <epic-id>", Short: "Adopt an existing worktree as an Epic", Long: "Register one existing clean local worktree as the repository anchor for an Epic without changing Git.", Example: "  ply workspace epic adopt ply-agentic-workflow-support --title \"Ply agentic workflow support\" --project ply --repo ply --worktree ../ply/ply_agentic_workflow_support --ref refs/heads/ply_agentic_workflow_support --expected-oid 54f3631cbea789f25a4134945c7ca16d343139df", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Epic ID, got %d arguments", len(args)))
		}
		_, err := workspace.ParseEpicAdoptInput(args[0], title, projectID, repoID, worktree, ref, expectedOID)
		return err
	}, RunE: func(cmd *cobra.Command, args []string) error {
		input, err := workspace.ParseEpicAdoptInput(args[0], title, projectID, repoID, worktree, ref, expectedOID)
		if err != nil {
			return err
		}
		result, err := services.adopt(input)
		if err != nil {
			return err
		}
		return renderEpicMutation(cmd, result)
	}}
	adopt.Flags().StringVar(&title, "title", "", "Epic display title")
	adopt.Flags().StringVar(&projectID, "project", "", "registered Project ID")
	adopt.Flags().StringVar(&repoID, "repo", "", "registered repository ID")
	adopt.Flags().StringVar(&worktree, "worktree", "", "existing Epic worktree path")
	adopt.Flags().StringVar(&ref, "ref", "", "exact full local Epic branch ref")
	adopt.Flags().StringVar(&expectedOID, "expected-oid", "", "expected current Epic commit object ID")
	for _, name := range []string{"title", "project", "repo", "worktree", "ref", "expected-oid"} {
		_ = adopt.MarkFlagRequired(name)
	}
	setWorkFlagErrors(adopt)
	var format string
	show := &cobra.Command{Use: "show <epic-id>", Short: "Show a registered Epic", Long: "Show persisted Epic identity and a fresh read-only observation of its adopted worktree.", Example: "  ply workspace epic show ply-agentic-workflow-support\n  ply workspace epic show ply-agentic-workflow-support --format json", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Epic ID, got %d arguments", len(args)))
		}
		if _, err := workspace.ParseEpicID(args[0]); err != nil {
			return err
		}
		return validateWorkFormat(format)
	}, RunE: func(cmd *cobra.Command, args []string) error {
		id, err := workspace.ParseEpicID(args[0])
		if err != nil {
			return err
		}
		result, err := services.show(id)
		if err != nil {
			return err
		}
		if format == "json" {
			bytes, err := workspace.MarshalEpicReadback(result)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(append(bytes, '\n'))
			return err
		}
		return renderEpicShow(cmd, result)
	}}
	show.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(show)
	list := &cobra.Command{Use: "list", Short: "List registered Epics", Long: "List Epics registered in the containing Ply workspace.", Example: "  ply workspace epic list", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
		}
		return nil
	}, RunE: func(cmd *cobra.Command, args []string) error {
		result, err := services.list()
		if err != nil {
			return err
		}
		if len(result.Epics) == 0 {
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "No Epics are registered in Ply workspace %s.\n", result.Workspace)
			return err
		}
		if _, err = fmt.Fprintf(cmd.OutOrStdout(), "Epics in Ply workspace %s:\n", result.Workspace); err != nil {
			return err
		}
		for _, epic := range result.Epics {
			repos := make([]string, len(epic.RepoBindings))
			for index, binding := range epic.RepoBindings {
				repos[index] = string(binding.RepoID)
			}
			sort.Strings(repos)
			if _, err = fmt.Fprintf(cmd.OutOrStdout(), "  %s: %s (Project %s; repositories: %s)\n", epic.ID, epic.Title, epic.ProjectID, strings.Join(repos, ", ")); err != nil {
				return err
			}
		}
		return nil
	}}
	setWorkFlagErrors(list)
	command.AddCommand(adopt, show, list)
	return command
}

func renderEpicMutation(command *cobra.Command, result workspace.EpicMutationResult) error {
	prefix := "Adopted Epic"
	if !result.Created {
		prefix = "Epic"
	}
	if result.Created {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "%s %s: %s\n", prefix, result.Epic.ID, result.Epic.Title); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "Epic %s is already adopted: %s\n", result.Epic.ID, result.Epic.Title); err != nil {
			return err
		}
	}
	return renderEpicBinding(command, result.Epic)
}
func renderEpicBinding(command *cobra.Command, epic workspace.EpicRecord) error {
	binding := epic.RepoBindings[0]
	_, err := fmt.Fprintf(command.OutOrStdout(), "Project / repository: %s / %s\nGit common directory: %s\nWorktree: %s\nWorktree ID: %s\nBranch: %s\nRecorded commit: %s\nRecorded tree: %s\nOrigin: adopted\nGit changed by this command: no\n", epic.ProjectID, binding.RepoID, binding.GitCommonDir, binding.Worktree.Locator, binding.Worktree.ID, binding.Worktree.Ref, binding.Worktree.OID, binding.Worktree.Tree)
	return err
}
func renderEpicShow(command *cobra.Command, result workspace.EpicReadbackResult) error {
	if _, err := fmt.Fprintf(command.OutOrStdout(), "Epic %s: %s\n", result.Epic.ID, result.Epic.Title); err != nil {
		return err
	}
	if err := renderEpicBindingWithoutEffect(command, result.Epic); err != nil {
		return err
	}
	summary := "clean; ref, commit, tree, common directory, and worktree inventory match the recorded binding"
	if len(result.Reasons) > 0 {
		summary = "stale; one or more current fields differ from the recorded binding"
		for _, fresh := range result.Freshness {
			if fresh.Locator == "unknown" || fresh.Ref == "unknown" || fresh.OID == "unknown" || fresh.Tree == "unknown" || fresh.GitCommonDir == "unknown" || fresh.Clean == "unknown" || fresh.InventoryMatch == "unknown" {
				summary = "unknown; one or more current fields could not be observed"
				break
			}
		}
	}
	if _, err := fmt.Fprintf(command.OutOrStdout(), "Fresh observation: %s\nReady for Task worktree creation: %s\n", summary, yesNo(result.Ready)); err != nil {
		return err
	}
	for _, reason := range result.Reasons {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "Reason: %s\n", workReasonText(reason)); err != nil {
			return err
		}
	}
	return nil
}
func renderEpicBindingWithoutEffect(command *cobra.Command, epic workspace.EpicRecord) error {
	binding := epic.RepoBindings[0]
	_, err := fmt.Fprintf(command.OutOrStdout(), "Project / repository: %s / %s\nGit common directory: %s\nWorktree: %s\nWorktree ID: %s\nBranch: %s\nRecorded commit: %s\nRecorded tree: %s\nOrigin: adopted\n", epic.ProjectID, binding.RepoID, binding.GitCommonDir, binding.Worktree.Locator, binding.Worktree.ID, binding.Worktree.Ref, binding.Worktree.OID, binding.Worktree.Tree)
	return err
}
func validateWorkFormat(value string) error {
	if value != "text" && value != "json" {
		return workspace.WorkInvalidArguments(fmt.Sprintf("unsupported output format %q (expected text or json)", value))
	}
	return nil
}
func setWorkFlagErrors(command *cobra.Command) {
	command.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error { return workspace.WorkInvalidArguments(err.Error()) })
}
func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
