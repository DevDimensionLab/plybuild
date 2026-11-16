package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every Project data category lives in <wrapper>/data/<category>; a
// definition with data_paths is rejected as an unknown field.
func TestProjectDataPathIsDataCategoryInWrapper(t *testing.T) {
	wrapper := physicalPath(t, t.TempDir())
	for _, category := range []string{DataPathBacklog, DataPathBacklogArchive, DataPathInbox, "plans"} {
		path, err := ProjectDataPath(ProjectRecord{ID: "sample", Wrapper: wrapper}, category)
		if err != nil || path != filepath.Join(wrapper, "data", category) {
			t.Fatalf("%s: path = %q, %v", category, path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(wrapper, "data")); !os.IsNotExist(err) {
		t.Fatalf("resolution created a directory: %v", err)
	}
}

func TestDefinitionWithDataPathsIsRejected(t *testing.T) {
	root, _, d := durableFixture(t)
	path := filepath.Join(root, DefinitionsFile)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(raw), `"wrapper":".",`, `"wrapper":".","data_paths":{"backlog":"data/backlog"},`, 1)
	if updated == string(raw) {
		t.Fatal("fixture definition not changed")
	}
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ShowProject(d, "sample"); err == nil || !strings.Contains(err.Error(), "data_paths") {
		t.Fatalf("show project = %v", err)
	}
}

func TestProjectDataPathAllowsSymlinkInsideWrapper(t *testing.T) {
	wrapper := physicalPath(t, t.TempDir())
	if err := os.Mkdir(filepath.Join(wrapper, "real"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(wrapper, "real"), filepath.Join(wrapper, "data")); err != nil {
		t.Fatal(err)
	}
	project := ProjectRecord{ID: "sample", Wrapper: wrapper}
	if path, err := ProjectDataPath(project, DataPathBacklog); err != nil || path != filepath.Join(wrapper, "data", "backlog") {
		t.Fatalf("path = %q, %v", path, err)
	}
}

func TestProjectDataPathRejectsInvalidPathsWithoutLeakingThem(t *testing.T) {
	wrapper := physicalPath(t, t.TempDir())
	outside := physicalPath(t, t.TempDir())
	if err := os.Symlink(outside, filepath.Join(wrapper, "escape")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(wrapper, "data"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(wrapper, "data", "backlog")); err != nil {
		t.Fatal(err)
	}
	defaultEscape := physicalPath(t, t.TempDir())
	if err := os.Symlink(outside, filepath.Join(defaultEscape, "data")); err != nil {
		t.Fatal(err)
	}
	for name, project := range map[string]ProjectRecord{
		"symlink-escape":      {Wrapper: wrapper},
		"symlink-data-escape": {Wrapper: defaultEscape},
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
