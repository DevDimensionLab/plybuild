package taskrun

import "github.com/devdimensionlab/plybuild/internal/workspace"

func deliveryCandidateQualified(ds *DeliveryState) bool {
	if ds == nil || len(ds.Candidates) == 0 {
		return false
	}
	c := ds.Candidates[len(ds.Candidates)-1]
	if ds.Attempt != nil && ds.Attempt.State == "attempted" {
		return false
	}
	return c.TaskResult.VerificationReceipt != nil && c.TaskResult.VerificationReceipt.SHA256 == c.Verification.SHA256
}

func LatestDeliveryHumanQA(records []workspace.TaskHumanQARecord, resultID workspace.TaskResultID, oid, tree string) (*workspace.TaskHumanQARecord, error) {
	return workspace.LatestTaskHumanQA(records, resultID, oid, tree)
}

// Automatic requalification cannot erase a human rejection of unchanged
// candidate bytes. The returned QA keeps its original TaskResult identity.
func LatestDeliveryAcceptanceHumanQA(records []workspace.TaskHumanQARecord, result workspace.TaskResultRecord, agreement *workspace.DeliveryAgreement, parallel ...bool) (*workspace.TaskHumanQARecord, error) {
	if len(parallel) > 0 && parallel[0] {
		return workspace.LatestParallelTaskHumanQA(records, result)
	}
	if agreement != nil && agreement.AutomaticAcceptance() {
		return workspace.LatestTaskCandidateHumanQA(records, result)
	}
	return LatestDeliveryHumanQA(records, result.ID, result.ResultOID, result.ResultTree)
}

func taskResultEvidenceRequest(r workspace.TaskResultRecord) workspace.TaskHandoffEvidenceRequest {
	return workspace.TaskHandoffEvidenceRequest{ActivityID: r.ActivityID, RunID: r.RunID, HandoffID: r.HandoffID, HandoffLocator: r.HandoffLocator, HandoffSHA256: r.HandoffSHA256, StartReceiptID: r.StartReceiptID, StartReceiptLocator: r.StartReceiptLocator, StartReceiptSHA256: r.StartReceiptSHA256, TerminalResultID: r.TerminalResultID, TerminalResultLocator: r.TerminalResultLocator, TerminalResultSHA256: r.TerminalResultSHA256, InspectionSHA256: r.InspectionSHA256}
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
