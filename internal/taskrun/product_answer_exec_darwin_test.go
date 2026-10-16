//go:build darwin

package taskrun

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSL01OwnedRunnerHonorsFinalStartGuard(t *testing.T) {
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	shell, e := readFile("/bin/sh", 16<<20, false)
	if e != nil {
		t.Fatal(e)
	}
	marker := filepath.Join(root, "unexpected-start")
	completion := ProviderCompletion{}
	entered := 0
	blocked := errors.New("parent pause won at final Start")
	p, e := workerOwnedExec(LaunchSpec{Executable: Executable{"/bin/sh", hash(shell)}, Argv: []string{"/bin/sh", "-c", `printf started > "$1"`, "stand-in", marker}, CWD: root, StreamsRoot: filepath.Join(root, "streams"), Timeout: time.Second, Completion: &completion, StartGuard: func(start func() error) error { entered++; return blocked }}, func(Process) error { return nil })
	if entered != 1 || !errors.Is(e, blocked) || p.State != "not_started" {
		t.Fatalf("final gate bypassed: calls=%d process=%+v err=%v", entered, p, e)
	}
	if _, e = os.Stat(marker); !os.IsNotExist(e) {
		t.Fatal("provider stand-in was started despite rejected final gate")
	}
	// A rejected guard never entered the OS-start closure. A later allowed
	// dispatch of this same reserved stream destination must still be possible.
	p, e = workerOwnedExec(LaunchSpec{Executable: Executable{"/bin/sh", hash(shell)}, Argv: []string{"/bin/sh", "-c", `printf started > "$1"`, "stand-in", marker}, CWD: root, StreamsRoot: filepath.Join(root, "streams"), Timeout: time.Second, Completion: &completion, StartGuard: func(start func() error) error { entered++; return start() }}, func(Process) error { return nil })
	if entered != 2 || e != nil || p.State != "exited" || p.ExitCode == nil || *p.ExitCode != 0 {
		t.Fatalf("allowed same-intent resume blocked: calls=%d process=%+v err=%v", entered, p, e)
	}
	if raw, e := os.ReadFile(marker); e != nil || string(raw) != "started" {
		t.Fatalf("allowed stand-in did not run: %q %v", raw, e)
	}
}

func TestSL01OwnedRunnerCapsAReservedTurnAtItsAbsoluteDeadline(t *testing.T) {
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	raw, e := readFile("/bin/sh", 16<<20, false)
	if e != nil {
		t.Fatal(e)
	}
	completion := ProviderCompletion{}
	before := time.Now()
	p, e := workerOwnedExec(LaunchSpec{Executable: Executable{"/bin/sh", hash(raw)}, Argv: []string{"/bin/sh", "-c", "sleep 2"}, CWD: root, StreamsRoot: filepath.Join(root, "streams"), Timeout: 3 * time.Second, Deadline: before.Add(100 * time.Millisecond), Completion: &completion}, func(Process) error { return nil })
	if e == nil || !completion.TimedOut || time.Since(before) > time.Second || p.State != "exited" || !quiescent(p) {
		t.Fatalf("absolute deadline did not cap reserved relative timeout: elapsed=%s process=%+v completion=%+v err=%v", time.Since(before), p, completion, e)
	}
}
