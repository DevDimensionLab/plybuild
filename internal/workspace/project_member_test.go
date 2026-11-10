package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

func TestProjectMemberAddPreservesIdentityAndDirtyRepository(t *testing.T) {
	root := createProjectWorkspace(t)
	deps := SystemDependencies()
	deps.Files = projectCwdFileSystem{FileSystem: deps.Files, cwd: root}
	wrapper := t.TempDir()
	first := physicalPath(t, t.TempDir())
	extra := physicalPath(t, t.TempDir())
	for _, repo := range []string{first, extra} {
		runGitForTest(t, repo, "init", "-q")
	}
	original, err := AddProject(deps, ProjectAddInput{ProjectID: "shared", Name: "Shared", Wrapper: wrapper, Repos: []RepoInput{{RepoID: "first", Path: first}}})
	if err != nil {
		t.Fatal(err)
	}
	// Reuse a repository registered by another Project, preserving its original locator.
	_, err = AddProject(deps, ProjectAddInput{ProjectID: "other", Name: "Other", Wrapper: t.TempDir(), Repos: []RepoInput{{RepoID: "extra", Path: extra}}})
	if err != nil {
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
	path := filepath.Join(root, MarkerDirectory, ProjectsFile)
	before, _ := os.ReadFile(path)
	retry, err := AddProjectMember(deps, "shared", RepoInput{RepoID: "extra", Path: extra})
	if err != nil || retry.Created {
		t.Fatalf("retry = %#v, %v", retry, err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(before) {
		t.Fatal("retry rewrote registry")
	}
	dirty, err := os.ReadFile(filepath.Join(extra, "dirty"))
	if err != nil || string(dirty) != "untracked source" {
		t.Fatalf("dirty source changed: %q %v", dirty, err)
	}
	other, err := ShowProject(deps, "other")
	if err != nil || !reflect.DeepEqual(other.Project.RepoIDs, []RepoID{"extra"}) {
		t.Fatalf("other = %#v %v", other, err)
	}
	for _, test := range []struct {
		id    ProjectID
		repo  RepoInput
		class ErrorClass
	}{
		{"shared", RepoInput{RepoID: "alias", Path: alias}, ErrorProjectConflict},
		{"shared", RepoInput{RepoID: "extra", Path: first}, ErrorProjectConflict},
		{"missing", RepoInput{RepoID: "extra", Path: extra}, ErrorProjectNotFound},
		{"shared", RepoInput{RepoID: "missing", Path: filepath.Join(root, "absent")}, ErrorProjectPath},
	} {
		_, err := AddProjectMember(deps, test.id, test.repo)
		assertProjectErrorClass(t, err, test.class)
		after, _ := os.ReadFile(path)
		if string(after) != string(before) {
			t.Fatal("failure changed registry")
		}
	}
}

func TestProjectMemberConcurrentAppendAndIdentityConflict(t *testing.T) {
	root := emptyStoreRoot(t)
	store := newSystemProjectStore()
	project, repos := storeFixture(root, "shared", "Shared", "first")
	if _, _, _, err := store.Add(root, project, repos); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for _, id := range []RepoID{"a", "b"} {
		wg.Add(1)
		go func(id RepoID) {
			defer wg.Done()
			_, _, _, err := store.AddMember(root, "shared", RepoRecord{ID: id, Locator: filepath.Join(root, string(id)), GitCommonDir: filepath.Join(root, "common-"+string(id))})
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
	projects, repositories, err := store.Snapshot(root)
	if err != nil || len(repositories) != 3 || !reflect.DeepEqual(projects[0].RepoIDs, []RepoID{"a", "b", "first"}) {
		t.Fatalf("snapshot = %#v %#v %v", projects, repositories, err)
	}
	// Competing IDs for the same new Git identity: exactly one can publish.
	errs = make(chan error, 2)
	for _, id := range []RepoID{"c", "d"} {
		wg.Add(1)
		go func(id RepoID) {
			defer wg.Done()
			_, _, _, err := store.AddMember(root, "shared", RepoRecord{ID: id, Locator: filepath.Join(root, string(id)), GitCommonDir: filepath.Join(root, "same-common")})
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
