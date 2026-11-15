package plan

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const journalSuffix = ".journal.jsonl"

// JournalEntry is one record in a plan's own journal, for work that belongs
// to the plan rather than to one Task: a step start or finish, or a note. At
// is the store's time; the rest is the reporter's claim, whose kinds and
// outcomes the caller checks.
type JournalEntry struct {
	At           string  `json:"at"`
	Type         string  `json:"type"`
	Title        string  `json:"title"`
	ActorID      string  `json:"actor_id"`
	Role         string  `json:"role"`
	StepID       *string `json:"step_id,omitempty"`
	StepKind     *string `json:"step_kind,omitempty"`
	ParentStepID *string `json:"parent_step_id,omitempty"`
	Outcome      *string `json:"outcome,omitempty"`
}

// AppendJournal appends e to the journal of the existing plan id.
func (s Store) AppendJournal(id string, e JournalEntry) (JournalEntry, error) {
	if !idPattern.MatchString(id) {
		return JournalEntry{}, invalidArguments("plan ID must be PLAN-<n>")
	}
	step := e.Type == "step_started" || e.Type == "step_finished"
	if !step && e.Type != "note" || step != (e.StepID != nil && e.StepKind != nil) {
		return JournalEntry{}, invalidArguments("a step start or finish needs a step ID and kind; a note has none")
	}
	if !validLine(e.Title, maxLineRunes) || !validLine(e.ActorID, maxLineRunes) || !validLine(e.Role, maxLineRunes) {
		return JournalEntry{}, invalidArguments("title, actor and role must be one line")
	}
	_, err := s.locked(false, func() (View, error) {
		if _, err := s.read(id); err != nil {
			return View{}, err
		}
		e.At = s.now().UTC().Format(time.RFC3339Nano)
		line, err := json.Marshal(e)
		if err != nil {
			return View{}, ioError("encode journal record", err, nothingWritten)
		}
		file, err := os.OpenFile(s.journalPath(id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return View{}, ioError("open journal of "+id, err, nothingWritten)
		}
		defer file.Close()
		if _, err = file.Write(append(line, '\n')); err != nil {
			return View{}, ioError("write journal of "+id, err, uncertainSuffix)
		}
		return View{}, nil
	})
	return e, err
}

// Journal reads the journal of plan id; a plan without one has none.
func (s Store) Journal(id string) ([]JournalEntry, error) {
	if !idPattern.MatchString(id) {
		return nil, invalidArguments("plan ID must be PLAN-<n>")
	}
	entries := []JournalEntry{}
	file, err := os.Open(s.journalPath(id))
	if errors.Is(err, fs.ErrNotExist) {
		return entries, nil
	}
	if err != nil {
		return nil, ioError("read journal of "+id, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var e JournalEntry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return entries, invalid("journal of " + id + " has an unreadable record")
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return entries, ioError("read journal of "+id, err)
	}
	return entries, nil
}

func (s Store) journalPath(id string) string { return filepath.Join(s.dir, id+journalSuffix) }
