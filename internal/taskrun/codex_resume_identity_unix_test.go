//go:build !windows

package taskrun

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// This uses an actual OS-held Codex-format thread lock. Only Herdr is modeled;
// it deliberately reports no hook metadata until after the first turn.
func TestCodexResumeDiscoversLiveThreadBeforeFirstTurn(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	f.In.ProviderExecutable = Executable{executable, hash(data)}
	f.In.Timeout = 15 * time.Second
	link := filepath.Join(f.R.WorkspaceRoot, "provider-bin", "codex")
	if err = os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err = os.Symlink(executable, link); err != nil {
		t.Fatal(err)
	}
	lock := filepath.Join(os.Getenv("HOME"), ".codex", "thread-writer-locks", f.Owner.Transport.AgentSessionID+".lock")
	if err = os.MkdirAll(filepath.Dir(lock), 0700); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(executable, "-test.run=^TestCodexResumeWriterHelper$")
	cmd.Env = append(os.Environ(), "PLY_TEST_CODEX_LOCK="+lock)
	cmd.Dir = f.Prepared.Preparation.Preparation.Plan.WorktreePath
	input, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	output, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { input.Close(); cmd.Wait() })
	reader := bufio.NewReader(output)
	if line, err := reader.ReadString('\n'); err != nil || line != "locked\n" {
		t.Fatalf("lock helper: %q %v", line, err)
	}
	// The writer-absence boundary is separately covered by the real liveness
	// regressions. Here the child models the new process started by Herdr.
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	workflowTestModel(t, f.workflowFixture, map[string]any{"session_absent_gets": 100000, "resume_processes": []any{map[string]any{"pid": cmd.Process.Pid, "name": "codex", "cwd": cmd.Dir, "argv": []string{"codex", "resume", f.Owner.Transport.AgentSessionID}}}})
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err != nil {
		pending, _ := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
		t.Fatalf("live same-thread writer was not discovered before input: %v; diagnostic: %s", err, pending.Detail)
	}
	if o.OwnerResume == nil || !o.OwnerResume.SessionObserved || o.OwnerResume.Step != ownerResumePromptSent {
		t.Fatalf("same session not observed: %+v", o.OwnerResume)
	}

	if n := workflowTestCalls(t, f.workflowFixture, "agent prompt"); n != 2 {
		t.Fatalf("expected one continuation, got %d", n)
	}
	state, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	processes := func(changes map[string]any) []any {
		process := map[string]any{"pid": cmd.Process.Pid, "name": "codex", "cwd": cmd.Dir, "argv": []string{"codex", "resume", o.Transport.AgentSessionID}}
		for key, value := range changes {
			process[key] = value
		}
		return []any{process}
	}
	reject := func(t *testing.T) {
		t.Helper()
		agent, err := workflowAgentRead(f.D, state, time.Time{})
		if err == nil {
			err = workflowSession(state, agent, false)
		}
		if err == nil {
			t.Fatal("unproven runtime identity accepted")
		}
		if workflowTestCalls(t, f.workflowFixture, "agent prompt") != 2 {
			t.Fatal("inspection sent input")
		}
	}
	for _, tc := range []struct {
		name    string
		changes map[string]any
	}{
		{"wrong_pid", map[string]any{"pid": os.Getpid()}},
		{"wrong_cwd", map[string]any{"cwd": "/"}},
		{"wrong_live_thread_argument", map[string]any{"argv": []string{"codex", "resume", "wrong-session"}}},
		{"missing_live_arguments", map[string]any{"argv": nil}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflowTestModel(t, f.workflowFixture, map[string]any{"resume_processes": processes(tc.changes)})
			reject(t)
		})
	}
	t.Run("multiple_foreground_writers", func(t *testing.T) {
		both := append(processes(nil), processes(nil)...)
		workflowTestModel(t, f.workflowFixture, map[string]any{"resume_processes": both})
		reject(t)
	})
	workflowTestModel(t, f.workflowFixture, map[string]any{"resume_processes": processes(nil)})
	for _, mode := range []string{"unlock", "shared"} {
		t.Run(mode+"_is_not_exclusive", func(t *testing.T) {
			input.Write([]byte(mode + "\n"))
			if line, err := reader.ReadString('\n'); err != nil || line != "locked\n" {
				t.Fatalf("helper mode: %q %v", line, err)
			}
			reject(t)
		})
	}
	input.Write([]byte("exclusive\n"))
	if _, err := reader.ReadString('\n'); err != nil {
		t.Fatal(err)
	}
	t.Run("another_process_holds_the_lock_file", func(t *testing.T) {
		other := exec.Command(executable, "-test.run=^TestCodexResumeWriterHelper$")
		other.Env = append(os.Environ(), "PLY_TEST_CODEX_LOCK="+lock, "PLY_TEST_CODEX_LOCK_MODE=open")
		input, err := other.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		output, err := other.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		if err = other.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() { input.Close(); other.Wait() }()
		if line, err := bufio.NewReader(output).ReadString('\n'); err != nil || line != "locked\n" {
			t.Fatalf("second holder: %q %v", line, err)
		}
		reject(t)
	})
	t.Run("missing_lock", func(t *testing.T) {
		if err := os.Rename(lock, lock+".saved"); err != nil {
			t.Fatal(err)
		}
		defer os.Rename(lock+".saved", lock)
		reject(t)
	})
	t.Run("inspection_denied", func(t *testing.T) {
		workflowTestModel(t, f.workflowFixture, map[string]any{"resume_inspection_denied": true})
		defer workflowTestModel(t, f.workflowFixture, map[string]any{"resume_inspection_denied": false})
		reject(t)
	})
	// An actual callback also works before hook publication, through only the
	// resumed control. The persisted prompt keeps the independent evidence.
	proof, err := os.ReadFile(filepath.Join(p.Directory, "prompt-attempt.json"))
	if err != nil || !bytes.Contains(proof, []byte("local_codex_exclusive_thread_lock_v1")) {
		t.Fatalf("runtime observation not preserved: %v", err)
	}
	resumed := f.D
	resumed.Executable = func() (string, error) { return p.ControlExecutable.Path, nil }
	o = ownerResumeReport(t, f.deliveryFixture, o, resumed, "runtime-lock-callback")
	if o.OwnerResume.OwnerCallback != "observed" {
		t.Fatalf("callback not genuine: %+v", o.OwnerResume)
	}
}

