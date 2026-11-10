package workspace

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/filelock"
	"gopkg.in/yaml.v3"
)

const (
	ProjectsFile = "projects.yaml"
	projectsLock = "projects.lock"
)

type ProjectStore interface {
	Add(workspaceRoot string, project ProjectRecord, repositories []RepoRecord) (ProjectRecord, []RepoRecord, bool, error)
	Snapshot(workspaceRoot string) ([]ProjectRecord, []RepoRecord, error)
}

type ProjectSnapshot struct {
	Projects []ProjectRecord
	Repos    []RepoRecord
}

type ProjectSnapshotLocker interface {
	WithSnapshotLock(string, func(ProjectSnapshot) error) error
}

type projectRegistry struct {
	FormatVersion int             `yaml:"format_version"`
	Projects      []ProjectRecord `yaml:"projects"`
	Repos         []RepoRecord    `yaml:"repos"`
}

type storeFaults struct {
	fail func(string) error
}

type systemProjectStore struct {
	faults      *storeFaults
	replaceFile func(string, []byte, fs.FileMode) error
}

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

func (store *systemProjectStore) Add(workspaceRoot string, project ProjectRecord, repositories []RepoRecord) (ProjectRecord, []RepoRecord, bool, error) {
	return store.mutate(workspaceRoot, func(registry projectRegistry) (projectRegistry, ProjectRecord, []RepoRecord, bool, error) {
		return applyProject(registry, project, repositories)
	})
}

func (store *systemProjectStore) AddMember(workspaceRoot string, id ProjectID, repository RepoRecord) (ProjectRecord, []RepoRecord, bool, error) {
	return store.mutate(workspaceRoot, func(registry projectRegistry) (projectRegistry, ProjectRecord, []RepoRecord, bool, error) {
		return applyProjectMember(registry, id, repository)
	})
}

// mutate serializes read, identity validation and atomic publication together.
func (store *systemProjectStore) mutate(workspaceRoot string, apply func(projectRegistry) (projectRegistry, ProjectRecord, []RepoRecord, bool, error)) (ProjectRecord, []RepoRecord, bool, error) {
	lock, err := store.acquire(workspaceRoot)
	if err != nil {
		return ProjectRecord{}, nil, false, err
	}

	var storedProject ProjectRecord
	var storedRepositories []RepoRecord
	var created bool
	operationErr := func() error {
		defs, _, err := readDefinitions(workspaceRoot)
		if err != nil {
			return err
		}
		if defs != nil {
			return ProjectInvalidArguments("edit " + DefinitionsFile + " to change authoritative Project/repository definitions")
		}
		registry, err := store.read(workspaceRoot)
		if err != nil {
			return err
		}
		updated, resultProject, resultRepositories, wasCreated, err := apply(registry)
		if err != nil {
			return err
		}
		storedProject = resultProject
		storedRepositories = resultRepositories
		created = wasCreated
		if !created {
			return nil
		}
		return store.publish(workspaceRoot, updated)
	}()
	releaseErr := store.release(lock)
	if operationErr != nil {
		if releaseErr != nil {
			return ProjectRecord{}, nil, false, projectIOError("release after failure", lock.Path(), fmt.Errorf("%v; %w", operationErr, releaseErr))
		}
		return ProjectRecord{}, nil, false, operationErr
	}
	if releaseErr != nil {
		return ProjectRecord{}, nil, false, releaseErr
	}
	return storedProject, storedRepositories, created, nil
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
	lock, err := filelock.Acquire(lockPath, filelock.Options{Fail: store.fail})
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
	defs, _, err := readDefinitions(workspaceRoot)
	if err != nil {
		return projectRegistry{}, err
	}
	if defs != nil {
		return definitionProjects(workspaceRoot, *defs)
	}
	registry, _, err := readLegacyProjectRegistry(workspaceRoot)
	return registry, err
}

