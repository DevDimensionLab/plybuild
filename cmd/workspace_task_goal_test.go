package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Uses only the supplied make-built CLI and ordinary public registration
// commands. No provider, process inventory, Herdr or sandbox is exercised.
func TestFirstGoalCandidateCLIRegistrationAndReadOnlyReadiness(t *testing.T) {
	firstGoalCLIJourney(t, false)
}

func TestExistingFirstGoalCandidateCLIRegistrationAndReadOnlyReadiness(t *testing.T) {
	firstGoalCLIJourney(t, true)
}

func firstGoalCLIJourney(t *testing.T, existing bool) {
	t.Helper()
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("opt-in CLI fixture: set PLY_CANDIDATE_BINARY to the actual make-built candidate; native acceptance requires it")
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
	doc := filepath.Join(f.root, "spec.md")
	body := "# Supplied Spec\nExact UTF-8: blåbær.\n"
	mustPlanWrite(t, doc, body, 0600)
	extra := filepath.Join(f.root, "design.txt")
	mustPlanWrite(t, extra, "Extra supplied design\n", 0600)
	args := []string{"workspace", "task", "goal", "register", "supplied", "--spec", "supplied-v1", "--document", doc, "--extra-document", "extra=" + extra, "--title", "Supplied product Goal", "--objective", "Observe supplied behavior", "--criterion", "The requested behavior is observable", "--project", "example", "--repo", "code", "--epic", "work", "--format", "json"}
	decode := func(raw string) workspace.FirstGoalResult {
		t.Helper()
		var out workspace.FirstGoalResult
		if e := json.Unmarshal([]byte(raw), &out); e != nil {
			t.Fatalf("invalid result: %v\n%s", e, raw)
		}
		return out
	}
	var preservedProblem *workspace.TaskRevisionRef
	var preservedTask string
	if existing {
		f.success(t, f.root, "workspace", "task", "create", "supplied", "--title", "Imported backlog label", "--description", "Historical imported product context", "--project", "example", "--repo", "code", "--epic", "work")
		preservedTask = f.success(t, f.root, "workspace", "task", "show", "supplied", "--format", "json")
		// The Task readback supplies its native current Problem reference.
		var taskReadback struct {
			Problem struct {
				Head *workspace.TaskRevisionRef `json:"head"`
			} `json:"problem"`
		}
		if e := json.Unmarshal([]byte(preservedTask), &taskReadback); e != nil {
			t.Fatal(e)
		}
		head := taskReadback.Problem.Head
		if head == nil || head.Revision != 1 {
			t.Fatalf("missing initial Problem head: %s", preservedTask)
		}
		problemDoc := filepath.Join(f.root, "imported-problem.md")
		mustPlanWrite(t, problemDoc, "# Imported backlog problem\nPreserve the imported context.\n", 0600)
		// Native publisher reads and hashes a source descriptor supplied by this isolated fixture.
		b, e := os.ReadFile(problemDoc)
		if e != nil {
			t.Fatal(e)
		}
		draft := map[string]any{"kind": "WorkspaceTaskProblemDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "fixture/imported-problem-2", "task_id": "supplied", "registry_upgrade": nil, "expected_previous": head, "origin": map[string]any{"kind": "authored"}, "title": "Imported backlog problem", "summary": "Preserve historical imported context", "problem_document_id": "problem", "documents": []any{map[string]any{"id": "problem", "source": map[string]any{"kind": "file", "locator": problemDoc, "sha256": fmt.Sprintf("sha256:%x", sha256.Sum256(b)), "size_bytes": len(b), "media_type": "text/markdown", "git_provenance": nil}}}, "sources": []any{}, "claims": []any{}, "deadline": nil, "change_reason": "Import the preserved product context", "recorder": map[string]any{"actor_claim": "isolated fixture", "control_surface": "CLI test", "recorded_at_utc": "2026-10-22T00:00:00Z"}}
		draftBytes, e := json.Marshal(draft)
		if e != nil {
			t.Fatal(e)
		}
		draftPath := filepath.Join(f.root, "problem-draft.json")
		mustPlanWrite(t, draftPath, string(draftBytes), 0600)
		f.success(t, f.root, "workspace", "task", "problem", "record", "supplied", "--file", draftPath, "--format", "json")
	}
	before := mr05Snapshot(t, f.root)
	preview := decode(f.success(t, f.root, args...))
	if existing {
		preservedProblem = preview.Problem
		if preservedProblem == nil || preservedProblem.Revision != 2 || len(preview.Effects) != 2 {
			t.Fatalf("existing Task preview lost basis/effects: %#v", preview)
		}
	}
	if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) || preview.Applied || preview.Run != nil || preview.CheckCwd != repo || len(preview.Sources) != 2 {
		t.Fatal("preview wrote files or lost bindings")
	}
	invalid := append(append([]string{}, args...), "--apply", "--confirmation", "wrong")
	_, invalidError, invalidExit := f.run(t, f.root, nil, invalid...)
	if invalidExit == 0 || !strings.Contains(invalidError, "confirmation") {
		t.Fatal("wrong confirmation was accepted", invalidError)
	}
	if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
		t.Fatal("invalid confirmation wrote state")
	}
	apply := append(append([]string{}, args...), "--apply", "--confirmation", preview.Confirmation)
	out := decode(f.success(t, f.root, apply...))
	if !out.Applied || out.Goal.Spec.Revision != 1 || out.Goal.Spec.ManifestSHA256 == "" || out.Contract == nil || len(out.Contract.RequiredInputs) != 5 || out.Run != nil || out.Readiness != "not_checked" {
		t.Fatalf("missing exact registered refs: %#v", out)
	}
	if existing {
		if !reflect.DeepEqual(out.Problem, preservedProblem) || out.Contract.Problem != *preservedProblem {
			t.Fatal("Goal replaced current Problem")
		}
		// Task identity remains visible with the original imported title and description.
		afterTask := f.success(t, f.root, "workspace", "task", "show", "supplied", "--format", "json")
		for _, label := range []string{"Imported backlog label", "Historical imported product context"} {
			if !strings.Contains(preservedTask, label) || !strings.Contains(afterTask, label) {
				t.Fatal("imported identity was replaced")
			}
		}
	}
	if len(out.Queue.Pending) != 1 || out.Queue.Pending[0].Goal == nil || *out.Queue.Pending[0].Goal != out.Goal {
		t.Fatal("published Goal was not queued exactly")
	}
	if out.Contract.Executor.Provider != nil || out.Contract.Executor.Model != nil || out.Contract.Executor.Effort != nil {
		t.Fatal("registration invented executor settings")
	}
	if got, e := os.ReadFile(doc); e != nil || string(got) != body {
		t.Fatal("source document changed")
	}
	registered := mr05Snapshot(t, f.root)
	again := decode(f.success(t, f.root, apply...))
	if again.Goal != out.Goal || !reflect.DeepEqual(registered, mr05Snapshot(t, f.root)) {
		t.Fatal("CLI retry duplicated native state")
	}
	f.success(t, f.root, "workspace", "task", "spec", "show", "supplied", "--spec", out.Goal.SpecID, "--revision", "1", "--format", "json")
	// Dirty return blocks the existing metadata-only execute preview before any
	// runtime inspection. Registration retains its successful result and Run=null.
	mustPlanWrite(t, filepath.Join(repo, "dirty.txt"), "dirty return\n", 0600)
	beforeCheck := mr05Snapshot(t, f.root)
	check := append([]string{}, out.NextArgv[1:]...)
	_, stderr, exit := f.run(t, repo, nil, check...)
	if exit == 0 || !strings.Contains(strings.ToLower(stderr), "clean") {
		t.Fatalf("expected honest dirty-return readiness blocker, got %d %s", exit, stderr)
	}
	if !reflect.DeepEqual(beforeCheck, mr05Snapshot(t, f.root)) {
		t.Fatal("read-only execute check wrote files")
	}
	historical := decode(f.success(t, f.root, apply...))
	if !historical.Applied || historical.Goal != out.Goal || historical.Run != nil {
		t.Fatal("readiness blocker erased registration success")
	}
}

