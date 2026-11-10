// Package assignment stores a Project's Assignments as one JSON file each,
// <directory>/ASSIGNMENT-<n>.json. An Assignment is written once by Add and
// never changed. Add numbers under an exclusive lock on <directory>/.lock;
// reads take none. Errors are *Error with a fixed Class and a detail that
// never contains absolute paths or the submitted text.
package assignment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// DataPathCategory is the data_paths category of a Project's Assignments.
const DataPathCategory = "assignments"

// Error classes shared with the CLI contract.
const (
	ClassInvalidArguments = "assignment_invalid_arguments"
	ClassNotFound         = "assignment_not_found"
	ClassInvalid          = "assignment_invalid"
	ClassIOError          = "assignment_io_error"
)

// Error is one `<class>: <detail>` failure.
type Error struct{ Class, Detail string }

func (e *Error) Error() string { return e.Class + ": " + e.Detail }

// Assignment is the exact content of one Assignment file.
type Assignment struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
	Repos     []string `json:"repos"`
	Topics    []string `json:"topics,omitempty"`
}

// Summary is one Assignment as shown by List, without its text.
type Summary struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	CreatedAt string   `json:"created_at"`
	Repos     []string `json:"repos"`
	Topics    []string `json:"topics,omitempty"`
}

// InvalidFile is a *.json file that cannot be read as an Assignment.
type InvalidFile struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

// ListResult holds readable Assignments and unreadable files; both are never nil.
type ListResult struct {
	Items   []Summary     `json:"items"`
	Invalid []InvalidFile `json:"invalid"`
}

// Store reads and writes the Assignments of one directory.
type Store struct {
	dir     string
	now     func() time.Time
	replace func(string, []byte, fs.FileMode) error
	claim   func() (string, error)
	syncDir func(string) error
	lock    filelock.Options
}

// Open returns the store for directory, the absolute Assignment directory
// itself (not a project root). Nothing is created.
func Open(directory string) Store {
	return Store{dir: directory, now: time.Now, replace: atomicfile.Replace, claim: newClaim, syncDir: syncDirectory}
}

const (
	lockName      = ".lock"
	fileSuffix    = ".json"
	maxNumber     = 999999999
	maxTitleRunes = 200
	maxTextBytes  = 64 * 1024
	maxFileBytes  = 1 << 20
	maxTags       = 32
	maxTopicRunes = 100
	timeLayout    = "2006-01-02T15:04:05Z"
)

