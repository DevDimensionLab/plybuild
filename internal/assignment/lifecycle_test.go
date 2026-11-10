package assignment

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

const (
	lifecycleAt = "2026-11-12T08:14:03Z"
	claim1      = "c_00000000000000000000000000000001"
	claim2      = "c_00000000000000000000000000000002"
	otherClaim  = "c_ffffffffffffffffffffffffffffffff"
	secretActor = "hemmelig agent"
	secretWhy   = "hemmelig grunn"
)

// lifecycleStore has ASSIGNMENT-1..n and deterministic claims c_<n as 32 hex>.
func lifecycleStore(t *testing.T, n int) (Store, string) {
	t.Helper()
	s, dir := testStore(t)
	var mu sync.Mutex
	next := 0
	s.claim = func() (string, error) {
		mu.Lock()
		defer mu.Unlock()
		next++
		return fmt.Sprintf("c_%032x", next), nil
	}
	for i := 1; i <= n; i++ {
		mustAdd(t, s, fmt.Sprintf("T%d", i))
	}
	return s, dir
}

func recordPath(dir, id string) string { return filepath.Join(dir, ".lifecycle", id+".json") }

func fileHash(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(contents))
}

func mustTake(t *testing.T, s Store, id string, revision int, actor string) Lifecycle {
	t.Helper()
	l, err := s.Take(id, revision, actor)
	if err != nil {
		t.Fatalf("take %s: %v", id, err)
	}
	return l
}

func mustLifecycle(t *testing.T, s Store, id string) Lifecycle {
	t.Helper()
	l, err := s.Lifecycle(id)
	if err != nil {
		t.Fatalf("lifecycle %s: %v", id, err)
	}
	return l
}

func viewIDs(t *testing.T, s Store, view View) []string {
	t.Helper()
	result, err := s.ListView(view)
	if err != nil {
		t.Fatalf("list %s: %v", view, err)
	}
	ids := []string{}
	for _, item := range result.Items {
		ids = append(ids, item.ID)
	}
	return ids
}

func requireView(t *testing.T, s Store, view View, want ...string) {
	t.Helper()
	if want == nil {
		want = []string{}
	}
	if got := viewIDs(t, s, view); !reflect.DeepEqual(got, want) {
		t.Fatalf("view %s = %v, want %v", view, got, want)
	}
}

func requireConflict(t *testing.T, err error, detail string) {
	t.Helper()
	if e := requireClass(t, err, ClassConflict); detail != "" && e.Detail != detail {
		t.Fatalf("detail %q, want %q", e.Detail, detail)
	}
}

func readRecordBytes(t *testing.T, dir, id string) string {
	t.Helper()
	contents, err := os.ReadFile(recordPath(dir, id))
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}

func TestNewAssignmentIsAvailableAtRevisionZero(t *testing.T) {
	s, dir := lifecycleStore(t, 3)
	l := mustLifecycle(t, s, "ASSIGNMENT-2")
	want := Lifecycle{State: StateAvailable, Revision: 0, Events: []Event{}}
	if !reflect.DeepEqual(l, want) {
		t.Fatalf("lifecycle %+v", l)
	}
	requireView(t, s, ViewAvailable, "ASSIGNMENT-1", "ASSIGNMENT-2", "ASSIGNMENT-3")
	requireView(t, s, ViewTaken)
	requireView(t, s, ViewArchive)
	requireView(t, s, ViewAll, "ASSIGNMENT-1", "ASSIGNMENT-2", "ASSIGNMENT-3")
	requireAbsent(t, filepath.Join(dir, ".lifecycle"))
	_, err := s.Lifecycle("ASSIGNMENT-4")
	requireClass(t, err, ClassNotFound)
	_, err = s.Lifecycle("ASSIGNMENT-01")
	requireClass(t, err, ClassInvalidArguments)
}

func TestTakeDoneJourneyArchivesWithoutOtherEffects(t *testing.T) {
	s, dir := lifecycleStore(t, 3)
	original := fileHash(t, filepath.Join(dir, "ASSIGNMENT-2.json"))
	taken := mustTake(t, s, "ASSIGNMENT-2", 0, "a")
	wantHolder := &Holder{Actor: "a", Claim: claim1, TakenAt: lifecycleAt}
	if taken.State != StateTaken || taken.Revision != 1 || !reflect.DeepEqual(taken.Holder, wantHolder) || taken.Done != nil ||
		!reflect.DeepEqual(taken.Events, []Event{{Seq: 1, At: lifecycleAt, Kind: EventTaken, Actor: "a", Claim: claim1}}) {
		t.Fatalf("take %+v", taken)
	}
	if shown := mustLifecycle(t, s, "ASSIGNMENT-2"); !reflect.DeepEqual(shown, taken) {
		t.Fatalf("show %+v", shown)
	}
	requireView(t, s, ViewAvailable, "ASSIGNMENT-1", "ASSIGNMENT-3")
	requireView(t, s, ViewTaken, "ASSIGNMENT-2")

	before := snapshotTree(t, filepath.Dir(dir))
	done, err := s.Done("ASSIGNMENT-2", claim1)
	if err != nil {
		t.Fatal(err)
	}
	if done.State != StateDone || done.Revision != 2 || done.Holder != nil || !reflect.DeepEqual(done.Done, &DoneBy{Actor: "a", At: lifecycleAt}) {
		t.Fatalf("done %+v", done)
	}
	requireView(t, s, ViewAvailable, "ASSIGNMENT-1", "ASSIGNMENT-3")
	requireView(t, s, ViewTaken)
	requireView(t, s, ViewArchive, "ASSIGNMENT-2")
	requireView(t, s, ViewAll, "ASSIGNMENT-1", "ASSIGNMENT-2", "ASSIGNMENT-3")
	result, err := s.ListView(ViewArchive)
	if err != nil {
		t.Fatal(err)
	}
	item := result.Items[0]
	if item.State != StateDone || item.Revision != 2 || item.DoneAt != lifecycleAt || item.Actor != "a" || item.TakenAt != "" || item.Title != "T2" || result.View != ViewArchive {
		t.Fatalf("archive item %+v", item)
	}
	changed := diffTrees(before, snapshotTree(t, filepath.Dir(dir)))
	if !reflect.DeepEqual(changed, []string{filepath.Join("assignments", ".lifecycle", "ASSIGNMENT-2.json")}) {
		t.Fatalf("done changed %v", changed)
	}
	if got := fileHash(t, filepath.Join(dir, "ASSIGNMENT-2.json")); got != original {
		t.Fatal("original changed")
	}
}

