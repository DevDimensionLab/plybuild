package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestProjectStorePublishesDeterministicRegistry(t *testing.T) {
	root := emptyStoreRoot(t)
	store := newSystemProjectStore()
	project, repositories := storeFixture(root, "ply", "Ply", "ply")
	storedProject, storedRepositories, created, err := store.Add(root, project, repositories)
	if err != nil {
		t.Fatal(err)
	}
	if !created || !reflect.DeepEqual(storedProject, project) || !reflect.DeepEqual(storedRepositories, repositories) {
		t.Fatalf("result = %#v, %#v, created %t", storedProject, storedRepositories, created)
	}
	want := "format_version: 1\nprojects:\n  - id: ply\n    name: Ply\n    wrapper: " + filepath.Join(root, "wrapper-ply") + "\n    repo_ids:\n      - ply\nrepos:\n  - id: ply\n    locator: " + filepath.Join(root, "repo-ply") + "\n    git_common_dir: " + filepath.Join(root, "common-ply") + "\n"
	registryPath := filepath.Join(root, MarkerDirectory, ProjectsFile)
	got, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("registry = %q, want %q", got, want)
	}
	if runtime.GOOS != "windows" {
		assertMode(t, registryPath, 0o644)
		assertMode(t, filepath.Join(root, MarkerDirectory, projectsLock), 0o600)
	}
}

func TestProjectStoreIdempotentAddAndSnapshotDoNotRewrite(t *testing.T) {
	root := emptyStoreRoot(t)
	store := newSystemProjectStore()
	project, repositories := storeFixture(root, "ply", "Ply", "ply")
	if _, _, _, err := store.Add(root, project, repositories); err != nil {
		t.Fatal(err)
	}
	registryPath := filepath.Join(root, MarkerDirectory, ProjectsFile)
	before, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	oldTime := time.Unix(946684800, 0)
	if err := os.Chtimes(registryPath, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	repositories[0].Locator = filepath.Join(root, "alternate-locator")
	storedProject, storedRepositories, created, err := store.Add(root, project, repositories)
	if err != nil {
		t.Fatal(err)
	}
	if created || storedProject.ID != "ply" || storedRepositories[0].Locator == repositories[0].Locator {
		t.Fatalf("idempotent result = %#v, %#v, created %t", storedProject, storedRepositories, created)
	}
	if _, _, err := store.Snapshot(root); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, after) || info.ModTime().Unix() != oldTime.Unix() {
		t.Fatalf("idempotent access changed registry: mtime %v, bytes equal %t", info.ModTime(), reflect.DeepEqual(before, after))
	}
}

