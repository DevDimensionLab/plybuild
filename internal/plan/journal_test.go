package plan

import (
	"reflect"
	"testing"
)

func TestJournalKeepsRecordsOfAnExistingPlanInOrder(t *testing.T) {
	store, _ := testStore(t)
	if _, err := store.AppendJournal("PLAN-1", JournalEntry{Type: "note", Title: "x", ActorID: "a", Role: "planner"}); err == nil {
		t.Fatal("journal of a missing plan accepted")
	}
	if _, err := store.Add([]byte(specExample), "a", "r"); err != nil {
		t.Fatal(err)
	}
	if got, err := store.Journal("PLAN-1"); err != nil || len(got) != 0 {
		t.Fatalf("empty journal: %v %v", got, err)
	}
	id, kind, outcome := "map", "planning", "complete"
	first := JournalEntry{Type: "step_started", Title: "Map", ActorID: "ply-plan", Role: "planner", StepID: &id, StepKind: &kind}
	second := JournalEntry{Type: "step_finished", Title: "Map", ActorID: "ply-plan", Role: "planner", StepID: &id, StepKind: &kind, Outcome: &outcome}
	for _, e := range []JournalEntry{first, second} {
		if written, err := store.AppendJournal("PLAN-1", e); err != nil || written.At == "" {
			t.Fatalf("append %+v %v", written, err)
		}
	}
	got, err := store.Journal("PLAN-1")
	if err != nil || len(got) != 2 || got[0].At > got[1].At {
		t.Fatalf("readback %+v %v", got, err)
	}
	got[0].At, got[1].At = "", ""
	if !reflect.DeepEqual(got, []JournalEntry{first, second}) {
		t.Fatalf("records changed: %+v", got)
	}
	if list, err := store.List(""); err != nil || len(list.Items) != 1 || len(list.Invalid) != 0 {
		t.Fatalf("journal file disturbs list: %+v %v", list, err)
	}
	if _, err := store.AppendJournal("PLAN-1", JournalEntry{Type: "step_started", Title: "x", ActorID: "a", Role: "planner"}); err == nil {
		t.Fatal("step without ID accepted")
	}
}
