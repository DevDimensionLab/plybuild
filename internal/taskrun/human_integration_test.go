package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Exact DeliveryState wire shape from the immutable controller baseline
// 680fecc. Its decoder rejects unknown fields, so a negative human answer must
// return ownership using this shape while retaining new facts in event files.
type legacyDeliveryStateBeforeHumanIntegration struct {
	Phase                string                `json:"phase"`
	OwnerClaim           string                `json:"owner_claim"`
	PermissionState      string                `json:"permission_state"`
	ActualPolicySHA256   *string               `json:"actual_policy_sha256"`
	ActualPolicyEvidence []Evidence            `json:"actual_policy_evidence"`
	Events               []DeliveryEventRecord `json:"events"`
	Candidates           []DeliveryCandidate   `json:"candidates"`
	Attempt              *DeliveryAttempt      `json:"attempt"`
	LastEventSHA256      *string               `json:"last_event_sha256"`
}

func humanIntegrationTestInput(t *testing.T, f deliveryFixture, o WorkflowRun, outcome, key string) HumanIntegrationInput {
	t.Helper()
	var answer workflowhandoff.DeliveryHumanAttestation
	if e := readValue(deliveryTestHuman(t, f, o, outcome), 1<<20, &answer); e != nil {
		t.Fatal(e)
	}
	in := HumanIntegrationInput{PlanID: "int_" + key, PlanSHA256: hash([]byte("plan-" + key)), Attestation: answer}
	path := writeAny(t, f.R.WorkspaceRoot, "decision-"+key+"-"+outcome+".json", map[string]any{"kind": "ply.integration.decision", "schema_version": 1, "integration_id": in.PlanID, "plan_sha256": in.PlanSHA256, "human": answer})
	in.Decision = FileBinding{path, hashFileTest(t, path)}
	return in
}

func humanIntegrationCandidate(t *testing.T, human bool) (deliveryFixture, WorkflowRun) {
	t.Helper()
	a := workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalBranch}
	if human {
		a.SchemaVersion = 2
		a.IntegrationOwner = workspace.IntegrationOwnerHuman
	}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, ParentRef: "main"})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "human integration candidate")
	o, e := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "human"))
	if e != nil {
		t.Fatal(e)
	}
	return f, o
}

