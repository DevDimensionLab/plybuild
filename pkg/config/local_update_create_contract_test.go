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

func (*recordingLocalConfigUpdateCreateFilesystem) OpenZipEntry(*zip.File) (io.ReadCloser, error) {
	return nil, errors.New("unexpected local-config update-create archive entry open")
}

func (*recordingLocalConfigUpdateCreateFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected local-config update-create archive open")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Close(filesystem.File) error {
	return recording.unexpected("close")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) CloseReader(io.Closer) error {
	return recording.unexpected("close reader")
}

type recordedLocalConfigUpdateCreateFile struct {
	closeCalls int
	closeErr   error
}

func (*recordedLocalConfigUpdateCreateFile) Write(data []byte) (int, error) {
	return len(data), nil
}

func (file *recordedLocalConfigUpdateCreateFile) Close() error {
	file.closeCalls++
	return file.closeErr
}

type recordingLocalConfigUpdateCreateFilesystem struct {
	createPaths          []string
	createFile           filesystem.File
	createErr            error
	unexpectedOperations []string
}

var _ filesystem.FileSystem = (*recordingLocalConfigUpdateCreateFilesystem)(nil)

func (recording *recordingLocalConfigUpdateCreateFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Create(path string) (filesystem.File, error) {
	recording.createPaths = append(recording.createPaths, path)
	return recording.createFile, recording.createErr
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected local-config-update-create " + operation)
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) dependencies() localConfigUpdateCreateDependencies {
	return localConfigUpdateCreateDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingLocalConfigUpdateCreateFilesystem) assertedCreatePaths() ([]string, error) {
	if len(recording.createPaths) == 0 {
		return nil, errors.New("recorded local-config-update create population is empty")
	}
	return recording.createPaths, nil
}

func TestLocalConfigDirUpdateLocalConfigCreateSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemLocalConfigUpdateCreateDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("LocalConfigDir.UpdateLocalConfig create selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("LocalConfigDir.UpdateLocalConfig create filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestLocalConfigUpdateCreatePreservesExactCallerComposedPathFileErrorAndCloseResult(t *testing.T) {
	localCfgDir := OpenLocalConfig(`/arbitrary config directory//with spaces/../backslash\segment-ø`)
	configFilePath := localCfgDir.FilePath()
	createError := errors.New("complete local-config-update create dependency error")
	closeError := errors.New("complete local-config-update close result")
	successFile := &recordedLocalConfigUpdateCreateFile{closeErr: closeError}
	combinedFile := &recordedLocalConfigUpdateCreateFile{}
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

	if len(tests) == 0 {
		t.Fatal("TestLocalConfigUpdateCreatePreservesExactCallerComposedPathFileErrorAndCloseResult test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingLocalConfigUpdateCreateFilesystem{
				createFile: test.createFile,
				createErr:  test.createErr,
			}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("local-config-update-create dependency lost its complete filesystem value: %#v", dependencies)
			}

			actualFile, actualErr := createLocalConfigUpdate(dependencies, configFilePath)

			if actualFile != test.createFile || actualErr != test.createErr {
				t.Fatalf("local-config-update create result was (%#v, %v), want exact dependency result (%#v, %v)",
					actualFile, actualErr, test.createFile, test.createErr)
			}
			paths, populationErr := recording.assertedCreatePaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(paths, []string{configFilePath}) {
				t.Fatalf("local-config-update create paths were %#v, want one exact path %#v", paths, []string{configFilePath})
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("local-config-update create invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
			if test.closeErr != nil {
				if err := actualFile.Close(); err != test.closeErr {
					t.Fatalf("local-config-update file Close returned %v, want exact result %v", err, test.closeErr)
				}
				if successFile.closeCalls != 1 {
					t.Fatalf("local-config-update file Close calls were %d, want 1", successFile.closeCalls)
				}
			}
		})
	}
}

func TestLocalConfigUpdateCreateDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	configFilePath := "/developer/home/.ply/must-not-be-accessed/local-config.yaml"

	file, err := createLocalConfigUpdate(localConfigUpdateCreateDependencies{}, configFilePath)

	if file != nil || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe local-config-update create dependency default returned (%#v, %v), want (nil, %v)",
			file, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedLocalConfigUpdateCreateRejectsEmptyCreatePopulation(t *testing.T) {
	recording := &recordingLocalConfigUpdateCreateFilesystem{}

	if _, err := recording.assertedCreatePaths(); err == nil {
		t.Fatal("empty recorded local-config-update create population passed")
	}
}
