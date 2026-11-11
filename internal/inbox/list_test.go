package inbox

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// L1: equal reply times break on conversation_id, then message_id.
func TestListOrderIsTotal(t *testing.T) {
	store := newTestStore(t, storeRoot(t), "owner-1")
	same := "2026-11-11T14:00:00Z"
	puts := []struct{ conv, msg int }{{2, 1}, {1, 2}, {1, 1}}
	for _, p := range puts {
		req := request(0)
		req.ConversationID, req.MessageID, req.RepliedAtUTC = conversationID(p.conv), messageID(p.msg), same
		mustPut(t, store, req)
	}
	newer := request(9)
	mustPut(t, store, newer)
	want := []string{
		conversationID(9) + "/" + messageID(9),
		conversationID(1) + "/" + messageID(1),
		conversationID(1) + "/" + messageID(2),
		conversationID(2) + "/" + messageID(1),
	}
	if got := listKeys(mustList(t, store, 20, "")); !reflect.DeepEqual(got, want) {
		t.Fatalf("order %v", got)
	}
}

// L2: stable paging over 45 items.
func TestListPagingIsStable(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	for i := 100; i < 145; i++ {
		mustPut(t, store, request(i))
	}
	first := mustList(t, store, 20, "")
	if len(first.Items) != 20 || first.NextCursor == nil || first.TotalCount != 45 || first.Cursor != nil || first.Limit != 20 {
		t.Fatalf("page 1: %d items, total %d", len(first.Items), first.TotalCount)
	}
	// A newer reply lands before the cursor; an older import lands after it.
	mustPut(t, store, request(500))
	older := request(1)
	older.WriteKind = WriteKindImport
	mustPut(t, store, older)
	second := mustList(t, store, 20, *first.NextCursor)
	third := mustList(t, store, 20, *second.NextCursor)
	if second.Cursor == nil || *second.Cursor != *first.NextCursor {
		t.Fatal("cursor not echoed")
	}
	seen := map[string]bool{}
	all := append(append(append([]string{}, listKeys(first)...), listKeys(second)...), listKeys(third)...)
	for _, key := range all {
		if seen[key] {
			t.Fatalf("repeated %s", key)
		}
		seen[key] = true
	}
	if len(second.Items) != 20 || len(third.Items) != 6 || third.NextCursor != nil {
		t.Fatalf("pages %d/%d next %v", len(second.Items), len(third.Items), third.NextCursor)
	}
	if seen[conversationID(500)+"/"+messageID(500)] {
		t.Fatal("newer item shifted into later pages")
	}
	if !seen[conversationID(1)+"/"+messageID(1)] || third.TotalCount != 47 || third.UnreadCount == nil || *third.UnreadCount != 47 {
		t.Fatalf("older import missing from pass or counts: total %d", third.TotalCount)
	}
	for i := 100; i < 145; i++ {
		if !seen[conversationID(i)+"/"+messageID(i)] {
			t.Fatalf("gap at %d", i)
		}
	}
	again := mustList(t, store, 50, "")
	if listKeys(again)[0] != conversationID(500)+"/"+messageID(500) || len(again.Items) != 47 || again.NextCursor != nil {
		t.Fatalf("next pass %v", listKeys(again)[:1])
	}
}

func TestListCursorIsScopeBound(t *testing.T) {
	root := storeRoot(t)
	ownerA := newTestStore(t, root, "owner-a")
	ownerB := newTestStore(t, root, "owner-b")
	otherProject, err := Open(storeRoot(t), "other", "owner-a")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		mustPut(t, ownerA, request(i))
		mustPut(t, ownerB, request(i))
		mustPut(t, otherProject, request(i))
	}
	page := mustList(t, ownerA, 1, "")
	if page.NextCursor == nil {
		t.Fatal("expected cursor")
	}
	for name, store := range map[string]*Store{"foreign owner": ownerB, "foreign project": otherProject} {
		_, err := store.List(1, *page.NextCursor)
		expectClass(t, err, ClassInvalidArguments)
		_ = name
	}
	garbled := []string{"!!", "e30", base64.RawURLEncoding.EncodeToString([]byte(`{"version":2}`)), *page.NextCursor + "=", strings.ToUpper(*page.NextCursor)}
	for _, cursor := range garbled {
		_, err := ownerA.List(1, cursor)
		expectClass(t, err, ClassInvalidArguments)
	}
	for _, limit := range []int{0, 51, -1} {
		_, err := ownerA.List(limit, "")
		expectClass(t, err, ClassInvalidArguments)
	}
}

