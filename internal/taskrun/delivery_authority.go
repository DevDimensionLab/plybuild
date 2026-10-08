package taskrun

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type DeliveryAuthority struct {
	RunID             string                           `json:"run_id"`
	RequestSHA256     string                           `json:"request_sha256"`
	MandateSHA256     string                           `json:"mandate_sha256"`
	PreparationID     string                           `json:"preparation_id"`
	OwnerClaim        string                           `json:"owner_claim"`
	SourceRef         string                           `json:"source_ref"`
	ExpectedParentOID string                           `json:"expected_parent_oid"`
	Agreement         workspace.DeliveryAgreement      `json:"agreement"`
	Authorization     *workspace.DeliveryAuthorization `json:"authorization"`
	Candidate         DeliveryCandidate                `json:"candidate"`
}

// deliveryAuthorityState validates durable authority without requiring the old
// provider process, executable or conversation to still exist. The accepted
// policy and exact frozen input bytes remain required.
func deliveryAuthorityState(d Dependencies, s workflowState) (*workspace.DeliveryAuthorization, error) {
	if !deliveryRun(s.Request) || s.Request.Delivery.Agreement == nil || s.Result.Delivery == nil || s.Acceptance == nil || s.StartSHA256 == nil || s.Result.Delivery.PermissionState != "recipient_confirmed_contract" {
		return nil, workflowError(4, "delivery requires an explicit frozen agreement and actual accepted authority; historical runs cannot acquire new effects")
	}
	if err := workflowhandoff.ValidateDeliveryMandate(s.Request.HandoffDraft, s.Request.Delivery.Agreement); err != nil {
		return nil, err
	}
	if _, err := workflowBound(s.Request.Delivery.Goal, 4<<20); err != nil {
		return nil, err
	}
	raw, err := workflowBound(*s.Acceptance, 1<<20)
	if err != nil {
		return nil, err
	}
	var acceptance DeliveryAcceptance
	if err = decode(raw, 1<<20, &acceptance); err != nil {
		return nil, err
	}
	if acceptance.Envelope != deliveryEnv("run-acceptance") || acceptance.Acceptance.Acceptance != "started" {
		return nil, workflowError(4, "delivery authority lacks positive native acceptance")
	}
	if err = workflowClaim(s, acceptance.RunID, acceptance.RequestSHA256, acceptance.SessionID); err != nil {
		return nil, err
	}
	if err = positiveDeliveryClaim(acceptance.Acceptance, acceptance.DeliveryPermission, s.Request); err != nil {
		return nil, err
	}
	if !equal(acceptance.RuntimeClaim.EffectivePolicySHA256, s.Result.Delivery.ActualPolicySHA256) || !equal(acceptance.DeliveryPermission.ActualPolicyEvidence, s.Result.Delivery.ActualPolicyEvidence) {
		return nil, workflowError(4, "actual delivery permission differs from its preserved acceptance")
	}
	mandate, err := workflowBound(FileBinding{s.Result.Handoff.Locator, s.Result.Handoff.SHA256}, 2<<20)
	if err != nil {
		return nil, err
	}
	var h struct {
		Source   string        `json:"source_draft_sha256"`
		Inputs   []FileBinding `json:"inputs"`
		Delivery struct {
			Agreement *workspace.DeliveryAgreement `json:"agreement"`
		} `json:"delivery_binding"`
	}
	if err = json.Unmarshal(mandate, &h); err != nil {
		return nil, err
	}
	if !equal(h.Delivery.Agreement, s.Request.Delivery.Agreement) {
		return nil, workflowError(4, "native mandate delivery agreement changed")
	}
	if h.Source != hash(s.Request.HandoffDraft) {
		return nil, workflowError(4, "native mandate no longer binds the frozen draft authority")
	}
	for _, input := range h.Inputs {
		if _, err = workflowBound(input, 64<<20); err != nil {
			return nil, err
		}
	}
	facts, err := workflowhandoff.ReadTaskRunReturnFacts(d.Workflow, s.Result.Handoff.Locator)
	if err != nil {
		return nil, err
	}
	if facts.Start == nil || facts.Start.SHA256 != *s.StartSHA256 {
		return nil, workflowError(4, "delivery native start evidence changed")
	}
	human, err := preservedHumanIntegration(s)
	if err != nil {
		return nil, err
	}
	return &workspace.DeliveryAuthorization{Agreement: *s.Request.Delivery.Agreement, AgreementSHA256: workspace.DeliveryAgreementDigest(*s.Request.Delivery.Agreement), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, MandateSHA256: s.Result.Handoff.SHA256, AllowedEffects: append([]string{}, acceptance.DeliveryPermission.AllowedEffects...), PermissionConfirmed: true, HumanIntegration: human, HumanIntegrationRequired: s.Request.Delivery.Agreement.HumanOwnedIntegration() && s.Request.Delivery.Agreement.Mode != workspace.DeliveryPullRequest || s.Result.Delivery.OwnershipRelease != nil}, nil
}

