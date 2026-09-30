package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func outputEpicBase(c *cobra.Command, format string, r workspace.EpicBaseResult, err error) error {
	if r.Kind == "" {
		return err
	}
	if format == "json" {
		b, e := workspace.MarshalTaskQueue(r)
		if e != nil {
			return e
		}
		if _, e = c.OutOrStdout().Write(append(b, '\n')); e != nil {
			return e
		}
	} else {
		next := fmt.Sprintf("ply workspace epic base show %s --repo %s", r.Target.EpicID, r.Target.RepoID)
		if r.Operation != nil {
			next = "ply workspace epic base operation show " + r.Operation.ID
		} else if r.Confirmation != nil {
			next = fmt.Sprintf("ply workspace epic base update %s --repo %s --apply --confirm %s", r.Target.EpicID, r.Target.RepoID, *r.Confirmation)
			if r.Plan.RegistryUpgrade != nil {
				next += " --upgrade-registry " + r.Plan.RegistryUpgrade.RegistrySHA256
			}
		}
		if _, e := fmt.Fprintf(c.OutOrStdout(), "Epic base: %s, freshness %s.\nTarget: %s / %s / %s\nWorkplace: %s\n", r.State, r.Freshness, r.Target.ProjectID, r.Target.RepoID, r.Target.EpicID, r.Target.ParentLocator); e != nil {
			return e
		}
		for _, v := range r.Versions {
			if _, e := fmt.Fprintf(c.OutOrStdout(), "Base %d: %s (tree %s)\n", v.Revision, v.OID, v.Tree); e != nil {
				return e
			}
		}
		if r.Plan != nil {
			if _, e := fmt.Fprintf(c.OutOrStdout(), "Affected selections: %d; record a new Spec revision on the new base, assess, select, and update the queue.\n", len(r.Plan.AffectedSelections)); e != nil {
				return e
			}
		}
		if _, e := fmt.Fprintln(c.OutOrStdout(), "Next action: "+next); e != nil {
			return e
		}
	}
	return err
}
func addWorkspaceEpicBaseCommands(epic *cobra.Command, d workspace.Dependencies) {
	base := queueParent("base", "Inspect and explicitly update the versioned Epic base", "  ply workspace epic base show agentic-workflow --repo ply\n  ply workspace epic base update agentic-workflow --repo ply --check")
	epic.AddCommand(base)
	for _, name := range []string{"show", "update"} {
		name := name
		var repo, format, confirm, upgrade string
		var check, apply bool
		c := &cobra.Command{Use: name + " <epic-id>", Short: map[string]string{"show": "Show adopted and current Epic bases", "update": "Preview or confirm a clean descendant as the new Epic base"}[name], Long: "Read the exact registered Epic and repository. Base updates publish append-only metadata after confirmation and never move Git refs. Old Specs, worktrees and results remain historical. New work requires a new Spec revision, ready assessment, explicit selection and queue binding on the current base.", Example: "  ply workspace epic base " + name + " agentic-workflow --repo ply --format json", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
			if e := validateWorkFormat(format); e != nil {
				return e
			}
			id, e := workspace.ParseEpicID(args[0])
			if e != nil {
				return e
			}
			if check && apply || !apply && (confirm != "" || upgrade != "") || apply && confirm == "" {
				return workspace.WorkInvalidArguments("--check and --apply are exclusive; confirmation is required only for apply")
			}
			in := workspace.EpicBaseInput{EpicID: id, RepoID: workspace.RepoID(repo), Apply: apply, Confirmation: confirm, Upgrade: upgrade}
			var r workspace.EpicBaseResult
			if name == "show" {
				r, e = workspace.ShowEpicBase(d, in)
			} else {
				r, e = workspace.UpdateEpicBase(d, in)
			}
			return outputEpicBase(c, format, r, e)
		}}
		c.Flags().StringVar(&repo, "repo", "", "registered repository ID")
		_ = c.MarkFlagRequired("repo")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if name == "update" {
			c.Flags().BoolVar(&check, "check", false, "read-only preview (the default)")
			c.Flags().BoolVar(&apply, "apply", false, "publish the confirmed metadata base")
			c.Flags().StringVar(&confirm, "confirm", "", "exact sha256: preview digest")
			c.Flags().StringVar(&upgrade, "upgrade-registry", "", "explicit legacy registry digest shown in the preview")
		}
		setWorkFlagErrors(c)
		base.AddCommand(c)
	}
	operation := queueParent("operation", "Inspect a preserved base update request", "  ply workspace epic base operation show ebu_<digest> --format json")
	base.AddCommand(operation)
	var format string
	show := &cobra.Command{Use: "show <operation-id>", Short: "Show the preserved base request and metadata outcome", Long: "Read-only recovery inspection after a lost response. An intent is not_committed; a committed request remains historical after later Git drift. Use the same explicit apply to recover an intent.", Example: "  ply workspace epic base operation show ebu_<digest> --format json", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(format); e != nil {
			return e
		}
		r, e := workspace.ShowEpicBaseOperation(d, args[0])
		return outputEpicBase(c, format, r, e)
	}}
	show.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(show)
	operation.AddCommand(show)
}