func TestOriginalNeverChangesAcrossLifecycle(t *testing.T) {
	s, dir := lifecycleStore(t, 2)
	path := filepath.Join(dir, "ASSIGNMENT-2.json")
	original := fileHash(t, path)
	info, _ := os.Stat(path)
	mustTake(t, s, "ASSIGNMENT-2", 0, "a")
	if _, err := s.Release("ASSIGNMENT-2", claim1); err != nil {
		t.Fatal(err)
	}
	mustTake(t, s, "ASSIGNMENT-2", 2, "b")
	if _, err := s.Done("ASSIGNMENT-2", claim2); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(path)
	if fileHash(t, path) != original || !after.ModTime().Equal(info.ModTime()) || !os.SameFile(info, after) {
		t.Fatal("original file was rewritten")
	}
}

func TestReleaseReturnsToAvailableAndMovesRevision(t *testing.T) {
	s, _ := lifecycleStore(t, 2)
	mustTake(t, s, "ASSIGNMENT-2", 0, "a")
	released, err := s.Release("ASSIGNMENT-2", claim1)
	if err != nil {
		t.Fatal(err)
	}
	if released.State != StateAvailable || released.Revision != 2 || released.Holder != nil ||
		!reflect.DeepEqual(released.Events[1], Event{Seq: 2, At: lifecycleAt, Kind: EventReleased, Actor: "a", Claim: claim1}) {
		t.Fatalf("release %+v", released)
	}
	requireView(t, s, ViewAvailable, "ASSIGNMENT-1", "ASSIGNMENT-2")
	_, err = s.Take("ASSIGNMENT-2", 1, "b")
	requireConflict(t, err, "stale revision 1, current 2")
	if l := mustTake(t, s, "ASSIGNMENT-2", 2, "b"); l.Revision != 3 || l.Holder.Claim == claim1 {
		t.Fatalf("retake %+v", l)
	}
}

func TestWrongClaimAndTerminalDoneAreConflictsWithoutWrite(t *testing.T) {
	s, dir := lifecycleStore(t, 1)
	mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	before := readRecordBytes(t, dir, "ASSIGNMENT-1")
	_, err := s.Done("ASSIGNMENT-1", otherClaim)
	requireConflict(t, err, "claim does not hold ASSIGNMENT-1")
	_, err = s.Release("ASSIGNMENT-1", otherClaim)
	requireConflict(t, err, "claim does not hold ASSIGNMENT-1")
	_, err = s.Take("ASSIGNMENT-1", 1, "b")
	requireConflict(t, err, "ASSIGNMENT-1 is taken")
	if readRecordBytes(t, dir, "ASSIGNMENT-1") != before {
		t.Fatal("record changed")
	}

	if _, err := s.Done("ASSIGNMENT-1", claim1); err != nil {
		t.Fatal(err)
	}
	before = readRecordBytes(t, dir, "ASSIGNMENT-1")
	_, err = s.Take("ASSIGNMENT-1", 2, "b")
	requireConflict(t, err, "ASSIGNMENT-1 is done")
	_, err = s.Release("ASSIGNMENT-1", claim1)
	requireConflict(t, err, "ASSIGNMENT-1 is done")
	_, err = s.Revoke("ASSIGNMENT-1", 2, "human", "gone")
	requireConflict(t, err, "ASSIGNMENT-1 is done")
	_, err = s.Done("ASSIGNMENT-1", otherClaim)
	requireConflict(t, err, "ASSIGNMENT-1 is done")
	if readRecordBytes(t, dir, "ASSIGNMENT-1") != before {
		t.Fatal("record changed after done")
	}

	s2, _ := lifecycleStore(t, 1)
	_, err = s2.Release("ASSIGNMENT-1", claim1)
	requireConflict(t, err, "claim does not hold ASSIGNMENT-1")
	_, err = s2.Done("ASSIGNMENT-1", claim1)
	requireConflict(t, err, "claim does not hold ASSIGNMENT-1")
	_, err = s2.Revoke("ASSIGNMENT-1", 0, "human", "gone")
	requireConflict(t, err, "ASSIGNMENT-1 is not taken")
}

func TestConcurrentTakesHaveExactlyOneWinner(t *testing.T) {
	s, dir := lifecycleStore(t, 1)
	const takers = 16
	var wg sync.WaitGroup
	results := make(chan error, takers)
	for i := 0; i < takers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := s.Take("ASSIGNMENT-1", 0, fmt.Sprintf("agent-%d", i))
			results <- err
		}(i)
	}
	wg.Wait()
	close(results)
	won := 0
	for err := range results {
		if err == nil {
			won++
			continue
		}
		requireConflict(t, err, "")
	}
	if won != 1 {
		t.Fatalf("%d winners", won)
	}
	var rec struct{ Events []json.RawMessage }
	if err := json.Unmarshal([]byte(readRecordBytes(t, dir, "ASSIGNMENT-1")), &rec); err != nil || len(rec.Events) != 1 {
		t.Fatalf("record events %d %v", len(rec.Events), err)
	}
}

func TestAddDuringTakesKeepsNumberingAndRecords(t *testing.T) {
	s, _ := lifecycleStore(t, 1)
	var wg sync.WaitGroup
	takes := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			if _, err := s.Add(fmt.Sprintf("N%d", i), "x", []string{"ply"}, nil); err != nil {
				t.Error(err)
			}
		}(i)
		go func(i int) {
			defer wg.Done()
			_, err := s.Take("ASSIGNMENT-1", 0, fmt.Sprintf("agent-%d", i))
			takes <- err
		}(i)
	}
	wg.Wait()
	close(takes)
	won := 0
	for err := range takes {
		if err == nil {
			won++
		}
	}
	if won != 1 {
		t.Fatalf("%d winners", won)
	}
	result, err := s.ListView(ViewAll)
	if err != nil || len(result.Invalid) != 0 || len(result.Items) != 9 {
		t.Fatalf("list %+v %v", result, err)
	}
	for i, item := range result.Items {
		if item.ID != fmt.Sprintf("ASSIGNMENT-%d", i+1) {
			t.Fatalf("item %d is %s", i, item.ID)
		}
	}
	requireView(t, s, ViewTaken, "ASSIGNMENT-1")
}

