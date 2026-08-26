package cmd

import (
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

type recordingPluginDiagramsOpenProcess struct {
	commands []process.Command
	err      error
}

func (recording *recordingPluginDiagramsOpenProcess) dependencies() pluginDiagramsOpenDependencies {
	return pluginDiagramsOpenDependencies{
		Process: process.Dependencies{Runner: recording},
	}
}

func (recording *recordingPluginDiagramsOpenProcess) Run(command process.Command) error {
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	return recording.err
}

func (recording *recordingPluginDiagramsOpenProcess) assertedCommands() ([]process.Command, error) {
	if len(recording.commands) == 0 {
		return nil, errors.New("recorded plugin-diagrams open process population is empty")
	}
	return recording.commands, nil
}

func TestPluginDiagramsOpenSelectsCompleteSystemProcessDependencies(t *testing.T) {
	dependencies := systemPluginDiagramsOpenDependencies()
	systemProcess := process.System()

	if dependencies.Process.Runner == nil {
		t.Fatal("plugin-diagrams open selected an incomplete process dependency")
	}
	if reflect.TypeOf(dependencies.Process.Runner) != reflect.TypeOf(systemProcess.Runner) {
		t.Fatalf("plugin-diagrams open process dependency is %T, want %T",
			dependencies.Process.Runner, systemProcess.Runner)
	}
}

func TestPluginDiagramsOpenPreservesExactRequestAttemptAndIgnoredError(t *testing.T) {
	outputPngFile := `/arbitrary output path//with spaces/../backslash\diagram-ø.png`
	processError := errors.New("complete ignored plugin-diagrams open process dependency error")
	recording := &recordingPluginDiagramsOpenProcess{err: processError}
	dependencies := recording.dependencies()
	if dependencies.Process.Runner != recording {
		t.Fatalf("plugin-diagrams open dependency lost its complete process value: %#v", dependencies)
	}

	openStructurizrDiagram(dependencies, outputPngFile)

	assertPluginDiagramsOpenErrorDirectlyIgnored(t)
	commands, populationErr := recording.assertedCommands()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	want := []process.Command{{
		Name: "open",
		Args: []string{outputPngFile},
	}}
	if !reflect.DeepEqual(commands, want) {
		t.Fatalf("recorded plugin-diagrams open commands differ:\n got: %#v\nwant: %#v", commands, want)
	}
	if len(commands) != 1 {
		t.Fatalf("plugin-diagrams open made %d process requests, want exactly 1", len(commands))
	}
	command := commands[0]
	if command.Dir != "" || command.Stdin != nil || command.Stdout != nil || command.Stderr != nil || command.Start {
		t.Fatalf("plugin-diagrams open changed working directory, streams, or synchronous mode: %#v", command)
	}
}

func TestPluginDiagramsOpenDependenciesDefaultToSafeNoProcess(t *testing.T) {
	openStructurizrDiagram(
		pluginDiagramsOpenDependencies{},
		"/developer/home/project/must-not-be-accessed/diagram.png",
	)
}

func TestRecordedPluginDiagramsOpenCommandsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingPluginDiagramsOpenProcess{}

	if _, err := recording.assertedCommands(); err == nil {
		t.Fatal("empty recorded plugin-diagrams open process population passed")
	}
}

func TestPluginDiagramsOpenRemainsAfterSuccessfulGraphvizAndContinuesIteration(t *testing.T) {
	function := parsedPluginDiagramsFunction(t, "runStructurizrDiagrams")
	var loop *ast.RangeStmt
	ast.Inspect(function.Body, func(node ast.Node) bool {
		candidate, ok := node.(*ast.RangeStmt)
		if ok && loop == nil {
			loop = candidate
		}
		return true
	})
	if loop == nil || loop.Body == nil {
		t.Fatal("private plugin-diagrams flow is missing its discovered-file iteration")
	}

	graphvizIndex := -1
	for index, statement := range loop.Body.List {
		assignment, ok := statement.(*ast.AssignStmt)
		if !ok || len(assignment.Rhs) != 1 {
			continue
		}
		call, ok := assignment.Rhs[0].(*ast.CallExpr)
		if !ok {
			continue
		}
		identifier, ok := call.Fun.(*ast.Ident)
		if ok && identifier.Name == "convertStructurizrDiagram" {
			graphvizIndex = index
			break
		}
	}
	if graphvizIndex < 0 || graphvizIndex+2 >= len(loop.Body.List) {
		t.Fatal("plugin-diagrams iteration is missing Graphviz conversion, error gate, or open helper")
	}
	assertPrivatePluginDiagramsGraphvizRequest(t, loop.Body.List[graphvizIndex])
	assertPluginDiagramsDotErrorGate(t, loop.Body.List[graphvizIndex+1])

	openCall, ok := loop.Body.List[graphvizIndex+2].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("plugin-diagrams open helper statement is %T, want direct expression", loop.Body.List[graphvizIndex+2])
	}
	call, ok := openCall.X.(*ast.CallExpr)
	if !ok {
		t.Fatalf("plugin-diagrams open helper expression is %T, want call", openCall.X)
	}
	identifier, ok := call.Fun.(*ast.Ident)
	if !ok || identifier.Name != "openStructurizrDiagram" {
		t.Fatalf("plugin-diagrams post-dot call is %#v, want openStructurizrDiagram", call.Fun)
	}
	if len(call.Args) != 2 {
		t.Fatalf("plugin-diagrams open helper received %d arguments, want complete dependency and output path", len(call.Args))
	}
	dependency, dependencyOK := call.Args[0].(*ast.Ident)
	output, outputOK := call.Args[1].(*ast.Ident)
	if !dependencyOK || dependency.Name != "openDependencies" || !outputOK || output.Name != "outputPngFile" {
		t.Fatalf("plugin-diagrams open helper arguments are %#v, want complete openDependencies and outputPngFile", call.Args)
	}
	if graphvizIndex+2 != len(loop.Body.List)-1 {
		t.Fatal("plugin-diagrams open helper no longer unconditionally continues to the next discovered file")
	}
}

