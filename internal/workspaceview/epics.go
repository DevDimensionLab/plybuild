package workspaceview

import (
	"sort"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const EpicListKind = "WorkspaceEpicListReadback@1"

type EpicFilters struct {
	ProjectID *workspace.ProjectID `json:"project_id"`
}

type EpicBase struct {
	RepoID      workspace.RepoID `json:"repo_id"`
	Revision    int              `json:"revision"`
	OID         string           `json:"oid"`
	Tree        string           `json:"tree"`
	Freshness   string           `json:"freshness"`
	SourceBasis string           `json:"source_basis"`
}

type EpicWorktree struct {
	RepoID      workspace.RepoID     `json:"repo_id"`
	WorktreeID  workspace.WorktreeID `json:"worktree_id"`
	Locator     string               `json:"locator"`
	Ref         string               `json:"ref"`
	Freshness   string               `json:"freshness"`
	SourceBasis string               `json:"source_basis"`
}

type ProgressCounts struct {
	Integrated int `json:"integrated"`
	Attention  int `json:"attention"`
	InProgress int `json:"in_progress"`
	Nothing    int `json:"nothing"`
}

type EpicRow struct {
	EpicID          workspace.EpicID         `json:"epic_id"`
	ProjectID       workspace.ProjectID      `json:"project_id"`
	Title           string                   `json:"title"`
	Lifecycle       workspace.LifecycleState `json:"lifecycle"`
	Base            []EpicBase               `json:"base"`
	Worktrees       []EpicWorktree           `json:"worktrees"`
	TaskCount       int                      `json:"task_count"`
	ProgressCounts  ProgressCounts           `json:"progress_counts"`
	LastActivityUTC *string                  `json:"last_activity_utc"`
}

type EpicList struct {
	Kind          string       `json:"kind"`
	SchemaVersion int          `json:"schema_version"`
	Workspace     WorkspaceRef `json:"workspace"`
	Scope         EpicFilters  `json:"scope"`
	Epics         []EpicRow    `json:"epics"`
	Freshness     string       `json:"freshness"`
	SourceBasis   string       `json:"source_basis"`
	Reasons       []string     `json:"reasons"`
}

func ReadEpics(d workspace.Dependencies, filters EpicFilters) (EpicList, error) {
	if err := (workspace.TaskListFilters{ProjectID: filters.ProjectID}).Validate(); err != nil {
		return EpicList{}, err
	}
	s, err := LoadSnapshot(d)
	if err != nil {
		return EpicList{}, err
	}
	if err := ValidateTaskFilters(s, workspace.TaskListFilters{ProjectID: filters.ProjectID}); err != nil {
		return EpicList{}, err
	}
	process := readProgressProcess(d, s)
	s.RefreshFreshness(d)
	return BuildEpics(s, filters, process)
}

func BuildEpics(s *Snapshot, filters EpicFilters, process map[workspace.TaskID]ProcessFacts) (EpicList, error) {
	tasks, err := BuildTaskList(s, workspace.TaskListFilters{ProjectID: filters.ProjectID}, process)
	if err != nil {
		return EpicList{}, err
	}
	return epicsFromTasks(s, filters, tasks), nil
}

func epicsFromTasks(s *Snapshot, filters EpicFilters, tasks TaskList) EpicList {
	out := EpicList{EpicListKind, 1, tasks.Workspace, filters, []EpicRow{}, tasks.Freshness, "registered", tasks.Reasons}
	tasksByEpic := map[workspace.EpicID][]TaskRow{}
	activityByEpic := workspace.RegisteredEpicActivityTimes(s.Registry)
	for _, task := range tasks.Tasks {
		tasksByEpic[task.ParentEpicID] = append(tasksByEpic[task.ParentEpicID], task)
	}
	for _, epic := range s.Registry.Epics {
		if filters.ProjectID != nil && epic.ProjectID != *filters.ProjectID {
			continue
		}
		row := EpicRow{EpicID: epic.ID, ProjectID: epic.ProjectID, Title: epic.Title, Lifecycle: s.Lifecycle.Epic(epic.ID), Base: []EpicBase{}, Worktrees: []EpicWorktree{}, LastActivityUTC: activityByEpic[epic.ID]}
		for _, base := range workspace.RegisteredEpicBases(s.Registry, epic) {
			row.Base = append(row.Base, EpicBase{base.RepoID, base.Revision, base.OID, base.Tree, s.Freshness, "registered"})
			row.Worktrees = append(row.Worktrees, EpicWorktree{base.RepoID, base.WorktreeID, base.Locator, base.Ref, s.Freshness, "registered"})
		}
		for _, task := range tasksByEpic[epic.ID] {
			row.TaskCount++
			switch task.Progress.State {
			case "integrated", "completed":
				row.ProgressCounts.Integrated++
			case "attention", "cleanup_pending":
				row.ProgressCounts.Attention++
			case "in_progress":
				row.ProgressCounts.InProgress++
			default:
				row.ProgressCounts.Nothing++
			}
			row.LastActivityUTC = latestTime(row.LastActivityUTC, task.Progress.LastActivityUTC)
		}
		for _, event := range s.Lifecycle.Events {
			if event.SubjectKind == "epic" && event.SubjectID == string(epic.ID) {
				row.LastActivityUTC = latestTime(row.LastActivityUTC, &event.RecordedAtUTC)
			}
		}
		sort.Slice(row.Base, func(i, j int) bool { return row.Base[i].RepoID < row.Base[j].RepoID })
		sort.Slice(row.Worktrees, func(i, j int) bool { return row.Worktrees[i].RepoID < row.Worktrees[j].RepoID })
		out.Epics = append(out.Epics, row)
	}
	sort.Slice(out.Epics, func(i, j int) bool { return out.Epics[i].EpicID < out.Epics[j].EpicID })
	return out
}
