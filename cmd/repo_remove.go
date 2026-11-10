package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/repoestablish/system"
	"github.com/devdimensionlab/plybuild/internal/reporemove"
	removesystem "github.com/devdimensionlab/plybuild/internal/reporemove/system"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// repoRemoveEnv is the caller context of ply repo remove: working directory,
// environment and the operation.
type repoRemoveEnv struct {
	files  workspace.FileSystem
	getenv func(string) string
	run    func(context.Context, reporemove.Request) (reporemove.Result, error)
}

func systemRepoRemoveEnv(runner process.Runner) repoRemoveEnv {
	return repoRemoveEnv{files: workspace.SystemDependencies().Files, getenv: os.Getenv,
		run: func(ctx context.Context, req reporemove.Request) (reporemove.Result, error) {
			return reporemove.Run(ctx, req, removesystem.New(runner))
		}}
}

func newRepoRemoveCommand(env repoRemoveEnv) *cobra.Command {
	var req reporemove.Request
	command := &cobra.Command{
		Use:   "remove <repo-name>",
		Short: "Remove a Project repository's local wrapper, definition and metadata without changing its remote",
		Long: "Undo the local part of 'ply repo create <repo-name>': clear the wrapper metadata, remove the repository from\n" +
			"ply-workspace.json and delete the wrapper directory <Project root>/<repo-name>, only when the checks prove that no\n" +
			"local data is lost. The remote repository is never changed; it is contacted by one read-only ls-remote only when\n" +
			"the repository has local history. Without --apply the command previews: it writes nothing and prints the checks,\n" +
			"the planned effects and a plan digest. --apply --expect <digest> repeats every check and removes only when the\n" +
			"plan is unchanged. Project, workspace and actor resolve as for 'ply repo create'; --remote defaults to the\n" +
			"Project defaults' remote, else origin.",
		Example: "  ply repo remove ply-code-map\n  PLY_ACTOR=ply-epic ply repo remove ply-code-map --apply --expect sha256:<digest>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return fmt.Errorf("%w: JSON output is not supported", reporemove.ErrInvalid)
			}
			if req.Apply != (req.Expect != "") {
				return fmt.Errorf("%w: --apply and --expect <plan digest> must be given together", reporemove.ErrInvalid)
			}
			if req.Workspace == "" {
				observed, err := workspace.ObserveContainingRoot(workspace.Dependencies{Files: env.files})
				if err != nil {
					return err
				}
				req.Workspace = observed.Root
			} else if !filepath.IsAbs(req.Workspace) {
				return fmt.Errorf("%w: --workspace must be an absolute path", reporemove.ErrInvalid)
			}
			if req.Project == "" {
				cwd, err := env.files.Getwd()
				if err == nil {
					cwd, err = env.files.EvalSymlinks(cwd)
				}
				if err != nil {
					return fmt.Errorf("%w: working directory: %v", reporemove.ErrInvalid, err)
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
			r, plan := result.Resolved, result.Plan
			if r.Wrapper != "" {
				report("Resolved Project %s: root %s, wrapper %s, repository %s, repository ID %s, remote %s\n",
					r.Project, r.ProjectRoot, r.Wrapper, r.RepoDir, r.RepoName, r.Remote)
			}
			switch {
			case result.Nothing:
				report("Nothing to remove: repository %s is not defined in Project %s and wrapper %s is absent.\n", r.RepoName, r.Project, r.Wrapper)
			case err == nil && !req.Apply:
				for _, check := range plan.Checks {
					if check.Detail == "" {
						report("check %s: ok\n", check.Name)
					} else {
						report("check %s: ok (%s)\n", check.Name, check.Detail)
					}
				}
				metadata := fmt.Sprintf("clear wrapper metadata (revision %d)", plan.MetadataRevision)
				if plan.MetadataCleared {
					metadata = fmt.Sprintf("wrapper metadata already cleared (revision %d)", plan.MetadataRevision)
				}
				report("plan: %s, remove definitions member %s, delete wrapper %s\n", metadata, r.RepoName, r.Wrapper)
				if len(plan.GitFiles) > 0 {
					report("local Git files deleted with the repository: %s\n", strings.Join(plan.GitFiles, ", "))
				}
				report("plan digest: %s\n", plan.Digest)
			}
			for _, step := range result.Steps {
				outcome := "unchanged"
				if step.Changed {
					outcome = "removed"
				}
				report("%s: %s\n", step.Step, outcome)
			}
			if plan.RemoteURL != "" {
				report("Remote repository %s is not changed.\n", plan.RemoteURL)
			} else {
				report("The remote repository is not changed.\n")
			}
			if writeErr != nil {
				return errors.Join(err, writeErr)
			}
			return err
		},
	}
	command.Flags().StringVar(&req.Project, "project", "", "Project ID (default: the Project whose root is the current directory)")
	command.Flags().StringVar(&req.Workspace, "workspace", "", "absolute workspace root (default: the workspace containing the current directory)")
	command.Flags().StringVar(&req.Remote, "remote", "", "remote name (default: Project defaults, else origin)")
	command.Flags().StringVar(&req.Actor, "actor", "", "actor recorded in the wrapper metadata clear (default: PLY_ACTOR)")
	command.Flags().BoolVar(&req.Apply, "apply", false, "remove after repeating every check; requires --expect")
	command.Flags().StringVar(&req.Expect, "expect", "", "plan digest printed by the preview")
	return command
}
