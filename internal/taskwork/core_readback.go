package taskwork

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"sort"
)

const TaskListReadbackKind = "WorkspaceTaskListReadback@1"

// TaskListFilters are optional independent AND filters, not a queue target.
// A non-nil pointer to an empty ID is invalid.
type TaskListFilters struct {
	ProjectID *workspace.ProjectID `json:"project_id"`
	RepoID    *workspace.RepoID    `json:"repo_id"`
	EpicID    *EpicID              `json:"epic_id"`
}

func (f TaskListFilters) Validate() error {
	if f.ProjectID != nil {
		if err := validateWorkIdentifier("project", string(*f.ProjectID)); err != nil {
			return err
		}
	}
	if f.RepoID != nil {
		if err := validateWorkIdentifier("repository", string(*f.RepoID)); err != nil {
			return err
		}
	}
	if f.EpicID != nil {
		if err := validateWorkIdentifier("Epic", string(*f.EpicID)); err != nil {
			return err
		}
	}
	return nil
}

func validateTaskListProjectFilters(root string, registry WorkItemRegistry, lifecycles WorkItemLifecycles, f TaskListFilters) error {
	if f.ProjectID == nil && f.RepoID == nil {
		return nil
	}
	projects, repos, err := registryProjects(root, registry, lifecycles)
	if err != nil {
		return err
	}
	if f.ProjectID != nil {
		found := false
		for _, p := range projects {
			if p.ID == *f.ProjectID {
				found = true
				break
			}
		}
		if !found {
			return workspace.WorkError(workspace.ErrorWorkNotFound, fmt.Sprintf("project %s is not registered", *f.ProjectID), nil)
		}
	}
	if f.RepoID != nil {
		found := false
		for _, r := range repos {
			if r.ID == *f.RepoID {
				found = true
				break
			}
		}
		if !found {
			return workspace.WorkError(workspace.ErrorWorkNotFound, fmt.Sprintf("repository %s is not registered", *f.RepoID), nil)
		}
	}
	return nil
}

// TaskListTitle carries read facts independently of human presentation text.
type TaskListTitle struct {
	Title  *string `json:"title"`
	Source string  `json:"title_source"`
	Status string  `json:"title_status"`
}

func MarshalTaskList(result TaskListResult) ([]byte, error) {
	type taskReadRecord struct {
		ID            TaskID              `json:"task_id"`
		ProjectID     workspace.ProjectID `json:"project_id"`
		RepoID        workspace.RepoID    `json:"repo_id"`
		ParentEpicID  EpicID              `json:"parent_epic_id"`
		WorktreeState WorkItemState       `json:"worktree_state"`
		Lifecycle     LifecycleState      `json:"lifecycle"`
		TaskListTitle
	}
	tasks := make([]taskReadRecord, 0, len(result.Tasks))
	for _, t := range result.Tasks {
		title, ok := result.Titles[t.ID]
		if !ok {
			return nil, workspace.WorkError(workspace.ErrorWorkIO, fmt.Sprintf("Task %s title read facts are missing", t.ID), nil)
		}
		tasks = append(tasks, taskReadRecord{t.ID, t.ProjectID, t.RepoID, t.ParentEpicID, t.WorktreeState, result.Lifecycles.Task(t.ID), title})
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return json.Marshal(struct {
		Kind          string                      `json:"kind"`
		SchemaVersion int                         `json:"schema_version"`
		Workspace     workspace.CoreReadWorkspace `json:"workspace"`
		Scope         TaskListFilters             `json:"scope"`
		Tasks         []taskReadRecord            `json:"tasks"`
	}{TaskListReadbackKind, workspace.CoreReadSchemaVersion, workspace.CoreReadWorkspace{Root: result.Workspace}, result.Filters, tasks})
}
