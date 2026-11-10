package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/agentstart"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

const goalAgentPrompt = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. The Goal decides return, notification, delivery and delegation; this start grants no authority beyond it and does not attest runtime permissions. Goal:\n"

type goalAgentPlan struct {
	Task      string          `json:"task"`
	Worktree  string          `json:"worktree"`
	Ref       string          `json:"ref"`
	OID       string          `json:"oid"`
	Workspace string          `json:"herdr_workspace"`
	Provider  string          `json:"provider,omitempty"`
	Goal      json.RawMessage `json:"goal"`
}

func goalProvider(goal json.RawMessage) (string, error) {
	var g struct {
		Executor struct {
			Provider *string `json:"provider"`
			Model    *string `json:"model"`
			Effort   *string `json:"effort"`
		} `json:"executor"`
	}
	if err := json.Unmarshal(goal, &g); err != nil {
		return "", err
	}
	provider := agentstart.ProviderCodex
	if g.Executor.Provider != nil {
		provider = *g.Executor.Provider
	}
	if provider != agentstart.ProviderClaude && provider != agentstart.ProviderCodex {
		return "", fmt.Errorf("goal-start supports Claude and Codex Goals only")
	}
	if g.Executor.Model != nil || g.Executor.Effort != nil {
		return "", fmt.Errorf("explicit model/effort overrides require a supported launcher; not silently discarded")
	}
	return provider, nil
}

func goalStartInput(plan goalAgentPlan, receipt string, herdr agentstart.HerdrBinary, precondition func() error) agentstart.StartInput {
	return agentstart.StartInput{
		Key:          plan.Task,
		ReceiptPath:  receipt,
		Provider:     plan.Provider,
		Cwd:          plan.Worktree,
		Workspace:    plan.Workspace,
		Label:        plan.Task,
		Prompt:       goalAgentPrompt + string(plan.Goal),
		Herdr:        herdr,
		Precondition: precondition,
	}
}

func startGoal(ctx context.Context, plan goalAgentPlan, receipt string, herdr agentstart.HerdrBinary, precondition func() error) (agentstart.Receipt, error) {
	return agentstart.StartAgent(ctx, goalStartInput(plan, receipt, herdr, precondition))
}

func resumeGoal(ctx context.Context, plan goalAgentPlan, receipt string, herdr agentstart.HerdrBinary, precondition func() error) (agentstart.Receipt, error) {
	return agentstart.ResumeAgent(ctx, goalStartInput(plan, receipt, herdr, precondition))
}

func checkGoalGit(wt *workspace.TaskWorktreeBinding) error {
	git := func(a ...string) (string, error) {
		b, e := exec.Command("git", append([]string{"--no-optional-locks", "-C", wt.Locator}, a...)...).Output()
		return string(b), e
	}
	status, e := git("status", "--porcelain")
	if e != nil {
		return e
	}
	if status != "" {
		return fmt.Errorf("Task worktree must be clean before initial start")
	}
	ref, e := git("symbolic-ref", "HEAD")
	if e != nil {
		return e
	}
	if ref != wt.Ref+"\n" {
		return fmt.Errorf("Task ref changed before initial start")
	}
	oid, e := git("rev-parse", "HEAD")
	if e != nil {
		return e
	}
	if oid != wt.OID+"\n" {
		return fmt.Errorf("Task HEAD changed before initial start")
	}
	return nil
}

