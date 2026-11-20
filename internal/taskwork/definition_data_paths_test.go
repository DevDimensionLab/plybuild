package taskwork

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A definition with data_paths is rejected as an unknown field.
func TestDefinitionWithDataPathsIsRejected(t *testing.T) {
	t.Parallel()
	root, _, d := durableFixture(t)
	path := filepath.Join(root, workspace.DefinitionsFile)
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
	if _, err := workspace.ShowProject(d.Dependencies, "sample"); err == nil || !strings.Contains(err.Error(), "data_paths") {
		t.Fatalf("show project = %v", err)
	}
}
