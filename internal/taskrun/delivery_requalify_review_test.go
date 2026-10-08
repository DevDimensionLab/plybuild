package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeliveryAlreadyQualifiedReuseStillRequiresPublicationChain(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "qualified required publication")
	review := deliveryTestReview(t, f, "exact")
	o, err := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review)
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
		_, err := WorkflowDeliveryRequalify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, review, o.Delivery.Attempt.ID)
		if err == nil {
			t.Fatal("already-qualified receipt reuse ignored missing required publication request")
		}
		return
	}
	t.Fatal("required Task publication missing from fixture")
}
