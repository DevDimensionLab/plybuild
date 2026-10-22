package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newTaskAdoptionCommand(s workspaceTaskServices, release bool) *cobra.Command {
	var in workspace.TaskAdoptionInput
	var project, repo, epic, operation, reason, confirm, format string
	var apply bool
	use, short := "adopt <task-id>", "Adopt an existing useful clean Task source without Git effects"
	if release {
		use = "release-adoption <adoption-id>"
		short = "Release one proven-unstarted adoption's metadata only"
	}
	c := &cobra.Command{Use: use, Short: short, Args: func(c *cobra.Command, args []string) error {
		if len(args) != 1 {
			return workspace.WorkInvalidArguments("exactly one Task ID required")
		}
		if release {
			operation = args[0]
		} else {
			id, e := workspace.ParseTaskID(args[0])
			if e != nil {
				return e
			}
			in.TaskID = id
		}
		if apply && confirm == "" || !apply && confirm != "" {
			return workspace.WorkInvalidArguments("--confirmation is required only with --apply")
		}
		return validateWorkFormat(format)
	}, RunE: func(c *cobra.Command, args []string) error {
		in.Target = workspace.QueueTargetInput{ProjectID: workspace.ProjectID(project), RepoID: workspace.RepoID(repo), EpicID: workspace.EpicID(epic)}
		var out workspace.TaskAdoptionReadback
		var err error
		if release {
			out, err = s.releaseAdoption(workspace.WorktreeOperationID(operation), reason, apply, confirm)
		} else {
			out, err = s.adopt(in, apply, confirm)
		}
		if err != nil && out.Kind == "" {
			return taskContentCommandError(c, format, err)
		}
		if format == "json" {
			b, e := json.Marshal(out)
			if e != nil {
				return e
			}
			if _, e = c.OutOrStdout().Write(append(b, '\n')); e != nil {
				return e
			}
		} else {
			fmt.Fprintf(c.OutOrStdout(), "Adoption: %s\nTask: %s\nSource: %s %s@%s\nParent: %s@%s\nConfirmation: %s\nGit effects: none\nAgent started: no\nQualified: no\n", out.Phase, out.Plan.Input.TaskID, out.Plan.Input.Path, out.Plan.Input.Ref, out.Plan.Source.OID, out.Plan.Target.ParentRef, out.Plan.Parent.OID, out.Confirmation)
			if release {
				fmt.Fprintln(c.OutOrStdout(), "Metadata only: yes\nWriter absence: not_required")
			}
			for _, r := range out.Reasons {
				fmt.Fprintln(c.OutOrStdout(), r)
			}
			fmt.Fprintf(c.OutOrStdout(), "Check cwd: %s\nNext argv: %q\n", out.CheckCwd, out.NextArgv)
		}
		return err
	}}
	c.Flags().BoolVar(&apply, "apply", false, "apply the exact confirmed metadata request")
	c.Flags().StringVar(&confirm, "confirmation", "", "exact preview SHA-256")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	if release {
		c.Flags().StringVar(&reason, "reason", "", "reason for metadata release")
		_ = c.MarkFlagRequired("reason")
	} else {
		c.Flags().StringVar(&project, "project", "", "registered Project ID")
		c.Flags().StringVar(&repo, "repo", "", "registered repository ID")
		c.Flags().StringVar(&epic, "epic", "", "registered Epic ID")
		c.Flags().StringVar(&in.Path, "path", "", "existing physical linked worktree path")
		c.Flags().StringVar(&in.Ref, "ref", "", "full local source branch ref")
		c.Flags().StringVar(&in.ExpectedOID, "expected-oid", "", "exact clean source commit")
		for _, n := range []string{"project", "repo", "epic", "path", "ref", "expected-oid"} {
			_ = c.MarkFlagRequired(n)
		}
	}
	setWorkFlagErrors(c)
	return c
}
