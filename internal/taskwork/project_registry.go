package taskwork

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/filelock"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// workIDsLock serializes the creation of Epic and Task IDs across the
// registries of all Projects. It is the only work file left in .ply.
const workIDsLock = "work-ids.lock"

// ProjectRegistry is one Project's registry snapshot.
type ProjectRegistry struct {
	Work     workspace.ProjectWork
	Registry WorkItemRegistry
}

// SnapshotWorkspace reads the registry of every defined Project, for reads
// that span the workspace. The snapshots are independent read-only files.
func SnapshotWorkspace(d Dependencies, root string) ([]ProjectRegistry, error) {
	if d.WorkItems == nil {
		return nil, workspace.WorkError(workspace.ErrorWorkIO, "work-item store dependency is required", nil)
	}
	works, err := workspace.DefinedWork(root)
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
func registryOf(registries []ProjectRegistry, project workspace.ProjectID) (ProjectRegistry, bool) {
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
func LocateEpic(d Dependencies, root string, id EpicID) (workspace.ProjectWork, error) {
	return locate(d, root, fmt.Sprintf("Epic %s", id), func(r WorkItemRegistry) bool { return hasEpic(r, id) })
}

// LocateTask returns the work of the Project whose registry holds Task id.
func LocateTask(d Dependencies, root string, id TaskID) (workspace.ProjectWork, error) {
	return locate(d, root, fmt.Sprintf("Task %s", id), func(r WorkItemRegistry) bool { return hasTask(r, id) })
}

func isWorkNotFound(err error) bool {
	var typed *workspace.WorkItemError
	return errors.As(err, &typed) && typed.Class == workspace.ErrorWorkNotFound
}

func locate(d Dependencies, root, item string, holds func(WorkItemRegistry) bool) (workspace.ProjectWork, error) {
	registries, err := SnapshotWorkspace(d, root)
	if err != nil {
		return workspace.ProjectWork{}, err
	}
	var owners []workspace.ProjectWork
	for _, r := range registries {
		if holds(r.Registry) {
			owners = append(owners, r.Work)
		}
	}
	switch len(owners) {
	case 0:
		return workspace.ProjectWork{}, workspace.WorkError(workspace.ErrorWorkNotFound, fmt.Sprintf("%s is not registered", item), nil)
	case 1:
		return owners[0], nil
	default:
		return workspace.ProjectWork{}, workspace.WorkError(workspace.ErrorWorkStoreConflict, fmt.Sprintf("%s is registered in Projects %s and %s", item, owners[0].Project, owners[1].Project), nil)
	}
}

// withNewWorkID registers a new Epic or Task ID in the Project of work. It
// holds the workspace ID lock, checks that no other Project's registry holds
// the ID, then runs operation, which writes only work's own registry. Lock
// order: definitions lock (when the caller holds it), ID lock, work lock.
func withNewWorkID(d Dependencies, work workspace.ProjectWork, item string, holds func(WorkItemRegistry) bool, operation func() error) error {
	path := filepath.Join(work.Workspace, workspace.MarkerDirectory, workIDsLock)
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
				return workspace.WorkError(workspace.ErrorWorkIdentityConflict, fmt.Sprintf("%s is already registered in Project %s; IDs are unique in the workspace", item, r.Work.Project), nil)
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
func withTaskWork(d Dependencies, root string, id TaskID, operation func(workspace.ProjectSnapshot, workspace.ProjectWork, WorkItemStoreSession) error) error {
	return d.ProjectLocks.WithSnapshotLock(root, func(projects workspace.ProjectSnapshot) error {
		work, err := LocateTask(d, root, id)
		if err != nil {
			return err
		}
		return d.WorkItems.WithLock(work, func(session WorkItemStoreSession) error { return operation(projects, work, session) })
	})
}

// withEpicWork is withTaskWork for the Project that owns Epic id.
func withEpicWork(d Dependencies, root string, id EpicID, operation func(workspace.ProjectSnapshot, workspace.ProjectWork, WorkItemStoreSession) error) error {
	return d.ProjectLocks.WithSnapshotLock(root, func(projects workspace.ProjectSnapshot) error {
		work, err := LocateEpic(d, root, id)
		if err != nil {
			return err
		}
		return d.WorkItems.WithLock(work, func(session WorkItemStoreSession) error { return operation(projects, work, session) })
	})
}
