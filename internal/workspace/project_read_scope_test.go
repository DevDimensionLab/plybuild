package workspace

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// gitCalls puts a git in PATH that records the arguments of each call and
// runs the real git, and returns a reader of the recorded calls.
func gitCalls(t *testing.T) func() []string {
	t.Helper()
	real, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	log := filepath.Join(bin, "calls")
	script := "#!/bin/sh\necho \"$*\" >> '" + log + "'\nexec '" + real + "' \"$@\"\n"
	if err := os.WriteFile(filepath.Join(bin, "git"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return func() []string {
		b, err := os.ReadFile(log)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			t.Fatal(err)
		}
		calls := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
		if err := os.Remove(log); err != nil {
			t.Fatal(err)
		}
		return calls
	}
}

// A read pays only for what its answer holds: the list has no Git facts, and
// show observes only the members of its own Project, so a broken Project
// stops neither.
func TestProjectReadsObserveOnlyWhatTheyShow(t *testing.T) {
	f := newDefinitionFixture(t)
	broken := strings.Replace(fixtureDefinitions, "\n  ]\n", `,
    {"id": "map", "name": "Map", "wrapper": "ply-code-map", "repos": [{"id": "map", "path": "ply-code-map"}, {"id": "twin", "path": "app/main"}]}
  ]
`, 1)
	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), broken)
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: f.root}
	calls := gitCalls(t)

	listed, err := ListProjects(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed.Projects) != 3 {
		t.Fatalf("listed %d projects; want 3", len(listed.Projects))
	}
	if got := calls(); len(got) != 0 {
		t.Errorf("list ran git %d times; want 0: %v", len(got), got)
	}

	shown, err := ShowProject(d, "ply")
	if err != nil {
		t.Fatal(err)
	}
	if len(shown.Repos) != 1 || shown.Repos[0].ID != "ply" || shown.Repos[0].GitCommonDir == "" {
		t.Fatalf("show ply: %+v", shown.Repos)
	}
	own := filepath.Join(f.root, "ply", "main")
	for _, call := range calls() {
		if !strings.Contains(call, own) {
			t.Errorf("show ply ran git outside %s: %s", own, call)
		}
	}

}

// A read of one member observes that member alone.
func TestProjectMemberReadObservesOnlyThatMember(t *testing.T) {
	f := newDefinitionFixture(t)
	two := strings.Replace(fixtureDefinitions, `"repos": [{"id": "ply", "path": "ply/main"}]`, `"repos": [{"id": "map", "path": "ply-code-map"}, {"id": "ply", "path": "ply/main"}]`, 1)
	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), two)
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: f.root}
	calls := gitCalls(t)
	shown, err := ShowProjectMember(d, "ply", "map")
	if err != nil {
		t.Fatal(err)
	}
	if len(shown.Project.RepoIDs) != 2 || len(shown.Repos) != 1 || shown.Repos[0].ID != "map" {
		t.Fatalf("member read: %+v %+v", shown.Project, shown.Repos)
	}
	other := filepath.Join(f.root, "ply", "main")
	for _, call := range calls() {
		if strings.Contains(call, other) {
			t.Errorf("member read ran git in %s: %s", other, call)
		}
	}
}
