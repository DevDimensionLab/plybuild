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
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

func (f definitionFixture) metadataDependencies() Dependencies {
	d := SystemDependencies()
	d.Files = projectCwdFileSystem{d.Files, f.root}
	d.RepoUses = f.history
	return d
}

func (f definitionFixture) update(t *testing.T, in ProjectMetadataInput) ProjectMetadataMutation {
	t.Helper()
	out, err := UpdateProjectMetadata(f.metadataDependencies(), in)
	if err != nil {
		t.Fatalf("%+v: %v", in, err)
	}
	return out
}

func requireAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("%s exists: %v", path, err)
	}
}

func TestProjectMetadataIsStoredInDefinitions(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	if err := os.Mkdir(filepath.Join(f.root, "ply", "planning"), 0o755); err != nil {
		t.Fatal(err)
	}
	outside := physicalPath(t, t.TempDir())
	wrapper := filepath.Join(f.root, "ply")
	changes := []ProjectMetadataInput{
		{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: "ply/planning"},
		{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionDocs, Locator: outside},
		{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: wrapper},
	}
	for _, in := range changes {
		if out := f.update(t, in); !out.Changed {
			t.Fatalf("%+v unchanged", in)
		}
	}
	defs := f.definitions(t)
	if want := []workspacedef.Companion{{Role: "planning", Path: "ply/planning"}, {Role: "docs", Path: outside}}; !reflect.DeepEqual(defs.Projects[0].Companions, want) {
		t.Fatalf("companions %+v", defs.Projects[0].Companions)
	}
	if want := (workspacedef.Repo{ID: "ply", Path: "ply/main", Wrapper: "ply"}); !reflect.DeepEqual(defs.Projects[0].Repos[0], want) {
		t.Fatalf("repository %+v", defs.Projects[0].Repos[0])
	}
	if defs.Projects[1].Companions != nil || defs.Projects[1].Repos[0].Wrapper != "" {
		t.Fatalf("other Project changed: %+v", defs.Projects[1])
	}
	requireAbsent(t, filepath.Join(f.ply, "project-metadata.json"))

	state := f.state(t)
	for _, in := range changes {
		if out := f.update(t, in); out.Changed {
			t.Fatalf("%+v changed again", in)
		}
	}
	f.requireUnchanged(t, state)
}

func TestProjectMetadataReadbackResolvesPhysicalPaths(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	planning := filepath.Join(f.root, "ply", "planning")
	if err := os.Mkdir(planning, 0o755); err != nil {
		t.Fatal(err)
	}
	f.update(t, ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: "ply/planning"})
	out := f.update(t, ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: "ply"})
	if want := []ProjectCompanion{{CompanionPlanning, planning}}; !reflect.DeepEqual(out.Project.Companions, want) {
		t.Fatalf("companions %+v", out.Project.Companions)
	}
	if len(out.Repos) != 1 || out.Repos[0].Wrapper != filepath.Join(f.root, "ply") {
		t.Fatalf("repositories %+v", out.Repos)
	}
}

func TestProjectMetadataSurvivesMissingDirectoriesAndClears(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	planning := filepath.Join(f.root, "ply", "planning")
	if err := os.Mkdir(planning, 0o755); err != nil {
		t.Fatal(err)
	}
	add := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: planning}
	f.update(t, add)
	f.update(t, ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: filepath.Join(f.root, "ply")})
	if err := os.Remove(planning); err != nil {
		t.Fatal(err)
	}
	show, err := ShowProject(f.metadataDependencies(), "ply")
	if err != nil || len(show.Project.Companions) != 1 || show.Project.Companions[0].Locator != planning {
		t.Fatalf("lost companion of a missing directory: %+v %v", show.Project, err)
	}
	remove := add
	remove.Operation = ProjectCompanionRemove
	if out := f.update(t, remove); !out.Changed || len(out.Project.Companions) != 0 {
		t.Fatalf("remove missing directory %+v", out)
	}
	clear := ProjectMetadataInput{Operation: ProjectRepoWrapperClear, ProjectID: "ply", RepoID: "ply"}
	if out := f.update(t, clear); !out.Changed || out.Repos[0].Wrapper != "" {
		t.Fatalf("clear %+v", out)
	}
	if defs := f.definitions(t); len(defs.Projects[0].Companions) != 0 || defs.Projects[0].Repos[0].Wrapper != "" {
		t.Fatalf("definitions %+v", defs.Projects[0])
	}
	raw, err := os.ReadFile(filepath.Join(f.root, workspacedef.File))
	if err != nil || strings.Contains(string(raw), "companions") || strings.Count(string(raw), `"wrapper"`) != 2 {
		t.Fatalf("empty metadata is written: %s %v", raw, err)
	}
	state := f.state(t)
	if out := f.update(t, clear); out.Changed {
		t.Fatalf("clear again %+v", out)
	}
	if out := f.update(t, remove); out.Changed {
		t.Fatalf("remove again %+v", out)
	}
	f.requireUnchanged(t, state)
}

