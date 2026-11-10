package assignment

// The lifecycle of an Assignment is an append-only event list in
// <directory>/.lifecycle/ASSIGNMENT-<n>.json. State and revision are derived
// by replay; the original Assignment file is never written. Transitions run
// under the same <directory>/.lock as Add; reads take no lock. Done is
// terminal, and the archive is only the view of done Assignments.

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
	"runtime"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// Lifecycle error classes shared with the CLI contract.
const (
	ClassConflict         = "assignment_conflict"
	ClassLifecycleInvalid = "assignment_lifecycle_invalid"
)

// State is derived from the lifecycle events; it is never stored.
type State string

const (
	StateAvailable State = "available"
	StateTaken     State = "taken"
	StateDone      State = "done"
)

// View selects the Assignments ListView returns. ViewArchive is the done ones.
type View string

const (
	ViewAvailable View = "available"
	ViewTaken     View = "taken"
	ViewArchive   View = "archive"
	ViewAll       View = "all"
)

// Views are all View values, the default first.
var Views = []View{ViewAvailable, ViewTaken, ViewArchive, ViewAll}

// stateViews is the one view besides ViewAll that shows each state.
var stateViews = map[State]View{StateAvailable: ViewAvailable, StateTaken: ViewTaken, StateDone: ViewArchive}

// EventKind is the transition one Event records.
type EventKind string

const (
	EventTaken    EventKind = "taken"
	EventReleased EventKind = "released"
	EventRevoked  EventKind = "revoked"
	EventDone     EventKind = "done"
)

// Event is one recorded transition. Claim is set for taken, released and
// done; Reason only for revoked.
type Event struct {
	Seq    int       `json:"seq"`
	At     string    `json:"at"`
	Kind   EventKind `json:"kind"`
	Actor  string    `json:"actor"`
	Claim  string    `json:"claim,omitempty"`
	Reason string    `json:"reason,omitempty"`
}

// Holder is the open claim of a taken Assignment.
type Holder struct {
	Actor   string `json:"actor"`
	Claim   string `json:"claim"`
	TakenAt string `json:"taken_at"`
}

// DoneBy says who marked a done Assignment done, and when.
type DoneBy struct {
	Actor string `json:"actor"`
	At    string `json:"at"`
}

// Lifecycle is the replayed record: Holder only when taken, Done only when done.
type Lifecycle struct {
	State    State   `json:"state"`
	Revision int     `json:"revision"`
	Holder   *Holder `json:"holder"`
	Done     *DoneBy `json:"done"`
	Events   []Event `json:"events"`
}

// ViewItem is one Assignment in a view, without text or claim.
type ViewItem struct {
	Summary
	State    State  `json:"state"`
	Revision int    `json:"revision"`
	TakenAt  string `json:"taken_at,omitempty"`
	DoneAt   string `json:"done_at,omitempty"`
	Actor    string `json:"actor,omitempty"`
}

// ViewResult holds one view and every unreadable file; both lists are never
// nil. A lifecycle entry's Error is "<class>: <detail>", with class
// assignment_lifecycle_invalid or, for an unreadable record, assignment_io_error.
type ViewResult struct {
	View    View          `json:"view"`
	Items   []ViewItem    `json:"items"`
	Invalid []InvalidFile `json:"invalid"`
}

type record struct {
	Format string            `json:"format"`
	ID     string            `json:"id"`
	Events []json.RawMessage `json:"events"`
}

const (
	lifecycleDir    = ".lifecycle"
	lifecycleFormat = "ply-assignment-lifecycle/v1"
	maxEvents       = 64
	maxActorRunes   = 100
	maxReasonRunes  = 200
	uncertainSuffix = "; result unknown, run show"
)

var claimPattern = regexp.MustCompile(`^c_[0-9a-f]{32}$`)

// recordFields and eventFields are the exact keys of a record and of each
// event kind; nothing else, such as a Task or backlog reference, is accepted.
var (
	recordFields = map[string]bool{"format": true, "id": true, "events": true}
	eventFields  = map[EventKind]map[string]bool{
		EventTaken:    {"seq": true, "at": true, "kind": true, "actor": true, "claim": true},
		EventReleased: {"seq": true, "at": true, "kind": true, "actor": true, "claim": true},
		EventDone:     {"seq": true, "at": true, "kind": true, "actor": true, "claim": true},
		EventRevoked:  {"seq": true, "at": true, "kind": true, "actor": true, "reason": true},
	}
	anyEventField = map[string]bool{"seq": true, "at": true, "kind": true, "actor": true, "claim": true, "reason": true}
)

func newClaim() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return "c_" + hex.EncodeToString(raw[:]), nil
}

