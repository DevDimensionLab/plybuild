// Package plan stores a Project's plans as one JSON file each,
// <directory>/PLAN-<n>.json, holding every revision of the plan. Writes take
// an exclusive lock on <directory>/.lock and replace the file atomically;
// reads take none and create nothing. Errors are *Error with a fixed Class
// and a detail that never contains absolute paths or submitted text.
package plan

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

// DataPathCategory is the data_paths category of a Project's plans.
const DataPathCategory = "plans"

// Error classes shared with the CLI contract.
const (
	ClassInvalidArguments = "plan_invalid_arguments"
	ClassNotFound         = "plan_not_found"
	ClassInvalid          = "plan_invalid"
	ClassConflict         = "plan_conflict"
	ClassIOError          = "plan_io_error"
	ClassFull             = "plan_full"
)

// Error is one `<class>: <detail>` failure.
type Error struct{ Class, Detail string }

func (e *Error) Error() string { return e.Class + ": " + e.Detail }

// Statuses are the states of a plan and of a step.
var Statuses = []string{"planned", "active", "waiting", "blocked", "done", "unknown"}

// Attention is a stated need for someone's action.
type Attention struct {
	Reason     string `json:"reason"`
	Who        string `json:"who"`
	NeedsHuman bool   `json:"needs_human"`
}

// Step is one ordered main step of a plan.
type Step struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Outcome    string     `json:"outcome"`
	Owner      string     `json:"owner"`
	Status     string     `json:"status"`
	NextAction string     `json:"next_action"`
	Attention  *Attention `json:"attention,omitempty"`
	Result     string     `json:"result,omitempty"`
	Refs       []string   `json:"refs,omitempty"`
	DependsOn  []string   `json:"depends_on,omitempty"`
}

// Plan is one complete snapshot as submitted by add and update.
type Plan struct {
	Title      string     `json:"title"`
	Goal       string     `json:"goal"`
	Scope      string     `json:"scope"`
	Owner      string     `json:"owner"`
	Selection  string     `json:"selection"`
	Status     string     `json:"status"`
	Summary    string     `json:"summary"`
	NextAction string     `json:"next_action"`
	Attention  *Attention `json:"attention,omitempty"`
	Details    string     `json:"details,omitempty"`
	Decisions  []string   `json:"decisions,omitempty"`
	Learning   string     `json:"learning,omitempty"`
	Refs       []string   `json:"refs,omitempty"`
	Steps      []Step     `json:"steps"`
}

// Progress counts the steps of one snapshot per status; nothing else.
type Progress struct {
	Total   int `json:"total"`
	Planned int `json:"planned"`
	Active  int `json:"active"`
	Waiting int `json:"waiting"`
	Blocked int `json:"blocked"`
	Done    int `json:"done"`
	Unknown int `json:"unknown"`
}

