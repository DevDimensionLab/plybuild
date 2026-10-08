package cmd

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/taskexecute"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteContinueCommand(d taskrun.Dependencies) *cobra.Command {
	var input taskexecute.ContinuationInput
	var format string
	c := &cobra.Command{Use: "continue RUN_ID", Short: "Continue the same accepted delivery across a compatible Ply update",
		Long: "Preserve a compatible control for the same live owner, Task worktree, frozen context and accepted authority. The installed CLI validates the old control and required contracts, preserves a durable compatibility proof, and returns the new exact callback executable. It never restarts an agent, sends input, replaces frozen bytes, grants authority or replays uncertain effects. Use --check to inspect without publication; the ordinary path requires no renewed human permission.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.WorkInvalidArguments("expected one run ID")
			}
			if input.ContextPath == "" {
				return workspace.WorkInvalidArguments("--context requires the original private execution context")
			}
			return executeFormat(format)
		},
		RunE: func(c *cobra.Command, args []string) error {
			input.RunID = args[0]
			out, err := taskexecute.ContinueDelivery(d, input)
			if format == "json" {
				raw, e := taskrun.Canonical(out)
				if e != nil {
					return e
				}
				if _, e = c.OutOrStdout().Write(append(raw, '\n')); e != nil {
					return e
				}
			} else if _, e := fmt.Fprintf(c.OutOrStdout(), "Continuation %s: %s\n", out.State, out.NextAction); e != nil {
				return e
			}
			return workflowRunError(err)
		},
	}
	c.Flags().StringVar(&input.ContextPath, "context", "", "exact original private execution context (required)")
	c.Flags().BoolVar(&input.Check, "check", false, "inspect compatibility and the same live session without publication")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.SetFlagErrorFunc(func(c *cobra.Command, err error) error { return workspace.WorkInvalidArguments(err.Error()) })
	return c
}
