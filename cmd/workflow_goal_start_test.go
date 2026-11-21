package cmd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/agentstart"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestWorkflowOffersExplicitGoalStart(t *testing.T) {
	c := newWorkflowCommand()
	got, _, err := c.Find([]string{"goal-start"})
	if err != nil || got == c {
		t.Fatal("registered Goal has no explicit agent-start operation")
	}
}

func TestGoalStartHelpPointsToOptionalCoordinatorReturn(t *testing.T) {
	c := newWorkflowGoalStartCommand(taskwork.Dependencies{})
	var out bytes.Buffer
	c.SetOut(&out)
	if err := c.Help(); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Koordinatorretur velges i Goal.constraints", "docs/goal-coordinator-return.md"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("goal-start help lacks %q: %s", want, out.String())
		}
	}
}

func TestGoalStartRejectsStaleRepositoryBeforeEffects(t *testing.T) {
	task := taskwork.TaskReadbackResult{WorktreeReady: true, Reasons: []string{"target_worktree_stale", "worktree_inventory_stale"}}
	task.Task.Worktree = &taskwork.TaskWorktreeBinding{}
	if err := validateGoalTask(taskwork.Dependencies{}, task); err == nil {
		t.Fatal("stale Git identity admitted")
	}
}
func TestGoalStartAcceptsObservedBindingAndRejectsDifferentRepository(t *testing.T) {
	p := func(s string) *string { return &s }
	yes := true
	task := taskwork.TaskReadbackResult{WorktreeReady: true, ProjectFreshness: "fresh"}
	task.Task.GitCommonDir = "/repo/.git"
	task.Task.Worktree = &taskwork.TaskWorktreeBinding{Locator: "/task", Ref: "refs/heads/task", OID: "oid", Tree: "tree", GitCommonDir: "/repo/.git"}
	task.Target = taskwork.TaskObservedTarget{Locator: p("/task"), Ref: p("refs/heads/task"), OID: p("oid"), Tree: p("tree"), GitCommonDir: p("/repo/.git"), Clean: &yes, InventoryMatch: &yes}
	task.Source = taskwork.TaskObservedSource{Ref: p("refs/heads/task"), OID: p("oid"), Tree: p("tree"), CheckedOutAt: p("/task")}
	if err := validateGoalTask(taskwork.Dependencies{}, task); err != nil {
		t.Fatal(err)
	}
	task.Source.CheckedOutAt = p("/elsewhere")
	if err := validateGoalTask(taskwork.Dependencies{}, task); err == nil {
		t.Fatal("Task branch checked out in another worktree admitted")
	}
	task.Source = taskwork.TaskObservedSource{}
	if err := validateGoalTask(taskwork.Dependencies{}, task); err == nil {
		t.Fatal("unobserved Task branch admitted")
	}
	task.Source = taskwork.TaskObservedSource{Ref: p("refs/heads/task"), OID: p("oid"), Tree: p("tree"), CheckedOutAt: p("/task")}
	task.Target.GitCommonDir = p("/clone/.git")
	if err := validateGoalTask(taskwork.Dependencies{}, task); err == nil {
		t.Fatal("separate clone with identical commit/ref admitted")
	}
	task.Target.GitCommonDir = p("/repo/.git")
	task.Target.Tree = p("other")
	if err := validateGoalTask(taskwork.Dependencies{}, task); err == nil {
		t.Fatal("bound commit with a different tree admitted")
	}
	task.Target.Tree = p("tree")
	task.Target.InventoryMatch = nil
	if err := validateGoalTask(taskwork.Dependencies{}, task); err == nil {
		t.Fatal("unknown inventory admitted")
	}
}

const goalPromptWithReportFile = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. Report results, actual tests, source commit and limitations to coordinator ply-admin wM:p1 through Herdr. Coordinator owns notifications. No push, merge, install, live configuration changes or server restart. Do not start other agents. Preserve your final report in this worktree and return its path. This start does not attest runtime permissions, QA, TaskResult or integration. Goal:\n"

const goalPromptWithPrivateCoordinator = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. Report results, actual tests, source commit and limitations to coordinator ply-admin wM:p1 through Herdr. Coordinator owns notifications. No push, merge, install, live configuration changes or server restart. Do not start other agents. This start does not attest runtime permissions, QA, TaskResult or integration. Goal:\n"

const goalPromptAttestingDeliveryProcess = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. The Goal decides return, notification, delivery and delegation; this start grants no authority beyond it and does not attest runtime permissions, QA, TaskResult or integration. Goal:\n"

const wantGoalPrompt = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. The Goal decides return, notification, delivery and delegation; this start grants no authority beyond it and does not attest runtime permissions. Goal:\n"

func TestGoalStartSelectsExecutorFromRegisteredGoal(t *testing.T) {
	str := func(s string) *string { return &s }
	for name, c := range map[string]struct {
		goal string
		want goalExecutor
		fail bool
	}{
		"no executor":            {`{}`, goalExecutor{Provider: "codex"}, false},
		"null executor":          {`{"executor":null}`, goalExecutor{Provider: "codex"}, false},
		"null fields":            {`{"executor":{"provider":null,"model":null,"effort":null}}`, goalExecutor{Provider: "codex"}, false},
		"codex":                  {`{"executor":{"provider":"codex"}}`, goalExecutor{Provider: "codex"}, false},
		"claude":                 {`{"executor":{"provider":"claude","model":null,"effort":null}}`, goalExecutor{Provider: "claude"}, false},
		"claude model effort":    {`{"executor":{"provider":"claude","model":"opus","effort":"max"}}`, goalExecutor{Provider: "claude", Model: str("opus"), Effort: str("max")}, false},
		"codex model effort":     {`{"executor":{"provider":"codex","model":"gpt-6-astra","effort":"high"}}`, goalExecutor{Provider: "codex", Model: str("gpt-6-astra"), Effort: str("high")}, false},
		"model without provider": {`{"executor":{"model":"x"}}`, goalExecutor{Provider: "codex", Model: str("x")}, false},
		"unknown provider":       {`{"executor":{"provider":"gemini"}}`, goalExecutor{}, true},
		"not json":               {`{`, goalExecutor{}, true},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := goalExecutorOf(json.RawMessage(c.goal))
			if c.fail {
				if err == nil {
					t.Fatalf("admitted %s", c.goal)
				}
				return
			}
			if err != nil || !reflect.DeepEqual(got, c.want) {
				t.Fatalf("%+v %v", got, err)
			}
		})
	}
}