// View is one stored revision with the values derived from it.
type View struct {
	ID         string   `json:"id"`
	Revision   int      `json:"revision"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	Actor      string   `json:"actor"`
	Reason     string   `json:"reason"`
	NeedsHuman bool     `json:"needs_human"`
	Plan       Plan     `json:"plan"`
	Progress   Progress `json:"progress"`
}

// HistoryEntry is the metadata of one revision.
type HistoryEntry struct {
	Revision int    `json:"revision"`
	At       string `json:"at"`
	Actor    string `json:"actor"`
	Reason   string `json:"reason"`
}

// ListItem is the latest revision of one plan as shown by List.
type ListItem struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Selection  string     `json:"selection"`
	Status     string     `json:"status"`
	Owner      string     `json:"owner"`
	Summary    string     `json:"summary"`
	UpdatedAt  string     `json:"updated_at"`
	Revision   int        `json:"revision"`
	NextAction string     `json:"next_action"`
	Attention  *Attention `json:"attention"`
	NeedsHuman bool       `json:"needs_human"`
	Progress   Progress   `json:"progress"`
}

// InvalidFile is a plan file that cannot be read; File is only its name.
type InvalidFile struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

// ListResult holds readable plans and unreadable files; both are never nil.
type ListResult struct {
	Items   []ListItem    `json:"items"`
	Invalid []InvalidFile `json:"invalid"`
}

// Store reads and writes the plans of one directory.
type Store struct {
	dir     string
	now     func() time.Time
	replace func(string, []byte, fs.FileMode) error
	lock    filelock.Options

	beforeOpen func(path string) // test seam: runs right before a read opens the plan file
}

// Open returns the store for directory, the absolute plans directory itself.
// Nothing is created until a write.
func Open(directory string) Store {
	return Store{dir: directory, now: time.Now, replace: atomicfile.Replace}
}

const (
	lockName          = ".lock"
	fileSuffix        = ".json"
	schemaVersion     = 1
	maxNumber         = 999999999
	maxSnapshotBytes  = 128 * 1024
	maxFileBytes      = 8 << 20
	maxSteps          = 128
	maxListEntries    = 128
	maxLineRunes      = 200
	maxReasonRunes    = 1000
	maxSummaryRunes   = 500
	timeLayout        = "2006-01-02T15:04:05Z"
	uncertainSuffix   = "; result unknown, run show or list before writing again"
	nothingWritten    = "; no plan file changed"
	selectionProposal = "proposal"
	selectionSelected = "selected"
	statusDone        = "done"
)

var (
	idPattern     = regexp.MustCompile(`^PLAN-([1-9][0-9]{0,8})$`)
	filePattern   = regexp.MustCompile(`^PLAN-([1-9][0-9]{0,8})\.json$`)
	stepIDPattern = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,63}$`)
)

// Add validates the snapshot, actor and reason before any effect, then stores
// revision 1 of a new plan numbered one above the highest PLAN-<n>.json name.
func (s Store) Add(snapshot []byte, actor, reason string) (View, error) {
	p, err := s.checkWrite(snapshot, actor, reason)
	if err != nil {
		return View{}, err
	}
	return s.locked(true, func() (View, error) {
		number, err := s.highestNumber()
		if err != nil {
			return View{}, err
		}
		if number >= maxNumber {
			return View{}, &Error{Class: ClassFull, Detail: "plan id space exhausted" + nothingWritten}
		}
		id := "PLAN-" + strconv.Itoa(number+1)
		return s.write(document{SchemaVersion: schemaVersion, ID: id}, p, actor, reason)
	})
}

// Update stores the snapshot as the next revision when expected is the
// current revision; otherwise it returns plan_conflict and writes nothing.
func (s Store) Update(id string, expected int, snapshot []byte, actor, reason string) (View, error) {
	if !idPattern.MatchString(id) {
		return View{}, invalidArguments("invalid plan id")
	}
	if expected < 1 {
		return View{}, invalidArguments("revision must be a positive number")
	}
	p, err := s.checkWrite(snapshot, actor, reason)
	if err != nil {
		return View{}, err
	}
	return s.locked(false, func() (View, error) {
		d, err := s.read(id)
		if err != nil {
			return View{}, err
		}
		if current := len(d.Revisions); current != expected {
			return View{}, &Error{Class: ClassConflict, Detail: fmt.Sprintf("%s is at revision %d, not %d; run show and decide again", id, current, expected)}
		}
		return s.write(d, p, actor, reason)
	})
}

// Show returns the latest revision, or revision when it is positive.
func (s Store) Show(id string, revision int) (View, error) {
	if !idPattern.MatchString(id) {
		return View{}, invalidArguments("invalid plan id")
	}
	if revision < 0 {
		return View{}, invalidArguments("revision must be a positive number")
	}
	d, err := s.read(id)
	if err != nil {
		return View{}, err
	}
	if revision == 0 {
		revision = len(d.Revisions)
	}
	if revision > len(d.Revisions) {
		return View{}, &Error{Class: ClassNotFound, Detail: fmt.Sprintf("%s has no revision %d", id, revision)}
	}
	return d.view(revision), nil
}

// History returns the metadata of every revision in ascending order.
func (s Store) History(id string) ([]HistoryEntry, error) {
	if !idPattern.MatchString(id) {
		return nil, invalidArguments("invalid plan id")
	}
	d, err := s.read(id)
	if err != nil {
		return nil, err
	}
	history := make([]HistoryEntry, 0, len(d.Revisions))
	for _, r := range d.Revisions {
		history = append(history, HistoryEntry{Revision: r.Revision, At: r.At, Actor: r.Actor, Reason: r.Reason})
	}
	return history, nil
}

