package taskexecute

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func boundFixture(t *testing.T) (taskrun.Dependencies, BoundStartInput, string, string) {
	t.Helper()
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic bound-start fixture", TimeoutSeconds: 30}}
	d, _, root, epic := launcherFixtureForProvider(t, "codex", a)
	home := filepath.Join(root, "fixture-home")
	if err := os.Mkdir(home, 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("HERDR_ENV", "")
	standin := strings.Replace(recoveryLauncherHerdr, "args = sys.argv[1:]", "args = sys.argv[1:]\nmodel.setdefault('calls', []).append(args)", 1)
	if err := os.WriteFile(filepath.Join(root, "bin", "herdr"), []byte(standin), 0700); err != nil {
		t.Fatal(err)
	}
	d.HerdrTimeout = 2 * time.Second
	return d, BoundStartInput{SpecID: "observable-goal", StartMode: "interactive", HerdrWorkspace: "w-fixture"}, root, epic
}

func boundPreview(t *testing.T, d taskrun.Dependencies, in BoundStartInput) BoundStartPreview {
	t.Helper()
	p, err := PreviewBoundStart(d, "fixture-request-0001", in)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestBoundStartPreviewAndStatusAreReadOnly(t *testing.T) {
	d, in, root, epic := boundFixture(t)
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	refs := fixtureGit(t, epic, "show-ref")
	trees := fixtureGit(t, epic, "worktree", "list", "--porcelain")
	p := boundPreview(t, d, in)
	if p.State != "ready" || p.Request.Binding == "" || p.Basis.Runtime.Runtime.Model != "fixture-model" || p.Basis.Goal.Delivery == nil || p.Basis.RuntimeAcceptance == "" || len(p.Basis.PolicySources) == 0 || p.Basis.PermissionDefinition == nil {
		t.Fatalf("incomplete preview: %+v", p)
	}
	s, err := StatusBoundStart(d, root, p.Request.RequestID)
	if err != nil || s.State != "not_recorded" || s.ExecutionEffect != "not_confirmed" {
		t.Fatalf("status: %+v %v", s, err)
	}
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) || refs != fixtureGit(t, epic, "show-ref") || trees != fixtureGit(t, epic, "worktree", "list", "--porcelain") {
		t.Fatal("read changed native files, refs or worktrees")
	}
	if _, err := os.Lstat(filepath.Join(root, "bin", "recovery-model.json")); !os.IsNotExist(err) {
		t.Fatal("read invoked Herdr")
	}
}

func TestBoundStartRejectsChangedBasisBeforeAnyIntent(t *testing.T) {
	for _, change := range []string{"binding", "missing-binding", "schema", "mode", "model", "permission", "trust", "provider", "herdr-binary", "herdr-target", "environment", "queue", "goal", "base", "agreement"} {
		t.Run(change, func(t *testing.T) {
			d, in, root, epic := boundFixture(t)
			p := boundPreview(t, d, in)
			r := p.Request
			switch change {
			case "binding":
				r.Binding = "sha256:" + strings.Repeat("0", 64)
			case "missing-binding":
				r.Binding = ""
			case "schema":
				r.SchemaVersion = 2
			case "mode":
				r.Input.StartMode = "headless"
			case "model", "permission", "trust":
				path := filepath.Join(root, "codex-home", "config.toml")
				raw, _ := os.ReadFile(path)
				if change == "model" {
					raw = []byte(strings.Replace(string(raw), "fixture-model", "changed-model", 1))
				}
				if change == "permission" {
					raw = []byte(strings.Replace(string(raw), "network = false", "network = true", 1))
				}
				if change == "trust" {
					raw = append(raw, []byte("\n[projects.'"+epic+"']\ntrust_level='trusted'\n")...)
				}
				if err := os.WriteFile(path, raw, 0600); err != nil {
					t.Fatal(err)
				}
			case "control", "provider", "herdr-binary":
				name := map[string]string{"control": "ply", "provider": "codex", "herdr-binary": "herdr"}[change]
				if err := os.WriteFile(filepath.Join(root, "bin", name), []byte("#!/bin/sh\n# changed executable\nexit 0\n"), 0700); err != nil {
					t.Fatal(err)
				}
			case "herdr-target":
				r.Input.HerdrWorkspace = "different-workspace"
			case "environment":
				t.Setenv("PATH", os.Getenv("PATH")+string(os.PathListSeparator)+root)
			case "queue":
				boundChangeQueue(t, d, root)
			case "goal", "agreement":
				boundChangeGoal(t, d, root, change)
			case "base":
				if err := os.WriteFile(filepath.Join(epic, "another.txt"), []byte("different base"), 0600); err != nil {
					t.Fatal(err)
				}
				fixtureGit(t, epic, "add", "another.txt")
				fixtureGit(t, epic, "commit", "-m", "changed base")
			}
			before := recoveryFiles(t, filepath.Join(root, ".ply"))
			trees := fixtureGit(t, epic, "worktree", "list", "--porcelain")
			out, err := ApplyBoundStart(d, r)
			if err == nil || out.State != "rejected" || out.ExecutionEffect != "none" || out.Error == nil {
				t.Fatalf("change accepted: %+v %v", out, err)
			}
			if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) || trees != fixtureGit(t, epic, "worktree", "list", "--porcelain") {
				t.Fatal("stale/invalid approval wrote native state")
			}
		})
	}
}