func deliveryEvidenceAuthorization(d Dependencies, request workspace.TaskHandoffEvidenceRequest) (*workspace.DeliveryAuthorization, error) {
	raw, err := workflowBound(FileBinding{request.HandoffLocator, request.HandoffSHA256}, 2<<20)
	if err != nil {
		return nil, err
	}
	runID, requestSHA, err := workflowhandoff.DeliveryCandidateRuntimeBinding(raw)
	if err != nil || runID == "" {
		return nil, err
	}
	var h struct {
		Workspace struct {
			Root string `json:"root"`
		} `json:"workspace_binding"`
		Delivery struct {
			ParentLocator string                       `json:"parent_handoff_locator"`
			ParentSHA     string                       `json:"parent_handoff_sha256"`
			Agreement     *workspace.DeliveryAgreement `json:"agreement"`
			CandidateKey  string                       `json:"candidate_key"`
		} `json:"delivery_binding"`
		Target struct {
			OID  string `json:"oid"`
			Tree string `json:"tree"`
		} `json:"target_binding"`
	}
	if err = json.Unmarshal(raw, &h); err != nil {
		return nil, err
	}
	s, err := workflowRead(h.Workspace.Root, runID)
	if err != nil {
		return nil, err
	}
	if s.Result.RequestSHA256 != requestSHA || h.Delivery.ParentLocator != s.Result.Handoff.Locator || h.Delivery.ParentSHA != s.Result.Handoff.SHA256 || s.Request.Delivery == nil || !equal(h.Delivery.Agreement, s.Request.Delivery.Agreement) {
		return nil, workflowError(4, "candidate does not bind its exact frozen workflow authority")
	}
	// Historical technical evidence remains readable, but it cannot carry live
	// delivery authority after correction, fail or a replacement candidate.
	current := false
	if ds := s.Result.Delivery; ds != nil {
		if len(ds.Candidates) > 0 && (ds.Phase == "awaiting_human_qa" || ds.Phase == "human_qa_passed" || ds.Phase == "integrating" || ds.Phase == "completed") {
			c := ds.Candidates[len(ds.Candidates)-1]
			current = c.Handoff.Locator == request.HandoffLocator && c.Handoff.SHA256 == request.HandoffSHA256 && c.OID == h.Target.OID && c.Tree == h.Target.Tree
		}
		// Qualification reads its native evidence before appending TaskResult
		// to workflow state. Only that exact in-flight verifier may bind it.
		if attempt := ds.Attempt; ds.Phase == "verifying" && attempt != nil && attempt.Kind == "verification" && attempt.ID == h.Delivery.CandidateKey && attempt.CandidateOID == h.Target.OID && attempt.CandidateTree == h.Target.Tree {
			current = true
		}
	}
	if !current {
		return nil, nil
	}
	a, err := deliveryAuthorityState(d, s)
	if err != nil {
		return nil, err
	}
	a.CandidateRunID = request.RunID
	return a, workspace.ValidateDeliveryAuthorization(*a)
}

