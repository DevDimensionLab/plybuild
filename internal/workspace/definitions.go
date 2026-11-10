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
	"strconv"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/repohistory"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
	"gopkg.in/yaml.v3"
)

type DefinitionRepoResult struct {
	Workspace string
	ProjectID ProjectID
	RepoID    RepoID
	Path      string // slash-separated, relative to Workspace
	Changed   bool
}

// AddDefinitionRepo adds an existing Git worktree root as the last repository
// member of a Project defined in ply-workspace.json. The current and the
// changed definitions are both validated against live Git while projects.lock
// is held; a new repository ID must also match its recorded history.
func AddDefinitionRepo(d Dependencies, project, repo, path string) (DefinitionRepoResult, error) {
	projectID, err := ParseProjectID(project)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	if err := validateIdentifier("repository", repo); err != nil {
		return DefinitionRepoResult{}, err
	}
	if path == "" {
		return DefinitionRepoResult{}, ProjectInvalidArguments(fmt.Sprintf("repository %s path is empty", repo))
	}
	if d.Files == nil {
		return DefinitionRepoResult{}, projectError(ErrorProjectIO, "dependencies: filesystem is required", nil)
	}
	cwd, err := projectPhysicalWorkingDirectory(d.Files)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	root, err := locateContainingWorkspace(d.Files, cwd)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	locator, err := canonicalDirectory(d.Files, cwd, path, fmt.Sprintf("repository %s", repo))
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	relative, err := filepath.Rel(root, locator)
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return DefinitionRepoResult{}, ProjectInvalidArguments(fmt.Sprintf("repository %s at %s must be inside workspace %s", repo, locator, root))
	}
	stored := filepath.ToSlash(relative)
	_, basis, err := workspacedef.Read(root)
	if err != nil {
		return DefinitionRepoResult{}, definitionUpdateError(err)
	}
	added := false
	mutate := func(defs *workspacedef.Definitions) (bool, error) {
		changed, err := workspacedef.AddRepo(defs, project, repo, stored)
		added = changed
		return changed, err
	}
	check := func(defs workspacedef.Definitions) error {
		if err := validateDefinitions(defs); err != nil {
			return err
		}
		registry, err := definitionProjects(root, defs)
		if err != nil || !added {
			return err
		}
		return checkHistoricalIdentity(root, registry, RepoID(repo))
	}
	result, err := workspacedef.Update(root, basis, mutate, check)
	if err != nil {
		return DefinitionRepoResult{}, definitionUpdateError(err)
	}
	return DefinitionRepoResult{Workspace: root, ProjectID: projectID, RepoID: RepoID(repo), Path: stored, Changed: result.Changed}, nil
}

// CheckDefinitions validates the current ply-workspace.json of the workspace
// at root like the check of AddDefinitionRepo, including live Git
// observation, without projects.lock and without writing.
func CheckDefinitions(d Dependencies, root string) error {
	observed, err := ObserveRoot(d, root)
	if err != nil {
		return err
	}
	defs, _, err := workspacedef.Read(observed.Root)
	if err != nil {
		return definitionUpdateError(err)
	}
	if defs == nil {
		return definitionUpdateError(&workspacedef.Error{Kind: workspacedef.KindMissing, Path: filepath.Join(observed.Root, workspacedef.File), Detail: "definitions do not exist"})
	}
	if err := validateDefinitions(*defs); err != nil {
		return err
	}
	_, err = definitionProjects(observed.Root, *defs)
	return err
}

func validateDefinitions(defs workspacedef.Definitions) error {
	projects := map[string]bool{}
	epics := map[string]bool{}
	for _, p := range defs.Projects {
		input := ProjectAddInput{ProjectID: ProjectID(p.ID), Name: p.Name, Wrapper: p.Wrapper}
		for _, r := range p.Repos {
			input.Repos = append(input.Repos, RepoInput{RepoID: RepoID(r.ID), Path: r.Path})
		}
		if err := validateProjectAddInput(input); err != nil {
			return err
		}
		if err := validateDataPaths(p.DataPaths); err != nil {
			return err
		}
		if projects[p.ID] {
			return ProjectInvalidArguments("duplicate Project in definitions")
		}
		projects[p.ID] = true
		for _, e := range p.Epics {
			if err := validateWorkIdentifier("Epic", e.ID); err != nil {
				return err
			}
			if err := validateWorkText("Epic title", e.Title, 128); err != nil {
				return err
			}
			member := false
			for _, r := range p.Repos {
				if r.ID == e.Repo {
					member = true
				}
			}
			if !member || e.Path == "" || !validFullBranchRef(e.Ref) || epics[e.ID] {
				return WorkInvalidArguments("invalid or duplicate Epic definition")
			}
			epics[e.ID] = true
		}
	}
	return nil
}

