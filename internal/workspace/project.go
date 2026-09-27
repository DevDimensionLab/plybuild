package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const identifierRule = "^[a-z][a-z0-9]*(-[a-z0-9]+)*$ (1..63 ASCII bytes)"

var identifierPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

type ProjectID string
type RepoID string

type ProjectAddInput struct {
	ProjectID ProjectID
	Name      string
	Wrapper   string
	Repos     []RepoInput
}

type RepoInput struct {
	RepoID RepoID
	Path   string
}

type ProjectRecord struct {
	ID      ProjectID `yaml:"id"`
	Name    string    `yaml:"name"`
	Wrapper string    `yaml:"wrapper"`
	RepoIDs []RepoID  `yaml:"repo_ids"`
}

type RepoRecord struct {
	ID           RepoID `yaml:"id"`
	Locator      string `yaml:"locator"`
	GitCommonDir string `yaml:"git_common_dir"`
}

type ProjectResult struct {
	Workspace string
	Project   ProjectRecord
	Repos     []RepoRecord
	Created   bool
}

type ProjectListResult struct {
	Workspace string
	Projects  []ProjectRecord
}

const (
	ErrorProjectInvalidArguments ErrorClass = "workspace_project_invalid_arguments"
	ErrorWorkspaceNotFound       ErrorClass = "workspace_not_found"
	ErrorProjectWorkspace        ErrorClass = "workspace_project_workspace_conflict"
	ErrorProjectPath             ErrorClass = "workspace_project_path_error"
	ErrorProjectRepoInvalid      ErrorClass = "workspace_project_repo_invalid"
	ErrorProjectRepoObservation  ErrorClass = "workspace_project_repo_observation_error"
	ErrorProjectDuplicateRepo    ErrorClass = "workspace_project_duplicate_repo"
	ErrorProjectNotFound         ErrorClass = "workspace_project_not_found"
	ErrorProjectConflict         ErrorClass = "workspace_project_conflict"
	ErrorProjectStoreConflict    ErrorClass = "workspace_project_store_conflict"
	ErrorProjectIO               ErrorClass = "workspace_project_io_error"
)

type ProjectError struct {
	Class  ErrorClass
	Detail string
	Err    error
}

func (err *ProjectError) Error() string {
	detail := err.Detail
	if detail == "" && err.Err != nil {
		detail = err.Err.Error()
	}
	return fmt.Sprintf("%s: %s", err.Class, detail)
}

func (err *ProjectError) Unwrap() error { return err.Err }

func ProjectInvalidArguments(detail string) error {
	return &ProjectError{Class: ErrorProjectInvalidArguments, Detail: detail}
}

func ParseProjectID(value string) (ProjectID, error) {
	if err := validateIdentifier("project", value); err != nil {
		return "", err
	}
	return ProjectID(value), nil
}

func ParseProjectAddInput(projectID, name, wrapper string, repositories []string) (ProjectAddInput, error) {
	id, err := ParseProjectID(projectID)
	if err != nil {
		return ProjectAddInput{}, err
	}
	input := ProjectAddInput{ProjectID: id, Name: name, Wrapper: wrapper}
	for _, repository := range repositories {
		repoID, path, found := strings.Cut(repository, "=")
		if !found || repoID == "" || path == "" {
			return ProjectAddInput{}, ProjectInvalidArguments(fmt.Sprintf("repository member %q must be <repo-id>=<path>", repository))
		}
		if err := validateIdentifier("repository", repoID); err != nil {
			return ProjectAddInput{}, err
		}
		input.Repos = append(input.Repos, RepoInput{RepoID: RepoID(repoID), Path: path})
	}
	if err := validateProjectAddInput(input); err != nil {
		return ProjectAddInput{}, err
	}
	return input, nil
}

func validateProjectAddInput(input ProjectAddInput) error {
	if err := validateIdentifier("project", string(input.ProjectID)); err != nil {
		return err
	}
	if err := validateName(input.Name); err != nil {
		return err
	}
	if input.Wrapper == "" {
		return ProjectInvalidArguments("--wrapper is required")
	}
	if len(input.Repos) == 0 {
		return ProjectInvalidArguments("at least one --repo <repo-id>=<path> is required")
	}
	seen := make(map[RepoID]struct{}, len(input.Repos))
	for _, repository := range input.Repos {
		if err := validateIdentifier("repository", string(repository.RepoID)); err != nil {
			return err
		}
		if repository.Path == "" {
			return ProjectInvalidArguments(fmt.Sprintf("repository %s path is empty", repository.RepoID))
		}
		if _, duplicate := seen[repository.RepoID]; duplicate {
			return projectError(ErrorProjectDuplicateRepo, fmt.Sprintf("repository ID %s appears more than once", repository.RepoID), nil)
		}
		seen[repository.RepoID] = struct{}{}
	}
	return nil
}

