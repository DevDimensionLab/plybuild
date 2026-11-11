package inbox

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var updateFixtures = flag.Bool("update-fixtures", false, "rewrite schemas/ply-inbox-v1 from the store")

const fixtureDirectory = "../../schemas/ply-inbox-v1"

// The fixtures are the producer contract plyapp tests against. They are
// produced by the real store and decoded back with the production types.
func TestFixturesMatchStoreBehaviour(t *testing.T) {
	putRequest := readFixture(t, "put-request.json")
	req, err := DecodePutRequest(putRequest)
	if err != nil {
		t.Fatal(err)
	}
	root := storeRoot(t)
	store := newTestStore(t, root, "owner-1")
	created := mustPut(t, store, req)
	older := request(1)
	older.WriteKind = WriteKindImport
	mustPut(t, store, older)
	titled := request(2)
	titled.RepliedAtUTC, titled.Title = "2026-11-11T14:10:00Z", strPtr("Review klar")
	mustPut(t, store, titled)

	stored := readFile(t, messagePath(root, "owner-1", req, created.Item.RepliedAtUTC))
	marked := mustMark(t, store, 0, MarkRead)
	shown, err := store.Show(req.ConversationID, req.MessageID)
	if err != nil {
		t.Fatal(err)
	}
	list := mustList(t, store, 2, "")
	conflicting := req
	conflicting.Body += "\nEndret."
	_, err = store.Put(conflicting)
	conflict := expectClass(t, err, ClassConflict)
	if err := os.Rename(filepath.Join(root, "owner-1", "read"), filepath.Join(root, "owner-1", "read-moved")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "owner-1", "read"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	unavailable := mustList(t, store, 20, "")

	outputs := map[string]any{"item.json": shown, "list.json": list, "list-unavailable.json": unavailable, "mark.json": marked, "error-conflict.json": NewErrorResult(conflict)}
	if *updateFixtures {
		writeFixture(t, "message.json", stored)
		for name, value := range outputs {
			encoded, err := encodeJSON(value)
			if err != nil {
				t.Fatal(err)
			}
			writeFixture(t, name, encoded)
		}
	}
	if !bytes.Equal(readFixture(t, "message.json"), stored) {
		t.Fatalf("message.json differs from the stored file:\n%s", stored)
	}
	for name, value := range outputs {
		encoded, err := encodeJSON(value)
		if err != nil {
			t.Fatal(err)
		}
		fixture := readFixture(t, name)
		if !bytes.Equal(fixture, encoded) {
			t.Fatalf("%s differs from store output:\n%s", name, encoded)
		}
		decoded := reflect.New(reflect.TypeOf(value))
		decoder := json.NewDecoder(bytes.NewReader(fixture))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(decoded.Interface()); err != nil {
			t.Fatalf("%s does not decode with %T: %v", name, value, err)
		}
		if !reflect.DeepEqual(decoded.Elem().Interface(), value) {
			t.Fatalf("%s round trip differs", name)
		}
	}
}

// S6: golden vector. The fixture message's payload hashes to its payload_sha256.
func TestFixtureMessageGoldenPayloadHash(t *testing.T) {
	contents := readFixture(t, "message.json")
	var probe struct {
		RepliedAtUTC   string `json:"replied_at_utc"`
		ConversationID string `json:"conversation_id"`
		MessageID      string `json:"message_id"`
	}
	if err := json.Unmarshal(contents, &probe); err != nil {
		t.Fatal(err)
	}
	name := compactTime(probe.RepliedAtUTC) + "_" + probe.ConversationID + "_" + probe.MessageID + ".json"
	message, err := decodeMessageFile(contents, name, "ply", "owner-1")
	if err != nil {
		t.Fatal(err)
	}
	// Computed independently (Python json.dumps sort_keys, compact, ensure_ascii=False) from the RFC 8785 payload.
	const golden = "sha256:366ee08325e6ffe905382270509dc804ebaccdcdbed05553bde5343700df00a2"
	hash, err := payloadSHA256("ply", "owner-1", payloadOf(message))
	if err != nil {
		t.Fatal(err)
	}
	if hash != message.PayloadSHA256 || hash != golden {
		t.Fatalf("payload hash %s, stored %s, golden %s", hash, message.PayloadSHA256, golden)
	}
}

func TestFixtureDerivedTitleCases(t *testing.T) {
	var cases struct {
		Kind          string `json:"kind"`
		Rule          string `json:"rule"`
		FallbackTitle string `json:"fallback_title"`
		MaxCodePoints int    `json:"max_code_points"`
		Cases         []struct {
			Name         string `json:"name"`
			Body         string `json:"body"`
			DerivedTitle string `json:"derived_title"`
		} `json:"cases"`
	}
	decoder := json.NewDecoder(bytes.NewReader(readFixture(t, "derived-title-cases.json")))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cases); err != nil {
		t.Fatal(err)
	}
	if cases.Kind != "InboxDerivedTitleCases@1" || cases.Rule != DerivedTitleRule || cases.FallbackTitle != FallbackTitle || cases.MaxCodePoints != maxTitleRunes || len(cases.Cases) < 10 {
		t.Fatalf("rule %s with %d cases", cases.Rule, len(cases.Cases))
	}
	for _, c := range cases.Cases {
		if got := DeriveTitle(c.Body); got != c.DerivedTitle {
			t.Errorf("%s: DeriveTitle = %q, fixture %q", c.Name, got, c.DerivedTitle)
		}
	}
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	return readFile(t, filepath.Join(fixtureDirectory, name))
}

func writeFixture(t *testing.T, name string, contents []byte) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(fixtureDirectory, name), contents, 0o644); err != nil {
		t.Fatal(err)
	}
}
