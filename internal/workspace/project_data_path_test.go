package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeDataPathsDefinition(t *testing.T, root, dataPaths string) {
	t.Helper()
	path := filepath.Join(root, DefinitionsFile)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(raw), `"wrapper":".",`, `"wrapper":".","data_paths":`+dataPaths+`,`, 1)
	if updated == string(raw) {
		t.Fatal("fixture definition not changed")
	}
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDefinitionDataPathsReachProjectProjections(t *testing.T) {
	root, _, d := durableFixture(t)
	writeDataPathsDefinition(t, root, `{"backlog":"workflow/backlog"}`)
	want := map[string]string{"backlog": "workflow/backlog"}

	shown, err := ShowProject(d, "sample")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(shown.Project.DataPaths, want) {
		t.Fatalf("show data paths = %#v", shown.Project.DataPaths)
	}
	listed, err := ListProjects(d)
	if err != nil || len(listed.Projects) != 1 || !reflect.DeepEqual(listed.Projects[0].DataPaths, want) {
		t.Fatalf("list projects = %#v, %v", listed.Projects, err)
	}
	if _, err := ListEpics(d); err != nil {
		t.Fatal(err)
	}
	registered, _, err := registrationProjects(d, root)
	if err != nil || len(registered) != 1 || !reflect.DeepEqual(registered[0].DataPaths, want) {
		t.Fatalf("registration projects = %#v, %v", registered, err)
	}
	if _, err := os.Stat(filepath.Join(root, MarkerDirectory, ProjectsFile)); !os.IsNotExist(err) {
		t.Fatalf("data paths duplicated in a second registry: %v", err)
	}
}

func TestDefinitionRejectsInvalidDataPaths(t *testing.T) {
	for name, dataPaths := range map[string]string{
		"empty":            `{"backlog":""}`,
		"absolute":         `{"backlog":"/tmp/backlog"}`,
		"parent":           `{"backlog":"../backlog"}`,
		"nested-parent":    `{"backlog":"workflow/../../backlog"}`,
		"invalid-category": `{"Back Log":"backlog"}`,
	} {
		t.Run(name, func(t *testing.T) {
			root, _, d := durableFixture(t)
			writeDataPathsDefinition(t, root, dataPaths)
			_, err := ShowProject(d, "sample")
			var projectErr *ProjectError
			if !errors.As(err, &projectErr) || projectErr.Class != ErrorProjectInvalidArguments {
				t.Fatalf("show project = %v", err)
			}
		})
	}
}

func TestProjectDataPathDefaultsToWrapperCategory(t *testing.T) {
	wrapper := physicalPath(t, t.TempDir())
	for name, dataPaths := range map[string]map[string]string{
		"missing-data-paths": nil,
		"empty-data-paths":   {},
		"missing-backlog":    {"other": "elsewhere"},
	} {
		t.Run(name, func(t *testing.T) {
			path, err := ProjectDataPath(ProjectRecord{ID: "sample", Wrapper: wrapper, DataPaths: dataPaths}, DataPathBacklog)
			if err != nil || path != filepath.Join(wrapper, "backlog") {
				t.Fatalf("path = %q, %v", path, err)
			}
		})
	}
}

func TestProjectDataPathUsesConfiguredRelativePath(t *testing.T) {
	wrapper := physicalPath(t, t.TempDir())
	project := ProjectRecord{ID: "sample", Wrapper: wrapper, DataPaths: map[string]string{DataPathBacklog: "workflow/backlog"}}
	path, err := ProjectDataPath(project, DataPathBacklog)
	if err != nil || path != filepath.Join(wrapper, "workflow", "backlog") {
		t.Fatalf("path = %q, %v", path, err)
	}
	if _, err := os.Stat(filepath.Join(wrapper, "workflow")); !os.IsNotExist(err) {
		t.Fatalf("resolution created a directory: %v", err)
	}
}

func TestProjectDataPathAllowsSymlinkInsideWrapper(t *testing.T) {
	wrapper := physicalPath(t, t.TempDir())
	if err := os.Mkdir(filepath.Join(wrapper, "real"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(wrapper, "real"), filepath.Join(wrapper, "alias")); err != nil {
		t.Fatal(err)
	}
	project := ProjectRecord{ID: "sample", Wrapper: wrapper, DataPaths: map[string]string{DataPathBacklog: "alias/backlog"}}
	if path, err := ProjectDataPath(project, DataPathBacklog); err != nil || path != filepath.Join(wrapper, "alias", "backlog") {
		t.Fatalf("path = %q, %v", path, err)
	}
}

func TestProjectDataPathRejectsInvalidPathsWithoutLeakingThem(t *testing.T) {
	wrapper := physicalPath(t, t.TempDir())
	outside := physicalPath(t, t.TempDir())
	if err := os.Symlink(outside, filepath.Join(wrapper, "escape")); err != nil {
		t.Fatal(err)
	}
	defaultEscape := physicalPath(t, t.TempDir())
	if err := os.Symlink(outside, filepath.Join(defaultEscape, "backlog")); err != nil {
		t.Fatal(err)
	}
	for name, project := range map[string]ProjectRecord{
		"empty":           {Wrapper: wrapper, DataPaths: map[string]string{DataPathBacklog: ""}},
		"absolute":        {Wrapper: wrapper, DataPaths: map[string]string{DataPathBacklog: outside}},
		"parent":          {Wrapper: wrapper, DataPaths: map[string]string{DataPathBacklog: "../backlog"}},
		"nested-parent":   {Wrapper: wrapper, DataPaths: map[string]string{DataPathBacklog: "a/../../backlog"}},
		"symlink-escape":  {Wrapper: wrapper, DataPaths: map[string]string{DataPathBacklog: "escape/backlog"}},
		"symlink-default": {Wrapper: defaultEscape},
	} {
		t.Run(name, func(t *testing.T) {
			project.ID = "sample"
			path, err := ProjectDataPath(project, DataPathBacklog)
			var projectErr *ProjectError
			if !errors.As(err, &projectErr) || projectErr.Class != ErrorProjectPath || path != "" {
				t.Fatalf("path = %q, %v", path, err)
			}
			if strings.Contains(err.Error(), outside) || strings.Contains(err.Error(), wrapper) || strings.Contains(err.Error(), defaultEscape) {
				t.Fatalf("error leaks a path: %v", err)
			}
		})
	}
	if entries, err := os.ReadDir(outside); err != nil || len(entries) != 0 {
		t.Fatalf("outside touched: %v, %v", entries, err)
	}
}

func TestProjectDataPathResolvesBacklogArchiveSeparately(t *testing.T) {
	wrapper := t.TempDir()
	dataPaths := map[string]string{DataPathBacklogArchive: "data/archive/backlog"}
	if err := validateDataPaths(dataPaths); err != nil {
		t.Fatal(err)
	}
	for configured, want := range map[bool]string{false: filepath.Join(wrapper, "backlog-archive"), true: filepath.Join(wrapper, "data", "archive", "backlog")} {
		project := ProjectRecord{ID: "sample", Wrapper: wrapper}
		if configured {
			project.DataPaths = dataPaths
		}
		if path, err := ProjectDataPath(project, DataPathBacklogArchive); err != nil || path != want {
			t.Fatalf("configured=%v: path = %q, %v", configured, path, err)
		}
	}
}
