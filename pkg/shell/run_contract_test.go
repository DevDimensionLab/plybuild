package shell

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type recordingShellRunProcess struct {
	commands []process.Command
	stdout   []byte
	stderr   []byte
	err      error
}

func (recording *recordingShellRunProcess) dependencies() runDependencies {
	return runDependencies{Process: process.Dependencies{Runner: recording}}
}

func (recording *recordingShellRunProcess) Run(command process.Command) error {
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	if command.Stdout != nil {
		_, _ = command.Stdout.Write(recording.stdout)
	}
	if command.Stderr != nil {
		_, _ = command.Stderr.Write(recording.stderr)
	}
	return recording.err
}

func (recording *recordingShellRunProcess) assertedCommands() ([]process.Command, error) {
	if len(recording.commands) == 0 {
		return nil, errors.New("recorded shell Run process population is empty")
	}
	return recording.commands, nil
}

func TestRunSelectsCompleteSystemProcessDependencies(t *testing.T) {
	dependencies := systemRunDependencies()
	systemProcess := process.System()

	if dependencies.Process.Runner == nil {
		t.Fatal("Run selected an incomplete process dependency")
	}
	if reflect.TypeOf(dependencies.Process.Runner) != reflect.TypeOf(systemProcess.Runner) {
		t.Fatalf("Run process dependency is %T, want %T",
			dependencies.Process.Runner, systemProcess.Runner)
	}
}

func TestRunPreservesExactRequestStreamsSynchronousAttemptBytesAndLegacyErrorResult(t *testing.T) {
	name := string([]byte{'c', 'o', 'm', 'p', 'l', 'e', 't', 'e', ' ', 't', 'o', 'o', 'l', '\\', 0xc3, 0xb8})
	arguments := []string{
		string([]byte{'f', 'i', 'r', 's', 't', 0x00, 'v', 'a', 'l', 'u', 'e'}),
		"--second=complete value",
		string([]byte{0xff, '\n', '\\', 'l', 'a', 's', 't'}),
	}
	stdoutBytes := []byte{'s', 't', 'd', 'o', 'u', 't', 0x00, 0xff, '\n'}
	stderrBytes := []byte{'s', 't', 'd', 'e', 'r', 'r', 0x00, 0xfe, '\n'}
	processError := errors.New("complete shell Run process dependency error")
	recording := &recordingShellRunProcess{
		stdout: stdoutBytes,
		stderr: stderrBytes,
		err:    processError,
	}
	dependencies := recording.dependencies()
	if dependencies.Process.Runner != recording {
		t.Fatalf("Run dependency lost its complete process value: %#v", dependencies)
	}

	output := runWithDependencies(dependencies, name, arguments...)

	commands, populationErr := recording.assertedCommands()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(commands) != 1 {
		t.Fatalf("Run made %d process requests, want exactly 1", len(commands))
	}
	command := commands[0]
	if command.Name != name || !reflect.DeepEqual(command.Args, arguments) {
		t.Fatalf("Run process request changed name or ordered argument bytes:\n got: %#v\nwant: name=%q args=%#v",
			command, name, arguments)
	}
	if command.Dir != "" || command.Stdin != nil || command.Start {
		t.Fatalf("Run process request changed synchronous defaults: %#v", command)
	}
	stdout, stdoutOK := command.Stdout.(*bytes.Buffer)
	stderr, stderrOK := command.Stderr.(*bytes.Buffer)
	if !stdoutOK || !stderrOK || stdout == stderr {
		t.Fatalf("Run process streams were stdout=%T stderr=%T with identities %p and %p",
			command.Stdout, command.Stderr, stdout, stderr)
	}
	if !bytes.Equal(output.StdOut.Bytes(), stdoutBytes) || !bytes.Equal(output.StdErr.Bytes(), stderrBytes) {
		t.Fatalf("Run output changed delivered bytes: stdout=%v stderr=%v", output.StdOut.Bytes(), output.StdErr.Bytes())
	}
	if output.Err != nil {
		t.Fatalf("legacy Run Output.Err was %v after exact dependency error %v, want nil", output.Err, processError)
	}
}

func TestRunDependenciesDefaultToSafeNoProcess(t *testing.T) {
	output := runWithDependencies(
		runDependencies{},
		"this-shell-run-process-must-not-exist",
		"complete argument",
	)

	if output.Err != nil || output.StdOut.Len() != 0 || output.StdErr.Len() != 0 {
		t.Fatalf("safe Run process default returned unexpected output: %#v", output)
	}
}

func TestRecordedShellRunCommandsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingShellRunProcess{}

	if _, err := recording.assertedCommands(); err == nil {
		t.Fatal("empty recorded shell Run process population passed")
	}
}
