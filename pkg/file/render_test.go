package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
)

type recordedRenderOperation struct {
	name string
	path string
	data []byte
}

type recordedRenderFile struct {
	operations *[]recordedRenderOperation
	writeErr   error
	closed     bool
}

func (file *recordedRenderFile) Write(data []byte) (int, error) {
	*file.operations = append(*file.operations, recordedRenderOperation{
		name: "write",
		data: append([]byte{}, data...),
	})
	if file.writeErr != nil {
		return 0, file.writeErr
	}
	return len(data), nil
}

func (file *recordedRenderFile) Close() error {
	file.closed = true
	*file.operations = append(*file.operations, recordedRenderOperation{name: "close"})
	return nil
}

type recordingRenderFilesystem struct {
	operations []recordedRenderOperation
	readData   []byte
	readErr    error
	createErr  error
	writeErr   error
	file       *recordedRenderFile
}

func (*recordingRenderFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected render working directory")
}

func (recording *recordingRenderFilesystem) ReadFile(path string) ([]byte, error) {
	recording.operations = append(recording.operations, recordedRenderOperation{name: "read", path: path})
	return append([]byte{}, recording.readData...), recording.readErr
}

func (*recordingRenderFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected render read directory")
}

func (*recordingRenderFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected render read directory entries")
}

func (*recordingRenderFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected render stat")
}

func (*recordingRenderFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected render mkdir")
}

func (*recordingRenderFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected render mkdir")
}

func (*recordingRenderFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected render write file")
}

func (*recordingRenderFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected render open file")
}

func (*recordingRenderFilesystem) Remove(string) error {
	return errors.New("unexpected render remove")
}

func (*recordingRenderFilesystem) RemoveAll(string) error {
	return errors.New("unexpected render recursive remove")
}

func (*recordingRenderFilesystem) Rename(string, string) error {
	return errors.New("unexpected render rename")
}

func (*recordingRenderFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected render glob")
}

func (*recordingRenderFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected render walk")
}

func (recording *recordingRenderFilesystem) Create(path string) (filesystem.File, error) {
	recording.operations = append(recording.operations, recordedRenderOperation{name: "create", path: path})
	if recording.createErr != nil {
		return nil, recording.createErr
	}
	recording.file = &recordedRenderFile{operations: &recording.operations, writeErr: recording.writeErr}
	return recording.file, nil
}

func (*recordingRenderFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected render copy")
}

func (recording *recordingRenderFilesystem) dependencies() renderDependencies {
	return renderDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingRenderFilesystem) assertedOperations() ([]recordedRenderOperation, error) {
	if len(recording.operations) == 0 {
		return nil, errors.New("recorded render operation population is empty")
	}
	return recording.operations, nil
}

func (recording *recordingRenderFilesystem) assertedReadPaths() ([]string, error) {
	var paths []string
	for _, operation := range recording.operations {
		if operation.name == "read" {
			paths = append(paths, operation.path)
		}
	}
	if len(paths) == 0 {
		return nil, errors.New("recorded render read population is empty")
	}
	return paths, nil
}

func (recording *recordingRenderFilesystem) assertedCreatePaths() ([]string, error) {
	var paths []string
	for _, operation := range recording.operations {
		if operation.name == "create" {
			paths = append(paths, operation.path)
		}
	}
	if len(paths) == 0 {
		return nil, errors.New("recorded render create population is empty")
	}
	return paths, nil
}

func (recording *recordingRenderFilesystem) assertedWrites() ([][]byte, error) {
	var writes [][]byte
	for _, operation := range recording.operations {
		if operation.name == "write" {
			writes = append(writes, append([]byte{}, operation.data...))
		}
	}
	if len(writes) == 0 {
		return nil, errors.New("recorded render write population is empty")
	}
	return writes, nil
}

func TestRenderSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemRenderDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Render selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Render filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestRenderPreservesCompleteDependencyPathsOrderAndRenderedBytes(t *testing.T) {
	inputPath := "/complete render/input path with spaces/template.render"
	outputPath := "/complete render/output path with spaces/rendered.txt"
	recording := &recordingRenderFilesystem{
		readData: []byte("artifact={{.ArtifactID}}\ngroup={{.GroupID}}\n"),
	}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("render dependency lost its complete filesystem value: %#v", dependencies)
	}

	err := render(dependencies, inputPath, outputPath, struct {
		ArtifactID string
		GroupID    string
	}{ArtifactID: "complete-artifact", GroupID: "complete.group"})

	if err != nil {
		t.Fatalf("render returned %v, want nil", err)
	}
	want := []recordedRenderOperation{
		{name: "read", path: inputPath},
		{name: "create", path: outputPath},
		{name: "write", data: []byte("artifact=")},
		{name: "write", data: []byte("complete-artifact")},
		{name: "write", data: []byte("\ngroup=")},
		{name: "write", data: []byte("complete.group")},
		{name: "write", data: []byte("\n")},
	}
	assertRecordedRenderOperations(t, recording, want)
	if recording.file == nil || recording.file.closed {
		t.Fatalf("render file was %#v; want a created file left open", recording.file)
	}
	readPaths, populationErr := recording.assertedReadPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	createPaths, populationErr := recording.assertedCreatePaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	writes, populationErr := recording.assertedWrites()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(readPaths, []string{inputPath}) || !reflect.DeepEqual(createPaths, []string{outputPath}) ||
		string(joinRenderWrites(writes)) != "artifact=complete-artifact\ngroup=complete.group\n" {
		t.Fatalf("render populations lost complete values: reads=%#v creates=%#v writes=%q", readPaths, createPaths, joinRenderWrites(writes))
	}
}

