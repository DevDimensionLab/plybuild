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

func (*recordingLocalConfigTouchCreateFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected local-config touch-create archive open")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Close(filesystem.File) error {
	return recording.unexpected("close")
}

type recordedLocalConfigTouchCreateFile struct {
	closeCalls int
	closeErr   error
}

func (*recordedLocalConfigTouchCreateFile) Write(data []byte) (int, error) {
	return len(data), nil
}

func (file *recordedLocalConfigTouchCreateFile) Close() error {
	file.closeCalls++
	return file.closeErr
}

type recordingLocalConfigTouchCreateFilesystem struct {
	createPaths          []string
	createFile           filesystem.File
	createErr            error
	unexpectedOperations []string
}

var _ filesystem.FileSystem = (*recordingLocalConfigTouchCreateFilesystem)(nil)

func (recording *recordingLocalConfigTouchCreateFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Create(path string) (filesystem.File, error) {
	recording.createPaths = append(recording.createPaths, path)
	return recording.createFile, recording.createErr
}

func (recording *recordingLocalConfigTouchCreateFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingLocalConfigTouchCreateFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected local-config-touch-create " + operation)
}

func (recording *recordingLocalConfigTouchCreateFilesystem) dependencies() localConfigTouchCreateDependencies {
	return localConfigTouchCreateDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingLocalConfigTouchCreateFilesystem) assertedCreatePaths() ([]string, error) {
	if len(recording.createPaths) == 0 {
		return nil, errors.New("recorded local-config-touch create population is empty")
	}
	return recording.createPaths, nil
}

func TestLocalConfigDirTouchFileCreateSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemLocalConfigTouchCreateDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("LocalConfigDir.TouchFile create selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("LocalConfigDir.TouchFile create filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestLocalConfigTouchCreatePreservesExactCallerComposedPathFileErrorAndCloseResult(t *testing.T) {
	localCfgDir := OpenLocalConfig(`/arbitrary config directory//with spaces/../backslash\segment-ø`)
	configFilePath := localCfgDir.FilePath()
	createError := errors.New("complete local-config-touch create dependency error")
	closeError := errors.New("complete local-config-touch close result")
	successFile := &recordedLocalConfigTouchCreateFile{closeErr: closeError}
	combinedFile := &recordedLocalConfigTouchCreateFile{}
	tests := []struct {
		name       string
		createFile filesystem.File
		createErr  error
		closeErr   error
	}{
		{name: "nil result"},
		{name: "successful file", createFile: successFile, closeErr: closeError},
		{name: "arbitrary error", createErr: createError},
		{name: "unusual file and error result", createFile: combinedFile, createErr: createError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingLocalConfigTouchCreateFilesystem{
				createFile: test.createFile,
				createErr:  test.createErr,
			}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("local-config-touch-create dependency lost its complete filesystem value: %#v", dependencies)
			}

			actualFile, actualErr := createLocalConfigTouch(dependencies, configFilePath)

			if actualFile != test.createFile || actualErr != test.createErr {
				t.Fatalf("local-config-touch create result was (%#v, %v), want exact dependency result (%#v, %v)",
					actualFile, actualErr, test.createFile, test.createErr)
			}
			paths, populationErr := recording.assertedCreatePaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(paths, []string{configFilePath}) {
				t.Fatalf("local-config-touch create paths were %#v, want one exact path %#v", paths, []string{configFilePath})
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("local-config-touch create invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
			if test.closeErr != nil {
				if err := actualFile.Close(); err != test.closeErr {
					t.Fatalf("local-config-touch file Close returned %v, want exact result %v", err, test.closeErr)
				}
				if successFile.closeCalls != 1 {
					t.Fatalf("local-config-touch file Close calls were %d, want 1", successFile.closeCalls)
				}
			}
		})
	}
}

func TestLocalConfigTouchCreateDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	configFilePath := "/developer/home/.ply/must-not-be-accessed/local-config.yaml"

	file, err := createLocalConfigTouch(localConfigTouchCreateDependencies{}, configFilePath)

	if file != nil || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe local-config-touch create dependency default returned (%#v, %v), want (nil, %v)",
			file, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedLocalConfigTouchCreateRejectsEmptyCreatePopulation(t *testing.T) {
	recording := &recordingLocalConfigTouchCreateFilesystem{}

	if _, err := recording.assertedCreatePaths(); err == nil {
		t.Fatal("empty recorded local-config-touch create population passed")
	}
}
