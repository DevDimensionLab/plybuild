package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

func TestProjectMemberAddPreservesIdentityAndDirtyRepository(t *testing.T) {
	deps, root := discoverFixture(t)
	first, extra, owned := filepath.Join(root, "first"), filepath.Join(root, "extra"), filepath.Join(root, "owned")
	discoverRepo(t, first, "")
	discoverRepo(t, extra, "")
	discoverRepo(t, owned, "")
	if _, err := Init(deps); err != nil {
		t.Fatal(err)
	}
	original, err := AddProject(deps, ProjectAddInput{ProjectID: "shared", Name: "Shared", Wrapper: ".", Repos: []RepoInput{{RepoID: "first", Path: first}}})
	if err != nil {
		t.Fatal(err)
	}
	// Project shared's wrapper is the workspace root; wrappers may not overlap.
	other := filepath.Join(physicalPath(t, t.TempDir()), "other")
	if err := os.Mkdir(other, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err = AddProject(deps, ProjectAddInput{ProjectID: "other", Name: "Other", Wrapper: other, Repos: []RepoInput{{RepoID: "owned", Path: owned}}}); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(extra, "dirty"), "untracked source")
	alias := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(extra, alias); err != nil {
		t.Fatal(err)
	}
	result, err := AddProjectMember(deps, "shared", RepoInput{RepoID: "extra", Path: alias})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Created || result.Project.Name != original.Project.Name || result.Project.Wrapper != original.Project.Wrapper || !reflect.DeepEqual(result.Project.RepoIDs, []RepoID{"extra", "first"}) || result.Repos[0].Locator != extra {
		t.Fatalf("result = %#v", result)
	}
	path := filepath.Join(root, DefinitionsFile)
	before, _ := os.ReadFile(path)
	retry, err := AddProjectMember(deps, "shared", RepoInput{RepoID: "extra", Path: extra})
	if err != nil || retry.Created {
		t.Fatalf("retry = %#v, %v", retry, err)
	}
	dirty, err := os.ReadFile(filepath.Join(extra, "dirty"))
	if err != nil || string(dirty) != "untracked source" {
		t.Fatalf("dirty source changed: %q %v", dirty, err)
	}
	shown, err := ShowProject(deps, "other")
	if err != nil || !reflect.DeepEqual(shown.Project.RepoIDs, []RepoID{"owned"}) {
		t.Fatalf("other = %#v %v", shown, err)
	}
	for _, test := range []struct {
		id    ProjectID
		repo  RepoInput
		check func(error)
	}{
		{"shared", RepoInput{RepoID: "alias", Path: alias}, func(err error) { assertProjectErrorClass(t, err, ErrorProjectConflict) }},
		{"shared", RepoInput{RepoID: "extra", Path: first}, func(err error) { assertWorkErrorClass(t, err, ErrorWorkIdentityConflict) }},
		{"missing", RepoInput{RepoID: "extra", Path: extra}, func(err error) { assertProjectErrorClass(t, err, ErrorProjectInvalidArguments) }},
		// A repository belongs to exactly one Project.
		{"other", RepoInput{RepoID: "extra", Path: extra}, func(err error) { assertProjectErrorClass(t, err, ErrorProjectDuplicateRepo) }},
		{"shared", RepoInput{RepoID: "missing", Path: filepath.Join(root, "absent")}, func(err error) { assertProjectErrorClass(t, err, ErrorProjectPath) }},
	} {
		_, err := AddProjectMember(deps, test.id, test.repo)
		test.check(err)
		if after, _ := os.ReadFile(path); string(after) != string(before) {
			t.Fatal("retry or failure changed the definitions")
		}
	}
}

