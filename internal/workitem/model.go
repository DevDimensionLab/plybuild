// Package workitem owns registered Task, Epic and worktree data and their
// pure relationship behavior. Storage, history and live effects belong to callers.
package workitem

type ProjectID string
type RepoID string

type EpicID string
type TaskID string
type WorktreeID string
type WorktreeOperationID string

type WorkItemState string

const (
	WorkItemUnbound                WorkItemState = "unbound"
	WorkItemCreating               WorkItemState = "creating"
	WorkItemReady                  WorkItemState = "worktree_ready"
	WorkItemRetired                WorkItemState = "retired"
	WorkItemReconciliationRequired WorkItemState = "reconciliation_required"
)

type EpicWorktreeBinding struct {
	ID           WorktreeID `yaml:"id" json:"id"`
	OwnerKind    string     `yaml:"owner_kind" json:"owner_kind"`
	OwnerID      EpicID     `yaml:"owner_id" json:"owner_id"`
	Origin       string     `yaml:"origin" json:"origin"`
	Locator      string     `yaml:"locator" json:"locator"`
	Ref          string     `yaml:"ref" json:"ref"`
	OID          string     `yaml:"oid" json:"oid"`
	Tree         string     `yaml:"tree" json:"tree"`
	GitCommonDir string     `yaml:"git_common_dir" json:"git_common_dir"`
}

type EpicRepoBinding struct {
	RepoID       RepoID              `yaml:"repo_id" json:"repo_id"`
	GitCommonDir string              `yaml:"git_common_dir" json:"git_common_dir"`
	Worktree     EpicWorktreeBinding `yaml:"worktree" json:"worktree"`
}

type EpicRecord struct {
	ID           EpicID            `yaml:"id" json:"epic_id"`
	Title        string            `yaml:"title" json:"title"`
	ProjectID    ProjectID         `yaml:"project_id" json:"project_id"`
	RepoBindings []EpicRepoBinding `yaml:"repo_bindings" json:"repo_bindings"`
}

type TaskWorktreeBinding struct {
	ID               WorktreeID          `yaml:"id" json:"id"`
	OwnerKind        string              `yaml:"owner_kind" json:"owner_kind"`
	OwnerID          TaskID              `yaml:"owner_id" json:"owner_id"`
	Origin           string              `yaml:"origin" json:"origin"`
	Locator          string              `yaml:"locator" json:"locator"`
	Ref              string              `yaml:"ref" json:"ref"`
	OID              string              `yaml:"oid" json:"oid"`
	Tree             string              `yaml:"tree" json:"tree"`
	GitCommonDir     string              `yaml:"git_common_dir" json:"git_common_dir"`
	ParentEpicID     EpicID              `yaml:"parent_epic_id" json:"parent_epic_id"`
	ParentWorktreeID WorktreeID          `yaml:"parent_worktree_id" json:"parent_worktree_id"`
	ParentRef        string              `yaml:"parent_ref" json:"parent_ref"`
	ParentOID        string              `yaml:"parent_oid" json:"parent_oid"`
	ParentTree       string              `yaml:"parent_tree" json:"parent_tree"`
	OperationID      WorktreeOperationID `yaml:"operation_id" json:"operation_id"`
	IntentDigest     string              `yaml:"intent_digest" json:"intent_digest"`
}

type TaskRecord struct {
	ID            TaskID               `yaml:"id" json:"task_id"`
	Title         string               `yaml:"title" json:"title"`
	Description   string               `yaml:"description" json:"description"`
	ParentEpicID  EpicID               `yaml:"parent_epic_id" json:"parent_epic_id"`
	ProjectID     ProjectID            `yaml:"project_id" json:"project_id"`
	RepoID        RepoID               `yaml:"repo_id" json:"repo_id"`
	GitCommonDir  string               `yaml:"git_common_dir" json:"git_common_dir"`
	WorktreeState WorkItemState        `yaml:"worktree_state" json:"worktree_state"`
	Worktree      *TaskWorktreeBinding `yaml:"worktree" json:"worktree"`
}
