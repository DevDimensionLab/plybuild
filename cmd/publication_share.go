package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/agentlaunch"
	"github.com/devdimensionlab/plybuild/internal/publication"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// shareSkillMention starts the share skill; only Claude has the Artifact tool.
const shareSkillMention = "/ply-publication-share"

// shareActor is PLY_ACTOR of the agent that shares.
const shareActor = "claude:publication-share"

func addShareCommands(command *cobra.Command, addFlags func(*cobra.Command) *cobra.Command, open func() (publication.Store, error),
	resolved func() workspace.ProjectRecord, d workspace.Dependencies, start func(context.Context, agentlaunch.Request) (agentlaunch.Result, error), format *string) {
	var path publication.Path
	var run, actor, herdrWorkspace, claude, herdr string
	var timeout time.Duration
	pathFlags := func(c *cobra.Command) {
		c.Flags().StringVar(&path.Repo, "repo", "", "repo ID (default: the whole Project)")
		c.Flags().StringVar(&path.Module, "module", "", "module path in the repo (default: the repo itself)")
		c.Flags().StringVar(&path.Category, "category", "", "category (required)")
		c.Flags().StringVar(&path.Agent, "agent", "", "agent of the path")
	}
	readback := func(cmd *cobra.Command, share publication.Share) error {
		if *format == "json" {
			return writeAssignmentJSON(cmd.OutOrStdout(), share)
		}
		_, err := io.WriteString(cmd.OutOrStdout(), shareText(path, share))
		return err
	}

	shareCommand := addFlags(&cobra.Command{
		Use:   "share --category <category> [--run <run>] --herdr-workspace <id>",
		Short: "Share one run as the path's Claude artifact through a Claude agent in a new Herdr tab",
		Long: "Reserve the share of the path in share.json, then start one Claude in a new visible Herdr tab in the run's directory with /ply-publication-share as its first prompt.\n" +
			"The agent publishes the run as a private Claude artifact and records the link with share-result; each path keeps one fixed link that a share of a newer run updates.\n" +
			"A run the path already shows is returned without an agent. An attempt that is pending or unknown is refused until share-result ends it. A failed or unknown start is recorded and never retried.\n" +
			"Without --run the newest run is shared. The caller must run in Herdr. The actor is --actor or PLY_ACTOR.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			project := resolved()
			if actor == "" {
				actor = os.Getenv("PLY_ACTOR")
			}
			v, request, err := prepareShare(s, d, project, path, run, herdrWorkspace, claude, herdr)
			if err != nil {
				return err
			}
			request.Timeout = timeout
			share, started, err := s.ReserveShare(string(project.ID), path, v.Run, actor)
			if err != nil {
				return err
			}
			if !started {
				return readback(cmd, share)
			}
			result, launchErr := start(cmd.Context(), request)
			state, detail, session := launchOutcome(result, launchErr)
			if share, err = s.FinishShareLaunch(path, v.Run, session, state, detail); err != nil {
				return err
			}
			if err := readback(cmd, share); err != nil {
				return err
			}
			if launchErr != nil {
				return &publication.Error{Class: publication.ClassShareLaunchFailed, Detail: "the start of the share agent is " + state}
			}
			return nil
		},
	})
	pathFlags(shareCommand)
	shareCommand.Flags().StringVar(&run, "run", "", "run ID (default: the newest)")
	shareCommand.Flags().StringVar(&actor, "actor", "", "who asks for the share (default: PLY_ACTOR)")
	shareCommand.Flags().StringVar(&herdrWorkspace, "herdr-workspace", "", "Herdr workspace ID of the new tab (required)")
	shareCommand.Flags().StringVar(&claude, "claude-bin", "claude", "Claude executable")
	shareCommand.Flags().StringVar(&herdr, "herdr-bin", "herdr", "Herdr executable")
	shareCommand.Flags().DurationVar(&timeout, "timeout", 45*time.Second, "total budget for starting the agent")

	var url, failure string
	resultCommand := addFlags(&cobra.Command{
		Use:   "share-result --category <category> --run <run> (--url <url> | --error <text>)",
		Short: "Record the artifact link or the failure of an attempt to share a run",
		Long: "End the pending or unknown attempt to share --run: --url records the artifact link (https://claude.ai/artifact/<id>), which must be the path's fixed link once it has one; --error records a failure of one line and keeps the previous link.\n" +
			"Recording the same result again returns the share unchanged. Used by the share agent, or by a person to end an attempt that will not finish.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := open()
			if err != nil {
				return err
			}
			if path.Repo != "" && !slices.Contains(resolved().RepoIDs, workspace.RepoID(path.Repo)) {
				return publicationInvalidArguments("--repo must be a repo ID of project " + string(resolved().ID))
			}
			share, err := s.RecordShareResult(path, run, url, failure)
			if err != nil {
				return err
			}
			return readback(cmd, share)
		},
	})
	pathFlags(resultCommand)
	resultCommand.Flags().StringVar(&run, "run", "", "run ID of the attempt (required)")
	resultCommand.Flags().StringVar(&url, "url", "", "artifact link")
	resultCommand.Flags().StringVar(&failure, "error", "", "why the share failed, one line")
}

