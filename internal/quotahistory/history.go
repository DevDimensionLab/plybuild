// Package quotahistory keeps the quota points this machine has observed:
// <dir>/<series>/<YYYY-MM>.jsonl, one JSON point per line, appended under
// <dir>/.lock and never rewritten. A point is written only when its value
// changed or the last point of the series is ten minutes older, and never
// when it is not newer than the last point.
package quotahistory

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// Kind names the read contract of ply quota history.
const Kind = "PlyQuotaHistory@1"

// StateFile is the directory below $XDG_STATE_HOME.
const StateFile = "ply/quota-history"

// Spec is one series with the unit of its used value.
type Spec struct{ ID, Unit string }

// Series are the recorded series in their fixed order.
var Series = []Spec{
	{"claude.five_hour", "percent"}, {"claude.weekly", "percent"},
	{"codex.five_hour", "percent"}, {"codex.weekly", "percent"},
	{"github.actions", "minutes"},
}

// ErrUnknownSeries is returned for a series that is not in Series.
var ErrUnknownSeries = errors.New("unknown quota series")

// Every is the time after which an unchanged value is written again.
const Every = 10 * time.Minute

const (
	lockName  = ".lock"
	monthForm = "2006-01"
	tail      = 4096
)

// Point is one observed value; At is the source's observation time.
type Point struct {
	At       time.Time  `json:"at"`
	Used     float64    `json:"used"`
	Limit    *float64   `json:"limit,omitempty"`
	ResetsAt *time.Time `json:"resets_at,omitempty"`
}

func (p Point) sameValue(q Point) bool {
	return p.Used == q.Used && equal(p.Limit, q.Limit) && equal(p.ResetsAt, q.ResetsAt)
}

func equal[T comparable](a, b *T) bool {
	return a == nil && b == nil || a != nil && b != nil && *a == *b
}

// SeriesView is one series of a read; Points is never nil.
type SeriesView struct {
	ID     string  `json:"id"`
	Unit   string  `json:"unit"`
	Points []Point `json:"points"`
}

// Invalid is a line that could not be read; File is relative to the directory.
type Invalid struct {
	File  string `json:"file"`
	Line  int    `json:"line"`
	Error string `json:"error"`
}

// History is the result of Read; both lists are never nil.
type History struct {
	Series  []SeriesView `json:"series"`
	Invalid []Invalid    `json:"invalid"`
}

// Store reads and writes the history in one directory.
type Store struct{ dir string }

// Open returns the store of dir; nothing is created until a point is written.
func Open(dir string) Store { return Store{dir: dir} }

// DefaultDir is $XDG_STATE_HOME/ply/quota-history, by default below ~/.local/state.
func DefaultDir(getenv func(string) string, home func() (string, error)) (string, error) {
	state := getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(state) {
		h, err := home()
		if err != nil {
			return "", err
		}
		state = filepath.Join(h, ".local", "state")
	}
	return filepath.Join(state, filepath.FromSlash(StateFile)), nil
}

func unit(id string) (string, bool) {
	for _, s := range Series {
		if s.ID == id {
			return s.Unit, true
		}
	}
	return "", false
}

// Record appends p to series unless the rule above leaves it out, and says
// whether it was written.
func (s Store) Record(series string, p Point) (bool, error) {
	if _, ok := unit(series); !ok {
		return false, fmt.Errorf("%w: %s", ErrUnknownSeries, series)
	}
	p.At = p.At.UTC().Truncate(time.Second)
	if p.ResetsAt != nil {
		reset := p.ResetsAt.UTC().Truncate(time.Second)
		p.ResetsAt = &reset
	}
	directory := filepath.Join(s.dir, series)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return false, err
	}
	lock, err := filelock.Acquire(filepath.Join(s.dir, lockName), filelock.Options{})
	if err != nil {
		return false, err
	}
	written, err := s.append(directory, p)
	if releaseErr := lock.Release(); releaseErr != nil && err == nil {
		err = releaseErr
	}
	return written, err
}

