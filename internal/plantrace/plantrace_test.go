package plantrace

import (
	"os"
	"slices"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/plan"
	"github.com/devdimensionlab/plybuild/internal/testutil"
)

func TestMain(m *testing.M) { os.Exit(testutil.Isolate(m)) }

func TestMergeEventsKeepsPlanJournalOrder(t *testing.T) {
	now := time.Date(2026, 11, 23, 12, 0, 0, 0, time.UTC)
	for _, tt := range []struct {
		name  string
		delta time.Duration
	}{
		{name: "clock moves backward", delta: -time.Second},
		{name: "same timestamp"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			first := planEvent("PLAN-1", 8, plan.JournalEntry{Type: "note", At: now.Format(time.RFC3339Nano)})
			second := planEvent("PLAN-1", 9, plan.JournalEntry{Type: "note", At: now.Add(tt.delta).Format(time.RFC3339Nano)})
			assertEventOrder(t, mergeEvents(nil, []Event{first, second}), first.EventID, second.EventID)
		})
	}
}

func TestMergeEventsOrdersJournalHeadsByInstant(t *testing.T) {
	start := "2026-11-23T12:00:00Z"
	halfSecond := "2026-11-23T12:00:00.5Z"
	sameInstant := "2026-11-23T13:00:00.5+01:00"
	taskEvents := []Event{
		{EventID: "jev_b", OccurredAt: &sameInstant, RecordedAt: &start},
		{EventID: "jev_a", RecordedAt: &halfSecond},
	}
	first := planEvent("PLAN-1", 0, plan.JournalEntry{Type: "note", At: start})
	second := planEvent("PLAN-1", 1, plan.JournalEntry{Type: "note", At: halfSecond})
	assertEventOrder(t, mergeEvents(taskEvents, []Event{first, second}), first.EventID, second.EventID, "jev_a", "jev_b")
}

func assertEventOrder(t *testing.T, events []Event, want ...string) {
	t.Helper()
	var got []string
	for _, event := range events {
		got = append(got, event.EventID)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("event order = %v, want %v", got, want)
	}
}

func TestUnreadableTimesAreReasons(t *testing.T) {
	good, bad := "2026-11-23T12:00:00Z", "yesterday"
	events := []Event{
		{EventID: "jev_a", TaskID: "task", OccurredAt: &bad, RecordedAt: &good},
		{EventID: "jev_b", TaskID: "task", OccurredAt: &good},
		planEvent("PLAN-1", 0, plan.JournalEntry{Type: "note", At: bad}),
	}
	got := timeReasons(events)
	want := []Reason{
		{"task", ReasonTimeUnreadable, `Event jev_a has an unreadable time "yesterday"; it is ordered first.`},
		{"", ReasonTimeUnreadable, `Event PLAN-1-journal-1 has an unreadable time "yesterday"; it is ordered first.`},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("timeReasons = %v, want %v", got, want)
	}
}
