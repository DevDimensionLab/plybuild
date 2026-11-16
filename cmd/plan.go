package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/plantrace"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// Fixed JSON contract names of the plan readbacks.
const (
	planReadbackKind  = "PlyPlanReadback@1"
	planListKind      = "PlyPlanList@1"
	planHistoryKind   = "PlyPlanHistory@1"
	planSchemaVersion = 1
)

// planStdinLimit reads one byte past the store's 128 KiB snapshot limit so
// the store, not the CLI, rejects an oversized snapshot.
const planStdinLimit = 128*1024 + 1

func newPlanCommand(d workspace.Dependencies) *cobra.Command {
	command := &cobra.Command{
		Use:   "plan",
		Short: "Add, update, show, list and trace the plans of a Project",
		Long: "A plan describes one larger piece of work for agents and humans: PLAN-<n>, a goal, scope, owner, status and ordered steps. It grants no mandate, starts nothing and is never derived from Task, Assignment or agent status.\n" +
			"Each plan is one JSON file holding every revision, in <project wrapper>/data/plans/.\n" +
			"Without --project the single Project whose wrapper contains the current directory is used.\n" +
			"add and update read one complete snapshot as JSON from stdin. To change a plan, read it with show --format json, edit .plan and write it back with update --revision <the revision read>.",
		Example: "  ply plan add --project ply --actor ply-epic --reason 'Brukeren valgte oppdraget' < plan.json\n" +
			"  ply plan list --project ply\n" +
			"  ply plan show PLAN-1 --project ply --format json\n" +
			"  ply plan update PLAN-1 --project ply --revision 1 --actor ply-epic --reason 'Spec ferdig' < updated-plan.json\n" +
			"  ply plan history PLAN-1 --project ply",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return planInvalidArguments("unknown plan command")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	// Flag errors from pflag repeat the raw input, so only a fixed detail is kept.
	command.SetFlagErrorFunc(func(*cobra.Command, error) error {
		return planInvalidArguments("invalid flag")
	})

	var project, format, actor, reason, status string
	var revision int
	var details bool
	addFlags := func(c *cobra.Command) *cobra.Command {
		c.Flags().StringVar(&project, "project", "", "Project ID (default: the Project whose wrapper contains the current directory)")
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		command.AddCommand(c)
		return c
	}
	validate := func(positional int, required ...string) cobra.PositionalArgs {
		return func(cmd *cobra.Command, args []string) error {
			if len(args) != positional {
				return planInvalidArguments(fmt.Sprintf("expected %d positional arguments, got %d", positional, len(args)))
			}
			if format != "text" && format != "json" {
				return planInvalidArguments("--format must be text or json")
			}
			for _, name := range required {
				if !cmd.Flags().Changed(name) {
					return planInvalidArguments("--" + name + " is required")
				}
			}
			if cmd.Flags().Changed("revision") && revision < 1 {
				return planInvalidArguments("--revision must be a positive number")
			}
			return nil
		}
	}
	var resolved workspace.ProjectRecord
	open := func() (plan.Store, error) {
		var err error
		if resolved, err = resolvePlanProject(d, project); err != nil {
			return plan.Store{}, err
		}
		directory, err := resolved.Data().Plans()
		if err != nil {
			return plan.Store{}, err
		}
		return plan.Open(directory), nil
	}
	readback := func(cmd *cobra.Command, v plan.View) error {
		if format == "json" {
			return writeAssignmentJSON(cmd.OutOrStdout(), struct {
				Kind          string              `json:"kind"`
				SchemaVersion int                 `json:"schema_version"`
				Project       workspace.ProjectID `json:"project"`
				plan.View
			}{planReadbackKind, planSchemaVersion, resolved.ID, v})
		}
		_, err := io.WriteString(cmd.OutOrStdout(), planText(v, details))
		return err
	}
	write := func(cmd *cobra.Command, store func(plan.Store, []byte) (plan.View, error)) error {
		s, err := open()
		if err != nil {
			return err
		}
		snapshot, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), planStdinLimit))
		if err != nil {
			return &plan.Error{Class: plan.ClassIOError, Detail: "read stdin; no plan file changed"}
		}
		v, err := store(s, snapshot)
		if err != nil {
			return err
		}
		if readback(cmd, v) != nil {
			// The writer's error may carry paths; the stored result is what matters.
			return &plan.Error{Class: plan.ClassIOError, Detail: fmt.Sprintf("%s revision %d was stored; writing its readback failed; run show before writing again", v.ID, v.Revision)}
		}
		return nil
	}

	addCommand := addFlags(&cobra.Command{
		Use:   "add --actor <actor> --reason <reason> < plan.json",
		Short: "Add a plan from a complete JSON snapshot on stdin",
		Long: "Add a plan as revision 1 from one complete JSON snapshot on stdin (at most 128 KiB). Unknown or duplicate keys and trailing data are refused. It prints the stored revision.\n" +
			"--actor names who writes (attribution, not authentication); --reason says why. If the answer is lost or the result is unknown, run list before trying again.",
		Args: validate(0, "actor", "reason"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return write(cmd, func(s plan.Store, snapshot []byte) (plan.View, error) {
				return s.Add(snapshot, actor, reason)
			})
		},
	})
	addCommand.Flags().StringVar(&actor, "actor", "", "who writes, one line (required)")
	addCommand.Flags().StringVar(&reason, "reason", "", "why, one line (required)")

	updateCommand := addFlags(&cobra.Command{
		Use:   "update <PLAN-n> --revision <r> --actor <actor> --reason <reason> < plan.json",
		Short: "Store a complete JSON snapshot on stdin as the next revision",
		Long: "Store the snapshot as the next revision when --revision is the current revision. Otherwise plan_conflict is returned and nothing is written: read the plan again and decide; there is no merge or retry.\n" +
			"An identical snapshot is stored as a new revision too. If the result is unknown, run show before writing again.",
		Args: validate(1, "revision", "actor", "reason"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return write(cmd, func(s plan.Store, snapshot []byte) (plan.View, error) {
				return s.Update(args[0], revision, snapshot, actor, reason)
			})
		},
	})
	updateCommand.Flags().IntVar(&revision, "revision", 0, "current revision read with show (required)")
	updateCommand.Flags().StringVar(&actor, "actor", "", "who writes, one line (required)")
	updateCommand.Flags().StringVar(&reason, "reason", "", "why, one line (required)")

	showCommand := addFlags(&cobra.Command{
		Use:   "show <PLAN-n> [--revision <r>] [--details]",
		Short: "Show the latest or one earlier revision of a plan",
		Long: "The text view is short: title, status, goal, summary, owner, next action, attention and the ordered steps with their next action and attention. --details adds scope, details, results, decisions, learning, dependencies and refs.\n" +
			"JSON always holds the whole revision.",
		Args: validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			v, err := s.Show(args[0], revision)
			if err != nil {
				return err
			}
			return readback(cmd, v)
		},
	})
	showCommand.Flags().IntVar(&revision, "revision", 0, "revision to show (default: the latest)")
	showCommand.Flags().BoolVar(&details, "details", false, "add scope, details, results, decisions, learning and refs to the text view")

	listCommand := addFlags(&cobra.Command{
		Use:   "list [--status <status>]",
		Short: "List the latest revision of every plan, most recently updated first",
		Long: "Files that cannot be read are always listed under invalid, also with --status.\n" +
			"Text lines are <id> <status>/<selection> revision <n> <updated_at> <title>, with the next action below and a mark when a human answer is needed.",
		Args: validate(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			result, err := s.List(status)
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					plan.ListResult
				}{planListKind, planSchemaVersion, resolved.ID, result})
			}
			_, err = io.WriteString(cmd.OutOrStdout(), planListText(result))
			return err
		},
	})
	listCommand.Flags().StringVar(&status, "status", "", "only plans with this status: "+strings.Join(plan.Statuses, ", "))

	command.AddCommand(newPlanJournalCommand(&project, open))

	addFlags(&cobra.Command{
		Use:   "trace <PLAN-n>",
		Short: "Read the recorded workflow of a plan: revisions, linked Tasks and their journal events",
		Long: "A Task belongs to the plan when the latest plan link in its journal names it (goal register or task create with --plan). Records of ply plan journal have an empty task_id.\n" +
			"JSON is " + plantrace.Kind + ": the plan revisions, the linked Tasks with their plan step, and every event of their journals with actor, role, step, outcome and time, ordered by occurred, else recorded time. Nothing is interpreted.",
		Args: validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			trace, err := plantrace.Read(d, resolved.ID, s, args[0])
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), trace)
			}
			var out strings.Builder
			fmt.Fprintf(&out, "%s %s (%s): %d revisions, %d Tasks, %d events\n", trace.Plan.ID, trace.Plan.Title, trace.Plan.Status, len(trace.Plan.Revisions), len(trace.Tasks), len(trace.Events))
			for _, e := range trace.Events {
				when, who := "unknown", "unknown"
				if t := at(e.OccurredAt, e.RecordedAt); t != nil {
					when = *t
				}
				if e.ActorID != nil {
					who = *e.ActorID + "/" + *e.Role
				}
				fmt.Fprintf(&out, "%s %s %s %s: %s\n", when, e.TaskID, who, e.Type, e.Title)
			}
			_, err = io.WriteString(cmd.OutOrStdout(), out.String())
			return err
		},
	})

	addFlags(&cobra.Command{
		Use:   "history <PLAN-n>",
		Short: "List revision, time, actor and reason of every revision",
		Long:  "Revisions are listed in ascending order, never truncated. Read one with show --revision.",
		Args:  validate(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			history, err := s.History(args[0])
			if err != nil {
				return err
			}
			if format == "json" {
				return writeAssignmentJSON(cmd.OutOrStdout(), struct {
					Kind          string              `json:"kind"`
					SchemaVersion int                 `json:"schema_version"`
					Project       workspace.ProjectID `json:"project"`
					ID            string              `json:"id"`
					Revisions     []plan.HistoryEntry `json:"revisions"`
				}{planHistoryKind, planSchemaVersion, resolved.ID, args[0], history})
			}
			var out strings.Builder
			for _, entry := range history {
				fmt.Fprintf(&out, "%d %s %s: %s\n", entry.Revision, entry.At, entry.Actor, entry.Reason)
			}
			_, err = io.WriteString(cmd.OutOrStdout(), out.String())
			return err
		},
	})
	return command
}

