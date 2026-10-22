package cmd

import (
	"encoding/json"
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
	args := []string{"workspace", "task", "goal", "register", "supplied", "--spec", "supplied-v1", "--document", doc, "--extra-document", "extra=" + extra, "--title", "Supplied product Goal", "--objective", "Observe supplied behavior", "--criterion", "The requested behavior is observable", "--project", "example", "--repo", "code", "--epic", "work", "--delivery", workspace.DeliveryLocalEpic, "--acceptance", "automatic", "--responsible-actor", "fixture operator", "--acceptance-timeout", "30", "--format", "json"}
	decode := func(raw string) workspace.FirstGoalResult {
		t.Helper()
		var out workspace.FirstGoalResult
		if e := json.Unmarshal([]byte(raw), &out); e != nil {
			t.Fatalf("invalid result: %v\n%s", e, raw)
		}
		return out
	}
	before := mr05Snapshot(t, f.root)
	preview := decode(f.success(t, f.root, args...))
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
