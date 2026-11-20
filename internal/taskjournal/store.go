package taskjournal

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func hashStream(f *os.File) (string, error) {
	a, e := f.Stat()
	if e != nil {
		return "", e
	}
	h := sha256.New()
	n, e := io.Copy(h, f)
	if e != nil {
		return "", e
	}
	z, e := f.Stat()
	if e != nil {
		return "", e
	}
	if n != a.Size() || a.Size() != z.Size() || a.ModTime() != z.ModTime() {
		return "", conflict("source changed during read")
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
func workspaceID(root, marker string) string { return "ws_" + digest([]string{root, marker}) }

// JournalDir is the directory of the Task journals of the Project of work.
func JournalDir(work workspace.ProjectWork) string {
	return filepath.Join(work.TaskProcess(), "v1")
}
func storePath(work workspace.ProjectWork, task string) string {
	return filepath.Join(JournalDir(work), task)
}

type stored struct {
	Record          Record
	Common          Common
	Event           *EventInput
	Observation     *ObservationInput
	Locator, SHA256 string
}

func readStore(work workspace.ProjectWork, task, ws string) ([]stored, error) {
	out := []stored{}
	p := storePath(work, task)
	d, e := openDir(p, false)
	if os.IsNotExist(e) {
		return out, nil
	}
	if e != nil {
		return out, e
	}
	defer d.Close()
	entries, e := d.ReadDir(-1)
	if e != nil {
		return out, e
	}
	sortEntries(entries)
	keys := map[string]bool{}
	var prev *string
	for _, entry := range entries {
		n := entry.Name()
		if n == "append.lock" || strings.HasPrefix(n, ".publish-") {
			continue
		}
		want := fmt.Sprintf("%09d.json", len(out)+1)
		if n != want {
			return out, conflict("journal sequence gap or unexpected file: " + n)
		}
		path := filepath.Join(p, n)
		r, b, e := readRecord(path, ws, task, len(out)+1, prev)
		if e != nil {
			return out, e
		}
		kind := "event"
		if strings.Contains(string(r.Input), `"kind":"ply.workspace.task-journal-observation-input"`) {
			kind = "observation"
		}
		c, ev, ob, canonical, e := parseInput(r.Input, kind, task)
		if e != nil {
			return out, e
		}
		if hash(canonical) != r.InputSHA256 || r.EventID != eventID(ws, task, c.PublicationKey) || keys[c.PublicationKey] {
			return out, conflict("journal identity or duplicate publication key")
		}
		if e = validateHistory(out, ev, ob); e != nil {
			return out, e
		}
		keys[c.PublicationKey] = true
		h := hash(b)
		out = append(out, stored{r, c, ev, ob, path, h})
		prev = &h
	}
	return out, nil
}

// readRecord reads one private record file and checks its envelope: binding to
// the workspace and Task, position in the hash chain and a UTC registration.
func readRecord(path, ws, task string, sequence int, prev *string) (Record, []byte, error) {
	var r Record
	b, e := readFile(path, 128<<10)
	if e != nil {
		return r, nil, e
	}
	st, e := os.Lstat(path)
	if e != nil || st.Mode().Perm() != 0600 {
		return r, nil, conflict("record must be private")
	}
	if e = decode(b, 128<<10, &r); e != nil {
		return r, nil, e
	}
	if r.Kind != "ply.workspace.task-journal-record" || r.SchemaVersion != 1 || r.WorkspaceID != ws || r.TaskID != task || r.Sequence != sequence || !equal(r.PreviousSHA256, prev) || hash(r.Input) != r.InputSHA256 {
		return r, nil, conflict("journal record binding or chain differs")
	}
	tm, e := parseTime(r.RecordedAt)
	if e != nil || tm.Location() != time.UTC {
		return r, nil, conflict("recorded_at is not UTC")
	}
	return r, b, nil
}
func eventID(ws, task, key string) string { return "jev_" + digest([]string{ws, task, key}) }
func (s Service) check(stage string) error {
	if s.fault != nil {
		return s.fault(stage)
	}
	return nil
}
func receipt(r stored, created bool) AppendResult {
	return AppendResult{"ply.workspace.task-journal-append", 1, r.Record.TaskID, r.Common.PublicationKey, r.Record.EventID, r.Record.InputSHA256, r.Record.RecordedAt, created, r.Locator, r.SHA256, false, nil}
}
func retry(records []stored, c Common, b []byte) (AppendResult, bool, error) {
	for _, r := range records {
		if r.Common.PublicationKey == c.PublicationKey {
			if r.Record.InputSHA256 != hash(b) {
				return AppendResult{}, true, conflict("publication key already names different input")
			}
			result := receipt(r, false)
			// Identity is already resolved. Mutable source loss is diagnostic and
			// must never turn a successful immutable retry into a failed append.
			for _, source := range r.Common.Sources {
				h, e := sourceHash(source.Locator)
				if e != nil || h != source.SHA256 {
					result.Warnings = append(result.Warnings, fmt.Sprintf("source %q is unavailable or changed; the original event is retained", source.Locator))
				}
			}
			return result, true, nil
		}
	}
	return AppendResult{}, false, nil
}

// Ply keeps a copy of every source under its digest, so a journal stays
// complete after the original is changed or removed.
const (
	sourceLimit      = 4 << 20
	eventSourceLimit = 16 << 20
)

func sourceCopyPath(work workspace.ProjectWork, sha string) string {
	return filepath.Join(work.TaskProcess(), "sources", sha)
}

// checkSourceSizes refuses sources over 4 MiB each or 16 MiB together.
func checkSourceSizes(sources []SourceRef) error {
	total := int64(0)
	for _, src := range sources {
		info, e := os.Stat(src.Locator)
		if e != nil {
			return e
		}
		if info.Size() > sourceLimit {
			return invalid("source " + src.Locator + " exceeds 4 MiB")
		}
		if total += info.Size(); total > eventSourceLimit {
			return invalid("the sources of one record exceed 16 MiB together")
		}
	}
	return nil
}

// keepSource copies the bytes of src into Ply's source store, once per digest.
func keepSource(work workspace.ProjectWork, src SourceRef) error {
	target := sourceCopyPath(work, src.SHA256)
	if h, e := sourceHash(target); e == nil && h == src.SHA256 {
		return nil
	}
	b, e := readFile(src.Locator, sourceLimit)
	if e != nil {
		return e
	}
	if hash(b) != src.SHA256 {
		return conflict("source bytes changed before append: " + src.Locator)
	}
	if e = os.MkdirAll(filepath.Dir(target), 0o700); e != nil {
		return e
	}
	nonce := make([]byte, 8)
	if _, e = rand.Read(nonce); e != nil {
		return e
	}
	temp := fmt.Sprintf("%s.%x.tmp", target, nonce)
	if e = os.WriteFile(temp, b, 0o600); e != nil {
		return e
	}
	if e = os.Rename(temp, target); e != nil {
		os.Remove(temp)
		return e
	}
	return nil
}

// Append publishes one contribution read from file. Records with role ply
// are Ply's own and come only from AppendPly.
func (s Service) Append(task, file, kind string) (AppendResult, error) {
	if _, e := taskwork.ParseTaskID(task); e != nil {
		return AppendResult{}, e
	}
	b, e := readFile(file, inputLimit)
	if e != nil {
		return AppendResult{}, e
	}
	c, _, _, _, e := parseInput(b, kind, task)
	if e != nil {
		return AppendResult{}, e
	}
	if c.Actor.Role == plyRole {
		return AppendResult{}, invalid("records with role ply are written by Ply, not from a file")
	}
	return s.appendInput(task, b, kind)
}

// AppendEvent publishes one event built by a caller instead of read from a
// file, under the same rules as Append.
func (s Service) AppendEvent(task string, in EventInput) (AppendResult, error) {
	if in.Actor.Role == plyRole {
		return AppendResult{}, invalid("records with role ply are written by Ply")
	}
	b, e := Canonical(in)
	if e != nil {
		return AppendResult{}, e
	}
	return s.appendInput(task, b, "event")
}

// attemptOperations are recorded once per attempt; every other operation is
// recorded once per distinct action, so an identical retry adds nothing.
var attemptOperations = map[string]bool{"goal_start_reserved": true, "goal_start_started": true, "goal_start_prompt_submitted": true}

// AppendPly records an operation Ply performed for actor, at Ply's own time.
// References with an empty value are left out. A plan link equal to the
// latest one, and a repeated identical operation, return the earlier record
// instead of a new one.
func (s Service) AppendPly(task, actor, title string, action PlyAction) (AppendResult, error) {
	if _, e := taskwork.ParseTaskID(task); e != nil {
		return AppendResult{}, e
	}
	var refs map[string]string
	for key, value := range action.Refs {
		if value != "" {
			if refs == nil {
				refs = map[string]string{}
			}
			refs[key] = value
		}
	}
	action.Refs = refs
	if actor == "" {
		actor = "unknown"
	}
	b, e := taskwork.ReadTaskJournalBasis(s.Workspace, taskwork.TaskID(task))
	if e != nil {
		return AppendResult{}, e
	}
	records, e := readStore(b.Work, task, workspaceID(b.Workspace.Root, b.Workspace.MarkerSHA256))
	if e != nil {
		return AppendResult{}, e
	}
	key := ""
	switch {
	case action.Operation == "plan_link":
		var latest *stored
		for i := range records {
			if r := records[i]; r.Event != nil && r.Event.Ply != nil && r.Event.Ply.Operation == "plan_link" {
				latest = &records[i]
			}
		}
		if latest != nil && equal(latest.Event.Ply, action) {
			return receipt(*latest, false), nil
		}
	case !attemptOperations[action.Operation]:
		key = "ply:" + action.Operation + ":" + digest(action)[:32]
		for _, r := range records {
			if r.Common.PublicationKey == key {
				return receipt(r, false), nil
			}
		}
	}
	if key == "" {
		nonce := make([]byte, 8)
		if _, e := rand.Read(nonce); e != nil {
			return AppendResult{}, e
		}
		key = fmt.Sprintf("ply:%s:%x", action.Operation, nonce)
	}
	now := s.Now().UTC().Format(time.RFC3339Nano)
	in := EventInput{Common: Common{Kind: "ply.workspace.task-journal-event-input", SchemaVersion: 1, PublicationKey: key, TaskID: task,
		Actor: Actor{ID: actor, Role: plyRole}, OccurredAt: &now, TimeBasis: TimeBasis{Kind: "observed", Clock: ptr("ply"), Precision: "nanosecond"}, Sources: []SourceRef{}},
		Type: "ply_action", Title: title, Relations: []Relation{}, Ply: &action}
	raw, e := Canonical(in)
	if e != nil {
		return AppendResult{}, e
	}
	return s.appendInput(task, raw, "event")
}

// PlanLink returns the latest plan link Ply recorded for task, or nil. A
// damaged journal is read up to its validated prefix, so the damage stops
// no one who only asks for the link; journal show reports it.
func (s Service) PlanLink(task string) (*PlyAction, error) {
	b, e := taskwork.ReadTaskJournalBasis(s.Workspace, taskwork.TaskID(task))
	if e != nil {
		return nil, e
	}
	records, _ := ReadContributions(b.Workspace, b.Work, task)
	return LatestPlanLink(records), nil
}

// Contribution is one published journal record: an event or an observation
// with its identity and Ply's registration time.
type Contribution struct {
	EventID     string
	RecordedAt  string
	Common      Common
	Event       *EventInput
	Observation *ObservationInput
}

// Journal names the journal of Task in the work of its Project.
type Journal struct {
	Work workspace.ProjectWork
	Task string
}

// Journals lists the Task journals of every Project in the workspace at root,
// Project by Project.
func Journals(root string) ([]Journal, error) {
	works, e := workspace.DefinedWork(root)
	if e != nil {
		return nil, e
	}
	out := []Journal{}
	for _, work := range works {
		entries, e := os.ReadDir(JournalDir(work))
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return nil, e
		}
		for _, entry := range entries {
			if _, err := taskwork.ParseTaskID(entry.Name()); err == nil && entry.IsDir() {
				out = append(out, Journal{work, entry.Name()})
			}
		}
	}
	return out, nil
}

// ReadContributions reads the validated journal records of task in the work
// of its Project without native facts, registry or source reads. A damaged
// journal returns its validated prefix with the error.
func ReadContributions(w workspace.WorkspaceObservation, work workspace.ProjectWork, task string) ([]Contribution, error) {
	records, e := readStore(work, task, workspaceID(w.Root, w.MarkerSHA256))
	out := make([]Contribution, 0, len(records))
	for _, r := range records {
		out = append(out, Contribution{r.Record.EventID, r.Record.RecordedAt, r.Common, r.Event, r.Observation})
	}
	return out, e
}

// LatestPlanLink returns the latest plan link among records, or nil.
func LatestPlanLink(records []Contribution) *PlyAction {
	var link *PlyAction
	for _, r := range records {
		if r.Event != nil && r.Event.Ply != nil && r.Event.Ply.Operation == "plan_link" {
			link = r.Event.Ply
		}
	}
	return link
}

func (s Service) appendInput(task string, raw []byte, kind string) (result AppendResult, err error) {
	c, ev, ob, b, e := parseInput(raw, kind, task)
	if e != nil {
		return result, e
	}
	basis, e := taskwork.ReadTaskJournalBasis(s.Workspace, taskwork.TaskID(task))
	if e != nil {
		return result, e
	}
	work := basis.Work
	ws := workspaceID(basis.Workspace.Root, basis.Workspace.MarkerSHA256)
	records, e := readStore(work, task, ws)
	if e != nil {
		return result, e
	}
	if r, found, e := retry(records, c, b); found {
		return r, e
	}
	if e = s.validateAppend(basis, records, c, ev, ob); e != nil {
		return result, e
	}
	if e = checkSourceSizes(c.Sources); e != nil {
		return result, e
	}
	// No directory or lock is created until the first append is fully validated.
	p := storePath(work, task)
	if e = s.check("before_directories"); e != nil {
		return result, e
	}
	d, e := createStore(work, task)
	if e != nil {
		return result, e
	}
	defer d.Close()
	for _, part := range []string{work.TaskProcess(), JournalDir(work), p} {
		if e = regularMode(part); e != nil {
			return result, e
		}
	}
	if e = s.check("lock"); e != nil {
		return result, e
	}
	l, e := lock(d)
	if e != nil {
		return result, e
	}
	published := false
	defer func() {
		u := unlock(l)
		if u == nil {
			u = s.check("unlock")
		}
		if u != nil {
			err = u
		}
		if err != nil && published {
			result = AppendResult{}
			err = fmt.Errorf("task_journal_publication_unknown: key %s; inspect journal show %s or retry the same input; record %s: %w", c.PublicationKey, task, p, err)
		}
	}()
	// Recheck identity, key and relations after serialization; retry precedes
	// mutable artifact checks and never rewrites previously published bytes.
	nowBasis, e := taskwork.ReadTaskJournalBasis(s.Workspace, taskwork.TaskID(task))
	if e != nil {
		return result, e
	}
	if nowBasis.Workspace != basis.Workspace || nowBasis.Work != basis.Work {
		return result, conflict("workspace identity changed")
	}
	records, e = readStore(work, task, ws)
	if e != nil {
		return result, e
	}
	if r, found, e := retry(records, c, b); found {
		return r, e
	}
	if e = s.validateAppend(nowBasis, records, c, ev, ob); e != nil {
		return result, e
	}
	for _, src := range c.Sources {
		if e = keepSource(work, src); e != nil {
			return result, e
		}
	}
	var prev *string
	if len(records) > 0 {
		prev = ptr(records[len(records)-1].SHA256)
	}
	now := s.Now().UTC()
	if c.OccurredAt != nil {
		// One clock serves Ply and the reporter here; a second covers precision.
		if at, e := parseTime(*c.OccurredAt); e == nil && at.After(now.Add(time.Second)) {
			return result, invalid("occurred_at is after Ply's registration time " + now.Format(time.RFC3339))
		}
	}
	r := Record{"ply.workspace.task-journal-record", 1, ws, task, eventID(ws, task, c.PublicationKey), len(records) + 1, prev, now.Format(time.RFC3339Nano), hash(b), b}
	bytes, e := Canonical(r)
	if e != nil {
		return result, e
	}
	nonce := make([]byte, 16)
	if _, e = rand.Read(nonce); e != nil {
		return result, e
	}
	temp := fmt.Sprintf(".publish-%x", nonce)
	defer removeAt(d, temp)
	if e = s.check("write"); e != nil {
		return result, e
	}
	if e = writeTemp(d, temp, bytes); e != nil {
		return result, e
	}
	if e = s.check("before_publish"); e != nil {
		return result, e
	}
	name := fmt.Sprintf("%09d.json", r.Sequence)
	if e = linkAt(d, temp, name); e != nil {
		return result, e
	}
	published = true
	if e = s.check("after_publish"); e != nil {
		return result, e
	}
	if e = d.Sync(); e != nil {
		return result, e
	}
	if e = s.check("directory_sync"); e != nil {
		return result, e
	}
	result = receipt(stored{r, c, ev, ob, filepath.Join(p, name), hash(bytes)}, true)
	return result, nil
}

// validateAppend checks an input against the journal records and the current
// sources. appendInput runs it before the lock and again under it.
func (s Service) validateAppend(basis taskwork.TaskJournalBasis, records []stored, c Common, ev *EventInput, ob *ObservationInput) error {
	snap, e := s.build(basis, records)
	if e != nil {
		return e
	}
	if e = validateReferences(c, ev, ob, snap); e != nil {
		return e
	}
	if ev != nil && ev.Candidate != nil {
		if e = taskwork.ValidateTaskJournalCandidate(s.Workspace, basis, ev.Candidate.OID, ev.Candidate.Tree); e != nil {
			return e
		}
	}
	for _, src := range c.Sources {
		h, e := sourceHash(src.Locator)
		if e != nil {
			return e
		}
		if h != src.SHA256 {
			return conflict("source bytes differ: " + src.Locator)
		}
	}
	return nil
}