// readLegacyProjectRegistry reads .ply/projects.yaml itself, whether or not
// definitions exist; exists is false when the file is absent.
func readLegacyProjectRegistry(workspaceRoot string) (projectRegistry, bool, error) {
	path := filepath.Join(workspaceRoot, MarkerDirectory, ProjectsFile)
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return projectRegistry{FormatVersion: FormatVersion, Projects: []ProjectRecord{}, Repos: []RepoRecord{}}, false, nil
	}
	if err != nil {
		return projectRegistry{}, false, projectStoreConflict(path, "cannot inspect project registry", err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return projectRegistry{}, false, projectStoreConflict(path, "project registry is not a regular file", nil)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return projectRegistry{}, false, projectStoreConflict(path, "cannot read project registry", err)
	}
	registry, err := decodeProjectRegistry(contents)
	if err != nil {
		return projectRegistry{}, false, projectStoreConflict(path, err.Error(), err)
	}
	return registry, true, nil
}

func (store *systemProjectStore) publish(workspaceRoot string, registry projectRegistry) error {
	contents, err := encodeProjectRegistry(registry)
	if err != nil {
		return projectIOError("marshal registry", filepath.Join(workspaceRoot, MarkerDirectory, ProjectsFile), err)
	}
	destination := filepath.Join(workspaceRoot, MarkerDirectory, ProjectsFile)
	replace := store.replaceFile
	if replace == nil {
		replace = atomicfile.Replace
	}
	if err := replace(destination, contents, 0o644); err != nil {
		var failure *atomicfile.Error
		if errors.As(err, &failure) {
			return projectError(ErrorProjectIO, failure.Error(), err)
		}
		return projectIOError("replace", destination, err)
	}
	return nil
}

func (store *systemProjectStore) fail(operation string) error {
	if store.faults == nil || store.faults.fail == nil {
		return nil
	}
	return store.faults.fail(operation)
}

