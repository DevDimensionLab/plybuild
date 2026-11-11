package inbox

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var fixedNow = time.Date(2026, 11, 11, 15, 0, 0, 0, time.UTC)

func conversationID(i int) string { return fmt.Sprintf("conv_%032x", i) }
func messageID(i int) string      { return fmt.Sprintf("msg_%032x", i) }

func strPtr(value string) *string { return &value }

// request returns a valid live request; i varies the identity and the reply minute.
func request(i int) PutRequest {
	return PutRequest{
		Kind:           PutRequestKind,
		ConversationID: conversationID(i),
		MessageID:      messageID(i),
		RepliedAtUTC:   time.Date(2026, 11, 11, 14, 0, 0, 0, time.UTC).Add(time.Duration(i) * time.Minute).Format(time.RFC3339),
		Body:           fmt.Sprintf("Svar nummer %d", i),
		AnswerSource:   "codex.item_completed",
		AgentLabel:     "quality",
		Provider:       "codex",
		WriteKind:      WriteKindLive,
	}
}

func newTestStore(t *testing.T, root, owner string) *Store {
	t.Helper()
	store, err := Open(root, "ply", owner)
	if err != nil {
		t.Fatal(err)
	}
	store.now = func() time.Time { return fixedNow }
	return store
}

func storeRoot(t *testing.T) string { return filepath.Join(t.TempDir(), "inbox") }

func mustPut(t *testing.T, store *Store, req PutRequest) PutResult {
	t.Helper()
	result, err := store.Put(req)
	if err != nil {
		t.Fatalf("put %s: %v", req.MessageID, err)
	}
	return result
}

func mustMark(t *testing.T, store *Store, i int, value MarkValue) MarkResult {
	t.Helper()
	result, err := store.Mark(conversationID(i), messageID(i), value)
	if err != nil {
		t.Fatalf("mark %d %s: %v", i, value, err)
	}
	return result
}

func mustList(t *testing.T, store *Store, limit int, cursor string) List {
	t.Helper()
	list, err := store.List(limit, cursor)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	return list
}

func errorClass(t *testing.T, err error) *Error {
	t.Helper()
	var inboxError *Error
	if !errors.As(err, &inboxError) {
		t.Fatalf("expected *Error, got %v", err)
	}
	return inboxError
}

func expectClass(t *testing.T, err error, class string) *Error {
	t.Helper()
	if err == nil {
		t.Fatalf("expected %s, got success", class)
	}
	inboxError := errorClass(t, err)
	if inboxError.Class != class {
		t.Fatalf("expected %s, got %v", class, err)
	}
	return inboxError
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func messagePath(root, owner string, req PutRequest, canonicalTime string) string {
	return filepath.Join(root, owner, "messages", compactTime(canonicalTime)+"_"+req.ConversationID+"_"+req.MessageID+".json")
}

func listKeys(list List) []string {
	keys := []string{}
	for _, item := range list.Items {
		keys = append(keys, item.ConversationID+"/"+item.MessageID)
	}
	return keys
}