func newWorkflowGoalStartCommand(d workspace.Dependencies) *cobra.Command {
	var spec, herdrWorkspace string
	var revision int
	var apply, resume bool
	c := &cobra.Command{Use: "goal-start TASK_ID", Short: "Start an explicitly assigned Goal in its bound Task worktree", Long: "Preview a registered Goal and existing Task worktree. --apply starts one visible Claude or Codex agent, chosen by the Goal executor provider (Codex when absent), with inherited runtime policy. A durable reservation prevents duplicate launches after failure or unknown effects; inspect its receipt and pane. When the start stopped before its prompt (for example on folder trust), settle the pane yourself and rerun with --apply --resume: it sends the same Goal prompt at most once to the same agent and never resends an unknown prompt.\n\nKoordinatorretur velges i Goal.constraints; se docs/goal-coordinator-return.md. Goal-start oppretter ingen koordinator og krever ingen retur.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		if resume && !apply {
			return fmt.Errorf("--resume requires --apply")
		}
		id, e := workspace.ParseTaskID(args[0])
		if e != nil {
			return e
		}
		task, e := workspace.ShowTask(d, id)
		if e != nil {
			return e
		}
		if e = validateGoalTask(task); e != nil {
			return e
		}
		content, e := workspace.ShowTaskSpec(d, workspace.TaskContentQuery{TaskID: id, SpecID: spec, Revision: revision})
		if e != nil {
			return e
		}
		raw, e := canonicaljson.Marshal(content.Value)
		if e != nil {
			return e
		}
		var read struct {
			Goal    json.RawMessage `json:"goal"`
			Reasons []string        `json:"reasons"`
		}
		if e = json.Unmarshal(raw, &read); e != nil {
			return e
		}
		if len(read.Goal) == 0 || string(read.Goal) == "null" || len(read.Reasons) > 0 {
			return fmt.Errorf("registered revision must be an intact Goal")
		}
		provider, e := goalProvider(read.Goal)
		if e != nil {
			return e
		}
		wt := task.Task.Worktree
		if e = checkGoalGit(wt); e != nil {
			return e
		}
		if herdrWorkspace == "" {
			return fmt.Errorf("--herdr-workspace is required")
		}
		plan := goalAgentPlan{Task: args[0], Worktree: wt.Locator, Ref: wt.Ref, OID: wt.OID, Workspace: herdrWorkspace, Provider: provider, Goal: read.Goal}
		if !apply {
			return json.NewEncoder(c.OutOrStdout()).Encode(plan)
		}
		if os.Getenv("HERDR_ENV") != "1" {
			return fmt.Errorf("goal-start --apply requires the actual Herdr terminal")
		}
		herdr, e := agentstart.ResolveHerdr()
		if e != nil {
			return e
		}
		run := startGoal
		if resume {
			run = resumeGoal
		}
		out, e := run(c.Context(), plan, filepath.Join(task.Workspace, ".ply", "task-process", "goal-start", args[0]+".json"), herdr, func() error { return checkGoalGit(wt) })
		if writeErr := json.NewEncoder(c.OutOrStdout()).Encode(out); writeErr != nil {
			return writeErr
		}
		return e
	}}
	c.Flags().StringVar(&spec, "spec", "", "registered Goal Spec ID")
	c.Flags().IntVar(&revision, "revision", 0, "registered Goal revision")
	c.Flags().StringVar(&herdrWorkspace, "herdr-workspace", "", "explicit Herdr workspace")
	c.Flags().BoolVar(&apply, "apply", false, "start the explicit Goal once")
	c.Flags().BoolVar(&resume, "resume", false, "with --apply, continue a start stopped before its prompt in the same agent")
	return c
}

func validateGoalTask(task workspace.TaskReadbackResult) error {
	if task.Task.Worktree == nil || !task.WorktreeReady {
		return fmt.Errorf("create this Task worktree explicitly before agent start")
	}
	wt, target := task.Task.Worktree, task.Target
	if task.ProjectFreshness != "fresh" || target.InventoryMatch == nil || !*target.InventoryMatch || target.Clean == nil || !*target.Clean {
		return fmt.Errorf("Task Git identity or inventory is stale or unknown")
	}
	match := func(got *string, want string) bool { return got != nil && *got == want }
	if !match(target.Locator, wt.Locator) || !match(target.Ref, wt.Ref) || !match(target.OID, wt.OID) || !match(target.Tree, wt.Tree) || !match(target.GitCommonDir, wt.GitCommonDir) || wt.GitCommonDir != task.Task.GitCommonDir {
		return fmt.Errorf("Task Git binding differs from the observed worktree")
	}
	return nil
}
