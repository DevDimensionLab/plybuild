package taskjournal

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
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
func storePath(root, task string) string {
	return filepath.Join(root, ".ply", "task-process", "v1", task)
}

type stored struct {
	Record          Record
	Common          Common
	Event           *EventInput
	Observation     *ObservationInput
	Locator, SHA256 string
}

func readStore(root, task, ws string) ([]stored, error) {
	out := []stored{}
	p := storePath(root, task)
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
		b, e := readFile(path, 128<<10)
		if e != nil {
			return out, e
		}
		st, e := os.Lstat(path)
		if e != nil || st.Mode().Perm() != 0600 {
			return out, conflict("record must be private")
		}
		var r Record
		if e = decode(b, 128<<10, &r); e != nil {
			return out, e
		}
		if r.Kind != "ply.workspace.task-journal-record" || r.SchemaVersion != 1 || r.WorkspaceID != ws || r.TaskID != task || r.Sequence != len(out)+1 || !equal(r.PreviousSHA256, prev) || hash(r.Input) != r.InputSHA256 {
			return out, conflict("journal record binding or chain differs")
		}
		tm, e := parseTime(r.RecordedAt)
		if e != nil || tm.Location() != time.UTC {
			return out, conflict("recorded_at is not UTC")
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

func (s Service) Append(task, file, kind string) (result AppendResult, err error) {
	if _, e := workspace.ParseTaskID(task); e != nil {
		return result, e
	}
	b, e := readFile(file, inputLimit)
	if e != nil {
		return result, e
	}
	c, ev, ob, b, e := parseInput(b, kind, task)
	if e != nil {
		return result, e
	}
	basis, e := workspace.ReadTaskJournalBasis(s.Workspace, workspace.TaskID(task))
	if e != nil {
		return result, e
	}
	root := basis.Workspace.Root
	ws := workspaceID(root, basis.Workspace.MarkerSHA256)
	records, e := readStore(root, task, ws)
	if e != nil {
		return result, e
	}
	if r, found, e := retry(records, c, b); found {
		return r, e
	}
	snap, runs, e := s.build(basis, records)
	if e != nil {
		return result, e
	}
	if e = validateReferences(c, ev, ob, snap, runs); e != nil {
		return result, e
	}
	if ev != nil && ev.Candidate != nil {
		if e = workspace.ValidateTaskJournalCandidate(s.Workspace, basis, ev.Candidate.OID, ev.Candidate.Tree); e != nil {
			return result, e
		}
	}
	for _, src := range c.Sources {
		h, e := sourceHash(src.Locator)
		if e != nil {
			return result, e
		}
		if h != src.SHA256 {
			return result, conflict("source bytes differ: " + src.Locator)
		}
	}
	// No directory or lock is created until the first append is fully validated.
	p := storePath(root, task)
	if e = s.check("before_directories"); e != nil {
		return result, e
	}
	d, e := createStore(root, task)
	if e != nil {
		return result, e
	}
	defer d.Close()
	for _, part := range []string{filepath.Join(root, ".ply", "task-process"), filepath.Join(root, ".ply", "task-process", "v1"), p} {
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
	nowBasis, e := workspace.ReadTaskJournalBasis(s.Workspace, workspace.TaskID(task))
	if e != nil {
		return result, e
	}
	if nowBasis.Workspace != basis.Workspace {
		return result, conflict("workspace identity changed")
	}
	records, e = readStore(root, task, ws)
	if e != nil {
		return result, e
	}
	if r, found, e := retry(records, c, b); found {
		return r, e
	}
	snap, runs, e = s.build(nowBasis, records)
	if e != nil {
		return result, e
	}
	if e = validateReferences(c, ev, ob, snap, runs); e != nil {
		return result, e
	}
	if ev != nil && ev.Candidate != nil {
		if e = workspace.ValidateTaskJournalCandidate(s.Workspace, basis, ev.Candidate.OID, ev.Candidate.Tree); e != nil {
			return result, e
		}
	}
	for _, src := range c.Sources {
		h, e := sourceHash(src.Locator)
		if e != nil {
			return result, e
		}
		if h != src.SHA256 {
			return result, conflict("source bytes changed before append")
		}
	}
	var prev *string
	if len(records) > 0 {
		prev = ptr(records[len(records)-1].SHA256)
	}
	r := Record{"ply.workspace.task-journal-record", 1, ws, task, eventID(ws, task, c.PublicationKey), len(records) + 1, prev, s.Now().UTC().Format(time.RFC3339Nano), hash(b), b}
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