// Lifecycle returns the replayed lifecycle of one valid Assignment. It takes
// no lock and creates nothing.
func (s Store) Lifecycle(id string) (Lifecycle, error) {
	if !idPattern.MatchString(id) {
		return Lifecycle{}, invalidArguments("invalid assignment id")
	}
	if _, err := s.read(id); err != nil {
		return Lifecycle{}, err
	}
	l, _, err := s.readLifecycle(id)
	return l, err
}

// Take gives the available Assignment one new claim for actor, if revision
// is still the current one.
func (s Store) Take(id string, revision int, actor string) (Lifecycle, error) {
	if err := checkTransition(s, id, revision); err != nil {
		return Lifecycle{}, err
	}
	if !validLine(actor, maxActorRunes) {
		return Lifecycle{}, invalidArguments("actor must be one line of 1 to 100 characters without control characters")
	}
	claim, err := s.claim()
	if err != nil {
		return Lifecycle{}, &Error{Class: ClassIOError, Detail: "generate claim"}
	}
	return s.transition(id, func(l Lifecycle) (*Event, error) {
		if err := requireAvailable(id, l, revision); err != nil {
			return nil, err
		}
		if l.Revision+2 > maxEvents {
			return nil, conflict("lifecycle event limit reached")
		}
		return &Event{Kind: EventTaken, Actor: actor, Claim: claim}, nil
	})
}

// Release gives a taken Assignment back unfinished. Repeating it right after
// it succeeded writes nothing.
func (s Store) Release(id, claim string) (Lifecycle, error) {
	return s.closeClaim(id, claim, EventReleased)
}

// Done marks a taken Assignment done, which puts it in the archive view.
// Repeating it writes nothing.
func (s Store) Done(id, claim string) (Lifecycle, error) {
	return s.closeClaim(id, claim, EventDone)
}

// Revoke makes a taken Assignment available again without its claim, if
// revision is still the current one. It is the recovery for a lost holder.
func (s Store) Revoke(id string, revision int, actor, reason string) (Lifecycle, error) {
	if err := checkTransition(s, id, revision); err != nil {
		return Lifecycle{}, err
	}
	if !validLine(actor, maxActorRunes) {
		return Lifecycle{}, invalidArguments("actor must be one line of 1 to 100 characters without control characters")
	}
	if !validLine(reason, maxReasonRunes) {
		return Lifecycle{}, invalidArguments("reason must be one line of 1 to 200 characters without control characters")
	}
	return s.transition(id, func(l Lifecycle) (*Event, error) {
		switch {
		case l.State == StateDone:
			return nil, conflict(id + " is done")
		case l.State != StateTaken:
			return nil, conflict(id + " is not taken")
		case l.Revision != revision:
			return nil, staleRevision(revision, l.Revision)
		}
		return &Event{Kind: EventRevoked, Actor: actor, Reason: reason}, nil
	})
}

func (s Store) closeClaim(id, claim string, kind EventKind) (Lifecycle, error) {
	if err := checkTransition(s, id, 0); err != nil {
		return Lifecycle{}, err
	}
	if !claimPattern.MatchString(claim) {
		return Lifecycle{}, invalidArguments("claim must be c_ followed by 32 lowercase hex digits")
	}
	return s.transition(id, func(l Lifecycle) (*Event, error) {
		if n := len(l.Events); n > 0 && l.Events[n-1].Kind == kind && l.Events[n-1].Claim == claim {
			return nil, nil
		}
		switch {
		case l.State == StateDone:
			return nil, conflict(id + " is done")
		case l.State != StateTaken || l.Holder.Claim != claim:
			return nil, conflict("claim does not hold " + id)
		}
		return &Event{Kind: kind, Actor: l.Holder.Actor, Claim: claim}, nil
	})
}

func checkTransition(s Store, id string, revision int) error {
	switch {
	case !filepath.IsAbs(s.dir):
		return invalidArguments("assignment directory must be absolute")
	case !idPattern.MatchString(id):
		return invalidArguments("invalid assignment id")
	case revision < 0:
		return invalidArguments("revision must be a non-negative integer")
	}
	return nil
}

func requireAvailable(id string, l Lifecycle, revision int) error {
	switch {
	case l.State != StateAvailable:
		return conflict(id + " is " + string(l.State))
	case l.Revision != revision:
		return staleRevision(revision, l.Revision)
	}
	return nil
}