func TestProjectMetadataRemoveMissingCompanionThroughSymlink(t *testing.T) {
	t.Parallel()
	for _, relative := range []string{"planning", filepath.Join("nested", "planning")} {
		t.Run(relative, func(t *testing.T) {
			f := newDefinitionFixture(t)
			wrapper := filepath.Join(f.root, "ply")
			alias := filepath.Join(f.root, "wrapper-alias")
			if err := os.Symlink(wrapper, alias); err != nil {
				t.Fatal(err)
			}
			physical := filepath.Join(wrapper, relative)
			if err := os.MkdirAll(physical, 0o755); err != nil {
				t.Fatal(err)
			}
			in := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: filepath.Join(alias, relative)}
			if out := f.update(t, in); !out.Changed || len(out.Project.Companions) != 1 || out.Project.Companions[0].Locator != physical {
				t.Fatalf("register alias %+v", out.Project)
			}
			if err := os.RemoveAll(filepath.Join(wrapper, strings.Split(relative, string(filepath.Separator))[0])); err != nil {
				t.Fatal(err)
			}
			in.Operation = ProjectCompanionRemove
			if out := f.update(t, in); !out.Changed || len(out.Project.Companions) != 0 {
				t.Fatalf("remove missing companion using its original alias %+v", out.Project)
			}
		})
	}
}

func TestProjectMetadataNormalizesAliases(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	planning := filepath.Join(f.root, "ply", "planning")
	if err := os.Mkdir(planning, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(planning, filepath.Join(f.root, "planning-alias")); err != nil {
		t.Fatal(err)
	}
	in := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: "planning-alias"}
	if out := f.update(t, in); len(out.Project.Companions) != 1 || out.Project.Companions[0].Locator != planning {
		t.Fatalf("physical metadata %+v", out.Project)
	}
	in.Locator = planning
	if out := f.update(t, in); out.Changed {
		t.Fatalf("alias made a second registration: %+v", out.Project)
	}
}

func TestProjectMetadataRejectsBeforeWriting(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct {
		edit  func(string) string
		in    ProjectMetadataInput
		check func(*testing.T, error)
	}{
		"invalid JSON": {edit: func(s string) string { return s + "}" },
			in: ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: "ply"}},
		"unknown field": {edit: func(s string) string {
			return strings.Replace(s, `"format_version": 1,`, `"format_version": 1, "extra": 1,`, 1)
		}, in: ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: "ply"}},
		"invalid companion role in definitions": {edit: func(s string) string {
			return strings.Replace(s, `"repos": [{"id": "app", "path": "app/main"}]`, `"repos": [{"id": "app", "path": "app/main"}], "companions": [{"role": "backup", "path": "app"}]`, 1)
		}, in: ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: "ply"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectInvalidArguments) }},
		"unknown Project": {in: ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "missing", Role: CompanionDocs, Locator: "ply"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectNotFound) }},
		"invalid role": {in: ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: "backup", Locator: "ply"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectInvalidArguments) }},
		"missing companion directory": {in: ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionDocs, Locator: "gone"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectPath) }},
		"repository not a member": {in: ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "app", Locator: "app"},
			check: func(t *testing.T, err error) { assertWorkErrorClass(t, err, ErrorWorkProjectConflict) }},
		"missing wrapper": {in: ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: "gone"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectPath) }},
		"wrapper without repository": {in: ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: "app"},
			check: func(t *testing.T, err error) { assertProjectErrorClass(t, err, ErrorProjectInvalidArguments) }},
	} {
		t.Run(name, func(t *testing.T) {
			f := newDefinitionFixture(t)
			if tc.edit != nil {
				mustWriteFile(t, filepath.Join(f.root, workspacedef.File), tc.edit(fixtureDefinitions))
			}
			state := f.state(t)
			_, err := UpdateProjectMetadata(f.metadataDependencies(), tc.in)
			if tc.check != nil {
				tc.check(t, err)
			} else {
				assertProjectErrorClass(t, err, ErrorProjectStoreConflict)
			}
			f.requireUnchanged(t, state)
		})
	}
}

func TestProjectMetadataRejectsNonUTF8PathBeforeFilesystemAccess(t *testing.T) {
	t.Parallel()
	in := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionDocs, Locator: string([]byte{0xff})}
	if err := in.Validate(); err == nil {
		t.Fatal("non-UTF-8 path could be replaced during JSON serialization")
	}
}

