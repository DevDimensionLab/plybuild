package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// problemTaskFixture registers Task "task" with problem revision 1, whose
// only document "problem" holds "# Explicit fixture problem\n".
func problemTaskFixture(t *testing.T) (taskwork.Dependencies, string) {
	t.Helper()
	d, root := viewFixture(t)
	repo := filepath.Join(root, "repository")
	if err := os.Mkdir(repo, 0700); err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) string {
		t.Helper()
		c := exec.Command("git", append([]string{"-C", repo}, args...)...)
		raw, err := c.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, raw)
		}
		return strings.TrimSpace(string(raw))
	}
	git("init", "-b", "epic")
	git("-c", "user.email=fixture@example.invalid", "-c", "user.name=Fixture", "commit", "--allow-empty", "-m", "fixture")
	if _, err := workspace.AddProject(d.Dependencies, workspace.ProjectAddInput{ProjectID: "project", Name: "Project", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "repo", Path: repo}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := taskwork.AdoptEpic(d, taskwork.EpicAdoptInput{EpicID: "epic", Title: "Epic", ProjectID: "project", RepoID: "repo", Worktree: repo, Ref: "refs/heads/epic", ExpectedOID: git("rev-parse", "HEAD")}); err != nil {
		t.Fatal(err)
	}
	if _, err := taskwork.CreateTask(d, taskwork.TaskCreateInput{TaskID: "task", Title: "Task", Description: "Synthetic diagnostic fixture", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo"}); err != nil {
		t.Fatal(err)
	}

	document := filepath.Join(root, "problem.md")
	contents := []byte("# Explicit fixture problem\n")
	if err := os.WriteFile(document, contents, 0600); err != nil {
		t.Fatal(err)
	}
	fields := map[string]any{
		"kind": "WorkspaceTaskProblemDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785",
		"publication_key": "fixture/problem", "task_id": "task", "registry_upgrade": nil,
		"recorder":          map[string]any{"actor_claim": "fixture", "control_surface": "isolated test", "recorded_at_utc": "2026-10-30T00:00:00Z"},
		"expected_previous": nil, "origin": map[string]any{"kind": "authored"}, "title": "Fixture problem", "summary": "Test content diagnostics.", "problem_document_id": "problem",
		"documents": []any{map[string]any{"id": "problem", "source": map[string]any{"kind": "file", "locator": document, "sha256": fmt.Sprintf("sha256:%x", sha256.Sum256(contents)), "size_bytes": len(contents), "media_type": "text/markdown", "git_provenance": nil}}},
		"sources":   []any{}, "claims": []any{}, "deadline": nil, "change_reason": "Explicit fixture setup.",
	}
	problemRaw, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	problemPath := filepath.Join(root, "problem.json")
	if err := os.WriteFile(problemPath, problemRaw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := taskwork.RecordTaskProblem(d, taskwork.TaskContentInput{TaskID: "task", File: problemPath}); err != nil {
		t.Fatal(err)
	}
	return d, root
}

func TestTaskCommandSeparatesIncompatibleDocumentsFromTaskIdentity(t *testing.T) {
	d, root := problemTaskFixture(t)
	work, err := taskwork.LocateTask(d, root, "task")
	if err != nil {
		t.Fatal(err)
	}
	r, err := d.WorkItems.Snapshot(work)
	if err != nil {
		t.Fatal(err)
	}
	draft := taskwork.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "fixture/queue", ProjectID: "project", RepoID: "repo", EpicID: "epic", Entries: []taskwork.TaskGoalQueueEntry{}, Recorder: taskwork.TaskGoalQueueRecorder{ActorClaim: "fixture", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-08T00:00:00Z"}}
	raw, err := json.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "queue.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = taskwork.SetTaskQueue(d, path); err != nil {
		t.Fatal(err)
	}
	r, err = d.WorkItems.Snapshot(work)
	if err != nil {
		t.Fatal(err)
	}
	p := &r.TaskContentPublications[0]
	for _, pair := range []struct {
		kind   string
		digest *string
	}{{"requests", &p.RequestSHA256}, {"manifests", &p.OutcomeRef.ManifestSHA256}} {
		raw, err := d.TaskContent.Read(work, pair.kind, *pair.digest)
		if err != nil {
			t.Fatal(err)
		}
		v, err := canonicaljson.DecodeStrict(raw)
		if err != nil {
			t.Fatal(err)
		}
		v = append(v.(canonicaljson.Object), canonicaljson.Member{Name: "future_policy", Value: "SECRET_FIXTURE_VALUE_NEVER_DISPLAY"})
		raw, err = canonicaljson.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		*pair.digest, err = d.TaskContent.Publish(work, pair.kind, raw)
		if err != nil {
			t.Fatal(err)
		}
	}
	r.TaskProblemRevisions[0].ManifestSHA256 = p.OutcomeRef.ManifestSHA256
	// Store the changed references in the registry's own physical format.
	if err = d.WorkItems.WithLock(work, func(session taskwork.WorkItemStoreSession) error { return session.Publish(r) }); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(work.Registry())
	if err != nil {
		t.Fatal(err)
	}
	out, err := runWorkspaceView(t, d, "task", "show", "task", "--format", "json")
	if err != nil {
		t.Fatalf("document semantics blocked Task observation: %v", err)
	}
	var observed map[string]any
	if err := json.Unmarshal([]byte(out), &observed); err != nil {
		t.Fatal(err)
	}
	problem, ok := observed["problem"].(map[string]any)
	if !ok || problem["integrity"] != "incompatible" {
		t.Fatalf("document incompatibility was hidden: %s", out)
	}
	if !strings.Contains(out, "task_content_reader_incompatible") || strings.Contains(out, "SECRET_FIXTURE_VALUE_NEVER_DISPLAY") {
		t.Fatalf("diagnostic missing or content leaked: %s", out)
	}
	after, err := os.ReadFile(work.Registry())
	if err != nil || string(before) != string(after) {
		t.Fatal("diagnostic changed registry", err)
	}
}
