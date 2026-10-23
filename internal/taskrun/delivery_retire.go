package taskrun

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type DeliveryRetireInput struct {
	RetainEffects bool   `json:"retain_effects,omitempty"`
	Observer      string `json:"observer"`
	Reason        string `json:"reason"`
}

// DeliveryRetirement is the coordinator's selected end of one interactive
// owner. It preserves the Run, its refusals and its source as history. It
// grants nothing to a later Task copy, which receives its own new mandate.
type DeliveryRetirement struct {
	Mode            string                        `json:"mode,omitempty"`
	Provider        string                        `json:"provider,omitempty"`
	PreparationID   string                        `json:"preparation_id,omitempty"`
	DeliverySource  *ReportRetainedSource         `json:"delivery_source,omitempty"`
	DeliveryEffects []DeliveryRetainedEffect      `json:"delivery_effects,omitempty"`
	Kind            string                        `json:"kind"`
	SchemaVersion   int                           `json:"schema_version"`
	RunID           string                        `json:"run_id"`
	RequestSHA256   string                        `json:"request_sha256"`
	NativeSession   string                        `json:"native_session"`
	Input           DeliveryRetireInput           `json:"input"`
	PaneClosed      bool                          `json:"pane_closed"`
	Absence         HumanOwnerRecoveryObservation `json:"absence"`
	ActiveStop      *FileBinding                  `json:"active_stop"`
	WriterLiveness  *TaskWriterLiveness           `json:"writer_liveness,omitempty"`
	ReportEffects   []ReportRetainedEffect        `json:"report_effects,omitempty"`
	ReportSource    *ReportRetainedSource         `json:"report_source,omitempty"`
	RecordedAtUTC   string                        `json:"recorded_at_utc"`
}

func deliveryRetirement(s workflowState) (*DeliveryRetirement, error) {
	if s.Result.Delivery == nil {
		return nil, nil
	}
	for _, event := range s.Result.Delivery.Events {
		if event.Kind != "owner_retirement" {
			continue
		}
		raw, err := workflowBound(event.Binding, 1<<20)
		if err != nil {
			return nil, err
		}
		var r DeliveryRetirement
		if err = decode(raw, 1<<20, &r); err != nil {
			return nil, err
		}
		if r.Kind != "ply.workflow.owner-retirement" || r.SchemaVersion != 1 || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 {
			return nil, workflowError(4, "owner retirement binding changed")
		}

		if r.Mode != "" && r.Mode != deliveryRetainEffects {
			return nil, workflowError(4, "unsupported retirement mode")
		}
		if r.Mode == deliveryRetainEffects && (!r.Input.RetainEffects || r.NativeSession != s.Result.Transport.AgentSessionID || r.Provider != s.Request.Runtime.Provider || r.PreparationID != s.Request.PreparationID || r.DeliverySource == nil || r.ActiveStop == nil) {
			return nil, workflowError(4, "retained retirement identity binding differs")
		}
		if r.Mode == deliveryRetainEffects {
			source := r.DeliverySource
			if source.Worktree != s.Observed.Target.WorktreeLocator || source.Ref != s.Observed.Target.Ref || source.GitCommonDir != s.Observed.Target.GitCommonDir {
				return nil, workflowError(4, "retained source identity binding differs")
			}
		}
		return &r, nil
	}
	return nil, nil
}

// WorkflowRetireDelivery is the coordinator's single operation for a selected
// replacement. It needs no callback from the retired owner. Ply closes the
// bound pane itself, then requires positive absence of every known Task
// writer before it records retirement and releases only this queue slot.
// Unresolved attempts and unknown or partial effects stop before any effect.
func WorkflowRetireDelivery(d Dependencies, root, id string, in DeliveryRetireInput) (WorkflowRun, error) {
	if in.RetainEffects {
		return workflowRetireRetainingEffects(d, root, id, in)
	}
	return workflowRetireDelivery(d, root, id, in, nil)
}

