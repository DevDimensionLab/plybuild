package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/repoestablish"
	"github.com/devdimensionlab/plybuild/internal/repoestablish/system"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// repoCreateNote is the known K6 boundary: only project metadata reads the definitions file.
const repoCreateNote = "Note: project metadata reads ply-workspace.json; 'ply workspace project show/list' and Task/queue still read the legacy registry and do not see this repository."

// repoCreateEnv is the caller context of ply repo create: working directory,
// environment and the operation.
type repoCreateEnv struct {
	files  workspace.FileSystem
	getenv func(string) string
	run    func(context.Context, repoestablish.Request) (repoestablish.Result, error)
}

func systemRepoCreateEnv(runner process.Runner) repoCreateEnv {
	return repoCreateEnv{files: workspace.SystemDependencies().Files, getenv: os.Getenv,
		run: func(ctx context.Context, req repoestablish.Request) (repoestablish.Result, error) {
			return repoestablish.Run(ctx, req, system.New(runner))
		}}
}

func newRepoCreateCommand(env repoCreateEnv) *cobra.Command {
	var req repoestablish.Request
	command := &cobra.Command{
		Use:   "create <repo-name>",
		Short: "Create a Project repository with its wrapper, Git, GitHub, definition and metadata",
		Long: "Establish the repository <repo-name> of a Project in fixed order: wrapper directory <Project root>/<repo-name>,\n" +
			"local Git repository, private GitHub repository, remote, workspace definition and wrapper metadata.\n" +
			"The command creates a private GitHub repository through gh when it is not visible.\n" +
			"The repository name is the repository ID, the wrapper name and the GitHub name. The Project is the one whose\n" +
			"root is the current directory unless --project is given; github_owner, leaf, branch and remote come from\n" +
			"ply-repo-defaults.json in the Project root unless overridden, else main, main and origin. The actor is --actor or PLY_ACTOR.\n" +
			"Completed steps are left unchanged on a rerun; a conflict is never overwritten.",
		Example: "  PLY_ACTOR=ply-epic ply repo create ply-code-map",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return fmt.Errorf("%w: JSON output is not supported", repoestablish.ErrInvalid)
			}
			if req.Workspace == "" {
				observed, err := workspace.ObserveContainingRoot(workspace.Dependencies{Files: env.files})
				if err != nil {
					return err
				}
				req.Workspace = observed.Root
			} else if !filepath.IsAbs(req.Workspace) {
				return fmt.Errorf("%w: --workspace must be an absolute path", repoestablish.ErrInvalid)
			}
			if req.Project == "" {
				cwd, err := env.files.Getwd()
				if err == nil {
					cwd, err = env.files.EvalSymlinks(cwd)
				}
				if err != nil {
					return fmt.Errorf("%w: working directory: %v", repoestablish.ErrInvalid, err)
				}
				if req.Project, err = system.ProjectAtRoot(req.Workspace, cwd); err != nil {
					return err
				}
			}
			if req.Actor == "" {
				req.Actor = env.getenv("PLY_ACTOR")
			}
			req.RepoName = args[0]
			result, err := env.run(cmd.Context(), req)
			var writeErr error
			report := func(format string, args ...interface{}) {
				if writeErr == nil {
					_, writeErr = fmt.Fprintf(cmd.OutOrStdout(), format, args...)
				}
			}
			if r := result.Resolved; r.Wrapper != "" {
				report("Resolved Project %s: root %s, wrapper %s, repository %s, owner %s, branch %s, remote %s, actor %s\n",
					r.Project, r.ProjectRoot, r.Wrapper, r.RepoDir, r.GitHubOwner, r.Branch, r.Remote, r.Actor)
			}
			for _, step := range result.Steps {
				outcome := "unchanged"
				if step.Changed {
					outcome = "created"
				}
				report("%s: %s\n", step.Step, outcome)
			}
			report("%s\n", repoCreateNote)
			if writeErr != nil {
				return errors.Join(err, writeErr)
			}
			return err
		},
	}
	command.Flags().StringVar(&req.Project, "project", "", "Project ID (default: the Project whose root is the current directory)")
	command.Flags().StringVar(&req.Workspace, "workspace", "", "absolute workspace root (default: the workspace containing the current directory)")
	command.Flags().StringVar(&req.Actor, "actor", "", "actor recorded in the wrapper metadata (default: PLY_ACTOR)")
	command.Flags().StringVar(&req.Overrides.GitHubOwner, "github-owner", "", "GitHub owner (default: github_owner in the Project defaults)")
	command.Flags().StringVar(&req.Overrides.Leaf, "leaf", "", "repository directory inside the wrapper (default: Project defaults, else main)")
	command.Flags().StringVar(&req.Overrides.Branch, "branch", "", "initial branch (default: Project defaults, else main)")
	command.Flags().StringVar(&req.Overrides.Remote, "remote", "", "remote name (default: Project defaults, else origin)")
	return command
}
