package workspace

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const (
	WorkItemsFile = "work-items.yaml"
	workItemsLock = "work-items.lock"
)

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

type WorktreeObservationRecord struct {
	Classification             string   `yaml:"classification" json:"classification"`
	ParentRefOID               *string  `yaml:"parent_ref_oid" json:"parent_ref_oid"`
	ParentRefTree              *string  `yaml:"parent_ref_tree" json:"parent_ref_tree"`
	ParentWorktreeLocator      *string  `yaml:"parent_worktree_locator" json:"parent_worktree_locator"`
	ParentWorktreeRef          *string  `yaml:"parent_worktree_ref" json:"parent_worktree_ref"`
	ParentWorktreeOID          *string  `yaml:"parent_worktree_oid" json:"parent_worktree_oid"`
	ParentWorktreeTree         *string  `yaml:"parent_worktree_tree" json:"parent_worktree_tree"`
	ParentWorktreeGitCommonDir *string  `yaml:"parent_worktree_git_common_dir" json:"parent_worktree_git_common_dir"`
	ParentWorktreeClean        *bool    `yaml:"parent_worktree_clean" json:"parent_worktree_clean"`
	SourceRefOID               *string  `yaml:"source_ref_oid" json:"source_ref_oid"`
	SourceRefTree              *string  `yaml:"source_ref_tree" json:"source_ref_tree"`
	SourceRefCheckedOutAt      *string  `yaml:"source_ref_checked_out_at" json:"source_ref_checked_out_at"`
	TargetKind                 string   `yaml:"target_kind" json:"target_kind"`
	TargetLocator              *string  `yaml:"target_locator" json:"target_locator"`
	TargetRef                  *string  `yaml:"target_ref" json:"target_ref"`
	TargetOID                  *string  `yaml:"target_oid" json:"target_oid"`
	TargetTree                 *string  `yaml:"target_tree" json:"target_tree"`
	TargetGitCommonDir         *string  `yaml:"target_git_common_dir" json:"target_git_common_dir"`
	TargetClean                *bool    `yaml:"target_clean" json:"target_clean"`
	InventoryMatch             *bool    `yaml:"inventory_match" json:"inventory_match"`
	Reasons                    []string `yaml:"reasons" json:"reasons"`
}

type WorktreeOperationRecord struct {
	ID               WorktreeOperationID       `yaml:"id" json:"id"`
	TaskID           TaskID                    `yaml:"task_id" json:"task_id"`
	WorktreeID       WorktreeID                `yaml:"worktree_id" json:"worktree_id"`
	IntentDigest     string                    `yaml:"intent_digest" json:"intent_digest"`
	State            string                    `yaml:"state" json:"state"`
	ProjectID        ProjectID                 `yaml:"project_id" json:"project_id"`
	RepoID           RepoID                    `yaml:"repo_id" json:"repo_id"`
	GitCommonDir     string                    `yaml:"git_common_dir" json:"git_common_dir"`
	ParentEpicID     EpicID                    `yaml:"parent_epic_id" json:"parent_epic_id"`
	ParentWorktreeID WorktreeID                `yaml:"parent_worktree_id" json:"parent_worktree_id"`
	ParentRef        string                    `yaml:"parent_ref" json:"parent_ref"`
	ParentOID        string                    `yaml:"parent_oid" json:"parent_oid"`
	ParentTree       string                    `yaml:"parent_tree" json:"parent_tree"`
	SourceRef        string                    `yaml:"source_ref" json:"source_ref"`
	TargetLocator    string                    `yaml:"target_locator" json:"target_locator"`
	LastObservation  WorktreeObservationRecord `yaml:"last_observation" json:"last_observation"`
}

type WorkItemRegistry struct {
	FormatVersion           int                       `yaml:"format_version"`
	Epics                   []EpicRecord              `yaml:"epics"`
	Tasks                   []TaskRecord              `yaml:"tasks"`
	WorktreeOperations      []WorktreeOperationRecord `yaml:"worktree_operations"`
	TaskResults             []TaskResultRecord        `yaml:"task_results,omitempty"`
	HumanQARecords          []TaskHumanQARecord       `yaml:"human_qa_records,omitempty"`
	IntegrationAuthorities  []IntegrationAuthority    `yaml:"integration_authorities,omitempty"`
	IntegrationIntents      []IntegrationIntent       `yaml:"integration_intents,omitempty"`
	IntegrationAttempts     []IntegrationAttempt      `yaml:"integration_attempts,omitempty"`
	IntegrationResults      []IntegrationResult       `yaml:"integration_results,omitempty"`
	TaskSpecPolicies        []TaskSpecPolicy          `yaml:"task_spec_policies"`
	TaskProblemRevisions    []TaskProblemReference    `yaml:"task_problem_revisions"`
	TaskSpecRevisions       []TaskSpecReference       `yaml:"task_spec_revisions"`
	TaskSpecAssessments     []TaskAssessmentReference `yaml:"task_spec_assessments"`
	TaskSolutionSelections  []TaskSelectionReference  `yaml:"task_solution_selections"`
	TaskResultSpecBindings  []TaskResultSpecBinding   `yaml:"task_result_spec_bindings"`
	TaskContentPublications []TaskContentPublication  `yaml:"task_content_publications"`
	EpicBaseVersions        []EpicBaseVersion         `yaml:"epic_base_versions"`
	EpicBaseUpdates         []EpicBaseUpdate          `yaml:"epic_base_updates"`
	TaskQueueEvents         []TaskQueueEvent          `yaml:"task_queue_events"`
	TaskPreparations        []TaskPreparation         `yaml:"task_preparations"`
	RawSHA256               string                    `yaml:"-"`
}

type EpicAdoptInput struct {
	EpicID      EpicID
	Title       string
	ProjectID   ProjectID
	RepoID      RepoID
	Worktree    string
	Ref         string
	ExpectedOID string
}

type TaskCreateInput struct {
	RegistryUpgrade *TaskRegistryUpgrade
	TaskID          TaskID
	Title           string
	Description     string
	ParentEpicID    EpicID
	ProjectID       ProjectID
	RepoID          RepoID
}

type TaskWorktreeCreateInput struct {
	TaskID            TaskID
	Branch            string
	Path              string
	ExpectedParentOID string
}

type EpicMutationResult struct {
	Workspace string
	Epic      EpicRecord
	Created   bool
}

type TaskMutationResult struct {
	Content   *TaskContentMutationResult
	Workspace string
	Task      TaskRecord
	Epic      EpicRecord
	Created   bool
}

type TaskWorktreeMutationResult struct {
	Workspace string
	Task      TaskRecord
	Epic      EpicRecord
	Operation WorktreeOperationRecord
	Outcome   string
}

type EpicListResult struct {
	Workspace string
	Epics     []EpicRecord
}
type TaskListResult struct {
	Lifecycles    WorkItemLifecycleSnapshot
	CurrentTitles map[TaskID]string
	Titles        map[TaskID]TaskListTitle
	Filters       TaskListFilters
	Workspace     string
	Tasks         []TaskRecord
	Epic          *EpicRecord
}

type WorkItemErrorClass string

