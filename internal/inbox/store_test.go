package inbox

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

// S1: a new put creates one strict InboxMessage@1 file with mode 0600.
func TestPutCreatesMessageFile(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	req := request(1)
	req.RepliedAtUTC = "2026-11-11T14:05:00.123456Z"
	result := mustPut(t, store, req)
	if result.Kind != PutResultKind || result.SchemaVersion != 1 || result.Effect != EffectCreated {
		t.Fatalf("result %+v", result)
	}
	path := filepath.Join(root, "owner-1", "messages", "20261111T140500.123456000Z_"+req.ConversationID+"_"+req.MessageID+".json")
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("file mode %v", info.Mode())
	}
	for _, dir := range []string{root, filepath.Join(root, "owner-1"), filepath.Join(root, "owner-1", "messages")} {
		if info, err := os.Lstat(dir); err != nil || info.Mode().Perm() != 0o700 {
			t.Fatalf("directory %s mode %v %v", filepath.Base(dir), info.Mode(), err)
		}
	}
	message, err := decodeMessageFile(readFile(t, path), filepath.Base(path), "ply", "owner-1")
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256([]byte(req.Body))
	want := Message{
		Kind: MessageKind, SchemaVersion: 1, ProjectID: "ply", Owner: "owner-1",
		ConversationID: req.ConversationID, MessageID: req.MessageID,
		RepliedAtUTC: "2026-11-11T14:05:00.123456000Z", DerivedTitle: req.Body, DerivedTitleRule: DerivedTitleRule,
		Body: req.Body, BodySHA256: hex.EncodeToString(digest[:]), AnswerSource: req.AnswerSource,
		AgentLabel: "quality", Provider: "codex", WriteKind: WriteKindLive,
		RecordedAtUTC: "2026-11-11T15:00:00.000000000Z", PayloadSHA256: result.Item.PayloadSHA256,
	}
	if !reflect.DeepEqual(message, want) || !reflect.DeepEqual(result.Item, want) {
		t.Fatalf("stored %+v\nresult %+v\nwant %+v", message, result.Item, want)
	}
	if !strings.HasPrefix(want.PayloadSHA256, "sha256:") || len(want.PayloadSHA256) != 71 {
		t.Fatalf("payload hash %q", want.PayloadSHA256)
	}
}

// S2, S3: replay is a no-op that keeps bytes, mark, first metadata and position.
func TestReplayIsUnchangedAndKeepsEverything(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	for i := 1; i <= 3; i++ {
		mustPut(t, store, request(i))
	}
	mustMark(t, store, 2, MarkRead)
	before := mustList(t, store, 20, "")
	path := messagePath(root, "owner-1", request(2), mustCanonical(t, request(2).RepliedAtUTC))
	bytesBefore := readFile(t, path)

	store.now = func() time.Time { return fixedNow.Add(time.Hour) }
	replay := request(2)
	replay.AgentLabel, replay.Provider, replay.WriteKind = "renamed", "claude", WriteKindImport
	result := mustPut(t, store, replay)
	if result.Effect != EffectUnchanged {
		t.Fatalf("effect %s", result.Effect)
	}
	if result.Item.AgentLabel != "quality" || result.Item.Provider != "codex" || result.Item.WriteKind != WriteKindLive || result.Item.RecordedAtUTC != "2026-11-11T15:00:00.000000000Z" {
		t.Fatalf("first-write metadata changed: %+v", result.Item)
	}
	if !bytes.Equal(bytesBefore, readFile(t, path)) {
		t.Fatal("message file bytes changed")
	}
	after := mustList(t, store, 20, "")
	if !reflect.DeepEqual(listKeys(before), listKeys(after)) || after.Items[1].Read == nil || !*after.Items[1].Read {
		t.Fatalf("order or mark changed: %v %v", listKeys(after), after.Items[1].Read)
	}
}

