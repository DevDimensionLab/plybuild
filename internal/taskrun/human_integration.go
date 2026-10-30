package taskrun

import (
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"time"
)

type HumanIntegrationAcceptance struct {
	Kind          string                                   `json:"kind"`
	SchemaVersion int                                      `json:"schema_version"`
	RunID         string                                   `json:"run_id"`
	RequestSHA256 string                                   `json:"request_sha256"`
	PlanID        string                                   `json:"plan_id"`
	PlanSHA256    string                                   `json:"plan_sha256"`
	Decision      FileBinding                              `json:"decision"`
	HumanQA       workspace.TaskHumanQARecord              `json:"human_qa"`
	Release       FileBinding                              `json:"release"`
	Accepted      bool                                     `json:"accepted"`
	Attestation   workflowhandoff.DeliveryHumanAttestation `json:"attestation"`
}

type deliveryOwnershipRelease struct {
	Kind               string                       `json:"kind"`
	SchemaVersion      int                          `json:"schema_version"`
	RunID              string                       `json:"run_id"`
	RequestSHA256      string                       `json:"request_sha256"`
	TaskID             workspace.TaskID             `json:"task_id"`
	TaskResultID       workspace.TaskResultID       `json:"task_result_id"`
	ResultOID          string                       `json:"result_oid"`
	ResultTree         string                       `json:"result_tree"`
	SourceRef          string                       `json:"source_ref"`
	SourceLocator      string                       `json:"source_locator"`
	OwnerClaim         string                       `json:"owner_claim"`
	SessionID          string                       `json:"session_id"`
	Reason             string                       `json:"reason"`
	RecordedAtUTC      string                       `json:"recorded_at_utc"`
	Origin             string                       `json:"origin"`
	InvokingExecutable *Executable                  `json:"invoking_executable,omitempty"`
	ProviderExit       *DeliveryStartupExitEvidence `json:"provider_exit,omitempty"`
}

// deliveryRelease checks the current candidate and the immutable release event.
// An idle UI, technical verification, or an old release never proves inactivity.
func deliveryRelease(s workflowState, c DeliveryCandidate) (*FileBinding, error) {
	if s.Result.Delivery == nil || s.Result.Delivery.OwnershipRelease == nil {
		return nil, nil
	}
	b := *s.Result.Delivery.OwnershipRelease
	raw, err := workflowBound(b, 1<<20)
	if err != nil {
		return nil, err
	}
	var r deliveryOwnershipRelease
	if err = decode(raw, 1<<20, &r); err != nil {
		return nil, err
	}
	_, timeErr := time.Parse(time.RFC3339Nano, r.RecordedAtUTC)
	if r.Kind != "PlyDeliveryOwnerRelease@1" || r.SchemaVersion != 1 || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.TaskID != c.TaskResult.TaskID || r.TaskResultID != c.TaskResult.ID || r.ResultOID != c.OID || r.ResultTree != c.Tree || r.SourceRef != c.TaskResult.SourceRef || r.SourceLocator != c.TaskResult.SourceLocator || r.OwnerClaim != s.Request.Delivery.OwnerClaim || r.SessionID != s.Result.SessionID || !plain(r.Reason, 1, 2000) || timeErr != nil {
		return nil, workflowError(4, "owner_release_invalid: release differs from the exact native candidate or owner")
	}
	for _, event := range s.Result.Delivery.Events {
		if event.Kind == "ownership_release" && equal(event.Binding, b) {
			return &b, nil
		}
	}
	return nil, workflowError(4, "owner_release_invalid: release is absent from native history")
}

// PreviewHumanIntegration is read-only, including when the owner is active.
