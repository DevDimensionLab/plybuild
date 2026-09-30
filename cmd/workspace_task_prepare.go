package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func outputPreparation(c *cobra.Command, format string, r workspace.TaskPreparationReadback, err error) error {
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
		path := r.Target.ParentLocator
		task := "none"
		if r.Plan != nil {
			path = r.Plan.WorktreePath
			task = string(r.Plan.TaskID)
		}
		next := "ply workspace task queue list"
		if r.Preparation != nil {
			next = "ply workspace task preparation show " + r.Preparation.ID
		} else if r.Confirmation != nil {
			selector := "--next"
			if r.Plan.Selector.TaskID != nil {
				selector = string(*r.Plan.Selector.TaskID)
			}
			next = fmt.Sprintf("ply workspace task prepare %s --project %s --repo %s --epic %s --apply --confirm %s", selector, r.Target.ProjectID, r.Target.RepoID, r.Target.EpicID, *r.Confirmation)
			if r.Plan.RegistryUpgrade != nil {
				next += " --upgrade-registry " + r.Plan.RegistryUpgrade.RegistrySHA256
			}
		}
		if _, e := fmt.Fprintf(c.OutOrStdout(), "Task preparation: %s (%s).\nTarget: %s / %s / %s\nTask: %s\nWorkplace: %s\nFreshness: %s\n", r.State, r.Disposition, r.Target.ProjectID, r.Target.RepoID, r.Target.EpicID, task, path, r.Freshness); e != nil {
			return e
		}
		for _, v := range r.Reasons {
			if _, e := fmt.Fprintf(c.OutOrStdout(), "%s: %s\n", v.Code, v.Message); e != nil {
				return e
			}
		}
		if _, e := fmt.Fprintln(c.OutOrStdout(), "Next action: "+next); e != nil {
			return e
		}
	}
	return err
}
func addWorkspaceTaskPrepareCommands(task *cobra.Command, d workspace.Dependencies) {
	var next, check, apply bool
	var confirm, upgrade, format string
	var target queueTargetFlags
	c := &cobra.Command{Use: "prepare [task-id]", Short: "Preview or confirm one exact queued Task worktree", Long: "Default and --check are read-only. Choose --next or one pending Task ID. --apply requires the exact preview digest and, for a legacy registry, --upgrade-registry. A retry returns the preserved request before ranking pending Tasks. Prepared does not authorize agent start, delivery, human QA or integration.", Example: "  ply workspace task prepare --next --check --format json\n  ply workspace task prepare --next --apply --confirm sha256:<digest>\n  ply workspace task prepare explain-errors --project ply --repo ply --epic agentic-workflow", Args: cobra.MaximumNArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(format); e != nil {
			return e
		}
		if next == (len(args) == 1) {
			return workspace.WorkInvalidArguments("choose exactly one of --next or <task-id>")
		}
		if check && apply || !apply && (confirm != "" || upgrade != "") || apply && confirm == "" {
			return workspace.WorkInvalidArguments("--check and --apply are exclusive; --confirm and --upgrade-registry are only for apply")
		}
		in, e := target.input()
		if e != nil {
			return e
		}
		selector := workspace.QueueSelector{Kind: "next"}
		if len(args) == 1 {
			id, e := workspace.ParseTaskID(args[0])
			if e != nil {
				return e
			}
			selector = workspace.QueueSelector{Kind: "task", TaskID: &id}
		}
		r, e := workspace.PrepareTask(d, workspace.TaskPrepareInput{Target: in, Selector: selector, Apply: apply, Confirmation: confirm, Upgrade: upgrade})
		return outputPreparation(c, format, r, e)
	}}
	target.bind(c)
	c.Flags().BoolVar(&next, "next", false, "use the current preparation or highest ranked ready Task")
	c.Flags().BoolVar(&check, "check", false, "read-only preview (the default)")
	c.Flags().BoolVar(&apply, "apply", false, "confirm the exact preview and create or recover its worktree")
	c.Flags().StringVar(&confirm, "confirm", "", "exact sha256: preview digest")
	c.Flags().StringVar(&upgrade, "upgrade-registry", "", "explicit legacy registry digest shown in the preview")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	task.AddCommand(c)
	p := queueParent("preparation", "Inspect preserved Task preparations", "  ply workspace task preparation show pre_<digest> --format json")
	task.AddCommand(p)
	var showFormat string
	show := &cobra.Command{Use: "show <preparation-id>", Short: "Show a preparation and a separate fresh worktree observation", Long: "Read the immutable request, outcome and queue disposition without recovery or metadata writes. Later commits and dirty files do not erase a historical prepared outcome.", Example: "  ply workspace task preparation show pre_<digest> --format json", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(showFormat); e != nil {
			return e
		}
		r, e := workspace.ShowTaskPreparation(d, args[0])
		return outputPreparation(c, showFormat, r, e)
	}}
	show.Flags().StringVar(&showFormat, "format", "text", "output format (text or json)")
	setWorkFlagErrors(show)
	p.AddCommand(show)
}
