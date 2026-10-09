package taskrun

import (
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// DeliveryStatus is a read-only projection, never a receipt or an authority
// input. Each outcome remains bound to its own candidate and evidence.
type DeliveryStatus struct {
	Source              DeliverySourceStatus               `json:"source"`
	Verification        *DeliveryVerificationStatus        `json:"verification"`
	QualifiedCandidate  *DeliveryQualifiedStatus           `json:"qualified_candidate"`
	HumanJudgment       *DeliveryHumanStatus               `json:"human_judgment"`
	Acceptance          *DeliveryAcceptanceStatus          `json:"acceptance,omitempty"`
	AcceptanceSelection *DeliveryAcceptanceSelectionStatus `json:"acceptance_selection,omitempty"`
	FinalDelivery       DeliveryFinalStatus                `json:"final_delivery"`
	Reasons             []Reason                           `json:"reasons"`
}

type DeliveryAcceptanceSelectionStatus struct {
	Choice             DeliveryAcceptanceChoice    `json:"choice"`
	Receipt            FileBinding                 `json:"receipt"`
	EffectiveAgreement workspace.DeliveryAgreement `json:"effective_agreement"`
}

type DeliverySourceStatus struct {
	OID   string `json:"oid"`
	Tree  string `json:"tree"`
	Clean bool   `json:"clean"`
	State string `json:"state"`
}

type DeliveryVerificationStatus struct {
	AttemptID     string       `json:"attempt_id"`
	CandidateOID  string       `json:"candidate_oid"`
	CandidateTree string       `json:"candidate_tree"`
	Outcome       string       `json:"outcome"`
	Exit          *int         `json:"exit"`
	Qualification string       `json:"qualification"`
	Current       bool         `json:"current"`
	InputsMatch   bool         `json:"inputs_match"`
	Receipt       *FileBinding `json:"receipt"`
}

type DeliveryQualifiedStatus struct {
	Key          string                 `json:"key"`
	OID          string                 `json:"oid"`
	Tree         string                 `json:"tree"`
	TaskResultID workspace.TaskResultID `json:"task_result_id"`
	Current      bool                   `json:"current"`
}

type DeliveryHumanStatus struct {
	CandidateOID string                    `json:"candidate_oid"`
	TaskResultID workspace.TaskResultID    `json:"task_result_id"`
	HumanQAID    workspace.HumanQARecordID `json:"human_qa_id"`
	Outcome      string                    `json:"outcome"`
	Current      bool                      `json:"current"`
}

type DeliveryAcceptanceStatus struct {
	workspace.DeliveryAcceptanceDecision
	Current          bool   `json:"current"`
	ResponsibleActor string `json:"responsible_actor,omitempty"`
}

type DeliveryFinalStatus struct {
	State        string       `json:"state"`
	CandidateOID string       `json:"candidate_oid,omitempty"`
	Evidence     *FileBinding `json:"evidence"`
}

func deliveryStatus(d Dependencies, s workflowState) *DeliveryStatus {
	if !deliveryRun(s.Request) || s.Result.Delivery == nil {
		return nil
	}
	out := &DeliveryStatus{Source: DeliverySourceStatus{State: "unknown"}, FinalDelivery: DeliveryFinalStatus{State: "not_delivered"}, Reasons: []Reason{}}
	if b := s.Result.Delivery.AcceptanceSelection; b != nil && s.EffectiveDeliveryAgreement != nil {
		_, choice, _, err := readAcceptanceSelection(s, *b)
		if err == nil {
			out.AcceptanceSelection = &DeliveryAcceptanceSelectionStatus{Choice: choice, Receipt: *b, EffectiveAgreement: *s.EffectiveDeliveryAgreement}
		}
	}
	automatic := deliveryEffectiveAgreement(s) != nil && deliveryEffectiveAgreement(s).AutomaticAcceptance()
	if automatic {
		out.Acceptance = &DeliveryAcceptanceStatus{DeliveryAcceptanceDecision: workspace.DeliveryAcceptanceDecision{Mode: "automatic", Outcome: "blocked", Reason: "Automatic acceptance requires a current qualified candidate and successful test evidence.", PolicySHA256: digest(*deliveryEffectiveAgreement(s).Acceptance)}, ResponsibleActor: deliveryEffectiveAgreement(s).Acceptance.ResponsibleActor}
	}
	defer func() {
		// Preserve historical qualification, but never present invalid or newer
		// unresolved verification evidence as the current technical candidate.
		if q := out.QualifiedCandidate; q != nil {
			v := out.Verification
			q.Current = q.Current && v != nil && v.Outcome == "passed" && v.InputsMatch && v.Qualification == "qualified" && v.AttemptID == q.Key
			if a := out.Acceptance; a != nil {
				a.Current = a.Current && q.Current
				if a.Outcome == "pass" && !a.Current {
					a.Outcome, a.Reason = "blocked", "Candidate or required acceptance inputs no longer match the preserved qualification."
				}
			}
		}
	}()
	problem := func(code, detail string) { out.Reasons = append(out.Reasons, Reason{code, detail}) }
	dependenciesMatch := true
	if s.StartSHA256 != nil {
		if err := workflowhandoff.ValidateDeliveryOwnerContinuity(d.Workflow, s.Result.Handoff.Locator); err != nil {
			dependenciesMatch = false
			problem("delivery_required_inputs_invalid", "Required Task dependencies cannot support current qualification: "+err.Error())
		}
	}
	x, err := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref)
	if err != nil {
		problem("delivery_source_unavailable", "Current source could not be observed: "+err.Error())
	} else if x.Ref != s.Observed.Target.Ref || x.GitCommonDir != s.Observed.Target.GitCommonDir {
		problem("delivery_source_stale", "Current source identity differs from the bound Task worktree")
	} else {
		out.Source = DeliverySourceStatus{OID: x.OID, Tree: x.Tree, Clean: x.Clean && len(x.InProgress) == 0, State: "observed"}
	}
	matches := func(oid, tree string) bool {
		return out.Source.State == "observed" && out.Source.Clean && out.Source.OID == oid && out.Source.Tree == tree
	}
	delivery := s.Result.Delivery
	if len(delivery.Candidates) > 0 {
		candidate := delivery.Candidates[len(delivery.Candidates)-1]
		current := matches(candidate.OID, candidate.Tree)
		if a := delivery.Attempt; a != nil && a.Kind == "verification" && (a.State != "recorded" || a.ID != candidate.Key) {
			current = false
		}
		for i := len(delivery.Events) - 1; i >= 0; i-- {
			if event := delivery.Events[i]; event.Kind == "verification" {
				current = current && event.ID == candidate.Key
				break
			}
		}
		out.QualifiedCandidate = &DeliveryQualifiedStatus{candidate.Key, candidate.OID, candidate.Tree, candidate.TaskResult.ID, current}
		qa := candidate.HumanQA
		qaCurrent := true
		var qaErr error
		if d.Workspace.WorkItems == nil {
			qaErr = fmt.Errorf("native work item registry is unavailable")
		} else {
			var registry workspace.WorkItemRegistry
			registry, qaErr = d.Workspace.WorkItems.Snapshot(s.Request.WorkspaceRoot)
			if qaErr == nil {
				qa, qaErr = LatestDeliveryAcceptanceHumanQA(registry.HumanQARecords, candidate.TaskResult, deliveryEffectiveAgreement(s))
			}
		}
		if qaErr == nil {
			qaErr = deliveryHumanEvidence(qa)
		}
		if qaErr != nil {
			qaCurrent = false
			problem("delivery_human_judgment_unknown", "Current native human judgment cannot be established: "+qaErr.Error())
		}
		if qa != nil {
			judgmentCurrent := qaCurrent && matches(qa.ResultOID, qa.ResultTree) && (qa.TaskResultID == candidate.TaskResult.ID || automatic && qa.TaskID == candidate.TaskResult.TaskID && qa.Outcome != "pass")
			out.HumanJudgment = &DeliveryHumanStatus{qa.ResultOID, qa.TaskResultID, qa.ID, qa.Outcome, judgmentCurrent}
		}
		decision, acceptanceErr := workspace.EvaluateDeliveryAcceptance(deliveryEffectiveAgreement(s), candidate.TaskResult, qa)
		out.Acceptance = &DeliveryAcceptanceStatus{DeliveryAcceptanceDecision: decision, Current: current && qaCurrent && acceptanceErr == nil}
		if automatic {
			out.Acceptance.ResponsibleActor = deliveryEffectiveAgreement(s).Acceptance.ResponsibleActor
		}
		if acceptanceErr != nil {
			out.Acceptance.Outcome, out.Acceptance.Reason = "blocked", acceptanceErr.Error()
			problem("delivery_acceptance_invalid", acceptanceErr.Error())
		}
		if qaErr != nil {
			out.Acceptance.Outcome, out.Acceptance.Reason = "blocked", qaErr.Error()
		}
		if delivery.Phase == "completed" || delivery.Phase == "closing" {
			binding := candidate.Integration
			if binding == nil {
				binding = candidate.PullRequest
			}
			out.FinalDelivery = DeliveryFinalStatus{State: "unknown", CandidateOID: candidate.OID, Evidence: binding}
			if binding == nil {
				problem("delivery_completion_evidence_missing", "Completed phase has no observed delivery receipt")
			} else if _, err := workflowBound(*binding, 4<<20); err != nil {
				problem("delivery_completion_evidence_invalid", err.Error())
			} else {
				// Delivery is an observed historical effect. A later input change or
				// human answer changes current acceptance, never the recorded effect.
				completedDecision, err := workspace.EvaluateDeliveryAcceptance(deliveryEffectiveAgreement(s), candidate.TaskResult, candidate.HumanQA)
				if err != nil || completedDecision.Outcome != "pass" {
					problem("delivery_completion_judgment_invalid", "Completion lacks the exact candidate's preserved acceptance")
				} else {
					out.FinalDelivery.State = "delivered"
				}
			}
		}
	}

	var binding *FileBinding
	attempt := delivery.Attempt
	if attempt != nil && attempt.Kind == "verification" {
		v := &DeliveryVerificationStatus{AttemptID: attempt.ID, CandidateOID: attempt.CandidateOID, CandidateTree: attempt.CandidateTree, Outcome: "unknown", Qualification: "pending", Current: matches(attempt.CandidateOID, attempt.CandidateTree)}
		out.Verification = v
		if attempt.State == "recorded" {
			v.Qualification = "not_qualified"
		}
		// Recover only observation here; this never executes or qualifies anything.
		initialBytes, err := readFile(filepath.Join(attempt.Path, "attempt.json"), 256<<10, true)
		var initial deliveryVerificationReceipt
		if err == nil {
			err = decode(initialBytes, 256<<10, &initial)
		}
		if err == nil {
			_, err = deliveryRecoverVerification(s, initial.Review.Locator)
		}
		if err != nil {
			problem("delivery_verification_evidence_unknown", fmt.Sprintf("Attempt %s: %s", attempt.ID, err))
			return out
		}
		b, _, err := deliveryReadBinding(filepath.Join(attempt.Path, "verification.json"), 256<<10)
		if err != nil {
			problem("delivery_verification_evidence_unknown", err.Error())
			return out
		}
		binding = &b
		// Recorded evidence is additionally bound by the immutable event chain.
		if attempt.State == "recorded" {
			found := false
			for _, event := range delivery.Events {
				if event.ID == attempt.ID && event.Kind == "verification" {
					found = event.Binding.SHA256 == b.SHA256
				}
			}
			if !found {
				problem("delivery_verification_evidence_invalid", "Verification receipt differs from its recorded event")
				return out
			}
		}
	} else {
		for i := len(delivery.Events) - 1; i >= 0; i-- {
			event := delivery.Events[i]
			if event.Kind == "verification" {
				binding = &event.Binding
				break
			}
		}
	}
	if binding == nil {
		return out
	}
	if out.Verification == nil {
		out.Verification = &DeliveryVerificationStatus{Outcome: "unknown", Qualification: "not_qualified"}
	}
	v := out.Verification
	raw, err := workflowBound(*binding, 256<<10)
	var receipt deliveryVerificationReceipt
	if err == nil {
		err = decode(raw, 256<<10, &receipt)
	}
	if err == nil && (!deliveryVerificationVersion(receipt) || receipt.RunID != s.Result.RunID || receipt.RequestSHA256 != s.Result.RequestSHA256) {
		err = integrity("verification receipt identity differs")
	}
	if err == nil {
		for _, b := range []FileBinding{receipt.AcceptanceSnapshot, receipt.Stdout, receipt.Stderr} {
			if _, err = workflowBound(b, 4<<20); err != nil {
				break
			}
		}
	}
	if err != nil {
		problem("delivery_verification_evidence_invalid", err.Error())
		return out
	}
	v.AttemptID, v.CandidateOID, v.CandidateTree, v.Receipt = receipt.AttemptID, receipt.CandidateOID, receipt.CandidateTree, binding
	v.Current = matches(receipt.CandidateOID, receipt.CandidateTree)
	v.InputsMatch = v.Current && dependenciesMatch
	// Live inputs determine whether evidence can still qualify this source.
	// They cannot erase a completed command's separately preserved outcome.
	for _, b := range []FileBinding{receipt.Acceptance, receipt.Review} {
		if _, err := workflowBound(b, 4<<20); err != nil {
			v.InputsMatch = false
			problem("delivery_verification_inputs_changed", "Completed verification is preserved; input changed and receipt reuse requires matching inputs: "+err.Error())
		}
	}
	if automatic {
		candidate := DeliveryCandidate{Key: receipt.AttemptID, OID: receipt.CandidateOID, Tree: receipt.CandidateTree, Verification: *binding}
		if err := validateAutomaticVerification(d, s, candidate); err != nil {
			v.InputsMatch = false
			problem("delivery_automatic_acceptance_invalid", err.Error())
		}
	}
	v.Exit = receipt.Exit
	if receipt.Exit != nil {
		v.Outcome = "failed"
		if *receipt.Exit == 0 && receipt.Error == "" {
			v.Outcome = "passed"
		}
	}
	if automatic && receipt.Outcome == "blocked" {
		v.Outcome = "blocked"
	}
	if automatic && (receipt.Outcome == "fail" || receipt.Outcome == "blocked") && v.Current {
		out.Acceptance.Outcome, out.Acceptance.Reason = receipt.Outcome, receipt.Error
		out.Acceptance.TaskResultID = ""
		out.Acceptance.ResultOID, out.Acceptance.ResultTree = receipt.CandidateOID, receipt.CandidateTree
		out.Acceptance.Current = false
	}
	for _, c := range delivery.Candidates {
		if c.Key == receipt.AttemptID && c.OID == receipt.CandidateOID && c.Tree == receipt.CandidateTree && c.Verification.SHA256 == digest(receipt) {
			v.Qualification = "qualified"
			break
		}
	}
	return out
}
