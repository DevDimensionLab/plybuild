package process

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
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
		Start:  true,
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

func TestSystemFalseStartModePerformsOneSynchronousRun(t *testing.T) {
	command := Command{
		Name:  "/bin/sh",
		Args:  []string{"-c", "exit 73"},
		Start: false,
	}

	err := Execute(System(), command)

	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("false start mode returned %T %v, want one synchronous *exec.ExitError", err, err)
	}
	if exitError.ExitCode() != 73 {
		t.Fatalf("false start mode exit code was %d, want 73", exitError.ExitCode())
	}
}

func TestSystemTrueStartModePerformsOneAsynchronousStart(t *testing.T) {
	command := Command{
		Name:  "/bin/sh",
		Args:  []string{"-c", "exit 74"},
		Start: true,
	}

	err := Execute(System(), command)

	if err != nil {
		t.Fatalf("true start mode waited for the child result or returned a start error: %v", err)
	}
}

func TestSystemTrueStartModeReturnsExactStartError(t *testing.T) {
	missingDirectory := filepath.Join(t.TempDir(), "missing working directory")
	command := Command{
		Name:  "/bin/sh",
		Args:  []string{"-c", "exit 0"},
		Dir:   missingDirectory,
		Start: true,
	}

	err := Execute(System(), command)

	pathError, ok := err.(*os.PathError)
	if !ok {
		t.Fatalf("true start mode returned %T %v, want the direct start *os.PathError", err, err)
	}
	if pathError.Op != "chdir" || pathError.Path != missingDirectory || !errors.Is(pathError.Err, os.ErrNotExist) {
		t.Fatalf("true start mode changed the direct start error: %#v", pathError)
	}
}
