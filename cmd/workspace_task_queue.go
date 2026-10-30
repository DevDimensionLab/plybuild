package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type queueTargetFlags struct{ project, repo, epic string }

func (f *queueTargetFlags) bind(c *cobra.Command) {
	c.Flags().StringVar(&f.project, "project", "", "registered Project ID (requires --repo and --epic)")
	c.Flags().StringVar(&f.repo, "repo", "", "registered repository ID (requires --project and --epic)")
	c.Flags().StringVar(&f.epic, "epic", "", "registered Epic ID (requires --project and --repo)")
}
func (f queueTargetFlags) input() (workspace.QueueTargetInput, error) {
	if (f.project != "" || f.repo != "" || f.epic != "") && (f.project == "" || f.repo == "" || f.epic == "") {
		return workspace.QueueTargetInput{}, workspace.WorkInvalidArguments("provide all of --project, --repo and --epic, or none")
	}
	return workspace.QueueTargetInput{ProjectID: workspace.ProjectID(f.project), RepoID: workspace.RepoID(f.repo), EpicID: workspace.EpicID(f.epic)}, nil
}
func queueParent(name, short, example string) *cobra.Command {
	var format string
	c := &cobra.Command{Use: name, Short: short, Long: short + ". Inspect command help before an explicit mutation. Prepared is not delivery, human QA, or integration approval.", Example: example, Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(format); e != nil {
			return e
		}
		return c.Help()
	}}
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
func outputQueue(c *cobra.Command, format string, r workspace.WorkspaceTaskQueueReadback, err error) error {
	if r.Kind != "" {
		if format == "json" {
			b, e := workspace.MarshalTaskQueue(r)
			if e != nil {
				return e
			}
			if _, e = c.OutOrStdout().Write(append(b, '\n')); e != nil {
				return e
			}
		} else {
			if _, e := fmt.Fprint(c.OutOrStdout(), workspace.TaskQueueText(r)); e != nil {
				return e
			}
		}
	}
	return err
}
func addWorkspaceTaskQueueCommands(task *cobra.Command, d workspace.Dependencies) {
	queue := queueParent("queue", "Inspect and explicitly prioritize pending Tasks", "  ply workspace task queue list --ready\n  ply workspace task queue set --file /absolute/queue.json")
	task.AddCommand(queue)
	var listFormat string
	var ready bool
	var schemaVersion int
	var target queueTargetFlags
	list := &cobra.Command{Use: "list", Short: "Show the prioritized queue and current preparation", Long: "Read the queue without creating locks or upgrading the registry. --ready filters pending entries and preserves ranks; the current preparation remains visible. Omit target flags only from the registered Epic worktree.", Example: "  ply workspace task queue list --ready --project ply --repo ply --epic agentic-workflow", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(listFormat); e != nil {
			return e
		}
		in, e := target.input()
		if e != nil {
			return e
		}
		r, e := workspace.ListTaskQueueVersion(d, in, ready, schemaVersion)
		return outputQueue(c, listFormat, r, e)
	}}
	list.Flags().IntVar(&schemaVersion, "schema-version", 1, "queue read schema (2 includes every active preparation)")
	target.bind(list)
	list.Flags().BoolVar(&ready, "ready", false, "show only ready pending Tasks while retaining the current preparation")
	list.Flags().StringVar(&listFormat, "format", "text", "output format (text or json)")
	setWorkFlagErrors(list)
	queue.AddCommand(list)
	var file, setFormat string
	set := &cobra.Command{Use: "set", Short: "Replace the pending order with exact goal or solution references", Long: "Read strict WorkspaceTaskQueueDraft@1 (selected solutions with a human priority claim) or @2 (goal references with planner provenance) from an absolute physical JSON file. Bind expected_revision and every exact revision/hash. Goal publication records facts and starts no work. A legacy registry requires registry_upgrade. The current preparation is retained and must not occur in entries.", Example: "  ply workspace task queue set --file /absolute/queue.json --format json", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		if e := validateWorkFormat(setFormat); e != nil {
			return e
		}
		r, e := workspace.SetTaskQueue(d, file)
		return outputQueue(c, setFormat, r, e)
	}}
	set.Flags().StringVar(&file, "file", "", "absolute physical JSON draft path (maximum 256 KiB)")
	set.Flags().StringVar(&setFormat, "format", "text", "output format (text or json)")
	_ = set.MarkFlagRequired("file")
	setWorkFlagErrors(set)
	queue.AddCommand(set)
	for _, kind := range []string{"advance", "release"} {
		kind := kind
		var id, reason, format string
		var removedSource bool
		var expectedSourceOID string
		var expected int
		var target queueTargetFlags
		c := &cobra.Command{Use: kind, Short: map[string]string{"advance": "Close a prepared choice without preparing the next Task", "release": "Release a known choice while preserving its worktree and history"}[kind], Long: "Close the exact current preparation with an expected queue revision and a reason. Advance requires a known prepared worktree. Release accepts known no-effect or prepared state. Partial and unknown effects must be resolved first. A released Task cannot be resumed or adopted in this version. No worktree or evidence is deleted; no next Task is prepared.", Example: "  ply workspace task queue " + kind + " --preparation pre_<digest> --expected-revision 2 --reason \"Human queue decision\"", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
			if e := validateWorkFormat(format); e != nil {
				return e
			}
			in, e := target.input()
			if e != nil {
				return e
			}
			if removedSource {
				r, e := workspace.ReleaseTaskQueueRemovedSource(d, in, id, expected, reason, expectedSourceOID)
				return outputQueue(c, format, r, e)
			}
			if expectedSourceOID != "" || c.Flags().Changed("expected-source-oid") {
				return workspace.WorkInvalidArguments("--expected-source-oid requires --removed-source")
			}
			r, e := workspace.CloseTaskQueue(d, in, id, expected, reason, kind)
			return outputQueue(c, format, r, e)
		}}
		target.bind(c)
		c.Flags().StringVar(&id, "preparation", "", "exact current preparation ID")
		c.Flags().IntVar(&expected, "expected-revision", 0, "exact current queue revision")
		c.Flags().StringVar(&reason, "reason", "", "human reason for closing this choice (maximum 2000 characters)")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if kind == "release" {
			c.Long += " With --removed-source and --expected-source-oid, release a known prepared choice only when its path is absent and its current ref is preserved without any checkout. This releases the queue slot; it does not complete the Task or Run, or claim that processes have stopped."
			c.Flags().BoolVar(&removedSource, "removed-source", false, "explicitly release a prepared choice after its source worktree was removed")
			c.Flags().StringVar(&expectedSourceOID, "expected-source-oid", "", "exact current preserved source ref OID (required with --removed-source)")
		}
		for _, f := range []string{"preparation", "expected-revision", "reason"} {
			_ = c.MarkFlagRequired(f)
		}
		setWorkFlagErrors(c)
		queue.AddCommand(c)
	}
	// Ready dispatch is separate from the independent registered-Task filters.
	for _, c := range task.Commands() {
		if c.Name() != "list" {
			continue
		}
		var ready bool
		var queueSchema int
		c.Flags().IntVar(&queueSchema, "schema-version", 1, "queue read schema for --ready (2 includes every active preparation)")
		old := c.RunE
		c.Flags().BoolVar(&ready, "ready", false, "show the ready queue instead of the registered-Task list")
		c.RunE = func(c *cobra.Command, args []string) error {
			format, _ := c.Flags().GetString("format")
			if e := validateWorkFormat(format); e != nil {
				return e
			}
			if !ready {
				return old(c, args)
			}
			project, _ := c.Flags().GetString("project")
			repo, _ := c.Flags().GetString("repo")
			epic, _ := c.Flags().GetString("epic")
			in, e := (queueTargetFlags{project, repo, epic}).input()
			if e != nil {
				return e
			}
			r, e := workspace.ListTaskQueueVersion(d, in, true, queueSchema)
			return outputQueue(c, format, r, e)
		}
		c.Long += " Use --ready for the prioritized queue; an explicit queue target requires all three of --project, --repo and --epic."
	}
	addWorkspaceTaskPrepareCommands(task, d)
}
