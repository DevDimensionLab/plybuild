package workspace

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ProjectMemberStore appends membership under the same lock as Project creation.
type ProjectMemberStore interface {
	AddMember(string, ProjectID, RepoRecord) (ProjectRecord, []RepoRecord, bool, error)
}

func ParseProjectMemberInput(projectID, repository string) (ProjectID, RepoInput, error) {
	id, err := ParseProjectID(projectID)
	if err != nil {
		return "", RepoInput{}, err
	}
	repoID, path, found := strings.Cut(repository, "=")
	if !found {
		return "", RepoInput{}, ProjectInvalidArguments("repository member must be <repo-id>=<path>")
	}
	repo := RepoInput{RepoID: RepoID(repoID), Path: path}
	if err := validateProjectMemberInput(id, repo); err != nil {
		return "", RepoInput{}, err
	}
	return id, repo, nil
}

// AddProjectMember registers metadata only; repository source and execution bindings are untouched.
// Created reports whether membership was added rather than already present.
func AddProjectMember(dependencies Dependencies, id ProjectID, repository RepoInput) (ProjectResult, error) {
	if err := validateProjectMemberInput(id, repository); err != nil {
		return ProjectResult{}, err
	}
	store, ok := dependencies.Projects.(ProjectMemberStore)
	if !ok || dependencies.Files == nil || dependencies.Repos == nil {
		return ProjectResult{}, projectError(ErrorProjectIO, "dependencies: filesystem, repository observer, and member store are required", nil)
	}
	cwd, err := projectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return ProjectResult{}, err
	}
	root, err := locateContainingWorkspace(dependencies.Files, cwd)
	if err != nil {
		return ProjectResult{}, err
	}
	locator, err := canonicalDirectory(dependencies.Files, cwd, repository.Path, fmt.Sprintf("repository %s", repository.RepoID))
	if err != nil {
		return ProjectResult{}, err
	}
	observation, err := dependencies.Repos.Observe(locator)
	if err != nil {
		var typed *ProjectError
		if errors.As(err, &typed) {
			return ProjectResult{}, err
		}
		return ProjectResult{}, projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe repository %s at %s: %v", repository.RepoID, locator, err), err)
	}
	if observation.Locator != locator || !validStoredPath(observation.Locator) || !validStoredPath(observation.GitCommonDir) {
		return ProjectResult{}, projectError(ErrorProjectRepoObservation, "observer returned malformed physical paths", nil)
	}
	project, repos, added, err := store.AddMember(root, id, RepoRecord{ID: repository.RepoID, Locator: observation.Locator, GitCommonDir: observation.GitCommonDir})
	if err != nil {
		return ProjectResult{}, err
	}
	return ProjectResult{Workspace: root, Project: project, Repos: repos, Created: added}, nil
}

func applyProjectMember(registry projectRegistry, id ProjectID, repository RepoRecord) (projectRegistry, ProjectRecord, []RepoRecord, bool, error) {
	index := -1
	for i, project := range registry.Projects {
		if project.ID == id {
			index = i
			break
		}
	}
	if index < 0 {
		return registry, ProjectRecord{}, nil, false, projectError(ErrorProjectNotFound, fmt.Sprintf("project %s is not registered", id), nil)
	}
	known := false
	for _, existing := range registry.Repos {
		if existing.ID == repository.ID {
			if existing.GitCommonDir != repository.GitCommonDir {
				return registry, ProjectRecord{}, nil, false, projectError(ErrorProjectConflict, fmt.Sprintf("repository %s is already registered with different Git common directory %s (requested %s)", repository.ID, existing.GitCommonDir, repository.GitCommonDir), nil)
			}
			known = true
		} else if existing.GitCommonDir == repository.GitCommonDir {
			return registry, ProjectRecord{}, nil, false, projectError(ErrorProjectConflict, fmt.Sprintf("Git common directory %s is already registered as repository %s", repository.GitCommonDir, existing.ID), nil)
		}
	}
	project := registry.Projects[index]
	for _, member := range project.RepoIDs {
		if member == repository.ID {
			return registry, project, repositoriesForProject(project, registry.Repos), false, nil
		}
	}
	updated := projectRegistry{FormatVersion: registry.FormatVersion, Projects: cloneProjects(registry.Projects), Repos: cloneRepos(registry.Repos)}
	project = updated.Projects[index]
	project.RepoIDs = append(project.RepoIDs, repository.ID)
	sort.Slice(project.RepoIDs, func(i, j int) bool { return project.RepoIDs[i] < project.RepoIDs[j] })
	updated.Projects[index] = project
	if !known {
		updated.Repos = append(updated.Repos, repository)
		sort.Slice(updated.Repos, func(i, j int) bool { return updated.Repos[i].ID < updated.Repos[j].ID })
	}
	return updated, project, repositoriesForProject(project, updated.Repos), true, nil
}

func validateProjectMemberInput(id ProjectID, repository RepoInput) error {
	if err := validateIdentifier("project", string(id)); err != nil {
		return err
	}
	if err := validateIdentifier("repository", string(repository.RepoID)); err != nil {
		return err
	}
	if repository.Path == "" {
		return ProjectInvalidArguments(fmt.Sprintf("repository %s path is empty", repository.RepoID))
	}
	return nil
}
