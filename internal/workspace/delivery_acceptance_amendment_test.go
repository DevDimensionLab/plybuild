package workspace

import (
	"strings"
	"testing"
)

func acceptanceAmendmentFixture(original DeliveryAgreement) DeliveryAcceptanceAmendment {
	return DeliveryAcceptanceAmendment{SchemaVersion: 1, RunID: "wfr_fixture", RequestSHA256: "sha256:" + strings.Repeat("1", 64), MandateSHA256: "sha256:" + strings.Repeat("1", 64), PreviousAgreementSHA256: DeliveryAgreementDigest(original), Policy: DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic fixture orchestrator", TimeoutSeconds: 60}, DecisionLocator: "/workspace/native-acceptance-choice.json", DecisionSHA256: "sha256:" + strings.Repeat("2", 64)}
}

func TestDeliveryAcceptanceAmendmentPreservesFrozenScopeAndEffects(t *testing.T) {
	original := automaticAgreementFixture()
	original.SchemaVersion, original.Acceptance, original.SourceRef = 1, nil, "refs/heads/task"
	before := DeliveryAgreementDigest(original)
	amendment := acceptanceAmendmentFixture(original)
	effective, err := ApplyDeliveryAcceptanceAmendment(original, amendment)
	if err != nil || !effective.AutomaticAcceptance() || DeliveryAgreementDigest(original) != before {
		t.Fatalf("explicit choice changed immutable agreement: %+v %v", effective, err)
	}
	scope := effective
	scope.SchemaVersion, scope.Acceptance = original.SchemaVersion, nil
	if !contentTypedEqual(scope, original) || !contentTypedEqual(effective.AllowedEffects(), original.AllowedEffects()) || effective.StopAfter() != original.StopAfter() {
		t.Fatal("acceptance choice changed scope, effects or stop boundary")
	}
	for name, mutate := range map[string]func(*DeliveryAgreement, *DeliveryAcceptanceAmendment){
		"unbound source": func(a *DeliveryAgreement, _ *DeliveryAcceptanceAmendment) { a.SourceRef = "" },
		"human owner": func(a *DeliveryAgreement, _ *DeliveryAcceptanceAmendment) {
			a.SchemaVersion, a.IntegrationOwner = 2, IntegrationOwnerHuman
		},
		"already automatic": func(a *DeliveryAgreement, _ *DeliveryAcceptanceAmendment) {
			a.SchemaVersion, a.Acceptance = 3, effective.Acceptance
		},
		"branch mode": func(a *DeliveryAgreement, _ *DeliveryAcceptanceAmendment) { a.Mode = DeliveryLocalBranch },
		"main target": func(a *DeliveryAgreement, _ *DeliveryAcceptanceAmendment) { a.TargetRef = "refs/heads/main" },
		"PR mode": func(a *DeliveryAgreement, _ *DeliveryAcceptanceAmendment) {
			a.Mode, a.Remote, a.GitHubRepository, a.TargetWorktree = DeliveryPullRequest, "origin", "fixture/repo", ""
		},
		"schema":          func(_ *DeliveryAgreement, a *DeliveryAcceptanceAmendment) { a.SchemaVersion = 2 },
		"run missing":     func(_ *DeliveryAgreement, a *DeliveryAcceptanceAmendment) { a.RunID = "" },
		"request missing": func(_ *DeliveryAgreement, a *DeliveryAcceptanceAmendment) { a.RequestSHA256 = "" },
		"mandate missing": func(_ *DeliveryAgreement, a *DeliveryAcceptanceAmendment) { a.MandateSHA256 = "" },
		"previous agreement": func(_ *DeliveryAgreement, a *DeliveryAcceptanceAmendment) {
			a.PreviousAgreementSHA256 = a.DecisionSHA256
		},
		"relative decision":       func(_ *DeliveryAgreement, a *DeliveryAcceptanceAmendment) { a.DecisionLocator = "decision.json" },
		"missing decision digest": func(_ *DeliveryAgreement, a *DeliveryAcceptanceAmendment) { a.DecisionSHA256 = "" },
		"invalid policy":          func(_ *DeliveryAgreement, a *DeliveryAcceptanceAmendment) { a.Policy.TimeoutSeconds = 0 },
	} {
		t.Run(name, func(t *testing.T) {
			a, choice := original, amendment
			mutate(&a, &choice)
			if _, err := ApplyDeliveryAcceptanceAmendment(a, choice); err == nil {
				t.Fatal("invalid acceptance amendment was accepted")
			}
		})
	}
}
