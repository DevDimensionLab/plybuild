package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/repohistory"
	"github.com/devdimensionlab/plybuild/internal/workspacedef"
)

type DefinitionRepoResult struct {
	Workspace string
	ProjectID ProjectID
	RepoID    RepoID
	Path      string // slash-separated, relative to Workspace
	Changed   bool
}

// AddDefinitionRepo adds an existing Git worktree root as the last repository
// member of a Project defined in ply-workspace.json, with its remote URL
// when one is given (filled once into a member without one). The current and the
// changed definitions are both validated against live Git while projects.lock
// is held; a new repository ID must also match its recorded history.
func AddDefinitionRepo(d Dependencies, project, repo, path, remote string) (DefinitionRepoResult, error) {
	projectID, err := ParseProjectID(project)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	if err := ValidateIdentifier("repository", repo); err != nil {
		return DefinitionRepoResult{}, err
	}
	if path == "" {
		return DefinitionRepoResult{}, ProjectInvalidArguments(fmt.Sprintf("repository %s path is empty", repo))
	}
	if d.Files == nil || d.RepoUses == nil {
		return DefinitionRepoResult{}, projectError(ErrorProjectIO, "dependencies: filesystem and work-item store are required", nil)
	}
	cwd, err := ProjectPhysicalWorkingDirectory(d.Files)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	root, err := LocateContainingWorkspace(d.Files, cwd)
	if err != nil {
		return DefinitionRepoResult{}, err
	}
	locator, err := CanonicalDirectory(d.Files, cwd, path, fmt.Sprintf("repository %s", repo))
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
		changed, err := workspacedef.AddRepo(defs, project, repo, stored, remote)
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
		work, err := DefinedProjectWork(root, projectID)
		if err != nil {
			return err
		}
		return checkHistoricalIdentity(d.RepoUses, work, registry, RepoID(repo))
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
		return MissingDefinitions(observed.Root)
	}
	if err := validateDefinitions(*defs); err != nil {
		return err
	}
	_, err = definitionProjects(observed.Root, *defs)
	return err
}

// validateDefinitions checks definitions without Git. A repository belongs to
// exactly one Project, so its work items live in one Project's registry.
func validateDefinitions(defs workspacedef.Definitions) error {
	projects := map[string]bool{}
	owners := map[string]string{}
	for _, p := range defs.Projects {
		input := ProjectAddInput{ProjectID: ProjectID(p.ID), Name: p.Name, Wrapper: p.Wrapper}
		for _, r := range p.Repos {
			input.Repos = append(input.Repos, RepoInput{RepoID: RepoID(r.ID), Path: r.Path})
			if owner, found := owners[r.ID]; found && owner != p.ID {
				return projectError(ErrorProjectDuplicateRepo, fmt.Sprintf("repository %s belongs to both Project %s and Project %s; a repository belongs to exactly one Project", r.ID, owner, p.ID), nil)
			}
			owners[r.ID] = p.ID
		}
		if err := validateProjectAddInput(input); err != nil {
			return err
		}
		for _, c := range p.Companions {
			if !validCompanionRole(CompanionRole(c.Role)) || c.Path == "" {
				return ProjectInvalidArguments(fmt.Sprintf("companion %q of Project %s must have role planning, docs or other and a path", c.Path, p.ID))
			}
		}
		if projects[p.ID] {
			return ProjectInvalidArguments("duplicate Project in definitions")
		}
		projects[p.ID] = true
	}
	return nil
}

// Resolve definitions using the same physical paths, Git observer and identity
// conflict checks as Project registration. No second copy is published in .ply.
func definitionProjects(root string, defs workspacedef.Definitions) (projectRegistry, error) {
	return resolveDefinitionProjects(root, defs, false)
}

// readableProjects resolves ply-workspace.json for reading: a repository whose
// directory or Git identity cannot be read is kept as Unavailable, so its
// Project stays readable. Every other failure, such as two members with one
// Git identity, still stops the read.
func readableProjects(root string) ([]ProjectRecord, []RepoRecord, error) {
	defs, _, err := ReadDefinitions(root)
	if err != nil {
		return nil, nil, err
	}
	if defs == nil {
		return nil, nil, MissingDefinitions(root)
	}
	registry, err := resolveDefinitionProjects(root, *defs, true)
	if err != nil {
		return nil, nil, err
	}
	return registry.Projects, registry.Repos, nil
}

