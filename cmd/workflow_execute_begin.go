package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteBeginCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.DeliveryBeginInput
	var contextPath, format string
	c := &cobra.Command{Use: "begin RUN_ID", Short: "Record actual recipient observations and accept the existing startup mandate",
		Long: "From the bound Task and control executable, report your actual provider, native session, policy and filesystem roots. Ply preserves the observation, hashes evidence, binds the frozen agreement and submits native acceptance. --permission-confirmed means the observed runtime covers all selected effects and targets; the launch request alone cannot establish it. Unknown authority stops dependent work. No handwritten acceptance JSON is needed. Repeat identical observations after a lost response; changed observations never overwrite acceptance.",
		Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
			if err := executeFormat(format); err != nil {
				return err
			}
			if contextPath == "" || in.Provider == "" || in.NativeSession == "" || in.Policy == "" || len(in.ReadRoots) == 0 || len(in.WriteRoots) == 0 {
				return workspace.WorkInvalidArguments("--context, --provider, --native-session, --policy, --read-root and --write-root are required actual observations")
			}
			ws, err := d.Workflow.Workspace.ObserveContaining()
			if err != nil {
				return err
			}
			result, err := taskrun.WorkflowDeliveryBegin(d, ws.Observation.Root, args[0], contextPath, in)
			if outErr := writeExecuteResult(c, format, result); outErr != nil {
				return outErr
			}
			return workflowRunError(err)
		}}
	c.Flags().StringVar(&contextPath, "context", "", "exact private execution context")
	c.Flags().StringVar(&in.Provider, "provider", "", "actual provider: codex or claude")
	c.Flags().StringVar(&in.NativeSession, "native-session", "", "actual native session observed by Herdr")
	c.Flags().StringVar(&in.Model, "model", "", "actual model when known; omit when unknown")
	c.Flags().StringVar(&in.Profile, "profile", "", "actual policy/mode label when available")
	c.Flags().StringVar(&in.Policy, "policy", "", "actual effective policy observation and its source, without secrets")
	c.Flags().StringArrayVar(&in.ReadRoots, "read-root", nil, "actual physical read root; repeat for multiple roots")
	c.Flags().StringArrayVar(&in.WriteRoots, "write-root", nil, "actual physical write root; repeat for multiple roots")
	c.Flags().BoolVar(&in.PermissionConfirmed, "permission-confirmed", false, "observed policy covers the exact mandate, selected effects and targets")
	c.Flags().StringVar(&in.NotificationRoute, "notification-route", "", "existing authorized route to bind natively; omit for internal delegates whose parent owns notifications")
	c.Flags().StringVar(&in.NotificationAuthority, "notification-authority", "", "source of the existing user authorization for this route")
	c.Flags().StringVar(&in.NotificationTimezone, "notification-timezone", "", "recipient's IANA timezone for the authorized route")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
