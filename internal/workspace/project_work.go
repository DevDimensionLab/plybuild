package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/projectdata"
)

// workIDsLock serializes the creation of Epic and Task IDs across the
// registries of all Projects. It is the only work file left in .ply.
const workIDsLock = "work-ids.lock"

// ProjectWork is the work of one Project: its work-item registry, the
// registry lock, the task-content store the registry references and the
// runtime stores of its work (lifecycle, task process, human input). They all
// live in the Project's data work directory, which projectdata owns. Values
// come only from projectWork, so no other code places these files.
type ProjectWork struct {
	Workspace string    // workspace root, whose definitions the store binds to
	Project   ProjectID // the owning Project
	Dir       string    // <wrapper>/data/work
}

// Registry is the path of the Project's work-item registry.
func (w ProjectWork) Registry() string { return filepath.Join(w.Dir, WorkItemsFile) }

func (w ProjectWork) lock() string { return filepath.Join(w.Dir, workItemsLock) }

func (w ProjectWork) content() string { return filepath.Join(w.Dir, "task-content", "v1") }

// Lifecycle is the path of the Project's work-item lifecycle history.
func (w ProjectWork) Lifecycle() string { return filepath.Join(w.Dir, WorkItemLifecycleFile) }

// HumanInput is the directory of the Project's human-input questions and
// their answers.
func (w ProjectWork) HumanInput() string { return filepath.Join(w.Dir, "human-input", "v1") }

// AgentWork is the directory of the Project's agent follow-ups.
func (w ProjectWork) AgentWork() string { return filepath.Join(w.Dir, "agent-work", "v1") }

// TaskProcess is the directory of the Project's Task process records: Task
// journals with their kept sources, and goal-start receipts.
func (w ProjectWork) TaskProcess() string { return filepath.Join(w.Dir, "task-process") }

// Create makes the work directory on first write. The wrapper must exist;
// data/ and data/work/ may be created.
func (w ProjectWork) Create() error {
	if _, err := os.Stat(filepath.Dir(filepath.Dir(w.Dir))); err != nil {
		return workError(ErrorWorkIO, fmt.Sprintf("wrapper of Project %s is unavailable", w.Project), err)
	}
	if err := os.MkdirAll(w.Dir, 0o755); err != nil {
		return workError(ErrorWorkIO, fmt.Sprintf("create work directory of Project %s", w.Project), err)
	}
	return nil
}

func workOf(root string, project ProjectRecord) (ProjectWork, error) {
	dir, err := project.Data().Work()
	var dataErr *projectdata.Error
	if errors.As(err, &dataErr) {
		return ProjectWork{}, projectError(ErrorProjectPath, dataErr.Detail, err)
	}
	if err != nil {
		return ProjectWork{}, err
	}
	return ProjectWork{Workspace: root, Project: project.ID, Dir: dir}, nil
}

// DefinedWork is the work of every Project in ply-workspace.json, sorted by
// Project ID. Wrappers resolve like definitionProjects resolves them, without
// observing Git, so reads stay available when a repository is not.
func DefinedWork(root string) ([]ProjectWork, error) {
	defs, _, err := readDefinitions(root)
	if err != nil {
		return nil, err
	}
	if defs == nil {
		return nil, missingDefinitions(root)
	}
	works := make([]ProjectWork, 0, len(defs.Projects))
	projects := make([]ProjectRecord, 0, len(defs.Projects))
	for _, p := range defs.Projects {
		wrapper, err := canonicalDirectory(systemFileSystem{}, root, p.Wrapper, "wrapper")
		if err != nil {
			return nil, err
		}
		for _, existing := range projects {
			if err := wrapperConflict(wrapper, existing.Wrapper, existing.ID); err != nil {
				return nil, err
			}
		}
		project := ProjectRecord{ID: ProjectID(p.ID), Name: p.Name, Wrapper: wrapper}
		projects = append(projects, project)
		work, err := workOf(root, project)
		if err != nil {
			return nil, err
		}
		works = append(works, work)
	}
	sort.Slice(works, func(i, j int) bool { return works[i].Project < works[j].Project })
	return works, nil
}

