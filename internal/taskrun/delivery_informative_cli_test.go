package taskrun

import (
	"bytes"

	"os"
	"path/filepath"

	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Real frozen Goal/Spec/Handoff and public CLI verification/integration/retention.
// Provider, runtime claim and review are synthetic; this is no Claude acceptance.

func assertInformativeRetained(t *testing.T, f deliveryFixture, o WorkflowRun, report []byte) {
	t.Helper()
	c := o.Delivery.Candidates[0]
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	resource := map[bool]string{true: "kept", false: "retired"}[closeout != nil && closeout.Plan.Keep]
	if err != nil || closeout == nil || closeout.State != "complete" || closeout.ResourceState != resource || closeout.WorktreeRemoved == closeout.Plan.Keep || closeout.BranchRemoved == closeout.Plan.Keep || closeout.Retention == nil {
		t.Fatalf("native closeout with retained evidence incomplete: %+v %v", closeout, err)
	}
	found := false
	for _, e := range closeout.Retention.Entries {
		if e.OriginLocator == filepath.Join(f.R.WorkspaceRoot, "historical-acceptance.sh") {
			t.Fatal("historical citation became a new required dependency")
		}
		if e.SHA256 != hash(report) {
			continue
		}
		found = true
		if e.SizeBytes != int64(len(report)) {
			t.Fatal("report size changed")
		}
		got, err := workspace.ReadTaskRetainedEvidence(f.R.WorkspaceRoot, c.TaskResult.TaskID, e.OriginLocator, e.SHA256)
		if err != nil || !bytes.Equal(got, report) {
			t.Fatal("exact original binding did not resolve retained report", err)
		}
		original, err := os.ReadFile(e.OriginLocator)
		if err != nil || !bytes.Equal(original, report) {
			t.Fatal("frozen report rewritten", err)
		}
	}
	if !found {
		t.Fatal("report missing from retained evidence")
	}
	if _, err := os.Lstat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); closeout.Plan.Keep && err != nil || !closeout.Plan.Keep && !os.IsNotExist(err) {
		t.Fatal("source state differs from the frozen cleanup grant", err)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("wrong integrated target")
	}
	base, err := workspace.ShowEpicBase(f.D.Workspace, workspace.EpicBaseInput{EpicID: "epic", RepoID: "ply"})
	if err != nil || deliveryCLILatestBase(base).OID != c.OID {
		t.Fatal("Epic base did not advance to exact candidate", err)
	}
	t.Logf("Native closeout: run=%s result=%s candidate=%s manifest=%s; report bytes/hash/size retained, source %s", o.RunID, c.TaskResult.ID, c.OID, closeout.Retention.ManifestSHA256, closeout.ResourceState)
}

// Set PLY_INFORMATIVE_OLD_READER to the unmodified base CLI. It creates the
// actual already-integrated retention failure, then the candidate uses native
// same-Run continuation. No private state is edited to construct this boundary.
