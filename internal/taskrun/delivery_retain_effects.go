package taskrun

import (
	"fmt"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/sirupsen/logrus"
)

const deliveryRetainEffects = "retain_effects"
const retainedResourcesCaution = "Resources retired; original effects retained, including uncertainty. No delivery or replay. Fresh Task mandate required."

// DeliveryRetainedEffect preserves the original claim and its immutable evidence;
// it is deliberately not an effect resolution.
type DeliveryRetainedEffect struct {
	Stop           FileBinding `json:"stop"`
	OriginalEffect string      `json:"original_effect"`
	Action         string      `json:"action"`
	Evidence       FileBinding `json:"evidence"`
}

type RetainedDeliveryObservation struct {
	Kind           string                        `json:"kind"`
	SchemaVersion  int                           `json:"schema_version"`
	RunID          string                        `json:"run_id"`
	RequestSHA256  string                        `json:"request_sha256"`
	NativeSession  string                        `json:"native_session"`
	Provider       string                        `json:"provider"`
	PreparationID  string                        `json:"preparation_id"`
	Retirement     FileBinding                   `json:"retirement"`
	Previous       FileBinding                   `json:"previous"`
	Source         ReportRetainedSource          `json:"source"`
	Effects        []DeliveryRetainedEffect      `json:"retained_effects"`
	ActiveStop     *FileBinding                  `json:"active_stop"`
	Absence        HumanOwnerRecoveryObservation `json:"absence"`
	WriterLiveness *TaskWriterLiveness           `json:"writer_liveness"`
	RecordedAtUTC  string                        `json:"recorded_at_utc"`
}

// These checks deliberately treat launchers as historical. This coordinator
// neither invokes a provider nor authenticates on behalf of the absent owner.
func deliveryRetentionGuard(d Dependencies, s workflowState) error {
	if !deliveryRun(s.Request) || headlessRun(s.Request) || s.Result.Delivery == nil || s.Request.Delivery.Agreement == nil || s.Request.Delivery.Agreement.Mode != workspace.DeliveryLocalEpic {
		return workflowError(4, "retain-effects requires a started stopped interactive local Epic delivery")
	}
	cwd, err := d.CWD()
	if err != nil {
		return err
	}
	if sourceContains(s.Observed.Target.WorktreeLocator, cwd) {
		return workflowError(4, "retain-effects must run outside the Task source; use the containing workspace")
	}
	if err := deliveryAcceptedAuthority(s); err != nil {
		return err
	}
	if err := workflowFreshArtifacts(d, s, false, false); err != nil {
		return err
	}
	facts, err := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, s.Result.Handoff.Locator)
	if err != nil {
		return err
	}
	if facts.Start == nil || s.StartSHA256 == nil || facts.Start.SHA256 != *s.StartSHA256 || s.Result.Transport.AgentSessionID == "" {
		return workflowError(4, "retain-effects requires the exact native start and session")
	}
	dl := s.Result.Delivery
	if facts.Terminal != nil || s.Result.FinalReturn.TerminalSHA256 != nil || dl.Phase == "completed" || dl.OwnershipRelease != nil || dl.HumanIntegration != nil {
		return workflowError(4, "terminal or released ownership cannot be retired")
	}
	if err = deliveryQuotaAttemptGuard(s); err != nil {
		return err
	}
	if dl.ActiveStop == nil {
		return workflowError(4, "retain-effects requires a preserved active coordinator stop")
	}
	p, err := workflowReportPreparation(d, s)
	if err != nil {
		return err
	}
	a := s.Request.Delivery.Agreement
	t := p.Plan.Target
	if a.ProjectID != t.ProjectID || a.RepoID != t.RepoID || a.EpicID != t.EpicID || a.TargetRef != t.ParentRef || a.TargetWorktree != t.ParentLocator || t.GitCommonDir != s.Observed.Target.GitCommonDir {
		return workflowError(4, "retained preparation target differs from frozen delivery")
	}
	return nil
}

