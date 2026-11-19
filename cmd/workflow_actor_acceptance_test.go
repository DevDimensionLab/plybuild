package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// startTask registers task on PLAN-1 step review, creates its worktree and
// starts its Codex agent with a fake Herdr, all without --actor or
// PLY_ACTOR, so every actor comes from Ply's own rule.
func startTask(t *testing.T, f planCLIFixture, oid, task string, extra ...string) {
	t.Helper()
	registerGoal(t, f, task, append([]string{"--plan", "PLAN-1", "--plan-step", "review"}, extra...)...)
	f.success(t, f.root, "workspace", "task", "worktree", "create", task, "--branch", task, "--path", filepath.Join(f.root, task), "--expected-parent-oid", oid)
	herdr, _ := fakeHerdrBinary(t)
	env := map[string]string{"HERDR_ENV": "1", "PATH": filepath.Dir(herdr.Path) + string(os.PathListSeparator) + os.Getenv("PATH")}
	if _, stderr, exit := f.run(t, f.root, env, "workflow", "goal-start", task, "--spec", task+"-v1", "--revision", "1", "--herdr-workspace", "w1", "--apply"); exit != 0 {
		t.Fatalf("goal-start: %d %s", exit, stderr)
	}
}

func traceOf(t *testing.T, f planCLIFixture) planTraceDoc {
	t.Helper()
	var trace planTraceDoc
	if e := json.Unmarshal([]byte(f.success(t, f.root, "plan", "trace", "PLAN-1", "--project", "example", "--format", "json")), &trace); e != nil {
		t.Fatal(e)
	}
	return trace
}

// Spec ply-actor-and-task-role: without --actor and PLY_ACTOR, Ply's
// operations get the plan owner and the agent's steps get the agent that
// goal-start started; a reviewer Task's steps get its role.
func TestActorAndRoleComeFromTheTaskWithoutActorOrEnvironment(t *testing.T) {
	f, oid := recordingFixture(t)
	before := mr05Snapshot(t, f.root)
	f.failure(t, f.root, nil, "--role must be developer or reviewer", "workspace", "task", "goal", "register", "bad", "--spec", "bad-v1", "--title", "T", "--objective", "O", "--criterion", "C",
		"--project", "example", "--repo", "code", "--epic", "work", "--role", "boss")
	if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
		t.Fatal("refused role wrote state")
	}
	startTask(t, f, oid, "rev", "--role", "reviewer")
	if got := operationsOf(t, plyActions(t, f, "rev"), "planner-a"); strings.Join(got, ",") != "goal_register,goal_start_prompt_submitted,goal_start_reserved,plan_link,worktree_create" {
		t.Fatalf("operations %v", got)
	}
	f.success(t, f.root, "workspace", "task", "journal", "step", "start", "rev", "--id", "r1", "--kind", "review", "--title", "Review")
	if _, stderr, exit := f.run(t, f.root, map[string]string{"PLY_ACTOR": "coordinator"}, "workspace", "task", "journal", "note", "rev", "--title", "Checked"); exit != 0 {
		t.Fatalf("note: %d %s", exit, stderr)
	}
	trace := traceOf(t, f)
	if len(trace.Tasks) != 1 || trace.Tasks[0].Role != "reviewer" {
		t.Fatalf("tasks %+v", trace.Tasks)
	}
	var step, note bool
	for _, ev := range trace.Events {
		if ev.ActorBasis != "" {
			t.Fatalf("new record with derived actor: %+v", ev)
		}
		step = step || ev.Type == "step_started" && *ev.ActorID == "codex:rev" && *ev.Role == "reviewer"
		note = note || ev.Type == "note" && *ev.ActorID == "coordinator" && *ev.Role == "reviewer"
	}
	if !step || !note {
		t.Fatalf("events %+v", trace.Events)
	}

	registerGoal(t, f, "loose")
	if got := operationsOf(t, plyActions(t, f, "loose"), "unknown"); strings.Join(got, ",") != "goal_register" {
		t.Fatalf("operations %v", got)
	}
}

// Records stored with actor unknown before the rule get their actor by the
// same rule when the trace is read, marked derived; the journal is unchanged.
func TestPlanTraceDerivesTheActorOfOlderUnknownRecords(t *testing.T) {
	f, oid := recordingFixture(t)
	startTask(t, f, oid, "old")
	f.success(t, f.root, "workspace", "task", "journal", "step", "start", "old", "--id", "i1", "--kind", "implementation", "--title", "Implement", "--actor", "unknown")
	f.success(t, f.root, "workspace", "task", "journal", "link", "old", "--plan", "PLAN-1", "--plan-step", "impl", "--actor", "unknown")
	var step, link bool
	for _, ev := range traceOf(t, f).Events {
		step = step || ev.Type == "step_started" && *ev.ActorID == "codex:old" && ev.ActorBasis == "derived" && *ev.Role == "developer"
		link = link || ev.PlyOperation != nil && *ev.PlyOperation == "plan_link" && *ev.ActorID == "planner-a" && ev.ActorBasis == "derived"
	}
	if !step || !link {
		t.Fatalf("events %+v", traceOf(t, f).Events)
	}
	if got := operationsOf(t, plyActions(t, f, "old")[5:], "unknown"); strings.Join(got, ",") != "plan_link" {
		t.Fatalf("stored records changed: %v", got)
	}
}
