package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/planningrepo"
	"github.com/spf13/cobra"
)

func newWorkflowEpicPlanCommand(dependencies planningrepo.Dependencies) *cobra.Command {
	var input planningrepo.Input
	var format string
	command := &cobra.Command{
		Use: "plan", Short: "Create a standalone local planning Git repository",
		Long: `Create <root>/planning from the bundled agentic/planning template.
Select every registered member of a Project, or repeat --repo for explicit Git
worktree roots. Explicit repositories require --root; Project mode defaults to
the registered wrapper. All paths are resolved to physical absolute paths.

Validation, complete rendering and Git identity checks precede creation. --check
previews the same checks without writes. Creation commits the generated files on
main using available Git identity. Every existing planning path stops creation;
there is no overwrite, update or continue mode. A partial failure is preserved.
No agent, workspace, Project, Epic, product worktree or remote is created.
Planning prose uses --language; CLI, help and errors remain English.`,
		Example: "  ply workflow epic plan --project ply --language nb\n  ply workflow epic plan --repo /work/example/main --root /work/example --goal \"Plan the next release\"\n  ply workflow epic plan --repo \"/work/trip/service main\" --repo /work/trip/api/main --root /work/trip --language nb --check --format json",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("plan expects no positional arguments")
			}
			if format != "text" && format != "json" {
				return fmt.Errorf("unsupported format %q: use text or json", format)
			}
			return planningrepo.Validate(input)
		},
		// Keep this leaf independent of Spring/Maven/profile initialization even
		// when embedded under another command tree.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := planningrepo.Create(dependencies, input)
			if err != nil {
				return err
			}
			if format == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			}
			if result.State == "preview" {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "Planning repository preview at %s (%s, %d repositories, %d files). No files created.\n", result.PlanningPath, result.Language, len(result.Repositories), len(result.Files))
			} else {
				_, err = fmt.Fprintf(cmd.OutOrStdout(), "Created planning repository at %s (%s).\nBranch %s, commit %s. Start with README.md.\n", result.PlanningPath, result.Language, result.Git.Branch, result.Git.OID)
			}
			return err
		},
	}
	flags := command.Flags()
	flags.StringVar(&input.Project, "project", "", "registered Project ID (all explicit members)")
	flags.StringArrayVar(&input.Repositories, "repo", nil, "existing Git worktree root (repeatable; requires --root)")
	flags.StringVar(&input.Root, "root", "", "existing root directory (default: Project wrapper)")
	flags.StringVar(&input.Language, "language", "en", "planning document language (en or nb)")
	flags.StringVar(&input.Goal, "goal", "", "optional initial planning goal, treated as text")
	flags.BoolVar(&input.Check, "check", false, "validate and render without creating files or Git state")
	flags.StringVar(&format, "format", "text", "output format (text or json)")
	return command
}