// ValidateDeliveryAuthority does not require human QA, so registration can
// preserve a qualified candidate while its actual answer is still pending.
func ValidateDeliveryAuthority(d Dependencies, root, id string, result workspace.TaskResultRecord, agreement workspace.DeliveryAgreement) (DeliveryAuthority, error) {
	var out DeliveryAuthority
	s, err := workflowRead(root, id)
	if err != nil {
		return out, err
	}
	a, err := deliveryAuthorityState(d, s)
	if err != nil {
		return out, err
	}
	if !equal(a.Agreement, agreement) {
		return out, workflowError(4, "mode or target differs from the frozen delivery agreement; publish a revised Spec and authorize a new execution")
	}
	if len(s.Result.Delivery.Candidates) == 0 {
		return out, workflowError(4, "delivery requires a technically qualified native candidate")
	}
	c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	if !equal(c.TaskResult, result) || c.OID != result.ResultOID || c.Tree != result.ResultTree || c.TaskResult.SourceRef != agreement.SourceRef || s.Observed.Target.Ref != agreement.SourceRef {
		return out, workflowError(4, "delivery does not name the exact current qualified TaskResult and source")
	}
	if s.Result.Delivery.Phase != "awaiting_human_qa" && s.Result.Delivery.Phase != "human_qa_passed" && s.Result.Delivery.Phase != "integrating" && s.Result.Delivery.Phase != "completed" {
		return out, workflowError(4, "delivery candidate has not reached technical qualification")
	}
	if _, err = workflowBound(c.Verification, 4<<20); err != nil {
		return out, err
	}
	registry, err := d.Workspace.WorkItems.Snapshot(root)
	if err != nil {
		return out, err
	}
	found := false
	for _, recorded := range registry.TaskResults {
		if recorded.ID == result.ID {
			found = equal(recorded, result)
		}
	}
	if !found {
		return out, workflowError(4, "qualified TaskResult is absent or changed in native registry")
	}
	request := taskResultEvidenceRequest(result)
	evidence, err := workflowhandoff.NewTaskHandoffEvidenceReader(d.Workflow).ReadTaskEvidence(request)
	if err != nil {
		return out, err
	}
	a.CandidateRunID = result.RunID
	if err = validateDeliveryTechnicalEvidence(result, evidence, a); err != nil {
		return out, err
	}
	for _, artifact := range result.Artifacts {
		if _, err = workflowBound(FileBinding{artifact.Locator, artifact.SHA256}, 64<<20); err != nil {
			return out, err
		}
	}
	x, err := deliveryTarget(d, s)
	if err != nil {
		return out, err
	}
	if x.OID != c.OID || x.Tree != c.Tree {
		return out, workflowError(4, "candidate changed after technical qualification")
	}
	// Native QA can be recorded outside the original owner conversation. Reuse
	// its exact candidate identity; never turn a technical pass into human pass.
	c.HumanQA, err = LatestDeliveryHumanQA(registry.HumanQARecords, result.ID, c.OID, c.Tree)
	if err != nil {
		return out, err
	}
	out = DeliveryAuthority{RunID: id, RequestSHA256: s.Result.RequestSHA256, MandateSHA256: s.Result.Handoff.SHA256, PreparationID: s.Request.PreparationID, OwnerClaim: s.Request.Delivery.OwnerClaim, SourceRef: agreement.SourceRef, ExpectedParentOID: s.Observed.Epic.OID, Agreement: agreement, Authorization: a, Candidate: c}
	return out, nil
}

func validateDeliveryTechnicalEvidence(result workspace.TaskResultRecord, evidence workspace.TaskHandoffEvidence, authority *workspace.DeliveryAuthorization) error {
	if (result.TechnicalGate != "passed" && result.TechnicalGate != "good_enough_with_known_debt") || !evidence.TaskSpecValid || !evidence.TaskRequirementsValid || !equal(evidence.DeliveryAuthorization, authority) {
		return workflowError(4, "delivery candidate evidence or accepted authority no longer validates")
	}
	return workspace.ValidateTaskResultEvidence(result, evidence)
}

// LatestDeliveryHumanQA shares the native integration chronology rule with
// delivery registration/execution. Registry ordering is an identity ordering,
// never the order in which actual human judgments occurred.
func LatestDeliveryHumanQA(records []workspace.TaskHumanQARecord, resultID workspace.TaskResultID, oid, tree string) (*workspace.TaskHumanQARecord, error) {
	return workspace.LatestTaskHumanQA(records, resultID, oid, tree)
}

func taskResultEvidenceRequest(r workspace.TaskResultRecord) workspace.TaskHandoffEvidenceRequest {
	return workspace.TaskHandoffEvidenceRequest{ActivityID: r.ActivityID, RunID: r.RunID, HandoffID: r.HandoffID, HandoffLocator: r.HandoffLocator, HandoffSHA256: r.HandoffSHA256, StartReceiptID: r.StartReceiptID, StartReceiptLocator: r.StartReceiptLocator, StartReceiptSHA256: r.StartReceiptSHA256, TerminalResultID: r.TerminalResultID, TerminalResultLocator: r.TerminalResultLocator, TerminalResultSHA256: r.TerminalResultSHA256, InspectionSHA256: r.InspectionSHA256}
}

