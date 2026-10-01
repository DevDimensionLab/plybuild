package taskrun

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestR1ReadOnlyDeterministicPreview(t *testing.T) {
	d, r, f := fixture(t)
	before := treeState(t, r.WorkspaceRoot)
	a, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	b, e := PreviewStart(d, f)
	if e != nil || !equal(a, b) {
		t.Fatalf("non-deterministic preview: %v", e)
	}
	if before != treeState(t, r.WorkspaceRoot) {
		t.Fatal("preview wrote files")
	}
	bad := r
	bad.PreparationSHA256 = hash([]byte("stale"))
	badFile := writeAny(t, r.WorkspaceRoot, "stale.json", bad)
	if _, e = PreviewStart(d, badFile); e == nil {
		t.Fatal("stale preparation accepted")
	}
	if _, e = os.Stat(storeRoot(r.WorkspaceRoot)); !os.IsNotExist(e) {
		t.Fatal("preflight wrote store")
	}
}
func TestR6CompleteSyntheticReturnAndHistoricalRetry(t *testing.T) {
	d, r, f := fixture(t)
	now := d.Now()
	d.Now = func() time.Time { return now }
	runner := d.Runner.(*fakeRunner)
	runner.inside = func(s LaunchSpec) error {
		if s.CWD != dMustPreparation(t, d, r).Plan.WorktreePath || len(s.Argv) != 7 || s.Argv[1] != "--cd" || s.Argv[5] != "--no-alt-screen" {
			return fmt.Errorf("incorrect argv: %v", s.Argv)
		}
		claim := acceptance(t, d, r)
		cp := writeAny(t, r.WorkspaceRoot, "accept.json", claim)
		a, e := Accept(d, r.WorkspaceRoot, RunID(r.RequestKey), cp)
		if e != nil {
			return e
		}
		if a.Acceptance.State != "started" {
			return fmt.Errorf("not accepted: %+v", a)
		}
		report := semanticReport(t, d, r)
		rp := writeAny(t, r.WorkspaceRoot, "report.json", report)
		out, e := SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), rp)
		if e != nil {
			return e
		}
		if out.Delivery.State != "received" || out.Collection.State == "qualified" {
			j, _ := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
			_, err := workflowhandoff.SubmitResultDocument(d.Workflow, workflowhandoff.SubmitInput{HandoffLocator: j.Binding.Handoff.Locator, DraftPath: terminalDraftPath(r, digest(report))})
			return fmt.Errorf("report before exit: %v; %+v", err, out)
		}
		// Retrying is storage-idempotent, while readback measures elapsed time
		// afresh until process exit. Exercise that boundary without wall-clock luck.
		if out.Budget.ElapsedSeconds == nil {
			return fmt.Errorf("report readback lacks elapsed time")
		}
		want := out
		want.Budget.ElapsedSeconds = ptr(*out.Budget.ElapsedSeconds + 1)
		beforeRetry := treeState(t, r.WorkspaceRoot)
		now = now.Add(time.Second)
		again, e := SubmitReport(d, r.WorkspaceRoot, RunID(r.RequestKey), rp)
		if e != nil || !equal(want, again) {
			return fmt.Errorf("report retry: %v", e)
		}
		if treeState(t, r.WorkspaceRoot) != beforeRetry {
			return fmt.Errorf("identical report retry changed preserved files")
		}
		return nil
	}
	p, e := PreviewStart(d, f)
	if e != nil {
		t.Fatal(e)
	}
	out, e := startAndObserve(t, d, r, f, *p.(Preview).Confirmation)
	if e != nil {
		t.Fatalf("%v; child=%v\n%+v", e, runner.lastErr, out)
	}
	if out.Collection.State != "qualified" {
		t.Fatalf("not qualified: %+v", out)
	}
	registry, e := d.Workspace.WorkItems.Snapshot(r.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	if len(registry.TaskResults) != 1 || len(registry.HumanQARecords) != 0 {
		t.Fatalf("unexpected registry results")
	}
	if runner.calls != 1 {
		t.Fatal("spawn count")
	}
	target := out.ObservedTarget.WorktreeLocator
	os.WriteFile(filepath.Join(target, "later.txt"), []byte("later independent work"), 0644)
	retry, e := Start(d, f, "historical retry needs no new authority")
	if e != nil || retry.Collection.TaskResultID == nil || *retry.Collection.TaskResultID != *out.Collection.TaskResultID || runner.calls != 1 {
		t.Fatalf("retry relaunched or lost result: %v %+v", e, retry)
	}
	before := treeState(t, r.WorkspaceRoot)
	if _, e = Show(d, r.WorkspaceRoot, RunID(r.RequestKey)); e != nil {
		t.Fatal(e)
	}
	if _, e = PreviewCollect(d, r.WorkspaceRoot, RunID(r.RequestKey)); e != nil {
		t.Fatal(e)
	}
	if before != treeState(t, r.WorkspaceRoot) {
		t.Fatal("readback mutated history")
	}
}
func dMustPreparation(t *testing.T, d Dependencies, r Request) workspace.TaskPreparation {
	t.Helper()
	p, e := workspace.ShowTaskPreparation(d.Workspace, r.PreparationID)
	if e != nil {
		t.Fatal(e)
	}
	return *p.Preparation
}
func treeState(t *testing.T, root string) string {
	t.Helper()
	m := map[string]string{}
	e := filepath.WalkDir(root, func(p string, x os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if x.Type().IsRegular() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			m[p] = hash(b)
		}
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	b, _ := json.Marshal(m)
	return string(b)
}