func deliveryRetainedEffects(s workflowState) ([]DeliveryRetainedEffect, error) {
	// Resolutions pop only their latest pending stop; earlier uncertainty survives.
	type pendingStop struct {
		binding FileBinding
		record  DeliveryStopRecord
	}
	pending := []pendingStop{}
	for _, e := range s.Result.Delivery.Events {
		switch e.Kind {
		case "coordinator_stop":
			raw, err := workflowBound(e.Binding, 2<<20)
			if err != nil {
				return nil, err
			}
			var stop DeliveryStopRecord
			if err = decode(raw, 2<<20, &stop); err != nil {
				return nil, err
			}
			if stop.Kind != "ply.workflow.coordinator-stop" || stop.SchemaVersion != 1 || stop.RunID != s.Result.RunID || stop.Input.EventID != e.ID {
				return nil, workflowError(4, "retained stop identity binding differs")
			}
			if _, err = workflowBound(stop.Evidence, 1<<20); err != nil {
				return nil, err
			}
			if stop.Source.SHA256 != stop.Evidence.SHA256 {
				return nil, workflowError(4, "retained stop evidence differs")
			}
			pending = append(pending, pendingStop{e.Binding, stop})
		case "stop_resolution":
			if len(pending) == 0 {
				return nil, workflowError(4, "retained stop resolution has no pending stop")
			}
			pending = pending[:len(pending)-1]
		}
	}
	out := []DeliveryRetainedEffect{}
	for _, p := range pending {
		stop := p.record
		if stop.Input.Provider != s.Request.Runtime.Provider || stop.Input.NativeSession != s.Result.Transport.AgentSessionID {
			return nil, workflowError(4, "retained stop provider/session binding differs")
		}
		out = append(out, DeliveryRetainedEffect{p.binding, stop.Input.EffectStatus, stop.Input.Action, stop.Evidence})
	}
	if len(out) == 0 || s.Result.Delivery.ActiveStop == nil || out[len(out)-1].Stop != *s.Result.Delivery.ActiveStop {
		return nil, workflowError(4, "retained active stop differs")
	}
	return out, nil
}

func deliveryRetentionObserve(d Dependencies, s workflowState) (ReportRetainedSource, []DeliveryRetainedEffect, taskWriterEvidence, error) {
	var empty taskWriterEvidence
	if err := deliveryRetentionGuard(d, s); err != nil {
		return ReportRetainedSource{}, nil, empty, err
	}
	effects, err := deliveryRetainedEffects(s)
	if err != nil {
		return ReportRetainedSource{}, nil, empty, err
	}
	source, err := workflowRetainedSource(d, s, true)
	if err != nil {
		return source, nil, empty, err
	}
	writers := (taskWriters{d, s}).cleanup()
	if writers.state != taskWritersAbsent {
		return source, nil, writers, writers.err
	}
	after, err := workflowRetainedSource(d, s, true)
	if err != nil {
		return source, nil, writers, err
	}
	if !equal(source, after) {
		return source, nil, writers, workflowError(4, "source changed during retained resource observation; retry after local changes settle")
	}
	if err = deliveryRetentionGuard(d, s); err != nil {
		return source, nil, writers, err
	}
	current, err := deliveryRetainedEffects(s)
	if err != nil {
		return source, nil, writers, err
	}
	if !equal(effects, current) {
		return source, nil, writers, workflowError(4, "stop effects changed during retained resource observation")
	}
	return source, effects, writers, nil
}

func deliveryRetentionWarn(s workflowState, b FileBinding, source ReportRetainedSource, effects []DeliveryRetainedEffect) {
	uncertain := []string{}
	for _, e := range effects {
		if e.OriginalEffect == DeliveryEffectPartial || e.OriginalEffect == DeliveryEffectUnknown {
			uncertain = append(uncertain, e.OriginalEffect)
		}
	}
	if len(uncertain) > 0 {
		logrus.WithFields(logrus.Fields{"run": s.Result.RunID, "original_effect": strings.Join(uncertain, ","), "retained_record": b.Locator, "source": source.SourceSHA256}).Warn(retainedResourcesCaution)
	}
}