func applyProject(registry projectRegistry, project ProjectRecord, repositories []RepoRecord) (projectRegistry, ProjectRecord, []RepoRecord, bool, error) {
	existingProjects := make(map[ProjectID]ProjectRecord, len(registry.Projects))
	projectByWrapper := make(map[string]ProjectRecord, len(registry.Projects))
	for _, existing := range registry.Projects {
		existingProjects[existing.ID] = existing
		projectByWrapper[existing.Wrapper] = existing
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
	if owner, found := projectByWrapper[project.Wrapper]; found {
		return registry, ProjectRecord{}, nil, false, projectError(ErrorProjectConflict, fmt.Sprintf("wrapper %s is already registered to project %s", project.Wrapper, owner.ID), nil)
	}

	updated := projectRegistry{FormatVersion: FormatVersion, Projects: cloneProjects(registry.Projects), Repos: cloneRepos(registry.Repos)}
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

func encodeProjectRegistry(registry projectRegistry) ([]byte, error) {
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(registry); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func decodeProjectRegistry(contents []byte) (projectRegistry, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		return projectRegistry{}, errors.New("project registry is invalid")
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		return projectRegistry{}, errors.New("project registry must contain exactly one YAML document")
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return projectRegistry{}, errors.New("project registry root must be a mapping")
	}
	root := document.Content[0]
	fields, err := exactMapping(root, map[string]yaml.Kind{
		"format_version": yaml.ScalarNode,
		"projects":       yaml.SequenceNode,
		"repos":          yaml.SequenceNode,
	})
	if err != nil {
		return projectRegistry{}, err
	}
	if fields["format_version"].Tag != "!!int" || fields["format_version"].Value != "1" {
		return projectRegistry{}, errors.New("unsupported format_version")
	}
	registry := projectRegistry{FormatVersion: FormatVersion, Projects: []ProjectRecord{}, Repos: []RepoRecord{}}
	for _, node := range fields["projects"].Content {
		project, err := decodeProjectRecord(node)
		if err != nil {
			return projectRegistry{}, err
		}
		registry.Projects = append(registry.Projects, project)
	}
	for _, node := range fields["repos"].Content {
		repository, err := decodeRepoRecord(node)
		if err != nil {
			return projectRegistry{}, err
		}
		registry.Repos = append(registry.Repos, repository)
	}
	if err := validateRegistry(registry); err != nil {
		return projectRegistry{}, err
	}
	return registry, nil
}

func decodeProjectRecord(node *yaml.Node) (ProjectRecord, error) {
	fields, err := exactMapping(node, map[string]yaml.Kind{
		"id": yaml.ScalarNode, "name": yaml.ScalarNode, "wrapper": yaml.ScalarNode, "repo_ids": yaml.SequenceNode,
	})
	if err != nil {
		return ProjectRecord{}, fmt.Errorf("invalid project record: %w", err)
	}
	for _, key := range []string{"id", "name", "wrapper"} {
		if fields[key].Tag != "!!str" {
			return ProjectRecord{}, fmt.Errorf("project %s must be a string", key)
		}
	}
	project := ProjectRecord{ID: ProjectID(fields["id"].Value), Name: fields["name"].Value, Wrapper: fields["wrapper"].Value, RepoIDs: []RepoID{}}
	for _, node := range fields["repo_ids"].Content {
		if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
			return ProjectRecord{}, errors.New("project repo_ids must contain strings")
		}
		project.RepoIDs = append(project.RepoIDs, RepoID(node.Value))
	}
	return project, nil
}

func decodeRepoRecord(node *yaml.Node) (RepoRecord, error) {
	fields, err := exactMapping(node, map[string]yaml.Kind{
		"id": yaml.ScalarNode, "locator": yaml.ScalarNode, "git_common_dir": yaml.ScalarNode,
	})
	if err != nil {
		return RepoRecord{}, fmt.Errorf("invalid repository record: %w", err)
	}
	for _, key := range []string{"id", "locator", "git_common_dir"} {
		if fields[key].Tag != "!!str" {
			return RepoRecord{}, fmt.Errorf("repository %s must be a string", key)
		}
	}
	return RepoRecord{ID: RepoID(fields["id"].Value), Locator: fields["locator"].Value, GitCommonDir: fields["git_common_dir"].Value}, nil
}

func exactMapping(node *yaml.Node, expected map[string]yaml.Kind) (map[string]*yaml.Node, error) {
	if node.Kind != yaml.MappingNode {
		return nil, errors.New("value must be a mapping")
	}
	if len(node.Content) != len(expected)*2 {
		return nil, errors.New("mapping has missing or unknown fields")
	}
	fields := make(map[string]*yaml.Node, len(expected))
	for index := 0; index < len(node.Content); index += 2 {
		key := node.Content[index]
		value := node.Content[index+1]
		kind, known := expected[key.Value]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" || !known {
			return nil, errors.New("mapping has an unknown field")
		}
		if _, duplicate := fields[key.Value]; duplicate {
			return nil, errors.New("mapping has a duplicate field")
		}
		if value.Kind != kind {
			return nil, fmt.Errorf("field %s has the wrong type", key.Value)
		}
		fields[key.Value] = value
	}
	return fields, nil
}

func validateRegistry(registry projectRegistry) error {
	projects := make(map[ProjectID]ProjectRecord, len(registry.Projects))
	wrapperOwners := make(map[string]ProjectID, len(registry.Projects))
	lastProject := ProjectID("")
	for index, project := range registry.Projects {
		if err := validateIdentifier("project", string(project.ID)); err != nil {
			return errors.New("project registry contains an invalid project ID")
		}
		if index > 0 && project.ID <= lastProject {
			return errors.New("project records must be uniquely sorted by ID")
		}
		lastProject = project.ID
		if err := validateName(project.Name); err != nil {
			return errors.New("project registry contains an invalid project name")
		}
		if !validStoredPath(project.Wrapper) {
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
			if err := validateIdentifier("repository", string(id)); err != nil {
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
		if err := validateIdentifier("repository", string(repository.ID)); err != nil {
			return errors.New("project registry contains an invalid repository ID")
		}
		if index > 0 && repository.ID <= lastRepo {
			return errors.New("repository records must be uniquely sorted by ID")
		}
		lastRepo = repository.ID
		if !validStoredPath(repository.Locator) || !validStoredPath(repository.GitCommonDir) {
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
