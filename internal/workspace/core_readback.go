package workspace

import (
	"encoding/json"
	"fmt"
	"sort"
)

const (
	ProjectListReadbackKind = "WorkspaceProjectListReadback@1"
	ProjectReadbackKind     = "WorkspaceProjectReadback@1"
	TaskListReadbackKind    = "WorkspaceTaskListReadback@1"
	CoreReadSchemaVersion   = 1
)

// TaskListFilters are optional independent AND filters, not a queue target.
// A non-nil pointer to an empty ID is invalid.
type TaskListFilters struct {
	ProjectID *ProjectID `json:"project_id"`
	RepoID    *RepoID    `json:"repo_id"`
	EpicID    *EpicID    `json:"epic_id"`
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

func validateTaskListProjectFilters(d Dependencies, root string, f TaskListFilters) error {
	if f.ProjectID == nil && f.RepoID == nil {
		return nil
	}
	if d.Projects == nil {
		return workError(ErrorWorkIO, "project store dependency is required", nil)
	}
	projects, repos, err := d.Projects.Snapshot(root)
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
			return workError(ErrorWorkNotFound, fmt.Sprintf("project %s is not registered", *f.ProjectID), nil)
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
			return workError(ErrorWorkNotFound, fmt.Sprintf("repository %s is not registered", *f.RepoID), nil)
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

type coreReadWorkspace struct {
	Root string `json:"root"`
}
type projectReadScope struct {
	ProjectID *ProjectID `json:"project_id"`
}
type projectReadRecord struct {
	ID         ProjectID          `json:"project_id"`
	Name       string             `json:"name"`
	Wrapper    string             `json:"wrapper"`
	RepoCount  int                `json:"repo_count"`
	Companions []ProjectCompanion `json:"companions"`
}
type repositoryReadRecord struct {
	ID           RepoID  `json:"repo_id"`
	Locator      string  `json:"locator"`
	GitCommonDir string  `json:"git_common_dir"`
	Wrapper      *string `json:"wrapper"`
}

func projectRead(p ProjectRecord, metadata ProjectMetadataSnapshot) projectReadRecord {
	return projectReadRecord{p.ID, p.Name, p.Wrapper, len(p.RepoIDs), metadata.Companions(p.ID)}
}

func MarshalProjectList(result ProjectListResult) ([]byte, error) {
	projects := make([]projectReadRecord, 0, len(result.Projects))
	for _, p := range result.Projects {
		projects = append(projects, projectRead(p, result.Metadata))
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].ID < projects[j].ID })
	return json.Marshal(struct {
		Kind          string              `json:"kind"`
		SchemaVersion int                 `json:"schema_version"`
		Workspace     coreReadWorkspace   `json:"workspace"`
		Scope         projectReadScope    `json:"scope"`
		Projects      []projectReadRecord `json:"projects"`
	}{ProjectListReadbackKind, CoreReadSchemaVersion, coreReadWorkspace{result.Workspace}, projectReadScope{}, projects})
}

func MarshalProject(result ProjectResult) ([]byte, error) {
	repos := make([]repositoryReadRecord, 0, len(result.Repos))
	for _, r := range result.Repos {
		repos = append(repos, repositoryReadRecord{r.ID, r.Locator, r.GitCommonDir, result.Metadata.RepositoryWrapper(r.ID)})
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].ID < repos[j].ID })
	return json.Marshal(struct {
		Kind          string                 `json:"kind"`
		SchemaVersion int                    `json:"schema_version"`
		Workspace     coreReadWorkspace      `json:"workspace"`
		Scope         projectReadScope       `json:"scope"`
		Project       projectReadRecord      `json:"project"`
		Repositories  []repositoryReadRecord `json:"repositories"`
	}{ProjectReadbackKind, CoreReadSchemaVersion, coreReadWorkspace{result.Workspace}, projectReadScope{&result.Project.ID}, projectRead(result.Project, result.Metadata), repos})
}

func MarshalTaskList(result TaskListResult) ([]byte, error) {
	type taskReadRecord struct {
		ID            TaskID         `json:"task_id"`
		ProjectID     ProjectID      `json:"project_id"`
		RepoID        RepoID         `json:"repo_id"`
		ParentEpicID  EpicID         `json:"parent_epic_id"`
		WorktreeState WorkItemState  `json:"worktree_state"`
		Lifecycle     LifecycleState `json:"lifecycle"`
		TaskListTitle
	}
	tasks := make([]taskReadRecord, 0, len(result.Tasks))
	for _, t := range result.Tasks {
		title, ok := result.Titles[t.ID]
		if !ok {
			return nil, workError(ErrorWorkIO, fmt.Sprintf("Task %s title read facts are missing", t.ID), nil)
		}
		tasks = append(tasks, taskReadRecord{t.ID, t.ProjectID, t.RepoID, t.ParentEpicID, t.WorktreeState, result.Lifecycles.Task(t.ID), title})
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return json.Marshal(struct {
		Kind          string            `json:"kind"`
		SchemaVersion int               `json:"schema_version"`
		Workspace     coreReadWorkspace `json:"workspace"`
		Scope         TaskListFilters   `json:"scope"`
		Tasks         []taskReadRecord  `json:"tasks"`
	}{TaskListReadbackKind, CoreReadSchemaVersion, coreReadWorkspace{result.Workspace}, result.Filters, tasks})
}