const (
	ErrorWorkInvalidArguments                  WorkItemErrorClass = "workspace_work_invalid_arguments"
	ErrorWorkWorkspaceConflict                 WorkItemErrorClass = "workspace_work_workspace_conflict"
	ErrorWorkNotFound                          WorkItemErrorClass = "workspace_work_not_found"
	ErrorWorkProjectConflict                   WorkItemErrorClass = "workspace_work_project_conflict"
	ErrorWorkIdentityConflict                  WorkItemErrorClass = "workspace_work_identity_conflict"
	ErrorWorkGitObservation                    WorkItemErrorClass = "workspace_work_git_observation_error"
	ErrorWorkGitEffect                         WorkItemErrorClass = "workspace_work_git_effect_error"
	ErrorWorkParentDirty                       WorkItemErrorClass = "workspace_work_parent_dirty"
	ErrorWorkParentStale                       WorkItemErrorClass = "workspace_work_parent_stale"
	ErrorWorkPathConflict                      WorkItemErrorClass = "workspace_work_path_conflict"
	ErrorWorkRefConflict                       WorkItemErrorClass = "workspace_work_ref_conflict"
	ErrorWorkReconciliation                    WorkItemErrorClass = "workspace_work_reconciliation_required"
	ErrorWorkStoreConflict                     WorkItemErrorClass = "workspace_work_store_conflict"
	ErrorWorkIO                                WorkItemErrorClass = "workspace_work_io_error"
	ErrorTaskResultSchemaInvalid               WorkItemErrorClass = "workspace_task_result_schema_invalid"
	ErrorTaskResultEvidenceConflict            WorkItemErrorClass = "workspace_task_result_evidence_conflict"
	ErrorTaskResultConflict                    WorkItemErrorClass = "workspace_task_result_conflict"
	ErrorTaskQASchemaInvalid                   WorkItemErrorClass = "workspace_task_qa_schema_invalid"
	ErrorTaskQAEvidenceConflict                WorkItemErrorClass = "workspace_task_qa_evidence_conflict"
	ErrorTaskQAConflict                        WorkItemErrorClass = "workspace_task_qa_conflict"
	ErrorTaskIntegrationBlocked                WorkItemErrorClass = "workspace_task_integration_blocked"
	ErrorTaskIntegrationConflict               WorkItemErrorClass = "workspace_task_integration_conflict"
	ErrorTaskIntegrationConfirmationMismatch   WorkItemErrorClass = "workspace_task_integration_confirmation_mismatch"
	ErrorTaskIntegrationReconciliationRequired WorkItemErrorClass = "workspace_task_integration_reconciliation_required"
	ErrorTaskIntegrationUnsupportedRepository  WorkItemErrorClass = "workspace_task_integration_unsupported_repository"
)

type WorkItemError struct {
	Class  WorkItemErrorClass
	Detail string
	Err    error
}

func (err *WorkItemError) Error() string {
	detail := err.Detail
	if detail == "" && err.Err != nil {
		detail = err.Err.Error()
	}
	return fmt.Sprintf("%s: %s", err.Class, detail)
}
func (err *WorkItemError) Unwrap() error { return err.Err }

func workError(class WorkItemErrorClass, detail string, err error) error {
	return &WorkItemError{Class: class, Detail: detail, Err: err}
}

func WorkInvalidArguments(detail string) error {
	return workError(ErrorWorkInvalidArguments, detail, nil)
}

func ParseEpicID(value string) (EpicID, error) {
	if err := validateWorkIdentifier("Epic", value); err != nil {
		return "", err
	}
	return EpicID(value), nil
}
func ParseTaskID(value string) (TaskID, error) {
	if err := validateWorkIdentifier("Task", value); err != nil {
		return "", err
	}
	return TaskID(value), nil
}

func ParseEpicAdoptInput(epicID, title, projectID, repoID, worktree, ref, expectedOID string) (EpicAdoptInput, error) {
	id, err := ParseEpicID(epicID)
	if err != nil {
		return EpicAdoptInput{}, err
	}
	if err := validateWorkText("Epic title", title, 128); err != nil {
		return EpicAdoptInput{}, err
	}
	if err := validateWorkIdentifier("Project", projectID); err != nil {
		return EpicAdoptInput{}, err
	}
	if err := validateWorkIdentifier("repository", repoID); err != nil {
		return EpicAdoptInput{}, err
	}
	if worktree == "" || ref == "" || expectedOID == "" {
		return EpicAdoptInput{}, WorkInvalidArguments("--worktree, --ref, and --expected-oid are required")
	}
	if !validFullBranchRef(ref) {
		return EpicAdoptInput{}, WorkInvalidArguments("--ref must be an exact full local refs/heads/* ref")
	}
	if !validOIDText(expectedOID) {
		return EpicAdoptInput{}, WorkInvalidArguments("--expected-oid must be a lowercase full Git object ID")
	}
	return EpicAdoptInput{EpicID: id, Title: title, ProjectID: ProjectID(projectID), RepoID: RepoID(repoID), Worktree: worktree, Ref: ref, ExpectedOID: expectedOID}, nil
}

func ParseTaskCreateInput(taskID, title, description, epicID, projectID, repoID string) (TaskCreateInput, error) {
	id, err := ParseTaskID(taskID)
	if err != nil {
		return TaskCreateInput{}, err
	}
	parent, err := ParseEpicID(epicID)
	if err != nil {
		return TaskCreateInput{}, err
	}
	if err := validateWorkText("Task title", title, 128); err != nil {
		return TaskCreateInput{}, err
	}
	if err := validateWorkText("Task description", description, 2048); err != nil {
		return TaskCreateInput{}, err
	}
	if err := validateWorkIdentifier("Project", projectID); err != nil {
		return TaskCreateInput{}, err
	}
	if err := validateWorkIdentifier("repository", repoID); err != nil {
		return TaskCreateInput{}, err
	}
	return TaskCreateInput{TaskID: id, Title: title, Description: description, ParentEpicID: parent, ProjectID: ProjectID(projectID), RepoID: RepoID(repoID)}, nil
}

func ParseTaskWorktreeCreateInput(taskID, branch, path, expectedParentOID string) (TaskWorktreeCreateInput, error) {
	id, err := ParseTaskID(taskID)
	if err != nil {
		return TaskWorktreeCreateInput{}, err
	}
	if !validShortBranchText(branch) {
		return TaskWorktreeCreateInput{}, WorkInvalidArguments("--branch must be a valid short local branch name")
	}
	if path == "" {
		return TaskWorktreeCreateInput{}, WorkInvalidArguments("--path is required")
	}
	if !validOIDText(expectedParentOID) {
		return TaskWorktreeCreateInput{}, WorkInvalidArguments("--expected-parent-oid must be a lowercase full Git object ID")
	}
	return TaskWorktreeCreateInput{TaskID: id, Branch: branch, Path: path, ExpectedParentOID: expectedParentOID}, nil
}

func validateWorkIdentifier(kind, value string) error {
	if len(value) < 1 || len(value) > 63 || !isASCII(value) || !identifierPattern.MatchString(value) {
		return WorkInvalidArguments(fmt.Sprintf("%s ID %q must match %s", kind, value, identifierRule))
	}
	return nil
}

