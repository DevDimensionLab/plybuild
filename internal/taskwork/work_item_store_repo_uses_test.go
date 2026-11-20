package taskwork

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/repohistory"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// historyRegistry has the logical shape of the real registry: bindings carry
// worktrees and other fields that the projection ignores.
const historyRegistry = `epics:
  - id: ply-master
    title: Ply master
    project_id: ply
    repo_bindings:
      - repo_id: ply
        git_common_dir: /w/ply/main/.git
        worktree:
          id: wt_1
          locator: /w/ply/main
          ref: refs/heads/master
      - repo_id: app
        git_common_dir: /w/app/main/.git
tasks:
  - id: k5
    epic_id: ply-master
    repo_id: ply
    git_common_dir: /w/ply/main/.git
    worktree:
      locator: /w/ply/k5
  - id: old
    repo_id: app-server
    git_common_dir: /w/app/server/.git
worktree_operations: []
`

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

// recordedBindings reads the repository identities of work like
// AddDefinitionRepo does.
func recordedBindings(work workspace.ProjectWork) ([]repohistory.Binding, error) {
	recorded, err := newSystemWorkItemStore().RecordedRepos(work)
	var bindings []repohistory.Binding
	for _, use := range recorded {
		bindings = append(bindings, repohistory.Binding{RepoID: string(use.RepoID), GitCommonDir: use.GitCommonDir})
	}
	return bindings, err
}

func TestDefinitionHistoryProjectsRegistry(t *testing.T) {
	t.Parallel()
	want := []repohistory.Binding{{RepoID: "ply", GitCommonDir: "/w/ply/main/.git"}, {RepoID: "app", GitCommonDir: "/w/app/main/.git"}, {RepoID: "ply", GitCommonDir: "/w/ply/main/.git"}, {RepoID: "app-server", GitCommonDir: "/w/app/server/.git"}}
	root := t.TempDir()
	writeWorkItems(t, root, envelope(RegistryFormat, historyRegistry))
	got, err := recordedBindings(testWork(t, root))
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("bindings %+v, error %v", got, err)
	}
}

func TestDefinitionHistoryAbsentAndEmptyCollections(t *testing.T) {
	t.Parallel()
	got, err := recordedBindings(testWork(t, t.TempDir()))
	if got != nil || err != nil {
		t.Fatalf("absent: %+v %v", got, err)
	}
	root := t.TempDir()
	writeWorkItems(t, root, "format_version: 11\nregistry: {}\n")
	got, err = recordedBindings(testWork(t, root))
	if len(got) != 0 || err != nil {
		t.Fatalf("no collections: %+v %v", got, err)
	}
}

func TestDefinitionHistoryEmptyCommonDirRecordsNoIdentity(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeWorkItems(t, root, envelope(RegistryFormat, `epics:
  - repo_bindings:
      - repo_id: Not-Valid
        git_common_dir: ""
      - repo_id: 7
tasks:
  - repo_id: "bad id"
    git_common_dir: ""
  - {}
`))
	got, err := recordedBindings(testWork(t, root))
	if err != nil || len(got) != 0 {
		t.Fatalf("bindings %+v, error %v", got, err)
	}
}

// A null in place of an Epic, a repository binding or a Task records no identity.
func TestDefinitionHistoryNullRecordsNoIdentity(t *testing.T) {
	t.Parallel()
	for name, body := range map[string]string{
		"epic":               "epics:\n  - null\n",
		"repository binding": "epics:\n  - repo_bindings:\n      - null\n",
		"task":               "tasks:\n  - null\n",
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			writeWorkItems(t, root, envelope(RegistryFormat, body))
			got, err := recordedBindings(testWork(t, root))
			if err != nil || len(got) != 0 {
				t.Fatalf("bindings %+v, error %v", got, err)
			}
		})
	}
}

func TestDefinitionHistoryRejectsUnsupportedRegistries(t *testing.T) {
	t.Parallel()
	single := func(body string) string { return envelope(RegistryFormat, body) }
	for name, contents := range map[string]string{
		"empty":                   "",
		"two documents":           single(historyRegistry) + "---\nformat_version: 11\n",
		"sequence root":           "- format_version: 11\n",
		"inline format 4":         "format_version: 4\n" + historyRegistry,
		"inline format 8":         "format_version: 8\n" + historyRegistry,
		"previous envelope 10":    envelope(10, "format_version: 4\n"+historyRegistry),
		"future format":           envelope(12, historyRegistry),
		"format 0":                envelope(0, historyRegistry),
		"string format":           "format_version: \"11\"\nregistry: {}\n",
		"missing format":          "registry: {}\n",
		"extra envelope key":      single(historyRegistry) + "integration_results: []\n",
		"envelope without inner":  "format_version: 11\n",
		"envelope registry list":  "format_version: 11\nregistry: []\n",
		"alias":                   single("epics: &e []\ntasks: *e\n"),
		"merge key":               single("base: &b {repo_id: ply}\ntasks:\n  - <<: *b\n    git_common_dir: /w/.git\n"),
		"duplicate repo_id":       single("tasks:\n  - repo_id: ply\n    repo_id: app\n    git_common_dir: /w/.git\n"),
		"epics not a sequence":    single("epics: {}\n"),
		"task not a mapping":      single("tasks:\n  - ply\n"),
		"bindings not a sequence": single("epics:\n  - repo_bindings: {}\n"),
		"relative common dir":     single("tasks:\n  - repo_id: ply\n    git_common_dir: ply/.git\n"),
		"unclean common dir":      single("tasks:\n  - repo_id: ply\n    git_common_dir: /w/ply/../ply/.git\n"),
		"null common dir":         single("tasks:\n  - repo_id: ply\n    git_common_dir: null\n"),
		"invalid repo id":         single("epics:\n  - repo_bindings:\n      - repo_id: Ply\n        git_common_dir: /w/.git\n"),
		"missing repo id":         single("tasks:\n  - git_common_dir: /w/.git\n"),
		"integer repo id":         single("tasks:\n  - repo_id: 7\n    git_common_dir: /w/.git\n"),
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			writeWorkItems(t, root, contents)
			_, err := recordedBindings(testWork(t, root))
			assertWorkErrorClass(t, err, workspace.ErrorWorkStoreConflict)
		})
	}
}

