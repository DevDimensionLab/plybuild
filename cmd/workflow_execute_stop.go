package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteObserveStopCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.DeliveryStopInput
	var format string
	c := &cobra.Command{Use: "observe-stop RUN_ID", Short: "Preserve a sourced provider stop without impersonating its owner", Long: "Coordinator entrypoint when startup, acceptance or report transport is denied. Supply the actual provider, bound native session (omit only before session binding), stage, reason and a private evidence file. Ply preserves and hashes the source, retains owner reservations and makes the stop visible. This is not runtime acceptance and sends no provider input. Reuse the same event ID after a lost reply; conflicting evidence cannot overwrite it.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowObserveDeliveryStop(d, ws.Observation.Root, args[0], in)
		if e := writeExecuteResult(c, format, out); e != nil {
			return e
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&in.EventID, "event-id", "", "stable ID for this actual stop; reuse after a lost reply")
	c.Flags().StringVar(&in.Provider, "provider", "", "actual provider")
	c.Flags().StringVar(&in.NativeSession, "native-session", "", "bound native session; absent only before binding")
	c.Flags().StringVar(&in.Stage, "stage", "", "denied stage: startup, acceptance, receipt, report, work, verification or delivery")
	c.Flags().StringVar(&in.Observer, "observer", "", "actual coordinator observing the denial")
	c.Flags().StringVar(&in.Reason, "reason", "", "actual observed reason, without guessed classifier causes")
	c.Flags().StringVar(&in.Action, "action", "", "actual rejected action; required after startup")
	c.Flags().StringVar(&in.EffectStatus, "effect-status", "", "known effects: not_started, none_observed, partial or unknown")
	c.Flags().StringVar(&in.Rule, "rule", "", "observed runtime rule name when available; do not infer it")
	c.Flags().StringVar(&in.EvidencePath, "evidence", "", "private file containing the observed denial")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}

func newWorkflowExecuteResolveStopCommand(d taskrun.Dependencies) *cobra.Command {
	var contextPath, stop, resolution, format string
	var permission bool
	c := &cobra.Command{Use: "resolve-stop RUN_ID", Short: "Record the same owner's repaired cause before supported continuation", Long: "Owner-only continuation after a sourced stop. Bind the exact active stop digest and state the actual cause repair with current permission confirmed. Preserves the original denial and session; grants no new acceptance, starts no agent, sends no input and releases no writer. If startup is still absent, follow the frozen callback guide: use execute check-in for mandate_receipt_v1, or execute begin with actual observations for legacy contracts. Unknown attempts and unresolved authority remain blocked.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowResolveDeliveryStop(d, ws.Observation.Root, args[0], contextPath, stop, resolution, permission)
		if e := writeExecuteResult(c, format, out); e != nil {
			return e
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&contextPath, "context", "", "exact original private context")
	c.Flags().StringVar(&stop, "stop", "", "SHA-256 of the current preserved stop observation")
	c.Flags().StringVar(&resolution, "resolution", "", "actual repaired cause and source of the current permission observation")
	c.Flags().BoolVar(&permission, "permission-confirmed", false, "current actual permission covers continuation within the original mandate")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