func TestDocumentFreeGoalRegistersAndStartsWithWholeAssignment(t *testing.T) {
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
	args := []string{"workspace", "task", "goal", "register", "plain", "--spec", "plain-v1", "--title", "Plain title", "--objective", "Plain objective", "--criterion", "Plain criterion", "--constraint", "Plain human limit", "--project", "example", "--repo", "code", "--epic", "work", "--provider", "claude", "--format", "json"}

	before := mr05Snapshot(t, f.root)
	blank := append(append([]string{}, args...), "--document", "")
	if _, stderr, exit := f.run(t, f.root, nil, blank...); exit == 0 || !strings.Contains(stderr, "absolute physical file paths") {
		t.Fatalf("explicit empty --document silently dropped: %d %s", exit, stderr)
	}
	var preview, applied workspace.FirstGoalResult
	if e := json.Unmarshal([]byte(f.success(t, f.root, args...)), &preview); e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) || preview.Applied || len(preview.Sources) != 0 {
		t.Fatalf("preview wrote state or invented sources: %#v", preview)
	}
	apply := append(append([]string{}, args...), "--apply", "--confirmation", preview.Confirmation)
	if e := json.Unmarshal([]byte(f.success(t, f.root, apply...)), &applied); e != nil {
		t.Fatal(e)
	}
	if !applied.Applied || applied.Contract == nil || len(applied.Queue.Pending) != 1 {
		t.Fatalf("not registered: %#v", applied)
	}
	worktree := filepath.Join(f.root, "plain")
	f.success(t, f.root, "workspace", "task", "worktree", "create", "plain", "--branch", "plain", "--path", worktree, "--expected-parent-oid", oid)

	herdr, log := fakeHerdrBinary(t)
	env := map[string]string{"HERDR_ENV": "1", "PATH": filepath.Dir(herdr.Path) + string(os.PathListSeparator) + os.Getenv("PATH")}
	start := []string{"workflow", "goal-start", "plain", "--spec", "plain-v1", "--revision", "1", "--herdr-workspace", "w1"}
	if _, stderr, exit := f.run(t, f.root, env, append(start, "--apply")...); exit != 0 {
		t.Fatalf("goal-start failed: %d %s", exit, stderr)
	}
	calls, e := os.ReadFile(log)
	if e != nil {
		t.Fatal(e)
	}
	if strings.Count(string(calls), "agent prompt") != 1 {
		t.Fatalf("expected one prompt: %s", calls)
	}
	for _, want := range []string{"Plain title", "Plain objective", "Plain criterion", "Plain human limit", `"provider":"claude"`, applied.Goal.Spec.ManifestSHA256} {
		if !strings.Contains(string(calls), want) {
			t.Fatalf("native start lost %q: %s", want, calls)
		}
	}
	if strings.Contains(string(calls), `"delivery"`) {
		t.Fatalf("native start carried a delivery policy: %s", calls)
	}
	if strings.Contains(string(calls), "task-goal-doc/goal/") {
		t.Fatalf("native start invented a Goal document: %s", calls)
	}
}