// prepareShare checks everything a share needs before any effect and returns
// the run to share with the start request of its agent.
func prepareShare(s publication.Store, d workspace.Dependencies, project workspace.ProjectRecord, path publication.Path, run, herdrWorkspace, claude, herdr string) (publication.View, agentlaunch.Request, error) {
	var request agentlaunch.Request
	if path.Repo != "" && !slices.Contains(project.RepoIDs, workspace.RepoID(path.Repo)) {
		return publication.View{}, request, publicationInvalidArguments("--repo must be a repo ID of project " + string(project.ID))
	}
	if herdrWorkspace == "" {
		return publication.View{}, request, publicationInvalidArguments("--herdr-workspace is required")
	}
	if os.Getenv("HERDR_ENV") != "1" {
		return publication.View{}, request, publicationInvalidArguments("the caller must run in Herdr (HERDR_ENV=1)")
	}
	v, err := s.Show(path.Repo, path.Module, path.Category, path.Agent, run)
	if err != nil {
		return v, request, err
	}
	observed, err := workspace.ObserveContaining(d)
	if err != nil {
		return v, request, err
	}
	cwd, err := filepath.EvalSymlinks(v.Path)
	if err != nil {
		return v, request, &publication.Error{Class: publication.ClassIOError, Detail: "could not resolve the run directory"}
	}
	executable, err := absoluteExecutable(claude)
	if err != nil {
		return v, request, publicationInvalidArguments("--claude-bin is not an executable")
	}
	runtime, err := absoluteExecutable(herdr)
	if err != nil {
		return v, request, publicationInvalidArguments("--herdr-bin is not an executable")
	}
	return v, agentlaunch.Request{Agent: "claude", Herdr: runtime, Executable: executable, Cwd: cwd, Workspace: herdrWorkspace,
		Label: "share " + publicationPathSubject(path), Env: []string{"PLY_ACTOR=" + shareActor},
		InitialPrompt: sharePrompt(string(project.ID), observed.Root, path, v.Run)}, nil
}

func absoluteExecutable(name string) (string, error) {
	executable, err := exec.LookPath(name)
	if err != nil {
		return "", err
	}
	return filepath.Abs(executable)
}

func sharePrompt(project, root string, p publication.Path, run string) string {
	prompt := shareSkillMention + " --project " + project + " --workspace " + root
	if p.Repo != "" {
		prompt += " --repo " + p.Repo
	}
	if p.Module != "" && p.Module != publication.RootModule {
		prompt += " --module " + p.Module
	}
	prompt += " --category " + p.Category
	if p.Agent != "" {
		prompt += " --agent " + p.Agent
	}
	return prompt + " --run " + run
}

// launchOutcome maps the start of the agent to the state of the attempt.
func launchOutcome(result agentlaunch.Result, err error) (string, string, *publication.ShareSession) {
	var session *publication.ShareSession
	if result.Pane != "" {
		session = &publication.ShareSession{Provider: "claude", Pane: result.Pane, Tab: result.Tab, ID: result.Session}
	}
	switch {
	case err == nil:
		return publication.ShareStatePending, "", session
	case result.EffectUnknown:
		return publication.ShareStateUnknown, err.Error(), session
	default:
		return publication.ShareStateFailed, err.Error(), session
	}
}

func shareText(p publication.Path, share publication.Share) string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s share %s of %s\n", publicationPathSubject(p), share.State, share.Run)
	if share.URL != nil {
		fmt.Fprintf(&out, "url: %s", *share.URL)
		if share.SharedRun != nil {
			fmt.Fprintf(&out, " (shows %s)", *share.SharedRun)
		}
		out.WriteString("\n")
	}
	if share.Error != nil {
		fmt.Fprintf(&out, "error: %s\n", *share.Error)
	}
	if share.Session != nil {
		fmt.Fprintf(&out, "agent: %s pane %s\n", share.Session.Provider, share.Session.Pane)
	}
	return out.String()
}
