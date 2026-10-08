package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteReleaseCommand(d taskrun.Dependencies) *cobra.Command {
	var contextPath, reason, format string
	c := &cobra.Command{Use: "release RUN_ID", Short: "Relinquish writing ownership of the exact qualified candidate", Long: "Owner-only revocation: preserve the exact context, native session, candidate and reason, then prevent agent integration callbacks from competing with human integration. A newer installed Ply can record this release without replacing the frozen control executable. This command grants no Git, PR merge, installation or cleanup effect.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if e := executeFormat(format); e != nil {
			return e
		}
		if contextPath == "" || reason == "" {
			return workspace.WorkInvalidArguments("--context and --reason are required")
		}
		ws, e := d.Workflow.Workspace.ObserveContaining()
		if e != nil {
			return e
		}
		r, e := taskrun.ReleaseDeliveryOwnership(d, ws.Observation.Root, args[0], contextPath, reason)
		if e != nil {
			return e
		}
		return writeExecuteResult(c, format, r)
	}}
	c.Flags().StringVar(&contextPath, "context", "", "exact private execution context")
	c.Flags().StringVar(&reason, "reason", "", "actual reason for handing this candidate to human integration")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
