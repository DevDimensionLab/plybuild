// Package inbox stores the canonical agent replies of one project owner and
// that owner's read marks. Each reply is one write-once JSON file whose name
// carries the full sort key; a read mark is a separate file whose presence
// means read. Put and Mark hold an exclusive lock on <root>/.lock for at most
// two seconds; reads take none and never create anything.
// Errors are *Error with a fixed Class and a detail that never contains
// absolute paths, body or title text.
package inbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"syscall"
	"time"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

// Error classes shared with the CLI contract.
const (
	ClassInvalidArguments = "inbox_invalid_arguments"
	ClassConflict         = "inbox_conflict"
	ClassItemNotFound     = "inbox_item_not_found"
	ClassItemInvalid      = "inbox_item_invalid"
	ClassBusy             = "inbox_busy"
	ClassIOError          = "inbox_io_error"
)

// Contract kinds, all at SchemaVersion.
const (
	SchemaVersion  = 1
	Mode           = "inbox_v1"
	PutRequestKind = "InboxPutRequest@1"
	PayloadKind    = "InboxPayload@1"
	MessageKind    = "InboxMessage@1"
	ReadMarkKind   = "InboxReadMark@1"
	PutResultKind  = "InboxPutResult@1"
	ItemKind       = "InboxItem@1"
	ListKind       = "InboxList@1"
	MarkResultKind = "InboxMarkResult@1"
	ErrorKind      = "InboxError@1"
)

// Put and Mark effects.
const (
	EffectCreated   = "created"
	EffectUnchanged = "unchanged"
	EffectMarked    = "marked"
	EffectNone      = "none"
)

// Write kinds are provenance only.
const (
	WriteKindLive   = "live"
	WriteKindImport = "import"
)

// Title sources of an item's display_title.
const (
	TitleSourceAgent   = "agent"
	TitleSourceDerived = "derived"
)

// List read state and unread count basis.
const (
	ReadStateAvailable      = "available"
	ReadStateUnavailable    = "unavailable"
	BasisValidatedItems     = "validated_items"
	BasisFilenameUnverified = "filename_unverified"
	BasisUnavailable        = "unavailable"
)

// MarkValue is the requested read state.
type MarkValue string

const (
	MarkRead   MarkValue = "read"
	MarkUnread MarkValue = "unread"
)

// Error is one `<class>: <detail>` failure. The conflict fields are set only
// for ClassConflict and name fields and hashes, never their values.
type Error struct {
	Class                  string
	Detail                 string
	ConflictingFields      []string
	ExistingPayloadSHA256  string
	SubmittedPayloadSHA256 string
}

func (e *Error) Error() string { return e.Class + ": " + e.Detail }

// Retryable reports whether the same request may succeed later.
func (e *Error) Retryable() bool { return e.Class == ClassBusy || e.Class == ClassIOError }

// ErrorResult is InboxError@1.
type ErrorResult struct {
	Kind                   string   `json:"kind"`
	SchemaVersion          int      `json:"schema_version"`
	Class                  string   `json:"class"`
	Detail                 string   `json:"detail"`
	Retryable              bool     `json:"retryable"`
	ConflictingFields      []string `json:"conflicting_fields,omitempty"`
	ExistingPayloadSHA256  string   `json:"existing_payload_sha256,omitempty"`
	SubmittedPayloadSHA256 string   `json:"submitted_payload_sha256,omitempty"`
}

// NewErrorResult is the InboxError@1 of err.
func NewErrorResult(err *Error) ErrorResult {
	return ErrorResult{ErrorKind, SchemaVersion, err.Class, err.Detail, err.Retryable(), err.ConflictingFields, err.ExistingPayloadSHA256, err.SubmittedPayloadSHA256}
}

// Message is InboxMessage@1, the exact content of one message file.
type Message struct {
	Kind             string  `json:"kind"`
	SchemaVersion    int     `json:"schema_version"`
	ProjectID        string  `json:"project_id"`
	Owner            string  `json:"owner"`
	ConversationID   string  `json:"conversation_id"`
	MessageID        string  `json:"message_id"`
	RepliedAtUTC     string  `json:"replied_at_utc"`
	Title            *string `json:"title"`
	DerivedTitle     string  `json:"derived_title"`
	DerivedTitleRule string  `json:"derived_title_rule"`
	Body             string  `json:"body"`
	BodySHA256       string  `json:"body_sha256"`
	Truncated        bool    `json:"truncated"`
	FullSHA256       *string `json:"full_sha256"`
	AnswerSource     string  `json:"answer_source"`
	AgentLabel       string  `json:"agent_label"`
	Provider         string  `json:"provider"`
	WriteKind        string  `json:"write_kind"`
	RecordedAtUTC    string  `json:"recorded_at_utc"`
	PayloadSHA256    string  `json:"payload_sha256"`
}

// Item is InboxItem@1: the message fields plus the shown title and read
// state. Read is nil when the mark state is unreadable.
type Item struct {
	Kind             string  `json:"kind"`
	SchemaVersion    int     `json:"schema_version"`
	ProjectID        string  `json:"project_id"`
	Owner            string  `json:"owner"`
	ConversationID   string  `json:"conversation_id"`
	MessageID        string  `json:"message_id"`
	RepliedAtUTC     string  `json:"replied_at_utc"`
	Title            *string `json:"title"`
	DerivedTitle     string  `json:"derived_title"`
	DerivedTitleRule string  `json:"derived_title_rule"`
	TitleSource      string  `json:"title_source"`
	DisplayTitle     string  `json:"display_title"`
	Body             string  `json:"body"`
	BodySHA256       string  `json:"body_sha256"`
	Truncated        bool    `json:"truncated"`
	FullSHA256       *string `json:"full_sha256"`
	AnswerSource     string  `json:"answer_source"`
	AgentLabel       string  `json:"agent_label"`
	Provider         string  `json:"provider"`
	WriteKind        string  `json:"write_kind"`
	RecordedAtUTC    string  `json:"recorded_at_utc"`
	PayloadSHA256    string  `json:"payload_sha256"`
	Read             *bool   `json:"read"`
}

