package weekrun

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/currentstate"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestAR01AnswerBlockRecoveryIgnoresDisplayWording(t *testing.T) {
	f := newMR05Week(t)
	f.tick(t)
	f.tick(t)
	r := f.n.runs["first"]
	r.Phase = "awaiting_product_answer"
	r.NextAction = taskrun.WorkflowAction{Actor: "human", Message: "Continuation outcome is unknown; inspect this Run."}
	f.n.runs["first"] = r
	f.tick(t)
	if f.s.Phase != "unknown" || f.s.Entries[0].State != "blocked" {
		t.Fatalf("fixture did not record an answer-derived block: %+v", f.s)
	}
	f.s.NextAction = "The display wording has changed."
	if err := publish(filepath.Join(f.dir, "state.json"), f.s, true); err != nil {
		t.Fatal(err)
	}
	r.Phase, r.Stage = "running", "reviewer"
	r.NextAction = taskrun.WorkflowAction{Actor: "worker", Message: "The answer use is confirmed; review is current."}
	f.n.runs["first"] = r
	before, err := os.ReadFile(filepath.Join(f.dir, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	starts, resumes, controls, releases := len(f.n.starts), len(f.n.resumes), len(f.n.controls), len(f.n.released)
	shown, err := f.m.Show(f.c.Workspace, f.s.ID)
	if err != nil || shown.Phase != "running" || shown.Entries[0].State != "running" {
		t.Fatalf("Show retained a resolved answer block: %+v: %v", shown, err)
	}
	after, err := os.ReadFile(filepath.Join(f.dir, "state.json"))
	if err != nil || !bytes.Equal(before, after) || starts != len(f.n.starts) || resumes != len(f.n.resumes) || controls != len(f.n.controls) || releases != len(f.n.released) {
		t.Fatal("Show changed history or invoked a native effect")
	}
	f.tick(t)
	if f.s.Phase != shown.Phase || f.s.Entries[0].State != shown.Entries[0].State {
		t.Fatalf("Tick disagrees with Show: %+v / %+v", f.s, shown)
	}
}

func TestAR01WeekSeesReceivedAnswerBeforeFirstWaitTick(t *testing.T) {
	f := newMR05Week(t)
	f.tick(t)
	f.tick(t)
	r := f.n.runs["first"]
	r.Phase = "awaiting_product_answer"
	r.NextAction = taskrun.WorkflowAction{Actor: "worker", Message: "Recorded answer; continuation has not started."}
	f.n.runs["first"] = r
	shown, err := f.m.Show(f.c.Workspace, f.s.ID)
	if err != nil || shown.Phase != "running" || !strings.Contains(shown.NextAction, r.NextAction.Message) {
		t.Fatalf("Show requires a prior cached wait label: %+v: %v", shown, err)
	}
	f.tick(t)
	if f.s.Phase != shown.Phase || f.s.NextAction != shown.NextAction {
		t.Fatalf("Tick disagrees with the pre-Tick read: %+v / %+v", f.s, shown)
	}
}

func TestAR01WeekKeepsUnresolvedDispatch(t *testing.T) {
	f := newMR05Week(t)
	f.tick(t)
	f.tick(t)
	f.s.Phase, f.s.NextAction = "unknown", "Inspect the reserved native effect"
	f.s.Dispatch = &Dispatch{Kind: "resume", TaskID: "first", RunID: f.s.Entries[0].RunID, Sequence: "reserved-dispatch"}
	if err := publish(filepath.Join(f.dir, "state.json"), f.s, true); err != nil {
		t.Fatal(err)
	}
	shown, err := f.m.Show(f.c.Workspace, f.s.ID)
	if err != nil || shown.Phase != f.s.Phase || shown.NextAction != f.s.NextAction || shown.Dispatch == nil {
		t.Fatalf("read hid an unresolved native effect: %+v: %v", shown, err)
	}
}

func TestAR01LegacyInitialRegistrationKeepsOriginalBytes(t *testing.T) {
	f := newMR05Week(t)
	f.s.GuidanceOrigin = ""
	for i := range f.s.Entries {
		f.s.Entries[i].BlockOrigin = ""
	}
	path := filepath.Join(f.dir, "state.json")
	if err := publish(path, f.s, true); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.Remove(taskrun.ParentControlPath(f.c.Workspace, f.b.ID)); err != nil {
		t.Fatal(err)
	}
	s, err := f.m.Start(f.c, false)
	if err != nil || s.StartedAt != f.s.StartedAt || s.Deadline != f.s.Deadline || len(f.n.starts) != 0 {
		t.Fatalf("historical initial shape cannot recover its original gate: %+v: %v", s, err)
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("registration recovery rewrote the historical state")
	}
}

func TestAR01OrdinaryBlockIgnoresNativeLookingDisplay(t *testing.T) {
	f := newMR05Week(t)
	f.tick(t)
	f.tick(t)
	f.s.Phase, f.s.Entries[0].State = "blocked", "blocked"
	f.s.Entries[0].BlockOrigin = currentstate.OtherOrigin
	f.s.GuidanceOrigin = currentstate.NativeRun
	f.s.NextAction = nativeNextPrefix + "worker: the display can name any actor"
	if err := publish(filepath.Join(f.dir, "state.json"), f.s, true); err != nil {
		t.Fatal(err)
	}
	shown, err := f.m.Show(f.c.Workspace, f.s.ID)
	if err != nil || shown.Phase != "blocked" || shown.Entries[0].State != "blocked" {
		t.Fatalf("Show cleared an ordinary block: %+v: %v", shown, err)
	}
	f.tick(t)
	if f.s.Phase != "blocked" || f.s.Entries[0].State != "blocked" {
		t.Fatalf("Tick cleared an ordinary block: %+v", f.s)
	}
}