func deliveryRetentionHistory(s workflowState) (*DeliveryRetirement, FileBinding, []RetainedDeliveryObservation, []FileBinding, error) {
	r, err := deliveryRetirement(s)
	if err != nil {
		return r, FileBinding{}, nil, nil, err
	}
	if r == nil {
		if s.Result.Delivery != nil {
			for _, event := range s.Result.Delivery.Events {
				if event.Kind == "retained_resource_observation" {
					return nil, FileBinding{}, nil, nil, workflowError(4, "retained observation has no retirement")
				}
			}
		}
		return nil, FileBinding{}, nil, nil, nil
	}
	var retirement FileBinding
	for _, e := range s.Result.Delivery.Events {
		if e.Kind == "owner_retirement" {
			retirement = e.Binding
			break
		}
	}
	observations := []RetainedDeliveryObservation{}
	bindings := []FileBinding{}
	prior := retirement
	for _, e := range s.Result.Delivery.Events {
		if e.Kind != "retained_resource_observation" {
			continue
		}
		raw, err := workflowBound(e.Binding, 2<<20)
		if err != nil {
			return nil, retirement, nil, nil, err
		}
		var o RetainedDeliveryObservation
		if err = decode(raw, 2<<20, &o); err != nil {
			return nil, retirement, nil, nil, err
		}
		if r.Mode != deliveryRetainEffects || o.Kind != "ply.workflow.retained-resource-observation" || o.SchemaVersion != 1 || o.RunID != s.Result.RunID || o.RequestSHA256 != s.Result.RequestSHA256 || o.NativeSession != r.NativeSession || o.Provider != r.Provider || o.PreparationID != r.PreparationID || o.Retirement != retirement || o.Previous != prior || e.ID != "retained-observation-"+fmt.Sprint(len(observations)+1) {
			return nil, retirement, nil, nil, workflowError(4, "retained observation chain binding differs")
		}
		if o.Source.Worktree != r.DeliverySource.Worktree || o.Source.Ref != r.DeliverySource.Ref || o.Source.GitCommonDir != r.DeliverySource.GitCommonDir {
			return nil, retirement, nil, nil, workflowError(4, "retained observation source identity differs")
		}
		observations = append(observations, o)
		bindings = append(bindings, e.Binding)
		prior = e.Binding
	}
	return r, retirement, observations, bindings, nil
}

