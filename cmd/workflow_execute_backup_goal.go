package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
	"strings"
)

func newWorkflowExecuteBackupGoalCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.BackupGoalInput
	var task, format, confirmation string
	var apply bool
	c := &cobra.Command{Use: "backup-goal RUN_ID", Short: "Preview or register a fresh backup Goal from an already retired delivery Run", Long: "Registration only for an already natively retired, unchanged pre-work Goal Spec@2 automatic local Epic Delivery Run. Preview is read-only. Apply creates fresh Task/Problem/Goal publications and appends to the same queue, preserving current and pending order. Explicit provider replaces the source executor; omitted model and effort resolve from fresh native provider defaults. Product requirements/document bindings remain binding, with a normative fresh runtime identity addendum. No runtime acceptance, tests, review, QA or delivery is inherited. Creates no branch/worktree and starts no provider. Live retirement-to-start, post-work/review quota, report/legacy, PR and cross-repository support remain later G04 work.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		in.TaskID = workspace.TaskID(task)
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
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&task, "task", "", "explicit fresh Task ID")
	c.Flags().StringVar(&in.SpecID, "spec", "", "explicit fresh Spec ID")
	c.Flags().StringVar(&in.Provider, "provider", "", "explicit backup provider (codex or claude)")
	c.Flags().StringVar(&in.Model, "model", "", "optional fresh model override; defaults from selected provider settings")
	c.Flags().StringVar(&in.Effort, "effort", "", "optional fresh effort override; defaults from selected provider settings")
	c.Flags().BoolVar(&apply, "apply", false, "publish and append the previewed fresh Goal; never start")
	c.Flags().StringVar(&confirmation, "confirmation", "", "exact fresh binding confirmation returned by preview")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	executeLogsToStderr(c)
	return c
}
