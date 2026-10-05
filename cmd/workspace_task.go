package cmd

import (
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workspaceTaskServices struct {
	create         func(workspace.TaskCreateInput) (workspace.TaskMutationResult, error)
	show           func(workspace.TaskID) (workspace.TaskReadbackResult, error)
	list           func(workspace.TaskListFilters) (workspace.TaskListResult, error)
	createWorktree func(workspace.TaskWorktreeCreateInput) (workspace.TaskWorktreeMutationResult, error)
	recordResult   func(workspace.TaskResultRecordInput) (workspace.TaskResultMutationResult, error)
	recordQA       func(workspace.TaskHumanQARecordInput) (workspace.TaskHumanQAMutationResult, error)
	integrate      func(workspace.TaskIntegrationInput) (workspace.TaskIntegrationResult, error)
}

func newWorkspaceTaskCommand(dependencies workspace.Dependencies) *cobra.Command {
	command := newWorkspaceTaskCommandWithServices(workspaceTaskServices{create: func(input workspace.TaskCreateInput) (workspace.TaskMutationResult, error) {
		return workspace.CreateTask(dependencies, input)
	}, show: func(id workspace.TaskID) (workspace.TaskReadbackResult, error) {
		return workspace.ShowTask(dependencies, id)
	}, list: func(filters workspace.TaskListFilters) (workspace.TaskListResult, error) {
		return workspace.ListTasksWithFilters(dependencies, filters)
	}, createWorktree: func(input workspace.TaskWorktreeCreateInput) (workspace.TaskWorktreeMutationResult, error) {
		return workspace.CreateTaskWorktree(dependencies, input)
	}, recordResult: func(input workspace.TaskResultRecordInput) (workspace.TaskResultMutationResult, error) {
		return workspace.RecordTaskResult(dependencies, input)
	}, recordQA: func(input workspace.TaskHumanQARecordInput) (workspace.TaskHumanQAMutationResult, error) {
		return workspace.RecordTaskHumanQA(dependencies, input)
	}, integrate: func(input workspace.TaskIntegrationInput) (workspace.TaskIntegrationResult, error) {
		if input.Apply {
			return workspace.ApplyTaskIntegration(dependencies, input)
		}
		return workspace.CheckTaskIntegration(dependencies, input)
	}})
	addWorkspaceTaskContentCommands(command, dependencies)
	addWorkspaceTaskQueueCommands(command, dependencies)
	addWorkspaceTaskRunCommands(command, dependencies)
	addWorkspaceTaskJournalCommands(command, dependencies)
	return command
}

func newWorkspaceTaskCommandWithServices(services workspaceTaskServices) *cobra.Command {
	command := &cobra.Command{Use: "task", Short: "Manage Tasks in a Ply workspace", Long: "Create, inspect, and advance repository-bound Tasks owned by workspace Epics.", Example: "  ply workspace task create workspace-work-item-bootstrap --title \"Workspace-owned Epic, Task, and worktree support\" --description \"Add explicit workspace work items and prepare a Task worktree from the Epic base.\" --epic ply-agentic-workflow-support --project ply --repo ply\n  ply workspace task integrate workspace-work-item-bootstrap --result trs_0123456789abcdef0123456789abcdef --qa hqa_0123456789abcdef0123456789abcdef --expected-result-oid a55b192334cafcd8527895372205205fa43f3cc5 --expected-parent-oid 54f3631cbea789f25a4134945c7ca16d343139df --check"}
	var title, description, epicID, projectID, repoID, createFormat, upgradeStore string
	create := &cobra.Command{Use: "create <task-id>", Short: "Create a repository-bound Task and record its initial problem", Long: "Create one Task under an existing Epic and record its initial problem. A format 1 or 2 registry requires --upgrade-store with its exact digest. Solution selection, worktree creation, agent start, QA and integration remain separate actions.", Example: "  ply workspace task create workspace-work-item-bootstrap --title \"Workspace-owned Epic, Task, and worktree support\" --description \"Add explicit workspace work items and prepare a Task worktree from the Epic base.\" --epic ply-agentic-workflow-support --project ply --repo ply", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Task ID, got %d arguments", len(args)))
		}
		_, err := workspace.ParseTaskCreateInput(args[0], title, description, epicID, projectID, repoID)
		return err
	}, RunE: func(cmd *cobra.Command, args []string) error {
		input, err := workspace.ParseTaskCreateInput(args[0], title, description, epicID, projectID, repoID)
		if err != nil {
			return err
		}
		if err := validateWorkFormat(createFormat); err != nil {
			return err
		}
		if upgradeStore != "" {
			input.RegistryUpgrade = &workspace.TaskRegistryUpgrade{RegistrySHA256: upgradeStore}
		}
		result, err := services.create(input)
		if createFormat == "json" && result.Content != nil {
			b, marshalErr := workspace.MarshalTaskContentMutation(*result.Content)
			if marshalErr != nil {
				return marshalErr
			}
			if _, writeErr := cmd.OutOrStdout().Write(append(b, '\n')); writeErr != nil {
				return fmt.Errorf("%w; inspect with ply workspace task publication show %s --key %s", writeErr, input.TaskID, result.Content.PublicationKey)
			}
			return err
		}
		if err != nil {
			return err
		}
		if createFormat == "json" {
			readback, readErr := services.show(input.TaskID)
			if readErr != nil {
				return readErr
			}
			b, marshalErr := workspace.MarshalTaskReadback(readback)
			if marshalErr != nil {
				return marshalErr
			}
			_, writeErr := cmd.OutOrStdout().Write(append(b, '\n'))
			return writeErr
		}
		if writeErr := renderTaskMutation(cmd, result); writeErr != nil {
			if result.Content != nil {
				return fmt.Errorf("%w; inspect with ply workspace task publication show %s --key %s", writeErr, input.TaskID, result.Content.PublicationKey)
			}
			return writeErr
		}
		return nil
	}}
	create.Flags().StringVar(&createFormat, "format", "text", "output format (text or json)")
	create.Flags().StringVar(&upgradeStore, "upgrade-store", "", "acknowledge the exact SHA-256 digest of a format 1 or 2 registry")
	create.Flags().StringVar(&title, "title", "", "Task display title")
	create.Flags().StringVar(&description, "description", "", "one-line Task description")
	create.Flags().StringVar(&epicID, "epic", "", "parent Epic ID")
	create.Flags().StringVar(&projectID, "project", "", "registered Project ID")
	create.Flags().StringVar(&repoID, "repo", "", "registered repository ID")
	for _, name := range []string{"title", "description", "epic", "project", "repo"} {
		_ = create.MarkFlagRequired(name)
	}
	setWorkFlagErrors(create)
	var format string
	show := &cobra.Command{Use: "show <task-id>", Short: "Show a registered Task", Long: "Show persisted Task, delivery, QA, and local integration state with a fresh read-only Git observation.", Example: "  ply workspace task show workspace-work-item-bootstrap\n  ply workspace task show workspace-work-item-bootstrap --format json", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Task ID, got %d arguments", len(args)))
		}
		if _, err := workspace.ParseTaskID(args[0]); err != nil {
			return err
		}
		return validateWorkFormat(format)
	}, RunE: func(cmd *cobra.Command, args []string) error {
		id, err := workspace.ParseTaskID(args[0])
		if err != nil {
			return err
		}
		result, err := services.show(id)
		if err != nil {
			return err
		}
		if format == "json" {
			bytes, err := workspace.MarshalTaskReadback(result)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(append(bytes, '\n'))
			return err
		}
		return renderTaskShow(cmd, result)
	}}
	show.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(show)
	var filterEpic, filterProject, filterRepo, listFormat string
	filters := func(cmd *cobra.Command) workspace.TaskListFilters {
		var f workspace.TaskListFilters
		if cmd.Flags().Changed("project") {
			id := workspace.ProjectID(filterProject)
			f.ProjectID = &id
		}
		if cmd.Flags().Changed("repo") {
			id := workspace.RepoID(filterRepo)
			f.RepoID = &id
		}
		if cmd.Flags().Changed("epic") {
			id := workspace.EpicID(filterEpic)
			f.EpicID = &id
		}
		return f
	}
	list := &cobra.Command{Use: "list", Short: "List registered Tasks", Long: "List all Tasks registered in the containing Ply workspace. Project, repository and Epic filters work independently and combine as exact AND filters, without an implicit working-directory selection.", Example: "  ply workspace task list\n  ply workspace task list --project ply --repo ply\n  ply workspace task list --epic ply-agentic-workflow-support --format json", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
		}
		// Ready retains its original argument validation and separate target semantics.
		if ready, _ := cmd.Flags().GetBool("ready"); ready {
			if filterEpic != "" {
				_, err := workspace.ParseEpicID(filterEpic)
				return err
			}
			return nil
		}
		if err := validateWorkFormat(listFormat); err != nil {
			return err
		}
		return filters(cmd).Validate()
	}, RunE: func(cmd *cobra.Command, args []string) error {
		result, err := services.list(filters(cmd))
		if err != nil {
			return err
		}
		if listFormat == "json" {
			b, err := workspace.MarshalTaskList(result)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(append(b, '\n'))
			return err
		}
		return renderTaskList(cmd, result)
	}}
	list.Flags().StringVar(&filterEpic, "epic", "", "filter by parent Epic ID; with --ready, requires --project and --repo")
	list.Flags().StringVar(&filterProject, "project", "", "filter by Project ID; with --ready, requires --repo and --epic")
	list.Flags().StringVar(&filterRepo, "repo", "", "filter by repository ID; with --ready, requires --project and --epic")
	list.Flags().StringVar(&listFormat, "format", "text", "output format (text or json)")
	setWorkFlagErrors(list)
	worktree := newWorkspaceTaskWorktreeCommand(services)
	command.AddCommand(create, show, list, worktree, newWorkspaceTaskResultCommand(services), newWorkspaceTaskQACommand(services), newWorkspaceTaskIntegrateCommand(services))
	return command
}

