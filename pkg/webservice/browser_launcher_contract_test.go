package webservice

import (
	"errors"
	"reflect"
	"runtime"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type recordingBrowserLauncherProcess struct {
	commands []process.Command
	err      error
}

func (recording *recordingBrowserLauncherProcess) dependencies(goos string) browserLauncherDependencies {
	processDependencies := process.System()
	processDependencies.Runner = recording
	return browserLauncherDependencies{
		Process: processDependencies,
		GOOS:    goos,
	}
}

func (recording *recordingBrowserLauncherProcess) Run(command process.Command) error {
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	return recording.err
}

func (recording *recordingBrowserLauncherProcess) assertedCommands() ([]process.Command, error) {
	if len(recording.commands) == 0 {
		return nil, errors.New("recorded browser-launcher process population is empty")
	}
	return recording.commands, nil
}

func TestOpenBrowserSelectsCompleteSystemProcessAndRuntimePlatform(t *testing.T) {
	dependencies := systemBrowserLauncherDependencies()
	systemProcess := process.System()

	if dependencies.Process.Runner == nil || dependencies.Process.Stdout != systemProcess.Stdout {
		t.Fatal("OpenBrowser selected an incomplete process dependency")
	}
	if reflect.TypeOf(dependencies.Process.Runner) != reflect.TypeOf(systemProcess.Runner) {
		t.Fatalf("OpenBrowser process dependency is %T, want %T",
			dependencies.Process.Runner, systemProcess.Runner)
	}
	if dependencies.GOOS != runtime.GOOS {
		t.Fatalf("OpenBrowser platform is %q, want runtime.GOOS %q", dependencies.GOOS, runtime.GOOS)
	}
}

func TestBrowserLauncherPreservesEveryPlatformRequestURLStartAttemptAndError(t *testing.T) {
	url := string([]byte{
		'h', 't', 't', 'p', 's', ':', '/', '/', 'e', 'x', 'a', 'm', 'p', 'l', 'e', '.', 'i', 'n', 'v', 'a', 'l', 'i', 'd', '/',
		's', 'p', 'a', 'c', 'e', ' ', '\\', 0x00, '?', 'q', '=', '%', '2', 'F', '&', 0xc3, 0xb8, '\n',
	})
	tests := []struct {
		name       string
		goos       string
		processErr error
		want       process.Command
	}{
		{
			name:       "linux xdg-open",
			goos:       "linux",
			processErr: errors.New("exact linux browser start error"),
			want:       process.Command{Name: "xdg-open", Args: []string{url}, Start: true},
		},
		{
			name:       "windows rundll32",
			goos:       "windows",
			processErr: errors.New("exact windows browser start error"),
			want: process.Command{
				Name: "rundll32", Args: []string{"url.dll,FileProtocolHandler", url}, Start: true,
			},
		},
		{
			name:       "darwin open",
			goos:       "darwin",
			processErr: errors.New("exact darwin browser start error"),
			want:       process.Command{Name: "open", Args: []string{url}, Start: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingBrowserLauncherProcess{err: test.processErr}
			dependencies := recording.dependencies(test.goos)
			if dependencies.Process.Runner != recording {
				t.Fatalf("browser-launcher dependency lost its complete process value: %#v", dependencies)
			}

			err := openBrowser(dependencies, url)

			if err != test.processErr {
				t.Fatalf("browser-launcher error was %v, want exact start error %v", err, test.processErr)
			}
			commands, populationErr := recording.assertedCommands()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(commands, []process.Command{test.want}) {
				t.Fatalf("recorded browser-launcher commands differ:\n got: %#v\nwant: %#v",
					commands, []process.Command{test.want})
			}
			command := commands[0]
			if command.Dir != "" || command.Stdin != nil || command.Stdout != nil || command.Stderr != nil {
				t.Fatalf("browser-launcher command changed directory or streams: %#v", command)
			}
		})
	}
}

func TestBrowserLauncherUnsupportedPlatformReturnsExactErrorWithoutProcessAttempt(t *testing.T) {
	recording := &recordingBrowserLauncherProcess{err: errors.New("must not be returned")}

	err := openBrowser(recording.dependencies("unsupported-complete-platform"), "arbitrary://url bytes")

	if err == nil || err.Error() != "unsupported platform" {
		t.Fatalf("unsupported platform error was %v, want exact text %q", err, "unsupported platform")
	}
	if len(recording.commands) != 0 {
		t.Fatalf("unsupported platform made process attempts: %#v", recording.commands)
	}
}

func TestBrowserLauncherDependenciesDefaultToSafeNoProcess(t *testing.T) {
	err := openBrowser(
		browserLauncherDependencies{GOOS: "linux"},
		"https://must-not-open.example.invalid/complete URL",
	)

	if err != nil {
		t.Fatalf("safe browser-launcher process default returned an error: %v", err)
	}
}

func TestRecordedBrowserLauncherCommandsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingBrowserLauncherProcess{}

	if _, err := recording.assertedCommands(); err == nil {
		t.Fatal("empty recorded browser-launcher process population passed")
	}
}