func TestHumanIntegrationRequiresReleaseAndPreservedPlan(t *testing.T) {
	for _, human := range []bool{false, true} {
		t.Run(map[bool]string{false: "legacy-explicit-takeover", true: "human-owner-v2"}[human], func(t *testing.T) {
			f, o := humanIntegrationCandidate(t, human)
			c := o.Delivery.Candidates[0]
			requestBefore, e := os.ReadFile(f.File)
			if e != nil {
				t.Fatal(e)
			}
			if human && (f.R.Delivery.LocalIntegration != "none" || len(f.R.Delivery.Agreement.AllowedEffects()) != 0 || f.R.Delivery.Agreement.StopAfter() != "qualified_candidate_before_human_integration") {
				t.Fatal("human owner acquired developer integration authority")
			}
			in := humanIntegrationTestInput(t, f, o, "pass", "exact")
			statePath := filepath.Join(o.Paths.RunRoot, "state.json")
			before, _ := os.ReadFile(statePath)
			p, e := PreviewHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
			if e != nil || p.Released || p.Reason != "owner_active" {
				t.Fatalf("active owner preview: %+v %v", p, e)
			}
			after, _ := os.ReadFile(statePath)
			if !bytes.Equal(before, after) {
				t.Fatal("preview mutated state")
			}
			if _, e = AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in); e == nil || !strings.Contains(e.Error(), "owner_active") {
				t.Fatalf("prompt replaced ownership release: %v", e)
			}
			// QA alone still cannot authorize a human-owned developer integration.
			o, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
			if e != nil {
				t.Fatal(e)
			}
			old, e := nativeDeliveryCandidate(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
			if e != nil {
				t.Fatal(e)
			}
			if human {
				if _, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); e == nil {
					t.Fatal("QA pass expanded human-owned developer mandate")
				}
			}
			// A newer installed CLI may perform this release-only revocation.
			newCLI := filepath.Join(f.R.WorkspaceRoot, "new-ply")
			if e = os.WriteFile(newCLI, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
				t.Fatal(e)
			}
			f.D.Executable = func() (string, error) { return newCLI, nil }
			o, e = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic source owner stops writing before human integration.")
			if e != nil {
				t.Fatal(e)
			}
			if e = CheckAgentDeliveryExecution(f.D, f.R.WorkspaceRoot, o.RunID); e == nil {
				t.Fatal("released agent could execute delivery")
			}
			if _, e = CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, old.Candidate.HumanQA.ID); e == nil {
				t.Fatal("old local API bypassed takeover")
			}
			if _, e = workflowhandoff.IntegrateDeliveryCandidate(f.D.Workflow, string(c.TaskResult.TaskID), c.TaskResult.ID, old.Candidate.HumanQA.ID, old.ExpectedParentOID, c.OID, workflowhandoff.DeliveryOwnerAuthority{RunID: o.RunID, RequestSHA256: old.RequestSHA256, ActorClaim: old.OwnerClaim, PreparationID: old.PreparationID, Authorization: old.Authorization}); e == nil {
				t.Fatal("stale native owner bypassed release")
			}
			bad := in
			bad.PlanSHA256 = hash([]byte("changed plan"))
			if _, e = AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, bad); e == nil {
				t.Fatal("changed plan consumed original decision")
			}
			accepted, e := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if e != nil || !accepted.Accepted || accepted.HumanQA.Outcome != "pass" {
				t.Fatalf("human acceptance: %+v %v", accepted, e)
			}
			if _, e = CompleteHumanLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, hash([]byte("different"))); e == nil {
				t.Fatal("different plan completed integration")
			}
			completed, e := CompleteHumanLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, in.PlanSHA256)
			if e != nil || !completed.Completed {
				t.Fatalf("human native completion: %+v %v", completed, e)
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
				t.Fatal("exact candidate not integrated")
			}
			negative := humanIntegrationTestInput(t, f, o, "fail", "exact")
			negative.Attestation.CompletedAtUTC = "2026-10-05T00:00:03Z"
			negative.Decision.Locator = writeAny(t, f.R.WorkspaceRoot, "negative-after-integration.json", map[string]any{"kind": "ply.integration.decision", "schema_version": 1, "integration_id": negative.PlanID, "plan_sha256": negative.PlanSHA256, "human": negative.Attestation})
			negative.Decision.SHA256 = hashFileTest(t, negative.Decision.Locator)
			if _, e = AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, negative); e == nil || !strings.Contains(e.Error(), "integration_already_observed") {
				t.Fatalf("negative decision restored authority after completed integration: %v", e)
			}
			if _, e = CompleteHumanLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, in.PlanSHA256); e != nil {
				t.Fatalf("same human plan could not resume: %v", e)
			}
			inventory, e := ReadInventory(f.D, f.R.WorkspaceRoot)
			if e != nil || len(inventory.Runs) != 1 || inventory.Runs[0].Freshness != "fresh" {
				t.Fatalf("human integration history lost event integrity: %+v %v", inventory, e)
			}
			trace, e := ReadTraceHistory(f.D, f.R.WorkspaceRoot, string(c.TaskResult.TaskID))
			if e != nil || len(trace.Runs) != 1 || trace.Runs[0].Coverage != "complete" {
				t.Fatalf("human integration trace: %+v %v", trace, e)
			}
			seenRelease, seenDecision := false, false
			for _, event := range trace.Runs[0].Entries {
				seenRelease = seenRelease || event.Kind == "source_ownership_release"
				seenDecision = seenDecision || event.Kind == "human_integration_decision"
			}
			if !seenRelease || !seenDecision {
				t.Fatal("human takeover history omitted native release or actual human decision")
			}
			requestAfter, _ := os.ReadFile(f.File)
			if !bytes.Equal(requestBefore, requestAfter) {
				t.Fatal("human takeover rewrote frozen request")
			}
		})
	}
}

