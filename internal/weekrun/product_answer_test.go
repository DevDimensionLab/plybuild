package weekrun

import (
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