// S4: a different immutable field is a conflict naming exactly that field.
func TestReplayWithDifferentPayloadConflicts(t *testing.T) {
	cases := []struct {
		field  string
		mutate func(*PutRequest)
	}{
		{"title", func(r *PutRequest) { r.Title = strPtr("Ny tittel") }},
		{"body", func(r *PutRequest) { r.Body = "Hemmelig annen tekst" }},
		{"truncated", func(r *PutRequest) { r.Truncated, r.FullSHA256 = true, strPtr(strings.Repeat("b", 64)) }},
		{"answer_source", func(r *PutRequest) { r.AnswerSource = "claude.result" }},
		{"replied_at_utc", func(r *PutRequest) { r.RepliedAtUTC = "2030-01-01T00:00:00Z" }},
	}
	for _, c := range cases {
		t.Run(c.field, func(t *testing.T) {
			root := storeRoot(t)
			store := newTestStore(t, root, "owner-1")
			original := request(1)
			first := mustPut(t, store, original)
			before := snapshot(t, root)
			changed := request(1)
			c.mutate(&changed)
			_, err := store.Put(changed)
			conflict := expectClass(t, err, ClassConflict)
			want := []string{c.field}
			if c.field == "truncated" {
				want = []string{"truncated", "full_sha256"}
			}
			if !reflect.DeepEqual(conflict.ConflictingFields, want) {
				t.Fatalf("fields %v, want %v", conflict.ConflictingFields, want)
			}
			if conflict.ExistingPayloadSHA256 != first.Item.PayloadSHA256 || conflict.SubmittedPayloadSHA256 == "" || conflict.SubmittedPayloadSHA256 == first.Item.PayloadSHA256 {
				t.Fatalf("hashes %q %q", conflict.ExistingPayloadSHA256, conflict.SubmittedPayloadSHA256)
			}
			if conflict.Retryable() || strings.Contains(conflict.Error(), "Hemmelig") || strings.Contains(conflict.Error(), "Ny tittel") {
				t.Fatalf("conflict %v", conflict)
			}
			if !reflect.DeepEqual(before, snapshot(t, root)) {
				t.Fatal("conflict wrote files")
			}
		})
	}
}

// S9: a put whose result was lost is reconciled by the next put.
func TestLostReplyRetryIsUnchanged(t *testing.T) {
	store := newTestStore(t, storeRoot(t), "owner-1")
	req := request(7)
	req.Truncated, req.FullSHA256 = true, strPtr(strings.Repeat("c", 64))
	_, _ = store.Put(req)
	again, err := Open(store.root, "ply", "owner-1")
	if err != nil {
		t.Fatal(err)
	}
	result := mustPut(t, again, req)
	if result.Effect != EffectUnchanged || result.Item.FullSHA256 == nil || *result.Item.FullSHA256 != strings.Repeat("c", 64) {
		t.Fatalf("retry %+v", result)
	}
}

// S10: an uncertain write is a retryable io error and the next put reconciles.
func TestUncertainWriteIsRetryableIOError(t *testing.T) {
	store := newTestStore(t, storeRoot(t), "owner-1")
	store.replace = func(path string, contents []byte, mode fs.FileMode) error {
		if err := atomicfile.Replace(path, contents, mode); err != nil {
			return err
		}
		return &atomicfile.Error{Op: "sync directory", Path: filepath.Dir(path), Replaced: true, Err: os.ErrInvalid}
	}
	_, err := store.Put(request(1))
	ioError := expectClass(t, err, ClassIOError)
	if !ioError.Retryable() || strings.Contains(ioError.Detail, store.root) {
		t.Fatalf("io error %v", ioError)
	}
	store.replace = atomicfile.Replace
	if result := mustPut(t, store, request(1)); result.Effect != EffectUnchanged {
		t.Fatalf("effect %s", result.Effect)
	}
}

// S11: a held lock gives inbox_busy; concurrent puts create exactly once.
func TestLockTimeoutAndConcurrentPuts(t *testing.T) {
	if defaultLockTimeout != 2*time.Second {
		t.Fatalf("lock timeout %v", defaultLockTimeout)
	}
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	mustPut(t, store, request(1))
	holder := newTestStore(t, root, "owner-2")
	unlock, err := holder.lock()
	if err != nil {
		t.Fatal(err)
	}
	store.lockTimeout = 150 * time.Millisecond
	started := time.Now()
	_, err = store.Put(request(2))
	busy := expectClass(t, err, ClassBusy)
	if !busy.Retryable() || time.Since(started) < 150*time.Millisecond {
		t.Fatalf("busy %v after %v", busy, time.Since(started))
	}
	_, err = store.Mark(conversationID(1), messageID(1), MarkRead)
	expectClass(t, err, ClassBusy)
	unlock()

	effects := make(chan string, 8)
	var wait sync.WaitGroup
	for i := 0; i < 8; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			worker := newTestStore(t, root, "owner-1")
			result, err := worker.Put(request(3))
			if err != nil {
				effects <- err.Error()
				return
			}
			effects <- result.Effect
		}()
	}
	wait.Wait()
	close(effects)
	counts := map[string]int{}
	for effect := range effects {
		counts[effect]++
	}
	if counts[EffectCreated] != 1 || counts[EffectUnchanged] != 7 {
		t.Fatalf("effects %v", counts)
	}
}