func assertPluginDiagramsOpenErrorDirectlyIgnored(t *testing.T) {
	t.Helper()
	function := parsedPluginDiagramsFunction(t, "openStructurizrDiagram")
	if function.Body == nil || len(function.Body.List) != 1 {
		t.Fatal("private plugin-diagrams open helper must contain only its ignored process attempt")
	}
	assignment, callName := assignedCall(function.Body.List[0])
	if callName != "process.Execute" || assignment == nil || len(assignment.Lhs) != 1 {
		t.Fatalf("plugin-diagrams open helper calls %q, want one directly ignored process.Execute", callName)
	}
	ignored, ok := assignment.Lhs[0].(*ast.Ident)
	if !ok || ignored.Name != "_" || assignment.Tok != token.ASSIGN {
		t.Fatal("plugin-diagrams open process error is no longer directly ignored")
	}
}

func parsedPluginDiagramsFunction(t *testing.T, name string) *ast.FuncDecl {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate plugin-diagrams open contract")
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, filepath.Join(filepath.Dir(testFile), "plugin_diagrams.go"), nil, 0)
	if err != nil {
		t.Fatalf("parse plugin-diagrams source: %v", err)
	}
	for _, declaration := range parsed.Decls {
		candidate, ok := declaration.(*ast.FuncDecl)
		if ok && candidate.Recv == nil && candidate.Name.Name == name {
			return candidate
		}
	}
	t.Fatalf("private plugin-diagrams function %q is missing", name)
	return nil
}

func assertPrivatePluginDiagramsGraphvizRequest(t *testing.T, statement ast.Stmt) {
	t.Helper()
	assignment, ok := statement.(*ast.AssignStmt)
	if !ok || len(assignment.Rhs) != 1 {
		t.Fatalf("plugin-diagrams conversion is %T, want one assigned helper call", statement)
	}
	call, ok := assignment.Rhs[0].(*ast.CallExpr)
	if !ok {
		t.Fatalf("plugin-diagrams conversion expression is %T, want helper call", assignment.Rhs[0])
	}
	identifier, ok := call.Fun.(*ast.Ident)
	if !ok || identifier.Name != "convertStructurizrDiagram" {
		t.Fatalf("plugin-diagrams conversion helper is %#v, want convertStructurizrDiagram", call.Fun)
	}
	if len(call.Args) != 3 {
		t.Fatalf("plugin-diagrams conversion received %d arguments, want dependency, input, and output", len(call.Args))
	}
	dependency, dependencyOK := call.Args[0].(*ast.Ident)
	input, inputOK := call.Args[1].(*ast.Ident)
	output, outputOK := call.Args[2].(*ast.Ident)
	if !dependencyOK || dependency.Name != "graphvizDependencies" || !inputOK || input.Name != "file" ||
		!outputOK || output.Name != "outputPngFile" {
		t.Fatalf("plugin-diagrams conversion helper arguments are %#v", call.Args)
	}
}

func assertPluginDiagramsDotErrorGate(t *testing.T, statement ast.Stmt) {
	t.Helper()
	gate, ok := statement.(*ast.IfStmt)
	if !ok || gate.Body == nil || len(gate.Body.List) != 1 {
		t.Fatalf("plugin-diagrams dot error gate is %T, want one direct if return", statement)
	}
	condition, ok := gate.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("plugin-diagrams dot error condition is %T, want err != nil", gate.Cond)
	}
	conditionError, leftOK := condition.X.(*ast.Ident)
	conditionNil, rightOK := condition.Y.(*ast.Ident)
	if !leftOK || conditionError.Name != "err" || condition.Op != token.NEQ ||
		!rightOK || conditionNil.Name != "nil" {
		t.Fatalf("plugin-diagrams dot error condition is %#v, want err != nil", gate.Cond)
	}
	returned, ok := gate.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returned.Results) != 1 {
		t.Fatal("plugin-diagrams dot error gate no longer returns the direct conversion error")
	}
	identifier, ok := returned.Results[0].(*ast.Ident)
	if !ok || identifier.Name != "err" {
		t.Fatalf("plugin-diagrams dot error return is %#v, want err", returned.Results[0])
	}
}