func TestGoalStartPassesGoalModelAndEffortToTheAgent(t *testing.T) {
	model, effort := "gpt-6-astra", "high"
	plan := goalAgentPlan{Task: "task-1", Worktree: "/wt", Workspace: "wM", Provider: "codex", Model: &model, Effort: &effort, Goal: json.RawMessage(`{}`)}
	in := goalStartInput(plan, "/r/task-1.json", agentstart.HerdrBinary{}, nil)
	if in.Model == nil || *in.Model != model || in.Effort == nil || *in.Effort != effort {
		t.Fatalf("%+v", in)
	}
}

func TestGoalStartMapsPlanToStartInputKeepingPromptBindingAndInheritedRuntime(t *testing.T) {
	plan := goalAgentPlan{Task: "task-1", Worktree: "/wt", Ref: "refs/heads/t", OID: "abc", Workspace: "wM", Provider: "claude", Goal: json.RawMessage(`{"a":1}`)}
	herdr := agentstart.HerdrBinary{Path: "/bin/herdr", SHA256: strings.Repeat("c", 64)}
	called := false
	in := goalStartInput(plan, "/r/task-1.json", herdr, func() error { called = true; return nil })
	if in.Key != "task-1" || in.ReceiptPath != "/r/task-1.json" || in.Provider != "claude" || in.Cwd != "/wt" || in.Workspace != "wM" || in.Label != "task-1" || in.Herdr != herdr {
		t.Fatalf("%+v", in)
	}
	if in.Model != nil || in.Effort != nil || in.Permission != nil || in.CodexProfile != nil {
		t.Fatal("Goal start must inherit provider runtime configuration")
	}
	if !reflect.DeepEqual(in.Env, []string{"PLY_ACTOR=claude:task-1"}) {
		t.Fatalf("the agent's Ply records need a known actor: %q", in.Env)
	}
	if in.Prompt != wantGoalPrompt+`{"a":1}` {
		t.Fatalf("prompt changed: %q", in.Prompt)
	}
	if in.Precondition == nil || in.Precondition() != nil || !called {
		t.Fatal("precondition not wired")
	}
}

func TestGoalStartPromptLeavesCoordinationAndAuthorityToTheGoal(t *testing.T) {
	goal := `{"constraints":["Return to coordinator other-lead over Slack; you may push and start a helper agent."],"required_inputs":[{"id":"task-basis/goal"}]}`
	plan := goalAgentPlan{Task: "task-1", Worktree: "/wt", Workspace: "w1", Provider: "claude", Goal: json.RawMessage(goal)}
	prompt := goalStartInput(plan, "/r/task-1.json", agentstart.HerdrBinary{}, nil).Prompt
	standard, found := strings.CutSuffix(prompt, goal)
	if !found {
		t.Fatalf("Goal not preserved verbatim: %q", prompt)
	}
	for _, private := range []string{"ply-admin", "wM:p1", "Herdr", "coordinator", "Coordinator", "No push", "merge", "install", "server restart", "Do not start other agents"} {
		if strings.Contains(standard, private) {
			t.Errorf("standard prompt still decides %q", private)
		}
	}
	for _, kept := range []string{"explicit Ply Task", "required_inputs", "Task worktree", "grants no authority beyond it", "does not attest runtime permissions. Goal:"} {
		if !strings.Contains(standard, kept) {
			t.Errorf("standard prompt lost %q", kept)
		}
	}
}