func newWorkspaceTaskResultCommand(services workspaceTaskServices) *cobra.Command {
	parent := &cobra.Command{Use: "result", Short: "Manage controlled Task results", Long: "Record immutable delivery evidence as typed results for workspace Tasks.", Example: "  ply workspace task result record workspace-work-item-bootstrap --file /absolute/task-result.json"}
	var file, format string
	record := &cobra.Command{Use: "record <task-id>", Short: "Record a controlled Task result", Long: "Validate immutable handoff evidence and record one typed result for an exact Task worktree commit.", Example: "  ply workspace task result record workspace-work-item-bootstrap --file /absolute/task-result.json\n  ply workspace task result record workspace-work-item-bootstrap --file /absolute/task-result.json --format json", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Task ID, got %d arguments", len(args)))
		}
		if _, err := workspace.ParseTaskID(args[0]); err != nil {
			return err
		}
		return validateWorkFormat(format)
	}, RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := workspace.ParseTaskID(args[0])
		result, err := services.recordResult(workspace.TaskResultRecordInput{TaskID: id, File: file})
		if err != nil {
			return err
		}
		if format == "json" {
			b, err := workspace.MarshalTaskResultReadback(result)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(append(b, '\n'))
			return err
		}
		first := "Recorded Task result"
		if !result.Created {
			first = "Task result"
		}
		if result.Created {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s.\n", first, result.Record.ID)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s already exists with identical content.\n", first, result.Record.ID)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Task: %s / %s\nTechnical gate: %s\nHandoff evidence: %s / %s (validated)\nStore transition: %s\nNext action: %s\n", result.Record.TaskID, result.Record.ResultOID, result.Record.TechnicalGate, result.Record.HandoffID, result.Record.TerminalResultID, result.Record.StoreTransition, taskResultNextActionText(result.Record))
		return err
	}}
	record.Flags().StringVar(&file, "file", "", "Task result draft JSON file")
	record.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	_ = record.MarkFlagRequired("file")
	setWorkFlagErrors(record)
	parent.AddCommand(record)
	return parent
}

