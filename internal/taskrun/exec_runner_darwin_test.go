//go:build darwin

package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFactoryPipeRunnerAndTimeout(t *testing.T) {
	for _, timeout := range []bool{false, true} {
		t.Run(map[bool]string{false: "complete", true: "timeout"}[timeout], func(t *testing.T) {
			root, e := filepath.EvalSymlinks(t.TempDir())
			if e != nil {
				t.Fatal(e)
			}
			script := "#!/bin/sh\nsleep 0.05\ncat <<'EVENTS'\n" + string(completedEvents) + "EVENTS\n"
			limit := 2 * time.Second
			if timeout {
				script = "#!/bin/sh\nexec sleep 30\n"
				limit = 50 * time.Millisecond
			}
			path := filepath.Join(root, "synthetic-provider")
			if e = os.WriteFile(path, []byte(script), 0700); e != nil {
				t.Fatal(e)
			}
			completion := &ProviderCompletion{}
			p, e := systemExecRunner().Run(LaunchSpec{Executable: Executable{path, hashFileTest(t, path)}, Argv: []string{path}, CWD: root, Timeout: limit, StreamsRoot: filepath.Join(root, "streams"), Completion: completion}, func(Process) error { return nil })
			if !timeout && (e != nil || !providerInactive(*completion)) {
				t.Fatalf("%v %+v %+v", e, p, completion)
			}
			if timeout && (!completion.TimedOut || providerInactive(*completion) || !quiescent(p)) {
				t.Fatalf("timeout did not remain nonqualifying and reaped: %v %+v %+v", e, p, completion)
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
