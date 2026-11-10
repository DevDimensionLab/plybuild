package cmd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/agentstart"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
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
	c := newWorkflowGoalStartCommand(workspace.Dependencies{})
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
	task := workspace.TaskReadbackResult{WorktreeReady: true, Reasons: []string{"target_worktree_stale", "worktree_inventory_stale"}}
	task.Task.Worktree = &workspace.TaskWorktreeBinding{}
	if err := validateGoalTask(task); err == nil {
		t.Fatal("stale Git identity admitted")
	}
}
func TestGoalStartAcceptsObservedBindingAndRejectsDifferentRepository(t *testing.T) {
	p := func(s string) *string { return &s }
	yes := true
	task := workspace.TaskReadbackResult{WorktreeReady: true, ProjectFreshness: "fresh"}
	task.Task.GitCommonDir = "/repo/.git"
	task.Task.Worktree = &workspace.TaskWorktreeBinding{Locator: "/task", Ref: "refs/heads/task", OID: "oid", Tree: "tree", GitCommonDir: "/repo/.git"}
	task.Target = workspace.TaskObservedTarget{Locator: p("/task"), Ref: p("refs/heads/task"), OID: p("oid"), Tree: p("tree"), GitCommonDir: p("/repo/.git"), Clean: &yes, InventoryMatch: &yes}
	if err := validateGoalTask(task); err != nil {
		t.Fatal(err)
	}
	task.Target.GitCommonDir = p("/clone/.git")
	if err := validateGoalTask(task); err == nil {
		t.Fatal("separate clone with identical commit/ref admitted")
	}
	task.Target.GitCommonDir = p("/repo/.git")
	task.Target.InventoryMatch = nil
	if err := validateGoalTask(task); err == nil {
		t.Fatal("unknown inventory admitted")
	}
}

const goalPromptWithReportFile = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. Report results, actual tests, source commit and limitations to coordinator ply-admin wM:p1 through Herdr. Coordinator owns notifications. No push, merge, install, live configuration changes or server restart. Do not start other agents. Preserve your final report in this worktree and return its path. This start does not attest runtime permissions, QA, TaskResult or integration. Goal:\n"

const goalPromptWithPrivateCoordinator = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. Report results, actual tests, source commit and limitations to coordinator ply-admin wM:p1 through Herdr. Coordinator owns notifications. No push, merge, install, live configuration changes or server restart. Do not start other agents. This start does not attest runtime permissions, QA, TaskResult or integration. Goal:\n"

const goalPromptAttestingDeliveryProcess = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. The Goal decides return, notification, delivery and delegation; this start grants no authority beyond it and does not attest runtime permissions, QA, TaskResult or integration. Goal:\n"

const wantGoalPrompt = "You own this explicit Ply Task. Read repository instructions and the preserved Goal below, including its required_inputs. Work only within its scope in this Task worktree. The Goal decides return, notification, delivery and delegation; this start grants no authority beyond it and does not attest runtime permissions. Goal:\n"

func TestGoalStartSelectsProviderFromRegisteredExecutor(t *testing.T) {
	for name, c := range map[string]struct {
		goal, want string
		fail       bool
	}{
		"no executor":            {`{}`, "codex", false},
		"null executor":          {`{"executor":null}`, "codex", false},
		"null fields":            {`{"executor":{"provider":null,"model":null,"effort":null}}`, "codex", false},
		"codex":                  {`{"executor":{"provider":"codex"}}`, "codex", false},
		"claude":                 {`{"executor":{"provider":"claude","model":null,"effort":null}}`, "claude", false},
		"unknown provider":       {`{"executor":{"provider":"gemini"}}`, "", true},
		"explicit model":         {`{"executor":{"provider":"claude","model":"opus"}}`, "", true},
		"explicit effort":        {`{"executor":{"provider":"codex","effort":"high"}}`, "", true},
		"model without provider": {`{"executor":{"model":"x"}}`, "", true},
		"not json":               {`{`, "", true},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := goalProvider(json.RawMessage(c.goal))
			if c.fail {
				if err == nil {
					t.Fatalf("admitted %s", c.goal)
				}
				return
			}
			if err != nil || got != c.want {
				t.Fatalf("%q %v", got, err)
			}
		})
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
	if in.Model != nil || in.Permission != nil || in.CodexProfile != nil || len(in.Env) != 0 {
		t.Fatal("Goal start must inherit provider runtime configuration")
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

func goalGitRepo(t *testing.T) (workspace.TaskWorktreeBinding, string) {
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
	return workspace.TaskWorktreeBinding{Locator: dir, Ref: "refs/heads/task", OID: git("rev-parse", "HEAD")}, dir
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
			if len(lines) < 4 || !strings.HasPrefix(lines[0], "tab create") || !strings.Contains(lines[1], "--kind "+provider+" --pane w1:p2") || strings.Contains(lines[1], " -- ") || lines[2] != "agent get w1:p2" || !strings.HasPrefix(lines[3], "agent prompt w1:p2 ") || strings.Count(string(raw), "agent prompt") != 1 {
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
	c := newWorkflowGoalStartCommand(workspace.Dependencies{})
	c.SetArgs([]string{"task-1", "--resume"})
	c.SetOut(new(strings.Builder))
	c.SetErr(new(strings.Builder))
	err := c.Execute()
	if err == nil || !strings.Contains(err.Error(), "--resume requires --apply") {
		t.Fatalf("resume without apply: %v", err)
	}
}
