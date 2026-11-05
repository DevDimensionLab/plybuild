// Package backlog stores a project backlog as one JSON file per item under
// <projectRoot>/backlog/. Item text and notes are kept byte for byte; notes only
// grow. Add, AddNote and SetStatus hold an exclusive lock on backlog/.lock; reads take none.
// Errors are *Error with a fixed Class and a detail that never contains
// absolute paths or the submitted text.
package backlog

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
)

// Error classes shared with the CLI contract.
const (
	ClassInvalidArguments = "backlog_invalid_arguments"
	ClassItemNotFound     = "backlog_item_not_found"
	ClassItemInvalid      = "backlog_item_invalid"
	ClassIOError          = "backlog_io_error"
)

// Error is one `<class>: <detail>` failure.
type Error struct {
	Class  string
	Detail string
}

func (e *Error) Error() string { return e.Class + ": " + e.Detail }

// Note is one appended addition to an item.
type Note struct {
	At   string `json:"at"`
	Text string `json:"text"`
}

// Status is the stored state of an item.
type Status string

// The only valid statuses; any transition between them is allowed.
const (
	StatusOpen  Status = "open"
	StatusReady Status = "ready"
	StatusDone  Status = "done"
)

func (s Status) valid() bool { return s == StatusOpen || s == StatusReady || s == StatusDone }

// Item is the exact content of one item file.
type Item struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	Status    Status `json:"status"`
	Notes     []Note `json:"notes"`
}

// Summary is one item as shown by List.
type Summary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	NoteCount int    `json:"note_count"`
	Status    Status `json:"status"`
}

// InvalidFile is a *.json file in backlog/ that cannot be read as an item.
type InvalidFile struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

// ListResult holds readable items and unreadable files; both are never nil.
type ListResult struct {
	Items   []Summary     `json:"items"`
	Invalid []InvalidFile `json:"invalid"`
}

// Store reads and writes the backlog of one project root.
type Store struct {
	root   string
	now    func() time.Time
	random io.Reader
}

// Open returns the store for projectRoot, which must be an absolute, existing directory.
func Open(projectRoot string) Store {
	return Store{root: projectRoot, now: time.Now, random: rand.Reader}
}

const (
	directoryName  = "backlog"
	lockName       = ".lock"
	fileSuffix     = ".json"
	maxTitleRunes  = 200
	maxTextBytes   = 64 * 1024
	maxIDAttempts  = 5
	randomIDBytes  = 4
	timestampStyle = time.RFC3339
)

var itemID = regexp.MustCompile(`^[0-9]{8}-[0-9a-f]{8}$`)

// List returns all readable items, newest first, and every unreadable *.json file.
func (s Store) List() (ListResult, error) {
	result := ListResult{Items: []Summary{}, Invalid: []InvalidFile{}}
	if err := s.checkRoot(); err != nil {
		return result, err
	}
	entries, err := os.ReadDir(s.directory())
	if errors.Is(err, fs.ErrNotExist) {
		return result, nil
	}
	if err != nil {
		return result, ioError("read backlog directory", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, fileSuffix) {
			continue
		}
		item, err := s.read(strings.TrimSuffix(name, fileSuffix))
		if err != nil {
			result.Invalid = append(result.Invalid, InvalidFile{File: name, Error: detail(err)})
			continue
		}
		updated := item.CreatedAt
		if len(item.Notes) > 0 {
			updated = item.Notes[len(item.Notes)-1].At
		}
		result.Items = append(result.Items, Summary{ID: item.ID, Title: item.Title, CreatedAt: item.CreatedAt, UpdatedAt: updated, NoteCount: len(item.Notes), Status: item.Status})
	}
	sort.Slice(result.Items, func(i, j int) bool {
		a, b := result.Items[i], result.Items[j]
		if a.CreatedAt != b.CreatedAt {
			return a.CreatedAt > b.CreatedAt
		}
		return a.ID < b.ID
	})
	return result, nil
}

// Read returns the item with id.
func (s Store) Read(id string) (Item, error) {
	if err := s.checkRoot(); err != nil {
		return Item{}, err
	}
	if !itemID.MatchString(id) {
		return Item{}, invalidArguments("invalid item id")
	}
	return s.read(id)
}

// Add creates a new item with a trimmed one-line title and verbatim text.
func (s Store) Add(title, text string) (Item, error) {
	if err := s.checkRoot(); err != nil {
		return Item{}, err
	}
	title = strings.TrimSpace(title)
	if title == "" || strings.ContainsAny(title, "\r\n") || !utf8.ValidString(title) || utf8.RuneCountInString(title) > maxTitleRunes {
		return Item{}, invalidArguments("title must be one non-empty line of at most 200 characters")
	}
	if err := checkText(text); err != nil {
		return Item{}, err
	}
	if err := os.MkdirAll(s.directory(), 0o755); err != nil {
		return Item{}, ioError("create backlog directory", err)
	}
	unlock, err := s.lock()
	if err != nil {
		return Item{}, err
	}
	defer unlock()
	now := s.now().UTC()
	for attempt := 0; attempt < maxIDAttempts; attempt++ {
		random := make([]byte, randomIDBytes)
		if _, err := io.ReadFull(s.random, random); err != nil {
			return Item{}, ioError("generate item id", err)
		}
		id := now.Format("20060102") + "-" + hex.EncodeToString(random)
		if _, err := os.Lstat(s.path(id)); err == nil {
			continue
		} else if !errors.Is(err, fs.ErrNotExist) {
			return Item{}, ioError("check item "+id, err)
		}
		item := Item{ID: id, Title: title, Text: text, CreatedAt: now.Format(timestampStyle), Status: StatusOpen, Notes: []Note{}}
		return item, s.write(item)
	}
	return Item{}, &Error{Class: ClassIOError, Detail: "could not allocate a unique item id"}
}

