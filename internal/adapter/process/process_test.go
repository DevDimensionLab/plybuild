package process

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type recordingRunner struct {
	command Command
	err     error
	runs    int
}

type recordingReader struct {
	reads int
}

func (reader *recordingReader) Read([]byte) (int, error) {
	reader.reads++
	return 0, io.EOF
}

func (runner *recordingRunner) Run(command Command) error {
	runner.runs++
	runner.command = command
	return runner.err
}

func TestDependenciesDefaultToSafeNoOp(t *testing.T) {
	stdin := &recordingReader{}
	err := Execute(Dependencies{}, Command{
		Name:  "this-process-must-not-exist",
		Stdin: stdin,
	})
	if err != nil {
		t.Fatalf("zero-value process dependencies returned an error: %v", err)
	}
	if stdin.reads != 0 {
		t.Fatalf("zero-value process dependencies read stdin %d times", stdin.reads)
	}
}

func TestExecutePassesCompleteCommandToDependency(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	want := Command{
		Name:   "tool",
		Args:   []string{"first", "second", "complete dependency value"},
		Dir:    "/working/directory",
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
	sentinel := errors.New("recorded result")
	runner := &recordingRunner{err: sentinel}
	dependencies := Dependencies{Runner: runner}

	err := Execute(dependencies, want)

	if err != sentinel {
		t.Fatalf("Execute returned %v, want exact error %v", err, sentinel)
	}
	if runner.runs != 1 {
		t.Fatalf("Execute made %d process attempts, want 1", runner.runs)
	}
	if !reflect.DeepEqual(runner.command, want) {
		t.Fatalf("dependency received an incomplete command:\n got: %#v\nwant: %#v", runner.command, want)
	}
	if runner.command.Stdin != stdin {
		t.Fatal("stdin reader identity did not reach the process dependency unchanged")
	}
}

func TestSystemExecutesArgumentsInRequestedWorkingDirectory(t *testing.T) {
	workingDirectory := t.TempDir()
	resolvedWorkingDirectory, err := filepath.EvalSymlinks(workingDirectory)
	if err != nil {
		t.Fatalf("resolve working directory: %v", err)
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	stdin := strings.NewReader("complete stdin value\nwith a second line\n")
	command := Command{
		Name: "/bin/sh",
		Args: []string{
			"-c",
			`pwd -P; printf '%s\n' "$1"; cat; printf '%s\n' 'system diagnostic' >&2`,
			"process-adapter-test",
			"complete argument value",
		},
		Dir:    workingDirectory,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	err = Execute(System(), command)

	if err != nil {
		t.Fatalf("system process returned an error: %v", err)
	}
	wantStdout := resolvedWorkingDirectory + "\ncomplete argument value\ncomplete stdin value\nwith a second line\n"
	if stdout.String() != wantStdout {
		t.Fatalf("system stdout is %q, want %q", stdout.String(), wantStdout)
	}
	if stderr.String() != "system diagnostic\n" {
		t.Fatalf("system stderr is %q", stderr.String())
	}
}