// transition runs one critical section under the Assignment lock: read and
// validate, let decide pick the event (nil for an idempotent repeat), append
// it and replace the record atomically.
func (s Store) transition(id string, decide func(Lifecycle) (*Event, error)) (Lifecycle, error) {
	if _, err := os.Lstat(s.dir); errors.Is(err, fs.ErrNotExist) {
		return Lifecycle{}, &Error{Class: ClassNotFound, Detail: id}
	}
	lock, err := filelock.Acquire(filepath.Join(s.dir, lockName), s.lock)
	if err != nil {
		return Lifecycle{}, &Error{Class: ClassIOError, Detail: "lock assignments"}
	}
	l, wrote, err := s.transitionLocked(id, decide)
	if releaseErr := lock.Release(); releaseErr != nil && err == nil {
		if wrote {
			return Lifecycle{}, &Error{Class: ClassIOError, Detail: "release assignments lock after writing lifecycle " + id + uncertainSuffix}
		}
		return Lifecycle{}, ioError("release assignments lock", releaseErr)
	}
	return l, err
}

func (s Store) transitionLocked(id string, decide func(Lifecycle) (*Event, error)) (Lifecycle, bool, error) {
	if _, err := s.read(id); err != nil {
		return Lifecycle{}, false, err
	}
	l, raw, err := s.readLifecycle(id)
	if err != nil {
		return Lifecycle{}, false, err
	}
	next, err := decide(l)
	if err != nil || next == nil {
		return l, false, err
	}
	next.Seq = l.Revision + 1
	next.At = s.now().UTC().Format(timeLayout)
	if n := len(l.Events); n > 0 && next.At < l.Events[n-1].At {
		next.At = l.Events[n-1].At // a clock behind the record must not break its order
	}
	after, err := replay(append(slices.Clone(l.Events), *next))
	if err != nil {
		return Lifecycle{}, false, &Error{Class: ClassIOError, Detail: "new lifecycle event does not replay"}
	}
	encodedEvent, err := json.Marshal(next)
	if err != nil {
		return Lifecycle{}, false, ioError("encode lifecycle "+id, err)
	}
	contents, err := encodeRecord(record{Format: lifecycleFormat, ID: id, Events: append(raw, encodedEvent)})
	if err != nil {
		return Lifecycle{}, false, ioError("encode lifecycle "+id, err)
	}
	if err := s.ensureLifecycleDir(); err != nil {
		return Lifecycle{}, false, err
	}
	if err := s.replace(s.recordPath(id), contents, 0o644); err != nil {
		var replaceError *atomicfile.Error
		if errors.As(err, &replaceError) && (replaceError.Replaced || replaceError.Op == "replace") {
			uncertain := ioError("write lifecycle "+id, err).(*Error)
			uncertain.Detail += uncertainSuffix
			return Lifecycle{}, false, uncertain
		}
		return Lifecycle{}, false, ioError("write lifecycle "+id, err)
	}
	return after, true, nil
}

func encodeRecord(r record) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(r); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s Store) ensureLifecycleDir() error {
	path := filepath.Join(s.dir, lifecycleDir)
	if err := os.Mkdir(path, 0o755); err != nil && !errors.Is(err, fs.ErrExist) {
		return ioError("create lifecycle directory", err)
	}
	if info, err := os.Lstat(path); err != nil || !info.IsDir() {
		return lifecycleInvalid("not a directory")
	}
	// The record's rename is synced by atomicfile; the directory entry of
	// .lifecycle itself is synced here, on every transition, so a failed
	// earlier sync is retried.
	if err := s.syncDir(s.dir); err != nil {
		return ioError("sync assignments directory", err)
	}
	return nil
}

// syncDirectory makes a directory's entries durable, as atomicfile does after
// a replace; like atomicfile it is a no-op on Windows.
func syncDirectory(path string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	directory, err := os.Open(path)
	if err != nil {
		return err
	}
	if err := directory.Sync(); err != nil {
		_ = directory.Close()
		return err
	}
	return directory.Close()
}

func (s Store) recordPath(id string) string {
	return filepath.Join(s.dir, lifecycleDir, id+fileSuffix)
}

// readLifecycle reads and replays one record; absent is available at
// revision 0. It also returns the raw events, for appending.
func (s Store) readLifecycle(id string) (Lifecycle, []json.RawMessage, error) {
	empty := Lifecycle{State: StateAvailable, Events: []Event{}}
	info, err := os.Lstat(filepath.Join(s.dir, lifecycleDir))
	if errors.Is(err, fs.ErrNotExist) {
		return empty, nil, nil
	}
	if err != nil {
		return Lifecycle{}, nil, ioError("inspect lifecycle directory", err)
	}
	if !info.IsDir() {
		return Lifecycle{}, nil, lifecycleInvalid("not a directory")
	}
	contents, err := readLimited(s.recordPath(id), "lifecycle "+id, lifecycleInvalid)
	if errors.Is(err, fs.ErrNotExist) {
		return empty, nil, nil
	}
	if err != nil {
		return Lifecycle{}, nil, err
	}
	return decodeRecord(id, contents)
}

