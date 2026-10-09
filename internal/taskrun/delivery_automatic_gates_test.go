package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func autoGateFixture(t *testing.T, requireHuman bool, owner string) (deliveryFixture, WorkflowRun, string, string) {
	t.Helper()
	agreement := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, IntegrationOwner: owner, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "fixture orchestrator", RequireHumanQA: requireHuman, TimeoutSeconds: 30}}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &agreement})
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "automatically accepted fixture candidate")
	binary := filepath.Join(f.R.WorkspaceRoot, "candidate-executable")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nprintf 'fixture candidate\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	script := "#!/bin/sh\nset -eu\ntest \"$(cat README.md)\" = 'automatically accepted fixture candidate'\ntest \"$(\"$PLY_CANDIDATE_BINARY\")\" = 'fixture candidate'\ntest -n \"$PLY_CANDIDATE_OID\"\ntest -n \"$PLY_CANDIDATE_TREE\"\n"
	if err := os.WriteFile(f.AcceptancePath, []byte(script), 0600); err != nil {
		t.Fatal(err)
	}
	review := deliveryTestReview(t, f, "automatic-gate")
	var err error
	o, err = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
	if err != nil {
		t.Fatalf("qualify automatic fixture: %v", err)
	}
	return f, o, binary, review
}

func autoGateInput(a DeliveryAuthority) workspace.TaskIntegrationInput {
	return workspace.TaskIntegrationInput{TaskID: a.Candidate.TaskResult.TaskID, TaskResultID: a.Candidate.TaskResult.ID, ExpectedResultOID: a.Candidate.OID, ExpectedParentOID: a.ExpectedParentOID, DeliveryAuthorization: a.Authorization, DeliveryOwner: &workspace.DeliveryIntegrationOwner{RunID: a.RunID, RequestSHA256: a.RequestSHA256, ActorClaim: a.OwnerClaim, PreparationID: a.PreparationID}}
}

func TestAutomaticDeliveryIntegratesWithoutHumanRecord(t *testing.T) {
	f, o, _, _ := autoGateFixture(t, false, "")
	c := o.Delivery.Candidates[0]
	if o.Delivery.Phase != "automatic_acceptance_passed" || c.HumanQA != nil || o.DeliveryStatus.HumanJudgment != nil || o.DeliveryStatus.Acceptance == nil || o.DeliveryStatus.Acceptance.Outcome != "pass" || !o.DeliveryStatus.Acceptance.Current {
		t.Fatalf("automatic qualification is not distinct from human judgment: %+v", o.DeliveryStatus)
	}
	a, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult, *f.R.Delivery.Agreement)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := workspace.CheckTaskIntegration(f.D.Workspace, autoGateInput(a))
	if err != nil || preview.Readback.Classification != "ready" {
		t.Fatalf("native preview requires fabricated human QA: %+v %v", preview.Readback, err)
	}
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil || o.Delivery.Phase != "completed" || o.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("automatic local delivery did not complete: %+v %v", o.DeliveryStatus, err)
	}
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.HumanQARecords) != 0 || o.Delivery.Candidates[0].HumanQA != nil {
		t.Fatal("automatic delivery manufactured a human record")
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("automatic delivery did not update the agreed local Epic")
	}
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || closeout.State != "complete" || !closeout.LifecycleCompleted || closeout.ResourceState != "kept" || closeout.WorktreeRemoved || closeout.BranchRemoved {
		t.Fatalf("automatic delivery did not complete retained native closeout: %+v %v", closeout, err)
	}
	observed, handled, err := ObserveAutomaticLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
	if err != nil || !handled || observed.Closeout == nil || !equal(observed.Closeout, closeout) {
		t.Fatalf("automatic observation omitted already completed native closeout: %+v %v", observed.Closeout, err)
	}
	lifecycle, err := workspace.ReadWorkItemLifecycle(f.D.Workspace, f.R.WorkspaceRoot)
	if err != nil || lifecycle.Task(c.TaskResult.TaskID) != workspace.LifecycleCompleted {
		t.Fatalf("automatic delivery did not complete the Task lifecycle: %+v %v", lifecycle, err)
	}
	base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || deliveryCLILatestBase(base).OID != c.OID {
		t.Fatalf("automatic local integration left the native Epic base stale: %+v %v", base, err)
	}
	before := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	if _, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err != nil {
		t.Fatalf("automatic completed delivery cannot be observed again: %v", err)
	}
	if got := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic"); got != before {
		t.Fatal("automatic delivery retry repeated the Git effect")
	}
}

