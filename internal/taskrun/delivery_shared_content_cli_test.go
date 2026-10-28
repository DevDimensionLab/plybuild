package taskrun

import (
	"os"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestSharedCloseoutCurrentAndPreviousCLIReadback(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o, report := informativeCLIFixture(t, binary)
	assertAutomaticCLIClosure(t, binary, f, o)
	assertInformativeRetained(t, f, o, report)
	c := o.Delivery.Candidates[0]
	closed, err := workspace.ReadTaskCloseoutAt(f.R.WorkspaceRoot, c.TaskResult.TaskID)
	if err != nil || closed == nil || closed.Retention == nil {
		t.Fatalf("missing closeout: %v", err)
	}
	shared := 0
	for _, e := range closed.Retention.Entries {
		origin, err := os.Stat(e.OriginLocator)
		if err != nil {
			continue
		}
		retained, err := os.Stat(e.RetainedLocator)
		if err == nil && os.SameFile(origin, retained) {
			shared++
		}
	}
	if shared == 0 {
		t.Fatal("built CLI still copied all immutable workspace input")
	}
	old := os.Getenv("PLY_CLOSEOUT_COMPAT_READER")
	if old == "" {
		t.Log("previous CLI not supplied; current CLI closure and sharing verified")
		return
	}
	readback, output, err := cleanupCLIAt(old, f.Parent, "execute", "follow", o.RunID, "--timeout", "5")
	if err != nil || readback.Delivery.Phase != "completed" || readback.DeliveryStatus.FinalDelivery.State != "delivered" {
		t.Fatalf("previous CLI could not read shared retention: %v\n%s", err, output)
	}
	t.Logf("current and previous CLI read completed closeout with %d shared immutable references", shared)
}