func TestProjectStoreRejectsStrictRegistryViolations(t *testing.T) {
	root := emptyStoreRoot(t)
	absolute := filepath.Join(root, "absolute")
	valid := fmt.Sprintf("format_version: 1\nprojects:\n  - id: ply\n    name: Ply\n    wrapper: %s\n    repo_ids:\n      - ply\nrepos:\n  - id: ply\n    locator: %s\n    git_common_dir: %s\n", absolute+"-wrapper", absolute+"-repo", absolute+"-common")
	tests := map[string]string{
		"empty":                "",
		"unknown field":        strings.Replace(valid, "format_version: 1", "format_version: 1\nunknown: true", 1),
		"duplicate field":      strings.Replace(valid, "format_version: 1", "format_version: 1\nformat_version: 1", 1),
		"multi document":       valid + "---\n" + valid,
		"wrong root":           "- value\n",
		"wrong type":           strings.Replace(valid, "projects:\n", "projects: wrong\n", 1),
		"wrong version":        strings.Replace(valid, "format_version: 1", "format_version: 2", 1),
		"invalid project ID":   strings.Replace(valid, "id: ply", "id: Ply", 1),
		"invalid name":         strings.Replace(valid, "name: Ply", "name: ' Ply'", 1),
		"relative wrapper":     strings.Replace(valid, absolute+"-wrapper", "relative", 1),
		"empty membership":     strings.Replace(valid, "repo_ids:\n      - ply", "repo_ids: []", 1),
		"missing repository":   strings.Replace(valid, "- ply\nrepos:", "- missing\nrepos:", 1),
		"unused repository":    strings.Replace(valid, "repos:\n", fmt.Sprintf("repos:\n  - id: extra\n    locator: %s\n    git_common_dir: %s\n", absolute+"-extra", absolute+"-extra-common"), 1),
		"unsorted project IDs": fmt.Sprintf("format_version: 1\nprojects:\n  - id: zed\n    name: Zed\n    wrapper: %s\n    repo_ids: [zed]\n  - id: alpha\n    name: Alpha\n    wrapper: %s\n    repo_ids: [alpha]\nrepos:\n  - id: alpha\n    locator: %s\n    git_common_dir: %s\n  - id: zed\n    locator: %s\n    git_common_dir: %s\n", absolute+"-z-wrapper", absolute+"-a-wrapper", absolute+"-a", absolute+"-a-common", absolute+"-z", absolute+"-z-common"),
		"duplicate common dir": fmt.Sprintf("format_version: 1\nprojects:\n  - id: ply\n    name: Ply\n    wrapper: %s\n    repo_ids: [one, two]\nrepos:\n  - id: one\n    locator: %s\n    git_common_dir: %s\n  - id: two\n    locator: %s\n    git_common_dir: %s\n", absolute+"-wrapper", absolute+"-one", absolute+"-common", absolute+"-two", absolute+"-common"),
	}
	for name, contents := range tests {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(root, MarkerDirectory, ProjectsFile)
			if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
				t.Fatal(err)
			}
			_, _, err := newSystemProjectStore().Snapshot(root)
			assertProjectErrorClass(t, err, ErrorProjectStoreConflict)
			got, readErr := os.ReadFile(path)
			if readErr != nil || string(got) != contents {
				t.Fatalf("invalid bytes changed: %q, %v", got, readErr)
			}
		})
	}
}

func TestProjectStoreRejectsRegistryAndLockSymlinks(t *testing.T) {
	for _, name := range []string{ProjectsFile, projectsLock} {
		t.Run(name, func(t *testing.T) {
			root := emptyStoreRoot(t)
			target := filepath.Join(root, "target")
			mustWriteFile(t, target, "target")
			if err := os.Symlink(target, filepath.Join(root, MarkerDirectory, name)); err != nil {
				t.Fatal(err)
			}
			project, repositories := storeFixture(root, "ply", "Ply", "ply")
			_, _, _, err := newSystemProjectStore().Add(root, project, repositories)
			assertProjectErrorClass(t, err, ErrorProjectStoreConflict)
			contents, readErr := os.ReadFile(target)
			if readErr != nil || string(contents) != "target" {
				t.Fatalf("symlink target changed: %q, %v", contents, readErr)
			}
		})
	}
}

