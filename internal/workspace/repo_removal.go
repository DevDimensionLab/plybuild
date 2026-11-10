package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

// RemoveDefinitionRepo removes the member {repo, path} of project from
// ply-workspace.json under projects.lock. expectedBasis is the digest the
// caller previewed; it is never replaced by a fresh read. The current and the
// changed definitions are validated like AddDefinitionRepo validates them.
func RemoveDefinitionRepo(d Dependencies, root, project, repo, path, expectedBasis string) (DefinitionRepoResult, error) {
	projectID, err := ParseProjectID(project)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	if err := validateIdentifier("repository", repo); err != nil {
		return DefinitionRepoResult{}, err
	}
	observed, err := ObserveRoot(d, root)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	mutate := func(defs *workspacedef.Definitions) (bool, error) {
		return workspacedef.RemoveRepo(defs, project, repo, path)
	}
	check := func(defs workspacedef.Definitions) error {
		if err := validateDefinitions(defs); err != nil {
			return err
		}
		_, err := definitionProjects(observed.Root, defs)
		return err
	}
	result, err := workspacedef.Update(observed.Root, expectedBasis, mutate, check)
	if err != nil {
		return DefinitionRepoResult{}, definitionUpdateError(err)
	}
	return DefinitionRepoResult{Workspace: observed.Root, ProjectID: projectID, RepoID: RepoID(repo), Path: path, Changed: result.Changed}, nil
}

// CheckRepoRemoval reports, without lock or write, whether repo can leave
// project: the definitions stay valid without it, and nothing else names it or
// lies inside its wrapper <Project root>/<repo>. References are
// ErrorProjectConflict; unreadable sources keep their own classes.
func CheckRepoRemoval(d Dependencies, root, project, repo string) error {
	if _, err := ParseProjectID(project); err != nil {
		return err
	}
	if err := validateIdentifier("repository", repo); err != nil {
		return err
	}
	if d.WorkItems == nil {
		return projectError(ErrorProjectIO, "work-item store dependency is required", nil)
	}
	observed, err := ObserveRoot(d, root)
	if err != nil {
		return err
	}
	root = observed.Root
	defs, _, err := workspacedef.Read(root)
	if err != nil {
		return definitionUpdateError(err)
	}
	if defs == nil {
		return definitionUpdateError(&workspacedef.Error{Kind: workspacedef.KindMissing, Path: filepath.Join(root, workspacedef.File), Detail: "definitions do not exist"})
	}
	referenced := func(format string, args ...interface{}) error {
		return projectError(ErrorProjectConflict, fmt.Sprintf("repository %s cannot be removed: ", repo)+fmt.Sprintf(format, args...), nil)
	}
	resolve := func(path string) string {
		if filepath.IsAbs(path) {
			return filepath.Clean(path)
		}
		return filepath.Join(root, filepath.FromSlash(path))
	}
	wrapper, member := "", ""
	for _, p := range defs.Projects {
		if p.ID == project {
			wrapper = filepath.Join(resolve(p.Wrapper), repo)
			for _, r := range p.Repos {
				if r.ID == repo {
					member = r.Path
				}
			}
		}
		for _, e := range p.Epics {
			if e.Repo == repo {
				return referenced("Epic %s in %s names it", e.ID, workspacedef.File)
			}
		}
	}
	if wrapper == "" {
		return ProjectInvalidArguments(fmt.Sprintf("Project %s is not defined", project))
	}
	candidate := *defs
	candidate.Projects = append([]workspacedef.Project{}, defs.Projects...)
	if _, err := workspacedef.RemoveRepo(&candidate, project, repo, member); err != nil {
		return definitionUpdateError(err)
	}
	if err := validateDefinitions(candidate); err != nil {
		return err
	}
	physicalWrapper, err := physicalLocation(d.Files, wrapper)
	if err != nil {
		return err
	}
	within := func(parent, path string) bool {
		rel, err := filepath.Rel(parent, path)
		return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
	}
	// inside compares the stored and the physical location, so that an alias
	// outside the wrapper that resolves into it is a reference too.
	inside := func(path string) (bool, error) {
		if within(wrapper, path) {
			return true, nil
		}
		physical, err := physicalLocation(d.Files, path)
		return err == nil && within(physicalWrapper, physical), err
	}
	for _, p := range defs.Projects {
		if found, err := inside(resolve(p.Wrapper)); err != nil || found {
			return firstErr(err, referenced("Project %s has its root inside wrapper %s", p.ID, wrapper))
		}
		for _, r := range p.Repos {
			if p.ID == project && r.ID == repo {
				continue
			}
			if found, err := inside(resolve(r.Path)); err != nil || found {
				return firstErr(err, referenced("repository %s of Project %s is inside wrapper %s", r.ID, p.ID, wrapper))
			}
		}
		for _, e := range p.Epics {
			if found, err := inside(resolve(e.Path)); err != nil || found {
				return firstErr(err, referenced("Epic %s is inside wrapper %s", e.ID, wrapper))
			}
		}
	}
	metadata, err := readProjectMetadata(root)
	if err != nil {
		return err
	}
	for id, companions := range metadata.companions {
		for _, companion := range companions {
			if found, err := inside(companion.Locator); err != nil || found {
				return firstErr(err, referenced("companion %s of Project %s is inside wrapper %s", companion.Locator, id, wrapper))
			}
		}
	}
	registry, err := d.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		return err
	}
	for _, e := range registry.Epics {
		for _, b := range e.RepoBindings {
			if string(b.RepoID) == repo {
				return referenced("Epic %s in %s binds it", e.ID, workItemsPath(root))
			}
		}
	}
	for _, task := range registry.Tasks {
		if string(task.RepoID) == repo {
			return referenced("Task %s in %s binds it", task.ID, workItemsPath(root))
		}
	}
	for _, operation := range registry.WorktreeOperations {
		if string(operation.RepoID) == repo {
			return referenced("worktree operation %s in %s names it", operation.ID, workItemsPath(root))
		}
	}
	for _, operation := range registry.EpicCreateOperations {
		if string(operation.Request.RepoID) == repo {
			return referenced("Epic create operation %s in %s names it", operation.ID, workItemsPath(root))
		}
	}
	legacy, exists, err := readLegacyProjectRegistry(root)
	if err != nil {
		return err
	}
	for _, r := range legacy.Repos {
		if exists && string(r.ID) == repo {
			return referenced("the legacy registry %s lists it", filepath.Join(root, MarkerDirectory, ProjectsFile))
		}
	}
	return nil
}

