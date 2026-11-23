package inbox

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

// S7: field rules reject with inbox_invalid_arguments and write nothing.
func TestPutRejectsFieldRuleViolations(t *testing.T) {
	bodyDigest := sha256.Sum256([]byte("Svar nummer 1"))
	cases := []struct {
		name   string
		mutate func(*PutRequest)
	}{
		{"kind", func(r *PutRequest) { r.Kind = "InboxPutRequest@2" }},
		{"conversation pattern", func(r *PutRequest) { r.ConversationID = "conv_ABC" }},
		{"conversation uppercase hex", func(r *PutRequest) { r.ConversationID = "conv_" + strings.Repeat("A", 32) }},
		{"message pattern", func(r *PutRequest) { r.MessageID = "msg_" + strings.Repeat("0", 31) }},
		{"label as identity", func(r *PutRequest) { r.MessageID = "quality" }},
		{"time without Z", func(r *PutRequest) { r.RepliedAtUTC = "2026-11-11T14:05:00+00:00" }},
		{"time with ten fraction digits", func(r *PutRequest) { r.RepliedAtUTC = "2026-11-11T14:05:00.1234567890Z" }},
		{"time lowercase z", func(r *PutRequest) { r.RepliedAtUTC = "2026-11-11T14:05:00z" }},
		{"time impossible date", func(r *PutRequest) { r.RepliedAtUTC = "2026-02-30T14:05:00Z" }},
		{"title newline", func(r *PutRequest) { r.Title = strPtr("to\nlinjer") }},
		{"title control", func(r *PutRequest) { r.Title = strPtr("bell\a") }},
		{"title format character", func(r *PutRequest) { r.Title = strPtr("zero\u200bwidth") }},
		{"title line separator", func(r *PutRequest) { r.Title = strPtr("a b") }},
		{"title trailing newline", func(r *PutRequest) { r.Title = strPtr("Hello\n") }},
		{"title leading next line", func(r *PutRequest) { r.Title = strPtr("\u0085Hello") }},
		{"title trailing line separator", func(r *PutRequest) { r.Title = strPtr("Hello\u2028") }},
		{"title empty", func(r *PutRequest) { r.Title = strPtr("") }},
		{"title blank", func(r *PutRequest) { r.Title = strPtr("   ") }},
		{"title 61", func(r *PutRequest) { r.Title = strPtr(strings.Repeat("t", 61)) }},
		{"body too long", func(r *PutRequest) { r.Body = strings.Repeat("é", maxBodyRunes+1) }},
		{"body invalid utf8", func(r *PutRequest) { r.Body = "a\xffb" }},
		{"digest without truncation", func(r *PutRequest) { r.FullSHA256 = strPtr(strings.Repeat("a", 64)) }},
		{"truncated without digest", func(r *PutRequest) { r.Truncated = true }},
		{"truncated with body digest", func(r *PutRequest) {
			r.Truncated, r.FullSHA256 = true, strPtr(hex.EncodeToString(bodyDigest[:]))
		}},
		{"digest uppercase", func(r *PutRequest) { r.Truncated, r.FullSHA256 = true, strPtr(strings.Repeat("A", 64)) }},
		{"digest short", func(r *PutRequest) { r.Truncated, r.FullSHA256 = true, strPtr(strings.Repeat("a", 63)) }},
		{"answer source", func(r *PutRequest) { r.AnswerSource = "Codex" }},
		{"agent label empty", func(r *PutRequest) { r.AgentLabel = "" }},
		{"agent label 65", func(r *PutRequest) { r.AgentLabel = strings.Repeat("l", 65) }},
		{"agent label newline", func(r *PutRequest) { r.AgentLabel = "a\nb" }},
		{"provider", func(r *PutRequest) { r.Provider = "Codex" }},
		{"write kind", func(r *PutRequest) { r.WriteKind = "backfill" }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := storeRoot(t)
			store := newTestStore(t, root, "owner-1")
			req := request(1)
			c.mutate(&req)
			_, err := store.Put(req)
			inboxError := expectClass(t, err, ClassInvalidArguments)
			if inboxError.Retryable() {
				t.Fatal("invalid arguments must not be retryable")
			}
			if strings.Contains(inboxError.Detail, "Svar") || strings.Contains(inboxError.Detail, root) {
				t.Fatalf("detail leaks input: %q", inboxError.Detail)
			}
			list := mustList(t, store, 20, "")
			if list.TotalCount != 0 {
				t.Fatalf("rejected put stored an item")
			}
		})
	}
}