func workflowRetireDelivery(d Dependencies, root, id string, in DeliveryRetireInput, beforeClose func(workflowState) error) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	if !plain(in.Observer, 1, 256) || strings.TrimSpace(in.Reason) == "" || len(in.Reason) > 1900 {
		return WorkflowRun{}, workflowError(2, "actual observer and bounded reason are required")
	}
	err := workflowUpdate(d, root, id, func(s *workflowState) error {
		if !deliveryRun(s.Request) {
			return workflowRetireReport(d, s, in)
		}
		if !deliveryRun(s.Request) || headlessRun(s.Request) || s.Result.Delivery == nil {
			return workflowError(4, "retirement requires an interactive delivery owner; use workflow worker for headless Runs")
		}
		prior, err := deliveryRetirement(*s)
		if err != nil {
			return err
		}
		if prior != nil {
			if prior.Input != in {
				return workflowError(4, "owner is already retired with a different recorded reason")
			}
			// A retry after a failed state publication indexes the durable event.
			s.Result.Delivery.Phase, s.Result.Round.State = "stopped", "stopped"
			return nil
		}
		dl := s.Result.Delivery
		switch {
		case dl.Phase == "completed":
			return workflowError(4, "delivery is complete; a new goal needs no retirement")
		case dl.OwnershipRelease != nil || dl.HumanIntegration != nil:
			return workflowError(4, "candidate ownership is already released to human integration; continue that plan")
		case dl.Attempt != nil && dl.Attempt.State != "recorded":
			return workflowError(4, "active_writer: "+dl.Attempt.Kind+" attempt is "+dl.Attempt.State+"; resolve its actual effect before replacement")
		}
		// Every unresolved stop counts, not only the latest pending one.
		stops, err := deliveryUnresolvedStops(*s)
		if err != nil {
			return err
		}
		for _, stop := range stops {
			if e := stop.Input.EffectStatus; e == DeliveryEffectUnknown || e == DeliveryEffectPartial {
				corrected, err := deliveryQuotaAllowsRetirement(*s, stop)
				if err != nil {
					return err
				}
				if corrected {
					continue
				}
				return workflowError(4, "unknown_effect: preserved stop "+stop.Input.EventID+" has "+e+" effects for action "+stop.Input.Action+"; establish them from actual observations before replacement")
			}
		}
		if s.OwnerResume != nil && (s.OwnerResumeStep == ownerResumeTabAttempted || s.OwnerResumeStep == ownerResumePromptAttempt) {
			return workflowError(4, "unknown_effect: owner resume tab or continuation input has an unconfirmed effect; reconcile that same reserved attempt before retirement")
		}
		if beforeClose != nil {
			if err := beforeClose(*s); err != nil {
				return err
			}
		}
		closed, err := deliveryRetireClosePane(d, *s)
		if err != nil {
			if beforeClose != nil && backupRetirementMustPreserve(err) {
				return &backupRetirementEffectError{Action: "inspect/close bound retirement pane", Err: err}
			}
			return err
		}
		absence, liveness, err := deliveryRetireAbsenceSelected(d, *s, beforeClose != nil)
		if err != nil {
			return err
		}
		if beforeClose != nil {
			if err := beforeClose(*s); err != nil {
				return err
			}
		}
		// Absence establishes no writer remains, but cannot establish what the
		// writer did between the first source check and the actual pane close.
		for _, stop := range stops {
			if e := stop.Input.EffectStatus; e == DeliveryEffectUnknown || e == DeliveryEffectPartial {
				corrected, err := deliveryQuotaAllowsRetirement(*s, stop)
				if err != nil {
					return workflowError(4, "quota_source_changed_after_close: writer absence observed; retirement and queue release remain blocked. "+err.Error())
				}
				if !corrected {
					return workflowError(4, "quota correction absent after close; retirement and queue release remain blocked")
				}
			}
		}
		r := DeliveryRetirement{Kind: "ply.workflow.owner-retirement", SchemaVersion: 1, RunID: id, RequestSHA256: s.Result.RequestSHA256, NativeSession: s.Result.Transport.AgentSessionID, Input: in, PaneClosed: closed, Absence: absence, ActiveStop: dl.ActiveStop, WriterLiveness: liveness, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		if _, err = deliveryAppendEvent(d, s, "owner-retirement", "owner_retirement", r); err != nil {
			return err
		}
		dl.Phase, s.Result.Round.State = "stopped", "stopped"
		return nil
	})
	if err == nil {
		err = deliveryRetireReleaseQueue(d, root, id, in.Reason)
	}
	return deliveryReadback(d, root, id, err)
}

