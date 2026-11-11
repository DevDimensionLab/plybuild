package inbox

import (
	"encoding/base64"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

var messageKeys = []string{"kind", "schema_version", "project_id", "owner", "conversation_id", "message_id", "replied_at_utc", "title", "derived_title",
	"derived_title_rule", "body", "body_sha256", "truncated", "full_sha256", "answer_source", "agent_label", "provider", "write_kind", "recorded_at_utc", "payload_sha256"}

// decodeMessageFile strictly decodes the message file name in project and
// owner: closed keys, field rules, name, scope, derived values and payload hash.
func decodeMessageFile(contents []byte, name, project, owner string) (Message, error) {
	f, err := closedObject(contents, messageKeys)
	if err != nil {
		return Message{}, itemInvalid(err.Error())
	}
	m := Message{SchemaVersion: SchemaVersion}
	ok := f.str("kind", &m.Kind) && f.integer("schema_version", SchemaVersion) && f.str("project_id", &m.ProjectID) && f.str("owner", &m.Owner) &&
		f.str("conversation_id", &m.ConversationID) && f.str("message_id", &m.MessageID) && f.str("replied_at_utc", &m.RepliedAtUTC) &&
		f.nullableStr("title", &m.Title) && f.str("derived_title", &m.DerivedTitle) && f.str("derived_title_rule", &m.DerivedTitleRule) &&
		f.str("body", &m.Body) && f.str("body_sha256", &m.BodySHA256) && f.boolean("truncated", &m.Truncated) && f.nullableStr("full_sha256", &m.FullSHA256) &&
		f.str("answer_source", &m.AnswerSource) && f.str("agent_label", &m.AgentLabel) && f.str("provider", &m.Provider) &&
		f.str("write_kind", &m.WriteKind) && f.str("recorded_at_utc", &m.RecordedAtUTC) && f.str("payload_sha256", &m.PayloadSHA256)
	if !ok {
		return Message{}, itemInvalid(f.problem)
	}
	switch {
	case m.Kind != MessageKind:
		return Message{}, itemInvalid("unsupported kind")
	case m.ProjectID != project || m.Owner != owner:
		return Message{}, itemInvalid("project or owner does not match the store")
	case !validCanonicalTime(m.RepliedAtUTC) || !validCanonicalTime(m.RecordedAtUTC):
		return Message{}, itemInvalid("times must be canonical")
	case name != messageKey(m)+fileSuffix:
		return Message{}, itemInvalid("content does not match file name")
	}
	p, err := checkPayload(PutRequest{ConversationID: m.ConversationID, MessageID: m.MessageID, RepliedAtUTC: m.RepliedAtUTC, Title: m.Title,
		Body: m.Body, Truncated: m.Truncated, FullSHA256: m.FullSHA256, AnswerSource: m.AnswerSource})
	if err != nil {
		return Message{}, itemInvalid("stored payload: " + err.(*Error).Detail)
	}
	if err := checkMetadata(m.AgentLabel, m.Provider, m.WriteKind); err != nil {
		return Message{}, itemInvalid("stored metadata: " + err.(*Error).Detail)
	}
	switch {
	case m.Title != nil && *p.title != *m.Title:
		return Message{}, itemInvalid("title is not trimmed")
	case m.DerivedTitleRule != DerivedTitleRule || m.DerivedTitle != DeriveTitle(m.Body):
		return Message{}, itemInvalid("derived title does not match its rule")
	case m.BodySHA256 != bodySHA256(m.Body):
		return Message{}, itemInvalid("body_sha256 does not match the body")
	}
	if hash, err := payloadSHA256(project, owner, p); err != nil || hash != m.PayloadSHA256 {
		return Message{}, itemInvalid("payload_sha256 does not match the payload")
	}
	return m, nil
}

var markKeys = []string{"kind", "schema_version", "owner", "conversation_id", "message_id", "marked_at_utc"}

// decodeMarkFile strictly decodes the read mark file name of owner.
func decodeMarkFile(contents []byte, name, owner string) (ReadMark, error) {
	f, err := closedObject(contents, markKeys)
	if err != nil {
		return ReadMark{}, itemInvalid("read mark: " + err.Error())
	}
	mark := ReadMark{SchemaVersion: SchemaVersion}
	ok := f.str("kind", &mark.Kind) && f.integer("schema_version", SchemaVersion) && f.str("owner", &mark.Owner) &&
		f.str("conversation_id", &mark.ConversationID) && f.str("message_id", &mark.MessageID) && f.str("marked_at_utc", &mark.MarkedAtUTC)
	switch {
	case !ok:
		return ReadMark{}, itemInvalid("read mark: " + f.problem)
	case mark.Kind != ReadMarkKind:
		return ReadMark{}, itemInvalid("read mark: unsupported kind")
	case mark.Owner != owner:
		return ReadMark{}, itemInvalid("read mark: owner does not match the store")
	case checkIdentity(mark.ConversationID, mark.MessageID) != nil || name != identity(mark.ConversationID, mark.MessageID)+fileSuffix:
		return ReadMark{}, itemInvalid("read mark: content does not match file name")
	case !validCanonicalTime(mark.MarkedAtUTC):
		return ReadMark{}, itemInvalid("read mark: marked_at_utc must be canonical")
	}
	return mark, nil
}

var cursorKeys = []string{"version", "project_id", "owner", "key"}

// encodeCursor is base64url without padding of the closed cursor object.
func (s *Store) encodeCursor(key string) string {
	encoded, _ := canonicaljson.Marshal(canonicaljson.Object{
		{Name: "version", Value: int64(1)}, {Name: "project_id", Value: s.project}, {Name: "owner", Value: s.owner}, {Name: "key", Value: key},
	})
	return base64.RawURLEncoding.EncodeToString(encoded)
}

// decodeCursor validates the cursor and its scope and returns its key.
func (s *Store) decodeCursor(cursor string) (string, error) {
	invalid := invalidArguments("--cursor is not a next_cursor of this project and owner")
	raw, err := base64.RawURLEncoding.Strict().DecodeString(cursor)
	if err != nil {
		return "", invalid
	}
	f, err := closedObject(raw, cursorKeys)
	if err != nil {
		return "", invalid
	}
	var project, owner, key string
	if !f.integer("version", 1) || !f.str("project_id", &project) || !f.str("owner", &owner) || !f.str("key", &key) {
		return "", invalid
	}
	if project != s.project || owner != s.owner {
		return "", invalid
	}
	if match := messageNamePattern.FindStringSubmatch(key + fileSuffix); match == nil || !validCompactTime(match[1]) {
		return "", invalid
	}
	return key, nil
}