func nativeDeliveryCandidate(d Dependencies, root, id string, resultID workspace.TaskResultID) (DeliveryAuthority, error) {
	s, err := workflowRead(root, id)
	if err != nil {
		return DeliveryAuthority{}, err
	}
	if !deliveryRun(s.Request) || s.Request.Delivery.Agreement == nil || s.Result.Delivery == nil || len(s.Result.Delivery.Candidates) == 0 {
		return DeliveryAuthority{}, workflowError(4, "explicit native delivery authority is missing")
	}
	c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	if c.TaskResult.ID != resultID {
		return DeliveryAuthority{}, workflowError(4, "delivery candidate differs")
	}
	return ValidateDeliveryAuthority(d, root, id, c.TaskResult, *s.Request.Delivery.Agreement)
}

func requireDeliveryPass(a DeliveryAuthority) error {
	qa := a.Candidate.HumanQA
	if qa == nil || qa.Outcome != "pass" || qa.TaskResultID != a.Candidate.TaskResult.ID || qa.ResultOID != a.Candidate.OID || qa.ResultTree != a.Candidate.Tree {
		return workflowError(4, "delivery requires actual human pass for this exact candidate")
	}
	for _, evidence := range qa.Evidence {
		raw, err := workflowBound(FileBinding{evidence.Locator, evidence.SHA256}, 64<<20)
		if err != nil {
			return err
		}
		if int64(len(raw)) != evidence.SizeBytes {
			return workflowError(4, "candidate human QA evidence size changed")
		}
	}
	return nil
}

type PullRequestDeliveryObservation struct {
	DeliveryID      string `json:"delivery_id"`
	Repository      string `json:"repository"`
	HeadRef         string `json:"head_ref"`
	HeadOID         string `json:"head_oid"`
	BaseRef         string `json:"base_ref"`
	URL             string `json:"url"`
	Number          int    `json:"number"`
	MetadataApplied bool   `json:"metadata_applied"`
}

