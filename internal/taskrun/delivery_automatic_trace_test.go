package taskrun

import (
	"encoding/json"
	"os"
	"testing"
)

func TestAutomaticDeliveryTracePreservesNativeVerificationAndRetainedCloseout(t *testing.T) {
	f, o, binary, _ := autoGateFixture(t, false, "")
	var err error
	o, err = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil {
		t.Fatal(err)
	}
	// The trace reads historical bytes and cannot depend on the original live
	// acceptance program or candidate executable after a completed delivery.
	for _, path := range []string{binary, f.AcceptancePath} {
		if err = os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
	run := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	checks := traceKind(run, "verification")
	if len(checks) != 1 || checks[0].Outcome == nil || *checks[0].Outcome != "passed" || checks[0].Verification == nil || !checks[0].Verification.InputsBound {
		t.Fatalf("trace lost the automatic verifier evidence/outcome: checks=%+v reasons=%+v", checks, run.Reasons)
	}
	var verification map[string]json.RawMessage
	if err = json.Unmarshal(checks[0].Data, &verification); err != nil || len(verification["automatic"]) == 0 {
		t.Fatalf("trace omitted the automatic invocation's policy and exact executable: %s (%v)", checks[0].Data, err)
	}
	if checks[0].Role != "ply" || checks[0].EvidenceClass != "controlled_verification" || checks[0].Candidate == nil || checks[0].Candidate.TaskResultID == nil || *checks[0].Candidate.TaskResultID != string(o.Delivery.Candidates[0].TaskResult.ID) {
		t.Fatalf("trace lost automatic verification provenance: %+v", checks[0])
	}
	if len(traceKind(run, "human_qa")) != 0 {
		t.Fatal("automatic trace invented a human answer")
	}
	integration := traceKind(run, "integration")
	if len(integration) != 1 || integration[0].Role != "ply" || integration[0].NativeID == nil {
		t.Fatalf("trace lost the native automatic integration: entries=%+v reasons=%+v", integration, run.Reasons)
	}
	closed := traceKind(run, "task_closeout")
	if len(closed) != 1 || closed[0].Role != "ply" || closed[0].Outcome == nil || *closed[0].Outcome != "complete" || closed[0].Candidate == nil || closed[0].Candidate.OID != o.Delivery.Candidates[0].OID {
		t.Fatalf("trace omitted native automatic Task closeout: entries=%+v reasons=%+v", closed, run.Reasons)
	}
	var closeout struct {
		ResourceState      string `json:"resource_state"`
		LifecycleCompleted bool   `json:"lifecycle_completed"`
		WorktreeRemoved    bool   `json:"worktree_removed"`
		BranchRemoved      bool   `json:"branch_removed"`
	}
	if err = json.Unmarshal(closed[0].Data, &closeout); err != nil || closeout.ResourceState != "kept" || !closeout.LifecycleCompleted || closeout.WorktreeRemoved || closeout.BranchRemoved {
		t.Fatalf("trace misstates retained automatic closeout: %s (%v)", closed[0].Data, err)
	}
	releases := traceKind(run, "source_ownership_release")
	if len(releases) != 1 || releases[0].Role != "ply" || releases[0].ActorClaim != "ply automatic delivery closeout" {
		t.Fatalf("trace misattributes native automatic owner release: %+v", releases)
	}
	if run.Coverage != "complete" {
		t.Fatalf("valid automatic evidence became incomplete: %+v", run.Reasons)
	}
}
