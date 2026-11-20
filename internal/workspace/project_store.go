package workspace

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/filelock"
)

const projectsLock = "projects.lock"

// ProjectStore reads the Projects defined in ply-workspace.json.
type ProjectStore interface {
	Snapshot(workspaceRoot string) ([]ProjectRecord, []RepoRecord, error)
}

type ProjectSnapshot struct {
	Projects []ProjectRecord
	Repos    []RepoRecord
}

type ProjectSnapshotLocker interface {
	WithSnapshotLock(string, func(ProjectSnapshot) error) error
}

// projectRegistry is the resolved form of the definitions.
type projectRegistry struct {
	Projects []ProjectRecord
	Repos    []RepoRecord
}

type systemProjectStore struct{}

func newSystemProjectStore() *systemProjectStore { return &systemProjectStore{} }

func (store *systemProjectStore) WithSnapshotLock(workspaceRoot string, operation func(ProjectSnapshot) error) error {
	lock, err := store.acquire(workspaceRoot)
	if err != nil {
		return err
	}
	operationErr := func() error {
		registry, err := store.read(workspaceRoot)
		if err != nil {
			return err
		}
		return operation(ProjectSnapshot{Projects: cloneProjects(registry.Projects), Repos: cloneRepos(registry.Repos)})
	}()
	releaseErr := store.release(lock)
	if operationErr != nil {
		if releaseErr != nil {
			return projectIOError("release after failure", lock.Path(), fmt.Errorf("%v; %w", operationErr, releaseErr))
		}
		return operationErr
	}
	return releaseErr
}

func (store *systemProjectStore) Snapshot(workspaceRoot string) ([]ProjectRecord, []RepoRecord, error) {
	registry, err := store.read(workspaceRoot)
	if err != nil {
		return nil, nil, err
	}
	return cloneProjects(registry.Projects), cloneRepos(registry.Repos), nil
}

func (store *systemProjectStore) acquire(workspaceRoot string) (*filelock.Lock, error) {
	lockPath := filepath.Join(workspaceRoot, MarkerDirectory, projectsLock)
	lock, err := filelock.Acquire(lockPath, filelock.Options{})
	var lockErr *filelock.Error
	if errors.As(err, &lockErr) {
		return nil, projectLockError(lockPath, lockErr)
	}
	return lock, err
}

func projectLockError(path string, err *filelock.Error) error {
	switch err.Op {
	case "inspect":
		return projectStoreConflict(path, err.Err.Error(), err.Err)
	case "not-regular":
		return projectStoreConflict(path, "project lock is not a regular file", nil)
	case "stat-opened":
		return projectStoreConflict(path, "cannot validate opened project lock", err.Err)
	case "lstat-path":
		return projectStoreConflict(path, "cannot validate project lock path", err.Err)
	case "changed":
		return projectStoreConflict(path, "project lock changed or is not a regular file", nil)
	case "open":
		return projectIOError("open lock", path, err.Err)
	default:
		return projectIOError("lock", path, err.Err)
	}
}

func (store *systemProjectStore) release(lock *filelock.Lock) error {
	if err := lock.Release(); err != nil {
		return projectReleaseError(err)
	}
	return nil
}

func projectReleaseError(err *filelock.ReleaseError) error {
	if err.Unlock != nil {
		unlockErr := err.Unlock
		if err.Close != nil {
			unlockErr = fmt.Errorf("%w; close: %v", unlockErr, err.Close)
		}
		return projectIOError("unlock", err.Path, unlockErr)
	}
	return projectIOError("close lock", err.Path, err.Close)
}

func (store *systemProjectStore) read(workspaceRoot string) (projectRegistry, error) {
	defs, _, err := ReadDefinitions(workspaceRoot)
	if err != nil {
		return projectRegistry{}, err
	}
	if defs == nil {
		return projectRegistry{}, MissingDefinitions(workspaceRoot)
	}
	return definitionProjects(workspaceRoot, *defs)
}

func applyProject(registry projectRegistry, project ProjectRecord, repositories []RepoRecord) (projectRegistry, ProjectRecord, []RepoRecord, bool, error) {
	existingProjects := make(map[ProjectID]ProjectRecord, len(registry.Projects))
	for _, existing := range registry.Projects {
		existingProjects[existing.ID] = existing
	}
	repoByID := make(map[RepoID]RepoRecord, len(registry.Repos))
	repoByCommon := make(map[string]RepoRecord, len(registry.Repos))
	for _, existing := range registry.Repos {
		repoByID[existing.ID] = existing
		repoByCommon[existing.GitCommonDir] = existing
	}

	if existing, found := existingProjects[project.ID]; found {
		field := projectDifference(existing, project, repoByID, repositories)
		if field != "" {
			return registry, ProjectRecord{}, nil, false, projectError(ErrorProjectConflict, fmt.Sprintf("project %s is already registered with different %s", project.ID, field), nil)
		}
		return registry, existing, repositoriesForProject(existing, registry.Repos), false, nil
	}
	for _, existing := range registry.Projects {
		if err := wrapperConflict(project.Wrapper, existing.Wrapper, existing.ID); err != nil {
			return registry, ProjectRecord{}, nil, false, err
		}
	}

	updated := projectRegistry{Projects: cloneProjects(registry.Projects), Repos: cloneRepos(registry.Repos)}
	for _, repository := range repositories {
		if existing, found := repoByID[repository.ID]; found {
			if existing.GitCommonDir != repository.GitCommonDir {
				return registry, ProjectRecord{}, nil, false, projectError(ErrorProjectConflict, fmt.Sprintf("repository %s is already registered with different Git common directory", repository.ID), nil)
			}
			continue
		}
		if existing, found := repoByCommon[repository.GitCommonDir]; found {
			return registry, ProjectRecord{}, nil, false, projectError(ErrorProjectConflict, fmt.Sprintf("Git common directory %s is already registered as repository %s", repository.GitCommonDir, existing.ID), nil)
		}
		updated.Repos = append(updated.Repos, repository)
		repoByID[repository.ID] = repository
		repoByCommon[repository.GitCommonDir] = repository
	}
	updated.Projects = append(updated.Projects, project)
	sort.Slice(updated.Projects, func(i, j int) bool { return updated.Projects[i].ID < updated.Projects[j].ID })
	sort.Slice(updated.Repos, func(i, j int) bool { return updated.Repos[i].ID < updated.Repos[j].ID })
	return updated, project, repositoriesForProject(project, updated.Repos), true, nil
}

