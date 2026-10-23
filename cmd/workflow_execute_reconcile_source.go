package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func addExecuteSourceReconciliation(parent *cobra.Command, d taskrun.Dependencies) {
	var in workspace.TaskSourceReconciliationInput
	var context, format string
	var check, cont, abort bool
	c := &cobra.Command{Use: "reconcile-source RUN_ID", Short: "Merge the exact moved Epic parent into this owned parallel Task source", Long: "Use the original frozen control, same Task cwd and owner --context. --check previews exact source, parent, review and generation without writes. Apply only its --confirm digest. Resolve and stage a known source conflict, then --continue --operation ID; --abort refuses unexpected edits. This changes only the Task source and appends its native effective basis. Fresh independent review and exact verification are required before serial return. Historical serial Runs cannot acquire this authority.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if err := executeFormat(format); err != nil {
			return err
		}
		if context == "" || check && in.Apply || cont && abort || check && (cont || abort) || !check && !in.Apply && !cont && !abort {
			return workspace.WorkInvalidArguments("choose --check, --apply --confirm DIGEST, --continue --operation ID, or --abort --operation ID, with --context")
		}
		if cont || abort {
			if in.OperationID == "" || in.Confirmation != "" {
				return workspace.WorkInvalidArguments("continue/abort requires --operation ID and no new confirmation")
			}
			in.Apply = true
			if cont {
				in.Action = "continue"
			} else {
				in.Action = "abort"
			}
		}
		ws, err := d.Workflow.Workspace.ObserveContaining()
		if err != nil {
			return err
		}
		out, err := taskrun.WorkflowDeliveryReconcileSource(d, ws.Observation.Root, args[0], context, in)
		if format == "json" {
			raw, e := taskrun.Canonical(out)
			if e != nil {
				return e
			}
			if _, e = c.OutOrStdout().Write(append(raw, '\n')); e != nil {
				return e
			}
		} else {
			if _, e := fmt.Fprintf(c.OutOrStdout(), "Source reconciliation: %s\nConfirmation: %s\nNext: %s\n", out.State, out.Confirmation, out.NextAction); e != nil {
				return e
			}
		}
		return workflowRunError(err)
	}}
	c.Flags().StringVar(&context, "context", "", "exact bound owner callback context")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.Flags().BoolVar(&check, "check", false, "preview without source or metadata writes")
	c.Flags().BoolVar(&in.Apply, "apply", false, "apply the exact source merge confirmation")
	c.Flags().StringVar(&in.Confirmation, "confirm", "", "fresh preview SHA256")
	c.Flags().StringVar(&in.OperationID, "operation", "", "same recorded native source operation")
	c.Flags().BoolVar(&cont, "continue", false, "commit the fully staged resolution for this known merge")
	c.Flags().BoolVar(&abort, "abort", false, "abort this known merge only when recorded conflict edits remain unchanged")
	setWorkFlagErrors(c)
	parent.AddCommand(c)
}
