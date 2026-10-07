package delivery

import (
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type validNativeEvidence struct{ value workspace.TaskHandoffEvidence }

func (r validNativeEvidence) ReadTaskEvidence(workspace.TaskHandoffEvidenceRequest) (workspace.TaskHandoffEvidence, error) {
	return r.value, nil
}

func TestDeliveryPreservesNativeVerifierExpectedExitAndBindingSemantics(t *testing.T) {
	basis := workspace.TaskSpecBasis{TaskID: "task"}
	result := workspace.TaskResultRecord{TechnicalGate: "passed", ReportedOutcome: "complete", ResultOID: strings.Repeat("a", 40), ResultTree: strings.Repeat("b", 40), SourceRef: "refs/heads/task", SourceLocator: "/registered/task", VerifierResults: []workspace.TaskVerifierResultRecord{{VerifierID: "negative-contract", Outcome: "passed", Exit: 1, BoundOIDOrSHA256: "sha256:" + strings.Repeat("c", 64)}}}
	// The native reader has already checked the verifier's declared expected
	// exit and digest binding. Delivery must preserve that result unchanged.
	evidence := workspace.TaskHandoffEvidence{SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: true, TaskSpecValid: true, TaskRequirementsValid: true, EvidenceCoverageValid: true, AcceptedStartOutcome: "started", ReportedOutcome: "complete", TaskSpecBasis: &basis, ExpectedVerifierIDs: []string{"negative-contract"}, VerifierResults: result.VerifierResults, Review: result.Review, ResultOID: result.ResultOID, ResultTree: result.ResultTree, TargetRef: result.SourceRef, TargetWorktree: result.SourceLocator}
	d := taskrun.Dependencies{Workspace: workspace.Dependencies{HandoffEvidence: validNativeEvidence{evidence}}}
	if e := checkEvidence(d, result, basis); e != nil {
		t.Fatalf("native passed verifier was reinterpreted: %v", e)
	}
}
