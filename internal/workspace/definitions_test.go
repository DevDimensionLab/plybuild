package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/repohistory"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

func TestProjectShowReturnsDefinedRemoteAndNullWhenAbsent(t *testing.T) {
	f := newDefinitionFixture(t)
	path := filepath.Join(f.root, workspacedef.File)
	definitions := strings.Replace(fixtureDefinitions, `"path": "ply/main"`, `"path": "ply/main", "remote": "git@github.com:example/ply.git"`, 1)
	mustWriteFile(t, path, definitions)
	dependencies := SystemDependencies()
	dependencies.Files = projectCwdFileSystem{FileSystem: dependencies.Files, cwd: f.root}
	for _, tc := range []struct {
		project ProjectID
		remote  any
	}{{"ply", "git@github.com:example/ply.git"}, {"app", nil}} {
		shown, err := ShowProject(dependencies, tc.project)
		if err != nil {
			t.Fatal(err)
		}
		encoded, err := MarshalProject(shown)
		if err != nil {
			t.Fatal(err)
		}
		var readback struct {
			Repositories []map[string]any `json:"repositories"`
		}
		if err := json.Unmarshal(encoded, &readback); err != nil {
			t.Fatal(err)
		}
		if len(readback.Repositories) != 1 || readback.Repositories[0]["remote"] != tc.remote {
			t.Fatalf("project %s remote: %s", tc.project, encoded)
		}
		if _, present := readback.Repositories[0]["remote"]; !present {
			t.Fatalf("project %s missing remote field: %s", tc.project, encoded)
		}
	}
}

func assertWorkErrorClass(t *testing.T, err error, class WorkItemErrorClass) {
	t.Helper()
	var typed *WorkItemError
	if !errors.As(err, &typed) || typed.Class != class {
		t.Fatalf("error = %v, want class %s", err, class)
	}
}

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

func writeWorkItems(t *testing.T, root, contents string, project ...ProjectID) {
	t.Helper()
	path := testWork(t, root, project...).Registry()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, path, contents)
}

func TestDefinitionHistoryProjectsRegistry(t *testing.T) {
	want := []repohistory.Binding{{RepoID: "ply", GitCommonDir: "/w/ply/main/.git"}, {RepoID: "app", GitCommonDir: "/w/app/main/.git"}, {RepoID: "ply", GitCommonDir: "/w/ply/main/.git"}, {RepoID: "app-server", GitCommonDir: "/w/app/server/.git"}}
	root := t.TempDir()
	writeWorkItems(t, root, envelope(RegistryFormat, historyRegistry))
	got, err := historicalBindings(testWork(t, root))
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("bindings %+v, error %v", got, err)
	}
}

func TestDefinitionHistoryAbsentAndEmptyCollections(t *testing.T) {
	got, err := historicalBindings(testWork(t, t.TempDir()))
	if got != nil || err != nil {
		t.Fatalf("absent: %+v %v", got, err)
	}
	root := t.TempDir()
	writeWorkItems(t, root, "format_version: 11\nregistry: {}\n")
	got, err = historicalBindings(testWork(t, root))
	if len(got) != 0 || err != nil {
		t.Fatalf("no collections: %+v %v", got, err)
	}
}

func TestDefinitionHistoryEmptyCommonDirRecordsNoIdentity(t *testing.T) {
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
	got, err := historicalBindings(testWork(t, root))
	if err != nil || len(got) != 4 {
		t.Fatalf("bindings %+v, error %v", got, err)
	}
	identities, err := repohistory.CommonDirs(got)
	if err != nil || len(identities) != 0 {
		t.Fatalf("identities %v, error %v", identities, err)
	}
}

func TestDefinitionHistoryRejectsUnsupportedRegistries(t *testing.T) {
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
			_, err := historicalBindings(testWork(t, root))
			assertWorkErrorClass(t, err, ErrorWorkStoreConflict)
		})
	}
}