// validateDataPaths checks the static rules of projects[].data_paths: each
// category is a Ply identifier and each path is a non-empty relative path
// without parent traversal.
func validateDataPaths(dataPaths map[string]string) error {
	for category, path := range dataPaths {
		if err := validateIdentifier("data path category", category); err != nil {
			return err
		}
		if detail := dataPathProblem(path); detail != "" {
			return ProjectInvalidArguments(fmt.Sprintf("data path %s %s", category, detail))
		}
	}
	return nil
}

func dataPathProblem(path string) string {
	switch {
	case path == "":
		return "is empty"
	case filepath.IsAbs(path):
		return "must be relative to the project wrapper"
	}
	for _, element := range strings.Split(filepath.ToSlash(path), "/") {
		if element == ".." {
			return "must not contain parent traversal"
		}
	}
	return ""
}

// Resolve definitions using the same physical paths, Git observer and identity
// conflict checks as Project registration. No second copy is published in .ply.
func definitionProjects(root string, defs workspacedef.Definitions) (projectRegistry, error) {
	files := systemFileSystem{}
	observer := newSystemRepoObserver(files)
	registry := projectRegistry{FormatVersion: FormatVersion, Projects: []ProjectRecord{}, Repos: []RepoRecord{}}
	for _, p := range defs.Projects {
		wrapper, err := canonicalDirectory(files, root, p.Wrapper, "wrapper")
		if err != nil {
			return registry, err
		}
		project := ProjectRecord{ID: ProjectID(p.ID), Name: p.Name, Wrapper: wrapper, DataPaths: p.DataPaths}
		repos := []RepoRecord{}
		for _, r := range p.Repos {
			path, err := canonicalDirectory(files, root, r.Path, "repository")
			if err != nil {
				return registry, err
			}
			observed, err := observer.Observe(path)
			if err != nil {
				return registry, err
			}
			project.RepoIDs = append(project.RepoIDs, RepoID(r.ID))
			repos = append(repos, RepoRecord{ID: RepoID(r.ID), Locator: observed.Locator, GitCommonDir: observed.GitCommonDir})
		}
		sort.Slice(project.RepoIDs, func(i, j int) bool { return project.RepoIDs[i] < project.RepoIDs[j] })
		var errApply error
		registry, _, _, _, errApply = applyProject(registry, project, repos)
		if errApply != nil {
			return registry, errApply
		}
	}
	return registry, validateRegistry(registry)
}

// metadataProjects is the Project lookup of project metadata only. When
// ply-workspace.json exists it alone is authoritative and is projected
// statically, without Git observation or existence checks; repositories carry
// no GitCommonDir, so the snapshot must not leave metadata. Without the file
// the legacy store answers as before. It takes no lock.
func metadataProjects(d Dependencies, root string) (ProjectSnapshot, error) {
	defs, _, err := workspacedef.Read(root)
	if err != nil {
		return ProjectSnapshot{}, definitionUpdateError(err)
	}
	if defs == nil {
		projects, repos, err := d.Projects.Snapshot(root)
		return ProjectSnapshot{Projects: projects, Repos: repos}, err
	}
	if err := validateDefinitions(*defs); err != nil {
		return ProjectSnapshot{}, err
	}
	resolve := func(path string) string {
		if filepath.IsAbs(path) {
			return filepath.Clean(path)
		}
		return filepath.Join(root, filepath.FromSlash(path))
	}
	snapshot := ProjectSnapshot{Projects: []ProjectRecord{}, Repos: []RepoRecord{}}
	locators := map[RepoID]string{}
	for _, p := range defs.Projects {
		project := ProjectRecord{ID: ProjectID(p.ID), Name: p.Name, Wrapper: resolve(p.Wrapper), DataPaths: p.DataPaths}
		for _, r := range p.Repos {
			id, locator := RepoID(r.ID), resolve(r.Path)
			project.RepoIDs = append(project.RepoIDs, id)
			if existing, found := locators[id]; found {
				if existing != locator {
					return ProjectSnapshot{}, projectError(ErrorProjectConflict, fmt.Sprintf("repository %s is defined with different paths in %s", id, workspacedef.File), nil)
				}
				continue
			}
			locators[id] = locator
			snapshot.Repos = append(snapshot.Repos, RepoRecord{ID: id, Locator: locator})
		}
		sort.Slice(project.RepoIDs, func(i, j int) bool { return project.RepoIDs[i] < project.RepoIDs[j] })
		snapshot.Projects = append(snapshot.Projects, project)
	}
	return snapshot, nil
}

