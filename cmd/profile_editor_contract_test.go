package cmd

import (
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

func (recording *recordingProfileEditorProcess) dependencies() profileEditorDependencies {
	dependencies := systemProfileEditorDependencies()
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

func TestProfileEditorSelectsCompleteSystemProcessDependencies(t *testing.T) {
	dependencies := systemProfileEditorDependencies()
	systemProcess := process.System()

	if dependencies.Process.Runner == nil || dependencies.Process.Stdout != systemProcess.Stdout {
		t.Fatal("profile editor selected an incomplete process dependency")
	}
	if reflect.TypeOf(dependencies.Process.Runner) != reflect.TypeOf(systemProcess.Runner) {
		t.Fatalf("profile-editor process dependency is %T, want %T",
			dependencies.Process.Runner, systemProcess.Runner)
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("EDITOR", test.editorValue)
			recording := &recordingProfileEditorProcess{err: test.processErr}
			dependencies := recording.dependencies()
			if dependencies.Process.Runner != recording {
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
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: nil,
			}}
			if !reflect.DeepEqual(commands, want) {
				t.Fatalf("recorded profile-editor commands differ:\n got: %#v\nwant: %#v", commands, want)
			}
			if commands[0].Stdin != os.Stdin || commands[0].Stdout != os.Stdout {
				t.Fatal("profile-editor standard stream identities changed")
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