func TestAutomaticDeliveryCloseoutInterruptionResumesWithoutIntegration(t *testing.T) {
	f, o, _, _ := autoGateFixture(t, false, "")
	c := o.Delivery.Candidates[0]
	d := f.D
	d.Fault = func(point string) error {
		if point == "automatic_closeout_after_release" {
			return errors.New("fixture interruption before native retained closeout")
		}
		return nil
	}
	interrupted, err := WorkflowDeliveryIntegrate(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err == nil || interrupted.Delivery.Phase != "closing" || interrupted.DeliveryStatus.FinalDelivery.State != "delivered" || interrupted.Delivery.OwnershipRelease == nil {
		t.Fatalf("closeout interruption hid integration or claimed completion: %+v %v", interrupted.Delivery, err)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("closeout fixture did not reach the actual Git integration")
	}
	before := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	done, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil || done.Delivery.Phase != "completed" {
		t.Fatalf("released automatic owner could not resume only closeout: %+v %v", done.Delivery, err)
	}
	if got := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic"); got != before {
		t.Fatal("native closeout resume repeated Git integration")
	}
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || closeout.State != "complete" || closeout.ResourceState != "kept" {
		t.Fatalf("resumed automatic delivery omitted retained closeout: %+v %v", closeout, err)
	}
}

func TestAutomaticDeliveryCLICloseoutRecoveryFromReleasedOwner(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, _, _ := autoGateFixture(t, false, "")
	c := o.Delivery.Candidates[0]
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	registered := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "automatic-closeout-recovery"))
	id := deliveryCLIString(t, registered, "id")
	d := f.D
	d.Fault = func(point string) error {
		if point == "automatic_closeout_after_release" {
			return errors.New("fixture interruption before native retained closeout")
		}
		return nil
	}
	if _, err := WorkflowDeliveryIntegrate(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
		t.Fatal("fixture did not interrupt before automatic native closeout")
	}
	before := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	checked := deliveryCLIOK(t, binary, cwd, "check", id)
	if deliveryCLIString(t, checked, "state") != "pending_closeout" || string(checked["native_closed"]) != "false" || string(checked["local_integration"]) == "null" {
		t.Fatalf("Delivery check concealed observed integration or misdirected closeout recovery: state=%s reasons=%s", checked["state"], checked["reasons"])
	}
	delivered := deliveryCLIOK(t, binary, cwd, "execute", id)
	if deliveryCLIString(t, delivered, "state") != "delivered" || string(delivered["native_closed"]) != "true" {
		t.Fatalf("Delivery CLI failed to resume a released automatic owner: %+v", delivered)
	}
	if got := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic"); got != before {
		t.Fatal("Delivery CLI closeout recovery repeated Git integration")
	}
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || closeout.State != "complete" || closeout.ResourceState != "kept" {
		t.Fatalf("Delivery CLI recovery omitted native retained closeout: %+v %v", closeout, err)
	}
	checked = deliveryCLIOK(t, binary, cwd, "check", id)
	if string(checked["native_closed"]) != "true" || deliveryCLIString(t, checked, "state") == "pending_closeout" {
		t.Fatalf("Delivery check hid already completed native closeout: state=%s native_closed=%s", checked["state"], checked["native_closed"])
	}
}

