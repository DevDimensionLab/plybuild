package quota

import (
	"testing"
	"time"
)

var testNow = time.Unix(1794070000, 0).UTC()

func at(offset time.Duration) *time.Time {
	t := testNow.Add(offset)
	return &t
}

func TestRemainingClampsUsedPercent(t *testing.T) {
	for used, want := range map[float64]float64{38: 62, 82: 18, 130: 0, -4: 100, 0: 100} {
		if got := (Window{ID: FiveHour, UsedPercent: used}).Remaining(); got != want {
			t.Errorf("used %v: remaining %v, want %v", used, got, want)
		}
	}
}

func TestWindowLabelsUseDuration(t *testing.T) {
	for minutes, want := range map[int]string{300: "5-hour", 10080: "Weekly", 1440: "1-day", 120: "2-hour", 90: "90-minute"} {
		if got := (WindowID{Minutes: minutes}).Label(); got != want {
			t.Errorf("%d minutes: %q, want %q", minutes, got, want)
		}
	}
}

func TestSelectNeverCombinesWindows(t *testing.T) {
	windows := []Window{
		{ID: FiveHour, UsedPercent: 38, ResetsAt: at(time.Hour)},
		{ID: Weekly, UsedPercent: 82, ResetsAt: at(48 * time.Hour)},
		{ID: Weekly, Pool: "Spark", UsedPercent: 99},
	}
	cases := []struct {
		choice Choice
		want   float64
		window WindowID
	}{{ChooseFiveHour, 62, FiveHour}, {ChooseWeekly, 18, Weekly}, {ChooseLowest, 18, Weekly}}
	for _, c := range cases {
		got := Select(c.choice, windows, testNow)
		if got.Kind != SelectedValue || got.Remaining != c.want || got.Window.ID != c.window || got.Window.Pool != "" {
			t.Errorf("choice %v: %+v, want %v of %v", c.choice, got, c.want, c.window)
		}
	}
}

func TestSelectPendingAndNotReported(t *testing.T) {
	pending := []Window{{ID: FiveHour, UsedPercent: 90, ResetsAt: at(-time.Second)}, {ID: Weekly, UsedPercent: 10}}
	if got := Select(ChooseLowest, pending, testNow); got.Kind != SelectedPending || got.Window.ID != FiveHour {
		t.Fatalf("lowest with passed reset = %+v, want pending 5-hour", got)
	}
	if got := Select(ChooseWeekly, pending, testNow); got.Kind != SelectedValue || got.Remaining != 90 {
		t.Fatalf("weekly ignores 5-hour pending: %+v", got)
	}
	poolOnly := []Window{{ID: Weekly, Pool: "Spark", UsedPercent: 1}}
	for _, windows := range [][]Window{nil, poolOnly} {
		if got := Select(ChooseLowest, windows, testNow); got.Kind != SelectedNotReported || got.Window != nil {
			t.Fatalf("Select(%v) = %+v, want not reported", windows, got)
		}
	}
	if got := Select(ChooseFiveHour, []Window{{ID: Weekly, UsedPercent: 5}}, testNow); got.Kind != SelectedNotReported {
		t.Fatalf("weekly-only five-hour = %+v", got)
	}
}

func TestSelectTiesKeepInputOrder(t *testing.T) {
	windows := []Window{{ID: Weekly, UsedPercent: 50}, {ID: FiveHour, UsedPercent: 50}}
	if got := Select(ChooseLowest, windows, testNow); got.Window.ID != Weekly {
		t.Fatalf("tie chose %v", got.Window.ID)
	}
}

func TestStaleBySource(t *testing.T) {
	codex := Snapshot{Source: SourceCodex, ObservedAt: testNow.Add(-121 * time.Second)}
	claude := Snapshot{Source: SourceClaude, ObservedAt: testNow.Add(-121 * time.Second)}
	if !codex.Stale(testNow) || claude.Stale(testNow) {
		t.Fatalf("stale codex=%v claude=%v", codex.Stale(testNow), claude.Stale(testNow))
	}
	claude.ObservedAt = testNow.Add(-3601 * time.Second)
	if !claude.Stale(testNow) {
		t.Fatal("claude older than an hour is not stale")
	}
}
