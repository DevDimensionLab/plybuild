package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

func workflowRunError(e error) error {
	if e == nil {
		return nil
	}
	var t *taskrun.Error
	if errors.As(e, &t) {
		if t.Exit == 3 {
			return &taskrun.Error{Code: t.Code, Detail: t.Detail, Exit: 4}
		}
		return e
	}
	return e
}
func newWorkflowRunCommand(d taskrun.Dependencies) *cobra.Command {
	run := &cobra.Command{Use: "run", Short: "Run a prepared Task in Herdr and review its reports", Long: "Start one authorized Claude or Codex session for an existing Task preparation, follow immutable round reports, and request bounded corrections in the same session. Reviewed reports do not publish TaskResult or attest provider inactivity. Human QA and integration remain separate gates.", Example: "  ply workflow run start --file /absolute/prepared-task-request.json --check\n  ply workflow run show wfr_<digest>", Args: cobra.NoArgs, RunE: func(c *cobra.Command, a []string) error { return c.Help() }}
	usage := func(detail string) error {
		return &taskrun.Error{Code: "workflow_run_invalid_arguments", Detail: detail, Exit: 2}
	}
	run.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return usage(e.Error()) })
	for _, name := range []string{"start", "show", "follow", "accept", "report", "review"} {
		name := name
		var file, format, confirm, contextPath string
		var check, apply bool
		var timeout int
		use := name
		if name != "start" {
			use += " RUN_ID"
		}
		short := map[string]string{"start": "Preview or confirm one prepared Task start in Herdr", "show": "Read preserved run and artifact freshness without effects", "follow": "Observe the same session until its round is ready for review", "accept": "Accept the actual bound runtime before Task writes", "report": "Preserve an immutable report for the current round", "review": "Accept, block or request one bounded correction"}[name]
		detail := map[string]string{
			"start":  " Select runtime.provider in the request: codex (default when omitted) or claude. runtime.model is a separate choice. Claude requires null config_profile and the manual permission mode; Codex keeps its managed permission profile. --check is the read-only default. --apply requires --confirm with the exact preview digest and the local HERDR_ENV=1 context. The request must bind an existing native Task preparation, Spec and permission evidence. An optional codex_project_trust grant confirms only the exact physical repository root for this process, before launch. Preview binds the trust effect and policy; sandbox and approvals remain enforced and no persistent config is changed. Project execution layers, explicit distrust or unknown policy stop before a tab is created. An optional Claude-only claude_project_trust grant with mode configuration persists folder trust for the exact prepared worktree before launch. Preview shows the config path and effect; failed trust writes warn and continue the same native start. Tool permissions remain unchanged. Repeated starts reuse the reservation without launching or writing trust again. No TaskResult is published.",
			"show":   " Reads saved state and current artifact hashes without Herdr calls, locks or recovery. Cached transport is not a fresh provider observation.",
			"follow": " Waits for the correct round report and two fresh idle/done observations of the exact same selected agent session. --timeout defaults to 60 seconds. Timeout ends observation, not the agent. A finished tab is transport status, not technical or human approval.",
			"accept": " This is the recipient's first action. Supply private Acceptance@2 fields using kind ply.workflow.run-acceptance, schema_version 1. The exact context, executable, cwd, native session and necessary permission facts must match. A null actual model is permitted; known mismatch is rejected.",
			"report": " Submit native TaskRun report fields using kind ply.workflow.round-report, schema_version 1, plus the current round, control_id and previous_report_sha256. Include cumulative usage and the frozen Spec's task-requirements artifact. Run after the last target write of this round. No native terminal is published until review accepts.",
			"review": " Bind the named coordinator's accepted, blocked or changes_requested decision to the current immutable report. Corrections reserve shared agreement A before one prompt attempt. A repeated decision never resends input. Accepted seals a native terminal but does not publish TaskResult, attest provider inactivity, perform human QA or integrate the candidate.",
		}[name]
		example := map[string]string{"start": " --file /absolute/prepared-task-request.json --check", "show": " wfr_<digest> --format json", "follow": " wfr_<digest> --timeout 60", "accept": " wfr_<digest> --context /absolute/current/context.json --file /absolute/acceptance.json", "report": " wfr_<digest> --context /absolute/current/context.json --file /absolute/round-report.json", "review": " wfr_<digest> --file /absolute/review.json"}[name]
		c := &cobra.Command{Use: use, Short: short, Long: short + "." + detail, Example: "  ply workflow run " + name + example}
		c.Args = func(c *cobra.Command, a []string) error {
			want := 1
			if name == "start" {
				want = 0
			}
			if len(a) != want {
				return usage(fmt.Sprintf("expected %d positional arguments", want))
			}
			if format != "text" && format != "json" {
				return usage("--format must be text or json")
			}
			if (name == "start" || name == "accept" || name == "report" || name == "review") && file == "" {
				return usage("--file is required")
			}
			if c.Flags().Changed("check") && apply {
				return usage("--check and --apply are mutually exclusive")
			}
			if c.Flags().Changed("confirm") && !apply || apply && confirm == "" {
				return usage("--confirm requires --apply; apply requires the exact confirmation")
			}
			if (name == "accept" || name == "report") && contextPath == "" {
				return usage("--context is required")
			}
			if name == "follow" && (timeout <= 0 || timeout > 86400) {
				return usage("--timeout must be between 1 and 86400 seconds")
			}
			if f := c.Flag("force"); f != nil && f.Changed {
				return usage("--force cannot authorize a workflow run")
			}
			return nil
		}
		c.RunE = func(c *cobra.Command, a []string) error {
			var result any
			var e error
			if name == "start" {
				if apply {
					result, e = taskrun.WorkflowStart(d, file, confirm)
				} else {
					result, e = taskrun.WorkflowPreviewStart(d, file)
				}
			} else {
				ws, err := d.Workflow.Workspace.ObserveContaining()
				if err != nil {
					return err
				}
				root := ws.Observation.Root
				switch name {
				case "show":
					result, e = taskrun.WorkflowShow(d, root, a[0])
				case "follow":
					result, e = taskrun.WorkflowFollow(d, root, a[0], time.Duration(timeout)*time.Second)
				case "accept":
					result, e = taskrun.WorkflowAccept(d, root, a[0], contextPath, file)
				case "report":
					result, e = taskrun.WorkflowReportRound(d, root, a[0], contextPath, file)
				case "review":
					result, e = taskrun.WorkflowReviewRun(d, root, a[0], file)
				}
			}
			if result == nil {
				result = map[string]any{"kind": "ply.workflow.run-error", "schema_version": 1, "message": fmt.Sprint(e)}
			}
			if format == "json" {
				b, err := taskrun.Canonical(result)
				if err != nil {
					return err
				}
				if _, err = c.OutOrStdout().Write(append(b, '\n')); err != nil {
					return err
				}
			} else {
				switch v := result.(type) {
				case taskrun.WorkflowRun:
					fmt.Fprintf(c.OutOrStdout(), "Run: %s\nAgent: %s\nTransport: %s (%s); round %d: %s\nReport review: %s\nTaskResult: not published. Provider inactivity: not attested. Human QA and integration: separate pending gates.\nNext: %s — %s\n", v.RunID, v.Provider, v.Transport.State, v.Transport.Observation, v.Round.Number, v.Round.State, v.FinalReturn.State, v.NextAction.Actor, v.NextAction.Message)
					writeWorkflowClaudeTrust(c.OutOrStdout(), v.ClaudeTrust)
					for _, r := range v.Reasons {
						fmt.Fprintln(c.OutOrStdout(), "Needs attention: "+r.Detail)
					}
				case taskrun.WorkflowPreview:
					fmt.Fprintf(c.OutOrStdout(), "Run: %s\nAgent: %s\nRead-only preview for the existing Task preparation.\n", v.RunID, v.Provider)
					for _, effect := range v.Effects {
						fmt.Fprintln(c.OutOrStdout(), effect)
					}
					if v.Confirmation != nil {
						fmt.Fprintf(c.OutOrStdout(), "Confirmation: %s\nApply: ply workflow run start --file %s --apply --confirm %s\n", *v.Confirmation, taskrun.ShellQuote(file), *v.Confirmation)
					}
					for _, r := range v.Reasons {
						fmt.Fprintln(c.OutOrStdout(), r.Detail)
					}
				}
			}
			return workflowRunError(e)
		}
		c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
		if name == "start" || name == "accept" || name == "report" || name == "review" {
			c.Flags().StringVar(&file, "file", "", "absolute physical JSON input path")
		}
		if name == "start" {
			c.Flags().BoolVar(&check, "check", true, "preview without writing (default)")
			c.Flags().BoolVar(&apply, "apply", false, "apply the exact confirmed preview")
			c.Flags().StringVar(&confirm, "confirm", "", "exact sha256 preview confirmation")
		}
		if name == "accept" || name == "report" {
			c.Flags().StringVar(&contextPath, "context", "", "exact private context path for the current round")
		}
		if name == "follow" {
			c.Flags().IntVar(&timeout, "timeout", 60, "positive observation timeout in seconds (does not stop the agent)")
		}
		c.SetFlagErrorFunc(func(c *cobra.Command, e error) error { return usage(e.Error()) })
		run.AddCommand(c)
	}
	return run
}
