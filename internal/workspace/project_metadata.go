package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

type CompanionRole string

const (
	CompanionPlanning       CompanionRole = "planning"
	CompanionDocs           CompanionRole = "docs"
	CompanionOther          CompanionRole = "other"
	ProjectCompanionAdd                   = "companion_add"
	ProjectCompanionRemove                = "companion_remove"
	ProjectRepoWrapperSet                 = "repo_wrapper_set"
	ProjectRepoWrapperClear               = "repo_wrapper_clear"
)

type ProjectCompanion struct {
	Role    CompanionRole `json:"role"`
	Locator string        `json:"locator"`
}

func validCompanionRole(role CompanionRole) bool {
	return role == CompanionPlanning || role == CompanionDocs || role == CompanionOther
}

// sortedCompanions orders companions by role, then locator, and is never nil.
func sortedCompanions(companions []ProjectCompanion) []ProjectCompanion {
	items := append([]ProjectCompanion{}, companions...)
	sort.Slice(items, func(i, j int) bool {
		if items[i].Role != items[j].Role {
			return items[i].Role < items[j].Role
		}
		return items[i].Locator < items[j].Locator
	})
	return items
}

// ProjectMetadataInput is one change of a Project's companions or of a
// repository wrapper in ply-workspace.json.
type ProjectMetadataInput struct {
	Operation string
	ProjectID ProjectID
	RepoID    RepoID
	Role      CompanionRole
	Locator   string
}

func (in ProjectMetadataInput) Validate() error {
	if _, err := ParseProjectID(string(in.ProjectID)); err != nil {
		return err
	}
	switch in.Operation {
	case ProjectCompanionAdd, ProjectCompanionRemove:
		if !validCompanionRole(in.Role) {
			return ProjectInvalidArguments("--role must be planning, docs, or other")
		}
		if in.RepoID != "" {
			return ProjectInvalidArguments("companion changes do not take a repository ID")
		}
	case ProjectRepoWrapperSet, ProjectRepoWrapperClear:
		if err := ValidateIdentifier("repository", string(in.RepoID)); err != nil {
			return err
		}
		if in.Role != "" {
			return ProjectInvalidArguments("repository wrapper changes do not take a companion role")
		}
	default:
		return ProjectInvalidArguments("unknown project metadata operation")
	}
	if in.Operation == ProjectRepoWrapperClear {
		if in.Locator != "" {
			return ProjectInvalidArguments("clearing a repository wrapper does not take a path")
		}
	} else if in.Locator == "" || !utf8.ValidString(in.Locator) || strings.IndexByte(in.Locator, 0) >= 0 {
		return ProjectInvalidArguments("metadata path must be nonempty valid UTF-8 and contain no NUL")
	}
	return nil
}

// ProjectMetadataMutation is the Project as project show reads it after the
// change, and whether ply-workspace.json changed.
type ProjectMetadataMutation struct {
	ProjectResult
	Changed bool
}

// UpdateProjectMetadata applies one validated change to ply-workspace.json
// under projects.lock; a relative path is resolved against the working
// directory. An unchanged registration writes nothing.
func UpdateProjectMetadata(d Dependencies, in ProjectMetadataInput) (ProjectMetadataMutation, error) {
	if err := in.Validate(); err != nil {
		return ProjectMetadataMutation{}, err
	}
	root, err := projectWorkspaceRoot(d)
	if err != nil {
		return ProjectMetadataMutation{}, err
	}
	if in.Locator != "" {
		cwd, err := ProjectPhysicalWorkingDirectory(d.Files)
		if err != nil {
			return ProjectMetadataMutation{}, err
		}
		if in.Locator, err = projectMetadataPath(d.Files, cwd, in.Locator, in.Operation == ProjectCompanionRemove); err != nil {
			return ProjectMetadataMutation{}, err
		}
		if !utf8.ValidString(in.Locator) {
			return ProjectMetadataMutation{}, ProjectInvalidArguments("physical metadata path must be valid UTF-8")
		}
	}
	changed, err := updateDefinitions(root, func(defs *workspacedef.Definitions) (bool, error) {
		return applyProjectMetadata(root, defs, in, true)
	})
	if err != nil {
		return ProjectMetadataMutation{}, err
	}
	shown, err := ShowProject(d, in.ProjectID)
	return ProjectMetadataMutation{shown, changed}, err
}