// wrapperConflict refuses wrapper when it overlaps the wrapper of Project
// owner: the same directory, or one inside the other. A Project's data lives
// in its wrapper, so overlapping wrappers would mix two Projects' data.
func wrapperConflict(wrapper, existing string, owner ProjectID) error {
	if !pathOverlap(wrapper, existing) {
		return nil
	}
	return projectError(ErrorProjectConflict, fmt.Sprintf("wrapper %s overlaps the wrapper of project %s", wrapper, owner), nil)
}

// projectWork is the work of one defined Project.
func projectWork(root string, id ProjectID) (ProjectWork, error) {
	works, err := DefinedWork(root)
	if err != nil {
		return ProjectWork{}, err
	}
	for _, work := range works {
		if work.Project == id {
			return work, nil
		}
	}
	return ProjectWork{}, projectError(ErrorProjectNotFound, fmt.Sprintf("project %s is not defined in %s", id, DefinitionsFile), nil)
}

// ProjectRegistry is one Project's registry snapshot.
type ProjectRegistry struct {
	Work     ProjectWork
	Registry WorkItemRegistry
}

// SnapshotWorkspace reads the registry of every defined Project, for reads
// that span the workspace. The snapshots are independent read-only files.
func SnapshotWorkspace(d Dependencies, root string) ([]ProjectRegistry, error) {
	if d.WorkItems == nil {
		return nil, workError(ErrorWorkIO, "work-item store dependency is required", nil)
	}
	works, err := DefinedWork(root)
	if err != nil {
		return nil, err
	}
	out := make([]ProjectRegistry, 0, len(works))
	for _, work := range works {
		registry, err := d.WorkItems.SnapshotRegistrations(work)
		if err != nil {
			return nil, err
		}
		out = append(out, ProjectRegistry{Work: work, Registry: registry})
	}
	return out, nil
}

// MergeRegistries joins Project registries into one read-only view for
// listing. It has no physical digest and no evidence references: content and
// evidence are read through the owning ProjectRegistry.
func MergeRegistries(registries []ProjectRegistry) WorkItemRegistry {
	out := emptyWorkItemRegistry()
	for _, p := range registries {
		r := p.Registry
		out.Epics = append(out.Epics, r.Epics...)
		out.Tasks = append(out.Tasks, r.Tasks...)
		out.WorktreeOperations = append(out.WorktreeOperations, r.WorktreeOperations...)
		out.TaskResults = append(out.TaskResults, r.TaskResults...)
		out.HumanQARecords = append(out.HumanQARecords, r.HumanQARecords...)
		out.IntegrationAuthorities = append(out.IntegrationAuthorities, r.IntegrationAuthorities...)
		out.IntegrationIntents = append(out.IntegrationIntents, r.IntegrationIntents...)
		out.IntegrationAttempts = append(out.IntegrationAttempts, r.IntegrationAttempts...)
		out.IntegrationResults = append(out.IntegrationResults, r.IntegrationResults...)
		out.TaskProblemRevisions = append(out.TaskProblemRevisions, r.TaskProblemRevisions...)
		out.TaskSpecRevisions = append(out.TaskSpecRevisions, r.TaskSpecRevisions...)
		out.TaskSpecAssessments = append(out.TaskSpecAssessments, r.TaskSpecAssessments...)
		out.TaskSolutionSelections = append(out.TaskSolutionSelections, r.TaskSolutionSelections...)
		out.TaskResultSpecBindings = append(out.TaskResultSpecBindings, r.TaskResultSpecBindings...)
		out.TaskContentPublications = append(out.TaskContentPublications, r.TaskContentPublications...)
		out.EpicBaseVersions = append(out.EpicBaseVersions, r.EpicBaseVersions...)
		out.EpicBaseUpdates = append(out.EpicBaseUpdates, r.EpicBaseUpdates...)
		out.TaskQueueEvents = append(out.TaskQueueEvents, r.TaskQueueEvents...)
		out.TaskPreparations = append(out.TaskPreparations, r.TaskPreparations...)
		out.EpicCreateOperations = append(out.EpicCreateOperations, r.EpicCreateOperations...)
	}
	sortWorkRegistry(&out)
	return out
}

// registryOf returns the registry of project, or false.
func registryOf(registries []ProjectRegistry, project ProjectID) (ProjectRegistry, bool) {
	for _, r := range registries {
		if r.Work.Project == project {
			return r, true
		}
	}
	return ProjectRegistry{}, false
}

