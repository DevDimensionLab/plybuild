package workspace

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type IntegrationStatusEntry struct {
	RecordKind         string  `yaml:"record_kind" json:"record_kind"`
	PathBase64         string  `yaml:"path_base64" json:"path_base64"`
	OriginalPathBase64 *string `yaml:"original_path_base64" json:"original_path_base64"`
	IndexState         *string `yaml:"index_state" json:"index_state"`
	WorktreeState      *string `yaml:"worktree_state" json:"worktree_state"`
	SubmoduleState     *string `yaml:"submodule_state" json:"submodule_state"`
}

type IntegrationInventoryEntry struct {
	Locator  string `yaml:"locator" json:"locator"`
	Ref      string `yaml:"ref" json:"ref"`
	OID      string `yaml:"oid" json:"oid"`
	Locked   bool   `yaml:"locked" json:"locked"`
	Prunable bool   `yaml:"prunable" json:"prunable"`
	Bare     bool   `yaml:"bare" json:"bare"`
	Detached bool   `yaml:"detached" json:"detached"`
}

type IntegrationReflogEntry struct {
	Ordinal  int    `yaml:"ordinal" json:"ordinal"`
	OID      string `yaml:"oid" json:"oid"`
	Selector string `yaml:"selector" json:"selector"`
	Action   string `yaml:"action" json:"action"`
}

type PlanWorktreeObservation struct {
	WorktreeLocator string                   `yaml:"worktree_locator" json:"worktree_locator"`
	Ref             string                   `yaml:"ref" json:"ref"`
	OID             string                   `yaml:"oid" json:"oid"`
	Tree            string                   `yaml:"tree" json:"tree"`
	GitCommonDir    string                   `yaml:"git_common_dir" json:"git_common_dir"`
	ObjectFormat    string                   `yaml:"object_format" json:"object_format"`
	RefFormat       string                   `yaml:"ref_format" json:"ref_format"`
	Symbolic        bool                     `yaml:"symbolic" json:"symbolic"`
	Clean           bool                     `yaml:"clean" json:"clean"`
	StatusEntries   []IntegrationStatusEntry `yaml:"status_entries" json:"status_entries"`
	InProgress      []string                 `yaml:"in_progress" json:"in_progress"`
}

type IntegrationObservation struct {
	ObservedAtUTC    string                      `yaml:"observed_at_utc" json:"observed_at_utc"`
	WorktreeLocator  string                      `yaml:"worktree_locator" json:"worktree_locator"`
	Ref              string                      `yaml:"ref" json:"ref"`
	OID              string                      `yaml:"oid" json:"oid"`
	Tree             string                      `yaml:"tree" json:"tree"`
	GitCommonDir     string                      `yaml:"git_common_dir" json:"git_common_dir"`
	ObjectFormat     string                      `yaml:"object_format" json:"object_format"`
	RefFormat        string                      `yaml:"ref_format" json:"ref_format"`
	Symbolic         bool                        `yaml:"symbolic" json:"symbolic"`
	Clean            bool                        `yaml:"clean" json:"clean"`
	StatusEntries    []IntegrationStatusEntry    `yaml:"status_entries" json:"status_entries"`
	InProgress       []string                    `yaml:"in_progress" json:"in_progress"`
	InventoryEntries []IntegrationInventoryEntry `yaml:"inventory_entries" json:"inventory_entries"`
}

type IntegrationAllowedEffect struct {
	Kind              string `yaml:"kind" json:"kind"`
	ParentRef         string `yaml:"parent_ref" json:"parent_ref"`
	ExpectedParentOID string `yaml:"expected_parent_oid" json:"expected_parent_oid"`
	ResultOID         string `yaml:"result_oid" json:"result_oid"`
	MaxOccurrences    int    `yaml:"max_occurrences" json:"max_occurrences"`
}

type WorkspaceTaskIntegrationPlan struct {
	Kind                  string                      `yaml:"kind" json:"kind"`
	SchemaVersion         int                         `yaml:"schema_version" json:"schema_version"`
	Format                string                      `yaml:"format" json:"format"`
	FormatVersion         int                         `yaml:"format_version" json:"format_version"`
	Canonicalization      string                      `yaml:"canonicalization" json:"canonicalization"`
	Workspace             IntegrationPlanWorkspace    `yaml:"workspace" json:"workspace"`
	Project               IntegrationPlanProject      `yaml:"project" json:"project"`
	Repository            IntegrationPlanRepository   `yaml:"repository" json:"repository"`
	Epic                  IntegrationPlanEpic         `yaml:"epic" json:"epic"`
	Task                  IntegrationPlanTask         `yaml:"task" json:"task"`
	TaskResult            IntegrationPlanTaskResult   `yaml:"task_result" json:"task_result"`
	HumanQA               IntegrationPlanHumanQA      `yaml:"human_qa" json:"human_qa"`
	RetryAfterResultID    *IntegrationResultID        `yaml:"retry_after_result_id" json:"retry_after_result_id"`
	WorkItemsSHA256       string                      `yaml:"work_items_sha256" json:"work_items_sha256"`
	StoreTransition       string                      `yaml:"store_transition" json:"store_transition"`
	ObservedSource        PlanWorktreeObservation     `yaml:"observed_source" json:"observed_source"`
	ObservedParent        PlanWorktreeObservation     `yaml:"observed_parent" json:"observed_parent"`
	ObservedInventory     []IntegrationInventoryEntry `yaml:"observed_inventory" json:"observed_inventory"`
	ObservedReflog        []IntegrationReflogEntry    `yaml:"observed_reflog" json:"observed_reflog"`
	TechnicalGateReady    bool                        `yaml:"technical_gate_ready" json:"technical_gate_ready"`
	HumanQAReady          bool                        `yaml:"human_qa_ready" json:"human_qa_ready"`
	AncestryReady         bool                        `yaml:"ancestry_ready" json:"ancestry_ready"`
	Readiness             string                      `yaml:"readiness" json:"readiness"`
	Reasons               []string                    `yaml:"reasons" json:"reasons"`
	Effect                IntegrationPlanEffect       `yaml:"effect" json:"effect"`
	TaskSpecGuard         *TaskSpecRelevance          `yaml:"task_spec_guard,omitempty" json:"task_spec_guard"`
	DeliveryAuthorization *DeliveryAuthorization      `yaml:"delivery_authorization,omitempty" json:"delivery_authorization,omitempty"`
}

type IntegrationPlanWorkspace struct {
	Root         string `yaml:"root" json:"root"`
	MarkerSHA256 string `yaml:"marker_sha256" json:"marker_sha256"`
}
type IntegrationPlanProject struct {
	ProjectID ProjectID `yaml:"project_id" json:"project_id"`
}
type IntegrationPlanRepository struct {
	RepoID            RepoID `yaml:"repo_id" json:"repo_id"`
	RegisteredLocator string `yaml:"registered_locator" json:"registered_locator"`
	GitCommonDir      string `yaml:"git_common_dir" json:"git_common_dir"`
	GitVersion        string `yaml:"git_version" json:"git_version"`
	ObjectFormat      string `yaml:"object_format" json:"object_format"`
	RefFormat         string `yaml:"ref_format" json:"ref_format"`
	Shallow           bool   `yaml:"shallow" json:"shallow"`
	PartialClone      bool   `yaml:"partial_clone" json:"partial_clone"`
	SparseCheckout    bool   `yaml:"sparse_checkout" json:"sparse_checkout"`
}
type IntegrationPlanEpic struct {
	EpicID             EpicID     `yaml:"epic_id" json:"epic_id"`
	ParentWorktreeID   WorktreeID `yaml:"parent_worktree_id" json:"parent_worktree_id"`
	ParentLocator      string     `yaml:"parent_locator" json:"parent_locator"`
	ParentRef          string     `yaml:"parent_ref" json:"parent_ref"`
	ExpectedParentOID  string     `yaml:"expected_parent_oid" json:"expected_parent_oid"`
	ExpectedParentTree string     `yaml:"expected_parent_tree" json:"expected_parent_tree"`
}
type IntegrationPlanTask struct {
	TaskID         TaskID     `yaml:"task_id" json:"task_id"`
	TaskWorktreeID WorktreeID `yaml:"task_worktree_id" json:"task_worktree_id"`
	SourceLocator  string     `yaml:"source_locator" json:"source_locator"`
	SourceRef      string     `yaml:"source_ref" json:"source_ref"`
	ResultOID      string     `yaml:"result_oid" json:"result_oid"`
	ResultTree     string     `yaml:"result_tree" json:"result_tree"`
}
type IntegrationPlanTaskResult struct {
	ID            TaskResultID `yaml:"id" json:"id"`
	DraftSHA256   string       `yaml:"draft_sha256" json:"draft_sha256"`
	TechnicalGate string       `yaml:"technical_gate" json:"technical_gate"`
}
type IntegrationPlanHumanQA struct {
	ID          HumanQARecordID `yaml:"id" json:"id"`
	DraftSHA256 string          `yaml:"draft_sha256" json:"draft_sha256"`
	Outcome     string          `yaml:"outcome" json:"outcome"`
}
type IntegrationPlanEffect struct {
	Kind              string   `yaml:"kind" json:"kind"`
	ParentRef         string   `yaml:"parent_ref" json:"parent_ref"`
	ExpectedParentOID string   `yaml:"expected_parent_oid" json:"expected_parent_oid"`
	ResultOID         string   `yaml:"result_oid" json:"result_oid"`
	MaxOccurrences    int      `yaml:"max_occurrences" json:"max_occurrences"`
	Argv              []string `yaml:"argv" json:"argv"`
}

type IntegrationAuthority struct {
	ID                 IntegrationAuthorityID       `yaml:"id" json:"id"`
	Mode               string                       `yaml:"mode" json:"mode"`
	CreatedAtUTC       string                       `yaml:"created_at_utc" json:"created_at_utc"`
	PlanSHA256         string                       `yaml:"plan_sha256" json:"plan_sha256"`
	Plan               WorkspaceTaskIntegrationPlan `yaml:"plan" json:"plan"`
	RetryAfterResultID *IntegrationResultID         `yaml:"retry_after_result_id" json:"retry_after_result_id"`
	TaskID             TaskID                       `yaml:"task_id" json:"task_id"`
	TaskResultID       TaskResultID                 `yaml:"task_result_id" json:"task_result_id"`
	HumanQARecordID    HumanQARecordID              `yaml:"human_qa_record_id" json:"human_qa_record_id"`
	AllowedEffect      IntegrationAllowedEffect     `yaml:"allowed_effect" json:"allowed_effect"`
	DeliveryOwner      *DeliveryIntegrationOwner    `yaml:"delivery_owner,omitempty" json:"delivery_owner,omitempty"`
}

