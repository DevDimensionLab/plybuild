package taskrun

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type HumanIntegrationOwnership struct {
	RunID         string                 `json:"run_id"`
	TaskResultID  workspace.TaskResultID `json:"task_result_id"`
	CandidateOID  string                 `json:"candidate_oid"`
	CandidateTree string                 `json:"candidate_tree"`
	Released      bool                   `json:"released"`
	Release       *FileBinding           `json:"release"`
	Reason        string                 `json:"reason"`
	NextAction    string                 `json:"next_action"`
}

type HumanIntegrationInput struct {
	PlanID      string                                   `json:"plan_id"`
	PlanSHA256  string                                   `json:"plan_sha256"`
	Decision    FileBinding                              `json:"decision"`
	Attestation workflowhandoff.DeliveryHumanAttestation `json:"attestation"`
}

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

func validateHumanIntegrationDecision(in HumanIntegrationInput) error {
	raw, err := workflowBound(in.Decision, 1<<20)
	if err != nil {
		return err
	}
	var decision struct {
		Kind          string                                   `json:"kind"`
		SchemaVersion int                                      `json:"schema_version"`
		IntegrationID string                                   `json:"integration_id"`
		PlanSHA256    string                                   `json:"plan_sha256"`
		Human         workflowhandoff.DeliveryHumanAttestation `json:"human"`
	}
	if err = decode(raw, 1<<20, &decision); err != nil {
		return err
	}
	if decision.Kind != "ply.integration.decision" || decision.SchemaVersion != 1 || decision.IntegrationID != in.PlanID || decision.PlanSHA256 != in.PlanSHA256 || !equal(decision.Human, in.Attestation) {
		return workflowError(4, "plan_changed: preserved human decision does not bind the selected plan and exact answer")
	}
	return nil
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
	if r.Origin != "bound_owner_release" && r.Origin != "observed_provider_exit" && r.Origin != "automatic_delivery_completed" || (r.Origin == "bound_owner_release" || r.Origin == "automatic_delivery_completed") && (r.InvokingExecutable == nil || r.ProviderExit != nil) || r.Origin == "observed_provider_exit" && (r.ProviderExit == nil || r.ProviderExit.PaneID != s.Result.Transport.PaneID || r.ProviderExit.TerminalID != s.Result.Transport.TerminalID) {
		return nil, workflowError(4, "owner_release_invalid: release provenance is incomplete")
	}
	if r.Origin == "automatic_delivery_completed" {
		if _, observed, err := automaticIntegratedResult(s, c.TaskResult.ID); err != nil || !observed {
			return nil, workflowError(4, "owner_release_invalid: automatic release lacks the exact observed local integration")
		}
		if a := s.Result.Delivery.Attempt; a != nil && a.State == "attempted" {
			return nil, workflowError(4, "owner_release_invalid: automatic delivery effect remains unresolved")
		}
	}
	for _, event := range s.Result.Delivery.Events {
		if event.Kind == "ownership_release" && equal(event.Binding, b) {
			return &b, nil
		}
	}
	return nil, workflowError(4, "owner_release_invalid: release is absent from native history")
}

// PreviewHumanIntegration is read-only, including when the owner is active.
func PreviewHumanIntegration(d Dependencies, root, id string, resultID workspace.TaskResultID) (HumanIntegrationOwnership, error) {
	var out HumanIntegrationOwnership
	a, err := nativeDeliveryCandidate(d, root, id, resultID)
	if err != nil {
		return out, err
	}
	s, err := workflowRead(root, id)
	if err != nil {
		return out, err
	}
	out = HumanIntegrationOwnership{RunID: id, TaskResultID: resultID, CandidateOID: a.Candidate.OID, CandidateTree: a.Candidate.Tree, Reason: "owner_active", NextAction: "The delivery owner must explicitly release the exact candidate with ply workflow execute release before human integration."}
	out.Release, err = deliveryRelease(s, a.Candidate)
	if err != nil {
		return out, err
	}
	out.Released = out.Release != nil
	if out.Released {
		out.Reason, out.NextAction = "owner_released", "Review the exact integration plan and preserve an actual human pass before effects."
	} else if _, e := deliveryRecoveryExit(d, s); e == nil {
		out.Reason, out.NextAction = "owner_exited_release_required", "The exact provider process has ended. Run ply integration release --delivery <selected-delivery-id> to preserve its release, then review the integration plan."
	}
	return out, nil
}

