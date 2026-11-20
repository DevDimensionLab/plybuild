package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/projectdata"
)

const (
	WorkItemsFile         = "work-items.yaml"
	workItemsLock         = "work-items.lock"
	WorkItemLifecycleFile = "work-item-lifecycle.json"
)

// ProjectWork is the work of one Project: its work-item registry, the
// registry lock, the task-content store the registry references and the
// runtime stores of its work (lifecycle, task process, human input). They all
// live in the Project's data work directory, which projectdata owns. Values
// come only from DefinedProjectWork, so no other code places these files.
type ProjectWork struct {
	Workspace string    // workspace root, whose definitions the store binds to
	Project   ProjectID // the owning Project
	Dir       string    // <wrapper>/data/work
}

// Registry is the path of the Project's work-item registry.
func (w ProjectWork) Registry() string { return filepath.Join(w.Dir, WorkItemsFile) }

func (w ProjectWork) Lock() string { return filepath.Join(w.Dir, workItemsLock) }

func (w ProjectWork) Content() string { return filepath.Join(w.Dir, "task-content", "v1") }

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
		return WorkError(ErrorWorkIO, fmt.Sprintf("wrapper of Project %s is unavailable", w.Project), err)
	}
	if err := os.MkdirAll(w.Dir, 0o755); err != nil {
		return WorkError(ErrorWorkIO, fmt.Sprintf("create work directory of Project %s", w.Project), err)
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
	defs, _, err := ReadDefinitions(root)
	if err != nil {
		return nil, err
	}
	if defs == nil {
		return nil, MissingDefinitions(root)
	}
	works := make([]ProjectWork, 0, len(defs.Projects))
	projects := make([]ProjectRecord, 0, len(defs.Projects))
	for _, p := range defs.Projects {
		wrapper, err := CanonicalDirectory(systemFileSystem{}, root, p.Wrapper, "wrapper")
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
	if !PathOverlap(wrapper, existing) {
		return nil
	}
	return projectError(ErrorProjectConflict, fmt.Sprintf("wrapper %s overlaps the wrapper of project %s", wrapper, owner), nil)
}

// DefinedProjectWork is the work of one defined Project.
func DefinedProjectWork(root string, id ProjectID) (ProjectWork, error) {
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

func PathOverlap(a, b string) bool { return PathWithin(a, b) || PathWithin(b, a) }

func PathWithin(parent, path string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