// CompletePullRequestDelivery closes only this Task's queue and records the
// observed PR receipt. It never writes a target ref, integration or Epic base.
func CompletePullRequestDelivery(d Dependencies, root, id string, resultID workspace.TaskResultID, receipt FileBinding, observed PullRequestDeliveryObservation) (WorkflowRun, error) {
	if e := CheckAgentDeliveryExecution(d, root, id); e != nil {
		return WorkflowRun{}, e
	}
	a, err := nativeDeliveryCandidate(d, root, id, resultID)
	if err != nil {
		return WorkflowRun{}, err
	}
	if a.Authorization.HumanIntegrationRequired || a.Authorization.HumanIntegration != nil {
		return WorkflowRun{}, workflowError(4, "human_integration_required: source ownership changed before PR completion")
	}
	if err = requireDeliveryPass(a); err != nil {
		return WorkflowRun{}, err
	}
	if a.Agreement.Mode != workspace.DeliveryPullRequest || !strings.EqualFold(observed.Repository, a.Agreement.GitHubRepository) || observed.HeadRef != a.SourceRef || observed.HeadOID != a.Candidate.OID || observed.BaseRef != a.Agreement.TargetRef || observed.Number < 1 || observed.URL == "" || !observed.MetadataApplied {
		return WorkflowRun{}, workflowError(4, "observed PR does not satisfy the exact delivery agreement and metadata boundary")
	}
	if _, err = workflowBound(receipt, 4<<20); err != nil {
		return WorkflowRun{}, err
	}
	s, err := workflowRead(root, id)
	if err != nil {
		return WorkflowRun{}, err
	}
	if s.Result.Delivery.Phase == "completed" && s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1].PullRequest != nil {
		prior := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1].PullRequest
		raw, err := workflowBound(*prior, 4<<20)
		if err != nil {
			return WorkflowRun{}, err
		}
		var recorded struct {
			Observed PullRequestDeliveryObservation `json:"observed"`
		}
		if err = json.Unmarshal(raw, &recorded); err != nil {
			return WorkflowRun{}, err
		}
		want, got := recorded.Observed, observed
		want.DeliveryID, got.DeliveryID = "", ""
		if !equal(want, got) {
			return WorkflowRun{}, workflowError(4, "completed PR receipt differs; preserve its original observed delivery")
		}
		return WorkflowShow(d, root, id)
	}
	qt := workspace.QueueTargetInput{ProjectID: a.Agreement.ProjectID, RepoID: a.Agreement.RepoID, EpicID: a.Agreement.EpicID}
	if _, err = workflowhandoff.CloseDeliveryTaskQueue(d.Workspace, qt, a.Candidate.TaskResult.TaskID, a.PreparationID, "Exact candidate delivered through verified PR; no local integration or merge."); err != nil {
		return WorkflowRun{}, err
	}
	err = workflowUpdate(d, root, id, func(s *workflowState) error {
		if len(s.Result.Delivery.Candidates) == 0 || s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1].TaskResult.ID != resultID {
			return workflowError(4, "candidate changed before PR closure")
		}
		if prior := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1].PullRequest; s.Result.Delivery.Phase == "completed" && prior != nil {
			raw, e := workflowBound(*prior, 4<<20)
			if e != nil {
				return e
			}
			var recorded struct {
				Observed PullRequestDeliveryObservation `json:"observed"`
			}
			if e = json.Unmarshal(raw, &recorded); e != nil {
				return e
			}
			want, got := recorded.Observed, observed
			want.DeliveryID, got.DeliveryID = "", ""
			if !equal(want, got) {
				return workflowError(4, "concurrent PR completion differs from the preserved receipt")
			}
			return nil
		}
		binding, e := deliveryAppendEvent(d, s, "pull-request-"+a.Candidate.Key, "pull_request", map[string]any{"kind": "PlyDeliveryPullRequest@1", "run_id": id, "request_sha256": a.RequestSHA256, "candidate": a.Candidate.Key, "receipt": receipt, "observed": observed, "recorded_at_utc": d.Now().UTC().Format(time.RFC3339Nano)})
		if e != nil {
			return e
		}
		c := &s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		c.PullRequest, c.HumanQA = &binding, a.Candidate.HumanQA
		s.Result.Delivery.Phase, s.Result.Round.State, s.Result.FinalReturn.State = "completed", "completed", "completed"
		s.Result.FinalReturn.ReportSHA256 = &binding.SHA256
		s.Result.NextAction = WorkflowAction{"user", "The exact passed candidate is delivered through the observed PR. The target branch and Epic base were not updated; PR merge remains a separate decision."}
		return nil
	})
	return deliveryReadback(d, root, id, err)
}

// CompleteLocalDelivery resumes the native at-most-once integration ledger,
// then preserves truthful queue/base/runtime completion independently of the
// original interactive owner process.
func CompleteLocalDelivery(d Dependencies, root, id string, resultID workspace.TaskResultID, qaID workspace.HumanQARecordID) (workflowhandoff.DeliveryIntegrationResult, error) {
	return completeLocalDelivery(d, root, id, resultID, qaID, false)
}

