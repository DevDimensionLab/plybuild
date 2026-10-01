package taskrun

import (
	"encoding/json"
	"errors"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"testing"
)

func TestR8FrozenDraftRecoversBeforeAndAfterTaskResultWrite(t *testing.T) {
	for _, point := range []string{"before_task_result", "after_task_result"} {
		t.Run(point, func(t *testing.T) {
			d, r, f := fixture(t)
			completeInside(t, &d, r, nil)
			d.Fault = func(p string) error {
				if p == point {
					return errors.New("synthetic crash at " + p)
				}
				return nil
			}
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			_, e = startAndObserve(t, d, r, f, *p.(Preview).Confirmation)
			if e == nil {
				t.Fatal("fault did not fire")
			}
			d.Fault = nil
			draftPath := filepath.Join(runPaths(r).TempRoot, "task-result-draft.json")
			before, e := os.ReadFile(draftPath)
			if e != nil {
				t.Fatal(e)
			}
			registry, e := d.Workspace.WorkItems.Snapshot(r.WorkspaceRoot)
			if e != nil {
				t.Fatal(e)
			}
			if point == "after_task_result" {
				if len(registry.TaskResults) != 1 {
					t.Fatal("missing published record")
				}
				target := dMustPreparation(t, d, r).Plan.WorktreePath
				os.WriteFile(filepath.Join(target, "later.txt"), []byte("later work"), 0644)
			}
			state := treeState(t, r.WorkspaceRoot)
			preview, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
			if e != nil {
				t.Fatal(e)
			}
			if state != treeState(t, r.WorkspaceRoot) {
				t.Fatal("collect check wrote files")
			}
			out, e := Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation)
			if e != nil {
				t.Fatal(e)
			}
			if out.Collection.State != "qualified" {
				t.Fatalf("recovery failed: %+v", out)
			}
			after, _ := os.ReadFile(draftPath)
			if string(before) != string(after) {
				t.Fatal("recovery changed frozen draft")
			}
			registry, e = d.Workspace.WorkItems.Snapshot(r.WorkspaceRoot)
			if e != nil || len(registry.TaskResults) != 1 {
				t.Fatalf("duplicate result: %v", e)
			}
		})
	}
}
func TestR8PartialSemanticReportRecoversWithoutNewCandidate(t *testing.T) {
	d, r, f := fixture(t)
	completeInside(t, &d, r, nil)
	once := false
	d.Fault = func(p string) error {
		if p == "after_submit_terminal" && !once {
			once = true
			return errors.New("synthetic lost report acknowledgement")
		}
		return nil
	}
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	_, _ = Start(d, f, *p.(Preview).Confirmation)
	d.Fault = nil
	// Parent observed an unknown helper outcome after the callback error. Recovery
	// may link the terminal, but cannot invent a successful process termination.
	before := treeState(t, r.WorkspaceRoot)
	preview, e := PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	if before != treeState(t, r.WorkspaceRoot) {
		t.Fatal("check wrote")
	}
	out, e := Collect(d, r.WorkspaceRoot, RunID(r.RequestKey), *preview.Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	if out.Delivery.State != "received" || out.Collection.State == "qualified" {
		t.Fatalf("partial return: %+v", out)
	}
}
func TestR9CachePrefixAndCorruption(t *testing.T) {
	d, r, f := fixture(t)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	_, _ = Start(d, f, *p.(Preview).Confirmation)
	j, e := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	prefix := journal{Request: r, Binding: j.Binding, Result: initial(r)}
	for i := 0; i < 2; i++ {
		if e = fold(&prefix, j.Events[i]); e != nil {
			t.Fatal(e)
		}
		prefix.Result.LastEventSHA256 = &j.Hashes[i]
	}
	path := filepath.Join(runPaths(r).RunRoot, "result.json")
	if e = replaceValue(path, prefix.Result); e != nil {
		t.Fatal(e)
	}
	before := treeState(t, r.WorkspaceRoot)
	if _, e = Show(d, r.WorkspaceRoot, RunID(r.RequestKey)); e != nil {
		t.Fatal(e)
	}
	if before != treeState(t, r.WorkspaceRoot) {
		t.Fatal("show repaired cache")
	}
	prefix.Result.Launch.State = "failed"
	replaceValue(path, prefix.Result)
	if _, e = Show(d, r.WorkspaceRoot, RunID(r.RequestKey)); e == nil {
		t.Fatal("contradictory cache accepted")
	}
}

func TestR12LaterSelectionPreservesHistoricalReturn(t *testing.T) {
	d, r, f := fixture(t)
	completeInside(t, &d, r, nil)
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	out, e := startAndObserve(t, d, r, f, *p.(Preview).Confirmation)
	if e != nil {
		t.Fatal(e)
	}
	plan := dMustPreparation(t, d, r).Plan
	in := taskSpecFixtureDraft(t, r.WorkspaceRoot, "WorkspaceTaskSolutionSelectionDraft@1", "fixture/withdraw", map[string]any{"expected_previous_selection": plan.Selection, "action": "withdraw", "solution": nil, "reason": "Synthetic later selection change", "human_decision": map[string]any{"actor_claim": "synthetic human", "decided_at_utc": "2026-10-01T01:00:00Z", "source": "explicit_human_instruction", "statement": "Synthetic withdrawal, not actual human QA."}})
	if _, e = workspace.WithdrawTaskSolution(d.Workspace, in); e != nil {
		t.Fatal(e)
	}
	later, e := Show(d, r.WorkspaceRoot, RunID(r.RequestKey))
	if e != nil {
		t.Fatal(e)
	}
	if later.Collection.TaskResultID == nil || *later.Collection.TaskResultID != *out.Collection.TaskResultID {
		t.Fatal("historical result changed")
	}
	if _, e = Start(d, f, "historical"); e != nil {
		t.Fatal(e)
	}
	if d.Runner.(*fakeRunner).calls != 1 {
		t.Fatal("restarted after selection changed")
	}
}
func TestR6KnownDebtGateRetainsExistingSemantics(t *testing.T) {
	d, r, f := fixture(t)
	completeInside(t, &d, r, func(p *Report) {
		var artifacts []map[string]any
		json.Unmarshal(p.Artifacts, &artifacts)
		staging := filepath.Dir(artifacts[0]["locator"].(string))
		control := writeAny(t, staging, "control.json", map[string]any{"synthetic": "reviewed rare coverage gap"})
		info, _ := os.Stat(control)
		artifacts = append([]map[string]any{{"artifact_id": "debt-control", "kind": "managed", "description": "Synthetic debt control", "media_type": "application/json", "classification": "workspace_internal", "size_bytes": info.Size(), "sha256": hashFileTest(t, control), "locator": control}}, artifacts...)
		p.Artifacts, _ = Canonical(artifacts)
		p.EvidenceGaps = json.RawMessage(`[{"type":"rare_coverage","detail":"Synthetic rare branch coverage gap.","artifact_ids":["debt-control"],"effect_ids":[]}]`)
		p.TechnicalAssessment = TechnicalAssessment{Gate: "good_enough_with_known_debt", RequiredVerifierIDs: []string{"test"}, AcceptedDebt: []workspace.TaskAcceptedDebtRecord{{ID: "rare-coverage", RiskClass: "coverage", Severity: "low", Summary: "Synthetic gap", Control: "Synthetic bounded followup", EvidenceArtifactIDs: []string{"debt-control"}}}}
	})
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	out, e := startAndObserve(t, d, r, f, *p.(Preview).Confirmation)
	if e != nil || out.Collection.State != "qualified" {
		t.Fatalf("known debt rejected: %v %+v child=%v", e, out, d.Runner.(*fakeRunner).lastErr)
	}
}
