package cmd

import (
	"archive/zip"
	"bytes"
	"errors"
	"go/ast"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type recordedPluginDiagramsGraphvizWrite struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingPluginDiagramsGraphvizEffects struct {
	commands             []process.Command
	writes               []recordedPluginDiagramsGraphvizWrite
	stdout               []byte
	stderr               []byte
	processErr           error
	writeErr             error
	sequence             []string
	unexpectedOperations []string
}

func (recording *recordingPluginDiagramsGraphvizEffects) dependencies() pluginDiagramsGraphvizDependencies {
	return pluginDiagramsGraphvizDependencies{
		Process: process.Dependencies{Runner: recording},
		Files:   filesystem.Dependencies{FileSystem: recording},
	}
}

func (recording *recordingPluginDiagramsGraphvizEffects) Run(command process.Command) error {
	command.Args = append([]string(nil), command.Args...)
	recording.commands = append(recording.commands, command)
	recording.sequence = append(recording.sequence, "process")
	if command.Stdout != nil {
		_, _ = command.Stdout.Write(recording.stdout)
	}
	if command.Stderr != nil {
		_, _ = command.Stderr.Write(recording.stderr)
	}
	return recording.processErr
}

func (recording *recordingPluginDiagramsGraphvizEffects) assertedCommands() ([]process.Command, error) {
	if len(recording.commands) == 0 {
		return nil, errors.New("recorded plugin-diagrams Graphviz process population is empty")
	}
	return recording.commands, nil
}

func (recording *recordingPluginDiagramsGraphvizEffects) assertedWrites() ([]recordedPluginDiagramsGraphvizWrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded plugin-diagrams Graphviz write population is empty")
	}
	return recording.writes, nil
}

func (recording *recordingPluginDiagramsGraphvizEffects) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected plugin-diagrams Graphviz " + operation)
}

func (recording *recordingPluginDiagramsGraphvizEffects) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingPluginDiagramsGraphvizEffects) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingPluginDiagramsGraphvizEffects) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingPluginDiagramsGraphvizEffects) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingPluginDiagramsGraphvizEffects) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir all")
}

func (recording *recordingPluginDiagramsGraphvizEffects) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedPluginDiagramsGraphvizWrite{
		path: path,
		data: append([]byte(nil), data...),
		mode: mode,
	})
	recording.sequence = append(recording.sequence, "write")
	return recording.writeErr
}

func (recording *recordingPluginDiagramsGraphvizEffects) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingPluginDiagramsGraphvizEffects) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, recording.unexpected("archive open")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingPluginDiagramsGraphvizEffects) RemoveAll(string) error {
	return recording.unexpected("remove all")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingPluginDiagramsGraphvizEffects) Close(filesystem.File) error {
	return recording.unexpected("close")
}

func (recording *recordingPluginDiagramsGraphvizEffects) CloseReader(io.Closer) error {
	return recording.unexpected("close reader")
}