func (s Store) append(directory string, p Point) (bool, error) {
	last, found, err := lastPoint(directory)
	if err != nil {
		return false, err
	}
	if found && (!p.At.After(last.At) || p.sameValue(last) && p.At.Sub(last.At) < Every) {
		return false, nil
	}
	line, err := json.Marshal(p)
	if err != nil {
		return false, err
	}
	file, err := os.OpenFile(filepath.Join(directory, p.At.Format(monthForm)+".jsonl"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return false, err
	}
	_, err = file.Write(append(line, '\n'))
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err == nil, err
}

// months returns the month files of directory, newest last.
func months(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		name := entry.Name()
		if _, err := time.Parse(monthForm, strings.TrimSuffix(name, ".jsonl")); err == nil && strings.HasSuffix(name, ".jsonl") && entry.Type().IsRegular() {
			names = append(names, name)
		}
	}
	return names, nil // ReadDir sorts by name
}

// lastPoint reads the last readable point from the end of the newest files.
func lastPoint(directory string) (Point, bool, error) {
	names, err := months(directory)
	if err != nil {
		return Point{}, false, err
	}
	for i := len(names) - 1; i >= 0; i-- {
		lines, err := tailLines(filepath.Join(directory, names[i]))
		if err != nil {
			return Point{}, false, err
		}
		for j := len(lines) - 1; j >= 0; j-- {
			if p, err := parse(lines[j]); err == nil {
				return p, true, nil
			}
		}
	}
	return Point{}, false, nil
}

// tailLines returns the complete lines in the last bytes of path.
func tailLines(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	offset := max(0, info.Size()-tail)
	data := make([]byte, info.Size()-offset)
	if _, err := file.ReadAt(data, offset); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	lines := bytes.Split(data, []byte{'\n'})
	if offset > 0 {
		lines = lines[1:] // the first line may be cut
	}
	return lines, nil
}

func parse(line []byte) (Point, error) {
	var p Point
	if err := json.Unmarshal(line, &p); err != nil {
		return Point{}, errors.New("not a JSON point")
	}
	if p.At.IsZero() {
		return Point{}, errors.New("point without at")
	}
	return p, nil
}

// Read returns the points of ids (every series when empty) observed from
// since through until, sorted and without repeated observation times.
func (s Store) Read(ids []string, since, until time.Time) (History, error) {
	history := History{Series: []SeriesView{}, Invalid: []Invalid{}}
	for _, id := range ids {
		if _, ok := unit(id); !ok {
			return history, fmt.Errorf("%w: %s", ErrUnknownSeries, id)
		}
	}
	for _, spec := range Series {
		if len(ids) > 0 && !slices.Contains(ids, spec.ID) {
			continue
		}
		points, err := s.readSeries(spec.ID, since, until, &history)
		if err != nil {
			return history, err
		}
		history.Series = append(history.Series, SeriesView{ID: spec.ID, Unit: spec.Unit, Points: points})
	}
	return history, nil
}

func (s Store) readSeries(id string, since, until time.Time, history *History) ([]Point, error) {
	points := []Point{}
	first := time.Date(since.UTC().Year(), since.UTC().Month(), 1, 0, 0, 0, 0, time.UTC)
	for month := first; !month.After(until.UTC()); month = month.AddDate(0, 1, 0) {
		read, err := s.readMonth(filepath.Join(id, month.Format(monthForm)+".jsonl"), history)
		if err != nil {
			return nil, err
		}
		for _, p := range read {
			if !p.At.Before(since) && !p.At.After(until) {
				points = append(points, p)
			}
		}
	}
	sort.SliceStable(points, func(i, j int) bool { return points[i].At.Before(points[j].At) })
	return slices.CompactFunc(points, func(a, b Point) bool { return a.At.Equal(b.At) }), nil
}

// readMonth returns the readable points of one month file and adds its other
// lines to history.Invalid; a missing file has no points.
func (s Store) readMonth(name string, history *History) ([]Point, error) {
	data, err := os.ReadFile(filepath.Join(s.dir, name))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var points []Point
	for n, line := range bytes.Split(data, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		if p, err := parse(line); err != nil {
			history.Invalid = append(history.Invalid, Invalid{File: filepath.ToSlash(name), Line: n + 1, Error: err.Error()})
		} else {
			points = append(points, p)
		}
	}
	return points, nil
}