func TestRevokeRecoversAndInvalidatesOldClaim(t *testing.T) {
	s, _ := lifecycleStore(t, 1)
	mustTake(t, s, "ASSIGNMENT-1", 0, "crashed")
	_, err := s.Revoke("ASSIGNMENT-1", 0, "human", "session lost")
	requireConflict(t, err, "stale revision 0, current 1")
	revoked, err := s.Revoke("ASSIGNMENT-1", 1, "human", "session lost")
	if err != nil {
		t.Fatal(err)
	}
	if revoked.State != StateAvailable || revoked.Revision != 2 ||
		!reflect.DeepEqual(revoked.Events[1], Event{Seq: 2, At: lifecycleAt, Kind: EventRevoked, Actor: "human", Reason: "session lost"}) {
		t.Fatalf("revoke %+v", revoked)
	}
	_, err = s.Done("ASSIGNMENT-1", claim1)
	requireConflict(t, err, "claim does not hold ASSIGNMENT-1")
	mustTake(t, s, "ASSIGNMENT-1", 2, "fresh")
	_, err = s.Done("ASSIGNMENT-1", claim1)
	requireConflict(t, err, "claim does not hold ASSIGNMENT-1")
	_, err = s.Release("ASSIGNMENT-1", claim1)
	requireConflict(t, err, "claim does not hold ASSIGNMENT-1")
	if l := mustLifecycle(t, s, "ASSIGNMENT-1"); l.Holder == nil || l.Holder.Actor != "fresh" || l.Holder.Claim != claim2 {
		t.Fatalf("holder %+v", l.Holder)
	}
}

func TestRetriedDoneAndReleaseAreIdempotent(t *testing.T) {
	s, dir := lifecycleStore(t, 2)
	mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	first, err := s.Done("ASSIGNMENT-1", claim1)
	if err != nil {
		t.Fatal(err)
	}
	before := readRecordBytes(t, dir, "ASSIGNMENT-1")
	again, err := s.Done("ASSIGNMENT-1", claim1)
	if err != nil || !reflect.DeepEqual(again, first) || readRecordBytes(t, dir, "ASSIGNMENT-1") != before {
		t.Fatalf("repeat done %+v %v", again, err)
	}

	mustTake(t, s, "ASSIGNMENT-2", 0, "b")
	released, err := s.Release("ASSIGNMENT-2", claim2)
	if err != nil {
		t.Fatal(err)
	}
	before = readRecordBytes(t, dir, "ASSIGNMENT-2")
	again, err = s.Release("ASSIGNMENT-2", claim2)
	if err != nil || !reflect.DeepEqual(again, released) || readRecordBytes(t, dir, "ASSIGNMENT-2") != before {
		t.Fatalf("repeat release %+v %v", again, err)
	}
	mustTake(t, s, "ASSIGNMENT-2", 2, "c")
	_, err = s.Release("ASSIGNMENT-2", claim2)
	requireConflict(t, err, "claim does not hold ASSIGNMENT-2")
}

func TestTakeRefusedAtEventLimit(t *testing.T) {
	s, dir := lifecycleStore(t, 1)
	revision := 0
	for revision+2 <= maxEvents {
		l := mustTake(t, s, "ASSIGNMENT-1", revision, "a")
		l, err := s.Release("ASSIGNMENT-1", l.Holder.Claim)
		if err != nil {
			t.Fatal(err)
		}
		revision = l.Revision
	}
	if revision != maxEvents {
		t.Fatalf("revision %d", revision)
	}
	_, err := s.Take("ASSIGNMENT-1", revision, "a")
	requireConflict(t, err, "lifecycle event limit reached")

	writeRaw(t, filepath.Join(dir, ".lifecycle"), "ASSIGNMENT-1.json", recordJSON("ASSIGNMENT-1", cycleEvents(31, true)...))
	if l := mustLifecycle(t, s, "ASSIGNMENT-1"); l.Revision != 63 || l.State != StateTaken {
		t.Fatalf("63 %+v", l)
	}
	writeRaw(t, filepath.Join(dir, ".lifecycle"), "ASSIGNMENT-1.json", recordJSON("ASSIGNMENT-1", cycleEvents(31, false)[:62]...))
	if l := mustLifecycle(t, s, "ASSIGNMENT-1"); l.Revision != 62 {
		t.Fatalf("62 %+v", l)
	}
	l := mustTake(t, s, "ASSIGNMENT-1", 62, "a")
	if l.Revision != 63 {
		t.Fatalf("take at 62 %+v", l)
	}
	if l, err := s.Done("ASSIGNMENT-1", l.Holder.Claim); err != nil || l.Revision != 64 {
		t.Fatalf("done at 63 %+v %v", l, err)
	}
}