func TestFirstGoalAfterNativeWorktreeCandidateCLI(t *testing.T) {
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
	f.success(t, f.root, "workspace", "task", "create", "bound", "--title", "Bound first", "--description", "Worktree before Goal", "--project", "example", "--repo", "code", "--epic", "work")
	problemDoc := filepath.Join(f.root, "bound-problem.md")
	mustPlanWrite(t, problemDoc, "# Bound problem\nRegister a Goal after the worktree.\n", 0600)
	problemBytes, e := os.ReadFile(problemDoc)
	if e != nil {
		t.Fatal(e)
	}
	draft, e := json.Marshal(map[string]any{"kind": "WorkspaceTaskProblemDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "fixture/bound-problem-1", "task_id": "bound", "registry_upgrade": nil, "expected_previous": nil, "origin": map[string]any{"kind": "authored"}, "title": "Bound problem", "summary": "Register a Goal after the worktree", "problem_document_id": "problem", "documents": []any{map[string]any{"id": "problem", "source": map[string]any{"kind": "file", "locator": problemDoc, "sha256": fmt.Sprintf("sha256:%x", sha256.Sum256(problemBytes)), "size_bytes": len(problemBytes), "media_type": "text/markdown", "git_provenance": nil}}}, "sources": []any{}, "claims": []any{}, "deadline": nil, "change_reason": "Record the bound Task problem", "recorder": map[string]any{"actor_claim": "isolated fixture", "control_surface": "CLI test", "recorded_at_utc": "2026-11-06T00:00:00Z"}})
	if e != nil {
		t.Fatal(e)
	}
	draftPath := filepath.Join(f.root, "bound-problem.json")
	mustPlanWrite(t, draftPath, string(draft), 0600)
	f.success(t, f.root, "workspace", "task", "problem", "record", "bound", "--file", draftPath, "--format", "json")
	worktree := filepath.Join(f.root, "bound")
	f.success(t, f.root, "workspace", "task", "worktree", "create", "bound", "--branch", "bound", "--path", worktree, "--expected-parent-oid", oid)
	binding := func() string {
		t.Helper()
		var shown struct {
			Resource struct {
				Persisted struct {
					Worktree      json.RawMessage `json:"worktree"`
					WorktreeState string          `json:"worktree_state"`
				} `json:"persisted"`
			} `json:"resource"`
		}
		if e := json.Unmarshal([]byte(f.success(t, f.root, "workspace", "task", "show", "bound", "--format", "json")), &shown); e != nil {
			t.Fatal(e)
		}
		p := shown.Resource.Persisted
		if p.WorktreeState != "worktree_ready" || len(p.Worktree) == 0 || string(p.Worktree) == "null" {
			t.Fatalf("Task has no ready native binding: %s %s", p.WorktreeState, p.Worktree)
		}
		return p.WorktreeState + string(p.Worktree)
	}
	bound := binding()

	args := []string{"workspace", "task", "goal", "register", "bound", "--spec", "bound-v1", "--title", "Bound title", "--objective", "Bound objective", "--criterion", "Bound criterion", "--project", "example", "--repo", "code", "--epic", "work", "--provider", "claude", "--format", "json"}
	before := mr05Snapshot(t, f.root)
	var preview, applied workspace.FirstGoalResult
	if e := json.Unmarshal([]byte(f.success(t, f.root, args...)), &preview); e != nil {
		t.Fatal(e)
	}
	if preview.Applied || !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
		t.Fatal("register preview wrote state")
	}
	apply := append(append([]string{}, args...), "--apply", "--confirmation", preview.Confirmation)
	if e := json.Unmarshal([]byte(f.success(t, f.root, apply...)), &applied); e != nil {
		t.Fatal(e)
	}
	if !applied.Applied || applied.Contract == nil || len(applied.Queue.Pending) != 1 || binding() != bound {
		t.Fatalf("registration lost Goal or binding: %#v", applied)
	}
	registered := mr05Snapshot(t, f.root)
	if f.success(t, f.root, apply...); !reflect.DeepEqual(registered, mr05Snapshot(t, f.root)) {
		t.Fatal("register retry duplicated native state")
	}

	herdr, log := fakeHerdrBinary(t)
	env := map[string]string{"HERDR_ENV": "1", "PATH": filepath.Dir(herdr.Path) + string(os.PathListSeparator) + os.Getenv("PATH")}
	start := []string{"workflow", "goal-start", "bound", "--spec", "bound-v1", "--revision", "1", "--herdr-workspace", "w1"}
	out, stderr, exit := f.run(t, f.root, env, start...)
	if exit != 0 {
		t.Fatalf("goal-start preview failed: %d %s", exit, stderr)
	}
	var plan goalAgentPlan
	if e := json.Unmarshal([]byte(out), &plan); e != nil {
		t.Fatal(e)
	}
	if _, e := os.Stat(log); !os.IsNotExist(e) {
		t.Fatal("goal-start preview called Herdr")
	}
	startOID := plan.OID
	if startOID != oid || plan.Worktree != worktree || plan.Ref != "refs/heads/bound" {
		t.Fatalf("preview lost the bound start: %#v", plan)
	}
	for i := 0; i < 2; i++ {
		if _, stderr, exit := f.run(t, f.root, env, append(start, "--apply")...); exit != 0 {
			t.Fatalf("goal-start apply %d failed: %d %s", i, exit, stderr)
		}
	}
	calls, e := os.ReadFile(log)
	if e != nil {
		t.Fatal(e)
	}
	if strings.Count(string(calls), "agent prompt") != 1 || !strings.Contains(string(calls), applied.Goal.Spec.ManifestSHA256) {
		t.Fatalf("expected one prompt for the registered Goal: %s", calls)
	}

	mustPlanWrite(t, filepath.Join(worktree, "result.txt"), "result\n", 0600)
	f.git(t, worktree, "add", "result.txt")
	f.git(t, worktree, "commit", "-m", "goal result")
	candidate := f.git(t, worktree, "rev-parse", "HEAD")
	if candidate == startOID {
		t.Fatal("fixture commit did not advance the Task")
	}
	integrate := []string{"workspace", "task", "integrate", "bound", "--expected-result-oid", candidate, "--expected-parent-oid", startOID, "--format", "json"}
	var check struct {
		Classification string `json:"classification"`
		Plan           struct {
			SHA256 string `json:"sha256"`
		} `json:"plan"`
	}
	if e := json.Unmarshal([]byte(f.success(t, f.root, append(integrate, "--check")...)), &check); e != nil {
		t.Fatal(e)
	}
	if check.Classification != "ready" || check.Plan.SHA256 == "" || f.git(t, repo, "rev-parse", "refs/heads/epic") != startOID {
		t.Fatalf("integration check: %#v", check)
	}
	f.success(t, f.root, append(integrate, "--apply", "--confirm", check.Plan.SHA256)...)
	if got := f.git(t, repo, "rev-parse", "refs/heads/epic"); got != candidate {
		t.Fatalf("Epic ref = %s, want candidate %s", got, candidate)
	}
}