// DeliveryIntegrationOwner records the existing delivery mandate's provenance.
// It does not replace the exact TaskResult, actual human QA, plan or effect guard.
type DeliveryIntegrationOwner struct {
	RunID         string `yaml:"run_id" json:"run_id"`
	RequestSHA256 string `yaml:"request_sha256" json:"request_sha256"`
	ActorClaim    string `yaml:"actor_claim" json:"actor_claim"`
	PreparationID string `yaml:"preparation_id" json:"preparation_id"`
}

type IntegrationIntent struct {
	ID                    IntegrationIntentID      `yaml:"id" json:"id"`
	IntentSHA256          string                   `yaml:"intent_sha256" json:"intent_sha256"`
	AuthorityID           IntegrationAuthorityID   `yaml:"authority_id" json:"authority_id"`
	State                 string                   `yaml:"state" json:"state"`
	WorkspaceRoot         string                   `yaml:"workspace_root" json:"workspace_root"`
	WorkspaceMarkerSHA256 string                   `yaml:"workspace_marker_sha256" json:"workspace_marker_sha256"`
	ProjectID             ProjectID                `yaml:"project_id" json:"project_id"`
	RepoID                RepoID                   `yaml:"repo_id" json:"repo_id"`
	RegisteredRepoLocator string                   `yaml:"registered_repo_locator" json:"registered_repo_locator"`
	GitCommonDir          string                   `yaml:"git_common_dir" json:"git_common_dir"`
	ObjectFormat          string                   `yaml:"object_format" json:"object_format"`
	RefFormat             string                   `yaml:"ref_format" json:"ref_format"`
	EpicID                EpicID                   `yaml:"epic_id" json:"epic_id"`
	ParentWorktreeID      WorktreeID               `yaml:"parent_worktree_id" json:"parent_worktree_id"`
	ParentLocator         string                   `yaml:"parent_locator" json:"parent_locator"`
	ParentRef             string                   `yaml:"parent_ref" json:"parent_ref"`
	ExpectedParentOID     string                   `yaml:"expected_parent_oid" json:"expected_parent_oid"`
	ExpectedParentTree    string                   `yaml:"expected_parent_tree" json:"expected_parent_tree"`
	TaskID                TaskID                   `yaml:"task_id" json:"task_id"`
	TaskWorktreeID        WorktreeID               `yaml:"task_worktree_id" json:"task_worktree_id"`
	SourceLocator         string                   `yaml:"source_locator" json:"source_locator"`
	SourceRef             string                   `yaml:"source_ref" json:"source_ref"`
	ResultOID             string                   `yaml:"result_oid" json:"result_oid"`
	ResultTree            string                   `yaml:"result_tree" json:"result_tree"`
	TaskResultID          TaskResultID             `yaml:"task_result_id" json:"task_result_id"`
	HumanQARecordID       HumanQARecordID          `yaml:"human_qa_record_id" json:"human_qa_record_id"`
	PlanSHA256            string                   `yaml:"plan_sha256" json:"plan_sha256"`
	Effect                IntegrationAllowedEffect `yaml:"effect" json:"effect"`
	MaxAttempts           int                      `yaml:"max_attempts" json:"max_attempts"`
}

type IntegrationCommandStream struct {
	SizeBytes     int64  `yaml:"size_bytes" json:"size_bytes"`
	CapturedBytes int64  `yaml:"captured_bytes" json:"captured_bytes"`
	Truncated     bool   `yaml:"truncated" json:"truncated"`
	SHA256        string `yaml:"sha256" json:"sha256"`
	Base64        string `yaml:"base64" json:"base64"`
}
type IntegrationCommandEvidence struct {
	Launched    bool                     `yaml:"launched" json:"launched"`
	ExitCode    *int                     `yaml:"exit_code" json:"exit_code"`
	LaunchError *string                  `yaml:"launch_error" json:"launch_error"`
	Stdout      IntegrationCommandStream `yaml:"stdout" json:"stdout"`
	Stderr      IntegrationCommandStream `yaml:"stderr" json:"stderr"`
}
type IntegrationNextAction struct {
	Kind   string   `yaml:"kind" json:"kind"`
	Reason string   `yaml:"reason" json:"reason"`
	Argv   []string `yaml:"argv" json:"argv"`
}

type IntegrationAttempt struct {
	ID              IntegrationAttemptID   `yaml:"id" json:"id"`
	AuthorityID     IntegrationAuthorityID `yaml:"authority_id" json:"authority_id"`
	IntentID        IntegrationIntentID    `yaml:"intent_id" json:"intent_id"`
	Ordinal         int                    `yaml:"ordinal" json:"ordinal"`
	State           string                 `yaml:"state" json:"state"`
	ResultID        *IntegrationResultID   `yaml:"result_id" json:"result_id"`
	PreparedAtUTC   string                 `yaml:"prepared_at_utc" json:"prepared_at_utc"`
	CommandIdentity string                 `yaml:"command_identity" json:"command_identity"`
	ReflogAction    string                 `yaml:"reflog_action" json:"reflog_action"`
	PreObservation  IntegrationObservation `yaml:"pre_observation" json:"pre_observation"`
}

type IntegrationResult struct {
	ID                 IntegrationResultID        `yaml:"id" json:"id"`
	AuthorityID        IntegrationAuthorityID     `yaml:"authority_id" json:"authority_id"`
	IntentID           IntegrationIntentID        `yaml:"intent_id" json:"intent_id"`
	AttemptID          *IntegrationAttemptID      `yaml:"attempt_id" json:"attempt_id"`
	RetryAfterResultID *IntegrationResultID       `yaml:"retry_after_result_id" json:"retry_after_result_id"`
	Outcome            string                     `yaml:"outcome" json:"outcome"`
	GitChanged         *bool                      `yaml:"git_changed" json:"git_changed"`
	RecordedAtUTC      string                     `yaml:"recorded_at_utc" json:"recorded_at_utc"`
	BeforeObservation  IntegrationObservation     `yaml:"before_observation" json:"before_observation"`
	AfterObservation   IntegrationObservation     `yaml:"after_observation" json:"after_observation"`
	Command            IntegrationCommandEvidence `yaml:"command" json:"command"`
	Reflog             []IntegrationReflogEntry   `yaml:"reflog" json:"reflog"`
	Recovered          bool                       `yaml:"recovered" json:"recovered"`
	RecoveryStatus     string                     `yaml:"recovery_status" json:"recovery_status"`
	NextAction         IntegrationNextAction      `yaml:"next_action" json:"next_action"`
}

type IntegrationRepositoryObservation struct {
	GitVersion, ObjectFormat, RefFormat   string
	Shallow, PartialClone, SparseCheckout bool
	Inventory                             []IntegrationInventoryEntry
}
type IntegrationWorktreeObservation struct {
	Locator, Ref, OID, Tree, GitCommonDir, ObjectFormat, RefFormat string
	Symbolic, Clean                                                bool
	StatusEntries                                                  []IntegrationStatusEntry
	InProgress                                                     []string
}
type IntegrationMergeInput struct {
	Repository                RepoRecord
	ParentWorktree, ResultOID string
	AttemptID                 IntegrationAttemptID
}

type TaskIntegrationGit interface {
	ObserveIntegrationRepository(RepoRecord) (IntegrationRepositoryObservation, error)
	ObserveIntegrationWorktree(string, string) (IntegrationWorktreeObservation, error)
	CheckAncestor(RepoRecord, string, string) (bool, error)
	ObserveParentReflog(string, string, int) ([]IntegrationReflogEntry, error)
	MergeFastForward(IntegrationMergeInput) GitCommandOutcome
}

type TaskIntegrationInput struct {
	TaskID                               TaskID
	TaskResultID                         TaskResultID
	HumanQARecordID                      HumanQARecordID
	ExpectedResultOID, ExpectedParentOID string
	RetryAfterResultID                   *IntegrationResultID
	Apply                                bool
	Confirmation                         string
	DeliveryOwner                        *DeliveryIntegrationOwner
	DeliveryAuthorization                *DeliveryAuthorization
}

type TaskIntegrationResult struct {
	Readback WorkspaceTaskIntegrationReadback
}

func ParseTaskIntegrationInput(taskID, resultID, qaID, resultOID, parentOID, retryAfter string, apply bool, confirmation string) (TaskIntegrationInput, error) {
	task, err := ParseTaskID(taskID)
	if err != nil {
		return TaskIntegrationInput{}, err
	}
	result, err := ParseTaskResultID(resultID)
	if err != nil {
		return TaskIntegrationInput{}, err
	}
	qa, err := ParseHumanQARecordID(qaID)
	if err != nil {
		return TaskIntegrationInput{}, err
	}
	input := TaskIntegrationInput{TaskID: task, TaskResultID: result, HumanQARecordID: qa, ExpectedResultOID: resultOID, ExpectedParentOID: parentOID, Apply: apply, Confirmation: confirmation}
	if retryAfter != "" {
		id := IntegrationResultID(retryAfter)
		input.RetryAfterResultID = &id
	}
	if err := validateIntegrationInput(input); err != nil {
		return TaskIntegrationInput{}, err
	}
	if confirmation != "" && !digestPattern.MatchString(confirmation) {
		return TaskIntegrationInput{}, WorkInvalidArguments("--confirm must be a sha256 digest")
	}
	return input, nil
}

