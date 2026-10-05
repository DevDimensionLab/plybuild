package taskrun

// The workflow transport shares the native validators and reservation lock, but
// has no TaskRun process record and never publishes a TaskResult.
import (
	"encoding/json"
	"path/filepath"
	"strings"
)

type WorkflowRequest struct {
	Envelope
	RequestKey        string             `json:"request_key"`
	WorkspaceRoot     string             `json:"workspace_root"`
	PreparationID     string             `json:"preparation_id"`
	PreparationSHA256 string             `json:"preparation_sha256"`
	HandoffDraft      json.RawMessage    `json:"handoff_draft"`
	Runtime           Runtime            `json:"runtime"`
	CodexProjectTrust *CodexProjectTrust `json:"codex_project_trust,omitempty"`
	Agreement         Agreement          `json:"agreement"`
	HumanAuthority    HumanAuthority     `json:"human_authority"`
	Herdr             struct {
		Executable  Executable `json:"executable"`
		WorkspaceID string     `json:"workspace_id"`
		TabLabel    string     `json:"tab_label"`
	} `json:"herdr"`
	Coordinator struct {
		ActorClaim        string `json:"actor_claim"`
		MayRequestChanges bool   `json:"may_request_changes"`
	} `json:"coordinator"`
	ReturnMode string            `json:"return_mode"`
	Delivery   *DeliveryContract `json:"delivery,omitempty"`
}
type WorkflowPaths struct {
	RunRoot string `json:"run_root"`
	Context string `json:"context"`
}
type WorkflowHandoff struct {
	ID      string `json:"id"`
	Locator string `json:"locator"`
	SHA256  string `json:"sha256"`
}
type WorkflowTransport struct {
	WorkspaceID    string `json:"workspace_id"`
	TabID          string `json:"tab_id"`
	PaneID         string `json:"pane_id"`
	TerminalID     string `json:"terminal_id"`
	AgentSessionID string `json:"agent_session_id"`
	State          string `json:"state"`
	Observation    string `json:"observation"`
	ObservedAt     string `json:"observed_at"`
	LaunchPending  *bool  `json:"launch_pending,omitempty"`
}
type WorkflowRound struct {
	Number               int     `json:"number"`
	State                string  `json:"state"`
	ReportPath           *string `json:"report_path"`
	ReportSHA256         *string `json:"report_sha256"`
	ControlID            *string `json:"control_id"`
	PreviousReportSHA256 *string `json:"previous_report_sha256"`
}
type WorkflowReviewState struct {
	ReviewID *string `json:"review_id"`
	Decision *string `json:"decision"`
}
type WorkflowBudget struct {
	Used              *BudgetUsage `json:"used"`
	Floor             BudgetUsage  `json:"floor"`
	WithinAgreement   *bool        `json:"within_agreement"`
	MeasurementSource string       `json:"measurement_source"`
}
type WorkflowFinal struct {
	State          string  `json:"state"`
	ReportSHA256   *string `json:"report_sha256"`
	TerminalSHA256 *string `json:"terminal_sha256"`
}
type WorkflowAction struct {
	Actor   string `json:"actor"`
	Message string `json:"message"`
}
type WorkflowRun struct {
	Envelope
	Provider        string              `json:"provider,omitempty"`
	RunID           string              `json:"run_id"`
	RequestSHA256   string              `json:"request_sha256"`
	SessionID       string              `json:"session_id"`
	Handoff         WorkflowHandoff     `json:"handoff"`
	Paths           WorkflowPaths       `json:"paths"`
	Transport       WorkflowTransport   `json:"transport"`
	Round           WorkflowRound       `json:"round"`
	Review          WorkflowReviewState `json:"review"`
	Budget          WorkflowBudget      `json:"budget"`
	FinalReturn     WorkflowFinal       `json:"final_return"`
	TaskResultState string              `json:"task_result_state"`
	RuntimeFacts    RuntimeFacts        `json:"runtime_facts"`
	Reasons         []Reason            `json:"reasons"`
	NextAction      WorkflowAction      `json:"next_action"`
	Delivery        *DeliveryState      `json:"delivery,omitempty"`
}
type WorkflowPreview struct {
	Envelope
	Provider      string              `json:"provider,omitempty"`
	RunID         string              `json:"run_id"`
	RequestSHA256 string              `json:"request_sha256"`
	Confirmation  *string             `json:"confirmation"`
	Reasons       []Reason            `json:"reasons"`
	Effects       []string            `json:"effects"`
	Observed      Observed            `json:"observed"`
	Paths         WorkflowPaths       `json:"paths"`
	CodexTrust    *workflowTrustFacts `json:"codex_project_trust,omitempty"`
}
type WorkflowReport struct {
	Report
	Round                int     `json:"round"`
	ControlID            *string `json:"control_id"`
	PreviousReportSHA256 *string `json:"previous_report_sha256"`
}
type WorkflowFinding struct {
	Observed   string `json:"observed"`
	Expected   string `json:"expected"`
	Acceptance string `json:"acceptance"`
}
type WorkflowReview struct {
	Envelope
	ReviewID      string            `json:"review_id"`
	RunID         string            `json:"run_id"`
	RequestSHA256 string            `json:"request_sha256"`
	HandoffSHA256 string            `json:"handoff_sha256"`
	Round         int               `json:"round"`
	ReportSHA256  string            `json:"report_sha256"`
	Reviewer      string            `json:"reviewer"`
	Decision      string            `json:"decision"`
	Findings      []WorkflowFinding `json:"findings"`
}
type workflowContext struct {
	Envelope
	RunID                string     `json:"run_id"`
	RequestSHA256        string     `json:"request_sha256"`
	SessionID            string     `json:"session_id"`
	HandoffSHA256        string     `json:"handoff_sha256"`
	Round                int        `json:"round"`
	ControlID            *string    `json:"control_id"`
	PreviousReportSHA256 *string    `json:"previous_report_sha256"`
	PlyExecutable        Executable `json:"ply_executable"`
}
type workflowRecord struct {
	Number       int          `json:"number"`
	Report       FileBinding  `json:"report"`
	Terminal     FileBinding  `json:"terminal"`
	TargetSHA256 string       `json:"target_sha256"`
	Review       *FileBinding `json:"review"`
}
type workflowState struct {
	Request       WorkflowRequest     `json:"request"`
	Observed      Observed            `json:"observed"`
	Result        WorkflowRun         `json:"result"`
	Phase         string              `json:"phase"`
	ContextSHA256 string              `json:"context_sha256"`
	Acceptance    *FileBinding        `json:"acceptance"`
	StartDraft    *FileBinding        `json:"start_draft"`
	StartSHA256   *string             `json:"start_sha256"`
	Records       []workflowRecord    `json:"records"`
	CodexTrust    *workflowTrustFacts `json:"codex_project_trust,omitempty"`
}

