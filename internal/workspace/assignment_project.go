package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// ResolveDefinedProject returns the Project record id from the workspace that
// contains the physical working directory, through the project metadata lookup:
// ply-workspace.json when present, otherwise the legacy registry. An empty id
// selects the single Project whose wrapper contains the working directory.
// Only the ProjectRecord leaves; repository records stay private.
func ResolveDefinedProject(d Dependencies, id ProjectID) (ProjectRecord, error) {
	if id != "" {
		if err := validateIdentifier("project", string(id)); err != nil {
			return ProjectRecord{}, err
		}
	}
	root, err := projectWorkspaceRoot(d)
	if err != nil {
		return ProjectRecord{}, err
	}
	snapshot, err := metadataProjects(d, root)
	if err != nil {
		return ProjectRecord{}, err
	}
	if id != "" {
		for _, project := range snapshot.Projects {
			if project.ID == id {
				return project, nil
			}
		}
		return ProjectRecord{}, projectError(ErrorProjectNotFound, fmt.Sprintf("project %s is not registered", id), nil)
	}
	candidates, err := projectsContainingWorkingDirectory(d.Files, snapshot.Projects)
	if err != nil {
		return ProjectRecord{}, err
	}
	if len(candidates) != 1 {
		return ProjectRecord{}, ProjectInvalidArguments(fmt.Sprintf("project is required: %d projects contain the current directory", len(candidates)))
	}
	return candidates[0], nil
}

// projectsContainingWorkingDirectory returns the Projects whose physical
// wrapper is or contains the physical working directory. A wrapper that does
// not exist is skipped.
func projectsContainingWorkingDirectory(files FileSystem, projects []ProjectRecord) ([]ProjectRecord, error) {
	cwd, err := projectPhysicalWorkingDirectory(files)
	if err != nil {
		return nil, err
	}
	var candidates []ProjectRecord
	for _, project := range projects {
		wrapper, err := files.EvalSymlinks(project.Wrapper)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, projectError(ErrorProjectPath, fmt.Sprintf("cannot resolve wrapper of project %s", project.ID), nil)
		}
		if inside, err := filepath.Rel(wrapper, cwd); err == nil && inside != ".." && !strings.HasPrefix(inside, ".."+string(filepath.Separator)) {
			candidates = append(candidates, project)
		}
	}
	return candidates, nil
}
