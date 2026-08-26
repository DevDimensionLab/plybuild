package cmd

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type recordingPluginDiagramsExportProcess struct {
	commands []process.Command
	err      error
}

func (recording *recordingPluginDiagramsExportProcess) dependencies() pluginDiagramsExportDependencies {
	return pluginDiagramsExportDependencies{
		Process: process.Dependencies{Runner: recording},
	}
}

func (recording *recordingPluginDiagramsExportProcess) Run(command process.Command) error {
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	return recording.err
}

func (recording *recordingPluginDiagramsExportProcess) assertedCommands() ([]process.Command, error) {
	if len(recording.commands) == 0 {
		return nil, errors.New("recorded plugin-diagrams export process population is empty")
	}
	return recording.commands, nil
}

func TestPluginDiagramsExportSelectsCompleteSystemProcessDependencies(t *testing.T) {
	dependencies := systemPluginDiagramsExportDependencies()
	systemProcess := process.System()

	if dependencies.Process.Runner == nil {
		t.Fatal("plugin-diagrams export selected an incomplete process dependency")
	}
	if reflect.TypeOf(dependencies.Process.Runner) != reflect.TypeOf(systemProcess.Runner) {
		t.Fatalf("plugin-diagrams export process dependency is %T, want %T",
			dependencies.Process.Runner, systemProcess.Runner)
	}
}

func TestPluginDiagramsExportPreservesExactRequestAttemptIgnoredErrorAndContinuesToDiscovery(t *testing.T) {
	useTemporaryWorkingDirectory(t)
	workspace := `/arbitrary workspace path//with spaces/../backslash\diagram-ø.dsl`
	processError := errors.New("complete ignored plugin-diagrams export process dependency error")
	recording := &recordingPluginDiagramsExportProcess{err: processError}
	dependencies := recording.dependencies()
	if dependencies.Process.Runner != recording {
		t.Fatalf("plugin-diagrams export dependency lost its complete process value: %#v", dependencies)
	}

	err := runStructurizrDiagrams(dependencies, workspace)

	if err != nil {
		t.Fatalf("plugin-diagrams returned %v after exact ignored process error %v", err, processError)
	}
	assertPluginDiagramsExportBetweenDeletionAndDiscovery(t)
	commands, populationErr := recording.assertedCommands()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	want := []process.Command{{
		Name: "structurizr-cli",
		Args: []string{"export", "-w", workspace, "-format", "dot", "-output", ".structurizr/"},
	}}
	if !reflect.DeepEqual(commands, want) {
		t.Fatalf("recorded plugin-diagrams export commands differ:\n got: %#v\nwant: %#v", commands, want)
	}
	if len(commands) != 1 {
		t.Fatalf("plugin-diagrams export made %d process requests, want exactly 1", len(commands))
	}
	command := commands[0]
	if command.Dir != "" || command.Stdin != nil || command.Stdout != nil || command.Stderr != nil || command.Start {
		t.Fatalf("plugin-diagrams export changed working directory, streams, or synchronous mode: %#v", command)
	}
}

func TestPluginDiagramsExportDependenciesDefaultToSafeNoProcess(t *testing.T) {
	useTemporaryWorkingDirectory(t)

	err := runStructurizrDiagrams(
		pluginDiagramsExportDependencies{},
		"/developer/home/project/must-not-be-accessed/workspace.dsl",
	)

	if err != nil {
		t.Fatalf("safe plugin-diagrams export default returned an error: %v", err)
	}
}

func TestRecordedPluginDiagramsExportCommandsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingPluginDiagramsExportProcess{}

	if _, err := recording.assertedCommands(); err == nil {
		t.Fatal("empty recorded plugin-diagrams export process population passed")
	}
}

func useTemporaryWorkingDirectory(t *testing.T) {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("get original working directory: %v", err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("change to temporary working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Errorf("restore original working directory: %v", err)
		}
	})
}

func assertPluginDiagramsExportBetweenDeletionAndDiscovery(t *testing.T) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate plugin-diagrams export contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "plugin_diagrams.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse plugin-diagrams source: %v", err)
	}
	var function *ast.FuncDecl
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if ok && candidate.Recv == nil && candidate.Name.Name == "runStructurizrDiagrams" {
			function = candidate
			break
		}
	}
	if function == nil || function.Body == nil || len(function.Body.List) < 4 {
		t.Fatal("private plugin-diagrams flow is missing its deletion, export, or discovery statement")
	}
	wantCalls := []string{"file.DeleteAll", "process.Execute", "file.FindAll"}
	for index, want := range wantCalls {
		assignment, callName := assignedCall(function.Body.List[index+1])
		if callName != want {
			t.Fatalf("plugin-diagrams statement %d calls %q, want %q", index+1, callName, want)
		}
		if index == 1 {
			ignored, ok := assignment.Lhs[0].(*ast.Ident)
			if !ok || len(assignment.Lhs) != 1 || ignored.Name != "_" || assignment.Tok != token.ASSIGN {
				t.Fatal("plugin-diagrams export process error is no longer directly ignored")
			}
		}
	}
}

func assignedCall(statement ast.Stmt) (*ast.AssignStmt, string) {
	assignment, ok := statement.(*ast.AssignStmt)
	if !ok || len(assignment.Rhs) != 1 {
		return nil, ""
	}
	call, ok := assignment.Rhs[0].(*ast.CallExpr)
	if !ok {
		return assignment, ""
	}
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return assignment, ""
	}
	receiver, ok := selector.X.(*ast.Ident)
	if !ok {
		return assignment, ""
	}
	return assignment, receiver.Name + "." + selector.Sel.Name
}
