package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/sirupsen/logrus"
)

const (
	gitProcessMutation = "1. git clone keeps URL before target directory"
	gitCommitMutation  = "7. git commit keeps target directory before message"
)

type recordedProcessCall struct {
	Name string
	Args []string
	Dir  string
}

type recordingProcess struct {
	commands []process.Command
	results  []recordedGitProcessResult
	sequence *[]string
}

type recordedGitProcessResult struct {
	stdout []byte
	stderr []byte
	err    error
}

type recordingGitLogHook struct {
	entries  []*logrus.Entry
	sequence *[]string
}

func (recording *recordingProcess) dependencies(dependencies process.Dependencies) process.Dependencies {
	dependencies.Runner = recording
	return dependencies
}

func (recording *recordingProcess) Run(command process.Command) error {
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "process")
	}
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	index := len(recording.commands) - 1
	if index >= len(recording.results) {
		return nil
	}
	result := recording.results[index]
	if command.Stdout != nil {
		_, _ = command.Stdout.Write(result.stdout)
	}
	if command.Stderr != nil {
		_, _ = command.Stderr.Write(result.stderr)
	}
	return result.err
}

func (recording *recordingProcess) assertedCalls() ([]recordedProcessCall, error) {
	commands, err := recording.assertedCommands()
	if err != nil {
		return nil, err
	}
	calls := make([]recordedProcessCall, 0, len(commands))
	for _, command := range commands {
		calls = append(calls, recordedProcessCall{
			Name: command.Name,
			Args: append([]string(nil), command.Args...),
			Dir:  command.Dir,
		})
	}
	return calls, nil
}

func (recording *recordingProcess) assertedCommands() ([]process.Command, error) {
	if len(recording.commands) == 0 {
		return nil, errors.New("recorded process call population is empty")
	}
	return recording.commands, nil
}

func (hook *recordingGitLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *recordingGitLogHook) Fire(entry *logrus.Entry) error {
	if hook.sequence != nil {
		*hook.sequence = append(*hook.sequence, "log")
	}
	copy := *entry
	hook.entries = append(hook.entries, &copy)
	return nil
}

func TestGitSelectsExactSystemRunnerWithoutStandardOutput(t *testing.T) {
	dependencies := systemGitDependencies()
	systemProcess := process.System()

	if dependencies.Runner == nil {
		t.Fatal("Git selected no system process runner")
	}
	if dependencies.Runner != systemProcess.Runner {
		t.Fatalf("Git process dependency is %T, want exact runner %T",
			dependencies.Runner, systemProcess.Runner)
	}
	if dependencies.Stdout != nil {
		t.Fatalf("Git inherited unused standard output %T, want nil", dependencies.Stdout)
	}
}

func TestGitCloneKeepsURLBeforeTargetDirectory(t *testing.T) {
	if gitProcessMutation == "" {
		t.Fatal("git-process mutation label is empty")
	}
	recording := &recordingProcess{results: []recordedGitProcessResult{{
		stdout: []byte("cloned\n"), stderr: []byte("git diagnostic\n"),
	}}}
	url := "ssh://git@example.invalid/team/repository.git"
	target := filepath.Join(t.TempDir(), "complete target directory")

	output := gitClone(recording.dependencies(process.Dependencies{}), url, target)

	assertRecordedProcessCalls(t, recording, []recordedProcessCall{{
		Name: "git",
		Args: []string{"clone", url, target},
		Dir:  "",
	}})
	if output.StdOut.String() != "cloned\n" || output.StdErr.String() != "git diagnostic\n" {
		t.Fatalf("shell output changed: stdout %q, stderr %q", output.StdOut.String(), output.StdErr.String())
	}
}

