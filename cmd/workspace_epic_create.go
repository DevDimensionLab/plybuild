package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceEpicCreateCommand(d workspace.Dependencies) *cobra.Command {
	var in workspace.EpicCreateInput
	var project, repo, format string
	c := &cobra.Command{Use: "create <epic-id>", Short: "Create one Epic branch and worktree from an explicit ordinary parent", Long: "Create one repo-specific Epic from an explicitly selected clean ordinary parent checkout. Persist intent before additive Git effects. Exact retries recover no-effect, branch-only or exact effects; unknown effects require inspection. With an active ply-workspace.json the Epic definition is added there before the runtime Epic. Output reports the observed worktree and config effects (not_attempted, no_effect, partial, completed or unknown). This creates no Task, commit, agent or product QA. Use adopt for an existing worktree.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(format); e != nil {
			return e
		}
		id, e := workspace.ParseEpicID(args[0])
		if e != nil {
			return e
		}
		in.EpicID = id
		in.ProjectID = workspace.ProjectID(project)
		in.RepoID = workspace.RepoID(repo)
		result, e := workspace.CreateEpic(d, in)
		if result.Outcome == "" && e != nil {
			result.Outcome = "rejected"
		}
		if format == "json" {
			if write := json.NewEncoder(c.OutOrStdout()).Encode(result); write != nil {
				return write
			}
		} else {
			text := fmt.Sprintf("Epic creation %s: %s\nEffects: worktree %s, config %s\n", id, result.Outcome, result.Effects.Worktree, result.Effects.Config)
			if result.Operation.ID != "" {
				text += fmt.Sprintf("Operation: %s\nWorktree: %s\n%s\n", result.Operation.ID, result.Operation.Request.Path, result.NextAction)
			}
			if _, write := fmt.Fprint(c.OutOrStdout(), text); write != nil {
				return write
			}
		}
		return e
	}}
	c.Flags().StringVar(&in.Title, "title", "", "Epic display title")
	c.Flags().StringVar(&project, "project", "", "registered Project ID")
	c.Flags().StringVar(&repo, "repo", "", "registered repository ID")
	c.Flags().StringVar(&in.ParentRef, "parent-ref", "", "explicit full local ordinary parent ref")
	c.Flags().StringVar(&in.ExpectedParentOID, "expected-parent-oid", "", "exact expected parent commit")
	c.Flags().StringVar(&in.Branch, "branch", "", "new short local branch name")
	c.Flags().StringVar(&in.Path, "path", "", "new worktree leaf under an existing physical parent")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	for _, n := range []string{"title", "project", "repo", "parent-ref", "expected-parent-oid", "branch", "path"} {
		_ = c.MarkFlagRequired(n)
	}
	setWorkFlagErrors(c)
	return c
}