func TestCodexResumeWriterHelper(t *testing.T) {
	path := os.Getenv("PLY_TEST_CODEX_LOCK")
	if path == "" {
		return
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if os.Getenv("PLY_TEST_CODEX_LOCK_MODE") != "open" {
		if err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
			t.Fatal(err)
		}
	}
	os.Stdout.WriteString("locked\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		mode := scanner.Text()
		if err = syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
			t.Fatal(err)
		}
		flags := 0
		switch mode {
		case "shared":
			flags = syscall.LOCK_SH
		case "exclusive":
			flags = syscall.LOCK_EX
		case "unlock":
		default:
			t.Fatal("unknown helper mode")
		}
		if flags != 0 {
			if err = syscall.Flock(int(file.Fd()), flags|syscall.LOCK_NB); err != nil {
				t.Fatal(err)
			}
		}
		os.Stdout.WriteString("locked\n")
	}
}

// A known no-input session_pending attempt has a coordinator-owned exit;
// positive Herdr absence must not hide a surviving unmanaged local writer.
func TestCodexResumeSessionPendingRetireChecksLocalWriters(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	workflowTestModel(t, f.workflowFixture, map[string]any{"session_absent_gets": 100000})
	f.In.Timeout = 1500 * time.Millisecond
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	o, err := WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation)
	if err == nil || o.OwnerResume.Step != ownerResumeStartAttempted {
		t.Fatalf("not pending: %+v %v", o.OwnerResume, err)
	}
	pending, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err == nil || pending.Reason != "session_pending" {
		t.Fatalf("not known session_pending: %+v %v", pending, err)
	}
	before := deliveryRecoveryFiles(t, o.Paths.RunRoot)
	owner := f.D
	retireAsCoordinator(&f.deliveryFixture)
	in := DeliveryRetireInput{Observer: "fixture root", Reason: "Bounded exit from known session_pending without owner callback"}
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid()}, {PID: 65001, ParentPID: 1, Session: "transcript_open"}}, nil
	}
	if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in); err == nil || !strings.Contains(err.Error(), "writer") {
		t.Fatalf("retirement ignored surviving writer: %v", err)
	}
	if retireQueueCurrent(t, f.deliveryFixture) == nil {
		t.Fatal("blocked retirement released queue")
	}
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	out, err := WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, in)
	if err != nil {
		t.Fatal(err)
	}
	if out.Delivery.Phase != "stopped" || retireQueueCurrent(t, f.deliveryFixture) != nil {
		t.Fatal("known pending exit did not release only this queue")
	}
	deliveryRecoveryUnchanged(t, before, filepath.Join(o.Paths.RunRoot, "state.json"))
	if _, err = os.Stat(f.Prepared.Preparation.Preparation.Plan.WorktreePath); err != nil {
		t.Fatal("work lost")
	}
	owner.Executable = func() (string, error) { return p.ControlExecutable.Path, nil }
	late := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: o.RunID, RequestSHA256: o.RequestSHA256, SessionID: o.SessionID, EventID: "late", PreviousEventSHA256: out.Delivery.LastEventSHA256, Phase: "working", Summary: "Late callback", Meaning: "Must stay closed", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	if _, err = WorkflowDeliveryReport(owner, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, writeAny(t, f.R.WorkspaceRoot, "late-retired.json", late)); err == nil || !strings.Contains(err.Error(), "owner_retired") {
		t.Fatalf("late callback accepted: %v", err)
	}
	if workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("pending exit sent input")
	}
}