// S12: corrupt or duplicate files are never overwritten and are refused.
func TestCorruptAndDuplicateFilesAreRefused(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	mustPut(t, store, request(1))
	path := messagePath(root, "owner-1", request(1), mustCanonical(t, request(1).RepliedAtUTC))
	if err := os.WriteFile(path, []byte(`{"kind":"InboxMessage@1"`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := store.Put(request(1))
	expectClass(t, err, ClassItemInvalid)
	if string(readFile(t, path)) != `{"kind":"InboxMessage@1"` {
		t.Fatal("corrupt file overwritten")
	}
	_, err = store.Show(conversationID(1), messageID(1))
	expectClass(t, err, ClassItemInvalid)

	mustPut(t, store, request(2))
	original := messagePath(root, "owner-1", request(2), mustCanonical(t, request(2).RepliedAtUTC))
	duplicate := filepath.Join(filepath.Dir(original), "20200101T000000.000000000Z_"+conversationID(2)+"_"+messageID(2)+".json")
	if err := os.WriteFile(duplicate, readFile(t, original), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = store.Put(request(2))
	expectClass(t, err, ClassItemInvalid)
	_, err = store.Mark(conversationID(2), messageID(2), MarkRead)
	expectClass(t, err, ClassItemInvalid)
	list := mustList(t, store, 20, "")
	invalid := map[string]bool{}
	for _, entry := range list.Invalid {
		invalid[entry.Entry] = true
	}
	if !invalid[filepath.Base(original)] || !invalid[filepath.Base(duplicate)] || list.TotalCount != 1 || list.UnreadCountExact {
		t.Fatalf("list %+v", list)
	}
	if _, err := os.Lstat(filepath.Join(root, "owner-1", "read", conversationID(2)+"_"+messageID(2)+".json")); !os.IsNotExist(err) {
		t.Fatal("refused mark wrote a mark")
	}
}

func TestShowUnknownItemAndReadState(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	_, err := store.Show(conversationID(1), messageID(1))
	expectClass(t, err, ClassItemNotFound)
	if _, err := os.Lstat(root); !os.IsNotExist(err) {
		t.Fatal("show created the store")
	}
	_, err = store.Show("conv_x", messageID(1))
	expectClass(t, err, ClassInvalidArguments)
	req := request(1)
	req.Title = strPtr("Agenttittel")
	mustPut(t, store, req)
	item, err := store.Show(conversationID(1), messageID(1))
	if err != nil {
		t.Fatal(err)
	}
	if item.Kind != ItemKind || item.TitleSource != TitleSourceAgent || item.DisplayTitle != "Agenttittel" || item.Read == nil || *item.Read {
		t.Fatalf("item %+v", item)
	}
	mustMark(t, store, 1, MarkRead)
	if item, _ = store.Show(conversationID(1), messageID(1)); item.Read == nil || !*item.Read {
		t.Fatal("read mark not shown")
	}
	writeMark(t, root, "owner-1", 1, `{"kind":"InboxReadMark@1"}`)
	if item, err = store.Show(conversationID(1), messageID(1)); err != nil || item.Read != nil {
		t.Fatalf("corrupt mark must give read null: %v %v", item.Read, err)
	}
}

func mustCanonical(t *testing.T, value string) string {
	t.Helper()
	canonical, err := canonicalTime(value)
	if err != nil {
		t.Fatal(err)
	}
	return canonical
}

func writeMark(t *testing.T, root, owner string, i int, contents string) string {
	t.Helper()
	path := filepath.Join(root, owner, "read", conversationID(i)+"_"+messageID(i)+".json")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Name() == lockName {
			return nil
		}
		if entry.IsDir() {
			out[path] = "dir"
			return nil
		}
		contents, err := os.ReadFile(path)
		out[path] = string(contents)
		return err
	})
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	return out
}

// S4: a different full_sha256 alone is a conflict on that field only.
func TestReplayWithDifferentFullDigestConflicts(t *testing.T) {
	store := newTestStore(t, storeRoot(t), "owner-1")
	req := request(1)
	req.Truncated, req.FullSHA256 = true, strPtr(strings.Repeat("a", 64))
	mustPut(t, store, req)
	req.FullSHA256 = strPtr(strings.Repeat("b", 64))
	_, err := store.Put(req)
	if conflict := expectClass(t, err, ClassConflict); !reflect.DeepEqual(conflict.ConflictingFields, []string{"full_sha256"}) {
		t.Fatalf("fields %v", conflict.ConflictingFields)
	}
}

// §3.3: an unreadable existing file is inbox_item_invalid and is never overwritten.
func TestUnreadableExistingFileIsItemInvalid(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("chmod 0 does not make a file unreadable for root")
	}
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	mustPut(t, store, request(1))
	path := messagePath(root, "owner-1", request(1), mustCanonical(t, request(1).RepliedAtUTC))
	if err := os.Chmod(path, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(path, fileMode) })
	_, err := store.Put(request(1))
	if invalid := expectClass(t, err, ClassItemInvalid); invalid.Retryable() {
		t.Fatal("unreadable file must not be retryable")
	}
}

func TestMarkUnreadWithoutMarksCreatesNothing(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	mustPut(t, store, request(1))
	if result := mustMark(t, store, 1, MarkUnread); result.Effect != EffectNone || result.Read {
		t.Fatalf("unread %+v", result)
	}
	if _, err := os.Lstat(filepath.Join(root, "owner-1", "read")); !os.IsNotExist(err) {
		t.Fatal("unread created the read directory")
	}
}
