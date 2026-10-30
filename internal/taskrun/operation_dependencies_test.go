package taskrun

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func startupPublishSiblingProblem(t *testing.T, d Dependencies, root string) {
	t.Helper()
	r, err := d.Workspace.WorkItems.SnapshotRegistrations(root)
	if err != nil {
		t.Fatal(err)
	}
	var previous workspace.TaskRevisionRef
	for _, p := range r.TaskProblemRevisions {
		if p.TaskID == "sibling" {
			previous = workspace.TaskRevisionRef{Revision: p.Revision, ManifestSHA256: p.ManifestSHA256}
		}
	}
	doc := filepath.Join(root, "sibling.md")
	body := []byte("# Revised independent problem\n")
	if err = os.WriteFile(doc, body, 0600); err != nil {
		t.Fatal(err)
	}
	draft := map[string]any{"kind": "WorkspaceTaskProblemDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "sibling/revision", "task_id": "sibling", "registry_upgrade": nil, "expected_previous": previous, "origin": map[string]any{"kind": "authored"}, "title": "Revised sibling", "summary": "Independent revised need", "problem_document_id": "problem", "documents": []any{map[string]any{"id": "problem", "source": map[string]any{"kind": "file", "locator": doc, "sha256": hash(body), "size_bytes": len(body), "media_type": "text/markdown", "git_provenance": nil}}}, "sources": []any{}, "claims": []any{}, "deadline": nil, "change_reason": "Independent publication", "recorder": map[string]any{"actor_claim": "fixture", "control_surface": "synthetic test", "recorded_at_utc": "2026-10-05T00:00:00Z"}}
	b, err := json.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "sibling-problem.json")
	if err = os.WriteFile(file, b, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.RecordTaskProblem(d.Workspace, workspace.TaskContentInput{TaskID: "sibling", File: file}); err != nil {
		t.Fatal(err)
	}
}

func TestWorkflowStartupConfirmationUsesRelevantContentDependencies(t *testing.T) {
	f := newDeliveryFixture(t, "codex")
	if _, err := workspace.CreateTask(f.D.Workspace, workspace.TaskCreateInput{TaskID: "sibling", Title: "Sibling", Description: "Independent content", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	recordFixtureProblem(t, f.D.Workspace, f.R.WorkspaceRoot, "sibling")
	before, err := workflowPreview(f.D, f.R, f.File, WorkflowPaths{})
	if err != nil {
		t.Fatal(err)
	}
	startupPublishSiblingProblem(t, f.D, f.R.WorkspaceRoot)
	after, err := workflowPreview(f.D, f.R, f.File, WorkflowPaths{})
	if err != nil {
		t.Fatal(err)
	}
	if before.Observed.RegistrySHA256 == after.Observed.RegistrySHA256 {
		t.Fatal("physical registry readback did not change")
	}
	if before.Confirmation == nil || after.Confirmation == nil || *before.Confirmation != *after.Confirmation {
		t.Fatal("unrelated content invalidated startup confirmation")
	}
	if _, err = WorkflowStart(f.D, f.File, *before.Confirmation); err != nil {
		t.Fatal(err)
	}
}