func validateWorkText(kind, value string, maximum int) error {
	if !utf8.ValidString(value) || value == "" || strings.TrimSpace(value) != value || utf8.RuneCountInString(value) > maximum {
		return WorkInvalidArguments(fmt.Sprintf("%s must be valid UTF-8, contain 1..%d characters, and have no surrounding whitespace", kind, maximum))
	}
	for _, character := range value {
		if unicode.IsControl(character) {
			return WorkInvalidArguments(kind + " must be one line without control characters")
		}
	}
	return nil
}

func validOIDText(value string) bool {
	if len(value) != 40 && len(value) != 64 {
		return false
	}
	for _, c := range value {
		if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

func validFullBranchRef(ref string) bool {
	if !strings.HasPrefix(ref, "refs/heads/") {
		return false
	}
	return validRefText(ref)
}

func validShortBranchText(branch string) bool {
	return branch != "" && branch != "@" && !strings.HasPrefix(branch, "-") && !strings.HasPrefix(branch, "refs/") && validFullBranchRef("refs/heads/"+branch)
}

func validRefText(ref string) bool {
	if ref == "" || strings.HasSuffix(ref, "/") || strings.HasSuffix(ref, ".") || strings.Contains(ref, "//") || strings.Contains(ref, "..") || strings.Contains(ref, "@{") {
		return false
	}
	for _, character := range ref {
		if character <= ' ' || character == 0x7f || strings.ContainsRune("~^:?*[\\", character) {
			return false
		}
	}
	for _, component := range strings.Split(ref, "/") {
		if component == "" || strings.HasPrefix(component, ".") || strings.HasSuffix(component, ".lock") {
			return false
		}
	}
	return true
}

type WorkItemIDSource interface {
	NewWorktreeID() (WorktreeID, error)
	NewOperationID() (WorktreeOperationID, error)
}
type cryptoWorkItemIDSource struct{ Reader io.Reader }

func (source cryptoWorkItemIDSource) reader() io.Reader {
	if source.Reader != nil {
		return source.Reader
	}
	return rand.Reader
}
func (source cryptoWorkItemIDSource) NewWorktreeID() (WorktreeID, error) {
	value, err := randomWorkItemID(source.reader(), "wt_")
	return WorktreeID(value), err
}
func (source cryptoWorkItemIDSource) NewOperationID() (WorktreeOperationID, error) {
	value, err := randomWorkItemID(source.reader(), "wop_")
	return WorktreeOperationID(value), err
}
func randomWorkItemID(reader io.Reader, prefix string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := io.ReadFull(reader, bytes); err != nil {
		return "", workError(ErrorWorkIO, "generate work-item identity", err)
	}
	return prefix + hex.EncodeToString(bytes), nil
}

func containingWorkItemWorkspace(dependencies Dependencies) (string, error) {
	if dependencies.Files == nil {
		return "", workError(ErrorWorkIO, "filesystem dependency is required", nil)
	}
	cwd, err := projectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return "", mapProjectWorkError(err)
	}
	root, err := locateContainingWorkspace(dependencies.Files, cwd)
	if err != nil {
		return "", mapProjectWorkError(err)
	}
	return root, nil
}

func mapProjectWorkError(err error) error {
	var project *ProjectError
	if errors.As(err, &project) {
		switch project.Class {
		case ErrorWorkspaceNotFound:
			return err
		case ErrorProjectWorkspace:
			return workError(ErrorWorkWorkspaceConflict, project.Detail, err)
		default:
			return workError(ErrorWorkProjectConflict, project.Detail, err)
		}
	}
	return workError(ErrorWorkIO, err.Error(), err)
}

func projectAndRepo(snapshot ProjectSnapshot, projectID ProjectID, repoID RepoID) (ProjectRecord, RepoRecord, error) {
	var project ProjectRecord
	foundProject := false
	for _, candidate := range snapshot.Projects {
		if candidate.ID == projectID {
			project = candidate
			foundProject = true
			break
		}
	}
	if !foundProject {
		return ProjectRecord{}, RepoRecord{}, workError(ErrorWorkNotFound, fmt.Sprintf("Project %s is not registered", projectID), nil)
	}
	member := false
	for _, id := range project.RepoIDs {
		if id == repoID {
			member = true
			break
		}
	}
	if !member {
		return ProjectRecord{}, RepoRecord{}, workError(ErrorWorkProjectConflict, fmt.Sprintf("repository %s is not a member of Project %s", repoID, projectID), nil)
	}
	for _, repository := range snapshot.Repos {
		if repository.ID == repoID {
			return project, repository, nil
		}
	}
	return ProjectRecord{}, RepoRecord{}, workError(ErrorWorkProjectConflict, fmt.Sprintf("repository %s is missing from the Project registry", repoID), nil)
}

func canonicalAbsentTarget(files FileSystem, cwd, value string) (string, error) {
	abs := value
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(cwd, abs)
	}
	abs = filepath.Clean(abs)
	leaf := filepath.Base(abs)
	if leaf == "." || leaf == string(filepath.Separator) || leaf == "" {
		return "", workError(ErrorWorkPathConflict, fmt.Sprintf("target path %s has no ordinary leaf", value), nil)
	}
	parent, err := canonicalDirectory(files, "", filepath.Dir(abs), "Task worktree parent")
	if err != nil {
		return "", workError(ErrorWorkPathConflict, err.Error(), err)
	}
	target := filepath.Join(parent, leaf)
	if info, err := files.Lstat(target); err == nil {
		return "", workError(ErrorWorkPathConflict, fmt.Sprintf("target path %s already exists as %s", target, info.Mode().Type()), nil)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return "", workError(ErrorWorkPathConflict, fmt.Sprintf("cannot inspect target path %s", target), err)
	}
	return target, nil
}