func newWorkspaceTaskQACommand(services workspaceTaskServices) *cobra.Command {
	parent := &cobra.Command{Use: "qa", Short: "Manage human QA records for Tasks", Long: "Record a human product QA outcome for an exact controlled Task result.", Example: "  ply workspace task qa record workspace-work-item-bootstrap --file /absolute/human-qa.json"}
	var file, format string
	record := &cobra.Command{Use: "record <task-id>", Short: "Record human QA for a Task result", Long: "Validate a human QA draft and record its outcome without starting integration.", Example: "  ply workspace task qa record workspace-work-item-bootstrap --file /absolute/human-qa.json\n  ply workspace task qa record workspace-work-item-bootstrap --file /absolute/human-qa.json --format json", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Task ID, got %d arguments", len(args)))
		}
		if _, err := workspace.ParseTaskID(args[0]); err != nil {
			return err
		}
		return validateWorkFormat(format)
	}, RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := workspace.ParseTaskID(args[0])
		result, err := services.recordQA(workspace.TaskHumanQARecordInput{TaskID: id, File: file})
		if err != nil {
			return err
		}
		if format == "json" {
			b, err := workspace.MarshalTaskHumanQAReadback(result)
			if err != nil {
				return err
			}
			_, err = cmd.OutOrStdout().Write(append(b, '\n'))
			return err
		}
		if result.Created {
			fmt.Fprintf(cmd.OutOrStdout(), "Recorded human QA %s.\n", result.Record.ID)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Human QA %s already exists with identical content.\n", result.Record.ID)
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Task result: %s / %s\nOutcome: %s\nHuman identity: locally claimed; not cryptographically attested\nGit changed by this command: no\nNext action: %s\n", result.Record.TaskResultID, result.Record.ResultOID, result.Record.Outcome, taskQANextActionText(result.Record))
		return err
	}}
	record.Flags().StringVar(&file, "file", "", "human QA draft JSON file")
	record.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	_ = record.MarkFlagRequired("file")
	setWorkFlagErrors(record)
	parent.AddCommand(record)
	return parent
}