func workflowRetireRetainingEffects(d Dependencies, root, id string, in DeliveryRetireInput) (WorkflowRun, error) {
	if err := containing(d, root); err != nil {
		return WorkflowRun{}, err
	}
	if !plain(in.Observer, 1, 256) || strings.TrimSpace(in.Reason) == "" || len(in.Reason) > 1900 {
		return WorkflowRun{}, workflowError(2, "actual observer and bounded reason are required")
	}
	err := workflowUpdate(d, root, id, func(s *workflowState) error {
		r, _, _, _, err := deliveryRetentionHistory(*s)
		if err != nil {
			return err
		}
		if r != nil {
			if r.Mode != deliveryRetainEffects || r.Input != in {
				return workflowError(4, "retirement input/mode is immutable")
			}
			s.Result.Delivery.Phase, s.Result.Round.State = "stopped", "stopped"
			return nil
		}
		source, effects, writers, err := deliveryRetentionObserve(d, *s)
		if err != nil {
			return err
		}
		r = &DeliveryRetirement{Kind: "ply.workflow.owner-retirement", SchemaVersion: 1, RunID: id, RequestSHA256: s.Result.RequestSHA256, NativeSession: s.Result.Transport.AgentSessionID, Input: in, Mode: deliveryRetainEffects, Provider: s.Request.Runtime.Provider, PreparationID: s.Request.PreparationID, DeliverySource: &source, DeliveryEffects: effects, ActiveStop: s.Result.Delivery.ActiveStop, Absence: writers.absence, WriterLiveness: writers.liveness, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		b, err := deliveryAppendEvent(d, s, "owner-retirement", "owner_retirement", r)
		if err != nil {
			return err
		}
		deliveryRetentionWarn(*s, b, source, effects)
		if err = d.fault("retain_effects_after_retirement"); err != nil {
			return err
		}
		s.Result.Delivery.Phase, s.Result.Round.State = "stopped", "stopped"
		return nil
	})
	if err == nil {
		err = deliveryRetainedRelease(d, root, id, in)
	}
	return deliveryReadback(d, root, id, err)
}

func deliveryRetainedRelease(d Dependencies, root, id string, in DeliveryRetireInput) error {
	return workflowUpdate(d, root, id, func(s *workflowState) error {
		r, retirement, observations, bindings, err := deliveryRetentionHistory(*s)
		if err != nil {
			return err
		}
		if r == nil || r.Mode != deliveryRetainEffects || r.Input != in {
			return workflowError(4, "retained retirement binding differs")
		}
		p, err := workflowReportPreparation(d, *s)
		if err != nil {
			return err
		}
		if p.Disposition == "released" || p.Disposition == "advanced" {
			return nil
		}
		source, effects, writers, err := deliveryRetentionObserve(d, *s)
		if err != nil {
			return err
		}
		oldSource, oldEffects, previous := r.DeliverySource, r.DeliveryEffects, retirement
		if len(observations) > 0 {
			o := observations[len(observations)-1]
			oldSource, oldEffects, previous = &o.Source, o.Effects, bindings[len(bindings)-1]
		}
		if !equal(oldSource, &source) || !equal(oldEffects, effects) {
			o := RetainedDeliveryObservation{Kind: "ply.workflow.retained-resource-observation", SchemaVersion: 1, RunID: id, RequestSHA256: s.Result.RequestSHA256, NativeSession: r.NativeSession, Provider: r.Provider, PreparationID: r.PreparationID, Retirement: retirement, Previous: previous, Source: source, Effects: effects, ActiveStop: s.Result.Delivery.ActiveStop, Absence: writers.absence, WriterLiveness: writers.liveness, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
			b, err := deliveryAppendEvent(d, s, fmt.Sprintf("retained-observation-%d", len(observations)+1), "retained_resource_observation", o)
			if err != nil {
				return err
			}
			deliveryRetentionWarn(*s, b, source, effects)
			if err = d.fault("retain_effects_after_observation"); err != nil {
				return err
			}
		}
		if err = deliveryRetentionGuard(d, *s); err != nil {
			return err
		}
		current, err := deliveryRetainedEffects(*s)
		if err != nil {
			return err
		}
		if !equal(current, effects) {
			return workflowError(4, "retained stop bindings changed before queue release")
		}
		if err = d.fault("retain_effects_before_release"); err != nil {
			return err
		}
		target := p.Plan.Target
		qtarget := workspace.QueueTargetInput{ProjectID: target.ProjectID, RepoID: target.RepoID, EpicID: target.EpicID}
		q, err := workspace.ListTaskQueue(d.Workspace, qtarget, false)
		if err != nil {
			return err
		}
		if q.Current == nil || q.Current.PreparationID != s.Request.PreparationID {
			return workflowError(4, "retired owner no longer owns exact current preparation; inspect queue without advancing a successor")
		}
		_, err = workspace.CloseTaskQueue(d.Workspace, qtarget, s.Request.PreparationID, q.Revision, "Resources retired; source/effects retained", "release")
		if err != nil {
			return err
		}
		return d.fault("retain_effects_after_release")
	})
}
