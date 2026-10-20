package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// ReportStopResolutionInput describes a coordinator's actual action-result or
// effect-inspection evidence. It never resumes the owner or confirms permission.
type ReportStopResolutionInput struct {
	StopSHA256   string `json:"stop_sha256"`
	SourceSHA256 string `json:"source_sha256"`
	EffectStatus string `json:"effect_status"`
	Basis        string `json:"basis"`
	Observer     string `json:"observer"`
	Reason       string `json:"reason"`
	EvidencePath string `json:"evidence_path"`
}

type reportStopResolution struct {
	Kind           string                    `json:"kind"`
	SchemaVersion  int                       `json:"schema_version"`
	RunID          string                    `json:"run_id"`
	RequestSHA256  string                    `json:"request_sha256"`
	NativeSession  string                    `json:"native_session"`
	Stop           FileBinding               `json:"stop"`
	Input          ReportStopResolutionInput `json:"input"`
	Evidence       FileBinding               `json:"evidence"`
	EvidenceSource FileBinding               `json:"evidence_source"`
	RecordedAtUTC  string                    `json:"recorded_at_utc"`
}

func reportResolutionPath(s workflowState, stop FileBinding) string {
	return filepath.Join(s.Result.Paths.RunRoot, "report-stop-resolutions", stop.SHA256[7:]+".json")
}

func reportOriginalEffect(s workflowState, stop FileBinding) (string, error) {
	raw, err := workflowBound(stop, 2<<20)
	if err != nil {
		return "", err
	}
	if s.ReportProviderStop != nil && stop == *s.ReportProviderStop {
		var r DeliveryStopRecord
		if err = decode(raw, 2<<20, &r); err != nil {
			return "", err
		}
		return r.Input.EffectStatus, nil
	}
	if s.ReportStop == nil || stop != *s.ReportStop {
		return "", workflowError(4, "resolution must bind an actual preserved report stop")
	}
	var r WorkflowReport
	if err = decode(raw, 2<<20, &r); err != nil {
		return "", err
	}
	var forbidden []json.RawMessage
	json.Unmarshal(r.ForbiddenEffectsObserved, &forbidden)
	if len(forbidden) > 0 {
		return DeliveryEffectPartial, nil
	}
	var reasons []struct {
		Type string `json:"type"`
	}
	json.Unmarshal(r.StopReasons, &reasons)
	for _, reason := range reasons {
		if reason.Type == "unknown_or_partial_effect" {
			return DeliveryEffectUnknown, nil
		}
	}
	return DeliveryEffectNoneObserved, nil
}

func readReportStopResolution(s workflowState, stop FileBinding) (*reportStopResolution, *FileBinding, error) {
	path := reportResolutionPath(s, stop)
	raw, err := readFile(path, 2<<20, true)
	if os.IsNotExist(err) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	var r reportStopResolution
	if err = decode(raw, 2<<20, &r); err != nil {
		return nil, nil, err
	}
	if r.Kind != "ply.workflow.report-stop-resolution" || r.SchemaVersion != 1 || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.NativeSession != s.Result.Transport.AgentSessionID || r.Stop != stop || r.Input.StopSHA256 != stop.SHA256 {
		return nil, nil, workflowError(4, "report stop resolution binding differs")
	}
	if _, err = workflowBound(r.Evidence, 1<<20); err != nil {
		return nil, nil, err
	}
	binding := FileBinding{path, hash(raw)}
	return &r, &binding, nil
}

