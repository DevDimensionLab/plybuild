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
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

type TaskResultID string
type HumanQARecordID string
type IntegrationAuthorityID string
type IntegrationIntentID string
type IntegrationAttemptID string
type IntegrationResultID string

type TaskArtifactRecord struct {
	ArtifactID string `yaml:"artifact_id" json:"artifact_id"`
	Role       string `yaml:"role" json:"role"`
	Locator    string `yaml:"locator" json:"locator"`
	SHA256     string `yaml:"sha256" json:"sha256"`
	SizeBytes  int64  `yaml:"size_bytes" json:"size_bytes"`
}

type TaskVerifierResultRecord struct {
	VerifierID       string   `yaml:"verifier_id" json:"verifier_id"`
	Outcome          string   `yaml:"outcome" json:"outcome"`
	Argv             []string `yaml:"argv" json:"argv"`
	CWD              string   `yaml:"cwd" json:"cwd"`
	Exit             int      `yaml:"exit" json:"exit"`
	BoundOIDOrSHA256 string   `yaml:"bound_oid_or_sha256" json:"bound_oid_or_sha256"`
	StdoutArtifactID *string  `yaml:"stdout_artifact_id" json:"stdout_artifact_id"`
	StderrArtifactID *string  `yaml:"stderr_artifact_id" json:"stderr_artifact_id"`
}

type TaskReviewEntry struct {
	ID                  string   `yaml:"id" json:"id"`
	Severity            string   `yaml:"severity" json:"severity"`
	Summary             string   `yaml:"summary" json:"summary"`
	EvidenceArtifactIDs []string `yaml:"evidence_artifact_ids" json:"evidence_artifact_ids"`
}

type TaskReviewRecord struct {
	Findings               []TaskReviewEntry `yaml:"findings" json:"findings"`
	Fixes                  []TaskReviewEntry `yaml:"fixes" json:"fixes"`
	OpenActionableFindings []TaskReviewEntry `yaml:"open_actionable_findings" json:"open_actionable_findings"`
}

type TaskAcceptedDebtRecord struct {
	ID                  string   `yaml:"id" json:"id"`
	RiskClass           string   `yaml:"risk_class" json:"risk_class"`
	Severity            string   `yaml:"severity" json:"severity"`
	Summary             string   `yaml:"summary" json:"summary"`
	Control             string   `yaml:"control" json:"control"`
	EvidenceArtifactIDs []string `yaml:"evidence_artifact_ids" json:"evidence_artifact_ids"`
}