func CheckTaskIntegration(dependencies Dependencies, input TaskIntegrationInput) (TaskIntegrationResult, error) {
	input.Apply = false
	if err := validateIntegrationInput(input); err != nil {
		return TaskIntegrationResult{}, err
	}
	if err := requireIntegrationDependencies(dependencies); err != nil {
		return TaskIntegrationResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return TaskIntegrationResult{}, err
	}
	projects, repos, err := dependencies.Projects.Snapshot(root)
	if err != nil {
		return TaskIntegrationResult{}, mapProjectWorkError(err)
	}
	registry, err := dependencies.WorkItems.Snapshot(root)
	if err != nil {
		return TaskIntegrationResult{}, err
	}
	if input.DeliveryOwner != nil || input.DeliveryAuthorization != nil {
		if recovered, found, err := checkPersistedDeliveryIntegration(dependencies, root, ProjectSnapshot{Projects: projects, Repos: repos}, registry, input); found || err != nil {
			return recovered, err
		}
	}
	if input.RetryAfterResultID == nil {
		if authority, intent, attempt, result := unresolvedIntegrationLeaf(registry, input); result != nil {
			if !contentTypedEqual(authority.DeliveryOwner, input.DeliveryOwner) || !contentTypedEqual(authority.Plan.DeliveryAuthorization, input.DeliveryAuthorization) {
				return TaskIntegrationResult{}, workError(ErrorTaskIntegrationConflict, "integration delivery owner differs", nil)
			}
			ctx, err := contextFromPlan(registry, authority.Plan)
			if err != nil {
				return TaskIntegrationResult{}, err
			}
			if err := revalidateIntegrationEvidence(dependencies, ctx.Result, ctx.QA); err != nil {
				return TaskIntegrationResult{}, err
			}
			ctx.TaskSpecRelevance = taskResultSpecRelevance(dependencies, root, ProjectSnapshot{Projects: projects, Repos: repos}, registry, ctx.Task, ctx.Result.ID)
			ctx.PersistedBefore = integrationPersistedView(registry, authority, intent, attempt, result)
			return TaskIntegrationResult{Readback: newIntegrationReadback(root, registry, ctx, authority.Plan, authority.PlanSHA256, authority, intent, attempt, result)}, nil
		}
	}
	plan, digest, ctx, err := buildIntegrationPlan(dependencies, root, ProjectSnapshot{Projects: projects, Repos: repos}, registry, input)
	if err != nil {
		return TaskIntegrationResult{}, err
	}
	return TaskIntegrationResult{Readback: newIntegrationReadback(root, registry, ctx, plan, digest, nil, nil, nil, nil)}, nil
}

func ApplyTaskIntegration(dependencies Dependencies, input TaskIntegrationInput) (TaskIntegrationResult, error) {
	input.Apply = true
	if err := validateIntegrationInput(input); err != nil {
		return TaskIntegrationResult{}, err
	}
	if input.Confirmation == "" || !digestPattern.MatchString(input.Confirmation) {
		return TaskIntegrationResult{}, WorkInvalidArguments("--confirm must be a sha256 digest")
	}
	if err := requireIntegrationDependencies(dependencies); err != nil {
		return TaskIntegrationResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return TaskIntegrationResult{}, err
	}
	var output WorkspaceTaskIntegrationReadback
	err = dependencies.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			if authority := findAuthorityByPlan(registry, input.Confirmation, input); authority != nil {
				intent := findIntentForAuthority(registry, authority.ID)
				if intent == nil {
					return workError(ErrorTaskIntegrationConflict, "integration authority has no exact intent", nil)
				}
				attempt := findAttemptForAuthority(registry, authority.ID)
				stored := findIntegrationResultForAuthority(registry, authority.ID)
				ctx, err := contextFromPlan(registry, authority.Plan)
				if err != nil {
					return err
				}
				ctx.TaskSpecRelevance = taskResultSpecRelevance(dependencies, root, projects, registry, ctx.Task, ctx.Result.ID)
				ctx.PersistedBefore = integrationPersistedView(registry, authority, intent, attempt, stored)
				if stored != nil {
					output = newIntegrationReadback(root, registry, ctx, authority.Plan, authority.PlanSHA256, authority, intent, attempt, stored)
					return nil
				}
				if attempt != nil {
					return reconcileIntegrationAttempt(dependencies, session, &registry, root, ctx, *authority, *intent, *attempt, &output)
				}
				if err := revalidateIntegrationEvidence(dependencies, ctx.Result, ctx.QA); err != nil {
					return err
				}
				ctx.PreconditionReason = authorityProjectBindingReason(projects, authority.Plan)
				ctx.TaskSpecRelevance = taskResultSpecRelevance(dependencies, root, projects, registry, ctx.Task, ctx.Result.ID)
				if ctx.TaskSpecRelevance != nil && ctx.TaskSpecRelevance.Relevance != "current" {
					ctx.PreconditionReason = "task_spec_result_not_current"
				}
				return continueIntegrationAuthority(dependencies, session, &registry, root, ctx, *authority, *intent, &output, true)
			}
			plan, digest, ctx, err := buildIntegrationPlan(dependencies, root, projects, registry, input)
			if err != nil {
				return err
			}
			if digest != input.Confirmation {
				return workError(ErrorTaskIntegrationConfirmationMismatch, fmt.Sprintf("fresh plan digest is %s", digest), nil)
			}
			ctx.PersistedBefore = integrationPersistedView(registry, nil, nil, nil, nil)
			if plan.Readiness != "ready" && plan.Readiness != "already_integrated" {
				output = newIntegrationReadback(root, registry, ctx, plan, digest, nil, nil, nil, nil)
				return nil
			}
			authorityID, err := dependencies.TaskLifecycleIDs.NewIntegrationAuthorityID()
			if err != nil {
				return err
			}
			intentID, err := dependencies.TaskLifecycleIDs.NewIntegrationIntentID()
			if err != nil {
				return err
			}
			now := dependencies.WorkClock.Now().UTC().Format(time.RFC3339Nano)
			authority := IntegrationAuthority{ID: authorityID, Mode: "human_cli_start", CreatedAtUTC: now, PlanSHA256: digest, Plan: plan, RetryAfterResultID: input.RetryAfterResultID, TaskID: input.TaskID, TaskResultID: input.TaskResultID, HumanQARecordID: input.HumanQARecordID, AllowedEffect: IntegrationAllowedEffect{Kind: "local_ff_only", ParentRef: plan.Epic.ParentRef, ExpectedParentOID: plan.Epic.ExpectedParentOID, ResultOID: plan.Task.ResultOID, MaxOccurrences: 1}}
			if input.DeliveryOwner != nil {
				authority.Mode = "delivery_owner_after_human_pass"
				if input.DeliveryAuthorization != nil && input.DeliveryAuthorization.HumanIntegration != nil {
					authority.Mode = "human_integration_plan"
				}
				copy := *input.DeliveryOwner
				authority.DeliveryOwner = &copy
			}
			intent := intentFromPlan(intentID, authorityID, plan, digest)
			intent.IntentSHA256 = integrationIntentDigest(intent)
			registry.IntegrationAuthorities = append(registry.IntegrationAuthorities, authority)
			registry.IntegrationIntents = append(registry.IntegrationIntents, intent)
			sortWorkRegistry(&registry)
			published, err := publishWorkItemRegistryRecover(session, registry)
			if err != nil {
				return err
			}
			registry = published
			return continueIntegrationAuthority(dependencies, session, &registry, root, ctx, authority, intent, &output, false)
		})
	})
	if err != nil {
		return TaskIntegrationResult{Readback: output}, mapMutationError(err)
	}
	return TaskIntegrationResult{Readback: output}, nil
}

type integrationContext struct {
	TaskSpecRelevance         *TaskSpecRelevance
	Task                      TaskRecord
	Epic                      EpicRecord
	Parent                    EpicRepoBinding
	Result                    TaskResultRecord
	QA                        TaskHumanQARecord
	Repo                      RepoRecord
	Repository                IntegrationRepositoryObservation
	Source, ParentObservation IntegrationWorktreeObservation
	Reflog                    []IntegrationReflogEntry
	SourceObserved            bool
	ParentObserved            bool
	RepositoryObserved        bool
	ReflogObserved            bool
	PersistedBefore           canonicaljson.Value
	PreconditionReason        string
}

func validateIntegrationInput(input TaskIntegrationInput) error {
	if input.DeliveryOwner != nil && !validDeliveryIntegrationOwner(input.DeliveryOwner) {
		return WorkInvalidArguments("invalid delivery owner provenance")
	}
	if input.DeliveryAuthorization != nil && ValidateDeliveryAuthorization(*input.DeliveryAuthorization) != nil {
		return WorkInvalidArguments("invalid native delivery authorization")
	}
	if input.DeliveryAuthorization != nil && input.DeliveryOwner != nil && (input.DeliveryOwner.RunID != input.DeliveryAuthorization.RunID || input.DeliveryOwner.RequestSHA256 != input.DeliveryAuthorization.RequestSHA256) {
		return WorkInvalidArguments("delivery owner differs from native authorized run")
	}
	if _, err := ParseTaskID(string(input.TaskID)); err != nil {
		return err
	}
	if _, err := ParseTaskResultID(string(input.TaskResultID)); err != nil {
		return err
	}
	if _, err := ParseHumanQARecordID(string(input.HumanQARecordID)); err != nil {
		return err
	}
	if !validOIDText(input.ExpectedResultOID) || !validOIDText(input.ExpectedParentOID) {
		return WorkInvalidArguments("expected OIDs must be full lowercase object IDs")
	}
	if input.RetryAfterResultID != nil && !lifecycleIDPatterns["integration result"].MatchString(string(*input.RetryAfterResultID)) {
		return WorkInvalidArguments("invalid prior integration result ID")
	}
	return nil
}
func requireIntegrationDependencies(d Dependencies) error {
	if d.Files == nil || d.Projects == nil || d.ProjectLocks == nil || d.WorkItems == nil || d.TaskLifecycleIDs == nil || d.WorkClock == nil || d.HandoffEvidence == nil || d.IntegrationGit == nil {
		return workError(ErrorWorkIO, "integration dependencies are required", nil)
	}
	return nil
}

