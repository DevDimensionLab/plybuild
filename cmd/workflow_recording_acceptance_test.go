package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
)

// recordingFixture is a workspace with one Epic and PLAN-1 (steps impl and
// review) for the workflow recording acceptance of backlog 20261115-08522a41.
func recordingFixture(t *testing.T) (planCLIFixture, string) {
	t.Helper()
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("opt-in CLI fixture: set PLY_CANDIDATE_BINARY to the actual make-built candidate")
	}
	f := newPlanCLI(t, binary)
	repo := filepath.Join(f.root, "repo")
	if e := os.Mkdir(repo, 0700); e != nil {
		t.Fatal(e)
	}
	f.git(t, repo, "init", "-b", "epic")
	mustPlanWrite(t, filepath.Join(repo, "README.md"), "fixture\n", 0600)
	f.git(t, repo, "add", "README.md")
	f.git(t, repo, "commit", "-m", "fixture")
	oid := f.git(t, repo, "rev-parse", "HEAD")
	f.success(t, f.root, "workspace", "init")
	f.success(t, f.root, "workspace", "project", "add", "example", "--name", "Example", "--wrapper", f.root, "--repo", "code="+repo)
	f.success(t, f.root, "workspace", "epic", "adopt", "work", "--title", "Work", "--project", "example", "--repo", "code", "--worktree", repo, "--ref", "refs/heads/epic", "--expected-oid", oid)
	snapshot := `{"title":"Recorded plan","goal":"g","scope":"s","owner":"planner-a","selection":"selected","status":"active","summary":"s","next_action":"n","steps":[` +
		`{"id":"impl","title":"Implement","outcome":"o","owner":"dev","status":"active","next_action":"n"},` +
		`{"id":"review","title":"Review","outcome":"o","owner":"rev","status":"planned","next_action":"n"}]}`
	if _, e := plan.Open(filepath.Join(f.root, "plans")).Add([]byte(snapshot), "planner-a", "fixture"); e != nil {
		t.Fatal(e)
	}
	return f, oid
}

func registerGoal(t *testing.T, f planCLIFixture, task string, extra ...string) {
	t.Helper()
	args := append([]string{"workspace", "task", "goal", "register", task, "--spec", task + "-v1", "--title", "T " + task, "--objective", "O", "--criterion", "C",
		"--project", "example", "--repo", "code", "--epic", "work", "--format", "json"}, extra...)
	var preview struct{ Confirmation string }
	if e := json.Unmarshal([]byte(f.success(t, f.root, args...)), &preview); e != nil {
		t.Fatal(e)
	}
	f.success(t, f.root, append(args, "--apply", "--confirmation", preview.Confirmation)...)
}

func plyActions(t *testing.T, f planCLIFixture, task string) []taskjournal.Event {
	t.Helper()
	var snap taskjournal.Snapshot
	if e := json.Unmarshal([]byte(f.success(t, f.root, "workspace", "task", "journal", "show", task, "--format", "json")), &snap); e != nil {
		t.Fatal(e)
	}
	var out []taskjournal.Event
	for _, ev := range snap.Events {
		if ev.Type == "ply_action" {
			out = append(out, ev)
		}
	}
	return out
}

// K1: plan and plan step on a Task, validated before anything is stored; K4:
// Ply's own records carry the actor and Ply's time.
func TestGoalRegisterLinksPlanAndRecordsActor(t *testing.T) {
	f, _ := recordingFixture(t)
	before := mr05Snapshot(t, f.root)
	base := []string{"workspace", "task", "goal", "register", "linked", "--spec", "linked-v1", "--title", "T", "--objective", "O", "--criterion", "C",
		"--project", "example", "--repo", "code", "--epic", "work"}
	f.failure(t, f.root, nil, "--plan PLAN-9 is not a plan of project example", append(base, "--plan", "PLAN-9")...)
	f.failure(t, f.root, nil, "--plan-step nope is not a step of PLAN-1", append(base, "--plan", "PLAN-1", "--plan-step", "nope")...)
	f.failure(t, f.root, nil, "--plan-step requires --plan", append(base, "--plan-step", "impl")...)
	if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
		t.Fatal("refused plan link wrote state")
	}
	registerGoal(t, f, "linked", "--plan", "PLAN-1", "--plan-step", "impl", "--actor", "planner-a")
	actions := plyActions(t, f, "linked")
	operations := []string{}
	for _, ev := range actions {
		var in taskjournal.EventInput
		if e := json.Unmarshal(ev.Data, &in); e != nil || in.Ply == nil {
			t.Fatalf("ply data %s: %v", ev.Data, e)
		}
		operations = append(operations, in.Ply.Operation)
		if ev.Actor == nil || ev.Actor.ID != "planner-a" || ev.Actor.Role != "ply" || ev.OccurredAt == nil || ev.RecordedAt == nil {
			t.Fatalf("native record without actor or time: %+v", ev)
		}
		if in.Ply.Operation == "plan_link" && (*in.Ply.Plan != "PLAN-1" || *in.Ply.PlanStep != "impl") {
			t.Fatalf("plan link %+v", in.Ply)
		}
	}
	if sort.Strings(operations); strings.Join(operations, ",") != "goal_register,plan_link" {
		t.Fatalf("operations %v", operations)
	}
}

