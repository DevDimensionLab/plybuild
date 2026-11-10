package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/internal/githubrepo"
	"github.com/devdimensionlab/plybuild/internal/gitrepo"
	"github.com/spf13/cobra"
)

const (
	repoInitTimeout = 30 * time.Second
	githubTimeout   = 60 * time.Second
)

func newRepoCommand(runner process.Runner) *cobra.Command {
	repoCommand := &cobra.Command{
		Use:   "repo",
		Short: "Manage local Git repositories",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	var branch string
	initCommand := &cobra.Command{
		Use:     "init <dir> --branch <name>",
		Short:   "Initialize an empty or missing directory as a Git repository",
		Long:    "Initialize an empty or missing directory as a Git repository on the given branch.\nAn existing repository on the same branch is left unchanged; anything else is a conflict.",
		Example: "  ply repo init ./service --branch main",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return fmt.Errorf("%w: JSON output is not supported", gitrepo.ErrInvalid)
			}
			dir, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("%w: %v", gitrepo.ErrInvalid, err)
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), repoInitTimeout)
			defer cancel()
			repo, changed, err := gitrepo.EnsureInit(ctx, runner, dir, branch)
			if err != nil {
				return err
			}
			format := "Git repository already initialized at %s on branch %s.\n"
			if changed {
				format = "Initialized Git repository at %s on branch %s.\n"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), format, repo.Toplevel, repo.Branch)
			return err
		},
	}
	initCommand.Flags().StringVar(&branch, "branch", "", "branch name for the repository (required)")
	_ = initCommand.MarkFlagRequired("branch")
	repoCommand.AddCommand(initCommand)
	remoteCommand := &cobra.Command{Use: "remote", Short: "Manage local Git remotes"}
	remoteCommand.AddCommand(&cobra.Command{
		Use:     "set <dir> <name> <url>",
		Short:   "Add a local Git remote, or verify that it is already set",
		Long:    "Add the named remote with the URL to the Git repository rooted at dir.\nA remote with exactly that local URL is left unchanged; anything else is a conflict and is never overwritten.",
		Example: "  ply repo remote set ./service origin git@github.com:org/service.git",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return fmt.Errorf("%w: JSON output is not supported", gitrepo.ErrInvalid)
			}
			dir, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("%w: %v", gitrepo.ErrInvalid, err)
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), repoInitTimeout)
			defer cancel()
			changed, err := gitrepo.EnsureRemote(ctx, runner, dir, args[1], args[2])
			if err != nil {
				return err
			}
			format := "Git remote %s already set in %s.\n"
			if changed {
				format = "Added Git remote %s to %s.\n"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), format, args[1], dir)
			return err
		},
	})
	repoCommand.AddCommand(remoteCommand)
	githubCommand := &cobra.Command{Use: "github", Short: "Manage GitHub repositories"}
	githubCommand.AddCommand(&cobra.Command{
		Use:     "ensure <owner>/<name>",
		Short:   "Create a private empty GitHub repository, or verify that it already is one",
		Long:    "Create the private github.com repository owner/name once when it is not visible, or verify that it is already private and empty.\nThe canonical SSH URL is printed on stdout; any other existing repository is a conflict and is never changed.",
		Example: "  ply repo remote set ./service origin \"$(ply repo github ensure org/service)\"",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return fmt.Errorf("%w: JSON output is not supported", githubrepo.ErrInvalid)
			}
			owner, name, ok := strings.Cut(args[0], "/")
			if !ok || strings.Contains(name, "/") {
				return fmt.Errorf("%w: expected <owner>/<name>", githubrepo.ErrInvalid)
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), githubTimeout)
			defer cancel()
			repo, created, err := githubrepo.EnsurePrivate(ctx, runner, owner, name)
			if err != nil {
				return err
			}
			format := "GitHub repository %s is already private and empty.\n"
			if created {
				format = "Created private GitHub repository %s.\n"
			}
			if _, err := fmt.Fprintf(cmd.ErrOrStderr(), format, args[0]); err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), repo.SSHURL)
			return err
		},
	})
	repoCommand.AddCommand(githubCommand)
	repoCommand.AddCommand(newRepoCreateCommand(systemRepoCreateEnv(runner)))
	repoCommand.AddCommand(newRepoRemoveCommand(systemRepoRemoveEnv(runner)))
	return repoCommand
}

func init() {
	RootCmd.AddCommand(newRepoCommand(process.SystemRunner().Runner))
}
