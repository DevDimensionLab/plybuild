package inbox

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

// PutRequest is InboxPutRequest@1. Project and owner come from the store.
type PutRequest struct {
	Kind           string  `json:"kind"`
	ConversationID string  `json:"conversation_id"`
	MessageID      string  `json:"message_id"`
	RepliedAtUTC   string  `json:"replied_at_utc"`
	Title          *string `json:"title"`
	Body           string  `json:"body"`
	Truncated      bool    `json:"truncated"`
	FullSHA256     *string `json:"full_sha256"`
	AnswerSource   string  `json:"answer_source"`
	AgentLabel     string  `json:"agent_label"`
	Provider       string  `json:"provider"`
	WriteKind      string  `json:"write_kind"`
}

const (
	maxRequestBytes = 512 * 1024
	maxBodyRunes    = 32000
	maxTitleRunes   = 60
	maxLabelRunes   = 64
	canonicalLayout = "2006-01-02T15:04:05.000000000Z"
	compactLayout   = "20060102T150405.000000000Z"
)

var (
	conversationPattern = regexp.MustCompile(`^conv_[0-9a-f]{32}$`)
	messagePattern      = regexp.MustCompile(`^msg_[0-9a-f]{32}$`)
	timePattern         = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]{1,9})?Z$`)
	canonicalPattern    = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{9}Z$`)
	sourcePattern       = regexp.MustCompile(`^[a-z][a-z0-9_.-]{0,63}$`)
	providerPattern     = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
	digestPattern       = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

var putRequestKeys = []string{"kind", "conversation_id", "message_id", "replied_at_utc", "title", "body", "truncated", "full_sha256", "answer_source", "agent_label", "provider", "write_kind"}

// DecodePutRequest reads one closed InboxPutRequest@1 of at most 512 KiB.
// Field rules are checked by Put.
func DecodePutRequest(raw []byte) (PutRequest, error) {
	if len(raw) > maxRequestBytes {
		return PutRequest{}, invalidArguments("request exceeds 512 KiB")
	}
	f, err := closedObject(raw, putRequestKeys)
	if err != nil {
		return PutRequest{}, invalidArguments("request: " + err.Error())
	}
	var req PutRequest
	ok := f.str("kind", &req.Kind) && f.str("conversation_id", &req.ConversationID) && f.str("message_id", &req.MessageID) &&
		f.str("replied_at_utc", &req.RepliedAtUTC) && f.nullableStr("title", &req.Title) && f.str("body", &req.Body) &&
		f.boolean("truncated", &req.Truncated) && f.nullableStr("full_sha256", &req.FullSHA256) && f.str("answer_source", &req.AnswerSource) &&
		f.str("agent_label", &req.AgentLabel) && f.str("provider", &req.Provider) && f.str("write_kind", &req.WriteKind)
	if !ok {
		return PutRequest{}, invalidArguments("request: " + f.problem)
	}
	return req, nil
}

// payload holds the immutable fields of one reply in canonical form.
type payload struct {
	conversationID, messageID, repliedAtUTC string
	title                                   *string
	body                                    string
	truncated                               bool
	fullSHA256                              *string
	answerSource                            string
}

// payloadFieldOrder is the order conflicting_fields are reported in.
var payloadFieldOrder = []string{"replied_at_utc", "title", "body", "truncated", "full_sha256", "answer_source"}

func (p payload) values() map[string]canonicaljson.Value {
	return map[string]canonicaljson.Value{
		"replied_at_utc": p.repliedAtUTC, "title": nullable(p.title), "body": p.body,
		"truncated": p.truncated, "full_sha256": nullable(p.fullSHA256), "answer_source": p.answerSource,
	}
}

func nullable(value *string) canonicaljson.Value {
	if value == nil {
		return nil
	}
	return *value
}