func intentDigest(operation WorktreeOperationRecord) (string, error) {
	value := canonicaljson.Object{
		{Name: "git_common_dir", Value: operation.GitCommonDir}, {Name: "kind", Value: "WorkspaceTaskWorktreeIntent@1"},
		{Name: "parent_epic_id", Value: string(operation.ParentEpicID)}, {Name: "parent_oid", Value: operation.ParentOID},
		{Name: "parent_ref", Value: operation.ParentRef}, {Name: "parent_tree", Value: operation.ParentTree},
		{Name: "parent_worktree_id", Value: string(operation.ParentWorktreeID)}, {Name: "project_id", Value: string(operation.ProjectID)},
		{Name: "repo_id", Value: string(operation.RepoID)}, {Name: "source_ref", Value: operation.SourceRef},
		{Name: "target_locator", Value: operation.TargetLocator}, {Name: "task_id", Value: string(operation.TaskID)},
	}
	bytes, err := canonicaljson.Marshal(value)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(bytes)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func findEpic(registry WorkItemRegistry, id EpicID) (*EpicRecord, int) {
	for i := range registry.Epics {
		if registry.Epics[i].ID == id {
			return &registry.Epics[i], i
		}
	}
	return nil, -1
}
func findTask(registry WorkItemRegistry, id TaskID) (*TaskRecord, int) {
	for i := range registry.Tasks {
		if registry.Tasks[i].ID == id {
			return &registry.Tasks[i], i
		}
	}
	return nil, -1
}
func findOperation(registry WorkItemRegistry, taskID TaskID) (*WorktreeOperationRecord, int) {
	for i := range registry.WorktreeOperations {
		if registry.WorktreeOperations[i].TaskID == taskID {
			return &registry.WorktreeOperations[i], i
		}
	}
	return nil, -1
}
func findEpicRepo(epic EpicRecord, repoID RepoID) (*EpicRepoBinding, int) {
	for i := range epic.RepoBindings {
		if epic.RepoBindings[i].RepoID == repoID {
			return &epic.RepoBindings[i], i
		}
	}
	return nil, -1
}

func sortWorkRegistry(registry *WorkItemRegistry) {
	sortTaskContentRegistry(registry)
	sortQueueRegistry(registry)
	sort.Slice(registry.Epics, func(i, j int) bool { return registry.Epics[i].ID < registry.Epics[j].ID })
	for i := range registry.Epics {
		sort.Slice(registry.Epics[i].RepoBindings, func(a, b int) bool {
			return registry.Epics[i].RepoBindings[a].RepoID < registry.Epics[i].RepoBindings[b].RepoID
		})
	}
	sort.Slice(registry.Tasks, func(i, j int) bool { return registry.Tasks[i].ID < registry.Tasks[j].ID })
	sort.Slice(registry.WorktreeOperations, func(i, j int) bool { return registry.WorktreeOperations[i].ID < registry.WorktreeOperations[j].ID })
	sort.Slice(registry.TaskResults, func(i, j int) bool { return registry.TaskResults[i].ID < registry.TaskResults[j].ID })
	sort.Slice(registry.HumanQARecords, func(i, j int) bool { return registry.HumanQARecords[i].ID < registry.HumanQARecords[j].ID })
	sort.Slice(registry.IntegrationAuthorities, func(i, j int) bool {
		return registry.IntegrationAuthorities[i].ID < registry.IntegrationAuthorities[j].ID
	})
	sort.Slice(registry.IntegrationIntents, func(i, j int) bool { return registry.IntegrationIntents[i].ID < registry.IntegrationIntents[j].ID })
	sort.Slice(registry.IntegrationAttempts, func(i, j int) bool { return registry.IntegrationAttempts[i].ID < registry.IntegrationAttempts[j].ID })
	sort.Slice(registry.IntegrationResults, func(i, j int) bool { return registry.IntegrationResults[i].ID < registry.IntegrationResults[j].ID })
}

func AdoptEpic(dependencies Dependencies, input EpicAdoptInput) (EpicMutationResult, error) {
	if _, err := ParseEpicAdoptInput(string(input.EpicID), input.Title, string(input.ProjectID), string(input.RepoID), input.Worktree, input.Ref, input.ExpectedOID); err != nil {
		return EpicMutationResult{}, err
	}
	if err := requireWorkMutationDependencies(dependencies); err != nil {
		return EpicMutationResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return EpicMutationResult{}, err
	}
	cwd, err := projectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return EpicMutationResult{}, mapProjectWorkError(err)
	}
	locator, err := canonicalDirectory(dependencies.Files, cwd, input.Worktree, "Epic worktree")
	if err != nil {
		return EpicMutationResult{}, workError(ErrorWorkPathConflict, err.Error(), err)
	}
	result := EpicMutationResult{Workspace: root}
	err = dependencies.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		_, repository, err := projectAndRepo(projects, input.ProjectID, input.RepoID)
		if err != nil {
			return err
		}
		return dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			if existing, _ := findEpic(registry, input.EpicID); existing != nil {
				binding, _ := findEpicRepo(*existing, input.RepoID)
				if existing.Title == input.Title && existing.ProjectID == input.ProjectID && binding != nil && binding.GitCommonDir == repository.GitCommonDir && binding.Worktree.Locator == locator && binding.Worktree.Ref == input.Ref && binding.Worktree.OID == input.ExpectedOID {
					observation, validationErr := validateEpicAdoptionObservation(dependencies, repository, locator, input.Ref, input.ExpectedOID)
					if validationErr != nil {
						return validationErr
					}
					if observation.Tree != binding.Worktree.Tree {
						return workError(ErrorWorkParentStale, fmt.Sprintf("Epic worktree %s tree differs from the recorded binding", locator), nil)
					}
					result.Epic = *existing
					return nil
				}
				return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Epic %s is already registered with different identity or binding", input.EpicID), nil)
			}
			observation, err := validateEpicAdoptionObservation(dependencies, repository, locator, input.Ref, input.ExpectedOID)
			if err != nil {
				return err
			}
			for _, epic := range registry.Epics {
				for _, binding := range epic.RepoBindings {
					if binding.GitCommonDir == repository.GitCommonDir && (binding.Worktree.Locator == locator || binding.Worktree.Ref == input.Ref) {
						return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Epic worktree or ref is already owned by Epic %s", epic.ID), nil)
					}
				}
			}
			id, err := dependencies.WorkIDs.NewWorktreeID()
			if err != nil {
				return err
			}
			epic := EpicRecord{ID: input.EpicID, Title: input.Title, ProjectID: input.ProjectID, RepoBindings: []EpicRepoBinding{{RepoID: input.RepoID, GitCommonDir: repository.GitCommonDir, Worktree: EpicWorktreeBinding{ID: id, OwnerKind: "epic", OwnerID: input.EpicID, Origin: "adopted", Locator: locator, Ref: input.Ref, OID: input.ExpectedOID, Tree: observation.Tree, GitCommonDir: repository.GitCommonDir}}}}
			registry.Epics = append(registry.Epics, epic)
			if registry.FormatVersion == 4 {
				registry.EpicBaseVersions = append(registry.EpicBaseVersions, adoptedEpicBase(epic, epic.RepoBindings[0]))
			}
			sortWorkRegistry(&registry)
			if err := session.Publish(registry); err != nil {
				return err
			}
			result.Epic = epic
			result.Created = true
			return nil
		})
	})
	if err != nil {
		return EpicMutationResult{}, mapMutationError(err)
	}
	return result, nil
}

func validateEpicAdoptionObservation(dependencies Dependencies, repository RepoRecord, locator, ref, expectedOID string) (GitWorktreeObservation, error) {
	repoObservation, err := dependencies.WorkGit.ObserveRepo(repository)
	if err != nil {
		return GitWorktreeObservation{}, err
	}
	if repoObservation.Locator != repository.Locator || repoObservation.GitCommonDir != repository.GitCommonDir {
		return GitWorktreeObservation{}, workError(ErrorWorkProjectConflict, fmt.Sprintf("current repository %s binding differs from Project registry", repository.ID), nil)
	}
	observation, err := dependencies.WorkGit.ObserveWorktree(locator)
	if err != nil {
		return GitWorktreeObservation{}, err
	}
	if observation.Ref != ref || observation.OID != expectedOID || observation.GitCommonDir != repository.GitCommonDir || !observation.InventoryMatch {
		return GitWorktreeObservation{}, workError(ErrorWorkGitObservation, fmt.Sprintf("Epic worktree %s does not match ref, commit, common directory, and inventory", locator), nil)
	}
	if !observation.Clean {
		return GitWorktreeObservation{}, workError(ErrorWorkParentDirty, fmt.Sprintf("Epic worktree %s is not clean", locator), nil)
	}
	refObservation, err := dependencies.WorkGit.ObserveRef(repository, ref)
	if err != nil {
		return GitWorktreeObservation{}, err
	}
	if !refObservation.Exists || refObservation.OID != expectedOID || refObservation.Tree != observation.Tree || refObservation.CheckedOutAt != locator {
		return GitWorktreeObservation{}, workError(ErrorWorkParentStale, fmt.Sprintf("Epic ref %s does not match expected commit %s", ref, expectedOID), nil)
	}
	return observation, nil
}

