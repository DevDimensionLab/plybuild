package weekrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestSL01ParentRestartRetainsAwaitingTask(t *testing.T) {
	f := newMR05Week(t)
	f.tick(t)
	r := f.n.runs["first"]
	r.Phase = "awaiting_product_answer"
	f.n.runs["first"] = r
	before := len(f.n.starts)
	f.reload(t)
	f.tick(t)
	if f.s.Phase != "awaiting_product_answer" || f.s.Entries[0].State != "awaiting_product_answer" || f.s.ActiveTaskID != "first" || len(f.n.starts) != before || len(f.n.resumes) != 0 || len(f.n.released) != 0 || f.s.Deadline != f.b.Deadline || !strings.Contains(f.s.NextAction, "product answer") {
		t.Fatalf("parent lost explicit wait: %+v starts=%v resumes=%v", f.s, f.n.starts, f.n.resumes)
	}
	if taskrun.PublicControlState(f.s.Phase) != "awaiting_product_answer" {
		t.Fatal("awaiting projected as unknown")
	}
}
func TestSL01WeekPassesExplicitQuestionContract(t *testing.T) {
	f := newMR05Week(t)
	f.c.ProductQuestions = &taskrun.ProductQuestionContract{SchemaVersion: 1, Mode: taskrun.BoundAnswerMode, MaxContinuationTurns: 1}
	if e := Validate(f.c); e != nil {
		t.Fatal(e)
	}
	w := workerContract(f.c)
	if w.SchemaVersion != 2 || w.ProductQuestions == nil || w.ProductQuestions.Mode != taskrun.BoundAnswerMode {
		t.Fatal(w)
	}
	f.c.SchemaVersion = 1
	f.c.RuntimeControl = nil
	if Validate(f.c) == nil {
		t.Fatal("legacy week silently upgraded")
	}
}

func TestSL01WeekZeroProviderRetriesRequiresNewQuestionMode(t *testing.T) {
	f := newMR05Week(t)
	f.control(t, "stop")
	f.tick(t)
	f.c.MaxAttempts, f.c.MaxProviderRetries = 1, 0
	if Validate(f.c) == nil {
		t.Fatal("legacy bound week accepted zero retries")
	}
	f.c.ProductQuestions = &taskrun.ProductQuestionContract{SchemaVersion: 1, Mode: taskrun.BoundAnswerMode, MaxContinuationTurns: 1}
	if e := Validate(f.c); e != nil {
		t.Fatalf("explicit SL01 zero retries rejected: %v", e)
	}
	f.c.Name = "zero-retry-week"
	f.n.c = f.c
	s, e := f.m.Start(f.c, false)
	if e != nil {
		t.Fatal(e)
	}
	b, _, e := readBinding(f.c.Workspace, s.ID)
	if e != nil || b.Contract.MaxAttempts != 1 || b.Contract.MaxProviderRetries != 0 {
		t.Fatalf("frozen week lost limits: %+v %v", b, e)
	}
	w := workerContract(b.Contract)
	if w.SchemaVersion != 2 || w.ProductQuestions == nil || w.MaxAttempts != 1 || w.MaxProviderRetries != 0 {
		t.Fatalf("NativeAdapter mapping lost explicit zero: %+v", w)
	}
}

// NA01: the bound wait label belongs only to a native Run that still asks for
// its answer. Continuation, review and completion replace it.
func TestNA01WeekLeavesAnswerWaitWhenNativeRunContinues(t *testing.T) {
	f := newMR05Week(t)
	f.tick(t)
	f.tick(t)
	r := f.n.runs["first"]
	r.Phase = "awaiting_product_answer"
	f.n.runs["first"] = r
	f.tick(t)
	if f.s.Entries[0].State != "awaiting_product_answer" || f.s.NextAction != productAnswerWait {
		t.Fatalf("unanswered question must still wait: %+v", f.s)
	}
	r.Phase, r.Stage = "running", "reviewer"
	f.n.runs["first"] = r
	f.tick(t)
	if f.s.Phase != "running" || f.s.Entries[0].State != "running" || strings.Contains(f.s.NextAction, "product answer") {
		t.Fatalf("continued Run still projects the answer wait: %+v", f.s)
	}
	f.n.complete("first")
	f.tick(t)
	if f.s.Entries[0].State != "completed" || strings.Contains(f.s.NextAction, "product answer") {
		t.Fatalf("completed Task still projects the answer wait: %+v", f.s)
	}
}

