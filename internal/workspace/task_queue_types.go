package workspace

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Format 4 keeps adopted bindings immutable and records append-only decisions.
type QueueTarget struct {
	ProjectID        ProjectID  `yaml:"project_id" json:"project_id"`
	RepoID           RepoID     `yaml:"repo_id" json:"repo_id"`
	EpicID           EpicID     `yaml:"epic_id" json:"epic_id"`
	GitCommonDir     string     `yaml:"git_common_dir" json:"git_common_dir"`
	ParentWorktreeID WorktreeID `yaml:"parent_worktree_id" json:"parent_worktree_id"`
	ParentLocator    string     `yaml:"parent_locator" json:"parent_locator"`
	ParentRef        string     `yaml:"parent_ref" json:"parent_ref"`
}

type QueueReason struct {
	Code    string `yaml:"code" json:"code"`
	Message string `yaml:"message" json:"message"`
}

type TaskQueueEvent struct {
	QueueID       string         `yaml:"queue_id" json:"queue_id"`
	Revision      int            `yaml:"revision" json:"revision"`
	Kind          string         `yaml:"kind" json:"kind"`
	RequestSHA256 string         `yaml:"request_sha256" json:"request_sha256"`
	Request       map[string]any `yaml:"request" json:"request"`
	RecordedAtUTC string         `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}

type EpicBaseVersion struct {
	ProjectID        ProjectID  `yaml:"project_id" json:"project_id"`
	RepoID           RepoID     `yaml:"repo_id" json:"repo_id"`
	EpicID           EpicID     `yaml:"epic_id" json:"epic_id"`
	Revision         int        `yaml:"revision" json:"revision"`
	PreviousRevision *int       `yaml:"previous_revision" json:"previous_revision"`
	WorktreeID       WorktreeID `yaml:"worktree_id" json:"worktree_id"`
	Locator          string     `yaml:"locator" json:"locator"`
	Ref              string     `yaml:"ref" json:"ref"`
	GitCommonDir     string     `yaml:"git_common_dir" json:"git_common_dir"`
	OID              string     `yaml:"oid" json:"oid"`
	Tree             string     `yaml:"tree" json:"tree"`
	UpdateID         *string    `yaml:"update_id" json:"update_id"`
}

type EpicAffectedSelection struct {
	TaskID    TaskID           `yaml:"task_id" json:"task_id"`
	Selection *TaskDecisionRef `yaml:"selection" json:"selection"`
}

type WorkspaceEpicBaseUpdatePlan struct {
	Kind                  string                   `yaml:"kind" json:"kind"`
	SchemaVersion         int                      `yaml:"schema_version" json:"schema_version"`
	Workspace             IntegrationPlanWorkspace `yaml:"workspace" json:"workspace"`
	Target                QueueTarget              `yaml:"target" json:"target"`
	RegistryUpgrade       *TaskRegistryUpgrade     `yaml:"registry_upgrade" json:"registry_upgrade"`
	RegistryFormatVersion int                      `yaml:"registry_format_version" json:"registry_format_version"`
	ExpectedRevision      int                      `yaml:"expected_revision" json:"expected_revision"`
	PreviousOID           string                   `yaml:"previous_oid" json:"previous_oid"`
	PreviousTree          string                   `yaml:"previous_tree" json:"previous_tree"`
	NextOID               string                   `yaml:"next_oid" json:"next_oid"`
	NextTree              string                   `yaml:"next_tree" json:"next_tree"`
	ObservedParent        PlanWorktreeObservation  `yaml:"observed_parent" json:"observed_parent"`
	AffectedSelections    []EpicAffectedSelection  `yaml:"affected_selections" json:"affected_selections"`
	CurrentPreparationID  *string                  `yaml:"current_preparation_id" json:"current_preparation_id"`
	UnresolvedOperations  []string                 `yaml:"unresolved_operations" json:"unresolved_operations"`
	ActivePreparationIDs  *[]string                `yaml:"active_preparation_ids,omitempty" json:"active_preparation_ids,omitempty"`
}

// A plan of the current schema has no member of the removed preparation;
// historical plans keep theirs exactly as recorded and digested.
const currentEpicBaseUpdatePlanSchema = 3

func (p WorkspaceEpicBaseUpdatePlan) MarshalJSON() ([]byte, error) {
	type wire WorkspaceEpicBaseUpdatePlan
	if p.SchemaVersion != currentEpicBaseUpdatePlanSchema {
		return json.Marshal(wire(p))
	}
	return json.Marshal(struct {
		wire
		CurrentPreparationID *string `json:"current_preparation_id,omitempty"`
	}{wire: wire(p)})
}

func (p WorkspaceEpicBaseUpdatePlan) MarshalYAML() (any, error) {
	type wire WorkspaceEpicBaseUpdatePlan
	if p.SchemaVersion != currentEpicBaseUpdatePlanSchema {
		return wire(p), nil
	}
	var n yaml.Node
	if err := n.Encode(wire(p)); err != nil {
		return nil, err
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == "current_preparation_id" {
			n.Content = append(n.Content[:i], n.Content[i+2:]...)
			break
		}
	}
	return &n, nil
}

type EpicBaseUpdateOutcome struct {
	Revision *int          `yaml:"revision" json:"revision"`
	Reasons  []QueueReason `yaml:"reasons" json:"reasons"`
}

type EpicBaseUpdate struct {
	ID            string                      `yaml:"id" json:"id"`
	PlanSHA256    string                      `yaml:"plan_sha256" json:"plan_sha256"`
	Plan          WorkspaceEpicBaseUpdatePlan `yaml:"plan" json:"plan"`
	Phase         string                      `yaml:"phase" json:"phase"`
	RecordedAtUTC string                      `yaml:"recorded_at_utc" json:"recorded_at_utc"`
	Outcome       *EpicBaseUpdateOutcome      `yaml:"outcome" json:"outcome"`
}

type QueueRegistryReadback struct {
	FormatVersion   int                  `yaml:"format_version" json:"format_version"`
	SHA256          string               `yaml:"sha256" json:"sha256"`
	RequiredUpgrade *TaskRegistryUpgrade `yaml:"required_upgrade" json:"required_upgrade"`
}

type QueuePending struct {
	SerialReadiness   string                  `yaml:"serial_readiness,omitempty" json:"serial_readiness,omitempty"`
	ParallelReadiness string                  `yaml:"parallel_readiness,omitempty" json:"parallel_readiness,omitempty"`
	Goal              *TaskGoalRef            `yaml:"goal,omitempty" json:"goal,omitempty"`
	Executor          *TaskExecutorAssignment `yaml:"executor,omitempty" json:"executor,omitempty"`
	Rank              int                     `yaml:"rank" json:"rank"`
	TaskID            TaskID                  `yaml:"task_id" json:"task_id"`
	Title             string                  `yaml:"title" json:"title"`
	Selection         *TaskDecisionRef        `yaml:"selection" json:"selection"`
	SpecID            *string                 `yaml:"spec_id" json:"spec_id"`
	SpecRevision      *int                    `yaml:"spec_revision" json:"spec_revision"`
	BaseRevision      *int                    `yaml:"base_revision" json:"base_revision"`
	ParentOID         *string                 `yaml:"parent_oid" json:"parent_oid"`
	State             string                  `yaml:"state" json:"state"`
	Reasons           []QueueReason           `yaml:"reasons" json:"reasons"`
}

type WorkspaceTaskQueueReadback struct {
	Kind          string                   `yaml:"kind" json:"kind"`
	SchemaVersion int                      `yaml:"schema_version" json:"schema_version"`
	Workspace     IntegrationPlanWorkspace `yaml:"workspace" json:"workspace"`
	Target        QueueTarget              `yaml:"target" json:"target"`
	Registry      QueueRegistryReadback    `yaml:"registry" json:"registry"`
	QueueID       string                   `yaml:"queue_id" json:"queue_id"`
	Revision      int                      `yaml:"revision" json:"revision"`
	Pending       []QueuePending           `yaml:"pending" json:"pending"`
	Reasons       []QueueReason            `yaml:"reasons" json:"reasons"`
}
