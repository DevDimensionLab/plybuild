package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Native records, Git, review sealing and executable acceptance are real;
// owner/reviewer/Herdr/absence observations use the inspected fake runtime.
func TestParallelCancelUnstartedThenFreshQualification(t *testing.T) {
	parallelCancelUnstartedThenFreshQualification(t, false)
}

func TestParallelDriverRefusalThenCancelFreshQualification(t *testing.T) {
	parallelCancelUnstartedThenFreshQualification(t, true)
}

func parallelCancelUnstartedThenFreshQualification(t *testing.T, driverRefusal bool) {
	t.Helper()
	f, owner, sibling, other := parallelOwnerPair(t)
	source := f.Prepared.Preparation.Plan.WorktreePath
	if e := os.WriteFile(filepath.Join(source, "owned.txt"), []byte("fixed owned source\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runGit(t, source, "add", "owned.txt")
	runGit(t, source, "commit", "-m", "Fixed owned candidate")
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		binary = filepath.Join(f.R.WorkspaceRoot, "simulated-candidate")
		if e := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0700); e != nil {
			t.Fatal(e)
		}
	}
	if !filepath.IsAbs(binary) {
		t.Fatal("candidate must be absolute; no installed fallback")
	}
	if e := os.WriteFile(f.AcceptancePath, []byte("#!/bin/sh\nset -eu\ngrep -q fixed owned.txt\n\"$PLY_CANDIDATE_BINARY\" version --format json\n"), 0600); e != nil {
		t.Fatal(e)
	}
	owner = parallelQualify(t, f, owner, binary)
	originalCandidate := owner.Delivery.Candidates[0]
	originalRequest, e := os.ReadFile(f.File)
	if e != nil {
		t.Fatal(e)
	}
	otherBefore, e := workflowRead(sibling.R.WorkspaceRoot, other.RunID)
	if e != nil {
		t.Fatal(e)
	}
	moveParent := func(message string) {
		t.Helper()
		runGit(t, f.Parent, "commit", "--allow-empty", "-m", message)
		p, e := workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
		if e != nil {
			t.Fatal(e)
		}
		if _, e = workspace.UpdateEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply", Apply: true, Confirmation: *p.Confirmation}); e != nil {
			t.Fatal(e)
		}
	}
	reviewRoundIn(t, owner.Transport.PaneID)
	moveParent("Parent one")
	p, e := WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, owner.RunID, owner.Paths.Context, workspace.TaskSourceReconciliationInput{})
	if e != nil {
		t.Fatal(e)
	}
	if driverRefusal {
		// Isolated repository-local config only. The driver is rejected before
		// any invocation, even though no merge attribute selects it here.
		runGit(t, source, "config", "--local", "merge.private-fixture.driver", "false")
	} else {
		f.D.Workspace.SourceReconciliationFault = func(stage, id string) error {
			if stage == "after_intent" {
				return errors.New("simulated pre-effect interruption")
			}
			return nil
		}
	}
	intent, e := WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, owner.RunID, owner.Paths.Context, workspace.TaskSourceReconciliationInput{Apply: true, Confirmation: p.Confirmation})
	if e == nil || intent.State != "intent" {
		t.Fatalf("pre-effect intent: %+v %v", intent, e)
	}
	if driverRefusal {
		if !strings.Contains(e.Error(), "source_merge_driver_unsupported") || !intent.Operation.Unstarted {
			t.Fatalf("driver refusal consumed positive pre-effect proof: %+v %v", intent, e)
		}
		runGit(t, source, "config", "--local", "--unset", "merge.private-fixture.driver")
	}
	f.D.Workspace.SourceReconciliationFault = nil
	moveParent("Parent two after unstarted intent")
	cancel := workspace.TaskSourceReconciliationInput{Apply: true, OperationID: intent.Operation.ID, Action: "cancel-unstarted"}
	if _, e = WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, owner.RunID, other.Paths.Context, cancel); e == nil {
		t.Fatal("foreign owner context cancelled intent")
	}
	before := gitOutput(t, source, "rev-parse", "HEAD")
	cancelled, e := WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, owner.RunID, owner.Paths.Context, cancel)
	if e != nil || cancelled.State != "aborted" || cancelled.Operation.ID != intent.Operation.ID || gitOutput(t, source, "rev-parse", "HEAD") != before {
		t.Fatalf("native same-owner no-effect cancellation: %+v %v", cancelled, e)
	}
	p, e = WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, owner.RunID, owner.Paths.Context, workspace.TaskSourceReconciliationInput{})
	if e != nil {
		t.Fatal(e)
	}
	merged, e := WorkflowDeliveryReconcileSource(f.D, f.R.WorkspaceRoot, owner.RunID, owner.Paths.Context, workspace.TaskSourceReconciliationInput{Apply: true, Confirmation: p.Confirmation})
	if e != nil || merged.State != "committed" || merged.Operation.Plan.AttemptOrdinal != 2 {
		t.Fatalf("fresh parent reconciliation: %+v %v", merged, e)
	}
	if _, e = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, owner.RunID, owner.Paths.Context, "", binary); e == nil {
		t.Fatal("fresh basis inherited prior review")
	}
	owner = parallelQualify(t, f, owner, binary)
	if len(owner.Delivery.Candidates) != 2 || !equal(owner.Delivery.Candidates[0], originalCandidate) || owner.Delivery.Candidates[1].Verification == originalCandidate.Verification {
		t.Fatal("fresh qualification replaced history or reused old evidence")
	}
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil || len(r.SourceReconciliations) != 2 || !equal(r.SourceReconciliations[0], *cancelled.Operation) {
		t.Fatalf("cancellation history changed: %v", e)
	}
	basis, e := workspace.EffectiveTaskDeliveryBasis(r, f.R.PreparationID)
	if e != nil || basis.Generation != 1 || basis.Parent.OID != gitOutput(t, f.Parent, "rev-parse", "HEAD") {
		t.Fatalf("fresh effective basis: %+v %v", basis, e)
	}
	raw, e := os.ReadFile(f.File)
	if e != nil || string(raw) != string(originalRequest) {
		t.Fatal("cancellation rebound frozen request/mandate")
	}
	otherAfter, e := workflowRead(sibling.R.WorkspaceRoot, other.RunID)
	if e != nil || !equal(otherBefore, otherAfter) {
		t.Fatal("source cancellation changed sibling Run")
	}
	t.Logf("simulated runtime: cancelled=%s fresh=%s generation=%d candidates=%d", cancelled.Operation.ID, merged.Operation.ID, basis.Generation, len(owner.Delivery.Candidates))
}