// ReleaseDeliveryOwnership is an explicit callback from the bound owner. It
// revokes further owner callbacks but does not terminate or guess process state.
// A later negative human answer explicitly returns correction ownership.
func ReleaseDeliveryOwnership(d Dependencies, root, id, contextPath, reason string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	if !plain(reason, 1, 2000) {
		return WorkflowRun{}, workflowError(2, "ownership release requires an explicit reason")
	}
	err := deliveryReleaseUpdate(d, root, id, func(s *workflowState) error {
		if !deliveryRun(s.Request) || s.Result.Delivery == nil || len(s.Result.Delivery.Candidates) == 0 {
			return workflowError(4, "ownership release requires an exact qualified delivery candidate")
		}
		// Revocation is deliberately callable from a newer installed CLI: the
		// immutable controller may predate release. It grants no target effect.
		if contextPath != s.Result.Paths.Context || s.ContextSHA256 == "" {
			return workflowError(4, "release requires the exact current private owner context")
		}
		if _, e := workflowBound(FileBinding{contextPath, s.ContextSHA256}, 1<<20); e != nil {
			return e
		}
		cwd, e := d.CWD()
		if e != nil {
			return e
		}
		if cwd != s.Observed.Target.WorktreeLocator {
			return workflowError(4, "release cwd differs from the exact Task worktree")
		}
		c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		if prior, e := deliveryRelease(*s, c); e != nil || prior != nil {
			return e
		}
		if _, e := workflowAgentGet(d, *s, false); e != nil {
			return e
		}
		if _, e := nativeDeliveryCandidate(d, root, id, c.TaskResult.ID); e != nil {
			return e
		}
		if a := s.Result.Delivery.Attempt; a != nil && a.State == "attempted" {
			return workflowError(4, "owner_active: an unresolved delivery attempt must be observed before release")
		}
		path, e := d.Executable()
		if e != nil {
			return e
		}
		path, e = filepath.EvalSymlinks(path)
		if e != nil {
			return e
		}
		bytes, e := readFile(path, 128<<20, false)
		if e != nil {
			return e
		}
		return preserveDeliveryRelease(d, s, c, reason, "bound_owner_release", &Executable{Path: path, SHA256: hash(bytes)}, nil)
	})
	return deliveryReadback(d, root, id, err)
}

func preserveDeliveryRelease(d Dependencies, s *workflowState, c DeliveryCandidate, reason, origin string, invoker *Executable, exit *DeliveryStartupExitEvidence) error {
	r := deliveryOwnershipRelease{Kind: "PlyDeliveryOwnerRelease@1", SchemaVersion: 1, RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ResultOID: c.OID, ResultTree: c.Tree, SourceRef: c.TaskResult.SourceRef, SourceLocator: c.TaskResult.SourceLocator, OwnerClaim: s.Request.Delivery.OwnerClaim, SessionID: s.Result.SessionID, Reason: reason, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano), Origin: origin, InvokingExecutable: invoker, ProviderExit: exit}
	b, e := deliveryAppendEvent(d, s, fmt.Sprintf("release-%s-%08d", c.Key, len(s.Result.Delivery.Events)+1), "ownership_release", r)
	if e != nil {
		return e
	}
	s.Result.Delivery.OwnershipRelease = &b
	s.Result.NextAction = WorkflowAction{"user", "The owner released the exact source candidate. Start ply integration from an existing return checkout and confirm its visible plan."}
	return nil
}

// Revocation uses the TaskRun -> workspace lock order. Native integration
// callers outside the delivery adapter also finish (or preserve uncertainty)
// before the release is published. New effects then see the released authority.
func deliveryReleaseUpdate(d Dependencies, root, id string, fn func(*workflowState) error) error {
	return withStore(root, func() error {
		s, e := workflowRead(root, id)
		if e != nil {
			return e
		}
		if s.Result.Delivery == nil || len(s.Result.Delivery.Candidates) == 0 {
			return workflowError(4, "ownership release requires an exact qualified delivery candidate")
		}
		c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		if a := s.Result.Delivery.Attempt; a != nil && a.State == "attempted" {
			return workflowError(4, "owner_active: a reserved delivery effect is unresolved")
		}
		// Publication includes push, PR creation, metadata and native queue
		// completion. A technical candidate alone never releases that writer.
		if a := s.Request.Delivery.Agreement; a != nil && a.Mode == workspace.DeliveryPullRequest && (s.Result.Delivery.Phase != "completed" || c.PullRequest == nil) {
			return workflowError(4, "owner_active: finish the exact PR publication before releasing its source")
		}
		if d.Workspace.WorkItems == nil {
			return workflowError(4, "native integration ownership cannot be observed")
		}
		return d.Workspace.WorkItems.WithLock(root, func(session workspace.WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			for _, authority := range r.IntegrationAuthorities {
				if authority.TaskID != c.TaskResult.TaskID || authority.TaskResultID != c.TaskResult.ID {
					continue
				}
				settled := false
				for _, result := range r.IntegrationResults {
					if result.AuthorityID == authority.ID {
						settled = result.Outcome != "partial" && result.Outcome != "unknown"
					}
				}
				if !settled {
					return workflowError(4, "owner_active: the native integration effect is unresolved")
				}
			}
			if e = fn(&s); e != nil {
				return e
			}
			return workflowSave(d, s)
		})
	})
}