func newWorkspaceTaskIntegrateCommand(services workspaceTaskServices) *cobra.Command {
	var resultID, qaID, resultOID, parentOID, confirm, retry, format string
	var check, apply bool
	command := &cobra.Command{Use: "integrate <task-id>", Short: "Integrate a Task into its Epic parent", Long: "Check or apply one confirmed local fast-forward from an exact Task result to its registered Epic parent.", Example: "  ply workspace task integrate workspace-work-item-bootstrap --result trs_0123456789abcdef0123456789abcdef --qa hqa_0123456789abcdef0123456789abcdef --expected-result-oid a55b192334cafcd8527895372205205fa43f3cc5 --expected-parent-oid 54f3631cbea789f25a4134945c7ca16d343139df --check\n  ply workspace task integrate workspace-work-item-bootstrap --result trs_0123456789abcdef0123456789abcdef --qa hqa_0123456789abcdef0123456789abcdef --expected-result-oid a55b192334cafcd8527895372205205fa43f3cc5 --expected-parent-oid 54f3631cbea789f25a4134945c7ca16d343139df --apply --confirm sha256:<64hex>", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Task ID, got %d arguments", len(args)))
		}
		if check == apply {
			return workspace.WorkInvalidArguments("exactly one of --check or --apply is required")
		}
		if check && confirm != "" {
			return workspace.WorkInvalidArguments("--confirm is forbidden with --check")
		}
		if apply && confirm == "" {
			return workspace.WorkInvalidArguments("--confirm is required with --apply")
		}
		if err := validateWorkFormat(format); err != nil {
			return err
		}
		_, err := workspace.ParseTaskIntegrationInput(args[0], resultID, qaID, resultOID, parentOID, retry, apply, confirm)
		return err
	}, RunE: func(cmd *cobra.Command, args []string) error {
		input, err := workspace.ParseTaskIntegrationInput(args[0], resultID, qaID, resultOID, parentOID, retry, apply, confirm)
		if err != nil {
			return err
		}
		out, err := services.integrate(input)
		if err != nil && out.Readback.Value == nil {
			return err
		}
		if format == "json" {
			b, marshalErr := workspace.MarshalTaskIntegrationReadback(out.Readback)
			if marshalErr != nil {
				return marshalErr
			}
			if _, writeErr := cmd.OutOrStdout().Write(append(b, '\n')); writeErr != nil {
				return fmt.Errorf("%w; inspect with ply workspace task show %s --format json", writeErr, input.TaskID)
			}
			return err
		}
		if _, writeErr := fmt.Fprint(cmd.OutOrStdout(), workspace.RenderTaskIntegrationText(out.Readback)); writeErr != nil {
			return fmt.Errorf("%w; inspect with ply workspace task show %s --format json", writeErr, input.TaskID)
		}
		return err
	}}
	command.Flags().StringVar(&resultID, "result", "", "controlled Task result ID")
	command.Flags().StringVar(&qaID, "qa", "", "human QA record ID")
	command.Flags().StringVar(&resultOID, "expected-result-oid", "", "expected exact Task result commit object ID")
	command.Flags().StringVar(&parentOID, "expected-parent-oid", "", "expected current Epic parent commit object ID")
	command.Flags().BoolVar(&check, "check", false, "check the exact integration plan without writing state")
	command.Flags().BoolVar(&apply, "apply", false, "apply one confirmed local fast-forward")
	command.Flags().StringVar(&confirm, "confirm", "", "confirmed integration plan SHA-256")
	command.Flags().StringVar(&retry, "retry-after", "", "exact prior no-effect integration result ID")
	command.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	for _, n := range []string{"result", "qa", "expected-result-oid", "expected-parent-oid"} {
		_ = command.MarkFlagRequired(n)
	}
	setWorkFlagErrors(command)
	return command
}

