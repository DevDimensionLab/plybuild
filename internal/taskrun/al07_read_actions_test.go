package taskrun

import (
	"encoding/json"

	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAL07StartupActionMatchesFrozenProtocol(t *testing.T) {
	for _, receipt := range []bool{false, true} {
		t.Run(map[bool]string{false: "legacy", true: "receipt"}[receipt], func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			if receipt {
				f = receiptFixture(t, "codex")
			}
			o := workflowTestStart(t, f.workflowFixture)
			if o.NextAction.Actor != "recipient" {
				t.Fatalf("wrong startup actor: %+v", o.NextAction)
			}
			if receipt && (!strings.Contains(o.NextAction.Message, "check-in") || strings.Contains(o.NextAction.Message, "Accept the bound runtime")) {
				t.Fatalf("receipt advised legacy acceptance: %+v", o.NextAction)
			}
			if !receipt && !strings.Contains(o.NextAction.Message, "Accept the bound runtime") {
				t.Fatalf("legacy instruction changed: %+v", o.NextAction)
			}
			row := deliveryInventoryRows(t, f)[o.RunID]
			if row.Delivery.NextAction.Reason != o.NextAction.Message {
				t.Fatalf("inventory disagrees: %+v", row.Delivery.NextAction)
			}
		})
	}
}

// Native integration and reservation provide the consumed recovery fixture.
// Phase changes model installer observations; these reads never run installation.

func al07AssertWorkspaceRecovery(t *testing.T, binary string, f deliveryFixture, operation, actor, phase string) {
	t.Helper()
	before := deliveryCLISnapshot(t, filepath.Join(f.R.WorkspaceRoot, ".ply"))
	c := exec.Command(binary, "workspace", "task", "show", "task", "--progress", "--format", "json")
	c.Dir = f.Parent
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("workspace progress: %v %s", err, out)
	}
	var read struct {
		Task struct {
			Progress struct {
				NextActions []struct {
					Actor, Kind, Reason string
					EvidenceIDs         []string `json:"evidence_ids"`
				} `json:"next_actions"`
			} `json:"progress"`
		} `json:"task"`
	}
	if err = json.Unmarshal(out, &read); err != nil {
		t.Fatal(err)
	}
	actions := read.Task.Progress.NextActions
	if len(actions) != 1 || actions[0].Actor != actor || len(actions[0].EvidenceIDs) != 1 || actions[0].EvidenceIDs[0] != operation {
		t.Fatalf("workspace disagrees with frozen Run actor: %s", out)
	}
	if phase == "install_unknown" && (actions[0].Kind != "inspect_installation" || !strings.Contains(actions[0].Reason, "Observe")) {
		t.Fatalf("workspace advises unknown install replay: %s", out)
	}
	if after := deliveryCLISnapshot(t, filepath.Join(f.R.WorkspaceRoot, ".ply")); after != before {
		t.Fatal("workspace read changed native sources")
	}
}