func validateIdentifier(kind, value string) error {
	if len(value) < 1 || len(value) > 63 || !isASCII(value) || !identifierPattern.MatchString(value) {
		return ProjectInvalidArguments(fmt.Sprintf("%s ID %q must match %s", kind, value, identifierRule))
	}
	return nil
}

func validateName(name string) error {
	if !utf8.ValidString(name) || name == "" || strings.TrimSpace(name) != name || utf8.RuneCountInString(name) > 128 {
		return ProjectInvalidArguments("project name must be valid UTF-8, contain 1..128 characters, and have no surrounding whitespace")
	}
	for _, character := range name {
		if unicode.IsControl(character) {
			return ProjectInvalidArguments("project name must be one line without control characters")
		}
	}
	return nil
}

func isASCII(value string) bool {
	for index := 0; index < len(value); index++ {
		if value[index] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func AddProject(dependencies Dependencies, input ProjectAddInput) (ProjectResult, error) {
	if err := validateProjectAddInput(input); err != nil {
		return ProjectResult{}, err
	}
	if dependencies.Files == nil || dependencies.Repos == nil || dependencies.Projects == nil {
		return ProjectResult{}, projectError(ErrorProjectIO, "dependencies: filesystem, repository observer, and project store are required", nil)
	}
	cwd, err := projectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return ProjectResult{}, err
	}
	workspaceRoot, err := locateContainingWorkspace(dependencies.Files, cwd)
	if err != nil {
		return ProjectResult{}, err
	}
	wrapper, err := canonicalDirectory(dependencies.Files, cwd, input.Wrapper, "wrapper")
	if err != nil {
		return ProjectResult{}, err
	}

	project := ProjectRecord{ID: input.ProjectID, Name: input.Name, Wrapper: wrapper}
	repositories := make([]RepoRecord, 0, len(input.Repos))
	commonDirectories := make(map[string]RepoID, len(input.Repos))
	for _, repository := range input.Repos {
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
			return ProjectResult{}, projectError(ErrorProjectRepoObservation, fmt.Sprintf("observe repository %s at %s: observer returned malformed physical paths", repository.RepoID, locator), nil)
		}
		if existing, duplicate := commonDirectories[observation.GitCommonDir]; duplicate {
			return ProjectResult{}, projectError(ErrorProjectDuplicateRepo, fmt.Sprintf("repositories %s and %s resolve to Git common directory %s", existing, repository.RepoID, observation.GitCommonDir), nil)
		}
		commonDirectories[observation.GitCommonDir] = repository.RepoID
		project.RepoIDs = append(project.RepoIDs, repository.RepoID)
		repositories = append(repositories, RepoRecord{ID: repository.RepoID, Locator: observation.Locator, GitCommonDir: observation.GitCommonDir})
	}
	sort.Slice(project.RepoIDs, func(i, j int) bool { return project.RepoIDs[i] < project.RepoIDs[j] })
	sort.Slice(repositories, func(i, j int) bool { return repositories[i].ID < repositories[j].ID })

	storedProject, storedRepositories, created, err := dependencies.Projects.Add(workspaceRoot, project, repositories)
	if err != nil {
		return ProjectResult{}, err
	}
	return ProjectResult{Workspace: workspaceRoot, Project: storedProject, Repos: storedRepositories, Created: created}, nil
}

func ShowProject(dependencies Dependencies, id ProjectID) (ProjectResult, error) {
	if err := validateIdentifier("project", string(id)); err != nil {
		return ProjectResult{}, err
	}
	workspaceRoot, projects, repositories, err := projectSnapshot(dependencies)
	if err != nil {
		return ProjectResult{}, err
	}
	for _, project := range projects {
		if project.ID == id {
			return ProjectResult{Workspace: workspaceRoot, Project: project, Repos: repositoriesForProject(project, repositories)}, nil
		}
	}
	return ProjectResult{}, projectError(ErrorProjectNotFound, fmt.Sprintf("project %s is not registered in Ply workspace %s", id, workspaceRoot), nil)
}

func ListProjects(dependencies Dependencies) (ProjectListResult, error) {
	workspaceRoot, projects, _, err := projectSnapshot(dependencies)
	if err != nil {
		return ProjectListResult{}, err
	}
	return ProjectListResult{Workspace: workspaceRoot, Projects: projects}, nil
}