func boundChangeQueue(t *testing.T, d taskrun.Dependencies, root string) {
	t.Helper()
	path := filepath.Join(root, "queue.json")
	var q workspace.WorkspaceTaskGoalQueueDraft
	raw, _ := os.ReadFile(path)
	if err := json.Unmarshal(raw, &q); err != nil {
		t.Fatal(err)
	}
	current, err := workspace.ListTaskQueue(d.Workspace, workspace.QueueTargetInput{}, false)
	if err != nil {
		t.Fatal(err)
	}
	q.PublicationKey = "fixture/bound-changed-queue"
	q.ExpectedRevision = current.Revision
	q.RegistryUpgrade = nil
	fixtureJSON(t, path, q)
	if _, err = workspace.SetTaskQueue(d.Workspace, path); err != nil {
		t.Fatal(err)
	}
}

func boundChangeGoal(t *testing.T, d taskrun.Dependencies, root, change string) {
	t.Helper()
	path := filepath.Join(root, "goal.json")
	raw, _ := os.ReadFile(path)
	var draft map[string]any
	if err := json.Unmarshal(raw, &draft); err != nil {
		t.Fatal(err)
	}
	p := boundPreview(t, d, BoundStartInput{SpecID: "observable-goal", StartMode: "interactive", HerdrWorkspace: "w-fixture"})
	draft["publication_key"] = "fixture/changed-goal"
	draft["expected_previous"] = p.Basis.Goal.Goal.Goal.Spec
	if change == "goal" {
		draft["objective"] = "A changed version of the native Goal"
	} else {
		draft["delivery"].(map[string]any)["acceptance"].(map[string]any)["responsible_actor"] = "different acceptance owner"
	}
	fixtureJSON(t, path, draft)
	if _, err := workspace.RecordTaskSpec(d.Workspace, workspace.TaskContentInput{TaskID: "goal", File: path}); err != nil {
		t.Fatal(err)
	}
}

