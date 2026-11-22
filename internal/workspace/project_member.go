package workspace

import (
	"fmt"
	"path/filepath"
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
	cwd, err := ProjectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return ProjectResult{}, err
	}
	root, err := LocateContainingWorkspace(dependencies.Files, cwd)
	if err != nil {
		return ProjectResult{}, err
	}
	locator, err := CanonicalDirectory(dependencies.Files, cwd, repository.Path, fmt.Sprintf("repository %s", repository.RepoID))
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

// MoveProjectMember gives the member repo of Project id the ID and path of
// to, with the URL of the origin at that path; the member keeps its wrapper,
// which must contain the new path. Files and Git are untouched. The move is
// refused while the Project's registry binds the member. Created reports
// whether the definitions changed.
func MoveProjectMember(dependencies Dependencies, id ProjectID, repo RepoID, to RepoInput) (ProjectResult, error) {
	if err := validateProjectMemberInput(id, to); err != nil {
		return ProjectResult{}, err
	}
	if err := ValidateIdentifier("repository", string(repo)); err != nil {
		return ProjectResult{}, err
	}
	remotes, ok := dependencies.Repos.(RemoteObserver)
	if !ok || dependencies.Files == nil {
		return ProjectResult{}, projectError(ErrorProjectIO, "dependencies: filesystem and a repository observer that reads remotes are required", nil)
	}
	cwd, err := ProjectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return ProjectResult{}, err
	}
	root, err := LocateContainingWorkspace(dependencies.Files, cwd)
	if err != nil {
		return ProjectResult{}, err
	}
	locator, err := CanonicalDirectory(dependencies.Files, cwd, to.Path, fmt.Sprintf("repository %s", to.RepoID))
	if err != nil {
		return ProjectResult{}, err
	}
	remote, err := remotes.Remote(locator)
	if err != nil {
		return ProjectResult{}, err
	}
	if dependencies.RepoUses == nil {
		return ProjectResult{}, projectError(ErrorProjectIO, "work-item store dependency is required", nil)
	}
	work, err := DefinedProjectWork(root, id)
	if err != nil {
		return ProjectResult{}, err
	}
	stored := definitionPath(root, locator)
	// The current definitions name the old path, which the move has usually
	// emptied, so only definitions with the moved member are checked against Git.
	check := func(defs workspacedef.Definitions) error {
		if err := validateDefinitions(defs); err != nil {
			return err
		}
		for _, p := range defs.Projects {
			for _, r := range p.Repos {
				if p.ID == string(id) && r.ID == string(to.RepoID) && r.Path == stored {
					_, err := definitionProjects(root, defs)
					return err
				}
			}
		}
		return nil
	}
	result, err := workspacedef.Apply(root, func(defs *workspacedef.Definitions) (bool, error) {
		changed, err := workspacedef.MoveRepo(defs, string(id), string(repo), string(to.RepoID), stored, remote)
		if err != nil || !changed {
			return changed, err
		}
		// Live registrations keep the member's ID and Git location; they are
		// read under projects.lock, where Tasks and Epics are registered.
		uses, err := dependencies.RepoUses.RepoUses(work)
		if err != nil {
			return false, err
		}
		if registration := repoRegistration(uses, string(repo), work.Registry()); registration != "" {
			return false, projectError(ErrorProjectConflict, fmt.Sprintf("repository %s cannot move: %s", repo, registration), nil)
		}
		for _, p := range defs.Projects {
			for _, r := range p.Repos {
				if p.ID != string(id) || r.ID != string(to.RepoID) || r.Wrapper == "" {
					continue
				}
				rel, err := filepath.Rel(ResolveDefinitionPath(root, r.Wrapper), locator)
				if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
					return false, ProjectInvalidArguments(fmt.Sprintf("repository wrapper %s must contain %s; set or clear the wrapper first", ResolveDefinitionPath(root, r.Wrapper), locator))
				}
			}
		}
		return true, nil
	}, check)
	if err != nil {
		return ProjectResult{}, definitionUpdateError(err)
	}
	project, repos, err := definedProject(root, id)
	if err != nil {
		return ProjectResult{}, err
	}
	return ProjectResult{Workspace: root, Project: project, Repos: repos, Created: result.Changed}, nil
}

// RemoveProjectMember removes the member repo from Project id in the
// definitions. Its directory, files and Git are untouched and need not exist.
// The removal is refused while the Project's registry binds the member.
func RemoveProjectMember(dependencies Dependencies, id ProjectID, repo RepoID) (DefinitionRepoResult, error) {
	if err := ValidateIdentifier("project", string(id)); err != nil {
		return DefinitionRepoResult{}, err
	}
	if err := ValidateIdentifier("repository", string(repo)); err != nil {
		return DefinitionRepoResult{}, err
	}
	if dependencies.Files == nil || dependencies.RepoUses == nil {
		return DefinitionRepoResult{}, projectError(ErrorProjectIO, "dependencies: filesystem and work-item store are required", nil)
	}
	cwd, err := ProjectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	root, err := LocateContainingWorkspace(dependencies.Files, cwd)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	work, err := DefinedProjectWork(root, id)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	path := ""
	result, err := workspacedef.Apply(root, func(defs *workspacedef.Definitions) (bool, error) {
		path = ""
		members := 0
		for _, p := range defs.Projects {
			for _, r := range p.Repos {
				if p.ID == string(id) {
					members++
				}
				if p.ID == string(id) && r.ID == string(repo) {
					path = r.Path
				}
			}
		}
		if path == "" {
			return false, ProjectInvalidArguments(fmt.Sprintf("repository %s is not a member of Project %s", repo, id))
		}
		if members == 1 {
			return false, ProjectInvalidArguments(fmt.Sprintf("repository %s is the last member of Project %s; a Project needs at least one repository", repo, id))
		}
		// Registrations are read under projects.lock, where Tasks and Epics are registered.
		uses, err := dependencies.RepoUses.RepoUses(work)
		if err != nil {
			return false, err
		}
		if registration := repoRegistration(uses, string(repo), work.Registry()); registration != "" {
			return false, projectError(ErrorProjectConflict, fmt.Sprintf("repository %s cannot be removed: %s", repo, registration), nil)
		}
		return workspacedef.RemoveRepo(defs, string(id), string(repo), path)
	}, validateDefinitions)
	if err != nil {
		return DefinitionRepoResult{}, definitionUpdateError(err)
	}
	return DefinitionRepoResult{Workspace: root, ProjectID: id, RepoID: repo, Path: path, Changed: result.Changed}, nil
}

func validateProjectMemberInput(id ProjectID, repository RepoInput) error {
	if err := ValidateIdentifier("project", string(id)); err != nil {
		return err
	}
	if err := ValidateIdentifier("repository", string(repository.RepoID)); err != nil {
		return err
	}
	if repository.Path == "" {
		return ProjectInvalidArguments(fmt.Sprintf("repository %s path is empty", repository.RepoID))
	}
	return nil
}