// WorkflowResolveReportStop appends one immutable, source-bound observation to
// an unknown stop. The original refusal remains active for all owner callbacks.
// A no-effect observation can qualify only retirement, never receipt or work.
func WorkflowResolveReportStop(d Dependencies, root, id string, in ReportStopResolutionInput) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	if !digestPattern.MatchString(in.StopSHA256) || !digestPattern.MatchString(in.SourceSHA256) || !plain(in.Observer, 1, 256) || !plain(in.Reason, 1, 1900) {
		return WorkflowRun{}, workflowError(2, "exact stop/source digests, actual observer and bounded reason are required")
	}
	if in.Basis != "action_result" && in.Basis != "effect_inspection" {
		return WorkflowRun{}, workflowError(4, "effect proof must come from the actual action result or effect inspection; writer absence alone cannot establish past effects")
	}
	if in.EffectStatus != DeliveryEffectNotStarted && in.EffectStatus != DeliveryEffectNoneObserved && in.EffectStatus != DeliveryEffectPartial {
		return WorkflowRun{}, workflowError(4, "resolution must preserve the actually observed effect status; unknown effects remain unresolved")
	}
	source, raw, err := deliveryReadBinding(in.EvidencePath, 1<<20)
	if err != nil {
		return WorkflowRun{}, err
	}
	err = workflowUpdate(d, root, id, func(s *workflowState) error {
		if deliveryRun(s.Request) || s.Request.ReturnMode != "reviewed_report_only" {
			return workflowError(4, "effect reconciliation requires a report owner")
		}
		if err := workflowFreshArtifacts(d, *s, false, false); err != nil {
			return err
		}
		var stop *FileBinding
		for _, binding := range []*FileBinding{s.ReportStop, s.ReportProviderStop} {
			if binding != nil && binding.SHA256 == in.StopSHA256 {
				stop = binding
			}
		}
		if stop == nil {
			return workflowError(4, "stop digest does not identify this Run's preserved observation")
		}
		prior, _, err := readReportStopResolution(*s, *stop)
		if err != nil {
			return err
		}
		if prior != nil {
			if prior.Input != in || prior.EvidenceSource != source {
				return workflowError(4, "report stop resolution is immutable; a conflicting claim cannot replace preserved effects")
			}
			return nil
		}
		if s.ReportRetirement != nil {
			return workflowError(4, "owner_retired: only identical reconciliation readback is allowed")
		}
		effect, err := reportOriginalEffect(*s, *stop)
		if err != nil {
			return err
		}
		if effect == DeliveryEffectPartial {
			return workflowError(4, "actual partial effects remain blocked; absence or a later no-effect claim cannot erase them")
		}
		if effect != "" && effect != DeliveryEffectUnknown {
			return workflowError(4, "original stop already records known effects; use its existing retirement path")
		}
		target, err := workflowTarget(d, *s)
		if err != nil {
			return err
		}
		if target != in.SourceSHA256 {
			return workflowError(4, "effect observation source differs from the exact current Task; inspect source and evidence before reconciliation")
		}
		evidencePath := filepath.Join(s.Result.Paths.RunRoot, "report-stop-resolutions", "evidence", source.SHA256[7:])
		if err = d.writeOnce(evidencePath, raw); err != nil {
			return err
		}
		after, err := workflowTarget(d, *s)
		if err != nil {
			return err
		}
		if after != target {
			return workflowError(4, "Task source changed during effect reconciliation")
		}
		record := reportStopResolution{Kind: "ply.workflow.report-stop-resolution", SchemaVersion: 1, RunID: id, RequestSHA256: s.Result.RequestSHA256, NativeSession: s.Result.Transport.AgentSessionID, Stop: *stop, Input: in, Evidence: FileBinding{evidencePath, source.SHA256}, EvidenceSource: source, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		if _, err = workflowKeep(d, reportResolutionPath(*s, *stop), record); err != nil {
			return err
		}
		s.Result.Round.State = "stopped"
		s.Result.NextAction = workflowReportStopAction(*s)
		return d.fault("report_after_stop_resolution")
	})
	return deliveryReadback(d, root, id, err)
}

func workflowReportEffectsPermitRetirement(d Dependencies, s workflowState) error {
	for _, stop := range []*FileBinding{s.ReportStop, s.ReportProviderStop} {
		if stop == nil {
			continue
		}
		effect, err := reportOriginalEffect(s, *stop)
		if err != nil {
			return err
		}
		if effect == DeliveryEffectPartial {
			return workflowError(4, "partial_effect: actual report effects remain blocked; preserve and reconcile those effects without retirement")
		}
		if effect == DeliveryEffectNoneObserved || effect == DeliveryEffectNotStarted {
			continue
		}
		resolution, _, err := readReportStopResolution(s, *stop)
		if err != nil {
			return err
		}
		if resolution == nil {
			return workflowError(4, "unknown_effect: preserve the refusal; use workflow run resolve-stop with actual effect evidence and exact source before retirement")
		}
		if resolution.Input.EffectStatus != DeliveryEffectNotStarted && resolution.Input.EffectStatus != DeliveryEffectNoneObserved {
			return workflowError(4, "partial_effect: the preserved reconciliation still records partial effects")
		}
		source, err := workflowTarget(d, s)
		if err != nil {
			return err
		}
		if source != resolution.Input.SourceSHA256 {
			return workflowError(4, "effect_source_changed: reconciliation no longer describes the exact Task source; no retirement or queue release")
		}
	}
	return nil
}
