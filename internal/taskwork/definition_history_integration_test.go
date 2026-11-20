package taskwork

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

const historyDefinitions = `{
  "format_version": 1,
  "projects": [
    {
      "id": "ply",
      "name": "Ply",
      "wrapper": "ply",
      "repos": [{"id": "ply", "path": "ply/main"}]
    },
    {"id": "app", "name": "App", "wrapper": "app", "repos": [{"id": "app", "path": "app/main"}]}
  ]
}
`

// newHistoryWorkspace is a workspace with Projects ply and app and the Git
// repositories ply/main, app/main, ply-code-map and retired. Project ply has
// no work-item registry yet.
func newHistoryWorkspace(t *testing.T) string {
	t.Helper()
	root := physicalPath(t, t.TempDir())
	if _, err := workspace.Init(workspace.Dependencies{Files: projectCwdFileSystem{workspace.SystemDependencies().Files, root}}); err != nil {
		t.Fatal(err)
	}
	for _, repository := range []string{"ply/main", "app/main", "ply-code-map", "retired"} {
		path := filepath.Join(root, filepath.FromSlash(repository))
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
		runGitForTest(t, path, "init", "-q")
	}
	mustWriteFile(t, filepath.Join(root, workspacedef.File), historyDefinitions)
	return root
}

// historyRegistries write registries of Project ply that record repository
// retired with Git common directory /elsewhere/retired/.git (or, for the
// published registry, root/repo.git), each with content that the strict
// registration read rejects.
var historyRegistries = map[string]func(*testing.T, string){
	"partial history": func(t *testing.T, root string) {
		writeWorkItems(t, root, fmt.Sprintf(`format_version: 11
registry:
  epics:
    - id: ply-master
      repo_bindings:
        - repo_id: ply
          git_common_dir: %[1]s/ply/main/.git
  tasks:
    - id: k5
      repo_id: ply-code-map
      git_common_dir: %[1]s/ply-code-map/.git
    - id: old
      repo_id: retired
      git_common_dir: /elsewhere/retired/.git
  worktree_operations: []
`, root), "ply")
	},
	"null records beside a repository binding": func(t *testing.T, root string) {
		writeWorkItems(t, root, envelope(RegistryFormat, "epics:\n  - null\n  - repo_bindings:\n      - null\n      - repo_id: retired\n        git_common_dir: /elsewhere/retired/.git\ntasks:\n  - null\n"), "ply")
	},
	"null records beside a Task": func(t *testing.T, root string) {
		writeWorkItems(t, root, envelope(RegistryFormat, "epics:\n  - null\n  - repo_bindings:\n      - null\ntasks:\n  - null\n  - repo_id: retired\n    git_common_dir: /elsewhere/retired/.git\n"), "ply")
	},
	"missing evidence object": publishDamaged(func(s string) string {
		return strings.Replace(s, "integration_results: []", "integration_results:\n    - sha256:"+strings.Repeat("0", 64), 1)
	}),
	"unknown field":     publishDamaged(func(s string) string { return strings.Replace(s, "  tasks: []", "  tasks: []\n  surprise: 1", 1) }),
	"invalid Epic base": publishDamaged(func(s string) string { return strings.Replace(s, strings.Repeat("a", 40), "nothex", 1) }),
}

// publishDamaged publishes a valid registry of Project ply whose Epic binds
// retired, then damages it so that only the strict read rejects it.
func publishDamaged(damage func(string) string) func(*testing.T, string) {
	return func(t *testing.T, root string) {
		t.Helper()
		work := testWork(t, root, "ply")
		registry := minimalEpicRegistry(root)
		registry.Epics[0].RepoBindings[0].RepoID = "retired"
		registry.EpicBaseVersions = []EpicBaseVersion{adoptedEpicBase(registry.Epics[0], registry.Epics[0].RepoBindings[0])}
		if err := newSystemWorkItemStore().WithLock(work, func(s WorkItemStoreSession) error { return s.Publish(registry) }); err != nil {
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
	}
}

// AddDefinitionRepo with the work-item store as RepoUses admits a new
// repository ID and refuses a historical ID with another Git common
// directory without touching the definitions, also when the registry holds
// null records or content the strict registration read rejects.
func TestAddDefinitionRepoReadsHistoryThroughTheWorkItemStore(t *testing.T) {
	t.Parallel()
	for name, write := range historyRegistries {
		t.Run(name, func(t *testing.T) {
			root := newHistoryWorkspace(t)
			write(t, root)
			d := workspace.Dependencies{Files: projectCwdFileSystem{workspace.SystemDependencies().Files, root}, RepoUses: newSystemWorkItemStore()}

			result, err := workspace.AddDefinitionRepo(d, "ply", "ply-code-map", "ply-code-map", "")
			if err != nil || !result.Changed {
				t.Fatalf("new repository: %+v, %v", result, err)
			}
			definitions := filepath.Join(root, workspacedef.File)
			before, err := os.ReadFile(definitions)
			if err != nil {
				t.Fatal(err)
			}

			_, err = workspace.AddDefinitionRepo(d, "ply", "retired", "retired", "")

			assertWorkErrorClass(t, err, workspace.ErrorWorkIdentityConflict)
			if !strings.Contains(err.Error(), "retired") {
				t.Fatalf("error %v", err)
			}
			if after, readErr := os.ReadFile(definitions); readErr != nil || string(after) != string(before) {
				t.Fatalf("definitions changed on refusal (%v):\n%s", readErr, after)
			}
		})
	}
}

// The store's RepoUses is the strict registration read: it rejects every
// registry whose repository identities RecordedRepos still reads.
func TestWorkItemStoreRepoUsesRejectsWhatRecordedReposReads(t *testing.T) {
	t.Parallel()
	for name, write := range historyRegistries {
		t.Run(name, func(t *testing.T) {
			root := newHistoryWorkspace(t)
			write(t, root)
			work := testWork(t, root, "ply")
			uses := SystemDependencies().RepoUses

			recorded, err := uses.RecordedRepos(work)
			found := false
			for _, use := range recorded {
				found = found || use.RepoID == "retired" && use.GitCommonDir != ""
			}
			if err != nil || !found {
				t.Fatalf("recorded %+v, %v", recorded, err)
			}

			_, err = uses.RepoUses(work)

			assertWorkErrorClass(t, err, workspace.ErrorWorkStoreConflict)
		})
	}
}