func TestTransitionsRejectInvalidArgumentsWithoutEffect(t *testing.T) {
	s, dir := lifecycleStore(t, 1)
	for name, call := range map[string]func() error{
		"take bad id":          func() error { _, err := s.Take("ASSIGNMENT-01", 0, "a"); return err },
		"take negative":        func() error { _, err := s.Take("ASSIGNMENT-1", -1, "a"); return err },
		"take blank actor":     func() error { _, err := s.Take("ASSIGNMENT-1", 0, " "); return err },
		"take two-line actor":  func() error { _, err := s.Take("ASSIGNMENT-1", 0, "a\nb"); return err },
		"take control actor":   func() error { _, err := s.Take("ASSIGNMENT-1", 0, "a\tb"); return err },
		"take long actor":      func() error { _, err := s.Take("ASSIGNMENT-1", 0, strings.Repeat("ø", 101)); return err },
		"release bad claim":    func() error { _, err := s.Release("ASSIGNMENT-1", "c_ABC"); return err },
		"done bad claim":       func() error { _, err := s.Done("ASSIGNMENT-1", claim1+"0"); return err },
		"revoke blank reason":  func() error { _, err := s.Revoke("ASSIGNMENT-1", 0, "h", " "); return err },
		"revoke long reason":   func() error { _, err := s.Revoke("ASSIGNMENT-1", 0, "h", strings.Repeat("r", 201)); return err },
		"revoke bad actor":     func() error { _, err := s.Revoke("ASSIGNMENT-1", 0, "", "r"); return err },
		"revoke negative":      func() error { _, err := s.Revoke("ASSIGNMENT-1", -1, "h", "r"); return err },
		"revoke two-line why":  func() error { _, err := s.Revoke("ASSIGNMENT-1", 0, "h", "a\nb"); return err },
		"relative store take":  func() error { _, err := Open("assignments").Take("ASSIGNMENT-1", 0, "a"); return err },
		"relative store done":  func() error { _, err := Open("assignments").Done("ASSIGNMENT-1", claim1); return err },
		"invalid view":         func() error { _, err := s.ListView(View("open")); return err },
		"invalid lifecycle id": func() error { _, err := s.Lifecycle("x"); return err },
	} {
		t.Run(name, func(t *testing.T) {
			requireClass(t, call(), ClassInvalidArguments)
			requireAbsent(t, filepath.Join(dir, ".lifecycle"))
		})
	}
	if l := mustTake(t, s, "ASSIGNMENT-1", 0, strings.Repeat("ø", 100)); l.Revision != 1 {
		t.Fatalf("100-rune actor %+v", l)
	}
	_, err := s.Take("ASSIGNMENT-2", 0, "a")
	requireClass(t, err, ClassNotFound)
	empty, _ := testStore(t)
	_, err = empty.Take("ASSIGNMENT-1", 0, "a")
	requireClass(t, err, ClassNotFound)
	requireAbsent(t, empty.dir)
}

func TestReplaceFailureOutcomes(t *testing.T) {
	t.Run("before rename has no effect", func(t *testing.T) {
		s, dir := lifecycleStore(t, 1)
		failing := s
		failing.replace = func(string, []byte, fs.FileMode) error {
			return &atomicfile.Error{Op: "sync", Path: dir, Err: &fs.PathError{Op: "sync", Path: dir, Err: fs.ErrPermission}}
		}
		_, err := failing.Take("ASSIGNMENT-1", 0, secretActor)
		if e := requireClass(t, err, ClassIOError); strings.Contains(e.Detail, "result unknown") {
			t.Fatalf("detail %q", e.Detail)
		}
		requireNoPath(t, s, err, secretActor, claim1)
		if l := mustLifecycle(t, s, "ASSIGNMENT-1"); l.State != StateAvailable || l.Revision != 0 {
			t.Fatalf("after failed take %+v", l)
		}
	})
	t.Run("rename error is unknown", func(t *testing.T) {
		s, dir := lifecycleStore(t, 1)
		failing := s
		failing.replace = func(string, []byte, fs.FileMode) error {
			return &atomicfile.Error{Op: "replace", Path: dir, Err: &os.LinkError{Op: "rename", Old: dir, New: dir, Err: fs.ErrPermission}}
		}
		_, err := failing.Take("ASSIGNMENT-1", 0, "a")
		if e := requireClass(t, err, ClassIOError); !strings.HasSuffix(e.Detail, "; result unknown, run show") {
			t.Fatalf("detail %q", e.Detail)
		}
		requireNoPath(t, s, err)
	})
	t.Run("directory sync after rename is unknown and visible", func(t *testing.T) {
		s, dir := lifecycleStore(t, 1)
		mustTake(t, s, "ASSIGNMENT-1", 0, "a")
		failing := s
		failing.replace = func(path string, contents []byte, mode fs.FileMode) error {
			if err := atomicfile.Replace(path, contents, mode); err != nil {
				t.Fatal(err)
			}
			return &atomicfile.Error{Op: "sync directory", Path: dir, Replaced: true, Err: &fs.PathError{Op: "sync", Path: dir, Err: fs.ErrInvalid}}
		}
		_, err := failing.Done("ASSIGNMENT-1", claim1)
		if e := requireClass(t, err, ClassIOError); !strings.HasSuffix(e.Detail, "; result unknown, run show") {
			t.Fatalf("detail %q", e.Detail)
		}
		requireNoPath(t, s, err, claim1)
		if l := mustLifecycle(t, s, "ASSIGNMENT-1"); l.State != StateDone || l.Revision != 2 {
			t.Fatalf("show after uncertain done %+v", l)
		}
	})
	t.Run("lock failure", func(t *testing.T) {
		s, dir := lifecycleStore(t, 1)
		s.lock = filelock.Options{Fail: func(op string) error {
			if op == "lock" {
				return errors.New("injected")
			}
			return nil
		}}
		_, err := s.Take("ASSIGNMENT-1", 0, "a")
		if e := requireClass(t, err, ClassIOError); e.Detail != "lock assignments" {
			t.Fatalf("detail %q", e.Detail)
		}
		requireAbsent(t, filepath.Join(dir, ".lifecycle"))
	})
	t.Run("claim generator failure", func(t *testing.T) {
		s, dir := lifecycleStore(t, 1)
		s.claim = func() (string, error) { return "", errors.New("no entropy") }
		_, err := s.Take("ASSIGNMENT-1", 0, "a")
		requireClass(t, err, ClassIOError)
		requireAbsent(t, filepath.Join(dir, ".lifecycle"))
	})
}

func TestDefaultClaimsAreRandomHex(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "assignments")
	s := Open(dir)
	mustAdd(t, s, "A")
	mustAdd(t, s, "B")
	a := mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	b := mustTake(t, s, "ASSIGNMENT-2", 0, "b")
	if !claimPattern.MatchString(a.Holder.Claim) || !claimPattern.MatchString(b.Holder.Claim) || a.Holder.Claim == b.Holder.Claim {
		t.Fatalf("claims %q %q", a.Holder.Claim, b.Holder.Claim)
	}
	info, err := os.Lstat(filepath.Join(dir, ".lifecycle"))
	if err != nil || !info.IsDir() || info.Mode().Perm() != 0o755 {
		t.Fatalf(".lifecycle %v %v", info, err)
	}
}

