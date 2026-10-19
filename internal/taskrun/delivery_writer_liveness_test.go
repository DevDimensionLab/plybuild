package taskrun

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// Real local processes (not modeled): the default observer must find a process
// by its cwd, by an open transcript and by a session argument, and never return
// command lines.
func TestSystemLocalProcessesObservesRealWriters(t *testing.T) {
	dir := t.TempDir()
	transcript := filepath.Join(dir, "session-0123.jsonl")
	if err := os.WriteFile(transcript, []byte("{}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	start := func(dir string, args ...string) *exec.Cmd {
		c := exec.Command("/bin/sh", args...)
		c.Dir = dir
		if err := c.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = c.Process.Kill(); _ = c.Wait() })
		return c
	}
	inCWD := start(dir, "-c", "sleep 30")
	holder := start("/", "-c", "exec 3<\"$0\"; sleep 30", transcript)
	named := start("/", "-c", "sleep 30; :", "resume01-session-fixture-0123")
	time.Sleep(300 * time.Millisecond)
	table, err := systemLocalProcesses("resume01-session-fixture-0123", transcript)
	if err != nil {
		t.Fatalf("local process inspection unavailable: %v", err)
	}
	byPID := map[int]LocalProcess{}
	for _, p := range table {
		byPID[p.PID] = p
	}
	physical, _ := filepath.EvalSymlinks(dir)
	if p := byPID[inCWD.Process.Pid]; p.CWD != physical && p.CWD != dir {
		t.Fatalf("cwd writer not observed: %+v", p)
	}
	if p := byPID[holder.Process.Pid]; p.Session != "transcript_open" {
		t.Fatalf("transcript holder not observed: %+v", p)
	}
	if p := byPID[named.Process.Pid]; p.Session != "session_argument" {
		t.Fatalf("session argument not observed: %+v", p)
	}
	if p := byPID[os.Getpid()]; p.PID == 0 {
		t.Fatal("the observing process itself is missing, so exclusions cannot be applied")
	}
}
