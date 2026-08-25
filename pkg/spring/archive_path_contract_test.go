package spring

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/file"
)

func (*recordingSpringArchivePathFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected Spring archive open")
}

func (recording *recordingSpringArchivePathFilesystem) Close(filesystem.File) error {
	return recording.unexpected("close")
}

func (recording *recordingSpringArchivePathFilesystem) CloseReader(io.ReadCloser) error {
	return recording.unexpected("close reader")
}

type recordingSpringArchivePathFilesystem struct {
	workingDirectory         string
	workingDirectoryErr      error
	workingDirectoryAttempts int
	unexpectedOperations     []string
}

var _ filesystem.FileSystem = (*recordingSpringArchivePathFilesystem)(nil)

func (recording *recordingSpringArchivePathFilesystem) WorkingDirectory() (string, error) {
	recording.workingDirectoryAttempts++
	return recording.workingDirectory, recording.workingDirectoryErr
}

func (recording *recordingSpringArchivePathFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingSpringArchivePathFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingSpringArchivePathFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingSpringArchivePathFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingSpringArchivePathFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingSpringArchivePathFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("recursive mkdir")
}

func (recording *recordingSpringArchivePathFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingSpringArchivePathFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingSpringArchivePathFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingSpringArchivePathFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingSpringArchivePathFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingSpringArchivePathFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingSpringArchivePathFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingSpringArchivePathFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingSpringArchivePathFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingSpringArchivePathFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected Spring archive-path " + operation)
}

func (recording *recordingSpringArchivePathFilesystem) dependencies() archivePathDependencies {
	return archivePathDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingSpringArchivePathFilesystem) assertedWorkingDirectoryAttempts() (int, error) {
	if recording.workingDirectoryAttempts == 0 {
		return 0, errors.New("recorded Spring archive-path working-directory population is empty")
	}
	return recording.workingDirectoryAttempts, nil
}

func TestArchivePathSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemArchivePathDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("archivePath selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("archivePath filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestArchivePathPreservesArbitraryWorkingDirectoryAndDirectUnixZipConstruction(t *testing.T) {
	workingDirectory := string([]byte{
		'/', 'a', 'r', 'b', 'i', 't', 'r', 'a', 'r', 'y', '/', '/', 'w', 'o', 'r', 'k', 'i', 'n', 'g', ' ',
		'd', 'i', 'r', '/', '.', '.', '/', 'b', 'a', 'c', 'k', 's', 'l', 'a', 's', 'h', '\\', 0xc3, 0xb8, '\n', 0x00,
	})
	recording := &recordingSpringArchivePathFilesystem{workingDirectory: workingDirectory}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("archive-path dependency lost its complete filesystem value: %#v", dependencies)
	}
	before := time.Now().Unix()

	archive, err := archivePathWithDependencies(dependencies)

	after := time.Now().Unix()
	if err != nil {
		t.Fatalf("archive-path composition returned an error: %v", err)
	}
	matchedTimestamp := false
	for timestamp := before; timestamp <= after; timestamp++ {
		if archive == file.Path("%s/spring-%d.zip", workingDirectory, timestamp) {
			matchedTimestamp = true
			break
		}
	}
	if !matchedTimestamp {
		t.Fatalf("archive path %q did not preserve directory %q and direct Unix timestamp in [%d, %d]",
			archive, workingDirectory, before, after)
	}
	assertOneSpringArchivePathWorkingDirectoryAttempt(t, recording)
	if len(recording.unexpectedOperations) != 0 {
		t.Fatalf("archive-path composition invoked unrelated filesystem operations: %#v",
			recording.unexpectedOperations)
	}
}

func TestArchivePathReturnsEmptyPathAndExactWorkingDirectoryErrorAfterOneAttempt(t *testing.T) {
	directoryError := errors.New("complete Spring archive-path working-directory error")
	recording := &recordingSpringArchivePathFilesystem{
		workingDirectory:    "/partial directory must be discarded",
		workingDirectoryErr: directoryError,
	}

	archive, err := archivePathWithDependencies(recording.dependencies())

	if archive != "" || err != directoryError {
		t.Fatalf("archive-path error result was (%q, %v), want (empty, exact %v)",
			archive, err, directoryError)
	}
	assertOneSpringArchivePathWorkingDirectoryAttempt(t, recording)
	if len(recording.unexpectedOperations) != 0 {
		t.Fatalf("archive-path error invoked unrelated filesystem operations: %#v",
			recording.unexpectedOperations)
	}
}

func TestArchivePathDependenciesDefaultToSafeEmptyResult(t *testing.T) {
	archive, err := archivePathWithDependencies(archivePathDependencies{})

	if archive != "" || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe archive-path dependency default returned (%q, %v), want (empty, %v)",
			archive, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedSpringArchivePathRejectsEmptyWorkingDirectoryPopulation(t *testing.T) {
	recording := &recordingSpringArchivePathFilesystem{}

	if _, err := recording.assertedWorkingDirectoryAttempts(); err == nil {
		t.Fatal("empty recorded Spring archive-path working-directory population passed")
	}
}

func assertOneSpringArchivePathWorkingDirectoryAttempt(t *testing.T, recording *recordingSpringArchivePathFilesystem) {
	t.Helper()
	attempts, err := recording.assertedWorkingDirectoryAttempts()
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 1 {
		t.Fatalf("archive-path composition made %d working-directory attempts, want 1", attempts)
	}
}