func resolveDefinitionProjects(root string, defs workspacedef.Definitions, keepUnavailable bool) (projectRegistry, error) {
	files := systemFileSystem{}
	observer := newSystemRepoObserver(files)
	registry := projectRegistry{Projects: []ProjectRecord{}, Repos: []RepoRecord{}}
	for _, p := range defs.Projects {
		wrapper, err := CanonicalDirectory(files, root, p.Wrapper, "wrapper")
		if err != nil {
			return registry, err
		}
		project := ProjectRecord{ID: ProjectID(p.ID), Name: p.Name, Wrapper: wrapper, Companions: definitionCompanions(root, p)}
		repos := []RepoRecord{}
		for _, r := range p.Repos {
			repo, err := definitionRepo(files, observer, root, r)
			if err != nil && !keepUnavailable {
				return registry, err
			}
			project.RepoIDs = append(project.RepoIDs, RepoID(r.ID))
			repos = append(repos, repo)
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

// definitionRepo resolves member r from its directory and Git identity. When
// that fails, the record is r as Unavailable with the failure as Diagnostic.
func definitionRepo(files FileSystem, observer RepoObserver, root string, r workspacedef.Repo) (RepoRecord, error) {
	repo := RepoRecord{ID: RepoID(r.ID), Locator: ResolveDefinitionPath(root, r.Path), Remote: r.Remote, Wrapper: definitionWrapper(root, r)}
	path, err := CanonicalDirectory(files, root, r.Path, "repository")
	if errors.Is(err, fs.ErrNotExist) {
		repo.Unavailable, repo.Diagnostic = RepoMissing, "the directory does not exist"
		return repo, err
	}
	if err != nil {
		repo.Unavailable, repo.Diagnostic = RepoUnreadable, err.Error()
		return repo, err
	}
	observed, err := observer.Observe(path)
	if err != nil {
		repo.Unavailable, repo.Diagnostic = RepoUnreadable, err.Error()
		return repo, err
	}
	repo.Locator, repo.GitCommonDir = observed.Locator, observed.GitCommonDir
	return repo, nil
}

// ObserveMember reads member repo the way project show does. It returns why
// the member is unavailable and the diagnostic, or "" when its directory is a
// Git worktree root.
func ObserveMember(repo RepoRecord) (RepoAvailability, string) {
	files := systemFileSystem{}
	observed, _ := definitionRepo(files, newSystemRepoObserver(files), "", workspacedef.Repo{ID: string(repo.ID), Path: repo.Locator})
	return observed.Unavailable, observed.Diagnostic
}

// updateDefinitions applies mutate to ply-workspace.json as it is under
// projects.lock and validates the result against live Git before it is
// written. It reports whether the file changed.
func updateDefinitions(root string, mutate func(*workspacedef.Definitions) (bool, error)) (bool, error) {
	check := func(defs workspacedef.Definitions) error {
		if err := validateDefinitions(defs); err != nil {
			return err
		}
		_, err := definitionProjects(root, defs)
		return err
	}
	result, err := workspacedef.Apply(root, mutate, check)
	if err != nil {
		return false, definitionUpdateError(err)
	}
	return result.Changed, nil
}

// definedProject resolves one Project of ply-workspace.json and its repositories.
func definedProject(root string, id ProjectID) (ProjectRecord, []RepoRecord, error) {
	projects, repos, err := newSystemProjectStore().Snapshot(root)
	if err != nil {
		return ProjectRecord{}, nil, err
	}
	for _, project := range projects {
		if project.ID == id {
			return project, repositoriesForProject(project, repos), nil
		}
	}
	return ProjectRecord{}, nil, projectError(ErrorProjectNotFound, fmt.Sprintf("project %s is not defined in %s", id, workspacedef.File), nil)
}

// definitionPath is the stored form of a physical path: relative to the
// workspace root when inside it, otherwise absolute.
func definitionPath(root, path string) string {
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return path
	}
	return filepath.ToSlash(relative)
}

// ResolveDefinitionPath is the physical form of a stored definitions path of
// the workspace at the physical root, without checking that it exists.
func ResolveDefinitionPath(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(root, filepath.FromSlash(path))
}

// definitionCompanions resolves the companions of p; nil when it has none.
func definitionCompanions(root string, p workspacedef.Project) []ProjectCompanion {
	var companions []ProjectCompanion
	for _, c := range p.Companions {
		companions = append(companions, ProjectCompanion{CompanionRole(c.Role), ResolveDefinitionPath(root, c.Path)})
	}
	return companions
}

// definitionKnowledge is p's knowledge pointer as written, with the locator of
// its repository among repos; nil when p has none.
func definitionKnowledge(p workspacedef.Project, repos []RepoRecord) *KnowledgeSource {
	if p.Knowledge == nil {
		return nil
	}
	source := &KnowledgeSource{Repo: RepoID(p.Knowledge.Repo), Path: p.Knowledge.Path}
	for _, r := range repos {
		if r.ID == source.Repo {
			source.Locator = r.Locator
		}
	}
	return source
}

// definitionWrapper resolves the wrapper of r; empty when it has none.
func definitionWrapper(root string, r workspacedef.Repo) string {
	if r.Wrapper == "" {
		return ""
	}
	return ResolveDefinitionPath(root, r.Wrapper)
}

func MissingDefinitions(root string) error {
	return definitionError(&workspacedef.Error{Kind: workspacedef.KindMissing, Path: filepath.Join(root, workspacedef.File), Detail: "definitions do not exist"})
}

// metadataProjects is the Project lookup without Git. When
// ply-workspace.json exists it alone is authoritative and is projected
// statically, without Git observation or existence checks; repositories carry
// no GitCommonDir. It takes no lock.
func metadataProjects(d Dependencies, root string) (ProjectSnapshot, error) {
	defs, _, err := workspacedef.Read(root)
	if err != nil {
		return ProjectSnapshot{}, definitionUpdateError(err)
	}
	if defs == nil {
		return ProjectSnapshot{}, MissingDefinitions(root)
	}
	if err := validateDefinitions(*defs); err != nil {
		return ProjectSnapshot{}, err
	}
	snapshot := ProjectSnapshot{Projects: []ProjectRecord{}, Repos: []RepoRecord{}}
	for _, p := range defs.Projects {
		project := ProjectRecord{ID: ProjectID(p.ID), Name: p.Name, Wrapper: ResolveDefinitionPath(root, p.Wrapper), Companions: definitionCompanions(root, p)}
		var repos []RepoRecord
		for _, r := range p.Repos {
			id := RepoID(r.ID)
			project.RepoIDs = append(project.RepoIDs, id)
			repos = append(repos, RepoRecord{ID: id, Locator: ResolveDefinitionPath(root, r.Path), Wrapper: definitionWrapper(root, r)})
		}
		snapshot.Repos = append(snapshot.Repos, repos...)
		sort.Slice(project.RepoIDs, func(i, j int) bool { return project.RepoIDs[i] < project.RepoIDs[j] })
		project.Knowledge = definitionKnowledge(p, repos)
		snapshot.Projects = append(snapshot.Projects, project)
	}
	return snapshot, nil
}

// checkHistoricalIdentity rejects a repository ID that the Project's
// work-item registry recorded with another Git common directory than the one
// observed now.
func checkHistoricalIdentity(uses RepoUses, work ProjectWork, registry projectRegistry, repo RepoID) error {
	recorded, err := uses.RepoUses(work)
	if err != nil {
		return err
	}
	var bindings []repohistory.Binding
	for _, use := range recorded {
		bindings = append(bindings, repohistory.Binding{RepoID: string(use.RepoID), GitCommonDir: use.GitCommonDir})
	}
	identities, err := repohistory.CommonDirs(bindings)
	var conflict *repohistory.ConflictError
	if errors.As(err, &conflict) {
		return WorkError(ErrorWorkIdentityConflict, fmt.Sprintf("work-item registry %s records repository %s with conflicting Git common directories", work.Registry(), conflict.RepoID), err)
	}
	if err != nil {
		return err
	}
	for _, r := range registry.Repos {
		if historical := identities[string(r.ID)]; r.ID == repo && historical != "" && historical != r.GitCommonDir {
			return WorkError(ErrorWorkIdentityConflict, fmt.Sprintf("repository %s was bound to Git common directory %s in %s, not %s; choose another repository ID", repo, historical, work.Registry(), r.GitCommonDir), nil)
		}
	}
	return nil
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
			return WorkError(ErrorWorkIdentityConflict, e.Detail, e)
		}
		return WorkError(ErrorWorkIdentityConflict, fmt.Sprintf("%s: %s; repeat the request", e.Path, e.Detail), e)
	case workspacedef.KindPublish:
		return WorkError(ErrorWorkIO, fmt.Sprintf("could not update %s; it is unchanged; repeat the request", e.Path), e)
	case workspacedef.KindReadback:
		return WorkError(ErrorWorkReconciliation, fmt.Sprintf("%s holds unexpected content after the repository definition update; inspect it before repeating the request", e.Path), e)
	case workspacedef.KindUsage:
		return ProjectInvalidArguments(e.Error())
	case workspacedef.KindMissing:
		return projectStoreConflict(e.Path, e.Detail+"; run ply workspace init", e)
	default:
		return projectStoreConflict(e.Path, e.Detail, e)
	}
}