func taskResultNextActionText(r workspace.TaskResultRecord) string {
	if r.TechnicalGate == "passed" || r.TechnicalGate == "good_enough_with_known_debt" {
		return fmt.Sprintf("Record human QA with `ply workspace task qa record %s --file <absolute-human-qa.json>`.", r.TaskID)
	}
	return "Start a separate result control before recording human QA."
}
func taskQANextActionText(r workspace.TaskHumanQARecord) string {
	if r.Outcome == "pass" {
		return fmt.Sprintf("Run `ply workspace task integrate %s --result %s --qa %s --expected-result-oid %s --expected-parent-oid <oid> --check`.", r.TaskID, r.TaskResultID, r.ID, r.ResultOID)
	}
	return "Start a separate delivery or QA clarification before integration."
}

func newWorkspaceTaskWorktreeCommand(services workspaceTaskServices) *cobra.Command {
	command := &cobra.Command{Use: "worktree", Short: "Manage Task worktrees", Long: "Create and inspect the single local worktree binding for a Task.", Example: "  ply workspace task worktree create workspace-work-item-bootstrap --branch ply_workspace_work_item_bootstrap --path ../ply/ply_workspace_work_item_bootstrap --expected-parent-oid 54f3631cbea789f25a4134945c7ca16d343139df"}
	var branch, path, expected string
	create := &cobra.Command{Use: "create <task-id>", Short: "Create and bind a Task worktree", Long: "Persist a Task worktree intent, then add one local branch and worktree from the exact Epic base.", Example: "  ply workspace task worktree create workspace-work-item-bootstrap --branch ply_workspace_work_item_bootstrap --path ../ply/ply_workspace_work_item_bootstrap --expected-parent-oid 54f3631cbea789f25a4134945c7ca16d343139df", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Task ID, got %d arguments", len(args)))
		}
		_, err := workspace.ParseTaskWorktreeCreateInput(args[0], branch, path, expected)
		return err
	}, RunE: func(cmd *cobra.Command, args []string) error {
		input, err := workspace.ParseTaskWorktreeCreateInput(args[0], branch, path, expected)
		if err != nil {
			return err
		}
		result, err := services.createWorktree(input)
		if err != nil {
			return err
		}
		return renderTaskWorktreeMutation(cmd, result)
	}}
	create.Flags().StringVar(&branch, "branch", "", "new short local branch name")
	create.Flags().StringVar(&path, "path", "", "new Task worktree path")
	create.Flags().StringVar(&expected, "expected-parent-oid", "", "expected current Epic parent commit object ID")
	for _, name := range []string{"branch", "path", "expected-parent-oid"} {
		_ = create.MarkFlagRequired(name)
	}
	setWorkFlagErrors(create)
	command.AddCommand(create)
	return command
}

