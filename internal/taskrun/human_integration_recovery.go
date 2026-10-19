package taskrun

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Recovery records absence and an explicit human revocation, never a fictitious
// process exit. Unmanaged writers remain the responsibility of that human.
type HumanOwnerRecoveryObservation struct {
	Transport         WorkflowTransport `json:"transport"`
	PaneInventorySHA  string            `json:"pane_inventory_sha256"`
	AgentInventorySHA string            `json:"agent_inventory_sha256"`
	ProcessLiveness   string            `json:"process_liveness"`
}

type HumanOwnerRecoveryPreview struct {
	Kind          string                         `json:"kind"`
	SchemaVersion int                            `json:"schema_version"`
	State         string                         `json:"state"`
	RunID         string                         `json:"run_id"`
	TaskResultID  workspace.TaskResultID         `json:"task_result_id"`
	CandidateOID  string                         `json:"candidate_oid"`
	SourceLocator string                         `json:"source_locator"`
	Confirmation  string                         `json:"confirmation_sha256"`
	Observation   *HumanOwnerRecoveryObservation `json:"observation,omitempty"`
	NextAction    string                         `json:"next_action"`
}

type HumanOwnerRecoveryAnswer struct {
	Confirmation   string `json:"confirmation_sha256"`
	Actor          string `json:"actor"`
	Answer         string `json:"answer"`
	Observation    string `json:"observation"`
	StartedAtUTC   string `json:"started_at_utc"`
	CompletedAtUTC string `json:"completed_at_utc"`
}

type humanOwnerRecovery struct {
	Answer      HumanOwnerRecoveryAnswer      `json:"human"`
	Observation HumanOwnerRecoveryObservation `json:"observation"`
}

func sourceContains(source, path string) bool {
	if path == "" {
		return false
	}
	if physical, err := filepath.EvalSymlinks(path); err == nil {
		path = physical
	}
	path = filepath.Clean(path)
	return path == source || strings.HasPrefix(path, source+string(filepath.Separator))
}

// The confirmation binds only this candidate and its native history. Complete
// inventories are rechecked after input, without requiring unrelated panes to
// remain byte-identical while the human reads the prompt.
func ownerRecoveryConfirmation(s workflowState, c DeliveryCandidate) string {
	return digest(struct {
		Run, Request   string
		Candidate      workspace.TaskResultRecord
		Transport      deliveryRecoveryTransport
		AgentSessionID string
		Previous       *string
	}{s.Result.RunID, s.Result.RequestSHA256, c.TaskResult, deliveryTransportIdentity(s.Result.Transport), s.Result.Transport.AgentSessionID, s.Result.Delivery.LastEventSHA256})
}

func validHumanOwnerRecoveryRelease(s workflowState, c DeliveryCandidate, r deliveryOwnershipRelease, index int) bool {
	if r.InvokingExecutable != nil || r.ProviderExit != nil || r.HumanRecovery == nil || !validOwnerRecoveryAnswer(r.HumanRecovery.Answer) {
		return false
	}
	o := r.HumanRecovery.Observation
	if deliveryTransportIdentity(o.Transport) != deliveryTransportIdentity(s.Result.Transport) || o.Transport.AgentSessionID != s.Result.Transport.AgentSessionID || o.ProcessLiveness != "unknown" || !digestPattern.MatchString(o.PaneInventorySHA) || !digestPattern.MatchString(o.AgentInventorySHA) {
		return false
	}
	prior := *s.Result.Delivery
	prior.LastEventSHA256 = nil
	if index > 0 {
		prior.LastEventSHA256 = ptr(prior.Events[index-1].Binding.SHA256)
	}
	s.Result.Delivery = &prior
	return r.HumanRecovery.Answer.Confirmation == ownerRecoveryConfirmation(s, c)
}

func PreviewHumanOwnerRecovery(d Dependencies, root, id string, resultID workspace.TaskResultID) (HumanOwnerRecoveryPreview, error) {
	p := HumanOwnerRecoveryPreview{Kind: "ply.integration.owner-recovery", SchemaVersion: 1, State: "blocked", RunID: id, TaskResultID: resultID}
	a, err := nativeDeliveryCandidate(d, root, id, resultID)
	if err != nil {
		return p, err
	}
	s, err := workflowRead(root, id)
	if err != nil {
		return p, err
	}
	p.CandidateOID, p.SourceLocator = a.Candidate.OID, a.Candidate.TaskResult.SourceLocator
	if prior, err := deliveryRelease(s, a.Candidate); err != nil {
		return p, err
	} else if prior != nil {
		p.State, p.NextAction = "released", "Ownership is already released. Inspect the human integration plan."
		return p, nil
	}
	if attempt := s.Result.Delivery.Attempt; attempt != nil && attempt.State == "attempted" {
		return p, workflowError(4, "owner_active: resolve the reserved delivery effect before recovery")
	}
	if a.Agreement.Mode == workspace.DeliveryPullRequest && (s.Result.Delivery.Phase != "completed" || a.Candidate.PullRequest == nil) {
		return p, workflowError(4, "owner_active: finish the exact PR publication before recovering its source")
	}
	observation, err := observeAbsentDeliveryOwner(d, s)
	if err != nil {
		return p, err
	}
	p.State, p.Observation, p.Confirmation = "ready", &observation, ownerRecoveryConfirmation(s, a.Candidate)
	p.NextAction = "Confirm all writers are stopped and explicitly release this exact candidate. This records no QA or integration."
	return p, nil
}

func validOwnerRecoveryAnswer(in HumanOwnerRecoveryAnswer) bool {
	started, e1 := time.Parse(time.RFC3339Nano, in.StartedAtUTC)
	completed, e2 := time.Parse(time.RFC3339Nano, in.CompletedAtUTC)
	return digestPattern.MatchString(in.Confirmation) && plain(in.Actor, 1, 256) && in.Answer == "released" && (in.Observation == "" || plain(in.Observation, 1, 2000)) && e1 == nil && e2 == nil && !completed.Before(started)
}

func RecoverHumanDeliveryOwnership(d Dependencies, root, id string, resultID workspace.TaskResultID, in HumanOwnerRecoveryAnswer) (WorkflowRun, error) {
	if !validOwnerRecoveryAnswer(in) {
		return WorkflowRun{}, workflowError(2, "owner_recovery_requires_human_release: an actual released answer must bind the preview and caller")
	}
	err := deliveryReleaseUpdate(d, root, id, func(s *workflowState) error {
		a, err := nativeDeliveryCandidate(d, root, id, resultID)
		if err != nil {
			return err
		}
		if prior, err := deliveryRelease(*s, a.Candidate); err != nil || prior != nil {
			return err
		}
		if ownerRecoveryConfirmation(*s, a.Candidate) != in.Confirmation {
			return workflowError(4, "plan_changed: candidate or native history changed since ownership recovery preview")
		}
		observation, err := observeAbsentDeliveryOwner(d, *s)
		if err != nil {
			return err
		}
		return preserveDeliveryRelease(d, s, a.Candidate, "Human explicitly revoked the absent owner's writing authority after checking known writers.", "human_absent_owner_recovery", nil, nil, &humanOwnerRecovery{in, observation})
	})
	return deliveryReadback(d, root, id, err)
}