func buildIntegrationPlan(d Dependencies, root string, projects ProjectSnapshot, registry WorkItemRegistry, input TaskIntegrationInput) (WorkspaceTaskIntegrationPlan, string, integrationContext, error) {
	if registry.FormatVersion < 2 || registry.FormatVersion > 4 {
		return WorkspaceTaskIntegrationPlan{}, "", integrationContext{}, workError(ErrorTaskIntegrationBlocked, "work-item store must be format 2, 3 or 4", nil)
	}
	task, _ := findTask(registry, input.TaskID)
	if task == nil || task.Worktree == nil || task.WorktreeState != WorkItemReady {
		return WorkspaceTaskIntegrationPlan{}, "", integrationContext{}, workError(ErrorTaskIntegrationBlocked, "Task does not have a ready worktree", nil)
	}
	epic, _ := findEpic(registry, task.ParentEpicID)
	if epic == nil {
		return WorkspaceTaskIntegrationPlan{}, "", integrationContext{}, workError(ErrorTaskIntegrationConflict, "Task parent Epic is missing", nil)
	}
	parent, _ := findEpicRepo(*epic, task.RepoID)
	if parent == nil {
		return WorkspaceTaskIntegrationPlan{}, "", integrationContext{}, workError(ErrorTaskIntegrationConflict, "Task parent repository binding is missing", nil)
	}
	tr := findTaskResult(registry, input.TaskResultID)
	qa := findHumanQA(registry, input.HumanQARecordID)
	if tr == nil || qa == nil {
		return WorkspaceTaskIntegrationPlan{}, "", integrationContext{}, workError(ErrorTaskIntegrationBlocked, "selected Task result or human QA record is missing", nil)
	}
	_, repo, err := projectAndRepo(projects, task.ProjectID, task.RepoID)
	if err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", integrationContext{}, err
	}
	ctx := integrationContext{Task: *task, Epic: *epic, Parent: *parent, Result: *tr, QA: *qa, Repo: repo}
	if err := revalidateIntegrationEvidence(d, *tr, *qa); err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", ctx, err
	}
	target := QueueTarget{ProjectID: task.ProjectID, RepoID: task.RepoID, EpicID: epic.ID, GitCommonDir: task.GitCommonDir, ParentWorktreeID: parent.Worktree.ID, ParentLocator: parent.Worktree.Locator, ParentRef: parent.Worktree.Ref}
	if err := validateIntegrationDelivery(d, root, registry, *tr, qa.ID, input.DeliveryAuthorization, target); err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", ctx, err
	}
	repository, err := d.IntegrationGit.ObserveIntegrationRepository(repo)
	if err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", ctx, workError(ErrorTaskIntegrationUnsupportedRepository, err.Error(), err)
	}
	ctx.Repository = repository
	ctx.RepositoryObserved = true
	source, err := d.IntegrationGit.ObserveIntegrationWorktree(task.Worktree.Locator, task.Worktree.Ref)
	if err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", ctx, workError(ErrorTaskIntegrationBlocked, err.Error(), err)
	}
	ctx.Source = source
	ctx.SourceObserved = true
	parentObs, err := d.IntegrationGit.ObserveIntegrationWorktree(parent.Worktree.Locator, parent.Worktree.Ref)
	if err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", ctx, workError(ErrorTaskIntegrationBlocked, err.Error(), err)
	}
	ctx.ParentObservation = parentObs
	ctx.ParentObserved = true
	ancestor, err := d.IntegrationGit.CheckAncestor(repo, input.ExpectedParentOID, input.ExpectedResultOID)
	if err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", ctx, workError(ErrorTaskIntegrationConflict, err.Error(), err)
	}
	reflog, err := d.IntegrationGit.ObserveParentReflog(parent.Worktree.Locator, parent.Worktree.Ref, 2)
	if err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", ctx, workError(ErrorTaskIntegrationBlocked, err.Error(), err)
	}
	ctx.Reflog = reflog
	ctx.ReflogObserved = true
	reasons := []string{}
	readiness := "ready"
	technical := tr.TechnicalGate == "passed" || tr.TechnicalGate == "good_enough_with_known_debt"
	human := qa.Outcome == "pass" && qa.TaskResultID == tr.ID && qa.ResultOID == tr.ResultOID && qa.ResultTree == tr.ResultTree
	conflict := func(reason string) { reasons = append(reasons, reason); readiness = "conflict" }
	block := func(reason string) {
		reasons = append(reasons, reason)
		if readiness != "conflict" {
			readiness = "blocked"
		}
	}
	if registry.FormatVersion == 4 {
		base := currentEpicBase(registry, *epic, *parent)
		if task.Worktree.ParentOID != base.OID || task.Worktree.ParentTree != base.Tree {
			block("epic_base_changed")
		}
		if pendingBaseUpdate(registry, epic.ID, task.RepoID) {
			block("epic_base_update_pending")
		}
	}
	if tr.TaskID != task.ID || tr.TaskWorktreeID != task.Worktree.ID || tr.ResultOID != input.ExpectedResultOID || tr.ResultTree == "" || tr.SourceLocator != task.Worktree.Locator || tr.SourceRef != task.Worktree.Ref {
		conflict("task_result_binding_mismatch")
	}
	if !technical {
		block("technical_gate_not_ready")
	}
	if !human {
		block("human_qa_not_ready")
	}
	if input.ExpectedParentOID != task.Worktree.ParentOID || input.ExpectedResultOID != tr.ResultOID {
		conflict("expected_oid_mismatch")
	}
	if ordinaryProductRef(parent.Worktree.Ref) && (input.DeliveryAuthorization == nil || input.DeliveryAuthorization.Agreement.Mode != DeliveryLocalBranch) {
		block("ordinary_main_parent_forbidden")
	}
	if source.OID != tr.ResultOID || source.Tree != tr.ResultTree || source.Ref != task.Worktree.Ref || source.GitCommonDir != task.GitCommonDir {
		conflict("source_binding_mismatch")
	}
	if source.GitCommonDir != repo.GitCommonDir || parentObs.GitCommonDir != repo.GitCommonDir || parent.GitCommonDir != repo.GitCommonDir || source.ObjectFormat != repository.ObjectFormat || parentObs.ObjectFormat != repository.ObjectFormat || source.RefFormat != repository.RefFormat || parentObs.RefFormat != repository.RefFormat {
		conflict("repository_binding_mismatch")
	}
	if !integrationInventoryContains(repository.Inventory, source) {
		conflict("source_inventory_mismatch")
	}
	if !integrationInventoryContains(repository.Inventory, parentObs) {
		conflict("parent_inventory_mismatch")
	}
	if !source.Clean || len(source.InProgress) > 0 {
		block("source_not_clean")
	}
	if !parentObs.Clean || len(parentObs.InProgress) > 0 {
		block("parent_not_clean")
	}
	if parentObs.OID == tr.ResultOID && parentObs.Tree == tr.ResultTree && readiness == "ready" {
		readiness = "already_integrated"
	} else if parentObs.OID != input.ExpectedParentOID || parentObs.Tree != task.Worktree.ParentTree {
		conflict("parent_moved")
	}
	if !ancestor {
		conflict("result_not_descendant")
	}
	if repository.ObjectFormat != "sha1" && repository.ObjectFormat != "sha256" {
		block("unsupported_object_format")
	}
	if repository.RefFormat != "files" {
		block("unsupported_ref_format")
	}
	if !objectIDsMatchFormat(repository.ObjectFormat, input.ExpectedParentOID, input.ExpectedResultOID, task.Worktree.ParentTree, tr.ResultTree, source.OID, source.Tree, parentObs.OID, parentObs.Tree) {
		conflict("object_id_length_mismatch")
	}
	if repository.Shallow || repository.PartialClone || repository.SparseCheckout {
		block("unsupported_repository_state")
	}
	if input.RetryAfterResultID != nil {
		prior := findIntegrationResult(registry, *input.RetryAfterResultID)
		if prior == nil || prior.Outcome != "no_effect" || prior.AuthorityID == "" {
			conflict("invalid_retry_predecessor")
		} else {
			priorAuthority := findAuthority(registry, prior.AuthorityID)
			if priorAuthority == nil || priorAuthority.TaskID != input.TaskID || priorAuthority.TaskResultID != input.TaskResultID || priorAuthority.HumanQARecordID != input.HumanQARecordID || priorAuthority.Plan.Epic.ParentRef != parent.Worktree.Ref {
				conflict("retry_binding_mismatch")
			}
			for _, candidate := range registry.IntegrationAuthorities {
				if candidate.RetryAfterResultID != nil && *candidate.RetryAfterResultID == *input.RetryAfterResultID {
					conflict("retry_predecessor_already_used")
				}
			}
		}
	}
	for _, candidate := range registry.IntegrationAuthorities {
		if candidate.Plan.Epic.ParentRef == parent.Worktree.Ref && findIntegrationResultForAuthority(registry, candidate.ID) == nil {
			conflict("parent_integration_already_active")
		}
	}
	for _, prior := range registry.IntegrationResults {
		priorAuthority := findAuthority(registry, prior.AuthorityID)
		if priorAuthority == nil || priorAuthority.TaskID != input.TaskID || priorAuthority.TaskResultID != input.TaskResultID || priorAuthority.HumanQARecordID != input.HumanQARecordID || priorAuthority.Plan.Epic.ParentRef != parent.Worktree.Ref {
			continue
		}
		if prior.Outcome == "partial" || prior.Outcome == "unknown" {
			block("prior_integration_requires_reconciliation")
		}
		if prior.Outcome == "no_effect" && input.RetryAfterResultID == nil {
			block("prior_no_effect_requires_retry")
		}
	}
	reasons = sortedReasons(reasons)
	markerBytes, err := d.Files.ReadFile(filepath.Join(root, MarkerDirectory, MarkerFile))
	if err != nil {
		return WorkspaceTaskIntegrationPlan{}, "", ctx, workError(ErrorWorkIO, "read workspace marker", err)
	}
	storeDigest := registry.RawSHA256
	if storeDigest == "" {
		encoded, e := encodeWorkItemRegistry(registry)
		if e != nil {
			return WorkspaceTaskIntegrationPlan{}, "", ctx, e
		}
		storeDigest = digestTaskBytes(encoded)
	}
	argv := []string{"git", "-c", "core.hooksPath=" + os.DevNull, "-c", "merge.autoStash=false", "-c", "gc.auto=0", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "submodule.recurse=false", "-C", parent.Worktree.Locator, "merge", "--ff-only", "--no-stat", "--no-autostash", tr.ResultOID}
	plan := WorkspaceTaskIntegrationPlan{Kind: "WorkspaceTaskIntegrationPlan@1", SchemaVersion: 1, Format: "json", FormatVersion: 1, Canonicalization: "RFC8785", Workspace: IntegrationPlanWorkspace{Root: root, MarkerSHA256: digestTaskBytes(markerBytes)}, Project: IntegrationPlanProject{ProjectID: task.ProjectID}, Repository: IntegrationPlanRepository{RepoID: task.RepoID, RegisteredLocator: repo.Locator, GitCommonDir: repo.GitCommonDir, GitVersion: repository.GitVersion, ObjectFormat: repository.ObjectFormat, RefFormat: repository.RefFormat, Shallow: repository.Shallow, PartialClone: repository.PartialClone, SparseCheckout: repository.SparseCheckout}, Epic: IntegrationPlanEpic{EpicID: epic.ID, ParentWorktreeID: parent.Worktree.ID, ParentLocator: parent.Worktree.Locator, ParentRef: parent.Worktree.Ref, ExpectedParentOID: task.Worktree.ParentOID, ExpectedParentTree: task.Worktree.ParentTree}, Task: IntegrationPlanTask{TaskID: task.ID, TaskWorktreeID: task.Worktree.ID, SourceLocator: task.Worktree.Locator, SourceRef: task.Worktree.Ref, ResultOID: tr.ResultOID, ResultTree: tr.ResultTree}, TaskResult: IntegrationPlanTaskResult{ID: tr.ID, DraftSHA256: tr.DraftSHA256, TechnicalGate: tr.TechnicalGate}, HumanQA: IntegrationPlanHumanQA{ID: qa.ID, DraftSHA256: qa.DraftSHA256, Outcome: qa.Outcome}, RetryAfterResultID: input.RetryAfterResultID, WorkItemsSHA256: storeDigest, StoreTransition: tr.StoreTransition, ObservedSource: planObservation(source), ObservedParent: planObservation(parentObs), ObservedInventory: repository.Inventory, ObservedReflog: reflog, TechnicalGateReady: technical, HumanQAReady: human, AncestryReady: ancestor, Readiness: readiness, Reasons: reasons, Effect: IntegrationPlanEffect{Kind: "local_ff_only", ParentRef: parent.Worktree.Ref, ExpectedParentOID: task.Worktree.ParentOID, ResultOID: tr.ResultOID, MaxOccurrences: 1, Argv: argv}}
	if taskRequiresSpec(registry, task.ID) {
		plan.Kind = "WorkspaceTaskIntegrationPlan@2"
		plan.SchemaVersion = 2
		plan.TaskSpecGuard = taskResultSpecRelevance(d, root, projects, registry, *task, tr.ID)
		ctx.TaskSpecRelevance = plan.TaskSpecGuard
		if plan.TaskSpecGuard.Relevance != "current" {
			plan.Readiness = "blocked"
			plan.Reasons = sortedReasons(append(plan.Reasons, plan.TaskSpecGuard.Reasons...))
		}
	}
	if input.DeliveryAuthorization != nil {
		plan.Kind, plan.SchemaVersion = "WorkspaceTaskIntegrationPlan@3", 3
		copy := *input.DeliveryAuthorization
		copy.AllowedEffects = append([]string{}, copy.AllowedEffects...)
		plan.DeliveryAuthorization = &copy
	}
	digest := integrationPlanDigest(plan)
	return plan, digest, ctx, nil
}

