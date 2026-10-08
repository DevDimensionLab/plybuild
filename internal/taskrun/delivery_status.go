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
	Source             DeliverySourceStatus        `json:"source"`
	Verification       *DeliveryVerificationStatus `json:"verification"`
	QualifiedCandidate *DeliveryQualifiedStatus    `json:"qualified_candidate"`
	HumanJudgment      *DeliveryHumanStatus        `json:"human_judgment"`
	FinalDelivery      DeliveryFinalStatus         `json:"final_delivery"`
	Reasons            []Reason                    `json:"reasons"`
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
	defer func() {
		// Preserve historical qualification, but never present invalid or newer
		// unresolved verification evidence as the current technical candidate.
		if q := out.QualifiedCandidate; q != nil {
			v := out.Verification
			q.Current = q.Current && v != nil && v.Outcome == "passed" && v.InputsMatch && v.Qualification == "qualified" && v.AttemptID == q.Key
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
		if qa := candidate.HumanQA; qa != nil {
			out.HumanJudgment = &DeliveryHumanStatus{qa.ResultOID, qa.TaskResultID, qa.ID, qa.Outcome, matches(qa.ResultOID, qa.ResultTree) && qa.TaskResultID == candidate.TaskResult.ID}
		}
		if delivery.Phase == "completed" {
			binding := candidate.Integration
			if binding == nil {
				binding = candidate.PullRequest
			}
			out.FinalDelivery = DeliveryFinalStatus{State: "unknown", CandidateOID: candidate.OID, Evidence: binding}
			if binding == nil {
				problem("delivery_completion_evidence_missing", "Completed phase has no observed delivery receipt")
			} else if _, err := workflowBound(*binding, 4<<20); err != nil {
				problem("delivery_completion_evidence_invalid", err.Error())
			} else if candidate.HumanQA == nil || candidate.HumanQA.Outcome != "pass" || candidate.HumanQA.TaskResultID != candidate.TaskResult.ID || candidate.HumanQA.ResultOID != candidate.OID || candidate.HumanQA.ResultTree != candidate.Tree {
				problem("delivery_completion_judgment_invalid", "Completion lacks the exact candidate's preserved human pass")
			} else {
				out.FinalDelivery.State = "delivered"
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
	if err == nil && (receipt.Kind != "PlyDeliveryVerification@1" || receipt.SchemaVersion != 1 || receipt.RunID != s.Result.RunID || receipt.RequestSHA256 != s.Result.RequestSHA256) {
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
	v.Exit = receipt.Exit
	if receipt.Exit != nil {
		v.Outcome = "failed"
		if *receipt.Exit == 0 && receipt.Error == "" {
			v.Outcome = "passed"
		}
	}
	for _, c := range delivery.Candidates {
		if c.Key == receipt.AttemptID && c.OID == receipt.CandidateOID && c.Tree == receipt.CandidateTree && c.Verification.SHA256 == digest(receipt) {
			v.Qualification = "qualified"
			break
		}
	}
	return out
}
