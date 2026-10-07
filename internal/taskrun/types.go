// Package taskrun owns one bounded provider attempt and its durable return.
// Workflow receipts and TaskResult validation remain owned by their existing packages.
package taskrun

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"time"
)

type Envelope struct {
	Kind          string `json:"kind"`
	SchemaVersion int    `json:"schema_version"`
}
type Executable struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type Evidence struct {
	Locator string `json:"locator"`
	SHA256  string `json:"sha256"`
	Role    string `json:"role"`
}
type Permission struct {
	AuthorityKind         string     `json:"authority_kind"`
	ProfileID             string     `json:"profile_id"`
	EffectivePolicySHA256 string     `json:"effective_policy_sha256"`
	Evidence              []Evidence `json:"evidence"`
}
type Runtime struct {
	Provider          string     `json:"provider"`
	Mode              string     `json:"mode"`
	Model             string     `json:"model"`
	Executable        Executable `json:"executable"`
	ConfigProfile     *string    `json:"config_profile"`
	PermissionBinding Permission `json:"permission_binding"`
	PlyExecutable     Executable `json:"ply_executable"`
}
type Agreement struct {
	Name                   string `json:"name"`
	MaxCorrectionRounds    int    `json:"max_correction_rounds"`
	MaxActiveSeconds       int    `json:"max_active_seconds"`
	MaxEnvironmentMeasures int    `json:"max_environment_measures"`
}
type ReturnPolicy struct {
	RecordTaskResult         bool   `json:"record_task_result"`
	RequireProcessQuiescence bool   `json:"require_process_quiescence"`
	HumanQA                  string `json:"human_qa"`
	Integration              string `json:"integration"`
}
type HumanAuthority struct {
	ActorClaim   string `json:"actor_claim"`
	StartSurface string `json:"start_surface"`
	Authorized   bool   `json:"authorized"`
}
type FactoryTest struct {
	AuthorizationPath   string `json:"authorization_path"`
	AuthorizationSHA256 string `json:"authorization_sha256"`
	Iteration           int    `json:"iteration"`
	TaskSlot            int    `json:"task_slot"`
	ReasoningEffort     string `json:"reasoning_effort"`
	TimeoutSeconds      int    `json:"timeout_seconds"`
}
type Request struct {
	Envelope
	RequestKey        string          `json:"request_key"`
	WorkspaceRoot     string          `json:"workspace_root"`
	PreparationID     string          `json:"preparation_id"`
	PreparationSHA256 string          `json:"preparation_sha256"`
	HandoffDraft      json.RawMessage `json:"handoff_draft"`
	Runtime           Runtime         `json:"runtime"`
	Agreement         Agreement       `json:"agreement"`
	ReturnPolicy      ReturnPolicy    `json:"return_policy"`
	HumanAuthority    HumanAuthority  `json:"human_authority"`
	FactoryTest       *FactoryTest    `json:"factory_test,omitempty"`
}
type Reason struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}
type FileBinding struct {
	Locator string `json:"locator"`
	SHA256  string `json:"sha256"`
}
type Observed struct {
	WorkspaceMarkerSHA256 string                            `json:"workspace_marker_sha256"`
	RegistrySHA256        string                            `json:"registry_sha256"`
	Target                workspace.PlanWorktreeObservation `json:"target"`
	Epic                  workspace.PlanWorktreeObservation `json:"epic"`
	RuntimeBindings       []FileBinding                     `json:"runtime_bindings"`
}
type Paths struct {
	RunRoot    string `json:"run_root"`
	TempRoot   string `json:"temp_root"`
	ReportPath string `json:"report_path"`
}
type Preview struct {
	Envelope
	RequestSHA256 string                    `json:"request_sha256"`
	RunID         string                    `json:"run_id"`
	Request       Request                   `json:"request"`
	Preparation   workspace.TaskPreparation `json:"preparation"`
	Observed      Observed                  `json:"observed"`
	Paths         Paths                     `json:"paths"`
	Reasons       []Reason                  `json:"reasons"`
	Confirmation  *string                   `json:"confirmation"`
	NextArgv      []string                  `json:"next_argv"`
}
type Handoff struct {
	ActivityID string `json:"activity_id"`
	RunID      string `json:"run_id"`
	HandoffID  string `json:"handoff_id"`
	Locator    string `json:"locator"`
	SHA256     string `json:"sha256"`
}
type Binding struct {
	Envelope
	RunID         string                    `json:"run_id"`
	RequestSHA256 string                    `json:"request_sha256"`
	Preparation   workspace.TaskPreparation `json:"preparation"`
	Handoff       Handoff                   `json:"handoff"`
	SessionID     string                    `json:"session_id"`
	RunRoot       string                    `json:"run_root"`
	TempRoot      string                    `json:"temp_root"`
	ReportPath    string                    `json:"report_path"`
}
type Context struct {
	Envelope
	RunID         string     `json:"run_id"`
	RequestSHA256 string     `json:"request_sha256"`
	SessionID     string     `json:"session_id"`
	WorkspaceRoot string     `json:"workspace_root"`
	BindingSHA256 string     `json:"binding_sha256"`
	PlyExecutable Executable `json:"ply_executable"`
}
type Event struct {
	Envelope
	RunID          string          `json:"run_id"`
	RequestSHA256  string          `json:"request_sha256"`
	Sequence       int             `json:"sequence"`
	PreviousSHA256 *string         `json:"previous_sha256"`
	Type           string          `json:"type"`
	RecordedAtUTC  string          `json:"recorded_at_utc"`
	Payload        json.RawMessage `json:"payload"`
}
type Launch struct {
	State    string `json:"state"`
	Attempts int    `json:"attempts"`
}
type AcceptanceState struct {
	State         string  `json:"state"`
	ReceiptSHA256 *string `json:"receipt_sha256"`
}
type Process struct {
	State         string  `json:"state"`
	PID           *int    `json:"pid"`
	StartIdentity *string `json:"start_identity"`
	ProcessGroup  *int    `json:"process_group"`
	ExitCode      *int    `json:"exit_code"`
	Signal        *string `json:"signal"`
	Quiescence    *bool   `json:"quiescence"`
}
type Delivery struct {
	State           string  `json:"state"`
	ReportSHA256    *string `json:"report_sha256"`
	TerminalSHA256  *string `json:"terminal_sha256"`
	ReportedOutcome *string `json:"reported_outcome"`
}
type Collection struct {
	State                 string  `json:"state"`
	TaskResultID          *string `json:"task_result_id"`
	TaskResultDraftSHA256 *string `json:"task_result_draft_sha256"`
}
type BudgetUsage struct {
	InitialExecutionStarted bool `json:"initial_execution_started"`
	CorrectionRounds        int  `json:"correction_rounds"`
	ActiveSeconds           int  `json:"active_seconds"`
	EnvironmentMeasures     int  `json:"environment_measures"`
}
type Budget struct {
	Reported        *BudgetUsage `json:"reported"`
	ElapsedSeconds  *int64       `json:"elapsed_seconds"`
	WithinAgreement *bool        `json:"within_agreement"`
}
type RuntimeFacts struct {
	RequestedModel string  `json:"requested_model"`
	ReportedModel  *string `json:"reported_model"`
	ModelState     string  `json:"model_state"`
	Source         *string `json:"source"`
	ClaimSHA256    *string `json:"claim_sha256"`
}
type TaskExecution struct {
	State                   string  `json:"state"`
	Source                  string  `json:"source"`
	StatusSHA256            *string `json:"status_sha256"`
	ObservedAtUTC           *string `json:"observed_at_utc"`
	HistoricalQualification bool    `json:"historical_qualification"`
}
type TaskStatus struct {
	Envelope
	RunID            string  `json:"run_id"`
	RequestSHA256    string  `json:"request_sha256"`
	SessionID        string  `json:"session_id"`
	BasisEventSHA256 string  `json:"basis_event_sha256"`
	Source           string  `json:"source"`
	ActorClaim       string  `json:"actor_claim"`
	ObservedAtUTC    string  `json:"observed_at_utc"`
	State            string  `json:"state"`
	WorktreePath     string  `json:"worktree_path"`
	VisibleGroupPath *string `json:"visible_group_path"`
	TaskLabel        *string `json:"task_label"`
	MatchingTasks    *int    `json:"matching_tasks"`
	NativeSessionID  *string `json:"native_session_id"`
}
type Result struct {
	Envelope
	RunID              string                             `json:"run_id"`
	RequestSHA256      string                             `json:"request_sha256"`
	LastEventSHA256    *string                            `json:"last_event_sha256"`
	Handoff            *Handoff                           `json:"handoff"`
	Launch             Launch                             `json:"launch"`
	Acceptance         AcceptanceState                    `json:"acceptance"`
	Process            Process                            `json:"process"`
	Delivery           Delivery                           `json:"delivery"`
	Collection         Collection                         `json:"collection"`
	Budget             Budget                             `json:"budget"`
	ObservedTarget     *workspace.PlanWorktreeObservation `json:"observed_target"`
	Reasons            []Reason                           `json:"reasons"`
	RuntimeFacts       RuntimeFacts                       `json:"runtime_facts"`
	TaskExecution      TaskExecution                      `json:"task_execution"`
	NextAction         string                             `json:"next_action"`
	ProviderCompletion *FileBinding                       `json:"provider_completion,omitempty"`
}
type RuntimeClaim struct {
	RuntimeID             *string `json:"runtime_id"`
	ModelID               *string `json:"model_id"`
	ProfileID             *string `json:"profile_id"`
	EffectivePolicySHA256 *string `json:"effective_policy_sha256"`
	NativeSessionID       *string `json:"native_session_id"`
}
type Acceptance struct {
	Envelope
	RunID         string          `json:"run_id"`
	RequestSHA256 string          `json:"request_sha256"`
	SessionID     string          `json:"session_id"`
	RuntimeClaim  RuntimeClaim    `json:"runtime_claim"`
	Sandbox       json.RawMessage `json:"sandbox"`
	Acceptance    string          `json:"acceptance"`
	Issues        json.RawMessage `json:"issues"`
}
type TechnicalAssessment struct {
	Gate                string                             `json:"gate"`
	RequiredVerifierIDs []string                           `json:"required_verifier_ids"`
	AcceptedDebt        []workspace.TaskAcceptedDebtRecord `json:"accepted_debt"`
}
type Report struct {
	Envelope
	RunID                    string              `json:"run_id"`
	RequestSHA256            string              `json:"request_sha256"`
	SessionID                string              `json:"session_id"`
	Outcome                  string              `json:"outcome"`
	Summary                  string              `json:"summary"`
	Meaning                  string              `json:"meaning"`
	BudgetUsage              *BudgetUsage        `json:"budget_usage"`
	StopReasons              json.RawMessage     `json:"stop_reasons"`
	ObservedEffects          json.RawMessage     `json:"observed_effects"`
	VerifierResults          json.RawMessage     `json:"verifier_results"`
	Review                   json.RawMessage     `json:"review"`
	Artifacts                json.RawMessage     `json:"artifacts"`
	EvidenceGaps             json.RawMessage     `json:"evidence_gaps"`
	ForbiddenEffectsObserved json.RawMessage     `json:"forbidden_effects_observed"`
	TechnicalAssessment      TechnicalAssessment `json:"technical_assessment"`
	ProcessObservation       *string             `json:"process_observation"`
	PreventiveFollowup       *string             `json:"preventive_followup"`
}
type CollectPreview struct {
	Envelope
	ProviderCompletion       *FileBinding                       `json:"provider_completion,omitempty"`
	RuntimeFacts             *RuntimeFacts                      `json:"runtime_facts,omitempty"`
	TaskExecution            TaskExecution                      `json:"task_execution"`
	BasisEventSHA256         *string                            `json:"basis_event_sha256"`
	ProposedTaskStatusSHA256 *string                            `json:"proposed_task_status_sha256"`
	RunID                    string                             `json:"run_id"`
	RequestSHA256            string                             `json:"request_sha256"`
	Process                  Process                            `json:"process"`
	ReportSHA256             *string                            `json:"report_sha256"`
	TerminalSHA256           *string                            `json:"terminal_sha256"`
	Target                   *workspace.PlanWorktreeObservation `json:"target"`
	InspectionSHA256         *string                            `json:"inspection_sha256"`
	Reasons                  []Reason                           `json:"reasons"`
	Confirmation             *string                            `json:"confirmation"`
	NextArgv                 []string                           `json:"next_argv"`
}