func newItem(m Message, read *bool) Item {
	source, display := TitleSourceDerived, m.DerivedTitle
	if m.Title != nil {
		source, display = TitleSourceAgent, *m.Title
	}
	return Item{ItemKind, SchemaVersion, m.ProjectID, m.Owner, m.ConversationID, m.MessageID, m.RepliedAtUTC, m.Title, m.DerivedTitle, m.DerivedTitleRule,
		source, display, m.Body, m.BodySHA256, m.Truncated, m.FullSHA256, m.AnswerSource, m.AgentLabel, m.Provider, m.WriteKind, m.RecordedAtUTC, m.PayloadSHA256, read}
}

// ReadMark is InboxReadMark@1; its presence means read.
type ReadMark struct {
	Kind           string `json:"kind"`
	SchemaVersion  int    `json:"schema_version"`
	Owner          string `json:"owner"`
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	MarkedAtUTC    string `json:"marked_at_utc"`
}

// PutResult is InboxPutResult@1.
type PutResult struct {
	Kind          string  `json:"kind"`
	SchemaVersion int     `json:"schema_version"`
	Effect        string  `json:"effect"`
	Item          Message `json:"item"`
}

// MarkResult is InboxMarkResult@1.
type MarkResult struct {
	Kind           string `json:"kind"`
	SchemaVersion  int    `json:"schema_version"`
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	Read           bool   `json:"read"`
	Effect         string `json:"effect"`
}

// InvalidEntry is a directory entry or page item that could not be used.
type InvalidEntry struct {
	Entry string `json:"entry"`
	Error string `json:"error"`
}

// List is InboxList@1. UnreadCount is nil when read marks are unavailable.
type List struct {
	Kind             string         `json:"kind"`
	SchemaVersion    int            `json:"schema_version"`
	ProjectID        string         `json:"project_id"`
	Owner            string         `json:"owner"`
	Limit            int            `json:"limit"`
	Cursor           *string        `json:"cursor"`
	NextCursor       *string        `json:"next_cursor"`
	TotalCount       int            `json:"total_count"`
	UnreadCount      *int           `json:"unread_count"`
	UnreadCountExact bool           `json:"unread_count_exact"`
	UnreadCountBasis string         `json:"unread_count_basis"`
	ReadState        string         `json:"read_state"`
	Items            []Item         `json:"items"`
	Invalid          []InvalidEntry `json:"invalid"`
}

// Store is the Inbox of one owner in one project.
type Store struct {
	root        string
	project     string
	owner       string
	now         func() time.Time
	replace     func(string, []byte, fs.FileMode) error
	lockTimeout time.Duration
}

const (
	lockName           = ".lock"
	fileSuffix         = ".json"
	defaultLockTimeout = 2 * time.Second
	maxFileBytes       = 1 << 20
	dirMode            = 0o700
	fileMode           = 0o600
)

var ownerPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,64}$`)

// Open returns the store of owner in project below the absolute Inbox root.
// Nothing is created until Put or Mark writes.
func Open(root, project, owner string) (*Store, error) {
	if !filepath.IsAbs(root) {
		return nil, invalidArguments("inbox directory must be an absolute path")
	}
	if err := CheckOwner(owner); err != nil {
		return nil, err
	}
	if project == "" || !oneLine(project) {
		return nil, invalidArguments("invalid project ID")
	}
	return &Store{root: root, project: project, owner: owner, now: time.Now, replace: atomicfile.Replace, lockTimeout: defaultLockTimeout}, nil
}

// CheckOwner accepts an opaque owner slug of 1-64 [A-Za-z0-9_.-] other than . and ..
func CheckOwner(owner string) error {
	if !ownerPattern.MatchString(owner) || owner == "." || owner == ".." {
		return invalidArguments("owner must match [A-Za-z0-9_.-]{1,64} and not be . or ..")
	}
	return nil
}

// encodeJSON is the encoding of files and JSON results.
func encodeJSON(value any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Encode writes value as a JSON result with the file encoding.
func Encode(out io.Writer, value any) error {
	contents, err := encodeJSON(value)
	if err != nil {
		return err
	}
	_, err = out.Write(contents)
	return err
}

func invalidArguments(detail string) *Error {
	return &Error{Class: ClassInvalidArguments, Detail: detail}
}
func itemInvalid(detail string) *Error { return &Error{Class: ClassItemInvalid, Detail: detail} }
func notFound(conversation, message string) *Error {
	return &Error{Class: ClassItemNotFound, Detail: conversation + "/" + message}
}

// ioError keeps only the operating-system cause so details never carry paths.
func ioError(operation string, err error) *Error {
	cause := "unexpected error"
	var pathError *fs.PathError
	var linkError *os.LinkError
	var atomicError *atomicfile.Error
	var errno syscall.Errno
	switch {
	case errors.As(err, &atomicError) && atomicError.Replaced:
		cause = "effect uncertain after " + atomicError.Op
	case errors.As(err, &pathError):
		cause = pathError.Err.Error()
	case errors.As(err, &linkError):
		cause = linkError.Err.Error()
	case errors.As(err, &errno):
		cause = errno.Error()
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
		cause = err.Error()
	}
	return &Error{Class: ClassIOError, Detail: operation + ": " + cause}
}
