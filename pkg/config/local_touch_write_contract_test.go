package config

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

type recordedLocalConfigTouchWrite struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingLocalConfigTouchFilesystem struct {
	writes               []recordedLocalConfigTouchWrite
	writeErr             error
	unexpectedOperations []string
}

var _ filesystem.FileSystem = (*recordingLocalConfigTouchFilesystem)(nil)

func (recording *recordingLocalConfigTouchFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingLocalConfigTouchFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingLocalConfigTouchFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingLocalConfigTouchFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigTouchFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigTouchFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedLocalConfigTouchWrite{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (recording *recordingLocalConfigTouchFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingLocalConfigTouchFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingLocalConfigTouchFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingLocalConfigTouchFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingLocalConfigTouchFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingLocalConfigTouchFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingLocalConfigTouchFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingLocalConfigTouchFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingLocalConfigTouchFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected local-config-touch " + operation)
}

func (recording *recordingLocalConfigTouchFilesystem) dependencies() localConfigTouchWriteDependencies {
	return localConfigTouchWriteDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingLocalConfigTouchFilesystem) assertedWrites() ([]recordedLocalConfigTouchWrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded local-config-touch write population is empty")
	}
	return recording.writes, nil
}

func TestLocalConfigDirTouchFileSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemLocalConfigTouchWriteDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("LocalConfigDir.TouchFile selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("LocalConfigDir.TouchFile filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestLocalConfigTouchWritePreservesExactCallerComposedPathBytesModeAndDependencyError(t *testing.T) {
	localCfgDir := OpenLocalConfig(`/arbitrary config directory//with spaces/../backslash\segment-ø`)
	configFilePath := localCfgDir.FilePath()
	writeError := errors.New("complete local-config-touch write dependency error")
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
			recording := &recordingLocalConfigTouchFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("local-config-touch dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := writeLocalConfigTouch(dependencies, configFilePath, test.data)

			if err != test.writeErr {
				t.Fatalf("local-config-touch write error was %v, want exact dependency error %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantWrites := []recordedLocalConfigTouchWrite{{
				path: configFilePath,
				data: test.data,
				mode: 0644,
			}}
			if !reflect.DeepEqual(writes, wantWrites) {
				t.Fatalf("recorded local-config-touch writes differ:\n got: %#v\nwant: %#v", writes, wantWrites)
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("local-config-touch invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
		})
	}
}

func TestLocalConfigTouchWriteDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	configFilePath := "/developer/home/.ply/must-not-be-accessed/local-config.yaml"

	err := writeLocalConfigTouch(localConfigTouchWriteDependencies{}, configFilePath, []byte("must not be written"))

	if err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe local-config-touch dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedLocalConfigTouchWriteRejectsEmptyWritePopulation(t *testing.T) {
	recording := &recordingLocalConfigTouchFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded local-config-touch write population passed")
	}
}