func TestPluginDiagramsGraphvizSelectsCompleteSystemProcessAndFilesystemDependencies(t *testing.T) {
	dependencies := systemPluginDiagramsGraphvizDependencies()
	systemProcess := process.System()
	systemFiles := filesystem.System()

	if dependencies.Process.Runner == nil || dependencies.Files.FileSystem == nil {
		t.Fatal("plugin-diagrams Graphviz selected incomplete process or filesystem dependencies")
	}
	if reflect.TypeOf(dependencies.Process.Runner) != reflect.TypeOf(systemProcess.Runner) {
		t.Fatalf("plugin-diagrams Graphviz process dependency is %T, want %T",
			dependencies.Process.Runner, systemProcess.Runner)
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("plugin-diagrams Graphviz filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestPluginDiagramsGraphvizPreservesExactRequestStreamsOutputWriteAndIgnoredWriteError(t *testing.T) {
	inputDotFile := `/arbitrary dot path//with spaces/../backslash\diagram-ø.dot`
	outputPngFile := `/arbitrary output path//with spaces/../backslash\diagram-β.png`
	stdoutBytes := []byte{0x00, 'P', 'N', 'G', '\n', 0x80, 0xff}
	stderrBytes := []byte("discarded Graphviz diagnostic \u00f8\n")
	writeError := errors.New("complete ignored plugin-diagrams Graphviz write dependency error")
	recording := &recordingPluginDiagramsGraphvizEffects{
		stdout:   stdoutBytes,
		stderr:   stderrBytes,
		writeErr: writeError,
	}
	dependencies := recording.dependencies()
	if dependencies.Process.Runner != recording || dependencies.Files.FileSystem != recording {
		t.Fatalf("plugin-diagrams Graphviz dependency lost its complete value: %#v", dependencies)
	}

	err := convertStructurizrDiagram(dependencies, inputDotFile, outputPngFile)

	if err != nil {
		t.Fatalf("plugin-diagrams Graphviz returned %v after exact ignored write error %v", err, writeError)
	}
	commands, populationErr := recording.assertedCommands()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(commands) != 1 {
		t.Fatalf("plugin-diagrams Graphviz made %d process requests, want exactly 1", len(commands))
	}
	command := commands[0]
	if command.Name != "dot" || !reflect.DeepEqual(command.Args, []string{inputDotFile, "-Tpng"}) {
		t.Fatalf("plugin-diagrams Graphviz process request changed executable or ordered arguments: %#v", command)
	}
	if command.Dir != "" || command.Stdin != nil || command.Start {
		t.Fatalf("plugin-diagrams Graphviz changed directory, stdin, or synchronous mode: %#v", command)
	}
	stdout, stdoutOK := command.Stdout.(*bytes.Buffer)
	stderr, stderrOK := command.Stderr.(*bytes.Buffer)
	if !stdoutOK || !stderrOK || stdout == stderr {
		t.Fatalf("plugin-diagrams Graphviz streams were stdout=%T stderr=%T with identities %p and %p",
			command.Stdout, command.Stderr, stdout, stderr)
	}
	if !bytes.Equal(stdout.Bytes(), stdoutBytes) || !bytes.Equal(stderr.Bytes(), stderrBytes) {
		t.Fatalf("plugin-diagrams Graphviz stream bytes changed: stdout=%v stderr=%v", stdout.Bytes(), stderr.Bytes())
	}
	writes, populationErr := recording.assertedWrites()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantWrites := []recordedPluginDiagramsGraphvizWrite{{
		path: outputPngFile,
		data: stdoutBytes,
		mode: 0644,
	}}
	if !reflect.DeepEqual(writes, wantWrites) {
		t.Fatalf("recorded plugin-diagrams Graphviz writes differ:\n got: %#v\nwant: %#v", writes, wantWrites)
	}
	if !reflect.DeepEqual(recording.sequence, []string{"process", "write"}) {
		t.Fatalf("plugin-diagrams Graphviz effect order was %#v, want process then write", recording.sequence)
	}
	if len(recording.unexpectedOperations) != 0 {
		t.Fatalf("plugin-diagrams Graphviz invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
	}
}

func TestPluginDiagramsGraphvizReturnsExactProcessErrorBeforeWriteOrOpen(t *testing.T) {
	processError := errors.New("complete plugin-diagrams Graphviz process dependency error")
	recording := &recordingPluginDiagramsGraphvizEffects{
		stdout:     []byte("partial stdout must not be written"),
		stderr:     []byte("exact discarded process-error stderr"),
		processErr: processError,
		writeErr:   errors.New("must not be observed"),
	}

	err := convertStructurizrDiagram(
		recording.dependencies(),
		"/complete process-error input.dot",
		"/developer/home/project/must-not-be-written-or-opened.png",
	)

	if err != processError {
		t.Fatalf("plugin-diagrams Graphviz process error was %v, want exact %v", err, processError)
	}
	commands, populationErr := recording.assertedCommands()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(commands) != 1 || commands[0].Name != "dot" {
		t.Fatalf("plugin-diagrams Graphviz process-error requests were %#v, want one dot request", commands)
	}
	if len(recording.writes) != 0 || !reflect.DeepEqual(recording.sequence, []string{"process"}) {
		t.Fatalf("plugin-diagrams Graphviz continued after process error: writes=%#v sequence=%#v",
			recording.writes, recording.sequence)
	}
	assertPluginDiagramsGraphvizPlacement(t)
}

func TestPluginDiagramsGraphvizDependenciesDefaultToSafeNoProcessOrWrite(t *testing.T) {
	outputPngFile := filepath.Join(t.TempDir(), "must-not-be-written.png")

	err := convertStructurizrDiagram(
		pluginDiagramsGraphvizDependencies{},
		"this-Graphviz-process-must-not-exist.dot",
		outputPngFile,
	)

	if err != nil {
		t.Fatalf("safe plugin-diagrams Graphviz default returned an error: %v", err)
	}
	if _, statErr := os.Stat(outputPngFile); !os.IsNotExist(statErr) {
		t.Fatalf("safe plugin-diagrams Graphviz default wrote the output path: %v", statErr)
	}
}

func TestRecordedPluginDiagramsGraphvizEffectsRejectEmptyPopulations(t *testing.T) {
	recording := &recordingPluginDiagramsGraphvizEffects{}

	if _, err := recording.assertedCommands(); err == nil {
		t.Fatal("empty recorded plugin-diagrams Graphviz process population passed")
	}
	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded plugin-diagrams Graphviz write population passed")
	}
}

func TestPluginDiagramsGraphvizRemainsOneConversionPerDiscoveredFileBeforeOpen(t *testing.T) {
	assertPluginDiagramsGraphvizPlacement(t)
}

func assertPluginDiagramsGraphvizPlacement(t *testing.T) {
	t.Helper()
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

	conversionIndex := -1
	conversionCalls := 0
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
		if !ok || identifier.Name != "convertStructurizrDiagram" {
			continue
		}
		conversionCalls++
		conversionIndex = index
		if len(call.Args) != 3 {
			t.Fatalf("plugin-diagrams Graphviz helper received %d arguments, want complete dependency, input, and output", len(call.Args))
		}
		dependency, dependencyOK := call.Args[0].(*ast.Ident)
		input, inputOK := call.Args[1].(*ast.Ident)
		output, outputOK := call.Args[2].(*ast.Ident)
		if !dependencyOK || dependency.Name != "graphvizDependencies" || !inputOK || input.Name != "file" ||
			!outputOK || output.Name != "outputPngFile" {
			t.Fatalf("plugin-diagrams Graphviz helper arguments are %#v", call.Args)
		}
	}
	if conversionCalls != 1 || conversionIndex < 0 {
		t.Fatalf("plugin-diagrams loop has %d Graphviz conversion calls, want exactly 1", conversionCalls)
	}
	if conversionIndex+2 >= len(loop.Body.List) {
		t.Fatal("plugin-diagrams iteration is missing the Graphviz error gate or following open helper")
	}
	assertPluginDiagramsDotErrorGate(t, loop.Body.List[conversionIndex+1])

	openStatement, ok := loop.Body.List[conversionIndex+2].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("plugin-diagrams post-Graphviz statement is %T, want direct open helper expression", loop.Body.List[conversionIndex+2])
	}
	openCall, ok := openStatement.X.(*ast.CallExpr)
	if !ok {
		t.Fatalf("plugin-diagrams post-Graphviz expression is %T, want direct open helper call", openStatement.X)
	}
	openName, nameOK := openCall.Fun.(*ast.Ident)
	if !nameOK || openName.Name != "openStructurizrDiagram" {
		t.Fatalf("plugin-diagrams post-Graphviz call is %#v, want unchanged openStructurizrDiagram", openStatement.X)
	}
	if conversionIndex+2 != len(loop.Body.List)-1 {
		t.Fatal("plugin-diagrams open helper no longer finishes each successful iteration")
	}
}