func TestBoundStartValidConcurrentReplayAndLostResponseHaveOneWriter(t *testing.T) {
	d, in, root, epic := boundFixture(t)
	p := boundPreview(t, d, in)
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan BoundStartStatus, 8)
	errs := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			out, err := ApplyBoundStart(d, p.Request)
			results <- out
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent valid apply: %v", err)
		}
	}
	for out := range results {
		if out.ExpectedRunID != boundRunID(p.Request) {
			t.Fatalf("different identity: %+v", out)
		}
	}
	var model struct {
		Tabs, Starts, Prompts int
		Calls                 [][]string
	}
	raw, err := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
	if err != nil || json.Unmarshal(raw, &model) != nil || model.Tabs != 1 || model.Starts != 1 || model.Prompts != 1 {
		t.Fatalf("native transport effects: %s %v", raw, err)
	}
	for _, call := range model.Calls {
		if len(call) > 1 && call[0] == "tab" && call[1] == "create" {
			joined := strings.Join(call, "\n")
			for _, required := range []string{"HOME=" + p.Basis.Environment.Home, "CODEX_HOME=" + p.Basis.Environment.CodexHome, "PATH=" + p.Basis.Environment.Path} {
				if !strings.Contains(joined, required) {
					t.Fatalf("tab did not receive approved environment: missing %s", required)
				}
			}
		}
	}
	if trees := fixtureGit(t, epic, "worktree", "list", "--porcelain"); strings.Count(trees, "worktree ") != 3 {
		t.Fatalf("duplicate Task worktree: %s", trees)
	}
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	// Own queue movement and later replaced runtime cannot turn replay into a
	// new start, nor make the reader require the original installed control.
	if err := os.Remove(filepath.Join(root, "bin", "codex")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "ply"), []byte("replacement"), 0700); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		out, err := ApplyBoundStart(d, p.Request)
		if err != nil || out.Run == nil || out.State != "existing_run" {
			t.Fatalf("lost reply replay: %+v %v", out, err)
		}
	}
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("replay changed durable state")
	}
	if after, _ := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json")); string(after) != string(raw) {
		t.Fatal("replay/status called Herdr")
	}
	changed := p.Request
	changed.Input.PermissionProfile = "different"
	if _, err := ApplyBoundStart(d, changed); err == nil || !strings.Contains(err.Error(), "bound_start_request_conflict") {
		t.Fatal("same request ID accepted different bytes")
	}
}

func TestBoundStartInterruptedAttemptNeverResumesOnReplay(t *testing.T) {
	for _, point := range []string{"bound_start_intent_published", "bound_start_after_intent", "execute_before_preparation", "workflow_after_reservation", "workflow_before_agent_send"} {
		t.Run(point, func(t *testing.T) {
			d, in, root, _ := boundFixture(t)
			p := boundPreview(t, d, in)
			hits := 0
			d.Fault = func(at string) error {
				if at == point {
					hits++
					return errors.New("synthetic lost response")
				}
				return nil
			}
			first, err := ApplyBoundStart(d, p.Request)
			if err == nil || hits != 1 || first.ExecutionEffect == "none" {
				t.Fatalf("no preserved interruption: %+v %v", first, err)
			}
			before := recoveryFiles(t, filepath.Join(root, ".ply"))
			d.Fault = nil
			for i := 0; i < 2; i++ {
				out, err := ApplyBoundStart(d, p.Request)
				if err != nil || out.ExpectedRunID != first.ExpectedRunID {
					t.Fatalf("replay: %+v %v", out, err)
				}
			}
			if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
				t.Fatal("unknown attempt was resumed")
			}
			if _, err := StatusBoundStart(d, root, p.Request.RequestID); err != nil {
				t.Fatal(err)
			}
		})
	}
}

type boundAfterProjectLock struct {
	workspace.ProjectSnapshotLocker
	once  sync.Once
	after func()
}

func (l *boundAfterProjectLock) WithSnapshotLock(root string, fn func(workspace.ProjectSnapshot) error) error {
	err := l.ProjectSnapshotLocker.WithSnapshotLock(root, fn)
	if err == nil {
		l.once.Do(l.after)
	}
	return err
}