func TestAutomaticDeliveryCLIPreviewObservesCompletedNativeCloseout(t *testing.T) {
	binary := deliveryCLIBinary(t)
	f, o, _, _ := autoGateFixture(t, false, "")
	cwd := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	registered := deliveryCLIOK(t, binary, cwd, "register", "--file", deliveryCLIRegisterInput(t, f, o, "automatic-native-closeout-preview"))
	id := deliveryCLIString(t, registered, "id")
	if _, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err != nil {
		t.Fatal(err)
	}
	// The Delivery receipt is still registered; read actual native completion
	// before execute has a chance to copy its result into that receipt.
	checked := deliveryCLIOK(t, binary, cwd, "check", id)
	if deliveryCLIString(t, checked, "state") != "ready" || string(checked["native_closed"]) != "true" || string(checked["local_integration"]) == "null" {
		t.Fatalf("preview concealed completed native closeout: state=%s native_closed=%s reasons=%s", checked["state"], checked["native_closed"], checked["reasons"])
	}
	before := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	delivered := deliveryCLIOK(t, binary, cwd, "execute", id)
	if deliveryCLIString(t, delivered, "state") != "delivered" || string(delivered["native_closed"]) != "true" {
		t.Fatalf("Delivery failed to record already completed native closeout: %+v", delivered)
	}
	if got := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic"); got != before {
		t.Fatal("observed completed native closeout repeated Git integration")
	}
}

func TestAutomaticDeliveryRevalidatesLiveInputsAtEveryGate(t *testing.T) {
	f, o, binary, review := autoGateFixture(t, false, "")
	c := o.Delivery.Candidates[0]
	a, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult, *f.R.Delivery.Agreement)
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range []struct{ name, path string }{{"script", f.AcceptancePath}, {"binary", binary}, {"review", review}} {
		t.Run(input.name, func(t *testing.T) {
			raw, err := os.ReadFile(input.path)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := os.WriteFile(input.path, raw, 0700); err != nil {
					t.Fatal(err)
				}
			}()
			if err = os.WriteFile(input.path, append(raw, []byte("\nchanged input\n")...), 0700); err != nil {
				t.Fatal(err)
			}
			if _, err = CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, ""); err == nil {
				t.Fatal("changed automatic acceptance input passed the upper delivery gate")
			}
			preview, err := workspace.CheckTaskIntegration(f.D.Workspace, autoGateInput(a))
			if err == nil && (preview.Readback.Classification == "ready" || preview.Readback.RecoveryStatus == "complete") {
				t.Fatal("replayed authority bypassed changed input through native integration")
			}
			shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if shown.DeliveryStatus.Acceptance.Outcome != "blocked" || shown.DeliveryStatus.Acceptance.Current || shown.DeliveryStatus.Verification.InputsMatch || shown.DeliveryStatus.Verification.Outcome != "passed" {
				t.Fatalf("status erased completed verification or retained current acceptance: %+v", shown.DeliveryStatus)
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != a.ExpectedParentOID {
				t.Fatal("changed acceptance input allowed a Git effect")
			}
		})
	}
}

func TestAutomaticDeliveryRuntimeDriftInvalidatesCurrentAcceptance(t *testing.T) {
	f, o, _, _ := autoGateFixture(t, false, "")
	if err := os.WriteFile(filepath.Join(f.R.WorkspaceRoot, "actual-policy.json"), []byte("changed actual runtime policy observation\n"), 0600); err != nil {
		t.Fatal(err)
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if shown.DeliveryStatus.QualifiedCandidate.Current || shown.DeliveryStatus.Verification.InputsMatch {
		t.Fatal("fixture did not invalidate candidate freshness through actual runtime policy drift")
	}
	if shown.DeliveryStatus.Acceptance.Current || shown.DeliveryStatus.Acceptance.Outcome != "blocked" {
		t.Fatalf("runtime drift left a current automatic acceptance pass: %+v", shown.DeliveryStatus.Acceptance)
	}
	if _, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
		t.Fatal("runtime drift acquired an integration effect")
	}
}

func TestAutomaticDeliveryPreservesRequiredHumanAndIntegrationOwnerGates(t *testing.T) {
	for _, test := range []struct {
		name         string
		requireHuman bool
		owner        string
	}{{"human product judgment", true, ""}, {"human integration owner", false, workspace.IntegrationOwnerHuman}} {
		t.Run(test.name, func(t *testing.T) {
			f, o, _, _ := autoGateFixture(t, test.requireHuman, test.owner)
			c := o.Delivery.Candidates[0]
			if _, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
				t.Fatal("automatic test overrode a preserved human gate")
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") == c.OID {
				t.Fatal("human-gated automatic candidate changed the target")
			}
			if test.requireHuman {
				if o.Delivery.Phase != "awaiting_human_qa" || o.DeliveryStatus.Acceptance.Outcome != "blocked" {
					t.Fatalf("required human QA was not visible: %+v", o.DeliveryStatus)
				}
				// An external native answer must be used even though the workflow's
				// cached candidate still has no HumanQA pointer.
				if _, err := workflowhandoff.RecordDeliveryHumanQA(f.D.Workflow, "task", c.TaskResult, "pass", deliveryTestHuman(t, f, o, "pass")); err != nil {
					t.Fatal(err)
				}
				if _, err := WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err != nil {
					t.Fatalf("actual required human judgment did not unlock delivery: %v", err)
				}
			}
		})
	}
}

