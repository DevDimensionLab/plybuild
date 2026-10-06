package taskrun

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// Recovery needs positive process evidence, not a screenshot or the absence of
// an agent get response. Extend the ordinary Herdr stand-in only for this suite.
const deliveryRecoveryInspectStandin = `
if cmd == ["pane", "get"]:
    pane = {"workspace_id":"w-fixture", "tab_id":"tab-fixture", "pane_id":"w-fixture:p1",
            "terminal_id":"terminal-fixture",
            "agent_status":"unknown", "foreground_cwd":model["foreground_cwd"]}
    pane.update(model.get("recovery_pane", {}))
    print(json.dumps({"id":"fixture", "result":{"pane":pane}}))
    raise SystemExit(0)
if cmd == ["pane", "process-info"]:
    info = {"pane_id":"w-fixture:p1", "shell_pid":4242, "foreground_process_group_id":4242,
            "foreground_processes":[{"pid":4242,"name":"zsh","cwd":model["foreground_cwd"]}]}
    info.update(model.get("recovery_process_info", {}))
    print(json.dumps({"id":"fixture", "result":{"process_info":info}}))
    raise SystemExit(0)
if cmd == ["agent", "start"] and model.get("recovery_armed"):
    model["bootstrap_mode"] = "normal"
    model["bootstrap_started"] = False
    model["bootstrap_gets"] = 0
    model["bootstrap_settled_observed"] = False
    model["agent_session_id"] = "replacement-session"
`

func deliveryRecoveryNewFixture(t *testing.T) deliveryFixture {
	t.Helper()
	f := newDeliveryFixture(t, "codex")
	herdr, err := os.ReadFile(f.R.Herdr.Executable.Path)
	if err != nil {
		t.Fatal(err)
	}
	herdr = bytes.Replace(herdr, []byte("cmd = args[:2]\n"), []byte("cmd = args[:2]\n"+deliveryRecoveryInspectStandin), 1)
	if err = os.WriteFile(f.R.Herdr.Executable.Path, herdr, 0700); err != nil {
		t.Fatal(err)
	}
	f.R.Herdr.Executable.SHA256 = hash(herdr)
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	f.D.StartupProcessTree = func(shellPID int) ([]StartupProcess, error) {
		if shellPID != 4242 {
			return nil, errors.New("fixture inspected an unexpected shell")
		}
		return []StartupProcess{{PID: 4242, ParentPID: 1, ProcessGroupID: 4242}}, nil
	}
	f.D.HerdrTimeout = 500 * time.Millisecond
	workflowTestModel(t, f.workflowFixture, map[string]any{
		"bootstrap_mode": "no_session", "interactive_ready": true,
		"bootstrap_expected_prompt": workflowExpectedBootstrapPrompt,
		"bootstrap_run_root":        workflowPaths(f.R, 0).RunRoot,
	})
	return f
}

func deliveryRecoveryFixture(t *testing.T) (deliveryFixture, WorkflowRun, DeliveryStartRecoveryInput) {
	t.Helper()
	f := deliveryRecoveryNewFixture(t)
	o, err := workflowReadinessStart(t, f.workflowFixture)
	if err == nil || o.RunID == "" || o.Transport.AgentSessionID != "" {
		t.Fatalf("fixture did not retain a pre-Task startup: %+v %v", o, err)
	}
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || s.Phase != "bootstrap_attempted" || s.Acceptance != nil || s.StartDraft != nil || s.StartSHA256 != nil {
		t.Fatalf("fixture startup state: %+v %v", s, err)
	}
	if _, err = os.Lstat(filepathForRound(s, "prompt-attempt.json")); !os.IsNotExist(err) {
		t.Fatalf("fixture attempted Task prompt: %v", err)
	}
	in := deliveryRecoveryReplaceProvider(t, f, o)
	f.D.HerdrTimeout = 3 * time.Second
	return f, o, in
}

