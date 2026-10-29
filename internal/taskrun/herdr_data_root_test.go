package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestWorkflowProjectDataRootFrozenAcrossConfigurationAndRounds(t *testing.T) {
	f := workflowTestFixture(t)
	data := filepath.Join(filepath.Dir(f.R.WorkspaceRoot), "project-data")
	selected, calls := data, 0
	f.D.Workspace.ProjectDataRoot = func(project string) (string, error) {
		calls++
		if project != "ply" {
			t.Fatalf("wrong project: %s", project)
		}
		return selected, nil
	}
	p, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	preview := p.(WorkflowPreview)
	want := filepath.Join(data, "workflow-runs", "runs", preview.RunID)
	if preview.Paths.RunRoot != want || calls != 1 {
		t.Fatalf("preview root=%s, resolver calls=%d", preview.Paths.RunRoot, calls)
	}
	if _, err = os.Stat(data); !os.IsNotExist(err) {
		t.Fatal("preview created data root")
	}
	selected = data + "-changed"
	if _, err = WorkflowStart(f.D, f.File, *preview.Confirmation); err == nil {
		t.Fatal("changed config reused old confirmation")
	}
	if _, err = os.Stat(workflowIndex(f.R.WorkspaceRoot, preview.RunID)); !os.IsNotExist(err) {
		t.Fatal("confirmation conflict reserved a run")
	}
	selected, calls = data, 0
	f.D.Fault = func(point string) error {
		if point == "workflow_before_reservation" {
			selected = data + "-changed"
		}
		return nil
	}
	o, err := WorkflowStart(f.D, f.File, *preview.Confirmation)
	if err != nil || o.Paths.RunRoot != want || calls != 1 {
		t.Fatalf("start did not freeze its one selection: %+v, calls=%d, %v", o.Paths, calls, err)
	}
	f.D.Workspace.ProjectDataRoot = func(string) (string, error) { return "", errors.New("mapping removed") }
	workflowAssertNoReplay(t, f, o)
	o = workflowTestSettle(t, f, workflowTestReport(t, f, workflowTestAccept(t, f, o)))
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "changes_requested"))
	if err != nil || o.Paths.RunRoot != want || o.Paths.Context != filepath.Join(want, "rounds", "001", "context.json") {
		t.Fatalf("correction moved run: %+v %v", o.Paths, err)
	}
	o = workflowTestSettle(t, f, workflowTestReport(t, f, o))
	o, err = WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if err != nil {
		t.Fatal(err)
	}
	rows, err := ReadInventory(f.D, f.R.WorkspaceRoot)
	if err != nil || len(rows.Runs) != 1 || rows.Runs[0].State != "completed" || rows.Runs[0].Freshness != "fresh" {
		t.Fatalf("inventory lost external completed run: %+v %v", rows, err)
	}
	traced := traceRun(t, traceRead(t, f.D, f.R.WorkspaceRoot), o.RunID)
	if len(traceKind(traced, "round_report")) != 2 || len(traceKind(traced, "round_review")) != 2 {
		t.Fatalf("trace lost external rounds: %+v", traced)
	}
	for _, path := range []string{selected, filepath.Join(workflowRoot(f.R.WorkspaceRoot), "runs")} {
		if _, err = os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("unexpected alternative run root: %s", path)
		}
	}
}

func TestWorkflowProjectDataRootRejectsUnsafePaths(t *testing.T) {
	f := workflowTestFixture(t)
	link := filepath.Join(f.R.WorkspaceRoot, "link")
	if err := os.Symlink(f.R.WorkspaceRoot, link); err != nil {
		t.Fatal(err)
	}
	for _, data := range []string{"relative", f.R.WorkspaceRoot + "/../unclean", filepath.Join(link, "missing"), f.File} {
		f.D.Workspace.ProjectDataRoot = func(string) (string, error) { return data, nil }
		if _, err := WorkflowPreviewStart(f.D, f.File); err == nil {
			t.Fatalf("unsafe data root accepted: %s", data)
		}
	}
}

func TestWorkflowProjectRootReadersRejectTampering(t *testing.T) {
	f := workflowTestFixture(t)
	f.D.Workspace.ProjectDataRoot = func(string) (string, error) { return filepath.Join(f.R.WorkspaceRoot, "data"), nil }
	o := workflowTestStart(t, f)
	index, state := workflowIndex(f.R.WorkspaceRoot, o.RunID), filepath.Join(o.Paths.RunRoot, "state.json")
	for _, target := range []string{"identity", "context", "root", "state", "private", "symlink"} {
		t.Run(target, func(t *testing.T) {
			path := index
			if target == "state" {
				path = state
			}
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.WriteFile(path, original, 0600); _ = os.Chmod(path, 0600) })
			var s workflowState
			if err = decode(original, 8<<20, &s); err != nil {
				t.Fatal(err)
			}
			switch target {
			case "identity":
				s.Result.RunID = "different"
			case "context":
				s.Result.Paths.Context = f.File
			case "root", "state":
				s.Result.Paths.RunRoot = f.R.WorkspaceRoot
			case "private":
				if err = os.Chmod(path, 0644); err != nil {
					t.Fatal(err)
				}
			case "symlink":
				if err = os.Rename(o.Paths.RunRoot, o.Paths.RunRoot+"-saved"); err != nil {
					t.Fatal(err)
				}
				if err = os.Symlink(o.Paths.RunRoot+"-saved", o.Paths.RunRoot); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = os.Remove(o.Paths.RunRoot); _ = os.Rename(o.Paths.RunRoot+"-saved", o.Paths.RunRoot) })
			}
			if target != "private" && target != "symlink" {
				if err = f.D.replaceValue(path, s); err != nil {
					t.Fatal(err)
				}
			}
			if _, err = workflowRead(f.R.WorkspaceRoot, o.RunID); err == nil {
				t.Fatal("reader accepted altered reservation/state")
			}
			row := readHerdrInventoryRun(f.D, f.R.WorkspaceRoot, o.RunID)
			if row.Freshness != "unknown" || !hasInventoryReason(row, "run_evidence_unavailable") {
				t.Fatalf("inventory accepted tampering: %+v", row)
			}
			reader := traceReader{captured: map[string]inventoryCapturedSource{}}
			reader.workflow(f.R.WorkspaceRoot, o.RunID)
			if !traceReason(reader.run, "trace_run_unbound") && !traceReason(reader.run, "trace_state_unavailable") {
				t.Fatalf("trace accepted tampering: %+v", reader.run)
			}
		})
	}
}

func testWorkflowPaths(r WorkflowRequest, round int) WorkflowPaths {
	return workflowPaths(filepath.Join(workflowRoot(r.WorkspaceRoot), "runs", workflowID(r)), round)
}
