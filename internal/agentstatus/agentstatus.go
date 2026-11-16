// Package agentstatus stores the status an agent reports about itself, one
// JSON file per actor in a Project's directory. A post always expires; reads
// compute expiry and change nothing. Set replaces the actor's file atomically
// under <directory>/.lock and removes other posts expired for over 24 hours.
// Errors are *Error with a fixed Class and a detail without paths or input.
package agentstatus

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// DataPathCategory is the data directory category of a Project's agent statuses.
const DataPathCategory = "agent-status"

// Error classes shared with the CLI contract.
const (
	ClassInvalidArguments = "agent_status_invalid_arguments"
	ClassIOError          = "agent_status_io_error"
)

// Error is one `<class>: <detail>` failure.
type Error struct{ Class, Detail string }

func (e *Error) Error() string { return e.Class + ": " + e.Detail }

// States are the reportable states with the longest duration of each.
var States = []string{"working", "waiting_background", "needs_answer", "done"}

var maxDuration = map[string]time.Duration{
	"working":            4 * time.Hour,
	"waiting_background": 4 * time.Hour,
	"needs_answer":       24 * time.Hour,
	"done":               4 * time.Hour,
}

// Item is one stored post; Expired is computed at the read.
type Item struct {
	Actor      string `json:"actor"`
	Pane       string `json:"pane"`
	State      string `json:"state"`
	Ref        string `json:"ref"`
	ReportedAt string `json:"reported_at"`
	ExpiresAt  string `json:"expires_at"`
	Expired    bool   `json:"expired"`
}

// InvalidFile is a file that cannot be read; File is only its name.
type InvalidFile struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

// ListResult holds every readable post by actor and every unreadable file.
type ListResult struct {
	ObservedAt string        `json:"observed_at"`
	Items      []Item        `json:"items"`
	Invalid    []InvalidFile `json:"invalid"`
}

// Store reads and writes the posts of one directory.
type Store struct {
	dir string
	now func() time.Time
}

// Open returns the store for directory, the absolute agent-status directory
// itself. Nothing is created until Set.
func Open(directory string) Store {
	return Store{dir: directory, now: time.Now}
}

const (
	lockName      = ".lock"
	fileSuffix    = ".json"
	schemaVersion = 1
	maxFileBytes  = 4096
	pruneAfter    = 24 * time.Hour
	timeLayout    = "2006-01-02T15:04:05Z"
)

