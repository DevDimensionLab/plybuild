package process

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
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

func TestSystemSelectsExactStandardOutput(t *testing.T) {
	dependencies := System()

	if dependencies.Stdout != os.Stdout {
		t.Fatalf("system standard output is %T, want exact os.Stdout", dependencies.Stdout)
	}
	standardOutput := SystemStdout()
	if standardOutput == nil {
		t.Fatal("system standard-output composition is empty")
	}
	if standardOutput != dependencies.Stdout || standardOutput != os.Stdout {
		t.Fatalf("system standard-output composition is %T, want exact existing os.Stdout identity", standardOutput)
	}
	assertSystemStdoutComposition(t)
}

func assertSystemStdoutComposition(t *testing.T) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate process standard-output contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "process.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse process adapter source: %v", err)
	}
	var function *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if ok && candidate.Recv == nil && candidate.Name.Name == "SystemStdout" {
			function = candidate
			break
		}
	}
	if function == nil || function.Type.Params == nil || len(function.Type.Params.List) != 0 ||
		function.Type.Results == nil || len(function.Type.Results.List) != 1 ||
		function.Body == nil || len(function.Body.List) != 1 {
		t.Fatal("SystemStdout signature or single-return body changed")
	}
	result, resultOK := function.Type.Results.List[0].Type.(*ast.SelectorExpr)
	var resultPackage *ast.Ident
	resultPackageOK := false
	if resultOK {
		resultPackage, resultPackageOK = result.X.(*ast.Ident)
	}
	if !resultOK || !resultPackageOK || resultPackage.Name != "io" || result.Sel.Name != "Writer" {
		t.Fatal("SystemStdout no longer returns io.Writer")
	}
	returned, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 1 {
		t.Fatal("SystemStdout no longer directly returns one existing composition")
	}
	stdoutCall, ok := returned.Results[0].(*ast.CallExpr)
	if !ok || len(stdoutCall.Args) != 2 {
		t.Fatal("SystemStdout no longer returns Stdout with the complete dependency and enabled selection")
	}
	stdout, stdoutOK := stdoutCall.Fun.(*ast.Ident)
	systemCall, systemCallOK := stdoutCall.Args[0].(*ast.CallExpr)
	enabled, enabledOK := stdoutCall.Args[1].(*ast.Ident)
	if !stdoutOK || stdout.Name != "Stdout" || !systemCallOK || len(systemCall.Args) != 0 ||
		!enabledOK || enabled.Name != "true" {
		t.Fatal("SystemStdout no longer delegates to the exact Stdout(System(), true) composition")
	}
	system, systemOK := systemCall.Fun.(*ast.Ident)
	if !systemOK || system.Name != "System" {
		t.Fatal("SystemStdout no longer selects the existing complete system dependency")
	}
}

func TestSystemRunnerSelectsExactSystemRunnerWithoutStandardOutput(t *testing.T) {
	runnerOnly := SystemRunner()
	complete := System()

	if runnerOnly.Runner == nil {
		t.Fatal("runner-only system dependency has no process runner")
	}
	if runnerOnly.Runner != complete.Runner {
		t.Fatalf("runner-only system dependency is %T, want exact runner %T",
			runnerOnly.Runner, complete.Runner)
	}
	if runnerOnly.Stdout != nil {
		t.Fatalf("runner-only system standard output is %T, want nil", runnerOnly.Stdout)
	}
	if complete.Stdout != os.Stdout {
		t.Fatalf("complete system standard output is %T, want exact os.Stdout", complete.Stdout)
	}
}

func TestStdoutForwardsInjectedWriterOnlyWhenEnabled(t *testing.T) {
	injected := &bytes.Buffer{}
	dependencies := Dependencies{Stdout: injected}

	if got := Stdout(dependencies, true); got != injected {
		t.Fatalf("enabled standard output is %T, want exact injected writer", got)
	}
	if got := Stdout(dependencies, false); got != nil {
		t.Fatalf("disabled standard output is %T, want nil", got)
	}
}

func TestStdoutDependenciesDefaultToNil(t *testing.T) {
	if got := Stdout(Dependencies{}, true); got != nil {
		t.Fatalf("zero-value standard output is %T, want nil", got)
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
	if runner.command.Stdout != stdout || runner.command.Stderr != stderr {
		t.Fatal("stdout or stderr writer identity did not reach the process dependency unchanged")
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
