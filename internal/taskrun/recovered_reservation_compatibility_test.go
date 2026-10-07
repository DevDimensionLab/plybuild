package taskrun

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestWorkflowReservationsReadRecoveredRunBeforeStartingAnotherTask(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	s, _ := deliveryInventoryRecovery(t, f)
	statePath := filepath.Join(s.Result.Paths.RunRoot, "state.json")
	before, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	other := workspace.PlanWorktreeObservation{
		WorktreeLocator: filepath.Join(f.R.WorkspaceRoot, "another-task"),
		GitCommonDir:    s.Observed.Target.GitCommonDir,
		Ref:             "refs/heads/another-task",
	}
	if err := workflowReservations(f.D, f.R.WorkspaceRoot, other); err != nil {
		t.Fatalf("a preserved recovered run must remain readable before another Task starts: %v", err)
	}
	if err := workflowReservations(f.D, f.R.WorkspaceRoot, s.Observed.Target); err == nil {
		t.Fatal("reading recovery metadata must not release the original active target")
	}
	after, err := os.ReadFile(statePath)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("reservation inspection changed the recovered run: %v", err)
	}
}
