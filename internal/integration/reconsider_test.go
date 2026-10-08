package integration

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These are integrity/journal fixtures below the native acceptance gate.
// Native human acceptance is covered separately by taskrun's CLI scenarios.
func reconsiderReceipt(t *testing.T, root, choice string) Receipt {
	t.Helper()
	p := Plan{Kind: "ply.integration.plan", SchemaVersion: 1, Workspace: root, Flow: "local_return", WorkflowRunID: "fixture-workflow", EffectKey: byteHash([]byte("one-candidate-target")), Selection: Selection{DeliveryID: choice}, TaskResult: workspace.TaskResultRecord{TaskID: "task", ID: "result", ResultOID: strings.Repeat("a", 40), ResultTree: strings.Repeat("b", 40)}}
	h := digest(p)
	r := Receipt{Kind: "ply.integration.receipt", SchemaVersion: 1, ID: operationID(h), Plan: p, PlanSHA256: h, State: "blocked", Events: []Event{}}
	r.Decision = Decision{Kind: "ply.integration.decision", SchemaVersion: 1, IntegrationID: r.ID, PlanSHA256: h, Human: reconsiderHuman(r, "pass")}
	r.DecisionFile = taskrun.FileBinding{Locator: filepath.Join(storeRoot(root), "decisions", r.ID+".json"), SHA256: digest(r.Decision)}
	if err := atomicJSON(r.DecisionFile.Locator, r.Decision, true); err != nil {
		t.Fatal(err)
	}
	if err := atomicJSON(receiptPath(root, r.ID), r, false); err != nil {
		t.Fatal(err)
	}
	return r
}

func reconsiderHuman(r Receipt, answer string) workflowhandoff.DeliveryHumanAttestation {
	return workflowhandoff.DeliveryHumanAttestation{Kind: "DeliveryHumanAttestation@1", SchemaVersion: 1, TaskID: string(r.Plan.TaskResult.TaskID), TaskResultID: string(r.Plan.TaskResult.ID), ResultOID: r.Plan.TaskResult.ResultOID, ResultTree: r.Plan.TaskResult.ResultTree, Answer: answer, Outcome: answer, ActorClaim: "Fixture human", StartSurface: "Fixture reconsideration", StartedAtUTC: "2026-10-08T12:00:00Z", CompletedAtUTC: "2026-10-08T12:01:00Z", Observation: "Synthetic journal fixture; no actual human product judgment."}
}

