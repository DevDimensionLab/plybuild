package workspace

import (
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

func (f definitionFixture) metadataDependencies() Dependencies {
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{d.Files, f.root}
	d.WorkClock = fixedWorkClock{time.Date(2026, 11, 8, 12, 0, 0, 0, time.UTC)}
	return d
}

func (f definitionFixture) setWrapper(project ProjectID, repo RepoID, wrapper string) (ProjectMetadataMutation, error) {
	input := ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: project, RepoID: repo, Locator: wrapper, ActorClaim: "codex"}
	return UpdateProjectMetadata(f.metadataDependencies(), input)
}

func (f definitionFixture) metadataPath() string {
	return filepath.Join(f.ply, ProjectMetadataFile)
}

// newMetadataFixture is the definition fixture as an Epic workspace: metadata
// writes bootstrap defined Epics, so the defined Epic is committed and adopted
// natively instead of being read from the partial K5 history registry.
func newMetadataFixture(t *testing.T) definitionFixture {
	t.Helper()
	f := newDefinitionFixture(t)
	if err := os.Remove(filepath.Join(f.ply, WorkItemsFile)); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(f.root, "ply", "main")
	runLocalGit(t, repo, "symbolic-ref", "HEAD", "refs/heads/master")
	runLocalGit(t, repo, "-c", "user.name=Ply tests", "-c", "user.email=tests@example.invalid", "commit", "-q", "--allow-empty", "-m", "fixture")
	if _, err := ObserveContaining(f.metadataDependencies()); err != nil {
		t.Fatal(err)
	}
	return f
}

// addLegacyProject seeds the legacy registry, which native Project add refuses
// to change while definitions exist.
func (f definitionFixture) addLegacyProject(t *testing.T, input ProjectAddInput) {
	t.Helper()
	definitions, aside := filepath.Join(f.root, workspacedef.File), filepath.Join(f.root, workspacedef.File+".aside")
	if err := os.Rename(definitions, aside); err != nil {
		t.Fatal(err)
	}
	_, err := AddProject(f.metadataDependencies(), input)
	if renameErr := os.Rename(aside, definitions); renameErr != nil {
		t.Fatal(renameErr)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func requireAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("%s exists: %v", path, err)
	}
}

func repositoryWrapper(readback ProjectMetadataReadback, repo RepoID) *string {
	for _, r := range readback.Repositories {
		if r.RepoID == repo {
			return r.Wrapper
		}
	}
	return nil
}

func TestDefinitionMetadataAddSetShowAndRerun(t *testing.T) {
	f := newMetadataFixture(t)
	if result, err := f.add(t, f.root, "ply", "ply-code-map", "./ply-code-map"); err != nil || !result.Changed {
		t.Fatalf("add %+v %v", result, err)
	}
	before := f.state(t)
	wrapper := filepath.Join(f.root, "ply-code-map")

	out, err := f.setWrapper("ply", "ply-code-map", wrapper)
	if err != nil || !out.Changed || out.Revision != 1 || out.Event == nil {
		t.Fatalf("set %+v %v", out, err)
	}
	show, err := ShowProjectMetadata(f.metadataDependencies(), "ply")
	if err != nil || show.Revision != 1 || len(show.Events) != 1 || show.Workspace.Root != f.root {
		t.Fatalf("show %+v %v", show, err)
	}
	if got := repositoryWrapper(show, "ply-code-map"); got == nil || *got != wrapper {
		t.Fatalf("wrapper %v in %+v", got, show.Repositories)
	}
	if got := repositoryWrapper(show, "ply"); got != nil || len(show.Repositories) != 2 {
		t.Fatalf("repositories %+v", show.Repositories)
	}
	after := f.state(t)
	if _, created := after[f.metadataPath()]; !created {
		t.Fatalf("metadata not published: %v", after)
	}
	delete(after, f.metadataPath())
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("other files changed:\nbefore %v\nafter  %v", before, after)
	}
	requireAbsent(t, filepath.Join(f.ply, ProjectsFile))

	published := f.state(t)
	rerun, err := f.setWrapper("ply", "ply-code-map", wrapper)
	if err != nil || rerun.Changed || rerun.Event != nil || rerun.Revision != 1 {
		t.Fatalf("rerun %+v %v", rerun, err)
	}
	f.requireUnchanged(t, published)
}

func TestDefinitionMetadataCompanions(t *testing.T) {
	f := newMetadataFixture(t)
	planning := filepath.Join(f.root, "ply")
	input := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "app", Role: CompanionPlanning, Locator: planning, ActorClaim: "codex"}
	if out, err := UpdateProjectMetadata(f.metadataDependencies(), input); err != nil || !out.Changed || len(out.Companions) != 1 {
		t.Fatalf("add %+v %v", out, err)
	}
	input.Operation = ProjectCompanionRemove
	if out, err := UpdateProjectMetadata(f.metadataDependencies(), input); err != nil || !out.Changed || len(out.Companions) != 0 || out.Revision != 2 {
		t.Fatalf("remove %+v %v", out, err)
	}
}