func TestDefinitionHistoryRejectsNonRegularRegistry(t *testing.T) {
	directory := t.TempDir()
	if err := os.MkdirAll(testWork(t, directory).Registry(), 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := historicalBindings(testWork(t, directory))
	assertWorkErrorClass(t, err, ErrorWorkStoreConflict)

	target := t.TempDir()
	writeWorkItems(t, target, envelope(RegistryFormat, historyRegistry))
	link := t.TempDir()
	if err := os.MkdirAll(testWork(t, link).Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(testWork(t, target).Registry(), testWork(t, link).Registry()); err != nil {
		t.Fatal(err)
	}
	_, err = historicalBindings(testWork(t, link))
	assertWorkErrorClass(t, err, ErrorWorkStoreConflict)
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("error %v", err)
	}
}

type definitionFixture struct {
	root string
	ply  string // .ply directory
}

const fixtureDefinitions = `{
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

// newDefinitionFixture builds a workspace shaped like the real one: an
// envelope 10, inner 4 registry whose history names live and retired repos.
func newDefinitionFixture(t *testing.T) definitionFixture {
	t.Helper()
	root := physicalPath(t, t.TempDir())
	if _, err := Init(Dependencies{Files: projectCwdFileSystem{systemFileSystem{}, root}}); err != nil {
		t.Fatal(err)
	}
	for _, repository := range []string{"ply/main", "app/main", "ply-code-map", "retired"} {
		path := filepath.Join(root, filepath.FromSlash(repository))
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatal(err)
		}
		runGitForTest(t, path, "init", "-q")
	}
	mustWriteFile(t, filepath.Join(root, workspacedef.File), fixtureDefinitions)
	history := fmt.Sprintf(`format_version: 11
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
`, root)
	writeWorkItems(t, root, history, "ply")
	mustWriteFile(t, filepath.Join(root, MarkerDirectory, projectsLock), "")
	return definitionFixture{root: root, ply: filepath.Join(root, MarkerDirectory)}
}

func (f definitionFixture) dependencies(cwd string) Dependencies {
	return Dependencies{Files: projectCwdFileSystem{systemFileSystem{}, cwd}}
}

func (f definitionFixture) add(t *testing.T, cwd, project, repo, path string) (DefinitionRepoResult, error) {
	t.Helper()
	return AddDefinitionRepo(f.dependencies(cwd), project, repo, path, "")
}

type treeState map[string]string

func snapshotTree(t *testing.T, root string) treeState {
	t.Helper()
	state := treeState{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		state[path] = fmt.Sprintf("%s|%s|%x", info.Mode(), info.ModTime().Format(time.RFC3339Nano), raw)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return state
}

func (f definitionFixture) state(t *testing.T) treeState {
	t.Helper()
	state := snapshotTree(t, f.ply)
	definitions := snapshotTree(t, filepath.Join(f.root, workspacedef.File))
	for path, value := range definitions {
		state[path] = value
	}
	return state
}

func (f definitionFixture) requireUnchanged(t *testing.T, before treeState) {
	t.Helper()
	if after := f.state(t); !reflect.DeepEqual(before, after) {
		t.Fatalf("workspace files changed:\nbefore %v\nafter  %v", before, after)
	}
}

func (f definitionFixture) definitions(t *testing.T) workspacedef.Definitions {
	t.Helper()
	defs, _, err := workspacedef.Read(f.root)
	if err != nil {
		t.Fatal(err)
	}
	return *defs
}

func TestAddDefinitionRepoAddsRerunsAndConflicts(t *testing.T) {
	f := newDefinitionFixture(t)
	ply := snapshotTree(t, f.ply)
	before := f.definitions(t)

	result, err := f.add(t, f.root, "ply", "ply-code-map", "ply-code-map")
	want := DefinitionRepoResult{Workspace: f.root, ProjectID: "ply", RepoID: "ply-code-map", Path: "ply-code-map", Changed: true}
	if err != nil || result != want {
		t.Fatalf("result %+v, error %v", result, err)
	}
	if after := snapshotTree(t, f.ply); !reflect.DeepEqual(ply, after) {
		t.Fatalf(".ply changed: %v", after)
	}
	before.Projects[0].Repos = append(before.Projects[0].Repos, workspacedef.Repo{ID: "ply-code-map", Path: "ply-code-map"})
	if after := f.definitions(t); !reflect.DeepEqual(before, after) {
		t.Fatalf("definitions %+v", after)
	}

	state := f.state(t)
	want.Changed = false
	if result, err = f.add(t, filepath.Join(f.root, "ply"), "ply", "ply-code-map", "../ply-code-map"); err != nil || result != want {
		t.Fatalf("rerun %+v, error %v", result, err)
	}
	f.requireUnchanged(t, state)

	_, err = f.add(t, f.root, "ply", "ply-code-map", "app/main")
	assertWorkErrorClass(t, err, ErrorWorkIdentityConflict)
	f.requireUnchanged(t, state)

	_, err = f.add(t, f.root, "app", "duplicate", "ply/main")
	assertProjectErrorClass(t, err, ErrorProjectConflict)
	f.requireUnchanged(t, state)
}

func TestAddDefinitionRepoRejectsDifferentHistoricalCommonDir(t *testing.T) {
	f := newDefinitionFixture(t)
	state := f.state(t)
	_, err := f.add(t, f.root, "ply", "retired", "retired")
	assertWorkErrorClass(t, err, ErrorWorkIdentityConflict)
	if !strings.Contains(err.Error(), "/elsewhere/retired/.git") {
		t.Fatalf("error %v", err)
	}
	f.requireUnchanged(t, state)
}

func TestAddDefinitionRepoRejectsConflictingHistory(t *testing.T) {
	f := newDefinitionFixture(t)
	writeWorkItems(t, f.root, envelope(RegistryFormat, "tasks:\n  - repo_id: retired\n    git_common_dir: /a/.git\n  - repo_id: retired\n    git_common_dir: /b/.git\n"), "app")
	state := f.state(t)
	_, err := f.add(t, f.root, "app", "retired", "retired")
	assertWorkErrorClass(t, err, ErrorWorkIdentityConflict)
	f.requireUnchanged(t, state)
}

func TestAddDefinitionRepoValidatesGitRoot(t *testing.T) {
	f := newDefinitionFixture(t)
	if err := os.MkdirAll(filepath.Join(f.root, "ply", "main", "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(f.root, "plain"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(f.root, "file"), "file\n")
	state := f.state(t)
	for path, class := range map[string]ErrorClass{
		"ply/main/sub": ErrorProjectRepoInvalid,
		"plain":        ErrorProjectRepoInvalid,
		"file":         ErrorProjectPath,
		"missing":      ErrorProjectPath,
	} {
		_, err := f.add(t, f.root, "ply", "candidate", path)
		assertProjectErrorClass(t, err, class)
		f.requireUnchanged(t, state)
	}
}

func TestAddDefinitionRepoStoresCanonicalPathAndAcceptsNestedUnbornRepo(t *testing.T) {
	f := newDefinitionFixture(t)
	nested := filepath.Join(f.root, "ply", "main", "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, nested, "init", "-q")
	if err := os.Symlink(filepath.Join(f.root, "ply-code-map"), filepath.Join(f.root, "link")); err != nil {
		t.Fatal(err)
	}
	result, err := f.add(t, f.root, "ply", "ply-code-map", "link")
	if err != nil || result.Path != "ply-code-map" || !result.Changed {
		t.Fatalf("symlink %+v, error %v", result, err)
	}
	result, err = f.add(t, filepath.Join(f.root, "ply", "main"), "ply", "nested", "nested")
	if err != nil || result.Path != "ply/main/nested" || !result.Changed {
		t.Fatalf("nested %+v, error %v", result, err)
	}
	repos := f.definitions(t).Projects[0].Repos
	if want := []workspacedef.Repo{{ID: "ply", Path: "ply/main"}, {ID: "ply-code-map", Path: "ply-code-map"}, {ID: "nested", Path: "ply/main/nested"}}; !reflect.DeepEqual(repos, want) {
		t.Fatalf("repos %+v", repos)
	}
}

func TestAddDefinitionRepoRejectsInvalidRequests(t *testing.T) {
	f := newDefinitionFixture(t)
	outside := physicalPath(t, t.TempDir())
	runGitForTest(t, outside, "init", "-q")
	state := f.state(t)
	for name, tc := range map[string]struct {
		project, repo, path string
	}{
		"project ID":      {"Ply", "x", "ply-code-map"},
		"repo ID":         {"ply", "X", "ply-code-map"},
		"empty path":      {"ply", "x", ""},
		"root":            {"ply", "x", "."},
		"outside root":    {"ply", "x", outside},
		"unknown Project": {"other", "x", "ply-code-map"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := f.add(t, f.root, tc.project, tc.repo, tc.path)
			assertProjectErrorClass(t, err, ErrorProjectInvalidArguments)
			f.requireUnchanged(t, state)
		})
	}
}

func TestAddDefinitionRepoRejectsInvalidCurrentDefinitions(t *testing.T) {
	for name, tc := range map[string]struct {
		edit  func(string) string
		class ErrorClass
	}{
		"schema": {func(s string) string {
			return strings.Replace(s, `"format_version": 1,`, `"format_version": 1, "extra": 1,`, 1)
		}, ErrorProjectStoreConflict},
		"format":       {func(s string) string { return strings.Replace(s, `"format_version": 1`, `"format_version": 2`, 1) }, ErrorProjectStoreConflict},
		"invalid JSON": {func(s string) string { return s + "}" }, ErrorProjectStoreConflict},
		"duplicate ID": {func(s string) string { return strings.Replace(s, `"id": "app", "name"`, `"id": "ply", "name"`, 1) }, ErrorProjectInvalidArguments},
		"data path": {func(s string) string {
			return strings.Replace(s, `"wrapper": "ply",`, `"wrapper": "ply", "data_paths": {},`, 1)
		}, ErrorProjectStoreConflict},
		"missing wrapper": {func(s string) string { return strings.Replace(s, `"wrapper": "app"`, `"wrapper": "gone"`, 1) }, ErrorProjectPath},
	} {
		t.Run(name, func(t *testing.T) {
			f := newDefinitionFixture(t)
			mustWriteFile(t, filepath.Join(f.root, workspacedef.File), tc.edit(fixtureDefinitions))
			state := f.state(t)
			_, err := f.add(t, f.root, "ply", "ply-code-map", "ply-code-map")
			assertProjectErrorClass(t, err, tc.class)
			f.requireUnchanged(t, state)
		})
	}
	t.Run("missing definitions", func(t *testing.T) {
		f := newDefinitionFixture(t)
		if err := os.Remove(filepath.Join(f.root, workspacedef.File)); err != nil {
			t.Fatal(err)
		}
		_, err := f.add(t, f.root, "ply", "ply-code-map", "ply-code-map")
		assertProjectErrorClass(t, err, ErrorProjectStoreConflict)
		if _, statErr := os.Lstat(filepath.Join(f.root, workspacedef.File)); !errors.Is(statErr, fs.ErrNotExist) {
			t.Fatalf("definitions created: %v", statErr)
		}
	})
}

func TestAddDefinitionRepoWaitsForProjectsLock(t *testing.T) {
	f := newDefinitionFixture(t)
	lock, err := filelock.Acquire(filepath.Join(f.ply, projectsLock), filelock.Options{})
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := AddDefinitionRepo(f.dependencies(f.root), "ply", "ply-code-map", "ply-code-map", "")
		done <- err
	}()
	select {
	case err := <-done:
		t.Fatalf("add did not wait for projects.lock: %v", err)
	case <-time.After(200 * time.Millisecond):
	}
	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestDefinitionUpdateErrorMapping(t *testing.T) {
	path := "/w/" + workspacedef.File
	replaceErr := errors.New("replace failed")
	for name, tc := range map[string]struct {
		err      error
		project  ErrorClass
		work     WorkItemErrorClass
		contains string
	}{
		"not regular":           {&workspacedef.Error{Kind: workspacedef.KindNotRegular, Path: path, Detail: "definitions must be a regular file"}, ErrorProjectStoreConflict, "", "definitions must be a regular file"},
		"missing":               {&workspacedef.Error{Kind: workspacedef.KindMissing, Path: path, Detail: "definitions do not exist"}, ErrorProjectStoreConflict, "", "definitions do not exist"},
		"not found":             {&workspacedef.Error{Kind: workspacedef.KindNotFound, Detail: "Project x is not defined"}, ErrorProjectInvalidArguments, "", "Project x is not defined"},
		"conflict":              {&workspacedef.Error{Kind: workspacedef.KindConflict, Path: path, Detail: "definitions changed since they were read"}, "", ErrorWorkIdentityConflict, "repeat the request"},
		"publish":               {&workspacedef.Error{Kind: workspacedef.KindPublish, Path: path, Err: replaceErr}, "", ErrorWorkIO, "repeat the request"},
		"readback":              {&workspacedef.Error{Kind: workspacedef.KindReadback, Path: path, Err: replaceErr}, "", ErrorWorkReconciliation, "inspect it before repeating"},
		"os error":              {&fs.PathError{Op: "open", Path: path, Err: fs.ErrPermission}, ErrorProjectIO, "", "permission denied"},
		"lock":                  {&filelock.Error{Op: "not-regular", Path: "/w/.ply/projects.lock"}, ErrorProjectStoreConflict, "", "project lock is not a regular file"},
		"release":               {&filelock.ReleaseError{Path: "/w/.ply/projects.lock", Close: replaceErr}, ErrorProjectIO, "", "close lock"},
		"release after failure": {errors.Join(&workspacedef.Error{Kind: workspacedef.KindConflict, Path: path}, &filelock.ReleaseError{Path: "/w/.ply/projects.lock", Unlock: replaceErr}), ErrorProjectIO, "", "release after failure"},
	} {
		t.Run(name, func(t *testing.T) {
			err := definitionUpdateError(tc.err)
			if tc.project != "" {
				assertProjectErrorClass(t, err, tc.project)
			} else {
				assertWorkErrorClass(t, err, tc.work)
			}
			if !strings.Contains(err.Error(), tc.contains) {
				t.Fatalf("error %v", err)
			}
		})
	}
	typed := ProjectInvalidArguments("from check")
	if definitionUpdateError(typed) != typed {
		t.Fatal("typed check error was not returned as is")
	}
}

func TestCheckDefinitionsValidatesCurrentFileWithoutLockOrWrite(t *testing.T) {
	f := newDefinitionFixture(t)
	lock, err := filelock.Acquire(filepath.Join(f.ply, projectsLock), filelock.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = lock.Release() }()
	state := f.state(t)
	if err := CheckDefinitions(f.dependencies("/"), f.root); err != nil {
		t.Fatal(err)
	}
	f.requireUnchanged(t, state)

	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), strings.Replace(fixtureDefinitions, `"wrapper": "app"`, `"wrapper": "gone"`, 1))
	assertProjectErrorClass(t, CheckDefinitions(f.dependencies("/"), f.root), ErrorProjectPath)
	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), strings.Replace(fixtureDefinitions, `"wrapper": "ply",`, `"wrapper": "ply", "data_paths": {},`, 1))
	assertProjectErrorClass(t, CheckDefinitions(f.dependencies("/"), f.root), ErrorProjectStoreConflict)
	if err := os.Remove(filepath.Join(f.root, workspacedef.File)); err != nil {
		t.Fatal(err)
	}
	assertProjectErrorClass(t, CheckDefinitions(f.dependencies("/"), f.root), ErrorProjectStoreConflict)
	assertProjectErrorClass(t, CheckDefinitions(f.dependencies("/"), filepath.Join(f.root, "ply")), ErrorWorkspaceNotFound)
}

// A repository belongs to exactly one Project: the same repository ID under
// two Projects is rejected wherever definitions are resolved.
func TestDefinitionsRejectRepositoryInTwoProjects(t *testing.T) {
	f := newDefinitionFixture(t)
	shared := strings.Replace(fixtureDefinitions, `"repos": [{"id": "app", "path": "app/main"}]`, `"repos": [{"id": "app", "path": "app/main"}, {"id": "ply", "path": "ply/main"}]`, 1)
	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), shared)
	assertProjectErrorClass(t, CheckDefinitions(f.dependencies("/"), f.root), ErrorProjectDuplicateRepo)
	_, _, err := newSystemProjectStore().Snapshot(f.root)
	assertProjectErrorClass(t, err, ErrorProjectDuplicateRepo)
	_, err = metadataProjects(f.dependencies("/"), f.root)
	assertProjectErrorClass(t, err, ErrorProjectDuplicateRepo)
}
