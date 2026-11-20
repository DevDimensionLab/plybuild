package cmd

import (
	"fmt"
	"io"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceEpicRetireCommand(d taskwork.Dependencies) *cobra.Command {
	var repo, targetRef, actor, format string
	var apply bool
	c := &cobra.Command{
		Use:   "retire <epic-id> --repo <repo-id> --target-ref <ref>",
		Short: "Remove an integrated Epic worktree and branch and archive the Epic",
		Long: "Check, or with --apply perform, the cleanup of one Epic repository binding: remove its clean worktree without force, delete its branch while it still points at the observed commit, and archive the Epic lifecycle.\n" +
			"The branch commit must be contained in --target-ref, such as refs/heads/main, and no Task of the Epic may still have a worktree. Uncommitted or untracked files stop the cleanup. The repository's main worktree is never removed. The registry keeps the binding as history. A rerun finishes what an interrupted run left.",
		Example: "  ply workspace epic retire code-map-d1 --repo ply-code-map --target-ref refs/heads/main\n  ply workspace epic retire code-map-d1 --repo ply-code-map --target-ref refs/heads/main --apply",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateWorkFormat(format); err != nil {
				return err
			}
			result, err := taskwork.RetireEpic(d, taskwork.EpicRetireInput{EpicID: taskwork.EpicID(args[0]), RepoID: workspace.RepoID(repo), TargetRef: targetRef, Actor: plyActor(actor), Apply: apply})
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string `json:"kind"`
					SchemaVersion int    `json:"schema_version"`
					taskwork.EpicRetireResult
				}{"WorkspaceEpicRetire@1", 1, result})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Epic %s worktree %s at %s, branch %s %s\nLifecycle: %s\n", result.EpicID, result.Worktree, result.Locator, result.Ref, result.Branch, result.Lifecycle)
			return err
		},
	}
	c.Flags().StringVar(&repo, "repo", "", "registered repository ID of the Epic binding (required)")
	c.Flags().StringVar(&targetRef, "target-ref", "", "full ref that must contain the Epic, such as refs/heads/main (required)")
	c.Flags().StringVar(&actor, "actor", "", "who retires, recorded in the Epic lifecycle (default: PLY_ACTOR)")
	c.Flags().BoolVar(&apply, "apply", false, "perform the cleanup (default: check only)")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	for _, name := range []string{"repo", "target-ref"} {
		_ = c.MarkFlagRequired(name)
	}
	setWorkFlagErrors(c)
	return c
}

// retireEpicAfterIntegration retires the Epic binding after its target was
// fast-forwarded. The integration stands either way; when cleanup cannot
// run, out says why and how to finish it later.
func retireEpicAfterIntegration(out io.Writer, d taskwork.Dependencies, input taskwork.EpicIntegrateInput, actor string) {
	result, err := taskwork.RetireEpic(d, taskwork.EpicRetireInput{EpicID: input.EpicID, RepoID: input.RepoID, TargetRef: input.TargetRef, Actor: actor, Apply: true})
	if err != nil {
		fmt.Fprintf(out, "Epic worktree kept: %v\nRetire it when it is safe: ply workspace epic retire %s --repo %s --target-ref %s --apply\n", err, input.EpicID, input.RepoID, input.TargetRef)
		return
	}
	fmt.Fprintf(out, "Epic worktree %s, branch %s, lifecycle %s\n", result.Worktree, result.Branch, result.Lifecycle)
}