func TestProjectMetadataParallelRegistrationsAreAllKept(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	var group sync.WaitGroup
	errorsOut := make(chan error, 3)
	for _, role := range []CompanionRole{CompanionPlanning, CompanionDocs, CompanionOther} {
		group.Add(1)
		go func(role CompanionRole) {
			defer group.Done()
			_, err := UpdateProjectMetadata(f.metadataDependencies(), ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: role, Locator: "ply"})
			errorsOut <- err
		}(role)
	}
	group.Wait()
	close(errorsOut)
	for err := range errorsOut {
		if err != nil {
			t.Fatal(err)
		}
	}
	if companions := f.definitions(t).Projects[0].Companions; len(companions) != 3 {
		t.Fatalf("lost changes: %+v", companions)
	}
}

func TestProjectMetadataWaitsForProjectsLock(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	lock, err := filelock.Acquire(filepath.Join(f.ply, projectsLock), filelock.Options{})
	if err != nil {
		t.Fatal(err)
	}
	type outcome struct {
		out ProjectMetadataMutation
		err error
	}
	done := make(chan outcome, 1)
	go func() {
		out, err := UpdateProjectMetadata(f.metadataDependencies(), ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: "ply"})
		done <- outcome{out, err}
	}()
	select {
	case got := <-done:
		t.Fatalf("metadata did not wait for projects.lock: %+v %v", got.out, got.err)
	case <-time.After(200 * time.Millisecond):
	}
	if err := lock.Release(); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-done:
		if got.err != nil || !got.out.Changed {
			t.Fatalf("after release %+v %v", got.out, got.err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("metadata did not finish after release")
	}
}

// plyapp-server reads project.companions[] {role, locator} and
// repositories[].wrapper (string or null) from project show.
func TestProjectMetadataInProjectShowAndListJSON(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	if _, err := f.add(t, f.root, "ply", "ply-code-map", "ply-code-map"); err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(f.root, "ply")
	f.update(t, ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: wrapper})
	f.update(t, ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: wrapper})
	show, err := ShowProject(f.metadataDependencies(), "ply")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := MarshalProject(show)
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Project struct {
			Companions []map[string]any `json:"companions"`
		} `json:"project"`
		Repositories []map[string]any `json:"repositories"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	if want := []map[string]any{{"role": "planning", "locator": wrapper}}; !reflect.DeepEqual(decoded.Project.Companions, want) {
		t.Fatalf("companions in %s", raw)
	}
	wrappers := map[any]any{}
	for _, r := range decoded.Repositories {
		value, present := r["wrapper"]
		if !present {
			t.Fatalf("repository without wrapper key in %s", raw)
		}
		wrappers[r["repo_id"]] = value
	}
	if want := map[any]any{"ply": wrapper, "ply-code-map": nil}; !reflect.DeepEqual(wrappers, want) {
		t.Fatalf("wrappers %v in %s", wrappers, raw)
	}

	list, err := ListProjects(f.metadataDependencies())
	if err != nil {
		t.Fatal(err)
	}
	raw, err = MarshalProjectList(list)
	if err != nil {
		t.Fatal(err)
	}
	var listed struct {
		Projects []struct {
			ID         ProjectID        `json:"project_id"`
			Companions []map[string]any `json:"companions"`
		} `json:"projects"`
	}
	if err := json.Unmarshal(raw, &listed); err != nil || len(listed.Projects) != 2 || listed.Projects[0].ID != "app" || listed.Projects[0].Companions == nil || len(listed.Projects[0].Companions) != 0 || len(listed.Projects[1].Companions) != 1 {
		t.Fatalf("project list metadata: %s %v", raw, err)
	}
}

// A metadata change answers with project show and whether it changed.
func TestProjectMetadataMutationJSONIsProjectShowWithChanged(t *testing.T) {
	t.Parallel()
	f := newDefinitionFixture(t)
	out := f.update(t, ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionDocs, Locator: "ply"})
	raw, err := MarshalProjectMetadataMutation(out)
	if err != nil {
		t.Fatal(err)
	}
	show, err := ShowProject(f.metadataDependencies(), "ply")
	if err != nil {
		t.Fatal(err)
	}
	shown, err := MarshalProject(show)
	if err != nil {
		t.Fatal(err)
	}
	var mutation, readback map[string]any
	if err := json.Unmarshal(raw, &mutation); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(shown, &readback); err != nil {
		t.Fatal(err)
	}
	if mutation["changed"] != true {
		t.Fatalf("changed missing in %s", raw)
	}
	delete(mutation, "changed")
	if !reflect.DeepEqual(mutation, readback) {
		t.Fatalf("mutation %s is not project show %s", raw, shown)
	}
}

func TestProjectsLockErrorMatchesSnapshotLock(t *testing.T) {
	t.Parallel()
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
