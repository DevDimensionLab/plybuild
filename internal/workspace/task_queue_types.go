package workspace

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

type QueueSelector struct {
	Kind   string  `yaml:"kind" json:"kind"`
	TaskID *TaskID `yaml:"task_id" json:"task_id"`
}

type QueueReason struct {
	Code    string `yaml:"code" json:"code"`
	Message string `yaml:"message" json:"message"`
}

type QueueEntry struct {
	TaskID    TaskID           `yaml:"task_id" json:"task_id"`
	Selection *TaskDecisionRef `yaml:"selection" json:"selection"`
}

type QueueHumanDecision struct {
	ActorClaim   string `yaml:"actor_claim" json:"actor_claim"`
	DecidedAtUTC string `yaml:"decided_at_utc" json:"decided_at_utc"`
	Source       string `yaml:"source" json:"source"`
	Statement    string `yaml:"statement" json:"statement"`
}

type WorkspaceTaskQueueDraft struct {
	Kind             string               `yaml:"kind" json:"kind"`
	SchemaVersion    int                  `yaml:"schema_version" json:"schema_version"`
	PublicationKey   string               `yaml:"publication_key" json:"publication_key"`
	ProjectID        ProjectID            `yaml:"project_id" json:"project_id"`
	RepoID           RepoID               `yaml:"repo_id" json:"repo_id"`
	EpicID           EpicID               `yaml:"epic_id" json:"epic_id"`
	ExpectedRevision int                  `yaml:"expected_revision" json:"expected_revision"`
	Entries          []QueueEntry         `yaml:"entries" json:"entries"`
	HumanDecision    QueueHumanDecision   `yaml:"human_decision" json:"human_decision"`
	RegistryUpgrade  *TaskRegistryUpgrade `yaml:"registry_upgrade" json:"registry_upgrade"`
}

type TaskQueueEvent struct {
	QueueID       string         `yaml:"queue_id" json:"queue_id"`
	Revision      int            `yaml:"revision" json:"revision"`
	Kind          string         `yaml:"kind" json:"kind"`
	RequestSHA256 string         `yaml:"request_sha256" json:"request_sha256"`
	Request       map[string]any `yaml:"request" json:"request"`
	RecordedAtUTC string         `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}

type WorkspaceTaskPreparePlan struct {
	Kind                  string                   `yaml:"kind" json:"kind"`
	SchemaVersion         int                      `yaml:"schema_version" json:"schema_version"`
	Workspace             IntegrationPlanWorkspace `yaml:"workspace" json:"workspace"`
	Target                QueueTarget              `yaml:"target" json:"target"`
	RegistryUpgrade       *TaskRegistryUpgrade     `yaml:"registry_upgrade" json:"registry_upgrade"`
	RegistryFormatVersion int                      `yaml:"registry_format_version" json:"registry_format_version"`
	QueueID               string                   `yaml:"queue_id" json:"queue_id"`
	QueueRevision         int                      `yaml:"queue_revision" json:"queue_revision"`
	Selector              QueueSelector            `yaml:"selector" json:"selector"`
	TaskID                TaskID                   `yaml:"task_id" json:"task_id"`
	Selection             TaskDecisionRef          `yaml:"selection" json:"selection"`
	Problem               TaskRevisionRef          `yaml:"problem" json:"problem"`
	SpecID                string                   `yaml:"spec_id" json:"spec_id"`
	Spec                  TaskRevisionRef          `yaml:"spec" json:"spec"`
	Assessment            TaskDecisionRef          `yaml:"assessment" json:"assessment"`
	RequiredInputs        []TaskSpecRequiredInput  `yaml:"required_inputs" json:"required_inputs"`
	BaseRevision          int                      `yaml:"base_revision" json:"base_revision"`
	ParentOID             string                   `yaml:"parent_oid" json:"parent_oid"`
	ParentTree            string                   `yaml:"parent_tree" json:"parent_tree"`
	Branch                string                   `yaml:"branch" json:"branch"`
	WorktreePath          string                   `yaml:"worktree_path" json:"worktree_path"`
	ObservedParent        PlanWorktreeObservation  `yaml:"observed_parent" json:"observed_parent"`
}

type TaskPreparationOutcome struct {
	Kind          string     `yaml:"kind" json:"kind"`
	WorktreeID    WorktreeID `yaml:"worktree_id" json:"worktree_id"`
	ObservedOID   string     `yaml:"observed_oid" json:"observed_oid"`
	ObservedTree  string     `yaml:"observed_tree" json:"observed_tree"`
	RecordedAtUTC string     `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}

type TaskPreparation struct {
	ID           string                   `yaml:"id" json:"id"`
	PlanSHA256   string                   `yaml:"plan_sha256" json:"plan_sha256"`
	Plan         WorkspaceTaskPreparePlan `yaml:"plan" json:"plan"`
	OperationID  WorktreeOperationID      `yaml:"operation_id" json:"operation_id"`
	CreatedAtUTC string                   `yaml:"created_at_utc" json:"created_at_utc"`
	Outcome      *TaskPreparationOutcome  `yaml:"outcome" json:"outcome"`
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

type QueueCurrent struct {
	Delivery      *DeliveryAgreement `yaml:"delivery,omitempty" json:"delivery,omitempty"`
	PreparationID string             `yaml:"preparation_id" json:"preparation_id"`
	TaskID        TaskID             `yaml:"task_id" json:"task_id"`
	State         string             `yaml:"state" json:"state"`
	Disposition   string             `yaml:"disposition" json:"disposition"`
	WorktreePath  string             `yaml:"worktree_path" json:"worktree_path"`
	Reasons       []QueueReason      `yaml:"reasons" json:"reasons"`
}

type QueuePending struct {
	Delivery     *DeliveryAgreement      `yaml:"delivery,omitempty" json:"delivery,omitempty"`
	Goal         *TaskGoalRef            `yaml:"goal,omitempty" json:"goal,omitempty"`
	Executor     *TaskExecutorAssignment `yaml:"executor,omitempty" json:"executor,omitempty"`
	Rank         int                     `yaml:"rank" json:"rank"`
	TaskID       TaskID                  `yaml:"task_id" json:"task_id"`
	Title        string                  `yaml:"title" json:"title"`
	Selection    *TaskDecisionRef        `yaml:"selection" json:"selection"`
	SpecID       *string                 `yaml:"spec_id" json:"spec_id"`
	SpecRevision *int                    `yaml:"spec_revision" json:"spec_revision"`
	BaseRevision *int                    `yaml:"base_revision" json:"base_revision"`
	ParentOID    *string                 `yaml:"parent_oid" json:"parent_oid"`
	State        string                  `yaml:"state" json:"state"`
	Reasons      []QueueReason           `yaml:"reasons" json:"reasons"`
}

type WorkspaceTaskQueueReadback struct {
	Kind          string                   `yaml:"kind" json:"kind"`
	SchemaVersion int                      `yaml:"schema_version" json:"schema_version"`
	Workspace     IntegrationPlanWorkspace `yaml:"workspace" json:"workspace"`
	Target        QueueTarget              `yaml:"target" json:"target"`
	Registry      QueueRegistryReadback    `yaml:"registry" json:"registry"`
	QueueID       string                   `yaml:"queue_id" json:"queue_id"`
	Revision      int                      `yaml:"revision" json:"revision"`
	Current       *QueueCurrent            `yaml:"current" json:"current"`
	Pending       []QueuePending           `yaml:"pending" json:"pending"`
	Reasons       []QueueReason            `yaml:"reasons" json:"reasons"`
}