func TestDefinitionHistoryRejectsNonRegularRegistry(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	if err := os.MkdirAll(testWork(t, directory).Registry(), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := recordedBindings(testWork(t, directory))
	assertWorkErrorClass(t, err, workspace.ErrorWorkStoreConflict)

	target := t.TempDir()
	writeWorkItems(t, target, envelope(RegistryFormat, historyRegistry))
	link := t.TempDir()
	if err := os.MkdirAll(testWork(t, link).Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(testWork(t, target).Registry(), testWork(t, link).Registry()); err != nil {
		t.Fatal(err)
	}
	_, err = recordedBindings(testWork(t, link))
	assertWorkErrorClass(t, err, workspace.ErrorWorkStoreConflict)
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("error %v", err)
	}
}

// Null records beside a binding do not hide the identity that binding records.
func TestDefinitionHistoryKeepsIdentityBesideNullRecords(t *testing.T) {
	t.Parallel()
	for name, body := range map[string]string{
		"repository binding": "epics:\n  - null\n  - repo_bindings:\n      - null\n      - repo_id: retired\n        git_common_dir: /elsewhere/retired/.git\ntasks:\n  - null\n",
		"task":               "epics:\n  - null\n  - repo_bindings:\n      - null\ntasks:\n  - null\n  - repo_id: retired\n    git_common_dir: /elsewhere/retired/.git\n",
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			writeWorkItems(t, root, envelope(RegistryFormat, body))
			got, err := recordedBindings(testWork(t, root))
			want := []repohistory.Binding{{RepoID: "retired", GitCommonDir: "/elsewhere/retired/.git"}}
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("bindings %+v, error %v", got, err)
			}
		})
	}
}

// RecordedRepos reads only the repository identities of the registry. A
// registry that the store rejects for other content still tells which Git
// common directory it records for a repository ID.
func TestDefinitionHistoryReadsRegistryContentTheStoreRejects(t *testing.T) {
	t.Parallel()
	for name, damage := range map[string]func(string) string{
		"missing evidence object": func(s string) string {
			return strings.Replace(s, "integration_results: []", "integration_results:\n    - sha256:"+strings.Repeat("0", 64), 1)
		},
		"unknown field":     func(s string) string { return strings.Replace(s, "  tasks: []", "  tasks: []\n  surprise: 1", 1) },
		"invalid Epic base": func(s string) string { return strings.Replace(s, strings.Repeat("a", 40), "nothex", 1) },
		"incomplete records": func(string) string {
			return envelope(RegistryFormat, "tasks:\n  - repo_id: retired\n    git_common_dir: /elsewhere/retired/.git\n")
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := createProjectWorkspace(t)
			work := testWork(t, root)
			registry := minimalEpicRegistry(root)
			registry.Epics[0].RepoBindings[0].RepoID = "retired"
			registry.EpicBaseVersions = []EpicBaseVersion{adoptedEpicBase(registry.Epics[0], registry.Epics[0].RepoBindings[0])}
			store := newSystemWorkItemStore()
			if err := store.WithLock(work, func(s WorkItemStoreSession) error { return s.Publish(registry) }); err != nil {
				t.Fatal(err)
			}
			valid, err := os.ReadFile(work.Registry())
			if err != nil {
				t.Fatal(err)
			}
			damaged := damage(string(valid))
			if damaged == string(valid) {
				t.Fatalf("registry was not damaged:\n%s", valid)
			}
			mustWriteFile(t, work.Registry(), damaged)
			_, err = store.SnapshotRegistrations(work)
			assertWorkErrorClass(t, err, workspace.ErrorWorkStoreConflict)

			got, err := recordedBindings(work)
			if err != nil || len(got) != 1 || got[0].RepoID != "retired" || got[0].GitCommonDir == "" {
				t.Fatalf("bindings %+v, error %v", got, err)
			}
		})
	}
}