func firstErr(err, otherwise error) error {
	if err != nil {
		return err
	}
	return otherwise
}

// physicalLocation resolves the absolute path like a metadata removal path: an
// existing path must resolve, and a missing one keeps its physical ancestor.
func physicalLocation(files FileSystem, path string) (string, error) {
	if _, err := files.Lstat(path); err == nil {
		physical, err := files.EvalSymlinks(path)
		if err != nil {
			return "", projectPathError("resolve reference", path, err)
		}
		return physical, nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", projectPathError("inspect reference", path, err)
	}
	return projectMetadataPath(files, "", path, true)
}

// ClearRepositoryWrapperAt clears the wrapper metadata of repo in the
// workspace at root with the expected revision. Unlike UpdateProjectMetadata
// it observes the root only and never bootstraps runtime Epics.
func ClearRepositoryWrapperAt(d Dependencies, root, project, repo, actor string, expectedRevision int) (ProjectMetadataMutation, error) {
	in := ProjectMetadataInput{Operation: ProjectRepoWrapperClear, ProjectID: ProjectID(project), RepoID: RepoID(repo), ActorClaim: actor, ExpectedRevision: &expectedRevision}
	if err := in.Validate(); err != nil {
		return ProjectMetadataMutation{}, err
	}
	if d.Projects == nil || d.WorkClock == nil {
		return ProjectMetadataMutation{}, projectError(ErrorProjectIO, "project store and clock dependencies are required", nil)
	}
	observed, err := ObserveRoot(d, root)
	if err != nil {
		return ProjectMetadataMutation{}, err
	}
	// A clear has no locator, so the working directory is never consulted.
	return applyProjectMetadata(d, observed.Root, observed.Root, in, workspaceMetadataPublisher{})
}
