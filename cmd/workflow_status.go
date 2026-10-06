package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspaceview"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newWorkflowStatusCommand(d workspace.Dependencies) *cobra.Command {
	return newWorkflowStatusCommandWithReader(func(options workspaceview.WorkflowStatusOptions) (workspaceview.WorkflowStatus, error) {
		return workspaceview.ReadWorkflowStatus(d, options)
	})
}

func newWorkflowStatusCommandWithReader(read func(workspaceview.WorkflowStatusOptions) (workspaceview.WorkflowStatus, error)) *cobra.Command {
	var format string
	var includeAll bool
	invalid := func(reason string) error { return fmt.Errorf("workflow_status_invalid_arguments: %s", reason) }
	command := &cobra.Command{
		Use:     "status",
		Short:   "Show registered work and next actions across the workspace",
		Long:    "Show registered work across all projects and Epics in the discovered workspace. Prioritize Needs you, Follow-up, In progress (recorded), and Ready next (registered). Untouched backlog, inactive work and registered completed work are hidden by default and counted; --all includes their history.\n\nExact project, repository and Epic filters are independently optional and combined with AND. Unknown IDs fail; known disjoint filters return an empty result. Run from any directory inside the workspace, including outside an Epic worktree.\n\nThis reads recorded sources without Git, provider or live process checks, initialization, locks or writes. Source freshness does not prove that an agent is alive. Ready next still requires execute preflight.\n\nFor this command, --json selects the same result data as --format json, including when --json precedes workflow. --json=false does not select a format; --json with explicit --format text is an error. JSON success is one WorkflowStatusReadback@1 object followed by a newline. Fatal errors write stderr with empty stdout. Other commands retain their existing --json logging behavior.",
		Example: "  ply workflow status\n  ply workflow status --project ply --epic cli\n  ply workflow status --all --json\n  ply --json workflow status\n  ply workflow status --format json",
		// Status must also remain safe if mounted directly under another root.
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return invalid("expected no positional arguments")
			}
			if format != "text" && format != "json" {
				return invalid("--format must be text or json")
			}
			if workflowStatusJSON(cmd) && cmd.Flags().Changed("format") && format == "text" {
				return invalid("--json conflicts with explicit --format text")
			}
			return viewTaskFilters(cmd).Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Keep the legacy global logging behavior outside this read operation.
			previousOutput := logrus.StandardLogger().Out
			logrus.SetOutput(cmd.ErrOrStderr())
			defer logrus.SetOutput(previousOutput)
			result, err := read(workspaceview.WorkflowStatusOptions{Filters: viewTaskFilters(cmd), IncludeAll: includeAll})
			if err != nil {
				return err
			}
			if format == "json" || workflowStatusJSON(cmd) {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), workflowStatusText(result))
			return err
		},
	}
	command.Flags().StringVar(&format, "format", "text", "output format (text or json); --json is a local result shortcut")
	command.Flags().BoolVar(&includeAll, "all", false, "include untouched backlog, inactive work and registered completed work")
	command.Flags().String("project", "", "exact registered project ID")
	command.Flags().String("repo", "", "exact registered repository ID")
	command.Flags().String("epic", "", "exact registered Epic ID")
	command.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return invalid(err.Error()) })
	return command
}

// Cobra merges inherited flags before Args and RunE, so the existing root flag
// works both before workflow and after status without a second flag or mutation.
func workflowStatusJSON(cmd *cobra.Command) bool {
	value, _ := cmd.Flags().GetBool("json")
	return value
}

