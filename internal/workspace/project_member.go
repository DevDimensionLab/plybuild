package workspace

import (
	"fmt"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

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

// AddProjectMember adds a repository member, with the URL of its origin, to a
// Project defined in ply-workspace.json; repository source is untouched.
// Created reports whether the definitions changed.
func AddProjectMember(dependencies Dependencies, id ProjectID, repository RepoInput) (ProjectResult, error) {
	if err := validateProjectMemberInput(id, repository); err != nil {
		return ProjectResult{}, err
	}
	remotes, ok := dependencies.Repos.(RemoteObserver)
	if !ok || dependencies.Files == nil {
		return ProjectResult{}, projectError(ErrorProjectIO, "dependencies: filesystem and a repository observer that reads remotes are required", nil)
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
	remote, err := remotes.Remote(locator)
	if err != nil {
		return ProjectResult{}, err
	}
	added, err := updateDefinitions(root, func(defs *workspacedef.Definitions) (bool, error) {
		return workspacedef.AddRepo(defs, string(id), string(repository.RepoID), definitionPath(root, locator), remote)
	})
	if err != nil {
		return ProjectResult{}, err
	}
	project, repos, err := definedProject(root, id)
	if err != nil {
		return ProjectResult{}, err
	}
	return ProjectResult{Workspace: root, Project: project, Repos: repos, Created: added}, nil
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