func projectDifference(existing, candidate ProjectRecord, repositoryByID map[RepoID]RepoRecord, candidates []RepoRecord) string {
	if existing.Name != candidate.Name {
		return "name"
	}
	if existing.Wrapper != candidate.Wrapper {
		return "wrapper"
	}
	if len(existing.RepoIDs) != len(candidate.RepoIDs) {
		return "repository members"
	}
	for index, id := range existing.RepoIDs {
		if candidate.RepoIDs[index] != id {
			return "repository members"
		}
	}
	candidateByID := make(map[RepoID]RepoRecord, len(candidates))
	for _, repository := range candidates {
		candidateByID[repository.ID] = repository
	}
	for _, id := range existing.RepoIDs {
		if repositoryByID[id].GitCommonDir != candidateByID[id].GitCommonDir {
			return fmt.Sprintf("repository %s", id)
		}
	}
	return ""
}

func validateRegistry(registry projectRegistry) error {
	projects := make(map[ProjectID]ProjectRecord, len(registry.Projects))
	wrapperOwners := make(map[string]ProjectID, len(registry.Projects))
	lastProject := ProjectID("")
	for index, project := range registry.Projects {
		if err := ValidateIdentifier("project", string(project.ID)); err != nil {
			return errors.New("project registry contains an invalid project ID")
		}
		if index > 0 && project.ID <= lastProject {
			return errors.New("project records must be uniquely sorted by ID")
		}
		lastProject = project.ID
		if err := validateName(project.Name); err != nil {
			return errors.New("project registry contains an invalid project name")
		}
		if !ValidStoredPath(project.Wrapper) {
			return errors.New("project registry contains an invalid wrapper path")
		}
		if owner, duplicate := wrapperOwners[project.Wrapper]; duplicate {
			return fmt.Errorf("project wrapper is duplicated by %s", owner)
		}
		wrapperOwners[project.Wrapper] = project.ID
		if len(project.RepoIDs) == 0 {
			return errors.New("project must reference at least one repository")
		}
		lastRepo := RepoID("")
		for repoIndex, id := range project.RepoIDs {
			if err := ValidateIdentifier("repository", string(id)); err != nil {
				return errors.New("project contains an invalid repository ID")
			}
			if repoIndex > 0 && id <= lastRepo {
				return errors.New("project repository IDs must be uniquely sorted")
			}
			lastRepo = id
		}
		projects[project.ID] = project
	}

	repositories := make(map[RepoID]RepoRecord, len(registry.Repos))
	commonDirectories := make(map[string]RepoID, len(registry.Repos))
	lastRepo := RepoID("")
	for index, repository := range registry.Repos {
		if err := ValidateIdentifier("repository", string(repository.ID)); err != nil {
			return errors.New("project registry contains an invalid repository ID")
		}
		if index > 0 && repository.ID <= lastRepo {
			return errors.New("repository records must be uniquely sorted by ID")
		}
		lastRepo = repository.ID
		if !ValidStoredPath(repository.Locator) || !ValidStoredPath(repository.GitCommonDir) {
			return errors.New("project registry contains an invalid repository path")
		}
		if owner, duplicate := commonDirectories[repository.GitCommonDir]; duplicate {
			return fmt.Errorf("Git common directory is duplicated by repository %s", owner)
		}
		commonDirectories[repository.GitCommonDir] = repository.ID
		repositories[repository.ID] = repository
	}
	referenced := make(map[RepoID]bool, len(repositories))
	for _, project := range registry.Projects {
		for _, id := range project.RepoIDs {
			if _, found := repositories[id]; !found {
				return fmt.Errorf("project %s references missing repository %s", project.ID, id)
			}
			referenced[id] = true
		}
	}
	for id := range repositories {
		if !referenced[id] {
			return fmt.Errorf("repository %s is not referenced by a project", id)
		}
	}
	return nil
}

func cloneProjects(input []ProjectRecord) []ProjectRecord {
	result := make([]ProjectRecord, len(input))
	for index, project := range input {
		result[index] = project
		result[index].RepoIDs = append([]RepoID(nil), project.RepoIDs...)
		result[index].Companions = append([]ProjectCompanion(nil), project.Companions...)
	}
	return result
}

func cloneRepos(input []RepoRecord) []RepoRecord { return append([]RepoRecord(nil), input...) }

func projectStoreConflict(path, detail string, err error) error {
	return projectError(ErrorProjectStoreConflict, fmt.Sprintf("%s: %s", path, detail), err)
}

func projectIOError(operation, path string, err error) error {
	detail := strings.TrimSpace(fmt.Sprintf("%s %s: %v", operation, path, err))
	return projectError(ErrorProjectIO, detail, err)
}
