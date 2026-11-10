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
	c := &cobra.Command{Use: "register TASK_ID", Short: "Preview or publish a supplied first Goal and append its exact queue reference", Long: "Supply an explicit title, objective, observable criteria, registered target; supplied UTF-8 documents are optional. Preview is read-only. Apply with its confirmation publishes the first Goal and appends its exact queue reference. A new Task gets an initial Problem; an existing unused Task retains its current Problem and identity. Unrelated queue entries are preserved. Identical retries recover native publications; changed real inputs conflict. Registration never starts an agent or creates a worktree. Choose subsequent operations explicitly.", Example: "  ply workspace task goal register explain-errors --spec explain-errors-v1 --title 'Explain errors' --objective 'Make failures actionable' --criterion 'Invalid input names the missing choice' --project ply --repo ply --epic workflow --format json\n  # Repeat the same flags with --apply --confirmation <preview confirmation>", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(format); e != nil {
			return e
		}
		task = args[0]
		in.TaskID = workspace.TaskID(task)
		in.Target = workspace.QueueTargetInput{ProjectID: workspace.ProjectID(project), RepoID: workspace.RepoID(repo), EpicID: workspace.EpicID(epic)}
		if c.Flags().Changed("document") {
			in.Documents = []workspace.GoalDocumentInput{{ID: "spec", Locator: document}}
		}
		for _, doc := range additional {
			id, path, ok := strings.Cut(doc, "=")
			if !ok {
				return workspace.WorkInvalidArguments("--extra-document requires id=/absolute/path")
			}
			in.Documents = append(in.Documents, workspace.GoalDocumentInput{ID: id, Locator: path})
		}
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
			fmt.Fprintf(c.OutOrStdout(), "Task: %s; Goal: %s@%d %s\nPhase: %s; registered: %t; queue revision: %d\nReturn: %s at %s\nConfirmation: %s\nNo execution started.\n", out.TaskID, out.Goal.SpecID, out.Goal.Spec.Revision, out.Goal.Spec.ManifestSHA256, out.Phase, out.Applied, out.Queue.Revision, out.Target.ParentRef, out.Target.ParentLocator, out.Confirmation)
			if out.Problem != nil {
				fmt.Fprintf(c.OutOrStdout(), "Problem: revision %d %s\n", out.Problem.Revision, out.Problem.ManifestSHA256)
			}
			for _, src := range out.Sources {
				fmt.Fprintf(c.OutOrStdout(), "Source: %s %s (%s, %d bytes)\n", src.ID, src.Locator, src.SHA256, src.SizeBytes)
			}
		}
		return err
	}}
	c.Flags().StringVar(&in.SpecID, "spec", "", "explicit fresh Spec ID")
	c.Flags().StringVar(&document, "document", "", "optional primary supplied UTF-8 Spec document (absolute physical path)")
	c.Flags().StringArrayVar(&additional, "extra-document", nil, "additional supplied design input as id=/absolute/path (repeatable)")
	c.Flags().StringVar(&in.Title, "title", "", "supplied product title")
	c.Flags().StringVar(&in.Objective, "objective", "", "supplied product objective")
	c.Flags().StringArrayVar(&in.Criteria, "criterion", nil, "observable acceptance criterion (repeatable; no Markdown inference)")
	c.Flags().StringArrayVar(&in.Constraints, "constraint", nil, "supplied constraint (repeatable)")
	c.Flags().StringVar(&project, "project", "", "explicit registered Project")
	c.Flags().StringVar(&repo, "repo", "", "explicit registered repository")
	c.Flags().StringVar(&epic, "epic", "", "explicit registered Epic")
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
