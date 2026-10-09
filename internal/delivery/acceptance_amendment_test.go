package delivery

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestDeliveryManifestBindsEffectiveAcceptanceAndNativeAmendment(t *testing.T) {
	policy := workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", TimeoutSeconds: 30}
	amendment := &workspace.DeliveryAcceptanceAmendment{SchemaVersion: 1, RunID: "fixture-run", RequestSHA256: "fixture-request", MandateSHA256: "fixture-mandate", PreviousAgreementSHA256: "fixture-original", Policy: policy, DecisionLocator: "/fixture/preserved-choice.json", DecisionSHA256: "fixture-choice"}
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, SourceRef: "refs/heads/task", TargetRef: "refs/heads/epic", TargetWorktree: "/fixture/epic", Acceptance: &policy}
	m := Manifest{Agreement: agreement, AcceptanceAmendment: amendment, RequestSHA256: "fixture-request", MandateSHA256: "fixture-mandate", PreparationID: "fixture-preparation", ExpectedParentOID: "fixture-parent"}
	validated := taskrun.DeliveryAuthority{Agreement: agreement, RequestSHA256: m.RequestSHA256, MandateSHA256: m.MandateSHA256, PreparationID: m.PreparationID, ExpectedParentOID: m.ExpectedParentOID, Authorization: &workspace.DeliveryAuthorization{Agreement: agreement, AcceptanceAmendment: amendment}}
	if err := manifestAuthorityMatches(m, validated); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		change func(*taskrun.DeliveryAuthority)
	}{
		{"missing native authority", func(a *taskrun.DeliveryAuthority) { a.Authorization = nil }},
		{"missing amendment", func(a *taskrun.DeliveryAuthority) { a.Authorization.AcceptanceAmendment = nil }},
		{"different choice origin", func(a *taskrun.DeliveryAuthority) {
			a.Authorization.AcceptanceAmendment.DecisionLocator = "/fixture/other-choice.json"
		}},
		{"different choice bytes", func(a *taskrun.DeliveryAuthority) {
			a.Authorization.AcceptanceAmendment.DecisionSHA256 = "other-choice"
		}},
		{"different previous agreement", func(a *taskrun.DeliveryAuthority) {
			a.Authorization.AcceptanceAmendment.PreviousAgreementSHA256 = "other-original"
		}},
		{"different policy", func(a *taskrun.DeliveryAuthority) { a.Agreement.Acceptance.RequireHumanQA = true }},
		{"different target", func(a *taskrun.DeliveryAuthority) { a.Agreement.TargetRef = "refs/heads/main" }},
		{"original policy substituted", func(a *taskrun.DeliveryAuthority) { a.Agreement.SchemaVersion, a.Agreement.Acceptance = 1, nil }},
	} {
		t.Run(test.name, func(t *testing.T) {
			var changed taskrun.DeliveryAuthority
			raw, err := json.Marshal(validated)
			if err != nil || json.Unmarshal(raw, &changed) != nil {
				t.Fatal("fixture authority could not be copied")
			}
			test.change(&changed)
			if err := manifestAuthorityMatches(m, changed); err == nil {
				t.Fatal("Delivery accepted a changed native policy choice or scope")
			}
		})
	}
}

func TestDeliveryHistoricalManifestDoesNotAcquireAcceptanceAmendment(t *testing.T) {
	agreement := workspace.DeliveryAgreement{SchemaVersion: 1, Mode: workspace.DeliveryLocalEpic, SourceRef: "refs/heads/task", TargetRef: "refs/heads/epic", TargetWorktree: "/fixture/epic"}
	m := Manifest{Agreement: agreement}
	a := taskrun.DeliveryAuthority{Agreement: agreement, Authorization: &workspace.DeliveryAuthorization{Agreement: agreement}}
	if err := manifestAuthorityMatches(m, a); err != nil {
		t.Fatalf("historical authority acquired an amendment requirement: %v", err)
	}
	raw, err := canonical(m)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte(`"acceptance_amendment"`)) {
		t.Fatal("historical manifest digest changed through an absent amendment field")
	}
}