func deliveryRecoveryReplaceProvider(t *testing.T, f deliveryFixture, o WorkflowRun) DeliveryStartRecoveryInput {
	t.Helper()
	// Model the updater removing the originally pinned executable. The runtime
	// name on PATH now resolves to another immutable executable, as with Homebrew.
	if err := os.Remove(f.R.Runtime.Executable.Path); err != nil {
		t.Fatal(err)
	}
	provider := filepath.Join(f.R.WorkspaceRoot, "synthetic-provider-v2")
	control := filepath.Join(f.R.WorkspaceRoot, "synthetic-ply-v2")
	for _, path := range []string{provider, control} {
		if err := os.WriteFile(path, []byte("#!/bin/sh\n# replacement fixture binary\nexit 0\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	providerLink := filepath.Join(f.R.WorkspaceRoot, "bin", "codex")
	if err := os.Remove(providerLink); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(provider, providerLink); err != nil {
		t.Fatal(err)
	}
	workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_armed": true})
	return DeliveryStartRecoveryInput{RunID: o.RunID, ProviderExecutable: Executable{provider, hashFileTest(t, provider)}, ControlExecutable: Executable{control, hashFileTest(t, control)}, Timeout: 3 * time.Second}
}

func deliveryRecoveryFiles(t *testing.T, root string) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[path] = data
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func deliveryRecoveryUnchanged(t *testing.T, before map[string][]byte, mutable string) {
	t.Helper()
	for path, expected := range before {
		if path == mutable {
			continue
		}
		actual, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(actual, expected) {
			t.Fatalf("original artifact changed: %s (%v)", path, err)
		}
	}
}

func deliveryRecoveryTestPreview(t *testing.T, f deliveryFixture, in DeliveryStartRecoveryInput) DeliveryStartRecoveryPreview {
	t.Helper()
	p, err := WorkflowPreviewDeliveryStartRecovery(f.D, f.R.WorkspaceRoot, in)
	if err != nil || p.State != "ready" || p.Confirmation == "" {
		t.Fatalf("recovery preview: %+v %v", p, err)
	}
	return p
}

func deliveryRecoveryPreserveControl(t *testing.T, in DeliveryStartRecoveryInput, p DeliveryStartRecoveryPreview) {
	t.Helper()
	b, err := os.ReadFile(in.ControlExecutable.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.MkdirAll(filepath.Dir(p.ControlExecutable.Path), 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(p.ControlExecutable.Path, b, 0700); err != nil {
		t.Fatal(err)
	}
}

func TestDeliveryStartRecoveryPreviewPreservesAllState(t *testing.T) {
	f, o, in := deliveryRecoveryFixture(t)
	before := deliveryRecoveryFiles(t, workflowRoot(f.R.WorkspaceRoot))
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	p := deliveryRecoveryTestPreview(t, f, in)
	if p.RunID != o.RunID || p.ProviderExecutable != in.ProviderExecutable || p.ControlExecutable.Path == in.ControlExecutable.Path {
		t.Fatalf("preview lost identity or failed to reserve a distinct control path: %+v", p)
	}
	if _, err = os.Lstat(p.RecoveryDirectory); !os.IsNotExist(err) {
		t.Fatalf("preview created recovery artifacts: %v", err)
	}
	deliveryRecoveryUnchanged(t, before, "")
	after := deliveryRecoveryFiles(t, workflowRoot(f.R.WorkspaceRoot))
	if len(after) != len(before) {
		t.Fatalf("preview added run files: before=%d after=%d", len(before), len(after))
	}
	registryAfter, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || registry.RawSHA256 != registryAfter.RawSHA256 {
		t.Fatalf("preview changed the preparation/queue: %v", err)
	}
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("preview launched or prompted a provider")
	}
}

func TestDeliveryStartRecoveryPreservesHistoryAndUsesNewCallbacks(t *testing.T) {
	f, original, in := deliveryRecoveryFixture(t)
	before := deliveryRecoveryFiles(t, workflowRoot(f.R.WorkspaceRoot))
	registry, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil {
		t.Fatal(err)
	}
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	o, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
	if err != nil {
		t.Fatalf("recovery failed: %+v %v", o, err)
	}
	if o.RunID != original.RunID || o.RequestSHA256 != original.RequestSHA256 || o.Handoff != original.Handoff || o.Paths.RunRoot != original.Paths.RunRoot || o.Transport.PaneID != original.Transport.PaneID || o.Transport.TabID != original.Transport.TabID || o.Transport.AgentSessionID != "replacement-session" {
		t.Fatalf("recovery changed Task identity or failed to bind replacement: %+v", o)
	}
	if o.Paths.Context == original.Paths.Context || !strings.HasPrefix(o.Paths.Context, p.RecoveryDirectory+string(os.PathSeparator)) {
		t.Fatalf("recovery did not isolate its private context: %s", o.Paths.Context)
	}
	deliveryRecoveryUnchanged(t, before, filepath.Join(original.Paths.RunRoot, "state.json"))
	registryAfter, err := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if err != nil || registry.RawSHA256 != registryAfter.RawSHA256 {
		t.Fatalf("recovery replaced the original Task/preparation/queue: %v", err)
	}
	a := deliveryTestAcceptance(t, f, o)
	file := writeAny(t, f.R.WorkspaceRoot, "replacement-acceptance.json", a)
	if _, err = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, original.Paths.Context, file); err == nil {
		t.Fatal("old private context accepted a new startup")
	}
	if _, err = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file); err == nil {
		t.Fatal("old control executable accepted a new startup")
	}
	f.D.Executable = func() (string, error) { return p.ControlExecutable.Path, nil }
	o, err = WorkflowAccept(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, file)
	if err != nil || o.Delivery.PermissionState != "recipient_confirmed_contract" {
		t.Fatalf("replacement runtime acceptance failed: %+v %v", o, err)
	}
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "replacement-working", Phase: "working", Summary: "Replacement startup accepted its exact runtime.", Meaning: "A synthetic fixture proves callbacks; no actual product work or human pass is claimed.", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	o, err = WorkflowDeliveryReport(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "replacement-report.json", r))
	if err != nil || o.Delivery.Phase != "working" || len(o.Delivery.Events) != 1 {
		t.Fatalf("new delivery callback failed: %+v %v", o, err)
	}
	if _, err = WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation); err != nil {
		t.Fatalf("exact retry must show its existing recovery: %v", err)
	}
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || workflowTestCalls(t, f.workflowFixture, "agent start") != 2 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 3 {
		t.Fatal("recovery or repeated apply duplicated a startup, bootstrap or Task prompt")
	}
}

