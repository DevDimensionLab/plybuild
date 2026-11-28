package publication

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func shareStore(t *testing.T) (Store, Path, []View) {
	t.Helper()
	root := filepath.Join(t.TempDir(), "publications")
	s := Open(root)
	clock := time.Date(2026, 11, 28, 10, 0, 0, 0, time.UTC)
	s.now = func() time.Time { clock = clock.Add(time.Second); return clock }
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>graf</html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	var runs []View
	for range 2 {
		v, err := s.Publish(Request{Project: "a", Repo: "a-app", Category: "code-health-history", Dir: dir, Entry: "index.html", Title: "graf"})
		if err != nil {
			t.Fatal(err)
		}
		runs = append(runs, v)
	}
	return s, Path{Repo: "a-app", Category: "code-health-history"}, runs
}

func requireClass(t *testing.T, err error, class string) {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) || e.Class != class {
		t.Fatalf("error %v, want %s", err, class)
	}
}

const artifact = "https://claude.ai/artifact/QGwrVaLGvjesvGro1vkR2y"

func TestShareReservesBeforeLaunchAndRefusesASecondAttempt(t *testing.T) {
	s, p, runs := shareStore(t)
	share, started, err := s.ReserveShare("a", p, runs[1].Run, "plyapp-server")
	if err != nil || !started || share.State != ShareStatePending || share.Run != runs[1].Run || share.URL != nil || share.SharedRun != nil ||
		share.StartedAt == "" || share.Actor != "plyapp-server" || share.Kind != ShareKind || share.Repo != "a-app" {
		t.Fatalf("reserve %+v %v %v", share, started, err)
	}
	if _, _, err := s.ReserveShare("a", p, runs[1].Run, "other"); err == nil {
		t.Fatal("second reservation while pending")
	} else {
		requireClass(t, err, ClassSharePending)
	}
	if _, err := s.FinishShareLaunch(p, runs[1].Run, &ShareSession{Provider: "claude", Pane: "p1", Tab: "t1"}, ShareStateUnknown, "submission unknown"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.ReserveShare("a", p, runs[1].Run, "other"); err == nil {
		t.Fatal("reservation while unknown")
	} else {
		requireClass(t, err, ClassSharePending)
	}
}

func TestShareReserveRequiresAnExistingRun(t *testing.T) {
	s, p, _ := shareStore(t)
	_, _, err := s.ReserveShare("a", p, "20991231T000000000Z-00000000", "x")
	requireClass(t, err, ClassNotFound)
	_, _, err = s.ReserveShare("a", Path{Repo: "a-app", Category: "code-map"}, "", "x")
	requireClass(t, err, ClassNotFound)
}

func TestShareResultKeepsOneLinkPerPath(t *testing.T) {
	s, p, runs := shareStore(t)
	if _, _, err := s.ReserveShare("a", p, runs[0].Run, ""); err != nil {
		t.Fatal(err)
	}
	session := &ShareSession{Provider: "claude", Pane: "p1", Tab: "t1", ID: "s1"}
	if share, err := s.FinishShareLaunch(p, runs[0].Run, session, ShareStatePending, ""); err != nil || share.State != ShareStatePending || share.Session == nil || share.Session.ID != "s1" {
		t.Fatalf("launch %+v %v", share, err)
	}
	if _, err := s.RecordShareResult(p, runs[1].Run, artifact, ""); err == nil {
		t.Fatal("result for another run")
	} else {
		requireClass(t, err, ClassShareConflict)
	}
	for _, bad := range []string{"http://claude.ai/artifact/x", "https://claude.ai/artifact/", "https://example.com/artifact/x", artifact + "/x"} {
		_, err := s.RecordShareResult(p, runs[0].Run, bad, "")
		requireClass(t, err, ClassInvalidArguments)
	}
	shared, err := s.RecordShareResult(p, runs[0].Run, artifact, "")
	if err != nil || shared.State != ShareStateShared || *shared.URL != artifact || *shared.SharedRun != runs[0].Run || shared.FinishedAt == nil || shared.Error != nil {
		t.Fatalf("shared %+v %v", shared, err)
	}
	if again, err := s.RecordShareResult(p, runs[0].Run, artifact, ""); err != nil || again.FinishedAt == nil || *again.FinishedAt != *shared.FinishedAt {
		t.Fatalf("replay %+v %v", again, err)
	}
	if same, started, err := s.ReserveShare("a", p, runs[0].Run, ""); err != nil || started || same.State != ShareStateShared {
		t.Fatalf("same run again %+v %v %v", same, started, err)
	}
	next, started, err := s.ReserveShare("a", p, runs[1].Run, "")
	if err != nil || !started || next.State != ShareStatePending || *next.URL != artifact || *next.SharedRun != runs[0].Run || next.Session != nil {
		t.Fatalf("newer run %+v %v %v", next, started, err)
	}
	if _, err := s.RecordShareResult(p, runs[1].Run, "https://claude.ai/artifact/Other", ""); err == nil {
		t.Fatal("a second link for the path")
	} else {
		requireClass(t, err, ClassShareConflict)
	}
	failed, err := s.RecordShareResult(p, runs[1].Run, "", "content check failed")
	if err != nil || failed.State != ShareStateFailed || *failed.Error != "content check failed" || *failed.URL != artifact || *failed.SharedRun != runs[0].Run {
		t.Fatalf("failed %+v %v", failed, err)
	}
	if retry, started, err := s.ReserveShare("a", p, runs[1].Run, ""); err != nil || !started || retry.Error != nil {
		t.Fatalf("retry after failure %+v %v %v", retry, started, err)
	}
}

func TestShareLaunchFailureEndsTheAttempt(t *testing.T) {
	s, p, runs := shareStore(t)
	if _, _, err := s.ReserveShare("a", p, runs[0].Run, ""); err != nil {
		t.Fatal(err)
	}
	failed, err := s.FinishShareLaunch(p, runs[0].Run, nil, ShareStateFailed, "launch failed: no shell")
	if err != nil || failed.State != ShareStateFailed || *failed.Error != "launch failed: no shell" || failed.FinishedAt == nil {
		t.Fatalf("%+v %v", failed, err)
	}
	if _, err := s.FinishShareLaunch(p, runs[0].Run, nil, ShareStatePending, ""); err == nil {
		t.Fatal("launch outcome after the attempt ended")
	} else {
		requireClass(t, err, ClassShareConflict)
	}
	if _, err := s.RecordShareResult(p, runs[0].Run, artifact, ""); err == nil {
		t.Fatal("result after the attempt ended")
	} else {
		requireClass(t, err, ClassShareConflict)
	}
}

func TestShareIsShownOnTheRunsOfItsPath(t *testing.T) {
	s, p, runs := shareStore(t)
	if v, err := s.Show(p.Repo, "", p.Category, "", ""); err != nil || v.Share != nil {
		t.Fatalf("before %+v %v", v.Share, err)
	}
	if _, _, err := s.ReserveShare("a", p, runs[0].Run, ""); err != nil {
		t.Fatal(err)
	}
	v, err := s.Show(p.Repo, "", p.Category, "", "")
	if err != nil || v.Run != runs[1].Run || v.Share == nil || v.Share.Run != runs[0].Run || v.Share.State != ShareStatePending {
		t.Fatalf("show %+v %v", v.Share, err)
	}
	list, err := s.List(Filter{})
	if err != nil || len(list.Items) != 1 || list.Items[0].Share == nil || list.Items[0].Share.Run != runs[0].Run {
		t.Fatalf("list %+v %v", list, err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(runs[0].Path), shareName), []byte("{"), 0o644); err != nil {
		t.Fatal(err)
	}
	list, err = s.List(Filter{})
	if err != nil || len(list.Items) != 1 || list.Items[0].Share != nil || len(list.Invalid) != 1 {
		t.Fatalf("unreadable share %+v %v", list, err)
	}
	if v, err := s.Show(p.Repo, "", p.Category, "", ""); err != nil || v.Share != nil {
		t.Fatalf("show unreadable share %+v %v", v.Share, err)
	}
}
