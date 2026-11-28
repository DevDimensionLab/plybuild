package quotahistory

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func point(when string, used float64) Point { return Point{At: at(when), Used: used} }

func record(t *testing.T, s Store, series string, p Point) bool {
	t.Helper()
	written, err := s.Record(series, p)
	if err != nil {
		t.Fatal(err)
	}
	return written
}

func read(t *testing.T, s Store, ids []string, since, until string) History {
	t.Helper()
	h, err := s.Read(ids, at(since), at(until))
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestRecordKeepsOnePointPerChangeOrTenMinutes(t *testing.T) {
	s := Open(t.TempDir())
	steps := []struct {
		p    Point
		want bool
	}{
		{point("2026-11-28T10:00:00Z", 40), true},
		{point("2026-11-28T10:05:00Z", 40), false},                                                         // same value within 10 min
		{point("2026-11-28T10:06:00Z", 41), true},                                                          // changed value
		{point("2026-11-28T10:06:00Z", 42), false},                                                         // not newer than the last point
		{point("2026-11-28T10:01:00Z", 43), false},                                                         // older than the last point
		{point("2026-11-28T10:16:00Z", 41), true},                                                          // ten minutes after the last point
		{Point{At: at("2026-11-28T10:17:00Z"), Used: 41, ResetsAt: ptr(at("2026-12-05T00:00:00Z"))}, true}, // reset changed
	}
	for i, step := range steps {
		if got := record(t, s, "claude.weekly", step.p); got != step.want {
			t.Fatalf("step %d: written %v, want %v", i, got, step.want)
		}
	}
	h := read(t, s, []string{"claude.weekly"}, "2026-11-28T00:00:00Z", "2026-11-29T00:00:00Z")
	if len(h.Series) != 1 || len(h.Series[0].Points) != 4 || h.Series[0].Unit != "percent" || h.Series[0].Points[3].ResetsAt == nil {
		t.Fatalf("%+v", h)
	}
	if _, err := s.Record("claude.daily", point("2026-11-28T10:00:00Z", 1)); err == nil {
		t.Fatal("unknown series recorded")
	}
}

func ptr[T any](v T) *T { return &v }

func TestRecordComparesAcrossTheMonthBoundary(t *testing.T) {
	dir := t.TempDir()
	s := Open(dir)
	limit := 3000.0
	if !record(t, s, "github.actions", Point{At: at("2026-11-30T23:55:00Z"), Used: 120, Limit: &limit}) {
		t.Fatal("first point")
	}
	if record(t, s, "github.actions", Point{At: at("2026-12-01T00:01:00Z"), Used: 120, Limit: &limit}) {
		t.Fatal("same value six minutes later, in a new month")
	}
	if !record(t, s, "github.actions", Point{At: at("2026-12-01T00:02:00Z"), Used: 0, Limit: &limit}) {
		t.Fatal("changed value in the new month")
	}
	for _, month := range []string{"2026-11.jsonl", "2026-12.jsonl"} {
		if _, err := os.Stat(filepath.Join(dir, "github.actions", month)); err != nil {
			t.Fatal(err)
		}
	}
	h := read(t, s, []string{"github.actions"}, "2026-11-30T00:00:00Z", "2026-12-02T00:00:00Z")
	if got := h.Series[0]; got.Unit != "minutes" || len(got.Points) != 2 || *got.Points[0].Limit != 3000 || got.Points[1].Used != 0 {
		t.Fatalf("%+v", got)
	}
}

func TestConcurrentRecordsKeepWholeLines(t *testing.T) {
	s := Open(t.TempDir())
	var group sync.WaitGroup
	for i := range 10 {
		group.Go(func() {
			if _, err := s.Record("codex.five_hour", Point{At: at("2026-11-28T10:00:00Z").Add(time.Duration(i) * time.Minute), Used: float64(i)}); err != nil {
				t.Error(err)
			}
		})
	}
	group.Wait()
	h := read(t, s, []string{"codex.five_hour"}, "2026-11-28T00:00:00Z", "2026-11-29T00:00:00Z")
	if len(h.Invalid) != 0 || len(h.Series[0].Points) == 0 {
		t.Fatalf("%+v", h)
	}
	for i := 1; i < len(h.Series[0].Points); i++ {
		if !h.Series[0].Points[i].At.After(h.Series[0].Points[i-1].At) {
			t.Fatalf("not sorted %+v", h.Series[0].Points)
		}
	}
}

func TestReadSelectsSortsAndReportsBadLines(t *testing.T) {
	dir := t.TempDir()
	s := Open(dir)
	if err := os.MkdirAll(filepath.Join(dir, "claude.five_hour"), 0o755); err != nil {
		t.Fatal(err)
	}
	lines := `{"at":"2026-11-28T12:00:00Z","used":30}
{"at":"2026-11-28T11:00:00Z","used":20}
not json
{"at":"2026-11-28T12:00:00Z","used":30}
{"at":"2026-11-20T12:00:00Z","used":5}
`
	if err := os.WriteFile(filepath.Join(dir, "claude.five_hour", "2026-11.jsonl"), []byte(lines), 0o644); err != nil {
		t.Fatal(err)
	}
	h := read(t, s, nil, "2026-11-28T00:00:00Z", "2026-11-28T12:00:00Z")
	if len(h.Series) != len(Series) || h.Series[0].ID != "claude.five_hour" || h.Series[4].ID != "github.actions" {
		t.Fatalf("series %+v", h.Series)
	}
	got := h.Series[0].Points
	if len(got) != 2 || got[0].Used != 20 || got[1].Used != 30 {
		t.Fatalf("points %+v", got)
	}
	if len(h.Invalid) != 1 || h.Invalid[0].File != "claude.five_hour/2026-11.jsonl" || h.Invalid[0].Line != 3 {
		t.Fatalf("invalid %+v", h.Invalid)
	}
	for _, other := range h.Series[1:] {
		if other.Points == nil || len(other.Points) != 0 {
			t.Fatalf("empty series %+v", other)
		}
	}
	if _, err := s.Read([]string{"claude.daily"}, at("2026-11-28T00:00:00Z"), at("2026-11-29T00:00:00Z")); !errors.Is(err, ErrUnknownSeries) {
		t.Fatalf("unknown series: %v", err)
	}
	if h, err := Open(filepath.Join(dir, "missing")).Read(nil, at("2026-11-28T00:00:00Z"), at("2026-11-29T00:00:00Z")); err != nil || len(h.Series) != len(Series) {
		t.Fatalf("missing directory %+v %v", h, err)
	}
}

func TestRecordFailsInAReadOnlyDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	if _, err := Open(filepath.Join(dir, "history")).Record("claude.weekly", point("2026-11-28T10:00:00Z", 1)); err == nil {
		t.Fatal("recorded in a read-only directory")
	}
}