func TestRenderInputErrorShortCircuitsBeforeCreateOrWrite(t *testing.T) {
	inputPath := "/complete render/input error/template.render"
	readError := errors.New("complete render input dependency error")
	recording := &recordingRenderFilesystem{
		readData: []byte("partial input bytes that must not be parsed"),
		readErr:  readError,
	}

	err := render(recording.dependencies(), inputPath, "/must-not-create/rendered.txt", struct{}{})

	if err != readError {
		t.Fatalf("render input error was %v, want exact dependency error %v", err, readError)
	}
	assertRecordedRenderOperations(t, recording, []recordedRenderOperation{{name: "read", path: inputPath}})
}

func TestRenderTemplateParsePanicPrecedesCreateOrWrite(t *testing.T) {
	inputPath := "/complete render/parse panic/template.render"
	recording := &recordingRenderFilesystem{readData: []byte("before {{")}
	var recovered interface{}

	func() {
		defer func() { recovered = recover() }()
		_ = render(recording.dependencies(), inputPath, "/must-not-create/rendered.txt", struct{}{})
	}()

	if recovered == nil {
		t.Fatal("invalid render template did not preserve template.Must panic behavior")
	}
	if !strings.Contains(fmt.Sprint(recovered), inputPath) {
		t.Fatalf("render parse panic %q lost the exact input template name %q", recovered, inputPath)
	}
	assertRecordedRenderOperations(t, recording, []recordedRenderOperation{{name: "read", path: inputPath}})
}

func TestRenderReturnsExactCreateErrorBeforeWrite(t *testing.T) {
	inputPath := "/complete render/create error/template.render"
	outputPath := "/complete render/create error/rendered.txt"
	createError := errors.New("complete render create dependency error")
	recording := &recordingRenderFilesystem{
		readData:  []byte("valid {{.Value}}"),
		createErr: createError,
	}

	err := render(recording.dependencies(), inputPath, outputPath, struct{ Value string }{Value: "data"})

	if err != createError {
		t.Fatalf("render create error was %v, want exact dependency error %v", err, createError)
	}
	assertRecordedRenderOperations(t, recording, []recordedRenderOperation{
		{name: "read", path: inputPath},
		{name: "create", path: outputPath},
	})
}

func TestRenderReturnsExactWriterExecutionErrorWithoutClosing(t *testing.T) {
	writeError := errors.New("complete render writer execution error")
	recording := &recordingRenderFilesystem{
		readData: []byte("complete rendered bytes"),
		writeErr: writeError,
	}

	err := render(recording.dependencies(), "/complete/render/input.render", "/complete/render/output.txt", struct{}{})

	if err != writeError {
		t.Fatalf("render execution error was %v, want exact writer error %v", err, writeError)
	}
	if recording.file == nil || recording.file.closed {
		t.Fatalf("render error file was %#v; want the created file left open", recording.file)
	}
	assertRecordedRenderOperations(t, recording, []recordedRenderOperation{
		{name: "read", path: "/complete/render/input.render"},
		{name: "create", path: "/complete/render/output.txt"},
		{name: "write", data: []byte("complete rendered bytes")},
	})
}

func TestRenderDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	err := render(
		renderDependencies{},
		"/developer/home/project/must-not-be-read.render",
		"/developer/home/project/must-not-be-created.txt",
		struct{}{},
	)

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe render dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedRenderRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingRenderFilesystem{}

	if _, err := recording.assertedOperations(); err == nil {
		t.Fatal("empty recorded render operation population passed")
	}
	if _, err := recording.assertedReadPaths(); err == nil {
		t.Fatal("empty recorded render read population passed")
	}
	if _, err := recording.assertedCreatePaths(); err == nil {
		t.Fatal("empty recorded render create population passed")
	}
	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded render write population passed")
	}
}

func assertRecordedRenderOperations(t *testing.T, recording *recordingRenderFilesystem, want []recordedRenderOperation) {
	t.Helper()
	operations, err := recording.assertedOperations()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(operations, want) {
		t.Fatalf("recorded render operations differ:\n got: %#v\nwant: %#v", operations, want)
	}
}

func joinRenderWrites(writes [][]byte) []byte {
	var joined []byte
	for _, write := range writes {
		joined = append(joined, write...)
	}
	return joined
}
