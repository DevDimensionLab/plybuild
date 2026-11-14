package agentstatus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testStore(t *testing.T) (Store, *time.Time) {
	t.Helper()
	now := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	s := Open(filepath.Join(t.TempDir(), "agent-status"))
	s.now = func() time.Time { return now }
	return s, &now
}

func requireClass(t *testing.T, err error, class string) {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) || e.Class != class {
		t.Fatalf("want %s, got %v", class, err)
	}
}

func TestSetThenListShowsStoredPost(t *testing.T) {
	s, _ := testStore(t)
	item, err := s.Set("claude:aa24", "w7:pCY", "waiting_background", "PLAN-36/arch", 30*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	want := Item{Actor: "claude:aa24", Pane: "w7:pCY", State: "waiting_background", Ref: "PLAN-36/arch",
		ReportedAt: "2026-11-14T12:00:00Z", ExpiresAt: "2026-11-14T12:30:00Z"}
	if item != want {
		t.Fatalf("set readback %+v", item)
	}
	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if list.ObservedAt != "2026-11-14T12:00:00Z" || len(list.Items) != 1 || list.Items[0] != want || len(list.Invalid) != 0 {
		t.Fatalf("list %+v", list)
	}
}

func TestSetReplacesSameActorAndKeepsOthers(t *testing.T) {
	s, _ := testStore(t)
	for _, set := range []struct{ actor, state string }{{"claude:a", "working"}, {"codex:0f", "done"}, {"claude:a", "needs_answer"}} {
		if _, err := s.Set(set.actor, "w1:p1", set.state, "", time.Hour); err != nil {
			t.Fatal(err)
		}
	}
	list, _ := s.List()
	if len(list.Items) != 2 || list.Items[0].Actor != "claude:a" || list.Items[0].State != "needs_answer" || list.Items[1].Actor != "codex:0f" {
		t.Fatalf("list %+v", list.Items)
	}
}

func TestExpiryIsShownAndOldExpiredPostsArePrunedOnSet(t *testing.T) {
	s, now := testStore(t)
	if _, err := s.Set("claude:old", "w1:p1", "done", "", time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Set("claude:recent", "w1:p2", "working", "", 4*time.Hour); err != nil {
		t.Fatal(err)
	}
	*now = now.Add(2 * time.Hour)
	list, _ := s.List()
	if len(list.Items) != 2 || !list.Items[0].Expired || list.Items[1].Expired {
		t.Fatalf("expiry %+v", list.Items)
	}
	*now = now.Add(23 * time.Hour) // old expired over 24 h ago, recent 21 h ago
	if _, err := s.Set("codex:01", "w1:p3", "working", "", time.Minute); err != nil {
		t.Fatal(err)
	}
	list, _ = s.List()
	if len(list.Items) != 2 || list.Items[0].Actor != "claude:recent" || list.Items[1].Actor != "codex:01" {
		t.Fatalf("pruning %+v", list.Items)
	}
}

func TestClearRemovesOwnPostAndEmptyListIsNotAnError(t *testing.T) {
	s, _ := testStore(t)
	if list, err := s.List(); err != nil || list.Items == nil || len(list.Items) != 0 || list.Invalid == nil {
		t.Fatalf("empty list %+v %v", list, err)
	}
	if removed, err := s.Clear("claude:a"); err != nil || removed {
		t.Fatalf("clear absent %v %v", removed, err)
	}
	if _, err := s.Set("claude:a", "w1:p1", "working", "", time.Minute); err != nil {
		t.Fatal(err)
	}
	if removed, err := s.Clear("claude:a"); err != nil || !removed {
		t.Fatalf("clear %v %v", removed, err)
	}
	if list, _ := s.List(); len(list.Items) != 0 {
		t.Fatalf("after clear %+v", list.Items)
	}
}

func TestInvalidInputIsRefusedBeforeAnyEffect(t *testing.T) {
	cases := []struct {
		name, actor, pane, state, ref string
		ttl                           time.Duration
	}{
		{"state", "claude:a", "w1:p1", "sleeping", "", time.Minute},
		{"no actor", "", "w1:p1", "working", "", time.Minute},
		{"free-text actor", "my agent", "w1:p1", "working", "", time.Minute},
		{"no pane", "claude:a", "", "working", "", time.Minute},
		{"pane form", "claude:a", "pane 1", "working", "", time.Minute},
		{"free-text ref", "claude:a", "w1:p1", "working", "waiting for review", time.Minute},
		{"long ref", "claude:a", "w1:p1", "working", fmt.Sprintf("%081d", 0), time.Minute},
		{"no duration", "claude:a", "w1:p1", "working", "", 0},
		{"over 4 h", "claude:a", "w1:p1", "waiting_background", "", 4*time.Hour + time.Second},
		{"needs_answer over 24 h", "claude:a", "w1:p1", "needs_answer", "", 24*time.Hour + time.Second},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, _ := testStore(t)
			_, err := s.Set(c.actor, c.pane, c.state, c.ref, c.ttl)
			requireClass(t, err, ClassInvalidArguments)
			if _, err := os.Lstat(s.dir); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("directory created")
			}
		})
	}
	s, _ := testStore(t)
	if _, err := s.Set("claude:a", "w1:p1", "needs_answer", "", 24*time.Hour); err != nil {
		t.Fatalf("needs_answer 24 h: %v", err)
	}
	if _, err := s.Set("claude:b", "w1:p1", "done", "", 4*time.Hour); err != nil {
		t.Fatalf("done 4 h: %v", err)
	}
}

func TestUnreadableFileIsListedAsInvalid(t *testing.T) {
	s, _ := testStore(t)
	if _, err := s.Set("claude:a", "w1:p1", "working", "", time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(s.dir, "broken.json"), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	list, err := s.List()
	if err != nil || len(list.Items) != 1 || len(list.Invalid) != 1 || list.Invalid[0].File != "broken.json" {
		t.Fatalf("list %+v %v", list, err)
	}
}

func TestListOfHundredPostsIsFast(t *testing.T) {
	s, _ := testStore(t)
	for i := 0; i < 100; i++ {
		if _, err := s.Set(fmt.Sprintf("codex:%02x", i), "w1:p1", "working", "", time.Hour); err != nil {
			t.Fatal(err)
		}
	}
	started := time.Now()
	list, err := s.List()
	elapsed := time.Since(started)
	if err != nil || len(list.Items) != 100 {
		t.Fatalf("list %d %v", len(list.Items), err)
	}
	if elapsed > 100*time.Millisecond {
		t.Fatalf("list took %s", elapsed)
	}
}