type Error struct {
	Code, Detail string
	Exit         int
}

func (e *Error) Error() string                           { return e.Code + ": " + e.Detail }
func failure(code string, exit int, detail string) error { return &Error{code, detail, exit} }
func invalid(detail string) error                        { return failure("task_run_schema_invalid", 2, detail) }
func conflict(detail string) error                       { return failure("task_run_conflict", 3, detail) }
func integrity(detail string) error                      { return failure("task_run_integrity", 3, detail) }
func ioError(e error) error                              { return failure("task_run_io", 4, fmt.Sprint(e)) }
func env(kind string) Envelope {
	version := 1
	if kind == "result" || kind == "collect-preview" {
		version = 2
	}
	return Envelope{"ply.workspace.task-run-" + kind, version}
}
func ptr[T any](v T) *T { return &v }

type LaunchSpec struct {
	Executable       Executable
	Argv             []string
	CWD, ContextPath string
	Timeout          time.Duration
	StreamsRoot      string
	Completion       *ProviderCompletion
}

// ProcessRunner must call started only after a successful spawn and release all
// run-store locks before blocking. There is no restart or timeout-kill operation.
type ProcessRunner interface {
	Check() error
	Run(LaunchSpec, func(Process) error) (Process, error)
}

// StartupProcess contains only the local identities needed to rule out a
// provider still running under the preserved terminal shell.
type StartupProcess struct {
	PID            int `json:"pid"`
	ParentPID      int `json:"parent_pid"`
	ProcessGroupID int `json:"process_group_id"`
}