// TestHelperHoldLock is not a test: it is the child process of
// TestKilledLockHolderDoesNotBlockTake.
func TestHelperHoldLock(t *testing.T) {
	path := os.Getenv("PLY_ASSIGNMENT_HOLD_LOCK")
	if path == "" {
		t.Skip("helper process only")
	}
	if _, err := filelock.Acquire(path, filelock.Options{}); err != nil {
		fmt.Println("error", err)
		os.Exit(2)
	}
	fmt.Println("locked")
	time.Sleep(time.Minute)
	os.Exit(3)
}

func TestKilledLockHolderDoesNotBlockTake(t *testing.T) {
	s, dir := lifecycleStore(t, 1)
	child := exec.Command(os.Args[0], "-test.run=^TestHelperHoldLock$")
	child.Env = append(os.Environ(), "PLY_ASSIGNMENT_HOLD_LOCK="+filepath.Join(dir, lockName))
	stdout, err := child.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := child.Start(); err != nil {
		t.Fatal(err)
	}
	line, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil || line != "locked\n" {
		_ = child.Process.Kill()
		t.Fatalf("child said %q %v", line, err)
	}
	result := make(chan error, 1)
	go func() {
		_, err := s.Take("ASSIGNMENT-1", 0, "a")
		result <- err
	}()
	select {
	case err := <-result:
		_ = child.Process.Kill()
		t.Fatalf("take finished while the lock was held: %v", err)
	case <-time.After(300 * time.Millisecond):
	}
	if err := child.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = child.Wait()
	select {
	case err := <-result:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("take still blocked after the lock holder was killed")
	}
	if l := mustLifecycle(t, s, "ASSIGNMENT-1"); l.State != StateTaken {
		t.Fatalf("after take %+v", l)
	}
}

func event(seq int, kind, actor, extra string) string {
	return fmt.Sprintf(`{"seq":%d,"at":%q,"kind":%q,"actor":%q%s}`, seq, lifecycleAt, kind, actor, extra)
}

func withClaim(claim string) string { return fmt.Sprintf(`,"claim":%q`, claim) }

func recordJSON(id string, events ...string) string {
	return fmt.Sprintf(`{"format":"ply-assignment-lifecycle/v1","id":%q,"events":[%s]}`, id, strings.Join(events, ","))
}

// cycleEvents is n take/release cycles, optionally followed by an open take.
func cycleEvents(n int, open bool) []string {
	var events []string
	for i := 0; i < n; i++ {
		claim := fmt.Sprintf("c_%032x", 1000+i)
		events = append(events, event(len(events)+1, "taken", "a", withClaim(claim)), event(len(events)+2, "released", "a", withClaim(claim)))
	}
	if open {
		events = append(events, event(len(events)+1, "taken", "a", withClaim(claim1)))
	}
	return events
}

func TestLifecycleReaderAcceptsValidHistory(t *testing.T) {
	s, dir := lifecycleStore(t, 1)
	writeRaw(t, filepath.Join(dir, ".lifecycle"), "ASSIGNMENT-1.json", recordJSON("ASSIGNMENT-1",
		event(1, "taken", "a", withClaim(claim1)),
		event(2, "revoked", "human", `,"reason":"gone"`),
		event(3, "taken", "b b", withClaim(claim2)),
		event(4, "done", "b b", withClaim(claim2)),
	))
	l := mustLifecycle(t, s, "ASSIGNMENT-1")
	if l.State != StateDone || l.Revision != 4 || !reflect.DeepEqual(l.Done, &DoneBy{Actor: "b b", At: lifecycleAt}) {
		t.Fatalf("lifecycle %+v", l)
	}
}

