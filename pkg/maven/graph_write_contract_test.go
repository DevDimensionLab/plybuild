package maven

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/file"
)

func (*recordingGraphStylesFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected graph-styles archive open")
}

func (recording *recordingGraphStylesFilesystem) Close(filesystem.File) error {
	return recording.unexpected("close")
}

func (recording *recordingGraphStylesFilesystem) CloseReader(io.Closer) error {
	return recording.unexpected("close reader")
}

type recordedGraphStylesWrite struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingGraphStylesFilesystem struct {
	writes               []recordedGraphStylesWrite
	writeErr             error
	unexpectedOperations []string
}

func (recording *recordingGraphStylesFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingGraphStylesFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingGraphStylesFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingGraphStylesFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingGraphStylesFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingGraphStylesFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingGraphStylesFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingGraphStylesFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedGraphStylesWrite{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (recording *recordingGraphStylesFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingGraphStylesFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingGraphStylesFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingGraphStylesFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingGraphStylesFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingGraphStylesFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingGraphStylesFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingGraphStylesFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingGraphStylesFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected Maven graph-styles " + operation)
}

func (recording *recordingGraphStylesFilesystem) dependencies() graphStylesWriteDependencies {
	return graphStylesWriteDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingGraphStylesFilesystem) assertedWrites() ([]recordedGraphStylesWrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded Maven graph-styles write population is empty")
	}
	return recording.writes, nil
}

func TestWriteGraphStylesSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemGraphStylesWriteDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("WriteGraphStyles selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("WriteGraphStyles filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestGraphStylesWritePreservesExactCallerComposedPathBytesModeAndDependencyError(t *testing.T) {
	projectPath := `/arbitrary project path//with spaces/../backslash\segment-ø`
	stylesFile := file.Path("%s/target/dependency-graph-styles.json", projectPath)
	writeError := errors.New("complete Maven graph-styles write dependency error")
	tests := []struct {
		name     string
		data     []byte
		writeErr error
	}{
		{name: "empty output", data: []byte{}},
		{name: "representative output", data: []byte("{\n    \"node-styles\": {}\n}")},
		{name: "non-ASCII output", data: []byte("{\n    \"node-styles\": {\"gruppe-ø\": {\"style\": \"填充\"}}\n}")},
		{
			name:     "arbitrary output and dependency error",
			data:     []byte{0x00, '{', '}', 0x0a, 0x80, 0xff},
			writeErr: writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingGraphStylesFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("Maven graph-styles dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := writeGraphStyles(dependencies, stylesFile, test.data)

			if err != test.writeErr {
				t.Fatalf("Maven graph-styles write error was %v, want exact dependency error %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantWrites := []recordedGraphStylesWrite{{
				path: stylesFile,
				data: test.data,
				mode: 0644,
			}}
			if !reflect.DeepEqual(writes, wantWrites) {
				t.Fatalf("recorded Maven graph-styles writes differ:\n got: %#v\nwant: %#v", writes, wantWrites)
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("Maven graph-styles invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
		})
	}
}

func TestGraphStylesWriteDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	stylesFile := "/developer/home/project/must-not-be-accessed/target/dependency-graph-styles.json"

	err := writeGraphStyles(graphStylesWriteDependencies{}, stylesFile, []byte("must not be written"))

	if err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe Maven graph-styles dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedGraphStylesWriteRejectsEmptyWritePopulation(t *testing.T) {
	recording := &recordingGraphStylesFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded Maven graph-styles write population passed")
	}
}