// List returns the latest revision of every readable plan, most recently
// updated first with ascending number as tiebreaker, optionally only those
// with status. Unreadable *.json files are always listed under Invalid.
func (s Store) List(status string) (ListResult, error) {
	if status != "" && !contains(Statuses, status) {
		return ListResult{}, invalidArguments("status must be one of " + strings.Join(Statuses, ", "))
	}
	result := ListResult{Items: []ListItem{}, Invalid: []InvalidFile{}}
	entries, err := os.ReadDir(s.dir)
	if errors.Is(err, fs.ErrNotExist) {
		return result, nil
	}
	if err != nil {
		return result, ioError("read plans directory", err)
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
		d, err := s.read(strings.TrimSuffix(name, fileSuffix))
		if err != nil {
			result.Invalid = append(result.Invalid, InvalidFile{File: name, Error: detail(err)})
			continue
		}
		v := d.view(len(d.Revisions))
		if status != "" && v.Plan.Status != status {
			continue
		}
		numbers[v.ID], _ = strconv.Atoi(match[1])
		result.Items = append(result.Items, ListItem{
			ID: v.ID, Title: v.Plan.Title, Selection: v.Plan.Selection, Status: v.Plan.Status, Owner: v.Plan.Owner,
			Summary: v.Plan.Summary, UpdatedAt: v.UpdatedAt, Revision: v.Revision, NextAction: v.Plan.NextAction,
			Attention: v.Plan.Attention, NeedsHuman: v.NeedsHuman, Progress: v.Progress,
		})
	}
	sort.Slice(result.Items, func(i, j int) bool {
		a, b := result.Items[i], result.Items[j]
		if a.UpdatedAt != b.UpdatedAt {
			return a.UpdatedAt > b.UpdatedAt
		}
		return numbers[a.ID] < numbers[b.ID]
	})
	sort.Slice(result.Invalid, func(i, j int) bool { return result.Invalid[i].File < result.Invalid[j].File })
	return result, nil
}

// document is the exact content of one plan file.
type document struct {
	SchemaVersion int        `json:"schema_version"`
	ID            string     `json:"id"`
	Revisions     []revision `json:"revisions"`
}

type revision struct {
	Revision int    `json:"revision"`
	At       string `json:"at"`
	Actor    string `json:"actor"`
	Reason   string `json:"reason"`
	Plan     Plan   `json:"plan"`
}

func (d document) view(n int) View {
	r := d.Revisions[n-1]
	v := View{ID: d.ID, Revision: n, CreatedAt: d.Revisions[0].At, UpdatedAt: r.At, Actor: r.Actor, Reason: r.Reason, Plan: r.Plan}
	v.NeedsHuman = r.Plan.Attention != nil && r.Plan.Attention.NeedsHuman
	v.Progress.Total = len(r.Plan.Steps)
	counts := map[string]*int{"planned": &v.Progress.Planned, "active": &v.Progress.Active, "waiting": &v.Progress.Waiting,
		"blocked": &v.Progress.Blocked, "done": &v.Progress.Done, "unknown": &v.Progress.Unknown}
	for _, step := range r.Plan.Steps {
		*counts[step.Status]++
		if step.Attention != nil && step.Attention.NeedsHuman {
			v.NeedsHuman = true
		}
	}
	return v
}

func (s Store) checkWrite(snapshot []byte, actor, reason string) (Plan, error) {
	if !filepath.IsAbs(s.dir) {
		return Plan{}, invalidArguments("plans directory must be absolute")
	}
	if !validLine(actor, maxLineRunes) {
		return Plan{}, invalidArguments("actor must be one non-blank line of at most 200 characters")
	}
	if !validLine(reason, maxReasonRunes) {
		return Plan{}, invalidArguments("reason must be one non-blank line of at most 1000 characters")
	}
	if len(snapshot) > maxSnapshotBytes {
		return Plan{}, invalidArguments("snapshot is larger than 128 KiB")
	}
	var p Plan
	if detail := decodeStrict(snapshot, planShape, &p); detail != "" {
		return Plan{}, invalidArguments("snapshot: " + detail)
	}
	if detail := checkPlan(p); detail != "" {
		return Plan{}, invalidArguments("snapshot: " + detail)
	}
	return p, nil
}

