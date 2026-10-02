//go:build darwin

package taskrun

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestFactoryPipeRunnerAndTimeout(t *testing.T) {
	for _, tc := range []struct {
		name    string
		delay   string
		timeout bool
	}{
		{name: "complete"},
		// Reproduce a slow synthetic startup before the first provider event.
		// Completion correctness must not depend on a two-second startup SLA.
		{name: "delayed_complete", delay: "/bin/sleep 2.1\n"},
		{name: "timeout", timeout: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, e := filepath.EvalSymlinks(t.TempDir())
			if e != nil {
				t.Fatal(e)
			}
			// printf is a shell builtin: the normal completion fixture needs no
			// extra sleep/cat processes before it can emit its terminal events.
			script := "#!/bin/sh\n" + tc.delay + "printf '%s' '" + string(completedEvents) + "'\n"
			// Allow bounded fixture startup slack; the separate timeout case
			// below still exercises an actual short runner deadline.
			limit := 15 * time.Second
			if tc.timeout {
				script = "#!/bin/sh\n: > ready\nexec /bin/sleep 30\n"
				limit = 100 * time.Millisecond
			}
			path := filepath.Join(root, "synthetic-provider")
			if e = os.WriteFile(path, []byte(script), 0700); e != nil {
				t.Fatal(e)
			}
			completion := &ProviderCompletion{}
			began := time.Now()
			var launched, ready time.Time
			p, e := systemExecRunner().Run(LaunchSpec{Executable: Executable{path, hashFileTest(t, path)}, Argv: []string{path}, CWD: root, Timeout: limit, StreamsRoot: filepath.Join(root, "streams"), Completion: completion}, func(Process) error {
				launched = time.Now()
				if !tc.timeout {
					return nil
				}
				// The runner starts its timeout after this callback. Observe actual
				// fixture startup first, so timeout tests a running provider.
				deadline := time.Now().Add(10 * time.Second)
				for {
					if _, err := os.Stat(filepath.Join(root, "ready")); err == nil {
						ready = time.Now()
						return nil
					} else if !os.IsNotExist(err) {
						return err
					}
					if time.Now().After(deadline) {
						return fmt.Errorf("synthetic provider did not become ready")
					}
					time.Sleep(10 * time.Millisecond)
				}
			})
			elapsed := time.Since(began)
			t.Logf("launch=%s total=%s limit=%s ready=%t", launched.Sub(began), elapsed, limit, !ready.IsZero())
			if p.State != "exited" || p.PID == nil || !quiescent(p) || !quiescent(completion.Process) {
				t.Fatalf("provider not reaped and quiescent: %v %+v %+v", e, p, completion)
			}
			var status unix.WaitStatus
			if _, err := unix.Wait4(*p.PID, &status, unix.WNOHANG, nil); !errors.Is(err, unix.ECHILD) {
				t.Fatalf("runner left an unreaped child: %v", err)
			}
			if !tc.timeout {
				if e != nil || !providerInactive(*completion) {
					t.Fatalf("completion did not qualify: %v %+v %+v", e, p, completion)
				}
				raw, err := os.ReadFile(completion.Stdout.Locator)
				if err != nil || string(raw) != string(completedEvents) || completion.Stdout.SHA256 != hash(completedEvents) || completion.Stderr.SHA256 != hash(nil) {
					t.Fatalf("incomplete or unbound provider output: %v %+v", err, completion)
				}
			} else if ready.IsZero() || time.Since(ready) < limit || e == nil || !completion.TimedOut || providerInactive(*completion) || p.Signal == nil || p.ExitCode != nil || !completion.StreamComplete || completion.SequenceValid {
				t.Fatalf("running provider timeout did not remain nonqualifying and reaped: %v %+v %+v", e, p, completion)
			}
		})
	}
}
func TestFactoryPipeStreamLimit(t *testing.T) {
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(root, "stream")
	f, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal(e)
	}
	rd, wr, e := os.Pipe()
	if e != nil {
		t.Fatal(e)
	}
	done := make(chan streamCopy, 1)
	go copyProviderStream(f, rd, done)
	_, e = wr.Write([]byte(strings.Repeat("x", providerStreamLimit+1)))
	if e != nil {
		t.Fatal(e)
	}
	wr.Close()
	result := <-done
	info, _ := os.Stat(path)
	if !result.truncated || info.Size() != providerStreamLimit || !result.complete {
		t.Fatal(result, info.Size())
	}
}

func TestReturnClosureRunnerReapsAfterCorrectedCommand(t *testing.T) {
	raw, err := os.ReadFile("testdata/qa_v4_return_events.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "synthetic-provider")
	if err = os.WriteFile(path, []byte("#!/bin/sh\nprintf '%s' "+ShellQuote(string(raw))+"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	completion := &ProviderCompletion{}
	p, err := systemExecRunner().Run(LaunchSpec{Executable: Executable{path, hashFileTest(t, path)}, Argv: []string{path}, CWD: root, Timeout: 15 * time.Second, StreamsRoot: filepath.Join(root, "streams"), Completion: completion}, func(Process) error { return nil })
	if err != nil || p.PID == nil || !providerInactive(*completion) {
		t.Fatalf("closed command failure did not yield known completion: %v %+v", err, completion)
	}
	var status unix.WaitStatus
	if _, err = unix.Wait4(*p.PID, &status, unix.WNOHANG, nil); !errors.Is(err, unix.ECHILD) {
		t.Fatalf("child not reaped: %v", err)
	}
	preserved, err := os.ReadFile(completion.Stdout.Locator)
	if err != nil || string(preserved) != string(raw) || completion.Stdout.SHA256 != hash(raw) {
		t.Fatal("tool failure data changed")
	}
}
