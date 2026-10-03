package taskrun

import (
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Opt-in native fixture export for the separately supplied CLI acceptance
// contract. The runner is an in-process fake; it never starts a provider.
func TestJournalNativeFixture(t *testing.T) {
	root := os.Getenv("PLY_JOURNAL_FIXTURE_ROOT")
	binary := os.Getenv("PLY_JOURNAL_TEST_BINARY")
	if root == "" || binary == "" {
		t.Skip("set private fixture root and actual CLI binary to export acceptance fixtures")
	}
	resume := os.Getenv("PLY_JOURNAL_RESUME_FIXTURE") == "1"
	var d Dependencies
	var first Request
	if resume {
		old, e := os.Getwd()
		if e != nil {
			t.Fatal(e)
		}
		t.Cleanup(func() { _ = os.Chdir(old) })
		if e = os.Chdir(root); e != nil {
			t.Fatal(e)
		}
		d = SystemDependencies(workspace.SystemDependencies())
		first, e = ReadRequest(filepath.Join(root, "request-1.json"))
		if e != nil {
			t.Fatal(e)
		}
	} else {
		if _, e := os.Lstat(root); !os.IsNotExist(e) {
			t.Fatal("fixture destination must be absent")
		}
		d, first, _ = journalFixture(t, root)
	}
	runs := []any{}
	var result workspace.TaskResultRecord
	for n := 1; n <= 2; n++ {
		r := first
		r.RequestKey = fmt.Sprintf("journal/run-%d", n)
		var draft map[string]any
		if e := json.Unmarshal(r.HandoffDraft, &draft); e != nil {
			t.Fatal(e)
		}
		draft["publication_key"] = r.RequestKey
		draft["activity_key"] = r.RequestKey
		r.HandoffDraft = mustCanonical(t, draft)
		local := d
		local.Runner = &fakeRunner{}
		local.ContextPath = func() string { return filepath.Join(runPaths(r).TempRoot, "context.json") }
		if !resume {
			completeInside(t, &local, r, nil)
			file := writeAny(t, root, fmt.Sprintf("request-%d.json", n), r)
			p, e := PreviewStart(local, file)
			if e != nil {
				t.Fatal(e)
			}
			preview := p.(Preview)
			if preview.Confirmation == nil {
				t.Fatalf("run preview: %+v", preview)
			}
			out, e := startAndObserve(t, local, r, file, *preview.Confirmation)
			if e != nil {
				t.Fatalf("native fixture run: %v; fake runner: %v", e, local.Runner.(*fakeRunner).lastErr)
			}
			if out.Collection.State != "qualified" {
				t.Fatalf("native return not qualified: %+v", out)
			}
		}
		j, e := ReadProcessJournal(local, root, RunID(r.RequestKey))
		if e != nil || len(j.Events) == 0 {
			t.Fatalf("native read validation: %v", e)
		}
		runs = append(runs, map[string]any{"run_id": RunID(r.RequestKey), "request_sha256": digest(r), "preparation_id": r.PreparationID, "preparation_sha256": r.PreparationSHA256})
	}
	registry, e := d.Workspace.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	if len(registry.TaskResults) != 2 {
		t.Fatal("two native results required")
	}
	result = registry.TaskResults[0]
	evidence := filepath.Join(root, "synthetic-human-qa.txt")
	body := []byte("Synthetic QA fixture only. No human product judgment was made.\n")
	if e = os.WriteFile(evidence, body, 0600); e != nil {
		t.Fatal(e)
	}
	qa := map[string]any{"kind": "WorkspaceTaskHumanQARecordDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "journal/synthetic-qa", "task_id": "journal-acceptance", "task_result_id": result.ID, "result_oid": result.ResultOID, "result_tree": result.ResultTree, "outcome": "pass", "actor": map[string]any{"actor_claim": "synthetic human, not actual approval", "start_surface": "isolated fixture", "started_at_utc": "2026-10-03T08:00:00Z", "completed_at_utc": "2026-10-03T08:01:00Z"}, "evidence": []any{map[string]any{"id": "report", "role": "report", "locator": evidence, "sha256": hash(body), "size_bytes": len(body)}}, "observation": "Synthetic pass for old candidate only.", "accepted_residual_risks": []any{}}
	if _, e = workspace.RecordTaskHumanQA(d.Workspace, workspace.TaskHumanQARecordInput{TaskID: "journal-acceptance", File: writeAny(t, root, "synthetic-qa.json", qa)}); e != nil {
		t.Fatal(e)
	}
	registry, e = d.Workspace.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	task := registry.Tasks[0]
	target := task.Worktree.Locator
	candidate := func(oid, tree string) any {
		return map[string]any{"repo_id": task.RepoID, "worktree_id": task.Worktree.ID, "oid": oid, "tree": tree}
	}
	a := candidate(result.ResultOID, result.ResultTree)
	if e = os.WriteFile(filepath.Join(target, "later-candidate.txt"), []byte("Candidate B, with no native technical or QA judgment.\n"), 0600); e != nil {
		t.Fatal(e)
	}
	runGit(t, target, "add", "later-candidate.txt")
	runGit(t, target, "-c", "user.name=Synthetic Journal Fixture", "-c", "user.email=journal@example.invalid", "commit", "-m", "synthetic candidate B")
	b := candidate(gitOutput(t, target, "rev-parse", "HEAD"), gitOutput(t, target, "rev-parse", "HEAD^{tree}"))
	command := exec.Command(binary, "workspace", "task", "journal", "show", "journal-acceptance", "--format", "json")
	command.Dir = root
	output, e := command.Output()
	if e != nil {
		t.Fatalf("actual CLI read: %v", e)
	}
	var snapshot struct {
		Events []struct {
			EventID string `json:"event_id"`
			Data    struct {
				Axis string `json:"axis"`
			} `json:"data"`
		} `json:"events"`
	}
	if e = json.Unmarshal(output, &snapshot); e != nil {
		t.Fatal(e)
	}
	var techID, qaID string
	for _, event := range snapshot.Events {
		if event.Data.Axis == "technical" {
			techID = event.EventID
		}
		if event.Data.Axis == "human" {
			qaID = event.EventID
		}
	}
	if techID == "" || qaID == "" {
		t.Fatal("native judgment events missing")
	}
	protected := []string{filepath.Join(root, "target"), target}
	entries, e := os.ReadDir(filepath.Join(root, ".ply"))
	if e != nil {
		t.Fatal(e)
	}
	for _, entry := range entries {
		if entry.Name() != "task-process" {
			protected = append(protected, filepath.Join(root, ".ply", entry.Name()))
		}
	}
	manifest := map[string]any{"workspace": root, "task_id": "journal-acceptance", "runs": runs, "candidate_a": a, "candidate_b": b, "native_qa_event_id": qaID, "native_technical_event_id": techID, "protected_paths": protected, "readonly_roots": []string{root, target}, "native_validation": "public workspace validators, native Start/Accept/SubmitReport/Collect using an in-process fake runner; ReadProcessJournal validated both chains"}
	writeAny(t, filepath.Dir(root), "manifest.json", manifest)
	if e = os.WriteFile(filepath.Join(filepath.Dir(root), "native-initial.json"), output, 0600); e != nil {
		t.Fatal(e)
	}
}

func mustCanonical(t *testing.T, v any) []byte {
	t.Helper()
	b, e := Canonical(v)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
