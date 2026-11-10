package workspace

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func TestProjectMetadataExplicitRegistrationSurvivesMissingDirectories(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	before, err := os.ReadFile(filepath.Join(root, MarkerDirectory, ProjectsFile))
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := ReadProjectMetadata(f.dependencies, root)
	if err != nil || snapshot.Exists || len(snapshot.Companions("ply")) != 0 || snapshot.RepositoryWrapper("ply") != nil {
		t.Fatalf("absent metadata = %#v %v", snapshot, err)
	}
	planning := filepath.Join(f.wrapper, "planning")
	if err := os.Mkdir(planning, 0o755); err != nil {
		t.Fatal(err)
	}
	input := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: planning, ActorClaim: "codex"}
	out, err := UpdateProjectMetadata(f.dependencies, input)
	if err != nil || !out.Changed || len(out.Companions) != 1 || out.Companions[0].Locator != planning {
		t.Fatalf("add = %#v %v", out, err)
	}
	if out, err := UpdateProjectMetadata(f.dependencies, input); err != nil || out.Changed || out.Revision != 1 {
		t.Fatalf("repeat = %#v %v", out, err)
	}
	input = ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: f.wrapper, ActorClaim: "codex"}
	if out, err := UpdateProjectMetadata(f.dependencies, input); err != nil || !out.Changed {
		t.Fatalf("set wrapper = %#v %v", out, err)
	}
	if err := os.Remove(planning); err != nil {
		t.Fatal(err)
	}
	snapshot, err = ReadProjectMetadata(f.dependencies, root)
	if err != nil || snapshot.Revision != 2 || len(snapshot.Companions("ply")) != 1 || snapshot.RepositoryWrapper("ply") == nil || *snapshot.RepositoryWrapper("ply") != f.wrapper {
		t.Fatalf("lost explicit metadata = %#v %v", snapshot, err)
	}
	input = ProjectMetadataInput{Operation: ProjectCompanionRemove, ProjectID: "ply", Role: CompanionPlanning, Locator: planning, ActorClaim: "human"}
	if out, err := UpdateProjectMetadata(f.dependencies, input); err != nil || !out.Changed || len(out.Companions) != 0 {
		t.Fatalf("remove missing directory = %#v %v", out, err)
	}
	after, err := os.ReadFile(filepath.Join(root, MarkerDirectory, ProjectsFile))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("metadata changed project identities: %v", err)
	}
}

func TestProjectMetadataRemoveMissingCompanionThroughSymlink(t *testing.T) {
	for _, relative := range []string{"planning", filepath.Join("nested", "planning")} {
		t.Run(relative, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			alias := filepath.Join(filepath.Dir(f.wrapper), "wrapper-alias")
			if err := os.Symlink(f.wrapper, alias); err != nil {
				t.Fatal(err)
			}
			physical := filepath.Join(f.wrapper, relative)
			if err := os.MkdirAll(physical, 0o755); err != nil {
				t.Fatal(err)
			}
			input := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: filepath.Join(alias, relative), ActorClaim: "actor"}
			out, err := UpdateProjectMetadata(f.dependencies, input)
			if err != nil || !out.Changed || len(out.Companions) != 1 || out.Companions[0].Locator != physical {
				t.Fatalf("register alias = %#v %v", out, err)
			}
			if err := os.Remove(physical); err != nil {
				t.Fatal(err)
			}
			if filepath.Dir(relative) != "." {
				if err := os.Remove(filepath.Dir(physical)); err != nil {
					t.Fatal(err)
				}
			}
			input.Operation = ProjectCompanionRemove
			out, err = UpdateProjectMetadata(f.dependencies, input)
			if err != nil || !out.Changed || out.Revision != 2 || len(out.Companions) != 0 || out.Event == nil || out.Event.Locator == nil || *out.Event.Locator != physical {
				t.Fatalf("remove missing companion using its original alias = %#v %v", out, err)
			}
		})
	}
}