var (
	actorPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{0,19}:[A-Za-z0-9._-]{1,80}$`)
	panePattern  = regexp.MustCompile(`^w[A-Za-z0-9]{1,20}:p[A-Za-z0-9]{1,20}$`)
	refPattern   = regexp.MustCompile(`^[A-Za-z0-9._:/-]{1,80}$`)
	filePattern  = regexp.MustCompile(`^[0-9a-f]{32}\.json$`)
)

// document is the exact content of one post file.
type document struct {
	SchemaVersion int    `json:"schema_version"`
	Actor         string `json:"actor"`
	Pane          string `json:"pane"`
	State         string `json:"state"`
	Ref           string `json:"ref"`
	ReportedAt    string `json:"reported_at"`
	ExpiresAt     string `json:"expires_at"`
}

// Set replaces the post of actor, expiring after duration, and returns it.
func (s Store) Set(actor, pane, state, ref string, duration time.Duration) (Item, error) {
	if err := s.checkDirectory(); err != nil {
		return Item{}, err
	}
	if err := checkActor(actor); err != nil {
		return Item{}, err
	}
	if !panePattern.MatchString(pane) {
		return Item{}, invalidArguments("pane must be a Herdr pane such as w7:pT")
	}
	limit, known := maxDuration[state]
	if !known {
		return Item{}, invalidArguments("state must be one of " + strings.Join(States, ", "))
	}
	if ref != "" && !refPattern.MatchString(ref) {
		return Item{}, invalidArguments("ref must be an ID of at most 80 letters, digits or . _ : / -")
	}
	if duration <= 0 || duration > limit {
		return Item{}, invalidArguments(fmt.Sprintf("duration for %s must be above zero and at most %s", state, limit))
	}
	now := s.now().UTC().Truncate(time.Second)
	d := document{SchemaVersion: schemaVersion, Actor: actor, Pane: pane, State: state, Ref: ref,
		ReportedAt: now.Format(timeLayout), ExpiresAt: now.Add(duration).Format(timeLayout)}
	contents, err := encode(d)
	if err != nil {
		return Item{}, ioError("encode post", err)
	}
	err = s.locked(func() error {
		s.prune(now, actor)
		return atomicfile.Replace(s.path(actor), contents, 0o644)
	})
	if err != nil {
		return Item{}, err
	}
	return d.item(now), nil
}

// Clear removes the post of actor and reports whether one was removed.
func (s Store) Clear(actor string) (bool, error) {
	if err := s.checkDirectory(); err != nil {
		return false, err
	}
	if err := checkActor(actor); err != nil {
		return false, err
	}
	if _, err := os.Lstat(s.dir); errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	removed := false
	err := s.locked(func() error {
		err := os.Remove(s.path(actor))
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		removed = err == nil
		return err
	})
	return removed, err
}

// List returns every readable post sorted by actor; it takes no lock.
func (s Store) List() (ListResult, error) {
	now := s.now().UTC().Truncate(time.Second)
	result := ListResult{ObservedAt: now.Format(timeLayout), Items: []Item{}, Invalid: []InvalidFile{}}
	documents, invalid, err := s.readAll()
	if err != nil {
		return result, err
	}
	for _, d := range documents {
		result.Items = append(result.Items, d.item(now))
	}
	result.Invalid = append(result.Invalid, invalid...)
	sort.Slice(result.Items, func(i, j int) bool { return result.Items[i].Actor < result.Items[j].Actor })
	sort.Slice(result.Invalid, func(i, j int) bool { return result.Invalid[i].File < result.Invalid[j].File })
	return result, nil
}

func (d document) item(now time.Time) Item {
	expires, _ := time.Parse(timeLayout, d.ExpiresAt)
	return Item{Actor: d.Actor, Pane: d.Pane, State: d.State, Ref: d.Ref, ReportedAt: d.ReportedAt,
		ExpiresAt: d.ExpiresAt, Expired: !now.Before(expires)}
}

// prune removes posts of other actors expired for over 24 hours. A failure
// leaves the post for the next Set; it never stops the caller's own report.
func (s Store) prune(now time.Time, actor string) {
	documents, _, err := s.readAll()
	if err != nil {
		return
	}
	for _, d := range documents {
		expires, _ := time.Parse(timeLayout, d.ExpiresAt)
		if d.Actor != actor && now.Sub(expires) > pruneAfter {
			_ = os.Remove(s.path(d.Actor))
		}
	}
}

func (s Store) readAll() ([]document, []InvalidFile, error) {
	entries, err := os.ReadDir(s.dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, ioError("read agent-status directory", err)
	}
	var documents []document
	var invalid []InvalidFile
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") || !strings.HasSuffix(name, fileSuffix) {
			continue
		}
		d, detail := s.read(name)
		if detail != "" {
			invalid = append(invalid, InvalidFile{File: name, Error: detail})
			continue
		}
		documents = append(documents, d)
	}
	return documents, invalid, nil
}

func (s Store) read(name string) (document, string) {
	if !filePattern.MatchString(name) {
		return document{}, "invalid file name"
	}
	path := filepath.Join(s.dir, name)
	info, err := os.Lstat(path)
	if err != nil {
		return document{}, "cannot inspect file"
	}
	if !info.Mode().IsRegular() {
		return document{}, "not a regular file"
	}
	if info.Size() > maxFileBytes {
		return document{}, "file too large"
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return document{}, "cannot read file"
	}
	var d document
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&d) != nil || decoder.More() {
		return document{}, "invalid JSON"
	}
	if d.SchemaVersion != schemaVersion {
		return document{}, "unknown schema_version"
	}
	_, known := maxDuration[d.State]
	_, reportedErr := time.Parse(timeLayout, d.ReportedAt)
	_, expiresErr := time.Parse(timeLayout, d.ExpiresAt)
	if checkActor(d.Actor) != nil || !panePattern.MatchString(d.Pane) || !known ||
		(d.Ref != "" && !refPattern.MatchString(d.Ref)) || reportedErr != nil || expiresErr != nil {
		return document{}, "invalid post"
	}
	if fileName(d.Actor) != name {
		return document{}, "actor does not match file name"
	}
	return d, ""
}

func (s Store) checkDirectory() error {
	if !filepath.IsAbs(s.dir) {
		return invalidArguments("agent-status directory must be absolute")
	}
	return nil
}

// locked runs write under the directory lock, creating the directory first.
func (s Store) locked(write func() error) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return ioError("create agent-status directory", err)
	}
	lock, err := filelock.Acquire(filepath.Join(s.dir, lockName), filelock.Options{})
	if err != nil {
		return ioError("lock agent-status", err)
	}
	err = write()
	if releaseErr := lock.Release(); releaseErr != nil && err == nil {
		return &Error{Class: ClassIOError, Detail: "the post was written; releasing the agent-status lock failed"}
	}
	if err != nil {
		return ioError("write post", err)
	}
	return nil
}

func checkActor(actor string) error {
	if !actorPattern.MatchString(actor) {
		return invalidArguments("actor must be <provider>:<id> such as claude:<session-id> or codex:<hex>")
	}
	return nil
}

// fileName keeps actor characters such as ':' out of file names.
func fileName(actor string) string {
	sum := sha256.Sum256([]byte(actor))
	return hex.EncodeToString(sum[:16]) + fileSuffix
}

func (s Store) path(actor string) string { return filepath.Join(s.dir, fileName(actor)) }

func encode(d document) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(d); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func invalidArguments(detail string) error {
	return &Error{Class: ClassInvalidArguments, Detail: detail}
}

// ioError keeps only the error kind, never the path inside err.
func ioError(action string, err error) error {
	kind := "failed"
	switch {
	case errors.Is(err, fs.ErrPermission):
		kind = "permission denied"
	case errors.Is(err, fs.ErrNotExist):
		kind = "not found"
	}
	return &Error{Class: ClassIOError, Detail: action + ": " + kind}
}
