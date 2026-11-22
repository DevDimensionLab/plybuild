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
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

func TestProjectShowReturnsDefinedRemoteAndNullWhenAbsent(t *testing.T) {
	t.Parallel()
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

type definitionFixture struct {
	root    string
	ply     string // .ply directory
	history historyUses
}

// historyUses serves, per Project, the registrations of a work-item
// registry: Epic bindings and Tasks with their repository identities.
type historyUses map[ProjectID][]RepoUse

func (h historyUses) RepoUses(work ProjectWork) ([]RepoUse, error) {
	return h[work.Project], nil
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

// newDefinitionFixture builds a workspace shaped like the real one, whose
// Project ply has registrations of live and retired repos.
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
	mustWriteFile(t, filepath.Join(root, MarkerDirectory, projectsLock), "")
	history := historyUses{"ply": {
		{Live: true, Kind: RepoUseEpic, ID: "ply-master", RepoID: "ply", GitCommonDir: root + "/ply/main/.git"},
		{Live: true, Kind: RepoUseTask, ID: "k5", RepoID: "ply-code-map", GitCommonDir: root + "/ply-code-map/.git"},
		{Live: true, Kind: RepoUseTask, ID: "old", RepoID: "retired", GitCommonDir: "/elsewhere/retired/.git"},
	}}
	return definitionFixture{root: root, ply: filepath.Join(root, MarkerDirectory), history: history}
}

func (f definitionFixture) dependencies(cwd string) Dependencies {
	return Dependencies{Files: projectCwdFileSystem{systemFileSystem{}, cwd}, RepoUses: f.history}
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
	t.Parallel()
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
	t.Parallel()
	f := newDefinitionFixture(t)
	state := f.state(t)
	_, err := f.add(t, f.root, "ply", "retired", "retired")
	assertWorkErrorClass(t, err, ErrorWorkIdentityConflict)
	if !strings.Contains(err.Error(), "/elsewhere/retired/.git") {
		t.Fatalf("error %v", err)
	}
	f.requireUnchanged(t, state)
}

// Finished work keeps its old Git location as history: it neither conflicts
// nor binds the repository ID to that location.
func TestAddDefinitionRepoIgnoresFinishedHistory(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	f.history["ply"] = []RepoUse{{Kind: RepoUseTask, ID: "old", RepoID: "retired", GitCommonDir: "/elsewhere/retired/.git"}, {Kind: RepoUseTask, ID: "older", RepoID: "retired", GitCommonDir: "/other/retired/.git"}}
	if _, err := f.add(t, f.root, "ply", "retired", "retired"); err != nil {
		t.Fatal(err)
	}
}

func TestAddDefinitionRepoRejectsConflictingHistory(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	f.history["app"] = []RepoUse{{Live: true, Kind: RepoUseTask, RepoID: "retired", GitCommonDir: "/a/.git"}, {Live: true, Kind: RepoUseTask, RepoID: "retired", GitCommonDir: "/b/.git"}}
	state := f.state(t)
	_, err := f.add(t, f.root, "app", "retired", "retired")
	assertWorkErrorClass(t, err, ErrorWorkIdentityConflict)
	f.requireUnchanged(t, state)
}

func TestAddDefinitionRepoValidatesGitRoot(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	f := newDefinitionFixture(t)
	shared := strings.Replace(fixtureDefinitions, `"repos": [{"id": "app", "path": "app/main"}]`, `"repos": [{"id": "app", "path": "app/main"}, {"id": "ply", "path": "ply/main"}]`, 1)
	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), shared)
	assertProjectErrorClass(t, CheckDefinitions(f.dependencies("/"), f.root), ErrorProjectDuplicateRepo)
	_, _, err := newSystemProjectStore().Snapshot(f.root)
	assertProjectErrorClass(t, err, ErrorProjectDuplicateRepo)
	_, err = metadataProjects(f.dependencies("/"), f.root)
	assertProjectErrorClass(t, err, ErrorProjectDuplicateRepo)
}
