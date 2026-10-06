package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspaceview"
	"github.com/spf13/cobra"
)

func viewFormat(c *cobra.Command) string { value, _ := c.Flags().GetString("format"); return value }
func viewString(c *cobra.Command, name string) string {
	value, _ := c.Flags().GetString(name)
	return value
}

func addViewFormat(c *cobra.Command) {
	c.Flags().String("format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
}

func viewNoArgs(c *cobra.Command, args []string) error {
	if len(args) != 0 {
		return workspace.WorkInvalidArguments("expected no positional arguments")
	}
	return validateWorkFormat(viewFormat(c))
}

func viewTaskFilters(c *cobra.Command) workspace.TaskListFilters {
	var f workspace.TaskListFilters
	if c.Flags().Changed("project") {
		value := workspace.ProjectID(viewString(c, "project"))
		f.ProjectID = &value
	}
	if c.Flags().Lookup("repo") != nil && c.Flags().Changed("repo") {
		value := workspace.RepoID(viewString(c, "repo"))
		f.RepoID = &value
	}
	if c.Flags().Lookup("epic") != nil && c.Flags().Changed("epic") {
		value := workspace.EpicID(viewString(c, "epic"))
		f.EpicID = &value
	}
	return f
}

func writeView(c *cobra.Command, value any, heading string, rows []string) error {
	if viewFormat(c) == "json" {
		return json.NewEncoder(c.OutOrStdout()).Encode(value)
	}
	if _, err := fmt.Fprintln(c.OutOrStdout(), heading); err != nil {
		return err
	}
	if len(rows) == 0 {
		_, err := fmt.Fprintln(c.OutOrStdout(), "  None.")
		return err
	}
	for _, row := range rows {
		if _, err := fmt.Fprintln(c.OutOrStdout(), "  "+row); err != nil {
			return err
		}
	}
	return nil
}

func viewTitle(title *string) string {
	if title == nil {
		return "Title unavailable"
	}
	return *title
}

func newWorkspaceStatusCommand(d workspace.Dependencies) *cobra.Command {
	c := &cobra.Command{Use: "status", Short: "Read change digests for registered workspace facts", Long: "Read content digests for Project, work-item, process-journal and run registers without locks or writes. Stable digests allow clients to skip unchanged registered data. They do not detect external Git changes or prove that a process is running. A changing or unreadable source has no cacheable digest.", Example: "  ply workspace status --format json", Args: viewNoArgs}
	c.RunE = func(c *cobra.Command, args []string) error {
		r, err := workspaceview.ReadChanges(d)
		if err != nil {
			return err
		}
		rows := []string{}
		for _, reg := range r.Registers {
			digest := "unknown"
			if reg.RegistrySHA256 != nil {
				digest = *reg.RegistrySHA256
			}
			rows = append(rows, fmt.Sprintf("%s: %s (%s)", reg.ID, digest, reg.Freshness))
		}
		return writeView(c, r, "Registered changes in Ply workspace "+r.Workspace.Root+":", rows)
	}
	addViewFormat(c)
	return c
}

func newWorkspaceOverviewCommand(d workspace.Dependencies) *cobra.Command {
	c := &cobra.Command{Use: "overview", Short: "Read Tasks, Epics and human attention together", Long: "Read registered Task progress, Epic summaries and human attention from one shared source capture. This is the inexpensive Home view: no live Git, Herdr or process checks, locks or writes. Progress keeps native facts separate from agent reports.", Example: "  ply workspace overview --format json\n  ply workspace overview --project ply", Args: func(c *cobra.Command, args []string) error {
		if err := viewNoArgs(c, args); err != nil {
			return err
		}
		return viewTaskFilters(c).Validate()
	}}
	c.RunE = func(c *cobra.Command, args []string) error {
		r, err := workspaceview.ReadOverview(d, workspaceview.EpicFilters{ProjectID: viewTaskFilters(c).ProjectID})
		if err != nil {
			return err
		}
		rows := []string{fmt.Sprintf("%d Tasks; %d Epics; %d human actions (%s registered facts)", len(r.Tasks), len(r.Epics), len(r.Attention), r.Freshness)}
		for _, item := range r.Attention {
			rows = append(rows, fmt.Sprintf("%s: %s — %s", item.SubjectID, item.Kind, item.Reason))
		}
		return writeView(c, r, "Overview of Ply workspace "+r.Workspace.Root+":", rows)
	}
	addViewFormat(c)
	c.Flags().String("project", "", "filter by Project ID")
	return c
}

func newWorkspaceAttentionCommand(d workspace.Dependencies) *cobra.Command {
	c := &cobra.Command{Use: "attention", Short: "List active Tasks that need a human action", Long: "Read explicit human actions for active Tasks in active Epics. Tasks without execution progress and parked, frozen or archived work are excluded. Registered results and process reports keep separate provenance; this command does not perform live Git or agent checks.", Example: "  ply workspace attention --format json\n  ply workspace attention --project ply", Args: func(c *cobra.Command, args []string) error {
		if err := viewNoArgs(c, args); err != nil {
			return err
		}
		return viewTaskFilters(c).Validate()
	}}
	c.RunE = func(c *cobra.Command, args []string) error {
		r, err := workspaceview.ReadAttention(d, viewTaskFilters(c))
		if err != nil {
			return err
		}
		rows := []string{}
		for _, item := range r.Items {
			rows = append(rows, fmt.Sprintf("%s: %s — %s (%s)", item.SubjectID, item.Kind, item.Reason, item.Severity))
		}
		return writeView(c, r, "Human attention in Ply workspace "+r.Workspace.Root+":", rows)
	}
	addViewFormat(c)
	c.Flags().String("project", "", "filter by Project ID")
	c.Flags().String("repo", "", "filter by repository ID")
	c.Flags().String("epic", "", "filter by Epic ID")
	return c
}

func newWorkspaceActivityCommand(d workspace.Dependencies) *cobra.Command {
	parent := &cobra.Command{Use: "journal", Short: "Read activity across the workspace", Long: "Read registered native facts and process contributions across Tasks without starting agents or changing state."}
	c := &cobra.Command{Use: "recent", Short: "List workspace activity in reverse recorded order", Long: "Read recorded activity across Tasks and Epic lifecycle changes. Use --since with an inclusive UTC timestamp or a positive duration, such as 24h. Undated events are shown last without a time filter and counted separately when filtered. A limit is explicit; zero means all events. This is a registered history, not a live Git log.", Example: "  ply workspace journal recent --since 24h --format json\n  ply workspace journal recent --project ply --limit 50"}
	c.Args = func(c *cobra.Command, args []string) error {
		if err := viewNoArgs(c, args); err != nil {
			return err
		}
		if err := viewTaskFilters(c).Validate(); err != nil {
			return err
		}
		if c.Flags().Changed("since") {
			if err := validateViewSince(viewString(c, "since")); err != nil {
				return err
			}
		}
		limit, _ := c.Flags().GetInt("limit")
		if limit < 0 {
			return workspace.WorkInvalidArguments("--limit must be zero or positive")
		}
		return nil
	}
	c.RunE = func(c *cobra.Command, args []string) error {
		s, err := workspaceview.LoadSnapshot(d)
		if err != nil {
			return err
		}
		limit, _ := c.Flags().GetInt("limit")
		r, err := workspaceview.ListActivity(d, s, workspaceview.ActivityOptions{Since: viewString(c, "since"), ProjectID: viewString(c, "project"), Limit: limit})
		if err != nil {
			return err
		}
		rows := []string{}
		for _, event := range r.Events {
			rows = append(rows, fmt.Sprintf("%s: %s", event.EventID, event.Summary))
		}
		return writeView(c, r, "Registered activity in Ply workspace "+s.Workspace.Root+":", rows)
	}
	addViewFormat(c)
	c.Flags().String("project", "", "filter by Project ID")
	c.Flags().String("since", "", "inclusive UTC timestamp or positive duration (for example 24h)")
	c.Flags().Int("limit", 0, "maximum events (0 means all)")
	parent.AddCommand(c)
	return parent
}

func validateViewSince(value string) error {
	if value == "" {
		return workspace.WorkInvalidArguments("--since must not be empty")
	}
	_, err := workspaceview.ParseActivitySince(value, time.Now())
	return err
}

func newWorkspaceRunsCommand(d workspace.Dependencies) *cobra.Command {
	parent := &cobra.Command{Use: "run", Short: "Inspect registered workspace agent runs", Long: "Read preserved run bindings and observed state without contacting Herdr or any provider."}
	c := &cobra.Command{Use: "list", Short: "List registered runs and their agent identity", Long: "List registered Codex and Claude runs across the workspace. Transport is herdr, terminal or headless. Cached active states are unknown as live process evidence. --active selects unresolved runs, including unknown ones; it does not claim they are currently running.", Example: "  ply workspace run list --format json\n  ply workspace run list --active --project ply", Args: func(c *cobra.Command, args []string) error {
		if err := viewNoArgs(c, args); err != nil {
			return err
		}
		return viewTaskFilters(c).Validate()
	}}
	c.RunE = func(c *cobra.Command, args []string) error {
		s, err := workspaceview.LoadSnapshot(d)
		if err != nil {
			return err
		}
		active, _ := c.Flags().GetBool("active")
		r, err := workspaceview.ListRuns(d, s, workspaceview.RunOptions{ActiveOnly: active, ProjectID: viewString(c, "project")})
		if err != nil {
			return err
		}
		rows := []string{}
		for _, run := range r.Runs {
			transport := "unknown"
			if run.Transport != nil {
				transport = *run.Transport
			}
			rows = append(rows, fmt.Sprintf("%s: %s (%s)", run.RunID, run.State, transport))
		}
		return writeView(c, r, "Registered runs in Ply workspace "+s.Workspace.Root+":", rows)
	}
	addViewFormat(c)
	c.Flags().String("project", "", "filter by Project ID")
	c.Flags().Bool("active", false, "show unresolved runs, including unknown live state")
	parent.AddCommand(c)
	return parent
}

func newWorkspaceWorktreesCommand(d workspace.Dependencies) *cobra.Command {
	parent := &cobra.Command{Use: "worktree", Short: "Inspect local worktrees and their registered owners", Long: "Inspect worktrees in registered repositories. This inventory never fetches, prunes, repairs or removes worktrees."}
	c := &cobra.Command{Use: "list", Short: "Observe worktrees in registered local repositories", Long: "Read the local Git worktree inventory, registered ownership and available base comparisons. This explicit observation can cost more than registered Home data. Missing or inaccessible facts remain unknown; no fetch, locks or cleanup are performed.", Example: "  ply workspace worktree list --format json\n  ply workspace worktree list --project ply --repo ply", Args: func(c *cobra.Command, args []string) error {
		if err := viewNoArgs(c, args); err != nil {
			return err
		}
		return viewTaskFilters(c).Validate()
	}}
	c.RunE = func(c *cobra.Command, args []string) error {
		s, err := workspaceview.LoadSnapshot(d)
		if err != nil {
			return err
		}
		r, err := workspaceview.ListWorktrees(d, s, workspaceview.WorktreeOptions{ProjectID: viewString(c, "project"), RepoID: viewString(c, "repo")})
		if err != nil {
			return err
		}
		rows := []string{}
		for _, wt := range r.Worktrees {
			rows = append(rows, wt.Locator)
		}
		return writeView(c, r, "Local worktrees in Ply workspace "+s.Workspace.Root+":", rows)
	}
	addViewFormat(c)
	c.Flags().String("project", "", "filter by Project ID")
	c.Flags().String("repo", "", "filter by repository ID")
	parent.AddCommand(c)
	return parent
}

func addWorkspaceTaskProgressCommands(task *cobra.Command, d workspace.Dependencies) {
	for _, c := range task.Commands() {
		if c.Name() != "list" && c.Name() != "show" {
			continue
		}
		c.Flags().Bool("progress", false, "read registered progress without live Git or agent checks")
		c.Long += " Use --progress for a versioned registered progress projection with explicit lifecycle and provenance; the ordinary mode keeps its existing contract."
		c.Example += "\n  ply workspace task " + c.Name()
		if c.Name() == "show" {
			c.Example += " my-task"
		}
		c.Example += " --progress --format json"
		originalArgs, originalRun := c.Args, c.RunE
		c.Args = func(c *cobra.Command, args []string) error {
			progress, _ := c.Flags().GetBool("progress")
			ready, _ := c.Flags().GetBool("ready")
			if progress && ready {
				return workspace.WorkInvalidArguments("--progress and --ready are mutually exclusive")
			}
			if originalArgs != nil {
				return originalArgs(c, args)
			}
			return nil
		}
		c.RunE = func(c *cobra.Command, args []string) error {
			progress, _ := c.Flags().GetBool("progress")
			if !progress {
				return originalRun(c, args)
			}
			if c.Name() == "show" {
				r, err := workspaceview.ReadTask(d, workspace.TaskID(args[0]))
				if err != nil {
					return err
				}
				return writeView(c, r, "Registered Task progress in Ply workspace "+r.Workspace.Root+":", []string{fmt.Sprintf("%s: %s (%s; %s)", r.Task.TaskID, viewTitle(r.Task.Title), r.Task.Progress.State, r.Task.Lifecycle)})
			}
			r, err := workspaceview.ReadTaskList(d, viewTaskFilters(c))
			if err != nil {
				return err
			}
			rows := []string{}
			for _, row := range r.Tasks {
				rows = append(rows, fmt.Sprintf("%s: %s (%s; %s)", row.TaskID, viewTitle(row.Title), row.Progress.State, row.Lifecycle))
			}
			return writeView(c, r, "Registered Task progress in Ply workspace "+r.Workspace.Root+":", rows)
		}
	}
}

func addWorkspaceEpicReadCommands(epic *cobra.Command, d workspace.Dependencies) {
	for _, c := range epic.Commands() {
		if c.Name() != "list" {
			continue
		}
		addViewFormat(c)
		c.Flags().String("project", "", "filter by Project ID")
		c.Long += " JSON includes current registered bases, worktrees, lifecycle and Task progress counts, including empty Epics. No live Git checks are performed."
		c.Example += "\n  ply workspace epic list --format json\n  ply workspace epic list --project ply"
		originalRun := c.RunE
		c.Args = func(c *cobra.Command, args []string) error {
			if err := viewNoArgs(c, args); err != nil {
				return err
			}
			return viewTaskFilters(c).Validate()
		}
		c.RunE = func(c *cobra.Command, args []string) error {
			if viewFormat(c) == "text" && !c.Flags().Changed("project") {
				return originalRun(c, args)
			}
			r, err := workspaceview.ReadEpics(d, workspaceview.EpicFilters{ProjectID: viewTaskFilters(c).ProjectID})
			if err != nil {
				return err
			}
			rows := []string{}
			for _, row := range r.Epics {
				rows = append(rows, fmt.Sprintf("%s: %s (%s; %d Tasks)", row.EpicID, row.Title, row.Lifecycle, row.TaskCount))
			}
			return writeView(c, r, "Epics in Ply workspace "+r.Workspace.Root+":", rows)
		}
	}
}