func workflowStatusText(result workspaceview.WorkflowStatus) string {
	var out strings.Builder
	fmt.Fprintf(&out, "Workflow status — %s\nRegistered facts · observed %s · source freshness: %s\n", result.Workspace.Root, result.ObservedAtUTC, result.Freshness)
	fmt.Fprintf(&out, "%d projects · %d Epics · %d of %d Tasks shown\n", len(result.Projects), len(result.Epics), result.Counts.Visible, result.Counts.Total)
	for _, group := range []struct{ category, title string }{
		{"needs_you", "Needs you"}, {"follow_up", "Follow-up"},
		{"in_progress", "In progress (recorded)"}, {"ready_next", "Ready next (registered)"},
		{"inactive", "Inactive"}, {"completed", "Completed (registered)"}, {"backlog", "Backlog"},
	} {
		count := result.Counts.Categories[group.category]
		if count == 0 {
			continue
		}
		fmt.Fprintf(&out, "\n%s (%d)\n", group.title, count)
		for _, row := range result.Items {
			if row.Category != group.category {
				continue
			}
			fmt.Fprintf(&out, "  %s / %s / %s / %s  %s\n", row.ProjectID, row.RepoID, row.EpicID, row.TaskID, workflowStatusLine(viewTitle(row.Title)))
			for _, action := range row.NextActions {
				historical := ""
				if !action.Current {
					historical = " · historical"
				}
				fmt.Fprintf(&out, "    %s: %s (%s · %s · %s%s)\n", workflowStatusActor(action.Actor), workflowStatusLine(action.Reason), action.Source, action.Kind, workflowStatusTime(action.SinceUTC), historical)
			}
			fmt.Fprintf(&out, "    Task lifecycle: %s · Epic lifecycle: %s · source freshness: %s · last activity: %s\n", row.TaskLifecycle, row.EpicLifecycle, row.Freshness, workflowStatusTime(row.LastActivityUTC))
			fmt.Fprintf(&out, "    Recorded progress: %s · technical: %s · human QA: %s · integration: %s\n", row.Progress.State, workflowStatusValue(row.Progress.TechnicalGate, "unknown"), workflowStatusValue(row.Progress.HumanQAOutcome, "unknown"), row.Progress.IntegrationClassification)
			for _, run := range row.Runs {
				fmt.Fprintf(&out, "    Run %s · %s · recorded state: %s · state freshness: %s · live state unknown\n", run.RunID, workflowStatusValue(run.Provider, "provider unknown"), run.State, run.StateFreshness)
				if run.Delivery != nil && run.Delivery.Report != nil {
					r := run.Delivery.Report
					fmt.Fprintf(&out, "      Reported %s: %s · %s\n", r.Phase, workflowStatusLine(r.Summary), workflowStatusTime(r.RecordedAtUTC))
				}
			}
			for _, queue := range row.Queues {
				rank := ""
				if queue.Rank != nil {
					rank = fmt.Sprintf(" · rank %d", *queue.Rank)
				}
				preflight := ""
				if queue.Readiness == "ready" {
					preflight = " · start preflight required"
				}
				fmt.Fprintf(&out, "    Queue %s · %s%s · registered readiness: %s%s\n", queue.QueueID, queue.State, rank, queue.Readiness, preflight)
				for _, reason := range queue.Reasons {
					fmt.Fprintf(&out, "      %s: %s\n", reason.Code, workflowStatusLine(reason.Message))
				}
			}
		}
	}
	if len(result.Items) == 0 {
		fmt.Fprintln(&out, "\nNo matching work to show.")
	}
	hidden := result.Counts.Hidden
	fmt.Fprintf(&out, "\nHidden: %d backlog · %d inactive · %d completed.", hidden.Backlog, hidden.Inactive, hidden.Completed)
	if !result.IncludeAll && hidden.Backlog+hidden.Inactive+hidden.Completed > 0 {
		fmt.Fprint(&out, " Use --all to inspect.")
	}
	fmt.Fprintln(&out)
	if len(result.Diagnostics) > 0 {
		fmt.Fprintf(&out, "\nDiagnostics (%d)\n", len(result.Diagnostics))
		for _, diagnostic := range result.Diagnostics {
			fmt.Fprintf(&out, "  %s (%s · %s · %s): %s\n", diagnostic.Code, diagnostic.Source, workflowStatusActor(diagnostic.Actor), diagnostic.Severity, workflowStatusLine(diagnostic.Reason))
			identities := []string{}
			for _, subject := range []struct {
				kind string
				id   *string
			}{{"Task", diagnostic.TaskID}, {"Run", diagnostic.RunID}, {"Queue", diagnostic.QueueID}} {
				if subject.id != nil {
					identities = append(identities, subject.kind+" "+workflowStatusLine(*subject.id))
				}
			}
			if len(identities) > 0 {
				fmt.Fprintf(&out, "    %s\n", strings.Join(identities, " · "))
			}
			if len(diagnostic.EvidenceIDs) > 0 {
				fmt.Fprintf(&out, "    Evidence: %s\n", workflowStatusLine(strings.Join(diagnostic.EvidenceIDs, ", ")))
			}
		}
	}
	return out.String()
}

func workflowStatusActor(actor string) string {
	switch actor {
	case "human":
		return "Human"
	case "agent":
		return "Agent"
	case "ply":
		return "Ply"
	default:
		return "Actor unknown"
	}
}

func workflowStatusTime(value *string) string { return workflowStatusValue(value, "time unknown") }
func workflowStatusValue(value *string, fallback string) string {
	if value == nil {
		return fallback
	}
	return *value
}
func workflowStatusLine(value string) string { return strings.Join(strings.Fields(value), " ") }
