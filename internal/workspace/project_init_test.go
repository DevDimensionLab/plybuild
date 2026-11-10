package workspace

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func initFixture(t *testing.T) (Dependencies, string, string) {
	t.Helper()
	root := createProjectWorkspace(t)
	wrapper := filepath.Join(root, "example")
	if err := os.Mkdir(wrapper, 0700); err != nil {
		t.Fatal(err)
	}
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{FileSystem: d.Files, cwd: wrapper}
	return d, root, wrapper
}
func initRepo(t *testing.T, path string) {
	t.Helper()
	runGitForTest(t, "", "init", "-b", "main", path)
	runGitForTest(t, path, "-c", "user.name=Fixture", "-c", "user.email=fixture@example.invalid", "commit", "--allow-empty", "-m", "fixture")
}
func TestProjectInitLinkedDiscoveryAndPreservedRetry(t *testing.T) {
	d, root, w := initFixture(t)
	main := filepath.Join(w, "main")
	initRepo(t, main)
	runGitForTest(t, main, "worktree", "add", "-b", "feature", filepath.Join(w, "aaa-feature"))
	mustWriteFile(t, filepath.Join(main, "dirty.txt"), "keep dirty\n")
	before := projectSnapshotFiles(t, root)
	p, err := InitProject(d, ProjectInitInput{ProjectID: "example"})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Repositories) != 1 || p.Repositories[0].Locator != main || p.Repositories[0].ID != "example" {
		t.Fatalf("plan %#v", p)
	}
	if !reflect.DeepEqual(before, projectSnapshotFiles(t, root)) {
		t.Fatal("preview wrote")
	}
	// Omitted diagnostics are not registration intent: a new note is benign.
	mustWriteFile(t, filepath.Join(w, "note.txt"), "unrelated wrapper note\n")
	benign, err := InitProject(d, ProjectInitInput{ProjectID: "example"})
	if err != nil || benign.Confirmation != p.Confirmation || len(benign.Omissions) != len(p.Omissions)+1 {
		t.Fatalf("benign omission changed intent: %#v %v", benign, err)
	}
	in := ProjectInitInput{ProjectID: "example", Apply: true, Confirmation: p.Confirmation}
	applied, err := InitProject(d, in)
	if err != nil || !applied.Applied {
		t.Fatalf("apply %#v %v", applied, err)
	}
	registered := projectSnapshotFiles(t, root)
	if _, err = InitProject(d, in); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(registered, projectSnapshotFiles(t, root)) {
		t.Fatal("retry changed state")
	}
	initRepo(t, filepath.Join(w, "extra"))
	changed := projectSnapshotFiles(t, root)
	_, err = InitProject(d, in)
	if err == nil || !strings.Contains(err.Error(), "project add") {
		t.Fatalf("changed membership %v", err)
	}
	if !reflect.DeepEqual(changed, projectSnapshotFiles(t, root)) {
		t.Fatal("conflict wrote")
	}
}
func TestProjectInitDistinctOmissionsAndDrift(t *testing.T) {
	d, root, w := initFixture(t)
	initRepo(t, filepath.Join(w, "api"))
	initRepo(t, filepath.Join(w, "web"))
	initRepo(t, filepath.Join(w, "nested", "main"))
	if err := os.Symlink(filepath.Join(w, "web"), filepath.Join(w, "alias")); err != nil {
		t.Fatal(err)
	}
	p, err := InitProject(d, ProjectInitInput{ProjectID: "example"})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Repositories) != 2 || p.Repositories[0].ID != "api" || p.Repositories[1].ID != "web" || len(p.Omissions) != 2 {
		t.Fatalf("plan %#v", p)
	}
	initRepo(t, filepath.Join(w, "third"))
	before := projectSnapshotFiles(t, root)
	_, err = InitProject(d, ProjectInitInput{ProjectID: "example", Apply: true, Confirmation: p.Confirmation})
	if err == nil || !strings.Contains(err.Error(), "confirmation") {
		t.Fatalf("drift %v", err)
	}
	if !reflect.DeepEqual(before, projectSnapshotFiles(t, root)) {
		t.Fatal("drift wrote")
	}
}
func TestProjectInitCompanionAndBoundaries(t *testing.T) {
	t.Run("registered companion", func(t *testing.T) {
		d, root, w := initFixture(t)
		main := filepath.Join(w, "main")
		initRepo(t, main)
		_, err := AddProject(d, ProjectAddInput{ProjectID: "example", Name: "Existing name", Wrapper: w, Repos: []RepoInput{{RepoID: "example", Path: main}}})
		if err != nil {
			t.Fatal(err)
		}
		docs := filepath.Join(w, "docs")
		initRepo(t, docs)
		_, err = UpdateProjectMetadata(d, ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "example", Role: CompanionDocs, Locator: docs, ActorClaim: "fixture"})
		if err != nil {
			t.Fatal(err)
		}
		before := projectSnapshotFiles(t, root)
		p, err := InitProject(d, ProjectInitInput{ProjectID: "example"})
		if err != nil {
			t.Fatal(err)
		}
		if p.Name != "Existing name" || len(p.Repositories) != 1 || len(p.Omissions) != 1 {
			t.Fatalf("plan %#v", p)
		}
		_, err = InitProject(d, ProjectInitInput{ProjectID: "example", Apply: true, Confirmation: p.Confirmation})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(before, projectSnapshotFiles(t, root)) {
			t.Fatal("companion retry changed metadata")
		}
	})
	for _, kind := range []string{"empty", "broken", "wrapper-root", "invalid-id"} {
		t.Run(kind, func(t *testing.T) {
			d, root, w := initFixture(t)
			switch kind {
			case "empty":
				initRepo(t, filepath.Join(w, "nested", "main"))
			case "broken":
				initRepo(t, filepath.Join(w, "main"))
				os.Mkdir(filepath.Join(w, "broken"), 0700)
				mustWriteFile(t, filepath.Join(w, "broken", ".git"), "gitdir: /missing/al34\n")
			case "wrapper-root":
				initRepo(t, w)
				initRepo(t, filepath.Join(w, "main"))
			case "invalid-id":
				initRepo(t, filepath.Join(w, "Bad_ID"))
				initRepo(t, filepath.Join(w, "web"))
			}
			before := projectSnapshotFiles(t, root)
			p, err := InitProject(d, ProjectInitInput{ProjectID: "example"})
			if err == nil || len(p.Blockers) == 0 || !strings.Contains(err.Error(), "project add") {
				t.Fatalf("boundary %#v %v", p, err)
			}
			if !reflect.DeepEqual(before, projectSnapshotFiles(t, root)) {
				t.Fatal("negative wrote")
			}
		})
	}
}
func projectSnapshotFiles(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !e.IsDir() && e.Type()&os.ModeSymlink == 0 {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			out[path] = string(b)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Reproduce identity drift after confirmed discovery but before ProjectAdd's
// observation. The existing store must never receive the changed identity.
func TestProjectInitIdentityReobservedBeforeMutation(t *testing.T) {
	d, root, w := initFixture(t)
	main := filepath.Join(w, "main")
	initRepo(t, main)
	p, err := InitProject(d, ProjectInitInput{ProjectID: "example"})
	if err != nil {
		t.Fatal(err)
	}
	before := projectSnapshotFiles(t, root)
	original := d.Repos
	calls := 0
	d.Repos = initObserverFunc(func(path string) (RepoObservation, error) {
		observed, err := original.Observe(path)
		calls++
		if calls == 2 {
			observed.GitCommonDir = filepath.Join(w, "changed-common-dir")
		}
		return observed, err
	})
	_, err = InitProject(d, ProjectInitInput{ProjectID: "example", Apply: true, Confirmation: p.Confirmation})
	if err == nil || !strings.Contains(err.Error(), "identity changed") || calls != 2 {
		t.Fatalf("reobservation calls=%d error=%v", calls, err)
	}
	if !reflect.DeepEqual(before, projectSnapshotFiles(t, root)) {
		t.Fatal("identity drift mutated state")
	}
}

type initObserverFunc func(string) (RepoObservation, error)

func (f initObserverFunc) Observe(path string) (RepoObservation, error) { return f(path) }
