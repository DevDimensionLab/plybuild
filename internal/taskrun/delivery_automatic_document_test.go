package taskrun

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func automaticDocumentFixture(t *testing.T, requireHuman, fail bool) (deliveryFixture, WorkflowRun, string) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "document fixture", RequireHumanQA: requireHuman, TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a})
	ownerPaneStandin(t, &f)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "documentation candidate")
	script := "#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'documentation candidate'\ntest -z \"${PLY_CANDIDATE_BINARY+x}\"\ntest -n \"$PLY_CANDIDATE_OID\"\ntest -n \"$PLY_CANDIDATE_TREE\"\n"
	if fail {
		script += "exit 7\n"
	}
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "document-only")
	return f, o, review
}

func TestAutomaticAcceptanceDocumentationWithoutBinary(t *testing.T) {
	f, o, review := automaticDocumentFixture(t, false, false)
	o, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, "")
	if err != nil || o.Delivery.Phase != "automatic_acceptance_passed" || len(o.Delivery.Candidates) != 1 {
		t.Fatalf("documentation candidate needs no executable product: %v %+v", err, o.DeliveryStatus)
	}
	c := o.Delivery.Candidates[0]
	raw, err := workflowBound(c.Verification, 4<<20)
	if err != nil {
		t.Fatal(err)
	}
	var receipt deliveryVerificationReceipt
	if err := decode(raw, 4<<20, &receipt); err != nil {
		t.Fatal(err)
	}
	if receipt.Automatic == nil || receipt.Automatic.CandidateBinary != (workflowhandoff.AutomaticBinaryBinding{}) || receipt.Automatic.CandidateBinarySnapshot != (workflowhandoff.AutomaticBinaryBinding{}) {
		t.Fatalf("documentation verification invented an executable: %+v", receipt.Automatic)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(receipt.AcceptanceSnapshot.Locator), "candidate-executable")); !os.IsNotExist(err) {
		t.Fatalf("documentation verifier created a binary snapshot: %v", err)
	}
	for _, defect := range []string{"original only", "snapshot only", "unbound binary environment", "changed policy", "missing candidate OID", "missing executor"} {
		t.Run(defect, func(t *testing.T) {
			v := *receipt.Automatic
			v.Environment = append([]string(nil), v.Environment...)
			binding := workflowhandoff.AutomaticBinaryBinding{Locator: filepath.Join(f.R.WorkspaceRoot, "unexpected-executable"), SHA256: "sha256:" + strings.Repeat("a", 64)}
			switch defect {
			case "original only":
				v.CandidateBinary = binding
			case "snapshot only":
				v.CandidateBinarySnapshot = binding
			case "unbound binary environment":
				v.Environment = append(v.Environment, "PLY_CANDIDATE_BINARY="+binding.Locator)
			case "changed policy":
				v.Policy.TimeoutSeconds++
			case "missing candidate OID":
				for i, entry := range v.Environment {
					if strings.HasPrefix(entry, "PLY_CANDIDATE_OID=") {
						v.Environment = append(v.Environment[:i], v.Environment[i+1:]...)
						break
					}
				}
			case "missing executor":
				v.ExecutorSessionID = ""
			}
			sort.Strings(v.Environment)
			v.InstructionsSHA256 = workflowhandoff.AutomaticInstructionsDigest(v, receipt.Argv, receipt.CWD, receipt.Acceptance.SHA256, receipt.CandidateOID, receipt.CandidateTree)
			if err := workflowhandoff.ValidateAutomaticInstructions(v, f.R.Delivery.Agreement, receipt.Argv, receipt.CWD, receipt.Acceptance.SHA256, receipt.CandidateOID, receipt.CandidateTree); err == nil {
				t.Fatal("documentation mode accepted incomplete or altered verification bindings")
			}
		})
	}
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	cleanupOwnerHandoff(t, f, o, err)
	cleanupOwnerLeaves(t, f)
	closed, err := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 5*time.Second)
	if err != nil || closed.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("documentation delivery failed after owner handoff: %v %+v", err, closed.DeliveryStatus)
	}
	run := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	checks := traceKind(run, "verification")
	if len(checks) != 1 || checks[0].Verification == nil || !checks[0].Verification.InputsBound || run.Coverage != "complete" {
		t.Fatalf("documentation verification lost its exact evidence: %+v %+v", checks, run.Reasons)
	}
	r, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(r.HumanQARecords) != 0 {
		t.Fatalf("documentation delivery invented human QA: %v", err)
	}
}

func TestAutomaticDocumentationRetainsFailureAndHumanGates(t *testing.T) {
	for _, requireHuman := range []bool{false, true} {
		name := "failed script"
		if requireHuman {
			name = "required human judgment"
		}
		t.Run(name, func(t *testing.T) {
			f, o, review := automaticDocumentFixture(t, requireHuman, !requireHuman)
			before := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			result, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, "")
			if !requireHuman && (err == nil || result.DeliveryStatus.Verification.Exit == nil || *result.DeliveryStatus.Verification.Exit != 7) {
				t.Fatalf("failed document check was not preserved: %v %+v", err, result.DeliveryStatus)
			}
			if requireHuman && err != nil {
				t.Fatalf("machine verification did not complete before human gate: %v", err)
			}
			if _, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
				t.Fatal("document verification bypassed its failed test or human gate")
			}
			if got := gitOutput(t, f.Parent, "rev-parse", "HEAD"); got != before {
				t.Fatal("blocked document candidate advanced the target")
			}
		})
	}
}

func TestAutomaticDocumentationCLIWithoutBinary(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	deliveryTestCandidate(t, f, "documentation CLI candidate")
	if err := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'documentation CLI candidate'\ntest -z \"${PLY_CANDIDATE_BINARY+x}\"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "document-cli")
	o = automaticCLIOK(t, binary, f, "execute", "verify", o.RunID, "--context", o.Paths.Context, "--review", review)
	if len(o.Delivery.Candidates) != 1 || o.DeliveryStatus.Acceptance.Outcome != "pass" {
		t.Fatalf("CLI did not qualify the documentation candidate: %+v", o.DeliveryStatus)
	}
	assertAutomaticCLIClosure(t, binary, f, o)
}
