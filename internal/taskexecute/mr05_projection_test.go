//go:build darwin

package taskexecute

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// The child runs the real native owner. Only the provider probe and attempt
// launcher are synthetic; they cannot attest provider or sandbox acceptance.
// StartWorker releases this direct child, so the parent explicitly wait4s it.
func mr05ProjectionOwner(t *testing.T) bool {
	if os.Getenv("PLY_MR05_PROJECTION_OWNER") != "1" {
		return false
	}
	root, id := os.Getenv("PLY_MR05_PROJECTION_ROOT"), os.Getenv("PLY_MR05_PROJECTION_RUN")
	defer func() { _ = os.WriteFile(filepath.Join(root, "owner.finished"), nil, 0600) }()
	d := taskrun.SystemDependencies(workspace.SystemDependencies())
	d.Executable = func() (string, error) { return os.Getenv("PLY_MR05_PROJECTION_CONTROL"), nil }
	err := taskrun.ServeWorker(d, root, id)
	if err != nil {
		t.Fatal(err)
	}
	run, err := taskrun.ShowWorker(d, root, id)
	if err != nil || run.Worker == nil || run.Worker.Phase != "blocked" || !strings.Contains(run.Worker.Reason, "exit status 73") {
		t.Fatalf("native owner did not observe the released synthetic attempt: %+v %v", run.Worker, err)
	}
	return true
}

type mr05ProjectionFixture struct {
	d       taskrun.Dependencies
	root    string
	result  Result
	request taskrun.WorkflowRequest
	goal    string
}

func mr05ProjectionQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }

func mr05ProjectionFixtureNew(t *testing.T) *mr05ProjectionFixture {
	t.Helper()
	// Native acceptance's TMPDIR may be inside the user's workspace. A real
	// workspace.Init must instead use an owned directory outside that workspace.
	scratch, err := os.MkdirTemp("/private/tmp", "ply-mr05-projection-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(scratch); err != nil {
			t.Error(err)
		}
	})
	t.Setenv("TMPDIR", scratch)
	t.Setenv("GOTMPDIR", scratch)
	home := filepath.Join(scratch, "home")
	if err := os.Mkdir(home, 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	testBinary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal(err)
	}
	a := workspace.DeliveryAgreement{SchemaVersion: 3, Mode: workspace.DeliveryLocalEpic, Acceptance: &workspace.DeliveryAcceptancePolicy{SchemaVersion: 1, Mode: "automatic", ResponsibleActor: "synthetic projection fixture only", TimeoutSeconds: 30}}
	d, in, root, _ := launcherFixtureForProvider(t, "codex", a)
	t.Setenv("HERDR_ENV", "")
	t.Setenv("HERDR_WORKSPACE_ID", "")
	in.Runtime.HerdrWorkspace = ""
	in.Worker = &taskrun.WorkerContract{SchemaVersion: 1, MaxAttempts: 1, AgentTimeoutSeconds: 60}
	in.ControlContract = taskrun.NewControlContract("standalone", "", "")
	adapter := filepath.Join(root, "attempt.py")
	// This stand-in inherits the actual owner lock. It does no provider work,
	// remains at a bounded barrier, and exits only when the test releases it.
	standIn := `import json, os, pathlib, sys, time
root = pathlib.Path(__file__).parent
args = sys.argv[1:]
assert len(args) == 8 and args[:3] == ['workflow','worker','_attempt']
assert args[4:6] == ['--workspace', str(root)] and args[6] == '--job'
os.fstat(3)
(root/'attempt.ready').write_text(json.dumps({'argv':args,'pid':os.getpid(),'parent':os.getppid()}))
deadline = time.monotonic() + 60
while not (root/'attempt.release').exists():
    if time.monotonic() >= deadline: raise SystemExit(74)
    time.sleep(0.02)
raise SystemExit(73)
`
	if err := os.WriteFile(adapter, []byte(standIn), 0600); err != nil {
		t.Fatal(err)
	}
	launcher := "#!/bin/sh\nset -eu\n" +
		"test \"$1/$2\" = workflow/worker\n" +
		"case \"$3\" in\n_supervise)\n" +
		"printf '%s' \"$$\" > " + mr05ProjectionQuote(filepath.Join(root, "owner.pid")) + "\n" +
		"export PLY_MR05_PROJECTION_OWNER=1\n" +
		"export PLY_MR05_PROJECTION_CONTROL=\"$0\"\n" +
		"export PLY_MR05_PROJECTION_ROOT=" + mr05ProjectionQuote(root) + "\n" +
		"export PLY_MR05_PROJECTION_RUN=\"$4\"\n" +
		"exec " + mr05ProjectionQuote(testBinary) + " -test.run='^TestMR05ProjectionNativeGoalDraft$' -test.timeout=90s\n;;\n" +
		"_attempt) exec " + mr05ProjectionQuote(python) + " -I -B " + mr05ProjectionQuote(adapter) + " \"$@\";;\n*) exit 92;;\nesac\n"
	if err := os.WriteFile(filepath.Join(root, "bin", "ply"), []byte(launcher), 0700); err != nil {
		t.Fatal(err)
	}
	// The native probe invokes this explicit test substitute. No actual Codex,
	// account, policy override, model call or execution attempt is available.
	probe := "#!/bin/sh\nset -eu\nfor arg do\n case \"$arg\" in\n sandbox) printf '%s\\n' synthetic-probe-only >> " + mr05ProjectionQuote(filepath.Join(root, "probe.calls")) + "; exit 0;;\n mcp) printf '%s\\n' synthetic-empty-mcp >> " + mr05ProjectionQuote(filepath.Join(root, "probe.calls")) + "; printf '[]\\n'; exit 0;;\n esac\ndone\nexit 93\n"
	if err := os.WriteFile(filepath.Join(root, "bin", "codex"), []byte(probe), 0700); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { mr05ProjectionReap(t, root) })
	out, err := Execute(d, in)
	if err != nil || out.Run == nil || out.Run.Worker == nil {
		t.Fatalf("real Goal execution failed: %+v %v", out, err)
	}
	f := &mr05ProjectionFixture{d: d, root: root, result: out}
	if _, err := readPreservedJSON(out.RequestPath, &f.request); err != nil {
		t.Fatal(err)
	}
	var draft struct {
		Binding struct {
			Project string `json:"project_id"`
			Repo    string `json:"repo_id"`
		} `json:"binding_request"`
		Task struct {
			Spec struct {
				Hash string `json:"manifest_sha256"`
			} `json:"spec"`
		} `json:"task_spec_binding"`
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(f.request.HandoffDraft, &draft) != nil || json.Unmarshal(f.request.HandoffDraft, &fields) != nil {
		t.Fatal("invalid genuine handoff draft")
	}
	if _, fabricated := fields["project_binding"]; fabricated || draft.Binding.Project != string(out.Goal.Target.ProjectID) || draft.Binding.Repo != string(out.Goal.Target.RepoID) || draft.Binding.Project != "fixture" || draft.Binding.Repo != "fixture" || draft.Task.Spec.Hash == "" {
		t.Fatal("native Goal execution did not preserve the registered target in binding_request")
	}
	f.goal = draft.Task.Spec.Hash
	deadline := time.Now().Add(25 * time.Second)
	for {
		if _, err := os.Stat(filepath.Join(root, "attempt.ready")); err == nil {
			break
		}
		_, finished := os.Stat(filepath.Join(root, "owner.finished"))
		if time.Now().After(deadline) || finished == nil {
			log, _ := os.ReadFile(filepath.Join(out.Run.Paths.RunRoot, "worker", "controller.log"))
			t.Fatalf("native owner did not reach the attempt barrier: %s", log)
		}
		time.Sleep(20 * time.Millisecond)
	}
	b, err := taskrun.ShowControl(d, root, out.Run.RunID)
	if err != nil || b.Phase != "running" || b.OwnerLock != "lock_held" || b.JobID == nil || b.UnavailableReason != nil {
		t.Fatalf("fixture is not at a real supported native boundary: %+v %v", b, err)
	}
	t.Logf("genuine binding_request=%s/%s; run=%s request=%s session=%s execution_goal=%s; synthetic probe/attempt only", draft.Binding.Project, draft.Binding.Repo, out.Run.RunID, out.Run.RequestSHA256, out.Run.Worker.SessionID, f.goal)
	return f
}