func CreateTask(dependencies Dependencies, input TaskCreateInput) (TaskMutationResult, error) {
	if _, err := ParseTaskCreateInput(string(input.TaskID), input.Title, input.Description, string(input.ParentEpicID), string(input.ProjectID), string(input.RepoID)); err != nil {
		return TaskMutationResult{}, err
	}
	if err := requireWorkMutationDependencies(dependencies); err != nil {
		return TaskMutationResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return TaskMutationResult{}, err
	}
	result := TaskMutationResult{Workspace: root}
	err = dependencies.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		_, repository, err := projectAndRepo(projects, input.ProjectID, input.RepoID)
		if err != nil {
			return err
		}
		return dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			epic, _ := findEpic(registry, input.ParentEpicID)
			if epic == nil {
				return workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s is not registered", input.ParentEpicID), nil)
			}
			binding, _ := findEpicRepo(*epic, input.RepoID)
			if epic.ProjectID != input.ProjectID || binding == nil || binding.GitCommonDir != repository.GitCommonDir {
				return workError(ErrorWorkProjectConflict, fmt.Sprintf("Task %s Project/repository does not match Epic %s", input.TaskID, input.ParentEpicID), nil)
			}
			if existing, _ := findTask(registry, input.TaskID); existing != nil {
				if existing.Title == input.Title && existing.Description == input.Description && existing.ParentEpicID == input.ParentEpicID && existing.ProjectID == input.ProjectID && existing.RepoID == input.RepoID && existing.GitCommonDir == repository.GitCommonDir {
					result.Task = *existing
					result.Epic = *epic
					if pub := taskPublication(registry, "task-create/"+string(input.TaskID)); pub != nil {
						content := newTaskContentMutation(root, input.TaskID, pub.PublicationKey, pub.IntentSHA256)
						content.PreState = taskContentState(registry, input.TaskID)
						content.PostState = content.PreState
						content.Classification = "existing"
						content.OutcomeRef = &pub.OutcomeRef
						result.Content = &content
						if err := validateTaskContentClosure(dependencies.TaskContent, root, registry, true); err != nil {
							return err
						}
						return syncTaskRegistry(root)
					}
					return nil
				}
				return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Task %s already exists with different fields", input.TaskID), nil)
			}
			task := TaskRecord{ID: input.TaskID, Title: input.Title, Description: input.Description, ParentEpicID: input.ParentEpicID, ProjectID: input.ProjectID, RepoID: input.RepoID, GitCommonDir: repository.GitCommonDir, WorktreeState: WorkItemUnbound, Worktree: nil}
			content, err := createTaskProblem(dependencies, root, session, registry, projects, task, input.RegistryUpgrade)
			result.Content = &content
			if err != nil {
				return err
			}
			result.Task = task
			result.Epic = *epic
			result.Created = true
			return nil
		})
	})
	if result.Content != nil {
		finished, finishErr := finishTaskContentMutation(dependencies, *result.Content, err)
		result.Content = &finished
		if finishErr != nil {
			return result, finishErr
		}
	}
	if err != nil {
		return result, mapMutationError(err)
	}
	return result, nil
}

func CreateTaskWorktree(dependencies Dependencies, input TaskWorktreeCreateInput) (TaskWorktreeMutationResult, error) {
	if _, err := ParseTaskWorktreeCreateInput(string(input.TaskID), input.Branch, input.Path, input.ExpectedParentOID); err != nil {
		return TaskWorktreeMutationResult{}, err
	}
	if err := requireWorkMutationDependencies(dependencies); err != nil {
		return TaskWorktreeMutationResult{}, err
	}
	if err := dependencies.WorkGit.ValidateBranch(input.Branch); err != nil {
		return TaskWorktreeMutationResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return TaskWorktreeMutationResult{}, err
	}
	cwd, err := projectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return TaskWorktreeMutationResult{}, mapProjectWorkError(err)
	}
	target, targetExists, err := canonicalTaskTarget(dependencies.Files, cwd, input.Path)
	if err != nil {
		return TaskWorktreeMutationResult{}, err
	}
	result := TaskWorktreeMutationResult{Workspace: root}
	err = dependencies.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			return createTaskWorktreeLocked(dependencies, root, projects, session, registry, input, target, targetExists, &result, nil)
		})
	})
	if err != nil {
		return TaskWorktreeMutationResult{}, mapMutationError(err)
	}
	return result, nil
}

func requireWorkMutationDependencies(dependencies Dependencies) error {
	if dependencies.Files == nil || dependencies.ProjectLocks == nil || dependencies.WorkItems == nil || dependencies.WorkGit == nil || dependencies.WorkIDs == nil {
		return workError(ErrorWorkIO, "filesystem, Project lock, work-item store, Git, and ID dependencies are required", nil)
	}
	return nil
}
func mapMutationError(err error) error {
	var work *WorkItemError
	if errors.As(err, &work) {
		return err
	}
	var project *ProjectError
	if errors.As(err, &project) {
		return mapProjectWorkError(err)
	}
	return workError(ErrorWorkIO, err.Error(), err)
}

func canonicalTaskTarget(files FileSystem, cwd, value string) (string, bool, error) {
	abs := value
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(cwd, abs)
	}
	abs = filepath.Clean(abs)
	leaf := filepath.Base(abs)
	if leaf == "" || leaf == "." || leaf == string(filepath.Separator) {
		return "", false, workError(ErrorWorkPathConflict, "Task worktree path has no ordinary leaf", nil)
	}
	parent, err := canonicalDirectory(files, "", filepath.Dir(abs), "Task worktree parent")
	if err != nil {
		return "", false, workError(ErrorWorkPathConflict, err.Error(), err)
	}
	target := filepath.Join(parent, leaf)
	info, err := files.Lstat(target)
	if errors.Is(err, fs.ErrNotExist) {
		return target, false, nil
	}
	if err != nil {
		return "", false, workError(ErrorWorkPathConflict, fmt.Sprintf("cannot inspect target path %s", target), err)
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		return "", false, workError(ErrorWorkPathConflict, fmt.Sprintf("target path %s is not an ordinary directory or is a symlink", target), nil)
	}
	return target, true, nil
}