// readLimited reads a regular file of at most 1 MiB without following a
// symlink; fs.ErrNotExist is returned as is.
func readLimited(path, name string, invalid func(string) error) ([]byte, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fs.ErrNotExist
	}
	if err != nil {
		return nil, ioError("inspect "+name, err)
	}
	if !info.Mode().IsRegular() {
		return nil, invalid("not a regular file")
	}
	if info.Size() > maxFileBytes {
		return nil, invalid("file too large")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, ioError("read "+name, err)
	}
	defer file.Close()
	if opened, err := file.Stat(); err != nil || !os.SameFile(info, opened) {
		return nil, invalid("not a regular file")
	}
	contents, err := io.ReadAll(io.LimitReader(file, maxFileBytes+1))
	if err != nil {
		return nil, ioError("read "+name, err)
	}
	if len(contents) > maxFileBytes {
		return nil, invalid("file too large")
	}
	return contents, nil
}

func decodeRecord(id string, contents []byte) (Lifecycle, []json.RawMessage, error) {
	var r record
	if err := decodeStrict(contents, recordFields, &r); err != nil {
		return Lifecycle{}, nil, err
	}
	switch {
	case r.Format != lifecycleFormat:
		return Lifecycle{}, nil, lifecycleInvalid("invalid format")
	case r.ID != id:
		return Lifecycle{}, nil, lifecycleInvalid("id does not match file name")
	case len(r.Events) == 0 || len(r.Events) > maxEvents:
		return Lifecycle{}, nil, lifecycleInvalid("invalid events")
	}
	events := make([]Event, 0, len(r.Events))
	for _, raw := range r.Events {
		var e Event
		if err := decodeStrict(raw, anyEventField, &e); err != nil {
			return Lifecycle{}, nil, err
		}
		if !validEvent(raw, e) {
			return Lifecycle{}, nil, lifecycleInvalid("invalid event")
		}
		events = append(events, e)
	}
	l, err := replay(events)
	return l, r.Events, err
}

// decodeStrict decodes one JSON object with exactly known, non-duplicate keys
// and nothing after it.
func decodeStrict(contents []byte, allowed map[string]bool, into any) error {
	keys := json.NewDecoder(bytes.NewReader(contents))
	if token, err := keys.Token(); err != nil || token != json.Delim('{') {
		return lifecycleInvalid("invalid JSON")
	}
	seen := map[string]bool{}
	for keys.More() {
		token, err := keys.Token()
		if err != nil {
			return lifecycleInvalid("invalid JSON")
		}
		key := token.(string)
		switch {
		case !allowed[key]:
			return lifecycleInvalid("unknown field")
		case seen[key]:
			return lifecycleInvalid("duplicate field")
		}
		seen[key] = true
		var value json.RawMessage
		if err := keys.Decode(&value); err != nil {
			return lifecycleInvalid("invalid JSON")
		}
	}
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(into); err != nil {
		return lifecycleInvalid("invalid JSON")
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return lifecycleInvalid("data after JSON value")
	}
	return nil
}

// validEvent checks one event's own keys and values; replay checks the order.
func validEvent(raw json.RawMessage, e Event) bool {
	fields, known := eventFields[e.Kind]
	if !known {
		return false
	}
	var present map[string]json.RawMessage
	if json.Unmarshal(raw, &present) != nil || len(present) != len(fields) {
		return false
	}
	for key := range present {
		if !fields[key] {
			return false
		}
	}
	if !validCreatedAt(e.At) || !validLine(e.Actor, maxActorRunes) {
		return false
	}
	if e.Kind == EventRevoked {
		return validLine(e.Reason, maxReasonRunes)
	}
	return claimPattern.MatchString(e.Claim)
}

// replay derives the lifecycle and rejects any history the transitions
// cannot have written.
func replay(events []Event) (Lifecycle, error) {
	l := Lifecycle{State: StateAvailable, Events: events}
	claims := map[string]bool{}
	for i, e := range events {
		valid := e.Seq == i+1 && (i == 0 || events[i-1].At <= e.At)
		switch e.Kind {
		case EventTaken:
			valid = valid && l.State == StateAvailable && !claims[e.Claim]
			claims[e.Claim] = true
			l.State, l.Holder = StateTaken, &Holder{Actor: e.Actor, Claim: e.Claim, TakenAt: e.At}
		case EventReleased, EventDone:
			valid = valid && l.State == StateTaken && l.Holder.Claim == e.Claim && l.Holder.Actor == e.Actor
			l.State, l.Holder = StateAvailable, nil
			if e.Kind == EventDone {
				l.State, l.Done = StateDone, &DoneBy{Actor: e.Actor, At: e.At}
			}
		case EventRevoked:
			valid = valid && l.State == StateTaken
			l.State, l.Holder = StateAvailable, nil
		}
		if !valid {
			return Lifecycle{}, lifecycleInvalid("invalid replay")
		}
	}
	l.Revision = len(events)
	return l, nil
}

