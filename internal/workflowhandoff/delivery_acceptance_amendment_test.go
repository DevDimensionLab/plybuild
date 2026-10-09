package workflowhandoff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func candidateAmendmentFixture() (workspace.DeliveryAgreement, workspace.DeliveryAcceptanceAmendment, canonicaljson.Object) {
	original := workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic, ProjectID: "ply", RepoID: "ply", EpicID: "epic", SourceRef: "refs/heads/task", TargetRef: "refs/heads/epic", TargetWorktree: "/workspace/epic"}
	a := workspace.DeliveryAcceptanceAmendment{SchemaVersion: 1, RunID: "fixture-run", RequestSHA256: digestBytes([]byte("request")), MandateSHA256: digestBytes([]byte("mandate")), PreviousAgreementSHA256: workspace.DeliveryAgreementDigest(original), Policy: workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic fixture orchestrator", TimeoutSeconds: 60}, DecisionLocator: "/workspace/acceptance-selection.json", DecisionSHA256: digestBytes([]byte("fixture decision"))}
	value := deliveryObject(map[string]any{"schema_version": 3, "delivery_binding": map[string]any{"mode": "candidate", "agreement": original, "workflow_run_id": a.RunID, "request_sha256": a.RequestSHA256, "parent_handoff_sha256": a.MandateSHA256, "acceptance_amendment": a}})
	return original, a, value
}

func TestDeliveryAcceptanceAmendmentCandidateRetainsFrozenAgreement(t *testing.T) {
	original, _, value := candidateAmendmentFixture()
	frozen, err := deliveryAgreement(value)
	if err != nil || frozen == nil || !canonicalEqual(bridgeValue(*frozen), bridgeValue(original)) {
		t.Fatalf("candidate rewrote its frozen agreement: %+v %v", frozen, err)
	}
	effective, err := effectiveDeliveryAgreement(value)
	if err != nil || effective == nil || !effective.AutomaticAcceptance() || !canonicalEqual(bridgeValue(effective.AllowedEffects()), bridgeValue(original.AllowedEffects())) {
		t.Fatalf("exact amendment was not derived: %+v %v", effective, err)
	}
	if workspace.DeliveryAgreementDigest(*frozen) != workspace.DeliveryAgreementDigest(original) {
		t.Fatal("deriving effective policy mutated the original")
	}
	// Both the owner mandate and the execution Spec stay at their original v1.
	if err := validateDeliverySpecAgreement(value, deliveryObject(map[string]any{"delivery": original})); err != nil {
		t.Fatal(err)
	}
	if err := validateDeliverySpecAgreement(value, deliveryObject(map[string]any{"delivery": effective})); err == nil {
		t.Fatal("amended candidate accepted rewritten execution Spec")
	}
}

func TestDeliveryAcceptanceAmendmentCandidateRequiresExactNativeBinding(t *testing.T) {
	for name, mutate := range map[string]func(canonicaljson.Object) canonicaljson.Object{
		"owner amendment": func(b canonicaljson.Object) canonicaljson.Object { return replaceObjectMember(b, "mode", "owner") },
		"different Run": func(b canonicaljson.Object) canonicaljson.Object {
			return replaceObjectMember(b, "workflow_run_id", "other-run")
		},
		"different request": func(b canonicaljson.Object) canonicaljson.Object {
			return replaceObjectMember(b, "request_sha256", digestBytes([]byte("different request")))
		},
		"different mandate": func(b canonicaljson.Object) canonicaljson.Object {
			return replaceObjectMember(b, "parent_handoff_sha256", digestBytes([]byte("different mandate")))
		},
		"unknown field": func(b canonicaljson.Object) canonicaljson.Object {
			v, _ := objectMember(b, "acceptance_amendment")
			return replaceObjectMember(b, "acceptance_amendment", append(v.(canonicaljson.Object), canonicaljson.Member{Name: "extra", Value: "not authorized"}))
		},
		"missing field": func(b canonicaljson.Object) canonicaljson.Object {
			v, _ := objectMember(b, "acceptance_amendment")
			members := canonicaljson.Object{}
			for _, member := range v.(canonicaljson.Object) {
				if member.Name != "decision_sha256" {
					members = append(members, member)
				}
			}
			return replaceObjectMember(b, "acceptance_amendment", members)
		},
		"null amendment": func(b canonicaljson.Object) canonicaljson.Object {
			return replaceObjectMember(b, "acceptance_amendment", nil)
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, _, value := candidateAmendmentFixture()
			b, _ := objectMember(value, "delivery_binding")
			value = replaceObjectMember(value, "delivery_binding", mutate(b.(canonicaljson.Object)))
			if _, err := effectiveDeliveryAgreement(value); err == nil {
				t.Fatal("candidate accepted a mismatched native acceptance binding")
			}
		})
	}
}

