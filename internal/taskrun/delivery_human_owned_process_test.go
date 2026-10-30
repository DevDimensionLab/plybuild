package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHumanOwnedVerificationNeedsNoReview(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	oid := deliveryTestCandidate(t, f, "human selected verification")
	got, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil {
		for cause := err; cause != nil; cause = errors.Unwrap(cause) {
			t.Logf("verification cause: %v", cause)
		}
		t.Fatalf("explicit verification required a review: %v", err)
	}
	if len(got.Delivery.Candidates) != 1 || got.Delivery.Candidates[0].OID != oid {
		t.Fatalf("actual candidate was not recorded: %+v", got.Delivery)
	}
	if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 {
		t.Fatal("verification started another agent")
	}
	result := got.Delivery.Candidates[0].TaskResult
	if result.HandoffID != "" || result.StartReceiptID != "" || result.TerminalResultID != "" || got.Delivery.Candidates[0].HumanQA != nil {
		t.Fatal("verification manufactured qualification evidence")
	}
}

func TestHumanOwnedFollowDoesNotCleanup(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	ownerPaneStandin(t, &f)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "observed candidate")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err != nil {
		t.Fatal(err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"agent_status": "idle"})
	_, _ = WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 200*time.Millisecond)
	if _, err := os.Stat(filepath.Join(f.Prepared.Preparation.Preparation.Plan.WorktreePath, "README.md")); err != nil {
		t.Fatalf("observation removed the source: %v", err)
	}
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
		t.Fatal("observation closed the owner's pane")
	}
}
