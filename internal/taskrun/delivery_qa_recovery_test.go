package taskrun

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Reproduce the historical writer's exact reservation and staging, then run
// the real lower parser. These are isolated synthetic artifacts, not human QA.
func rejectedQAFixture(t *testing.T) (deliveryFixture, WorkflowRun, DeliveryQARecoveryInput, string, []string) {
	t.Helper()
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalEpic}})
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "candidate waiting for its actual human answer")
	var err error
	o, err = WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "qa-recovery"))
	if err != nil {
		t.Fatal(err)
	}
	c := o.Delivery.Candidates[0]
	valid := deliveryTestHuman(t, f, o, "pass")
	var a workflowhandoff.DeliveryHumanAttestation
	if err = readValue(valid, 1<<20, &a); err != nil {
		t.Fatal(err)
	}
	a.StartedAtUTC = strings.TrimSuffix(a.StartedAtUTC, "Z") + "+00:00"
	a.CompletedAtUTC = strings.TrimSuffix(a.CompletedAtUTC, "Z") + "+00:00"
	bad := writeAny(t, f.R.WorkspaceRoot, "legacy-offset-human.json", a)
	binding, raw, err := deliveryReadBinding(bad, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	inputSHA := digest([]any{c.TaskResult.ID, "pass", binding})
	id := "qa-" + inputSHA[7:]
	path := filepath.Join(o.Paths.RunRoot, "delivery", "attempts", id)
	request := map[string]any{"kind": "PlyDeliveryHumanQARequest@1", "run_id": o.RunID, "request_sha256": o.RequestSHA256, "candidate": c.TaskResult.ID, "candidate_oid": c.OID, "candidate_tree": c.Tree, "outcome": "pass", "evidence": binding}
	if err = f.D.writeValue(filepath.Join(path, "attempt.json"), request); err != nil {
		t.Fatal(err)
	}
	err = workflowUpdate(f.D, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error {
		s.Result.Delivery.Attempt = &DeliveryAttempt{ID: id, Kind: "human_qa", State: "attempted", Path: path, CandidateOID: c.OID, CandidateTree: c.Tree, InputSHA256: inputSHA}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	h, err := f.D.Workflow.Store.ReadByLocator(c.TaskResult.HandoffLocator)
	if err != nil {
		t.Fatal(err)
	}
	staging := filepath.Join(h.Handoff.ReplyRoot, "staging", "human-qa-"+binding.SHA256[7:])
	attestation := filepath.Join(staging, "attestation.json")
	if err = f.D.writeOnce(attestation, raw); err != nil {
		t.Fatal(err)
	}
	actor := workspace.TaskHumanActorRecord{ActorClaim: a.ActorClaim, StartSurface: a.StartSurface, StartedAtUTC: a.StartedAtUTC, CompletedAtUTC: a.CompletedAtUTC}
	draft := map[string]any{"kind": "WorkspaceTaskHumanQARecordDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "delivery-qa/" + string(c.TaskResult.ID) + "/" + binding.SHA256[7:], "task_id": c.TaskResult.TaskID, "task_result_id": c.TaskResult.ID, "result_oid": c.OID, "result_tree": c.Tree, "outcome": "pass", "actor": actor, "evidence": []workspace.TaskHumanQAEvidenceRecord{{ID: "human-attestation", Role: "report", Locator: attestation, SHA256: binding.SHA256, SizeBytes: int64(len(raw))}}, "observation": a.Observation, "accepted_residual_risks": []workspace.TaskResidualRiskRecord{}}
	draftPath := filepath.Join(staging, "qa-draft.json")
	if err = f.D.writeValue(draftPath, draft); err != nil {
		t.Fatal(err)
	}
	_, err = workspace.RecordTaskHumanQA(f.D.Workspace, workspace.TaskHumanQARecordInput{TaskID: c.TaskResult.TaskID, File: draftPath})
	if err == nil || !strings.Contains(err.Error(), "invalid human QA actor") {
		t.Fatalf("historical lower validation mismatch missing: %v", err)
	}
	if _, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", valid); err == nil || !strings.Contains(err.Error(), "another delivery effect is unresolved") {
		t.Fatalf("historical reservation did not block corrected input: %v", err)
	}
	return f, o, DeliveryQARecoveryInput{RunID: o.RunID, ContextPath: o.Paths.Context, AttemptID: id, RuntimeEvidencePath: deliveryTestRuntimeObservation(t, f, o)}, valid, []string{bad, attestation, draftPath, filepath.Join(path, "attempt.json"), f.R.Runtime.PlyExecutable.Path}
}

func TestDeliveryQARecoveryInterruptedAndConcurrent(t *testing.T) {
	for _, point := range []string{"delivery_qa_recovery_after_receipt", "delivery_qa_recovery_after_event", "delivery_qa_recovery_after_state"} {
		t.Run(point, func(t *testing.T) {
			f, o, in, _, _ := rejectedQAFixture(t)
			p, err := WorkflowPreviewDeliveryQARecovery(f.D, f.R.WorkspaceRoot, in)
			if err != nil {
				t.Fatal(err)
			}
			fault := f.D
			fault.Fault = func(at string) error {
				if at == point {
					return errors.New("synthetic interrupted recovery")
				}
				return nil
			}
			if _, err = WorkflowRecoverDeliveryQA(fault, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
				t.Fatal("fault did not interrupt recovery")
			}
			var wg sync.WaitGroup
			errs := make(chan error, 2)
			for i := 0; i < 2; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, e := WorkflowRecoverDeliveryQA(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
					errs <- e
				}()
			}
			wg.Wait()
			close(errs)
			for e := range errs {
				if e != nil {
					t.Fatal(e)
				}
			}
			s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			n := 0
			for _, e := range s.Result.Delivery.Events {
				if e.ID == "recover-"+in.AttemptID {
					n++
				}
			}
			if n != 1 || s.Result.Delivery.Attempt != nil || s.Result.Delivery.Candidates[0].HumanQA != nil {
				t.Fatal("recovery duplicated, retained rejection or invented QA")
			}
		})
	}
}

func TestDeliveryQARecoveryRejectsUnprovenEffectsAndAuthority(t *testing.T) {
	for _, kind := range []string{"wrong-attempt", "context", "session", "policy", "stale-observation", "dirty-candidate", "changed-evidence", "changed-draft", "published-qa", "other-effect", "changed-control", "changed-candidate"} {
		t.Run(kind, func(t *testing.T) {
			f, o, in, valid, paths := rejectedQAFixture(t)
			d := f.D
			switch kind {
			case "wrong-attempt":
				in.AttemptID = "qa-" + strings.Repeat("0", 64)
			case "context":
				in.ContextPath = valid
			case "session", "policy", "stale-observation":
				var v DeliveryRuntimeObservation
				if err := readValue(in.RuntimeEvidencePath, 1<<20, &v); err != nil {
					t.Fatal(err)
				}
				if kind == "session" {
					v.RuntimeClaim.NativeSessionID = ptr("different-session")
				}
				if kind == "policy" {
					v.DeliveryPermission.PermissionConfirmed = false
				}
				if kind == "stale-observation" {
					v.ObservedAtUTC = "2000-01-01T00:00:00Z"
				}
				in.RuntimeEvidencePath = writeAny(t, f.R.WorkspaceRoot, "wrong-observation.json", v)
			case "dirty-candidate":
				if err := os.WriteFile(filepath.Join(f.Prepared.Preparation.Preparation.Plan.WorktreePath, "dirty.txt"), []byte("unreviewed\n"), 0600); err != nil {
					t.Fatal(err)
				}
			case "changed-evidence", "changed-draft", "changed-control":
				index := 0
				if kind == "changed-draft" {
					index = 2
				}
				if kind == "changed-control" {
					index = 4
				}
				if err := os.WriteFile(paths[index], []byte("changed artifact\n"), 0600); err != nil {
					t.Fatal(err)
				}
			case "published-qa":
				if _, err := workflowhandoff.RecordDeliveryHumanQA(d.Workflow, "task", o.Delivery.Candidates[0].TaskResult, "pass", valid); err != nil {
					t.Fatal(err)
				}
			case "other-effect":
				if err := workflowUpdate(d, f.R.WorkspaceRoot, o.RunID, func(s *workflowState) error { s.Result.Delivery.Attempt.Kind = "integration"; return nil }); err != nil {
					t.Fatal(err)
				}
			case "changed-candidate":
				deliveryTestCandidate(t, f, "later candidate without qualification")
			}
			statePath := filepath.Join(o.Paths.RunRoot, "state.json")
			before, err := os.ReadFile(statePath)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = WorkflowPreviewDeliveryQARecovery(d, f.R.WorkspaceRoot, in); err == nil {
				t.Fatal("unsafe recovery preview allowed")
			}
			if _, err = WorkflowRecoverDeliveryQA(d, f.R.WorkspaceRoot, in, "sha256:"+strings.Repeat("0", 64)); err == nil {
				t.Fatal("unsafe recovery applied")
			}
			after, _ := os.ReadFile(statePath)
			if !bytes.Equal(before, after) {
				t.Fatal("blocked recovery changed state")
			}
		})
	}
}

func TestDeliveryQARecoveryPreservesEvidenceAndAllowsOrdinaryQA(t *testing.T) {
	f, o, in, corrected, paths := rejectedQAFixture(t)
	frozen := map[string][]byte{}
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		frozen[path] = raw
	}
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	before, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	p, err := WorkflowPreviewDeliveryQARecovery(f.D, f.R.WorkspaceRoot, in)
	if err != nil || p.State != "ready" {
		t.Fatalf("known rejected QA has no supported recovery: %+v %v", p, err)
	}
	after, _ := os.ReadFile(statePath)
	if !bytes.Equal(before, after) {
		t.Fatal("preview mutated state")
	}
	p, err = WorkflowRecoverDeliveryQA(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil || p.State != "recovered" || p.Receipt == nil {
		t.Fatalf("recovery failed: %+v %v", p, err)
	}
	if strings.HasPrefix(p.NextAction, "Confirm") {
		t.Fatal("completed recovery still asks to confirm the completed transition")
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || s.Result.Delivery.Attempt != nil || s.Result.Delivery.Candidates[0].HumanQA != nil {
		t.Fatalf("recovery fabricated human QA or failed to retire rejection: %v", err)
	}
	for path, want := range frozen {
		raw, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(raw, want) {
			t.Fatalf("original bytes changed at %s: %v", path, err)
		}
	}
	if _, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context); err == nil {
		t.Fatal("recovery granted integration without QA")
	}
	o, err = WorkflowDeliveryQA(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "pass", corrected)
	if err != nil || o.Delivery.Candidates[0].HumanQA == nil {
		t.Fatalf("ordinary original control QA failed after recovery: %v", err)
	}
	before, _ = os.ReadFile(statePath)
	p, err = WorkflowRecoverDeliveryQA(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil || p.State != "existing" {
		t.Fatalf("same recovery not idempotent after later QA: %v", err)
	}
	after, _ = os.ReadFile(statePath)
	if !bytes.Equal(before, after) {
		t.Fatal("old recovery changed later QA state")
	}
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil || o.Delivery.Phase != "completed" {
		t.Fatalf("normal integration failed: %v", err)
	}
}