// checkHistoricalIdentity rejects a repository ID that the work-item registry
// recorded with another Git common directory than the one observed now.
func checkHistoricalIdentity(root string, registry projectRegistry, repo RepoID) error {
	bindings, err := historicalBindings(root)
	if err != nil {
		return err
	}
	identities, err := repohistory.CommonDirs(bindings)
	var conflict *repohistory.ConflictError
	if errors.As(err, &conflict) {
		return workError(ErrorWorkIdentityConflict, fmt.Sprintf("work-item registry %s records repository %s with conflicting Git common directories", workItemsPath(root), conflict.RepoID), err)
	}
	if err != nil {
		return err
	}
	for _, r := range registry.Repos {
		if historical := identities[string(r.ID)]; r.ID == repo && historical != "" && historical != r.GitCommonDir {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("repository %s was bound to Git common directory %s in %s, not %s; choose another repository ID", repo, historical, workItemsPath(root), r.GitCommonDir), nil)
		}
	}
	return nil
}

// historicalBindings projects the repository identities recorded in the Epic
// and Task records of .ply/work-items.yaml, inline formats 1-8 or envelopes 9
// and 10. All other registry content is ignored.
func historicalBindings(root string) ([]repohistory.Binding, error) {
	path := workItemsPath(root)
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, workError(ErrorWorkStoreConflict, fmt.Sprintf("cannot inspect work-item registry %s", path), err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return nil, workError(ErrorWorkStoreConflict, fmt.Sprintf("work-item registry %s is not a regular file", path), nil)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, workError(ErrorWorkIO, fmt.Sprintf("cannot read work-item registry %s", path), err)
	}
	bindings, err := projectHistory(contents)
	if err != nil {
		return nil, workError(ErrorWorkStoreConflict, fmt.Sprintf("invalid work-item registry %s: %v", path, err), err)
	}
	return bindings, nil
}

func projectHistory(contents []byte) ([]repohistory.Binding, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	var document yaml.Node
	if err := decoder.Decode(&document); errors.Is(err, io.EOF) {
		return nil, errors.New("registry is empty")
	} else if err != nil {
		return nil, err
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, errors.New("registry must contain exactly one YAML document")
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return nil, errors.New("registry root must be a mapping")
	}
	registry := document.Content[0]
	if err := validateQueueYAMLMap(registry); err != nil {
		return nil, err
	}
	version, ok := yamlInt(yamlField(registry, "format_version"))
	if ok && (version == 9 || version == 10) {
		inner := yamlField(registry, "registry")
		if len(registry.Content) != 4 || inner == nil || inner.Kind != yaml.MappingNode {
			return nil, errors.New("registry envelope must contain exactly format_version and a registry mapping")
		}
		registry = inner
		version, ok = yamlInt(yamlField(registry, "format_version"))
	}
	if !ok || version < 1 || version > 8 {
		return nil, errors.New("unsupported format_version")
	}
	var bindings []repohistory.Binding
	epics, err := yamlMappings(registry, "epics")
	if err != nil {
		return nil, err
	}
	for _, epic := range epics {
		repoBindings, err := yamlMappings(epic, "repo_bindings")
		if err != nil {
			return nil, err
		}
		for _, node := range repoBindings {
			binding, err := historyBinding(node)
			if err != nil {
				return nil, err
			}
			bindings = append(bindings, binding)
		}
	}
	tasks, err := yamlMappings(registry, "tasks")
	if err != nil {
		return nil, err
	}
	for _, node := range tasks {
		binding, err := historyBinding(node)
		if err != nil {
			return nil, err
		}
		bindings = append(bindings, binding)
	}
	return bindings, nil
}

