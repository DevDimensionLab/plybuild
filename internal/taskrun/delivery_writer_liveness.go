package taskrun

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const writerLivenessMethod = "local_process_cwd_and_session_v1"

const writerLivenessLimitation = "Observed locally: processes whose current directory (cwd) is inside the Task worktree, processes holding the provider transcript open, and processes naming the session in their arguments. A process that left the Task worktree cwd and neither holds the transcript nor names the session is not observable."

// systemLocalProcesses reads bounded local identities with lsof and ps. Command
// lines are only matched against the session id; they are never returned.
func systemLocalProcesses(sessionID, transcript string) ([]LocalProcess, error) {
	run := func(allowNone bool, name string, args ...string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, name, args...)
		command.WaitDelay = time.Second
		var output workflowLimitedOutput
		command.Stdout, command.Stderr = &output, io.Discard
		err := command.Run()
		if output.overflow {
			return "", errors.New("local process inspection exceeded its limit")
		}
		var exit *exec.ExitError
		if allowNone && errors.As(err, &exit) && exit.ExitCode() == 1 && len(output.data) == 0 {
			return "", nil
		}
		return string(output.data), err
	}
	processes := map[int]*LocalProcess{}
	table, err := run(false, "lsof", "-n", "-P", "-w", "-d", "cwd", "-F", "pRn")
	if err != nil {
		return nil, err
	}
	var current *LocalProcess
	for _, line := range strings.Split(table, "\n") {
		if line == "" {
			continue
		}
		switch line[0] {
		case 'p':
			pid, err := strconv.Atoi(line[1:])
			if err != nil || pid <= 0 {
				return nil, errors.New("local process inspection returned an invalid identity")
			}
			current = &LocalProcess{PID: pid}
			processes[pid] = current
		case 'R':
			if current != nil {
				current.ParentPID, _ = strconv.Atoi(line[1:])
			}
		case 'n':
			if current != nil {
				current.CWD = line[1:]
			}
		}
	}
	mark := func(pid int, evidence string) {
		if p, ok := processes[pid]; ok {
			p.Session = evidence
		} else {
			processes[pid] = &LocalProcess{PID: pid, Session: evidence}
		}
	}
	if transcript != "" {
		holders, err := run(true, "lsof", "-n", "-P", "-w", "-F", "p", "--", transcript)
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(holders, "\n") {
			if pid, err := strconv.Atoi(strings.TrimPrefix(line, "p")); strings.HasPrefix(line, "p") && err == nil {
				mark(pid, "transcript_open")
			}
		}
	}
	if sessionID != "" {
		commands, err := run(false, "/bin/ps", "-axo", "pid=,command=")
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(commands, "\n") {
			fields := strings.Fields(line)
			if len(fields) > 1 && strings.Contains(line, sessionID) {
				if pid, err := strconv.Atoi(fields[0]); err == nil {
					mark(pid, "session_argument")
				}
			}
		}
	}
	out := make([]LocalProcess, 0, len(processes))
	for _, p := range processes {
		out = append(out, *p)
	}
	return out, nil
}

// invokerShellPID returns the shell of the Herdr-identified invoking pane, or
// zero when it cannot be established. Only that exact shell is excludable.
func invokerShellPID(d Dependencies, s workflowState, pane string) int {
	if pane == "" {
		return 0
	}
	raw, err := workflowCall(d, s.Request, "pane", "process-info", "--pane", pane)
	if err != nil {
		return 0
	}
	var info struct {
		Info struct {
			PaneID   string `json:"pane_id"`
			ShellPID *int   `json:"shell_pid"`
		} `json:"process_info"`
	}
	if json.Unmarshal(raw, &info) != nil || info.Info.PaneID != pane || info.Info.ShellPID == nil || *info.Info.ShellPID <= 0 {
		return 0
	}
	return *info.Info.ShellPID
}

// observeWriterLiveness reconciles possible writers before a new provider is
// started. It excludes only this command, its own children and, when it is
// this command's direct parent, the Herdr-identified invoker shell.
func observeWriterLiveness(d Dependencies, s workflowState, transcript string, invokerShell int) (TaskWriterLiveness, error) {
	out := TaskWriterLiveness{Method: writerLivenessMethod, Scope: s.Observed.Target.WorktreeLocator, RemainingPIDs: []int{}, ExcludedPIDs: []int{}, Limitation: writerLivenessLimitation}
	observe := d.LocalProcesses
	if observe == nil {
		observe = systemLocalProcesses
	}
	table, err := observe(s.Result.Transport.AgentSessionID, transcript)
	if err != nil {
		return out, err
	}
	if len(table) == 0 {
		return out, errors.New("local process inspection returned no processes")
	}
	self := os.Getpid()
	parent := map[int]int{}
	children := map[int][]int{}
	for _, p := range table {
		parent[p.PID] = p.ParentPID
		children[p.ParentPID] = append(children[p.ParentPID], p.PID)
	}
	excluded := map[int]bool{}
	queue := []int{self}
	for len(queue) > 0 {
		pid := queue[0]
		queue = queue[1:]
		if excluded[pid] {
			continue
		}
		excluded[pid] = true
		queue = append(queue, children[pid]...)
	}
	if invokerShell > 0 && parent[self] == invokerShell {
		excluded[invokerShell] = true
	}
	for _, p := range table {
		if excluded[p.PID] {
			out.ExcludedPIDs = append(out.ExcludedPIDs, p.PID)
			continue
		}
		if p.Session != "" || p.CWD != "" && sourceContains(s.Observed.Target.WorktreeLocator, p.CWD) {
			out.RemainingPIDs = append(out.RemainingPIDs, p.PID)
		}
	}
	sort.Ints(out.RemainingPIDs)
	sort.Ints(out.ExcludedPIDs)
	return out, nil
}