func TestHumanIntegrationNegativeDecisionReturnsCorrectionsAndRevokesPass(t *testing.T) {
	for _, outcome := range []string{"fail", "blocked"} {
		t.Run(outcome, func(t *testing.T) {
			f, o := humanIntegrationCandidate(t, true)
			c := o.Delivery.Candidates[0]
			before := gitOutput(t, f.Parent, "rev-parse", "HEAD")
			if _, e := ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic release."); e != nil {
				t.Fatal(e)
			}
			pass := humanIntegrationTestInput(t, f, o, "pass", "negative")
			if _, e := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, pass); e != nil {
				t.Fatal(e)
			}
			in := humanIntegrationTestInput(t, f, o, outcome, "negative")
			in.Attestation.CompletedAtUTC = "2026-10-05T00:00:02Z"
			in.Decision.Locator = writeAny(t, f.R.WorkspaceRoot, "later-negative.json", map[string]any{"kind": "ply.integration.decision", "schema_version": 1, "integration_id": in.PlanID, "plan_sha256": in.PlanSHA256, "human": in.Attestation})
			in.Decision.SHA256 = hashFileTest(t, in.Decision.Locator)
			a, e := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in)
			if e != nil || a.Accepted || a.HumanQA.Outcome != outcome {
				t.Fatalf("negative answer: %+v %v", a, e)
			}
			if _, e = CompleteHumanLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, pass.PlanSHA256); e == nil {
				t.Fatal("negative judgment left previous plan authority")
			}
			current, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if e != nil {
				t.Fatal(e)
			}
			if current.Result.Delivery.OwnershipRelease != nil || current.Result.NextAction.Actor != "recipient" {
				t.Fatal("negative answer did not return correction ownership")
			}
			if current.Result.Delivery.HumanIntegration != nil {
				t.Fatal("negative answer retained a takeover field unknown to the immutable old controller")
			}
			encoded, err := json.Marshal(current.Result.Delivery)
			if err != nil {
				t.Fatal(err)
			}
			var legacy legacyDeliveryStateBeforeHumanIntegration
			if err = decode(encoded, 8<<20, &legacy); err != nil {
				t.Fatalf("original strict controller cannot decode returned correction state: %v", err)
			}
			if again, err := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil || !equal(again, a) {
				t.Fatalf("immutable negative decision did not remain idempotent: %+v %v", again, err)
			}
			if gitOutput(t, f.Parent, "rev-parse", "HEAD") != before {
				t.Fatal("negative answer integrated source")
			}
			deliveryTestCandidate(t, f, "corrected candidate")
			if _, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "corrected")); e != nil {
				t.Fatalf("same owner cannot qualify correction: %v", e)
			}
		})
	}
}

func TestHumanIntegrationRecoversOnlySameReservedDecision(t *testing.T) {
	f, o := humanIntegrationCandidate(t, true)
	if _, e := ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic release."); e != nil {
		t.Fatal(e)
	}
	in := humanIntegrationTestInput(t, f, o, "pass", "recovery")
	f.D.Fault = func(point string) error {
		if point == "human_integration_after_reservation" {
			return errors.New("injected reservation interruption")
		}
		return nil
	}
	if _, e := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in); e == nil {
		t.Fatal("fault did not interrupt")
	}
	f.D.Fault = nil
	other := humanIntegrationTestInput(t, f, o, "pass", "other")
	if _, e := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, other); e == nil {
		t.Fatal("different plan replaced reserved answer")
	}
	first, e := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if e != nil {
		t.Fatal(e)
	}
	second, e := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if e != nil || first.HumanQA.ID != second.HumanQA.ID {
		t.Fatalf("same decision did not reuse QA: %v", e)
	}
}

