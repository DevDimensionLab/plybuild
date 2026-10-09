package workflowhandoff

import (
	"encoding/json"
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func deliveryAcceptanceAmendment(o canonicaljson.Object) (*workspace.DeliveryAcceptanceAmendment, error) {
	v, _ := objectMember(o, "delivery_binding")
	b, _ := v.(canonicaljson.Object)
	v, found := objectMember(b, "acceptance_amendment")
	if !found {
		return nil, nil
	}
	if deliveryMode(o) != "candidate" {
		return nil, fmt.Errorf("an acceptance amendment belongs only to a new candidate; the frozen owner cannot be amended")
	}
	raw, err := canonicaljson.Marshal(v)
	if err != nil {
		return nil, err
	}
	var a workspace.DeliveryAcceptanceAmendment
	if err := json.Unmarshal(raw, &a); err != nil || !canonicalEqual(v, bridgeValue(a)) {
		return nil, fmt.Errorf("acceptance amendment has missing or unknown binding fields")
	}
	return &a, nil
}

func effectiveDeliveryCandidateAgreement(original *workspace.DeliveryAgreement, amendment *workspace.DeliveryAcceptanceAmendment, runID, requestSHA, mandateSHA string) (*workspace.DeliveryAgreement, error) {
	if amendment == nil {
		return original, nil
	}
	if original == nil || amendment.RunID != runID || amendment.RequestSHA256 != requestSHA || amendment.MandateSHA256 != mandateSHA {
		return nil, fmt.Errorf("acceptance amendment differs from the candidate's original native owner and Run")
	}
	effective, err := workspace.ApplyDeliveryAcceptanceAmendment(*original, *amendment)
	if err != nil {
		return nil, err
	}
	return &effective, nil
}

// The frozen agreement remains in delivery_binding.agreement. Only candidate
// verification and acceptance use this separately derived effective agreement.
func effectiveDeliveryAgreement(o canonicaljson.Object) (*workspace.DeliveryAgreement, error) {
	original, err := deliveryAgreement(o)
	if err != nil {
		return nil, err
	}
	amendment, err := deliveryAcceptanceAmendment(o)
	if err != nil {
		return nil, err
	}
	v, _ := objectMember(o, "delivery_binding")
	b, _ := v.(canonicaljson.Object)
	return effectiveDeliveryCandidateAgreement(original, amendment, objectString(b, "workflow_run_id"), objectString(b, "request_sha256"), objectString(b, "parent_handoff_sha256"))
}

func validateDeliveryAmendmentEvidence(d Dependencies, amendment *workspace.DeliveryAcceptanceAmendment) error {
	if amendment == nil {
		return nil
	}
	raw, err := deliveryReadFile(d.Files, amendment.DecisionLocator, 1<<20)
	if err != nil {
		return err
	}
	if digestBytes(raw) != amendment.DecisionSHA256 {
		return fmt.Errorf("preserved acceptance amendment decision bytes differ")
	}
	return nil
}
