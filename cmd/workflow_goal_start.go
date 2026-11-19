package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/agentstart"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/taskactor"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
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
	Model     *string         `json:"model,omitempty"`
	Effort    *string         `json:"effort,omitempty"`
	Goal      json.RawMessage `json:"goal"`
	Documents string          `json:"document_instruction,omitempty"`
	Journal   string          `json:"journal_instruction,omitempty"`
}

// journalInstruction asks the agent of a Task linked to a plan to record its
// steps, so that ply plan trace shows them; a Task without link gets none.
func journalInstruction(task string, link *taskjournal.PlyAction) string {
	if link == nil || link.Plan == nil {
		return ""
	}
	belongs := *link.Plan
	if link.PlanStep != nil {
		belongs += " step " + *link.PlanStep
	}
	return fmt.Sprintf("\n\nThis Task belongs to %s. Record each work step in the Task journal: ply workspace task journal step start %s --id <id> --kind <kind> --title <title> when it starts, and ply workspace task journal step finish %s with the same --id, --kind and --title plus --outcome <outcome> when it ends. Use --parent <id> for a step inside another, such as a correction or review. Kinds: %s; final reporting and cleanup use kind coordination. Outcomes: %s. When you return, record ply workspace task journal note %s --title <title>. Ply sets PLY_ACTOR in this terminal, so these entries name you; do not pass --actor.",
		belongs, task, task, strings.Join(taskjournal.StepKinds, ", "), strings.Join(taskjournal.Outcomes, ", "), task)
}

// goalDocumentInstruction names the ply command that reads each document of
// the Goal Spec revision; a Goal without documents gets none.
func goalDocumentInstruction(task, spec string, revision int, goal json.RawMessage) (string, error) {
	var g struct {
		RequiredInputs []struct {
			ID string `json:"id"`
		} `json:"required_inputs"`
	}
	if err := json.Unmarshal(goal, &g); err != nil {
		return "", err
	}
	var reads []string
	for _, in := range g.RequiredInputs {
		if id, ok := strings.CutPrefix(in.ID, workspace.GoalDocumentInputID); ok {
			reads = append(reads, fmt.Sprintf("\n  ply workspace task spec show %s --spec %s --revision %d --document %s", task, spec, revision, id))
		}
	}
	if len(reads) == 0 {
		return "", nil
	}
	return "\n\nRead each Goal document with ply:" + strings.Join(reads, ""), nil
}

// goalExecutor is the registered Goal executor; model and effort are native
// values the provider's CLI receives at start.
type goalExecutor struct {
	Provider string
	Model    *string
	Effort   *string
}

func goalExecutorOf(goal json.RawMessage) (goalExecutor, error) {
	var g struct {
		Executor struct {
			Provider *string `json:"provider"`
			Model    *string `json:"model"`
			Effort   *string `json:"effort"`
		} `json:"executor"`
	}
	if err := json.Unmarshal(goal, &g); err != nil {
		return goalExecutor{}, err
	}
	provider := agentstart.ProviderCodex
	if g.Executor.Provider != nil {
		provider = *g.Executor.Provider
	}
	if provider != agentstart.ProviderClaude && provider != agentstart.ProviderCodex {
		return goalExecutor{}, fmt.Errorf("goal-start supports Claude and Codex Goals only")
	}
	return goalExecutor{Provider: provider, Model: g.Executor.Model, Effort: g.Executor.Effort}, nil
}

// goalAgentActor is the actor the started agent's Ply commands record.
func goalAgentActor(plan goalAgentPlan) string {
	return taskactor.Agent(plan.Provider, plan.Task)
}

