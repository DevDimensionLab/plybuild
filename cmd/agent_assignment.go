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
	"github.com/spf13/cobra"
)

func newAgentAssignmentCommand(start func(context.Context, agentlaunch.Request) (agentlaunch.Result, error)) *cobra.Command {
	var project, cwd, workspace, label, claude, herdr string
	var timeout time.Duration
	command := &cobra.Command{
		Use: "assignment ASSIGNMENT-ID", Short: "Start Claude with the Assignment skill in a new visible Herdr tab",
		Long: "Start one Claude process with /ply-assignment as its initial prompt. The skill reads and takes the Assignment. No Task binding or durable launch receipt is created. Never retry an uncertain launch without inspecting the returned pane. The pilot supports sh/bash/zsh; foreground shell observation does not prove shell initialization has finished.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !regexp.MustCompile(`^ASSIGNMENT-[1-9][0-9]*$`).MatchString(args[0]) {
				return fmt.Errorf("expected ASSIGNMENT-N")
			}
			if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`).MatchString(project) {
				return fmt.Errorf("--project requires an explicit project identifier")
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
			executable, err := exec.LookPath(claude)
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
			result, startErr := start(cmd.Context(), agentlaunch.Request{Herdr: runtime, Executable: executable, Cwd: physical, Workspace: workspace, Label: label, InitialPrompt: "/ply-assignment " + args[0] + " --project " + project, Timeout: timeout})
			if err = json.NewEncoder(cmd.OutOrStdout()).Encode(result); err != nil {
				return err
			}
			return startErr
		},
	}
	command.Flags().StringVar(&project, "project", "", "Explicit Ply project identifier")
	command.Flags().StringVar(&cwd, "cwd", ".", "Existing working directory; no worktree is created")
	command.Flags().StringVar(&workspace, "herdr-workspace", "", "Explicit Herdr workspace ID")
	command.Flags().StringVar(&label, "label", "", "New tab label (defaults to Assignment ID)")
	command.Flags().StringVar(&claude, "claude-bin", "claude", "Claude executable")
	command.Flags().StringVar(&herdr, "herdr-bin", "herdr", "Herdr executable")
	command.Flags().DurationVar(&timeout, "timeout", 45*time.Second, "Total launch/observation budget")
	_ = command.MarkFlagRequired("project")
	_ = command.MarkFlagRequired("herdr-workspace")
	return command
}