func TestDeliveryStartRecoveryConcurrentApplyStartsOnce(t *testing.T) {
	f, _, in := deliveryRecoveryFixture(t)
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent apply: %v", err)
		}
	}
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || workflowTestCalls(t, f.workflowFixture, "agent start") != 2 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 3 {
		t.Fatal("concurrent recovery replayed an effect")
	}
}

func TestDeliveryStartRecoveryLateOriginalSenderCannotReplaceAcceptedState(t *testing.T) {
	f := deliveryRecoveryNewFixture(t)
	p, err := WorkflowPreviewStart(f.D, f.File)
	if err != nil {
		t.Fatal(err)
	}
	paused, release, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})
	var once sync.Once
	resume := func() { once.Do(func() { close(release) }) }
	old := f.D
	old.Fault = func(point string) error {
		if point == "workflow_before_agent_send" {
			close(paused)
			<-release
		}
		return nil
	}
	var originalErr error
	go func() {
		_, originalErr = WorkflowStart(old, f.File, *p.(WorkflowPreview).Confirmation)
		close(finished)
	}()
	t.Cleanup(func() {
		resume()
		select {
		case <-finished:
		case <-time.After(10 * time.Second):
			t.Error("original fixture sender did not finish")
		}
	})
	select {
	case <-paused:
	case <-finished:
		t.Fatalf("original sender failed before the interleaving: %v", originalErr)
	case <-time.After(10 * time.Second):
		t.Fatal("original sender did not reach its send reservation")
	}
	o, err := WorkflowShow(f.D, f.R.WorkspaceRoot, workflowID(f.R))
	if err != nil {
		t.Fatal(err)
	}
	in := deliveryRecoveryReplaceProvider(t, f, o)
	f.D.HerdrTimeout = 3 * time.Second
	recovery := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, recovery)
	o, err = WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, recovery.Confirmation)
	if err != nil {
		t.Fatal(err)
	}
	f.D.Executable = func() (string, error) { return recovery.ControlExecutable.Path, nil }
	o = deliveryTestAccept(t, f, o)
	if o.Delivery.PermissionState != "recipient_confirmed_contract" {
		t.Fatal("replacement fixture did not accept its runtime")
	}
	statePath := filepath.Join(o.Paths.RunRoot, "state.json")
	before, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	resume()
	select {
	case <-finished:
	case <-time.After(10 * time.Second):
		t.Fatal("late original sender did not stop")
	}
	if originalErr == nil {
		t.Fatal("old sender did not notice that its startup was replaced")
	}
	after, err := os.ReadFile(statePath)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("late original sender overwrote accepted replacement state: %v", err)
	}
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
		t.Fatal("late original sender replayed an effect after replacement acceptance")
	}
}

