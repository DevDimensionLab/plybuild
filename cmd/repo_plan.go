package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/repoestablish"
	"github.com/devdimensionlab/plybuild/internal/repoestablish/system"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
	"github.com/devdimensionlab/plybuild/internal/wrapperplan"
	"github.com/spf13/cobra"
)

func newRepoPlanCommand(runner process.Runner) *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Create the planning repository of the repository whose wrapper is the current directory",
		Long: "Create the local planning repository <repo-id>-plan in the current directory, as the last step of ply repo create does.\n" +
			"The current directory must be the wrapper of exactly one repository in ply-workspace.json; its ID names the\n" +
			"planning repository. The planning repository is made from the template ~/.agents/templates/plan-repo, committed\n" +
			"once on main and has no remote. An existing planning repository is left unchanged; a conflict is never overwritten.",
		Example: "  cd ply_projects/ply_structurizr && ply repo plan",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return fmt.Errorf("%w: JSON output is not supported", repoestablish.ErrInvalid)
			}
			observed, err := workspace.ObserveContaining(workspace.SystemDependencies())
			if err != nil {
				return err
			}
			cwd, err := os.Getwd()
			if err == nil {
				cwd, err = filepath.EvalSymlinks(cwd)
			}
			if err != nil {
				return fmt.Errorf("%w: working directory: %v", repoestablish.ErrInvalid, err)
			}
			repo, err := repoAtWrapper(observed.Root, cwd)
			if err != nil {
				return err
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("%w: home directory for the planning template: %v", repoestablish.ErrInvalid, err)
			}
			changed, err := system.EnsurePlan(cmd.Context(), runner, filepath.Join(home, system.PlanTemplate), cwd, repo)
			if err != nil {
				return err
			}
			outcome := "unchanged"
			if changed {
				outcome = "created"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Planning repository %s: %s\n", wrapperplan.Dir(cwd, repo), outcome)
			return err
		},
	}
}

// repoAtWrapper returns the ID of the one defined repository whose wrapper is
// the physical directory cwd.
func repoAtWrapper(root, cwd string) (string, error) {
	defs, _, err := workspacedef.Read(root)
	if err != nil {
		return "", err
	}
	if defs == nil {
		return "", fmt.Errorf("%w: %s is missing in %s", repoestablish.ErrInvalid, workspacedef.File, root)
	}
	var ids []string
	for _, p := range defs.Projects {
		for _, r := range p.Repos {
			if r.Wrapper != "" && workspace.ResolveDefinitionPath(root, r.Wrapper) == cwd {
				ids = append(ids, r.ID)
			}
		}
	}
	if len(ids) != 1 {
		return "", fmt.Errorf("%w: %d repositories (%s) have wrapper %s; run from the wrapper of exactly one repository",
			repoestablish.ErrInvalid, len(ids), strings.Join(ids, ", "), cwd)
	}
	return ids[0], nil
}
