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

func TestWorkspaceProjectInitCommandFlags(t *testing.T) {
	c := newWorkspaceProjectCommand(workspace.SystemDependencies())
	init, _, err := c.Find([]string{"init"})
	if err != nil || init.Name() != "init" {
		t.Fatal("missing init command", err)
	}
	for _, name := range []string{"wrapper", "name", "apply", "confirmation", "format"} {
		if init.Flags().Lookup(name) == nil {
			t.Fatal("missing flag", name)
		}
	}
}

// Synthetic fixtures only; can also exercise an installed CLI via the same env.
func TestProjectInitCandidateCLI(t *testing.T) {
	binary := os.Getenv("PLY_CANDIDATE_BINARY")
	if binary == "" {
		t.Skip("set PLY_CANDIDATE_BINARY to the make-built or installed CLI; required in native acceptance")
	}
	f := newPlanCLI(t, binary)
	success := func(cwd string, args ...string) string {
		t.Helper()
		stdout, stderr, exit := f.run(t, cwd, nil, args...)
		t.Logf("argv=%q cwd=%q stdout=%q stderr=%q exit=%d", append([]string{binary}, args...), cwd, stdout, stderr, exit)
		if exit != 0 {
			t.Fatalf("candidate command failed: %s", stderr)
		}
		return stdout
	}

	w := filepath.Join(f.root, "example")
	main := filepath.Join(w, "main")
	if err := os.MkdirAll(main, 0700); err != nil {
		t.Fatal(err)
	}
	f.git(t, main, "init", "-b", "main")
	f.git(t, main, "commit", "--allow-empty", "-m", "fixture")
	f.git(t, main, "worktree", "add", "-b", "feature", filepath.Join(w, "feature"))
	mustPlanWrite(t, filepath.Join(main, "dirty.txt"), "preserved dirty source\n", 0600)
	success(f.root, "workspace", "init")
	before := mr05Snapshot(t, f.root)
	args := []string{"workspace", "project", "init", "example", "--format", "json"}
	decode := func(raw string) workspace.ProjectInitPlan {
		t.Helper()
		var p workspace.ProjectInitPlan
		if err := json.Unmarshal([]byte(raw), &p); err != nil {
			t.Fatal(err, raw)
		}
		return p
	}
	p := decode(success(w, args...))
	if p.Applied || p.Wrapper != w || p.Name != "example" || len(p.Repositories) != 1 || p.Repositories[0].Locator != main || p.Confirmation == "" {
		t.Fatalf("preview %#v", p)
	}
	if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
		t.Fatal("preview wrote")
	}
	bad := append(append([]string{}, args...), "--apply", "--confirmation", "wrong")
	stdout, stderr, exit := f.run(t, w, nil, bad...)
	t.Logf("argv=%q cwd=%q stdout=%q stderr=%q exit=%d", append([]string{binary}, bad...), w, stdout, stderr, exit)
	if exit == 0 || !strings.Contains(stderr, "confirmation") || !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
		t.Fatal("invalid apply changed state", stderr)
	}
	apply := append(append([]string{}, args...), "--apply", "--confirmation", p.Confirmation)
	applied := decode(success(w, apply...))
	if !applied.Applied {
		t.Fatal("not applied")
	}
	success(w, "workspace", "project", "show", "example", "--format", "json")
	after := mr05Snapshot(t, f.root)
	if !decode(success(w, apply...)).Applied || !reflect.DeepEqual(after, mr05Snapshot(t, f.root)) {
		t.Fatal("retry changed files")
	}
	for path, value := range before {
		if after[path] != value {
			t.Fatalf("existing file changed: %s", path)
		}
	}
	text := success(w, "workspace", "project", "init", "example")
	if !strings.Contains(text, main) || !strings.Contains(text, "git common directory") || !strings.Contains(text, "--confirmation") {
		t.Fatal("text preview incomplete", text)
	}
	// A new direct repo changes membership; the old Project remains exact.
	extra := filepath.Join(w, "extra")
	if err := os.Mkdir(extra, 0700); err != nil {
		t.Fatal(err)
	}
	f.git(t, extra, "init", "-b", "main")
	changed := mr05Snapshot(t, f.root)
	stdout, stderr, exit = f.run(t, w, nil, apply...)
	t.Logf("argv=%q cwd=%q stdout=%q stderr=%q exit=%d", append([]string{binary}, apply...), w, stdout, stderr, exit)
	if exit == 0 || !strings.Contains(stderr, "project add-member") || !reflect.DeepEqual(changed, mr05Snapshot(t, f.root)) {
		t.Fatal("existing membership conflict lost state", stderr)
	}
}
