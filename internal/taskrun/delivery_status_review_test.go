package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeliveryStatusDoesNotPresentMissingPublicationChainAsCurrent(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "qualified required publication")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, deliveryTestReview(t, f, "exact"))
	if err != nil {
		t.Fatal(err)
	}
	r, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range r.TaskContentPublications {
		if p.TaskID != o.Delivery.Candidates[0].TaskResult.TaskID {
			continue
		}
		path := filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "requests", "sha256", strings.TrimPrefix(p.RequestSHA256, "sha256:")+".json")
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
		if err != nil {
			t.Fatal(err)
		}
		status := observedDeliveryStatus(t, shown)
		if status["qualified_candidate"].(map[string]any)["current"] != false {
			t.Fatal("missing required publication request still labels qualification current")
		}
		if status["verification"].(map[string]any)["outcome"] != "passed" {
			t.Fatal("missing bookkeeping dependency erased completed verifier evidence")
		}
		return
	}
	t.Fatal("required Task publication missing from fixture")
}
