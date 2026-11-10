package webservice

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type recordingBrowserLauncherProcess struct {
	commands []process.Command
	err      error
}

func (recording *recordingBrowserLauncherProcess) dependencies(
	dependencies browserLauncherDependencies,
) browserLauncherDependencies {
	dependencies.Process.Runner = recording
	return dependencies
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

func TestOpenBrowserSelectsExactSystemRunnerWithoutProcessStandardOutputAndRuntimePlatform(t *testing.T) {
	dependencies := systemBrowserLauncherDependencies()
	systemProcess := process.System()

	if dependencies.Process.Runner == nil {
		t.Fatal("OpenBrowser selected no system process runner")
	}
	if dependencies.Process.Runner != systemProcess.Runner {
		t.Fatalf("OpenBrowser process dependency is %T, want %T",
			dependencies.Process.Runner, systemProcess.Runner)
	}
	if dependencies.Process.Stdout != nil {
		t.Fatalf("OpenBrowser inherited unused process standard output %T, want nil",
			dependencies.Process.Stdout)
	}
	if dependencies.GOOS != runtime.GOOS {
		t.Fatalf("OpenBrowser platform is %q, want runtime.GOOS %q", dependencies.GOOS, runtime.GOOS)
	}
	assertPublicOpenBrowserComposition(t)
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

	if len(tests) == 0 {
		t.Fatal("TestBrowserLauncherPreservesEveryPlatformRequestURLStartAttemptAndError test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingBrowserLauncherProcess{err: test.processErr}
			callerOwnedProcessStdout := &bytes.Buffer{}
			dependencies := recording.dependencies(browserLauncherDependencies{
				Process: process.Dependencies{Stdout: callerOwnedProcessStdout},
				GOOS:    test.goos,
			})
			if dependencies.Process.Runner != recording ||
				dependencies.Process.Stdout != callerOwnedProcessStdout || dependencies.GOOS != test.goos {
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
			if len(commands) != 1 {
				t.Fatalf("browser-launcher made %d process requests, want exactly 1", len(commands))
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

func assertPublicOpenBrowserComposition(t *testing.T) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate browser-launcher contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "api.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse browser-launcher source: %v", err)
	}
	var function *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if ok && candidate.Recv == nil && candidate.Name.Name == "OpenBrowser" {
			function = candidate
			break
		}
	}
	if function == nil || function.Type.Params == nil || function.Type.Results == nil ||
		len(function.Type.Params.List) != 1 || len(function.Type.Results.List) != 1 ||
		function.Body == nil || len(function.Body.List) != 1 {
		t.Fatal("public OpenBrowser signature or single-return body changed")
	}
	parameter := function.Type.Params.List[0]
	parameterType, parameterTypeOK := parameter.Type.(*ast.Ident)
	resultType, resultTypeOK := function.Type.Results.List[0].Type.(*ast.Ident)
	if len(parameter.Names) != 1 || parameter.Names[0].Name != "url" ||
		!parameterTypeOK || parameterType.Name != "string" || !resultTypeOK || resultType.Name != "error" {
		t.Fatal("public OpenBrowser signature is no longer func OpenBrowser(url string) error")
	}
	returned, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 1 {
		t.Fatal("public OpenBrowser no longer directly returns one private launcher call")
	}
	call, ok := returned.Results[0].(*ast.CallExpr)
	if !ok || len(call.Args) != 2 {
		t.Fatal("public OpenBrowser no longer directly returns openBrowser with two arguments")
	}
	callee, calleeOK := call.Fun.(*ast.Ident)
	if !calleeOK || callee.Name != "openBrowser" {
		t.Fatal("public OpenBrowser no longer directly returns the private openBrowser helper")
	}
	selector, selectorOK := call.Args[0].(*ast.CallExpr)
	if !selectorOK {
		t.Fatal("public OpenBrowser no longer selects its system dependencies at the call site")
	}
	selectorName, selectorNameOK := selector.Fun.(*ast.Ident)
	url, urlOK := call.Args[1].(*ast.Ident)
	if !selectorOK || !selectorNameOK || selectorName.Name != "systemBrowserLauncherDependencies" ||
		len(selector.Args) != 0 || !urlOK || url.Name != "url" {
		t.Fatal("public OpenBrowser no longer selects system dependencies beside the exact caller URL")
	}
}

func TestBrowserLauncherUnsupportedPlatformReturnsExactErrorWithoutProcessAttempt(t *testing.T) {
	recording := &recordingBrowserLauncherProcess{err: errors.New("must not be returned")}

	err := openBrowser(recording.dependencies(browserLauncherDependencies{
		GOOS: "unsupported-complete-platform",
	}), "arbitrary://url bytes")

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
