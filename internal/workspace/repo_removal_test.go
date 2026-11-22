package workspace

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// registrationsStore answers RepoUses with fixed uses or a fixed error.
type registrationsStore struct {
	uses []RepoUse
	err  error
}

func (s registrationsStore) RepoUses(ProjectWork) ([]RepoUse, error) {
	return s.uses, s.err
}

// newEpicFixture is the definition fixture as an Epic workspace: the defined
// Epic is committed and adopted natively instead of being read from the
// partial K5 history.
func newEpicFixture(t *testing.T) definitionFixture {
	t.Helper()
	f := newDefinitionFixture(t)
	delete(f.history, "ply")
	repo := filepath.Join(f.root, "ply", "main")
	runLocalGit(t, repo, "symbolic-ref", "HEAD", "refs/heads/master")
	runLocalGit(t, repo, "-c", "user.name=Ply tests", "-c", "user.email=tests@example.invalid", "commit", "-q", "--allow-empty", "-m", "fixture")
	if _, err := ObserveContaining(f.metadataDependencies()); err != nil {
		t.Fatal(err)
	}
	return f
}

// newRemovalFixture is the Epic fixture with repository ply-code-map
// defined in Project ply at the layout repo create produces: ply/ply-code-map/main.
func newRemovalFixture(t *testing.T) definitionFixture {
	t.Helper()
	f := newEpicFixture(t)
	repo := filepath.Join(f.root, "ply", "ply-code-map", "main")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, repo, "init", "-q")
	if _, err := f.add(t, f.root, "ply", "ply-code-map", "ply/ply-code-map/main"); err != nil {
		t.Fatal(err)
	}
	return f
}

func (f definitionFixture) editDefinitions(t *testing.T, edit func(*workspacedef.Definitions)) {
	t.Helper()
	defs := f.definitions(t)
	edit(&defs)
	raw, err := json.MarshalIndent(defs, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(f.root, workspacedef.File), string(raw)+"\n")
}

// alias creates the directory target and a symlink to it at root/name.
func (f definitionFixture) alias(t *testing.T, name, target string) {
	t.Helper()
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(f.root, name)); err != nil {
		t.Fatal(err)
	}
}

func projectClass(err error) ErrorClass {
	var project *ProjectError
	if errors.As(err, &project) {
		return project.Class
	}
	return ""
}