func TestGitPullAndInitKeepExactArgumentsAndWorkingDirectory(t *testing.T) {
	target := filepath.Join(t.TempDir(), "repository")
	tests := []struct {
		name string
		run  func(process.Dependencies) Output
		want recordedProcessCall
	}{
		{
			name: "pull",
			run: func(dependencies process.Dependencies) Output {
				return gitPull(dependencies, target)
			},
			want: recordedProcessCall{Name: "git", Args: []string{"-C", target, "pull", "origin"}},
		},
		{
			name: "init",
			run: func(dependencies process.Dependencies) Output {
				return gitInit(dependencies, target)
			},
			want: recordedProcessCall{Name: "git", Args: []string{"-C", target, "init"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingProcess{}
			test.run(recording.dependencies(process.Dependencies{}))
			assertRecordedProcessCalls(t, recording, []recordedProcessCall{test.want})
		})
	}
}

func TestGitAddAndCommitKeepsTargetBeforeCompleteMessage(t *testing.T) {
	if gitCommitMutation == "" {
		t.Fatal("git-commit mutation label is empty")
	}
	recording := &recordingProcess{}
	target := filepath.Join(t.TempDir(), "repository with spaces")
	message := "complete dependency value: com.example:library:1.2.3"

	gitAddAndCommit(recording.dependencies(process.Dependencies{}), target, message)

	assertRecordedProcessCalls(t, recording, []recordedProcessCall{
		{Name: "git", Args: []string{"-C", target, "add", "."}},
		{Name: "git", Args: []string{"-C", target, "commit", "-m", fmt.Sprintf("\"%s\"", message)}},
	})
}

func TestGitWrappersPreserveExactRequestsLogsStreamsBytesAttemptsAndLegacyErrors(t *testing.T) {
	url := string([]byte{'s', 's', 'h', ':', '/', '/', 'g', 'i', 't', '@', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'i', 'n', 'v', 'a', 'l', 'i', 'd', '/', 0xc3, 0xb8, '\\', 0x00})
	target := string([]byte{'/', 'a', 'r', 'b', 'i', 't', 'r', 'a', 'r', 'y', ' ', 'g', 'i', 't', ' ', 't', 'a', 'r', 'g', 'e', 't', '/', 0xff, '\\', '\n'})
	message := string([]byte{'c', 'o', 'm', 'p', 'l', 'e', 't', 'e', ' ', 'm', 'e', 's', 's', 'a', 'g', 'e', ' ', 0xc3, 0xb8, '\\', 0x00, '\n'})
	oneRequestError := errors.New("complete one-request Git process dependency error")
	addError := errors.New("complete Git add process dependency error")
	commitError := errors.New("complete Git commit process dependency error")
	tests := []struct {
		name        string
		run         func(process.Dependencies) Output
		wantArgs    [][]string
		results     []recordedGitProcessResult
		outputIndex int
	}{
		{
			name: "clone",
			run: func(dependencies process.Dependencies) Output {
				return gitClone(dependencies, url, target)
			},
			wantArgs: [][]string{{"clone", url, target}},
			results: []recordedGitProcessResult{{
				stdout: []byte{'c', 'l', 'o', 'n', 'e', 0x00, 0xff, '\n'},
				stderr: []byte{'c', 'l', 'o', 'n', 'e', ' ', 'e', 'r', 'r', 0x00, 0xfe, '\n'},
				err:    oneRequestError,
			}},
		},
		{
			name: "pull",
			run: func(dependencies process.Dependencies) Output {
				return gitPull(dependencies, target)
			},
			wantArgs: [][]string{{"-C", target, "pull", "origin"}},
			results: []recordedGitProcessResult{{
				stdout: []byte{'p', 'u', 'l', 'l', 0x00, 0xfd, '\n'},
				stderr: []byte{'p', 'u', 'l', 'l', ' ', 'e', 'r', 'r', 0x00, 0xfc, '\n'},
				err:    oneRequestError,
			}},
		},
		{
			name: "init",
			run: func(dependencies process.Dependencies) Output {
				return gitInit(dependencies, target)
			},
			wantArgs: [][]string{{"-C", target, "init"}},
			results: []recordedGitProcessResult{{
				stdout: []byte{'i', 'n', 'i', 't', 0x00, 0xfb, '\n'},
				stderr: []byte{'i', 'n', 'i', 't', ' ', 'e', 'r', 'r', 0x00, 0xfa, '\n'},
				err:    oneRequestError,
			}},
		},
		{
			name: "add and commit",
			run: func(dependencies process.Dependencies) Output {
				return gitAddAndCommit(dependencies, target, message)
			},
			wantArgs: [][]string{
				{"-C", target, "add", "."},
				{"-C", target, "commit", "-m", fmt.Sprintf("\"%s\"", message)},
			},
			results: []recordedGitProcessResult{
				{
					stdout: []byte{'a', 'd', 'd', 0x00, 0xf9, '\n'},
					stderr: []byte{'a', 'd', 'd', ' ', 'e', 'r', 'r', 0x00, 0xf8, '\n'},
					err:    addError,
				},
				{
					stdout: []byte{'c', 'o', 'm', 'm', 'i', 't', 0x00, 0xf7, '\n'},
					stderr: []byte{'c', 'o', 'm', 'm', 'i', 't', ' ', 'e', 'r', 'r', 0x00, 0xf6, '\n'},
					err:    commitError,
				},
			},
			outputIndex: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sequence := []string{}
			hook := &recordingGitLogHook{sequence: &sequence}
			testLogger := logrus.New()
			testLogger.SetOutput(io.Discard)
			testLogger.SetLevel(logrus.DebugLevel)
			testLogger.AddHook(hook)
			previousLog := log
			log = testLogger
			t.Cleanup(func() {
				log = previousLog
			})
			recording := &recordingProcess{results: test.results, sequence: &sequence}
			callerOwnedStdout := &bytes.Buffer{}
			dependencies := recording.dependencies(process.Dependencies{Stdout: callerOwnedStdout})
			if dependencies.Runner != recording || dependencies.Stdout != callerOwnedStdout {
				t.Fatalf("Git dependency lost its complete process value: %#v", dependencies)
			}

			output := test.run(dependencies)

			commands, populationErr := recording.assertedCommands()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if len(commands) != len(test.wantArgs) {
				t.Fatalf("Git made %d process requests, want exactly %d", len(commands), len(test.wantArgs))
			}
			streamIdentities := make(map[*bytes.Buffer]bool)
			for index, command := range commands {
				if command.Name != "git" || !reflect.DeepEqual(command.Args, test.wantArgs[index]) {
					t.Fatalf("Git process request %d changed executable or ordered argument bytes:\n got: %#v\nwant: name=%q args=%#v",
						index, command, "git", test.wantArgs[index])
				}
				if command.Dir != "" || command.Stdin != nil || command.Start {
					t.Fatalf("Git process request %d changed synchronous defaults: %#v", index, command)
				}
				stdout, stdoutOK := command.Stdout.(*bytes.Buffer)
				stderr, stderrOK := command.Stderr.(*bytes.Buffer)
				if !stdoutOK || !stderrOK || stdout == stderr || streamIdentities[stdout] || streamIdentities[stderr] {
					t.Fatalf("Git process request %d streams were stdout=%T stderr=%T with reused identities %p and %p",
						index, command.Stdout, command.Stderr, stdout, stderr)
				}
				streamIdentities[stdout] = true
				streamIdentities[stderr] = true
				if !bytes.Equal(stdout.Bytes(), test.results[index].stdout) ||
					!bytes.Equal(stderr.Bytes(), test.results[index].stderr) {
					t.Fatalf("Git process request %d stream bytes changed: stdout=%v stderr=%v",
						index, stdout.Bytes(), stderr.Bytes())
				}
			}
			wantSequence := make([]string, 0, len(commands)*2)
			for range commands {
				wantSequence = append(wantSequence, "log", "process")
			}
			if !reflect.DeepEqual(sequence, wantSequence) {
				t.Fatalf("Git log/process order was %#v, want %#v", sequence, wantSequence)
			}
			if len(hook.entries) != len(commands) {
				t.Fatalf("Git emitted %d log entries, want exactly %d", len(hook.entries), len(commands))
			}
			for index, entry := range hook.entries {
				wantMessage := "running: git " + strings.Join(test.wantArgs[index], " ")
				if entry.Level != logrus.DebugLevel || entry.Message != wantMessage {
					t.Fatalf("Git log entry %d was level=%s message=%q, want debug %q",
						index, entry.Level, entry.Message, wantMessage)
				}
			}
			wantOutput := test.results[test.outputIndex]
			if !bytes.Equal(output.StdOut.Bytes(), wantOutput.stdout) ||
				!bytes.Equal(output.StdErr.Bytes(), wantOutput.stderr) {
				t.Fatalf("Git returned output changed delivered bytes: stdout=%v stderr=%v",
					output.StdOut.Bytes(), output.StdErr.Bytes())
			}
			if output.Err != nil {
				t.Fatalf("legacy Git Output.Err was %v after direct dependency errors, want nil", output.Err)
			}
		})
	}
}

