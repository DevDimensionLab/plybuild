package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/repohistory"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// DefinitionsFile lives beside .ply, independently of disposable runtime data.
const DefinitionsFile = workspacedef.File

// readDefinitions reads the definitions through workspacedef and applies the
// workspace Project, repository, data path and Epic rules.
func readDefinitions(root string) (*workspacedef.Definitions, string, error) {
	defs, digest, err := workspacedef.Read(root)
	var defErr *workspacedef.Error
	if errors.As(err, &defErr) {
		return nil, "", definitionError(defErr)
	}
	if err != nil || defs == nil {
		return nil, digest, err
	}
	if err = validateDefinitions(*defs); err != nil {
		return nil, "", err
	}
	return defs, digest, nil
}

type definitionRootFiles struct {
	FileSystem
	root string
}

func (f definitionRootFiles) Getwd() (string, error) { return f.root, nil }

// Config discovery can recreate only the workspace marker. Init retains its
// existing atomic publication, ancestor validation and concurrent-init guards.
func initDefinitionWorkspace(files FileSystem, root string) (bool, error) {
	defs, _, err := readDefinitions(root)
	if err != nil || defs == nil {
		return false, err
	}
	_, err = Init(Dependencies{Files: definitionRootFiles{files, root}})
	return err == nil, err
}

// Runtime Epic bindings are observations, not definitions. Bootstrap through
// native adoption, with fresh IDs and a current OID revalidated under its locks.
func ensureDefinitionEpics(d Dependencies, root string) error {
	defs, digest, err := readDefinitions(root)
	if err != nil || defs == nil {
		return err
	}
	if err = requireWorkMutationDependencies(d); err != nil {
		return err
	}
	registry, err := d.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		return err
	}
	projects, repos, err := d.Projects.Snapshot(root)
	if err != nil {
		return err
	}
	snapshot := ProjectSnapshot{Projects: projects, Repos: repos}
	defined := map[EpicID]bool{}
	for _, p := range defs.Projects {
		for _, e := range p.Epics {
			defined[EpicID(e.ID)] = true
		}
	}
	for _, e := range registry.Epics {
		if !defined[e.ID] {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("runtime Epic %s is absent from %s", e.ID, DefinitionsFile), nil)
		}
	}
	for _, p := range defs.Projects {
		for _, e := range p.Epics {
			epicID, projectID, repoID := EpicID(e.ID), ProjectID(p.ID), RepoID(e.Repo)
			if op, _ := findEpicCreate(registry, epicID); op != nil && op.State != "completed" {
				continue
			}
			_, repository, err := projectAndRepo(snapshot, projectID, repoID)
			if err != nil {
				return err
			}
			path, err := canonicalDirectory(d.Files, root, e.Path, "Epic worktree")
			existing, _ := findEpic(registry, epicID)
			if err != nil && existing != nil {
				path, err = archivedEpicLocator(d.Files, root, e.Path, epicID, err)
			}
			if err != nil {
				return err
			}
			if existing != nil {
				binding, _ := findEpicRepo(*existing, repoID)
				if existing.ProjectID != projectID || existing.Title != e.Title || len(existing.RepoBindings) != 1 || binding == nil || binding.Worktree.Locator != path || binding.Worktree.Ref != e.Ref || binding.GitCommonDir != repository.GitCommonDir {
					return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Epic %s differs from %s; reconcile runtime before changing a bound definition", e.ID, DefinitionsFile), nil)
				}
				continue
			}
			observed, err := d.WorkGit.ObserveWorktree(path)
			if err != nil {
				return err
			}
			input := EpicAdoptInput{EpicID: epicID, Title: e.Title, ProjectID: projectID, RepoID: repoID, Worktree: path, Ref: e.Ref, ExpectedOID: observed.OID}
			if _, err = adoptEpicAtRoot(d, input, root); err != nil {
				return err
			}
		}
	}
	_, after, err := readDefinitions(root)
	if err != nil {
		return err
	}
	if after != digest {
		return workError(ErrorWorkIdentityConflict, "definitions changed during Epic bootstrap; retry", nil)
	}
	return nil
}

func archivedEpicLocator(files FileSystem, root, input string, epicID EpicID, pathErr error) (string, error) {
	absolute := filepath.Clean(input)
	if !filepath.IsAbs(absolute) {
		absolute = filepath.Join(root, absolute)
	}
	if _, err := files.Lstat(absolute); !errors.Is(err, fs.ErrNotExist) {
		return "", pathErr
	}
	lifecycles, err := readWorkItemLifecycle(root)
	if err != nil {
		return "", err
	}
	if lifecycles.Epic(epicID) != LifecycleArchived {
		return "", pathErr
	}
	parent, err := files.EvalSymlinks(filepath.Dir(absolute))
	if err != nil {
		return "", pathErr
	}
	return filepath.Join(parent, filepath.Base(absolute)), nil
}

func checkDefinedAdoption(d Dependencies, root string, in EpicAdoptInput, locator string) error {
	defs, _, err := readDefinitions(root)
	if err != nil || defs == nil {
		return err
	}
	for _, p := range defs.Projects {
		for _, e := range p.Epics {
			if e.ID != string(in.EpicID) {
				continue
			}
			path, err := canonicalDirectory(d.Files, root, e.Path, "Epic worktree")
			if err != nil {
				return err
			}
			if p.ID == string(in.ProjectID) && e.Title == in.Title && e.Repo == string(in.RepoID) && path == locator && e.Ref == in.Ref {
				return nil
			}
		}
	}
	return WorkInvalidArguments("Epic adoption must match authoritative " + DefinitionsFile)
}

func registrationProjects(d Dependencies, root string) ([]ProjectRecord, []RepoRecord, error) {
	defs, _, err := readDefinitions(root)
	if err != nil {
		return nil, nil, err
	}
	if defs == nil {
		return d.Projects.Snapshot(root)
	}
	registry, err := d.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		return nil, nil, err
	}
	var bindings []repohistory.Binding
	for _, epic := range registry.Epics {
		for _, binding := range epic.RepoBindings {
			bindings = append(bindings, repohistory.Binding{RepoID: string(binding.RepoID), GitCommonDir: binding.GitCommonDir})
		}
	}
	for _, task := range registry.Tasks {
		bindings = append(bindings, repohistory.Binding{RepoID: string(task.RepoID), GitCommonDir: task.GitCommonDir})
	}
	identities, err := repohistory.CommonDirs(bindings)
	if err != nil {
		return nil, nil, WorkInvalidArguments("conflicting historical repository identity")
	}
	projects := []ProjectRecord{}
	repos := []RepoRecord{}
	seen := map[RepoID]bool{}
	path := func(value string) string {
		if filepath.IsAbs(value) {
			return filepath.Clean(value)
		}
		return filepath.Join(root, value)
	}
	for _, p := range defs.Projects {
		project := ProjectRecord{ID: ProjectID(p.ID), Name: p.Name, Wrapper: path(p.Wrapper), DataPaths: p.DataPaths}
		for _, r := range p.Repos {
			project.RepoIDs = append(project.RepoIDs, RepoID(r.ID))
			if !seen[RepoID(r.ID)] {
				repos = append(repos, RepoRecord{ID: RepoID(r.ID), Locator: path(r.Path), GitCommonDir: identities[r.ID]})
				seen[RepoID(r.ID)] = true
			}
		}
		projects = append(projects, project)
	}
	return projects, repos, nil
}