func TestDefinitionMetadataRejectsBeforeWriting(t *testing.T) {
	type request struct {
		project ProjectID
		repo    RepoID
		wrapper string
	}
	plyRequest := request{"ply", "ply", "ply"}
	legacyPly := func(t *testing.T, f definitionFixture) {
		input := ProjectAddInput{ProjectID: "ply", Name: "Ply", Wrapper: "ply", Repos: []RepoInput{{RepoID: "ply", Path: "ply/main"}}}
		f.addLegacyProject(t, input)
	}
	for name, tc := range map[string]struct {
		edit    func(string) string
		legacy  func(*testing.T, definitionFixture)
		request request
		check   func(*testing.T, error)
	}{
		"invalid JSON":                     {edit: func(s string) string { return s + "}" }, request: plyRequest},
		"invalid JSON with legacy Project": {edit: func(s string) string { return s + "}" }, legacy: legacyPly, request: plyRequest},
		"unknown field": {edit: func(s string) string {
			return strings.Replace(s, `"format_version": 1,`, `"format_version": 1, "extra": 1,`, 1)
		}, request: plyRequest},
		"format version": {edit: func(s string) string { return strings.Replace(s, `"format_version": 1`, `"format_version": 2`, 1) }, request: plyRequest},
		"duplicate Project": {edit: func(s string) string { return strings.Replace(s, `"id": "app", "name"`, `"id": "ply", "name"`, 1) }, request: plyRequest,
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectInvalidArguments) }},
		"Project only in legacy": {legacy: func(t *testing.T, f definitionFixture) {
			input := ProjectAddInput{ProjectID: "legacy", Name: "Legacy", Wrapper: "retired", Repos: []RepoInput{{RepoID: "retired", Path: "retired"}}}
			f.addLegacyProject(t, input)
		}, request: request{"legacy", "retired", "retired"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectNotFound) }},
		"repository not a member": {request: request{"ply", "app", "app"},
			check: func(t *testing.T, err error) { assertWorkErrorClass(t, err, ErrorWorkProjectConflict) }},
		"missing wrapper": {request: request{"ply", "ply", "gone"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectPath) }},
		"wrapper without repository": {request: request{"ply", "ply", "app"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectInvalidArguments) }},
		"repository with different paths": {edit: func(s string) string {
			return strings.Replace(s, `"repos": [{"id": "app", "path": "app/main"}]`, `"repos": [{"id": "app", "path": "app/main"}, {"id": "ply", "path": "app/main"}]`, 1)
		}, request: plyRequest, check: func(t *testing.T, err error) {
			assertProjectErrorClass(t, err, ErrorProjectConflict)
			if !strings.Contains(err.Error(), "repository ply is defined with different paths in "+workspacedef.File) {
				t.Fatalf("error %v", err)
			}
		}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newMetadataFixture(t)
			if tc.legacy != nil {
				tc.legacy(t, f)
			}
			if tc.edit != nil {
				mustWriteFile(t, filepath.Join(f.root, workspacedef.File), tc.edit(fixtureDefinitions))
			}
			state := f.state(t)
			_, err := f.setWrapper(tc.request.project, tc.request.repo, filepath.Join(f.root, tc.request.wrapper))
			if tc.check != nil {
				tc.check(t, err)
			} else {
				assertProjectErrorClass(t, err, ErrorProjectStoreConflict)
			}
			f.requireUnchanged(t, state)
			requireAbsent(t, f.metadataPath())
		})
	}
}

func TestDefinitionMetadataIgnoresLegacyRegistry(t *testing.T) {
	t.Run("corrupt", func(t *testing.T) {
		f := newMetadataFixture(t)
		legacy := filepath.Join(f.ply, ProjectsFile)
		mustWriteFile(t, legacy, "{{{ not yaml")
		if out, err := f.setWrapper("ply", "ply", filepath.Join(f.root, "ply")); err != nil || !out.Changed {
			t.Fatalf("set %+v %v", out, err)
		}
		if show, err := ShowProjectMetadata(f.metadataDependencies(), "ply"); err != nil || show.Revision != 1 {
			t.Fatalf("show %+v %v", show, err)
		}
		if raw, err := os.ReadFile(legacy); err != nil || string(raw) != "{{{ not yaml" {
			t.Fatalf("legacy registry changed: %q %v", raw, err)
		}
	})
	t.Run("same Project ID", func(t *testing.T) {
		f := newMetadataFixture(t)
		input := ProjectAddInput{ProjectID: "ply", Name: "Legacy Ply", Wrapper: "app", Repos: []RepoInput{{RepoID: "app", Path: "app/main"}}}
		f.addLegacyProject(t, input)
		show, err := ShowProjectMetadata(f.metadataDependencies(), "ply")
		if err != nil || len(show.Repositories) != 1 || show.Repositories[0].RepoID != "ply" {
			t.Fatalf("show %+v %v", show, err)
		}
		_, err = f.setWrapper("ply", "app", filepath.Join(f.root, "app"))
		assertWorkErrorClass(t, err, ErrorWorkProjectConflict)
	})
}