var (
	idPattern   = regexp.MustCompile(`^ASSIGNMENT-([1-9][0-9]{0,8})$`)
	filePattern = regexp.MustCompile(`^ASSIGNMENT-([1-9][0-9]{0,8})\.json$`)
	repoPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,99}$`)
)

// Add validates everything before any effect, then stores a new Assignment
// numbered one above the highest ASSIGNMENT-<n>.json name in the directory.
func (s Store) Add(title, text string, repos, topics []string) (Assignment, error) {
	if !filepath.IsAbs(s.dir) {
		return Assignment{}, invalidArguments("assignment directory must be absolute")
	}
	title = strings.TrimSpace(title)
	if !validTitle(title) {
		return Assignment{}, invalidArguments("title must be one non-empty line of at most 200 characters")
	}
	if !validText(text) {
		return Assignment{}, invalidArguments("text must be non-blank UTF-8 of at most 64 KiB")
	}
	repos = trimAll(repos)
	if len(repos) == 0 || !validTags(repos, validRepo) {
		return Assignment{}, invalidArguments("repos must be 1 to 32 distinct lowercase ASCII slugs of at most 100 characters")
	}
	topics = trimAll(topics)
	if !validTags(topics, validTopic) {
		return Assignment{}, invalidArguments("topics must be at most 32 distinct non-blank lines of at most 100 characters")
	}
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return Assignment{}, ioError("create assignments directory", err)
	}
	lock, err := filelock.Acquire(filepath.Join(s.dir, lockName), s.lock)
	if err != nil {
		return Assignment{}, &Error{Class: ClassIOError, Detail: "lock assignments"}
	}
	assignment, err := s.addLocked(title, text, repos, topics)
	if releaseErr := lock.Release(); releaseErr != nil && err == nil {
		return Assignment{}, ioError("release assignments lock after writing "+assignment.ID, releaseErr)
	}
	return assignment, err
}

func (s Store) addLocked(title, text string, repos, topics []string) (Assignment, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return Assignment{}, ioError("read assignments directory", err)
	}
	highest := 0
	for _, entry := range entries {
		if match := filePattern.FindStringSubmatch(entry.Name()); match != nil {
			if n, _ := strconv.Atoi(match[1]); n > highest {
				highest = n
			}
		}
	}
	if highest >= maxNumber {
		return Assignment{}, &Error{Class: ClassIOError, Detail: "id space exhausted"}
	}
	assignment := Assignment{
		ID:        fmt.Sprintf("ASSIGNMENT-%d", highest+1),
		Title:     title,
		Text:      text,
		CreatedAt: s.now().UTC().Format(timeLayout),
		Repos:     repos,
		Topics:    topics,
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(assignment); err != nil {
		return Assignment{}, ioError("encode "+assignment.ID, err)
	}
	if err := s.replace(s.path(assignment.ID), buffer.Bytes(), 0o644); err != nil {
		return Assignment{}, ioError("write "+assignment.ID, err)
	}
	return assignment, nil
}

// Read returns one Assignment. It takes no lock and creates nothing.
func (s Store) Read(id string) (Assignment, error) {
	if !idPattern.MatchString(id) {
		return Assignment{}, invalidArguments("invalid assignment id")
	}
	return s.read(id)
}

// List returns every readable Assignment by ascending number and every
// unreadable *.json file by name. One invalid file does not stop the listing.
func (s Store) List() (ListResult, error) {
	result := ListResult{Items: []Summary{}, Invalid: []InvalidFile{}}
	entries, err := os.ReadDir(s.dir)
	if errors.Is(err, fs.ErrNotExist) {
		return result, nil
	}
	if err != nil {
		return result, ioError("read assignments directory", err)
	}
	numbers := map[string]int{}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") || !strings.HasSuffix(name, fileSuffix) {
			continue
		}
		match := filePattern.FindStringSubmatch(name)
		if match == nil {
			result.Invalid = append(result.Invalid, InvalidFile{File: name, Error: "invalid file name"})
			continue
		}
		a, err := s.read(strings.TrimSuffix(name, fileSuffix))
		if err != nil {
			result.Invalid = append(result.Invalid, InvalidFile{File: name, Error: detail(err)})
			continue
		}
		numbers[a.ID], _ = strconv.Atoi(match[1])
		result.Items = append(result.Items, Summary{ID: a.ID, Title: a.Title, CreatedAt: a.CreatedAt, Repos: a.Repos, Topics: a.Topics})
	}
	sort.Slice(result.Items, func(i, j int) bool { return numbers[result.Items[i].ID] < numbers[result.Items[j].ID] })
	sort.Slice(result.Invalid, func(i, j int) bool { return result.Invalid[i].File < result.Invalid[j].File })
	return result, nil
}

func (s Store) path(id string) string { return filepath.Join(s.dir, id+fileSuffix) }

func (s Store) read(id string) (Assignment, error) {
	info, err := os.Lstat(s.path(id))
	if errors.Is(err, fs.ErrNotExist) {
		return Assignment{}, &Error{Class: ClassNotFound, Detail: id}
	}
	if err != nil {
		return Assignment{}, ioError("inspect "+id, err)
	}
	if !info.Mode().IsRegular() {
		return Assignment{}, invalid("not a regular file")
	}
	if info.Size() > maxFileBytes {
		return Assignment{}, invalid("file too large")
	}
	file, err := os.Open(s.path(id))
	if errors.Is(err, fs.ErrNotExist) {
		return Assignment{}, &Error{Class: ClassNotFound, Detail: id}
	}
	if err != nil {
		return Assignment{}, ioError("read "+id, err)
	}
	defer file.Close()
	if opened, err := file.Stat(); err != nil || !os.SameFile(info, opened) {
		return Assignment{}, invalid("not a regular file")
	}
	contents, err := io.ReadAll(io.LimitReader(file, maxFileBytes+1))
	if err != nil {
		return Assignment{}, ioError("read "+id, err)
	}
	if len(contents) > maxFileBytes {
		return Assignment{}, invalid("file too large")
	}
	return decode(id, contents)
}

// fileFields are the exact top-level keys of an Assignment file.
var fileFields = map[string]bool{"id": true, "title": true, "text": true, "created_at": true, "repos": true, "topics": true}

// checkKeys rejects keys that encoding/json would accept by case folding or
// by letting the last duplicate win, so no content is silently replaced.
func checkKeys(contents []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(contents))
	if token, err := decoder.Token(); err != nil || token != json.Delim('{') {
		return invalid("invalid JSON")
	}
	seen := map[string]bool{}
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return invalid("invalid JSON")
		}
		key := token.(string)
		switch {
		case !fileFields[key]:
			return invalid("unknown field")
		case seen[key]:
			return invalid("duplicate field")
		}
		seen[key] = true
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return invalid("invalid JSON")
		}
	}
	return nil
}

func decode(id string, contents []byte) (Assignment, error) {
	if err := checkKeys(contents); err != nil {
		return Assignment{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	var a Assignment
	if err := decoder.Decode(&a); err != nil {
		if strings.HasPrefix(err.Error(), "json: unknown field ") {
			return Assignment{}, invalid("unknown field")
		}
		return Assignment{}, invalid("invalid JSON")
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return Assignment{}, invalid("data after JSON value")
	}
	switch {
	case a.ID == "":
		return Assignment{}, invalid("missing id")
	case a.ID != id:
		return Assignment{}, invalid("id does not match file name")
	case a.Title == "":
		return Assignment{}, invalid("missing title")
	case a.Text == "":
		return Assignment{}, invalid("missing text")
	case a.CreatedAt == "":
		return Assignment{}, invalid("missing created_at")
	case a.Repos == nil:
		return Assignment{}, invalid("missing repos")
	case a.Title != strings.TrimSpace(a.Title) || !validTitle(a.Title):
		return Assignment{}, invalid("invalid title")
	case !validText(a.Text):
		return Assignment{}, invalid("invalid text")
	case !validCreatedAt(a.CreatedAt):
		return Assignment{}, invalid("invalid created_at")
	case len(a.Repos) == 0 || !validTags(a.Repos, validRepo):
		return Assignment{}, invalid("invalid repos")
	case !validTags(a.Topics, validTopic):
		return Assignment{}, invalid("invalid topics")
	}
	if len(a.Topics) == 0 {
		a.Topics = nil
	}
	return a, nil
}

func validTitle(title string) bool {
	return title != "" && !strings.ContainsAny(title, "\r\n") && utf8.ValidString(title) && utf8.RuneCountInString(title) <= maxTitleRunes
}

func validText(text string) bool {
	return utf8.ValidString(text) && len(text) <= maxTextBytes && strings.TrimSpace(text) != ""
}

func validCreatedAt(value string) bool {
	at, err := time.Parse(time.RFC3339, value)
	return err == nil && at.UTC().Format(timeLayout) == value
}

func trimAll(values []string) []string {
	var out []string
	for _, value := range values {
		out = append(out, strings.TrimSpace(value))
	}
	return out
}

// validTags allows at most 32 valid values without case-insensitive duplicates.
func validTags(values []string, valid func(string) bool) bool {
	for i, value := range values {
		if !valid(value) || slices.ContainsFunc(values[:i], func(kept string) bool { return strings.EqualFold(kept, value) }) {
			return false
		}
	}
	return len(values) <= maxTags
}

func validRepo(value string) bool { return repoPattern.MatchString(value) }

func validTopic(value string) bool {
	return value != "" && value == strings.TrimSpace(value) && utf8.ValidString(value) &&
		utf8.RuneCountInString(value) <= maxTopicRunes && strings.IndexFunc(value, unicode.IsControl) < 0
}

func invalidArguments(detail string) error {
	return &Error{Class: ClassInvalidArguments, Detail: detail}
}
func invalid(detail string) error { return &Error{Class: ClassInvalid, Detail: detail} }

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
	var e *Error
	if errors.As(err, &e) {
		return e.Detail
	}
	return err.Error()
}
