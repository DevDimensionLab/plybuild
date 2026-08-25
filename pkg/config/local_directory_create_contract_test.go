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

func (*recordingLocalConfigDirectoryCreateFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected local-config directory-create archive open")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Close(filesystem.File) error {
	return recording.unexpected("close")
}

type recordedLocalConfigDirectoryCreate struct {
	path string
	mode fs.FileMode
}

type recordingLocalConfigDirectoryCreateFilesystem struct {
	mkdirs               []recordedLocalConfigDirectoryCreate
	mkdirErr             error
	unexpectedOperations []string
}

var _ filesystem.FileSystem = (*recordingLocalConfigDirectoryCreateFilesystem)(nil)

func (recording *recordingLocalConfigDirectoryCreateFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Mkdir(path string, mode fs.FileMode) error {
	recording.mkdirs = append(recording.mkdirs, recordedLocalConfigDirectoryCreate{path: path, mode: mode})
	return recording.mkdirErr
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("recursive mkdir")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected local-config-directory-create " + operation)
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) dependencies() localConfigDirectoryCreateDependencies {
	return localConfigDirectoryCreateDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingLocalConfigDirectoryCreateFilesystem) assertedMkdirs() ([]recordedLocalConfigDirectoryCreate, error) {
	if len(recording.mkdirs) == 0 {
		return nil, errors.New("recorded local-config-directory create population is empty")
	}
	return recording.mkdirs, nil
}

func TestLocalConfigDirCheckOrCreateConfigDirCreateSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemLocalConfigDirectoryCreateDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("LocalConfigDir.CheckOrCreateConfigDir create selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("LocalConfigDir.CheckOrCreateConfigDir create filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestLocalConfigDirectoryCreatePreservesExactCallerSelectedPathModeAttemptAndError(t *testing.T) {
	localCfgDir := OpenLocalConfig(`/arbitrary config directory//with spaces/../backslash\segment-ø`)
	dir := localCfgDir.Implementation().Path
	mkdirError := errors.New("complete local-config-directory create dependency error")
	tests := []struct {
		name string
		err  error
	}{
		{name: "success"},
		{name: "arbitrary error", err: mkdirError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingLocalConfigDirectoryCreateFilesystem{mkdirErr: test.err}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("local-config-directory-create dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := createLocalConfigDirectory(dependencies, dir)

			if err != test.err {
				t.Fatalf("local-config-directory create result was %v, want exact dependency result %v", err, test.err)
			}
			mkdirs, populationErr := recording.assertedMkdirs()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			want := []recordedLocalConfigDirectoryCreate{{path: dir, mode: 0755}}
			if !reflect.DeepEqual(mkdirs, want) {
				t.Fatalf("local-config-directory creates were %#v, want one exact create %#v", mkdirs, want)
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("local-config-directory create invoked unrelated filesystem operations: %#v",
					recording.unexpectedOperations)
			}
		})
	}
}

func TestLocalConfigDirectoryCreateDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	dir := "/developer/home/.ply/must-not-be-accessed"

	err := createLocalConfigDirectory(localConfigDirectoryCreateDependencies{}, dir)

	if err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe local-config-directory create dependency default returned %v, want %v",
			err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedLocalConfigDirectoryCreateRejectsEmptyMkdirPopulation(t *testing.T) {
	recording := &recordingLocalConfigDirectoryCreateFilesystem{}

	if _, err := recording.assertedMkdirs(); err == nil {
		t.Fatal("empty recorded local-config-directory create population passed")
	}
}
