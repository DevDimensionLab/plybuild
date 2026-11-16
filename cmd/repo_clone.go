package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/gitrepo"
	"github.com/devdimensionlab/plybuild/internal/repoestablish/system"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
	"github.com/spf13/cobra"
)

const repoCloneTimeout = 10 * time.Minute

func newRepoCloneCommand(runner process.Runner) *cobra.Command {
	var project string
	command := &cobra.Command{
		Use:   "clone [<repo-id>...]",
		Short: "Clone the repositories of a Project that are missing locally",
		Long: "Clone each repository of the Project (or only the named ones) whose path in ply-workspace.json does not exist,\n" +
			"from its remote in ply-workspace.json. A repository whose path exists is reported as present and left as it is;\n" +
			"one without remote is skipped. The Project is the one whose root is the current directory unless --project is given.\n" +
			"A failed clone is reported, its partial directory removed, and the others still run; the command then fails.",
		Example: "  ply repo clone --project ply\n  ply repo clone --project ply ply-code-map",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return fmt.Errorf("%w: JSON output is not supported", gitrepo.ErrInvalid)
			}
			observed, err := workspace.ObserveContaining(workspace.SystemDependencies())
			if err != nil {
				return err
			}
			root := observed.Root
			if project == "" {
				cwd, err := os.Getwd()
				if err == nil {
					cwd, err = filepath.EvalSymlinks(cwd)
				}
				if err != nil {
					return fmt.Errorf("%w: working directory: %v", gitrepo.ErrInvalid, err)
				}
				if project, err = system.ProjectAtRoot(root, cwd); err != nil {
					return err
				}
			}
			defs, _, err := workspacedef.Read(root)
			if err != nil {
				return err
			}
			if defs == nil {
				return fmt.Errorf("%w: %s is missing in %s", gitrepo.ErrInvalid, workspacedef.File, root)
			}
			index := slices.IndexFunc(defs.Projects, func(p workspacedef.Project) bool { return p.ID == project })
			if index < 0 {
				return fmt.Errorf("%w: Project %s is not defined", gitrepo.ErrInvalid, project)
			}
			repos := defs.Projects[index].Repos
			for _, name := range args {
				if !slices.ContainsFunc(repos, func(r workspacedef.Repo) bool { return r.ID == name }) {
					return fmt.Errorf("%w: %s is not a repository of Project %s", gitrepo.ErrInvalid, name, project)
				}
			}
			var failed []string
			for _, repo := range repos {
				if len(args) > 0 && !slices.Contains(args, repo.ID) {
					continue
				}
				outcome := cloneRepo(cmd.Context(), runner, root, repo)
				if strings.HasPrefix(outcome, "failed") {
					failed = append(failed, repo.ID)
				}
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", repo.ID, outcome); err != nil {
					return err
				}
			}
			if len(failed) > 0 {
				return fmt.Errorf("%w: clone failed for %s", gitrepo.ErrUnknown, strings.Join(failed, ", "))
			}
			return nil
		},
	}
	command.Flags().StringVar(&project, "project", "", "Project ID (default: the Project whose root is the current directory)")
	return command
}

// cloneRepo clones one repository into its path and returns the outcome line.
func cloneRepo(ctx context.Context, runner process.Runner, root string, repo workspacedef.Repo) string {
	target := filepath.Join(root, filepath.FromSlash(repo.Path))
	if _, err := os.Lstat(target); err == nil {
		return "present"
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Sprintf("failed: inspect %s: %v", repo.Path, err)
	}
	if repo.Remote == "" {
		return "skipped, no remote"
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Sprintf("failed: create parent of %s: %v", repo.Path, err)
	}
	ctx, cancel := context.WithTimeout(ctx, repoCloneTimeout)
	defer cancel()
	var stderr strings.Builder
	err := runner.Run(process.Command{Name: "git", Args: []string{"clone", "--quiet", "--", repo.Remote, target}, Context: ctx, Dir: root, Stderr: &stderr})
	if err != nil {
		_ = os.RemoveAll(target)
		return fmt.Sprintf("failed: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	return "cloned into " + repo.Path
}