func TestDeliveryStartRecoveryOldReadinessCannotAdoptReplacementStartup(t *testing.T) {
	f, o, in := deliveryRecoveryFixture(t)
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	f.D.Fault = func(point string) error {
		if point == "delivery_recovery_before_agent_send" {
			return errors.New("replacement sender paused after durable reservation")
		}
		return nil
	}
	if _, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("fixture did not stop after reserving the replacement")
	}
	f.D.Fault = nil
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || state.Recovery == nil || state.Phase != "agent_start_attempted" {
		t.Fatalf("replacement fixture is not waiting for its own sender: %+v %v", state, err)
	}
	before := deliveryRecoveryFiles(t, o.Paths.RunRoot)
	// This is the boundary an old readiness loop crosses after its preceding
	// observation completed, while a concurrent recovery acquired the run.
	if err = workflowBootstrap(f.D, f.R.WorkspaceRoot, o.RunID, time.Now().Add(time.Second), "original"); err == nil {
		t.Fatal("old readiness observation adopted the replacement generation")
	}
	deliveryRecoveryUnchanged(t, before, "")
	if len(deliveryRecoveryFiles(t, o.Paths.RunRoot)) != len(before) || workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("old readiness loop reserved or sent input on behalf of the replacement")
	}
}

func TestDeliveryStartRecoveryForegroundShellDoesNotProveProviderExit(t *testing.T) {
	shell := StartupProcess{PID: 4242, ParentPID: 1, ProcessGroupID: 4242}
	background := StartupProcess{PID: 4343, ParentPID: 4242, ProcessGroupID: 4343}
	for _, tc := range []struct {
		name string
		tree []StartupProcess
		err  error
	}{
		{"background_or_suspended_child_in_another_group", []StartupProcess{shell, background}, nil},
		{"provider_grandchild_in_another_group", []StartupProcess{shell, background, {PID: 4545, ParentPID: 4343, ProcessGroupID: 4545}}, nil},
		{"missing_shell_process", []StartupProcess{background}, nil},
		{"incomplete_tree", []StartupProcess{shell, {PID: 4545, ParentPID: 4444, ProcessGroupID: 4545}}, nil},
		{"process_observation_unavailable", nil, errors.New("synthetic process inspection denied")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, _, in := deliveryRecoveryFixture(t)
			f.D.StartupProcessTree = func(int) ([]StartupProcess, error) {
				return append([]StartupProcess(nil), tc.tree...), tc.err
			}
			before := deliveryRecoveryFiles(t, workflowRoot(f.R.WorkspaceRoot))
			p, err := WorkflowPreviewDeliveryStartRecovery(f.D, f.R.WorkspaceRoot, in)
			if err == nil && p.State == "ready" {
				t.Fatal("foreground-only shell observation hid an active or unobserved descendant")
			}
			deliveryRecoveryUnchanged(t, before, "")
			if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
				t.Fatal("unknown or live process tree caused a new provider effect")
			}
		})
	}
}

func TestDeliveryStartRecoveryRechecksProcessTreeBeforeApply(t *testing.T) {
	f, _, in := deliveryRecoveryFixture(t)
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	f.D.StartupProcessTree = func(int) ([]StartupProcess, error) {
		return []StartupProcess{{PID: 4242, ParentPID: 1, ProcessGroupID: 4242}, {PID: 4343, ParentPID: 4242, ProcessGroupID: 4343}}, nil
	}
	if _, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
		t.Fatal("an earlier exit preview allowed a provider start while a background process existed")
	}
	if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("recovery failed to check descendants before startup effects")
	}
}