func workflowEnv(kind string) Envelope { return Envelope{"ply.workflow." + kind, 1} }
func workflowID(r WorkflowRequest) string {
	return "wfr_" + strings.TrimPrefix(digest([]string{r.WorkspaceRoot, r.RequestKey}), "sha256:")
}
func workflowRoot(root string) string { return filepath.Join(root, ".ply", "workflow-runs") }
func workflowPaths(r WorkflowRequest, round int) WorkflowPaths {
	p := filepath.Join(workflowRoot(r.WorkspaceRoot), "runs", workflowID(r))
	return WorkflowPaths{p, filepath.Join(p, "rounds", workflowRoundName(round), "context.json")}
}
func workflowNativeRequest(r WorkflowRequest) Request {
	return Request{Envelope: env("request"), RequestKey: r.RequestKey, WorkspaceRoot: r.WorkspaceRoot, PreparationID: r.PreparationID, PreparationSHA256: r.PreparationSHA256, HandoffDraft: r.HandoffDraft, Runtime: r.Runtime, Agreement: r.Agreement, HumanAuthority: HumanAuthority{r.HumanAuthority.ActorClaim, "human_ordinary_terminal", r.HumanAuthority.Authorized}, ReturnPolicy: ReturnPolicy{true, true, "separate", "separate"}}
}
func workflowError(exit int, detail string) error {
	return &Error{Code: "workflow_run", Detail: detail, Exit: exit}
}
