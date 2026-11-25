package machine

import (
	"bufio"
	"bytes"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
)

const maxProcesses = 10

// Process is one row of ps: the program name only, never arguments or path.
type Process struct {
	PID        int     `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	RSSBytes   int64   `json:"rss_bytes"`
	// DisplayName and Category come from the label table; null without a label.
	DisplayName *string `json:"display_name"`
	Category    *string `json:"category"`
}

// ProcessReader is the system seam for --processes.
type ProcessReader interface {
	Processes() ([]Process, error)
}

// StatusWithProcesses is Status plus the processes that use most CPU. The
// fields exist only when processes were asked for.
type StatusWithProcesses struct {
	Status
	Processes       []Process `json:"processes"`
	ProcessesReason *string   `json:"processes_reason"`
	// ProcessLabelsReason is set when the label table exists but cannot be read.
	ProcessLabelsReason *string `json:"process_labels_reason"`
}

// WithProcesses adds the processes to s; a read error leaves them null with a reason.
func WithProcesses(s Status, processes []Process, err error) StatusWithProcesses {
	if err != nil {
		return StatusWithProcesses{Status: s, ProcessesReason: ptr(ReasonReadFailed)}
	}
	return StatusWithProcesses{Status: s, Processes: processes}
}

func (s StatusWithProcesses) Text() string {
	var b strings.Builder
	b.WriteString(s.Status.Text())
	if s.Processes == nil {
		fmt.Fprintf(&b, "\n  processes unknown (%s)", *s.ProcessesReason)
	}
	for _, p := range s.Processes {
		fmt.Fprintf(&b, "\n  %.1f%%  %.0f MiB  %d %s%s", p.CPUPercent, float64(p.RSSBytes)/(1<<20), p.PID, p.Name, p.label())
	}
	if s.ProcessLabelsReason != nil {
		fmt.Fprintf(&b, "\n  process labels unknown (%s)", *s.ProcessLabelsReason)
	}
	return b.String()
}

func (p Process) label() string {
	switch {
	case p.DisplayName == nil:
		return ""
	case p.Category == nil:
		return " (" + *p.DisplayName + ")"
	}
	return " (" + *p.DisplayName + ", " + *p.Category + ")"
}

// parseProcesses reads `ps -Ao pid=,pcpu=,rss=,comm=` output: three numbers,
// then comm, which may contain spaces. It skips lines it cannot read and the
// ps process itself, and keeps the ten that use most CPU.
func parseProcesses(out []byte, self int) []Process {
	processes := []Process{}
	lines := bufio.NewScanner(bytes.NewReader(out))
	for lines.Scan() {
		if p, ok := parseProcess(lines.Text()); ok && p.PID != self {
			processes = append(processes, p)
		}
	}
	sort.Slice(processes, func(i, j int) bool {
		a, b := processes[i], processes[j]
		if a.CPUPercent != b.CPUPercent {
			return a.CPUPercent > b.CPUPercent
		}
		if a.RSSBytes != b.RSSBytes {
			return a.RSSBytes > b.RSSBytes
		}
		return a.PID < b.PID
	})
	if len(processes) > maxProcesses {
		processes = processes[:maxProcesses]
	}
	return processes
}

// parseProcess reads one ps line; ok is false when it is not pid, %cpu,
// rss and comm.
func parseProcess(line string) (p Process, ok bool) {
	fields, comm := numberFields(line)
	pid, errPID := strconv.Atoi(fields[0])
	cpu, errCPU := strconv.ParseFloat(fields[1], 64)
	rss, errRSS := strconv.ParseInt(fields[2], 10, 64)
	if comm == "" || errPID != nil || errCPU != nil || errRSS != nil {
		return p, false
	}
	return Process{PID: pid, Name: path.Base(comm), CPUPercent: cpu, RSSBytes: rss * 1024}, true
}

// numberFields splits off the first three space-separated fields and returns
// the rest of the line unchanged, so comm keeps its own spaces. comm is empty
// when the line is too short.
func numberFields(line string) (fields [3]string, comm string) {
	rest := line
	for i := range fields {
		rest = strings.TrimLeft(rest, " ")
		end := strings.IndexByte(rest, ' ')
		if end < 0 {
			return fields, ""
		}
		fields[i], rest = rest[:end], rest[end:]
	}
	return fields, strings.TrimSpace(rest)
}