// Historical state written by the earlier controller stays byte-identical; the
// read projects the current native facts without dispatching anything.
func TestNA01WeekShowProjectsStaleAnswerWaitReadOnly(t *testing.T) {
	f := newMR05Week(t)
	f.tick(t)
	f.tick(t)
	f.s.Phase, f.s.NextAction = "running", productAnswerWait
	f.s.Entries[0].State = "awaiting_product_answer"
	if e := publish(filepath.Join(f.dir, "state.json"), f.s, true); e != nil {
		t.Fatal(e)
	}
	before, e := os.ReadFile(filepath.Join(f.dir, "state.json"))
	if e != nil {
		t.Fatal(e)
	}
	r := f.n.runs["first"]
	r.Phase, r.Stage = "running", "reviewer"
	f.n.runs["first"] = r
	starts, resumes, controls := len(f.n.starts), len(f.n.resumes), len(f.n.controls)
	for i := 0; i < 2; i++ {
		s, e := f.m.Show(f.c.Workspace, f.s.ID)
		if e != nil {
			t.Fatal(e)
		}
		if s.Entries[0].State != "running" || strings.Contains(s.NextAction, "product answer") {
			t.Fatalf("read kept the stale answer wait: %+v", s)
		}
	}
	after, e := os.ReadFile(filepath.Join(f.dir, "state.json"))
	if e != nil || string(after) != string(before) || len(f.n.starts) != starts || len(f.n.resumes) != resumes || len(f.n.controls) != controls {
		t.Fatalf("read changed state or dispatched: %v starts=%v resumes=%v controls=%v", e, f.n.starts, f.n.resumes, f.n.controls)
	}
	f.n.complete("first")
	f.s.ActiveTaskID, f.s.Phase = "", "completed"
	f.s.Entries[0].State, f.s.Entries[1].State = "completed", "completed"
	if e = publish(filepath.Join(f.dir, "state.json"), f.s, true); e != nil {
		t.Fatal(e)
	}
	s, e := f.m.Show(f.c.Workspace, f.s.ID)
	if e != nil || strings.Contains(s.NextAction, "product answer") {
		t.Fatalf("completed week still asks for an answer: %+v %v", s, e)
	}
}

// NA01 review focus: a cached answer wait read before the next Tick agrees with
// the native Run's current boundary, not merely without the old wait text.
func TestNA01WeekShowAgreesWithNativeBoundaryBeforeTick(t *testing.T) {
	for _, c := range []struct {
		name, phase, entry, next string
		run                      func(*Run)
	}{
		{"completed", "running", "completed", "completed", func(r *Run) { r.Completed, r.Phase = true, "completed" }},
		{"paused", "paused", "running", "parent pause", func(r *Run) {
			r.Phase, r.NextAction = "paused", taskrun.WorkflowAction{Actor: "parent_control", Message: "The product answer and same Task are retained under parent pause."}
		}},
		{"stopped", "stopped", "running", "parent stop", func(r *Run) {
			r.Phase, r.NextAction = "stopped", taskrun.WorkflowAction{Actor: "parent_control", Message: "The product answer and same Task are retained under parent stop."}
		}},
		{"blocked", "blocked", "running", "Worker blocked", func(r *Run) {
			r.Phase, r.NextAction = "blocked", taskrun.WorkflowAction{Actor: "human", Message: "Worker blocked with preserved work: continuation_unknown"}
		}},
		{"unknown", "unknown", "unknown", "observation failed", func(r *Run) { r.ID = "missing-run" }},
	} {
		t.Run(c.name, func(t *testing.T) {
			f := newMR05Week(t)
			f.tick(t)
			f.tick(t)
			f.s.Phase, f.s.NextAction = "awaiting_product_answer", productAnswerWait
			f.s.Entries[0].State = "awaiting_product_answer"
			if e := publish(filepath.Join(f.dir, "state.json"), f.s, true); e != nil {
				t.Fatal(e)
			}
			r := f.n.runs["first"]
			c.run(&r)
			f.n.runs["first"] = r
			s, e := f.m.Show(f.c.Workspace, f.s.ID)
			if e != nil || s.Phase != c.phase || s.Entries[0].State != c.entry || !strings.Contains(s.NextAction, c.next) || s.NextAction == productAnswerWait {
				t.Fatalf("cached week %s/%s %q, want %s/%s containing %q: %v", s.Phase, s.Entries[0].State, s.NextAction, c.phase, c.entry, c.next, e)
			}
		})
	}
}

func TestNA01WeekTickLeavesAnswerWaitAtNativeBoundary(t *testing.T) {
	for _, phase := range []string{"paused", "blocked"} {
		t.Run(phase, func(t *testing.T) {
			f := newMR05Week(t)
			f.tick(t)
			f.tick(t)
			r := f.n.runs["first"]
			r.Phase = "awaiting_product_answer"
			f.n.runs["first"] = r
			f.tick(t)
			r.Phase, r.NextAction = phase, taskrun.WorkflowAction{Actor: "human", Message: "current native boundary " + phase}
			f.n.runs["first"] = r
			f.tick(t)
			if f.s.Phase != phase || f.s.Entries[0].State != "running" || !strings.Contains(f.s.NextAction, "current native boundary") {
				t.Fatalf("tick kept the answer wait at %s: %s/%s %q", phase, f.s.Phase, f.s.Entries[0].State, f.s.NextAction)
			}
		})
	}
}
