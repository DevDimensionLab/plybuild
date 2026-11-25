package machine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// LabelsKind is the kind of ply machine process-label list.
const LabelsKind = "MachineProcessLabels@1"

// LabelsFile is the label table below the user's home; one per machine.
const LabelsFile = ".config/ply/process-labels.json"

const labelsLock = "process-labels.lock"

// Categories are the values set accepts. Readers pass a stored value on as is.
var Categories = []string{"security", "agent", "development", "browser", "sync", "system"}

const (
	maxLabelName    = 200
	maxDisplayName  = 60
	codeInvalidArgs = "invalid_arguments"
	codeWriteFailed = "write_failed"
)

// The errors a table that exists but cannot be read gives; their text is the reason.
var (
	ErrLabelsRead   = errors.New(ReasonReadFailed)
	ErrLabelsFormat = errors.New(ReasonUnknownFormat)
)

// Label is a readable name and an optional category for one process name.
type Label struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Category    *string `json:"category"`
}

type labelTable struct {
	Labels []Label `json:"labels"`
}

// ReadLabels reads the table at path without locking or creating anything. A
// missing file is an empty table.
func ReadLabels(path string) (map[string]Label, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]Label{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %v; fix or delete the file", ErrLabelsRead, path, err)
	}
	var table labelTable
	if err := json.Unmarshal(b, &table); err != nil {
		return nil, fmt.Errorf("%w: %s: %v; fix or delete the file", ErrLabelsFormat, path, err)
	}
	labels := make(map[string]Label, len(table.Labels))
	for _, l := range table.Labels {
		labels[l.Name] = l
	}
	return labels, nil
}

// ListLabels is the table sorted by name.
func ListLabels(path string) ([]Label, error) {
	labels, err := ReadLabels(path)
	if err != nil {
		return nil, err
	}
	return sorted(labels), nil
}

// SetLabel adds or replaces the label for l.Name.
func SetLabel(path string, l Label) (Label, error) {
	if err := validLabel(l); err != nil {
		return Label{}, err
	}
	return l, change(path, func(labels map[string]Label) { labels[l.Name] = l })
}

// RemoveLabel removes the label for name; removed is false when there was none.
func RemoveLabel(path, name string) (removed bool, err error) {
	return removed, change(path, func(labels map[string]Label) {
		_, removed = labels[name]
		delete(labels, name)
	})
}

// change reads, edits and replaces the table under the lock beside it. A table
// that cannot be read is left as it is.
func change(path string, edit func(map[string]Label)) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("%s: %v", codeWriteFailed, err)
	}
	lock, err := filelock.Acquire(filepath.Join(dir, labelsLock), filelock.Options{})
	if err != nil {
		return fmt.Errorf("%s: %v", codeWriteFailed, err)
	}
	defer lock.Release()
	labels, err := ReadLabels(path)
	if err != nil {
		return err
	}
	edit(labels)
	b, err := json.MarshalIndent(labelTable{Labels: sorted(labels)}, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %v", codeWriteFailed, err)
	}
	if err := atomicfile.Replace(path, append(b, '\n'), 0o600); err != nil {
		return fmt.Errorf("%s: %v", codeWriteFailed, err)
	}
	return nil
}

func sorted(labels map[string]Label) []Label {
	list := make([]Label, 0, len(labels))
	for _, l := range labels {
		list = append(list, l)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	return list
}

func validLabel(l Label) error {
	invalid := func(format string, args ...any) error {
		return fmt.Errorf("%s: "+format, append([]any{codeInvalidArgs}, args...)...)
	}
	if !oneLine(l.Name, maxLabelName) {
		return invalid("the process name must be one line of 1 to %d characters", maxLabelName)
	}
	if !oneLine(l.DisplayName, maxDisplayName) {
		return invalid("--display-name must be one line of 1 to %d characters", maxDisplayName)
	}
	if l.Category != nil && !slices.Contains(Categories, *l.Category) {
		return invalid("--category must be one of %s", strings.Join(Categories, ", "))
	}
	return nil
}

func oneLine(s string, max int) bool {
	return s != "" && utf8.RuneCountInString(s) <= max && !strings.ContainsFunc(s, unicode.IsControl)
}

// WithLabels names the processes from the table; a table that cannot be read
// leaves every row without a name and gives the reason.
func WithLabels(s StatusWithProcesses, labels map[string]Label, err error) StatusWithProcesses {
	if s.Processes == nil {
		return s
	}
	if err != nil {
		reason := ReasonReadFailed
		if errors.Is(err, ErrLabelsFormat) {
			reason = ReasonUnknownFormat
		}
		s.ProcessLabelsReason = &reason
		return s
	}
	for i, p := range s.Processes {
		if l, ok := labels[p.Name]; ok {
			s.Processes[i].DisplayName, s.Processes[i].Category = &l.DisplayName, l.Category
		}
	}
	return s
}
