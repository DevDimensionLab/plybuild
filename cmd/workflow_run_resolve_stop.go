package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

func newWorkflowRunResolveStopCommand(d taskrun.Dependencies) *cobra.Command {
	var in taskrun.ReportStopResolutionInput
	var format string
	c := &cobra.Command{Use: "resolve-stop RUN_ID", Short: "Append actual report-effect evidence for a selected pre-work retirement", Long: "Coordinator-only report reconciliation. Bind the preserved stop digest and exact source_sha256 from workflow run show, identify an actual action_result or effect_inspection and provide its private evidence file. Writer absence alone is not effect proof. Original refusal and owner callback blocks remain unchanged. A sourced not_started/none_observed observation can qualify retirement only; partial effects stay blocked. This starts no process, sends no callback, grants no receipt/work authority and releases no reservation. Identical retries preserve the append-only observation.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowResolveReportStop(d, ws.Observation.Root, args[0], in)
		if e := writeExecuteResult(c, format, out); e != nil {
			return e
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&in.StopSHA256, "stop", "", "SHA-256 of the original preserved report/provider stop")
	c.Flags().StringVar(&in.SourceSHA256, "source", "", "exact current report_startup.source_sha256 from run show")
	c.Flags().StringVar(&in.EffectStatus, "effect-status", "", "actually observed effects: not_started, none_observed or partial")
	c.Flags().StringVar(&in.Basis, "basis", "", "actual evidence basis: action_result or effect_inspection; never writer absence alone")
	c.Flags().StringVar(&in.Observer, "observer", "", "actual coordinator observing the clarified effect")
	c.Flags().StringVar(&in.Reason, "reason", "", "concrete observation and why it resolves the prior uncertainty")
	c.Flags().StringVar(&in.EvidencePath, "evidence", "", "private source file containing the actual effect evidence")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
