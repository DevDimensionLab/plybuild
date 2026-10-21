package taskrun

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// The opt-in readback adds only observations; native retirement events,
// publication keys and execute intents remain the durable recovery records.
type BackupGoalResult struct {
	workspace.BackupGoalRegistrationResult
	Retirement *BackupRetirement `json:"retirement,omitempty"`
}

type BackupRetirement struct {
	RunID         string                            `json:"run_id"`
	RequestSHA256 string                            `json:"request_sha256"`
	PreparationID string                            `json:"preparation_id"`
	Source        workspace.PlanWorktreeObservation `json:"source"`
	State         string                            `json:"state"`
	NextAction    string                            `json:"next_action"`
	Owner         WorkflowTransport                 `json:"owner"`
	Stop          DeliveryStopRecord                `json:"stop"`
	StopBinding   FileBinding                       `json:"stop_binding"`
	Selected      DeliveryRetireInput               `json:"selected"`
	Binding       string                            `json:"binding"`
	Completed     bool                              `json:"completed"`
}

type backupOwnerIdentity struct {
	WorkspaceID    string `json:"workspace_id"`
	TabID          string `json:"tab_id"`
	PaneID         string `json:"pane_id"`
	TerminalID     string `json:"terminal_id"`
	AgentSessionID string `json:"agent_session_id"`
}

func backupRetirement(s workflowState, selected *DeliveryRetireInput, retired *DeliveryRetirement) (*BackupRetirement, string, error) {
	if selected == nil {
		return nil, retired.RecordedAtUTC, nil
	}
	if !plain(selected.Observer, 1, 256) || !plain(selected.Reason, 1, 1900) {
		return nil, "", workflowError(2, "selected retirement requires actual coordinator and bounded reason")
	}
	if retired != nil && retired.Input != *selected {
		return nil, "", workflowError(4, "selected retirement differs from the recorded coordinator/reason")
	}
	dl := s.Result.Delivery
	if s.OwnerResume != nil {
		return nil, "", workflowError(4, "backup_goal_resume_unsupported: source has owner resume history")
	}
	if dl.ActiveStop == nil {
		return nil, "", workflowError(4, "backup_goal_prework_stop_required: selected retirement needs a preserved known not_started stop")
	}
	stops, err := deliveryUnresolvedStops(s)
	if err != nil {
		return nil, "", err
	}
	for _, stop := range stops {
		if stop.Input.Rule != "" || stop.Input.EffectStatus != DeliveryEffectNotStarted || stop.Input.Stage == "delivery" {
			return nil, "", workflowError(4, "backup_goal_stop_unsupported: preserve rule refusals and unknown/partial effects; no retirement performed")
		}
	}
	raw, err := workflowBound(*dl.ActiveStop, 2<<20)
	if err != nil {
		return nil, "", err
	}
	var stop DeliveryStopRecord
	if err = decode(raw, 2<<20, &stop); err != nil {
		return nil, "", err
	}
	if len(stops) == 0 || !equal(stop, stops[len(stops)-1]) || stop.Input.NativeSession == "" || stop.Input.NativeSession != s.Result.Transport.AgentSessionID || stop.Input.Provider != s.Request.Runtime.Provider || strings.TrimSpace(stop.Input.Action) == "" {
		return nil, "", workflowError(4, "backup_goal_stop_binding_conflict: requires the exact original stopped owner/action")
	}
	if retired != nil && (retired.ActiveStop == nil || *retired.ActiveStop != *dl.ActiveStop || retired.NativeSession != s.Result.Transport.AgentSessionID) {
		return nil, "", workflowError(4, "backup_goal_stop_binding_conflict: retirement stop/owner changed")
	}
	p := &BackupRetirement{RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, PreparationID: s.Request.PreparationID, Source: s.Observed.Target, Owner: s.Result.Transport, Stop: stop, StopBinding: *dl.ActiveStop, Selected: *selected}
	// Bind semantic identity, not coordinator observations. Native follow can
	// change state, observation time and launch-pending without replacing an
	// owner. Keep the complete transport above for truthful readback.
	raw, err = Canonical(struct {
		RunID         string                            `json:"run_id"`
		RequestSHA256 string                            `json:"request_sha256"`
		PreparationID string                            `json:"preparation_id"`
		Source        workspace.PlanWorktreeObservation `json:"source"`
		Owner         backupOwnerIdentity               `json:"owner"`
		StopBinding   FileBinding                       `json:"stop_binding"`
		Selected      DeliveryRetireInput               `json:"selected"`
	}{p.RunID, p.RequestSHA256, p.PreparationID, p.Source,
		backupOwnerIdentity{p.Owner.WorkspaceID, p.Owner.TabID, p.Owner.PaneID, p.Owner.TerminalID, p.Owner.AgentSessionID}, p.StopBinding, p.Selected})
	if err != nil {
		return nil, "", err
	}
	p.Binding = hash(raw)
	p.Completed = retired != nil
	p.State, p.NextAction = "selected", "Apply identical input and confirmation to retire, register and optionally start."
	if retired != nil {
		p.State, p.NextAction = "retired", "Retirement is preserved. Recover the same fresh registration/publication keys; no old owner restart."
	}
	// Stable original evidence time survives our own retirement/publication phases.
	return p, stop.ObservedAtUTC, nil
}