func renderTaskMutation(command *cobra.Command, result workspace.TaskMutationResult) error {
	if result.Created {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "Created Task %s: %s\n", result.Task.ID, result.Task.Title); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "Task %s already exists: %s\n", result.Task.ID, result.Task.Title); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(command.OutOrStdout(), "Description: %s\nParent Epic: %s (%s)\nProject / repository: %s / %s\nGit common directory: %s\nWorktree: not created\nState: unbound\nAgent started by this command: no\n", result.Task.Description, result.Epic.ID, result.Epic.Title, result.Task.ProjectID, result.Task.RepoID, result.Task.GitCommonDir)
	return err
}
func renderTaskShow(command *cobra.Command, result workspace.TaskReadbackResult) error {
	if _, err := fmt.Fprintf(command.OutOrStdout(), "Task %s: %s\nDescription: %s\nParent Epic: %s (%s)\nProject / repository: %s / %s\nGit common directory: %s\n", result.Task.ID, result.Task.Title, result.Task.Description, result.Epic.ID, result.Epic.Title, result.Task.ProjectID, result.Task.RepoID, result.Task.GitCommonDir); err != nil {
		return err
	}
	if result.Task.WorktreeState == workspace.WorkItemUnbound {
		if _, err := fmt.Fprint(command.OutOrStdout(), "Worktree: not created\nState: unbound\nFresh observation: not requested\nReady for handoff: no\nReason: Task worktree has not been created.\nAgent started by this command: no\n"); err != nil {
			return err
		}
	} else if err := renderTaskWorktreeBlock(command, result.Task, result.Operation, result.ReadyForHandoff, result.Reasons); err != nil {
		return err
	}
	if result.Integration != nil {
		_, err := fmt.Fprint(command.OutOrStdout(), workspace.RenderTaskIntegrationText(*result.Integration))
		return err
	}
	return nil
}
func renderTaskWorktreeBlock(command *cobra.Command, task workspace.TaskRecord, operation *workspace.WorktreeOperationRecord, ready bool, reasons []string) error {
	locator := ""
	worktreeID := ""
	branch := ""
	parentRef := ""
	parentOID := ""
	parentTree := ""
	operationID := ""
	digest := ""
	if task.Worktree != nil {
		locator = task.Worktree.Locator
		worktreeID = string(task.Worktree.ID)
		branch = task.Worktree.Ref
		parentRef = task.Worktree.ParentRef
		parentOID = task.Worktree.ParentOID
		parentTree = task.Worktree.ParentTree
		operationID = string(task.Worktree.OperationID)
		digest = task.Worktree.IntentDigest
	} else if operation != nil {
		locator = operation.TargetLocator
		worktreeID = string(operation.WorktreeID)
		branch = operation.SourceRef
		parentRef = operation.ParentRef
		parentOID = operation.ParentOID
		parentTree = operation.ParentTree
		operationID = string(operation.ID)
		digest = operation.IntentDigest
	}
	summary := "clean; source and parent match the recorded start base"
	if !ready {
		summary = "stale; one or more current fields differ from the recorded binding"
		for _, reason := range reasons {
			if reason == "parent_observation_unknown" || reason == "target_observation_unknown" || reason == "worktree_inventory_unknown" || reason == "project_binding_unknown" {
				summary = "unknown; one or more current fields could not be observed"
				break
			}
		}
	}
	if _, err := fmt.Fprintf(command.OutOrStdout(), "Worktree: %s\nWorktree ID: %s\nBranch: %s\nParent: %s@%s\nStart tree: %s\nOperation: %s (%s)\nState: %s\nFresh observation: %s\nReady for handoff: %s\n", locator, worktreeID, branch, parentRef, parentOID, parentTree, operationID, digest, task.WorktreeState, summary, yesNo(ready)); err != nil {
		return err
	}
	for _, reason := range reasons {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "Reason: %s\n", workReasonText(reason)); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(command.OutOrStdout(), "Agent started by this command: no\n")
	return err
}
func renderTaskList(command *cobra.Command, result workspace.TaskListResult) error {
	if result.Filters.ProjectID != nil || result.Filters.RepoID != nil {
		labels := []string{}
		if result.Filters.ProjectID != nil {
			labels = append(labels, "Project "+string(*result.Filters.ProjectID))
		}
		if result.Filters.RepoID != nil {
			labels = append(labels, "repository "+string(*result.Filters.RepoID))
		}
		if result.Filters.EpicID != nil {
			labels = append(labels, "Epic "+string(*result.Filters.EpicID))
		}
		scope := strings.Join(labels, ", ")
		if len(result.Tasks) == 0 {
			_, err := fmt.Fprintf(command.OutOrStdout(), "No Tasks for %s are registered in Ply workspace %s.\n", scope, result.Workspace)
			return err
		}
		if _, err := fmt.Fprintf(command.OutOrStdout(), "Tasks for %s in Ply workspace %s:\n", scope, result.Workspace); err != nil {
			return err
		}
		for _, task := range result.Tasks {
			if _, err := fmt.Fprintf(command.OutOrStdout(), "  %s: %s (Epic %s; %s / %s; %s)\n", task.ID, currentTaskTitle(result, task), task.ParentEpicID, task.ProjectID, task.RepoID, task.WorktreeState); err != nil {
				return err
			}
		}
		return nil
	}
	if len(result.Tasks) == 0 {
		if result.Epic != nil {
			_, err := fmt.Fprintf(command.OutOrStdout(), "No Tasks for Epic %s are registered in Ply workspace %s.\n", result.Epic.ID, result.Workspace)
			return err
		}
		_, err := fmt.Fprintf(command.OutOrStdout(), "No Tasks are registered in Ply workspace %s.\n", result.Workspace)
		return err
	}
	if result.Epic != nil {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "Tasks for Epic %s in Ply workspace %s:\n", result.Epic.ID, result.Workspace); err != nil {
			return err
		}
		for _, task := range result.Tasks {
			if _, err := fmt.Fprintf(command.OutOrStdout(), "  %s: %s (%s / %s; %s)\n", task.ID, currentTaskTitle(result, task), task.ProjectID, task.RepoID, task.WorktreeState); err != nil {
				return err
			}
		}
		return nil
	}
	if _, err := fmt.Fprintf(command.OutOrStdout(), "Tasks in Ply workspace %s:\n", result.Workspace); err != nil {
		return err
	}
	for _, task := range result.Tasks {
		if _, err := fmt.Fprintf(command.OutOrStdout(), "  %s: %s (Epic %s; %s / %s; %s)\n", task.ID, currentTaskTitle(result, task), task.ParentEpicID, task.ProjectID, task.RepoID, task.WorktreeState); err != nil {
			return err
		}
	}
	return nil
}
func renderTaskWorktreeMutation(command *cobra.Command, result workspace.TaskWorktreeMutationResult) error {
	line := "Created Task worktree for %s.\n"
	if result.Outcome == "recovered" {
		line = "Recovered Task worktree for %s.\n"
	} else if result.Outcome == "already_ready" {
		line = "Task worktree for %s is already ready.\n"
	}
	if _, err := fmt.Fprintf(command.OutOrStdout(), line, result.Task.ID); err != nil {
		return err
	}
	return renderTaskWorktreeBlock(command, result.Task, &result.Operation, true, nil)
}

