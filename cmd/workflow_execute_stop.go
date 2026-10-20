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
	var contextPath, stop, resolution, action, format string
	c := &cobra.Command{Use: "resolve-stop RUN_ID", Short: "Record the same owner's clarified cause and changed action before continuing", Long: "Owner-only continuation after a sourced stop. Bind the exact active stop digest, state the actual clarified cause and name the changed action this owner will take under the same Run mandate. No permission claim is needed: runtime reviews that action when it happens. The refused action itself and unknown or partial effects stay blocked. Preserves the original denial and session; grants no new acceptance, starts no agent, sends no input and releases no writer. If startup is still absent, follow the frozen callback guide: use execute check-in for mandate_receipt_v1, or execute begin with actual observations for legacy contracts. Unknown attempts and unresolved authority remain blocked.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowResolveDeliveryStop(d, ws.Observation.Root, args[0], contextPath, stop, resolution, action)
		if e := writeExecuteResult(c, format, out); e != nil {
			return e
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&contextPath, "context", "", "exact original private context")
	c.Flags().StringVar(&stop, "stop", "", "SHA-256 of the current preserved stop observation")
	c.Flags().StringVar(&resolution, "resolution", "", "actual clarified cause and who clarified it")
	c.Flags().StringVar(&action, "action", "", "changed action this owner will take now; never the refused action")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}

func newWorkflowExecuteReconcileQuotaStopCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.DeliveryQuotaReconcileInput
	var format string
	c := &cobra.Command{Use: "reconcile-quota-stop RUN_ID", Short: "Audit a mistaken pre-execution Claude quota effect for retirement only", Long: "Coordinator annotation of one exact unresolved work stop. Ply locates the bound Claude native transcript itself, requires the supplied source hash, exact usage_limit_reached event and successful synchronous Bash result, and rejects tools or assistant activity after that selected boundary. At the boundary no outstanding tool, async effect or actual interruption may remain. Matched historical returns are source only; their effects are not assessed. The source is copied privately and the original stop remains immutable and active. This grants no owner continuation, acceptance, input or writer release. Only ordinary execute retire consumes the correction, rechecks the unchanged native source, all other stops and attempts, and proves every writer absent. No arbitrary evidence file or bypass flag is accepted. Identical event retries reuse the same observation; conflicting input is rejected.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowReconcileQuotaStop(d, ws.Observation.Root, args[0], in)
		if e := writeExecuteResult(c, format, out); e != nil {
			return e
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&in.EventID, "event-id", "", "stable immutable correction event ID")
	c.Flags().StringVar(&in.StopSHA256, "stop", "", "exact current stop SHA-256")
	c.Flags().StringVar(&in.Provider, "provider", "", "bound provider (claude)")
	c.Flags().StringVar(&in.NativeSession, "native-session", "", "exact bound native session")
	c.Flags().StringVar(&in.Action, "action", "", "exact stopped continuation action, not an earlier tool or unfinished code")
	c.Flags().StringVar(&in.QuotaEvent, "quota-event", "", "exact native quota rejection UUID")
	c.Flags().StringVar(&in.LastToolResult, "last-tool-result", "", "UUID of the last completed synchronous Bash result before quota")
	c.Flags().StringVar(&in.SourceSHA256, "source-sha256", "", "SHA-256 of the complete native transcript including trailing events")
	c.Flags().StringVar(&in.Observer, "observer", "", "actual coordinator correcting the observation")
	c.Flags().StringVar(&in.Reason, "reason", "", "why partial/unknown described unfinished work rather than this unstarted continuation")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}

func newWorkflowExecuteRetireCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.DeliveryRetireInput
	var format string
	c := &cobra.Command{Use: "retire RUN_ID", Short: "End a selected interactive owner and free its queue slot for a fresh Task copy", Long: "Coordinator operation for a selected replacement. Run it from the containing workspace; it needs no callback or private context from the retired owner. Ply closes the exact bound Herdr pane, then requires the pane, its managed agent and every known Task writer to be positively absent. It records retirement, closes the old owner's callbacks and releases only this Task's queue slot. The Run, its refusals, reviews and Task worktree remain history; a fresh Task copy gets its own new mandate and inherits no acceptance. An unresolved attempt, unknown or partial effect, or remaining writer stops with a diagnosis before anything is recorded or released. Repeating the same input after a lost reply is safe.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowRetireDelivery(d, ws.Observation.Root, args[0], in)
		if e := writeExecuteResult(c, format, out); e != nil {
			return e
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&in.Observer, "observer", "", "actual coordinator selecting the replacement")
	c.Flags().StringVar(&in.Reason, "reason", "", "actual reason for the selected replacement (maximum 1900 characters)")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