func operationsOf(t *testing.T, actions []taskjournal.Event, actor string) []string {
	t.Helper()
	operations := []string{}
	for _, ev := range actions {
		var in taskjournal.EventInput
		if e := json.Unmarshal(ev.Data, &in); e != nil || in.Ply == nil {
			t.Fatalf("ply data %s: %v", ev.Data, e)
		}
		if ev.Actor == nil || ev.Actor.ID != actor || ev.Actor.Role != "ply" || ev.OccurredAt == nil {
			t.Fatalf("native record without actor or time: %+v", ev)
		}
		operations = append(operations, in.Ply.Operation)
	}
	sort.Strings(operations)
	return operations
}

func TestTaskCreateLinksPlanWithActorFromEnvironment(t *testing.T) {
	f, _ := recordingFixture(t)
	args := []string{"workspace", "task", "create", "made", "--title", "Made", "--description", "D", "--epic", "work", "--project", "example", "--repo", "code", "--plan", "PLAN-1", "--plan-step", "review"}
	if _, stderr, exit := f.run(t, f.root, map[string]string{"PLY_ACTOR": "planner-b"}, args...); exit != 0 {
		t.Fatalf("task create: %d %s", exit, stderr)
	}
	if got := operationsOf(t, plyActions(t, f, "made"), "planner-b"); strings.Join(got, ",") != "plan_link,task_create" {
		t.Fatalf("operations %v", got)
	}
}

func TestGoalStartRecordsReservationAndReachedState(t *testing.T) {
	f, oid := recordingFixture(t)
	registerGoal(t, f, "started", "--plan", "PLAN-1", "--plan-step", "impl", "--actor", "planner-a")
	f.success(t, f.root, "workspace", "task", "worktree", "create", "started", "--branch", "started", "--path", filepath.Join(f.root, "started"), "--expected-parent-oid", oid)
	herdr, _ := fakeHerdrBinary(t)
	env := map[string]string{"HERDR_ENV": "1", "PLY_ACTOR": "planner-a", "PATH": filepath.Dir(herdr.Path) + string(os.PathListSeparator) + os.Getenv("PATH")}
	if _, stderr, exit := f.run(t, f.root, env, "workflow", "goal-start", "started", "--spec", "started-v1", "--revision", "1", "--herdr-workspace", "w1", "--apply"); exit != 0 {
		t.Fatalf("goal-start: %d %s", exit, stderr)
	}
	if got := operationsOf(t, plyActions(t, f, "started"), "planner-a"); strings.Join(got, ",") != "goal_register,goal_start_prompt_submitted,goal_start_reserved,plan_link" {
		t.Fatalf("operations %v", got)
	}
}

// integrateTask commits one change in the Task worktree and integrates it
// into the Epic; it returns the integrated commit.
func integrateTask(t *testing.T, f planCLIFixture, task, parent string, env map[string]string) string {
	t.Helper()
	worktree := filepath.Join(f.root, task)
	mustPlanWrite(t, filepath.Join(worktree, task+".txt"), task+"\n", 0600)
	f.git(t, worktree, "add", task+".txt")
	f.git(t, worktree, "commit", "-m", task)
	candidate := f.git(t, worktree, "rev-parse", "HEAD")
	integrate := []string{"workspace", "task", "integrate", task, "--expected-result-oid", candidate, "--expected-parent-oid", parent, "--format", "json"}
	var check struct {
		Plan struct {
			SHA256 string `json:"sha256"`
		} `json:"plan"`
	}
	if e := json.Unmarshal([]byte(f.success(t, f.root, append(integrate, "--check")...)), &check); e != nil {
		t.Fatal(e)
	}
	if _, stderr, exit := f.run(t, f.root, env, append(integrate, "--apply", "--confirm", check.Plan.SHA256)...); exit != 0 {
		t.Fatalf("integrate: %d %s", exit, stderr)
	}
	return candidate
}

func TestTaskIntegrateRecordsActorAndResult(t *testing.T) {
	f, oid := recordingFixture(t)
	registerGoal(t, f, "merged", "--plan", "PLAN-1", "--actor", "planner-a")
	f.success(t, f.root, "workspace", "task", "worktree", "create", "merged", "--branch", "merged", "--path", filepath.Join(f.root, "merged"), "--expected-parent-oid", oid)
	candidate := integrateTask(t, f, "merged", oid, map[string]string{"PLY_ACTOR": "planner-a"})
	var integration *taskjournal.EventInput
	for _, ev := range plyActions(t, f, "merged") {
		var in taskjournal.EventInput
		if e := json.Unmarshal(ev.Data, &in); e == nil && in.Ply != nil && in.Ply.Operation == "integration" {
			integration = &in
		}
	}
	if integration == nil || integration.Actor.ID != "planner-a" || integration.Ply.Refs["result_oid"] != candidate || integration.Ply.Refs["epic"] != "work" {
		t.Fatalf("integration record %+v", integration)
	}
}