func TestCodexResumeUnknownTransportCannotRetire(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	failed := f.D
	failed.Fault = func(point string) error {
		if point == "owner_resume_before_tab_send" {
			return workflowError(5, "unknown tab effect fixture")
		}
		return nil
	}
	if _, err = WorkflowResumeOwner(failed, f.R.WorkspaceRoot, f.In, p.Confirmation); err == nil {
		t.Fatal("fault not reached")
	}
	retireAsCoordinator(&f.deliveryFixture)
	if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, f.Owner.RunID, DeliveryRetireInput{Observer: "fixture root", Reason: "Must not forget an unknown tab effect"}); err == nil || !strings.Contains(err.Error(), "unknown_effect") {
		t.Fatalf("unknown tab attempt retired: %v", err)
	}
	if retireQueueCurrent(t, f.deliveryFixture) == nil {
		t.Fatal("unknown attempt released queue")
	}
	if workflowTestCalls(t, f.workflowFixture, "pane close") != 0 {
		t.Fatal("unknown attempt closed transport")
	}
}

func TestCodexResumeInspectionFailureDoesNotPollDeniedRead(t *testing.T) {
	f := newOwnerResumeFixture(t, "codex", false)
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
	workflowTestModel(t, f.workflowFixture, map[string]any{"session_absent_gets": 100000, "resume_inspection_denied": true})
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err != nil {
		t.Fatal(err)
	}
	ownerResumePreserve(t, f, p)
	f.In.Timeout = 1500 * time.Millisecond
	if _, err = WorkflowResumeOwner(f.D, f.R.WorkspaceRoot, f.In, p.Confirmation); err == nil {
		t.Fatal("denied inspection permitted continuation")
	}
	if n := workflowTestCalls(t, f.workflowFixture, "pane process-info"); n != 1 {
		t.Fatalf("denied read automatically repeated %d times", n)
	}
	if workflowTestCalls(t, f.workflowFixture, "agent prompt") != 1 {
		t.Fatal("denied inspection sent input")
	}
}

// An unrelated fixture's real process must not look like this fixture's
// competing writer merely because a test adapter reused a native session ID.
func TestOwnerResumeFixturesDoNotShareWriterIdentity(t *testing.T) {
	t.Setenv("HERDR_WORKSPACE_ID", "")
	binary := deliveryCLIBinary(t)
	first := newOwnerResumeFixture(t, "codex", false)
	second := newOwnerResumeFixture(t, "claude", false)
	holder := exec.Command("/bin/sh", "-c", "read ignored", first.Owner.Transport.AgentSessionID)
	holder.Dir = "/"
	input, err := holder.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err = holder.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { input.Close(); holder.Wait() }()
	task := second.Prepared.Preparation.Preparation.Plan.WorktreePath
	result, output, err := ownerResumeCLI(binary, task, "--resume", "--check", "--format", "json")
	if err != nil || deliveryCLIString(t, result, "state") != "ready" {
		t.Fatalf("unrelated fixture writer collided: first=%s second=%s err=%v output=%s", first.Owner.Transport.AgentSessionID, second.Owner.Transport.AgentSessionID, err, output)
	}
}
