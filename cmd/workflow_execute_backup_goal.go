package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteBackupGoalCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.BackupGoalInput
	var task, format, confirmation string
	var apply, start, retire bool
	var selected taskrun.DeliveryRetireInput
	var runtime taskexecute.RuntimeOptions
	c := &cobra.Command{Use: "backup-goal RUN_ID", Short: "Preview, register or explicitly start a fresh backup Goal from a retired delivery Run", Long: "Registration by default for an already natively retired, unchanged pre-work Goal Spec@2 automatic local Epic Delivery Run. Preview is read-only. Apply creates fresh Task/Problem/Goal publications and appends to the same queue, preserving current and pending order. Explicit provider replaces the source executor; omitted model and effort resolve from fresh native provider defaults. Product requirements/document bindings remain binding, with a normative fresh runtime identity addendum. No runtime acceptance, tests, review, QA or delivery is inherited. Without --start, creates no branch/worktree and starts no provider. --start previews the selected native runtime and return; --start --apply registers and delegates interactive startup to native execute in Herdr. Identical retries inspect the same execution; uncertain startup is never automatically restarted or resumed. --retire explicitly selects native retirement before fresh registration, requiring --coordinator and --reason plus a preserved source-bound known not_started pre-work stop. Preview binds the exact old owner/stop and fresh choices; apply closes that owner and establishes writer absence before releasing only its current slot. Identical input recovers the same native receipts and publications after each phase. Rule refusals and unknown effects block this prefix. Post-work/review quota, report/legacy, PR and cross-repository support remain later G04 work.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		if !retire && (c.Flags().Changed("coordinator") || c.Flags().Changed("reason")) {
			return workspace.WorkInvalidArguments("coordinator/reason require --retire")
		}
		if retire {
			in.Retire = &selected
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		in.TaskID = workspace.TaskID(task)
		if !start && (c.Flags().Changed("herdr-workspace") || c.Flags().Changed("permission-profile")) {
			return workspace.WorkInvalidArguments("runtime options require --start")
		}
		if start {
			out, err := taskexecute.BackupGoal(d, ws.Observation.Root, args[0], in, runtime, apply, confirmation)
			if format == "json" {
				if e := json.NewEncoder(c.OutOrStdout()).Encode(out); e != nil {
					return e
				}
			} else {
				fmt.Fprintf(c.OutOrStdout(), "Backup startup: %s\nTask: %s; Goal: %s@%d %s\nRegistered: %t\nConfirmation: %s\nReadiness: %s\nNext: %s\n", out.State, out.Registration.TaskID, out.Registration.Goal.SpecID, out.Registration.Goal.Spec.Revision, out.Registration.Goal.Spec.ManifestSHA256, out.Registration.Applied, out.Confirmation, out.Readiness, out.NextAction)
				fmt.Fprintf(c.OutOrStdout(), "Source: %s\nTarget: %s at %s\nDefaults: %s\n", out.Registration.Provenance, out.Registration.Target.ParentRef, out.Registration.Target.ParentLocator, out.Registration.ExecutorSource)
				if out.Registration.Retirement != nil {
					writeBackupRetirement(c, out.Registration.Retirement)
				}
				if out.Runtime != nil {
					fmt.Fprintf(c.OutOrStdout(), "Provider: %s; model: %s; effort: %s; Herdr: %s\n", out.Runtime.Runtime.Provider, out.Runtime.Runtime.Model, out.Runtime.ReasoningEffort, out.Runtime.HerdrWorkspace)
				}
				if out.Execution != nil {
					if e := writeExecuteResult(c, format, *out.Execution); e != nil {
						return e
					}
				}
			}
			return workflowRunError(err)
		}
		out, err := taskrun.WorkflowBackupGoal(d, ws.Observation.Root, args[0], in, apply, confirmation)
		if format == "json" {
			if e := json.NewEncoder(c.OutOrStdout()).Encode(out); e != nil {
				return e
			}
		} else if err == nil {
			_, e := fmt.Fprintf(c.OutOrStdout(), "Fresh backup Task: %s; Goal: %s@%d\nProvider: %s; model: %s; effort: %s\nDefaults: %s\nSource: %s\nTarget: %s at %s\nApplied: %t; queue revision: %d\nConfirmation: %s\nNo provider started. Next read-only check (from return worktree): %s\n", out.TaskID, out.Goal.SpecID, out.Goal.Spec.Revision, *out.Executor.Provider, *out.Executor.Model, *out.Executor.Effort, out.ExecutorSource, out.Provenance, out.Target.ParentRef, out.Target.ParentLocator, out.Applied, out.Queue.Revision, out.Confirmation, strings.Join(out.NextArgv, " "))
			if e != nil {
				return e
			}
		}
		if out.Retirement != nil && format != "json" {
			writeBackupRetirement(c, out.Retirement)
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&task, "task", "", "explicit fresh Task ID")
	c.Flags().StringVar(&in.SpecID, "spec", "", "explicit fresh Spec ID")
	c.Flags().StringVar(&in.Provider, "provider", "", "explicit backup provider (codex or claude)")
	c.Flags().StringVar(&in.Model, "model", "", "optional fresh model override; defaults from selected provider settings")
	c.Flags().StringVar(&in.Effort, "effort", "", "optional fresh effort override; defaults from selected provider settings")
	c.Flags().BoolVar(&retire, "retire", false, "explicitly retire the known stopped pre-work owner before backup registration")
	c.Flags().StringVar(&selected.Observer, "coordinator", "", "actual coordinator selecting --retire")
	c.Flags().StringVar(&selected.Reason, "reason", "", "actual reason for the selected --retire replacement")
	c.Flags().BoolVar(&start, "start", false, "opt in to native interactive startup; preview unless --apply is also supplied")
	c.Flags().StringVar(&runtime.HerdrWorkspace, "herdr-workspace", "", "native Herdr destination with --start (default: current workspace)")
	c.Flags().StringVar(&runtime.PermissionProfile, "permission-profile", "", "existing native permission profile with --start; no configuration changes")
	c.Flags().BoolVar(&apply, "apply", false, "publish and append the previewed fresh Goal; --start also requests native startup")
	c.Flags().StringVar(&confirmation, "confirmation", "", "exact fresh binding confirmation returned by preview")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	executeLogsToStderr(c)
	return c
}

func writeBackupRetirement(c *cobra.Command, r *taskrun.BackupRetirement) {
	fmt.Fprintf(c.OutOrStdout(), "Retirement: %s; completed=%t; Run: %s; pane: %s; session: %s; stop: %s; action: %s; coordinator: %s; reason: %s\nNext retirement action: %s\n", r.State, r.Completed, r.RunID, r.Owner.PaneID, r.Owner.AgentSessionID, r.StopBinding.SHA256, r.Stop.Input.Action, r.Selected.Observer, r.Selected.Reason, r.NextAction)
}