func hasEpic(r WorkItemRegistry, id EpicID) bool {
	if epic, _ := findEpic(r, id); epic != nil {
		return true
	}
	op, _ := findEpicCreate(r, id)
	return op != nil
}

func hasTask(r WorkItemRegistry, id TaskID) bool {
	task, _ := findTask(r, id)
	return task != nil
}

// LocateEpic returns the work of the Project whose registry holds Epic id,
// registered or pending creation.
func LocateEpic(d Dependencies, root string, id EpicID) (ProjectWork, error) {
	return locate(d, root, fmt.Sprintf("Epic %s", id), func(r WorkItemRegistry) bool { return hasEpic(r, id) })
}

// LocateTask returns the work of the Project whose registry holds Task id.
func LocateTask(d Dependencies, root string, id TaskID) (ProjectWork, error) {
	return locate(d, root, fmt.Sprintf("Task %s", id), func(r WorkItemRegistry) bool { return hasTask(r, id) })
}

func isWorkNotFound(err error) bool {
	var typed *WorkItemError
	return errors.As(err, &typed) && typed.Class == ErrorWorkNotFound
}

func locate(d Dependencies, root, item string, holds func(WorkItemRegistry) bool) (ProjectWork, error) {
	registries, err := SnapshotWorkspace(d, root)
	if err != nil {
		return ProjectWork{}, err
	}
	var owners []ProjectWork
	for _, r := range registries {
		if holds(r.Registry) {
			owners = append(owners, r.Work)
		}
	}
	switch len(owners) {
	case 0:
		return ProjectWork{}, workError(ErrorWorkNotFound, fmt.Sprintf("%s is not registered", item), nil)
	case 1:
		return owners[0], nil
	default:
		return ProjectWork{}, workError(ErrorWorkStoreConflict, fmt.Sprintf("%s is registered in Projects %s and %s", item, owners[0].Project, owners[1].Project), nil)
	}
}

// withNewWorkID registers a new Epic or Task ID in the Project of work. It
// holds the workspace ID lock, checks that no other Project's registry holds
// the ID, then runs operation, which writes only work's own registry. Lock
// order: definitions lock (when the caller holds it), ID lock, work lock.
func withNewWorkID(d Dependencies, work ProjectWork, item string, holds func(WorkItemRegistry) bool, operation func() error) error {
	path := filepath.Join(work.Workspace, MarkerDirectory, workIDsLock)
	lock, err := filelock.Acquire(path, filelock.Options{})
	var lockErr *filelock.Error
	if errors.As(err, &lockErr) {
		return workItemLockError(path, lockErr)
	}
	if err != nil {
		return err
	}
	opErr := func() error {
		registries, err := SnapshotWorkspace(d, work.Workspace)
		if err != nil {
			return err
		}
		for _, r := range registries {
			if r.Work.Project != work.Project && holds(r.Registry) {
				return workError(ErrorWorkIdentityConflict, fmt.Sprintf("%s is already registered in Project %s; IDs are unique in the workspace", item, r.Work.Project), nil)
			}
		}
		return operation()
	}()
	var releaseErr error
	if err := lock.Release(); err != nil {
		releaseErr = workItemReleaseError(err)
	}
	if opErr != nil {
		return opErr
	}
	return releaseErr
}

// withTaskWork holds the definitions lock, locates the Project that owns Task
// id and holds that Project's work lock while operation runs.
func withTaskWork(d Dependencies, root string, id TaskID, operation func(ProjectSnapshot, ProjectWork, WorkItemStoreSession) error) error {
	return d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		work, err := LocateTask(d, root, id)
		if err != nil {
			return err
		}
		return d.WorkItems.WithLock(work, func(session WorkItemStoreSession) error { return operation(projects, work, session) })
	})
}

// withEpicWork is withTaskWork for the Project that owns Epic id.
func withEpicWork(d Dependencies, root string, id EpicID, operation func(ProjectSnapshot, ProjectWork, WorkItemStoreSession) error) error {
	return d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		work, err := LocateEpic(d, root, id)
		if err != nil {
			return err
		}
		return d.WorkItems.WithLock(work, func(session WorkItemStoreSession) error { return operation(projects, work, session) })
	})
}