func TestProjectMemberConcurrentAppendAndIdentityConflict(t *testing.T) {
	deps, root := discoverFixture(t)
	for _, name := range []string{"first", "a", "b", "c"} {
		discoverRepo(t, filepath.Join(root, name), "")
	}
	if _, err := Init(deps); err != nil {
		t.Fatal(err)
	}
	if _, err := AddProject(deps, ProjectAddInput{ProjectID: "shared", Name: "Shared", Wrapper: ".", Repos: []RepoInput{{RepoID: "first", Path: "first"}}}); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for _, id := range []RepoID{"a", "b"} {
		wg.Add(1)
		go func(id RepoID) {
			defer wg.Done()
			_, err := AddProjectMember(deps, "shared", RepoInput{RepoID: id, Path: string(id)})
			errs <- err
		}(id)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	shown, err := ShowProject(deps, "shared")
	if err != nil || !reflect.DeepEqual(shown.Project.RepoIDs, []RepoID{"a", "b", "first"}) {
		t.Fatalf("show = %#v %v", shown, err)
	}
	// Competing IDs for the same new Git identity: exactly one can publish.
	errs = make(chan error, 2)
	for _, id := range []RepoID{"c", "d"} {
		wg.Add(1)
		go func(id RepoID) {
			defer wg.Done()
			_, err := AddProjectMember(deps, "shared", RepoInput{RepoID: id, Path: "c"})
			errs <- err
		}(id)
	}
	wg.Wait()
	close(errs)
	conflicts := 0
	for err := range errs {
		if err != nil {
			assertProjectErrorClass(t, err, ErrorProjectConflict)
			conflicts++
		}
	}
	if conflicts != 1 {
		t.Fatalf("conflicts = %d", conflicts)
	}
}

func TestProjectMemberMoveChangesIDPathAndRemote(t *testing.T) {
	deps, root := discoverFixture(t)
	first, old := filepath.Join(root, "first"), filepath.Join(root, "old")
	discoverRepo(t, first, "")
	discoverRepo(t, old, "git@github.com:o/old.git")
	if _, err := Init(deps); err != nil {
		t.Fatal(err)
	}
	if _, err := AddProject(deps, ProjectAddInput{ProjectID: "shared", Name: "Shared", Wrapper: ".", Repos: []RepoInput{{RepoID: "first", Path: first}, {RepoID: "old", Path: old}}}); err != nil {
		t.Fatal(err)
	}
	moved := filepath.Join(root, "group", "moved")
	if err := os.MkdirAll(filepath.Dir(moved), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(old, moved); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, moved, "remote", "set-url", "origin", "git@github.com:o/moved.git")
	result, err := MoveProjectMember(deps, "shared", "old", RepoInput{RepoID: "moved", Path: moved})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Created || !reflect.DeepEqual(result.Project.RepoIDs, []RepoID{"first", "moved"}) {
		t.Fatalf("result = %#v", result)
	}
	for _, repo := range result.Repos {
		if repo.ID == "moved" && (repo.Locator != moved || repo.Remote != "git@github.com:o/moved.git") {
			t.Fatalf("moved = %#v", repo)
		}
	}
	// The same ID with a new path is a move too.
	setWrapper := func(wrapper string) {
		t.Helper()
		if _, err := updateDefinitions(root, func(defs *workspacedef.Definitions) (bool, error) {
			defs.Projects[0].Repos[1].Wrapper = wrapper
			return true, nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	setWrapper("group")
	inside := filepath.Join(root, "group", "inside")
	if err := os.Rename(moved, inside); err != nil {
		t.Fatal(err)
	}
	if result, err := MoveProjectMember(deps, "shared", "moved", RepoInput{RepoID: "moved", Path: inside}); err != nil || !result.Created {
		t.Fatalf("path move = %#v, %v", result, err)
	}
	path := filepath.Join(root, DefinitionsFile)
	before, _ := os.ReadFile(path)
	if retry, err := MoveProjectMember(deps, "shared", "moved", RepoInput{RepoID: "moved", Path: inside}); err != nil || retry.Created {
		t.Fatalf("retry = %#v, %v", retry, err)
	}
	outside, plain := filepath.Join(root, "outside"), filepath.Join(root, "group", "plain")
	discoverRepo(t, outside, "")
	if err := os.Mkdir(plain, 0o755); err != nil {
		t.Fatal(err)
	}
	bound := deps
	bound.WorkItems = registrationsStore{deps.WorkItems, WorkItemRegistry{Tasks: []TaskRecord{{ID: "k", RepoID: "moved"}}}}
	for name, test := range map[string]struct {
		deps  Dependencies
		repo  RepoID
		to    RepoInput
		check func(error)
	}{
		"taken ID":           {deps, "moved", RepoInput{RepoID: "first", Path: inside}, func(err error) { assertWorkErrorClass(t, err, ErrorWorkIdentityConflict) }},
		"unknown repository": {deps, "absent", RepoInput{RepoID: "absent", Path: inside}, func(err error) { assertProjectErrorClass(t, err, ErrorProjectInvalidArguments) }},
		"not a repository":   {deps, "moved", RepoInput{RepoID: "moved", Path: plain}, func(err error) { assertProjectErrorClass(t, err, ErrorProjectRepoInvalid) }},
		"outside wrapper":    {deps, "moved", RepoInput{RepoID: "moved", Path: outside}, func(err error) { assertProjectErrorClass(t, err, ErrorProjectInvalidArguments) }},
		"bound old ID":       {bound, "moved", RepoInput{RepoID: "renamed", Path: inside}, func(err error) { assertProjectErrorClass(t, err, ErrorProjectConflict) }},
	} {
		_, err := MoveProjectMember(test.deps, "shared", test.repo, test.to)
		test.check(err)
		if after, _ := os.ReadFile(path); string(after) != string(before) {
			t.Fatalf("%s changed the definitions", name)
		}
	}
	// Bindings keep the old Git location, so they stop a move that keeps the ID too.
	if err := os.Rename(inside, moved); err != nil {
		t.Fatal(err)
	}
	_, err = MoveProjectMember(bound, "shared", "moved", RepoInput{RepoID: "moved", Path: moved})
	assertProjectErrorClass(t, err, ErrorProjectConflict)
	if after, _ := os.ReadFile(path); string(after) != string(before) {
		t.Fatal("bound move changed the definitions")
	}
}