func taskResultEvidenceRequest(result TaskResultRecord) TaskHandoffEvidenceRequest {
	return TaskHandoffEvidenceRequest{
		ActivityID: result.ActivityID, RunID: result.RunID, HandoffID: result.HandoffID, HandoffLocator: result.HandoffLocator, HandoffSHA256: result.HandoffSHA256,
		StartReceiptID: result.StartReceiptID, StartReceiptLocator: result.StartReceiptLocator, StartReceiptSHA256: result.StartReceiptSHA256,
		TerminalResultID: result.TerminalResultID, TerminalResultLocator: result.TerminalResultLocator, TerminalResultSHA256: result.TerminalResultSHA256,
		InspectionSHA256: result.InspectionSHA256,
	}
}

func revalidateIntegrationEvidence(d Dependencies, result TaskResultRecord, qa TaskHumanQARecord) error {
	evidence, err := d.HandoffEvidence.ReadTaskEvidence(taskResultEvidenceRequest(result))
	if err != nil {
		return workError(ErrorTaskIntegrationBlocked, "immutable handoff evidence can no longer be validated", err)
	}
	if err = ValidateTaskResultEvidence(result, evidence); err != nil {
		return err
	}
	if err = rehashQAEvidence(d.Files, qa.Evidence); err != nil {
		return workError(ErrorTaskIntegrationBlocked, "human QA evidence can no longer be validated", err)
	}
	return nil
}

// ValidateTaskResultEvidence compares a preserved result to revalidated native
// evidence using the original technical gate, including controlled evidence debt.
// It does not read files or establish runtime permission or human acceptance.
func ValidateTaskResultEvidence(result TaskResultRecord, evidence TaskHandoffEvidence) error {
	required := make([]string, len(result.VerifierResults))
	for i := range result.VerifierResults {
		required[i] = result.VerifierResults[i].VerifierID
	}
	draft := taskResultDraft{Handoff: taskResultEvidenceRequest(result), TechnicalGate: result.TechnicalGate, RequiredVerifierIDs: required, AcceptedDebt: result.AcceptedDebt, EvidenceArtifacts: result.Artifacts}
	if err := validateTaskEvidence(draft, evidence); err != nil || result.ReportedOutcome != evidence.ReportedOutcome || !reflect.DeepEqual(result.VerifierResults, evidence.VerifierResults) || !reflect.DeepEqual(result.Review, evidence.Review) {
		return workError(ErrorTaskIntegrationBlocked, "persisted Task result differs from revalidated handoff evidence", err)
	}
	return nil
}

func continueIntegrationAuthority(d Dependencies, session WorkItemStoreSession, registry *WorkItemRegistry, root string, ctx integrationContext, authority IntegrationAuthority, intent IntegrationIntent, output *WorkspaceTaskIntegrationReadback, recovered bool) error {
	target := QueueTarget{ProjectID: ctx.Task.ProjectID, RepoID: ctx.Task.RepoID, EpicID: ctx.Epic.ID, GitCommonDir: ctx.Task.GitCommonDir, ParentWorktreeID: ctx.Parent.Worktree.ID, ParentLocator: ctx.Parent.Worktree.Locator, ParentRef: ctx.Parent.Worktree.Ref}
	if e := validateIntegrationDelivery(d, root, *registry, ctx.Result, authority.HumanQARecordID, authority.Plan.DeliveryAuthorization, target); e != nil {
		return e
	}
	current, driftReason := observeAuthorityPrecondition(d, root, ctx, authority)
	if driftReason != "" {
		resultID, err := d.TaskLifecycleIDs.NewIntegrationResultID()
		if err != nil {
			return err
		}
		f := false
		now := d.WorkClock.Now().UTC().Format(time.RFC3339Nano)
		obs := integrationObservationBestEffort(now, current.ParentObservation, current.Repository.Inventory)
		res := IntegrationResult{ID: resultID, AuthorityID: authority.ID, IntentID: intent.ID, AttemptID: nil, RetryAfterResultID: authority.RetryAfterResultID, Outcome: "conflict", GitChanged: &f, RecordedAtUTC: now, BeforeObservation: obs, AfterObservation: obs, Command: emptyCommandEvidence(), Reflog: current.Reflog, Recovered: recovered, RecoveryStatus: "conflict", NextAction: IntegrationNextAction{Kind: "refresh_or_reverify_parent", Reason: driftReason, Argv: []string{}}}
		registry.IntegrationResults = append(registry.IntegrationResults, res)
		sortWorkRegistry(registry)
		published, err := publishWorkItemRegistryRecover(session, *registry)
		if err != nil {
			return err
		}
		*registry = published
		*output = newIntegrationReadback(root, *registry, current, authority.Plan, authority.PlanSHA256, &authority, &intent, nil, &res)
		return nil
	}
	ctx = current
	if authority.Plan.Readiness == "already_integrated" {
		resultID, err := d.TaskLifecycleIDs.NewIntegrationResultID()
		if err != nil {
			return err
		}
		f := false
		now := d.WorkClock.Now().UTC().Format(time.RFC3339Nano)
		obs := integrationObservation(now, ctx.ParentObservation, ctx.Repository.Inventory)
		res := IntegrationResult{ID: resultID, AuthorityID: authority.ID, IntentID: intent.ID, AttemptID: nil, RetryAfterResultID: authority.RetryAfterResultID, Outcome: "already_integrated", GitChanged: &f, RecordedAtUTC: now, BeforeObservation: obs, AfterObservation: obs, Command: emptyCommandEvidence(), Reflog: ctx.Reflog, Recovered: recovered, RecoveryStatus: "complete", NextAction: IntegrationNextAction{Kind: "none", Reason: "The Task is already integrated.", Argv: []string{}}}
		registry.IntegrationResults = append(registry.IntegrationResults, res)
		sortWorkRegistry(registry)
		published, err := publishWorkItemRegistryRecover(session, *registry)
		if err != nil {
			return err
		}
		*registry = published
		*output = newIntegrationReadback(root, *registry, ctx, authority.Plan, authority.PlanSHA256, &authority, &intent, nil, &res)
		return nil
	}
	attemptID, err := d.TaskLifecycleIDs.NewIntegrationAttemptID()
	if err != nil {
		return err
	}
	now := d.WorkClock.Now().UTC().Format(time.RFC3339Nano)
	attempt := IntegrationAttempt{ID: attemptID, AuthorityID: authority.ID, IntentID: intent.ID, Ordinal: 1, State: "prepared", ResultID: nil, PreparedAtUTC: now, CommandIdentity: "git-merge-ff-only", ReflogAction: "ply-workspace-task-integrate:" + string(attemptID), PreObservation: integrationObservation(now, ctx.ParentObservation, ctx.Repository.Inventory)}
	registry.IntegrationAttempts = append(registry.IntegrationAttempts, attempt)
	sortWorkRegistry(registry)
	published, err := publishWorkItemRegistryRecover(session, *registry)
	if err != nil {
		return err
	}
	*registry = published
	command := d.IntegrationGit.MergeFastForward(IntegrationMergeInput{Repository: ctx.Repo, ParentWorktree: ctx.Parent.Worktree.Locator, ResultOID: ctx.Result.ResultOID, AttemptID: attemptID})
	return finishIntegrationAttempt(d, session, registry, root, ctx, authority, intent, attempt, command, false, output)
}

