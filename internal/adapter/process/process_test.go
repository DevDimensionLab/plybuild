package process

import (
	"bytes"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

type recordingRunner struct {
	command Command
	err     error
}

func (runner *recordingRunner) Run(command Command) error {
	runner.command = command
	return runner.err
}

func TestDependenciesDefaultToSafeNoOp(t *testing.T) {
	err := Execute(Dependencies{}, Command{Name: "this-process-must-not-exist"})
	if err != nil {
		t.Fatalf("zero-value process dependencies returned an error: %v", err)
	}
}

func TestExecutePassesCompleteCommandToDependency(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	want := Command{
		Name:   "tool",
		Args:   []string{"first", "second", "complete dependency value"},
		Dir:    "/working/directory",
		Stdout: stdout,
		Stderr: stderr,
	}
	sentinel := errors.New("recorded result")
	runner := &recordingRunner{err: sentinel}
	dependencies := Dependencies{Runner: runner}

	err := Execute(dependencies, want)

	if !errors.Is(err, sentinel) {
		t.Fatalf("Execute returned %v, want %v", err, sentinel)
	}
	if !reflect.DeepEqual(runner.command, want) {
		t.Fatalf("dependency received an incomplete command:\n got: %#v\nwant: %#v", runner.command, want)
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
	command := Command{
		Name: "/bin/sh",
		Args: []string{
			"-c",
			`pwd -P; printf '%s\n' "$1"; printf '%s\n' 'system diagnostic' >&2`,
			"process-adapter-test",
			"complete argument value",
		},
		Dir:    workingDirectory,
		Stdout: stdout,
		Stderr: stderr,
	}

	err = Execute(System(), command)

	if err != nil {
		t.Fatalf("system process returned an error: %v", err)
	}
	wantStdout := resolvedWorkingDirectory + "\ncomplete argument value\n"
	if stdout.String() != wantStdout {
		t.Fatalf("system stdout is %q, want %q", stdout.String(), wantStdout)
	}
	if stderr.String() != "system diagnostic\n" {
		t.Fatalf("system stderr is %q", stderr.String())
	}
}
