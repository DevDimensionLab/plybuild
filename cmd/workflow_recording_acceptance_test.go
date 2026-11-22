package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/taskjournal"
	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
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
	if _, e := plan.Open(filepath.Join(f.root, "data", "plans")).Add([]byte(snapshot), "planner-a", "fixture"); e != nil {
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
	f.success(t, f.root, "workspace", "task", "worktree", "create", "started", "--branch", "started", "--path", filepath.Join(f.root, "started"), "--expected-parent-oid", oid, "--actor", "planner-a")
	herdr, _ := fakeHerdrBinary(t)
	env := map[string]string{"HERDR_ENV": "1", "PLY_ACTOR": "planner-a", "PATH": filepath.Dir(herdr.Path) + string(os.PathListSeparator) + os.Getenv("PATH")}
	preview := f.success(t, f.root, "workflow", "goal-start", "started", "--spec", "started-v1", "--revision", "1", "--herdr-workspace", "w1")
	if !strings.Contains(preview, "This Task belongs to PLAN-1 step impl.") {
		t.Fatalf("preview lacks the journal instruction: %s", preview)
	}
	if _, stderr, exit := f.run(t, f.root, env, "workflow", "goal-start", "started", "--spec", "started-v1", "--revision", "1", "--herdr-workspace", "w1", "--apply"); exit != 0 {
		t.Fatalf("goal-start: %d %s", exit, stderr)
	}
	if got := operationsOf(t, plyActions(t, f, "started"), "planner-a"); strings.Join(got, ",") != "goal_register,goal_start_prompt_submitted,goal_start_reserved,plan_link,worktree_create" {
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

// planTraceDoc is the part of PlyPlanTrace@1 the acceptance reads.
type planTraceDoc struct {
	Kind string `json:"kind"`
	Plan struct {
		ID        string `json:"id"`
		Revisions []struct {
			Revision int    `json:"revision"`
			Actor    string `json:"actor"`
		} `json:"revisions"`
	} `json:"plan"`
	Tasks []struct {
		TaskID   string  `json:"task_id"`
		PlanStep *string `json:"plan_step"`
		Role     string  `json:"role"`
	} `json:"tasks"`
	Events []struct {
		TaskID       string  `json:"task_id"`
		Type         string  `json:"type"`
		ActorID      *string `json:"actor_id"`
		Role         *string `json:"role"`
		StepID       *string `json:"step_id"`
		StepKind     *string `json:"step_kind"`
		ParentStepID *string `json:"parent_step_id"`
		PlanStep     *string `json:"plan_step"`
		PlyOperation *string `json:"ply_operation"`
		RecordedAt   *string `json:"recorded_at"`
		ActorBasis   string  `json:"actor_basis"`
	} `json:"events"`
}

// K8: one read gives the plan revisions, its linked Tasks and all their
// journal events with roles, steps and plan step; unlinked Tasks stay out.
func TestPlanTraceReadsLinkedTasksAndTheirJournals(t *testing.T) {
	f, _ := recordingFixture(t)
	registerGoal(t, f, "first", "--plan", "PLAN-1", "--plan-step", "impl", "--actor", "planner-a")
	registerGoal(t, f, "second", "--plan", "PLAN-1", "--plan-step", "review", "--actor", "planner-a")
	registerGoal(t, f, "other", "--actor", "planner-a")
	step := map[string]any{"kind": "ply.workspace.task-journal-event-input", "schema_version": 1, "publication_key": "dev:1", "task_id": "first",
		"actor": map[string]any{"id": "dev-1", "role": "developer", "session_id": nil}, "activity_id": nil, "run_binding": nil,
		"occurred_at": "2026-11-15T10:00:00Z", "time_basis": map[string]any{"kind": "reported", "clock": "agent", "precision": "second", "uncertainty": nil},
		"sources": []any{}, "type": "step_started", "title": "Implement", "detail": "", "step": map[string]any{"id": "impl-1", "kind": "implementation", "parent_step_id": nil},
		"outcome": nil, "candidate": nil, "relations": []any{}, "waiting": nil}
	b, _ := json.Marshal(step)
	input := filepath.Join(f.root, "step.json")
	mustPlanWrite(t, input, string(b), 0600)
	f.success(t, f.root, "workspace", "task", "journal", "record", "first", "--file", input)

	var trace planTraceDoc
	if e := json.Unmarshal([]byte(f.success(t, f.root, "plan", "trace", "PLAN-1", "--project", "example", "--format", "json")), &trace); e != nil {
		t.Fatal(e)
	}
	if trace.Kind != "PlyPlanTrace@1" || trace.Plan.ID != "PLAN-1" || len(trace.Plan.Revisions) != 1 || trace.Plan.Revisions[0].Actor != "planner-a" {
		t.Fatalf("plan %+v", trace)
	}
	tasks := map[string]string{}
	for _, task := range trace.Tasks {
		tasks[task.TaskID] = *task.PlanStep
	}
	if len(tasks) != 2 || tasks["first"] != "impl" || tasks["second"] != "review" {
		t.Fatalf("tasks %+v", trace.Tasks)
	}
	var developer, register bool
	for _, ev := range trace.Events {
		if ev.TaskID == "other" {
			t.Fatalf("unlinked Task in trace: %+v", ev)
		}
		if ev.Type == "step_started" && ev.Role != nil && *ev.Role == "developer" && *ev.StepID == "impl-1" && *ev.StepKind == "implementation" && *ev.PlanStep == "impl" {
			developer = true
		}
		if ev.PlyOperation != nil && *ev.PlyOperation == "goal_register" && ev.TaskID == "second" && *ev.ActorID == "planner-a" && *ev.Role == "ply" && ev.RecordedAt != nil {
			register = true
		}
	}
	if !developer || !register {
		t.Fatalf("events %+v", trace.Events)
	}
	f.failure(t, f.root, nil, "PLAN-9", "plan", "trace", "PLAN-9", "--project", "example", "--format", "json")
}

// K8 budget: a plan with 10 Tasks and 500 events is read within 1 s.
func TestPlanTraceReadsTenTasksAndFiveHundredEventsWithinASecond(t *testing.T) {
	f, _ := recordingFixture(t)
	for i := 0; i < 10; i++ {
		registerGoal(t, f, fmt.Sprintf("load-%d", i), "--plan", "PLAN-1", "--plan-step", "impl", "--actor", "planner-a")
	}
	t.Chdir(f.root)
	journal := taskjournal.New(taskwork.SystemDependencies())
	for i := 0; i < 10; i++ {
		for j := 0; j < 48; j++ { // with the two records of goal register: 50 per Task
			if _, e := journal.AppendPly(fmt.Sprintf("load-%d", i), "dev", "load", taskjournal.PlyAction{Operation: "lifecycle", Refs: map[string]string{"n": fmt.Sprint(j)}}); e != nil {
				t.Fatal(e)
			}
		}
	}
	start := time.Now()
	var trace planTraceDoc
	if e := json.Unmarshal([]byte(f.success(t, f.root, "plan", "trace", "PLAN-1", "--project", "example", "--format", "json")), &trace); e != nil {
		t.Fatal(e)
	}
	elapsed := time.Since(start)
	ply := 0
	for _, ev := range trace.Events {
		if ev.PlyOperation != nil {
			ply++
		}
	}
	t.Logf("plan trace: %d Tasks, %d events (%d ply) in %s", len(trace.Tasks), len(trace.Events), ply, elapsed)
	if len(trace.Tasks) != 10 || ply != 500 || elapsed > time.Second {
		t.Fatalf("%d Tasks, %d ply events in %s", len(trace.Tasks), ply, elapsed)
	}
}

// An existing Task is linked afterwards; step commands need no JSON file and
// a correction inside an implementation keeps its parent (K7).
func TestJournalLinkAndStepCommands(t *testing.T) {
	f, _ := recordingFixture(t)
	registerGoal(t, f, "later", "--actor", "planner-a")
	f.failure(t, f.root, nil, "--plan-step nope is not a step of PLAN-1", "workspace", "task", "journal", "link", "later", "--plan", "PLAN-1", "--plan-step", "nope")
	f.success(t, f.root, "workspace", "task", "journal", "link", "later", "--plan", "PLAN-1", "--plan-step", "impl", "--actor", "planner-a")
	journal := []string{"workspace", "task", "journal"}
	f.success(t, f.root, append(journal, "step", "start", "later", "--id", "impl", "--kind", "implementation", "--title", "Implement", "--actor", "dev-1")...)
	f.success(t, f.root, append(journal, "step", "start", "later", "--id", "fix-1", "--kind", "correction", "--parent", "impl", "--title", "Fix finding", "--actor", "dev-1")...)
	f.success(t, f.root, append(journal, "step", "finish", "later", "--id", "fix-1", "--kind", "correction", "--parent", "impl", "--title", "Fix finding", "--outcome", "complete", "--actor", "dev-1")...)
	f.success(t, f.root, append(journal, "step", "finish", "later", "--id", "impl", "--kind", "implementation", "--title", "Implement", "--outcome", "complete", "--actor", "dev-1")...)
	f.success(t, f.root, append(journal, "note", "later", "--title", "Returned", "--role", "developer", "--actor", "dev-1")...)
	var trace planTraceDoc
	if e := json.Unmarshal([]byte(f.success(t, f.root, "plan", "trace", "PLAN-1", "--project", "example", "--format", "json")), &trace); e != nil {
		t.Fatal(e)
	}
	if len(trace.Tasks) != 1 || trace.Tasks[0].TaskID != "later" || *trace.Tasks[0].PlanStep != "impl" {
		t.Fatalf("tasks %+v", trace.Tasks)
	}
	var child, note bool
	for _, ev := range trace.Events {
		if ev.StepID != nil && *ev.StepID == "fix-1" && ev.ParentStepID != nil && *ev.ParentStepID == "impl" {
			child = true
		}
		note = note || ev.Type == "note" && *ev.ActorID == "dev-1"
	}
	if !child || !note {
		t.Fatalf("events %+v", trace.Events)
	}
}

// K6: Epic -> main is recorded on the Tasks it brings to the target for the
// first time, with the exact before and after commits.
func TestEpicIntegrateRecordsOnNewlyIntegratedTasks(t *testing.T) {
	f, oid := recordingFixture(t)
	repo := filepath.Join(f.root, "repo")
	f.git(t, repo, "branch", "main", oid)
	env := map[string]string{"PLY_ACTOR": "planner-a"}
	registerGoal(t, f, "shipped", "--plan", "PLAN-1", "--actor", "planner-a")
	f.success(t, f.root, "workspace", "task", "worktree", "create", "shipped", "--branch", "shipped", "--path", filepath.Join(f.root, "shipped"), "--expected-parent-oid", oid)
	candidate := integrateTask(t, f, "shipped", oid, env)
	if _, stderr, exit := f.run(t, f.root, env, "workspace", "epic", "integrate", "work", "--repo", "code", "--target-ref", "refs/heads/main", "--expected-target-oid", oid, "--expected-epic-oid", candidate, "--apply"); exit != 0 {
		t.Fatalf("epic integrate: %d %s", exit, stderr)
	}
	var epic *taskjournal.EventInput
	for _, ev := range plyActions(t, f, "shipped") {
		var in taskjournal.EventInput
		if e := json.Unmarshal(ev.Data, &in); e == nil && in.Ply != nil && in.Ply.Operation == "epic_integrate" {
			epic = &in
		}
	}
	if epic == nil || epic.Actor.ID != "planner-a" || epic.Ply.Refs["before_oid"] != oid || epic.Ply.Refs["after_oid"] != candidate || epic.Ply.Refs["effect"] != "fast_forwarded" || epic.Ply.Refs["target_ref"] != "refs/heads/main" {
		t.Fatalf("epic integration record %+v", epic)
	}
}

// Integration retires the Task worktree and branch unless they are kept; a
// blocked integration leaves them alone and says nothing about cleanup.
func TestTaskIntegrateRetiresTheWorktree(t *testing.T) {
	for _, c := range []struct {
		task       string
		extra      []string
		locked     bool
		integrated bool
		gone       bool
	}{
		{"retired", nil, false, true, true},
		{"kept", []string{"--keep-worktree"}, false, true, false},
		{"blocked", nil, true, false, false},
	} {
		f, parent := recordingFixture(t)
		repo := filepath.Join(f.root, "repo")
		registerGoal(t, f, c.task, "--actor", "planner-a")
		f.success(t, f.root, "workspace", "task", "worktree", "create", c.task, "--branch", c.task, "--path", filepath.Join(f.root, c.task), "--expected-parent-oid", parent)
		worktree := filepath.Join(f.root, c.task)
		mustPlanWrite(t, filepath.Join(worktree, c.task+".txt"), c.task+"\n", 0600)
		f.git(t, worktree, "add", c.task+".txt")
		f.git(t, worktree, "commit", "-m", c.task)
		candidate := f.git(t, worktree, "rev-parse", "HEAD")
		if c.locked { // a locked worktree blocks the integration
			f.git(t, repo, "worktree", "lock", worktree)
		}
		integrate := []string{"workspace", "task", "integrate", c.task, "--expected-result-oid", candidate, "--expected-parent-oid", parent, "--format", "json"}
		var check struct {
			Plan struct {
				SHA256 string `json:"sha256"`
			} `json:"plan"`
		}
		if e := json.Unmarshal([]byte(f.success(t, f.root, append(integrate, "--check")...)), &check); e != nil {
			t.Fatal(e)
		}
		_, stderr, exit := f.run(t, f.root, nil, append(append(integrate, "--apply", "--confirm", check.Plan.SHA256), c.extra...)...)
		if exit != 0 || (f.git(t, repo, "rev-parse", "refs/heads/epic") == candidate) != c.integrated {
			t.Fatalf("%s: integration %d %s", c.task, exit, stderr)
		}
		_, statErr := os.Stat(worktree)
		branches := f.git(t, repo, "branch", "--list", c.task)
		if gone := os.IsNotExist(statErr) && branches == ""; gone != c.gone {
			t.Fatalf("%s: worktree gone %v (branch %q), want %v; stderr %s", c.task, gone, branches, c.gone, stderr)
		}
		if !c.integrated && strings.Contains(stderr, "Worktree") {
			t.Fatalf("%s: cleanup attempted after a blocked integration: %s", c.task, stderr)
		}
	}
}

// Epic -> main retires the Epic worktree and branch and archives the Epic;
// --keep-worktree keeps them.
func TestEpicIntegrateRetiresTheEpic(t *testing.T) {
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("opt-in CLI fixture: set PLY_CANDIDATE_BINARY to the actual make-built candidate")
	}
	for _, keep := range []bool{false, true} {
		f := newPlanCLI(t, binary)
		repo := filepath.Join(f.root, "repo")
		if e := os.Mkdir(repo, 0700); e != nil {
			t.Fatal(e)
		}
		f.git(t, repo, "init", "-b", "main")
		mustPlanWrite(t, filepath.Join(repo, "README.md"), "fixture\n", 0600)
		f.git(t, repo, "add", "README.md")
		f.git(t, repo, "commit", "-m", "fixture")
		base := f.git(t, repo, "rev-parse", "HEAD")
		epic := filepath.Join(f.root, "epic")
		f.git(t, repo, "worktree", "add", "-b", "epic", epic, base)
		f.success(t, f.root, "workspace", "init")
		mustPlanWrite(t, filepath.Join(f.root, "ply-workspace.json"), `{"format_version": 1, "projects": [{"id": "example", "name": "Example", "wrapper": ".",
  "repos": [{"id": "code", "path": "repo"}]}]}
`, 0600)
		f.success(t, f.root, "workspace", "epic", "adopt", "work", "--title", "Work", "--project", "example", "--repo", "code", "--worktree", epic, "--ref", "refs/heads/epic", "--expected-oid", base)
		mustPlanWrite(t, filepath.Join(epic, "epic.txt"), "epic\n", 0600)
		f.git(t, epic, "add", "epic.txt")
		f.git(t, epic, "commit", "-m", "epic work")
		head := f.git(t, epic, "rev-parse", "HEAD")
		args := []string{"workspace", "epic", "integrate", "work", "--repo", "code", "--target-ref", "refs/heads/main", "--expected-target-oid", base, "--expected-epic-oid", head, "--apply"}
		if keep {
			args = append(args, "--keep-worktree")
		}
		_, stderr, exit := f.run(t, f.root, map[string]string{"PLY_ACTOR": "planner-a"}, args...)
		if exit != 0 || f.git(t, repo, "rev-parse", "refs/heads/main") != head {
			t.Fatalf("keep=%v: integrate %d %s", keep, exit, stderr)
		}
		_, statErr := os.Stat(epic)
		lifecycle := f.success(t, f.root, "workspace", "epic", "lifecycle", "show", "work")
		archived := strings.Contains(lifecycle, "archived")
		if gone := os.IsNotExist(statErr) && f.git(t, repo, "branch", "--list", "epic") == ""; gone == keep || archived == keep {
			t.Fatalf("keep=%v: worktree gone %v, archived %v; stderr %s; lifecycle %s", keep, gone, archived, stderr, lifecycle)
		}
	}
}

// The plan's owner records steps and notes on the plan itself; plan trace
// shows them with an empty task_id among the Task events.
func TestPlanJournalRecordsPlanStepsInTheTrace(t *testing.T) {
	f, _ := recordingFixture(t)
	journal := []string{"plan", "journal"}
	f.failure(t, f.root, nil, "kind", append(journal, "step", "start", "PLAN-1", "--project", "example", "--id", "map", "--kind", "nope", "--title", "Map")...)
	f.success(t, f.root, append(journal, "step", "start", "PLAN-1", "--project", "example", "--id", "map", "--kind", "planning", "--title", "Map", "--actor", "ply-plan")...)
	f.success(t, f.root, append(journal, "step", "finish", "PLAN-1", "--project", "example", "--id", "map", "--kind", "planning", "--title", "Map", "--outcome", "complete", "--actor", "ply-plan")...)
	f.success(t, f.root, append(journal, "note", "PLAN-1", "--project", "example", "--title", "User chose option 1", "--actor", "ply-plan")...)
	f.success(t, f.root, append(journal, "note", "PLAN-1", "--project", "example", "--title", "Reviewed", "--actor", "rev-1", "--role", "reviewer")...)
	// The raw trace keeps each event's at, actor and role for the failure message
	// (backlog 20261121-058d5a58: the last event was seen once out of order).
	raw := f.success(t, f.root, "plan", "trace", "PLAN-1", "--project", "example", "--format", "json")
	var trace planTraceDoc
	if e := json.Unmarshal([]byte(raw), &trace); e != nil {
		t.Fatal(e)
	}
	var types []string
	for _, ev := range trace.Events {
		if ev.TaskID == "" && *ev.Role == "planner" && *ev.ActorID == "ply-plan" {
			types = append(types, ev.Type)
		}
	}
	if !reflect.DeepEqual(types, []string{"step_started", "step_finished", "note"}) {
		t.Fatalf("plan events %v; trace %s", types, raw)
	}
	last := trace.Events[len(trace.Events)-1]
	if last.TaskID != "" || *last.ActorID != "rev-1" || *last.Role != "reviewer" {
		t.Fatalf("last plan event is not rev-1; trace %s", raw)
	}
}

func TestTaskCleanBuildEmptiesTheRegisteredWorktreeBuild(t *testing.T) {
	f, oid := recordingFixture(t)
	registerGoal(t, f, "cleaned")
	f.failure(t, f.root, nil, "has no worktree", "workspace", "task", "clean-build", "cleaned")
	worktree := filepath.Join(f.root, "cleaned")
	f.success(t, f.root, "workspace", "task", "worktree", "create", "cleaned", "--branch", "cleaned", "--path", worktree, "--expected-parent-oid", oid)
	if e := os.MkdirAll(filepath.Join(worktree, ".build", "out"), 0700); e != nil {
		t.Fatal(e)
	}
	mustPlanWrite(t, filepath.Join(worktree, ".build", "out", "x"), "x\n", 0600)
	if out := f.success(t, f.root, "workspace", "task", "clean-build", "cleaned"); !strings.Contains(out, "Removed 1 entries") {
		t.Fatalf("output %q", out)
	}
	if entries, e := os.ReadDir(filepath.Join(worktree, ".build")); e != nil || len(entries) != 0 {
		t.Fatalf(".build %v %v", entries, e)
	}
}

// A damaged journal does not stop goal-start of a Task without plan link;
// the valid prefix is read for the link.
func TestGoalStartReadsTheValidPrefixOfADamagedJournal(t *testing.T) {
	f, oid := recordingFixture(t)
	registerGoal(t, f, "damaged")
	f.success(t, f.root, "workspace", "task", "worktree", "create", "damaged", "--branch", "damaged", "--path", filepath.Join(f.root, "damaged"), "--expected-parent-oid", oid)
	works, e := workspace.DefinedWork(f.root)
	if e != nil || len(works) != 1 {
		t.Fatalf("works %v %v", works, e)
	}
	if e := os.Remove(filepath.Join(taskjournal.JournalDir(works[0]), "damaged", "000000001.json")); e != nil {
		t.Fatal(e)
	}
	preview := f.success(t, f.root, "workflow", "goal-start", "damaged", "--spec", "damaged-v1", "--revision", "1", "--herdr-workspace", "w1")
	if strings.Contains(preview, "journal_instruction") {
		t.Fatalf("instruction without plan link: %s", preview)
	}
}