// ReleaseExitedDeliveryOwnership needs no private callback context. It
// positively observes the exact terminal's sole shell and complete process
// tree; missing UI state or an idle agent does not establish process exit.
func ReleaseExitedDeliveryOwnership(d Dependencies, root, id string, resultID workspace.TaskResultID) (WorkflowRun, error) {
	err := deliveryReleaseUpdate(d, root, id, func(s *workflowState) error {
		a, e := nativeDeliveryCandidate(d, root, id, resultID)
		if e != nil {
			return e
		}
		if prior, e := deliveryRelease(*s, a.Candidate); e != nil || prior != nil {
			return e
		}
		if attempt := s.Result.Delivery.Attempt; attempt != nil && attempt.State == "attempted" {
			return workflowError(4, "owner_active: resolve the reserved delivery effect before release")
		}
		exit, e := deliveryRecoveryExit(d, *s)
		if e != nil {
			return workflowError(4, "owner_active: exact provider exit is unproved; the active owner must explicitly release, or resolve its preserved terminal and process state")
		}
		return preserveDeliveryRelease(d, s, a.Candidate, "The exact provider has ended; human-started ownership release preserves its native process observation.", "observed_provider_exit", nil, &exit)
	})
	return deliveryReadback(d, root, id, err)
}

// CheckAgentDeliveryExecution protects delivery execute and execute integrate.
// Human pass does not expand the developer's frozen completion boundary.
func CheckAgentDeliveryExecution(d Dependencies, root, id string) error {
	s, e := workflowRead(root, id)
	if e != nil {
		return e
	}
	if !deliveryRun(s.Request) || s.Result.Delivery == nil {
		return workflowError(4, "native delivery owner is missing")
	}
	if s.Result.Delivery.OwnershipRelease != nil || s.Result.Delivery.HumanIntegration != nil && s.Result.Delivery.HumanIntegration.Accepted {
		return workflowError(4, "human_integration_required: source ownership was released; continue the preserved human integration plan")
	}
	if a := s.Request.Delivery.Agreement; a != nil && a.HumanOwnedIntegration() && a.Mode != workspace.DeliveryPullRequest {
		return workflowError(4, "human_integration_required: the developer mandate stops before integration; release the candidate and use ply integration")
	}
	return nil
}