func mr05ProjectionReap(t *testing.T, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, "attempt.release"), nil, 0600); err != nil {
		t.Error(err)
	}
	var pid int
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		raw, err := os.ReadFile(filepath.Join(root, "owner.pid"))
		if err == nil {
			pid, err = strconv.Atoi(string(raw))
			if err != nil || pid <= 1 {
				t.Fatalf("invalid owned child PID: %q", raw)
			}
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if pid == 0 {
		t.Error("no native owner PID was observed; cleanup is unconfirmed")
		return
	}
	var status syscall.WaitStatus
	for time.Now().Before(deadline) {
		n, err := syscall.Wait4(pid, &status, syscall.WNOHANG, nil)
		if err != nil {
			t.Errorf("cannot positively reap owned native child %d: %v", pid, err)
			return
		}
		if n == pid {
			if !status.Exited() || status.ExitStatus() != 0 {
				t.Errorf("native owner exited unexpectedly: %v", status)
			}
			t.Logf("reaped native owner pid=%d exit=%d", pid, status.ExitStatus())
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	// wait4 has positively established our still-unreaped direct child. Its
	// StartWorker session owns the local attempt, so stop only that process group.
	if pgid, err := syscall.Getpgid(pid); err == nil && pgid == pid {
		_ = syscall.Kill(-pid, syscall.SIGKILL)
		_, _ = syscall.Wait4(pid, &status, 0, nil)
	}
	t.Error("owned native child required forced cleanup; acceptance fails")
}

func mr05ProjectionInventory(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		value := info.Mode().String()
		if !entry.IsDir() {
			raw, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			value += digestBytes(raw)
		}
		out[path] = value
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func (f *mr05ProjectionFixture) assertBasis(t *testing.T, b taskrun.ControlBasis) {
	t.Helper()
	for name, pair := range map[string]struct {
		got  *string
		want string
	}{
		"project": {b.ProjectID, string(f.result.Goal.Target.ProjectID)}, "repo": {b.RepoID, string(f.result.Goal.Target.RepoID)},
		"run": {b.RunID, f.result.Run.RunID}, "request": {b.RunRequestHash, f.result.Run.RequestSHA256},
		"worker": {b.WorkerSession, f.result.Run.Worker.SessionID}, "execution_goal": {b.Goal, f.goal},
		"frozen_control": {b.FrozenControl, f.request.Runtime.PlyExecutable.SHA256},
	} {
		if pair.got == nil || *pair.got != pair.want {
			t.Errorf("known native %s lost: got=%v want=%q", name, pair.got, pair.want)
		}
	}
	if b.Workspace != f.root || b.Scope != (taskrun.ControlTarget{Kind: "run", ID: f.result.Run.RunID}) || b.UnavailableReason != nil || b.StateRevision == "" || !reflect.DeepEqual(b.Contract, f.request.ControlContract) {
		t.Errorf("public identity/eligibility changed: %+v", b)
	}
}

func (f *mr05ProjectionFixture) assertPreview(t *testing.T, p taskrun.ControlPreview) {
	t.Helper()
	if p.State != "ready" || p.Error != nil || p.Basis == nil || p.Request == nil {
		t.Fatalf("expected genuinely ready complete preview: %+v", p)
	}
	f.assertBasis(t, *p.Basis)
	if taskrun.ValidateControlRequest(*p.Request) != nil || p.Request.Workspace != f.root || p.Request.Input.Target != p.Basis.Scope || p.Request.Input.Action != "pause" {
		t.Fatal("preview lost complete approved request")
	}
}

// A separate historical interactive record has no known project identity.
// It is never used to produce the positive native Goal/worker assertions.
func mr05ProjectionLegacy(t *testing.T, root string) string {
	t.Helper()
	var emptySlices func(reflect.Value)
	emptySlices = func(v reflect.Value) {
		switch v.Kind() {
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				emptySlices(v.Field(i))
			}
		case reflect.Slice:
			if v.IsNil() {
				v.Set(reflect.MakeSlice(v.Type(), 0, 0))
			}
		}
	}
	r := taskrun.WorkflowRequest{Envelope: taskrun.Envelope{Kind: "ply.workflow.herdr-run-request", SchemaVersion: 1}, WorkspaceRoot: root, RequestKey: "projection-legacy-fixture", Runtime: taskrun.Runtime{Mode: "interactive", Provider: "codex"}, HandoffDraft: json.RawMessage(`{}`)}
	emptySlices(reflect.ValueOf(&r).Elem())
	digest := func(v any) string {
		raw, err := taskrun.Canonical(v)
		if err != nil {
			t.Fatal(err)
		}
		return digestBytes(raw)
	}
	id := "wfr_" + strings.TrimPrefix(digest([]string{root, r.RequestKey}), "sha256:")
	run := taskrun.WorkflowRun{RunID: id, SessionID: "ply:" + id, RequestSHA256: digest(r), Paths: taskrun.WorkflowPaths{RunRoot: filepath.Join(root, ".ply", "workflow-runs", "runs", id)}}
	emptySlices(reflect.ValueOf(&run).Elem())
	observed := taskrun.Observed{}
	emptySlices(reflect.ValueOf(&observed).Elem())
	fixtureJSON(t, filepath.Join(root, ".ply", "workflow-runs", "requests", id+".json"), map[string]any{"request": r, "observed": observed, "result": run, "phase": "", "context_sha256": "", "acceptance": nil, "start_draft": nil, "start_sha256": nil, "records": []any{}})
	return id
}

func mr05ProjectionAssertLegacy(t *testing.T, b taskrun.ControlBasis, p taskrun.ControlPreview) {
	t.Helper()
	if b.ProjectID != nil || b.RepoID != nil || b.WorkerSession != nil || len(b.Actions) != 0 || b.UnavailableReason == nil || *b.UnavailableReason != "runtime_unsupported" {
		t.Fatalf("unsupported legacy identity was invented or enabled: %+v", b)
	}
	if p.State != "rejected" || p.Request != nil || p.Error == nil || p.Error.Code != "control_unavailable" || p.Basis == nil || p.Basis.ProjectID != nil || p.Basis.RepoID != nil {
		t.Fatalf("legacy preview became approvable: %+v", p)
	}
}

func TestMR05ProjectionNativeGoalDraft(t *testing.T) {
	if mr05ProjectionOwner(t) {
		return
	}
	f := mr05ProjectionFixtureNew(t)
	legacy := mr05ProjectionLegacy(t, f.root)
	before := mr05ProjectionInventory(t, f.root)
	var revision string
	for i := 0; i < 2; i++ {
		b, err := taskrun.ShowControl(f.d, f.root, f.result.Run.RunID)
		if err != nil {
			t.Fatal(err)
		}
		f.assertBasis(t, b)
		p, err := taskrun.WorkerControls(f.d).Preview(f.root, f.result.Run.RunID, "pause", "projection-preview-0001")
		if err != nil {
			t.Fatal(err)
		}
		f.assertPreview(t, p)
		if b.StateRevision != p.Basis.StateRevision || (i > 0 && revision != b.StateRevision) {
			t.Fatal("read changed native state revision")
		}
		revision = b.StateRevision
	}
	t.Run("legacy_identity_unavailable", func(t *testing.T) {
		b, err := taskrun.ShowControl(f.d, f.root, legacy)
		if err != nil {
			t.Fatal(err)
		}
		p, err := taskrun.WorkerControls(f.d).Preview(f.root, legacy, "pause", "legacy-preview-0001")
		if err == nil || taskrun.ControlFailureOf(err).Code != "control_unavailable" {
			t.Fatal("legacy preview changed stable error", err)
		}
		mr05ProjectionAssertLegacy(t, b, p)
	})
	for _, tc := range []struct{ id, code string }{{"wfr_" + strings.Repeat("a", 64), "control_not_found"}, {"week_" + strings.Repeat("a", 32), "control_invalid_arguments"}} {
		if _, err := taskrun.ShowControl(f.d, f.root, tc.id); err == nil || taskrun.ControlFailureOf(err).Code != tc.code {
			t.Errorf("missing/wrong scope changed error: %v", err)
		}
	}
	if !reflect.DeepEqual(before, mr05ProjectionInventory(t, f.root)) {
		t.Fatal("show/preview changed native files, Git, launch/probe inventory or reserved a control")
	}
}

func TestMR05ProjectionCandidateCLI(t *testing.T) {
	// Resolve inputs before launcherFixture changes CWD. No build/skip fallback.
	binary, python := os.Getenv("PLY_CANDIDATE_BINARY"), os.Getenv("PLY_MR05_SCHEMA_PYTHON")
	for _, path := range []string{binary, python} {
		physical, err := filepath.EvalSymlinks(path)
		info, statErr := os.Stat(path)
		if !filepath.IsAbs(path) || err != nil || statErr != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 || (path == binary && physical != path) {
			t.Fatalf("required exact candidate/schema interpreter unavailable: %q %v %v", path, err, statErr)
		}
	}
	schema, err := filepath.Abs("../../schemas/workflow-control-v1.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	candidate, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	f := mr05ProjectionFixtureNew(t)
	legacy := mr05ProjectionLegacy(t, f.root)
	before := mr05ProjectionInventory(t, f.root)
	worktrees := fixtureGit(t, filepath.Join(f.root, "main"), "worktree", "list", "--porcelain")
	call := func(t *testing.T, success, schemaResult bool, args ...string) []byte {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, binary, args...)
		cmd.Dir = f.result.Goal.WorktreePath
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		raw, err := cmd.Output()
		var exit any
		if cmd.ProcessState != nil {
			exit = cmd.ProcessState.ExitCode()
		}
		t.Logf("candidate=%s sha256=%s argv=%q cwd=%s exit=%v stdout=%s stderr=%s", binary, digestBytes(candidate), args, cmd.Dir, exit, raw, stderr.String())
		if (err == nil) != success {
			t.Fatalf("unexpected actual CLI result: %v", err)
		}
		if schemaResult {
			validate := exec.CommandContext(ctx, python, "-I", "-B", "-c", `import importlib.metadata,json,jsonschema,sys
assert importlib.metadata.version('jsonschema') == '4.26.0'
jsonschema.Draft202012Validator(json.load(open(sys.argv[1])),format_checker=jsonschema.FormatChecker()).validate(json.load(sys.stdin))`, schema)
			validate.Stdin = bytes.NewReader(raw)
			if output, err := validate.CombinedOutput(); err != nil {
				t.Fatalf("actual CLI violates delivered schema: %v %s", err, output)
			}
		}
		return append(raw, []byte(stderr.String())...)
	}
	var revision string
	for i := 0; i < 2; i++ {
		for _, verb := range []string{"show", "preview"} {
			args := []string{"workflow", "control", verb, f.result.Run.RunID, "--workspace", f.root, "--format", "json"}
			if verb == "preview" {
				args = append(args, "--action", "pause", "--request-id", "projection-cli-0001")
			}
			raw := call(t, true, true, args...)
			var basis taskrun.ControlBasis
			if verb == "show" {
				if err := json.Unmarshal(raw, &basis); err != nil || basis.Kind != "ply.workflow.control-capabilities" {
					t.Fatalf("invalid show: %s %v", raw, err)
				}
				f.assertBasis(t, basis)
			} else {
				var p taskrun.ControlPreview
				if err := json.Unmarshal(raw, &p); err != nil || p.Kind != "ply.workflow.control-preview" {
					t.Fatalf("invalid preview: %s %v", raw, err)
				}
				f.assertPreview(t, p)
				basis = *p.Basis
			}
			if revision != "" && revision != basis.StateRevision {
				t.Fatal("CLI reads changed native revision")
			}
			revision = basis.StateRevision
		}
	}
	t.Run("legacy_identity_unavailable", func(t *testing.T) {
		var b taskrun.ControlBasis
		if err := json.Unmarshal(call(t, true, true, "workflow", "control", "show", legacy, "--workspace", f.root, "--format", "json"), &b); err != nil {
			t.Fatal(err)
		}
		var p taskrun.ControlPreview
		raw := call(t, false, true, "workflow", "control", "preview", legacy, "--workspace", f.root, "--format", "json", "--action", "pause", "--request-id", "legacy-cli-0001")
		// The rejected preview also emits a diagnostic on stderr. Decode the
		// first JSON value while the schema assertion above uses stdout alone.
		if err := json.NewDecoder(bytes.NewReader(raw)).Decode(&p); err != nil {
			t.Fatal(err)
		}
		mr05ProjectionAssertLegacy(t, b, p)
	})
	for _, tc := range []struct{ id, code string }{{"wfr_" + strings.Repeat("a", 64), "control_not_found"}, {"week_" + strings.Repeat("a", 32), "control_invalid_arguments"}} {
		raw := call(t, false, false, "workflow", "control", "show", tc.id, "--workspace", f.root, "--format", "json")
		if !strings.Contains(string(raw), tc.code) {
			t.Fatalf("missing/wrong scope changed public error: %s", raw)
		}
	}
	if !reflect.DeepEqual(before, mr05ProjectionInventory(t, f.root)) || worktrees != fixtureGit(t, filepath.Join(f.root, "main"), "worktree", "list", "--porcelain") {
		t.Fatal("CLI show/preview changed native/Git/provider inventory")
	}
	after, err := os.ReadFile(binary)
	if err != nil || !bytes.Equal(candidate, after) {
		t.Fatal("exact candidate bytes changed", err)
	}
	t.Logf("exact candidate CLI and delivered schema checked; binary SHA=%s", digestBytes(candidate))
}
