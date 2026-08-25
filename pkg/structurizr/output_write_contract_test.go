package structurizr

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
)

type recordedStructurizrOutputWrite struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingStructurizrOutputFilesystem struct {
	writes               []recordedStructurizrOutputWrite
	writeErr             error
	unexpectedOperations []string
}

func (recording *recordingStructurizrOutputFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingStructurizrOutputFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingStructurizrOutputFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingStructurizrOutputFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingStructurizrOutputFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedStructurizrOutputWrite{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (recording *recordingStructurizrOutputFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingStructurizrOutputFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingStructurizrOutputFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingStructurizrOutputFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingStructurizrOutputFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingStructurizrOutputFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingStructurizrOutputFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingStructurizrOutputFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingStructurizrOutputFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected structurizr-output " + operation)
}

func (recording *recordingStructurizrOutputFilesystem) dependencies() structurizrOutputWriteDependencies {
	return structurizrOutputWriteDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingStructurizrOutputFilesystem) assertedWrites() ([]recordedStructurizrOutputWrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded structurizr-output write population is empty")
	}
	return recording.writes, nil
}

func TestRunWithOutputToFileSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemStructurizrOutputWriteDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("RunWithOutputToFile selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("RunWithOutputToFile filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestStructurizrOutputWritePreservesExactPathBytesModeAndDependencyError(t *testing.T) {
	outputFile := `/arbitrary output path//with spaces/../backslash\diagram-ø.png`
	writeError := errors.New("complete structurizr-output write dependency error")
	tests := []struct {
		name     string
		data     []byte
		writeErr error
	}{
		{name: "empty output", data: []byte{}},
		{name: "representative output", data: []byte("digraph complete { a -> b }\n")},
		{name: "non-ASCII output", data: []byte("diagram ø β 日本語\n")},
		{
			name:     "arbitrary output and dependency error",
			data:     []byte{0x00, 'P', 'N', 'G', 0x0a, 0x80, 0xff},
			writeErr: writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingStructurizrOutputFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("structurizr-output dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := writeStructurizrOutput(dependencies, outputFile, test.data)

			if err != test.writeErr {
				t.Fatalf("structurizr-output write error was %v, want exact dependency error %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantWrites := []recordedStructurizrOutputWrite{{
				path: outputFile,
				data: test.data,
				mode: 0644,
			}}
			if !reflect.DeepEqual(writes, wantWrites) {
				t.Fatalf("recorded structurizr-output writes differ:\n got: %#v\nwant: %#v", writes, wantWrites)
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("structurizr-output invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
		})
	}
}

func TestStructurizrOutputWriteDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	outputFile := "/developer/home/project/must-not-be-accessed/diagram.png"

	err := writeStructurizrOutput(structurizrOutputWriteDependencies{}, outputFile, []byte("must not be written"))

	if err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe structurizr-output dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedStructurizrOutputWriteRejectsEmptyWritePopulation(t *testing.T) {
	recording := &recordingStructurizrOutputFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded structurizr-output write population passed")
	}
}