func TestInvalidLifecycleFailsClosed(t *testing.T) {
	valid := func(events ...string) string { return recordJSON("ASSIGNMENT-2", events...) }
	take := event(1, "taken", "a", withClaim(claim1))
	type setup func(t *testing.T, dir string)
	put := func(contents string) setup {
		return func(t *testing.T, dir string) {
			writeRaw(t, filepath.Join(dir, ".lifecycle"), "ASSIGNMENT-2.json", contents)
		}
	}
	cases := map[string]struct {
		setup  setup
		file   string
		detail string
		hidden bool
	}{
		"bad JSON":            {put(`{"format":` + secretActor), ".lifecycle/ASSIGNMENT-2.json", "invalid JSON", true},
		"unknown key":         {put(strings.Replace(valid(take), `"id"`, `"task":"x","id"`, 1)), ".lifecycle/ASSIGNMENT-2.json", "unknown field", true},
		"unknown event key":   {put(valid(strings.Replace(take, `"seq"`, `"branch":"x","seq"`, 1))), ".lifecycle/ASSIGNMENT-2.json", "unknown field", true},
		"case-folded key":     {put(strings.Replace(valid(take), `"id"`, `"ID"`, 1)), ".lifecycle/ASSIGNMENT-2.json", "unknown field", true},
		"duplicate key":       {put(strings.Replace(valid(take), `"id":"ASSIGNMENT-2"`, `"id":"ASSIGNMENT-2","id":"ASSIGNMENT-2"`, 1)), ".lifecycle/ASSIGNMENT-2.json", "duplicate field", true},
		"duplicate event key": {put(valid(strings.Replace(take, `"actor":"a"`, `"actor":"a","actor":"b"`, 1))), ".lifecycle/ASSIGNMENT-2.json", "duplicate field", true},
		"trailing data":       {put(valid(take) + "{}"), ".lifecycle/ASSIGNMENT-2.json", "data after JSON value", true},
		"wrong format":        {put(strings.Replace(valid(take), "/v1", "/v2", 1)), ".lifecycle/ASSIGNMENT-2.json", "invalid format", true},
		"wrong id":            {put(recordJSON("ASSIGNMENT-1", take)), ".lifecycle/ASSIGNMENT-2.json", "id does not match file name", true},
		"no events":           {put(valid()), ".lifecycle/ASSIGNMENT-2.json", "invalid events", true},
		"null events":         {put(strings.Replace(valid(take), "["+take+"]", "null", 1)), ".lifecycle/ASSIGNMENT-2.json", "invalid events", true},
		"released first":      {put(valid(event(1, "released", "a", withClaim(claim1)))), ".lifecycle/ASSIGNMENT-2.json", "invalid replay", true},
		"bad seq":             {put(valid(event(2, "taken", "a", withClaim(claim1)))), ".lifecycle/ASSIGNMENT-2.json", "invalid replay", true},
		"other actor done":    {put(valid(take, event(2, "done", "b", withClaim(claim1)))), ".lifecycle/ASSIGNMENT-2.json", "invalid replay", true},
		"other claim release": {put(valid(take, event(2, "released", "a", withClaim(claim2)))), ".lifecycle/ASSIGNMENT-2.json", "invalid replay", true},
		"taken twice":         {put(valid(take, event(2, "taken", "b", withClaim(claim2)))), ".lifecycle/ASSIGNMENT-2.json", "invalid replay", true},
		"reused claim":        {put(valid(take, event(2, "released", "a", withClaim(claim1)), event(3, "taken", "a", withClaim(claim1)))), ".lifecycle/ASSIGNMENT-2.json", "invalid replay", true},
		"event after done":    {put(valid(take, event(2, "done", "a", withClaim(claim1)), event(3, "taken", "b", withClaim(claim2)))), ".lifecycle/ASSIGNMENT-2.json", "invalid replay", true},
		"revoke with claim":   {put(valid(take, event(2, "revoked", "h", withClaim(claim1)+`,"reason":"r"`))), ".lifecycle/ASSIGNMENT-2.json", "invalid event", true},
		"revoke without why":  {put(valid(take, event(2, "revoked", "h", ""))), ".lifecycle/ASSIGNMENT-2.json", "invalid event", true},
		"take without claim":  {put(valid(event(1, "taken", "a", ""))), ".lifecycle/ASSIGNMENT-2.json", "invalid event", true},
		"reason on take":      {put(valid(event(1, "taken", "a", withClaim(claim1)+`,"reason":"r"`))), ".lifecycle/ASSIGNMENT-2.json", "invalid event", true},
		"unknown kind":        {put(valid(event(1, "archived", "a", withClaim(claim1)))), ".lifecycle/ASSIGNMENT-2.json", "invalid event", true},
		"bad claim":           {put(valid(event(1, "taken", "a", withClaim("c_X")))), ".lifecycle/ASSIGNMENT-2.json", "invalid event", true},
		"bad actor":           {put(valid(event(1, "taken", "a\nb", withClaim(claim1)))), ".lifecycle/ASSIGNMENT-2.json", "invalid event", true},
		"local time":          {put(valid(strings.Replace(take, "08:14:03Z", "09:14:03+01:00", 1))), ".lifecycle/ASSIGNMENT-2.json", "invalid event", true},
		"float seq":           {put(valid(strings.Replace(take, `"seq":1`, `"seq":1.0`, 1))), ".lifecycle/ASSIGNMENT-2.json", "invalid JSON", true},
		"time goes backwards": {put(valid(take, strings.Replace(event(2, "done", "a", withClaim(claim1)), "08:14:03Z", "08:14:02Z", 1))), ".lifecycle/ASSIGNMENT-2.json", "invalid replay", true},
		"too many events":     {put(valid(append(cycleEvents(32, false), event(65, "taken", "a", withClaim(claim1)))...)), ".lifecycle/ASSIGNMENT-2.json", "invalid events", true},
		"record is directory": {func(t *testing.T, dir string) {
			if err := os.MkdirAll(recordPath(dir, "ASSIGNMENT-2"), 0o755); err != nil {
				t.Fatal(err)
			}
		}, ".lifecycle/ASSIGNMENT-2.json", "not a regular file", true},
		"symlink": {func(t *testing.T, dir string) {
			writeRaw(t, dir, "target", valid(take))
			if err := os.MkdirAll(filepath.Join(dir, ".lifecycle"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(filepath.Join(dir, "target"), recordPath(dir, "ASSIGNMENT-2")); err != nil {
				t.Fatal(err)
			}
		}, ".lifecycle/ASSIGNMENT-2.json", "not a regular file", true},
		"oversized": {put(strings.Replace(valid(take), `"id"`, strings.Repeat(" ", 1<<20)+`"id"`, 1)), ".lifecycle/ASSIGNMENT-2.json", "file too large", true},
		"orphan record": {func(t *testing.T, dir string) {
			writeRaw(t, filepath.Join(dir, ".lifecycle"), "ASSIGNMENT-9.json", recordJSON("ASSIGNMENT-9", take))
		}, ".lifecycle/ASSIGNMENT-9.json", "no valid assignment", false},
		"foreign file": {func(t *testing.T, dir string) {
			writeRaw(t, filepath.Join(dir, ".lifecycle"), "notes.txt", "x")
		}, ".lifecycle/notes.txt", "unexpected entry", false},
		"lifecycle is file": {func(t *testing.T, dir string) {
			writeRaw(t, dir, ".lifecycle", "x")
		}, ".lifecycle", "not a directory", true},
		"lifecycle is symlink": {func(t *testing.T, dir string) {
			real := filepath.Join(filepath.Dir(dir), "real")
			if err := os.MkdirAll(real, 0o755); err != nil {
				t.Fatal(err)
			}
			writeRaw(t, real, "ASSIGNMENT-2.json", valid(take))
			if err := os.Symlink(real, filepath.Join(dir, ".lifecycle")); err != nil {
				t.Fatal(err)
			}
		}, ".lifecycle", "not a directory", true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			s, dir := lifecycleStore(t, 3)
			tc.setup(t, dir)
			before := snapshotTree(t, filepath.Dir(dir))
			result, err := s.ListView(ViewAll)
			if err != nil {
				t.Fatal(err)
			}
			want := []InvalidFile{{File: tc.file, Error: ClassLifecycleInvalid + ": " + tc.detail}}
			if !reflect.DeepEqual(result.Invalid, want) {
				t.Fatalf("invalid %+v, want %+v", result.Invalid, want)
			}
			if !tc.hidden {
				requireView(t, s, ViewAll, "ASSIGNMENT-1", "ASSIGNMENT-2", "ASSIGNMENT-3")
				return
			}
			for _, view := range []View{ViewAvailable, ViewTaken, ViewArchive, ViewAll} {
				if ids := viewIDs(t, s, view); slices.Contains(ids, "ASSIGNMENT-2") {
					t.Fatalf("view %s = %v shows ASSIGNMENT-2", view, ids)
				}
			}
			calls := map[string]func() error{
				"show":    func() error { _, err := s.Lifecycle("ASSIGNMENT-2"); return err },
				"take":    func() error { _, err := s.Take("ASSIGNMENT-2", 0, "a"); return err },
				"take 1":  func() error { _, err := s.Take("ASSIGNMENT-2", 1, "a"); return err },
				"release": func() error { _, err := s.Release("ASSIGNMENT-2", claim1); return err },
				"revoke":  func() error { _, err := s.Revoke("ASSIGNMENT-2", 1, "h", "r"); return err },
				"done":    func() error { _, err := s.Done("ASSIGNMENT-2", claim1); return err },
			}
			for call, run := range calls {
				err := run()
				if e := requireClass(t, err, ClassLifecycleInvalid); e.Detail != tc.detail {
					t.Fatalf("%s detail %q", call, e.Detail)
				}
				requireNoPath(t, s, err, secretActor)
			}
			if changed := diffTrees(before, snapshotTree(t, filepath.Dir(dir))); len(changed) != 0 {
				t.Fatalf("changed %v", changed)
			}
		})
	}
}

func TestInvalidOriginalHidesItsRecordAndIsReportedOnce(t *testing.T) {
	s, dir := lifecycleStore(t, 2)
	mustTake(t, s, "ASSIGNMENT-2", 0, "a")
	writeRaw(t, dir, "ASSIGNMENT-2.json", "{broken")
	result, err := s.ListView(ViewAll)
	if err != nil {
		t.Fatal(err)
	}
	want := []InvalidFile{
		{File: ".lifecycle/ASSIGNMENT-2.json", Error: ClassLifecycleInvalid + ": no valid assignment"},
		{File: "ASSIGNMENT-2.json", Error: "invalid JSON"},
	}
	if !reflect.DeepEqual(result.Invalid, want) || len(result.Items) != 1 {
		t.Fatalf("list %+v", result)
	}
	_, err = s.Release("ASSIGNMENT-2", claim1)
	requireClass(t, err, ClassInvalid)
}

func TestOldReaderSeesAllAssignments(t *testing.T) {
	s, _ := lifecycleStore(t, 3)
	mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	mustTake(t, s, "ASSIGNMENT-2", 0, "b")
	if _, err := s.Done("ASSIGNMENT-2", claim2); err != nil {
		t.Fatal(err)
	}
	result, err := s.List()
	if err != nil || len(result.Items) != 3 || len(result.Invalid) != 0 {
		t.Fatalf("v1 list %+v %v", result, err)
	}
	if got := mustAdd(t, s, "N"); got.ID != "ASSIGNMENT-4" {
		t.Fatalf("v1 add numbered %s", got.ID)
	}
}

func TestViewItemsCarryStateWithoutClaims(t *testing.T) {
	s, _ := lifecycleStore(t, 3)
	mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	mustTake(t, s, "ASSIGNMENT-2", 0, "b")
	if _, err := s.Done("ASSIGNMENT-2", claim2); err != nil {
		t.Fatal(err)
	}
	result, err := s.ListView(ViewAll)
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(result)
	if strings.Contains(string(encoded), "c_0") || strings.Contains(string(encoded), `"claim"`) {
		t.Fatalf("claim in list: %s", encoded)
	}
	var shape struct {
		View  string
		Items []map[string]any
	}
	_ = json.Unmarshal(encoded, &shape)
	keys := func(m map[string]any) []string {
		var out []string
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}
	want := [][]string{
		{"actor", "created_at", "id", "repos", "revision", "state", "taken_at", "title"},
		{"actor", "created_at", "done_at", "id", "repos", "revision", "state", "title"},
		{"created_at", "id", "repos", "revision", "state", "title"},
	}
	for i, item := range shape.Items {
		if !reflect.DeepEqual(keys(item), want[i]) {
			t.Fatalf("item %d keys %v", i, keys(item))
		}
	}
	if shape.View != "all" {
		t.Fatalf("view %q", shape.View)
	}
}

func TestErrorDetailsNeverCarryActorReasonClaimOrPath(t *testing.T) {
	s, _ := lifecycleStore(t, 2)
	mustTake(t, s, "ASSIGNMENT-1", 0, secretActor)
	var errs []error
	collect := func(_ Lifecycle, err error) { errs = append(errs, err) }
	collect(s.Take("ASSIGNMENT-1", 1, secretActor))
	collect(s.Take("ASSIGNMENT-1", 0, secretActor+"\n"))
	collect(s.Done("ASSIGNMENT-1", otherClaim))
	collect(s.Release("ASSIGNMENT-1", otherClaim))
	collect(s.Release("ASSIGNMENT-1", "c_"+secretActor))
	collect(s.Revoke("ASSIGNMENT-1", 0, secretActor, secretWhy))
	collect(s.Revoke("ASSIGNMENT-1", 1, secretActor, secretWhy+"\n"))
	collect(s.Revoke("ASSIGNMENT-2", 0, secretActor, secretWhy))
	collect(s.Take("ASSIGNMENT-7", 0, secretActor))
	for _, err := range errs {
		if err == nil {
			t.Fatal("expected error")
		}
		requireNoPath(t, s, err, secretActor, secretWhy, otherClaim, claim1, "hemmelig")
	}
}

func TestLifecycleIsolation(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps", "-f", "{{.ImportPath}}", ".").Output()
	if err != nil {
		t.Fatalf("go list: %v", err)
	}
	const module = "github.com/devdimensionlab/plybuild/"
	allowed := map[string]bool{module + "internal/assignment": true, module + "internal/atomicfile": true, module + "internal/filelock": true}
	for _, dep := range strings.Fields(string(out)) {
		if strings.HasPrefix(dep, module) && !allowed[dep] {
			t.Errorf("internal/assignment depends on %s", dep)
		}
	}
	if !reflect.DeepEqual(recordFields, map[string]bool{"format": true, "id": true, "events": true}) {
		t.Fatalf("record keys %v", recordFields)
	}
	wantEvent := map[EventKind][]string{
		EventTaken:    {"actor", "at", "claim", "kind", "seq"},
		EventReleased: {"actor", "at", "claim", "kind", "seq"},
		EventDone:     {"actor", "at", "claim", "kind", "seq"},
		EventRevoked:  {"actor", "at", "kind", "reason", "seq"},
	}
	for kind, fields := range eventFields {
		var got []string
		for field := range fields {
			got = append(got, field)
		}
		sort.Strings(got)
		if !reflect.DeepEqual(got, wantEvent[kind]) {
			t.Fatalf("event %s keys %v", kind, got)
		}
	}
	if len(eventFields) != len(wantEvent) {
		t.Fatalf("event kinds %v", eventFields)
	}

	root := t.TempDir()
	writeRaw(t, root, "sibling.txt", "x")
	dir := filepath.Join(root, "project", "assignments")
	s := Open(dir)
	s.now = func() time.Time { return fixedNow }
	mustAdd(t, s, "T")
	before := snapshotTree(t, root)
	l := mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	if _, err := s.Done("ASSIGNMENT-1", l.Holder.Claim); err != nil {
		t.Fatal(err)
	}
	for _, path := range diffTrees(before, snapshotTree(t, root)) {
		if !strings.HasPrefix(path, filepath.Join("project", "assignments")+string(filepath.Separator)) {
			t.Fatalf("changed outside assignments: %s", path)
		}
	}
}

// snapshotTree maps every path under root to its type, size, mode and modification time.
func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		value := fmt.Sprintf("%v %d %d", info.Mode(), info.Size(), info.ModTime().UnixNano())
		if info.Mode().IsRegular() {
			contents, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			value += fmt.Sprintf(" %x", sha256.Sum256(contents))
		}
		if !entry.IsDir() {
			out[rel] = value
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// diffTrees lists non-directory paths created, removed or changed, ignoring the lock file.
func diffTrees(before, after map[string]string) []string {
	changed := []string{}
	for path, value := range after {
		if before[path] != value && filepath.Base(path) != lockName {
			changed = append(changed, path)
		}
	}
	for path := range before {
		if _, ok := after[path]; !ok {
			changed = append(changed, path)
		}
	}
	sort.Strings(changed)
	return changed
}

func TestUnreadableRecordIsReportedAndHidesOnlyItsAssignment(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root reads files without permission")
	}
	s, dir := lifecycleStore(t, 3)
	mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	mustTake(t, s, "ASSIGNMENT-2", 0, "b")
	if err := os.Chmod(recordPath(dir, "ASSIGNMENT-2"), 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(recordPath(dir, "ASSIGNMENT-2"), 0o644) })
	result, err := s.ListView(ViewAll)
	if err != nil {
		t.Fatalf("list failed for one unreadable record: %v", err)
	}
	want := []InvalidFile{{File: ".lifecycle/ASSIGNMENT-2.json", Error: ClassIOError + ": read lifecycle ASSIGNMENT-2: permission denied"}}
	if !reflect.DeepEqual(result.Invalid, want) {
		t.Fatalf("invalid %+v", result.Invalid)
	}
	requireNoPath(t, s, errors.New(result.Invalid[0].Error))
	requireView(t, s, ViewAll, "ASSIGNMENT-1", "ASSIGNMENT-3")
	requireView(t, s, ViewAvailable, "ASSIGNMENT-3")
	requireView(t, s, ViewTaken, "ASSIGNMENT-1")
	_, err = s.Lifecycle("ASSIGNMENT-2")
	requireClass(t, err, ClassIOError)
}