// historyBinding reads repo_id and git_common_dir. A missing or empty common
// directory records no identity, so its repo_id is not validated.
func historyBinding(node *yaml.Node) (repohistory.Binding, error) {
	common, id := yamlField(node, "git_common_dir"), yamlField(node, "repo_id")
	if common == nil || isYAMLString(common) && common.Value == "" {
		return repohistory.Binding{}, nil
	}
	if !isYAMLString(common) || !validAbsoluteCleanPath(common.Value) {
		return repohistory.Binding{}, errors.New("git_common_dir must be an absolute clean path")
	}
	if id == nil || !isYAMLString(id) {
		return repohistory.Binding{}, errors.New("repo_id must be a string")
	}
	if err := validateIdentifier("repository", id.Value); err != nil {
		return repohistory.Binding{}, err
	}
	return repohistory.Binding{RepoID: id.Value, GitCommonDir: common.Value}, nil
}

// yamlField returns the value of key in a mapping whose keys are already
// validated as unique strings.
func yamlField(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

func yamlMappings(mapping *yaml.Node, key string) ([]*yaml.Node, error) {
	sequence := yamlField(mapping, key)
	if sequence == nil {
		return nil, nil
	}
	if sequence.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%s must be a sequence of mappings", key)
	}
	for _, item := range sequence.Content {
		if item.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("%s must be a sequence of mappings", key)
		}
	}
	return sequence.Content, nil
}

func yamlInt(node *yaml.Node) (int, bool) {
	if node == nil || node.Kind != yaml.ScalarNode || node.Tag != "!!int" {
		return 0, false
	}
	value, err := strconv.Atoi(node.Value)
	return value, err == nil
}

func isYAMLString(node *yaml.Node) bool {
	return node.Kind == yaml.ScalarNode && node.Tag == "!!str"
}

// definitionUpdateError gives definitions, lock and publication failures the
// existing Project and work-item error classes. Typed check errors pass as is.
func definitionUpdateError(err error) error {
	switch e := err.(type) {
	case *ProjectError, *WorkItemError:
		return err
	case *workspacedef.Error:
		return definitionError(e)
	case *filelock.Error, *filelock.ReleaseError:
		return projectsLockError(err)
	case interface{ Unwrap() []error }:
		var release *filelock.ReleaseError
		if errs := e.Unwrap(); len(errs) == 2 && errors.As(errs[1], &release) {
			return projectIOError("release after failure", release.Path, fmt.Errorf("%v; %w", definitionUpdateError(errs[0]), projectReleaseError(release)))
		}
	}
	return projectError(ErrorProjectIO, fmt.Sprintf("access definitions: %v", err), err)
}

// projectsLockError gives the lock and release failures of
// workspacedef.WithLock the classes and texts of
// systemProjectStore.WithSnapshotLock. Other errors pass as is.
func projectsLockError(err error) error {
	switch e := err.(type) {
	case *filelock.Error:
		return projectLockError(e.Path, e)
	case *filelock.ReleaseError:
		return projectReleaseError(e)
	case interface{ Unwrap() []error }:
		var release *filelock.ReleaseError
		if errs := e.Unwrap(); len(errs) == 2 && errors.As(errs[1], &release) {
			return projectIOError("release after failure", release.Path, fmt.Errorf("%v; %w", errs[0], projectReleaseError(release)))
		}
	}
	return err
}

func definitionError(e *workspacedef.Error) error {
	switch e.Kind {
	case workspacedef.KindNotFound:
		return ProjectInvalidArguments(e.Detail)
	case workspacedef.KindConflict:
		if e.Path == "" {
			return workError(ErrorWorkIdentityConflict, e.Detail, e)
		}
		return workError(ErrorWorkIdentityConflict, fmt.Sprintf("%s: %s; repeat the request", e.Path, e.Detail), e)
	case workspacedef.KindPublish:
		return workError(ErrorWorkIO, fmt.Sprintf("could not update %s; it is unchanged; repeat the request", e.Path), e)
	case workspacedef.KindReadback:
		return workError(ErrorWorkReconciliation, fmt.Sprintf("%s holds unexpected content after the repository definition update; inspect it before repeating the request", e.Path), e)
	case workspacedef.KindUsage:
		return ProjectInvalidArguments(e.Error())
	default:
		return projectStoreConflict(e.Path, e.Detail, e)
	}
}