type TaskRecorderRecord struct {
	ActorClaim     string `yaml:"actor_claim" json:"actor_claim"`
	ControlSurface string `yaml:"control_surface" json:"control_surface"`
	RecordedAtUTC  string `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}

type TaskResultRecord struct {
	ID                    TaskResultID               `yaml:"id" json:"id"`
	PublicationKey        string                     `yaml:"publication_key" json:"publication_key"`
	DraftSHA256           string                     `yaml:"draft_sha256" json:"draft_sha256"`
	StoreTransition       string                     `yaml:"store_transition" json:"store_transition"`
	TaskID                TaskID                     `yaml:"task_id" json:"task_id"`
	TaskWorktreeID        WorktreeID                 `yaml:"task_worktree_id" json:"task_worktree_id"`
	ProjectID             ProjectID                  `yaml:"project_id" json:"project_id"`
	RepoID                RepoID                     `yaml:"repo_id" json:"repo_id"`
	GitCommonDir          string                     `yaml:"git_common_dir" json:"git_common_dir"`
	SourceLocator         string                     `yaml:"source_locator" json:"source_locator"`
	SourceRef             string                     `yaml:"source_ref" json:"source_ref"`
	ResultOID             string                     `yaml:"result_oid" json:"result_oid"`
	ResultTree            string                     `yaml:"result_tree" json:"result_tree"`
	ActivityID            string                     `yaml:"activity_id" json:"activity_id"`
	RunID                 string                     `yaml:"run_id" json:"run_id"`
	HandoffID             string                     `yaml:"handoff_id" json:"handoff_id"`
	HandoffLocator        string                     `yaml:"handoff_locator" json:"handoff_locator"`
	HandoffSHA256         string                     `yaml:"handoff_sha256" json:"handoff_sha256"`
	StartReceiptID        string                     `yaml:"start_receipt_id" json:"start_receipt_id"`
	StartReceiptLocator   string                     `yaml:"start_receipt_locator" json:"start_receipt_locator"`
	StartReceiptSHA256    string                     `yaml:"start_receipt_sha256" json:"start_receipt_sha256"`
	TerminalResultID      string                     `yaml:"terminal_result_id" json:"terminal_result_id"`
	TerminalResultLocator string                     `yaml:"terminal_result_locator" json:"terminal_result_locator"`
	TerminalResultSHA256  string                     `yaml:"terminal_result_sha256" json:"terminal_result_sha256"`
	InspectionSHA256      string                     `yaml:"inspection_sha256" json:"inspection_sha256"`
	ReportedOutcome       string                     `yaml:"reported_outcome" json:"reported_outcome"`
	TechnicalGate         string                     `yaml:"technical_gate" json:"technical_gate"`
	VerifierResults       []TaskVerifierResultRecord `yaml:"verifier_results" json:"verifier_results"`
	Review                TaskReviewRecord           `yaml:"review" json:"review"`
	AcceptedDebt          []TaskAcceptedDebtRecord   `yaml:"accepted_debt" json:"accepted_debt"`
	Artifacts             []TaskArtifactRecord       `yaml:"artifacts" json:"artifacts"`
	Recorder              TaskRecorderRecord         `yaml:"recorder" json:"recorder"`
}

type TaskHumanActorRecord struct {
	ActorClaim     string `yaml:"actor_claim" json:"actor_claim"`
	StartSurface   string `yaml:"start_surface" json:"start_surface"`
	StartedAtUTC   string `yaml:"started_at_utc" json:"started_at_utc"`
	CompletedAtUTC string `yaml:"completed_at_utc" json:"completed_at_utc"`
}

type TaskHumanQAEvidenceRecord struct {
	ID        string `yaml:"id" json:"id"`
	Role      string `yaml:"role" json:"role"`
	Locator   string `yaml:"locator" json:"locator"`
	SHA256    string `yaml:"sha256" json:"sha256"`
	SizeBytes int64  `yaml:"size_bytes" json:"size_bytes"`
}

type TaskResidualRiskRecord struct {
	ID       string `yaml:"id" json:"id"`
	Severity string `yaml:"severity" json:"severity"`
	Summary  string `yaml:"summary" json:"summary"`
}

type TaskHumanQARecord struct {
	ID                    HumanQARecordID             `yaml:"id" json:"id"`
	PublicationKey        string                      `yaml:"publication_key" json:"publication_key"`
	DraftSHA256           string                      `yaml:"draft_sha256" json:"draft_sha256"`
	TaskID                TaskID                      `yaml:"task_id" json:"task_id"`
	TaskResultID          TaskResultID                `yaml:"task_result_id" json:"task_result_id"`
	ResultOID             string                      `yaml:"result_oid" json:"result_oid"`
	ResultTree            string                      `yaml:"result_tree" json:"result_tree"`
	Outcome               string                      `yaml:"outcome" json:"outcome"`
	Actor                 TaskHumanActorRecord        `yaml:"actor" json:"actor"`
	Evidence              []TaskHumanQAEvidenceRecord `yaml:"evidence" json:"evidence"`
	Observation           string                      `yaml:"observation" json:"observation"`
	AcceptedResidualRisks []TaskResidualRiskRecord    `yaml:"accepted_residual_risks" json:"accepted_residual_risks"`
}

// TaskHandoffEvidence is the redacted, rehashed projection supplied by the workflow-handoff
// package. It deliberately contains no reply capability or secret material.
type TaskHandoffEvidence struct {
	DeliveryAuthorization                                                      *DeliveryAuthorization
	TaskRequirementsValid                                                      bool
	TaskSpecBasis                                                              *TaskSpecBasis
	TaskSpecValid                                                              bool
	AcceptedStartOutcome                                                       string
	ActivityID, RunID, HandoffID, HandoffLocator, HandoffSHA256                string
	StartReceiptID, StartReceiptLocator, StartReceiptSHA256                    string
	TerminalResultID, TerminalResultLocator, TerminalResultSHA256              string
	InspectionSHA256, ReportedOutcome                                          string
	TargetWorktree, TargetRef, ResultOID, ResultTree, GitCommonDir             string
	SchemaValid, DigestValid, LifecycleValid, BindingValid                     bool
	CapabilityValid, PrincipalSessionValid, PolicyValid, EvidenceCoverageValid bool
	VerifierResults                                                            []TaskVerifierResultRecord
	ExpectedVerifierIDs                                                        []string
	Review                                                                     TaskReviewRecord
	Artifacts                                                                  []TaskArtifactRecord
	EvidenceGaps                                                               []string
}

type TaskHandoffEvidenceRequest struct {
	ActivityID, RunID, HandoffID, HandoffLocator, HandoffSHA256   string
	StartReceiptID, StartReceiptLocator, StartReceiptSHA256       string
	TerminalResultID, TerminalResultLocator, TerminalResultSHA256 string
	InspectionSHA256                                              string
}

type TaskHandoffEvidenceReader interface {
	ReadTaskEvidence(TaskHandoffEvidenceRequest) (TaskHandoffEvidence, error)
}

type TaskLifecycleIDSource interface {
	NewTaskResultID() (TaskResultID, error)
	NewHumanQARecordID() (HumanQARecordID, error)
	NewIntegrationAuthorityID() (IntegrationAuthorityID, error)
	NewIntegrationIntentID() (IntegrationIntentID, error)
	NewIntegrationAttemptID() (IntegrationAttemptID, error)
	NewIntegrationResultID() (IntegrationResultID, error)
}

type WorkClock interface{ Now() time.Time }

type cryptoTaskLifecycleIDSource struct{ Reader io.Reader }

func (source cryptoTaskLifecycleIDSource) reader() io.Reader {
	if source.Reader != nil {
		return source.Reader
	}
	return rand.Reader
}
func (source cryptoTaskLifecycleIDSource) next(prefix string) (string, error) {
	value, err := randomWorkItemID(source.reader(), prefix)
	return value, err
}
func (source cryptoTaskLifecycleIDSource) NewTaskResultID() (TaskResultID, error) {
	v, e := source.next("trs_")
	return TaskResultID(v), e
}
func (source cryptoTaskLifecycleIDSource) NewHumanQARecordID() (HumanQARecordID, error) {
	v, e := source.next("hqa_")
	return HumanQARecordID(v), e
}
func (source cryptoTaskLifecycleIDSource) NewIntegrationAuthorityID() (IntegrationAuthorityID, error) {
	v, e := source.next("iauth_")
	return IntegrationAuthorityID(v), e
}
func (source cryptoTaskLifecycleIDSource) NewIntegrationIntentID() (IntegrationIntentID, error) {
	v, e := source.next("iint_")
	return IntegrationIntentID(v), e
}
func (source cryptoTaskLifecycleIDSource) NewIntegrationAttemptID() (IntegrationAttemptID, error) {
	v, e := source.next("iat_")
	return IntegrationAttemptID(v), e
}
func (source cryptoTaskLifecycleIDSource) NewIntegrationResultID() (IntegrationResultID, error) {
	v, e := source.next("ires_")
	return IntegrationResultID(v), e
}

type systemWorkClock struct{}

func (systemWorkClock) Now() time.Time { return time.Now() }

type TaskResultRecordInput struct {
	TaskID TaskID
	File   string
}
type TaskHumanQARecordInput struct {
	TaskID TaskID
	File   string
}
type TaskResultMutationResult struct {
	RegistryVersion int
	SpecBinding     *TaskResultSpecBinding
	Workspace       string
	Record          TaskResultRecord
	Created         bool
}
type TaskHumanQAMutationResult struct {
	Workspace string
	Record    TaskHumanQARecord
	Created   bool
}

type taskResultDraft struct {
	Canonical                                                       []byte
	Digest, PublicationKey                                          string
	TaskID                                                          TaskID
	TaskWorktreeID                                                  WorktreeID
	Handoff                                                         TaskHandoffEvidenceRequest
	ProjectID                                                       ProjectID
	RepoID                                                          RepoID
	GitCommonDir, WorktreeLocator, SourceRef, ResultOID, ResultTree string
	TechnicalGate                                                   string
	RequiredVerifierIDs                                             []string
	AcceptedDebt                                                    []TaskAcceptedDebtRecord
	EvidenceArtifacts                                               []TaskArtifactRecord
	Recorder                                                        TaskRecorderRecord
}

type taskHumanQADraft struct {
	Canonical                      []byte
	Digest, PublicationKey         string
	TaskID                         TaskID
	TaskResultID                   TaskResultID
	ResultOID, ResultTree, Outcome string
	Actor                          TaskHumanActorRecord
	Evidence                       []TaskHumanQAEvidenceRecord
	Observation                    string
	AcceptedResidualRisks          []TaskResidualRiskRecord
}

var lifecycleIDPatterns = map[string]*regexp.Regexp{
	"task result":           regexp.MustCompile(`^trs_[0-9a-f]{32}$`),
	"human QA":              regexp.MustCompile(`^hqa_[0-9a-f]{32}$`),
	"integration authority": regexp.MustCompile(`^iauth_[0-9a-f]{32}$`),
	"integration intent":    regexp.MustCompile(`^iint_[0-9a-f]{32}$`),
	"integration attempt":   regexp.MustCompile(`^iat_[0-9a-f]{32}$`),
	"integration result":    regexp.MustCompile(`^ires_[0-9a-f]{32}$`),
}

var handoffEvidenceIDPatterns = map[string]*regexp.Regexp{
	"activity": regexp.MustCompile(`^act_[0-9a-f]{32}$`),
	"run":      regexp.MustCompile(`^run_[0-9a-f]{32}$`),
	"handoff":  regexp.MustCompile(`^hnd_[0-9a-f]{32}$`),
	"start":    regexp.MustCompile(`^rcp_[0-9a-f]{32}$`),
	"terminal": regexp.MustCompile(`^res_[0-9a-f]{32}$`),
}

func ParseTaskResultID(value string) (TaskResultID, error) {
	if !lifecycleIDPatterns["task result"].MatchString(value) {
		return "", WorkInvalidArguments("invalid Task result ID")
	}
	return TaskResultID(value), nil
}
func ParseHumanQARecordID(value string) (HumanQARecordID, error) {
	if !lifecycleIDPatterns["human QA"].MatchString(value) {
		return "", WorkInvalidArguments("invalid human QA record ID")
	}
	return HumanQARecordID(value), nil
}

func RecordTaskResult(dependencies Dependencies, input TaskResultRecordInput) (TaskResultMutationResult, error) {
	if _, err := ParseTaskID(string(input.TaskID)); err != nil {
		return TaskResultMutationResult{}, err
	}
	draft, err := readTaskResultDraft(dependencies.Files, input.File)
	if err != nil {
		return TaskResultMutationResult{}, err
	}
	if draft.TaskID != input.TaskID {
		return TaskResultMutationResult{}, workError(ErrorTaskResultSchemaInvalid, "draft task_id does not match command Task", nil)
	}
	if err := requireTaskDeliveryDependencies(dependencies, true); err != nil {
		return TaskResultMutationResult{}, err
	}
	evidence, err := dependencies.HandoffEvidence.ReadTaskEvidence(draft.Handoff)
	if err != nil {
		return TaskResultMutationResult{}, workError(ErrorTaskResultEvidenceConflict, err.Error(), err)
	}
	if err := validateTaskEvidence(draft, evidence); err != nil {
		return TaskResultMutationResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return TaskResultMutationResult{}, err
	}
	result := TaskResultMutationResult{Workspace: root}
	err = dependencies.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			result.RegistryVersion = registry.FormatVersion
			task, _ := findTask(registry, input.TaskID)
			if task == nil {
				return workError(ErrorWorkNotFound, fmt.Sprintf("Task %s is not registered", input.TaskID), nil)
			}
			if task.WorktreeState != WorkItemReady || task.Worktree == nil {
				return workError(ErrorTaskResultConflict, "Task does not have a ready worktree", nil)
			}
			for _, existing := range registry.TaskResults {
				if existing.PublicationKey == draft.PublicationKey {
					if existing.DraftSHA256 == draft.Digest {
						result.Record = existing
						result.SpecBinding = taskResultSpecLink(registry, existing.ID)
						return validateTaskContentClosure(dependencies.TaskContent, root, registry, false)
					}
					return workError(ErrorTaskResultConflict, "Task result publication_key already exists with different content", nil)
				}
				if existing.TaskID == draft.TaskID && existing.TerminalResultID == draft.Handoff.TerminalResultID && existing.ResultOID == draft.ResultOID && existing.DraftSHA256 != draft.Digest {
					return workError(ErrorTaskResultConflict, "Task result evidence is already recorded with different content", nil)
				}
			}
			_, repo, err := projectAndRepo(projects, task.ProjectID, task.RepoID)
			if err != nil {
				return err
			}
			if err := validateTaskResultBindings(*task, repo, draft, evidence, dependencies); err != nil {
				return err
			}
			if taskRequiresSpec(registry, task.ID) {
				if evidence.TaskSpecBasis == nil || !evidence.TaskSpecValid || evidence.AcceptedStartOutcome != "started" || evidence.TaskSpecBasis.TaskID != task.ID {
					return contentError("task_spec_result_basis_missing", "a controlled started Task Spec basis is required", nil)
				}
				if _, err := historicalTaskSpec(dependencies, root, registry, *evidence.TaskSpecBasis); err != nil {
					return err
				}
			}
			id, err := dependencies.TaskLifecycleIDs.NewTaskResultID()
			if err != nil {
				return err
			}
			transition := "none"
			if registry.FormatVersion == 1 {
				transition = "format_1_to_2"
				upgradeRegistryToV2(&registry)
			}
			record := TaskResultRecord{ID: id, PublicationKey: draft.PublicationKey, DraftSHA256: draft.Digest, StoreTransition: transition, TaskID: draft.TaskID, TaskWorktreeID: draft.TaskWorktreeID, ProjectID: draft.ProjectID, RepoID: draft.RepoID, GitCommonDir: draft.GitCommonDir, SourceLocator: draft.WorktreeLocator, SourceRef: draft.SourceRef, ResultOID: draft.ResultOID, ResultTree: draft.ResultTree, ActivityID: evidence.ActivityID, RunID: evidence.RunID, HandoffID: evidence.HandoffID, HandoffLocator: evidence.HandoffLocator, HandoffSHA256: evidence.HandoffSHA256, StartReceiptID: evidence.StartReceiptID, StartReceiptLocator: evidence.StartReceiptLocator, StartReceiptSHA256: evidence.StartReceiptSHA256, TerminalResultID: evidence.TerminalResultID, TerminalResultLocator: evidence.TerminalResultLocator, TerminalResultSHA256: evidence.TerminalResultSHA256, InspectionSHA256: evidence.InspectionSHA256, ReportedOutcome: evidence.ReportedOutcome, TechnicalGate: draft.TechnicalGate, VerifierResults: evidence.VerifierResults, Review: evidence.Review, AcceptedDebt: draft.AcceptedDebt, Artifacts: draft.EvidenceArtifacts, Recorder: draft.Recorder}
			registry.TaskResults = append(registry.TaskResults, record)
			if evidence.TaskSpecBasis != nil {
				link := TaskResultSpecBinding{TaskResultID: id, TaskID: task.ID, Basis: *evidence.TaskSpecBasis, HandoffSHA256: evidence.HandoffSHA256, StartReceiptSHA256: evidence.StartReceiptSHA256}
				registry.TaskResultSpecBindings = append(registry.TaskResultSpecBindings, link)
				result.SpecBinding = &link
			}
			sortWorkRegistry(&registry)
			if _, err := publishWorkItemRegistryRecover(session, registry); err != nil {
				return err
			}
			result.Record, result.Created = record, true
			return nil
		})
	})
	if err != nil {
		return TaskResultMutationResult{}, mapMutationError(err)
	}
	return result, nil
}

func RecordTaskHumanQA(dependencies Dependencies, input TaskHumanQARecordInput) (TaskHumanQAMutationResult, error) {
	if _, err := ParseTaskID(string(input.TaskID)); err != nil {
		return TaskHumanQAMutationResult{}, err
	}
	draft, err := readTaskHumanQADraft(dependencies.Files, input.File)
	if err != nil {
		return TaskHumanQAMutationResult{}, err
	}
	if draft.TaskID != input.TaskID {
		return TaskHumanQAMutationResult{}, workError(ErrorTaskQASchemaInvalid, "draft task_id does not match command Task", nil)
	}
	if err := requireTaskDeliveryDependencies(dependencies, false); err != nil {
		return TaskHumanQAMutationResult{}, err
	}
	if err := rehashQAEvidence(dependencies.Files, draft.Evidence); err != nil {
		return TaskHumanQAMutationResult{}, err
	}
	root, err := containingWorkItemWorkspace(dependencies)
	if err != nil {
		return TaskHumanQAMutationResult{}, err
	}
	result := TaskHumanQAMutationResult{Workspace: root}
	err = dependencies.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return dependencies.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			registry, err := session.Snapshot()
			if err != nil {
				return err
			}
			if registry.FormatVersion < 2 || registry.FormatVersion > 4 {
				return workError(ErrorTaskQAConflict, "Task result store has not been upgraded to format 2 or 3", nil)
			}
			task, _ := findTask(registry, draft.TaskID)
			if task == nil {
				return workError(ErrorWorkNotFound, "Task is not registered", nil)
			}
			if _, _, err := projectAndRepo(projects, task.ProjectID, task.RepoID); err != nil {
				return err
			}
			selected := findTaskResult(registry, draft.TaskResultID)
			if selected == nil || selected.TaskID != draft.TaskID || selected.ResultOID != draft.ResultOID || selected.ResultTree != draft.ResultTree {
				return workError(ErrorTaskQAEvidenceConflict, "human QA does not bind the selected Task result/OID/tree", nil)
			}
			for _, existing := range registry.HumanQARecords {
				if existing.PublicationKey == draft.PublicationKey {
					if existing.DraftSHA256 == draft.Digest {
						result.Record = existing
						return nil
					}
					return workError(ErrorTaskQAConflict, "human QA publication_key already exists with different content", nil)
				}
			}
			id, err := dependencies.TaskLifecycleIDs.NewHumanQARecordID()
			if err != nil {
				return err
			}
			record := TaskHumanQARecord{ID: id, PublicationKey: draft.PublicationKey, DraftSHA256: draft.Digest, TaskID: draft.TaskID, TaskResultID: draft.TaskResultID, ResultOID: draft.ResultOID, ResultTree: draft.ResultTree, Outcome: draft.Outcome, Actor: draft.Actor, Evidence: draft.Evidence, Observation: draft.Observation, AcceptedResidualRisks: draft.AcceptedResidualRisks}
			registry.HumanQARecords = append(registry.HumanQARecords, record)
			sortWorkRegistry(&registry)
			if _, err := publishWorkItemRegistryRecover(session, registry); err != nil {
				return err
			}
			result.Record, result.Created = record, true
			return nil
		})
	})
	if err != nil {
		return TaskHumanQAMutationResult{}, mapMutationError(err)
	}
	return result, nil
}

func requireTaskDeliveryDependencies(dependencies Dependencies, evidence bool) error {
	if dependencies.Files == nil || dependencies.ProjectLocks == nil || dependencies.WorkItems == nil || dependencies.TaskLifecycleIDs == nil {
		return workError(ErrorWorkIO, "Task lifecycle dependencies are required", nil)
	}
	if evidence && (dependencies.HandoffEvidence == nil || dependencies.WorkGit == nil) {
		return workError(ErrorWorkIO, "Task result evidence and Git dependencies are required", nil)
	}
	return nil
}

func validateTaskResultBindings(task TaskRecord, repo RepoRecord, draft taskResultDraft, evidence TaskHandoffEvidence, dependencies Dependencies) error {
	wt := task.Worktree
	if wt.ID != draft.TaskWorktreeID || task.ProjectID != draft.ProjectID || task.RepoID != draft.RepoID || task.GitCommonDir != draft.GitCommonDir || repo.GitCommonDir != draft.GitCommonDir || wt.Locator != draft.WorktreeLocator || wt.Ref != draft.SourceRef {
		return workError(ErrorTaskResultConflict, "Task result source binding differs from the registered Task", nil)
	}
	observed, err := dependencies.WorkGit.ObserveWorktree(wt.Locator)
	if err != nil {
		return workError(ErrorTaskResultEvidenceConflict, err.Error(), err)
	}
	if observed.Ref != draft.SourceRef || observed.OID != draft.ResultOID || observed.Tree != draft.ResultTree || observed.GitCommonDir != draft.GitCommonDir || !observed.Clean || !observed.InventoryMatch {
		return workError(ErrorTaskResultEvidenceConflict, "Task result does not match a clean exact source worktree", nil)
	}
	if evidence.TargetWorktree != draft.WorktreeLocator || evidence.TargetRef != draft.SourceRef || evidence.ResultOID != draft.ResultOID || evidence.ResultTree != draft.ResultTree || evidence.GitCommonDir != draft.GitCommonDir {
		return workError(ErrorTaskResultEvidenceConflict, "terminal target differs from the Task result source", nil)
	}
	return nil
}

func validateTaskEvidence(draft taskResultDraft, evidence TaskHandoffEvidence) error {
	request := draft.Handoff
	if evidence.ActivityID != request.ActivityID || evidence.RunID != request.RunID || evidence.HandoffID != request.HandoffID || evidence.HandoffLocator != request.HandoffLocator || evidence.HandoffSHA256 != request.HandoffSHA256 || evidence.StartReceiptID != request.StartReceiptID || evidence.StartReceiptLocator != request.StartReceiptLocator || evidence.StartReceiptSHA256 != request.StartReceiptSHA256 || evidence.TerminalResultID != request.TerminalResultID || evidence.TerminalResultLocator != request.TerminalResultLocator || evidence.TerminalResultSHA256 != request.TerminalResultSHA256 || evidence.InspectionSHA256 != request.InspectionSHA256 {
		return workError(ErrorTaskResultEvidenceConflict, "immutable handoff evidence binding differs from the draft", nil)
	}
	if !evidence.SchemaValid || !evidence.DigestValid || !evidence.LifecycleValid || !evidence.BindingValid || !evidence.CapabilityValid || !evidence.PrincipalSessionValid {
		return workError(ErrorTaskResultEvidenceConflict, "handoff evidence integrity or binding is invalid", nil)
	}
	if len(draft.RequiredVerifierIDs) != len(evidence.ExpectedVerifierIDs) {
		return workError(ErrorTaskResultEvidenceConflict, "required verifier IDs differ from the Handoff", nil)
	}
	for i := range draft.RequiredVerifierIDs {
		if draft.RequiredVerifierIDs[i] != evidence.ExpectedVerifierIDs[i] {
			return workError(ErrorTaskResultEvidenceConflict, "required verifier IDs differ from the Handoff", nil)
		}
	}
	green := draft.TechnicalGate == "passed" || draft.TechnicalGate == "good_enough_with_known_debt"
	if green && evidence.TaskSpecBasis != nil && !evidence.TaskRequirementsValid {
		return workError(ErrorTaskResultEvidenceConflict, "a green selected-Spec result requires complete passing functional requirement coverage", nil)
	}
	if green && (!evidence.PolicyValid || evidence.ReportedOutcome != "complete") {
		return workError(ErrorTaskResultEvidenceConflict, "green gate requires a policy-valid complete result", nil)
	}
	if green && len(evidence.Review.OpenActionableFindings) != 0 {
		return workError(ErrorTaskResultEvidenceConflict, "handoff review has open actionable findings", nil)
	}
	seen := map[string]int{}
	for _, verifier := range evidence.VerifierResults {
		if verifier.Outcome == "passed" {
			seen[verifier.VerifierID]++
		}
	}
	if green {
		for _, required := range draft.RequiredVerifierIDs {
			if seen[required] != 1 {
				return workError(ErrorTaskResultEvidenceConflict, "required verifier is missing, duplicated, or failed: "+required, nil)
			}
		}
		if len(evidence.VerifierResults) != len(draft.RequiredVerifierIDs) {
			return workError(ErrorTaskResultEvidenceConflict, "terminal verifier set differs from the Handoff", nil)
		}
	}
	if draft.TechnicalGate == "passed" && (!evidence.EvidenceCoverageValid || len(draft.AcceptedDebt) != 0) {
		return workError(ErrorTaskResultEvidenceConflict, "passed gate requires complete evidence and no accepted debt", nil)
	}
	if draft.TechnicalGate == "good_enough_with_known_debt" && (evidence.EvidenceCoverageValid || len(draft.AcceptedDebt) == 0) {
		return workError(ErrorTaskResultEvidenceConflict, "good-enough gate requires named evidence debt", nil)
	}
	if draft.TechnicalGate == "good_enough_with_known_debt" {
		if len(draft.AcceptedDebt) != len(evidence.EvidenceGaps) {
			return workError(ErrorTaskResultEvidenceConflict, "accepted debt must cover evidence gaps one-to-one", nil)
		}
		artifacts := map[string]string{}
		for _, artifact := range draft.EvidenceArtifacts {
			artifacts[artifact.ArtifactID] = artifact.Role
		}
		for _, debt := range draft.AcceptedDebt {
			for _, id := range debt.EvidenceArtifactIDs {
				if artifacts[id] != "debt_control" {
					return workError(ErrorTaskResultEvidenceConflict, "accepted debt requires managed debt_control evidence", nil)
				}
			}
		}
		for _, entries := range [][]TaskReviewEntry{evidence.Review.Findings, evidence.Review.Fixes, evidence.Review.OpenActionableFindings} {
			for _, finding := range entries {
				if finding.Severity == "high" || finding.Severity == "critical" {
					return workError(ErrorTaskResultEvidenceConflict, "high or critical review risk cannot be accepted as debt", nil)
				}
			}
		}
	}
	if len(draft.EvidenceArtifacts) != len(evidence.Artifacts) {
		return workError(ErrorTaskResultEvidenceConflict, "draft artifact set differs from terminal evidence", nil)
	}
	for i, item := range draft.EvidenceArtifacts {
		observed := evidence.Artifacts[i]
		roleMatches := item.Role == observed.Role || (item.Role == "debt_control" && observed.Role == "other")
		if item.ArtifactID != observed.ArtifactID || !roleMatches || item.Locator != observed.Locator || item.SHA256 != observed.SHA256 || item.SizeBytes != observed.SizeBytes {
			return workError(ErrorTaskResultEvidenceConflict, "draft artifact metadata differs from terminal evidence", nil)
		}
	}
	return nil
}

func upgradeRegistryToV2(registry *WorkItemRegistry) {
	registry.FormatVersion = 2
	if registry.TaskResults == nil {
		registry.TaskResults = []TaskResultRecord{}
	}
	if registry.HumanQARecords == nil {
		registry.HumanQARecords = []TaskHumanQARecord{}
	}
	if registry.IntegrationAuthorities == nil {
		registry.IntegrationAuthorities = []IntegrationAuthority{}
	}
	if registry.IntegrationIntents == nil {
		registry.IntegrationIntents = []IntegrationIntent{}
	}
	if registry.IntegrationAttempts == nil {
		registry.IntegrationAttempts = []IntegrationAttempt{}
	}
	if registry.IntegrationResults == nil {
		registry.IntegrationResults = []IntegrationResult{}
	}
}

func findTaskResult(registry WorkItemRegistry, id TaskResultID) *TaskResultRecord {
	for i := range registry.TaskResults {
		if registry.TaskResults[i].ID == id {
			return &registry.TaskResults[i]
		}
	}
	return nil
}
func findHumanQA(registry WorkItemRegistry, id HumanQARecordID) *TaskHumanQARecord {
	for i := range registry.HumanQARecords {
		if registry.HumanQARecords[i].ID == id {
			return &registry.HumanQARecords[i]
		}
	}
	return nil
}

func readTaskResultDraft(files FileSystem, input string) (taskResultDraft, error) {
	bytes, value, err := readStrictTaskDraft(files, input, 256<<10, "Task result draft", resultSchemaError)
	if err != nil {
		return taskResultDraft{}, err
	}
	fields, err := exactTaskObject(value, "Task result draft", "kind", "schema_version", "format", "format_version", "canonicalization", "publication_key", "task_id", "task_worktree_id", "handoff", "source", "technical_assessment", "evidence_artifacts", "recorder")
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	if err := validateTaskEnvelope(fields, "WorkspaceTaskResultRecordDraft@1"); err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	d := taskResultDraft{Canonical: bytes, Digest: digestTaskBytes(bytes)}
	d.PublicationKey, _ = taskString(fields, "publication_key")
	if err := validatePublicationKey(d.PublicationKey); err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	taskID, _ := taskString(fields, "task_id")
	d.TaskID, err = ParseTaskID(taskID)
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	wt, _ := taskString(fields, "task_worktree_id")
	if !worktreeIDPattern.MatchString(wt) {
		return taskResultDraft{}, resultSchemaError(errors.New("invalid task_worktree_id"))
	}
	d.TaskWorktreeID = WorktreeID(wt)
	h, err := exactTaskObject(fields["handoff"], "handoff", "activity_id", "run_id", "handoff_id", "handoff_locator", "handoff_sha256", "start_receipt_id", "start_receipt_locator", "start_receipt_sha256", "terminal_result_id", "terminal_result_locator", "terminal_result_sha256", "inspection_sha256")
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	d.Handoff = TaskHandoffEvidenceRequest{ActivityID: mustTaskString(h, "activity_id"), RunID: mustTaskString(h, "run_id"), HandoffID: mustTaskString(h, "handoff_id"), HandoffLocator: mustTaskString(h, "handoff_locator"), HandoffSHA256: mustTaskString(h, "handoff_sha256"), StartReceiptID: mustTaskString(h, "start_receipt_id"), StartReceiptLocator: mustTaskString(h, "start_receipt_locator"), StartReceiptSHA256: mustTaskString(h, "start_receipt_sha256"), TerminalResultID: mustTaskString(h, "terminal_result_id"), TerminalResultLocator: mustTaskString(h, "terminal_result_locator"), TerminalResultSHA256: mustTaskString(h, "terminal_result_sha256"), InspectionSHA256: mustTaskString(h, "inspection_sha256")}
	if !handoffEvidenceIDPatterns["activity"].MatchString(d.Handoff.ActivityID) || !handoffEvidenceIDPatterns["run"].MatchString(d.Handoff.RunID) || !handoffEvidenceIDPatterns["handoff"].MatchString(d.Handoff.HandoffID) || !handoffEvidenceIDPatterns["start"].MatchString(d.Handoff.StartReceiptID) || !handoffEvidenceIDPatterns["terminal"].MatchString(d.Handoff.TerminalResultID) || !validAbsoluteCleanPath(d.Handoff.HandoffLocator) || !validAbsoluteCleanPath(d.Handoff.StartReceiptLocator) || !validAbsoluteCleanPath(d.Handoff.TerminalResultLocator) || !digestPattern.MatchString(d.Handoff.HandoffSHA256) || !digestPattern.MatchString(d.Handoff.StartReceiptSHA256) || !digestPattern.MatchString(d.Handoff.TerminalResultSHA256) || !digestPattern.MatchString(d.Handoff.InspectionSHA256) {
		return taskResultDraft{}, resultSchemaError(errors.New("invalid handoff evidence binding"))
	}
	s, err := exactTaskObject(fields["source"], "source", "project_id", "repo_id", "git_common_dir", "worktree_locator", "source_ref", "result_oid", "result_tree")
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	d.ProjectID = ProjectID(mustTaskString(s, "project_id"))
	d.RepoID = RepoID(mustTaskString(s, "repo_id"))
	d.GitCommonDir = mustTaskString(s, "git_common_dir")
	d.WorktreeLocator = mustTaskString(s, "worktree_locator")
	d.SourceRef = mustTaskString(s, "source_ref")
	d.ResultOID = mustTaskString(s, "result_oid")
	d.ResultTree = mustTaskString(s, "result_tree")
	if validateWorkIdentifier("Project", string(d.ProjectID)) != nil || validateWorkIdentifier("repository", string(d.RepoID)) != nil || !validStoredPath(d.GitCommonDir) || !validStoredPath(d.WorktreeLocator) || !validFullBranchRef(d.SourceRef) || !validOIDText(d.ResultOID) || !validOIDText(d.ResultTree) || len(d.ResultOID) != len(d.ResultTree) {
		return taskResultDraft{}, resultSchemaError(errors.New("invalid source binding"))
	}
	ta, err := exactTaskObject(fields["technical_assessment"], "technical_assessment", "gate", "required_verifier_ids", "accepted_debt")
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	d.TechnicalGate = mustTaskString(ta, "gate")
	if !setString("passed", "good_enough_with_known_debt", "failed", "unknown")[d.TechnicalGate] {
		return taskResultDraft{}, resultSchemaError(errors.New("invalid technical gate"))
	}
	d.RequiredVerifierIDs, err = taskStringArray(ta["required_verifier_ids"], true)
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	d.AcceptedDebt, err = parseAcceptedDebt(ta["accepted_debt"])
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	d.EvidenceArtifacts, err = parseTaskArtifacts(fields["evidence_artifacts"])
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	r, err := exactTaskObject(fields["recorder"], "recorder", "actor_claim", "control_surface", "recorded_at_utc")
	if err != nil {
		return taskResultDraft{}, resultSchemaError(err)
	}
	d.Recorder = TaskRecorderRecord{ActorClaim: mustTaskString(r, "actor_claim"), ControlSurface: mustTaskString(r, "control_surface"), RecordedAtUTC: mustTaskString(r, "recorded_at_utc")}
	if !validTaskText(d.Recorder.ActorClaim, 1, 256) || !validTaskText(d.Recorder.ControlSurface, 1, 256) || !validTaskUTC(d.Recorder.RecordedAtUTC) {
		return taskResultDraft{}, resultSchemaError(errors.New("invalid recorder"))
	}
	return d, nil
}

func readTaskHumanQADraft(files FileSystem, input string) (taskHumanQADraft, error) {
	bytes, value, err := readStrictTaskDraft(files, input, 256<<10, "human QA draft", qaSchemaError)
	if err != nil {
		return taskHumanQADraft{}, err
	}
	f, err := exactTaskObject(value, "human QA draft", "kind", "schema_version", "format", "format_version", "canonicalization", "publication_key", "task_id", "task_result_id", "result_oid", "result_tree", "outcome", "actor", "evidence", "observation", "accepted_residual_risks")
	if err != nil {
		return taskHumanQADraft{}, qaSchemaError(err)
	}
	if err := validateTaskEnvelope(f, "WorkspaceTaskHumanQARecordDraft@1"); err != nil {
		return taskHumanQADraft{}, qaSchemaError(err)
	}
	d := taskHumanQADraft{Canonical: bytes, Digest: digestTaskBytes(bytes), PublicationKey: mustTaskString(f, "publication_key"), ResultOID: mustTaskString(f, "result_oid"), ResultTree: mustTaskString(f, "result_tree"), Outcome: mustTaskString(f, "outcome"), Observation: mustTaskString(f, "observation")}
	if err := validatePublicationKey(d.PublicationKey); err != nil {
		return taskHumanQADraft{}, qaSchemaError(err)
	}
	d.TaskID, err = ParseTaskID(mustTaskString(f, "task_id"))
	if err != nil {
		return taskHumanQADraft{}, qaSchemaError(err)
	}
	d.TaskResultID, err = ParseTaskResultID(mustTaskString(f, "task_result_id"))
	if err != nil {
		return taskHumanQADraft{}, qaSchemaError(err)
	}
	if !validOIDText(d.ResultOID) || !validOIDText(d.ResultTree) || !setString("pass", "fail", "blocked")[d.Outcome] || !validTaskText(d.Observation, 1, 2000) {
		return taskHumanQADraft{}, qaSchemaError(errors.New("invalid human QA outcome, result, or observation"))
	}
	a, err := exactTaskObject(f["actor"], "actor", "actor_claim", "start_surface", "started_at_utc", "completed_at_utc")
	if err != nil {
		return taskHumanQADraft{}, qaSchemaError(err)
	}
	d.Actor = TaskHumanActorRecord{ActorClaim: mustTaskString(a, "actor_claim"), StartSurface: mustTaskString(a, "start_surface"), StartedAtUTC: mustTaskString(a, "started_at_utc"), CompletedAtUTC: mustTaskString(a, "completed_at_utc")}
	started, se := time.Parse(time.RFC3339Nano, d.Actor.StartedAtUTC)
	completed, ce := time.Parse(time.RFC3339Nano, d.Actor.CompletedAtUTC)
	if !validTaskText(d.Actor.ActorClaim, 1, 256) || !validTaskText(d.Actor.StartSurface, 1, 256) || se != nil || ce != nil || !strings.HasSuffix(d.Actor.StartedAtUTC, "Z") || !strings.HasSuffix(d.Actor.CompletedAtUTC, "Z") || completed.Before(started) {
		return taskHumanQADraft{}, qaSchemaError(errors.New("invalid human QA actor"))
	}
	d.Evidence, err = parseQAEvidence(f["evidence"])
	if err != nil {
		return taskHumanQADraft{}, qaSchemaError(err)
	}
	d.AcceptedResidualRisks, err = parseResidualRisks(f["accepted_residual_risks"])
	if err != nil {
		return taskHumanQADraft{}, qaSchemaError(err)
	}
	if d.Outcome == "pass" {
		found := false
		for _, e := range d.Evidence {
			if e.Role == "report" {
				found = true
			}
		}
		if !found {
			return taskHumanQADraft{}, qaSchemaError(errors.New("pass requires report evidence"))
		}
	}
	return d, nil
}

func readStrictTaskDraft(files FileSystem, input string, limit int, label string, schemaError func(error) error) ([]byte, canonicaljson.Value, error) {
	if files == nil {
		return nil, nil, workError(ErrorWorkIO, "filesystem dependency is required", nil)
	}
	cwd, err := files.Getwd()
	if err != nil {
		return nil, nil, workError(ErrorWorkIO, "get working directory", err)
	}
	path := input
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	path = filepath.Clean(path)
	info, err := files.Lstat(path)
	if err != nil {
		return nil, nil, workError(ErrorWorkIO, "inspect "+label, err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return nil, nil, schemaError(errors.New(label + " must be a regular non-symlink file"))
	}
	physical, err := files.EvalSymlinks(path)
	if err != nil || physical != path {
		return nil, nil, schemaError(errors.New(label + " path must be physical"))
	}
	if info.Size() > int64(limit) {
		return nil, nil, schemaError(errors.New(label + " exceeds 256 KiB"))
	}
	raw, err := files.ReadFile(path)
	if err != nil {
		return nil, nil, workError(ErrorWorkIO, "read "+label, err)
	}
	value, err := canonicaljson.DecodeStrict(raw)
	if err != nil {
		return nil, nil, schemaError(err)
	}
	canonical, err := canonicaljson.Marshal(value)
	if err != nil {
		return nil, nil, schemaError(err)
	}
	return canonical, value, nil
}

func exactTaskObject(value canonicaljson.Value, context string, names ...string) (map[string]canonicaljson.Value, error) {
	o, ok := value.(canonicaljson.Object)
	if !ok {
		return nil, fmt.Errorf("%s must be an object", context)
	}
	want := map[string]bool{}
	for _, n := range names {
		want[n] = true
	}
	got := map[string]canonicaljson.Value{}
	for _, m := range o {
		if !want[m.Name] {
			return nil, fmt.Errorf("%s has unknown field %s", context, m.Name)
		}
		if _, exists := got[m.Name]; exists {
			return nil, fmt.Errorf("%s has duplicate field %s", context, m.Name)
		}
		got[m.Name] = m.Value
	}
	for _, n := range names {
		if _, ok := got[n]; !ok {
			return nil, fmt.Errorf("%s is missing field %s", context, n)
		}
	}
	return got, nil
}
func taskString(f map[string]canonicaljson.Value, n string) (string, error) {
	v, ok := f[n].(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", n)
	}
	return v, nil
}
func mustTaskString(f map[string]canonicaljson.Value, n string) string {
	v, _ := taskString(f, n)
	return v
}
func validateTaskEnvelope(f map[string]canonicaljson.Value, kind string) error {
	if mustTaskString(f, "kind") != kind || mustTaskString(f, "format") != "json" || mustTaskString(f, "canonicalization") != "RFC8785" {
		return errors.New("invalid draft envelope")
	}
	for _, n := range []string{"schema_version", "format_version"} {
		if v, ok := f[n].(int64); !ok || v != 1 {
			return errors.New("invalid draft version")
		}
	}
	return nil
}
func digestTaskBytes(b []byte) string {
	d := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(d[:])
}
func validatePublicationKey(v string) error {
	if len(v) < 1 || len(v) > 128 || !regexp.MustCompile(`^[a-z0-9][a-z0-9._/-]*$`).MatchString(v) || strings.Contains(v, "//") || strings.HasSuffix(v, "/") {
		return errors.New("invalid publication_key")
	}
	for _, p := range strings.Split(v, "/") {
		if p == ".." {
			return errors.New("invalid publication_key")
		}
	}
	return nil
}
func validAbsoluteCleanPath(v string) bool {
	return filepath.IsAbs(v) && filepath.Clean(v) == v && !strings.ContainsRune(v, 0)
}
func validTaskText(v string, min, max int) bool {
	if !utf8.ValidString(v) || v != strings.TrimSpace(v) {
		return false
	}
	count := 0
	for _, r := range v {
		if unicode.IsControl(r) {
			return false
		}
		count++
	}
	return count >= min && count <= max
}
func validTaskUTC(v string) bool {
	_, err := time.Parse(time.RFC3339Nano, v)
	return err == nil && strings.HasSuffix(v, "Z")
}
func setString(v ...string) map[string]bool {
	m := map[string]bool{}
	for _, s := range v {
		m[s] = true
	}
	return m
}
func taskStringArray(v canonicaljson.Value, allowEmpty bool) ([]string, error) {
	a, ok := v.([]canonicaljson.Value)
	if !ok {
		return nil, errors.New("expected array")
	}
	if !allowEmpty && len(a) == 0 {
		return nil, errors.New("array must not be empty")
	}
	r := make([]string, len(a))
	for i, x := range a {
		s, ok := x.(string)
		if !ok || !validTaskText(s, 1, 256) {
			return nil, errors.New("array contains invalid string")
		}
		r[i] = s
		if i > 0 && r[i-1] >= s {
			return nil, errors.New("array must be sorted and unique")
		}
	}
	return r, nil
}
func parseTaskArtifacts(v canonicaljson.Value) ([]TaskArtifactRecord, error) {
	a, ok := v.([]canonicaljson.Value)
	if !ok {
		return nil, errors.New("evidence_artifacts must be an array")
	}
	r := make([]TaskArtifactRecord, 0, len(a))
	last := ""
	var total int64
	for _, x := range a {
		f, e := exactTaskObject(x, "artifact", "artifact_id", "role", "locator", "sha256", "size_bytes")
		if e != nil {
			return nil, e
		}
		z := TaskArtifactRecord{ArtifactID: mustTaskString(f, "artifact_id"), Role: mustTaskString(f, "role"), Locator: mustTaskString(f, "locator"), SHA256: mustTaskString(f, "sha256")}
		var ok bool
		z.SizeBytes, ok = f["size_bytes"].(int64)
		total += z.SizeBytes
		if !ok || validatePublicationKey(z.ArtifactID) != nil || z.ArtifactID <= last || !setString("verifier_stdout", "verifier_stderr", "review", "debt_control", "other")[z.Role] || !validAbsoluteCleanPath(z.Locator) || !digestPattern.MatchString(z.SHA256) || z.SizeBytes < 0 || z.SizeBytes > 64<<20 || total > 256<<20 {
			return nil, errors.New("invalid or unsorted artifact")
		}
		last = z.ArtifactID
		r = append(r, z)
	}
	return r, nil
}
func parseAcceptedDebt(v canonicaljson.Value) ([]TaskAcceptedDebtRecord, error) {
	a, ok := v.([]canonicaljson.Value)
	if !ok {
		return nil, errors.New("accepted_debt must be an array")
	}
	r := make([]TaskAcceptedDebtRecord, 0, len(a))
	last := ""
	for _, x := range a {
		f, e := exactTaskObject(x, "accepted debt", "id", "risk_class", "severity", "summary", "control", "evidence_artifact_ids")
		if e != nil {
			return nil, e
		}
		z := TaskAcceptedDebtRecord{ID: mustTaskString(f, "id"), RiskClass: mustTaskString(f, "risk_class"), Severity: mustTaskString(f, "severity"), Summary: mustTaskString(f, "summary"), Control: mustTaskString(f, "control")}
		z.EvidenceArtifactIDs, e = taskStringArray(f["evidence_artifact_ids"], false)
		if e != nil || validatePublicationKey(z.ID) != nil || z.ID <= last || !setString("coverage", "evidence", "journal")[z.RiskClass] || !setString("low", "medium")[z.Severity] || !validTaskText(z.Summary, 1, 600) || !validTaskText(z.Control, 1, 2000) {
			return nil, errors.New("invalid or unsorted accepted debt")
		}
		last = z.ID
		r = append(r, z)
	}
	return r, nil
}
func parseQAEvidence(v canonicaljson.Value) ([]TaskHumanQAEvidenceRecord, error) {
	a, ok := v.([]canonicaljson.Value)
	if !ok || len(a) == 0 {
		return nil, errors.New("evidence must be a non-empty array")
	}
	r := make([]TaskHumanQAEvidenceRecord, 0, len(a))
	last := ""
	var total int64
	for _, x := range a {
		f, e := exactTaskObject(x, "QA evidence", "id", "role", "locator", "sha256", "size_bytes")
		if e != nil {
			return nil, e
		}
		z := TaskHumanQAEvidenceRecord{ID: mustTaskString(f, "id"), Role: mustTaskString(f, "role"), Locator: mustTaskString(f, "locator"), SHA256: mustTaskString(f, "sha256")}
		var ok bool
		z.SizeBytes, ok = f["size_bytes"].(int64)
		total += z.SizeBytes
		if !ok || validatePublicationKey(z.ID) != nil || z.ID <= last || !setString("script", "report", "screenshot", "log", "other")[z.Role] || !validAbsoluteCleanPath(z.Locator) || !digestPattern.MatchString(z.SHA256) || z.SizeBytes < 0 || z.SizeBytes > 64<<20 || total > 256<<20 {
			return nil, errors.New("invalid or unsorted QA evidence")
		}
		last = z.ID
		r = append(r, z)
	}
	return r, nil
}
func parseResidualRisks(v canonicaljson.Value) ([]TaskResidualRiskRecord, error) {
	a, ok := v.([]canonicaljson.Value)
	if !ok {
		return nil, errors.New("accepted_residual_risks must be an array")
	}
	r := make([]TaskResidualRiskRecord, 0, len(a))
	last := ""
	for _, x := range a {
		f, e := exactTaskObject(x, "residual risk", "id", "severity", "summary")
		if e != nil {
			return nil, e
		}
		z := TaskResidualRiskRecord{ID: mustTaskString(f, "id"), Severity: mustTaskString(f, "severity"), Summary: mustTaskString(f, "summary")}
		if validatePublicationKey(z.ID) != nil || z.ID <= last || !setString("low", "medium")[z.Severity] || !validTaskText(z.Summary, 1, 600) {
			return nil, errors.New("invalid or unsorted residual risk")
		}
		last = z.ID
		r = append(r, z)
	}
	return r, nil
}
func rehashQAEvidence(files FileSystem, evidence []TaskHumanQAEvidenceRecord) error {
	var total int64
	for _, item := range evidence {
		if !filepath.IsAbs(item.Locator) || filepath.Clean(item.Locator) != item.Locator {
			return workError(ErrorTaskQAEvidenceConflict, "QA evidence locator must be absolute and clean", nil)
		}
		info, err := files.Lstat(item.Locator)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
			return workError(ErrorTaskQAEvidenceConflict, "QA evidence is not a regular non-symlink file", err)
		}
		physical, err := files.EvalSymlinks(item.Locator)
		if err != nil || physical != item.Locator {
			return workError(ErrorTaskQAEvidenceConflict, "QA evidence locator is not physical", err)
		}
		b, err := files.ReadFile(item.Locator)
		if err != nil {
			return workError(ErrorTaskQAEvidenceConflict, "cannot read QA evidence", err)
		}
		total += int64(len(b))
		if int64(len(b)) != item.SizeBytes || digestTaskBytes(b) != item.SHA256 || len(b) > 64<<20 || total > 256<<20 {
			return workError(ErrorTaskQAEvidenceConflict, "QA evidence metadata differs from bytes", nil)
		}
	}
	return nil
}
func resultSchemaError(err error) error {
	return workError(ErrorTaskResultSchemaInvalid, err.Error(), err)
}
func qaSchemaError(err error) error { return workError(ErrorTaskQASchemaInvalid, err.Error(), err) }

func MarshalTaskResultReadback(result TaskResultMutationResult) ([]byte, error) {
	if result.RegistryVersion == 3 {
		status := "not_recorded"
		if result.SpecBinding != nil {
			status = "bound"
		}
		return canonicaljson.Marshal(contentObject(map[string]canonicaljson.Value{"kind": "WorkspaceTaskResultRecordReadback@2", "schema_version": int64(2), "format": "json", "format_version": int64(1), "canonicalization": "RFC8785", "workspace": result.Workspace, "record": taskResultCanonical(result.Record), "created": result.Created, "next_action": taskResultNextAction(result.Record), "spec_binding": contentRefValue(result.SpecBinding), "basis_status": status, "next_transition_authorized": false}))
	}
	return canonicaljson.Marshal(canonicaljson.Object{{Name: "kind", Value: "WorkspaceTaskResultRecordReadback@1"}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"}, {Name: "workspace", Value: result.Workspace}, {Name: "record", Value: taskResultCanonical(result.Record)}, {Name: "created", Value: result.Created}, {Name: "next_action", Value: taskResultNextAction(result.Record)}})
}
func MarshalTaskHumanQAReadback(result TaskHumanQAMutationResult) ([]byte, error) {
	return canonicaljson.Marshal(canonicaljson.Object{{Name: "kind", Value: "WorkspaceTaskHumanQARecordReadback@1"}, {Name: "schema_version", Value: int64(1)}, {Name: "format", Value: "json"}, {Name: "format_version", Value: int64(1)}, {Name: "canonicalization", Value: "RFC8785"}, {Name: "workspace", Value: result.Workspace}, {Name: "record", Value: taskHumanQACanonical(result.Record)}, {Name: "created", Value: result.Created}, {Name: "next_action", Value: taskQANextAction(result.Record)}})
}
func taskResultNextAction(r TaskResultRecord) string {
	if r.TechnicalGate == "passed" || r.TechnicalGate == "good_enough_with_known_debt" {
		return fmt.Sprintf("Record human QA with `ply workspace task qa record %s --file <absolute-human-qa.json>`.", r.TaskID)
	}
	return "Start a separate result control before recording human QA."
}
func taskQANextAction(r TaskHumanQARecord) string {
	if r.Outcome == "pass" {
		return fmt.Sprintf("Run `ply workspace task integrate %s --result %s --qa %s --expected-result-oid %s --expected-parent-oid <oid> --check`.", r.TaskID, r.TaskResultID, r.ID, r.ResultOID)
	}
	return "Start a separate delivery or QA clarification before integration."
}
func taskResultCanonical(r TaskResultRecord) canonicaljson.Value {
	b, _ := canonicaljson.DecodeStrict(mustJSON(r))
	return b
}
func taskHumanQACanonical(r TaskHumanQARecord) canonicaljson.Value {
	b, _ := canonicaljson.DecodeStrict(mustJSON(r))
	return b
}
func mustJSON(v any) []byte { return marshalCanonicalStruct(v) }

// marshalCanonicalStruct is intentionally limited to the record types in this file. YAML
// remains the store wire; this bridge produces their canonical JSON readback without maps.
func marshalCanonicalStruct(v any) []byte {
	value := recordCanonicalValue(v)
	b, _ := canonicaljson.Marshal(value)
	return b
}
func recordCanonicalValue(v any) canonicaljson.Value {
	switch x := v.(type) {
	case TaskResultRecord:
		verifiers := []canonicaljson.Value{}
		for _, z := range x.VerifierResults {
			verifiers = append(verifiers, recordCanonicalValue(z))
		}
		debts := []canonicaljson.Value{}
		for _, z := range x.AcceptedDebt {
			debts = append(debts, recordCanonicalValue(z))
		}
		artifacts := []canonicaljson.Value{}
		for _, z := range x.Artifacts {
			artifacts = append(artifacts, recordCanonicalValue(z))
		}
		return canonicaljson.Object{{Name: "id", Value: string(x.ID)}, {Name: "publication_key", Value: x.PublicationKey}, {Name: "draft_sha256", Value: x.DraftSHA256}, {Name: "store_transition", Value: x.StoreTransition}, {Name: "task_id", Value: string(x.TaskID)}, {Name: "task_worktree_id", Value: string(x.TaskWorktreeID)}, {Name: "project_id", Value: string(x.ProjectID)}, {Name: "repo_id", Value: string(x.RepoID)}, {Name: "git_common_dir", Value: x.GitCommonDir}, {Name: "source_locator", Value: x.SourceLocator}, {Name: "source_ref", Value: x.SourceRef}, {Name: "result_oid", Value: x.ResultOID}, {Name: "result_tree", Value: x.ResultTree}, {Name: "activity_id", Value: x.ActivityID}, {Name: "run_id", Value: x.RunID}, {Name: "handoff_id", Value: x.HandoffID}, {Name: "handoff_locator", Value: x.HandoffLocator}, {Name: "handoff_sha256", Value: x.HandoffSHA256}, {Name: "start_receipt_id", Value: x.StartReceiptID}, {Name: "start_receipt_locator", Value: x.StartReceiptLocator}, {Name: "start_receipt_sha256", Value: x.StartReceiptSHA256}, {Name: "terminal_result_id", Value: x.TerminalResultID}, {Name: "terminal_result_locator", Value: x.TerminalResultLocator}, {Name: "terminal_result_sha256", Value: x.TerminalResultSHA256}, {Name: "inspection_sha256", Value: x.InspectionSHA256}, {Name: "reported_outcome", Value: x.ReportedOutcome}, {Name: "technical_gate", Value: x.TechnicalGate}, {Name: "verifier_results", Value: verifiers}, {Name: "review", Value: recordCanonicalValue(x.Review)}, {Name: "accepted_debt", Value: debts}, {Name: "artifacts", Value: artifacts}, {Name: "recorder", Value: canonicaljson.Object{{Name: "actor_claim", Value: x.Recorder.ActorClaim}, {Name: "control_surface", Value: x.Recorder.ControlSurface}, {Name: "recorded_at_utc", Value: x.Recorder.RecordedAtUTC}}}}
	case TaskHumanQARecord:
		e := []canonicaljson.Value{}
		for _, z := range x.Evidence {
			e = append(e, recordCanonicalValue(z))
		}
		rr := []canonicaljson.Value{}
		for _, z := range x.AcceptedResidualRisks {
			rr = append(rr, recordCanonicalValue(z))
		}
		return canonicaljson.Object{{Name: "id", Value: string(x.ID)}, {Name: "publication_key", Value: x.PublicationKey}, {Name: "draft_sha256", Value: x.DraftSHA256}, {Name: "task_id", Value: string(x.TaskID)}, {Name: "task_result_id", Value: string(x.TaskResultID)}, {Name: "result_oid", Value: x.ResultOID}, {Name: "result_tree", Value: x.ResultTree}, {Name: "outcome", Value: x.Outcome}, {Name: "actor", Value: canonicaljson.Object{{Name: "actor_claim", Value: x.Actor.ActorClaim}, {Name: "start_surface", Value: x.Actor.StartSurface}, {Name: "started_at_utc", Value: x.Actor.StartedAtUTC}, {Name: "completed_at_utc", Value: x.Actor.CompletedAtUTC}}}, {Name: "evidence", Value: e}, {Name: "observation", Value: x.Observation}, {Name: "accepted_residual_risks", Value: rr}}
	case TaskVerifierResultRecord:
		return canonicaljson.Object{{Name: "verifier_id", Value: x.VerifierID}, {Name: "outcome", Value: x.Outcome}, {Name: "argv", Value: stringValues(x.Argv)}, {Name: "cwd", Value: x.CWD}, {Name: "exit", Value: int64(x.Exit)}, {Name: "bound_oid_or_sha256", Value: x.BoundOIDOrSHA256}, {Name: "stdout_artifact_id", Value: pointerValue(x.StdoutArtifactID)}, {Name: "stderr_artifact_id", Value: pointerValue(x.StderrArtifactID)}}
	case TaskReviewRecord:
		f := []canonicaljson.Value{}
		for _, z := range x.Findings {
			f = append(f, recordCanonicalValue(z))
		}
		fx := []canonicaljson.Value{}
		for _, z := range x.Fixes {
			fx = append(fx, recordCanonicalValue(z))
		}
		o := []canonicaljson.Value{}
		for _, z := range x.OpenActionableFindings {
			o = append(o, recordCanonicalValue(z))
		}
		return canonicaljson.Object{{Name: "findings", Value: f}, {Name: "fixes", Value: fx}, {Name: "open_actionable_findings", Value: o}}
	case TaskReviewEntry:
		return canonicaljson.Object{{Name: "id", Value: x.ID}, {Name: "severity", Value: x.Severity}, {Name: "summary", Value: x.Summary}, {Name: "evidence_artifact_ids", Value: stringValues(x.EvidenceArtifactIDs)}}
	case TaskAcceptedDebtRecord:
		return canonicaljson.Object{{Name: "id", Value: x.ID}, {Name: "risk_class", Value: x.RiskClass}, {Name: "severity", Value: x.Severity}, {Name: "summary", Value: x.Summary}, {Name: "control", Value: x.Control}, {Name: "evidence_artifact_ids", Value: stringValues(x.EvidenceArtifactIDs)}}
	case TaskArtifactRecord:
		return canonicaljson.Object{{Name: "artifact_id", Value: x.ArtifactID}, {Name: "role", Value: x.Role}, {Name: "locator", Value: x.Locator}, {Name: "sha256", Value: x.SHA256}, {Name: "size_bytes", Value: x.SizeBytes}}
	case TaskHumanQAEvidenceRecord:
		return canonicaljson.Object{{Name: "id", Value: x.ID}, {Name: "role", Value: x.Role}, {Name: "locator", Value: x.Locator}, {Name: "sha256", Value: x.SHA256}, {Name: "size_bytes", Value: x.SizeBytes}}
	case TaskResidualRiskRecord:
		return canonicaljson.Object{{Name: "id", Value: x.ID}, {Name: "severity", Value: x.Severity}, {Name: "summary", Value: x.Summary}}
	default:
		return nil
	}
}
