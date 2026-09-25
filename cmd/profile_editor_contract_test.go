package cmd

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/config"
)

type recordingProfileEditorProcess struct {
	commands []process.Command
	err      error
}

func (recording *recordingProfileEditorProcess) dependencies(
	dependencies profileEditorDependencies,
) profileEditorDependencies {
	dependencies.Process.Runner = recording
	return dependencies
}

func (recording *recordingProfileEditorProcess) Run(command process.Command) error {
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	return recording.err
}

func (recording *recordingProfileEditorProcess) assertedCommands() ([]process.Command, error) {
	if len(recording.commands) == 0 {
		return nil, errors.New("recorded profile-editor process population is empty")
	}
	return recording.commands, nil
}

func TestProfileEditorSelectsExactSystemRunnerWithoutProcessStandardOutput(t *testing.T) {
	dependencies := systemProfileEditorDependencies()
	systemProcess := process.System()

	if dependencies.Process.Runner == nil {
		t.Fatal("profile editor selected no system process runner")
	}
	if dependencies.Process.Runner != systemProcess.Runner {
		t.Fatalf("profile-editor process dependency is %T, want %T",
			dependencies.Process.Runner, systemProcess.Runner)
	}
	if dependencies.Process.Stdout != nil {
		t.Fatalf("profile editor inherited unused process standard output %T, want nil",
			dependencies.Process.Stdout)
	}
	if dependencies.Stdin != os.Stdin {
		t.Fatalf("profile-editor stdin is %T, want exact os.Stdin", dependencies.Stdin)
	}
	if dependencies.Stdout != os.Stdout {
		t.Fatalf("profile-editor stdout is %T, want exact os.Stdout", dependencies.Stdout)
	}
}

func TestProfileEditorPreservesExactEditorFallbackPathStreamsAttemptAndError(t *testing.T) {
	processError := errors.New("complete profile-editor process dependency error")
	tests := []struct {
		name        string
		editorValue string
		wantEditor  string
		processErr  error
	}{
		{
			name:        "non-empty EDITOR remains one exact executable",
			editorValue: `complete editor --must-remain-one-executable\segment-ø`,
			wantEditor:  `complete editor --must-remain-one-executable\segment-ø`,
		},
		{
			name:       "empty EDITOR falls back to vim",
			wantEditor: "vim",
			processErr: processError,
		},
	}

	if len(tests) == 0 {
		t.Fatal("TestProfileEditorPreservesExactEditorFallbackPathStreamsAttemptAndError test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("EDITOR", test.editorValue)
			recording := &recordingProfileEditorProcess{err: test.processErr}
			callerOwnedProcessStdout := &bytes.Buffer{}
			stdin := bytes.NewBufferString("complete profile-editor stdin bytes\n")
			stdout := &bytes.Buffer{}
			dependencies := recording.dependencies(profileEditorDependencies{
				Process: process.Dependencies{Stdout: callerOwnedProcessStdout},
				Stdin:   stdin,
				Stdout:  stdout,
			})
			if dependencies.Process.Runner != recording ||
				dependencies.Process.Stdout != callerOwnedProcessStdout ||
				dependencies.Stdin != stdin || dependencies.Stdout != stdout {
				t.Fatalf("profile-editor dependency lost its complete process value: %#v", dependencies)
			}
			localConfig := config.OpenLocalConfig(`/arbitrary profile directory//with spaces/../backslash\segment-β`)
			configPath := localConfig.FilePath()
			editor := selectProfileEditor(os.Getenv("EDITOR"))

			err := runProfileEditor(dependencies, editor, configPath)

			if err != test.processErr {
				t.Fatalf("profile-editor process error was %v, want exact dependency error %v", err, test.processErr)
			}
			commands, populationErr := recording.assertedCommands()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			want := []process.Command{{
				Name:   test.wantEditor,
				Args:   []string{configPath},
				Dir:    "",
				Stdin:  stdin,
				Stdout: stdout,
				Stderr: nil,
				Start:  false,
			}}
			if !reflect.DeepEqual(commands, want) {
				t.Fatalf("recorded profile-editor commands differ:\n got: %#v\nwant: %#v", commands, want)
			}
			if len(commands) != 1 {
				t.Fatalf("profile editor made %d process requests, want exactly 1", len(commands))
			}
			if len(commands[0].Args) != 1 || commands[0].Args[0] != configPath {
				t.Fatalf("profile editor changed its single exact config-path argument: %#v", commands[0].Args)
			}
			if commands[0].Stdin != stdin || commands[0].Stdout != stdout ||
				commands[0].Stdout == callerOwnedProcessStdout {
				t.Fatal("profile-editor dependency-selected command stream identities changed")
			}
		})
	}
}

func TestProfileEditorDependenciesDefaultToSafeNoProcess(t *testing.T) {
	err := runProfileEditor(
		profileEditorDependencies{},
		"this-profile-editor-process-must-not-exist",
		"/developer/home/.ply/must-not-be-accessed/local-config.yaml",
	)

	if err != nil {
		t.Fatalf("safe profile-editor process default returned an error: %v", err)
	}
}

func TestRecordedProfileEditorCommandsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingProfileEditorProcess{}

	if _, err := recording.assertedCommands(); err == nil {
		t.Fatal("empty recorded profile-editor process population passed")
	}
}
