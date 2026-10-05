package cmd

import (
	"github.com/devdimensionlab/plybuild/internal/planningrepo"
	"github.com/spf13/cobra"
)

func newWorkflowEpicCommand(dependencies planningrepo.Dependencies) *cobra.Command {
	command := &cobra.Command{
		Use: "epic", Short: "Prepare local Epic planning",
		Long:    "Prepare a standalone planning repository for one or more future Epics without starting an agent.",
		Example: "  ply workflow epic plan --project ply --language nb\n  ply workflow epic plan --repo /work/example/main --root /work/example --check",
	}
	command.AddCommand(newWorkflowEpicPlanCommand(dependencies))
	return command
}