func reconcileIntegrationAttempt(d Dependencies, session WorkItemStoreSession, registry *WorkItemRegistry, root string, ctx integrationContext, authority IntegrationAuthority, intent IntegrationIntent, attempt IntegrationAttempt, output *WorkspaceTaskIntegrationReadback) error {
	return finishIntegrationAttempt(d, session, registry, root, ctx, authority, intent, attempt, GitCommandOutcome{}, true, output)
}
func finishIntegrationAttempt(d Dependencies, session WorkItemStoreSession, registry *WorkItemRegistry, root string, ctx integrationContext, authority IntegrationAuthority, intent IntegrationIntent, attempt IntegrationAttempt, command GitCommandOutcome, recovered bool, output *WorkspaceTaskIntegrationReadback) error {
	now := d.WorkClock.Now().UTC().Format(time.RFC3339Nano)
	resultID, err := d.TaskLifecycleIDs.NewIntegrationResultID()
	if err != nil {
		return err
	}
	source, sourceErr := d.IntegrationGit.ObserveIntegrationWorktree(ctx.Task.Worktree.Locator, ctx.Task.Worktree.Ref)
	parent, parentErr := d.IntegrationGit.ObserveIntegrationWorktree(ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref)
	repository, repositoryErr := d.IntegrationGit.ObserveIntegrationRepository(ctx.Repo)
	reflog, reflogErr := d.IntegrationGit.ObserveParentReflog(ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref, 2)
	if reflog == nil {
		reflog = []IntegrationReflogEntry{}
	}
	outcome := "unknown"
	var changed *bool
	recovery := "unknown"
	next := IntegrationNextAction{Kind: "read_only_recovery_control", Reason: "State could not be classified safely.", Argv: []string{"ply", "workspace", "task", "show", string(ctx.Task.ID), "--format", "json"}}
	if sourceErr == nil && parentErr == nil && repositoryErr == nil && reflogErr == nil {
		if exactPostEffect(ctx, source, parent, repository, reflog, attempt.ID) {
			v := true
			changed = &v
			outcome = "exact_effect"
			recovery = "complete"
			next = IntegrationNextAction{Kind: "none", Reason: "The Task is integrated.", Argv: []string{}}
		} else if sameIntegrationWorktree(parent, observationWorktree(attempt.PreObservation)) && sameIntegrationWorktree(source, ctx.Source) && sameRepositoryObservation(repository, ctx.Repository) && !reflogMentionsAttempt(reflog, attempt.ID) {
			v := false
			changed = &v
			outcome = "no_effect"
			recovery = "safe-no-effect"
			next = IntegrationNextAction{Kind: "retry_after_no_effect", Reason: "Run a new check that names the no-effect result.", Argv: []string{"ply", "workspace", "task", "integrate", string(ctx.Task.ID), "--result", string(ctx.Result.ID), "--qa", string(ctx.QA.ID), "--expected-result-oid", ctx.Result.ResultOID, "--expected-parent-oid", ctx.Task.Worktree.ParentOID, "--retry-after", string(resultID), "--check"}}
		} else {
			outcome = "partial"
			recovery = "reconciliation-required"
		}
	}
	ctx.Source = source
	ctx.ParentObservation = parent
	ctx.Repository = repository
	ctx.Reflog = reflog
	ctx.SourceObserved = sourceErr == nil
	ctx.ParentObserved = parentErr == nil
	ctx.RepositoryObserved = repositoryErr == nil
	ctx.ReflogObserved = reflogErr == nil
	attemptID := attempt.ID
	afterInventory := []IntegrationInventoryEntry{}
	if repositoryErr == nil {
		afterInventory = repository.Inventory
	}
	res := IntegrationResult{ID: resultID, AuthorityID: authority.ID, IntentID: intent.ID, AttemptID: &attemptID, RetryAfterResultID: authority.RetryAfterResultID, Outcome: outcome, GitChanged: changed, RecordedAtUTC: now, BeforeObservation: attempt.PreObservation, AfterObservation: integrationObservationBestEffort(now, parent, afterInventory), Command: commandEvidence(command, recovered), Reflog: reflog, Recovered: recovered, RecoveryStatus: recovery, NextAction: next}
	for i := range registry.IntegrationAttempts {
		if registry.IntegrationAttempts[i].ID == attempt.ID {
			registry.IntegrationAttempts[i].State = "result_recorded"
			registry.IntegrationAttempts[i].ResultID = &resultID
			attempt = registry.IntegrationAttempts[i]
		}
	}
	registry.IntegrationResults = append(registry.IntegrationResults, res)
	sortWorkRegistry(registry)
	published, err := publishWorkItemRegistryRecover(session, *registry)
	if err != nil {
		// Git has already been attempted. Content corruption can prevent the
		// result write, but cannot erase the independently observed Git effect.
		observed, readErr := session.Snapshot()
		var observedAttempt *IntegrationAttempt
		if readErr == nil {
			observedAttempt = findAttemptForAuthority(observed, authority.ID)
		}
		*output = newIntegrationReadback(root, observed, ctx, authority.Plan, authority.PlanSHA256, &authority, &intent, observedAttempt, nil)
		if readErr != nil {
			replaceReadbackMember(&output.Value, "persisted_after", canonicaljson.Object{{Name: "format_version", Value: nil}, {Name: "sha256", Value: nil}, {Name: "authority_id", Value: nil}, {Name: "intent_id", Value: nil}, {Name: "attempt_id", Value: nil}, {Name: "integration_result_id", Value: nil}, {Name: "freshness", Value: "unknown"}})
		}
		output.Classification, output.GitChanged, output.RecoveryStatus = outcome, changed, "unknown"
		output.NextAction = IntegrationNextAction{Kind: "read_only_recovery_control", Reason: "The Git effect was observed, but its result could not be persisted. Inspect the Task before recovery.", Argv: []string{"ply", "workspace", "task", "show", string(ctx.Task.ID), "--format", "json"}}
		replaceReadbackMember(&output.Value, "classification", outcome)
		replaceReadbackMember(&output.Value, "git_changed", pointerValue(changed))
		replaceReadbackMember(&output.Value, "recovery_status", output.RecoveryStatus)
		replaceReadbackMember(&output.Value, "next_action", structCanonical(output.NextAction))
		return workError(ErrorWorkIO, fmt.Sprintf("observed integration outcome %s; result persistence failed; inspect with ply workspace task show %s --format json", outcome, ctx.Task.ID), err)
	}
	*registry = published
	*output = newIntegrationReadback(root, *registry, ctx, authority.Plan, authority.PlanSHA256, &authority, &intent, &attempt, &res)
	return nil
}

func observeAuthorityPrecondition(d Dependencies, root string, ctx integrationContext, authority IntegrationAuthority) (integrationContext, string) {
	ctx.SourceObserved = false
	ctx.ParentObserved = false
	ctx.RepositoryObserved = false
	ctx.ReflogObserved = false
	if ctx.PreconditionReason != "" {
		ctx.ParentObservation = IntegrationWorktreeObservation{StatusEntries: []IntegrationStatusEntry{}, InProgress: []string{}}
		ctx.Repository.Inventory = []IntegrationInventoryEntry{}
		ctx.Reflog = []IntegrationReflogEntry{}
		return ctx, ctx.PreconditionReason
	}
	marker, err := d.Files.ReadFile(filepath.Join(root, MarkerDirectory, MarkerFile))
	if err != nil || digestTaskBytes(marker) != authority.Plan.Workspace.MarkerSHA256 {
		return ctx, "The workspace marker changed after confirmation."
	}
	repository, err := d.IntegrationGit.ObserveIntegrationRepository(ctx.Repo)
	if err != nil {
		return ctx, "The repository can no longer be observed safely."
	}
	source, err := d.IntegrationGit.ObserveIntegrationWorktree(ctx.Task.Worktree.Locator, ctx.Task.Worktree.Ref)
	if err != nil {
		return ctx, "The Task source changed after confirmation."
	}
	parent, err := d.IntegrationGit.ObserveIntegrationWorktree(ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref)
	if err != nil {
		return ctx, "The Epic parent changed after confirmation."
	}
	reflog, err := d.IntegrationGit.ObserveParentReflog(ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref, 2)
	if err != nil {
		return ctx, "The Epic parent reflog changed after confirmation."
	}
	ancestor, err := d.IntegrationGit.CheckAncestor(ctx.Repo, authority.Plan.Epic.ExpectedParentOID, authority.Plan.Task.ResultOID)
	if err != nil || !ancestor {
		return ctx, "The confirmed ancestry is no longer valid."
	}
	current := ctx
	current.Repository, current.Source, current.ParentObservation, current.Reflog = repository, source, parent, reflog
	current.SourceObserved = true
	current.ParentObserved = true
	current.RepositoryObserved = true
	current.ReflogObserved = true
	if !sameRepositoryObservation(repository, ctx.Repository) || !sameIntegrationWorktree(source, worktreeFromPlan(authority.Plan.ObservedSource)) || !sameIntegrationWorktree(parent, worktreeFromPlan(authority.Plan.ObservedParent)) || !reflect.DeepEqual(reflog, authority.Plan.ObservedReflog) {
		return current, "Repository, source, parent, inventory, or reflog state changed after confirmation."
	}
	return current, ""
}

func authorityProjectBindingReason(projects ProjectSnapshot, plan WorkspaceTaskIntegrationPlan) string {
	_, repo, err := projectAndRepo(projects, plan.Project.ProjectID, plan.Repository.RepoID)
	if err != nil || repo.Locator != plan.Repository.RegisteredLocator || repo.GitCommonDir != plan.Repository.GitCommonDir {
		return "The registered Project repository changed after confirmation."
	}
	return ""
}

func sameRepositoryObservation(a, b IntegrationRepositoryObservation) bool {
	return a.GitVersion == b.GitVersion && a.ObjectFormat == b.ObjectFormat && a.RefFormat == b.RefFormat && a.Shallow == b.Shallow && a.PartialClone == b.PartialClone && a.SparseCheckout == b.SparseCheckout && reflect.DeepEqual(a.Inventory, b.Inventory)
}

func integrationInventoryContains(inventory []IntegrationInventoryEntry, worktree IntegrationWorktreeObservation) bool {
	matches := 0
	for _, entry := range inventory {
		if entry.Locator != worktree.Locator {
			continue
		}
		matches++
		if entry.Ref != worktree.Ref || entry.OID != worktree.OID || entry.Locked || entry.Prunable || entry.Bare || entry.Detached {
			return false
		}
	}
	return matches == 1
}