func goalStartInput(plan goalAgentPlan, receipt string, herdr agentstart.HerdrBinary, precondition func() error) agentstart.StartInput {
	return agentstart.StartInput{
		Key:          plan.Task,
		ReceiptPath:  receipt,
		Provider:     plan.Provider,
		Model:        plan.Model,
		Effort:       plan.Effort,
		Cwd:          plan.Worktree,
		Workspace:    plan.Workspace,
		Label:        plan.Task,
		Env:          []string{"PLY_ACTOR=" + goalAgentActor(plan)},
		Prompt:       goalAgentPrompt + string(plan.Goal) + plan.Documents + plan.Journal,
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
	var spec, herdrWorkspace, actor string
	var revision int
	var apply, resume bool
	c := &cobra.Command{Use: "goal-start TASK_ID", Short: "Start an explicitly assigned Goal in its bound Task worktree", Long: "Preview a registered Goal and existing Task worktree. --apply starts one visible Claude or Codex agent, chosen by the Goal executor provider (Codex when absent). The executor model and effort, when set, are passed to the agent CLI (Claude --model and --effort, Codex --model and -c model_reasoning_effort) as native values the agent CLI checks at start. The agent tab gets PLY_ACTOR=<provider>:<task>, so the Ply commands the agent runs record it as actor. The prompt names the ply workspace task spec show --document command that reads each Goal document; everything else is inherited runtime policy. A durable reservation prevents duplicate launches after failure or unknown effects; inspect its receipt and pane. --apply answers two startup dialogs when the pane shows exactly that dialog: Yes in Claude's folder trust for the Task worktree (receipt folder_trust_accepted) and Skip in Codex's update prompt (receipt codex_update_skipped). When the start stopped before its prompt on any other screen, settle the pane yourself and rerun with --apply --resume: it sends the same Goal prompt at most once to the same agent and never resends an unknown prompt.\n\nKoordinatorretur velges i Goal.constraints; se docs/goal-coordinator-return.md. Goal-start oppretter ingen koordinator og krever ingen retur.", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
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
		executor, e := goalExecutorOf(read.Goal)
		if e != nil {
			return e
		}
		provider := executor.Provider
		wt := task.Task.Worktree
		if e = checkGoalGit(wt); e != nil {
			return e
		}
		if herdrWorkspace == "" {
			return fmt.Errorf("--herdr-workspace is required")
		}
		link, e := taskjournal.New(d).PlanLink(args[0])
		if e != nil {
			return e
		}
		documents, e := goalDocumentInstruction(args[0], spec, revision, read.Goal)
		if e != nil {
			return e
		}
		plan := goalAgentPlan{Task: args[0], Worktree: wt.Locator, Ref: wt.Ref, OID: wt.OID, Workspace: herdrWorkspace, Provider: provider, Model: executor.Model, Effort: executor.Effort, Goal: read.Goal, Documents: documents, Journal: journalInstruction(args[0], link)}
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
		work, e := workspace.LocateTask(d, task.Workspace, id)
		if e != nil {
			return e
		}
		goalRefs := map[string]string{"spec_id": spec, "spec_revision": fmt.Sprint(revision), "provider": provider, "resume": fmt.Sprint(resume)}
		if e = recordPly(d, args[0], namedActor(actor), "nothing was started", taskjournal.PlyAction{Operation: "goal_start_reserved", Refs: goalRefs}); e != nil {
			return e
		}
		out, e := run(c.Context(), plan, workspace.GoalStartReceiptPath(work, id), herdr, func() error { return checkGoalGit(wt) })
		reached := taskjournal.PlyAction{Operation: "goal_start_started", Refs: map[string]string{"spec_id": spec, "spec_revision": fmt.Sprint(revision), "reached": out.Reached}}
		if out.Reached == "prompt_submitted" {
			reached.Operation = "goal_start_prompt_submitted"
		}
		for key, value := range map[string]string{"pending": out.Pending, "pane": out.Pane, "session": out.Session} {
			if value != "" {
				reached.Refs[key] = value
			}
		}
		recordErr := recordPly(d, args[0], namedActor(actor), "the start is in its receipt", reached)
		if writeErr := json.NewEncoder(c.OutOrStdout()).Encode(out); writeErr != nil {
			return writeErr
		}
		return errors.Join(e, recordErr)
	}}
	c.Flags().StringVar(&spec, "spec", "", "registered Goal Spec ID")
	c.Flags().IntVar(&revision, "revision", 0, "registered Goal revision")
	c.Flags().StringVar(&herdrWorkspace, "herdr-workspace", "", "explicit Herdr workspace")
	c.Flags().BoolVar(&apply, "apply", false, "start the explicit Goal once")
	c.Flags().BoolVar(&resume, "resume", false, "with --apply, continue a start stopped before its prompt in the same agent")
	c.Flags().StringVar(&actor, "actor", "", "who starts, recorded in the Task journal (default: PLY_ACTOR, else the owner of the Task's plan)")
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
