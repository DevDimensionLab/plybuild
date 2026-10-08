package taskrun

import (
	"os"
	"time"
)

// DeliveryRuntimeObservation is the current owner's explicit observation, at
// the same trust boundary as initial acceptance. Herdr attests session identity,
// not effective permissions. An installed binary cannot supply this claim.
type DeliveryRuntimeObservation struct {
	Envelope
	RunID              string                       `json:"run_id"`
	RequestSHA256      string                       `json:"request_sha256"`
	SessionID          string                       `json:"session_id"`
	Context            FileBinding                  `json:"context"`
	Acceptance         FileBinding                  `json:"acceptance"`
	RuntimeClaim       RuntimeClaim                 `json:"runtime_claim"`
	DeliveryPermission DeliveryPermissionAcceptance `json:"delivery_permission"`
	ObservedAtUTC      string                       `json:"observed_at_utc"`
}

type DeliveryLauncherObservation struct {
	Binding Executable `json:"binding"`
	State   string     `json:"state"`
}

func deliveryLauncherObservation(executable Executable) *DeliveryLauncherObservation {
	out := &DeliveryLauncherObservation{Binding: executable, State: "present"}
	if err := verifyExecutable(executable); err != nil {
		out.State = "changed_or_unavailable"
		if os.IsNotExist(err) {
			out.State = "missing"
		}
	}
	return out
}

func deliveryRuntimeObservation(s workflowState, binding FileBinding) error {
	return deliveryRuntimeObservationBinding(s, binding, true)
}

func deliveryRuntimeObservationBinding(s workflowState, binding FileBinding, policy bool) error {
	raw, err := workflowBound(binding, 1<<20)
	if err != nil {
		return deliveryContinuityError("runtime_observation", "current owner evidence is missing or altered: "+err.Error())
	}
	var observation DeliveryRuntimeObservation
	if err = decode(raw, 1<<20, &observation); err != nil {
		return deliveryContinuityError("runtime_observation", "unsupported current owner observation: "+err.Error())
	}
	if observation.Envelope != (Envelope{"ply.workflow.delivery-runtime-observation", 1}) || observation.RunID != s.Result.RunID || observation.RequestSHA256 != s.Result.RequestSHA256 || observation.SessionID != s.Result.SessionID || observation.Context != (FileBinding{s.Result.Paths.Context, s.ContextSHA256}) || s.Acceptance == nil || observation.Acceptance != *s.Acceptance {
		return deliveryContinuityError("runtime_observation", "current owner evidence must bind this exact run, original context and acceptance")
	}
	if _, err = time.Parse(time.RFC3339Nano, observation.ObservedAtUTC); err != nil {
		return deliveryContinuityError("runtime_observation", "current owner observation requires its actual UTC observation time")
	}
	c := observation.RuntimeClaim
	if c.NativeSessionID == nil || *c.NativeSessionID != s.Result.Transport.AgentSessionID || s.Result.Delivery == nil || !equal(c.EffectivePolicySHA256, s.Result.Delivery.ActualPolicySHA256) || !equal(observation.DeliveryPermission.ActualPolicyEvidence, s.Result.Delivery.ActualPolicyEvidence) {
		return deliveryContinuityError("continuation_authority", "observed current session or effective policy differs from accepted authority; resolve the actual change with the original owner")
	}
	if policy {
		if err = positiveDeliveryClaim(Acceptance{RuntimeClaim: c}, observation.DeliveryPermission, s.Request); err != nil {
			return deliveryContinuityError("continuation_authority", "current owner cannot confirm unchanged runtime authority: "+err.Error())
		}
	}
	return nil
}

func deliveryRuntimeObservationFresh(d Dependencies, binding FileBinding) error {
	raw, err := workflowBound(binding, 1<<20)
	if err != nil {
		return err
	}
	var observation DeliveryRuntimeObservation
	if err = decode(raw, 1<<20, &observation); err != nil {
		return err
	}
	observed, err := time.Parse(time.RFC3339Nano, observation.ObservedAtUTC)
	if err != nil || observed.After(d.Now().Add(30*time.Second)) || observed.Before(d.Now().Add(-10*time.Minute)) {
		return deliveryContinuityError("runtime_observation", "current owner evidence is stale or future-dated; observe the actual session and policy again without changing the original acceptance")
	}
	return nil
}

func deliveryContinuationConfirmation(s workflowState, source, control Executable, capability string, observation *FileBinding) string {
	values := []any{s, source, control, capability}
	if observation != nil {
		values = append(values, observation)
	}
	return digest(values)
}

func deliveryLiveOwner(d Dependencies, s workflowState) error {
	a, err := workflowAgentGet(d, s, false)
	if err != nil {
		return deliveryContinuityError("continuation_session", "observe the original provider/session without restarting: "+err.Error())
	}
	if a.ForegroundCWD != s.Observed.Target.WorktreeLocator || a.LaunchPending != nil && *a.LaunchPending {
		return deliveryContinuityError("continuation_session", "live owner foreground cwd is missing or differs, or a provider launch is still pending")
	}
	switch a.Status {
	case "working", "idle", "done", "blocked":
		return nil
	default:
		return deliveryContinuityError("continuation_session", "live owner status is missing or unknown; inspect the same Herdr session")
	}
}