func completeLocalDelivery(d Dependencies, root, id string, resultID workspace.TaskResultID, qaID workspace.HumanQARecordID, human bool) (workflowhandoff.DeliveryIntegrationResult, error) {
	var out workflowhandoff.DeliveryIntegrationResult
	if !human {
		if e := CheckAgentDeliveryExecution(d, root, id); e != nil {
			return out, e
		}
		if e := d.fault("local_delivery_after_agent_gate"); e != nil {
			return out, e
		}
	}
	a, err := nativeDeliveryCandidate(d, root, id, resultID)
	if err != nil {
		return out, err
	}
	if !human && (a.Authorization.HumanIntegrationRequired || a.Authorization.HumanIntegration != nil) {
		return out, workflowError(4, "human_integration_required: source ownership changed before local delivery execution")
	}
	if err = requireDeliveryPass(a); err != nil {
		return out, err
	}
	if a.Agreement.Mode == workspace.DeliveryPullRequest || a.Candidate.HumanQA.ID != qaID {
		return out, workflowError(4, "local delivery mode or candidate human QA differs")
	}
	// The reservation spans merge, queue and base effects. A released owner or
	// later human answer may not race the workspace's final authority check.
	inputSHA := digest([]any{resultID, qaID, a.Candidate.OID, a.ExpectedParentOID, a.Authorization, human})
	attemptID := "local-integration-" + inputSHA[7:]
	completed := false
	err = workflowUpdate(d, root, id, func(s *workflowState) error {
		current, e := nativeDeliveryCandidate(d, root, id, resultID)
		if e != nil {
			return e
		}
		if !equal(current.Authorization, a.Authorization) || current.Candidate.HumanQA == nil || current.Candidate.HumanQA.ID != qaID {
			return workflowError(4, "delivery authority changed before effect reservation")
		}
		if s.Result.Delivery.Phase == "completed" && s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1].Integration != nil {
			completed = true
			return nil
		}
		if prior := s.Result.Delivery.Attempt; prior != nil && prior.State == "attempted" && (prior.Kind != "integration" || prior.InputSHA256 != inputSHA) {
			return workflowError(4, "another delivery effect is unresolved")
		}
		path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", attemptID)
		if _, e = workflowKeep(d, filepath.Join(path, "attempt.json"), map[string]any{"kind": "PlyDeliveryIntegrationRequest@1", "run_id": id, "request_sha256": a.RequestSHA256, "candidate": resultID, "human_qa": qaID, "expected_parent_oid": a.ExpectedParentOID, "result_oid": a.Candidate.OID, "authorization": a.Authorization, "human_integration": human}); e != nil {
			return e
		}
		s.Result.Delivery.Attempt = &DeliveryAttempt{ID: attemptID, Kind: "integration", State: "attempted", Path: path, CandidateOID: a.Candidate.OID, CandidateTree: a.Candidate.Tree, InputSHA256: inputSHA}
		s.Result.Delivery.Phase = "integrating"
		return nil
	})
	if err != nil {
		return out, err
	}
	if !completed {
		if err = d.fault("delivery_after_integration_reservation"); err != nil {
			return out, err
		}
	}
	var integrationErr error
	out, integrationErr = workflowhandoff.IntegrateDeliveryCandidate(d.Workflow, string(a.Candidate.TaskResult.TaskID), resultID, qaID, a.ExpectedParentOID, a.Candidate.OID, workflowhandoff.DeliveryOwnerAuthority{RunID: id, RequestSHA256: a.RequestSHA256, ActorClaim: a.OwnerClaim, PreparationID: a.PreparationID, Authorization: a.Authorization})
	err = workflowUpdate(d, root, id, func(s *workflowState) error {
		c := &s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		if c.TaskResult.ID != resultID {
			return workflowError(4, "candidate changed before local completion")
		}
		if s.Result.Delivery.Phase == "completed" && c.Integration != nil {
			return nil
		}
		if s.Result.Delivery.Attempt == nil || s.Result.Delivery.Attempt.InputSHA256 != inputSHA {
			return workflowError(4, "local integration reservation changed")
		}
		observation := map[string]any{"kind": "PlyDeliveryIntegration@1", "run_id": id, "request_sha256": a.RequestSHA256, "candidate": c.Key, "result": out}
		if integrationErr != nil {
			observation["error"] = integrationErr.Error()
		}
		binding, e := deliveryAppendEvent(d, s, fmt.Sprintf("%s-observation-%08d", attemptID, len(s.Result.Delivery.Events)+1), "integration", observation)
		if e != nil {
			return e
		}
		if integrationErr != nil || !out.Completed {
			s.Result.NextAction = WorkflowAction{"recipient", "Inspect the preserved native integration, queue and base state. Resume only this exact reserved operation; ownership cannot be released while its effects are unresolved."}
			return nil
		}
		c.Integration, c.HumanQA = &binding, a.Candidate.HumanQA
		s.Result.Delivery.Attempt.State = "recorded"
		s.Result.Delivery.Phase, s.Result.Round.State, s.Result.FinalReturn.State = "completed", "completed", "completed"
		s.Result.FinalReturn.ReportSHA256 = &binding.SHA256
		s.Result.NextAction = WorkflowAction{"user", fmt.Sprintf("The exact passed candidate is integrated locally in %s and the registered base is current.", a.Agreement.TargetRef)}
		return nil
	})
	if err == nil {
		err = integrationErr
	}
	if err == nil && !out.Completed {
		err = workflowError(5, "local delivery is not yet complete; preserved native state determines the next action")
	}
	return out, err
}
