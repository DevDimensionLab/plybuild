package taskrun

import (
	"encoding/json"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// DeliveryContract is an explicit v2 ownership contract. It never upgrades an
// existing v1 run, its Agreement A, or the permissions of a provider process.
type DeliveryContract struct {
	OwnerClaim          string                       `json:"owner_claim"`
	Goal                FileBinding                  `json:"goal"`
	AcceptancePath      string                       `json:"acceptance_path"`
	AllowSubagents      bool                         `json:"allow_subagents"`
	AllowLocalInstall   bool                         `json:"allow_local_install"`
	LocalIntegration    string                       `json:"local_integration"`
	NotificationContext *FileBinding                 `json:"notification_context"`
	ReasoningEffort     string                       `json:"reasoning_effort,omitempty"`
	Agreement           *workspace.DeliveryAgreement `json:"agreement,omitempty"`
}

type DeliveryRequestInput struct {
	RequestKey, WorkspaceRoot, PreparationID, PreparationSHA256 string
	HandoffDraft                                                json.RawMessage
	Runtime                                                     Runtime
	HumanAuthority                                              HumanAuthority
	HerdrExecutable                                             Executable
	HerdrWorkspaceID, TabLabel                                  string
	CodexProjectTrust                                           *CodexProjectTrust
	ClaudeProjectTrust                                          *ClaudeProjectTrust
	Delivery                                                    DeliveryContract
}

type DeliveryPermissionAcceptance struct {
	LaunchContractSHA256    string     `json:"launch_contract_sha256"`
	PermissionConfirmed     bool       `json:"permission_confirmed"`
	ActualPolicyEvidence    []Evidence `json:"actual_policy_evidence"`
	DeliveryAgreementSHA256 string     `json:"delivery_agreement_sha256,omitempty"`
	AllowedEffects          []string   `json:"allowed_effects,omitempty"`
}

type DeliveryAcceptance struct {
	Acceptance
	DeliveryPermission DeliveryPermissionAcceptance `json:"delivery_permission"`
}

// DeliveryReport preserves progress and incomplete evidence. A report cannot
// itself claim a qualified candidate, a human verdict, or completed integration.
type DeliveryReport struct {
	Envelope
	RunID               string                   `json:"run_id"`
	RequestSHA256       string                   `json:"request_sha256"`
	SessionID           string                   `json:"session_id"`
	EventID             string                   `json:"event_id"`
	PreviousEventSHA256 *string                  `json:"previous_event_sha256"`
	Phase               string                   `json:"phase"`
	Summary             string                   `json:"summary"`
	Meaning             string                   `json:"meaning"`
	Question            *string                  `json:"question"`
	Evidence            []FileBinding            `json:"evidence"`
	VerifierResults     []DeliveryVerifierResult `json:"verifier_results"`
	ProcessObservation  *string                  `json:"process_observation"`
	PreventiveFollowup  *string                  `json:"preventive_followup"`
}

type DeliveryVerifierResult struct {
	ID       string        `json:"id"`
	Outcome  string        `json:"outcome"`
	Argv     []string      `json:"argv"`
	CWD      string        `json:"cwd"`
	Exit     *int          `json:"exit"`
	Evidence []FileBinding `json:"evidence"`
	Reason   string        `json:"reason"`
}

type DeliveryEventRecord struct {
	ID      string      `json:"id"`
	Kind    string      `json:"kind"`
	Binding FileBinding `json:"binding"`
}

type DeliveryCandidate struct {
	Number       int                          `json:"number"`
	Key          string                       `json:"key"`
	OID          string                       `json:"oid"`
	Tree         string                       `json:"tree"`
	Verification FileBinding                  `json:"verification"`
	Handoff      WorkflowHandoff              `json:"handoff"`
	TaskResult   workspace.TaskResultRecord   `json:"task_result"`
	HumanQA      *workspace.TaskHumanQARecord `json:"human_qa"`
	Integration  *FileBinding                 `json:"integration"`
	PullRequest  *FileBinding                 `json:"pull_request,omitempty"`
}

type DeliveryAttempt struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	State         string `json:"state"`
	Path          string `json:"path"`
	CandidateOID  string `json:"candidate_oid"`
	CandidateTree string `json:"candidate_tree"`
	InputSHA256   string `json:"input_sha256"`
}

type DeliveryState struct {
	Phase                string                `json:"phase"`
	OwnerClaim           string                `json:"owner_claim"`
	PermissionState      string                `json:"permission_state"`
	ActualPolicySHA256   *string               `json:"actual_policy_sha256"`
	ActualPolicyEvidence []Evidence            `json:"actual_policy_evidence"`
	Events               []DeliveryEventRecord `json:"events"`
	Candidates           []DeliveryCandidate   `json:"candidates"`
	Attempt              *DeliveryAttempt      `json:"attempt"`
	LastEventSHA256      *string               `json:"last_event_sha256"`
}

func deliveryEnv(kind string) Envelope   { return Envelope{"ply.workflow." + kind, 2} }
func deliveryRun(r WorkflowRequest) bool { return r.SchemaVersion == 2 && r.Delivery != nil }
