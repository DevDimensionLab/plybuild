package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

// WorkerContract is an explicit one-Spec execution choice, never an upgrade of
// an existing interactive owner. The controller owns native effects; providers
// receive only Task-file authority and their own disposable output directory.
type WorkerContract struct {
	SchemaVersion       int `json:"schema_version"`
	MaxAttempts         int `json:"max_attempts"`
	AgentTimeoutSeconds int `json:"agent_timeout_seconds"`
	MaxProviderRetries  int `json:"max_provider_retries,omitempty"`
}

type WorkerStatus struct {
	SessionID           string       `json:"session_id"`
	Phase               string       `json:"phase"`
	Stage               string       `json:"stage"`
	Generation          int          `json:"generation"`
	OwnerPID            int          `json:"owner_pid"`
	OwnerStartedAt      string       `json:"owner_started_at"`
	OwnerState          string       `json:"owner_state,omitempty"`
	ControlAction       string       `json:"control_action,omitempty"`
	ControlRequestedAt  string       `json:"control_requested_at,omitempty"`
	Iteration           int          `json:"iteration"`
	Job                 *WorkerJob   `json:"job,omitempty"`
	DeveloperSessions   []string     `json:"developer_sessions"`
	Review              *FileBinding `json:"review,omitempty"`
	ReviewProof         *FileBinding `json:"review_proof,omitempty"`
	CandidateBinary     string       `json:"candidate_binary,omitempty"`
	CandidateOID        string       `json:"candidate_oid,omitempty"`
	CandidateTree       string       `json:"candidate_tree,omitempty"`
	Pretest             *FileBinding `json:"pretest,omitempty"`
	Feedback            string       `json:"feedback,omitempty"`
	Reason              string       `json:"reason,omitempty"`
	UpdatedAt           string       `json:"updated_at"`
	RetryNotBefore      string       `json:"retry_not_before,omitempty"`
	ProviderRetries     int          `json:"provider_retries,omitempty"`
	LastProviderFailure string       `json:"last_provider_failure,omitempty"`
	QueueReleased       bool         `json:"queue_released,omitempty"`
}

type WorkerJob struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Directory  string `json:"directory"`
	IntentHash string `json:"intent_sha256"`
}

type workerOwner struct {
	root, runID string
	generation  int
	lock        *os.File
}

func headlessRun(r WorkflowRequest) bool {
	return r.Envelope == deliveryEnv("headless-run-request")
}

func deliveryOwnerSession(s workflowState) string {
	if headlessRun(s.Request) && s.Result.Worker != nil {
		return s.Result.Worker.SessionID
	}
	return s.Result.Transport.AgentSessionID
}

func workerDirectory(s workflowState) string {
	return filepath.Join(s.Result.Paths.RunRoot, "worker")
}

type workerControl struct {
	RunID       string `json:"run_id"`
	Action      string `json:"action"`
	RequestedAt string `json:"requested_at"`
}

