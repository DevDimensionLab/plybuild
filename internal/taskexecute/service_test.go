package taskexecute

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func fixtureGit(t *testing.T, cwd string, args ...string) string {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = cwd
	b, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, b)
	}
	return strings.TrimSpace(string(b))
}

func fixtureJSON(t *testing.T, path string, value any) {
	t.Helper()
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, b, 0600); err != nil {
		t.Fatal(err)
	}
}

// This is an isolated real registry/Git journey. Only provider transport is
// replaced; its reservation fault occurs before any agent command is sent.
func launcherFixture(t *testing.T) (taskrun.Dependencies, Input, string, string) {
	return launcherFixtureForProvider(t, "claude")
}

func launcherFixtureForProvider(t *testing.T, provider string, delivery ...workspace.DeliveryAgreement) (taskrun.Dependencies, Input, string, string) {
	t.Helper()
	root := physicalTemp(t)
	repo, epic := filepath.Join(root, "main"), filepath.Join(root, "epic")
	if err := os.Mkdir(repo, 0700); err != nil {
		t.Fatal(err)
	}
	fixtureGit(t, repo, "init", "-b", "main")
	fixtureGit(t, repo, "config", "user.email", "test@example.invalid")
	fixtureGit(t, repo, "config", "user.name", "Ply execution fixture")
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("base\n"), 0644); err != nil {
		t.Fatal(err)
	}
	fixtureGit(t, repo, "add", "README.md")
	if provider == "codex" {
		if err := os.WriteFile(filepath.Join(repo, "AGENTS.md"), []byte("# Fixture instructions\nKeep this isolated test repository local.\n"), 0644); err != nil {
			t.Fatal(err)
		}
		fixtureGit(t, repo, "add", "AGENTS.md")
	}
	fixtureGit(t, repo, "commit", "-m", "base")
	oid := fixtureGit(t, repo, "rev-parse", "HEAD")
	fixtureGit(t, repo, "worktree", "add", "-b", "epic", epic, oid)
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	if err = os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	w := workspace.SystemDependencies()
	if _, err = workspace.Init(w); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.AddProject(w, workspace.ProjectAddInput{ProjectID: "fixture", Name: "Execution fixture", Wrapper: root, Repos: []workspace.RepoInput{{RepoID: "fixture", Path: repo}}}); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "epic", Title: "Fixture Epic", ProjectID: "fixture", RepoID: "fixture", Worktree: epic, Ref: "refs/heads/epic", ExpectedOID: oid}); err != nil {
		t.Fatal(err)
	}
	if _, err = workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "goal", Title: "Goal fixture", Description: "Test the execute composition.", ParentEpicID: "epic", ProjectID: "fixture", RepoID: "fixture"}); err != nil {
		t.Fatal(err)
	}
	r, err := w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	problem := workspace.TaskRevisionRef{Revision: r.TaskProblemRevisions[0].Revision, ManifestSHA256: r.TaskProblemRevisions[0].ManifestSHA256}
	doc := filepath.Join(root, "design.md")
	body := []byte("# Fixture goal\nAdd one locally verified behavior. No external effects.\n")
	if err = os.WriteFile(doc, body, 0600); err != nil {
		t.Fatal(err)
	}
	draft := map[string]any{
		"kind": "WorkspaceTaskSpecDraft@2", "schema_version": 2, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "fixture/goal", "task_id": "goal", "registry_upgrade": nil,
		"recorder": map[string]any{"actor_claim": "fixture planner", "control_surface": "synthetic test", "recorded_at_utc": "2026-10-05T00:00:00Z"},
		"spec_id":  "observable-goal", "expected_previous": nil, "problem": problem, "title": "Goal fixture", "contract_kind": "goal", "objective": "Observe one complete execution startup without a real provider.",
		"documents": []any{map[string]any{"id": "design", "source": map[string]any{"kind": "file", "locator": doc, "sha256": digestBytes(body), "size_bytes": len(body), "media_type": "text/markdown", "git_provenance": nil}}},
		"design":    []any{map[string]any{"document_id": "design", "section": nil}}, "requirements": []any{map[string]any{"id": "behavior", "acceptance": "Observe exact provider and worktree binding."}}, "constraints": []string{"No external effects."},
		"executor": map[string]any{"provider": provider, "model": "fixture-model", "effort": "medium"}, "change_reason": "Isolated automated fixture; no real human product pass.",
	}
	if len(delivery) != 0 {
		a := delivery[0]
		if a.SchemaVersion == 0 {
			a.SchemaVersion = 1
		}
		a.ProjectID, a.RepoID, a.EpicID = "fixture", "fixture", "epic"
		if a.TargetRef == "" {
			a.TargetRef = "refs/heads/epic"
		}
		if a.Mode != workspace.DeliveryPullRequest {
			a.TargetWorktree = epic
		}
		draft["delivery"] = a
	}
	file := filepath.Join(root, "goal.json")
	fixtureJSON(t, file, draft)
	s, err := workspace.RecordTaskSpec(w, workspace.TaskContentInput{TaskID: "goal", File: file})
	if err != nil {
		t.Fatal(err)
	}
	r, err = w.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	q := workspace.WorkspaceTaskGoalQueueDraft{Kind: "WorkspaceTaskQueueDraft@2", SchemaVersion: 2, PublicationKey: "fixture/queue", ProjectID: "fixture", RepoID: "fixture", EpicID: "epic", Entries: []workspace.TaskGoalQueueEntry{{TaskID: "goal", Goal: &workspace.TaskGoalRef{SpecID: "observable-goal", Spec: workspace.TaskRevisionRef{Revision: *s.OutcomeRef.Revision, ManifestSHA256: s.OutcomeRef.ManifestSHA256}}}}, Recorder: workspace.TaskGoalQueueRecorder{ActorClaim: "fixture planner", ControlSurface: "synthetic test", RecordedAtUTC: "2026-10-05T00:00:00Z"}}
	if r.FormatVersion < 4 {
		q.RegistryUpgrade = &workspace.TaskRegistryUpgrade{FromVersion: r.FormatVersion, RegistrySHA256: r.RawSHA256}
	}
	file = filepath.Join(root, "queue.json")
	fixtureJSON(t, file, q)
	if _, err = workspace.SetTaskQueue(w, file); err != nil {
		t.Fatal(err)
	}
	binDir := filepath.Join(root, "bin")
	if err = os.Mkdir(binDir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{provider, "herdr", "ply"} {
		if err = os.WriteFile(filepath.Join(binDir, name), []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	if provider == "codex" {
		configHome := filepath.Join(root, "codex-home")
		if err = os.Mkdir(configHome, 0700); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(configHome, "config.toml"), []byte("model = 'fixture-model'\nmodel_reasoning_effort = 'medium'\ndefault_permissions = 'fixture'\n[permissions.fixture]\nnetwork = false\n"), 0600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("CODEX_HOME", configHome)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("HERDR_ENV", "1")
	t.Setenv("HERDR_WORKSPACE_ID", "w-fixture")
	if err = os.Chdir(epic); err != nil {
		t.Fatal(err)
	}
	d := taskrun.SystemDependencies(w)
	d.Executable = func() (string, error) { return filepath.Join(binDir, "ply"), nil }
	return d, Input{Next: true, Runtime: RuntimeOptions{ProviderExecutable: filepath.Join(binDir, provider), HerdrExecutable: filepath.Join(binDir, "herdr"), HerdrWorkspace: "w-fixture"}}, root, epic
}

func TestExecuteCodexPreservesNativeTrustForRepositoryInstructions(t *testing.T) {
	d, in, root, _ := launcherFixtureForProvider(t, "codex")
	configPath := filepath.Join(root, "codex-home", "config.toml")
	configBefore, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	reservations := 0
	d.Fault = func(point string) error {
		if point == "workflow_after_reservation" {
			reservations++
			return errors.New("synthetic stop before native transport effects")
		}
		return nil
	}
	result, err := Execute(d, in)
	if reservations != 1 || err == nil || !strings.Contains(err.Error(), "synthetic stop before native transport effects") {
		t.Fatalf("Codex repository instructions blocked the native trust path before reservation: reservations=%d, error=%v", reservations, err)
	}
	if result.Run == nil || result.Run.RunID == "" || result.Run.Transport.TabID != "" {
		t.Fatal("test did not stop at one reservation before native startup")
	}
	var request taskrun.WorkflowRequest
	if _, err = readPreservedJSON(result.RequestPath, &request); err != nil {
		t.Fatal(err)
	}
	if request.CodexProjectTrust != nil {
		t.Fatal("goal execution must leave repository trust to the native provider")
	}
	if request.ClaudeProjectTrust != nil || result.ClaudeProjectTrust != nil {
		t.Fatal("Codex execution received Claude configuration authority")
	}
	if grant := recordedClaudeTrust(t, filepath.Join(filepath.Dir(result.RequestPath), "intent.json")); grant != nil {
		t.Fatal("Codex intent recorded Claude configuration authority")
	}
	if request.Runtime.Provider != "codex" || request.Runtime.Model != "fixture-model" || request.Runtime.PermissionBinding.ProfileID != "fixture" || request.Delivery.ReasoningEffort != "medium" {
		t.Fatal("native trust changed the assigned runtime or permission profile")
	}
	if _, err = os.Stat(filepath.Join(result.Goal.WorktreePath, "AGENTS.md")); err != nil {
		t.Fatal("prepared Task did not preserve repository instructions")
	}
	configAfter, err := os.ReadFile(configPath)
	if err != nil || string(configBefore) != string(configAfter) {
		t.Fatal("execute changed persistent provider configuration")
	}
}

func TestExecuteChecksBeforeWorktreeOrLaunchAndReusesReservation(t *testing.T) {
	d, in, root, epic := launcherFixture(t)
	before := fixtureGit(t, epic, "worktree", "list", "--porcelain")
	missing := in
	missing.Runtime.ProviderExecutable = filepath.Join(root, "missing-claude")
	if _, err := Execute(d, missing); err == nil {
		t.Fatal("missing binary did not block")
	}
	if got := fixtureGit(t, epic, "worktree", "list", "--porcelain"); got != before {
		t.Fatal("missing binary created worktree")
	}
	in.Check = true
	preview, err := Execute(d, in)
	if err != nil {
		t.Fatal(err)
	}
	if preview.State != "ready" || preview.Runtime.Runtime.Provider != "claude" {
		t.Fatalf("bad preview: %+v", preview)
	}
	if _, err = os.Stat(filepath.Dir(preview.Goal.AcceptancePath)); !os.IsNotExist(err) {
		t.Fatal("check wrote run area")
	}
	if got := fixtureGit(t, epic, "worktree", "list", "--porcelain"); got != before {
		t.Fatal("check created worktree")
	}
	wrong := d
	wrong.CWD = func() (string, error) { return filepath.Join(root, "main"), nil }
	if _, err = Execute(wrong, in); err == nil {
		t.Fatal("wrong return worktree accepted")
	}
	in.Check = false
	calls := 0
	d.Fault = func(point string) error {
		if point == "workflow_after_reservation" {
			calls++
			return errors.New("fixture interrupted before native effects")
		}
		return nil
	}
	first, err := Execute(d, in)
	if err == nil || calls != 1 {
		t.Fatalf("did not reach one reserved startup: calls=%d, error=%v, result=%+v", calls, err, first)
	}
	if first.Run == nil || first.Run.RunID == "" || first.Run.Transport.TabID != "" {
		t.Fatalf("bad preserved startup: %+v", first)
	}
	var request taskrun.WorkflowRequest
	if _, err = readPreservedJSON(first.RequestPath, &request); err != nil {
		t.Fatal(err)
	}
	if request.Runtime.Provider != "claude" || request.Runtime.PermissionBinding.ProfileID != "auto" || request.Runtime.PermissionBinding.AuthorityKind != "launch_contract_pending_runtime_acceptance" || request.CodexProjectTrust != nil {
		t.Fatalf("new Claude request must ask for native auto mode without granting runtime tool authority: %+v", request.Runtime)
	}
	after := fixtureGit(t, epic, "worktree", "list", "--porcelain")
	if after == before {
		t.Fatal("no Task worktree created")
	}
	// Replacement of the public installation cannot invalidate readback of this
	// reserved run or lead the retry to another provider/worktree.
	installed, _ := d.Executable()
	if err = os.WriteFile(installed, []byte("#!/bin/sh\necho replaced\n"), 0700); err != nil {
		t.Fatal(err)
	}
	second, err := Execute(d, in)
	if err != nil {
		t.Fatal(err)
	}
	if second.State != "existing" || second.Run.RunID != first.Run.RunID || calls != 1 {
		t.Fatalf("retry replayed: %+v calls=%d", second, calls)
	}
	if got := fixtureGit(t, epic, "worktree", "list", "--porcelain"); got != after {
		t.Fatal("retry created another worktree")
	}
}

func TestExecuteKeepsEarlierManualLaunchWhenAutoBecomesDefault(t *testing.T) {
	d, in, _, _ := launcherFixture(t)
	in.Check = true
	preview, err := Execute(d, in)
	if err != nil {
		t.Fatal(err)
	}
	// Construct a historical intent in this private fixture before publishing
	// any request. Real preserved launches are never rewritten by this test.
	intent := launchIntent{
		Kind: "ply.workflow.execute-intent", SchemaVersion: 1,
		Goal: *preview.Goal, Runtime: *preview.Runtime,
		Human: workspace.QueueHumanDecision{ActorClaim: "Fixture caller", DecidedAtUTC: "2026-10-05T00:00:00Z", Source: "explicit_human_instruction", Statement: "Fixture of a preserved earlier manual-mode launch."},
	}
	intent.Runtime.Policy.PermissionProfile = "manual"
	intent.Runtime.Runtime.PermissionBinding.ProfileID = "manual"
	directory := filepath.Dir(preview.RequestPath)
	intent.Runtime, err = preserveRuntime(intent.Runtime, directory)
	if err != nil {
		t.Fatal(err)
	}
	intentBytes, err := taskrun.Canonical(intent)
	if err != nil {
		t.Fatal(err)
	}
	intentPath := filepath.Join(directory, "intent.json")
	if err = writeOnce(intentPath, intentBytes, 0600); err != nil {
		t.Fatal(err)
	}
	before := map[string][]byte{}
	for _, path := range []string{intentPath, filepath.Join(directory, "launch-policy.json"), intent.Runtime.Runtime.PlyExecutable.Path} {
		before[path], err = os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
	}
	checked, err := Execute(d, in)
	if err != nil || checked.ClaudeProjectTrust != nil {
		t.Fatalf("checking an earlier intent introduced Claude trust: preview=%+v error=%v", checked.ClaudeProjectTrust, err)
	}
	calls := 0
	d.Fault = func(point string) error {
		if point == "workflow_after_reservation" {
			calls++
			return errors.New("fixture interrupted before native effects")
		}
		return nil
	}
	in.Check = false
	in.Runtime.PermissionProfile = "auto"
	first, err := Execute(d, in)
	if err == nil || !strings.Contains(err.Error(), "fixture interrupted before native effects") || calls != 1 || first.Run == nil || first.Run.Transport.TabID != "" {
		t.Fatalf("historical launch did not reach its one reservation: calls=%d error=%v", calls, err)
	}
	var request taskrun.WorkflowRequest
	before[first.RequestPath], err = readPreservedJSON(first.RequestPath, &request)
	if err != nil {
		t.Fatal(err)
	}
	if request.Runtime.PermissionBinding.ProfileID != "manual" {
		t.Fatal("new default changed the preserved launch request")
	}
	if request.ClaudeProjectTrust != nil || first.ClaudeProjectTrust != nil || strings.Contains(string(before[first.RequestPath]), `"claude_project_trust"`) {
		t.Fatal("earlier launch acquired new Claude configuration authority")
	}
	second, err := Execute(d, in)
	if err != nil || second.State != "existing" || second.Run.RunID != first.Run.RunID || calls != 1 {
		t.Fatalf("historical reservation was not reused: calls=%d error=%v result=%+v", calls, err, second)
	}
	for path, original := range before {
		if after, err := os.ReadFile(path); err != nil || string(after) != string(original) {
			t.Fatalf("preserved launch artifact changed: %s", path)
		}
	}
}

func TestExecuteCanRefreshUnstartedIntentAfterPlannerQueueUpdate(t *testing.T) {
	d, in, root, _ := launcherFixture(t)
	d.Fault = func(point string) error {
		if point == "execute_before_preparation" {
			return errors.New("fixture interrupted before any preparation")
		}
		return nil
	}
	first, err := Execute(d, in)
	if err == nil {
		t.Fatal("missing fixture interruption")
	}
	intentPath := filepath.Join(filepath.Dir(first.RequestPath), "intent.json")
	before, err := os.ReadFile(intentPath)
	if err != nil {
		t.Fatal(err)
	}
	queuePath := filepath.Join(root, "queue.json")
	var q workspace.WorkspaceTaskGoalQueueDraft
	raw, err := os.ReadFile(queuePath)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &q); err != nil {
		t.Fatal(err)
	}
	current, err := workspace.ListTaskQueue(d.Workspace, workspace.QueueTargetInput{}, false)
	if err != nil {
		t.Fatal(err)
	}
	q.PublicationKey = "fixture/planner-next"
	q.ExpectedRevision = current.Revision
	q.RegistryUpgrade = nil
	fixtureJSON(t, queuePath, q)
	if _, err = workspace.SetTaskQueue(d.Workspace, queuePath); err != nil {
		t.Fatal(err)
	}
	second, err := Execute(d, in)
	if err == nil || !strings.Contains(err.Error(), "fixture interrupted before any preparation") {
		t.Fatalf("planner metadata permanently blocked unstarted intent: %v", err)
	}
	if second.RequestPath == first.RequestPath {
		t.Fatal("fresh selection overwrote original intent namespace")
	}
	after, err := os.ReadFile(intentPath)
	if err != nil || string(after) != string(before) {
		t.Fatal("original intent was changed")
	}
}

func TestExecuteRejectsChangedPreservedIntent(t *testing.T) {
	for _, change := range []string{"goal", "provider", "model", "effort", "unknown", "duplicate"} {
		t.Run(change, func(t *testing.T) {
			d, in, _, epic := launcherFixture(t)
			d.Fault = func(point string) error {
				if point == "execute_before_preparation" {
					return errors.New("interrupted before preparation")
				}
				return nil
			}
			first, err := Execute(d, in)
			if err == nil {
				t.Fatal("missing interruption")
			}
			path := filepath.Join(filepath.Dir(first.RequestPath), "intent.json")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var intent launchIntent
			if err = json.Unmarshal(raw, &intent); err != nil {
				t.Fatal(err)
			}
			switch change {
			case "goal":
				intent.Goal.Goal.Objective = "different goal"
			case "provider":
				intent.Runtime.Runtime.Provider = "codex"
			case "model":
				intent.Runtime.Runtime.Model = "different-model"
			case "effort":
				intent.Runtime.ReasoningEffort = "low"
			}
			raw, err = json.Marshal(intent)
			if err != nil {
				t.Fatal(err)
			}
			if change == "unknown" {
				raw = append([]byte(`{"unexpected":true,`), raw[1:]...)
			}
			if change == "duplicate" {
				raw = append([]byte(`{"kind":"ply.workflow.execute-intent",`), raw[1:]...)
			}
			if err = os.WriteFile(path, raw, 0600); err != nil {
				t.Fatal(err)
			}
			before := fixtureGit(t, epic, "worktree", "list", "--porcelain")
			in.Check = true
			if _, err = Execute(d, in); err == nil {
				t.Fatal("changed intent accepted")
			}
			if after := fixtureGit(t, epic, "worktree", "list", "--porcelain"); after != before {
				t.Fatal("rejected intent created a worktree")
			}
		})
	}
}
