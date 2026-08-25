package config

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
)

type recordingLocalConfigDirectoryStatFilesystem struct {
	statPaths            []string
	statInfo             fs.FileInfo
	statErr              error
	unexpectedOperations []string
}

var _ filesystem.FileSystem = (*recordingLocalConfigDirectoryStatFilesystem)(nil)

func (recording *recordingLocalConfigDirectoryStatFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.statPaths = append(recording.statPaths, path)
	return recording.statInfo, recording.statErr
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected local-config-directory-stat " + operation)
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) dependencies() localConfigDirectoryStatDependencies {
	return localConfigDirectoryStatDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingLocalConfigDirectoryStatFilesystem) assertedStatPaths() ([]string, error) {
	if len(recording.statPaths) == 0 {
		return nil, errors.New("recorded local-config-directory stat population is empty")
	}
	return recording.statPaths, nil
}

type localConfigDirectoryFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	system  interface{}
}

func (info *localConfigDirectoryFileInfo) Name() string       { return info.name }
func (info *localConfigDirectoryFileInfo) Size() int64        { return info.size }
func (info *localConfigDirectoryFileInfo) Mode() fs.FileMode  { return info.mode }
func (info *localConfigDirectoryFileInfo) ModTime() time.Time { return info.modTime }
func (info *localConfigDirectoryFileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info *localConfigDirectoryFileInfo) Sys() interface{}   { return info.system }

func TestLocalConfigDirCheckOrCreateConfigDirSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemLocalConfigDirectoryStatDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("LocalConfigDir.CheckOrCreateConfigDir selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("LocalConfigDir.CheckOrCreateConfigDir filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestLocalConfigDirectoryStatPreservesExactDirectoryFileInfoAndError(t *testing.T) {
	dir := `/arbitrary config directory//with spaces/../backslash\segment-ø`
	existingInfo := &localConfigDirectoryFileInfo{
		name: "existing-ø", size: 41, mode: fs.ModeDir | 0751,
		modTime: time.Unix(1700000301, 0), system: "complete existing metadata",
	}
	partialInfo := &localConfigDirectoryFileInfo{
		name: "partial-β", size: 82, mode: 0641,
		modTime: time.Unix(1700000302, 0), system: "complete partial metadata",
	}
	notExistError := &os.PathError{Op: "stat", Path: dir, Err: fs.ErrNotExist}
	otherError := errors.New("complete local-config-directory stat dependency error")
	tests := []struct {
		name string
		info fs.FileInfo
		err  error
	}{
		{name: "nil result"},
		{name: "existing directory", info: existingInfo},
		{name: "not-exist error", err: notExistError},
		{name: "partial metadata and arbitrary error", info: partialInfo, err: otherError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingLocalConfigDirectoryStatFilesystem{statInfo: test.info, statErr: test.err}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("local-config-directory-stat dependency lost its complete filesystem value: %#v", dependencies)
			}

			actualInfo, actualErr := statLocalConfigDirectory(dependencies, dir)

			if actualInfo != test.info || actualErr != test.err {
				t.Fatalf("local-config-directory stat result was (%#v, %v), want exact dependency result (%#v, %v)",
					actualInfo, actualErr, test.info, test.err)
			}
			paths, populationErr := recording.assertedStatPaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(paths, []string{dir}) {
				t.Fatalf("local-config-directory stat paths were %#v, want one exact path %#v", paths, []string{dir})
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("local-config-directory stat invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
		})
	}
}

func TestLocalConfigDirectoryStatDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	dir := "/developer/home/.ply/must-not-be-accessed"

	info, err := statLocalConfigDirectory(localConfigDirectoryStatDependencies{}, dir)

	if info != nil || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe local-config-directory stat dependency default returned (%#v, %v), want (nil, %v)",
			info, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedLocalConfigDirectoryStatRejectsEmptyStatPopulation(t *testing.T) {
	recording := &recordingLocalConfigDirectoryStatFilesystem{}

	if _, err := recording.assertedStatPaths(); err == nil {
		t.Fatal("empty recorded local-config-directory stat population passed")
	}
}
