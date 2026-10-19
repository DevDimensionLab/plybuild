//go:build !windows

package taskrun

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// Only bounded, read-only process inspection; raw streams and command lines
// never enter diagnostics. The shared probe changes no file bytes and cannot
// succeed while an exclusive writer owns the same inode.
func codexIdentityCommand(deadline time.Time, name string, args ...string) (string, error) {
	limit := time.Now().Add(3 * time.Second)
	if !deadline.IsZero() && deadline.Before(limit) {
		limit = deadline
	}
	ctx, cancel := context.WithDeadline(context.Background(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.WaitDelay = time.Second
	var out workflowLimitedOutput
	cmd.Stdout, cmd.Stderr = &out, io.Discard
	if err := cmd.Run(); err != nil || out.overflow {
		return "", workflowError(4, "Codex runtime process inspection unavailable; restore read access before observing this same attempt; no input sent")
	}
	return string(out.data), nil
}

func codexLockHolder(deadline time.Time, path string, pid int) error {
	out, err := codexIdentityCommand(deadline, "lsof", "-n", "-P", "-w", "-F", "p", "--", path)
	if err != nil {
		return err
	}
	found := false
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		// macOS lsof emits descriptor markers even when only p is selected.
		if strings.HasPrefix(line, "f") {
			continue
		}
		// This observer has the same inode open read-only for its lock probe.
		if line == "p"+strconv.Itoa(os.Getpid()) {
			continue
		}
		if line != "p"+strconv.Itoa(pid) {
			return workflowError(4, "Codex thread lock has another or unknown holder; no input sent")
		}
		found = true
	}
	if !found {
		return workflowError(4, "Codex thread lock has no live holder; no input sent")
	}
	return nil
}

func observeCodexResumeIdentity(d Dependencies, s workflowState, a workflowAgent, deadline time.Time) (*CodexResumeIdentity, error) {
	if s.Request.Runtime.Provider != "codex" || s.OwnerResume == nil || a.Session != nil {
		return nil, nil
	}
	if err := workflowLocation(s, a); err != nil {
		return nil, err
	}
	record, _, err := workflowOwnerResumeRecord(s)
	if err != nil || record == nil {
		return nil, err
	}
	raw, err := workflowCallUntil(d, s.Request, deadline, "pane", "process-info", "--pane", a.PaneID)
	if err != nil {
		return nil, err
	}
	var response struct {
		Info struct {
			PaneID    string `json:"pane_id"`
			Processes []struct {
				PID  int      `json:"pid"`
				Name string   `json:"name"`
				CWD  string   `json:"cwd"`
				Argv []string `json:"argv"`
			} `json:"foreground_processes"`
		} `json:"process_info"`
	}
	if json.Unmarshal(raw, &response) != nil || response.Info.PaneID != a.PaneID {
		return nil, workflowError(4, "Codex resume pane process identity unavailable; no input sent")
	}
	pid := 0
	for _, process := range response.Info.Processes {
		if process.Name != "codex" {
			continue
		}
		if pid != 0 || process.PID <= 0 || process.CWD != s.Observed.Target.WorktreeLocator {
			return nil, workflowError(4, "Codex resume foreground process is ambiguous or outside the Task; no input sent")
		}
		// This is observed live argv, corroborated below by OS executable and
		// exclusive ownership. Requested launch arguments alone prove nothing.
		if len(process.Argv) < 3 || process.Argv[1] != "resume" || process.Argv[2] != record.NativeSessionID {
			return nil, workflowError(4, "live Codex process is not resuming the bound thread; no input sent")
		}
		pid = process.PID
	}
	if pid == 0 {
		return nil, workflowError(5, "Codex resume has no observed foreground process; no input sent")
	}
	executable := record.ProviderExecutable
	if err = verifyExecutable(executable); err != nil {
		return nil, err
	}
	command, err := codexIdentityCommand(deadline, "lsof", "-n", "-P", "-w", "-a", "-p", strconv.Itoa(pid), "-d", "txt", "-F", "n")
	if err != nil {
		return nil, err
	}
	loaded := false
	for _, line := range strings.Split(command, "\n") {
		if line == "n"+executable.Path {
			loaded = true
		}
	}
	if !loaded {
		return nil, workflowError(4, "live Codex process executable differs from the reserved provider; no input sent")
	}
	// Use the same home as the preserved transcript, including a non-default
	// home. Never infer it from a transcript's date layout or inspect contents.
	if _, err := deliveryProviderSession(s); err != nil {
		return nil, err
	}
	path := filepath.Join(ownerResumeCodexHome(s), "thread-writer-locks", record.NativeSessionID+".lock")
	if err = physical(path, false); err != nil {
		return nil, workflowError(5, "Codex thread writer lock unavailable; no input sent")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, workflowError(5, "Codex thread writer lock unreadable; no input sent")
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return nil, workflowError(4, "Codex thread lock is not a regular file; no input sent")
	}
	if err = codexLockHolder(deadline, path, pid); err != nil {
		return nil, err
	}
	err = unix.Flock(int(file.Fd()), unix.LOCK_SH|unix.LOCK_NB)
	if err == nil {
		_ = unix.Flock(int(file.Fd()), unix.LOCK_UN)
		return nil, workflowError(4, "Codex thread lock is not exclusively owned; no input sent")
	}
	if !errors.Is(err, unix.EWOULDBLOCK) && !errors.Is(err, unix.EAGAIN) {
		return nil, workflowError(5, "Codex thread lock ownership inspection failed; no input sent")
	}
	if err = codexLockHolder(deadline, path, pid); err != nil {
		return nil, err
	}
	current, err := os.Stat(path)
	if err != nil || !os.SameFile(info, current) {
		return nil, workflowError(4, "Codex thread lock changed during observation; no input sent")
	}
	return &CodexResumeIdentity{"local_codex_exclusive_thread_lock_v1", record.NativeSessionID, a.PaneID, pid, executable, path, d.Now().UTC().Format(time.RFC3339Nano)}, nil
}