func TestProjectStoreConflictAndRepositoryReuseRules(t *testing.T) {
	root := emptyStoreRoot(t)
	store := newSystemProjectStore()
	first, firstRepos := storeFixture(root, "one", "One", "shared")
	if _, _, _, err := store.Add(root, first, firstRepos); err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name    string
		project ProjectRecord
		repos   []RepoRecord
	}{
		{name: "same ID different name", project: ProjectRecord{ID: "one", Name: "Changed", Wrapper: first.Wrapper, RepoIDs: []RepoID{"shared"}}, repos: firstRepos},
		{name: "same wrapper other ID", project: ProjectRecord{ID: "two", Name: "Two", Wrapper: first.Wrapper, RepoIDs: []RepoID{"two"}}, repos: []RepoRecord{{ID: "two", Locator: filepath.Join(root, "repo-two"), GitCommonDir: filepath.Join(root, "common-two")}}},
		{name: "repository ID rebound", project: ProjectRecord{ID: "two", Name: "Two", Wrapper: filepath.Join(root, "wrapper-two"), RepoIDs: []RepoID{"shared"}}, repos: []RepoRecord{{ID: "shared", Locator: filepath.Join(root, "other"), GitCommonDir: filepath.Join(root, "other-common")}}},
		{name: "common directory alias", project: ProjectRecord{ID: "two", Name: "Two", Wrapper: filepath.Join(root, "wrapper-two"), RepoIDs: []RepoID{"alias"}}, repos: []RepoRecord{{ID: "alias", Locator: filepath.Join(root, "alias"), GitCommonDir: firstRepos[0].GitCommonDir}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, _, _, err := store.Add(root, test.project, test.repos)
			assertProjectErrorClass(t, err, ErrorProjectConflict)
		})
	}

	second := ProjectRecord{ID: "two", Name: "Two", Wrapper: filepath.Join(root, "wrapper-two"), RepoIDs: []RepoID{"shared"}}
	stored, repositories, created, err := store.Add(root, second, []RepoRecord{{ID: "shared", Locator: filepath.Join(root, "alternate"), GitCommonDir: firstRepos[0].GitCommonDir}})
	if err != nil || !created || stored.ID != "two" || repositories[0].Locator != firstRepos[0].Locator {
		t.Fatalf("repository reuse = %#v, %#v, %t, %v", stored, repositories, created, err)
	}
}

func TestProjectStoreFaultsPreserveOldOrCompleteNewRegistryAndRetry(t *testing.T) {
	tests := []struct {
		name          string
		operation     string
		shortWrite    bool
		failCleanup   bool
		expectPublish bool
	}{
		{name: "lock open", operation: "lock-open"},
		{name: "lock", operation: "lock"},
		{name: "temp open", operation: "temp-open"},
		{name: "write", operation: "write"},
		{name: "short write", shortWrite: true},
		{name: "chmod", operation: "chmod"},
		{name: "file sync", operation: "file-sync"},
		{name: "close", operation: "close"},
		{name: "replace", operation: "replace"},
		{name: "cleanup", operation: "write", failCleanup: true},
		{name: "directory sync", operation: "directory-sync", expectPublish: true},
		{name: "unlock", operation: "unlock", expectPublish: true},
		{name: "lock close", operation: "lock-close", expectPublish: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := emptyStoreRoot(t)
			healthy := newSystemProjectStore()
			first, firstRepos := storeFixture(root, "one", "One", "one")
			if _, _, _, err := healthy.Add(root, first, firstRepos); err != nil {
				t.Fatal(err)
			}
			registryPath := filepath.Join(root, MarkerDirectory, ProjectsFile)
			oldBytes, err := os.ReadFile(registryPath)
			if err != nil {
				t.Fatal(err)
			}
			second, secondRepos := storeFixture(root, "two", "Two", "two")
			faulty := &systemProjectStore{faults: &storeFaults{shortWrite: test.shortWrite}}
			faulty.faults.fail = func(operation string) error {
				if operation == test.operation || (test.failCleanup && operation == "cleanup") {
					return errInjected
				}
				return nil
			}
			_, _, _, err = faulty.Add(root, second, secondRepos)
			assertProjectErrorClass(t, err, ErrorProjectIO)
			afterFailure, readErr := os.ReadFile(registryPath)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if test.expectPublish == reflect.DeepEqual(afterFailure, oldBytes) {
				t.Fatalf("published = %t, bytes unchanged = %t", test.expectPublish, reflect.DeepEqual(afterFailure, oldBytes))
			}
			_, _, created, retryErr := healthy.Add(root, second, secondRepos)
			if retryErr != nil || created == test.expectPublish {
				t.Fatalf("retry created = %t, error = %v; post-replace failure = %t", created, retryErr, test.expectPublish)
			}
		})
	}
}

