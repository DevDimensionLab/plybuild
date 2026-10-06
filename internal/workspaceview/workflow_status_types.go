package workspaceview

import (
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const WorkflowStatusKind = "WorkflowStatusReadback@1"

type WorkflowStatusOptions struct {
	Filters    workspace.TaskListFilters
	IncludeAll bool
}

type WorkflowStatus struct {
	Kind          string                    `json:"kind"`
	SchemaVersion int                       `json:"schema_version"`
	Workspace     WorkspaceRef              `json:"workspace"`
	Filters       workspace.TaskListFilters `json:"filters"`
	IncludeAll    bool                      `json:"include_all"`
	ObservedAtUTC string                    `json:"observed_at_utc"`
	SourceBasis   string                    `json:"source_basis"`
	Freshness     string                    `json:"freshness"`
	Projects      []WorkflowStatusProject   `json:"projects"`
	Epics         []EpicRow                 `json:"epics"`
	Items         []WorkflowStatusItem      `json:"items"`
	Counts        WorkflowStatusCounts      `json:"counts"`
	Diagnostics   []WorkflowDiagnostic      `json:"diagnostics"`
}

type WorkflowStatusProject struct {
	ProjectID workspace.ProjectID `json:"project_id"`
	Name      string              `json:"name"`
	RepoIDs   []workspace.RepoID  `json:"repo_ids"`
}

type WorkflowStatusCounts struct {
	Total      int                  `json:"total"`
	Visible    int                  `json:"visible"`
	Hidden     WorkflowStatusHidden `json:"hidden"`
	Categories map[string]int       `json:"categories"`
}

type WorkflowStatusHidden struct {
	Inactive  int `json:"inactive"`
	Completed int `json:"completed"`
	Backlog   int `json:"backlog"`
}

type WorkflowStatusItem struct {
	TaskID    workspace.TaskID    `json:"task_id"`
	ProjectID workspace.ProjectID `json:"project_id"`
	RepoID    workspace.RepoID    `json:"repo_id"`
	EpicID    workspace.EpicID    `json:"epic_id"`
	workspace.TaskListTitle
	WorktreeState   workspace.WorkItemState  `json:"worktree_state"`
	TaskLifecycle   workspace.LifecycleState `json:"task_lifecycle"`
	EpicLifecycle   workspace.LifecycleState `json:"epic_lifecycle"`
	Category        string                   `json:"category"`
	Progress        Progress                 `json:"progress"`
	LastActivityUTC *string                  `json:"last_activity_utc"`
	Freshness       string                   `json:"freshness"`
	NextActions     []WorkflowStatusAction   `json:"next_actions"`
	Runs            []taskrun.InventoryRun   `json:"runs"`
	Queues          []WorkflowStatusQueue    `json:"queues"`
}

// Current distinguishes actionable needs from facts retained by --all. The
// underlying actor, reason, evidence and time remain intact for historical work.
type WorkflowStatusAction struct {
	Action
	Current bool    `json:"current"`
	RunID   *string `json:"run_id"`
	QueueID *string `json:"queue_id"`
}

type WorkflowStatusQueue struct {
	QueueID     string                     `json:"queue_id"`
	Revision    int                        `json:"revision"`
	Target      workspace.QueueTarget      `json:"target"`
	State       string                     `json:"state"`
	Rank        *int                       `json:"rank"`
	SpecID      *string                    `json:"spec_id"`
	Spec        *workspace.TaskRevisionRef `json:"spec"`
	Goal        *workspace.TaskGoalRef     `json:"goal"`
	Selection   *workspace.TaskDecisionRef `json:"selection"`
	Readiness   string                     `json:"readiness"`
	Current     *workspace.QueueCurrent    `json:"current"`
	Freshness   string                     `json:"freshness"`
	SinceUTC    *string                    `json:"since_utc"`
	Reasons     []workspace.QueueReason    `json:"reasons"`
	EvidenceIDs []string                   `json:"evidence_ids"`
}

type WorkflowDiagnostic struct {
	Code        string   `json:"code"`
	Source      string   `json:"source"`
	Actor       string   `json:"actor"`
	Severity    string   `json:"severity"`
	Reason      string   `json:"reason"`
	SinceUTC    *string  `json:"since_utc"`
	TaskID      *string  `json:"task_id"`
	RunID       *string  `json:"run_id"`
	QueueID     *string  `json:"queue_id"`
	EvidenceIDs []string `json:"evidence_ids"`
}