func TestDeliveryAcceptanceAmendmentDecisionEvidenceIsPhysicalAndHashBound(t *testing.T) {
	_, choice, _ := candidateAmendmentFixture()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	choice.DecisionLocator = filepath.Join(root, "decision.json")
	if err := os.WriteFile(choice.DecisionLocator, []byte("fixture decision"), 0600); err != nil {
		t.Fatal(err)
	}
	d := SystemDependencies()
	if err := validateDeliveryAmendmentEvidence(d, &choice); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(choice.DecisionLocator, []byte("changed decision"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := validateDeliveryAmendmentEvidence(d, &choice); err == nil {
		t.Fatal("changed decision bytes were accepted")
	}
	if err := os.Remove(choice.DecisionLocator); err != nil {
		t.Fatal(err)
	}
	if err := validateDeliveryAmendmentEvidence(d, &choice); err == nil {
		t.Fatal("missing decision was accepted")
	}
	actual := filepath.Join(filepath.Dir(choice.DecisionLocator), "actual.json")
	if err := os.WriteFile(actual, []byte("fixture decision"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(actual, choice.DecisionLocator); err != nil {
		t.Fatal(err)
	}
	if err := validateDeliveryAmendmentEvidence(d, &choice); err == nil {
		t.Fatal("symlink replaced native decision identity")
	}
}

func TestDeliveryAcceptanceAmendmentDecisionRequiresExistingSandboxReadAuthority(t *testing.T) {
	f := prepareDeliveryFixture(t)
	snapshot, err := f.d.Store.ReadByLocator(f.locator)
	if err != nil {
		t.Fatal(err)
	}
	_, choice, _ := candidateAmendmentFixture()
	outside, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	choice.DecisionLocator = filepath.Join(outside, "native-choice.json")
	b, _ := objectMember(snapshot.Handoff.Value, "delivery_binding")
	binding := replaceObjectMember(b.(canonicaljson.Object), "mode", "candidate")
	binding = append(binding, canonicaljson.Member{Name: "acceptance_amendment", Value: bridgeValue(choice)})
	snapshot.Handoff.Value = replaceObjectMember(snapshot.Handoff.Value, "delivery_binding", binding)
	for _, outsideAllowed := range []bool{false, true} {
		roots := []string{f.root}
		if outsideAllowed {
			roots = append(roots, outside)
		}
		raw, _ := json.Marshal(map[string]any{"read_roots": roots, "write_roots": []string{f.root}, "temp_root": filepath.Join(f.root, "recipient-temp"), "matches_contract": true})
		v, err := canonicaljson.DecodeStrict(raw)
		if err == nil {
			err = validateSandboxContract(f.d.Files, snapshot, v.(canonicaljson.Object))
		}
		if outsideAllowed && err != nil {
			t.Fatalf("explicit existing read authority was rejected: %v", err)
		}
		if !outsideAllowed && (err == nil || !strings.Contains(err.Error(), "sandbox read roots do not cover "+choice.DecisionLocator)) {
			t.Fatalf("decision locator expanded accepted read authority: %v", err)
		}
	}
}
