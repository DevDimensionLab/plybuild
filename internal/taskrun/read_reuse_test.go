package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Real preserved store bytes exercise validation as history grows. The only
// stand-in counts the domain reader; it still reads and validates that store.
func TestInventoryReadsCloseoutHistoryOnce(t *testing.T) {
	for _, size := range []struct{ runs, receipts int }{{1, 1}, {16, 128}} {
		t.Run(fmt.Sprintf("%d-runs-%d-receipts", size.runs, size.receipts), func(t *testing.T) {
			f := newDeliveryFixture(t, "codex")
			root := f.R.WorkspaceRoot
			plan := workspace.TaskCloseoutPlan{WorkspaceRoot: root, TaskID: "task", TaskResultID: "result", ResultOID: "oid", ResultTree: "tree"}
			plan.Source.Locator, plan.Source.Ref = filepath.Join(root, "source"), "refs/heads/task"
			receipts := []workspace.TaskCloseoutReceipt{}
			for i := 0; i < size.receipts; i++ {
				p := plan
				if i > 0 {
					p.TaskID = workspace.TaskID(fmt.Sprintf("history-%d", i))
				}
				b, _ := json.Marshal(p)
				receipts = append(receipts, workspace.TaskCloseoutReceipt{Kind: "WorkspaceTaskCloseoutReceipt@1", SchemaVersion: 1, OperationID: fmt.Sprintf("closeout-%d", i), Plan: p, PlanSHA256: hash(b), State: "cleanup_pending", ResourceState: "retained", CreatedAtUTC: "2026-10-24T12:00:00Z", UpdatedAtUTC: "2026-10-24T12:00:00Z"})
			}
			store := filepath.Join(root, ".ply", "task-closeouts.json")
			if err := os.WriteFile(store, mustJSON(map[string]any{"kind": "WorkspaceTaskCloseoutStore@1", "schema_version": 1, "receipts": receipts}), 0600); err != nil {
				t.Fatal(err)
			}
			for i := 0; i < size.runs; i++ {
				s := deliveryInventoryState(t, f, fmt.Sprintf("history-%d", i))
				deliveryInventoryReport(t, f, &s, fmt.Sprintf("report-%d", i), "working")
				verification := writeAny(t, root, fmt.Sprintf("verification-%d.json", i), map[string]any{"fixture": true})
				s.Result.Delivery.Candidates = []DeliveryCandidate{{OID: plan.ResultOID, Tree: plan.ResultTree, Verification: FileBinding{verification, hashFileTest(t, verification)}, TaskResult: workspace.TaskResultRecord{ID: plan.TaskResultID, TaskID: plan.TaskID, SourceLocator: plan.Source.Locator, SourceRef: plan.Source.Ref, Artifacts: []workspace.TaskArtifactRecord{}, VerifierResults: []workspace.TaskVerifierResultRecord{}, AcceptedDebt: []workspace.TaskAcceptedDebtRecord{}, Review: workspace.TaskReviewRecord{Findings: []workspace.TaskReviewEntry{}, Fixes: []workspace.TaskReviewEntry{}, OpenActionableFindings: []workspace.TaskReviewEntry{}}}}}
				if err := workflowSave(f.D, s); err != nil {
					t.Fatal(err)
				}
			}
			reads := 0
			read := func(root string) ([]workspace.TaskCloseoutReceipt, error) {
				reads++
				return workspace.ReadTaskCloseoutsAt(root)
			}
			for request := 1; request <= 2; request++ {
				out, err := readInventory(f.D, root, read)
				if err != nil || len(out.Runs) != size.runs {
					t.Fatalf("inventory: %v %+v", err, out)
				}
				if reads != request {
					t.Errorf("request %d closeout store reads=%d; want total %d", request, reads, request)
				}
				for _, row := range out.Runs {
					if row.Closeout == nil || row.Closeout.OperationID != "closeout-0" || row.Freshness != "fresh" {
						t.Error("exact receipt lost", row.RunID, row.Closeout, row.Reasons)
					}
				}
			}
			// A later request must read damaged history rather than reuse prior facts.
			if err := os.WriteFile(store, []byte("damaged"), 0600); err != nil {
				t.Fatal(err)
			}
			out, err := readInventory(f.D, root, read)
			if err != nil {
				t.Fatal(err)
			}
			if reads != 3 {
				t.Errorf("damaged history read %d times total; want 3", reads)
			}
			for _, row := range out.Runs {
				if row.Closeout != nil || row.Freshness == "fresh" || row.StatusFollowupVeto == "clear" {
					t.Error("damaged source became fresh/clear", row)
				}
			}
			// Absence and subsequent creation are also request-local facts.
			if err := os.Remove(store); err != nil {
				t.Fatal(err)
			}
			out, err = readInventory(f.D, root, read)
			if err != nil || reads != 4 {
				t.Fatalf("absence: reads=%d err=%v", reads, err)
			}
			for _, row := range out.Runs {
				if row.Closeout != nil || row.Freshness != "fresh" {
					t.Error("absent store is not empty", row)
				}
			}
			// Valid store with a different exact source must still fail closed.
			receipts[0].Plan.Source.Ref = "refs/heads/another-task"
			b, _ := json.Marshal(receipts[0].Plan)
			receipts[0].PlanSHA256 = hash(b)
			if err := os.WriteFile(store, mustJSON(map[string]any{"kind": "WorkspaceTaskCloseoutStore@1", "schema_version": 1, "receipts": receipts}), 0600); err != nil {
				t.Fatal(err)
			}
			out, err = readInventory(f.D, root, read)
			if err != nil || reads != 5 {
				t.Fatalf("changed source: reads=%d err=%v", reads, err)
			}
			for _, row := range out.Runs {
				if row.Closeout != nil || row.Freshness == "fresh" || row.StatusFollowupVeto == "clear" {
					t.Error("foreign exact source was accepted", row)
				}
			}
		})
	}
}