func TestUnreadableLifecycleDirectoryFailsTheList(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root reads directories without permission")
	}
	s, dir := lifecycleStore(t, 2)
	mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	path := filepath.Join(dir, ".lifecycle")
	if err := os.Chmod(path, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0o755) })
	_, err := s.ListView(ViewAll)
	if e := requireClass(t, err, ClassIOError); e.Detail != "read lifecycle directory: permission denied" {
		t.Fatalf("detail %q", e.Detail)
	}
	requireNoPath(t, s, err)
}

func TestUnsearchableLifecycleDirectoryFailsTheList(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root searches directories without permission")
	}
	s, dir := lifecycleStore(t, 2)
	mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	path := filepath.Join(dir, ".lifecycle")
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, 0o755) })
	result, err := s.ListView(ViewAll)
	if e := requireClass(t, err, ClassIOError); e.Detail != "search lifecycle directory: permission denied" {
		t.Fatalf("detail %q, result %+v", e.Detail, result)
	}
	requireNoPath(t, s, err)
}

func TestTransitionSyncsAssignmentDirectoryBeforeWritingRecord(t *testing.T) {
	s, dir := lifecycleStore(t, 1)
	var calls []string
	s.syncDir = func(path string) error {
		info, err := os.Lstat(filepath.Join(path, ".lifecycle"))
		calls = append(calls, fmt.Sprintf("sync %s lifecycle-dir=%v", path, err == nil && info.IsDir()))
		return nil
	}
	s.replace = func(path string, contents []byte, mode fs.FileMode) error {
		calls = append(calls, "replace "+filepath.Base(path))
		return atomicfile.Replace(path, contents, mode)
	}
	l := mustTake(t, s, "ASSIGNMENT-1", 0, "a")
	if _, err := s.Done("ASSIGNMENT-1", l.Holder.Claim); err != nil {
		t.Fatal(err)
	}
	sync := "sync " + dir + " lifecycle-dir=true"
	want := []string{sync, "replace ASSIGNMENT-1.json", sync, "replace ASSIGNMENT-1.json"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls %q", calls)
	}

	failing, failingDir := lifecycleStore(t, 1)
	failing.syncDir = func(path string) error { return &fs.PathError{Op: "sync", Path: path, Err: fs.ErrPermission} }
	_, err := failing.Take("ASSIGNMENT-1", 0, "a")
	if e := requireClass(t, err, ClassIOError); e.Detail != "sync assignments directory: permission denied" {
		t.Fatalf("detail %q", e.Detail)
	}
	requireNoPath(t, failing, err)
	requireAbsent(t, recordPath(failingDir, "ASSIGNMENT-1"))
	failing.syncDir = syncDirectory
	if l := mustTake(t, failing, "ASSIGNMENT-1", 0, "a"); l.Revision != 1 {
		t.Fatalf("retry after failed sync %+v", l)
	}
}
