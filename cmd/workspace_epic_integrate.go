package cmd

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceEpicIntegrateCommand(d taskwork.Dependencies) *cobra.Command {
	var actor string
	var keepWorktree bool
	var c *cobra.Command
	c = newWorkspaceEpicIntegrateCommandWith(func(input taskwork.EpicIntegrateInput) (taskwork.EpicIntegrateResult, error) {
		result, err := taskwork.IntegrateEpic(d, input)
		if err == nil && result.Applied && result.Effect == "fast_forwarded" {
			err = recordEpicIntegration(d, input, result, namedActor(actor))
			if err == nil && !keepWorktree {
				retireEpicAfterIntegration(c.ErrOrStderr(), d, input, plyActor(actor))
			}
		}
		return result, err
	})
	c.Long += " After a fast-forward with --apply the Epic worktree and branch are retired and the Epic archived as with epic retire, unless --keep-worktree is given; when that cannot run, stderr gives the retire command."
	c.Flags().StringVar(&actor, "actor", "", "who integrates, recorded in the journals of the Tasks it brings to the target (default: PLY_ACTOR, else the owner of each Task's plan)")
	c.Flags().BoolVar(&keepWorktree, "keep-worktree", false, "keep the Epic worktree and branch after --apply instead of retiring them")
	return c
}

func newWorkspaceEpicIntegrateCommandWith(integrate func(taskwork.EpicIntegrateInput) (taskwork.EpicIntegrateResult, error)) *cobra.Command {
	var repo, targetRef, targetOID, epicOID, format string
	var apply bool
	c := &cobra.Command{Use: "integrate <epic-id>", Short: "Fast-forward a target branch such as main to the Epic", Long: "Check, or with --apply perform, one local fast-forward of --target-ref to the exact Epic commit. Both commits must match the expected object IDs. The target must not have commits outside the Epic; Ply never merges or rebases. A target checked out in a worktree must be clean and moves with git merge --ff-only; otherwise only the ref moves. The reflog action is ply-workspace-epic-integrate:<epic-id>. Nothing is pushed.", Example: "  ply workspace epic integrate code-map-d1 --repo ply-code-map --target-ref refs/heads/main --expected-target-oid <main-oid> --expected-epic-oid <epic-oid>\n  ply workspace epic integrate code-map-d1 --repo ply-code-map --target-ref refs/heads/main --expected-target-oid <main-oid> --expected-epic-oid <epic-oid> --apply", Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments(fmt.Sprintf("expected exactly one Epic ID, got %d arguments", len(args)))
		}
		return validateWorkFormat(format)
	}, RunE: func(cmd *cobra.Command, args []string) error {
		result, err := integrate(taskwork.EpicIntegrateInput{EpicID: taskwork.EpicID(args[0]), RepoID: workspace.RepoID(repo), TargetRef: targetRef, ExpectedTargetOID: targetOID, ExpectedEpicOID: epicOID, Apply: apply})
		if err != nil {
			return err
		}
		if format == "json" {
			return writeAssignmentJSON(cmd.OutOrStdout(), struct {
				Kind          string `json:"kind"`
				SchemaVersion int    `json:"schema_version"`
				taskwork.EpicIntegrateResult
			}{"WorkspaceEpicIntegrate@1", 1, result})
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Epic %s %s at %s\nTarget %s at %s\nEffect: %s\n", result.EpicID, result.EpicRef, result.EpicOID, result.TargetRef, result.TargetOID, result.Effect)
		return err
	}}
	c.Flags().StringVar(&repo, "repo", "", "registered repository ID of the Epic binding")
	c.Flags().StringVar(&targetRef, "target-ref", "", "full local branch ref to fast-forward, such as refs/heads/main")
	c.Flags().StringVar(&targetOID, "expected-target-oid", "", "exact current target commit")
	c.Flags().StringVar(&epicOID, "expected-epic-oid", "", "exact Epic commit to integrate")
	c.Flags().BoolVar(&apply, "apply", false, "perform the fast-forward (default: check only)")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	for _, name := range []string{"repo", "target-ref", "expected-target-oid", "expected-epic-oid"} {
		_ = c.MarkFlagRequired(name)
	}
	setWorkFlagErrors(c)
	return c
}
