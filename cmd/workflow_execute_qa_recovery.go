package cmd

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

type workflowQARecoveryActions struct {
	root    func() (string, error)
	preview func(string, taskrun.DeliveryQARecoveryInput) (taskrun.DeliveryQARecoveryPreview, error)
	recover func(string, taskrun.DeliveryQARecoveryInput, string) (taskrun.DeliveryQARecoveryPreview, error)
}

func newWorkflowExecuteQARecoveryCommand(d taskrun.Dependencies) *cobra.Command {
	return workflowExecuteQARecoveryCommand(workflowQARecoveryActions{
		root: func() (string, error) {
			ws, err := d.Workflow.Workspace.ObserveContaining()
			if err != nil {
				return "", err
			}
			return ws.Observation.Root, nil
		},
		preview: func(root string, input taskrun.DeliveryQARecoveryInput) (taskrun.DeliveryQARecoveryPreview, error) {
			return taskrun.WorkflowPreviewDeliveryQARecovery(d, root, input)
		},
		recover: func(root string, input taskrun.DeliveryQARecoveryInput, confirmation string) (taskrun.DeliveryQARecoveryPreview, error) {
			return taskrun.WorkflowRecoverDeliveryQA(d, root, input, confirmation)
		},
	})
}

func workflowExecuteQARecoveryCommand(actions workflowQARecoveryActions) *cobra.Command {
	var input taskrun.DeliveryQARecoveryInput
	var check bool
	var confirmation, format string
	c := &cobra.Command{
		Use:     "recover-qa RUN_ID",
		Short:   "Recover one human QA reservation rejected before publication",
		Long:    "Recover only the exact original QA attempt whose +00:00 timestamp was rejected by the known parser before any human QA publication. Requires the same live owner, original context, preserved evidence and a current runtime observation. --check inspects without writes; applying requires its exact --confirm digest and --attempt. Recovery preserves the rejected bytes and callback control. It records no human verdict, restarts no provider and grants no integration authority. After recovery, submit the actual human answer with a supported timestamp through the original bound control's ordinary qa callback.",
		Example: "  ply workflow execute recover-qa wfr_<digest> --attempt 'qa-<64-hex-attempt-digest>' --context /absolute/context.json --runtime-evidence /absolute/current-runtime.json --check\n  ply workflow execute recover-qa wfr_<digest> --attempt 'qa-<64-hex-attempt-digest>' --context /absolute/context.json --runtime-evidence /absolute/current-runtime.json --confirm sha256:<preview-digest>",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.WorkInvalidArguments("expected one run ID")
			}
			if input.AttemptID == "" || input.ContextPath == "" || input.RuntimeEvidencePath == "" {
				return workspace.WorkInvalidArguments("--attempt, --context and --runtime-evidence are required")
			}
			if check && confirmation != "" {
				return workspace.WorkInvalidArguments("--check cannot be combined with --confirm")
			}
			if !check {
				if !strings.HasPrefix(confirmation, "sha256:") || len(confirmation) != len("sha256:")+64 {
					return workspace.WorkInvalidArguments("apply requires --confirm with the exact sha256 digest from --check")
				}
				if _, err := hex.DecodeString(strings.TrimPrefix(confirmation, "sha256:")); err != nil {
					return workspace.WorkInvalidArguments("--confirm must be the exact sha256 digest from --check")
				}
			}
			return executeFormat(format)
		},
		RunE: func(c *cobra.Command, args []string) error {
			root, err := actions.root()
			if err != nil {
				return workflowRunError(err)
			}
			input.RunID = args[0]
			var out taskrun.DeliveryQARecoveryPreview
			if check {
				out, err = actions.preview(root, input)
			} else {
				out, err = actions.recover(root, input, confirmation)
			}
			if format == "json" {
				if outputErr := writeExecuteResult(c, format, out); outputErr != nil {
					return outputErr
				}
			} else {
				if _, outputErr := fmt.Fprintf(c.OutOrStdout(), "QA recovery %s: %s\n", out.State, out.NextAction); outputErr != nil {
					return outputErr
				}
				if out.Confirmation != "" {
					if _, outputErr := fmt.Fprintf(c.OutOrStdout(), "Confirmation: %s\n", out.Confirmation); outputErr != nil {
						return outputErr
					}
				}
			}
			return workflowRunError(err)
		},
	}
	c.Flags().StringVar(&input.AttemptID, "attempt", "", "exact original QA reservation ID; recovery never selects a newer attempt (required)")
	c.Flags().StringVar(&input.ContextPath, "context", "", "exact original private execution context (required)")
	c.Flags().StringVar(&input.RuntimeEvidencePath, "runtime-evidence", "", "private current observation of the same owner, session and runtime authority (required)")
	c.Flags().BoolVar(&check, "check", false, "inspect the known timestamp rejection and publication absence without writes")
	c.Flags().StringVar(&confirmation, "confirm", "", "exact sha256 digest from the read-only preview (required to apply)")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	c.SetFlagErrorFunc(func(c *cobra.Command, err error) error { return workspace.WorkInvalidArguments(err.Error()) })
	return c
}