func TestDefinitionMetadataShowNotFound(t *testing.T) {
	f := newMetadataFixture(t)
	_, err := ShowProjectMetadata(f.metadataDependencies(), "other")
	assertProjectErrorClass(t, err, ErrorProjectNotFound)
	if want := fmt.Sprintf("project other is not registered in Ply workspace %s", f.root); !strings.Contains(err.Error(), want) {
		t.Fatalf("error %v", err)
	}
}

func TestDefinitionMetadataShowRejectsInvalidDefinitionsWithoutLegacyFallback(t *testing.T) {
	f := newMetadataFixture(t)
	input := ProjectAddInput{ProjectID: "ply", Name: "Ply", Wrapper: "ply", Repos: []RepoInput{{RepoID: "ply", Path: "ply/main"}}}
	f.addLegacyProject(t, input)
	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), fixtureDefinitions+"}")
	_, err := ShowProjectMetadata(f.metadataDependencies(), "ply")
	assertProjectErrorClass(t, err, ErrorProjectStoreConflict)
}

func TestDefinitionMetadataShowRejectsCorruptMetadataStore(t *testing.T) {
	f := newMetadataFixture(t)
	mustWriteFile(t, f.metadataPath(), "{")
	_, err := ShowProjectMetadata(f.metadataDependencies(), "ply")
	assertWorkErrorClass(t, err, ErrorWorkStoreConflict)
	if !strings.Contains(err.Error(), "cannot read metadata "+f.metadataPath()) {
		t.Fatalf("error %v", err)
	}
	if raw, err := os.ReadFile(f.metadataPath()); err != nil || string(raw) != "{" {
		t.Fatalf("metadata store changed: %q %v", raw, err)
	}
}

func TestDefinitionMetadataWaitsForProjectsLock(t *testing.T) {
	f := newMetadataFixture(t)
	lock, err := filelock.Acquire(filepath.Join(f.ply, projectsLock), filelock.Options{})
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan struct{})
	type outcome struct {
		out ProjectMetadataMutation
		err error
	}
	done := make(chan outcome, 1)
	go func() {
		close(started)
		out, err := f.setWrapper("ply", "ply", filepath.Join(f.root, "ply"))
		done <- outcome{out, err}
	}()
	<-started
	select {
	case got := <-done:
		t.Fatalf("metadata did not wait for projects.lock: %+v %v", got.out, got.err)
	case <-time.After(200 * time.Millisecond):
	}
	requireAbsent(t, f.metadataPath())
	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-done:
		if got.err != nil || !got.out.Changed || got.out.Revision != 1 {
			t.Fatalf("after release %+v %v", got.out, got.err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("metadata did not finish after release")
	}
}

func TestProjectsLockErrorMatchesSnapshotLock(t *testing.T) {
	path := "/w/.ply/projects.lock"
	cause := errors.New("cause")
	lockErr := &filelock.Error{Op: "not-regular", Path: path}
	releaseErr := &filelock.ReleaseError{Path: path, Unlock: cause}
	operation := ProjectInvalidArguments("from callback")
	for name, tc := range map[string]struct{ err, want error }{
		"lock":                  {lockErr, projectLockError(path, lockErr)},
		"release":               {releaseErr, projectReleaseError(releaseErr)},
		"release after failure": {errors.Join(operation, releaseErr), projectIOError("release after failure", path, fmt.Errorf("%v; %w", operation, projectReleaseError(releaseErr)))},
	} {
		t.Run(name, func(t *testing.T) {
			got := projectsLockError(tc.err)
			var gotTyped, wantTyped *ProjectError
			if !errors.As(got, &gotTyped) || !errors.As(tc.want, &wantTyped) || gotTyped.Class != wantTyped.Class || got.Error() != tc.want.Error() {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
	if got := projectsLockError(operation); got != operation {
		t.Fatalf("callback error changed: %v", got)
	}
	if got := projectsLockError(nil); got != nil {
		t.Fatalf("nil became %v", got)
	}
}
