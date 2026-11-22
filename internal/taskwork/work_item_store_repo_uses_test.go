package taskwork

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func envelope(version int, body string) string {
	return fmt.Sprintf("format_version: %d\nregistry:\n  %s", version, strings.ReplaceAll(strings.TrimSuffix(body, "\n"), "\n", "\n  ")+"\n")
}

func writeWorkItems(t *testing.T, root, contents string, project ...workspace.ProjectID) {
	t.Helper()
	path := testWork(t, root, project...).Registry()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, path, contents)
}

// The registration read rejects YAML aliases and merge keys, which would let one
// node stand for several records.
func TestWorkItemStoreRepoUsesRejectsAliasesAndMergeKeys(t *testing.T) {
	t.Parallel()
	store := newSystemWorkItemStore()
	valid, err := encodeWorkItemRegistry(emptyWorkItemRegistry())
	if err != nil {
		t.Fatal(err)
	}
	control := t.TempDir()
	writeWorkItems(t, control, string(physicalRegistry(valid)))
	if _, err := store.RepoUses(testWork(t, control)); err != nil {
		t.Fatalf("valid registry: %v", err)
	}
	for _, tc := range []struct {
		name, logical, reason string
	}{
		{"alias", strings.Replace(strings.Replace(string(valid), "epics: []", "epics: &e []", 1), "tasks: []", "tasks: *e", 1), "aliases are unsupported"},
		{"merge key", string(valid) + "<<: {}\n", "duplicate or non-string request key"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			writeWorkItems(t, root, string(physicalRegistry([]byte(tc.logical))))
			_, err := store.RepoUses(testWork(t, root))
			assertWorkErrorClass(t, err, workspace.ErrorWorkStoreConflict)
			if d := RegistryError(err); d == nil || d.Code != "registry_invalid_content" || d.Detail != tc.reason || d.Artifact != testWork(t, root).Registry() {
				t.Fatalf("diagnostic %+v (%v)", d, err)
			}
		})
	}
}

// The registration read rejects a registry that is a directory or a symlink.
func TestWorkItemStoreRepoUsesRejectsNonRegularRegistry(t *testing.T) {
	t.Parallel()
	store := newSystemWorkItemStore()
	directory := t.TempDir()
	if err := os.MkdirAll(testWork(t, directory).Registry(), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := store.RepoUses(testWork(t, directory))
	assertWorkErrorClass(t, err, workspace.ErrorWorkStoreConflict)

	target := t.TempDir()
	writeWorkItems(t, target, envelope(RegistryFormat, "epics: []\ntasks: []\nworktree_operations: []\n"))
	link := t.TempDir()
	if err := os.MkdirAll(testWork(t, link).Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(testWork(t, target).Registry(), testWork(t, link).Registry()); err != nil {
		t.Fatal(err)
	}
	_, err = store.RepoUses(testWork(t, link))
	assertWorkErrorClass(t, err, workspace.ErrorWorkStoreConflict)
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("error %v", err)
	}
}