// AcceptHumanIntegration preserves a caller's actual answer and exact plan. It
// requires released ownership independently of the claim made at the prompt.
// No callback token or live provider process is needed after explicit release.
func AcceptHumanIntegration(d Dependencies, root, id string, in HumanIntegrationInput) (HumanIntegrationAcceptance, error) {
	var out HumanIntegrationAcceptance
	if !plain(in.PlanID, 1, 256) || !digestPattern.MatchString(in.PlanSHA256) {
		return out, workflowError(2, "human integration requires its exact plan ID and digest")
	}
	if e := validateHumanIntegrationDecision(in); e != nil {
		return out, e
	}
	resultID, e := workspace.ParseTaskResultID(in.Attestation.TaskResultID)
	if e != nil {
		return out, e
	}
	a, e := nativeDeliveryCandidate(d, root, id, resultID)
	if e != nil {
		return out, e
	}
	raw, e := Canonical(in.Attestation)
	if e != nil {
		return out, e
	}
	if _, e = workflowhandoff.ValidateDeliveryHumanAttestation(raw, string(a.Candidate.TaskResult.TaskID), a.Candidate.TaskResult, in.Attestation.Outcome); e != nil {
		return out, e
	}
	inputSHA := digest(in)
	var attestation FileBinding
	already := false
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if _, e := nativeDeliveryCandidate(d, root, id, resultID); e != nil {
			return e
		}
		currentCandidate := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		if in.Attestation.Outcome != "pass" && currentCandidate.Integration != nil {
			return workflowError(4, "integration_already_observed: a later negative answer cannot restore developer ownership after local integration")
		}
		// A negative decision clears the effective takeover fields so the
		// original immutable controller can resume corrections. The exact
		// immutable event still makes repeats idempotent without new authority.
		for i := len(s.Result.Delivery.Events) - 1; i >= 0; i-- {
			event := s.Result.Delivery.Events[i]
			if event.Kind != "human_integration" {
				continue
			}
			var prior HumanIntegrationAcceptance
			raw, e := workflowBound(event.Binding, 1<<20)
			if e != nil {
				return e
			}
			if e = decode(raw, 1<<20, &prior); e != nil {
				return e
			}
			if !prior.Accepted && prior.RunID == id && prior.RequestSHA256 == s.Result.RequestSHA256 && prior.PlanID == in.PlanID && prior.PlanSHA256 == in.PlanSHA256 && equal(prior.Decision, in.Decision) && equal(prior.Attestation, in.Attestation) {
				out, already = prior, true
				return nil
			}
		}
		if prior := s.Result.Delivery.HumanIntegration; prior != nil && prior.PlanID == in.PlanID && prior.PlanSHA256 == in.PlanSHA256 && equal(prior.Decision, in.Decision) && equal(prior.Attestation, in.Attestation) {
			out = *prior
			already = true
			return nil
		}
		c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		release, e := deliveryRelease(*s, c)
		if e != nil {
			return e
		}
		if release == nil {
			return workflowError(4, "owner_active: the bound developer must release source ownership before human integration")
		}
		if prior := s.Result.Delivery.HumanIntegration; prior != nil && prior.Accepted && (prior.PlanID != in.PlanID || prior.PlanSHA256 != in.PlanSHA256) {
			return workflowError(4, "plan_changed: an accepted integration already owns this candidate; resume or resolve that operation before another plan")
		}
		if prior := s.Result.Delivery.Attempt; prior != nil && prior.State == "attempted" && (prior.Kind != "human_integration" || prior.InputSHA256 != inputSHA) {
			return workflowError(4, "another delivery effect is unresolved")
		}
		path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "human-integration", inputSHA[7:])
		if _, e = workflowKeep(d, filepath.Join(path, "request.json"), in); e != nil {
			return e
		}
		attestation, e = workflowKeep(d, filepath.Join(path, "attestation.json"), in.Attestation)
		if e != nil {
			return e
		}
		s.Result.Delivery.Attempt = &DeliveryAttempt{ID: "human-integration-" + inputSHA[7:], Kind: "human_integration", State: "attempted", Path: path, CandidateOID: c.OID, CandidateTree: c.Tree, InputSHA256: inputSHA}
		return nil
	})
	if e != nil || already {
		return out, e
	}
	if e = d.fault("human_integration_after_reservation"); e != nil {
		return out, e
	}
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if s.Result.Delivery.Attempt == nil || s.Result.Delivery.Attempt.InputSHA256 != inputSHA {
			return workflowError(4, "human integration reservation changed")
		}
		current, e := nativeDeliveryCandidate(d, root, id, resultID)
		if e != nil {
			return e
		}
		release, e := deliveryRelease(*s, current.Candidate)
		if e != nil {
			return e
		}
		if release == nil {
			return workflowError(4, "owner_active: source ownership is no longer released")
		}
		if e = validateHumanIntegrationDecision(in); e != nil {
			return e
		}
		qa, e := workflowhandoff.RecordDeliveryHumanQA(d.Workflow, string(current.Candidate.TaskResult.TaskID), current.Candidate.TaskResult, in.Attestation.Outcome, attestation.Locator)
		if e != nil {
			return e
		}
		out = HumanIntegrationAcceptance{Kind: "PlyHumanIntegrationDecision@1", SchemaVersion: 1, RunID: id, RequestSHA256: s.Result.RequestSHA256, PlanID: in.PlanID, PlanSHA256: in.PlanSHA256, Decision: in.Decision, HumanQA: qa, Release: *release, Accepted: in.Attestation.Outcome == "pass", Attestation: in.Attestation}
		if _, e = deliveryAppendEvent(d, s, "human-integration-"+inputSHA[7:], "human_integration", out); e != nil {
			return e
		}
		s.Result.Delivery.HumanIntegration = &out
		// A published PR's completion keeps its original publication QA. The
		// later merge decision is a separate candidate-bound native QA record.
		if s.Result.Delivery.Phase != "completed" || current.Candidate.PullRequest == nil {
			s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1].HumanQA = &qa
		}
		s.Result.Delivery.Attempt.State = "recorded"
		if out.Accepted {
			if s.Result.Delivery.Phase != "completed" {
				s.Result.Delivery.Phase = "human_qa_passed"
			}
			s.Result.NextAction = WorkflowAction{"user", "The exact human answer and plan are preserved. Continue this human integration operation; the developer has no integration authority."}
		} else {
			s.Result.Delivery.OwnershipRelease = nil
			s.Result.Delivery.HumanIntegration = nil
			s.Result.Delivery.Phase = "awaiting_human_qa"
			s.Result.NextAction = WorkflowAction{"recipient", "The human rejected or blocked this plan. Correct the recorded issue and qualify and release the exact candidate for a new human decision."}
		}
		s.Result.Round.State = s.Result.Delivery.Phase
		return nil
	})
	return out, e
}