func workerReadControl(s workflowState, o *WorkflowRun) error {
	w := *o.Worker
	o.Worker = &w
	w.ControlAction, w.ControlRequestedAt = "", ""
	var control workerControl
	if err := readValue(filepath.Join(workerDirectory(s), "control.json"), 64<<10, &control); err == nil {
		if control.RunID != o.RunID {
			return fmt.Errorf("worker control belongs to another Run")
		}
		if control.Action != "continue" {
			w.ControlAction, w.ControlRequestedAt = control.Action, control.RequestedAt
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Older controllers retain pause/wait actions after their transitions. Project
// only that known guidance; persisted history and unrelated diagnostics remain.
func workerCurrentAction(o WorkflowRun) WorkflowAction {
	a, w := o.NextAction, o.Worker
	if w == nil || len(o.Reasons) != 0 || (a != (WorkflowAction{"human", "Work is preserved; use worker resume for this same Run."}) && a.Actor != "worker_wait") {
		return a
	}
	if w.Phase == "blocked" || (w.Reason != "" && !(w.Phase == "waiting" && w.Reason == "usage_limit: same Run waits until "+w.RetryNotBefore)) {
		return a
	}
	if (w.ControlAction == "resume" && (w.Phase == "paused" || w.Phase == "stopped" || w.Phase == "running" || w.Phase == "waiting")) ||
		(w.ControlAction == "pause" && (w.Phase == "running" || w.Phase == "waiting")) ||
		(w.ControlAction == "stop" && (w.Phase == "running" || w.Phase == "waiting" || w.Phase == "paused")) {
		return WorkflowAction{"worker", w.ControlAction + " requested for this same Run; consumption at a safe boundary is not yet confirmed. Existing attempts retain ownership until their completion is known."}
	}
	switch w.Phase {
	case "running":
		if w.Reason == "" && w.RetryNotBefore == "" {
			return WorkflowAction{"worker", "Observe the same Run; cached worker state does not attest provider liveness or attempt completion."}
		}
	case "waiting":
		if w.RetryNotBefore != "" && (w.Reason == "" || w.Reason == "usage_limit: same Run waits until "+w.RetryNotBefore) {
			return WorkflowAction{"worker_wait", "The same Run waits until " + w.RetryNotBefore + "; continuing still requires confirmed attempt completion."}
		}
	}
	return a
}

func validateWorkerRequest(r WorkflowRequest) error {
	w := r.Worker
	if w == nil || w.SchemaVersion != 1 || w.MaxAttempts < 1 || w.MaxAttempts > 20 || w.AgentTimeoutSeconds < 10 || w.AgentTimeoutSeconds > 86400 || w.MaxProviderRetries < 0 || w.MaxProviderRetries > 1000 {
		return workflowError(2, "headless requires an explicit bounded one-Spec worker contract")
	}
	if r.Delivery == nil || r.ReturnMode != "delivery_owner" || r.Agreement != (Agreement{}) || r.Coordinator.ActorClaim != "" || r.Coordinator.MayRequestChanges {
		return workflowError(2, "headless requires its own delivery-owner contract")
	}
	if !key(r.RequestKey) || !strings.HasPrefix(r.PreparationID, "pre_") || len(r.PreparationID) != 68 || !digestPattern.MatchString(r.PreparationSHA256) || !r.HumanAuthority.Authorized || r.HumanAuthority.StartSurface != "human_authorized_worker" || !plain(r.HumanAuthority.ActorClaim, 1, 256) {
		return workflowError(2, "invalid worker authority or preparation")
	}
	if r.Herdr.Executable != (Executable{}) || r.Herdr.WorkspaceID != "" || r.Herdr.TabLabel != "" || r.CodexProjectTrust != nil || r.ClaudeProjectTrust != nil {
		return workflowError(2, "headless cannot claim a Herdr transport or interactive trust grant")
	}
	if r.Runtime.Mode != "headless" || r.Runtime.Provider != "codex" || !plain(r.Runtime.Model, 1, 256) || r.Runtime.ConfigProfile != nil || r.Runtime.PermissionBinding.ProfileID != "ply-worker" {
		return workflowError(2, "worker supports an explicitly selected Codex headless runtime")
	}
	if err := physical(r.WorkspaceRoot, false); err != nil {
		return err
	}
	if _, err := runtimeBindings(r.Runtime); err != nil {
		return err
	}
	p := r.Runtime.PermissionBinding
	if p.AuthorityKind != "launch_contract_pending_runtime_acceptance" || !digestPattern.MatchString(p.EffectivePolicySHA256) || len(p.Evidence) == 0 {
		return workflowError(2, "worker requires a preserved launch contract and runtime probe")
	}
	if err := deliveryEvidence(p.Evidence); err != nil {
		return err
	}
	c := r.Delivery
	if c.Agreement == nil || !c.Agreement.AutomaticAcceptance() || c.Agreement.HumanOwnedIntegration() || c.Agreement.Acceptance.RequireHumanQA || string(c.Agreement.Mode) != "local_epic_integration" {
		return workflowError(2, "worker requires an explicit automatic local Epic agreement without a human gate")
	}
	if !plain(c.OwnerClaim, 1, 256) || c.AllowLocalInstall || c.NotificationContext != nil || c.LocalIntegration != DeliveryLocalIntegrationBoundary(c.Agreement) || !plain(c.ReasoningEffort, 1, 64) {
		return workflowError(2, "worker owner, local return or effect boundary is invalid")
	}
	if !filepath.IsAbs(c.AcceptancePath) || filepath.Clean(c.AcceptancePath) != c.AcceptancePath {
		return workflowError(2, "worker requires the exact acceptance entrypoint")
	}
	if _, err := workflowBound(c.Goal, 4<<20); err != nil {
		return err
	}
	if _, err := workflowhandoff.ValidateDeliveryTaskRunDraft(r.HandoffDraft); err != nil {
		return err
	}
	return workflowhandoff.ValidateDeliveryMandate(r.HandoffDraft, c.Agreement)
}

func workerLiveOwner(d Dependencies, s workflowState) error {
	o, w := d.WorkerOwner, s.Result.Worker
	if o == nil || w == nil || o.root != s.Request.WorkspaceRoot || o.runID != s.Result.RunID || o.generation != w.Generation || o.lock == nil {
		return workflowError(4, "worker_owner_required: callbacks require the current native worker generation and held process lock")
	}
	info, err := o.lock.Stat()
	if err != nil {
		return workflowError(4, "worker ownership lock is closed")
	}
	expected, err := os.Stat(filepath.Join(workerDirectory(s), "owner.lock"))
	if err != nil || !os.SameFile(info, expected) {
		return workflowError(4, "worker ownership lock identity changed")
	}
	return nil
}

func workerValidateIdentity(s workflowState) error {
	if !headlessRun(s.Request) || s.Request.Worker == nil || s.Result.Worker == nil || s.Result.Worker.SessionID != "ply-worker:"+s.Result.RunID || s.Request.Runtime.Mode != "headless" || s.Request.Runtime.Provider != "codex" || s.Request.Delivery == nil || s.Result.Delivery == nil || s.Request.Herdr != (WorkflowRequest{}).Herdr || s.Result.Transport.WorkspaceID != "" || s.Result.Transport.TabID != "" || s.Result.Transport.PaneID != "" || s.Result.Transport.TerminalID != "" || s.Result.Transport.AgentSessionID != "" {
		return fmt.Errorf("run is not an explicitly registered headless worker")
	}
	return nil
}