func TestHumanIntegrationRequiresPositiveProcessExitWithoutCallback(t *testing.T) {
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalEpic}})
	herdr, e := os.ReadFile(f.R.Herdr.Executable.Path)
	if e != nil {
		t.Fatal(e)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+deliveryRecoveryInspectStandin), 1)
	if e = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); e != nil {
		t.Fatal(e)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	f.D.StartupProcessTree = func(int) ([]StartupProcess, error) {
		return []StartupProcess{{PID: 4242, ParentPID: 1, ProcessGroupID: 4242}}, nil
	}
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "ended provider candidate")
	o, e = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "ended"))
	if e != nil {
		t.Fatal(e)
	}
	c := o.Delivery.Candidates[0]
	if _, e = ReleaseExitedDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID); e == nil {
		t.Fatal("idle managed provider was treated as exited")
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_armed": true, "recovery_started": false})
	f.D.StartupProcessTree = func(int) ([]StartupProcess, error) {
		return []StartupProcess{{PID: 4242, ParentPID: 1, ProcessGroupID: 4242}, {PID: 4343, ParentPID: 4242, ProcessGroupID: 4343}}, nil
	}
	if _, e = ReleaseExitedDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID); e == nil {
		t.Fatal("background writer was treated as exited")
	}
	f.D.StartupProcessTree = func(int) ([]StartupProcess, error) {
		return []StartupProcess{{PID: 4242, ParentPID: 1, ProcessGroupID: 4242}}, nil
	}
	p, e := PreviewHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
	if e != nil || p.Released || p.Reason != "owner_exited_release_required" {
		t.Fatalf("ended provider preview: %+v %v", p, e)
	}
	if _, e = ReleaseExitedDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID); e != nil {
		t.Fatal(e)
	}
	p, e = PreviewHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID)
	if e != nil || !p.Released || p.Release == nil {
		t.Fatalf("ended provider release: %+v %v", p, e)
	}
	if _, e = AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, humanIntegrationTestInput(t, f, o, "pass", "ended")); e != nil {
		t.Fatalf("human needed original callback/provider: %v", e)
	}
}

func TestLegacyOwnershipRejectsActiveAndAcceptsExactNativeRelease(t *testing.T) {
	f, o := humanIntegrationCandidate(t, false)
	c := o.Delivery.Candidates[0]
	if _, e := CheckLegacyTaskOwnership(f.D, f.R.WorkspaceRoot, c.TaskResult); e == nil {
		t.Fatal("legacy reconciliation ignored active native writer")
	}
	if _, e := ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic explicit source release."); e != nil {
		t.Fatal(e)
	}
	observation, e := CheckLegacyTaskOwnership(f.D, f.R.WorkspaceRoot, c.TaskResult)
	if e != nil || !observation.Ready || len(observation.Owners) != 1 || observation.Owners[0].State != "explicitly_released" {
		t.Fatalf("exact legacy release: %+v %v", observation, e)
	}
	bad := c.TaskResult
	bad.ResultOID = strings.Repeat("1", 40)
	if _, e = CheckLegacyTaskOwnership(f.D, f.R.WorkspaceRoot, bad); e == nil {
		t.Fatal("changed legacy result accepted")
	}
}

func TestAgentLocalCompletionCannotConsumeConcurrentHumanTakeover(t *testing.T) {
	f, o := humanIntegrationCandidate(t, false)
	var e error
	o, e = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if e != nil {
		t.Fatal(e)
	}
	c := o.Delivery.Candidates[0]
	before := gitOutput(t, f.Parent, "rev-parse", "HEAD")
	in := humanIntegrationTestInput(t, f, o, "pass", "concurrent-takeover")
	f.D.Fault = func(point string) error {
		if point != "local_delivery_after_agent_gate" {
			return nil
		}
		if _, e := ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic concurrent ownership release."); e != nil {
			return e
		}
		_, e := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in)
		return e
	}
	if _, e = CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, c.HumanQA.ID); e == nil || !strings.Contains(e.Error(), "human_integration_required") {
		t.Fatalf("agent consumed concurrent human authority: %v", e)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != before {
		t.Fatal("agent path integrated after a concurrent human takeover")
	}
}

