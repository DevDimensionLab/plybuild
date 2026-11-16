package projectdata

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func physical(t *testing.T) string {
	t.Helper()
	path, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// Every category of Project data lives in <wrapper>/data/<category>, and
// resolving a path creates nothing.
func TestEveryCategoryIsADirectoryUnderWrapperData(t *testing.T) {
	wrapper := physical(t)
	area := For("sample", wrapper)
	for category, resolve := range map[string]func() (string, error){
		"agent-status":    area.AgentStatus,
		"assignments":     area.Assignments,
		"backlog":         area.Backlog,
		"backlog-archive": area.BacklogArchive,
		"inbox":           area.Inbox,
		"intelligence":    area.Intelligence,
		"plans":           area.Plans,
		"publications":    area.Publications,
		"work":            area.Work,
	} {
		path, err := resolve()
		if err != nil || path != filepath.Join(wrapper, "data", category) {
			t.Fatalf("%s: path = %q, %v", category, path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(wrapper, "data")); !os.IsNotExist(err) {
		t.Fatalf("resolution created a directory: %v", err)
	}
}

func TestSymlinkInsideWrapperIsAllowed(t *testing.T) {
	wrapper := physical(t)
	if err := os.Mkdir(filepath.Join(wrapper, "real"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(wrapper, "real"), filepath.Join(wrapper, "data")); err != nil {
		t.Fatal(err)
	}
	if path, err := For("sample", wrapper).Backlog(); err != nil || path != filepath.Join(wrapper, "data", "backlog") {
		t.Fatalf("path = %q, %v", path, err)
	}
}

func TestPathLeavingWrapperIsRejectedWithoutLeakingPaths(t *testing.T) {
	outside := physical(t)
	categoryEscape := physical(t)
	if err := os.Mkdir(filepath.Join(categoryEscape, "data"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(categoryEscape, "data", "backlog")); err != nil {
		t.Fatal(err)
	}
	dataEscape := physical(t)
	if err := os.Symlink(outside, filepath.Join(dataEscape, "data")); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(physical(t), "missing")
	for name, wrapper := range map[string]string{"category": categoryEscape, "data": dataEscape, "missing-wrapper": missing} {
		t.Run(name, func(t *testing.T) {
			path, err := For("sample", wrapper).Backlog()
			var dataErr *Error
			if !errors.As(err, &dataErr) || dataErr.Class != ClassPath || path != "" {
				t.Fatalf("path = %q, %v", path, err)
			}
			if !strings.HasPrefix(err.Error(), "workspace_project_path_error: ") || !strings.Contains(err.Error(), "sample") {
				t.Fatalf("error = %v", err)
			}
			for _, leaked := range []string{outside, wrapper} {
				if strings.Contains(err.Error(), leaked) {
					t.Fatalf("error leaks a path: %v", err)
				}
			}
		})
	}
	if entries, err := os.ReadDir(outside); err != nil || len(entries) != 0 {
		t.Fatalf("outside touched: %v, %v", entries, err)
	}
}