func TestProjectMetadataReadAndNoOpDoNotWriteAndCASPreservesState(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	if err := os.Remove(filepath.Join(root, MarkerDirectory, projectsLock)); err != nil {
		t.Fatal(err)
	}
	if out, err := ShowProjectMetadata(f.dependencies, "ply"); err != nil || out.Source.Exists || len(out.Repositories) != 1 || out.Repositories[0].Wrapper != nil {
		t.Fatalf("absent read = %#v %v", out, err)
	}
	if _, err := os.Lstat(filepath.Join(root, MarkerDirectory, projectsLock)); !os.IsNotExist(err) {
		t.Fatalf("read created lock: %v", err)
	}
	zero := 0
	input := ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: f.wrapper, ActorClaim: "actor", ExpectedRevision: &zero}
	out, err := UpdateProjectMetadata(f.dependencies, input)
	if err != nil || out.Revision != 1 {
		t.Fatalf("set %#v %v", out, err)
	}
	input.Operation, input.Locator = ProjectRepoWrapperClear, ""
	if _, err := UpdateProjectMetadata(f.dependencies, input); err == nil {
		t.Fatal("stale expected revision accepted")
	}
	input.ExpectedRevision = nil
	input.Operation, input.Locator = ProjectRepoWrapperSet, f.wrapper
	before, err := os.ReadFile(out.Source.Locator)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(out.Source.Locator)
	if err != nil {
		t.Fatal(err)
	}
	out, err = UpdateProjectMetadata(f.dependencies, input)
	if err != nil || out.Changed || out.Revision != 1 {
		t.Fatalf("no-op = %#v %v", out, err)
	}
	after, err := os.ReadFile(out.Source.Locator)
	if err != nil {
		t.Fatal(err)
	}
	afterInfo, err := os.Stat(out.Source.Locator)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) || !info.ModTime().Equal(afterInfo.ModTime()) {
		t.Fatal("no-op rewrote metadata")
	}
	input.Operation, input.Locator = ProjectRepoWrapperClear, ""
	out, err = UpdateProjectMetadata(f.dependencies, input)
	if err != nil || !out.Changed || out.Revision != 2 || out.Repositories[0].Wrapper != nil {
		t.Fatalf("clear %#v %v", out, err)
	}
	if out, err := UpdateProjectMetadata(f.dependencies, input); err != nil || out.Changed || out.Revision != 2 {
		t.Fatalf("clear no-op %#v %v", out, err)
	}
}

func TestProjectMetadataRejectsCorruptionAndSymlinks(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	input := ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: f.wrapper, ActorClaim: "actor"}
	out, err := UpdateProjectMetadata(f.dependencies, input)
	if err != nil {
		t.Fatal(err)
	}
	good, err := os.ReadFile(out.Source.Locator)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range [][]byte{
		{},
		[]byte(strings.Replace(string(good), `"schema_version": 1`, `"schema_version": 99`, 1)),
		[]byte(strings.Replace(string(good), `"schema_version": 1`, `"schema_version": 1, "schema_version": 1`, 1)),
		[]byte(strings.Replace(string(good), `"actor_claim": "actor"`, `"actor_claim": "other"`, 1)),
		[]byte(`{"kind":"WorkspaceProjectMetadataStore@1","schema_version":1,"events":null}`),
	} {
		if err := os.WriteFile(out.Source.Locator, bad, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadProjectMetadata(f.dependencies, root); err == nil {
			t.Fatal("corrupt store read as empty")
		}
		if _, err := UpdateProjectMetadata(f.dependencies, input); err == nil {
			t.Fatal("corrupt store overwritten")
		}
		after, err := os.ReadFile(out.Source.Locator)
		if err != nil || !bytes.Equal(bad, after) {
			t.Fatalf("corruption changed %v", err)
		}
	}
	if err := os.Remove(out.Source.Locator); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "absent"), out.Source.Locator); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadProjectMetadata(f.dependencies, root); err == nil {
		t.Fatal("dangling symlink accepted")
	}
	if _, err := UpdateProjectMetadata(f.dependencies, input); err == nil {
		t.Fatal("dangling symlink replaced")
	}
}

func TestProjectMetadataParallelRegistrationsUseOneHistory(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	var group sync.WaitGroup
	errorsOut := make(chan error, 3)
	for _, role := range []CompanionRole{CompanionPlanning, CompanionDocs, CompanionOther} {
		group.Add(1)
		go func(role CompanionRole) {
			defer group.Done()
			_, err := UpdateProjectMetadata(f.dependencies, ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: role, Locator: f.wrapper, ActorClaim: "parallel-agent"})
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
	snapshot, err := ReadProjectMetadata(f.dependencies, root)
	if err != nil || snapshot.Revision != 3 || len(snapshot.Companions("ply")) != 3 {
		t.Fatalf("lost changes: %#v %v", snapshot, err)
	}
}

func TestProjectMetadataAtomicFailureAndRetry(t *testing.T) {
	for _, stage := range []string{"replace", "directory-sync"} {
		t.Run(stage, func(t *testing.T) {
			f := newWorkItemJourneyFixture(t)
			root := filepath.Dir(f.wrapper)
			input := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionDocs, Locator: f.wrapper, ActorClaim: "actor"}
			injected := errors.New("injected")
			_, err := updateProjectMetadata(f.dependencies, input, workspaceMetadataPublisher{fault: func(point string) error {
				if point == stage {
					return injected
				}
				return nil
			}})
			if !errors.Is(err, injected) {
				t.Fatalf("fault=%v", err)
			}
			snapshot, err := ReadProjectMetadata(f.dependencies, root)
			if err != nil {
				t.Fatal(err)
			}
			want := 0
			if stage == "directory-sync" {
				want = 1
			}
			if snapshot.Revision != want || len(snapshot.Companions("ply")) != want {
				t.Fatalf("partial state %#v", snapshot)
			}
			out, err := UpdateProjectMetadata(f.dependencies, input)
			if err != nil || out.Revision != 1 || out.Changed != (stage == "replace") {
				t.Fatalf("retry %#v %v", out, err)
			}
		})
	}
}