func TestProjectStoreHandleCloseReleasesAdvisoryLock(t *testing.T) {
	root := emptyStoreRoot(t)
	store := &systemProjectStore{}
	lock, err := store.acquire(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := lock.Close(); err != nil {
		t.Fatal(err)
	}

	project, repositories := storeFixture(root, "one", "One", "one")
	done := make(chan error, 1)
	go func() {
		_, _, _, err := store.Add(root, project, repositories)
		done <- err
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("project add remained blocked after the lock handle closed")
	}
}

func TestProjectStoreConcurrentAddsSerializeWithoutLostUpdates(t *testing.T) {
	t.Run("different", func(t *testing.T) {
		root := emptyStoreRoot(t)
		store := newSystemProjectStore()
		first, firstRepos := storeFixture(root, "one", "One", "one")
		second, secondRepos := storeFixture(root, "two", "Two", "two")
		results := concurrentAdds(t, store, root, []storeAdd{{first, firstRepos}, {second, secondRepos}})
		for _, result := range results {
			if result.err != nil || !result.created {
				t.Fatalf("different add = %#v", result)
			}
		}
		projects, _, err := store.Snapshot(root)
		if err != nil || len(projects) != 2 {
			t.Fatalf("snapshot projects = %#v, %v", projects, err)
		}
	})

	t.Run("identical", func(t *testing.T) {
		root := emptyStoreRoot(t)
		store := newSystemProjectStore()
		project, repositories := storeFixture(root, "one", "One", "one")
		results := concurrentAdds(t, store, root, []storeAdd{{project, repositories}, {project, repositories}})
		created := 0
		for _, result := range results {
			if result.err != nil {
				t.Fatal(result.err)
			}
			if result.created {
				created++
			}
		}
		if created != 1 {
			t.Fatalf("created count = %d", created)
		}
	})

	t.Run("conflicting", func(t *testing.T) {
		root := emptyStoreRoot(t)
		store := newSystemProjectStore()
		first, repositories := storeFixture(root, "one", "One", "one")
		second := first
		second.Name = "Other"
		results := concurrentAdds(t, store, root, []storeAdd{{first, repositories}, {second, repositories}})
		created, conflicts := 0, 0
		for _, result := range results {
			if result.created {
				created++
			}
			var typed *ProjectError
			if errors.As(result.err, &typed) && typed.Class == ErrorProjectConflict {
				conflicts++
			}
		}
		if created != 1 || conflicts != 1 {
			t.Fatalf("results = %#v", results)
		}
	})
}

type storeAdd struct {
	project ProjectRecord
	repos   []RepoRecord
}

type storeAddResult struct {
	created bool
	err     error
}

func concurrentAdds(t *testing.T, store ProjectStore, root string, additions []storeAdd) []storeAddResult {
	t.Helper()
	start := make(chan struct{})
	results := make(chan storeAddResult, len(additions))
	var wait sync.WaitGroup
	for _, addition := range additions {
		addition := addition
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			_, _, created, err := store.Add(root, addition.project, addition.repos)
			results <- storeAddResult{created: created, err: err}
		}()
	}
	close(start)
	done := make(chan struct{})
	go func() { wait.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("concurrent store calls timed out")
	}
	close(results)
	output := make([]storeAddResult, 0, len(additions))
	for result := range results {
		output = append(output, result)
	}
	return output
}

func emptyStoreRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, MarkerDirectory), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func storeFixture(root, projectID, name, repoID string) (ProjectRecord, []RepoRecord) {
	project := ProjectRecord{ID: ProjectID(projectID), Name: name, Wrapper: filepath.Join(root, "wrapper-"+projectID), RepoIDs: []RepoID{RepoID(repoID)}}
	repositories := []RepoRecord{{ID: RepoID(repoID), Locator: filepath.Join(root, "repo-"+repoID), GitCommonDir: filepath.Join(root, "common-"+repoID)}}
	return project, repositories
}
