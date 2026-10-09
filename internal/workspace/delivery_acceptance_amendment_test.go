package workspace

import (
	"path/filepath"
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

func newAmendedIntegrationFixture(t *testing.T) (integrationJourneyFixture, TaskIntegrationInput) {
	t.Helper()
	f := newIntegrationJourneyFixtureHistory(t, false, false, fixtureAgreement(DeliveryLocalEpic, "refs/heads/epic"))
	verifier := TaskVerifierResultRecord{VerifierID: "acceptance", Outcome: "passed", Argv: []string{"/bin/sh", "/private/acceptance.sh"}, CWD: f.taskPath, Exit: 0, BoundOIDOrSHA256: f.resultOID}
	if err := f.dependencies.WorkItems.WithLock(f.root, func(s WorkItemStoreSession) error {
		r, err := s.Snapshot()
		if err != nil {
			return err
		}
		findTaskResult(r, f.resultID).VerifierResults = []TaskVerifierResultRecord{verifier}
		r.HumanQARecords = []TaskHumanQARecord{}
		return s.Publish(r)
	}); err != nil {
		t.Fatal(err)
	}
	f.evidence.evidence.VerifierResults = []TaskVerifierResultRecord{verifier}
	f.evidence.evidence.ExpectedVerifierIDs = []string{verifier.VerifierID}
	auth := f.evidence.evidence.DeliveryAuthorization
	amendment := acceptanceAmendmentFixture(auth.Agreement)
	amendment.RunID, amendment.RequestSHA256, amendment.MandateSHA256 = auth.RunID, auth.RequestSHA256, auth.MandateSHA256
	amendment.DecisionLocator = filepath.Join(f.root, "native-acceptance-choice.json")
	effective, err := ApplyDeliveryAcceptanceAmendment(auth.Agreement, amendment)
	if err != nil {
		t.Fatal(err)
	}
	auth.Agreement, auth.AgreementSHA256, auth.AcceptanceAmendment = effective, DeliveryAgreementDigest(effective), &amendment
	in := deliveryIntegrationInput(f)
	in.HumanQARecordID = ""
	in.DeliveryOwner = &DeliveryIntegrationOwner{RunID: auth.RunID, RequestSHA256: auth.RequestSHA256, ActorClaim: "fixture orchestrator", PreparationID: "fixture-preparation"}
	return f, in
}

func TestDeliveryAcceptanceAmendmentNativeIntegrationKeepsOriginalSpec(t *testing.T) {
	f, in := newAmendedIntegrationFixture(t)
	r, err := f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil {
		t.Fatal(err)
	}
	original, err := deliveryAgreementForResult(f.dependencies, f.root, r, *findTaskResult(r, f.resultID))
	if err != nil || original == nil || original.SchemaVersion != 1 || original.Acceptance != nil {
		t.Fatalf("fixture lost original frozen agreement: %+v %v", original, err)
	}
	check, err := CheckTaskIntegration(f.dependencies, in)
	if err != nil || check.Readback.Classification != "ready" || check.Readback.HumanQA != nil || check.Readback.Acceptance == nil || check.Readback.Acceptance.Outcome != "pass" {
		t.Fatalf("explicit amended automatic candidate was blocked: %+v %v", check.Readback, err)
	}
	in.Apply, in.Confirmation = true, integrationDigest(t, check.Readback)
	if _, err = ApplyTaskIntegration(f.dependencies, in); err != nil {
		t.Fatal(err)
	}
	r, err = f.dependencies.WorkItems.Snapshot(f.root)
	if err != nil || len(r.HumanQARecords) != 0 || len(r.IntegrationAuthorities) != 1 || len(r.IntegrationResults) != 1 {
		t.Fatalf("native amended integration did not preserve exact machine history: %v", err)
	}
	auth := r.IntegrationAuthorities[0].Plan.DeliveryAuthorization
	if auth == nil || !contentTypedEqual(auth, in.DeliveryAuthorization) || auth.AcceptanceAmendment.PreviousAgreementSHA256 != DeliveryAgreementDigest(*original) {
		t.Fatal("native published authority lost explicit amendment binding")
	}
	after, err := deliveryAgreementForResult(f.dependencies, f.root, r, *findTaskResult(r, f.resultID))
	if err != nil || !contentTypedEqual(original, after) {
		t.Fatalf("native integration rewrote immutable Spec agreement: %v", err)
	}
}

func TestDeliveryAcceptanceAmendmentNativeGateRequiresExactRevalidatedBinding(t *testing.T) {
	for _, defect := range []string{"omitted amendment", "changed target", "missing native amendment", "different decision", "changed native policy"} {
		t.Run(defect, func(t *testing.T) {
			f, in := newAmendedIntegrationFixture(t)
			auth := *in.DeliveryAuthorization
			amendment := *auth.AcceptanceAmendment
			auth.AcceptanceAmendment, in.DeliveryAuthorization = &amendment, &auth
			switch defect {
			case "omitted amendment":
				auth.AcceptanceAmendment = nil
			case "changed target":
				auth.Agreement.TargetRef = "refs/heads/another-epic"
				auth.AgreementSHA256 = DeliveryAgreementDigest(auth.Agreement)
			case "missing native amendment":
				f.evidence.evidence.DeliveryAuthorization.AcceptanceAmendment = nil
			case "different decision":
				amendment.DecisionSHA256 = "sha256:" + strings.Repeat("3", 64)
			case "changed native policy":
				f.evidence.evidence.DeliveryAuthorization.PermissionConfirmed = false
			}
			git := &countingIntegrationGit{TaskIntegrationGit: f.dependencies.IntegrationGit}
			f.dependencies.IntegrationGit = git
			got, err := CheckTaskIntegration(f.dependencies, in)
			if err == nil && got.Readback.Classification == "ready" {
				t.Fatalf("%s bypassed native amendment validation", defect)
			}
			if git.count() != 0 || runLocalGit(t, filepath.Join(f.wrapper, "epic"), "rev-parse", "HEAD") != f.oid {
				t.Fatal("rejected acceptance amendment changed target")
			}
		})
	}
}