// locked runs write under the directory lock, creating the directory only
// when create is set. A lock failure means nothing was written; a release
// failure after a write keeps the written result known.
func (s Store) locked(create bool, write func() (View, error)) (View, error) {
	if create {
		if err := os.MkdirAll(s.dir, 0o755); err != nil {
			return View{}, ioError("create plans directory", err, nothingWritten)
		}
	} else if _, err := os.Lstat(s.dir); errors.Is(err, fs.ErrNotExist) {
		return View{}, &Error{Class: ClassNotFound, Detail: "no plans directory"}
	}
	lock, err := filelock.Acquire(filepath.Join(s.dir, lockName), s.lock)
	if err != nil {
		return View{}, ioError("lock plans", err, nothingWritten)
	}
	v, err := write()
	if releaseErr := lock.Release(); releaseErr != nil && err == nil {
		return View{}, &Error{Class: ClassIOError, Detail: fmt.Sprintf("%s revision %d was written; releasing the plans lock failed; run show before writing again", v.ID, v.Revision)}
	}
	return v, err
}

func (s Store) highestNumber() (int, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return 0, ioError("read plans directory", err, nothingWritten)
	}
	highest := 0
	for _, entry := range entries {
		if match := filePattern.FindStringSubmatch(entry.Name()); match != nil {
			if n, _ := strconv.Atoi(match[1]); n > highest {
				highest = n
			}
		}
	}
	return highest, nil
}

// write appends p as the next revision of d and replaces d's file.
func (s Store) write(d document, p Plan, actor, reason string) (View, error) {
	at := s.now().UTC().Format(timeLayout)
	if n := len(d.Revisions); n > 0 && at < d.Revisions[n-1].At {
		at = d.Revisions[n-1].At // a clock behind the file must not break its order
	}
	d.Revisions = append(d.Revisions, revision{Revision: len(d.Revisions) + 1, At: at, Actor: actor, Reason: reason, Plan: p})
	next := len(d.Revisions)
	contents, err := encode(d)
	if err != nil {
		return View{}, ioError(fmt.Sprintf("encode %s", d.ID), err, nothingWritten)
	}
	if len(contents) > maxFileBytes {
		return View{}, &Error{Class: ClassFull, Detail: fmt.Sprintf("%s revision %d would exceed 8 MiB%s", d.ID, next, nothingWritten)}
	}
	if err := s.replace(s.path(d.ID), contents, 0o644); err != nil {
		var replaceError *atomicfile.Error
		if errors.As(err, &replaceError) && (replaceError.Replaced || replaceError.Op == "replace") {
			return View{}, ioError(fmt.Sprintf("write %s revision %d", d.ID, next), err, uncertainSuffix)
		}
		return View{}, ioError(fmt.Sprintf("write %s revision %d", d.ID, next), err, nothingWritten)
	}
	return d.view(next), nil
}

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

func (s Store) path(id string) string { return filepath.Join(s.dir, id+fileSuffix) }

// errNotRegular is a plan path that is a symlink or another non-regular file.
var errNotRegular = errors.New("not a regular file")

// read returns the valid document of id. It takes no lock and creates
// nothing. It checks the file it opened, so a concurrent atomic replacement
// yields the whole old or new file, never a refusal.
func (s Store) read(id string) (document, error) {
	if s.beforeOpen != nil {
		s.beforeOpen(s.path(id))
	}
	file, err := openPlanFile(s.path(id))
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return document{}, &Error{Class: ClassNotFound, Detail: id}
	case errors.Is(err, errNotRegular):
		return document{}, invalid("not a regular file")
	case err != nil:
		return document{}, ioError("read "+id, err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return document{}, ioError("inspect "+id, err)
	}
	if !info.Mode().IsRegular() {
		return document{}, invalid("not a regular file")
	}
	if info.Size() > maxFileBytes {
		return document{}, invalid("file too large")
	}
	contents, err := io.ReadAll(io.LimitReader(file, maxFileBytes+1))
	if err != nil {
		return document{}, ioError("read "+id, err)
	}
	if len(contents) > maxFileBytes {
		return document{}, invalid("file too large")
	}
	return decodeDocument(id, contents)
}