func TestBoundStartRejectsNativeWriteBetweenContentAndWorktree(t *testing.T) {
	d, in, root, epic := boundFixture(t)
	p := boundPreview(t, d, in)
	original := d.Workspace
	d.Workspace.ProjectLocks = &boundAfterProjectLock{ProjectSnapshotLocker: original.ProjectLocks, after: func() {
		if _, err := workspace.CreateTask(original, workspace.TaskCreateInput{TaskID: "concurrent", Title: "Concurrent native work", Description: "A fixture writer at the native publication boundary.", ParentEpicID: "epic", ProjectID: "fixture", RepoID: "fixture"}); err != nil {
			t.Fatal(err)
		}
	}}
	out, err := ApplyBoundStart(d, p.Request)
	if err == nil || !strings.Contains(err.Error(), "bound_start_stale") || out.State != "retained_attempt" {
		t.Fatalf("native publication race accepted: %+v %v", out, err)
	}
	if strings.Count(fixtureGit(t, epic, "worktree", "list", "--porcelain"), "worktree ") != 2 {
		t.Fatal("native race created a Task worktree")
	}
	if _, err := os.Lstat(filepath.Join(root, "bin", "recovery-model.json")); !os.IsNotExist(err) {
		t.Fatal("native race invoked Herdr")
	}
}

func TestBoundStartRequestRejectsDuplicateMissingNullAndUnknownFields(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	p := boundPreview(t, d, in)
	canonical, _ := taskrun.Canonical(p.Request)
	for _, change := range []string{"valid", "duplicate", "missing", "null", "unknown"} {
		raw := string(canonical)
		switch change {
		case "duplicate":
			raw = `{"kind":"ply.workflow.bound-start-request",` + raw[1:]
		case "missing":
			raw = strings.Replace(raw, `"permission_profile":"",`, "", 1)
		case "null":
			raw = strings.Replace(raw, `"permission_profile":""`, `"permission_profile":null`, 1)
		case "unknown":
			raw = `{"unexpected":true,` + raw[1:]
		}
		path := filepath.Join(root, "request-"+change+".json")
		if err := os.WriteFile(path, []byte(raw), 0600); err != nil {
			t.Fatal(err)
		}
		_, err := ReadBoundStartRequest(path)
		if (err == nil) != (change == "valid") {
			t.Fatalf("%s request: %v", change, err)
		}
	}
}

func TestBoundStartNativeLockAndReservationRacesRejectChangedBasis(t *testing.T) {
	for _, point := range []string{"bound_start_before_lock", "execute_before_preparation", "workflow_before_reservation", "workflow_before_tab_send", "workflow_before_agent_send"} {
		t.Run(point, func(t *testing.T) {
			d, in, root, epic := boundFixture(t)
			p := boundPreview(t, d, in)
			hits := 0
			d.Fault = func(at string) error {
				if at == point {
					hits++
					if point == "bound_start_before_lock" || point == "execute_before_preparation" {
						boundChangeQueue(t, d, root)
					} else {
						path := filepath.Join(root, "codex-home", "config.toml")
						raw, _ := os.ReadFile(path)
						return os.WriteFile(path, append(raw, []byte("\n# changed at native boundary\n")...), 0600)
					}
				}
				return nil
			}
			out, err := ApplyBoundStart(d, p.Request)
			if err == nil || hits != 1 || !strings.Contains(err.Error(), "bound_start_stale") {
				t.Fatalf("race accepted: %+v %v", out, err)
			}
			var model struct{ Tabs, Starts int }
			raw, _ := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
			_ = json.Unmarshal(raw, &model)
			if model.Starts != 0 || (point != "workflow_before_agent_send" && model.Tabs != 0) {
				t.Fatalf("race started transport: %s", raw)
			}
			if point == "bound_start_before_lock" || point == "execute_before_preparation" {
				if strings.Count(fixtureGit(t, epic, "worktree", "list", "--porcelain"), "worktree ") != 2 {
					t.Fatal("race created Task worktree")
				}
			}
		})
	}
}