// deliveryUnresolvedStops replays stop history: each resolution addresses the
// latest pending stop, so earlier stops stay unresolved beneath it.
func deliveryUnresolvedStops(s workflowState) ([]DeliveryStopRecord, error) {
	var pending []DeliveryStopRecord
	for _, event := range s.Result.Delivery.Events {
		switch event.Kind {
		case "coordinator_stop":
			raw, err := workflowBound(event.Binding, 2<<20)
			if err != nil {
				return nil, err
			}
			var stop DeliveryStopRecord
			if err = decode(raw, 2<<20, &stop); err != nil {
				return nil, err
			}
			pending = append(pending, stop)
		case "stop_resolution":
			if len(pending) > 0 {
				pending = pending[:len(pending)-1]
			}
		}
	}
	return pending, nil
}

// deliveryRetireClosePane closes only the exact bound pane. An already absent
// pane is not proof of exit; deliveryRetireAbsence decides that afterwards.
func deliveryRetireClosePane(d Dependencies, s workflowState) (bool, error) {
	raw, err := workflowCall(d, s.Request, "pane", "get", s.Result.Transport.PaneID)
	var call *HerdrCallError
	if errors.As(err, &call) && call.ReadOnly && call.ProviderCode == "pane_not_found" {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var pane struct {
		Pane workflowAgent `json:"pane"`
	}
	if err = json.Unmarshal(raw, &pane); err != nil {
		return false, workflowError(4, "owner_unknown: invalid bound pane observation")
	}
	t := s.Result.Transport
	if pane.Pane.WorkspaceID != t.WorkspaceID || pane.Pane.TabID != t.TabID || pane.Pane.PaneID != t.PaneID || pane.Pane.TerminalID != t.TerminalID {
		return false, workflowError(4, "owner_unknown: bound pane identity differs; no process was stopped")
	}
	if _, err = workflowCall(d, s.Request, "pane", "close", t.PaneID); err != nil {
		return true, err
	}
	return true, nil
}

func deliveryRetireAbsence(d Dependencies, s workflowState) (HumanOwnerRecoveryObservation, *TaskWriterLiveness, error) {
	return deliveryRetireAbsenceSelected(d, s, false)
}

func deliveryRetireAbsenceSelected(d Dependencies, s workflowState, selected bool) (HumanOwnerRecoveryObservation, *TaskWriterLiveness, error) {
	limit := d.HerdrTimeout
	if limit <= 0 {
		limit = 15 * time.Second
	}
	deadline := time.Now().Add(limit)
	for {
		observed := (taskWriters{d, s}).retirement()
		if observed.state == taskWritersAbsent {
			return observed.absence, observed.liveness, nil
		}
		if observed.state == taskWritersUnknown {
			if selected && backupRetirementMustPreserve(observed.inspectionError) {
				return observed.absence, observed.liveness, &backupRetirementEffectError{Action: "inspect known writers after retirement close", Err: observed.inspectionError}
			}
			return observed.absence, observed.liveness, observed.err
		}
		if !time.Now().Before(deadline) {
			return observed.absence, observed.liveness, workflowError(4, "owner_active: writer absence is not established after closing the bound pane; nothing was recorded or released ("+observed.err.Error()+")")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func deliveryRetireReleaseQueue(d Dependencies, root, id, reason string) error {
	s, err := workflowRead(root, id)
	if err != nil {
		return err
	}
	if !deliveryRun(s.Request) {
		return workflowReleaseReportQueue(d, s)
	}
	a := s.Request.Delivery.Agreement
	if a == nil {
		return workflowError(4, "retired owner recorded; this historical Run has no structured queue target to release")
	}
	var draft struct {
		Binding struct {
			TaskID workspace.TaskID `json:"task_id"`
		} `json:"task_spec_binding"`
	}
	if err = json.Unmarshal(s.Request.HandoffDraft, &draft); err != nil || draft.Binding.TaskID == "" {
		return workflowError(4, "retired owner recorded; Task binding is missing")
	}
	wd := workspace.WithTaskContentScope(d.Workspace, draft.Binding.TaskID)
	target := workspace.QueueTargetInput{ProjectID: a.ProjectID, RepoID: a.RepoID, EpicID: a.EpicID}
	q, err := workspace.ListTaskQueueVersion(wd, target, false, 2)
	if err != nil {
		return err
	}
	if q.HasPreparation(s.Request.PreparationID) {
		_, err = workspace.CloseTaskQueue(wd, target, s.Request.PreparationID, q.Revision, "Selected owner replacement: "+reason, "release")
		return err
	}
	readback, err := workspace.ShowTaskPreparation(wd, s.Request.PreparationID)
	if err != nil {
		return err
	}
	if readback.Disposition != "released" {
		return workflowError(4, "retired owner recorded; its queue slot is "+readback.Disposition+", not current or released")
	}
	return nil
}