func TestGitProcessDependenciesDefaultToNoProcess(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe-default-repository")
	tests := []struct {
		name string
		run  func() Output
	}{
		{name: "clone", run: func() Output { return gitClone(process.Dependencies{}, "must-not-clone", target) }},
		{name: "pull", run: func() Output { return gitPull(process.Dependencies{}, target) }},
		{name: "init", run: func() Output { return gitInit(process.Dependencies{}, target) }},
		{name: "add and commit", run: func() Output {
			return gitAddAndCommit(process.Dependencies{}, target, "must not commit")
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := test.run()
			if output.Err != nil || output.StdOut.Len() != 0 || output.StdErr.Len() != 0 {
				t.Fatalf("safe Git process default returned unexpected output: %#v", output)
			}
		})
	}
	if _, err := os.Stat(filepath.Join(target, ".git")); !os.IsNotExist(err) {
		t.Fatalf("safe process default created git state: %v", err)
	}
}

func TestRecordedGitCallsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingProcess{}

	if _, err := recording.assertedCalls(); err == nil {
		t.Fatal("empty recorded process population passed")
	}
}

func assertRecordedProcessCalls(t *testing.T, recording *recordingProcess, want []recordedProcessCall) {
	t.Helper()
	calls, err := recording.assertedCalls()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("recorded process calls differ:\n got: %#v\nwant: %#v", calls, want)
	}
}