func TestProjectMetadataRejectsInvalidActorSubjectAndPathsWithoutPublishing(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	for _, input := range []ProjectMetadataInput{
		{Operation: ProjectCompanionAdd, ProjectID: "missing", Role: CompanionPlanning, Locator: f.wrapper, ActorClaim: "actor"},
		{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: "backup", Locator: f.wrapper, ActorClaim: "actor"},
		{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: f.wrapper, ActorClaim: ""},
		{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: filepath.Join(f.wrapper, "missing"), ActorClaim: "actor"},
		{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "unknown", Locator: f.wrapper, ActorClaim: "actor"},
		{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: t.TempDir(), ActorClaim: "actor"},
	} {
		if _, err := UpdateProjectMetadata(f.dependencies, input); err == nil {
			t.Fatalf("invalid input accepted %#v", input)
		}
	}
	snapshot, err := ReadProjectMetadata(f.dependencies, root)
	if err != nil || snapshot.Exists {
		t.Fatalf("invalid input wrote metadata: %#v %v", snapshot, err)
	}
}

func TestProjectMetadataRejectsNonUTF8PathBeforeFilesystemAccess(t *testing.T) {
	in := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionDocs, Locator: string([]byte{0xff}), ActorClaim: "actor"}
	if err := in.Validate(); err == nil {
		t.Fatal("non-UTF-8 path could be replaced during JSON serialization")
	}
}

func TestProjectMetadataNormalizesExplicitPathsWithoutScanning(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	root := filepath.Dir(f.wrapper)
	planning := filepath.Join(f.wrapper, "planning")
	if err := os.Mkdir(planning, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(planning, filepath.Join(root, "planning-alias")); err != nil {
		t.Fatal(err)
	}
	in := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: "planning-alias", ActorClaim: "actor"}
	out, err := UpdateProjectMetadata(f.dependencies, in)
	if err != nil || len(out.Companions) != 1 || out.Companions[0].Locator != planning {
		t.Fatalf("physical metadata = %#v %v", out, err)
	}
	in.Locator = planning
	if repeated, err := UpdateProjectMetadata(f.dependencies, in); err != nil || repeated.Changed || repeated.Revision != 1 {
		t.Fatalf("alias made a second registration: %#v %v", repeated, err)
	}
}

func TestProjectMetadataAppearsInExistingListAndShowJSON(t *testing.T) {
	f := newWorkItemJourneyFixture(t)
	for _, input := range []ProjectMetadataInput{
		{Operation: ProjectCompanionAdd, ProjectID: "ply", Role: CompanionPlanning, Locator: f.wrapper, ActorClaim: "actor"},
		{Operation: ProjectRepoWrapperSet, ProjectID: "ply", RepoID: "ply", Locator: f.wrapper, ActorClaim: "actor"},
	} {
		if _, err := UpdateProjectMetadata(f.dependencies, input); err != nil {
			t.Fatal(err)
		}
	}
	show, err := ShowProject(f.dependencies, "ply")
	if err != nil {
		t.Fatal(err)
	}
	b, err := MarshalProject(show)
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Project struct {
			Companions []ProjectCompanion `json:"companions"`
		} `json:"project"`
		Repositories []ProjectRepositoryMetadata `json:"repositories"`
	}
	if err := json.Unmarshal(b, &decoded); err != nil || len(decoded.Project.Companions) != 1 || len(decoded.Repositories) != 1 || decoded.Repositories[0].Wrapper == nil || *decoded.Repositories[0].Wrapper != f.wrapper {
		t.Fatalf("project show metadata: %s %v", b, err)
	}
	list, err := ListProjects(f.dependencies)
	if err != nil {
		t.Fatal(err)
	}
	b, err = MarshalProjectList(list)
	if err != nil {
		t.Fatal(err)
	}
	var listed struct {
		Projects []struct {
			Companions []ProjectCompanion `json:"companions"`
		} `json:"projects"`
	}
	if err := json.Unmarshal(b, &listed); err != nil || len(listed.Projects) != 1 || len(listed.Projects[0].Companions) != 1 {
		t.Fatalf("project list metadata: %s %v", b, err)
	}
}
