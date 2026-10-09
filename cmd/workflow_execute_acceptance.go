package cmd

import (
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func newWorkflowExecuteAcceptanceCommand(d taskrun.Dependencies) *cobra.Command {
	var contextPath, file, format string
	c := &cobra.Command{
		Use: "acceptance RUN_ID", Short: "Record an explicit automatic acceptance choice for the same local Epic run",
		Long:    "Read a private DeliveryAcceptanceChoice@1 JSON file recording who requested automatic acceptance, their exact policy-choice answer, the previous agreement and the current candidate. Preserve the original frozen agreement and prior attempts while selecting automatic acceptance for this existing local Epic run. The same accepted owner and original context are required. The choice supplies no test result, human QA or runtime permission and changes no source, return target or delivery effect. Verify the candidate under the selected policy before integration. If the immutable control predates this command, first use workflow execute continue with the private or installed CLI and then use its returned compatible control.",
		Example: "  /absolute/compatible/ply-control workflow execute acceptance wfr_<digest> --context /absolute/context.json --file /absolute/acceptance-choice.json --format json",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return workspace.WorkInvalidArguments("expected one run ID")
			}
			if !filepath.IsAbs(contextPath) || filepath.Clean(contextPath) != contextPath {
				return workspace.WorkInvalidArguments("--context requires the exact absolute original private execution context")
			}
			if !filepath.IsAbs(file) || filepath.Clean(file) != file {
				return workspace.WorkInvalidArguments("--file requires an absolute path to a private DeliveryAcceptanceChoice@1 JSON file")
			}
			return executeFormat(format)
		},
		RunE: func(c *cobra.Command, args []string) error {
			observed, err := d.Workflow.Workspace.ObserveContaining()
			if err != nil {
				return err
			}
			result, err := taskrun.WorkflowDeliverySelectAcceptance(d, observed.Observation.Root, args[0], contextPath, file)
			if outErr := writeExecuteResult(c, format, result); outErr != nil {
				return outErr
			}
			return workflowRunError(err)
		},
	}
	c.Flags().StringVar(&contextPath, "context", "", "exact original private execution context (required)")
	c.Flags().StringVar(&file, "file", "", "absolute private DeliveryAcceptanceChoice@1 JSON file (required)")
	c.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	setWorkFlagErrors(c)
	return c
}
