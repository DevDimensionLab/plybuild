package workspace

import "fmt"

// DeliveryAcceptancePolicy is an explicit choice in a frozen v3 agreement.
// It grants no runtime permissions and is restricted to local Epic delivery.
type DeliveryAcceptancePolicy struct {
	SchemaVersion    int    `yaml:"schema_version" json:"schema_version"`
	Mode             string `yaml:"mode" json:"mode"`
	ResponsibleActor string `yaml:"responsible_actor" json:"responsible_actor"`
	RequireHumanQA   bool   `yaml:"require_human_qa" json:"require_human_qa"`
	TimeoutSeconds   int    `yaml:"timeout_seconds" json:"timeout_seconds"`
}

// DeliveryAcceptanceDecision names the source of acceptance without presenting
// a machine result as a human judgment. Native evidence must be revalidated by
// callers before this decision can authorize an effect.
type DeliveryAcceptanceDecision struct {
	Mode            string          `yaml:"mode" json:"mode"`
	Outcome         string          `yaml:"outcome" json:"outcome"`
	Reason          string          `yaml:"reason" json:"reason"`
	TaskResultID    TaskResultID    `yaml:"task_result_id" json:"task_result_id"`
	ResultOID       string          `yaml:"result_oid" json:"result_oid"`
	ResultTree      string          `yaml:"result_tree" json:"result_tree"`
	PolicySHA256    string          `yaml:"policy_sha256,omitempty" json:"policy_sha256,omitempty"`
	EvidenceSHA256  string          `yaml:"evidence_sha256,omitempty" json:"evidence_sha256,omitempty"`
	HumanQARecordID HumanQARecordID `yaml:"human_qa_record_id,omitempty" json:"human_qa_record_id,omitempty"`
}

func validateDeliveryAcceptancePolicy(p DeliveryAcceptancePolicy) error {
	if p.SchemaVersion != 1 || p.Mode != "automatic" || !validTaskText(p.ResponsibleActor, 1, 256) || p.TimeoutSeconds < 1 || p.TimeoutSeconds > 86400 {
		return fmt.Errorf("automatic acceptance requires schema_version 1, mode automatic, responsible_actor and timeout_seconds 1..86400")
	}
	return nil
}

func (a DeliveryAgreement) AutomaticAcceptance() bool {
	return a.SchemaVersion == 3 && a.Mode == DeliveryLocalEpic && a.Acceptance != nil && a.Acceptance.Mode == "automatic"
}

// EvaluateDeliveryAcceptance is the shared decision for preview, status and
// execution. Automatic callers select latestQA with LatestTaskCandidateHumanQA;
// historical callers use LatestTaskHumanQA. A negative answer always vetoes the
// current candidate without changing the identity of the original human record.
func EvaluateDeliveryAcceptance(agreement *DeliveryAgreement, result TaskResultRecord, latestQA *TaskHumanQARecord) (DeliveryAcceptanceDecision, error) {
	out := DeliveryAcceptanceDecision{Mode: "human", Outcome: "blocked", Reason: "The exact candidate requires an actual human pass.", TaskResultID: result.ID, ResultOID: result.ResultOID, ResultTree: result.ResultTree}
	if agreement != nil {
		if err := ValidateDeliveryAgreement(*agreement); err != nil {
			return out, err
		}
		if agreement.AutomaticAcceptance() {
			out.Mode = "automatic"
			out.PolicySHA256 = queueDigest(*agreement.Acceptance)
			out.EvidenceSHA256 = result.TerminalResultSHA256
		}
	}
	if latestQA != nil {
		sameResult := latestQA.TaskResultID == result.ID
		if latestQA.ResultOID != result.ResultOID || latestQA.ResultTree != result.ResultTree || !sameResult && (out.Mode != "automatic" || latestQA.TaskID != result.TaskID) {
			return out, fmt.Errorf("human QA differs from the exact acceptance candidate")
		}
		if !sameResult && latestQA.Outcome == "pass" {
			// An older positive answer still belongs to its original TaskResult.
			// Automatic verification cannot reattribute it to a replacement result.
			latestQA = nil
		}
	}
	if latestQA != nil {
		out.HumanQARecordID = latestQA.ID
		if latestQA.Outcome != "pass" {
			if latestQA.Outcome != "fail" && latestQA.Outcome != "blocked" {
				return out, fmt.Errorf("candidate human QA has an unsupported outcome")
			}
			out.Outcome = latestQA.Outcome
			out.Reason = "The latest exact candidate human answer is " + latestQA.Outcome + "."
			return out, nil
		}
	}
	if out.Mode == "automatic" {
		if (result.TechnicalGate != "passed" && result.TechnicalGate != "good_enough_with_known_debt") || result.ReportedOutcome != "complete" || len(result.VerifierResults) == 0 || len(result.Review.OpenActionableFindings) != 0 || !digestPattern.MatchString(result.TerminalResultSHA256) {
			out.Reason = "Automatic acceptance requires preserved passing candidate verification and review without open actionable findings."
			return out, nil
		}
		for _, v := range result.VerifierResults {
			if v.Outcome != "passed" || v.Exit != 0 || v.BoundOIDOrSHA256 != result.ResultOID || len(v.Argv) == 0 || !validAbsoluteCleanPath(v.CWD) {
				out.Reason = "Automatic acceptance lacks a passing verifier invocation for this exact candidate."
				return out, nil
			}
		}
		if agreement.Acceptance.RequireHumanQA && latestQA == nil {
			out.Reason = "Automatic verification passed; the agreement also requires an actual human pass."
			return out, nil
		}
		out.Outcome, out.Reason = "pass", "The exact candidate passed its preserved automatic verification and review."
		return out, nil
	}
	if latestQA != nil {
		out.Outcome, out.Reason = "pass", "The exact candidate has an actual human pass."
	}
	return out, nil
}

func automaticIntegrationAuthorization(a *DeliveryAuthorization) bool {
	return a != nil && a.Agreement.AutomaticAcceptance() && ValidateDeliveryAuthorization(*a) == nil
}

func optionalHumanQA(q TaskHumanQARecord) *TaskHumanQARecord {
	if q.ID == "" {
		return nil
	}
	return &q
}

func integrationAcceptance(a *DeliveryAuthorization, result TaskResultRecord, qa TaskHumanQARecord) (DeliveryAcceptanceDecision, error) {
	var agreement *DeliveryAgreement
	if a != nil {
		agreement = &a.Agreement
	}
	return EvaluateDeliveryAcceptance(agreement, result, optionalHumanQA(qa))
}