type Dependencies struct {
	Executable      func() (string, error)
	Files           FileSystem
	Workspace       workspace.Dependencies
	Workflow        workflowhandoff.Dependencies
	Runner          ProcessRunner
	ExecRunner      ProcessRunner
	Now             func() time.Time
	CallbackContext *string
	TaskStatusPath  string
	ContextPath     func() string
	CWD             func() (string, error)
	Fault           func(string) error
	// StartupProcessTree observes the shell and all of its descendants. A nil
	// observer uses a bounded local process-table read, without command lines.
	StartupProcessTree func(shellPID int) ([]StartupProcess, error)
	// HerdrTimeout bounds transport calls and the shared start/readiness window,
	// never the provider's task lifetime.
	// Zero selects the documented adapter timeouts.
	HerdrTimeout time.Duration
}

func (d Dependencies) fault(point string) error {
	if d.Fault != nil {
		return d.Fault(point)
	}
	return nil
}

// FileSystem injects the durable publication boundary. Production always uses
// synced, exclusive publication; tests can model I/O failures without OS grants.
type FileSystem interface {
	WriteOnce(string, []byte) error
	Replace(string, []byte) error
}
type systemFileSystem struct{}

func (systemFileSystem) WriteOnce(p string, b []byte) error { return writeOnce(p, b) }
func (systemFileSystem) Replace(p string, b []byte) error   { return publishFile(p, b, true) }
func (d Dependencies) files() FileSystem {
	if d.Files != nil {
		return d.Files
	}
	return systemFileSystem{}
}
func (d Dependencies) writeOnce(p string, b []byte) error { return d.files().WriteOnce(p, b) }
func (d Dependencies) writeValue(p string, v any) error {
	b, e := Canonical(v)
	if e != nil {
		return e
	}
	return d.writeOnce(p, b)
}
func (d Dependencies) replaceValue(p string, v any) error {
	b, e := Canonical(v)
	if e != nil {
		return e
	}
	return d.files().Replace(p, b)
}