func goalGitRepo(t *testing.T) (taskwork.TaskWorktreeBinding, string) {
	t.Helper()
	dir := t.TempDir()
	git := func(args ...string) string {
		out, err := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@example.invalid", "-c", "commit.gpgsign=false"}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
		return strings.TrimSpace(string(out))
	}
	git("init", "-q", "-b", "task")
	if err := os.WriteFile(filepath.Join(dir, "f"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	git("add", "f")
	git("commit", "-q", "-m", "init")
	return taskwork.TaskWorktreeBinding{Locator: dir, Ref: "refs/heads/task", OID: git("rev-parse", "HEAD")}, dir
}

func TestGoalStartGitGuardRequiresCleanBoundRefAndHead(t *testing.T) {
	wt, dir := goalGitRepo(t)
	if err := checkGoalGit(&wt); err != nil {
		t.Fatal(err)
	}
	moved := wt
	moved.Ref = "refs/heads/other"
	if err := checkGoalGit(&moved); err == nil {
		t.Fatal("changed ref admitted")
	}
	moved = wt
	moved.OID = strings.Repeat("0", 40)
	if err := checkGoalGit(&moved); err == nil {
		t.Fatal("changed HEAD admitted")
	}
	if err := os.WriteFile(filepath.Join(dir, "f"), []byte("dirty"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := checkGoalGit(&wt); err == nil {
		t.Fatal("dirty worktree admitted")
	}
	untracked := filepath.Join(dir, "new")
	if err := os.WriteFile(filepath.Join(dir, "f"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(untracked, []byte("u"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := checkGoalGit(&wt); err == nil {
		t.Fatal("untracked file admitted")
	}
}

func fakeHerdrBinary(t *testing.T) (agentstart.HerdrBinary, string) {
	t.Helper()
	dir := t.TempDir()
	log := filepath.Join(dir, "calls.log")
	body := `#!/bin/sh
echo "$*" >> ` + log + `
state=` + dir + `
case "$1 $2" in
"tab create")
  p=; for a in "$@"; do [ "$p" = --cwd ] && echo "$a" > $state/cwd; p=$a; done
  echo '{"result":{"root_pane":{"pane_id":"w1:p2","tab_id":"w1:t2","terminal_id":"term-2","workspace_id":"w1"}}}';;
"pane process-info") echo '{"result":{"process_info":{"pane_id":"w1:p2","shell_pid":40,"foreground_process_group_id":40,"foreground_processes":[{"pid":40}]}}}';;
"agent start") echo "$3" > $state/name; echo "$5" > $state/kind; [ -f $state/fail-start ] && exit 1; echo '{"result":{"ok":true}}';;
"agent get")
  n=$(cat $state/name); k=$(cat $state/kind); c=$(cat $state/cwd)
  echo '{"result":{"agent":{"workspace_id":"w1","tab_id":"w1:t2","pane_id":"w1:p2","terminal_id":"term-2","foreground_cwd":"'$c'","name":"'$n'","agent":"'$k'","agent_status":"idle","interactive_ready":true,"launch_pending":false,"agent_session":{"agent":"'$k'","kind":"id","value":"s1"}}}}';;
*) echo '{"result":{"ok":true}}';;
esac
`
	path := filepath.Join(dir, "herdr")
	if err := os.WriteFile(path, []byte(body), 0700); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(body))
	return agentstart.HerdrBinary{Path: path, SHA256: hex.EncodeToString(sum[:])}, log
}

func TestGoalStartStartsEitherProviderOnceThroughSharedOwner(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			herdr, log := fakeHerdrBinary(t)
			wt := t.TempDir()
			receipt := filepath.Join(t.TempDir(), "goal-start", "task-1.json")
			plan := goalAgentPlan{Task: "task-1", Worktree: wt, Workspace: "w1", Provider: provider, Goal: json.RawMessage(`{"executor":{"provider":"` + provider + `"}}`)}
			r, err := startGoal(context.Background(), plan, receipt, herdr, func() error { return nil })
			if err != nil || r.Reached != agentstart.ReachedPromptSubmitted || r.Provider != provider {
				t.Fatalf("%+v %v", r, err)
			}
			raw, _ := os.ReadFile(log)
			lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
			if len(lines) < 5 || !strings.HasPrefix(lines[0], "tab create") || lines[1] != "pane process-info --pane w1:p2" || !strings.Contains(lines[2], "--kind "+provider+" --pane w1:p2") || strings.Contains(lines[2], " -- ") || lines[3] != "agent get w1:p2" || !strings.HasPrefix(lines[4], "agent prompt w1:p2 ") || strings.Count(string(raw), "agent prompt") != 1 {
				t.Fatalf("argv: %q", lines)
			}
			if _, err = startGoal(context.Background(), plan, receipt, herdr, func() error { return nil }); err != nil {
				t.Fatal(err)
			}
			again, _ := os.ReadFile(log)
			if string(again) != string(raw) {
				t.Fatal("same Goal start repeated an effect")
			}
			plan.Goal = json.RawMessage(`{"executor":{"provider":"` + provider + `"},"changed":true}`)
			_, err = startGoal(context.Background(), plan, receipt, herdr, func() error { return nil })
			var e *agentstart.Error
			if !errors.As(err, &e) || e.Code != agentstart.CodeInputMismatch {
				t.Fatalf("changed Goal not refused: %v", err)
			}
			final, _ := os.ReadFile(log)
			if string(final) != string(raw) {
				t.Fatal("mismatch had an effect")
			}
		})
	}
}

func TestGoalStartRefusesReceiptFromEarlierPromptWithoutEffect(t *testing.T) {
	for name, earlier := range map[string]string{"report file": goalPromptWithReportFile, "private coordinator": goalPromptWithPrivateCoordinator, "delivery process attestation": goalPromptAttestingDeliveryProcess} {
		t.Run(name, func(t *testing.T) {
			herdr, log := fakeHerdrBinary(t)
			receipt := filepath.Join(t.TempDir(), "goal-start", "task-1.json")
			plan := goalAgentPlan{Task: "task-1", Worktree: t.TempDir(), Workspace: "w1", Provider: "claude", Goal: json.RawMessage(`{"executor":{"provider":"claude"}}`)}
			started := goalStartInput(plan, receipt, herdr, func() error { return nil })
			started.Prompt = earlier + string(plan.Goal)
			if r, err := agentstart.StartAgent(context.Background(), started); err != nil || r.Reached != agentstart.ReachedPromptSubmitted {
				t.Fatalf("%+v %v", r, err)
			}
			calls, _ := os.ReadFile(log)
			stored, _ := os.ReadFile(receipt)
			_, err := startGoal(context.Background(), plan, receipt, herdr, func() error { return nil })
			var e *agentstart.Error
			if !errors.As(err, &e) || e.Code != agentstart.CodeInputMismatch {
				t.Fatalf("earlier prompt not refused: %v", err)
			}
			callsAfter, _ := os.ReadFile(log)
			storedAfter, _ := os.ReadFile(receipt)
			if string(callsAfter) != string(calls) || string(storedAfter) != string(stored) {
				t.Fatal("changed prompt sent again or rewrote the started receipt")
			}
		})
	}
}

func TestGoalStartDoesNotCarryRawHerdrMechanics(t *testing.T) {
	raw, err := os.ReadFile("workflow_goal_start.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"startGoalAgent", "goalAgentCommand", `"herdr"`, `"agent", "prompt"`, "tab\", \"create"} {
		if strings.Contains(string(raw), forbidden) {
			t.Errorf("caller still owns %s", forbidden)
		}
	}
}

func TestGoalStartResumeContinuesStoppedStartWithOnePrompt(t *testing.T) {
	herdr, log := fakeHerdrBinary(t)
	failStart := filepath.Join(filepath.Dir(log), "fail-start")
	if err := os.WriteFile(failStart, nil, 0600); err != nil {
		t.Fatal(err)
	}
	receipt := filepath.Join(t.TempDir(), "goal-start", "task-1.json")
	plan := goalAgentPlan{Task: "task-1", Worktree: t.TempDir(), Workspace: "w1", Provider: "claude", Goal: json.RawMessage(`{"executor":{"provider":"claude"}}`)}
	r, err := startGoal(context.Background(), plan, receipt, herdr, func() error { return nil })
	var e *agentstart.Error
	if !errors.As(err, &e) || e.Code != agentstart.ClassAgentResponse || r.Pending != agentstart.StageAgentStart {
		t.Fatalf("start did not stop at agent start: %+v %v", r, err)
	}
	if err = os.Remove(failStart); err != nil {
		t.Fatal(err)
	}
	stopped, _ := os.ReadFile(log)
	r, err = resumeGoal(context.Background(), plan, receipt, herdr, func() error { return nil })
	if err != nil || r.Reached != agentstart.ReachedPromptSubmitted || r.Failure != nil {
		t.Fatalf("%+v %v", r, err)
	}
	resumed, _ := os.ReadFile(log)
	added, _ := strings.CutPrefix(string(resumed), string(stopped))
	if want := "agent get w1:p2\nagent prompt w1:p2 " + wantGoalPrompt + string(plan.Goal) + " --wait --until working --timeout 10000\n"; added != want || strings.Count(string(resumed), "agent prompt") != 1 {
		t.Fatalf("resume calls: %q", added)
	}
	stored, _ := os.ReadFile(receipt)
	again, err := resumeGoal(context.Background(), plan, receipt, herdr, func() error { return nil })
	if err != nil || again != r {
		t.Fatalf("repeated resume: %+v %v", again, err)
	}
	final, _ := os.ReadFile(log)
	storedAfter, _ := os.ReadFile(receipt)
	if string(final) != string(resumed) || string(storedAfter) != string(stored) {
		t.Fatal("repeated resume had an effect")
	}
}

func TestGoalStartResumeRequiresApplyBeforeReading(t *testing.T) {
	c := newWorkflowGoalStartCommand(taskwork.Dependencies{})
	c.SetArgs([]string{"task-1", "--resume"})
	c.SetOut(new(strings.Builder))
	c.SetErr(new(strings.Builder))
	err := c.Execute()
	if err == nil || !strings.Contains(err.Error(), "--resume requires --apply") {
		t.Fatalf("resume without apply: %v", err)
	}
}

// A Task linked to a plan gets the journal commands after its Goal; any
// other Task gets the unchanged prompt.
func TestGoalStartPromptAsksAPlanLinkedTaskToRecordItsSteps(t *testing.T) {
	planID, step := "PLAN-39", "impl-d10"
	if journalInstruction("task-1", nil) != "" {
		t.Fatal("a Task without plan link got a journal instruction")
	}
	goal := `{"a":1}`
	plan := goalAgentPlan{Task: "task-1", Provider: "claude", Goal: json.RawMessage(goal), Journal: journalInstruction("task-1", &taskjournal.PlyAction{Operation: "plan_link", Plan: &planID, PlanStep: &step})}
	prompt := goalStartInput(plan, "/r/task-1.json", agentstart.HerdrBinary{}, nil).Prompt
	rest, found := strings.CutPrefix(prompt, wantGoalPrompt+goal)
	if !found {
		t.Fatalf("Goal not kept first: %q", prompt)
	}
	for _, want := range []string{"PLAN-39 step impl-d10", "ply workspace task journal step start task-1 --id", "step finish task-1", "--parent", "waiting", "declined", "journal note task-1", "PLY_ACTOR", "final reporting and cleanup use kind coordination"} {
		if !strings.Contains(rest, want) {
			t.Errorf("instruction lacks %q: %q", want, rest)
		}
	}
}

// Goal documents are preserved Ply data; the prompt names the ply command
// that reads each one, and a Goal without documents gets no instruction.
func TestGoalStartPromptNamesTheReaderOfEachGoalDocument(t *testing.T) {
	goal := json.RawMessage(`{"required_inputs":[{"id":"task-basis/goal"},{"id":"task-goal-doc/goal/spec"},{"id":"task-goal-doc/goal/jvm-plan"},{"id":"task-goal-doc/problem/problem"}]}`)
	got, err := goalDocumentInstruction("task-1", "spec-v1", 2, goal)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"ply workspace task spec show task-1 --spec spec-v1 --revision 2 --document spec", "ply workspace task spec show task-1 --spec spec-v1 --revision 2 --document jvm-plan"} {
		if !strings.Contains(got, want) {
			t.Errorf("instruction lacks %q: %q", want, got)
		}
	}
	if !strings.HasSuffix(got, "\n  ply workspace task spec show task-1 --spec spec-v1 --revision 2 --document jvm-plan") {
		t.Errorf("each command must stand alone on its line: %q", got)
	}
	if strings.Contains(got, "--document problem") {
		t.Errorf("problem document is not a Goal document: %q", got)
	}
	if got, err = goalDocumentInstruction("task-1", "spec-v1", 2, json.RawMessage(`{"required_inputs":[{"id":"task-basis/goal"}]}`)); err != nil || got != "" {
		t.Fatalf("%q %v", got, err)
	}
	plan := goalAgentPlan{Task: "task-1", Provider: "codex", Goal: goal, Documents: "\n\nDOCS", Journal: "\n\nJOURNAL"}
	if prompt := goalStartInput(plan, "/r", agentstart.HerdrBinary{}, nil).Prompt; prompt != wantGoalPrompt+string(goal)+"\n\nDOCS\n\nJOURNAL" {
		t.Fatalf("%q", prompt)
	}
}

// goalStartFixture gives Task "moved" the registered Goal moved-v1 and a ready
// worktree bound at the Epic's second commit, in a workspace of its own. It
// returns the worktree, the Epic's first commit and the bound commit.
func goalStartFixture(t *testing.T) (d taskwork.Dependencies, git func(dir string, args ...string) string, worktree, first, bound string) {
	t.Helper()
	d, root := viewFixture(t)
	git = func(dir string, args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@example.invalid", "-c", "commit.gpgsign=false", "-c", "core.hooksPath=" + os.DevNull}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
		return strings.TrimSpace(string(out))
	}
	run := func(args ...string) string {
		t.Helper()
		out, err := runWorkspaceView(t, d, args...)
		if err != nil {
			t.Fatalf("workspace %v: %v", args, err)
		}
		return out
	}
	repo := filepath.Join(root, "repo")
	if err := os.Mkdir(repo, 0700); err != nil {
		t.Fatal(err)
	}
	git(repo, "init", "-q", "-b", "epic")
	git(repo, "commit", "-q", "--allow-empty", "-m", "first")
	first = git(repo, "rev-parse", "HEAD")
	git(repo, "commit", "-q", "--allow-empty", "-m", "second")
	bound = git(repo, "rev-parse", "HEAD")
	run("project", "add", "example", "--name", "Example", "--wrapper", root, "--repo", "code="+repo)
	run("epic", "adopt", "work", "--title", "Work", "--project", "example", "--repo", "code", "--worktree", repo, "--ref", "refs/heads/epic", "--expected-oid", bound)
	register := []string{"task", "goal", "register", "moved", "--spec", "moved-v1", "--title", "T", "--objective", "O", "--criterion", "C", "--project", "example", "--repo", "code", "--epic", "work", "--format", "json"}
	var preview struct{ Confirmation string }
	if err := json.Unmarshal([]byte(run(register...)), &preview); err != nil {
		t.Fatal(err)
	}
	run(append(register, "--apply", "--confirmation", preview.Confirmation)...)
	worktree = filepath.Join(root, "moved")
	run("task", "worktree", "create", "moved", "--branch", "moved", "--path", worktree, "--expected-parent-oid", bound)
	return d, git, worktree, first, bound
}

func previewGoalStart(t *testing.T, d taskwork.Dependencies) (goalAgentPlan, error) {
	t.Helper()
	c := newWorkflowGoalStartCommand(d)
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	c.SilenceErrors, c.SilenceUsage = true, true
	c.SetArgs([]string{"moved", "--spec", "moved-v1", "--revision", "1", "--herdr-workspace", "w1"})
	if err := c.Execute(); err != nil {
		return goalAgentPlan{}, err
	}
	var plan goalAgentPlan
	if err := json.Unmarshal(out.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	return plan, nil
}

// Backlog 20261117-28d0dcb6: a Task whose agent has committed gets its next
// Goal in the same worktree, starting at the commit the worktree is at.
func TestGoalStartAcceptsCleanWorktreeAdvancedFromItsBinding(t *testing.T) {
	d, git, worktree, _, bound := goalStartFixture(t)
	plan, err := previewGoalStart(t, d)
	if err != nil || plan.OID != bound {
		t.Fatalf("bound Task worktree refused: %+v %v", plan, err)
	}
	git(worktree, "commit", "-q", "--allow-empty", "-m", "work")
	head := git(worktree, "rev-parse", "HEAD")
	plan, err = previewGoalStart(t, d)
	if err != nil || plan.OID != head || plan.Ref != "refs/heads/moved" || plan.Worktree != worktree {
		t.Fatalf("advanced Task worktree refused: %+v %v", plan, err)
	}
}

func TestGoalStartRejectsWorktreeThatIsNotACleanDescendantOnItsRef(t *testing.T) {
	for name, c := range map[string]struct {
		drift func(t *testing.T, git func(string, ...string) string, worktree, first string)
		want  string
	}{
		"dirty": {func(t *testing.T, git func(string, ...string) string, worktree, first string) {
			if err := os.WriteFile(filepath.Join(worktree, "untracked"), []byte("u"), 0600); err != nil {
				t.Fatal(err)
			}
		}, "Task Git identity or inventory is stale or unknown"},
		"head without the bound commit as ancestor": {func(t *testing.T, git func(string, ...string) string, worktree, first string) {
			git(worktree, "reset", "-q", "--hard", git(worktree, "commit-tree", "-m", "unrelated", "HEAD^{tree}"))
		}, "Task Git binding differs from the observed worktree"},
		"head behind the bound commit": {func(t *testing.T, git func(string, ...string) string, worktree, first string) {
			git(worktree, "reset", "-q", "--hard", first)
		}, "Task Git binding differs from the observed worktree"},
		"other ref": {func(t *testing.T, git func(string, ...string) string, worktree, first string) {
			git(worktree, "switch", "-q", "-c", "other")
		}, "Task Git binding differs from the observed worktree"},
		"branch checked out elsewhere": {func(t *testing.T, git func(string, ...string) string, worktree, first string) {
			git(worktree, "switch", "-q", "--detach")
			git(worktree, "worktree", "add", "-q", filepath.Join(filepath.Dir(worktree), "elsewhere"), "moved")
		}, "Task Git"},
		"unobservable worktree": {func(t *testing.T, git func(string, ...string) string, worktree, first string) {
			if err := os.RemoveAll(worktree); err != nil {
				t.Fatal(err)
			}
		}, "Task Git identity or inventory is stale or unknown"},
	} {
		t.Run(name, func(t *testing.T) {
			d, git, worktree, first, _ := goalStartFixture(t)
			git(worktree, "commit", "-q", "--allow-empty", "-m", "work")
			c.drift(t, git, worktree, first)
			if plan, err := previewGoalStart(t, d); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("drifted Task worktree: %+v %v", plan, err)
			}
		})
	}
}

// Backlog 20261121-4e694844: Git can be forced to check the Task branch out
// in a second worktree while the Task worktree keeps it. The branch can then
// not be observed, and the start is refused before any effect.
func TestGoalStartRejectsTaskBranchForcedOutInASecondWorktree(t *testing.T) {
	for name, c := range map[string]struct {
		advanced, dirty bool
		force           func(git func(string, ...string) string, worktree, second string)
	}{
		"worktree add --force at the bound commit": {false, false, func(git func(string, ...string) string, worktree, second string) {
			git(worktree, "worktree", "add", "-q", "--force", second, "moved")
		}},
		"checkout --ignore-other-worktrees after a commit": {true, false, func(git func(string, ...string) string, worktree, second string) {
			git(worktree, "worktree", "add", "-q", "--detach", second)
			git(second, "checkout", "-q", "--ignore-other-worktrees", "moved")
		}},
		// ASSIGNMENT-300: an unregistered file must not hide Git's cause.
		"worktree add --force with an unregistered file": {false, true, func(git func(string, ...string) string, worktree, second string) {
			git(worktree, "worktree", "add", "-q", "--force", second, "moved")
		}},
	} {
		t.Run(name, func(t *testing.T) {
			d, git, worktree, _, _ := goalStartFixture(t)
			if c.advanced {
				git(worktree, "commit", "-q", "--allow-empty", "-m", "work")
			}
			second := filepath.Join(filepath.Dir(worktree), "second")
			c.force(git, worktree, second)
			if c.dirty {
				if err := os.WriteFile(filepath.Join(worktree, "untracked"), nil, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			for _, dir := range []string{worktree, second} {
				if ref := git(dir, "symbolic-ref", "HEAD"); ref != "refs/heads/moved" {
					t.Fatalf("fixture: %s is on %s", dir, ref)
				}
			}
			want := "Task branch could not be observed as checked out in the Task worktree only: " + string(workspace.ErrorWorkGitObservation) + ": ref refs/heads/moved is checked out more than once"
			if _, err := previewGoalStart(t, d); err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("preview with the Task branch in two worktrees: %v", err)
			}
			journal := func() string {
				t.Helper()
				out, err := runWorkspaceView(t, d, "task", "journal", "show", "moved", "--format", "json")
				if err != nil {
					t.Fatal(err)
				}
				var snap taskjournal.Snapshot
				if err = json.Unmarshal([]byte(out), &snap); err != nil {
					t.Fatal(err)
				}
				var events []string
				for _, event := range snap.Events {
					events = append(events, event.EventID+" "+event.Title)
				}
				return strings.Join(events, "\n")
			}
			before := journal()
			log := fakeHerdrOnPath(t)
			if r, err := applyGoalStart(t, d, "1"); err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("start with the Task branch in two worktrees: %+v %v", r, err)
			}
			if calls, err := os.ReadFile(log); !os.IsNotExist(err) {
				t.Fatalf("the refused start reached Herdr: %s %v", calls, err)
			}
			if after := journal(); after != before {
				t.Fatalf("the refused start wrote to the Task journal: %s", after)
			}
			if receipts, err := filepath.Glob(filepath.Join(filepath.Dir(worktree), "data", "work", "task-process", "goal-start", "*")); err != nil || len(receipts) != 0 {
				t.Fatalf("the refused start reserved a receipt: %v %v", receipts, err)
			}
		})
	}
}

// Backlog 20261121-68d70343: when Git cannot observe the Task branch because
// it is forced out in a second worktree, task show reports the failed
// observation with its Git cause instead of a stale start base.
func TestTaskShowReportsAFailedSourceObservationWithItsGitCause(t *testing.T) {
	d, git, worktree, _, _ := goalStartFixture(t)
	git(worktree, "worktree", "add", "-q", "--force", filepath.Join(filepath.Dir(worktree), "second"), "moved")
	// A dirty Task worktree adds a reason that must not carry the Git cause.
	if err := os.WriteFile(filepath.Join(worktree, "untracked"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	cause := string(workspace.ErrorWorkGitObservation) + ": ref refs/heads/moved is checked out more than once"
	out, err := runWorkspaceView(t, d, "task", "show", "moved", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var shown struct {
		Resource struct {
			Reasons  []string `json:"reasons"`
			Observed struct {
				Source map[string]any `json:"source"`
			} `json:"observed"`
		} `json:"resource"`
	}
	if err = json.Unmarshal([]byte(out), &shown); err != nil {
		t.Fatal(err)
	}
	if reasons := shown.Resource.Reasons; !slices.Contains(reasons, "source_observation_unknown") || slices.Contains(reasons, "source_ref_stale") {
		t.Errorf("reasons = %q, want source_observation_unknown and no source_ref_stale", reasons)
	}
	if got := shown.Resource.Observed.Source["diagnostic"]; got != cause {
		t.Errorf("observed.source.diagnostic = %#v, want %q", got, cause)
	}
	text, err := runWorkspaceView(t, d, "task", "show", "moved")
	if err != nil {
		t.Fatal(err)
	}
	if want := "Reason: The current Task branch state could not be observed: " + cause + "\n"; !strings.Contains(text, want) || strings.Count(text, cause) != 1 || !strings.Contains(text, "Reason: The Task worktree is not clean.\n") || strings.Contains(text, "no longer matches the recorded start base") || !strings.Contains(text, "Fresh observation: unknown;") {
		t.Errorf("task show text lacks %q or still reports a stale start base:\n%s", want, text)
	}
}

// applyGoalStart runs goal-start --apply for one revision of moved-v1 with
// the fake Herdr of the test.
func applyGoalStart(t *testing.T, d taskwork.Dependencies, revision string) (agentstart.Receipt, error) {
	t.Helper()
	return applyGoalStartIn(t, d, revision, "w1")
}

func applyGoalStartIn(t *testing.T, d taskwork.Dependencies, revision, herdrWorkspace string) (agentstart.Receipt, error) {
	t.Helper()
	c := newWorkflowGoalStartCommand(d)
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	c.SilenceErrors, c.SilenceUsage = true, true
	c.SetArgs([]string{"moved", "--spec", "moved-v1", "--revision", revision, "--herdr-workspace", herdrWorkspace, "--apply", "--actor", "planner"})
	err := c.Execute()
	var r agentstart.Receipt
	if out.Len() > 0 {
		if e := json.Unmarshal(out.Bytes(), &r); e != nil {
			t.Fatal(e)
		}
	}
	return r, err
}

// fakeHerdrOnPath lets goal-start --apply of this test reach a fake Herdr and
// returns its call log.
func fakeHerdrOnPath(t *testing.T) string {
	t.Helper()
	herdr, log := fakeHerdrBinary(t)
	t.Setenv("HERDR_ENV", "1")
	t.Setenv("PATH", filepath.Dir(herdr.Path)+string(os.PathListSeparator)+os.Getenv("PATH"))
	return log
}

// registerSecondGoalRevision records revision 2 of moved-v1 with a changed
// objective and the given executor provider.
func registerSecondGoalRevision(t *testing.T, d taskwork.Dependencies, worktree string, provider any) {
	t.Helper()
	out, err := runWorkspaceView(t, d, "task", "spec", "show", "moved", "--spec", "moved-v1", "--revision", "1", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var shown struct {
		Goal struct {
			Goal struct {
				Spec json.RawMessage `json:"spec"`
			} `json:"goal"`
		} `json:"goal"`
		Revision map[string]any `json:"revision"`
	}
	if err = json.Unmarshal([]byte(out), &shown); err != nil {
		t.Fatal(err)
	}
	draft := shown.Revision
	for _, published := range []string{"previous", "recorded_at_utc", "revision", "source_draft_sha256"} {
		delete(draft, published)
	}
	draft["kind"], draft["registry_upgrade"], draft["expected_previous"] = "WorkspaceTaskSpecDraft@2", nil, shown.Goal.Goal.Spec
	draft["publication_key"], draft["objective"], draft["change_reason"] = "fixture/moved-goal-2", "O after the answers", "Answer the questions"
	draft["executor"] = map[string]any{"provider": provider, "model": nil, "effort": nil}
	raw, err := json.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(filepath.Dir(worktree), "goal-2.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = runWorkspaceView(t, d, "task", "spec", "record", "moved", "--file", path, "--format", "json"); err != nil {
		t.Fatal(err)
	}
}

// The guard before each Herdr effect holds the start to the commit observed
// in the worktree, so an agent starts in an advanced worktree. The same Goal
// is still sent once, also after further commits.
func TestGoalStartApplyStartsOnceInAdvancedWorktree(t *testing.T) {
	d, git, worktree, _, _ := goalStartFixture(t)
	git(worktree, "commit", "-q", "--allow-empty", "-m", "work")
	log := fakeHerdrOnPath(t)
	r, err := applyGoalStart(t, d, "1")
	if err != nil || r.Reached != agentstart.ReachedPromptSubmitted {
		t.Fatalf("advanced Task worktree not started: %+v %v", r, err)
	}
	calls, _ := os.ReadFile(log)
	if strings.Count(string(calls), "agent prompt") != 1 || !strings.Contains(string(calls), "--cwd "+worktree) {
		t.Fatalf("calls: %s", calls)
	}
	git(worktree, "commit", "-q", "--allow-empty", "-m", "more work")
	if again, err := applyGoalStart(t, d, "1"); err != nil || again != r {
		t.Fatalf("repeated start: %+v %v", again, err)
	}
	if after, _ := os.ReadFile(log); string(after) != string(calls) {
		t.Fatal("repeated start had an effect")
	}
}

// Backlog 20261117-28d0dcb6: after the agent of revision 1 has committed and
// revision 2 is registered, the same Task starts one agent for revision 2.
// Each revision is still sent once.
func TestGoalStartStartsTheNextGoalRevisionInTheSameTask(t *testing.T) {
	d, git, worktree, _, _ := goalStartFixture(t)
	log := fakeHerdrOnPath(t)
	first, err := applyGoalStart(t, d, "1")
	if err != nil || first.Reached != agentstart.ReachedPromptSubmitted {
		t.Fatalf("revision 1 not started: %+v %v", first, err)
	}
	git(worktree, "commit", "-q", "--allow-empty", "-m", "work")
	registerSecondGoalRevision(t, d, worktree, nil)
	second, err := applyGoalStart(t, d, "2")
	if err != nil || second.Reached != agentstart.ReachedPromptSubmitted || second.Key == first.Key || second.InputSHA256 == first.InputSHA256 {
		t.Fatalf("revision 2 not started beside revision 1: %+v %v", second, err)
	}
	calls, _ := os.ReadFile(log)
	var agents []string
	for _, line := range strings.Split(string(calls), "\n") {
		if name, found := strings.CutPrefix(line, "agent start "); found {
			agents = append(agents, strings.Fields(name)[0])
		}
	}
	if strings.Count(string(calls), "agent prompt") != 2 || !strings.Contains(string(calls), "O after the answers") || len(agents) != 2 || agents[0] == agents[1] {
		t.Fatalf("expected one prompt per revision to agents of their own: %s", calls)
	}
	for revision, want := range map[string]agentstart.Receipt{"1": first, "2": second} {
		if again, err := applyGoalStart(t, d, revision); err != nil || again != want {
			t.Fatalf("repeated start of revision %s: %+v %v", revision, again, err)
		}
	}
	if after, _ := os.ReadFile(log); string(after) != string(calls) {
		t.Fatal("a repeated start had an effect")
	}
}

// A Task started before receipts were kept per Goal revision has one receipt
// that names no revision: it still answers a repeated start and refuses any
// other, saying that the next revision needs a new Task, so nothing is
// started twice.
func TestGoalStartKeepsTheSingleReceiptOfAnOlderStart(t *testing.T) {
	d, git, worktree, _, _ := goalStartFixture(t)
	log := fakeHerdrOnPath(t)
	started, err := applyGoalStart(t, d, "1")
	if err != nil {
		t.Fatal(err)
	}
	receipts, err := filepath.Glob(filepath.Join(filepath.Dir(worktree), "data", "work", "task-process", "goal-start", "moved", "*.json"))
	if err != nil || len(receipts) != 1 {
		t.Fatalf("receipt of revision 1: %v %v", receipts, err)
	}
	started.Key = "moved"
	older, err := json.Marshal(started)
	if err != nil {
		t.Fatal(err)
	}
	single := filepath.Dir(receipts[0]) + ".json"
	if err = os.WriteFile(single, older, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.RemoveAll(filepath.Dir(receipts[0])); err != nil {
		t.Fatal(err)
	}
	calls, _ := os.ReadFile(log)
	if again, err := applyGoalStart(t, d, "1"); err != nil || again != started {
		t.Fatalf("repeated older start: %+v %v", again, err)
	}
	git(worktree, "commit", "-q", "--allow-empty", "-m", "work")
	registerSecondGoalRevision(t, d, worktree, nil)
	_, err = applyGoalStart(t, d, "2")
	var mismatch *agentstart.Error
	if !errors.As(err, &mismatch) || mismatch.Code != agentstart.CodeInputMismatch {
		t.Fatalf("revision 2 beside a receipt that names no revision: %v", err)
	}
	for _, says := range []string{"started by an older Ply", "receipt names no Goal revision", "revision 2", "needs a new Task"} {
		if !strings.Contains(err.Error(), says) {
			t.Fatalf("the refusal of revision 2 does not say %q: %v", says, err)
		}
	}
	if after, _ := os.ReadFile(log); string(after) != string(calls) {
		t.Fatal("a start beside the older receipt had an effect")
	}
	if _, err = os.Stat(filepath.Dir(receipts[0])); !os.IsNotExist(err) {
		t.Fatalf("a receipt per revision was reserved beside the older receipt: %v", err)
	}
}

// A receipt per Goal revision refuses other input for its revision as
// before, without the explanation meant for a receipt of an older Ply.
func TestGoalStartRefusesOtherInputForARevisionWithItsOwnReceipt(t *testing.T) {
	d, _, _, _, _ := goalStartFixture(t)
	fakeHerdrOnPath(t)
	if _, err := applyGoalStart(t, d, "1"); err != nil {
		t.Fatal(err)
	}
	_, err := applyGoalStartIn(t, d, "1", "w2")
	var mismatch *agentstart.Error
	if !errors.As(err, &mismatch) || mismatch.Code != agentstart.CodeInputMismatch || strings.Contains(err.Error(), "older Ply") {
		t.Fatalf("revision 1 in another Herdr workspace: %v", err)
	}
}

// A journal record without a named actor gets the agent of the highest
// started Goal revision.
func TestTaskAgentIsTheAgentOfTheHighestStartedGoalRevision(t *testing.T) {
	d, git, worktree, _, _ := goalStartFixture(t)
	fakeHerdrOnPath(t)
	t.Setenv("PLY_ACTOR", "")
	step := func(id string) string {
		t.Helper()
		if _, err := runWorkspaceView(t, d, "task", "journal", "step", "start", "moved", "--id", id, "--kind", "implementation", "--title", "Step "+id); err != nil {
			t.Fatal(err)
		}
		out, err := runWorkspaceView(t, d, "task", "journal", "show", "moved", "--format", "json")
		if err != nil {
			t.Fatal(err)
		}
		var snap taskjournal.Snapshot
		if err = json.Unmarshal([]byte(out), &snap); err != nil {
			t.Fatal(err)
		}
		for _, event := range snap.Events {
			if event.Title == "Step "+id && event.Actor != nil {
				return event.Actor.ID
			}
		}
		t.Fatalf("step %s not in the journal: %s", id, out)
		return ""
	}
	if _, err := applyGoalStart(t, d, "1"); err != nil {
		t.Fatal(err)
	}
	if actor := step("s1"); actor != "codex:moved" {
		t.Fatalf("actor after revision 1: %s", actor)
	}
	git(worktree, "commit", "-q", "--allow-empty", "-m", "work")
	registerSecondGoalRevision(t, d, worktree, "claude")
	if _, err := applyGoalStart(t, d, "2"); err != nil {
		t.Fatal(err)
	}
	if actor := step("s2"); actor != "claude:moved" {
		t.Fatalf("actor after revision 2: %s", actor)
	}
}