func TestDeliveryStartRecoveryRejectsUnsafeStartEvidence(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*testing.T, deliveryFixture, WorkflowRun)
	}{
		{"provider_active", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_pane": map[string]any{"agent": "codex", "agent_status": "working"}})
		}},
		{"foreground_process_unknown", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_process_info": map[string]any{"foreground_processes": []any{}}})
		}},
		{"provider_process_remains", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_process_info": map[string]any{"foreground_process_group_id": 4343, "foreground_processes": []any{map[string]any{"pid": 4343, "name": "codex", "cwd": f.Prepared.Preparation.Preparation.Plan.WorktreePath}}}})
		}},
		{"terminal_changed", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_pane": map[string]any{"terminal_id": "another-terminal"}})
		}},
		{"shell_cwd_changed", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_process_info": map[string]any{"foreground_processes": []any{map[string]any{"pid": 4242, "name": "zsh", "cwd": f.Parent}}}})
		}},
		{"orphan_prompt_attempt", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			writeAny(t, filepath.Join(o.Paths.RunRoot, "rounds", "000"), "prompt-attempt.json", map[string]any{"argv": []string{"agent", "prompt", "possibly-sent"}})
		}},
		{"orphan_acceptance", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			writeAny(t, o.Paths.RunRoot, "acceptance.json", map[string]any{"acceptance": "started"})
		}},
		{"orphan_start_draft", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			writeAny(t, o.Paths.RunRoot, "start-draft.json", map[string]any{"started": true})
		}},
		{"dirty_task", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			if err := os.WriteFile(filepath.Join(f.Prepared.Preparation.Preparation.Plan.WorktreePath, "untracked-work.txt"), []byte("possible recipient effects"), 0600); err != nil {
				t.Fatal(err)
			}
		}},
		{"task_base_changed", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			runGit(t, f.Prepared.Preparation.Preparation.Plan.WorktreePath, "commit", "--allow-empty", "-m", "unexpected base change")
		}},
		{"goal_input_changed", func(t *testing.T, f deliveryFixture, o WorkflowRun) {
			if err := os.WriteFile(f.R.Delivery.Goal.Locator, []byte("Changed frozen goal\n"), 0600); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, o, in := deliveryRecoveryFixture(t)
			tc.change(t, f, o)
			before := deliveryRecoveryFiles(t, workflowRoot(f.R.WorkspaceRoot))
			if p, err := WorkflowPreviewDeliveryStartRecovery(f.D, f.R.WorkspaceRoot, in); err == nil && p.State == "ready" {
				t.Fatalf("unsafe recovery became start-ready: %+v", p)
			}
			deliveryRecoveryUnchanged(t, before, "")
			if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
				t.Fatal("rejected recovery sent provider effects")
			}
		})
	}
}

func TestDeliveryStartRecoveryRechecksPreviewBeforeApply(t *testing.T) {
	for _, change := range []string{"provider_bytes", "control_bytes", "task_dirty", "orphan_prompt", "provider_returned"} {
		t.Run(change, func(t *testing.T) {
			f, o, in := deliveryRecoveryFixture(t)
			p := deliveryRecoveryTestPreview(t, f, in)
			deliveryRecoveryPreserveControl(t, in, p)
			var path string
			switch change {
			case "provider_bytes":
				path = in.ProviderExecutable.Path
			case "control_bytes":
				path = p.ControlExecutable.Path
			case "task_dirty":
				path = filepath.Join(f.Prepared.Preparation.Preparation.Plan.WorktreePath, "drift.txt")
			case "orphan_prompt":
				path = filepath.Join(o.Paths.RunRoot, "rounds", "000", "prompt-attempt.json")
			case "provider_returned":
				workflowTestModel(t, f.workflowFixture, map[string]any{"recovery_pane": map[string]any{"agent": "codex", "agent_status": "working"}})
			}
			if path != "" {
				if err := os.WriteFile(path, []byte("Changed after preview\n"), 0700); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation); err == nil {
				t.Fatal("preview authorized a changed runtime or target")
			}
			if workflowTestCalls(t, f.workflowFixture, "agent start") != 1 || workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
				t.Fatal("drift was detected after another startup effect")
			}
		})
	}
}

func TestDeliveryStartRecoveryLostStartupResponseNeverReplays(t *testing.T) {
	f, _, in := deliveryRecoveryFixture(t)
	p := deliveryRecoveryTestPreview(t, f, in)
	deliveryRecoveryPreserveControl(t, in, p)
	workflowTestModel(t, f.workflowFixture, map[string]any{"lost_command": "agent start", "observations": []any{workflowPendingObservation()}})
	o, err := WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
	if err == nil || o.RunID != in.RunID {
		t.Fatalf("lost startup response did not preserve same unknown run: %+v %v", o, err)
	}
	_, _ = WorkflowRecoverDeliveryStart(f.D, f.R.WorkspaceRoot, in, p.Confirmation)
	_, _ = WorkflowShow(f.D, f.R.WorkspaceRoot, in.RunID)
	if workflowTestCalls(t, f.workflowFixture, "tab create") != 1 || workflowTestCalls(t, f.workflowFixture, "agent start") != 2 {
		t.Fatal("unknown recovery startup was replayed")
	}
	// A lost start response must not cause another Task input to be delivered.
	for _, call := range workflowBootstrapPrompts(t, f.workflowFixture) {
		if call.Argv[3] != workflowExpectedBootstrapPrompt {
			t.Fatal("unknown replacement startup received Task input")
		}
	}
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil || state.Recovery == nil || state.Result.StartupRecovery == nil {
		t.Fatalf("unknown replacement lost its durable recovery binding: %v", err)
	}
}