// ListView returns the valid Assignments in view, by ascending number, and
// every unreadable Assignment or lifecycle file. An Assignment whose
// lifecycle cannot be read appears in no view.
func (s Store) ListView(view View) (ViewResult, error) {
	if !slices.Contains(Views, view) {
		return ViewResult{}, invalidArguments("view must be available, taken, archive or all")
	}
	base, err := s.List()
	if err != nil {
		return ViewResult{}, err
	}
	result := ViewResult{View: view, Items: []ViewItem{}, Invalid: base.Invalid}
	report := func(file string, err error) {
		result.Invalid = append(result.Invalid, InvalidFile{File: file, Error: err.Error()})
	}
	lifecycles, err := s.listLifecycles(base.Items, report)
	if err != nil {
		return ViewResult{}, err
	}
	for _, summary := range base.Items {
		l, readable := lifecycles[summary.ID]
		if !readable {
			continue
		}
		item := ViewItem{Summary: summary, State: l.State, Revision: l.Revision}
		switch {
		case l.Holder != nil:
			item.TakenAt, item.Actor = l.Holder.TakenAt, l.Holder.Actor
		case l.Done != nil:
			item.DoneAt, item.Actor = l.Done.At, l.Done.Actor
		}
		if view == ViewAll || stateViews[l.State] == view {
			result.Items = append(result.Items, item)
		}
	}
	slices.SortFunc(result.Invalid, func(a, b InvalidFile) int { return strings.Compare(a.File, b.File) })
	return result, nil
}

// listLifecycles returns the lifecycle of each readable Assignment in items
// and reports every invalid or unreadable lifecycle entry. Only a lifecycle
// directory that cannot be inspected or listed fails the whole list.
func (s Store) listLifecycles(items []Summary, report func(string, error)) (map[string]Lifecycle, error) {
	lifecycles := map[string]Lifecycle{}
	path := filepath.Join(s.dir, lifecycleDir)
	info, err := os.Lstat(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err != nil:
		return nil, ioError("inspect lifecycle directory", err)
	case !info.IsDir():
		report(lifecycleDir, lifecycleInvalid("not a directory"))
		return lifecycles, nil
	}
	entries, err := os.ReadDir(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, ioError("read lifecycle directory", err)
	}
	// A directory that can be listed but not searched would make every
	// record look unreadable; that is one directory failure, not many.
	if _, err := os.Lstat(path + string(filepath.Separator) + "."); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, ioError("search lifecycle directory", err)
	}
	// Any failure to read one record, invalid or unreadable, hides only its Assignment.
	for _, item := range items {
		l, _, err := s.readLifecycle(item.ID)
		if err != nil {
			report(lifecycleDir+"/"+item.ID+fileSuffix, err)
			continue
		}
		lifecycles[item.ID] = l
	}
	for _, entry := range entries {
		name := entry.Name()
		match := filePattern.FindStringSubmatch(name)
		switch {
		case strings.HasPrefix(name, "."):
		case match == nil:
			report(lifecycleDir+"/"+name, lifecycleInvalid("unexpected entry"))
		case !slices.ContainsFunc(items, func(item Summary) bool { return item.ID+fileSuffix == name }):
			report(lifecycleDir+"/"+name, lifecycleInvalid("no valid assignment"))
		}
	}
	return lifecycles, nil
}

// validLine is non-blank, trimmed UTF-8 of at most maxRunes characters
// without control characters.
func validLine(value string, maxRunes int) bool {
	return value != "" && value == strings.TrimSpace(value) && utf8.ValidString(value) &&
		utf8.RuneCountInString(value) <= maxRunes && strings.IndexFunc(value, unicode.IsControl) < 0
}

func conflict(detail string) error { return &Error{Class: ClassConflict, Detail: detail} }

func staleRevision(given, current int) error {
	return conflict(fmt.Sprintf("stale revision %d, current %d", given, current))
}

func lifecycleInvalid(detail string) error {
	return &Error{Class: ClassLifecycleInvalid, Detail: detail}
}