// AddNote appends a verbatim note to the item with id and returns the whole item.
func (s Store) AddNote(id, text string) (Item, error) {
	if err := s.checkRoot(); err != nil {
		return Item{}, err
	}
	if !itemID.MatchString(id) {
		return Item{}, invalidArguments("invalid item id")
	}
	if err := checkText(text); err != nil {
		return Item{}, err
	}
	if _, err := os.Stat(s.directory()); errors.Is(err, fs.ErrNotExist) {
		return Item{}, notFound(id)
	}
	unlock, err := s.lock()
	if err != nil {
		return Item{}, err
	}
	defer unlock()
	item, err := s.read(id)
	if err != nil {
		return Item{}, err
	}
	item.Notes = append(item.Notes, Note{At: s.now().UTC().Format(timestampStyle), Text: text})
	return item, s.write(item)
}

// SetStatus sets the item's status and returns the whole item.
// When the item already has that status, no item file is written.
func (s Store) SetStatus(id string, status Status) (Item, error) {
	if err := s.checkRoot(); err != nil {
		return Item{}, err
	}
	if !itemID.MatchString(id) {
		return Item{}, invalidArguments("invalid item id")
	}
	if !status.valid() {
		return Item{}, invalidArguments("status must be open, ready or done")
	}
	if _, err := os.Stat(s.directory()); errors.Is(err, fs.ErrNotExist) {
		return Item{}, notFound(id)
	}
	unlock, err := s.lock()
	if err != nil {
		return Item{}, err
	}
	defer unlock()
	item, err := s.read(id)
	if err != nil || item.Status == status {
		return item, err
	}
	item.Status = status
	return item, s.write(item)
}

func (s Store) checkRoot() error {
	if !filepath.IsAbs(s.root) {
		return invalidArguments("project root must be an absolute path")
	}
	if info, err := os.Stat(s.root); err != nil || !info.IsDir() {
		return invalidArguments("project root must be an existing directory")
	}
	return nil
}

func (s Store) directory() string     { return filepath.Join(s.root, directoryName) }
func (s Store) path(id string) string { return filepath.Join(s.directory(), id+fileSuffix) }

func (s Store) read(id string) (Item, error) {
	contents, err := os.ReadFile(s.path(id))
	if errors.Is(err, fs.ErrNotExist) {
		return Item{}, notFound(id)
	}
	if err != nil {
		return Item{}, ioError("read item "+id, err)
	}
	var item Item
	if err := json.Unmarshal(contents, &item); err != nil {
		var syntax *json.SyntaxError
		if errors.As(err, &syntax) {
			return Item{}, itemInvalid(fmt.Sprintf("invalid JSON at offset %d", syntax.Offset))
		}
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return Item{}, itemInvalid(fmt.Sprintf("invalid JSON at offset %d", len(contents)))
		}
		return Item{}, itemInvalid("invalid item JSON")
	}
	switch {
	case item.ID == "":
		return Item{}, itemInvalid("missing id")
	case item.Title == "":
		return Item{}, itemInvalid("missing title")
	case item.Text == "":
		return Item{}, itemInvalid("missing text")
	case item.ID != id:
		return Item{}, itemInvalid("id does not match file name")
	case item.Status == "":
		return Item{}, itemInvalid("missing status")
	case !item.Status.valid():
		return Item{}, itemInvalid("invalid status")
	}
	if item.Notes == nil {
		item.Notes = []Note{}
	}
	return item, nil
}

func (s Store) write(item Item) error {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(item); err != nil {
		return ioError("encode item "+item.ID, err)
	}
	if err := atomicfile.Replace(s.path(item.ID), buffer.Bytes(), 0o644); err != nil {
		return ioError("write item "+item.ID, err)
	}
	return nil
}

func (s Store) lock() (func(), error) {
	file, err := os.OpenFile(filepath.Join(s.directory(), lockName), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, ioError("open backlog lock", err)
	}
	if err := lockFile(file); err != nil {
		_ = file.Close()
		return nil, ioError("lock backlog", err)
	}
	return func() { _ = unlockFile(file); _ = file.Close() }, nil
}

func checkText(text string) error {
	if !utf8.ValidString(text) || len(text) > maxTextBytes || strings.TrimSpace(text) == "" {
		return invalidArguments("text must be non-blank UTF-8 of at most 64 KiB")
	}
	return nil
}

func invalidArguments(detail string) error {
	return &Error{Class: ClassInvalidArguments, Detail: detail}
}
func itemInvalid(detail string) error { return &Error{Class: ClassItemInvalid, Detail: detail} }
func notFound(id string) error        { return &Error{Class: ClassItemNotFound, Detail: id} }

// ioError keeps only the operating-system cause so details never carry absolute paths.
func ioError(operation string, err error) error {
	cause := "unexpected error"
	var pathError *fs.PathError
	var linkError *os.LinkError
	var errno syscall.Errno
	switch {
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

func detail(err error) string {
	var backlogError *Error
	if errors.As(err, &backlogError) {
		return backlogError.Detail
	}
	return err.Error()
}