// payloadSHA256 is "sha256:" + hex SHA-256 of the RFC 8785 InboxPayload@1.
func payloadSHA256(project, owner string, p payload) (string, error) {
	object := canonicaljson.Object{
		{Name: "kind", Value: PayloadKind}, {Name: "project_id", Value: project}, {Name: "owner", Value: owner},
		{Name: "conversation_id", Value: p.conversationID}, {Name: "message_id", Value: p.messageID},
	}
	values := p.values()
	for _, name := range payloadFieldOrder {
		object = append(object, canonicaljson.Member{Name: name, Value: values[name]})
	}
	encoded, err := canonicaljson.Marshal(object)
	if err != nil {
		return "", invalidArguments("payload is not canonical JSON")
	}
	digest := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func payloadOf(m Message) payload {
	return payload{m.ConversationID, m.MessageID, m.RepliedAtUTC, m.Title, m.Body, m.Truncated, m.FullSHA256, m.AnswerSource}
}

// conflictingFields names the payload fields in which a and b differ.
func conflictingFields(a, b payload) []string {
	left, right := a.values(), b.values()
	fields := []string{}
	for _, name := range payloadFieldOrder {
		if left[name] != right[name] {
			fields = append(fields, name)
		}
	}
	return fields
}

// checkPayload validates the request's immutable fields and returns them in
// canonical form: the reply time with nine fraction digits and the trimmed title.
func checkPayload(req PutRequest) (payload, error) {
	if !conversationPattern.MatchString(req.ConversationID) {
		return payload{}, invalidArguments("conversation_id must match conv_ and 32 lowercase hex")
	}
	if !messagePattern.MatchString(req.MessageID) {
		return payload{}, invalidArguments("message_id must match msg_ and 32 lowercase hex")
	}
	replied, err := canonicalTime(req.RepliedAtUTC)
	if err != nil {
		return payload{}, err
	}
	title := req.Title
	if title != nil {
		// Only surrounding spaces are trimmed; control characters and line breaks anywhere are rejected.
		trimmed := strings.TrimSpace(*title)
		if !oneLine(*title) || !validTitle(trimmed) {
			return payload{}, invalidArguments("title must be null or one line of 1-60 code points without control or format characters")
		}
		title = &trimmed
	}
	if !utf8.ValidString(req.Body) || utf8.RuneCountInString(req.Body) > maxBodyRunes {
		return payload{}, invalidArguments("body must be UTF-8 of at most 32000 code points")
	}
	switch {
	case !req.Truncated && req.FullSHA256 != nil:
		return payload{}, invalidArguments("full_sha256 must be null unless truncated")
	case req.Truncated && (req.FullSHA256 == nil || !digestPattern.MatchString(*req.FullSHA256)):
		return payload{}, invalidArguments("full_sha256 must be 64 lowercase hex when truncated")
	case req.Truncated && *req.FullSHA256 == bodySHA256(req.Body):
		return payload{}, invalidArguments("full_sha256 must differ from the body digest when truncated")
	}
	if !sourcePattern.MatchString(req.AnswerSource) {
		return payload{}, invalidArguments("invalid answer_source")
	}
	return payload{req.ConversationID, req.MessageID, replied, title, req.Body, req.Truncated, req.FullSHA256, req.AnswerSource}, nil
}

// checkMetadata validates the first-write presentation fields.
func checkMetadata(label, provider, writeKind string) error {
	switch {
	case label == "" || utf8.RuneCountInString(label) > maxLabelRunes || !oneLine(label):
		return invalidArguments("agent_label must be one line of 1-64 code points without control characters")
	case !providerPattern.MatchString(provider):
		return invalidArguments("invalid provider")
	case writeKind != WriteKindLive && writeKind != WriteKindImport:
		return invalidArguments("write_kind must be live or import")
	}
	return nil
}

func validTitle(title string) bool {
	if title == "" || utf8.RuneCountInString(title) > maxTitleRunes || !utf8.ValidString(title) {
		return false
	}
	for _, r := range title {
		if unicode.Is(unicode.Cf, r) {
			return false
		}
	}
	return oneLine(title)
}

// oneLine rejects invalid UTF-8, control characters and line or paragraph separators.
func oneLine(value string) bool {
	if !utf8.ValidString(value) {
		return false
	}
	for _, r := range value {
		if unicode.Is(unicode.Cc, r) || r == ' ' || r == ' ' {
			return false
		}
	}
	return true
}

// canonicalTime parses an RFC 3339 UTC time with a Z offset and 0-9 fraction
// digits and returns it with exactly nine.
func canonicalTime(value string) (string, error) {
	if !timePattern.MatchString(value) {
		return "", invalidArguments("replied_at_utc must be RFC 3339 with Z and at most 9 fraction digits")
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return "", invalidArguments("replied_at_utc is not a valid time")
	}
	return parsed.UTC().Format(canonicalLayout), nil
}

// compactTime turns a canonical time into the file name form; it assumes canonical input.
func compactTime(canonical string) string {
	parsed, err := time.Parse(canonicalLayout, canonical)
	if err != nil {
		return ""
	}
	return parsed.Format(compactLayout)
}

func validCanonicalTime(value string) bool {
	if !canonicalPattern.MatchString(value) {
		return false
	}
	canonical, err := canonicalTime(value)
	return err == nil && canonical == value
}

func bodySHA256(body string) string {
	digest := sha256.Sum256([]byte(body))
	return hex.EncodeToString(digest[:])
}

// fields reads typed members of one closed JSON object.
type fields struct {
	values  map[string]canonicaljson.Value
	problem string
}

// closedObject strictly decodes raw as one object with exactly keys.
func closedObject(raw []byte, keys []string) (fields, error) {
	value, err := canonicaljson.DecodeStrict(raw)
	if err != nil {
		return fields{}, errString("invalid JSON")
	}
	object, ok := value.(canonicaljson.Object)
	if !ok {
		return fields{}, errString("not a JSON object")
	}
	values := map[string]canonicaljson.Value{}
	for _, member := range object {
		values[member.Name] = member.Value
	}
	for _, key := range keys {
		if _, present := values[key]; !present {
			return fields{}, errString("missing " + key)
		}
	}
	if len(values) != len(keys) {
		return fields{}, errString("unknown field")
	}
	return fields{values: values}, nil
}

type errString string

func (e errString) Error() string { return string(e) }

func (f *fields) str(name string, target *string) bool {
	value, ok := f.values[name].(string)
	if !ok {
		f.problem = name + " must be a string"
		return false
	}
	*target = value
	return true
}

func (f *fields) nullableStr(name string, target **string) bool {
	switch value := f.values[name].(type) {
	case nil:
		*target = nil
		return true
	case string:
		*target = &value
		return true
	}
	f.problem = name + " must be a string or null"
	return false
}

func (f *fields) boolean(name string, target *bool) bool {
	value, ok := f.values[name].(bool)
	if !ok {
		f.problem = name + " must be a boolean"
		return false
	}
	*target = value
	return true
}

func (f *fields) integer(name string, want int64) bool {
	value, ok := f.values[name].(int64)
	if !ok || value != want {
		f.problem = "unsupported " + name
		return false
	}
	return true
}
