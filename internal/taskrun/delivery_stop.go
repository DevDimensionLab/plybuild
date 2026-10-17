package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DeliveryStopInput struct {
	Action        string `json:"action,omitempty"`
	EffectStatus  string `json:"effect_status,omitempty"`
	Rule          string `json:"rule,omitempty"`
	EventID       string `json:"event_id"`
	Provider      string `json:"provider"`
	NativeSession string `json:"native_session"`
	Stage         string `json:"stage"`
	Observer      string `json:"observer"`
	Reason        string `json:"reason"`
	EvidencePath  string `json:"evidence_path"`
}

func WorkflowObserveDeliveryStop(d Dependencies, root, id string, in DeliveryStopInput) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	if in.EventID == "" || len(in.EventID) > 128 || strings.TrimSpace(in.Observer) == "" || len(in.Observer) > 256 || strings.TrimSpace(in.Reason) == "" || len(in.Reason) > 4096 {
		return WorkflowRun{}, workflowError(2, "event ID, actual observer and bounded reason are required")
	}
	switch in.Stage {
	case "startup", "acceptance", "report":
	case "receipt", "work", "verification", "delivery":
		if strings.TrimSpace(in.Action) == "" || in.EffectStatus == "" {
			return WorkflowRun{}, workflowError(2, "action and known effect status are required for an action refusal")
		}
	default:
		return WorkflowRun{}, workflowError(2, "unsupported denial stage")
	}
	if len(in.Action) > 2000 || len(in.Rule) > 256 {
		return WorkflowRun{}, workflowError(2, "denial action or rule exceeds limit")
	}
	if in.EffectStatus != "" && in.EffectStatus != "not_started" && in.EffectStatus != "none_observed" && in.EffectStatus != "partial" && in.EffectStatus != "unknown" {
		return WorkflowRun{}, workflowError(2, "unsupported known effect status")
	}
	evidence, raw, err := deliveryReadBinding(in.EvidencePath, 1<<20)
	if err != nil {
		return WorkflowRun{}, err
	}
	err = workflowUpdate(d, root, id, func(s *workflowState) error {
		if !deliveryRun(s.Request) || headlessRun(s.Request) || s.Result.Delivery == nil {
			return workflowError(4, "stop observation requires an interactive delivery")
		}
		if in.Provider != s.Request.Runtime.Provider || in.NativeSession != s.Result.Transport.AgentSessionID || in.Stage != "startup" && in.NativeSession == "" {
			return workflowError(4, "observation must name the bound provider and native session")
		}
		for _, event := range s.Result.Delivery.Events {
			if event.ID != in.EventID {
				continue
			}
			if event.Kind != "coordinator_stop" {
				return workflowError(4, "event ID is already used")
			}
			var old DeliveryStopRecord
			b, e := workflowBound(event.Binding, 2<<20)
			if e != nil {
				return e
			}
			if e = decode(b, 2<<20, &old); e != nil {
				return e
			}
			if !equal(old.Input, in) || old.Source != evidence {
				return workflowError(4, "stop event already has different immutable evidence")
			}
			return nil
		}
		if s.Result.Delivery.Phase == "completed" || s.Result.Delivery.OwnershipRelease != nil {
			return workflowError(4, "delivery owner is already closed")
		}
		path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "stop-evidence", strings.TrimPrefix(evidence.SHA256, "sha256:"))
		if err := d.writeOnce(path, raw); err != nil {
			return err
		}
		record := DeliveryStopRecord{Kind: "ply.workflow.coordinator-stop", SchemaVersion: 1, RunID: id, Input: in, Source: evidence, Evidence: FileBinding{path, evidence.SHA256}, ObservedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		binding, err := deliveryAppendEvent(d, s, in.EventID, "coordinator_stop", record)
		if err != nil {
			return err
		}
		s.Result.Delivery.ActiveStop = &binding
		// Do not change acceptance, ownership, transport liveness or attempts.
		return nil
	})
	return deliveryReadback(d, root, id, err)
}

