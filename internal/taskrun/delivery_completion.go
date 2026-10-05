package taskrun

import (
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

func WorkflowDeliveryQA(d Dependencies, root, id, contextPath, outcome, evidencePath string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	if outcome != "pass" && outcome != "fail" && outcome != "blocked" {
		return WorkflowRun{}, workflowError(2, "human QA outcome must be pass, fail or blocked")
	}
	evidence, evidenceRaw, e := deliveryReadBinding(evidencePath, 1<<20)
	if e != nil {
		return WorkflowRun{}, e
	}
	var state workflowState
	var candidate DeliveryCandidate
	alreadyRecorded := false
	inputSHA, attemptID := "", ""
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if e := deliveryCallback(d, *s, contextPath, true); e != nil {
			return e
		}
		if e := workflowDeliveryFresh(d, *s, true); e != nil {
			return e
		}
		if len(s.Result.Delivery.Candidates) == 0 {
			return workflowError(4, "human QA requires a qualified exact candidate")
		}
		candidate = s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		// Invalid input has no side effect to recover. Validate before reserving
		// so a corrected real answer is not blocked by an unknown-attempt guard.
		if _, e := workflowhandoff.ValidateDeliveryHumanAttestation(evidenceRaw, string(candidate.TaskResult.TaskID), candidate.TaskResult, outcome); e != nil {
			return e
		}
		x, e := deliveryTarget(d, *s)
		if e != nil {
			return e
		}
		if x.OID != candidate.OID || x.Tree != candidate.Tree {
			return workflowError(4, "human QA cannot transfer to a changed candidate")
		}
		inputSHA = digest([]any{candidate.TaskResult.ID, outcome, evidence})
		attemptID = "qa-" + inputSHA[len("sha256:"):]
		for _, prior := range s.Result.Delivery.Events {
			if prior.ID == attemptID && prior.Kind == "human_qa" {
				alreadyRecorded = true
				return nil
			}
		}
		if prior := s.Result.Delivery.Attempt; prior != nil && prior.State == "attempted" && (prior.Kind != "human_qa" || prior.InputSHA256 != inputSHA) {
			return workflowError(4, "another delivery effect is unresolved")
		}
		path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", attemptID)
		request := map[string]any{"kind": "PlyDeliveryHumanQARequest@1", "run_id": id, "request_sha256": s.Result.RequestSHA256, "candidate": candidate.TaskResult.ID, "candidate_oid": candidate.OID, "candidate_tree": candidate.Tree, "outcome": outcome, "evidence": evidence}
		if e := d.writeValue(filepath.Join(path, "attempt.json"), request); e != nil {
			return e
		}
		s.Result.Delivery.Attempt = &DeliveryAttempt{ID: attemptID, Kind: "human_qa", State: "attempted", Path: path, CandidateOID: candidate.OID, CandidateTree: candidate.Tree, InputSHA256: inputSHA}
		state = *s
		return nil
	})
	if e != nil {
		o, _ := WorkflowShow(d, root, id)
		return o, e
	}
	if alreadyRecorded {
		return WorkflowShow(d, root, id)
	}
	if e = d.fault("delivery_after_qa_reservation"); e != nil {
		o, _ := WorkflowShow(d, root, id)
		return o, e
	}
	qa, e := workflowhandoff.RecordDeliveryHumanQA(d.Workflow, string(candidate.TaskResult.TaskID), candidate.TaskResult, outcome, evidencePath)
	if e != nil {
		o, _ := WorkflowShow(d, root, id)
		return o, e
	}
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if s.Result.Delivery.Attempt == nil || s.Result.Delivery.Attempt.InputSHA256 != inputSHA {
			return workflowError(4, "human QA reservation changed")
		}
		if _, e := workflowBound(evidence, 1<<20); e != nil {
			return e
		}
		value := map[string]any{"kind": "PlyDeliveryHumanQA@1", "run_id": id, "request_sha256": state.Result.RequestSHA256, "candidate": candidate.Key, "record": qa, "source": evidence}
		if _, e := deliveryAppendEvent(d, s, attemptID, "human_qa", value); e != nil {
			return e
		}
		s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1].HumanQA = &qa
		s.Result.Delivery.Attempt.State = "recorded"
		s.Result.Delivery.Phase = "working"
		s.Result.NextAction = WorkflowAction{"recipient", "Use the preserved human feedback to correct this delivery, then qualify a new candidate; the old verdict stays with its original bytes."}
		if outcome == "pass" {
			s.Result.Delivery.Phase = "human_qa_passed"
			s.Result.NextAction = WorkflowAction{"recipient", "Complete the already authorized local integration for this exact passed candidate and update the Epic base."}
		} else if outcome == "blocked" {
			s.Result.Delivery.Phase = "awaiting_human_qa"
			s.Result.NextAction = WorkflowAction{"recipient", "Resolve the recorded product-QA block within the selected goal and request an actual human judgment when ready."}
		}
		s.Result.Round.State = s.Result.Delivery.Phase
		return nil
	})
	return deliveryReadback(d, root, id, e)
}

