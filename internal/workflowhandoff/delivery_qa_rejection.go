package workflowhandoff

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type DeliveryRejectedHumanQA struct {
	AttestationPath string
	DraftPath       string
	PublicationKey  string
}

// ValidateRejectedDeliveryHumanQA reads preserved staging to prove the bounded
// legacy +00:00 timestamp rejection. It never repairs evidence or publishes QA.
// The caller must separately prove that no native QA record was published.
func ValidateRejectedDeliveryHumanQA(d Dependencies, taskID string, result workspace.TaskResultRecord, outcome, evidencePath string) (DeliveryRejectedHumanQA, error) {
	var empty DeliveryRejectedHumanQA
	raw, err := deliveryReadFile(d.Files, evidencePath, 256<<10)
	if err != nil {
		return empty, err
	}
	value, err := canonicaljson.DecodeStrict(raw)
	if err != nil {
		return empty, err
	}
	object, ok := value.(canonicaljson.Object)
	if !ok {
		return empty, fmt.Errorf("human QA rejection requires an exact attestation object")
	}
	originalTimes := map[string]string{}
	changed := false
	for i := range object {
		if object[i].Name != "started_at_utc" && object[i].Name != "completed_at_utc" {
			continue
		}
		original, ok := object[i].Value.(string)
		if !ok {
			return empty, fmt.Errorf("human QA rejection requires string timestamps")
		}
		originalTimes[object[i].Name] = original
		switch {
		case strings.HasSuffix(original, "+00:00"):
			object[i].Value = strings.TrimSuffix(original, "+00:00") + "Z"
			changed = true
		case strings.HasSuffix(original, "Z"):
		default:
			return empty, fmt.Errorf("human QA rejection is not the known +00:00 UTC spelling mismatch")
		}
	}
	if !changed {
		return empty, fmt.Errorf("human QA rejection has no +00:00 UTC spelling mismatch")
	}
	normalized, err := canonicaljson.Marshal(object)
	if err != nil {
		return empty, err
	}
	corrected, err := ValidateDeliveryHumanAttestation(normalized, taskID, result, outcome)
	if err != nil {
		return empty, fmt.Errorf("human QA rejection contains other invalid provenance: %w", err)
	}
	legacy := corrected
	legacy.StartedAtUTC = originalTimes["started_at_utc"]
	legacy.CompletedAtUTC = originalTimes["completed_at_utc"]
	control, err := d.Store.ReadByLocator(result.HandoffLocator)
	if err != nil {
		return empty, err
	}
	if taskSpecVersion(control.Handoff.Value) != 3 || deliveryMode(control.Handoff.Value) != "candidate" || control.Handoff.SHA256 != result.HandoffSHA256 {
		return empty, fmt.Errorf("human QA rejection requires the exact qualified delivery candidate")
	}
	area := filepath.Join(control.Handoff.ReplyRoot, "staging", "human-qa-"+digestBytes(raw)[7:])
	proof := DeliveryRejectedHumanQA{AttestationPath: filepath.Join(area, "attestation.json"), DraftPath: filepath.Join(area, "qa-draft.json"), PublicationKey: "delivery-qa/" + string(result.ID) + "/" + digestBytes(raw)[7:]}
	staged, err := deliveryReadFile(d.Files, proof.AttestationPath, 256<<10)
	if err != nil {
		return empty, err
	}
	if !bytes.Equal(staged, raw) {
		return empty, fmt.Errorf("staged human attestation differs from the original evidence")
	}
	expected, err := deliveryHumanQADraft(result, legacy, proof.AttestationPath, raw)
	if err != nil {
		return empty, err
	}
	draft, err := deliveryReadFile(d.Files, proof.DraftPath, 256<<10)
	if err != nil {
		return empty, err
	}
	if !bytes.Equal(draft, expected) {
		return empty, fmt.Errorf("staged human QA draft differs from its original derivation")
	}
	if workspace.ValidateTaskHumanQADraftBytes(draft) == nil {
		return empty, fmt.Errorf("staged human QA draft has no native schema rejection")
	}
	correctedDraft, err := deliveryHumanQADraft(result, corrected, proof.AttestationPath, raw)
	if err != nil {
		return empty, err
	}
	if err = workspace.ValidateTaskHumanQADraftBytes(correctedDraft); err != nil {
		return empty, fmt.Errorf("staged human QA draft has another native schema rejection: %w", err)
	}
	return proof, nil
}