func WorkflowBackupGoal(d Dependencies, root, id string, in BackupGoalInput, apply bool, confirmation string) (out BackupGoalResult, err error) {
	if in.Retire == nil {
		return workflowBackupGoal(d, root, id, in, apply, confirmation)
	}
	out, err = workflowBackupGoal(d, root, id, in, false, "")
	if err != nil || !apply {
		return out, err
	}
	if confirmation != out.Confirmation {
		return out, workflowError(4, "backup_goal_confirmation_conflict: selected retirement/fresh input differs from preview")
	}
	// Recheck all source and fresh registration guards under the native Run lock
	// immediately before close. No close can precede collision/default validation.
	_, err = workflowRetireDelivery(d, root, id, *in.Retire, func(workflowState) error {
		p, e := workflowBackupGoal(d, root, id, in, false, "")
		if e != nil {
			return e
		}
		if p.Confirmation != confirmation {
			return workflowError(4, "backup_goal_confirmation_conflict: source/input changed before retirement")
		}
		return nil
	})
	if err != nil {
		out.Retirement.State, out.Retirement.NextAction = "retirement_blocked", "Inspect the native retirement diagnosis and preserved stop before dependent registration/start; do not replay a denied or unknown action."
		var effect *backupRetirementEffectError
		if errors.As(err, &effect) {
			// Preserve an actual failed native action, never a quota guess or
			// a fabricated successful absence. A retry must not repeat it.
			s, readErr := workflowRead(root, id)
			if readErr != nil {
				return out, fmt.Errorf("%w; failure preservation blocked: %v", err, readErr)
			}
			evidence := filepath.Join(s.Result.Paths.RunRoot, "delivery", "backup-retirement-failure.json")
			raw, _ := Canonical(struct {
				Action string `json:"action"`
				Error  string `json:"error"`
			}{effect.Action, effect.Err.Error()})
			if e := d.writeOnce(evidence, raw); e != nil {
				return out, fmt.Errorf("%w; failure preservation blocked: %v", err, e)
			}
			effectStatus := DeliveryEffectUnknown
			var call *HerdrCallError
			if errors.As(effect.Err, &call) && call.ReadOnly {
				effectStatus = DeliveryEffectNotStarted
			}
			_, e := WorkflowObserveDeliveryStop(d, root, id, DeliveryStopInput{EventID: "backup-retirement-runtime-stop", Provider: s.Request.Runtime.Provider, NativeSession: s.Result.Transport.AgentSessionID, Stage: "delivery", Observer: in.Retire.Observer, Reason: "Selected retirement failed; dependent registration/start stopped", Action: effect.Action, EffectStatus: effectStatus, EvidencePath: evidence})
			if e != nil {
				return out, fmt.Errorf("%w; failure preservation blocked: %v", err, e)
			}
		}
		return out, err
	}
	out.Retirement.Completed = true
	out.Retirement.State, out.Retirement.NextAction = "retired", "Retirement completed. Repeat identical input to recover registration; preserve the same Task/Spec IDs."
	if d.Fault != nil {
		if err = d.Fault("backup_after_retirement"); err != nil {
			return out, err
		}
	}
	return workflowBackupGoal(d, root, id, in, true, confirmation)
}

// Used only by the selected prefix to distinguish a real native effect failure
// from pre-effect validation and publication failures.
type backupRetirementEffectError struct {
	Action string
	Err    error
}

func (e *backupRetirementEffectError) Error() string { return e.Err.Error() }
func (e *backupRetirementEffectError) Unwrap() error { return e.Err }

func backupRetirementMustPreserve(err error) bool {
	var call *HerdrCallError
	// Known validation and transient read diagnostics remain retryable. Only
	// an actual denied read or uncertain mutating call becomes a native stop.
	return errors.As(err, &call) && (!call.ReadOnly || call.ProviderCode == "permission_denied")
}
