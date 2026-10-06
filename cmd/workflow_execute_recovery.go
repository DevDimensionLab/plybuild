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
		Short:   "Replace an exited Codex startup before its first Task prompt",
		Long:    "Explicitly recover one exited Codex startup in its existing Herdr tab. Require the original Task and return worktrees, no Task-prompt attempt or runtime acceptance, and fresh pane and process-tree evidence of an idle shell without a background provider. Unavailable process inspection stops recovery before startup. Bind the installed Codex and Ply control executables without changing the goal, model, effort, permissions or original records. One replacement is allowed; repeated calls inspect it and never start another agent. Use --check to inspect readiness without writes or startup.",
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
			result, err := taskexecute.RecoverStartup(d, taskexecute.RecoveryInput{RunID: args[0], Check: check, Timeout: time.Duration(timeout) * time.Second})
			if outErr := writeExecuteResult(c, format, result); outErr != nil {
				return outErr
			}
			return workflowRunError(err)
		},
	}
	c.Flags().BoolVar(&check, "check", false, "preview the replacement binding and inspect the existing pane without writes or startup")
	c.Flags().IntVar(&timeout, "timeout", 60, "wait timeout in seconds; never stops the agent")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.SetFlagErrorFunc(func(c *cobra.Command, err error) error { return workspace.WorkInvalidArguments(err.Error()) })
	return c
}
