package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceTaskGoalCommand(d workspace.Dependencies) *cobra.Command {
	goal := &cobra.Command{Use: "goal", Short: "Register supplied first Goals without constructing native envelopes"}
	var in workspace.FirstGoalInput
	var project, repo, epic, task, document, format, confirmation, provider, model, effort string
	var additional []string
	var apply bool
	c := &cobra.Command{Use: "register TASK_ID", Short: "Preview or publish a supplied first Goal and append its exact queue reference", Long: "Supply a UTF-8 Spec document, explicit title, objective, observable criteria, registered target and delivery policy. Preview is read-only. Apply with its confirmation publishes Task/Problem/Goal and appends the exact Goal, preserving unrelated queue entries. Identical retries recover native publications; changed real inputs conflict. Registration never starts an agent, creates a worktree or claims human acceptance. Run the returned next_argv from check_cwd for separate read-only execution readiness; registration success is independent of readiness.", Example: "  ply workspace task goal register explain-errors --spec explain-errors-v1 --document /absolute/spec.md --title 'Explain errors' --objective 'Make failures actionable' --criterion 'Invalid input names the missing choice' --project ply --repo ply --epic workflow --delivery local_epic_integration --acceptance human --format json\n  # Repeat the same flags with --apply --confirmation <preview confirmation>", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(format); e != nil {
			return e
		}
		task = args[0]
		in.TaskID = workspace.TaskID(task)
		in.Target = workspace.QueueTargetInput{ProjectID: workspace.ProjectID(project), RepoID: workspace.RepoID(repo), EpicID: workspace.EpicID(epic)}
		in.Documents = []workspace.GoalDocumentInput{{ID: "spec", Locator: document}}
		for _, doc := range additional {
			id, path, ok := strings.Cut(doc, "=")
			if !ok {
				return workspace.WorkInvalidArguments("--extra-document requires id=/absolute/path")
			}
			in.Documents = append(in.Documents, workspace.GoalDocumentInput{ID: id, Locator: path})
		}
		// Omitted executor choices retain native execute-time provider defaults.
		if provider != "" {
			in.Executor.Provider = &provider
		}
		if model != "" {
			in.Executor.Model = &model
		}
		if effort != "" {
			in.Executor.Effort = &effort
		}
		out, err := workspace.RegisterFirstGoal(d, in, apply, confirmation)
		if format == "json" {
			if e := json.NewEncoder(c.OutOrStdout()).Encode(out); e != nil {
				return e
			}
		} else if out.Kind != "" {
			fmt.Fprintf(c.OutOrStdout(), "Task: %s; Goal: %s@%d %s\nPhase: %s; registered: %t; queue revision: %d\nReturn: %s at %s\nDelivery: %s; acceptance: %s\nConfirmation: %s\nNo execution started; readiness: %s.\nFrom %s, run: %s\n", out.TaskID, out.Goal.SpecID, out.Goal.Spec.Revision, out.Goal.Spec.ManifestSHA256, out.Phase, out.Applied, out.Queue.Revision, out.Target.ParentRef, out.Target.ParentLocator, out.Delivery.Mode, in.Delivery.Acceptance, out.Confirmation, out.Readiness, out.CheckCwd, strings.Join(out.NextArgv, " "))
			for _, src := range out.Sources {
				fmt.Fprintf(c.OutOrStdout(), "Source: %s %s (%s, %d bytes)\n", src.ID, src.Locator, src.SHA256, src.SizeBytes)
			}
		}
		return err
	}}
	c.Flags().StringVar(&in.SpecID, "spec", "", "explicit fresh Spec ID")
	c.Flags().StringVar(&document, "document", "", "primary supplied UTF-8 Spec document (absolute physical path)")
	c.Flags().StringArrayVar(&additional, "extra-document", nil, "additional supplied design input as id=/absolute/path (repeatable)")
	c.Flags().StringVar(&in.Title, "title", "", "supplied product title")
	c.Flags().StringVar(&in.Objective, "objective", "", "supplied product objective")
	c.Flags().StringArrayVar(&in.Criteria, "criterion", nil, "observable acceptance criterion (repeatable; no Markdown inference)")
	c.Flags().StringArrayVar(&in.Constraints, "constraint", nil, "supplied constraint (repeatable)")
	c.Flags().StringVar(&project, "project", "", "explicit registered Project")
	c.Flags().StringVar(&repo, "repo", "", "explicit registered repository")
	c.Flags().StringVar(&epic, "epic", "", "explicit registered Epic")
	c.Flags().StringVar(&in.Delivery.Mode, "delivery", "", "delivery mode: local_epic_integration, local_branch_integration or pull_request")
	c.Flags().StringVar(&in.Delivery.Acceptance, "acceptance", "", "explicit human or automatic acceptance policy")
	c.Flags().StringVar(&in.Delivery.ResponsibleActor, "responsible-actor", "", "responsible actor for automatic acceptance")
	c.Flags().IntVar(&in.Delivery.TimeoutSeconds, "acceptance-timeout", 0, "explicit automatic acceptance timeout in seconds")
	c.Flags().BoolVar(&in.Delivery.RequireHumanQA, "require-human-qa", false, "also require actual human QA under automatic acceptance")
	c.Flags().StringVar(&in.Delivery.GitHubRepository, "github-repository", "", "exact owner/repo for pull_request")
	c.Flags().StringVar(&in.Delivery.Remote, "remote", "", "selected remote for pull_request")
	c.Flags().StringVar(&in.Delivery.TargetRef, "target-ref", "", "exact full target ref for pull_request")
	c.Flags().StringVar(&provider, "provider", "", "optional native codex or claude assignment")
	c.Flags().StringVar(&model, "model", "", "optional native model override")
	c.Flags().StringVar(&effort, "effort", "", "optional native effort override")
	c.Flags().BoolVar(&apply, "apply", false, "publish and append only after exact preview confirmation")
	c.Flags().StringVar(&confirmation, "confirmation", "", "exact preview confirmation")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	goal.AddCommand(c)
	return goal
}