func projectSnapshot(dependencies Dependencies) (string, []ProjectRecord, []RepoRecord, error) {
	if dependencies.Files == nil || dependencies.Projects == nil {
		return "", nil, nil, projectError(ErrorProjectIO, "dependencies: filesystem and project store are required", nil)
	}
	cwd, err := projectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return "", nil, nil, err
	}
	workspaceRoot, err := locateContainingWorkspace(dependencies.Files, cwd)
	if err != nil {
		return "", nil, nil, err
	}
	projects, repositories, err := dependencies.Projects.Snapshot(workspaceRoot)
	return workspaceRoot, projects, repositories, err
}

func projectPhysicalWorkingDirectory(files FileSystem) (string, error) {
	workingDirectory, err := files.Getwd()
	if err != nil {
		return "", projectPathError("getwd", "", err)
	}
	return canonicalDirectory(files, "", workingDirectory, "working directory")
}

func canonicalDirectory(files FileSystem, base, input, kind string) (string, error) {
	path := input
	if !filepath.IsAbs(path) && base != "" {
		path = filepath.Join(base, path)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", projectPathError("abs", input, err)
	}
	absolute = filepath.Clean(absolute)
	physical, err := files.EvalSymlinks(absolute)
	if err != nil {
		return "", projectPathError("eval-symlinks", absolute, err)
	}
	physical, err = filepath.Abs(physical)
	if err != nil {
		return "", projectPathError("abs", physical, err)
	}
	physical = filepath.Clean(physical)
	info, err := files.Stat(physical)
	if err != nil {
		return "", projectPathError("stat", physical, err)
	}
	if !info.IsDir() {
		return "", projectPathError("stat", physical, fmt.Errorf("%s is not a directory", kind))
	}
	return physical, nil
}

func locateContainingWorkspace(files FileSystem, cwd string) (string, error) {
	for ancestor := cwd; ; ancestor = filepath.Dir(ancestor) {
		reserved := filepath.Join(ancestor, MarkerDirectory)
		info, err := files.Lstat(reserved)
		switch {
		case errors.Is(err, fs.ErrNotExist):
		case err != nil:
			return "", projectError(ErrorProjectWorkspace, fmt.Sprintf("inspect %s: %v", reserved, err), err)
		case !info.IsDir() || info.Mode()&fs.ModeSymlink != 0:
			return "", projectError(ErrorProjectWorkspace, fmt.Sprintf("%s: %s", reserved, reasonReservedPath), nil)
		default:
			markerPath := filepath.Join(reserved, MarkerFile)
			markerInfo, markerErr := files.Lstat(markerPath)
			switch {
			case errors.Is(markerErr, fs.ErrNotExist):
			case markerErr != nil:
				return "", projectError(ErrorProjectWorkspace, fmt.Sprintf("inspect %s: %v", markerPath, markerErr), markerErr)
			case !markerInfo.Mode().IsRegular() || markerInfo.Mode()&fs.ModeSymlink != 0:
				return "", projectError(ErrorProjectWorkspace, fmt.Sprintf("%s: %s", markerPath, reasonMarkerType), nil)
			default:
				if reason := readAndValidateMarker(files, markerPath, ancestor); reason != "" {
					return "", projectError(ErrorProjectWorkspace, fmt.Sprintf("%s: %s", markerPath, reason), nil)
				}
				return ancestor, nil
			}
		}
		if ancestor == filepath.Dir(ancestor) {
			break
		}
	}
	return "", projectError(ErrorWorkspaceNotFound, fmt.Sprintf("no Ply workspace contains %s", cwd), nil)
}

func validStoredPath(path string) bool {
	return path != "" && filepath.IsAbs(path) && filepath.Clean(path) == path
}

func repositoriesForProject(project ProjectRecord, repositories []RepoRecord) []RepoRecord {
	byID := make(map[RepoID]RepoRecord, len(repositories))
	for _, repository := range repositories {
		byID[repository.ID] = repository
	}
	result := make([]RepoRecord, 0, len(project.RepoIDs))
	for _, id := range project.RepoIDs {
		result = append(result, byID[id])
	}
	return result
}

func projectPathError(operation, path string, err error) error {
	detail := operation
	if path != "" {
		detail += " " + path
	}
	if err != nil {
		detail += ": " + err.Error()
	}
	return projectError(ErrorProjectPath, detail, err)
}

func projectError(class ErrorClass, detail string, err error) error {
	return &ProjectError{Class: class, Detail: detail, Err: err}
}