func decodeDocument(id string, contents []byte) (document, error) {
	// The version is checked first so a newer file is reported as such, not by its new fields.
	var version struct {
		SchemaVersion json.RawMessage `json:"schema_version"`
	}
	if err := json.Unmarshal(contents, &version); err != nil {
		return document{}, invalid("invalid JSON")
	}
	if string(version.SchemaVersion) != strconv.Itoa(schemaVersion) {
		return document{}, invalid("unknown schema_version")
	}
	var d document
	if detail := decodeStrict(contents, documentShape, &d); detail != "" {
		return document{}, invalid(detail)
	}
	if d.ID != id {
		return document{}, invalid("id does not match file name")
	}
	if len(d.Revisions) == 0 {
		return document{}, invalid("no revisions")
	}
	for i, r := range d.Revisions {
		switch {
		case r.Revision != i+1:
			return document{}, invalid("revisions are not numbered 1, 2, 3, ...")
		case !validTime(r.At):
			return document{}, invalid(fmt.Sprintf("revision %d: invalid at", r.Revision))
		case !validLine(r.Actor, maxLineRunes):
			return document{}, invalid(fmt.Sprintf("revision %d: invalid actor", r.Revision))
		case !validLine(r.Reason, maxReasonRunes):
			return document{}, invalid(fmt.Sprintf("revision %d: invalid reason", r.Revision))
		}
		if detail := checkPlan(r.Plan); detail != "" {
			return document{}, invalid(fmt.Sprintf("revision %d: %s", r.Revision, detail))
		}
	}
	return d, nil
}

// shape lists the exact keys of a JSON object (required or optional) or the
// element shape of an array. A leaf shape accepts any non-null value; its type
// is checked by encoding/json.
type shape struct {
	fields   map[string]field
	elements *shape
}

type field struct {
	shape    *shape
	required bool
}

var (
	leaf           = &shape{}
	attentionShape = &shape{fields: map[string]field{
		"reason": {leaf, true}, "who": {leaf, true}, "needs_human": {leaf, true},
	}}
	stepShape = &shape{fields: map[string]field{
		"id": {leaf, true}, "title": {leaf, true}, "outcome": {leaf, true}, "owner": {leaf, true},
		"status": {leaf, true}, "next_action": {leaf, true}, "attention": {attentionShape, false},
		"result": {leaf, false}, "refs": {leaf, false}, "depends_on": {leaf, false},
	}}
	planShape = &shape{fields: map[string]field{
		"title": {leaf, true}, "goal": {leaf, true}, "scope": {leaf, true}, "owner": {leaf, true},
		"selection": {leaf, true}, "status": {leaf, true}, "summary": {leaf, true}, "next_action": {leaf, true},
		"attention": {attentionShape, false}, "details": {leaf, false}, "decisions": {leaf, false},
		"learning": {leaf, false}, "refs": {leaf, false}, "steps": {&shape{elements: stepShape}, true},
	}}
	documentShape = &shape{fields: map[string]field{
		"schema_version": {leaf, true}, "id": {leaf, true},
		"revisions": {&shape{elements: &shape{fields: map[string]field{
			"revision": {leaf, true}, "at": {leaf, true}, "actor": {leaf, true}, "reason": {leaf, true}, "plan": {planShape, true},
		}}}, true},
	}}
)

// decodeStrict decodes exactly one JSON value of shape into target. It rejects
// invalid UTF-8, unknown, case-folded, duplicate or missing keys, null and
// trailing data, which encoding/json would otherwise accept or drop silently.
func decodeStrict(contents []byte, s *shape, target any) string {
	if !utf8.Valid(contents) {
		return "invalid UTF-8"
	}
	decoder := json.NewDecoder(bytes.NewReader(contents))
	if detail := checkValue(decoder, s); detail != "" {
		return detail
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return "data after JSON value"
	}
	if err := json.Unmarshal(contents, target); err != nil {
		return "wrong JSON type"
	}
	return ""
}

