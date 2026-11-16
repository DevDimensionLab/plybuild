package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/devdimensionlab/plybuild/internal/agentlaunch"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

// assignmentSkillMention is how each supported provider invokes a skill from its initial prompt.
var assignmentSkillMention = map[string]string{"claude": "/ply-assignment", "codex": "$ply-assignment"}

func newAgentAssignmentCommand(start func(context.Context, agentlaunch.Request) (agentlaunch.Result, error)) *cobra.Command {
	var provider, project, cwd, herdrWorkspace, label, claude, codex, herdr string
	var timeout time.Duration
	command := &cobra.Command{
		Use: "assignment ASSIGNMENT-ID", Short: "Start Claude or Codex with the Assignment skill in a new visible Herdr tab",
		Long: "Start one Claude (/ply-assignment) or Codex ($ply-assignment) process with the Assignment skill as its initial prompt. The skill reads and takes the Assignment in the workspace containing the current directory, passed as --workspace, so --cwd may be outside it. No Task binding or durable launch receipt is created. Never retry an uncertain launch without inspecting the returned pane. The pilot supports sh/bash/zsh; foreground shell observation does not prove shell initialization has finished.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !regexp.MustCompile(`^ASSIGNMENT-[1-9][0-9]*$`).MatchString(args[0]) {
				return fmt.Errorf("expected ASSIGNMENT-N")
			}
			mention, ok := assignmentSkillMention[provider]
			if !ok {
				return fmt.Errorf("--provider must be claude or codex")
			}
			if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`).MatchString(project) {
				return fmt.Errorf("--project requires an explicit project identifier")
			}
			observed, err := workspace.ObserveContaining(workspace.SystemDependencies())
			if err != nil {
				return err
			}
			physical, err := filepath.Abs(cwd)
			if err == nil {
				physical, err = filepath.EvalSymlinks(physical)
			}
			if err != nil {
				return err
			}
			stat, err := os.Stat(physical)
			if err != nil {
				return err
			}
			if !stat.IsDir() {
				return fmt.Errorf("--cwd must be a directory")
			}
			binary := claude
			if provider == "codex" {
				binary = codex
			}
			executable, err := exec.LookPath(binary)
			if err != nil {
				return err
			}
			executable, err = filepath.Abs(executable)
			if err != nil {
				return err
			}
			runtime, err := exec.LookPath(herdr)
			if err != nil {
				return err
			}
			runtime, err = filepath.Abs(runtime)
			if err != nil {
				return err
			}
			if label == "" {
				label = args[0]
			}
			result, startErr := start(cmd.Context(), agentlaunch.Request{Agent: provider, Herdr: runtime, Executable: executable, Cwd: physical, Workspace: herdrWorkspace, Label: label, InitialPrompt: mention + " " + args[0] + " --project " + project + " --workspace " + observed.Root, Timeout: timeout})
			if err = json.NewEncoder(cmd.OutOrStdout()).Encode(result); err != nil {
				return err
			}
			return startErr
		},
	}
	command.Flags().StringVar(&provider, "provider", "claude", "Agent provider: claude or codex")
	command.Flags().StringVar(&project, "project", "", "Explicit Ply project identifier")
	command.Flags().StringVar(&cwd, "cwd", ".", "Existing working directory; no worktree is created")
	command.Flags().StringVar(&herdrWorkspace, "herdr-workspace", "", "Explicit Herdr workspace ID")
	command.Flags().StringVar(&label, "label", "", "New tab label (defaults to Assignment ID)")
	command.Flags().StringVar(&claude, "claude-bin", "claude", "Claude executable")
	command.Flags().StringVar(&codex, "codex-bin", "codex", "Codex executable")
	command.Flags().StringVar(&herdr, "herdr-bin", "herdr", "Herdr executable")
	command.Flags().DurationVar(&timeout, "timeout", 45*time.Second, "Total launch/observation budget")
	_ = command.MarkFlagRequired("project")
	_ = command.MarkFlagRequired("herdr-workspace")
	return command
}