type DeliveryStopRecord struct {
	Kind          string            `json:"kind"`
	SchemaVersion int               `json:"schema_version"`
	RunID         string            `json:"run_id"`
	Input         DeliveryStopInput `json:"input"`
	Source        FileBinding       `json:"source"`
	Evidence      FileBinding       `json:"evidence"`
	ObservedAtUTC string            `json:"observed_at_utc"`
}

type DeliveryStopResolution struct {
	Kind                string      `json:"kind"`
	RunID               string      `json:"run_id"`
	Stop                FileBinding `json:"stop"`
	NativeSession       string      `json:"native_session"`
	Resolution          string      `json:"resolution"`
	PermissionConfirmed bool        `json:"permission_confirmed"`
}

// Resolution records the current owner's explicit cause repair. It starts no
// process, resends no input, grants no acceptance and releases no writer.
func WorkflowResolveDeliveryStop(d Dependencies, root, id, contextPath, stopSHA, resolution string, permissionConfirmed bool) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	if strings.TrimSpace(resolution) == "" || len(resolution) > 4096 || !permissionConfirmed {
		return WorkflowRun{}, workflowError(4, "resolution requires the actual repaired cause and confirmed current permission; do not repeat a denied action")
	}
	err := workflowUpdate(d, root, id, func(s *workflowState) error {
		if !deliveryRun(s.Request) || s.Result.Delivery == nil || headlessRun(s.Request) {
			return workflowError(4, "resolution requires the same interactive owner")
		}
		if err := workflowCallback(d, *s, contextPath); err != nil {
			return err
		}
		if err := deliveryLiveOwner(d, *s); err != nil {
			return err
		}
		if s.Result.Delivery.OwnershipRelease != nil || s.Result.Delivery.Phase == "completed" {
			return workflowError(4, "owner is closed")
		}
		if a := s.Result.Delivery.Attempt; a != nil && a.State != "recorded" {
			return workflowError(4, "an unresolved attempt blocks continuation")
		}
		// A completed resolution can be replayed while other denials remain.
		// Only a new transition must address the current pending stop.
		for _, event := range s.Result.Delivery.Events {
			if event.Kind != "stop_resolution" {
				continue
			}
			var old DeliveryStopResolution
			b, e := workflowBound(event.Binding, 1<<20)
			if e != nil {
				return e
			}
			if e = decode(b, 1<<20, &old); e != nil {
				return e
			}
			if old.Stop.SHA256 == stopSHA && old.Resolution == resolution && old.NativeSession == s.Result.Transport.AgentSessionID {
				return nil
			}
		}
		if s.Result.Delivery.ActiveStop == nil {
			return workflowError(4, "no matching unresolved stop")
		}
		stop := *s.Result.Delivery.ActiveStop
		if stop.SHA256 != stopSHA {
			return workflowError(4, "resolution must bind the current stop digest")
		}
		if _, err := workflowBound(stop, 2<<20); err != nil {
			return err
		}
		r := DeliveryStopResolution{Kind: "ply.workflow.stop-resolution", RunID: id, Stop: stop, NativeSession: s.Result.Transport.AgentSessionID, Resolution: resolution, PermissionConfirmed: permissionConfirmed}
		if _, err := deliveryAppendEvent(d, s, "resolve-"+strings.TrimPrefix(stopSHA, "sha256:"), "stop_resolution", r); err != nil {
			return err
		}
		s.Result.Delivery.ActiveStop = nil
		// An owner still needs actual acceptance if it was never recorded.
		if s.StartSHA256 == nil {
			s.Result.Delivery.Phase = "awaiting_acceptance"
			s.Result.Round.State = "awaiting_acceptance"
			s.Result.NextAction = WorkflowAction{"recipient", "Use execute begin with actual observations before Task writes."}
		} else if deliveryStarted(*s) {
			s.Result.Delivery.Phase = "working"
			s.Result.Round.State = "working"
			s.Result.NextAction = WorkflowAction{"recipient", "Continue the same accepted mandate after the preserved cause repair."}
		}
		return nil
	})
	return deliveryReadback(d, root, id, err)
}

func deliveryObservedStop(s workflowState) bool {
	return s.Result.Delivery != nil && s.Result.Delivery.ActiveStop != nil
}