func revokedFixture(t *testing.T, r Receipt) Receipt {
	t.Helper()
	h := reconsiderHuman(r, "blocked")
	decision := Decision{Kind: "ply.integration.decision", SchemaVersion: 1, IntegrationID: r.ID, PlanSHA256: r.PlanSHA256, Human: h}
	file := taskrun.FileBinding{Locator: filepath.Join(storeRoot(r.Plan.Workspace), "reconsiderations", r.ID, "decision.json"), SHA256: digest(decision)}
	if err := atomicJSON(file.Locator, decision, true); err != nil {
		t.Fatal(err)
	}
	native := taskrun.HumanIntegrationAcceptance{Kind: "PlyHumanIntegrationDecision@1", SchemaVersion: 1, RunID: r.Plan.WorkflowRunID, PlanID: r.ID, PlanSHA256: r.PlanSHA256, Decision: file, Attestation: h, Accepted: false, HumanQA: workspace.TaskHumanQARecord{TaskID: r.Plan.TaskResult.TaskID, TaskResultID: r.Plan.TaskResult.ID, ResultOID: h.ResultOID, ResultTree: h.ResultTree, Outcome: h.Answer}}
	c := &ReconsiderationReceipt{Kind: "ply.integration.reconsideration", SchemaVersion: 1, IntegrationID: r.ID, PlanSHA256: r.PlanSHA256, Decision: decision, DecisionFile: file, Native: &native, Evidence: ReconsiderationEvidence{Reason: "never_dispatched"}, RecordedAtUTC: h.CompletedAtUTC}
	proof := taskrun.FileBinding{Locator: filepath.Join(filepath.Dir(file.Locator), "revocation.json"), SHA256: digest(*c)}
	if err := atomicJSON(proof.Locator, *c, true); err != nil {
		t.Fatal(err)
	}
	c.Proof = &proof
	r.Reconsideration, r.State = c, "cancelled"
	if err := atomicJSON(receiptPath(r.Plan.Workspace, r.ID), r, false); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestReadReceiptRejectsCancelledPlanWithoutDurableNegativeRevocation(t *testing.T) {
	r := reconsiderReceipt(t, adapterPhysicalTemp(t), "original")
	r.State = "cancelled"
	if err := atomicJSON(receiptPath(r.Plan.Workspace, r.ID), r, false); err != nil {
		t.Fatal(err)
	}
	if _, err := readReceipt(r.Plan.Workspace, r.ID); err == nil {
		t.Fatal("cancelled state without actual preserved negative decision bypassed the active-plan guard")
	}
}

func TestRevokedPlanReservationCanBeSupersededWithoutRewritingHistory(t *testing.T) {
	root := adapterPhysicalTemp(t)
	original := reconsiderReceipt(t, root, "original")
	if err := reserve(original); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(storeRoot(root), "effects", strings.TrimPrefix(original.Plan.EffectKey, "sha256:")+".json")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	replacement := reconsiderReceipt(t, root, "changed-method")
	if err = reserve(replacement); err == nil {
		t.Fatal("live prior reservation was replaced")
	}
	original = revokedFixture(t, original)
	if err = reserve(replacement); err != nil {
		t.Fatalf("durable actual negative revocation could not release rejected plan: %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(before) != string(after) {
		t.Fatalf("original effect reservation was overwritten: %v", err)
	}
	if err = reserve(replacement); err != nil {
		t.Fatalf("current owner cannot resume: %v", err)
	}
	if err = reserve(original); err == nil {
		t.Fatal("revoked original plan reclaimed replacement's reservation")
	}
	third := reconsiderReceipt(t, root, "different-cleanup")
	if err = reserve(third); err == nil {
		t.Fatal("third owner bypassed current active replacement")
	}
	replacement = revokedFixture(t, replacement)
	if err = reserve(third); err != nil {
		t.Fatalf("second exact revocation did not extend the immutable chain: %v", err)
	}
	if err = os.Remove(original.Reconsideration.Proof.Locator); err != nil {
		t.Fatal(err)
	}
	if err = reserve(third); err == nil {
		t.Fatal("current owner trusted chain with missing predecessor proof")
	}
}

func TestReconsiderationNeverCancelsUnknownOrAlreadyIntegratedEffect(t *testing.T) {
	for _, name := range []string{"unknown-dispatch", "pending", "merged", "unstarted-but-merged", "wrong-rejection-state", "install-started", "integrated", "closeout-started", "local-intent"} {
		t.Run(name, func(t *testing.T) {
			r := rejectedMergeReceipt(t)
			r.Decision.Human.Answer = "pass"
			a := &rejectedMergeFixture{observation: PRMergeObservation{State: "ready"}}
			switch name {
			case "unknown-dispatch":
				r.PR.NoEffect = false
			case "pending":
				a.observation.State = "remote_pending"
			case "merged":
				a.observation.State = "merged"
			case "unstarted-but-merged":
				r.IntegrationStarted = false
				a.observation.State = "merged"
			case "wrong-rejection-state":
				r.PR.State = "effect_unknown"
			case "install-started":
				r.InstallationStarted = true
			case "integrated":
				r.Integrated = true
			case "closeout-started":
				r.Closeout = &workspace.TaskCloseoutReceipt{}
			case "local-intent":
				r.Plan.Flow = "local_return"
			}
			s := NewService(taskrun.Dependencies{}, Options{PR: a})
			if _, err := s.reconsiderEvidence(r); err == nil {
				t.Fatal("reconsideration accepted without proven no-effect boundary")
			}
			if a.merges != 0 {
				t.Fatal("reconsideration observation dispatched a merge")
			}
		})
	}
}

func TestCancelledReceiptRejectsChangedNativeRevocationOrNoEffectProof(t *testing.T) {
	for _, name := range []string{"native-kind", "native-run", "native-task", "native-pass", "proof-bytes", "dispatched", "no-native", "decision-answer"} {
		t.Run(name, func(t *testing.T) {
			r := revokedFixture(t, reconsiderReceipt(t, adapterPhysicalTemp(t), "original"))
			switch name {
			case "native-kind":
				r.Reconsideration.Native.Kind = "unobserved"
			case "native-run":
				r.Reconsideration.Native.RunID = "different-run"
			case "native-task":
				r.Reconsideration.Native.HumanQA.TaskID = "different-task"
			case "native-pass":
				r.Reconsideration.Native.Accepted = true
			case "proof-bytes":
				if err := os.Remove(r.Reconsideration.Proof.Locator); err != nil {
					t.Fatal(err)
				}
			case "dispatched":
				r.IntegrationStarted = true
			case "no-native":
				r.Reconsideration.Native = nil
			case "decision-answer":
				r.Reconsideration.Decision.Human.Answer = "pass"
			}
			// Keep the proof's own digest aligned in identity cases to exercise
			// semantic validation, not merely the generic checksum rejection.
			if strings.HasPrefix(name, "native-") {
				proof := *r.Reconsideration
				proof.Proof = nil
				if err := os.Remove(r.Reconsideration.Proof.Locator); err != nil {
					t.Fatal(err)
				}
				if err := atomicJSON(r.Reconsideration.Proof.Locator, proof, true); err != nil {
					t.Fatal(err)
				}
				r.Reconsideration.Proof.SHA256 = digest(proof)
			}
			if err := atomicJSON(receiptPath(r.Plan.Workspace, r.ID), r, false); err != nil {
				t.Fatal(err)
			}
			if _, err := readReceipt(r.Plan.Workspace, r.ID); err == nil {
				t.Fatal("cancelled plan trusted invalid native revocation evidence")
			}
		})
	}
}

func TestReconsiderationRecoversAfterNativeNegativeBeforeRevocationPublication(t *testing.T) {
	r := revokedFixture(t, reconsiderReceipt(t, adapterPhysicalTemp(t), "original"))
	proofPath := r.Reconsideration.Proof.Locator
	if err := os.Remove(proofPath); err != nil {
		t.Fatal(err)
	}
	r.Reconsideration.Proof, r.State = nil, "reconsidering"
	if err := atomicJSON(receiptPath(r.Plan.Workspace, r.ID), r, false); err != nil {
		t.Fatal(err)
	}
	before := digest(r.Reconsideration.Decision)
	s := NewService(taskrun.Dependencies{}, Options{Fault: func(phase string) error {
		if phase == "after_reconsideration_native" {
			return errors.New("injected stop before proof publication")
		}
		return nil
	}})
	if err := s.finishReconsider(taskrun.Dependencies{}, &r); err == nil {
		t.Fatal("interruption was not exercised")
	}
	stored, err := readReceipt(r.Plan.Workspace, r.ID)
	if err != nil || stored.State == "cancelled" || stored.Reconsideration.Proof != nil {
		t.Fatalf("incomplete revocation released authority: state=%s err=%v", stored.State, err)
	}
	s.options.Fault = nil
	if err = s.finishReconsider(taskrun.Dependencies{}, &stored); err != nil {
		t.Fatal(err)
	}
	stored, err = readReceipt(r.Plan.Workspace, r.ID)
	if err != nil || stored.State != "cancelled" || stored.Reconsideration.Proof == nil || digest(stored.Reconsideration.Decision) != before {
		t.Fatalf("same preserved negative answer did not complete revocation: state=%s err=%v", stored.State, err)
	}
	if err = s.finishReconsider(taskrun.Dependencies{}, &stored); err != nil || len(stored.Events) != 1 {
		t.Fatalf("recovery rewrote the negative decision or duplicated cancellation: %v events=%d", err, len(stored.Events))
	}
}