func precheckParentAndResources(dependencies Dependencies, repository RepoRecord, binding EpicRepoBinding, sourceRef, target string) error {
	repo, err := dependencies.WorkGit.ObserveRepo(repository)
	if err != nil {
		return err
	}
	if repo.GitCommonDir != repository.GitCommonDir || repo.Locator != repository.Locator {
		return workError(ErrorWorkProjectConflict, "registered repository binding is stale", nil)
	}
	parent, err := dependencies.WorkGit.ObserveWorktree(binding.Worktree.Locator)
	if err != nil {
		return err
	}
	if !parent.Clean {
		return workError(ErrorWorkParentDirty, fmt.Sprintf("Epic parent worktree %s is not clean", binding.Worktree.Locator), nil)
	}
	if parent.Ref != binding.Worktree.Ref || parent.OID != binding.Worktree.OID || parent.Tree != binding.Worktree.Tree || parent.GitCommonDir != binding.GitCommonDir || !parent.InventoryMatch {
		return workError(ErrorWorkParentStale, "Epic parent worktree differs from the recorded binding", nil)
	}
	parentRef, err := dependencies.WorkGit.ObserveRef(repository, binding.Worktree.Ref)
	if err != nil {
		return err
	}
	if !parentRef.Exists || parentRef.OID != binding.Worktree.OID || parentRef.Tree != binding.Worktree.Tree || parentRef.CheckedOutAt != binding.Worktree.Locator {
		return workError(ErrorWorkParentStale, "Epic parent ref differs from the recorded binding", nil)
	}
	source, err := dependencies.WorkGit.ObserveRef(repository, sourceRef)
	if err != nil {
		return err
	}
	if source.Exists {
		return workError(ErrorWorkRefConflict, fmt.Sprintf("source ref %s already exists", sourceRef), nil)
	}
	if _, err := dependencies.Files.Lstat(target); !errors.Is(err, fs.ErrNotExist) {
		if err == nil {
			return workError(ErrorWorkPathConflict, fmt.Sprintf("target path %s already exists", target), nil)
		}
		return workError(ErrorWorkPathConflict, fmt.Sprintf("cannot inspect target path %s", target), err)
	}
	inventory, err := dependencies.WorkGit.ListWorktrees(repository)
	if err != nil {
		return err
	}
	if _, found := worktreeEntry(inventory, target); found {
		return workError(ErrorWorkPathConflict, fmt.Sprintf("target path %s is already in Git worktree inventory", target), nil)
	}
	return nil
}

func emptyCreateObservation(classification string) WorktreeObservationRecord {
	return WorktreeObservationRecord{Classification: classification, TargetKind: "absent", Reasons: []string{}}
}
func stringPointer(value string) *string { return &value }
func boolPointer(value bool) *bool       { return &value }

func observeTaskCreate(dependencies Dependencies, repository RepoRecord, binding EpicRepoBinding, operation WorktreeOperationRecord) (string, WorktreeObservationRecord) {
	observation := emptyCreateObservation("partial_or_unknown")
	reasons := []string{}
	parentRef, refErr := dependencies.WorkGit.ObserveRef(repository, operation.ParentRef)
	parent, parentErr := dependencies.WorkGit.ObserveWorktree(binding.Worktree.Locator)
	source, sourceErr := dependencies.WorkGit.ObserveRef(repository, operation.SourceRef)
	_, inventoryErr := dependencies.WorkGit.ListWorktrees(repository)
	if refErr == nil && parentRef.Exists {
		observation.ParentRefOID = stringPointer(parentRef.OID)
		observation.ParentRefTree = stringPointer(parentRef.Tree)
	} else {
		reasons = append(reasons, "parent_observation_unknown")
	}
	parentExact := false
	if parentErr == nil {
		observation.ParentWorktreeLocator = stringPointer(parent.Locator)
		observation.ParentWorktreeRef = stringPointer(parent.Ref)
		observation.ParentWorktreeOID = stringPointer(parent.OID)
		observation.ParentWorktreeTree = stringPointer(parent.Tree)
		observation.ParentWorktreeGitCommonDir = stringPointer(parent.GitCommonDir)
		observation.ParentWorktreeClean = boolPointer(parent.Clean)
		parentExact = parent.Locator == binding.Worktree.Locator && parent.Ref == operation.ParentRef && parent.OID == operation.ParentOID && parent.Tree == operation.ParentTree && parent.GitCommonDir == operation.GitCommonDir && parent.Clean && parent.InventoryMatch
	} else {
		reasons = append(reasons, "parent_observation_unknown")
	}
	refExact := refErr == nil && parentRef.Exists && parentRef.OID == operation.ParentOID && parentRef.Tree == operation.ParentTree
	sourceAbsent := sourceErr == nil && !source.Exists
	branchOnly := false
	if sourceErr == nil && source.Exists {
		observation.SourceRefOID = stringPointer(source.OID)
		observation.SourceRefTree = stringPointer(source.Tree)
		if source.CheckedOutAt != "" {
			observation.SourceRefCheckedOutAt = stringPointer(source.CheckedOutAt)
		}
		branchOnly = source.OID == operation.ParentOID && source.Tree == operation.ParentTree && source.CheckedOutAt == ""
	} else if sourceErr != nil {
		reasons = append(reasons, "source_ref_stale")
	}
	info, pathErr := dependencies.Files.Lstat(operation.TargetLocator)
	targetAbsent := errors.Is(pathErr, fs.ErrNotExist)
	targetExact := false
	if targetAbsent {
		observation.TargetKind = "absent"
	} else if pathErr != nil {
		observation.TargetKind = "unknown"
		reasons = append(reasons, "target_observation_unknown")
	} else if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		observation.TargetKind = "other"
		observation.TargetLocator = stringPointer(operation.TargetLocator)
		reasons = append(reasons, "target_worktree_stale")
	} else {
		observation.TargetKind = "worktree"
		target, targetErr := dependencies.WorkGit.ObserveWorktree(operation.TargetLocator)
		if targetErr != nil {
			reasons = append(reasons, "target_observation_unknown")
		} else {
			observation.TargetLocator = stringPointer(target.Locator)
			observation.TargetRef = stringPointer(target.Ref)
			observation.TargetOID = stringPointer(target.OID)
			observation.TargetTree = stringPointer(target.Tree)
			observation.TargetGitCommonDir = stringPointer(target.GitCommonDir)
			observation.TargetClean = boolPointer(target.Clean)
			observation.InventoryMatch = boolPointer(target.InventoryMatch)
			targetExact = target.Locator == operation.TargetLocator && target.Ref == operation.SourceRef && target.OID == operation.ParentOID && target.Tree == operation.ParentTree && target.GitCommonDir == operation.GitCommonDir && target.Clean && target.InventoryMatch
		}
	}
	if inventoryErr != nil {
		reasons = append(reasons, "worktree_inventory_unknown")
	}
	sort.Strings(reasons)
	observation.Reasons = uniqueReasonStrings(reasons)
	switch {
	case parentExact && refExact && sourceErr == nil && source.Exists && source.OID == operation.ParentOID && source.Tree == operation.ParentTree && source.CheckedOutAt == operation.TargetLocator && targetExact:
		observation.Classification = "exact_effect"
		observation.Reasons = []string{}
		return "exact_effect", observation
	case parentExact && refExact && sourceAbsent && targetAbsent && inventoryErr == nil:
		observation.Classification = "no_effect"
		observation.Reasons = []string{}
		return "no_effect", observation
	case parentExact && refExact && branchOnly && targetAbsent && inventoryErr == nil:
		observation.Classification = "branch_only"
		observation.Reasons = []string{}
		return "branch_only", observation
	default:
		observation.Classification = "partial_or_unknown"
		if len(observation.Reasons) == 0 {
			observation.Reasons = []string{"task_worktree_reconciliation_required"}
		}
		return "partial_or_unknown", observation
	}
}

func uniqueReasonStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	result := []string{values[0]}
	for _, value := range values[1:] {
		if value != result[len(result)-1] {
			result = append(result, value)
		}
	}
	return result
}
func finalizeReady(session WorkItemStoreSession, registry *WorkItemRegistry, taskIndex, operationIndex int, epic EpicRecord, operation WorktreeOperationRecord, observation WorktreeObservationRecord) error {
	registry.WorktreeOperations[operationIndex].State = "ready"
	registry.WorktreeOperations[operationIndex].LastObservation = observation
	registry.Tasks[taskIndex].WorktreeState = WorkItemReady
	registry.Tasks[taskIndex].Worktree = &TaskWorktreeBinding{ID: operation.WorktreeID, OwnerKind: "task", OwnerID: operation.TaskID, Origin: "ply_created", Locator: operation.TargetLocator, Ref: operation.SourceRef, OID: operation.ParentOID, Tree: operation.ParentTree, GitCommonDir: operation.GitCommonDir, ParentEpicID: operation.ParentEpicID, ParentWorktreeID: operation.ParentWorktreeID, ParentRef: operation.ParentRef, ParentOID: operation.ParentOID, ParentTree: operation.ParentTree, OperationID: operation.ID, IntentDigest: operation.IntentDigest}
	for i := range registry.TaskPreparations {
		p := &registry.TaskPreparations[i]
		if p.OperationID == operation.ID && p.Outcome == nil {
			p.Outcome = &TaskPreparationOutcome{Kind: "prepared", WorktreeID: operation.WorktreeID, ObservedOID: operation.ParentOID, ObservedTree: operation.ParentTree, RecordedAtUTC: p.CreatedAtUTC}
		}
	}
	return session.Publish(*registry)
}
func persistReconciliation(session WorkItemStoreSession, registry *WorkItemRegistry, taskIndex, operationIndex int, observation WorktreeObservationRecord) error {
	registry.Tasks[taskIndex].WorktreeState = WorkItemReconciliationRequired
	registry.Tasks[taskIndex].Worktree = nil
	registry.WorktreeOperations[operationIndex].State = "reconciliation_required"
	observation.Classification = "partial_or_unknown"
	if len(observation.Reasons) == 0 {
		observation.Reasons = []string{"task_worktree_reconciliation_required"}
	}
	registry.WorktreeOperations[operationIndex].LastObservation = observation
	return session.Publish(*registry)
}
func reconciliationError(taskID TaskID) error {
	return workError(ErrorWorkReconciliation, fmt.Sprintf("Task %s worktree state requires reconciliation. Run `ply workspace task show %s --format json` to inspect the preserved intent and observation.", taskID, taskID), nil)
}

func ListEpics(dependencies Dependencies) (EpicListResult, error) {
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return EpicListResult{}, err
	}
	if dependencies.WorkItems == nil {
		return EpicListResult{}, workError(ErrorWorkIO, "work-item store dependency is required", nil)
	}
	registry, err := dependencies.WorkItems.Snapshot(root)
	if err != nil {
		return EpicListResult{}, err
	}
	return EpicListResult{Workspace: root, Epics: append([]EpicRecord(nil), registry.Epics...)}, nil
}
func ListTasks(dependencies Dependencies, epicID *EpicID) (TaskListResult, error) {
	return ListTasksWithFilters(dependencies, TaskListFilters{EpicID: epicID})
}

// ListTasksWithFilters reads registrations, independently of queue membership and
// checkout availability. In particular, it does not observe Git or take write locks.
func ListTasksWithFilters(dependencies Dependencies, filters TaskListFilters) (TaskListResult, error) {
	if err := filters.Validate(); err != nil {
		return TaskListResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return TaskListResult{}, err
	}
	if dependencies.WorkItems == nil {
		return TaskListResult{}, workError(ErrorWorkIO, "work-item store dependency is required", nil)
	}
	registry, err := dependencies.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		return TaskListResult{}, err
	}
	if err := validateTaskListProjectFilters(dependencies, root, filters); err != nil {
		return TaskListResult{}, err
	}
	lifecycles, err := readWorkItemLifecycle(root)
	if err != nil {
		return TaskListResult{}, err
	}
	result := TaskListResult{Workspace: root, Filters: filters, CurrentTitles: map[TaskID]string{}, Titles: map[TaskID]TaskListTitle{}, Lifecycles: lifecycles}
	epicID := filters.EpicID
	if epicID != nil {
		epic, _ := findEpic(registry, *epicID)
		if epic == nil {
			return TaskListResult{}, workError(ErrorWorkNotFound, fmt.Sprintf("Epic %s is not registered", *epicID), nil)
		}
		copyEpic := *epic
		result.Epic = &copyEpic
	}
	for _, task := range registry.Tasks {
		if (epicID == nil || task.ParentEpicID == *epicID) &&
			(filters.ProjectID == nil || task.ProjectID == *filters.ProjectID) &&
			(filters.RepoID == nil || task.RepoID == *filters.RepoID) {
			result.Tasks = append(result.Tasks, task)
			title := readTaskListTitle(dependencies, root, registry, task)
			if title.Source == "problem_revision" {
				if title.Title != nil {
					result.CurrentTitles[task.ID] = *title.Title
				} else {
					result.CurrentTitles[task.ID] = "Problem content unavailable"
				}
			}
			result.Titles[task.ID] = title
		}
	}
	return result, nil
}