func preservedHumanIntegration(s workflowState) (*workspace.HumanIntegrationAuthorization, error) {
	if s.Result.Delivery == nil || s.Result.Delivery.HumanIntegration == nil || !s.Result.Delivery.HumanIntegration.Accepted {
		return nil, nil
	}
	h := s.Result.Delivery.HumanIntegration
	if h.Kind != "PlyHumanIntegrationDecision@1" || h.SchemaVersion != 1 || h.RunID != s.Result.RunID || h.RequestSHA256 != s.Result.RequestSHA256 {
		return nil, workflowError(4, "human integration decision run identity changed")
	}
	if len(s.Result.Delivery.Candidates) == 0 {
		return nil, workflowError(4, "human integration candidate is missing")
	}
	c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	release, e := deliveryRelease(s, c)
	if e != nil {
		return nil, e
	}
	if release == nil || !equal(*release, h.Release) || h.HumanQA.TaskResultID != c.TaskResult.ID || h.HumanQA.ResultOID != c.OID || h.HumanQA.ResultTree != c.Tree || h.HumanQA.Outcome != "pass" {
		return nil, workflowError(4, "human integration no longer binds released candidate and actual pass")
	}
	if e = validateHumanIntegrationDecision(HumanIntegrationInput{PlanID: h.PlanID, PlanSHA256: h.PlanSHA256, Decision: h.Decision, Attestation: h.Attestation}); e != nil {
		return nil, e
	}
	raw, e := Canonical(h.Attestation)
	if e != nil {
		return nil, e
	}
	if _, e = workflowhandoff.ValidateDeliveryHumanAttestation(raw, string(c.TaskResult.TaskID), c.TaskResult, "pass"); e != nil {
		return nil, e
	}
	found := false
	for _, ev := range s.Result.Delivery.Events {
		if ev.Kind == "human_integration" {
			b, e := workflowBound(ev.Binding, 1<<20)
			if e != nil {
				return nil, e
			}
			var recorded HumanIntegrationAcceptance
			if e = decode(b, 1<<20, &recorded); e != nil {
				return nil, e
			}
			if equal(recorded, *h) {
				found = true
			}
		}
	}
	if !found {
		return nil, workflowError(4, "human integration is absent from native decision history")
	}
	return &workspace.HumanIntegrationAuthorization{SchemaVersion: 1, PlanID: h.PlanID, PlanSHA256: h.PlanSHA256, DecisionLocator: h.Decision.Locator, DecisionSHA256: h.Decision.SHA256, ReleaseLocator: h.Release.Locator, ReleaseSHA256: h.Release.SHA256, HumanQARecordID: h.HumanQA.ID, TaskResultID: c.TaskResult.ID}, nil
}

// CompleteHumanLocalDelivery grants only the accepted native plan. It never
// upgrades the old request and cannot publish or merge a PR.
func CompleteHumanLocalDelivery(d Dependencies, root, id string, resultID workspace.TaskResultID, planSHA256 string) (workflowhandoff.DeliveryIntegrationResult, error) {
	var out workflowhandoff.DeliveryIntegrationResult
	a, e := nativeDeliveryCandidate(d, root, id, resultID)
	if e != nil {
		return out, e
	}
	h := a.Authorization.HumanIntegration
	if h == nil || h.PlanSHA256 != planSHA256 || h.TaskResultID != resultID {
		return out, workflowError(4, "human_qa_required: exact preserved human plan acceptance is missing")
	}
	if e = requireDeliveryPass(a); e != nil {
		return out, e
	}
	if a.Candidate.HumanQA.ID != h.HumanQARecordID {
		return out, workflowError(4, "human_qa_required: a later human answer supersedes this plan")
	}
	return completeLocalDelivery(d, root, id, resultID, h.HumanQARecordID, true)
}
