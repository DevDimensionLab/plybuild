package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"gopkg.in/yaml.v3"
)

func TestTaskCommandDiagnosesRequiredNewerPublicationWithoutContentValues(t *testing.T) {
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
	common := filepath.Join(repo, ".git")
	if _, _, _, err := d.Projects.Add(root, workspace.ProjectRecord{ID: "project", Name: "Project", Wrapper: root, RepoIDs: []workspace.RepoID{"repo"}}, []workspace.RepoRecord{{ID: "repo", Locator: repo, GitCommonDir: common}}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.AdoptEpic(d, workspace.EpicAdoptInput{EpicID: "epic", Title: "Epic", ProjectID: "project", RepoID: "repo", Worktree: repo, Ref: "refs/heads/epic", ExpectedOID: git("rev-parse", "HEAD")}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.CreateTask(d, workspace.TaskCreateInput{TaskID: "task", Title: "Task", Description: "Synthetic diagnostic fixture", ParentEpicID: "epic", ProjectID: "project", RepoID: "repo"}); err != nil {
		t.Fatal(err)
	}
	r, err := d.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	draft := workspace.WorkspaceTaskQueueDraft{Kind: "WorkspaceTaskQueueDraft@1", SchemaVersion: 1, PublicationKey: "fixture/queue", ProjectID: "project", RepoID: "repo", EpicID: "epic", Entries: []workspace.QueueEntry{}, HumanDecision: workspace.QueueHumanDecision{ActorClaim: "fixture", DecidedAtUTC: "2026-10-08T00:00:00Z", Source: "human_cli", Statement: "Synthetic setup only."}, RegistryUpgrade: &workspace.TaskRegistryUpgrade{FromVersion: 3, RegistrySHA256: r.RawSHA256}}
	raw, err := json.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "queue.json")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.SetTaskQueue(d, path); err != nil {
		t.Fatal(err)
	}
	r, err = d.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	p := &r.TaskContentPublications[0]
	for _, pair := range []struct {
		kind   string
		digest *string
	}{{"requests", &p.RequestSHA256}, {"manifests", &p.OutcomeRef.ManifestSHA256}} {
		raw, err := d.TaskContent.Read(root, pair.kind, *pair.digest)
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
		*pair.digest, err = d.TaskContent.Publish(root, pair.kind, raw)
		if err != nil {
			t.Fatal(err)
		}
	}
	r.TaskProblemRevisions[0].ManifestSHA256 = p.OutcomeRef.ManifestSHA256
	raw, err = yaml.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(root, ".ply", "work-items.yaml"), raw, 0600); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	out, err := runWorkspaceView(t, d, "task", "show", "task", "--format", "json")
	if err == nil {
		t.Fatal("required newer semantics were accepted")
	}
	var diagnostic map[string]any
	if decodeErr := json.Unmarshal([]byte(out), &diagnostic); decodeErr != nil {
		t.Fatalf("failure must have machine diagnostic: %v; %q; %v", decodeErr, out, err)
	}
	if diagnostic["code"] != "task_content_reader_incompatible" || diagnostic["artifact_schema"] != "WorkspaceTaskProblemRevision@1" {
		t.Fatalf("wrong compatibility diagnostic: %s", out)
	}
	for _, key := range []string{"operation", "artifact", "reader_capability", "cause", "next_action"} {
		if value, ok := diagnostic[key].(string); !ok || value == "" {
			t.Fatalf("missing %s: %s", key, out)
		}
	}
	if !strings.Contains(out, "future_policy") || strings.Contains(out+err.Error(), "SECRET_FIXTURE_VALUE_NEVER_DISPLAY") {
		t.Fatalf("diagnostic lost field cause or leaked content: %s %v", out, err)
	}
	after, err := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if err != nil || string(before) != string(after) {
		t.Fatal("diagnostic changed registry", err)
	}
}
