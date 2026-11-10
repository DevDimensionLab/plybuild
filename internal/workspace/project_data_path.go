package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// DataPathBacklog is the data_paths category of the project backlog.
const DataPathBacklog = "backlog"

// DataPathBacklogArchive is the data_paths category of the project's backlog archive.
const DataPathBacklogArchive = "backlog-archive"

// ProjectDataPath resolves the directory of category for project: its
// configured data path, or <wrapper>/<category> when none is configured.
// Nothing is created. The existing part of the path must stay inside the
// wrapper after symlinks are followed. Error details carry no paths.
func ProjectDataPath(project ProjectRecord, category string) (string, error) {
	relative, configured := project.DataPaths[category]
	if !configured {
		relative = category
	}
	if detail := dataPathProblem(relative); detail != "" {
		return "", projectError(ErrorProjectPath, fmt.Sprintf("data path %s of project %s %s", category, project.ID, detail), nil)
	}
	path := filepath.Join(project.Wrapper, relative)
	wrapper, err := filepath.EvalSymlinks(project.Wrapper)
	if err != nil {
		return "", projectError(ErrorProjectPath, fmt.Sprintf("cannot resolve wrapper of project %s", project.ID), nil)
	}
	existing := path
	for {
		if _, err := os.Lstat(existing); err == nil {
			break
		} else if !errors.Is(err, fs.ErrNotExist) {
			return "", projectError(ErrorProjectPath, fmt.Sprintf("cannot inspect data path %s of project %s", category, project.ID), nil)
		}
		existing = filepath.Dir(existing)
	}
	resolved, err := filepath.EvalSymlinks(existing)
	if err != nil {
		return "", projectError(ErrorProjectPath, fmt.Sprintf("cannot resolve data path %s of project %s", category, project.ID), nil)
	}
	if inside, err := filepath.Rel(wrapper, resolved); err != nil || inside == ".." || strings.HasPrefix(inside, ".."+string(filepath.Separator)) {
		return "", projectError(ErrorProjectPath, fmt.Sprintf("data path %s of project %s leaves the project wrapper", category, project.ID), nil)
	}
	return path, nil
}
