package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
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