func TestBoundStartRejectsGitDriftAtTransportBoundary(t *testing.T) {
	for _, which := range []string{"epic", "task", "missing-task"} {
		t.Run(which, func(t *testing.T) {
			d, in, root, epic := boundFixture(t)
			p := boundPreview(t, d, in)
			d.Fault = func(point string) error {
				if which == "epic" && point == "workflow_before_tab_send" {
					if err := os.WriteFile(filepath.Join(epic, "drift.txt"), []byte("new base"), 0600); err != nil {
						return err
					}
					fixtureGit(t, epic, "add", "drift.txt")
					fixtureGit(t, epic, "commit", "-m", "drift after reservation")
				}
				if which == "task" && point == "workflow_before_agent_send" {
					return os.WriteFile(filepath.Join(p.Basis.Goal.WorktreePath, "drift.txt"), []byte("changed Task"), 0600)
				}
				if which == "missing-task" && point == "workflow_before_tab_send" {
					return os.Rename(p.Basis.Goal.WorktreePath, p.Basis.Goal.WorktreePath+"-displaced")
				}
				return nil
			}
			out, err := ApplyBoundStart(d, p.Request)
			if err == nil || !strings.Contains(err.Error(), "bound_start_stale") {
				t.Fatalf("transport accepted Git drift: %+v %v", out, err)
			}
			var model struct{ Tabs, Starts int }
			raw, _ := os.ReadFile(filepath.Join(root, "bin", "recovery-model.json"))
			_ = json.Unmarshal(raw, &model)
			if model.Starts != 0 || which != "task" && model.Tabs != 0 {
				t.Fatalf("Git drift started transport: %s", raw)
			}
		})
	}
}

func TestBoundStartDoesNotAdoptAnEarlierUnpreparedRuntime(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	config := filepath.Join(root, "codex-home", "config.toml")
	raw, _ := os.ReadFile(config)
	if err := os.WriteFile(config, append(raw, []byte("\n[permissions.wider]\nnetwork = true\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HERDR_ENV", "1")
	d.Fault = func(point string) error {
		if point == "execute_before_preparation" {
			return errors.New("earlier ordinary start interrupted")
		}
		return nil
	}
	if _, err := Execute(d, Input{SpecID: in.SpecID, Runtime: RuntimeOptions{HerdrWorkspace: in.HerdrWorkspace, PermissionProfile: "wider"}}); err == nil {
		t.Fatal("missing earlier interruption")
	}
	d.Fault = nil
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	if p, err := PreviewBoundStart(d, "another-approved-request", in); err == nil {
		t.Fatalf("new approval adopted prior wider runtime: %+v", p)
	}
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("rejected initial preview changed earlier intent")
	}
}

func TestBoundStartRejectsIgnoredTaskPolicyBeforeTransport(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	exclude := filepath.Join(root, "main", ".git", "info", "exclude")
	if err := os.WriteFile(exclude, []byte(".codex/\n"), 0600); err != nil {
		t.Fatal(err)
	}
	p := boundPreview(t, d, in)
	d.Fault = func(point string) error {
		if point == "workflow_before_tab_send" {
			if err := os.Mkdir(filepath.Join(p.Basis.Goal.WorktreePath, ".codex"), 0700); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(p.Basis.Goal.WorktreePath, ".codex", "config.toml"), []byte("[permissions.fixture]\nnetwork=true\n"), 0600)
		}
		return nil
	}
	out, err := ApplyBoundStart(d, p.Request)
	if err == nil || !strings.Contains(err.Error(), "bound_start_stale") {
		t.Fatalf("ignored Task policy accepted: state=%s error=%v", out.State, err)
	}
	if _, err := os.Lstat(filepath.Join(root, "bin", "recovery-model.json")); !os.IsNotExist(err) {
		t.Fatal("ignored Task policy started transport")
	}
}