func TestPutAcceptsBoundaryValues(t *testing.T) {
	store := newTestStore(t, storeRoot(t), "owner-1")
	cases := []func(*PutRequest){
		func(r *PutRequest) { r.Body = strings.Repeat("é", maxBodyRunes) },
		func(r *PutRequest) { r.Body = "" },
		func(r *PutRequest) { r.Title = strPtr(strings.Repeat("🙂", 60)) },
		func(r *PutRequest) { r.Title = strPtr("  Trimmet tittel  ") },
		func(r *PutRequest) { r.Truncated, r.FullSHA256 = true, strPtr(strings.Repeat("a", 64)) },
		func(r *PutRequest) { r.WriteKind = WriteKindImport },
		func(r *PutRequest) { r.RepliedAtUTC = "2026-11-11T14:05:00.123456789Z" },
	}
	for i, mutate := range cases {
		req := request(i)
		mutate(&req)
		if result := mustPut(t, store, req); result.Effect != EffectCreated {
			t.Fatalf("case %d: effect %s", i, result.Effect)
		}
	}
	item, err := store.Show(conversationID(3), messageID(3))
	if err != nil {
		t.Fatal(err)
	}
	if item.Title == nil || *item.Title != "Trimmet tittel" {
		t.Fatalf("title not trimmed: %v", item.Title)
	}
}

// S5: equal instants written differently hash the same.
func TestEqualInstantsHaveOneCanonicalTime(t *testing.T) {
	store := newTestStore(t, storeRoot(t), "owner-1")
	req := request(1)
	req.RepliedAtUTC = "2026-11-11T14:05:05Z"
	first := mustPut(t, store, req)
	if first.Item.RepliedAtUTC != "2026-11-11T14:05:05.000000000Z" {
		t.Fatalf("canonical time %s", first.Item.RepliedAtUTC)
	}
	req.RepliedAtUTC = "2026-11-11T14:05:05.000Z"
	second := mustPut(t, store, req)
	if second.Effect != EffectUnchanged || second.Item.PayloadSHA256 != first.Item.PayloadSHA256 {
		t.Fatalf("second put %s %s", second.Effect, second.Item.PayloadSHA256)
	}
}

// C3 (store side): the request decoder is closed and bounded.
func TestDecodePutRequestIsStrict(t *testing.T) {
	valid := `{"kind":"InboxPutRequest@1","conversation_id":"conv_00000000000000000000000000000001","message_id":"msg_00000000000000000000000000000001","replied_at_utc":"2026-11-11T14:05:00.123456Z","title":null,"body":"Hei","truncated":false,"full_sha256":null,"answer_source":"codex.item_completed","agent_label":"quality","provider":"codex","write_kind":"live"}`
	if req, err := DecodePutRequest([]byte(valid)); err != nil || req.Title != nil || req.Body != "Hei" {
		t.Fatalf("valid request: %+v %v", req, err)
	}
	cases := map[string]string{
		"unknown key":     strings.Replace(valid, `"body":"Hei"`, `"body":"Hei","extra":1`, 1),
		"missing title":   strings.Replace(valid, `"title":null,`, ``, 1),
		"duplicate key":   strings.Replace(valid, `"body":"Hei"`, `"body":"Hei","body":"Hei"`, 1),
		"trailing data":   valid + ` {}`,
		"null body":       strings.Replace(valid, `"body":"Hei"`, `"body":null`, 1),
		"number title":    strings.Replace(valid, `"title":null`, `"title":1`, 1),
		"string boolean":  strings.Replace(valid, `"truncated":false`, `"truncated":"false"`, 1),
		"not an object":   `[]`,
		"empty":           ``,
		"lone surrogate":  strings.Replace(valid, `"body":"Hei"`, `"body":"\ud800"`, 1),
		"oversize":        strings.Replace(valid, `"body":"Hei"`, `"body":"`+strings.Repeat("x", maxRequestBytes)+`"`, 1),
		"invalid utf8":    strings.Replace(valid, `"body":"Hei"`, "\"body\":\"\xff\"", 1),
		"wrong kind type": strings.Replace(valid, `"kind":"InboxPutRequest@1"`, `"kind":null`, 1),
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := DecodePutRequest([]byte(input))
			inboxError := expectClass(t, err, ClassInvalidArguments)
			if strings.Contains(inboxError.Detail, "Hei") {
				t.Fatalf("detail leaks body: %q", inboxError.Detail)
			}
		})
	}
}