func workReasonText(reason string) string {
	values := map[string]string{"task_worktree_unbound": "Task worktree has not been created.", "task_worktree_creation_in_progress": "Task worktree creation has not reached a verified terminal state.", "task_worktree_reconciliation_required": "Task worktree state requires an explicit reconciliation activity.", "project_binding_stale": "The current Project repository binding differs from the recorded binding.", "project_binding_unknown": "The current Project repository binding could not be observed.", "epic_ref_stale": "The Epic branch no longer matches the recorded commit and tree.", "epic_worktree_stale": "The Epic worktree no longer matches the recorded binding.", "epic_worktree_dirty": "The Epic worktree is not clean.", "epic_observation_unknown": "The current Epic worktree state could not be observed completely.", "parent_ref_stale": "The parent branch no longer matches the recorded start base.", "parent_worktree_stale": "The parent worktree no longer matches the recorded binding.", "parent_worktree_dirty": "The parent worktree is not clean.", "parent_observation_unknown": "The current parent worktree state could not be observed completely.", "source_ref_stale": "The Task branch no longer matches the recorded start base.", "source_ref_checked_out_elsewhere": "The Task branch is checked out in a different worktree.", "target_worktree_stale": "The Task worktree no longer matches the recorded binding.", "target_worktree_dirty": "The Task worktree is not clean.", "target_observation_unknown": "The current Task worktree state could not be observed completely.", "worktree_inventory_stale": "Git worktree inventory differs from the recorded binding.", "worktree_inventory_unknown": "Git worktree inventory could not be observed completely."}
	if value, ok := values[reason]; ok {
		return value
	}
	return reason
}

func currentTaskTitle(result workspace.TaskListResult, task workspace.TaskRecord) string {
	if title, ok := result.CurrentTitles[task.ID]; ok {
		return title
	}
	return task.Title
}