// The caller owns Project -> work-items locks. Reservation joins the first intent write.
func createTaskWorktreeLocked(dependencies Dependencies, root string, projects ProjectSnapshot, session WorkItemStoreSession, registry WorkItemRegistry, input TaskWorktreeCreateInput, target string, targetExists bool, result *TaskWorktreeMutationResult, reserve func(*WorkItemRegistry, WorktreeOperationRecord) error) error {
	task, taskIndex := findTask(registry, input.TaskID)
	if task == nil {
		return workError(ErrorWorkNotFound, fmt.Sprintf("Task %s is not registered", input.TaskID), nil)
	}
	epic, _ := findEpic(registry, task.ParentEpicID)
	if epic == nil {
		return workError(ErrorWorkStoreConflict, fmt.Sprintf("Task %s refers to missing Epic %s", task.ID, task.ParentEpicID), nil)
	}
	adopted, _ := findEpicRepo(*epic, task.RepoID)
	binding := adopted
	if adopted != nil {
		current := currentEpicBinding(registry, *epic, *adopted)
		binding = &current
	}
	if binding == nil {
		return workError(ErrorWorkStoreConflict, fmt.Sprintf("Task %s has no matching Epic repository binding", task.ID), nil)
	}
	_, repository, err := projectAndRepo(projects, task.ProjectID, task.RepoID)
	if err != nil {
		return err
	}
	if repository.GitCommonDir != task.GitCommonDir || binding.GitCommonDir != task.GitCommonDir {
		return workError(ErrorWorkProjectConflict, fmt.Sprintf("Task %s repository common directory changed", task.ID), nil)
	}
	sourceRef := "refs/heads/" + input.Branch
	operation, operationIndex := findOperation(registry, task.ID)
	intentCreated := false
	if prep := preparationForTask(registry, task.ID); prep != nil {
		disposition := preparationDisposition(registry, prep.ID)
		if disposition != "current" || prep.Outcome != nil {
			result.Task = *task
			result.Epic = *epic
			if operation != nil {
				result.Operation = *operation
			}
			result.Outcome = "historical_" + disposition
			return nil
		}
		if reserve == nil {
			return queueError("task_preparation_recovery_required", "use ply workspace task prepare --next --apply --confirm "+prep.PlanSHA256+" to recover the preserved preparation")
		}
	}
	if operation != nil {
		historical, ok := historicalEpicBinding(registry, *epic, *adopted, operation.ParentOID, operation.ParentTree)
		if !ok {
			return queueError("epic_base_conflict", "intent parent is not a historical base")
		}
		binding = &historical
	}
	if pendingBaseUpdate(registry, task.ParentEpicID, task.RepoID) {
		return queueError("epic_base_pending", "resolve the base update before creating a worktree")
	}

	if operation != nil {
		if operation.SourceRef != sourceRef || operation.TargetLocator != target || operation.ParentOID != input.ExpectedParentOID {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Task %s already has a different worktree intent", task.ID), nil)
		}
		if task.WorktreeState == WorkItemReconciliationRequired && reserve == nil {
			return reconciliationError(task.ID)
		}
	} else if task.WorktreeState != WorkItemUnbound || task.Worktree != nil {
		return workError(ErrorWorkStoreConflict, fmt.Sprintf("Task %s has inconsistent worktree state", task.ID), nil)
	}
	if operation == nil {
		if targetExists {
			return workError(ErrorWorkPathConflict, fmt.Sprintf("target path %s already exists", target), nil)
		}
		if input.ExpectedParentOID != binding.Worktree.OID {
			return workError(ErrorWorkParentStale, fmt.Sprintf("expected parent %s differs from recorded Epic commit %s", input.ExpectedParentOID, binding.Worktree.OID), nil)
		}
		if err := precheckParentAndResources(dependencies, repository, *binding, sourceRef, target); err != nil {
			return err
		}
		for _, candidate := range registry.WorktreeOperations {
			if candidate.GitCommonDir == task.GitCommonDir && (candidate.SourceRef == sourceRef || candidate.TargetLocator == target) {
				return workError(ErrorWorkIdentityConflict, fmt.Sprintf("branch or path is reserved by Task %s", candidate.TaskID), nil)
			}
		}
		worktreeID, err := dependencies.WorkIDs.NewWorktreeID()
		if err != nil {
			return err
		}
		operationID, err := dependencies.WorkIDs.NewOperationID()
		if err != nil {
			return err
		}
		created := WorktreeOperationRecord{ID: operationID, TaskID: task.ID, WorktreeID: worktreeID, State: "creating", ProjectID: task.ProjectID, RepoID: task.RepoID, GitCommonDir: task.GitCommonDir, ParentEpicID: task.ParentEpicID, ParentWorktreeID: binding.Worktree.ID, ParentRef: binding.Worktree.Ref, ParentOID: binding.Worktree.OID, ParentTree: binding.Worktree.Tree, SourceRef: sourceRef, TargetLocator: target, LastObservation: emptyCreateObservation("no_effect")}
		digest, err := intentDigest(created)
		if err != nil {
			return workError(ErrorWorkIO, "compute Task worktree intent digest", err)
		}
		created.IntentDigest = digest
		registry.Tasks[taskIndex].WorktreeState = WorkItemCreating
		registry.WorktreeOperations = append(registry.WorktreeOperations, created)
		if reserve != nil {
			if err := reserve(&registry, created); err != nil {
				return err
			}
		}
		sortWorkRegistry(&registry)
		if err := session.Publish(registry); err != nil {
			return err
		}
		operation, operationIndex = findOperation(registry, task.ID)
		task = &registry.Tasks[taskIndex]
		intentCreated = true
	}
	classification, observation := observeTaskCreate(dependencies, repository, *binding, *operation)
	if classification == "exact_effect" {
		if task.WorktreeState == WorkItemReady && task.Worktree != nil && operation.State == "ready" {
			result.Task = *task
			result.Epic = *epic
			result.Operation = *operation
			result.Outcome = "already_ready"
			return nil
		}
		if err := finalizeReady(session, &registry, taskIndex, operationIndex, *epic, *operation, observation); err != nil {
			return err
		}
		result.Task = registry.Tasks[taskIndex]
		result.Epic = *epic
		result.Operation = registry.WorktreeOperations[operationIndex]
		result.Outcome = "recovered"
		return nil
	}
	if task.WorktreeState == WorkItemReady {
		return workError(ErrorWorkReconciliation, fmt.Sprintf("Task %s recorded ready but Git no longer matches. Run `ply workspace task show %s --format json` to inspect the preserved intent and observation.", task.ID, task.ID), nil)
	}
	currentBase := currentEpicBase(registry, *epic, *adopted)
	if currentBase.OID != operation.ParentOID || currentBase.Tree != operation.ParentTree {
		return queueError("epic_base_changed", "intent base is no longer current; no new Git effect is permitted")
	}
	gitInput := GitCreateInput{Repository: repository, Branch: input.Branch, SourceRef: sourceRef, TargetLocator: target, ParentOID: operation.ParentOID}
	var outcome GitCommandOutcome
	switch classification {
	case "no_effect":
		outcome = dependencies.WorkGit.CreateBranchAndWorktree(gitInput)
	case "branch_only":
		outcome = dependencies.WorkGit.AddWorktreeForExistingBranch(gitInput)
	default:
		if err := persistReconciliation(session, &registry, taskIndex, operationIndex, observation); err != nil {
			return err
		}
		return reconciliationError(task.ID)
	}
	postClass, postObservation := observeTaskCreate(dependencies, repository, *binding, *operation)
	switch postClass {
	case "exact_effect":
		if err := finalizeReady(session, &registry, taskIndex, operationIndex, *epic, *operation, postObservation); err != nil {
			return err
		}
		result.Task = registry.Tasks[taskIndex]
		result.Epic = *epic
		result.Operation = registry.WorktreeOperations[operationIndex]
		if outcome.Err == nil && intentCreated && classification == "no_effect" {
			result.Outcome = "created"
		} else {
			result.Outcome = "recovered"
		}
		return nil
	case "no_effect", "branch_only":
		registry.WorktreeOperations[operationIndex].LastObservation = postObservation
		if err := session.Publish(registry); err != nil {
			return err
		}
		return workError(ErrorWorkGitEffect, fmt.Sprintf("Git worktree creation for Task %s did not complete (exit %d); identical retry is safe", task.ID, outcome.Exit), outcome.Err)
	default:
		if err := persistReconciliation(session, &registry, taskIndex, operationIndex, postObservation); err != nil {
			return err
		}
		return reconciliationError(task.ID)
	}
}