func TestBoundStartSeparateRequestsCannotAliasAnotherApprovedRun(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	config := filepath.Join(root, "codex-home", "config.toml")
	raw, _ := os.ReadFile(config)
	if err := os.WriteFile(config, append(raw, []byte("\n[permissions.wider]\nnetwork=true\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	first := boundPreview(t, d, in)
	in.PermissionProfile = "wider"
	second, err := PreviewBoundStart(d, "different-approved-request", in)
	if err != nil {
		t.Fatal(err)
	}
	firstRetained := make(chan struct{})
	release := make(chan struct{})
	done := make(chan error, 1)
	d1 := d
	d1.Fault = func(point string) error {
		if point == "bound_start_after_intent" {
			close(firstRetained)
			<-release
		}
		return nil
	}
	go func() { _, err := ApplyBoundStart(d1, first.Request); done <- err }()
	<-firstRetained
	d2 := d
	d2.Fault = func(point string) error {
		if point == "bound_start_after_intent" {
			return errors.New("second request stopped before effects")
		}
		return nil
	}
	_, secondErr := ApplyBoundStart(d2, second.Request)
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if secondErr == nil {
		t.Fatal("missing second interruption")
	}
	status, err := StatusBoundStart(d, root, second.Request.RequestID)
	if err != nil || status.State != "retained_attempt" || status.Run != nil {
		t.Fatalf("stopped request adopted another policy's native Run: state=%s run=%+v error=%v", status.State, status.Run, err)
	}
}

func TestBoundStartRejectsUnbindablePermissionDefinition(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	config := filepath.Join(root, "codex-home", "config.toml")
	raw, _ := os.ReadFile(config)
	if err := os.WriteFile(config, append(raw, []byte("fractional_value=0.5\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	before := recoveryFiles(t, filepath.Join(root, ".ply"))
	p, err := PreviewBoundStart(d, "unbindable-permissions", in)
	if err == nil || p.State != "rejected" || p.Error == nil {
		t.Fatalf("non-canonical policy produced approval: state=%s binding=%s error=%v", p.State, p.Request.Binding, err)
	}
	if !reflect.DeepEqual(before, recoveryFiles(t, filepath.Join(root, ".ply"))) {
		t.Fatal("invalid policy preview wrote state")
	}
}

func boundPublishSiblingProblem(t *testing.T, d taskrun.Dependencies, root string) {
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
	draft := map[string]any{"kind": "WorkspaceTaskProblemDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "sibling/revision", "task_id": "sibling", "registry_upgrade": nil, "expected_previous": previous, "origin": map[string]any{"kind": "authored"}, "title": "Revised sibling", "summary": "Independent revised need", "problem_document_id": "problem", "documents": []any{map[string]any{"id": "problem", "source": map[string]any{"kind": "file", "locator": doc, "sha256": digestBytes(body), "size_bytes": len(body), "media_type": "text/markdown", "git_provenance": nil}}}, "sources": []any{}, "claims": []any{}, "deadline": nil, "change_reason": "Independent publication", "recorder": map[string]any{"actor_claim": "fixture", "control_surface": "synthetic test", "recorded_at_utc": "2026-10-05T00:00:00Z"}}
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

func TestBoundStartUnrelatedContentDoesNotChangeConfirmation(t *testing.T) {
	d, in, root, _ := boundFixture(t)
	if _, err := workspace.CreateTask(d.Workspace, workspace.TaskCreateInput{TaskID: "sibling", Title: "Sibling", Description: "Independent content", ParentEpicID: "epic", ProjectID: "fixture", RepoID: "fixture"}); err != nil {
		t.Fatal(err)
	}
	p := boundPreview(t, d, in)
	boundPublishSiblingProblem(t, d, root)
	after := boundPreview(t, d, in)
	if after.Basis.RegistrySHA256 == p.Basis.RegistrySHA256 {
		t.Fatal("fixture did not publish a real registry change")
	}
	if after.Request.Binding != p.Request.Binding {
		t.Fatal("unrelated publication changed confirmation")
	}
	if _, err := ApplyBoundStart(d, p.Request); err != nil {
		t.Fatal(err)
	}
}