func checkValue(decoder *json.Decoder, s *shape) string {
	token, err := decoder.Token()
	if err != nil {
		return "invalid JSON"
	}
	switch {
	case token == nil:
		return "null value"
	case token == json.Delim('{'):
		return checkObject(decoder, s)
	case token == json.Delim('['):
		if s.fields != nil {
			return "wrong JSON type"
		}
		element := s.elements
		if element == nil {
			element = leaf
		}
		for decoder.More() {
			if detail := checkValue(decoder, element); detail != "" {
				return detail
			}
		}
		if _, err := decoder.Token(); err != nil {
			return "invalid JSON"
		}
	case s.fields != nil || s.elements != nil:
		return "wrong JSON type"
	}
	return ""
}

func checkObject(decoder *json.Decoder, s *shape) string {
	if s.fields == nil {
		return "wrong JSON type"
	}
	seen := map[string]bool{}
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return "invalid JSON"
		}
		key := token.(string)
		f, known := s.fields[key]
		switch {
		case !known:
			return "unknown field"
		case seen[key]:
			return "duplicate field"
		}
		seen[key] = true
		if detail := checkValue(decoder, f.shape); detail != "" {
			return detail
		}
	}
	if _, err := decoder.Token(); err != nil {
		return "invalid JSON"
	}
	for key, f := range s.fields {
		if f.required && !seen[key] {
			return "missing field " + key
		}
	}
	return ""
}

// checkPlan returns what is wrong with p, naming fields and step positions
// but never their submitted values.
func checkPlan(p Plan) string {
	checks := []struct {
		ok     bool
		detail string
	}{
		{validLine(p.Title, maxLineRunes), "title must be one non-blank line of at most 200 characters"},
		{validBlock(p.Goal, true, 0), "goal must be non-blank text"},
		{validBlock(p.Scope, true, 0), "scope must be non-blank text"},
		{validLine(p.Owner, maxLineRunes), "owner must be one non-blank line of at most 200 characters"},
		{p.Selection == selectionProposal || p.Selection == selectionSelected, "selection must be proposal or selected"},
		{contains(Statuses, p.Status), "status must be one of " + strings.Join(Statuses, ", ")},
		{validBlock(p.Summary, true, maxSummaryRunes), "summary must be non-blank text of at most 500 characters"},
		{validBlock(p.NextAction, false, 0), "next_action must be text"},
		{validBlock(p.Details, false, 0), "details must be text"},
		{validBlock(p.Learning, false, 0), "learning must be text"},
		{validList(p.Decisions), "decisions must be at most 128 non-blank texts"},
		{validList(p.Refs), "refs must be at most 128 non-blank texts"},
		{len(p.Steps) >= 1 && len(p.Steps) <= maxSteps, "steps must hold 1 to 128 steps"},
	}
	for _, c := range checks {
		if !c.ok {
			return c.detail
		}
	}
	if detail := checkAttention(p.Status, p.Attention); detail != "" {
		return detail
	}
	positions := map[string]int{}
	for i, step := range p.Steps {
		if detail := checkStep(step); detail != "" {
			return fmt.Sprintf("step %d: %s", i+1, detail)
		}
		if _, seen := positions[step.ID]; seen {
			return fmt.Sprintf("step %d: id is not unique", i+1)
		}
		positions[step.ID] = i
	}
	if detail := checkDependencies(p.Steps, positions); detail != "" {
		return detail
	}
	if p.Status == statusDone {
		for i, step := range p.Steps {
			if step.Status != statusDone {
				return fmt.Sprintf("status done requires every step done; step %d is not", i+1)
			}
		}
	}
	return ""
}

