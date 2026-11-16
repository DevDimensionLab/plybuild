package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A Project's data lives in its wrapper, so no wrapper may lie inside
// another: registering one is refused, and definitions edited to contain one
// do not resolve. Project sample's wrapper is the workspace root.
func TestNestedWrappersAreRefused(t *testing.T) {
	root, _, d := durableFixture(t)
	nested := filepath.Join(root, "nested")
	repo := filepath.Join(nested, "main")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	runLocalGit(t, repo, "init", "-b", "main")
	input := ProjectAddInput{ProjectID: "nested", Name: "Nested", Wrapper: nested, Repos: []RepoInput{{RepoID: "nested", Path: repo}}}
	_, err := AddProject(d, input)
	assertNestedWrapperConflict(t, "add", err)

	path := filepath.Join(root, DefinitionsFile)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(raw), `]}]}`, `]},{"id":"nested","name":"Nested","wrapper":"nested","repos":[{"id":"nested","path":"nested/main"}]}]}`, 1)
	if updated == string(raw) {
		t.Fatalf("fixture definition not changed: %s", raw)
	}
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = DefinedWork(root)
	assertNestedWrapperConflict(t, "defined work", err)
	_, err = ListProjects(d)
	assertNestedWrapperConflict(t, "list", err)
}

func assertNestedWrapperConflict(t *testing.T, operation string, err error) {
	t.Helper()
	var projectErr *ProjectError
	if !errors.As(err, &projectErr) || projectErr.Class != ErrorProjectConflict || !strings.Contains(err.Error(), "overlaps the wrapper of project") {
		t.Fatalf("%s: %v", operation, err)
	}
}