func objectIDsMatchFormat(format string, values ...string) bool {
	want := 0
	switch format {
	case "sha1":
		want = 40
	case "sha256":
		want = 64
	default:
		return false
	}
	for _, value := range values {
		if len(value) != want || !validOIDText(value) {
			return false
		}
	}
	return true
}

func exactPostEffect(ctx integrationContext, source, parent IntegrationWorktreeObservation, repository IntegrationRepositoryObservation, reflog []IntegrationReflogEntry, attemptID IntegrationAttemptID) bool {
	expectedRepository := ctx.Repository
	expectedRepository.Inventory = inventoryAfterFastForward(ctx.Repository.Inventory, ctx.Parent.Worktree.Locator, ctx.Parent.Worktree.Ref, ctx.Result.ResultOID)
	return source.Locator == ctx.Task.Worktree.Locator && source.Ref == ctx.Task.Worktree.Ref && source.OID == ctx.Result.ResultOID && source.Tree == ctx.Result.ResultTree && source.GitCommonDir == ctx.Task.GitCommonDir && source.Clean && len(source.StatusEntries) == 0 && len(source.InProgress) == 0 &&
		parent.Locator == ctx.Parent.Worktree.Locator && parent.Ref == ctx.Parent.Worktree.Ref && parent.OID == ctx.Result.ResultOID && parent.Tree == ctx.Result.ResultTree && parent.GitCommonDir == ctx.Parent.GitCommonDir && parent.Clean && len(parent.StatusEntries) == 0 && len(parent.InProgress) == 0 &&
		sameRepositoryObservation(repository, expectedRepository) && reflogProves(reflog, attemptID, ctx.Task.Worktree.ParentOID, ctx.Result.ResultOID)
}

func inventoryAfterFastForward(before []IntegrationInventoryEntry, parentLocator, parentRef, resultOID string) []IntegrationInventoryEntry {
	after := append([]IntegrationInventoryEntry(nil), before...)
	for i := range after {
		if after[i].Locator == parentLocator && after[i].Ref == parentRef {
			after[i].OID = resultOID
		}
	}
	return after
}

func intentFromPlan(id IntegrationIntentID, authority IntegrationAuthorityID, p WorkspaceTaskIntegrationPlan, digest string) IntegrationIntent {
	return IntegrationIntent{ID: id, AuthorityID: authority, State: "prepared", WorkspaceRoot: p.Workspace.Root, WorkspaceMarkerSHA256: p.Workspace.MarkerSHA256, ProjectID: p.Project.ProjectID, RepoID: p.Repository.RepoID, RegisteredRepoLocator: p.Repository.RegisteredLocator, GitCommonDir: p.Repository.GitCommonDir, ObjectFormat: p.Repository.ObjectFormat, RefFormat: p.Repository.RefFormat, EpicID: p.Epic.EpicID, ParentWorktreeID: p.Epic.ParentWorktreeID, ParentLocator: p.Epic.ParentLocator, ParentRef: p.Epic.ParentRef, ExpectedParentOID: p.Epic.ExpectedParentOID, ExpectedParentTree: p.Epic.ExpectedParentTree, TaskID: p.Task.TaskID, TaskWorktreeID: p.Task.TaskWorktreeID, SourceLocator: p.Task.SourceLocator, SourceRef: p.Task.SourceRef, ResultOID: p.Task.ResultOID, ResultTree: p.Task.ResultTree, TaskResultID: p.TaskResult.ID, HumanQARecordID: p.HumanQA.ID, PlanSHA256: digest, Effect: IntegrationAllowedEffect{Kind: "local_ff_only", ParentRef: p.Epic.ParentRef, ExpectedParentOID: p.Epic.ExpectedParentOID, ResultOID: p.Task.ResultOID, MaxOccurrences: 1}, MaxAttempts: 1}
}
func integrationPlanDigest(p WorkspaceTaskIntegrationPlan) string {
	b, _ := canonicaljson.Marshal(integrationPlanCanonical(p))
	return digestTaskBytes(b)
}
func integrationIntentDigest(i IntegrationIntent) string {
	full := integrationIntentCanonical(i).(canonicaljson.Object)
	content := make(canonicaljson.Object, 0, len(full)-2)
	for _, member := range full {
		if member.Name != "id" && member.Name != "intent_sha256" {
			content = append(content, member)
		}
	}
	b, _ := canonicaljson.Marshal(content)
	return digestTaskBytes(b)
}
func planObservation(o IntegrationWorktreeObservation) PlanWorktreeObservation {
	return PlanWorktreeObservation{WorktreeLocator: o.Locator, Ref: o.Ref, OID: o.OID, Tree: o.Tree, GitCommonDir: o.GitCommonDir, ObjectFormat: o.ObjectFormat, RefFormat: o.RefFormat, Symbolic: o.Symbolic, Clean: o.Clean, StatusEntries: o.StatusEntries, InProgress: o.InProgress}
}
func integrationObservation(at string, o IntegrationWorktreeObservation, inventory []IntegrationInventoryEntry) IntegrationObservation {
	return IntegrationObservation{ObservedAtUTC: at, WorktreeLocator: o.Locator, Ref: o.Ref, OID: o.OID, Tree: o.Tree, GitCommonDir: o.GitCommonDir, ObjectFormat: o.ObjectFormat, RefFormat: o.RefFormat, Symbolic: o.Symbolic, Clean: o.Clean, StatusEntries: o.StatusEntries, InProgress: o.InProgress, InventoryEntries: inventory}
}
func integrationObservationBestEffort(at string, o IntegrationWorktreeObservation, inventory []IntegrationInventoryEntry) IntegrationObservation {
	if o.StatusEntries == nil {
		o.StatusEntries = []IntegrationStatusEntry{}
	}
	if o.InProgress == nil {
		o.InProgress = []string{}
	}
	if inventory == nil {
		inventory = []IntegrationInventoryEntry{}
	}
	return integrationObservation(at, o, inventory)
}
func observationWorktree(o IntegrationObservation) IntegrationWorktreeObservation {
	return IntegrationWorktreeObservation{Locator: o.WorktreeLocator, Ref: o.Ref, OID: o.OID, Tree: o.Tree, GitCommonDir: o.GitCommonDir, ObjectFormat: o.ObjectFormat, RefFormat: o.RefFormat, Symbolic: o.Symbolic, Clean: o.Clean, StatusEntries: o.StatusEntries, InProgress: o.InProgress}
}
func sameIntegrationWorktree(a, b IntegrationWorktreeObservation) bool {
	return a.Locator == b.Locator && a.Ref == b.Ref && a.OID == b.OID && a.Tree == b.Tree && a.GitCommonDir == b.GitCommonDir && a.ObjectFormat == b.ObjectFormat && a.RefFormat == b.RefFormat && a.Symbolic == b.Symbolic && a.Clean == b.Clean && reflect.DeepEqual(a.StatusEntries, b.StatusEntries) && reflect.DeepEqual(a.InProgress, b.InProgress)
}
func reflogProves(entries []IntegrationReflogEntry, id IntegrationAttemptID, before, after string) bool {
	return len(entries) >= 2 && entries[0].OID == after && entries[1].OID == before && strings.HasPrefix(entries[0].Action, "ply-workspace-task-integrate:"+string(id)+":")
}
func reflogMentionsAttempt(entries []IntegrationReflogEntry, id IntegrationAttemptID) bool {
	prefix := "ply-workspace-task-integrate:" + string(id) + ":"
	for _, entry := range entries {
		if strings.HasPrefix(entry.Action, prefix) {
			return true
		}
	}
	return false
}
func emptyCommandEvidence() IntegrationCommandEvidence {
	return IntegrationCommandEvidence{Launched: false, ExitCode: nil, LaunchError: nil, Stdout: commandStream(nil), Stderr: commandStream(nil)}
}
func commandEvidence(o GitCommandOutcome, recovered bool) IntegrationCommandEvidence {
	if recovered {
		return emptyCommandEvidence()
	}
	launched := o.Err == nil || o.Exit >= 0
	var exit *int
	if launched {
		x := o.Exit
		exit = &x
	}
	var launch *string
	if o.Err != nil && o.Exit < 0 {
		x := "launch_failed"
		launch = &x
	}
	return IntegrationCommandEvidence{Launched: launched, ExitCode: exit, LaunchError: launch, Stdout: commandStream(o.Stdout), Stderr: commandStream(o.Stderr)}
}
func commandStream(b []byte) IntegrationCommandStream {
	d := sha256.Sum256(b)
	captured := b
	truncated := false
	if len(captured) > 64<<10 {
		captured = captured[:64<<10]
		truncated = true
	}
	return IntegrationCommandStream{SizeBytes: int64(len(b)), CapturedBytes: int64(len(captured)), Truncated: truncated, SHA256: "sha256:" + hex.EncodeToString(d[:]), Base64: base64.StdEncoding.EncodeToString(captured)}
}
func findAuthorityByPlan(r WorkItemRegistry, digest string, input TaskIntegrationInput) *IntegrationAuthority {
	for i := range r.IntegrationAuthorities {
		a := &r.IntegrationAuthorities[i]
		retryMatches := (a.RetryAfterResultID == nil) == (input.RetryAfterResultID == nil)
		if retryMatches && a.RetryAfterResultID != nil {
			retryMatches = *a.RetryAfterResultID == *input.RetryAfterResultID
		}
		if a.PlanSHA256 == digest && a.TaskID == input.TaskID && a.TaskResultID == input.TaskResultID && a.HumanQARecordID == input.HumanQARecordID && retryMatches && contentTypedEqual(a.DeliveryOwner, input.DeliveryOwner) && contentTypedEqual(a.Plan.DeliveryAuthorization, input.DeliveryAuthorization) && a.Plan.Task.ResultOID == input.ExpectedResultOID && a.Plan.Epic.ExpectedParentOID == input.ExpectedParentOID {
			return a
		}
	}
	return nil
}

func validDeliveryIntegrationOwner(o *DeliveryIntegrationOwner) bool {
	return o != nil && validTaskText(o.RunID, 1, 256) && digestPattern.MatchString(o.RequestSHA256) && validTaskText(o.ActorClaim, 1, 256) && validTaskText(o.PreparationID, 1, 256)
}