func TestAutomaticDeliveryLatestNativeHumanRejectionBlocks(t *testing.T) {
	f, o, _, _ := autoGateFixture(t, false, "")
	c := o.Delivery.Candidates[0]
	a, err := ValidateDeliveryAuthority(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult, *f.R.Delivery.Agreement)
	if err != nil {
		t.Fatal(err)
	}
	for i, outcome := range []string{"fail", "blocked"} {
		answer := workflowProviderDocument(t, deliveryTestHuman(t, f, o, outcome))
		answer["completed_at_utc"] = []string{"2026-10-05T00:00:02Z", "2026-10-05T00:00:03Z"}[i]
		if _, err = workflowhandoff.RecordDeliveryHumanQA(f.D.Workflow, "task", c.TaskResult, outcome, writeAny(t, f.R.WorkspaceRoot, "external-"+outcome+".json", answer)); err != nil {
			t.Fatal(err)
		}
		if _, err = CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, ""); err == nil || !strings.Contains(err.Error(), outcome) {
			t.Fatalf("latest actual human %s did not block automatic delivery: %v", outcome, err)
		}
		preview, err := workspace.CheckTaskIntegration(f.D.Workspace, autoGateInput(a))
		if err == nil && (preview.Readback.Classification == "ready" || preview.Readback.RecoveryStatus == "complete") {
			t.Fatal("earlier automatic acceptance bypassed newer native human answer")
		}
		shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
		if err != nil || shown.DeliveryStatus.Acceptance.Outcome != outcome || shown.DeliveryStatus.HumanJudgment == nil || shown.DeliveryStatus.HumanJudgment.Outcome != outcome {
			t.Fatalf("status failed to preserve latest native human %s: %+v %v", outcome, shown.DeliveryStatus, err)
		}
		if gitOutput(t, f.Parent, "rev-parse", "HEAD") != a.ExpectedParentOID {
			t.Fatal("negative human judgment allowed an integration effect")
		}
	}
}

func TestAutomaticDeliveryUnchangedCandidateRetainsHumanRejectionAfterReverification(t *testing.T) {
	for _, outcome := range []string{"fail", "blocked"} {
		t.Run(outcome, func(t *testing.T) {
			f, o, binary, review := autoGateFixture(t, false, "")
			c := o.Delivery.Candidates[0]
			if _, err := workflowhandoff.RecordDeliveryHumanQA(f.D.Workflow, "task", c.TaskResult, outcome, deliveryTestHuman(t, f, o, outcome)); err != nil {
				t.Fatal(err)
			}
			// A new verifier invocation does not correct the product bytes the
			// human rejected, even if native qualification assigns another ID.
			reverified, err := WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, binary)
			if err != nil {
				if gitOutput(t, f.Parent, "rev-parse", "HEAD") == c.OID {
					t.Fatal("rejected unchanged candidate acquired an integration effect")
				}
				return
			}
			if reverified.Delivery.Phase == "automatic_acceptance_passed" || strings.Contains(reverified.NextAction.Message, "Complete its authorized native local Epic delivery") {
				t.Fatalf("unchanged candidate with human %s was announced as accepted: phase=%s next=%+v", outcome, reverified.Delivery.Phase, reverified.NextAction)
			}
			if _, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, reverified.Paths.Context); err == nil {
				t.Fatalf("unchanged candidate bypassed human %s through a fresh TaskResult ID", outcome)
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") == c.OID {
				t.Fatal("rejected unchanged candidate acquired an integration effect")
			}
		})
	}
}
