package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newWorkflowRunCheckInCommand(d taskrun.Dependencies) *cobra.Command {
	receipt := true
	var in taskrun.DeliveryBeginInput
	var contextPath, format string
	c := &cobra.Command{Use: "check-in RUN_ID", Short: "Record the chosen session and assignment receipt",
		Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
			if err := executeFormat(format); err != nil {
				return err
			}
			if contextPath == "" || in.Provider == "" || in.NativeSession == "" {
				return workspace.WorkInvalidArguments("--context, --provider and --native-session are required")
			}
			ws, err := d.Workflow.Workspace.ObserveContaining()
			if err != nil {
				return err
			}
			result, err := taskrun.WorkflowDeliveryReceipt(d, ws.Observation.Root, args[0], contextPath, in, receipt)
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
	c.Flags().StringVar(&in.NotificationRoute, "notification-route", "", "existing authorized route to bind natively; omit for internal delegates whose parent owns notifications")
	c.Flags().StringVar(&in.NotificationAuthority, "notification-authority", "", "source of the existing user authorization for this route")
	c.Flags().StringVar(&in.NotificationTimezone, "notification-timezone", "", "recipient's IANA timezone for the authorized route")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}

func newWorkflowRunRetireCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.DeliveryRetireInput
	var format string
	c := &cobra.Command{Use: "retire RUN_ID", Short: "End a selected interactive owner and free its queue slot for a fresh Task copy", Long: "Coordinator operation for a selected replacement. With --retain-effects, only started/stopped interactive local Epic deliveries with already absent owners are supported. Source and all uncertain effects remain retained; no pane is closed, no effect is resolved and no delivery or replay is granted. Retry with identical inputs to finish a pending exact preparation release after current source and writer guards pass; completed retries read historical observations. Run it from the containing workspace; it needs no callback or private context from the retired owner. For reports, Ply observes the already absent owner and every known writer without closing a live pane. Pre-work reports require clean unchanged source and reconciled effects. Stopped post-start reports retain actual source (including dirty/unique work), refusal history and unknown/partial effects without resolving or replaying them. Their resource retirement creates no acceptance, technical pass, TaskResult, QA or delivery. Accepted terminal history cannot be retired. Installed Ply/provider launcher bytes are historical, and no callback is sent. For delivery, Ply closes the exact bound Herdr pane, then requires the pane, its managed agent and every known Task writer to be positively absent. It records retirement, closes the old owner's callbacks and releases only this Task's queue slot. The Run, its refusals, reviews and Task worktree remain history; a fresh Task copy gets its own new mandate and inherits no acceptance. For delivery and pre-work reports unresolved attempts or unknown/partial effects block retirement; any remaining or unobservable writer stops with a diagnosis before anything is recorded or released. Repeating the same input after a lost reply is safe.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
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
	c.Flags().BoolVar(&in.RetainEffects, "retain-effects", false, "retire only resources of an already absent stopped local Epic owner; retain source and unresolved effects")
	c.Flags().StringVar(&in.Observer, "observer", "", "actual coordinator selecting the replacement")
	c.Flags().StringVar(&in.Reason, "reason", "", "actual reason for the selected replacement (maximum 1900 characters)")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	// The workflow run alias must route first-retirement warnings to stderr too.
	executeLogsToStderr(c)
	return c
}

func executeLogsToStderr(c *cobra.Command) {
	for _, child := range c.Commands() {
		executeLogsToStderr(child)
	}
	if run := c.RunE; run != nil {
		c.RunE = func(c *cobra.Command, args []string) error {
			previousOutput := logrus.StandardLogger().Out
			logrus.SetOutput(c.ErrOrStderr())
			defer logrus.SetOutput(previousOutput)
			return run(c, args)
		}
	}
}

func executeFormat(format string) error {
	if format != "text" && format != "json" {
		return workspace.WorkInvalidArguments("--format must be text or json")
	}
	return nil
}

func writeExecuteResult(c *cobra.Command, format string, result any) error {
	if format == "json" {
		b, err := taskrun.Canonical(result)
		if err != nil {
			return err
		}
		_, err = c.OutOrStdout().Write(append(b, '\n'))
		return err
	}
	if v, ok := result.(taskrun.WorkflowRun); ok {
		fmt.Fprintf(c.OutOrStdout(), "Run: %s\nAgent: %s\nTransport: %s (%s)\n", v.RunID, v.Provider, v.Transport.State, v.Transport.Observation)
		if v.Delivery != nil {
			fmt.Fprintf(c.OutOrStdout(), "Delivery: %s\n", v.Delivery.Phase)
		}
		for _, reason := range v.Reasons {
			fmt.Fprintln(c.OutOrStdout(), "Needs attention: "+reason.Detail)
		}
		_, err := fmt.Fprintf(c.OutOrStdout(), "Next: %s — %s\n", v.NextAction.Actor, v.NextAction.Message)
		return err
	}
	_, err := fmt.Fprintln(c.OutOrStdout(), result)
	return err
}

func newWorkflowRunObserveStopCommand(d taskrun.Dependencies) *cobra.Command {
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