// L3: counts are exact only when every counted message was validated.
func TestListCountsAndInvalidEntries(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	for i := 1; i <= 3; i++ {
		mustPut(t, store, request(i))
	}
	mustMark(t, store, 1, MarkRead)
	exact := mustList(t, store, 20, "")
	if !exact.UnreadCountExact || exact.UnreadCountBasis != BasisValidatedItems || *exact.UnreadCount != 2 || exact.ReadState != ReadStateAvailable {
		t.Fatalf("all validated: %+v", exact)
	}
	paged := mustList(t, store, 1, "")
	if paged.UnreadCountExact || paged.UnreadCountBasis != BasisFilenameUnverified || *paged.UnreadCount != 2 || paged.TotalCount != 3 {
		t.Fatalf("paged counts: %+v", paged)
	}

	messages := filepath.Join(root, "owner-1", "messages")
	if err := os.WriteFile(filepath.Join(messages, "garbled.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(messages, ".atomicfile-1.tmp"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	// A corrupt mark for item 2 makes its state unknown, never read.
	writeMark(t, root, "owner-1", 2, `{"kind":"InboxReadMark@1","schema_version":1,"owner":"someone-else","conversation_id":"`+conversationID(2)+`","message_id":"`+messageID(2)+`","marked_at_utc":"2026-11-11T15:00:00.000000000Z"}`)
	list := mustList(t, store, 20, "")
	if list.UnreadCountExact || list.UnreadCountBasis != BasisFilenameUnverified || *list.UnreadCount != 1 || list.TotalCount != 3 {
		t.Fatalf("counts with invalid: %+v", list)
	}
	entries := map[string]string{}
	for _, entry := range list.Invalid {
		entries[entry.Entry] = entry.Error
	}
	if len(entries) != 2 || entries["garbled.json"] == "" || !strings.HasPrefix(entries[conversationID(2)+"_"+messageID(2)+".json"], ClassItemInvalid+": ") {
		t.Fatalf("invalid %v", entries)
	}
	for _, item := range list.Items {
		if item.ConversationID == conversationID(2) && item.Read != nil {
			t.Fatalf("corrupt mark yields read %v", *item.Read)
		}
	}
	// Item 3's mark points nowhere if its message vanishes; the orphan is reported and kept.
	orphan := writeMark(t, root, "owner-1", 9, `{"kind":"InboxReadMark@1","schema_version":1,"owner":"owner-1","conversation_id":"`+conversationID(9)+`","message_id":"`+messageID(9)+`","marked_at_utc":"2026-11-11T15:00:00.000000000Z"}`)
	list = mustList(t, store, 20, "")
	if _, err := os.Lstat(orphan); err != nil || len(list.Invalid) != 3 {
		t.Fatalf("orphan mark: %v %+v", err, list.Invalid)
	}
}

func TestListReadStateUnavailable(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	mustPut(t, store, request(1))
	if err := os.WriteFile(filepath.Join(root, "owner-1", "read"), []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	list := mustList(t, store, 20, "")
	if list.ReadState != ReadStateUnavailable || list.UnreadCount != nil || list.UnreadCountExact || list.UnreadCountBasis != BasisUnavailable || list.Items[0].Read != nil {
		t.Fatalf("unavailable: %+v", list)
	}
	item, err := store.Show(conversationID(1), messageID(1))
	if err != nil || item.Read != nil {
		t.Fatalf("show read %v %v", item.Read, err)
	}
}

// L4: an item with corrupt content stays in place as invalid and paging advances.
func TestListCorruptPageItemStaysInPlace(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	for i := 1; i <= 3; i++ {
		mustPut(t, store, request(i))
	}
	corrupt := messagePath(root, "owner-1", request(2), mustCanonical(t, request(2).RepliedAtUTC))
	if err := os.WriteFile(corrupt, bytes.Replace(readFile(t, corrupt), []byte(`"owner-1"`), []byte(`"owner-2"`), 1), 0o600); err != nil {
		t.Fatal(err)
	}
	first := mustList(t, store, 1, "")
	second := mustList(t, store, 1, *first.NextCursor)
	third := mustList(t, store, 1, *second.NextCursor)
	if len(second.Items) != 0 || len(second.Invalid) != 1 || second.Invalid[0].Entry != filepath.Base(corrupt) || second.NextCursor == nil {
		t.Fatalf("page 2 %+v", second)
	}
	if len(third.Items) != 1 || third.Items[0].MessageID != messageID(1) || third.NextCursor != nil {
		t.Fatalf("page 3 %+v", third)
	}
}

// L5: marking is idempotent, durable and never touches messages or order.
func TestMarkIsIdempotentAndNeverReorders(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	for i := 1; i <= 3; i++ {
		mustPut(t, store, request(i))
	}
	path := messagePath(root, "owner-1", request(2), mustCanonical(t, request(2).RepliedAtUTC))
	before, order := readFile(t, path), listKeys(mustList(t, store, 20, ""))
	steps := []struct {
		value  MarkValue
		effect string
		read   bool
	}{{MarkRead, EffectMarked, true}, {MarkRead, EffectNone, true}, {MarkUnread, EffectMarked, false}, {MarkUnread, EffectNone, false}, {MarkRead, EffectMarked, true}}
	for _, step := range steps {
		result := mustMark(t, store, 2, step.value)
		if result.Kind != MarkResultKind || result.Effect != step.effect || result.Read != step.read || result.ConversationID != conversationID(2) {
			t.Fatalf("mark %s: %+v", step.value, result)
		}
	}
	markInfo, err := os.Lstat(filepath.Join(root, "owner-1", "read", conversationID(2)+"_"+messageID(2)+".json"))
	if err != nil || markInfo.Mode().Perm() != 0o600 {
		t.Fatalf("mark file %v %v", markInfo, err)
	}
	if !bytes.Equal(before, readFile(t, path)) || !reflect.DeepEqual(order, listKeys(mustList(t, store, 20, ""))) {
		t.Fatal("mark changed message or order")
	}
	restarted, err := Open(root, "ply", "owner-1")
	if err != nil {
		t.Fatal(err)
	}
	if item, err := restarted.Show(conversationID(2), messageID(2)); err != nil || item.Read == nil || !*item.Read {
		t.Fatalf("mark lost on restart: %v", err)
	}
	_, err = store.Mark(conversationID(9), messageID(9), MarkRead)
	expectClass(t, err, ClassItemNotFound)
	_, err = store.Mark(conversationID(1), messageID(1), MarkValue("seen"))
	expectClass(t, err, ClassInvalidArguments)
	writeMark(t, root, "owner-1", 3, "garbage")
	_, err = store.Mark(conversationID(3), messageID(3), MarkUnread)
	expectClass(t, err, ClassItemInvalid)
}

// L6: no code path but explicit unread removes a mark.
func TestNoPruningAfterManyPuts(t *testing.T) {
	store := newTestStore(t, storeRoot(t), "owner-1")
	for i := 0; i < 1000; i++ {
		req := request(i)
		if i%2 == 1 {
			req.WriteKind = WriteKindImport
		}
		mustPut(t, store, req)
		if i%10 == 0 {
			mustMark(t, store, i, MarkRead)
		}
		if i%7 == 0 {
			mustPut(t, store, req)
		}
	}
	for i := 0; i < 1000; i += 10 {
		item, err := store.Show(conversationID(i), messageID(i))
		if err != nil || item.Read == nil || !*item.Read {
			t.Fatalf("mark %d lost: %v", i, err)
		}
	}
	if list := mustList(t, store, 50, ""); *list.UnreadCount != 900 || list.TotalCount != 1000 {
		t.Fatalf("counts %d/%d", *list.UnreadCount, list.TotalCount)
	}
}

// L7: owners are isolated, and owner values and symlinks cannot escape.
func TestOwnerIsolationAndEscapes(t *testing.T) {
	root := storeRoot(t)
	a := newTestStore(t, root, "owner-a")
	b := newTestStore(t, root, "owner-b")
	mustPut(t, a, request(1))
	mustMark(t, a, 1, MarkRead)
	if list := mustList(t, b, 20, ""); list.TotalCount != 0 || len(list.Items) != 0 {
		t.Fatal("owner b sees owner a")
	}
	_, err := b.Show(conversationID(1), messageID(1))
	expectClass(t, err, ClassItemNotFound)
	_, err = b.Mark(conversationID(1), messageID(1), MarkRead)
	expectClass(t, err, ClassItemNotFound)
	mustPut(t, b, request(1))
	if item, _ := b.Show(conversationID(1), messageID(1)); item.Read == nil || *item.Read {
		t.Fatal("owner b inherited owner a's mark")
	}

	for _, owner := range []string{".", "..", "", "a/b", `a\b`, strings.Repeat("o", 65), "ø"} {
		if _, err := Open(root, "ply", owner); err == nil {
			t.Fatalf("owner %q accepted", owner)
		} else {
			expectClass(t, err, ClassInvalidArguments)
		}
	}
	if _, err := Open("relative/inbox", "ply", "owner-a"); err == nil {
		t.Fatal("relative root accepted")
	}
}

func TestSymlinkedComponentsNeverEscape(t *testing.T) {
	for _, component := range []string{"owner", "messages", "read"} {
		t.Run(component, func(t *testing.T) {
			root := storeRoot(t)
			outside := t.TempDir()
			store := newTestStore(t, root, "owner-1")
			if err := os.MkdirAll(filepath.Join(root, "owner-1"), 0o700); err != nil {
				t.Fatal(err)
			}
			link := map[string]string{"owner": filepath.Join(root, "owner-1"), "messages": filepath.Join(root, "owner-1", "messages"), "read": filepath.Join(root, "owner-1", "read")}[component]
			if component != "owner" {
				mustPut(t, store, request(1))
				if err := os.RemoveAll(link); err != nil {
					t.Fatal(err)
				}
			} else if err := os.Remove(link); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(outside, link); err != nil {
				t.Fatal(err)
			}
			_, putErr := store.Put(request(2))
			_, markErr := store.Mark(conversationID(1), messageID(1), MarkRead)
			list, listErr := store.List(20, "")
			_, _ = store.Show(conversationID(1), messageID(1))
			entries, err := os.ReadDir(outside)
			if err != nil || len(entries) != 0 {
				t.Fatalf("escaped into symlink target: %v %v", entries, err)
			}
			if component == "read" {
				if putErr != nil || markErr == nil || listErr != nil || list.ReadState != ReadStateUnavailable {
					t.Fatalf("read symlink: put %v mark %v list %v %s", putErr, markErr, listErr, list.ReadState)
				}
				return
			}
			if putErr == nil || markErr == nil || listErr == nil {
				t.Fatalf("symlink accepted: put %v mark %v list %v", putErr, markErr, listErr)
			}
			expectClass(t, putErr, ClassItemInvalid)
		})
	}
}

// L8: a missing store is an empty exact list and reads create nothing.
func TestMissingStoreIsEmptyAndReadsCreateNothing(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	list := mustList(t, store, 20, "")
	if list.Kind != ListKind || list.TotalCount != 0 || list.UnreadCount == nil || *list.UnreadCount != 0 || !list.UnreadCountExact ||
		list.UnreadCountBasis != BasisValidatedItems || list.ReadState != ReadStateAvailable || list.NextCursor != nil || list.Items == nil || list.Invalid == nil {
		t.Fatalf("empty list %+v", list)
	}
	_, _ = store.Show(conversationID(1), messageID(1))
	_, _ = store.Mark(conversationID(1), messageID(1), MarkRead)
	if _, err := os.Lstat(root); !os.IsNotExist(err) {
		t.Fatal("read created the store")
	}
	mustPut(t, newTestStore(t, root, "owner-2"), request(1))
	if list := mustList(t, store, 20, ""); list.TotalCount != 0 || !list.UnreadCountExact {
		t.Fatal("missing owner directory is not empty")
	}
	if _, err := os.Lstat(filepath.Join(root, "owner-1")); !os.IsNotExist(err) {
		t.Fatal("list created the owner directory")
	}
}

// C5: Ply-side budget on the local disk.
func TestBudgetWithTenThousandItems(t *testing.T) {
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	mustPut(t, store, request(0))
	for i := 1; i < 10000; i++ {
		req := request(i)
		message, err := store.newMessage(req)
		if err != nil {
			t.Fatal(err)
		}
		contents, err := encodeJSON(message)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(messagePath(root, "owner-1", req, message.RepliedAtUTC), contents, 0o600); err != nil {
			t.Fatal(err)
		}
		if i%3 == 0 {
			writeMark(t, root, "owner-1", i, `{"kind":"InboxReadMark@1","schema_version":1,"owner":"owner-1","conversation_id":"`+req.ConversationID+`","message_id":"`+req.MessageID+`","marked_at_utc":"2026-11-11T15:00:00.000000000Z"}`)
		}
	}
	measure := func(name string, run func() error) {
		started := time.Now()
		if err := run(); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if elapsed := time.Since(started); elapsed >= time.Second {
			t.Fatalf("%s took %v", name, elapsed)
		}
	}
	measure("list", func() error {
		list, err := store.List(50, "")
		if err == nil && (list.TotalCount != 10000 || len(list.Items) != 50 || len(list.Invalid) != 0) {
			t.Fatalf("total %d items %d invalid %v", list.TotalCount, len(list.Items), list.Invalid)
		}
		return err
	})
	measure("put", func() error { _, err := store.Put(request(20000)); return err })
	measure("mark", func() error { _, err := store.Mark(conversationID(20000), messageID(20000), MarkRead); return err })
	measure("show", func() error { _, err := store.Show(conversationID(20000), messageID(20000)); return err })
}
