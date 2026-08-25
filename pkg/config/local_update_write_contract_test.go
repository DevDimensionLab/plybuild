package config

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
)

func (*recordingLocalConfigUpdateFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected local-config update-write archive open")
}

func (recording *recordingLocalConfigUpdateFilesystem) Close(filesystem.File) error {
	return recording.unexpected("close")
}

func (recording *recordingLocalConfigUpdateFilesystem) CloseReader(io.ReadCloser) error {
	return recording.unexpected("close reader")
}

type recordedLocalConfigUpdateWrite struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingLocalConfigUpdateFilesystem struct {
	writes               []recordedLocalConfigUpdateWrite
	writeErr             error
	unexpectedOperations []string
}

var _ filesystem.FileSystem = (*recordingLocalConfigUpdateFilesystem)(nil)

func (recording *recordingLocalConfigUpdateFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingLocalConfigUpdateFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingLocalConfigUpdateFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingLocalConfigUpdateFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingLocalConfigUpdateFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingLocalConfigUpdateFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigUpdateFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigUpdateFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedLocalConfigUpdateWrite{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (recording *recordingLocalConfigUpdateFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingLocalConfigUpdateFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingLocalConfigUpdateFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingLocalConfigUpdateFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingLocalConfigUpdateFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingLocalConfigUpdateFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingLocalConfigUpdateFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingLocalConfigUpdateFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingLocalConfigUpdateFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected local-config-update " + operation)
}

func (recording *recordingLocalConfigUpdateFilesystem) dependencies() localConfigUpdateWriteDependencies {
	return localConfigUpdateWriteDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingLocalConfigUpdateFilesystem) assertedWrites() ([]recordedLocalConfigUpdateWrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded local-config-update write population is empty")
	}
	return recording.writes, nil
}

func TestLocalConfigDirUpdateLocalConfigSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemLocalConfigUpdateWriteDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("LocalConfigDir.UpdateLocalConfig selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("LocalConfigDir.UpdateLocalConfig filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestLocalConfigUpdateWritePreservesExactCallerComposedPathBytesModeAndDependencyError(t *testing.T) {
	localCfgDir := OpenLocalConfig(`/arbitrary config directory//with spaces/../backslash\segment-ø`)
	configFilePath := localCfgDir.FilePath()
	writeError := errors.New("complete local-config-update write dependency error")
	tests := []struct {
		name     string
		data     []byte
		writeErr error
	}{
		{name: "empty output", data: []byte{}},
		{name: "representative YAML output", data: []byte("cloudConfig:\n  git:\n    url: https://example.invalid/config.git\n")},
		{name: "non-ASCII YAML output", data: []byte("sourceProvider:\n  host: vert-ø-β-日本語\n")},
		{
			name:     "arbitrary output and dependency error",
			data:     []byte{0x00, 'y', 'a', 'm', 'l', 0x0a, 0x80, 0xff},
			writeErr: writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingLocalConfigUpdateFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("local-config-update dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := writeLocalConfigUpdate(dependencies, configFilePath, test.data)

			if err != test.writeErr {
				t.Fatalf("local-config-update write error was %v, want exact dependency error %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantWrites := []recordedLocalConfigUpdateWrite{{
				path: configFilePath,
				data: test.data,
				mode: 0644,
			}}
			if !reflect.DeepEqual(writes, wantWrites) {
				t.Fatalf("recorded local-config-update writes differ:\n got: %#v\nwant: %#v", writes, wantWrites)
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("local-config-update invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
		})
	}
}

func TestLocalConfigUpdateWriteDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	configFilePath := "/developer/home/.ply/must-not-be-accessed/local-config.yaml"

	err := writeLocalConfigUpdate(localConfigUpdateWriteDependencies{}, configFilePath, []byte("must not be written"))

	if err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe local-config-update dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedLocalConfigUpdateWriteRejectsEmptyWritePopulation(t *testing.T) {
	recording := &recordingLocalConfigUpdateFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded local-config-update write population passed")
	}
}