func deliveryLoadObservedStop(s *workflowState) error {
	if s.Result.Delivery == nil {
		return nil
	}
	if err := deliveryRecoverStopEvents(s); err != nil {
		return err
	}
	s.Result.Delivery.ActiveStop = nil
	var pending []FileBinding
	for _, event := range s.Result.Delivery.Events {
		switch event.Kind {
		case "coordinator_stop":
			b, err := workflowBound(event.Binding, 2<<20)
			if err != nil {
				return err
			}
			var stop DeliveryStopRecord
			if err = decode(b, 2<<20, &stop); err != nil {
				return err
			}
			if stop.Kind != "ply.workflow.coordinator-stop" || stop.SchemaVersion != 1 || stop.RunID != s.Result.RunID || stop.Input.EventID != event.ID {
				return workflowError(4, "coordinator stop binding changed")
			}
			if _, err = workflowBound(stop.Evidence, 1<<20); err != nil {
				return err
			}
			binding := event.Binding
			pending = append(pending, binding)
			s.Result.Delivery.ActiveStop = &binding
		case "stop_resolution":
			b, err := workflowBound(event.Binding, 1<<20)
			if err != nil {
				return err
			}
			var resolution DeliveryStopResolution
			if err = decode(b, 1<<20, &resolution); err != nil {
				return err
			}
			if resolution.Kind != "ply.workflow.stop-resolution" || resolution.RunID != s.Result.RunID || !resolution.PermissionConfirmed || s.Result.Delivery.ActiveStop == nil || resolution.Stop != *s.Result.Delivery.ActiveStop {
				return workflowError(4, "stop resolution binding changed")
			}
			pending = pending[:len(pending)-1]
			s.Result.Delivery.ActiveStop = nil
			if len(pending) > 0 {
				binding := pending[len(pending)-1]
				s.Result.Delivery.ActiveStop = &binding
			}
		}
	}
	return nil
}

// Event bytes are durable before state.json indexes them. Recover only trailing
// coordinator stops into the read projection, so a failed state publication
// cannot hide a denial. The next locked update persists the recovered index;
// identical observe-stop retries then reuse the original timestamp and bytes.
// Other event kinds retain their existing operation-specific recovery paths.
func deliveryRecoverStopEvents(s *workflowState) error {
	for {
		path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "events", fmt.Sprintf("%08d.json", len(s.Result.Delivery.Events)+1))
		raw, err := readFile(path, 2<<20, true)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		var header struct {
			Kind string `json:"kind"`
		}
		if err = json.Unmarshal(raw, &header); err != nil {
			return err
		}
		if header.Kind != "ply.workflow.coordinator-stop" {
			return nil
		}
		var stop DeliveryStopRecord
		if err = decode(raw, 2<<20, &stop); err != nil {
			return err
		}
		for _, prior := range s.Result.Delivery.Events {
			if prior.ID == stop.Input.EventID {
				return workflowError(4, "unindexed stop reuses a recorded event ID")
			}
		}
		binding := FileBinding{path, hash(raw)}
		s.Result.Delivery.Events = append(s.Result.Delivery.Events, DeliveryEventRecord{stop.Input.EventID, "coordinator_stop", binding})
		s.Result.Delivery.LastEventSHA256 = &binding.SHA256
		// deliveryLoadObservedStop validates the full recovered record and its
		// preserved evidence before exposing this state to any caller.
	}
}

func deliveryProjectStop(s workflowState, o *WorkflowRun) {
	if !deliveryObservedStop(s) {
		return
	}
	projected := *o.Delivery
	projected.Phase = "stopped"
	o.Delivery = &projected
	o.Round.State = "stopped"
	o.NextAction = WorkflowAction{"coordinator", "A sourced provider stop is preserved at " + s.Result.Delivery.ActiveStop.Locator + ". Inspect and repair its cause within the existing mandate. The same owner must use execute resolve-stop with this observation digest and actual resolution before continuing; begin is still required if acceptance is absent. Do not resend denied input, infer permission or release the writer."}
}
