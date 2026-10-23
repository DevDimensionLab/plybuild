package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteDispositionCommand(d taskrun.Dependencies) *cobra.Command {
	var contextPath, file, operation, format string
	c := &cobra.Command{Use: "disposition RUN_ID", Short: "Bound same-owner native operations while preserving original stops and unknown effects", Long: "Append-only scoped disposition, never resolve-stop. Bind original Run/request/mandate/context/session, every unresolved stop, all unknown dimensions, actual evidence and per-operation read/write/child assessments. Only prepare, report, separate review, native verify and agreed local integrate are supported. Unknown runtime input/writes/children stay unknown. Current source, acceptance bytes, candidate, review, target and native attempt guards remain required. No cleanup, retirement, writer transfer, permission expansion or refused-action replay. --check-operation observes eligibility without executing that operation. Reuse identical --file bytes after a lost response; changed content needs a new event ID.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if e := executeFormat(format); e != nil {
			return e
		}
		ws, e := d.Workflow.Workspace.ObserveContaining()
		if e != nil {
			return e
		}
		out, e := taskrun.WorkflowDeliveryDisposition(d, ws.Observation.Root, args[0], contextPath, file, operation)
		if err := writeExecuteResult(c, format, out); err != nil {
			return err
		}
		return workflowRunError(e)
	}}
	c.Flags().StringVar(&contextPath, "context", "", "exact original private owner context")
	c.Flags().StringVar(&file, "file", "", "private immutable scoped disposition request JSON")
	c.Flags().StringVar(&operation, "check-operation", "", "observe prepare, report, review, verify or integrate eligibility; no effect")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
