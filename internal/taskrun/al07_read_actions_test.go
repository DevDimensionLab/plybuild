package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
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
func al07ReservedCloseout(t *testing.T) (deliveryFixture, WorkflowRun, workspace.TaskCloseoutReceipt) {
	t.Helper()
	f, o, _, _ := autoGateFixture(t, false, "")
	d := f.D
	d.Fault = func(point string) error {
		if point == "automatic_closeout_after_release" {
			return errors.New("fixture interruption")
		}
		return nil
	}
	o, err := WorkflowDeliveryIntegrate(d, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if err == nil || o.Delivery.OwnershipRelease == nil {
		t.Fatalf("missing interrupted integration: %v", err)
	}
	c := o.Delivery.Candidates[0]
	release := *o.Delivery.OwnershipRelease
	in := workspace.TaskCloseoutInput{OperationID: "al07-preserved-operation", TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ExpectedResultOID: c.OID, ExpectedResultTree: c.Tree, TargetRef: f.R.Delivery.Agreement.TargetRef, Keep: true, Ownership: workspace.TaskCloseoutOwnership{Released: true, TaskID: c.TaskResult.TaskID, TaskResultID: c.TaskResult.ID, ResultOID: c.OID, ResultTree: c.Tree, EvidenceLocator: release.Locator, EvidenceSHA256: release.SHA256}}
	p, err := workspace.PreviewTaskCloseout(f.D.Workspace, in)
	if err != nil || !p.Ready {
		t.Fatalf("reservation preview: %+v %v", p, err)
	}
	r, err := workspace.ReserveTaskCloseout(f.D.Workspace, in, p.PlanSHA256, "awaiting_closeout", func() error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	return f, o, r
}

func TestAL07ConsumedInventoryRecovery(t *testing.T) {
	f, o, r := al07ReservedCloseout(t)
	for _, phase := range []string{"awaiting_closeout", "installation_pending", "install_failed", "install_unknown"} {
		if _, err := workspace.SetTaskCloseoutPhase(f.D.Workspace, r.Plan.TaskID, r.OperationID, r.PlanSHA256, phase); err != nil {
			t.Fatal(err)
		}
		row := deliveryInventoryRows(t, f)[o.RunID]
		a := row.Delivery.NextAction
		if row.Freshness != "fresh" || a.Actor != "agent" || a.RecordedActor != "coordinator" || len(a.EvidenceIDs) != 1 || a.EvidenceIDs[0] != r.OperationID {
			t.Fatalf("%s: misdirected recovery: %+v reasons %+v", phase, a, row.Reasons)
		}
		if phase == "install_unknown" && (a.Kind != "inspect_installation" || !strings.Contains(a.Reason, "Observe")) {
			t.Fatalf("unknown installation advises replay: %+v", a)
		}
		shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
		if err != nil || shown.NextAction.Actor != "coordinator" || (phase == "install_unknown" && !strings.Contains(shown.NextAction.Message, "Observe")) {
			t.Fatalf("show misdirected recovery: %+v %v", shown.NextAction, err)
		}
	}
	if _, err := workspace.SetTaskCloseoutPhase(f.D.Workspace, r.Plan.TaskID, r.OperationID, r.PlanSHA256, "awaiting_closeout"); err != nil {
		t.Fatal(err)
	}
	p := r.Plan
	in := workspace.TaskCloseoutInput{OperationID: r.OperationID, TaskID: p.TaskID, TaskResultID: p.TaskResultID, ExpectedResultOID: p.ResultOID, ExpectedResultTree: p.ResultTree, TargetRef: p.TargetRef, Keep: p.Keep, Ownership: p.Ownership}
	if _, err := workspace.ApplyTaskCloseout(f.D.Workspace, in, r.PlanSHA256); err != nil {
		t.Fatal(err)
	}
	row := deliveryInventoryRows(t, f)[o.RunID]
	if row.State != "completed" || row.Unresolved || row.Delivery.NextAction.Kind != "task_closed" || row.Delivery.NextAction.Actor == "human" {
		t.Fatalf("closed Task still requires work: %+v", row)
	}
}

func TestAL07CandidateCLIReadActionsDoNotWrite(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o, r := al07ReservedCloseout(t)
	al07AssertWorkspaceRecovery(t, binary, f, r.OperationID, "agent", "awaiting_closeout")
	before := deliveryCLISnapshot(t, filepath.Join(f.R.WorkspaceRoot, ".ply"))
	for _, format := range []string{"json", "text"} {
		c := exec.Command(binary, "workflow", "run", "show", o.RunID, "--format", format)
		c.Dir = f.Parent
		out, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("show %s: %v %s", format, err, out)
		}
		if format == "json" {
			var run WorkflowRun
			if err = json.Unmarshal(out, &run); err != nil || run.NextAction.Actor != "coordinator" {
				t.Fatalf("JSON action: %s %v", out, err)
			}
		} else if !bytes.Contains(out, []byte("Next: coordinator")) || !bytes.Contains(out, []byte("closeout")) {
			t.Fatalf("text actor: %s", out)
		}
	}
	if after := deliveryCLISnapshot(t, filepath.Join(f.R.WorkspaceRoot, ".ply")); after != before {
		t.Fatal("read wrote native sources")
	}
}

func TestAL07RetiredReviewHistoryHasNoAvailableReviewer(t *testing.T) {
	f, o, _ := reviewRoundFixture(t)
	reviewRoundCandidate(t, f, "retired candidate")
	d := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, err := WorkflowDeliveryRequestReview(d, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "fixture root"})
	if err != nil {
		t.Fatal(err)
	}
	round := reviewRoundLatest(t, requested)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	// Synthetic observed retirement, used only to test the read projection.
	record := DeliveryRetirement{Kind: "ply.workflow.owner-retirement", SchemaVersion: 1, RunID: o.RunID, RequestSHA256: s.Result.RequestSHA256, NativeSession: s.Result.Transport.AgentSessionID, PaneClosed: true}
	if _, err = deliveryAppendEvent(f.D, &s, "owner-retirement", "owner_retirement", record); err != nil {
		t.Fatal(err)
	}
	s.Result.Delivery.Phase = "stopped"
	if err = workflowSave(f.D, s); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(round.Request.Locator)
	if err != nil {
		t.Fatal(err)
	}
	shown, err := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	review := shown.DeliveryStatus.Review
	if shown.NextAction.Actor != "coordinator" || review.Next.Actor != "none" || strings.Contains(review.Next.Message, "is open") {
		t.Fatalf("retired reviewer available: overall %+v review %+v", shown.NextAction, review)
	}
	if len(review.Rounds) != 1 || review.Rounds[0].State != "requested" || review.Rounds[0].Decision != "" || review.Rounds[0].ReviewerPane != round.ReviewerPane {
		t.Fatalf("review history changed: %+v", review)
	}
	after, err := os.ReadFile(round.Request.Locator)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("review read changed history")
	}
}

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
