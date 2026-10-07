package cmd

import (
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteRecoverStartCommand(d taskrun.Dependencies) *cobra.Command {
	var check bool
	var timeout int
	var format string
	c := &cobra.Command{
		Use:     "recover-start RUN_ID",
		Short:   "Recover an interrupted Codex startup before its first Task prompt",
		Long:    "Recover a preserved run by ID. For normal use, repeat ply workflow execute --spec SPEC_ID --restart instead. Recovery preserves the Task, worktree, goal and original records. Fresh checks can reuse an idle terminal or replace a closed terminal; missing transport is not proof of process exit. Only the new attempt may receive the Task prompt or accept callbacks. Each explicit recovery preserves a separate attempt; unknown effects and possible Task input prevent duplicate starts. Model and permissions remain unchanged. Use --check to preview without writes or startup.",
		Example: "  ply workflow execute recover-start wfr_<digest> --check\n  ply workflow execute recover-start wfr_<digest>\n  ply workflow execute show wfr_<digest>",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.WorkInvalidArguments("expected one run ID")
			}
			if err := executeFormat(format); err != nil {
				return err
			}
			if timeout < 1 || timeout > 86400 {
				return workspace.WorkInvalidArguments("--timeout must be between 1 and 86400 seconds")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			result, err := taskexecute.RecoverStartup(d, taskexecute.RecoveryInput{RunID: args[0], Check: check, Restart: true, Timeout: time.Duration(timeout) * time.Second})
			if outErr := writeExecuteResult(c, format, result); outErr != nil {
				return outErr
			}
			return workflowRunError(err)
		},
	}
	c.Flags().BoolVar(&check, "check", false, "preview the next startup attempt and inspect preserved transport without writes or startup")
	c.Flags().IntVar(&timeout, "timeout", 60, "wait timeout in seconds; never stops the agent")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.SetFlagErrorFunc(func(c *cobra.Command, err error) error { return workspace.WorkInvalidArguments(err.Error()) })
	return c
}
