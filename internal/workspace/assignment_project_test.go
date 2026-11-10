package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

func cwdDependencies(cwd string) Dependencies {
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{d.Files, cwd}
	return d
}

func requireProjectDetail(t *testing.T, err error, class ErrorClass, detail string) {
	t.Helper()
	assertProjectErrorClass(t, err, class)
	if got := err.(*ProjectError).Detail; got != detail {
		t.Fatalf("detail %q, want %q", got, detail)
	}
}

func TestResolveDefinedProjectFromDefinitionsWithoutLegacyRegistry(t *testing.T) {
	f := newDefinitionFixture(t)
	requireAbsent(t, filepath.Join(f.ply, ProjectsFile))
	project, err := ResolveDefinedProject(cwdDependencies(f.root), "app")
	want := ProjectRecord{ID: "app", Name: "App", Wrapper: filepath.Join(f.root, "app"), RepoIDs: []RepoID{"app"}}
	if err != nil || !reflect.DeepEqual(project, want) {
		t.Fatalf("project %+v %v", project, err)
	}
	_, err = ResolveDefinedProject(cwdDependencies(f.root), "nope")
	requireProjectDetail(t, err, ErrorProjectNotFound, "project nope is not registered")
	_, err = ResolveDefinedProject(cwdDependencies(f.root), "Not Valid")
	assertProjectErrorClass(t, err, ErrorProjectInvalidArguments)
	requireAbsent(t, filepath.Join(f.ply, ProjectsFile))
}

func TestResolveDefinedProjectFromLegacyRegistry(t *testing.T) {
	root := physicalPath(t, t.TempDir())
	d := cwdDependencies(root)
	if _, err := Init(d); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(root, "legacy", "main")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, repo, "init", "-q")
	input := ProjectAddInput{ProjectID: "legacy", Name: "Legacy", Wrapper: "legacy", Repos: []RepoInput{{RepoID: "core", Path: "legacy/main"}}}
	if _, err := AddProject(d, input); err != nil {
		t.Fatal(err)
	}
	requireAbsent(t, filepath.Join(root, workspacedef.File))
	project, err := ResolveDefinedProject(d, "legacy")
	if err != nil || project.Wrapper != filepath.Join(root, "legacy") || !reflect.DeepEqual(project.RepoIDs, []RepoID{"core"}) {
		t.Fatalf("project %+v %v", project, err)
	}
	selected, err := ResolveDefinedProject(cwdDependencies(repo), "")
	if err != nil || selected.ID != "legacy" {
		t.Fatalf("selected %+v %v", selected, err)
	}
}

func TestResolveDefinedProjectSelectsTheProjectContainingTheWorkingDirectory(t *testing.T) {
	f := newDefinitionFixture(t)
	for cwd, want := range map[string]ProjectID{"app": "app", "app/main": "app", "ply/main": "ply"} {
		project, err := ResolveDefinedProject(cwdDependencies(filepath.Join(f.root, cwd)), "")
		if err != nil || project.ID != want {
			t.Fatalf("cwd %s: %+v %v", cwd, project, err)
		}
	}
	link := filepath.Join(f.root, "link-to-app")
	if err := os.Symlink(filepath.Join(f.root, "app", "main"), link); err != nil {
		t.Fatal(err)
	}
	if project, err := ResolveDefinedProject(cwdDependencies(link), ""); err != nil || project.ID != "app" {
		t.Fatalf("via symlink: %+v %v", project, err)
	}
	_, err := ResolveDefinedProject(cwdDependencies(f.root), "")
	requireProjectDetail(t, err, ErrorProjectInvalidArguments, "project is required: 0 projects contain the current directory")
}

func TestResolveDefinedProjectRefusesNestedWrappersAndSkipsMissingOnes(t *testing.T) {
	f := newDefinitionFixture(t)
	definitions := `{"format_version": 1, "projects": [
  {"id": "outer", "name": "Outer", "wrapper": ".", "repos": [{"id": "ply", "path": "ply/main"}], "epics": []},
  {"id": "app", "name": "App", "wrapper": "app", "repos": [{"id": "app", "path": "app/main"}], "epics": []},
  {"id": "gone", "name": "Gone", "wrapper": "gone", "repos": [{"id": "gone", "path": "gone/main"}], "epics": []}
]}
`
	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), definitions)
	_, err := ResolveDefinedProject(cwdDependencies(filepath.Join(f.root, "app", "main")), "")
	requireProjectDetail(t, err, ErrorProjectInvalidArguments, "project is required: 2 projects contain the current directory")
	if strings.Contains(err.Error(), f.root) {
		t.Fatalf("error carries a path: %v", err)
	}
	project, err := ResolveDefinedProject(cwdDependencies(filepath.Join(f.root, "ply")), "")
	if err != nil || project.ID != "outer" {
		t.Fatalf("outer only: %+v %v", project, err)
	}
}

func TestProjectDataPathOfAssignments(t *testing.T) {
	f := newDefinitionFixture(t)
	project, err := ResolveDefinedProject(cwdDependencies(f.root), "app")
	if err != nil {
		t.Fatal(err)
	}
	if path, err := ProjectDataPath(project, "assignments"); err != nil || path != filepath.Join(f.root, "app", "assignments") {
		t.Fatalf("default %s %v", path, err)
	}
	project.DataPaths = map[string]string{"assignments": "data/assignments"}
	if path, err := ProjectDataPath(project, "assignments"); err != nil || path != filepath.Join(f.root, "app", "data", "assignments") {
		t.Fatalf("configured %s %v", path, err)
	}
	requireAbsent(t, filepath.Join(f.root, "app", "data"))
}