func TestCheckRepoRemovalRejectsReferencesWithoutWriting(t *testing.T) {
	t.Parallel()
	const repo RepoID = "ply-code-map"
	for name, tc := range map[string]struct {
		setup    func(*testing.T, definitionFixture, *Dependencies)
		project  string
		repo     RepoID
		class    ErrorClass
		workItem WorkItemErrorClass
		reason   string
	}{
		"clean": {setup: func(*testing.T, definitionFixture, *Dependencies) {}},
		"finished work": {setup: func(_ *testing.T, f definitionFixture, d *Dependencies) {
			locator := filepath.Join(f.root, "ply", "ply-code-map", "epic")
			d.RepoUses = registrationsStore{uses: []RepoUse{{Kind: RepoUseEpic, ID: "e", RepoID: repo, Worktree: locator}, {Kind: RepoUseTask, ID: "k", RepoID: repo}, {Kind: RepoUseWorktreeOperation, ID: "o", RepoID: repo}, {Kind: RepoUseEpicCreate, ID: "c", RepoID: repo}}}
		}},
		"last repository": {project: "app", repo: "app", class: ErrorProjectInvalidArguments, reason: "at least one",
			setup: func(*testing.T, definitionFixture, *Dependencies) {}},
		"Epic binding": {class: ErrorProjectConflict, reason: "Epic e in", setup: func(_ *testing.T, _ definitionFixture, d *Dependencies) {
			d.RepoUses = registrationsStore{uses: []RepoUse{{Live: true, Kind: RepoUseEpic, ID: "e", RepoID: repo}}}
		}},
		"Task binding with empty common dir": {class: ErrorProjectConflict, reason: "Task k in", setup: func(_ *testing.T, _ definitionFixture, d *Dependencies) {
			d.RepoUses = registrationsStore{uses: []RepoUse{{Live: true, Kind: RepoUseTask, ID: "k", RepoID: repo, GitCommonDir: ""}}}
		}},
		"worktree operation": {class: ErrorProjectConflict, reason: "worktree operation o", setup: func(_ *testing.T, _ definitionFixture, d *Dependencies) {
			d.RepoUses = registrationsStore{uses: []RepoUse{{Live: true, Kind: RepoUseWorktreeOperation, ID: "o", RepoID: repo}}}
		}},
		"Epic create operation": {class: ErrorProjectConflict, reason: "Epic create operation o", setup: func(_ *testing.T, _ definitionFixture, d *Dependencies) {
			d.RepoUses = registrationsStore{uses: []RepoUse{{Live: true, Kind: RepoUseEpicCreate, ID: "o", RepoID: repo}}}
		}},
		"invalid work-item registry": {workItem: ErrorWorkStoreConflict, reason: "invalid work-item registry", setup: func(_ *testing.T, _ definitionFixture, d *Dependencies) {
			d.RepoUses = registrationsStore{err: WorkError(ErrorWorkStoreConflict, "invalid work-item registry", nil)}
		}},
		"repository path inside wrapper": {class: ErrorProjectConflict, reason: "repository nested of Project app is inside", setup: func(t *testing.T, f definitionFixture, _ *Dependencies) {
			f.editDefinitions(t, func(d *workspacedef.Definitions) {
				d.Projects[1].Repos = append(d.Projects[1].Repos, workspacedef.Repo{ID: "nested", Path: "ply/ply-code-map/nested"})
			})
		}},
		"Epic path inside wrapper": {class: ErrorProjectConflict, reason: "Epic e is inside", setup: func(_ *testing.T, f definitionFixture, d *Dependencies) {
			locator := filepath.Join(f.root, "ply", "ply-code-map", "epic")
			d.RepoUses = registrationsStore{uses: []RepoUse{{Live: true, Kind: RepoUseEpic, ID: "e", RepoID: "app", Worktree: locator}}}
		}},
		"Project root inside wrapper": {class: ErrorProjectConflict, reason: "Project inner has its root inside", setup: func(t *testing.T, f definitionFixture, _ *Dependencies) {
			f.editDefinitions(t, func(d *workspacedef.Definitions) {
				d.Projects = append(d.Projects, workspacedef.Project{ID: "inner", Name: "Inner", Wrapper: "ply/ply-code-map/inner", Repos: []workspacedef.Repo{{ID: "retired", Path: "retired"}}})
			})
		}},
		"Project root alias resolving inside wrapper": {class: ErrorProjectConflict, reason: "Project inner has its root inside", setup: func(t *testing.T, f definitionFixture, _ *Dependencies) {
			f.alias(t, "inner-alias", filepath.Join(f.root, "ply", "ply-code-map", "inner"))
			f.editDefinitions(t, func(d *workspacedef.Definitions) {
				d.Projects = append(d.Projects, workspacedef.Project{ID: "inner", Name: "Inner", Wrapper: "inner-alias", Repos: []workspacedef.Repo{{ID: "retired", Path: "retired"}}})
			})
		}},
		"repository alias resolving inside wrapper": {class: ErrorProjectConflict, reason: "repository nested of Project app is inside", setup: func(t *testing.T, f definitionFixture, _ *Dependencies) {
			f.alias(t, "nested-alias", filepath.Join(f.root, "ply", "ply-code-map", "nested"))
			f.editDefinitions(t, func(d *workspacedef.Definitions) {
				d.Projects[1].Repos = append(d.Projects[1].Repos, workspacedef.Repo{ID: "nested", Path: "nested-alias"})
			})
		}},
		"unresolvable alias fails closed": {class: ErrorProjectPath, reason: "dangling-alias", setup: func(t *testing.T, f definitionFixture, _ *Dependencies) {
			if err := os.Symlink(filepath.Join(f.root, "missing"), filepath.Join(f.root, "dangling-alias")); err != nil {
				t.Fatal(err)
			}
			f.editDefinitions(t, func(d *workspacedef.Definitions) {
				d.Projects[1].Repos = append(d.Projects[1].Repos, workspacedef.Repo{ID: "dangling", Path: "dangling-alias"})
			})
		}},
		"companion inside wrapper": {class: ErrorProjectConflict, reason: "of Project app is inside", setup: func(t *testing.T, f definitionFixture, _ *Dependencies) {
			docs := filepath.Join(f.root, "ply", "ply-code-map", "docs")
			if err := os.Mkdir(docs, 0o755); err != nil {
				t.Fatal(err)
			}
			input := ProjectMetadataInput{Operation: ProjectCompanionAdd, ProjectID: "app", Role: CompanionDocs, Locator: docs}
			if _, err := UpdateProjectMetadata(f.metadataDependencies(), input); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newRemovalFixture(t)
			d := f.metadataDependencies()
			tc.setup(t, f, &d)
			project, id := "ply", repo
			if tc.project != "" {
				project, id = tc.project, tc.repo
			}
			before := f.state(t)

			err := CheckRepoRemoval(d, f.root, project, string(id))

			f.requireUnchanged(t, before)
			var work *WorkItemError
			switch {
			case tc.class == "" && tc.workItem == "":
				if err != nil {
					t.Fatalf("clean check: %v", err)
				}
			case tc.workItem != "":
				if !errors.As(err, &work) || work.Class != tc.workItem {
					t.Fatalf("error %v, want class %s", err, tc.workItem)
				}
			case projectClass(err) != tc.class:
				t.Fatalf("error %v, want class %s", err, tc.class)
			}
			if err != nil && !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("error %v, want reason %q", err, tc.reason)
			}
			if name != "legacy registry lists repository" {
			}
		})
	}
}

func TestRemoveDefinitionRepoUsesExpectedBasis(t *testing.T) {
	t.Parallel()
	f := newRemovalFixture(t)
	d := f.metadataDependencies()
	_, basis, err := workspacedef.Read(f.root)
	if err != nil {
		t.Fatal(err)
	}
	before := f.state(t)

	_, err = RemoveDefinitionRepo(d, f.root, "ply", "ply-code-map", "ply/ply-code-map/main", "sha256:stale")

	var work *WorkItemError
	if !errors.As(err, &work) || work.Class != ErrorWorkIdentityConflict {
		t.Fatalf("stale basis error %v", err)
	}
	f.requireUnchanged(t, before)

	result, err := RemoveDefinitionRepo(d, f.root, "ply", "ply-code-map", "ply/ply-code-map/main", basis)

	if err != nil || !result.Changed {
		t.Fatalf("remove %+v %v", result, err)
	}
	for _, r := range f.definitions(t).Projects[0].Repos {
		if r.ID == "ply-code-map" {
			t.Fatalf("member kept: %+v", f.definitions(t).Projects[0].Repos)
		}
	}
}
