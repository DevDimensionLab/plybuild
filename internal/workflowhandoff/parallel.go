package workflowhandoff

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func deliveryHandoff(o canonicaljson.Object) bool { v := taskSpecVersion(o); return v == 3 || v == 4 }

func validateParallelHandoff(o canonicaljson.Object) error {
	v, found := objectMember(o, "parallel_binding")
	effective, e := parallelEffectiveBasis(o)
	if e != nil {
		return e
	}
	if taskSpecVersion(o) == 4 && deliveryMode(o) == "candidate" && effective == nil {
		return fmt.Errorf("parallel candidate requires effective basis")
	}

	if taskSpecVersion(o) != 4 {
		if found {
			return fmt.Errorf("parallel binding requires Handoff@4")
		}
		return nil
	}
	m, err := exactObject(v, "parallel_binding", "preparation_id", "admission_mode", "source_reconciliation")
	if err != nil {
		return err
	}
	if objectMapString(m, "admission_mode") != workspace.AdmissionParallel || objectMapString(m, "source_reconciliation") != workspace.SourceReconciliationMergeParent || len(objectMapString(m, "preparation_id")) != 68 {
		return fmt.Errorf("invalid exact parallel authority")
	}
	return nil
}

// ParallelPreparation is a validated frozen authority projection, never a new
// grant to an older Handoff. Callers still validate the exact native member.
func ParallelPreparation(raw []byte) (string, error) {
	d, err := decodeHandoffDraft(raw)
	if err != nil {
		return "", err
	}
	if taskSpecVersion(d.Value) != 4 {
		return "", nil
	}
	v, _ := objectMember(d.Value, "parallel_binding")
	b, err := canonicaljson.Marshal(v)
	if err != nil {
		return "", err
	}
	var binding struct {
		PreparationID string `json:"preparation_id"`
	}
	if err = json.Unmarshal(b, &binding); err != nil {
		return "", err
	}
	return binding.PreparationID, nil
}

func parallelEffectiveBasis(o canonicaljson.Object) (*workspace.EffectiveDeliveryBasis, error) {
	v, found := objectMember(o, "effective_basis")
	if !found {
		return nil, nil
	}
	if taskSpecVersion(o) != 4 || deliveryMode(o) != "candidate" {
		return nil, fmt.Errorf("effective basis requires a parallel candidate")
	}
	raw, err := canonicaljson.Marshal(v)
	if err != nil {
		return nil, err
	}
	var basis workspace.EffectiveDeliveryBasis
	if err = json.Unmarshal(raw, &basis); err != nil {
		return nil, err
	}
	if !canonicalEqual(v, bridgeValue(basis)) || basis.Generation < 0 || len(basis.PreparationID) != 68 {
		return nil, fmt.Errorf("invalid effective basis shape")
	}
	return &basis, nil
}