func checkStep(step Step) string {
	checks := []struct {
		ok     bool
		detail string
	}{
		{stepIDPattern.MatchString(step.ID), "id must match [a-z][a-z0-9_-]{0,63}"},
		{validLine(step.Title, maxLineRunes), "title must be one non-blank line of at most 200 characters"},
		{validBlock(step.Outcome, true, 0), "outcome must be non-blank text"},
		{validLine(step.Owner, maxLineRunes), "owner must be one non-blank line of at most 200 characters"},
		{contains(Statuses, step.Status), "status must be one of " + strings.Join(Statuses, ", ")},
		{validBlock(step.NextAction, false, 0), "next_action must be text"},
		{validBlock(step.Result, false, 0), "result must be text"},
		{validList(step.Refs), "refs must be at most 128 non-blank texts"},
		{len(step.DependsOn) <= maxSteps, "depends_on must hold at most 128 step ids"},
	}
	for _, c := range checks {
		if !c.ok {
			return c.detail
		}
	}
	return checkAttention(step.Status, step.Attention)
}

func checkAttention(status string, a *Attention) string {
	if a == nil {
		if status == "waiting" || status == "blocked" || status == "unknown" {
			return "status " + status + " requires attention"
		}
		return ""
	}
	if !validBlock(a.Reason, true, 0) || !validBlock(a.Who, true, 0) {
		return "attention needs a non-blank reason and who"
	}
	return ""
}

// checkDependencies rejects unknown, repeated and self targets and cycles.
// Dependencies only describe relations; nothing here orders or gates steps.
func checkDependencies(steps []Step, positions map[string]int) string {
	for i, step := range steps {
		seen := map[string]bool{}
		for _, target := range step.DependsOn {
			j, known := positions[target]
			switch {
			case !known:
				return fmt.Sprintf("step %d: depends_on names no step of the plan", i+1)
			case j == i:
				return fmt.Sprintf("step %d: depends on itself", i+1)
			case seen[target]:
				return fmt.Sprintf("step %d: depends_on repeats a step", i+1)
			}
			seen[target] = true
		}
	}
	const (
		unvisited = iota
		visiting
		visited
	)
	state := make([]int, len(steps))
	var visit func(int) bool
	visit = func(i int) bool {
		state[i] = visiting
		for _, target := range steps[i].DependsOn {
			j := positions[target]
			if state[j] == visiting || (state[j] == unvisited && !visit(j)) {
				return false
			}
		}
		state[i] = visited
		return true
	}
	for i := range steps {
		if state[i] == unvisited && !visit(i) {
			return fmt.Sprintf("step %d: depends_on forms a cycle", i+1)
		}
	}
	return ""
}

// validLine is one non-blank line without control characters.
func validLine(value string, maxRunes int) bool {
	return strings.TrimSpace(value) != "" && utf8.ValidString(value) && utf8.RuneCountInString(value) <= maxRunes &&
		strings.IndexFunc(value, unicode.IsControl) < 0
}

// validBlock is text that may hold line feeds and tabs but no other control
// characters; maxRunes 0 means no limit of its own.
func validBlock(value string, required bool, maxRunes int) bool {
	if required && strings.TrimSpace(value) == "" {
		return false
	}
	if maxRunes > 0 && utf8.RuneCountInString(value) > maxRunes {
		return false
	}
	return utf8.ValidString(value) && strings.IndexFunc(value, func(r rune) bool { return unicode.IsControl(r) && r != '\n' && r != '\t' }) < 0
}

func validList(values []string) bool {
	if len(values) > maxListEntries {
		return false
	}
	for _, value := range values {
		if !validBlock(value, true, 0) {
			return false
		}
	}
	return true
}

func validTime(value string) bool {
	at, err := time.Parse(time.RFC3339, value)
	return err == nil && at.UTC().Format(timeLayout) == value
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func invalidArguments(detail string) error {
	return &Error{Class: ClassInvalidArguments, Detail: detail}
}

func invalid(detail string) error { return &Error{Class: ClassInvalid, Detail: detail} }

// ioError keeps only the operating-system cause so details never carry
// absolute paths; suffix states what is known about the effect.
func ioError(operation string, err error, suffix ...string) error {
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
	return &Error{Class: ClassIOError, Detail: operation + ": " + cause + strings.Join(suffix, "")}
}

func detail(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.Detail
	}
	return err.Error()
}
