package taskrun

import (
	"os"
	"path/filepath"
	"testing"
)

// Composition tests cross the same interface as action callers. Herdr and
// LocalProcesses are the existing fixture adapters, not a second policy model.
func TestAR03ComposedWriterEvidence(t *testing.T) {
	f := newOwnerResumeFixture(t, "claude", false)
	s, err := workflowRead(f.R.WorkspaceRoot, f.Owner.RunID)
	if err != nil {
		t.Fatal(err)
	}
	task := s.Observed.Target.WorktreeLocator
	for _, tc := range []struct {
		name  string
		table []LocalProcess
		err   error
		want  taskWriterState
	}{
		{"absent", []LocalProcess{{PID: os.Getpid(), CWD: task}}, nil, taskWritersAbsent},
		{"outside_herdr", []LocalProcess{{PID: os.Getpid()}, {PID: 4242, CWD: filepath.Join(task, "src")}}, nil, taskWritersPresent},
		{"session_elsewhere", []LocalProcess{{PID: os.Getpid()}, {PID: 4242, CWD: "/", Session: "transcript_open"}}, nil, taskWritersPresent},
		{"denied", nil, os.ErrPermission, taskWritersUnknown},
		{"empty", nil, nil, taskWritersUnknown},
		{"incomplete_missing_self", []LocalProcess{{PID: 1, CWD: "/"}}, nil, taskWritersUnknown},
		{"duplicate_identity", []LocalProcess{{PID: os.Getpid()}, {PID: os.Getpid()}}, nil, taskWritersUnknown},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := f.D
			calls := 0
			d.LocalProcesses = func(string, string) ([]LocalProcess, error) { calls++; return tc.table, tc.err }
			o := (taskWriters{d, s}).cleanup()
			if o.state != tc.want || (o.err == nil) != (tc.want == taskWritersAbsent) {
				t.Fatalf("state=%v err=%v; want %v", o.state, o.err, tc.want)
			}
			if calls != 1 {
				t.Fatalf("inspection calls=%d", calls)
			}
			if tc.want == taskWritersAbsent && o.absence.ProcessLiveness != "no_task_or_session_process_observed" {
				t.Fatal("absence lost composed basis")
			}
			historical, err := observeAbsentDeliveryOwner(d, s)
			if err != nil || historical.ProcessLiveness != "unknown" || calls != 1 {
				t.Fatal("human recovery gained local authority", err)
			}
		})
	}
	t.Run("incomplete_herdr", func(t *testing.T) {
		workflowTestModel(t, f.workflowFixture, map[string]any{"inventory": []any{map[string]any{"pane_id": "incomplete"}}})
		d := f.D
		d.LocalProcesses = func(string, string) ([]LocalProcess, error) {
			t.Fatal("inspected processes after unknown Herdr evidence")
			return nil, nil
		}
		o := (taskWriters{d, s}).cleanup()
		if o.state != taskWritersUnknown || o.herdrPresent {
			t.Fatalf("incomplete inventory=%+v", o)
		}
	})
}

func TestAR03HeadlessCleanupUsesActualSession(t *testing.T) {
	s := workflowState{}
	s.Request.Envelope = deliveryEnv("headless-run-request")
	s.Result.Worker = &WorkerStatus{SessionID: "actual-headless-session"}
	d := Dependencies{LocalProcesses: func(session, transcript string) ([]LocalProcess, error) {
		if session != "actual-headless-session" {
			t.Errorf("inspected session %q", session)
		}
		return []LocalProcess{{PID: os.Getpid()}}, nil
	}}
	o := (taskWriters{d, s}).cleanup()
	if o.state != taskWritersAbsent {
		t.Fatalf("cleanup=%+v", o)
	}
}

// A direct parent is excludable only as a plain invoking shell; positive
// session evidence is still a writer even if Herdr calls the pane plain.
func TestAR03ResumeInvokerSessionEvidenceStillBlocks(t *testing.T) {
	f := newOwnerResumeFixture(t, "claude", false)
	f.In.InvokerPaneID = "invoker-pane"
	task := f.Prepared.Preparation.Preparation.Plan.WorktreePath
	workflowTestModel(t, f.workflowFixture, map[string]any{"invoker_shell_pid": 900, "inventory": []any{map[string]any{"workspace_id": "w-fixture", "tab_id": "invoker-tab", "pane_id": "invoker-pane", "terminal_id": "invoker-terminal", "cwd": task}}})
	f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: 900, CWD: task, Session: "transcript_open"}, {PID: os.Getpid(), ParentPID: 900, CWD: task}}, nil
	}
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err == nil || p.Reason != "writer_process_alive" {
		t.Fatalf("session writer was excluded: %+v %v", p, err)
	}
}

func TestAR03ResumeDoesNotReplayDeniedManagedAgentInspection(t *testing.T) {
	f := newOwnerResumeFixture(t, "claude", false)
	workflowTestModel(t, f.workflowFixture, map[string]any{"managed_inspection_denied": true})
	before := workflowTestCalls(t, f.workflowFixture, "agent get")
	p, err := WorkflowPreviewOwnerResume(f.D, f.R.WorkspaceRoot, f.In)
	if err == nil || p.Reason != "inspection_unavailable" {
		t.Fatalf("denied inspection permitted resume: %+v %v", p, err)
	}
	if n := workflowTestCalls(t, f.workflowFixture, "agent get") - before; n != 1 {
		t.Fatalf("replayed denied inspection %d times", n)
	}
}