// resolvePlanProject maps Project lookup failures to fixed details without
// paths, as resolveAssignmentProject does.
func resolvePlanProject(d workspace.Dependencies, value string) (workspace.ProjectRecord, error) {
	project, err := resolveAssignmentProject(d, value)
	var projectError *workspace.ProjectError
	if err != nil && !errors.As(err, &projectError) {
		return workspace.ProjectRecord{}, &plan.Error{Class: plan.ClassIOError, Detail: "could not resolve the project"}
	}
	return project, err
}

func planText(v plan.View, details bool) string {
	var out strings.Builder
	p := v.Plan
	fmt.Fprintf(&out, "%s %s\n", v.ID, p.Title)
	fmt.Fprintf(&out, "status: %s (%s)\n", p.Status, p.Selection)
	writeField(&out, "", "goal", p.Goal)
	writeField(&out, "", "summary", p.Summary)
	fmt.Fprintf(&out, "owner: %s; revision %d, updated %s by %s (created %s)\n", p.Owner, v.Revision, v.UpdatedAt, v.Actor, v.CreatedAt)
	writeField(&out, "", "next", p.NextAction)
	writeAttention(&out, "", p.Attention)
	if v.NeedsHuman {
		out.WriteString("needs human: yes\n")
	}
	if details {
		writeField(&out, "", "scope", p.Scope)
		writeField(&out, "", "details", p.Details)
		writeList(&out, "decisions", p.Decisions)
		writeField(&out, "", "learning", p.Learning)
		writeField(&out, "", "refs", strings.Join(p.Refs, ", "))
	}
	progress := v.Progress
	fmt.Fprintf(&out, "steps (%d: %d planned, %d active, %d waiting, %d blocked, %d done, %d unknown):\n",
		progress.Total, progress.Planned, progress.Active, progress.Waiting, progress.Blocked, progress.Done, progress.Unknown)
	for i, step := range p.Steps {
		fmt.Fprintf(&out, "  %d. [%s] %s: %s (owner %s)\n", i+1, step.Status, step.ID, step.Title, step.Owner)
		writeField(&out, "     ", "next", step.NextAction)
		writeAttention(&out, "     ", step.Attention)
		if details {
			writeField(&out, "     ", "outcome", step.Outcome)
			writeField(&out, "     ", "result", step.Result)
			writeField(&out, "     ", "depends on", strings.Join(step.DependsOn, ", "))
			writeField(&out, "     ", "refs", strings.Join(step.Refs, ", "))
		}
	}
	return out.String()
}