// SetRepositoryWrapperAt sets the physical wrapper of repo in the workspace
// at root when the repository has none. The same wrapper is no change;
// another wrapper is a conflict and is kept.
func SetRepositoryWrapperAt(d Dependencies, root, project, repo, wrapper string) (bool, error) {
	in := ProjectMetadataInput{Operation: ProjectRepoWrapperSet, ProjectID: ProjectID(project), RepoID: RepoID(repo), Locator: wrapper}
	if err := in.Validate(); err != nil {
		return false, err
	}
	observed, err := ObserveRoot(d, root)
	if err != nil {
		return false, err
	}
	return updateDefinitions(observed.Root, func(defs *workspacedef.Definitions) (bool, error) {
		return applyProjectMetadata(observed.Root, defs, in, false)
	})
}

// applyProjectMetadata applies in, whose locator is physical, to defs. A set
// replaces another wrapper only when replace is true.
func applyProjectMetadata(root string, defs *workspacedef.Definitions, in ProjectMetadataInput, replace bool) (bool, error) {
	var project *workspacedef.Project
	for i := range defs.Projects {
		if defs.Projects[i].ID == string(in.ProjectID) {
			project = &defs.Projects[i]
		}
	}
	if project == nil {
		return false, projectError(ErrorProjectNotFound, fmt.Sprintf("project %s is not registered", in.ProjectID), nil)
	}
	stored := definitionPath(root, in.Locator)
	switch in.Operation {
	case ProjectCompanionAdd, ProjectCompanionRemove:
		companion := workspacedef.Companion{Role: string(in.Role), Path: stored}
		for i, existing := range project.Companions {
			if existing == companion {
				if in.Operation == ProjectCompanionAdd {
					return false, nil
				}
				project.Companions = append(append([]workspacedef.Companion{}, project.Companions[:i]...), project.Companions[i+1:]...)
				return true, nil
			}
		}
		if in.Operation == ProjectCompanionRemove {
			return false, nil
		}
		project.Companions = append(project.Companions, companion)
		return true, nil
	}
	var repo *workspacedef.Repo
	for i := range project.Repos {
		if project.Repos[i].ID == string(in.RepoID) {
			repo = &project.Repos[i]
		}
	}
	if repo == nil {
		return false, WorkError(ErrorWorkProjectConflict, fmt.Sprintf("repository %s is not a member of Project %s", in.RepoID, in.ProjectID), nil)
	}
	if in.Operation == ProjectRepoWrapperClear {
		changed := repo.Wrapper != ""
		repo.Wrapper = ""
		return changed, nil
	}
	if repo.Wrapper == stored {
		return false, nil
	}
	if repo.Wrapper != "" && !replace {
		return false, projectError(ErrorProjectConflict, fmt.Sprintf("repository %s already has wrapper %s, not %s", in.RepoID, ResolveDefinitionPath(root, repo.Wrapper), in.Locator), nil)
	}
	rel, err := filepath.Rel(in.Locator, ResolveDefinitionPath(root, repo.Path))
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false, ProjectInvalidArguments("repository wrapper must contain the registered repository locator")
	}
	repo.Wrapper = stored
	return true, nil
}

func projectMetadataPath(files FileSystem, cwd, input string, remove bool) (string, error) {
	if !remove {
		return CanonicalDirectory(files, cwd, input, "metadata")
	}
	path := input
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	path = filepath.Clean(path)
	ancestor := path
	missing := []string{}
	for {
		physical, err := files.EvalSymlinks(ancestor)
		if err == nil {
			// Removed children still belong to their physical parent. Preserve
			// that identity when the caller repeats the alias used at add time.
			return filepath.Join(append([]string{physical}, missing...)...), nil
		}
		if !errors.Is(err, fs.ErrNotExist) || filepath.Dir(ancestor) == ancestor {
			return "", projectPathError("resolve metadata removal", path, err)
		}
		missing = append([]string{filepath.Base(ancestor)}, missing...)
		ancestor = filepath.Dir(ancestor)
	}
}