type releaseDuringIntegrationGit struct {
	workspace.TaskIntegrationGit
	beforeMerge func()
}

func (g releaseDuringIntegrationGit) MergeFastForward(in workspace.IntegrationMergeInput) workspace.GitCommandOutcome {
	g.beforeMerge()
	return g.TaskIntegrationGit.MergeFastForward(in)
}

func TestLocalIntegrationReservesOwnerThroughNativeEffects(t *testing.T) {
	f, o := humanIntegrationCandidate(t, false)
	var err error
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if err != nil {
		t.Fatal(err)
	}
	c := o.Delivery.Candidates[0]
	called := false
	wd := *f.D.Workflow.TaskWorkspace
	wd.IntegrationGit = releaseDuringIntegrationGit{TaskIntegrationGit: wd.IntegrationGit, beforeMerge: func() {
		called = true
		if _, e := ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic release racing the actual merge."); e == nil || !strings.Contains(e.Error(), "owner_active") {
			t.Fatalf("in-flight native effect did not prevent release: %v", e)
		}
		s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
		if e != nil || s.Result.Delivery.OwnershipRelease != nil {
			t.Fatalf("release published before integration settled: %v", e)
		}
	}}
	f.D.Workflow.TaskWorkspace = &wd
	result, err := CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, c.HumanQA.ID)
	if err != nil || !result.Completed || !called {
		t.Fatalf("reserved completion: %+v %v merge=%v", result, err, called)
	}
	if _, err = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic release after exact native effects completed."); err != nil {
		t.Fatalf("settled effect retained active reservation: %v", err)
	}
}

