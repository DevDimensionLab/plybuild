package taskrun

import (
	"bytes"

	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Real public CLI verification/integration/closeout. The acceptance prints a
// native-looking JSON citation of removed test scratch to stdout and stderr.

func assertVerifierLogsRetained(t *testing.T, f deliveryFixture, o WorkflowRun, log []byte) {
	t.Helper()
	c := o.Delivery.Candidates[0]
	closeout, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closeout == nil || closeout.State != "complete" || closeout.ResourceState != "retired" || !closeout.WorktreeRemoved || !closeout.BranchRemoved || closeout.Retention == nil {
		t.Fatalf("native closeout with retained evidence incomplete: %+v %v", closeout, err)
	}
	logs := 0
	for _, e := range closeout.Retention.Entries {
		if strings.HasPrefix(e.OriginLocator, filepath.Join(f.R.WorkspaceRoot, "removed-test-scratch")) {
			t.Fatal("historical log citation became a required dependency")
		}
		if e.SHA256 != hash(log) {
			continue
		}
		logs++
		got, err := workspace.ReadTaskRetainedEvidence(f.R.WorkspaceRoot, c.TaskResult.TaskID, e.OriginLocator, e.SHA256)
		if err != nil || !bytes.Equal(got, log) || e.SizeBytes != int64(len(log)) {
			t.Fatal("exact log binding did not resolve retained bytes", err)
		}
		if original, err := os.ReadFile(e.OriginLocator); err == nil && !bytes.Equal(original, log) {
			t.Fatal("frozen log rewritten")
		}
	}
	// stdout and stderr: the original verifier logs and their candidate copies.
	if logs < 2 {
		t.Fatalf("verifier logs missing from retained evidence: %d", logs)
	}
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != c.OID {
		t.Fatal("wrong integrated target")
	}
	t.Logf("Native closeout: run=%s result=%s candidate=%s manifest=%s; %d exact log entries retained", o.RunID, c.TaskResult.ID, c.OID, closeout.Retention.ManifestSHA256, logs)
}

// Set PLY_INFORMATIVE_OLD_READER to the unmodified base CLI. It integrates;
// its coordinator follow then reaches the actual log-citation retention
// failure. The candidate's follow resumes the same Run's closeout without a
// repeated integration. No private state is edited to construct this boundary.