func WorkflowDeliveryIntegrate(d Dependencies, root, id, contextPath string) (WorkflowRun, error) {
	if e := containing(d, root); e != nil {
		return WorkflowRun{}, e
	}
	var state workflowState
	var candidate DeliveryCandidate
	var inputSHA, attemptID string
	e := workflowUpdate(d, root, id, func(s *workflowState) error {
		if deliveryRun(s.Request) && s.Result.Delivery != nil && s.Result.Delivery.Phase == "completed" {
			state = *s
			return nil
		}
		if e := deliveryCallback(d, *s, contextPath, true); e != nil {
			return e
		}
		if e := workflowDeliveryFresh(d, *s, true); e != nil {
			return e
		}
		if s.Request.Delivery.LocalIntegration != "after_human_pass" || len(s.Result.Delivery.Candidates) == 0 {
			return workflowError(4, "local integration is not covered by this delivery")
		}
		candidate = s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		if candidate.HumanQA == nil || candidate.HumanQA.Outcome != "pass" || candidate.HumanQA.TaskResultID != candidate.TaskResult.ID || candidate.HumanQA.ResultOID != candidate.OID || candidate.HumanQA.ResultTree != candidate.Tree || (s.Result.Delivery.Phase != "human_qa_passed" && s.Result.Delivery.Phase != "integrating") {
			return workflowError(4, "integration requires the current candidate's actual human pass")
		}
		inputSHA = digest([]any{candidate.TaskResult.ID, candidate.HumanQA.ID, candidate.OID, s.Observed.Epic.OID})
		attemptID = "integrate-" + candidate.Key
		if prior := s.Result.Delivery.Attempt; prior != nil && prior.State == "attempted" && (prior.Kind != "integration" || prior.InputSHA256 != inputSHA) {
			return workflowError(4, "another delivery effect is unresolved")
		}
		path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", attemptID)
		value := map[string]any{"kind": "PlyDeliveryIntegrationRequest@1", "run_id": id, "request_sha256": s.Result.RequestSHA256, "candidate": candidate.TaskResult.ID, "human_qa": candidate.HumanQA.ID, "expected_parent_oid": s.Observed.Epic.OID, "result_oid": candidate.OID}
		if e := d.writeValue(filepath.Join(path, "attempt.json"), value); e != nil {
			return e
		}
		s.Result.Delivery.Attempt = &DeliveryAttempt{ID: attemptID, Kind: "integration", State: "attempted", Path: path, CandidateOID: candidate.OID, CandidateTree: candidate.Tree, InputSHA256: inputSHA}
		s.Result.Delivery.Phase = "integrating"
		state = *s
		return nil
	})
	if e != nil {
		o, _ := WorkflowShow(d, root, id)
		return o, e
	}
	if state.Result.Delivery.Phase == "completed" {
		return WorkflowShow(d, root, id)
	}
	if e = d.fault("delivery_after_integration_reservation"); e != nil {
		o, _ := WorkflowShow(d, root, id)
		return o, e
	}
	result, integrationErr := workflowhandoff.IntegrateDeliveryCandidate(d.Workflow, string(candidate.TaskResult.TaskID), candidate.TaskResult.ID, candidate.HumanQA.ID, state.Observed.Epic.OID, candidate.OID, workflowhandoff.DeliveryOwnerAuthority{RunID: id, RequestSHA256: state.Result.RequestSHA256, ActorClaim: state.Request.Delivery.OwnerClaim, PreparationID: state.Request.PreparationID})
	e = workflowUpdate(d, root, id, func(s *workflowState) error {
		if s.Result.Delivery.Attempt == nil || s.Result.Delivery.Attempt.InputSHA256 != inputSHA {
			return workflowError(4, "integration reservation changed")
		}
		// Native integration owns at-most-once effect recovery. Each observation is
		// preserved, including unknown or interrupted outcomes, before readback.
		value := map[string]any{"kind": "PlyDeliveryIntegration@1", "run_id": id, "request_sha256": s.Result.RequestSHA256, "candidate": candidate.Key, "result": result}
		if integrationErr != nil {
			value["error"] = integrationErr.Error()
		}
		eventID := fmt.Sprintf("%s-observation-%08d", attemptID, len(s.Result.Delivery.Events)+1)
		binding, err := deliveryAppendEvent(d, s, eventID, "integration", value)
		if err != nil {
			return err
		}
		if integrationErr != nil || !result.Completed {
			s.Result.NextAction = WorkflowAction{"recipient", "Inspect the preserved native integration and base-update state; resume only the same bound operation, without repeating an unknown effect."}
			return nil
		}
		s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1].Integration = &binding
		s.Result.Delivery.Attempt.State = "recorded"
		s.Result.Delivery.Phase = "completed"
		s.Result.Round.State = "completed"
		s.Result.FinalReturn.State = "completed"
		s.Result.FinalReturn.ReportSHA256 = &binding.SHA256
		s.Result.NextAction = WorkflowAction{"user", "The exact passed candidate is integrated locally and the Epic base is current. The selected delivery is complete."}
		return nil
	})
	if e == nil {
		e = integrationErr
	}
	if e == nil && !result.Completed {
		e = workflowError(5, "local delivery is not yet complete; preserved native state determines the next action")
	}
	return deliveryReadback(d, root, id, e)
}