func TestHumanCloseoutRetainsOnlyItsOwnSourceReservation(t *testing.T) {
	f, o := humanIntegrationCandidate(t, true)
	c := o.Delivery.Candidates[0]
	if _, err := ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic release."); err != nil {
		t.Fatal(err)
	}
	in := humanIntegrationTestInput(t, f, o, "pass", "pending-closeout")
	if _, err := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
		t.Fatal(err)
	}
	if _, err := CompleteHumanLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, in.PlanSHA256); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if err = workflowReservations(f.D, f.R.WorkspaceRoot, s.Observed.Target); err == nil || !strings.Contains(err.Error(), "closeout") {
		t.Fatalf("completed integration released source before closeout: %v", err)
	}
	other := s.Observed.Target
	other.WorktreeLocator += "-next-task"
	other.Ref += "-next-task"
	if err = workflowReservations(f.D, f.R.WorkspaceRoot, other); err != nil {
		t.Fatalf("pending closeout blocked an unrelated Task source: %v", err)
	}
	own := *s.Result.Delivery.OwnershipRelease
	closeInput := workspace.TaskCloseoutInput{OperationID: in.PlanID, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ExpectedResultOID: c.OID, ExpectedResultTree: c.Tree, TargetRef: f.R.Delivery.Agreement.TargetRef, Keep: true, Ownership: workspace.TaskCloseoutOwnership{Released: true, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ResultOID: c.OID, ResultTree: c.Tree, EvidenceLocator: own.Locator, EvidenceSHA256: own.SHA256}}
	p, err := workspace.PreviewTaskCloseout(f.D.Workspace, closeInput)
	if err != nil || !p.Ready {
		t.Fatalf("closeout preview: %+v %v", p, err)
	}
	if _, err = workspace.ReserveTaskCloseout(f.D.Workspace, closeInput, p.PlanSHA256, "installation_pending", func() error {
		_, err := CheckLegacyTaskOwnership(f.D, f.R.WorkspaceRoot, c.TaskResult)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	for _, phase := range []string{"installation_pending", "install_failed", "install_unknown"} {
		if _, err = workspace.SetTaskCloseoutPhase(f.D.Workspace, c.TaskResult.TaskID, in.PlanID, p.PlanSHA256, phase); err != nil {
			t.Fatal(err)
		}
		shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
		if err != nil || shown.Closeout == nil || shown.Closeout.State != "cleanup_pending" || shown.Closeout.Phase != phase || !strings.Contains(shown.NextAction.Message, "installation") {
			t.Fatalf("pending installation lost in workflow readback: %+v %v", shown, err)
		}
		inventory, err := ReadInventory(f.D, f.R.WorkspaceRoot)
		if err != nil || len(inventory.Runs) != 1 || inventory.Runs[0].Closeout == nil || inventory.Runs[0].Closeout.Phase != phase || !strings.Contains(inventory.Runs[0].Delivery.NextAction.Kind, "installation") {
			t.Fatalf("pending installation lost in inventory: %+v %v", inventory, err)
		}
		if _, err = os.Stat(c.TaskResult.SourceLocator); err != nil {
			t.Fatal("pending installation removed source", err)
		}
	}
	if _, err = workspace.SetTaskCloseoutPhase(f.D.Workspace, c.TaskResult.TaskID, in.PlanID, p.PlanSHA256, "awaiting_closeout"); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.ApplyTaskCloseout(f.D.Workspace, closeInput, p.PlanSHA256); err != nil {
		t.Fatal(err)
	}
	if err = workflowReservations(f.D, f.R.WorkspaceRoot, s.Observed.Target); err != nil {
		t.Fatalf("completed exact closeout retained workflow reservation: %v", err)
	}
}

func TestLocalIntegrationInterruptedReservationBlocksReleaseAndResumes(t *testing.T) {
	f, o := humanIntegrationCandidate(t, false)
	var err error
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if err != nil {
		t.Fatal(err)
	}
	c := o.Delivery.Candidates[0]
	f.D.Fault = func(point string) error {
		if point == "delivery_after_integration_reservation" {
			return errors.New("synthetic interruption after durable reservation")
		}
		return nil
	}
	if _, err = CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, c.HumanQA.ID); err == nil {
		t.Fatal("interrupted completion succeeded")
	}
	if _, err = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic attempted release of unresolved effect."); err == nil || !strings.Contains(err.Error(), "owner_active") {
		t.Fatalf("unresolved reservation released ownership: %v", err)
	}
	f.D.Fault = nil
	if out, err := CompleteLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, c.HumanQA.ID); err != nil || !out.Completed {
		t.Fatalf("exact operation did not resume: %+v %v", out, err)
	}
	r, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || len(r.IntegrationAttempts) != 1 {
		t.Fatalf("native integration repeated or disappeared: %d %v", len(r.IntegrationAttempts), err)
	}
}

func TestPullRequestOwnershipCannotReleaseBeforePublication(t *testing.T) {
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryPullRequest, GitHubRepository: "fixture/repository", Remote: "origin"}, ParentRef: "main"})
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "synthetic PR candidate")
	var err error
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "pr-release"))
	if err != nil {
		t.Fatal(err)
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", deliveryTestHuman(t, f, o, "pass"))
	if err != nil {
		t.Fatal(err)
	}
	c := o.Delivery.Candidates[0]
	if _, err = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic early release."); err == nil || !strings.Contains(err.Error(), "PR publication") {
		t.Fatalf("qualified candidate released unsettled PR publication: %v", err)
	}
	if _, err = ReleaseExitedDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID); err == nil || !strings.Contains(err.Error(), "PR publication") {
		t.Fatalf("exit path bypassed unsettled PR publication: %v", err)
	}
	receipt := writeAny(t, f.R.WorkspaceRoot, "synthetic-pr-receipt.json", map[string]any{"fixture": "synthetic exact observation; no GitHub effect"})
	_, err = CompletePullRequestDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, FileBinding{receipt, hashFileTest(t, receipt)}, PullRequestDeliveryObservation{Repository: f.R.Delivery.Agreement.GitHubRepository, HeadRef: f.R.Delivery.Agreement.SourceRef, HeadOID: c.OID, BaseRef: f.R.Delivery.Agreement.TargetRef, URL: "https://example.invalid/fixture/pull/1", Number: 1, MetadataApplied: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic release after completed publication."); err != nil {
		t.Fatalf("completed publication cannot release: %v", err)
	}
}