// A delivery owner may resume after the merge, queue close or base update. Its
// original authority remains the sole effect ledger even when the current Epic
// basis has advanced to the delivered candidate. Never build a second plan for
// the same return, and never present stored before-state as a fresh observation.
func checkPersistedDeliveryIntegration(d Dependencies, root string, projects ProjectSnapshot, registry WorkItemRegistry, input TaskIntegrationInput) (TaskIntegrationResult, bool, error) {
	var authority *IntegrationAuthority
	for i := range registry.IntegrationAuthorities {
		a := &registry.IntegrationAuthorities[i]
		if a.TaskID != input.TaskID || a.TaskResultID != input.TaskResultID || a.HumanQARecordID != input.HumanQARecordID || a.Plan.Task.ResultOID != input.ExpectedResultOID || a.Plan.Epic.ExpectedParentOID != input.ExpectedParentOID || !contentTypedEqual(a.RetryAfterResultID, input.RetryAfterResultID) {
			continue
		}
		if !contentTypedEqual(a.DeliveryOwner, input.DeliveryOwner) || !contentTypedEqual(a.Plan.DeliveryAuthorization, input.DeliveryAuthorization) || authority != nil {
			return TaskIntegrationResult{}, false, workError(ErrorTaskIntegrationConflict, "delivery return authority is different or ambiguous", nil)
		}
		authority = a
	}
	if authority == nil {
		return TaskIntegrationResult{}, false, nil
	}
	if reason := authorityProjectBindingReason(projects, authority.Plan); reason != "" {
		return TaskIntegrationResult{}, true, workError(ErrorTaskIntegrationConflict, reason, nil)
	}
	ctx, err := contextFromPlan(registry, authority.Plan)
	if err != nil {
		return TaskIntegrationResult{}, true, err
	}
	if err = revalidateIntegrationEvidence(d, ctx.Result, ctx.QA); err != nil {
		return TaskIntegrationResult{}, true, err
	}
	intent := findIntentForAuthority(registry, authority.ID)
	if intent == nil {
		return TaskIntegrationResult{}, true, workError(ErrorTaskIntegrationConflict, "delivery return authority has no intent", nil)
	}
	attempt, result := findAttemptForAuthority(registry, authority.ID), findIntegrationResultForAuthority(registry, authority.ID)
	if attempt == nil && result == nil {
		// An unattempted intent is a proposed effect, so its original pass must
		// still be current. Completed or uncertain attempts remain historical
		// observations and are not relabeled as a fresh permission to execute.
		target := QueueTarget{ProjectID: ctx.Task.ProjectID, RepoID: ctx.Task.RepoID, EpicID: ctx.Epic.ID, GitCommonDir: ctx.Task.GitCommonDir, ParentWorktreeID: ctx.Parent.Worktree.ID, ParentLocator: ctx.Parent.Worktree.Locator, ParentRef: ctx.Parent.Worktree.Ref}
		if err = validateIntegrationDelivery(d, root, registry, ctx.Result, authority.HumanQARecordID, authority.Plan.DeliveryAuthorization, target); err != nil {
			return TaskIntegrationResult{}, true, err
		}
	}
	ctx.Source, err = d.IntegrationGit.ObserveIntegrationWorktree(authority.Plan.Task.SourceLocator, authority.Plan.Task.SourceRef)
	if err != nil {
		return TaskIntegrationResult{}, true, err
	}
	ctx.SourceObserved = true
	ctx.ParentObservation, err = d.IntegrationGit.ObserveIntegrationWorktree(authority.Plan.Epic.ParentLocator, authority.Plan.Epic.ParentRef)
	if err != nil {
		return TaskIntegrationResult{}, true, err
	}
	ctx.ParentObserved = true
	if result != nil && result.RecoveryStatus == "complete" {
		for _, observed := range []IntegrationWorktreeObservation{ctx.Source, ctx.ParentObservation} {
			if observed.OID != input.ExpectedResultOID || observed.Tree != ctx.Result.ResultTree || observed.GitCommonDir != authority.Plan.Repository.GitCommonDir || !observed.Clean || !observed.Symbolic || len(observed.InProgress) != 0 {
				return TaskIntegrationResult{}, true, workError(ErrorTaskIntegrationConflict, "completed delivery return has subsequent source or parent drift", nil)
			}
		}
	}
	ctx.TaskSpecRelevance = taskResultSpecRelevance(d, root, projects, registry, ctx.Task, ctx.Result.ID)
	ctx.PersistedBefore = integrationPersistedView(registry, authority, intent, attempt, result)
	return TaskIntegrationResult{Readback: newIntegrationReadback(root, registry, ctx, authority.Plan, authority.PlanSHA256, authority, intent, attempt, result)}, true, nil
}
func findIntentForAuthority(r WorkItemRegistry, id IntegrationAuthorityID) *IntegrationIntent {
	for i := range r.IntegrationIntents {
		if r.IntegrationIntents[i].AuthorityID == id {
			return &r.IntegrationIntents[i]
		}
	}
	return nil
}
func findAuthority(r WorkItemRegistry, id IntegrationAuthorityID) *IntegrationAuthority {
	for i := range r.IntegrationAuthorities {
		if r.IntegrationAuthorities[i].ID == id {
			return &r.IntegrationAuthorities[i]
		}
	}
	return nil
}
func findAttemptForAuthority(r WorkItemRegistry, id IntegrationAuthorityID) *IntegrationAttempt {
	for i := range r.IntegrationAttempts {
		if r.IntegrationAttempts[i].AuthorityID == id {
			return &r.IntegrationAttempts[i]
		}
	}
	return nil
}
func findIntegrationResultForAuthority(r WorkItemRegistry, id IntegrationAuthorityID) *IntegrationResult {
	for i := range r.IntegrationResults {
		if r.IntegrationResults[i].AuthorityID == id {
			return &r.IntegrationResults[i]
		}
	}
	return nil
}
func findIntegrationResult(r WorkItemRegistry, id IntegrationResultID) *IntegrationResult {
	for i := range r.IntegrationResults {
		if r.IntegrationResults[i].ID == id {
			return &r.IntegrationResults[i]
		}
	}
	return nil
}

func unresolvedIntegrationLeaf(r WorkItemRegistry, input TaskIntegrationInput) (*IntegrationAuthority, *IntegrationIntent, *IntegrationAttempt, *IntegrationResult) {
	usedAsPredecessor := map[IntegrationResultID]bool{}
	for _, authority := range r.IntegrationAuthorities {
		if authority.RetryAfterResultID != nil {
			usedAsPredecessor[*authority.RetryAfterResultID] = true
		}
	}
	var selectedAuthority *IntegrationAuthority
	var selectedResult *IntegrationResult
	for i := range r.IntegrationResults {
		result := &r.IntegrationResults[i]
		if usedAsPredecessor[result.ID] || (result.Outcome != "no_effect" && result.Outcome != "partial" && result.Outcome != "unknown") {
			continue
		}
		authority := findAuthority(r, result.AuthorityID)
		if authority == nil || authority.TaskID != input.TaskID || authority.TaskResultID != input.TaskResultID || authority.HumanQARecordID != input.HumanQARecordID || authority.Plan.Task.ResultOID != input.ExpectedResultOID || authority.Plan.Epic.ExpectedParentOID != input.ExpectedParentOID {
			continue
		}
		if selectedResult != nil {
			return nil, nil, nil, nil
		}
		selectedAuthority, selectedResult = authority, result
	}
	if selectedResult == nil {
		return nil, nil, nil, nil
	}
	return selectedAuthority, findIntentForAuthority(r, selectedAuthority.ID), findAttemptForAuthority(r, selectedAuthority.ID), selectedResult
}
func contextFromPlan(r WorkItemRegistry, p WorkspaceTaskIntegrationPlan) (integrationContext, error) {
	task, _ := findTask(r, p.Task.TaskID)
	epic, _ := findEpic(r, p.Epic.EpicID)
	tr := findTaskResult(r, p.TaskResult.ID)
	qa := findHumanQA(r, p.HumanQA.ID)
	if task == nil || epic == nil || tr == nil || qa == nil {
		return integrationContext{}, workError(ErrorTaskIntegrationConflict, "durable plan references missing state", nil)
	}
	parent, _ := findEpicRepo(*epic, task.RepoID)
	if parent == nil {
		return integrationContext{}, workError(ErrorTaskIntegrationConflict, "durable plan parent binding is missing", nil)
	}
	return integrationContext{Task: *task, Epic: *epic, Parent: *parent, Result: *tr, QA: *qa, Repo: RepoRecord{ID: p.Repository.RepoID, Locator: p.Repository.RegisteredLocator, GitCommonDir: p.Repository.GitCommonDir}, Repository: IntegrationRepositoryObservation{GitVersion: p.Repository.GitVersion, ObjectFormat: p.Repository.ObjectFormat, RefFormat: p.Repository.RefFormat, Shallow: p.Repository.Shallow, PartialClone: p.Repository.PartialClone, SparseCheckout: p.Repository.SparseCheckout, Inventory: p.ObservedInventory}, Source: worktreeFromPlan(p.ObservedSource), ParentObservation: worktreeFromPlan(p.ObservedParent), Reflog: p.ObservedReflog}, nil
}
func worktreeFromPlan(p PlanWorktreeObservation) IntegrationWorktreeObservation {
	return IntegrationWorktreeObservation{Locator: p.WorktreeLocator, Ref: p.Ref, OID: p.OID, Tree: p.Tree, GitCommonDir: p.GitCommonDir, ObjectFormat: p.ObjectFormat, RefFormat: p.RefFormat, Symbolic: p.Symbolic, Clean: p.Clean, StatusEntries: p.StatusEntries, InProgress: p.InProgress}
}

// Canonical plan conversion is explicit so the digest is independent of Go's JSON encoder.
func integrationPlanCanonical(p WorkspaceTaskIntegrationPlan) canonicaljson.Value {
	return structCanonical(p)
}
func integrationIntentCanonical(i IntegrationIntent) canonicaljson.Value { return structCanonical(i) }
func structCanonical(v any) canonicaljson.Value {
	// The plan/intent structs contain only deterministic scalar, slice, pointer and struct
	// fields. YAML round-tripping would lose JSON names; this reflective converter uses tags.
	return canonicalFromReflect(v)
}