// writeField writes a non-empty value, indenting its further lines.
func writeField(out *strings.Builder, indent, name, value string) {
	if value == "" {
		return
	}
	value = strings.ReplaceAll(strings.TrimRight(value, "\n"), "\n", "\n"+indent+"  ")
	fmt.Fprintf(out, "%s%s: %s\n", indent, name, value)
}

func writeList(out *strings.Builder, name string, values []string) {
	if len(values) == 0 {
		return
	}
	fmt.Fprintf(out, "%s:\n", name)
	for _, value := range values {
		writeField(out, "  ", "-", value)
	}
}

func writeAttention(out *strings.Builder, indent string, a *plan.Attention) {
	if a == nil {
		return
	}
	needs := "no"
	if a.NeedsHuman {
		needs = "yes"
	}
	writeField(out, indent, "attention", fmt.Sprintf("%s (who: %s, needs human: %s)", a.Reason, a.Who, needs))
}

func planListText(result plan.ListResult) string {
	var out strings.Builder
	if len(result.Items) == 0 {
		out.WriteString("no plans\n")
	}
	for _, item := range result.Items {
		mark := ""
		if item.NeedsHuman {
			mark = " [needs human]"
		}
		fmt.Fprintf(&out, "%s %s/%s revision %d %s %s%s\n", item.ID, item.Status, item.Selection, item.Revision, item.UpdatedAt, item.Title, mark)
		writeField(&out, "  ", "next", item.NextAction)
	}
	for _, invalid := range result.Invalid {
		fmt.Fprintf(&out, "warning: invalid plan file %s: %s\n", invalid.File, invalid.Error)
	}
	return out.String()
}

// at is the first of the times that is known.
func at(times ...*string) *string {
	for _, t := range times {
		if t != nil {
			return t
		}
	}
	return nil
}

func planInvalidArguments(detail string) error {
	return &plan.Error{Class: plan.ClassInvalidArguments, Detail: detail}
}

func init() {
	RootCmd.AddCommand(newPlanCommand(workspace.SystemDependencies()))
}