type closeoutRetirementFailureStore struct{ workspace.WorkItemStore }

func (s closeoutRetirementFailureStore) WithLock(root string, action func(workspace.WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session workspace.WorkItemStoreSession) error {
		return action(closeoutRetirementFailureSession{session})
	})
}

type closeoutRetirementFailureSession struct{ workspace.WorkItemStoreSession }

func (s closeoutRetirementFailureSession) Publish(registry workspace.WorkItemRegistry) error {
	for _, task := range registry.Tasks {
		if task.WorktreeState == workspace.WorkItemRetired {
			return errors.New("synthetic retirement publication interruption")
		}
	}
	return s.WorkItemStoreSession.Publish(registry)
}

func TestWorkflowCloseoutReadbackPreservesInterruptedRemoval(t *testing.T) {
	f, o := humanIntegrationCandidate(t, true)
	c := o.Delivery.Candidates[0]
	if _, err := ReleaseDeliveryOwnership(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "Synthetic owner release."); err != nil {
		t.Fatal(err)
	}
	in := humanIntegrationTestInput(t, f, o, "pass", "interrupted-removal")
	if _, err := AcceptHumanIntegration(f.D, f.R.WorkspaceRoot, o.RunID, in); err != nil {
		t.Fatal(err)
	}
	if _, err := CompleteHumanLocalDelivery(f.D, f.R.WorkspaceRoot, o.RunID, c.TaskResult.ID, in.PlanSHA256); err != nil {
		t.Fatal(err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	own := *s.Result.Delivery.OwnershipRelease
	closeInput := workspace.TaskCloseoutInput{OperationID: in.PlanID, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ExpectedResultOID: c.OID, ExpectedResultTree: c.Tree, TargetRef: f.R.Delivery.Agreement.TargetRef, Ownership: workspace.TaskCloseoutOwnership{Released: true, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ResultOID: c.OID, ResultTree: c.Tree, EvidenceLocator: own.Locator, EvidenceSHA256: own.SHA256}}
	p, err := workspace.PreviewTaskCloseout(f.D.Workspace, closeInput)
	if err != nil || !p.Ready {
		t.Fatalf("closeout preview: %+v %v", p, err)
	}
	interrupted := f.D.Workspace
	interrupted.WorkItems = closeoutRetirementFailureStore{interrupted.WorkItems}
	if _, err := workspace.ApplyTaskCloseout(interrupted, closeInput, p.PlanSHA256); err == nil || !strings.Contains(err.Error(), "synthetic retirement") {
		t.Fatalf("fault did not interrupt after actual source removal: %v", err)
	}
	if _, err := os.Stat(c.TaskResult.SourceLocator); !os.IsNotExist(err) {
		t.Fatal("fixture source was not actually removed", err)
	}
	f.D.CWD = func() (string, error) { return f.Parent, nil }
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil || shown.Closeout == nil || shown.Closeout.State != "cleanup_pending" || !shown.Closeout.WorktreeRemoved || shown.Closeout.LifecycleCompleted || !strings.Contains(shown.NextAction.Message, "closeout") {
		t.Fatalf("pending observed removal lost its resume action: %+v %v", shown, err)
	}
	for _, reason := range shown.Reasons {
		if reason.Code == "workflow_run_drift" {
			t.Fatalf("durably removed source mistaken for binding drift: %+v", reason)
		}
	}
	if _, err := workspace.ApplyTaskCloseout(f.D.Workspace, closeInput, p.PlanSHA256); err != nil {
		t.Fatalf("same interrupted closeout could not resume: %v", err)
	}
}
