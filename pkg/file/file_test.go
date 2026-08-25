package file

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/sirupsen/logrus"
)

type recordingFileExistenceFilesystem struct {
	statPaths []string
	statInfo  fs.FileInfo
	statErr   error
}

type recordingFileReadFilesystem struct {
	readPaths []string
	readData  []byte
	readErr   error
}

type recordedFileOverwrite struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingFileOverwriteFilesystem struct {
	writes   []recordedFileOverwrite
	writeErr error
}

type recordedFileCreate struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingFileCreateFilesystem struct {
	writes   []recordedFileCreate
	writeErr error
}

type recordedDirectoryCreate struct {
	path string
	mode fs.FileMode
}

type recordingDirectoryCreateFilesystem struct {
	statPaths []string
	creates   []recordedDirectoryCreate
	statInfo  fs.FileInfo
	statErr   error
	mkdirErr  error
}

type recordedFileOpenOperation struct {
	name  string
	path  string
	data  []byte
	flags int
	mode  fs.FileMode
}

type recordingFileOpenFilesystem struct {
	operations []recordedFileOpenOperation
	statInfo   fs.FileInfo
	statErr    error
	writeErr   error
	opened     *os.File
	openErr    error
}

type recordingFileDeleteFilesystem struct {
	removePaths []string
	removeErr   error
}

type recordingFileDeleteAllFilesystem struct {
	removeAllPaths []string
	removeAllErr   error
}

type recordedFileMove struct {
	source      string
	destination string
}

type recordingFileMoveFilesystem struct {
	moves     []recordedFileMove
	renameErr error
}

type recordingClearDirFilesystem struct {
	globPatterns   []string
	globMatches    []string
	globErr        error
	removeAllPaths []string
	removeAllErrs  map[string]error
}

type messageOnlyFormatter struct{}

func (messageOnlyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

type existingFileInfo struct {
	name string
}

func (info existingFileInfo) Name() string  { return info.name }
func (existingFileInfo) Size() int64        { return 0 }
func (existingFileInfo) Mode() fs.FileMode  { return 0644 }
func (existingFileInfo) ModTime() time.Time { return time.Time{} }
func (existingFileInfo) IsDir() bool        { return false }
func (existingFileInfo) Sys() interface{}   { return nil }

func (*recordingFileExistenceFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected file-existence read")
}

func (recording *recordingFileReadFilesystem) ReadFile(path string) ([]byte, error) {
	recording.readPaths = append(recording.readPaths, path)
	return append([]byte{}, recording.readData...), recording.readErr
}

func (*recordingFileReadFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected file-read stat")
}

func (*recordingFileReadFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected file-read mkdir")
}

func (*recordingFileReadFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected file-read write")
}

func (*recordingFileReadFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected file-read open file")
}

func (*recordingFileReadFilesystem) Remove(string) error {
	return errors.New("unexpected file-read remove")
}

func (*recordingFileReadFilesystem) RemoveAll(string) error {
	return errors.New("unexpected file-read recursive remove")
}

func (*recordingFileReadFilesystem) Rename(string, string) error {
	return errors.New("unexpected file-read rename")
}

func (*recordingFileReadFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected file-read glob")
}

func (*recordingFileReadFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected file-read create")
}

func (*recordingFileReadFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected file-read copy")
}

func (recording *recordingFileReadFilesystem) dependencies() openDependencies {
	return openDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileReadFilesystem) assertedReadPaths() ([]string, error) {
	if len(recording.readPaths) == 0 {
		return nil, errors.New("recorded file-read population is empty")
	}
	return recording.readPaths, nil
}

func (recording *recordingFileExistenceFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.statPaths = append(recording.statPaths, path)
	return recording.statInfo, recording.statErr
}

func (*recordingFileExistenceFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected file-existence mkdir")
}

func (*recordingFileExistenceFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected file-existence write")
}

func (*recordingFileExistenceFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected file-existence open file")
}

func (*recordingFileExistenceFilesystem) Remove(string) error {
	return errors.New("unexpected file-existence remove")
}

func (*recordingFileExistenceFilesystem) RemoveAll(string) error {
	return errors.New("unexpected file-existence recursive remove")
}

func (*recordingFileExistenceFilesystem) Rename(string, string) error {
	return errors.New("unexpected file-existence rename")
}

func (*recordingFileExistenceFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected file-existence glob")
}

func (*recordingFileExistenceFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected file-existence create")
}

func (*recordingFileExistenceFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected file-existence copy")
}

func (recording *recordingFileExistenceFilesystem) dependencies() existsDependencies {
	return existsDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileExistenceFilesystem) assertedStatPaths() ([]string, error) {
	if len(recording.statPaths) == 0 {
		return nil, errors.New("recorded file-existence stat population is empty")
	}
	return recording.statPaths, nil
}

func (*recordingFileOverwriteFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected file-overwrite read")
}

func (*recordingFileOverwriteFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected file-overwrite stat")
}

func (*recordingFileOverwriteFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected file-overwrite mkdir")
}

func (recording *recordingFileOverwriteFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedFileOverwrite{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (*recordingFileOverwriteFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected file-overwrite open file")
}

func (*recordingFileOverwriteFilesystem) Remove(string) error {
	return errors.New("unexpected file-overwrite remove")
}

func (*recordingFileOverwriteFilesystem) RemoveAll(string) error {
	return errors.New("unexpected file-overwrite recursive remove")
}

func (*recordingFileOverwriteFilesystem) Rename(string, string) error {
	return errors.New("unexpected file-overwrite rename")
}

func (*recordingFileOverwriteFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected file-overwrite glob")
}

func (*recordingFileOverwriteFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected file-overwrite create")
}

func (*recordingFileOverwriteFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected file-overwrite copy")
}

func (recording *recordingFileOverwriteFilesystem) dependencies() overwriteDependencies {
	return overwriteDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileOverwriteFilesystem) assertedWrites() ([]recordedFileOverwrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded file-overwrite population is empty")
	}
	return recording.writes, nil
}

func (*recordingFileCreateFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected file-create read")
}

func (*recordingFileCreateFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected file-create stat")
}

func (*recordingFileCreateFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected file-create mkdir")
}

func (recording *recordingFileCreateFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedFileCreate{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (*recordingFileCreateFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected file-create open file")
}

func (*recordingFileCreateFilesystem) Remove(string) error {
	return errors.New("unexpected file-create remove")
}

func (*recordingFileCreateFilesystem) RemoveAll(string) error {
	return errors.New("unexpected file-create recursive remove")
}

func (*recordingFileCreateFilesystem) Rename(string, string) error {
	return errors.New("unexpected file-create rename")
}

func (*recordingFileCreateFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected file-create glob")
}

func (*recordingFileCreateFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected file-create create")
}

func (*recordingFileCreateFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected file-create copy")
}

func (recording *recordingFileCreateFilesystem) dependencies() createFileDependencies {
	return createFileDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileCreateFilesystem) assertedWrites() ([]recordedFileCreate, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded file-create population is empty")
	}
	return recording.writes, nil
}

func (*recordingDirectoryCreateFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected directory-create read")
}

func (recording *recordingDirectoryCreateFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.statPaths = append(recording.statPaths, path)
	return recording.statInfo, recording.statErr
}

func (recording *recordingDirectoryCreateFilesystem) MkdirAll(path string, mode fs.FileMode) error {
	recording.creates = append(recording.creates, recordedDirectoryCreate{path: path, mode: mode})
	return recording.mkdirErr
}

func (*recordingDirectoryCreateFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected directory-create write")
}

func (*recordingDirectoryCreateFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected directory-create open file")
}

func (*recordingDirectoryCreateFilesystem) Remove(string) error {
	return errors.New("unexpected directory-create remove")
}

func (*recordingDirectoryCreateFilesystem) RemoveAll(string) error {
	return errors.New("unexpected directory-create recursive remove")
}

func (*recordingDirectoryCreateFilesystem) Rename(string, string) error {
	return errors.New("unexpected directory-create rename")
}

func (*recordingDirectoryCreateFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected directory-create glob")
}

func (*recordingDirectoryCreateFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected directory-create create")
}

func (*recordingDirectoryCreateFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected directory-create copy")
}

func (recording *recordingDirectoryCreateFilesystem) dependencies() createDirectoryDependencies {
	return createDirectoryDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingDirectoryCreateFilesystem) assertedStatPaths() ([]string, error) {
	if len(recording.statPaths) == 0 {
		return nil, errors.New("recorded directory-create stat population is empty")
	}
	return recording.statPaths, nil
}

func (recording *recordingDirectoryCreateFilesystem) assertedCreates() ([]recordedDirectoryCreate, error) {
	if len(recording.creates) == 0 {
		return nil, errors.New("recorded directory-create creation population is empty")
	}
	return recording.creates, nil
}

func (*recordingFileOpenFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected append-open read")
}

func (recording *recordingFileOpenFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.operations = append(recording.operations, recordedFileOpenOperation{name: "stat", path: path})
	return recording.statInfo, recording.statErr
}

func (*recordingFileOpenFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected append-open mkdir")
}

func (recording *recordingFileOpenFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.operations = append(recording.operations, recordedFileOpenOperation{
		name: "write", path: path, data: append([]byte{}, data...), mode: mode,
	})
	return recording.writeErr
}

func (*recordingFileOpenFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected append-open create")
}

func (*recordingFileOpenFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected append-open copy")
}

func (recording *recordingFileOpenFilesystem) OpenFile(path string, flags int, mode fs.FileMode) (*os.File, error) {
	recording.operations = append(recording.operations, recordedFileOpenOperation{
		name: "open-file", path: path, flags: flags, mode: mode,
	})
	return recording.opened, recording.openErr
}

func (*recordingFileOpenFilesystem) Remove(string) error {
	return errors.New("unexpected append-open remove")
}

func (*recordingFileOpenFilesystem) RemoveAll(string) error {
	return errors.New("unexpected append-open recursive remove")
}

func (*recordingFileOpenFilesystem) Rename(string, string) error {
	return errors.New("unexpected append-open rename")
}

func (*recordingFileOpenFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected append-open glob")
}

func (recording *recordingFileOpenFilesystem) dependencies() openFileDependencies {
	return openFileDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileOpenFilesystem) assertedOperations() ([]recordedFileOpenOperation, error) {
	if len(recording.operations) == 0 {
		return nil, errors.New("recorded append-open population is empty")
	}
	return recording.operations, nil
}

func (*recordingFileDeleteFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected single-file-delete read")
}

func (*recordingFileDeleteFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected single-file-delete stat")
}

func (*recordingFileDeleteFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected single-file-delete mkdir")
}

func (*recordingFileDeleteFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected single-file-delete write")
}

func (*recordingFileDeleteFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected single-file-delete open file")
}

func (*recordingFileDeleteFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected single-file-delete create")
}

func (*recordingFileDeleteFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected single-file-delete copy")
}

func (recording *recordingFileDeleteFilesystem) Remove(path string) error {
	recording.removePaths = append(recording.removePaths, path)
	return recording.removeErr
}

func (*recordingFileDeleteFilesystem) RemoveAll(string) error {
	return errors.New("unexpected single-file-delete recursive remove")
}

func (*recordingFileDeleteFilesystem) Rename(string, string) error {
	return errors.New("unexpected single-file-delete rename")
}

func (*recordingFileDeleteFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected single-file-delete glob")
}

func (recording *recordingFileDeleteFilesystem) dependencies() deleteSingleFileDependencies {
	return deleteSingleFileDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileDeleteFilesystem) assertedRemovePaths() ([]string, error) {
	if len(recording.removePaths) == 0 {
		return nil, errors.New("recorded single-file-delete population is empty")
	}
	return recording.removePaths, nil
}

func (*recordingFileDeleteAllFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected recursive-delete read")
}

func (*recordingFileDeleteAllFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected recursive-delete stat")
}

func (*recordingFileDeleteAllFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected recursive-delete mkdir")
}

func (*recordingFileDeleteAllFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected recursive-delete write")
}

func (*recordingFileDeleteAllFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected recursive-delete open file")
}

func (*recordingFileDeleteAllFilesystem) Remove(string) error {
	return errors.New("unexpected recursive-delete single remove")
}

func (recording *recordingFileDeleteAllFilesystem) RemoveAll(path string) error {
	recording.removeAllPaths = append(recording.removeAllPaths, path)
	return recording.removeAllErr
}

func (*recordingFileDeleteAllFilesystem) Rename(string, string) error {
	return errors.New("unexpected recursive-delete rename")
}

func (*recordingFileDeleteAllFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected recursive-delete glob")
}

func (*recordingFileDeleteAllFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected recursive-delete create")
}

func (*recordingFileDeleteAllFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected recursive-delete copy")
}

func (recording *recordingFileDeleteAllFilesystem) dependencies() deleteAllDependencies {
	return deleteAllDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileDeleteAllFilesystem) assertedRemoveAllPaths() ([]string, error) {
	if len(recording.removeAllPaths) == 0 {
		return nil, errors.New("recorded recursive-delete population is empty")
	}
	return recording.removeAllPaths, nil
}

func (*recordingFileMoveFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected file-move read")
}

func (*recordingFileMoveFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected file-move stat")
}

func (*recordingFileMoveFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected file-move mkdir")
}

func (*recordingFileMoveFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected file-move write")
}

func (*recordingFileMoveFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected file-move open file")
}

func (*recordingFileMoveFilesystem) Remove(string) error {
	return errors.New("unexpected file-move single remove")
}

func (*recordingFileMoveFilesystem) RemoveAll(string) error {
	return errors.New("unexpected file-move recursive remove")
}

func (recording *recordingFileMoveFilesystem) Rename(source, destination string) error {
	recording.moves = append(recording.moves, recordedFileMove{source: source, destination: destination})
	return recording.renameErr
}

func (*recordingFileMoveFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected file-move glob")
}

func (*recordingFileMoveFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected file-move create")
}

func (*recordingFileMoveFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected file-move copy")
}

func (recording *recordingFileMoveFilesystem) dependencies() moveDependencies {
	return moveDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileMoveFilesystem) assertedMoves() ([]recordedFileMove, error) {
	if len(recording.moves) == 0 {
		return nil, errors.New("recorded file-move population is empty")
	}
	return recording.moves, nil
}

func (*recordingClearDirFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected clear-directory read")
}

func (*recordingClearDirFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected clear-directory stat")
}

func (*recordingClearDirFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected clear-directory mkdir")
}

func (*recordingClearDirFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected clear-directory write")
}

func (*recordingClearDirFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected clear-directory open file")
}

func (*recordingClearDirFilesystem) Remove(string) error {
	return errors.New("unexpected clear-directory single remove")
}

func (recording *recordingClearDirFilesystem) RemoveAll(path string) error {
	recording.removeAllPaths = append(recording.removeAllPaths, path)
	return recording.removeAllErrs[path]
}

func (*recordingClearDirFilesystem) Rename(string, string) error {
	return errors.New("unexpected clear-directory rename")
}

func (recording *recordingClearDirFilesystem) Glob(pattern string) ([]string, error) {
	recording.globPatterns = append(recording.globPatterns, pattern)
	return append([]string{}, recording.globMatches...), recording.globErr
}

func (*recordingClearDirFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected clear-directory create")
}

func (*recordingClearDirFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected clear-directory copy")
}

func (recording *recordingClearDirFilesystem) dependencies() clearDirDependencies {
	return clearDirDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingClearDirFilesystem) assertedGlobPatterns() ([]string, error) {
	if len(recording.globPatterns) == 0 {
		return nil, errors.New("recorded clear-directory glob population is empty")
	}
	return recording.globPatterns, nil
}

func (recording *recordingClearDirFilesystem) assertedRemoveAllPaths() ([]string, error) {
	if len(recording.removeAllPaths) == 0 {
		return nil, errors.New("recorded clear-directory removal population is empty")
	}
	return recording.removeAllPaths, nil
}

func captureFileDebugLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	output := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(output)
	logger.SetFormatter(messageOnlyFormatter{})
	previousLogger := log
	log = logger
	t.Cleanup(func() { log = previousLogger })
	return output
}

func TestClearDirSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemClearDirDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("ClearDir selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("ClearDir filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestClearDirPreservesGlobResultRemovalExclusionAndLoggingOrder(t *testing.T) {
	directory := "/complete clear-directory/path with spaces"
	firstRemoval := filepath.Join(directory, "z-selected-first")
	excluded := filepath.Join(directory, "keep-one-and-keep-two")
	secondRemoval := filepath.Join(directory, "a-selected-second")
	recording := &recordingClearDirFilesystem{
		globMatches: []string{firstRemoval, excluded, secondRemoval},
	}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("clear-directory dependency lost its complete filesystem value: %#v", dependencies)
	}
	logOutput := captureFileDebugLog(t)

	err := clearDir(dependencies, directory, []string{"keep-one", "keep-two", "absent"})

	if err != nil {
		t.Fatalf("clear-directory returned %v, want nil", err)
	}
	patterns, populationErr := recording.assertedGlobPatterns()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantPatterns := []string{filepath.Join(directory, "*")}
	if !reflect.DeepEqual(patterns, wantPatterns) {
		t.Fatalf("clear-directory glob patterns were %#v, want %#v", patterns, wantPatterns)
	}
	removePaths, populationErr := recording.assertedRemoveAllPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantRemovePaths := []string{firstRemoval, secondRemoval}
	if !reflect.DeepEqual(removePaths, wantRemovePaths) {
		t.Fatalf("clear-directory removal paths were %#v, want ordered %#v", removePaths, wantRemovePaths)
	}
	wantLogs := []string{
		"Removing: " + firstRemoval,
		"Skipping removal of: " + excluded,
		"Skipping removal of: " + excluded,
		"Removing: " + secondRemoval,
	}
	gotLogs := strings.Split(strings.TrimSuffix(logOutput.String(), "\n"), "\n")
	if !reflect.DeepEqual(gotLogs, wantLogs) {
		t.Fatalf("clear-directory logs differ:\n got: %#v\nwant: %#v", gotLogs, wantLogs)
	}
}

func TestClearDirReturnsExactGlobErrorBeforeLoggingOrRemoval(t *testing.T) {
	directory := "/complete clear-directory/glob error"
	globError := errors.New("complete clear-directory glob dependency error")
	recording := &recordingClearDirFilesystem{
		globMatches: []string{filepath.Join(directory, "must-not-be-removed")},
		globErr:     globError,
	}
	logOutput := captureFileDebugLog(t)

	err := clearDir(recording.dependencies(), directory, []string{"exclude"})

	if err != globError {
		t.Fatalf("clear-directory glob error was %v, want exact dependency error %v", err, globError)
	}
	patterns, populationErr := recording.assertedGlobPatterns()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(patterns, []string{filepath.Join(directory, "*")}) {
		t.Fatalf("clear-directory glob patterns were %#v", patterns)
	}
	if len(recording.removeAllPaths) != 0 {
		t.Fatalf("clear-directory removed paths after a glob error: %#v", recording.removeAllPaths)
	}
	if logOutput.Len() != 0 {
		t.Fatalf("clear-directory logged after a glob error: %q", logOutput.String())
	}
}

func TestClearDirReturnsExactFirstRemovalErrorAndShortCircuits(t *testing.T) {
	directory := "/complete clear-directory/removal error"
	first := filepath.Join(directory, "first")
	failing := filepath.Join(directory, "failing")
	afterFailure := filepath.Join(directory, "must-not-be-removed")
	removeError := errors.New("complete clear-directory recursive-remove dependency error")
	recording := &recordingClearDirFilesystem{
		globMatches:   []string{first, failing, afterFailure},
		removeAllErrs: map[string]error{failing: removeError},
	}
	logOutput := captureFileDebugLog(t)

	err := clearDir(recording.dependencies(), directory, nil)

	if err != removeError {
		t.Fatalf("clear-directory removal error was %v, want exact dependency error %v", err, removeError)
	}
	wantPaths := []string{first, failing}
	if !reflect.DeepEqual(recording.removeAllPaths, wantPaths) {
		t.Fatalf("clear-directory removal paths were %#v, want first-error short circuit %#v", recording.removeAllPaths, wantPaths)
	}
	wantLogs := []string{"Removing: " + first, "Removing: " + failing}
	gotLogs := strings.Split(strings.TrimSuffix(logOutput.String(), "\n"), "\n")
	if !reflect.DeepEqual(gotLogs, wantLogs) {
		t.Fatalf("clear-directory error logs differ:\n got: %#v\nwant: %#v", gotLogs, wantLogs)
	}
}

func TestClearDirDependenciesDefaultToSafeNoPathAccess(t *testing.T) {
	directory := "/developer/home/project/must-not-be-accessed"

	err := clearDir(clearDirDependencies{}, directory, []string{"keep"})

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe clear-directory dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedClearDirRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingClearDirFilesystem{}

	if _, err := recording.assertedGlobPatterns(); err == nil {
		t.Fatal("empty recorded clear-directory glob population passed")
	}
	if _, err := recording.assertedRemoveAllPaths(); err == nil {
		t.Fatal("empty recorded clear-directory removal population passed")
	}
}

func TestMoveSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemMoveDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Move selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Move filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestMovePreservesCompletePathsAndExactDependencyError(t *testing.T) {
	source := "/complete file-move/source path with spaces/source.txt"
	destination := "/complete file-move/destination path with spaces/destination.txt"
	renameError := errors.New("complete file-move dependency error")
	tests := []struct {
		name      string
		renameErr error
	}{
		{name: "successful rename"},
		{name: "dependency error", renameErr: renameError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileMoveFilesystem{renameErr: test.renameErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("file-move dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := move(dependencies, source, destination)

			if err != test.renameErr {
				t.Fatalf("file-move error was %v, want exact dependency error %v", err, test.renameErr)
			}
			moves, populationErr := recording.assertedMoves()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			want := []recordedFileMove{{source: source, destination: destination}}
			if !reflect.DeepEqual(moves, want) {
				t.Fatalf("file-move dependency received values %#v, want %#v", moves, want)
			}
		})
	}
}

func TestMoveDependenciesDefaultToSafeNoPathAccess(t *testing.T) {
	source := "/developer/home/project/must-not-be-accessed/source.txt"
	destination := "/developer/home/project/must-not-be-accessed/destination.txt"

	err := move(moveDependencies{}, source, destination)

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe file-move dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedFileMoveRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileMoveFilesystem{}

	if _, err := recording.assertedMoves(); err == nil {
		t.Fatal("empty recorded file-move population passed")
	}
}

func TestDeleteSingleFileSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemDeleteSingleFileDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("DeleteSingleFile selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("DeleteSingleFile filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestDeleteSingleFilePreservesCompletePathAndExactDependencyError(t *testing.T) {
	path := "/complete single-file-delete/path with spaces/file.txt"
	removeError := errors.New("complete single-file-delete dependency error")
	tests := []struct {
		name      string
		removeErr error
	}{
		{name: "successful removal"},
		{name: "dependency error", removeErr: removeError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileDeleteFilesystem{removeErr: test.removeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("single-file-delete dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := deleteSingleFile(dependencies, path)

			if err != test.removeErr {
				t.Fatalf("single-file-delete error was %v, want exact dependency error %v", err, test.removeErr)
			}
			removePaths, populationErr := recording.assertedRemovePaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(removePaths, []string{path}) {
				t.Fatalf("single-file-delete dependency received paths %#v, want %#v", removePaths, []string{path})
			}
		})
	}
}

func TestDeleteSingleFileDependenciesDefaultToSafeNoPathAccess(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-be-removed.txt")

	err := deleteSingleFile(deleteSingleFileDependencies{}, target)

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe single-file-delete dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe single-file-delete dependency default accessed or mutated the target: %v", statErr)
	}
}

func TestRecordedFileDeleteRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileDeleteFilesystem{}

	if _, err := recording.assertedRemovePaths(); err == nil {
		t.Fatal("empty recorded single-file-delete population passed")
	}
}

func TestDeleteAllSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemDeleteAllDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("DeleteAll selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("DeleteAll filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestDeleteAllPreservesCompletePathAndExactDependencyError(t *testing.T) {
	path := "/complete recursive-delete/path with spaces"
	removeError := errors.New("complete recursive-delete dependency error")
	tests := []struct {
		name         string
		removeAllErr error
	}{
		{name: "successful removal"},
		{name: "dependency error", removeAllErr: removeError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileDeleteAllFilesystem{removeAllErr: test.removeAllErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("recursive-delete dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := deleteAll(dependencies, path)

			if err != test.removeAllErr {
				t.Fatalf("recursive-delete error was %v, want exact dependency error %v", err, test.removeAllErr)
			}
			removeAllPaths, populationErr := recording.assertedRemoveAllPaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(removeAllPaths, []string{path}) {
				t.Fatalf("recursive-delete dependency received paths %#v, want %#v", removeAllPaths, []string{path})
			}
		})
	}
}

func TestDeleteAllDependenciesDefaultToSafeNoPathAccess(t *testing.T) {
	path := "/developer/home/project/must-not-be-accessed"

	err := deleteAll(deleteAllDependencies{}, path)

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe recursive-delete dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedFileDeleteAllRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileDeleteAllFilesystem{}

	if _, err := recording.assertedRemoveAllPaths(); err == nil {
		t.Fatal("empty recorded recursive-delete population passed")
	}
}

func TestOpenFileSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemOpenFileDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("OpenFile selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("OpenFile filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestOpenFileMissingPathCreatesExactEmptyFileBeforeAppendOpen(t *testing.T) {
	path := "/complete append-open/path with spaces/file.txt"
	missingError := &os.PathError{Op: "stat", Path: path, Err: fs.ErrNotExist}
	openError := errors.New("complete append-open dependency error")
	opened := &os.File{}
	recording := &recordingFileOpenFilesystem{
		statErr: missingError,
		opened:  opened,
		openErr: openError,
	}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("append-open dependency lost its complete filesystem value: %#v", dependencies)
	}

	actual, err := openFile(dependencies, path)

	if actual != opened || err != openError {
		t.Fatalf("append-open result was (%p, %v), want exact dependency result (%p, %v)", actual, err, opened, openError)
	}
	want := []recordedFileOpenOperation{
		{name: "stat", path: path},
		{name: "write", path: path, data: []byte{}, mode: 0644},
		{name: "open-file", path: path, flags: os.O_APPEND | os.O_WRONLY, mode: 0644},
	}
	assertRecordedFileOpenOperations(t, recording, want)
}

func TestOpenFileExistingAndNonMissingStatErrorsSkipCreation(t *testing.T) {
	path := "/complete append-open/path with spaces/existing.txt"
	statError := errors.New("permission denied while probing complete append-open path")
	tests := []struct {
		name     string
		statInfo fs.FileInfo
		statErr  error
	}{
		{name: "existing path", statInfo: existingFileInfo{name: "existing.txt"}},
		{name: "non-missing stat error", statErr: statError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opened := &os.File{}
			recording := &recordingFileOpenFilesystem{
				statInfo: test.statInfo,
				statErr:  test.statErr,
				opened:   opened,
			}

			actual, err := openFile(recording.dependencies(), path)

			if actual != opened || err != nil {
				t.Fatalf("append-open result was (%p, %v), want (%p, nil)", actual, err, opened)
			}
			want := []recordedFileOpenOperation{
				{name: "stat", path: path},
				{name: "open-file", path: path, flags: os.O_APPEND | os.O_WRONLY, mode: 0644},
			}
			assertRecordedFileOpenOperations(t, recording, want)
		})
	}
}

func TestOpenFileCreateErrorShortCircuitsBeforeAppendOpen(t *testing.T) {
	path := "/complete append-open/create-error/file.txt"
	missingError := &os.PathError{Op: "stat", Path: path, Err: fs.ErrNotExist}
	createError := errors.New("complete empty-file creation dependency error")
	recording := &recordingFileOpenFilesystem{statErr: missingError, writeErr: createError}

	opened, err := openFile(recording.dependencies(), path)

	if opened != nil || err != createError {
		t.Fatalf("append-open create-error result was (%p, %v), want (nil, %v)", opened, err, createError)
	}
	want := []recordedFileOpenOperation{
		{name: "stat", path: path},
		{name: "write", path: path, data: []byte{}, mode: 0644},
	}
	assertRecordedFileOpenOperations(t, recording, want)
}

func TestOpenFileDependenciesDefaultToSafeNoPathAccess(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist.txt")

	opened, err := openFile(openFileDependencies{}, target)

	if opened != nil || !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe append-open dependency default returned (%p, %v), want (nil, %v)", opened, err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe append-open dependency default accessed or mutated the target: %v", statErr)
	}
}

func TestRecordedFileOpenRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileOpenFilesystem{}

	if _, err := recording.assertedOperations(); err == nil {
		t.Fatal("empty recorded append-open population passed")
	}
}

func assertRecordedFileOpenOperations(
	t *testing.T,
	recording *recordingFileOpenFilesystem,
	want []recordedFileOpenOperation,
) {
	t.Helper()
	operations, err := recording.assertedOperations()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(operations, want) {
		t.Fatalf("recorded append-open operations differ:\n got: %#v\nwant: %#v", operations, want)
	}
}

func TestOpenSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemOpenDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Open selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Open filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestOpenPreservesCompletePathBytesAndDependencyError(t *testing.T) {
	path := "/complete file-read/path with spaces/file.bin"
	readError := errors.New("complete file-read dependency error")
	tests := []struct {
		name     string
		readData []byte
		readErr  error
		wantData []byte
	}{
		{name: "empty bytes", readData: []byte{}, wantData: []byte{}},
		{
			name:     "multiline bytes",
			readData: []byte("first complete line\n\nthird complete line\n"),
			wantData: []byte("first complete line\n\nthird complete line\n"),
		},
		{name: "binary bytes", readData: []byte{0x00, 0xff, 0x7f, 0x0a}, wantData: []byte{0x00, 0xff, 0x7f, 0x0a}},
		{
			name:     "partial bytes with dependency error",
			readData: []byte("partial bytes that must be discarded"),
			readErr:  readError,
			wantData: []byte{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileReadFilesystem{readData: test.readData, readErr: test.readErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("file-read dependency lost its complete filesystem value: %#v", dependencies)
			}

			data, err := open(dependencies, path)
			if err != test.readErr {
				t.Fatalf("open error was %v, want exact dependency error %v", err, test.readErr)
			}
			if data == nil || !reflect.DeepEqual(data, test.wantData) {
				t.Fatalf("open bytes were %#v, want non-nil %#v", data, test.wantData)
			}
			readPaths, populationErr := recording.assertedReadPaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(readPaths, []string{path}) {
				t.Fatalf("file-read dependency received paths %#v, want %#v", readPaths, []string{path})
			}
		})
	}
}

func TestOpenDependenciesDefaultToSafeEmptyErrorWithoutPathAccess(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-be-read.txt")

	data, err := open(openDependencies{}, target)

	if data == nil || len(data) != 0 {
		t.Fatalf("safe file-read dependency default returned bytes %#v, want non-nil empty bytes", data)
	}
	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe file-read dependency default returned error %v, want %v", err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedFileReadRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileReadFilesystem{}

	if _, err := recording.assertedReadPaths(); err == nil {
		t.Fatal("empty recorded file-read population passed")
	}
}

func TestExistsSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemExistsDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Exists selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Exists filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestExistsPreservesCompletePathAndLegacyStatResults(t *testing.T) {
	path := "/complete file-existence/path with spaces/file.txt"
	otherError := errors.New("permission denied while probing complete path")
	tests := []struct {
		name     string
		info     fs.FileInfo
		err      error
		expected bool
	}{
		{
			name: "missing path is absent",
			err:  &os.PathError{Op: "stat", Path: path, Err: fs.ErrNotExist},
		},
		{
			name:     "existing path is present",
			info:     existingFileInfo{name: "file.txt"},
			expected: true,
		},
		{
			name:     "another stat error is present",
			err:      otherError,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileExistenceFilesystem{statInfo: test.info, statErr: test.err}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("file-existence dependency lost its complete filesystem value: %#v", dependencies)
			}

			if actual := exists(dependencies, path); actual != test.expected {
				t.Fatalf("exists result was %t, want %t for stat error %v", actual, test.expected, test.err)
			}
			statPaths, err := recording.assertedStatPaths()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(statPaths, []string{path}) {
				t.Fatalf("file-existence dependency received paths %#v, want %#v", statPaths, []string{path})
			}
		})
	}
}

func TestExistsDependenciesDefaultToSafeAbsent(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist.txt")

	if exists(existsDependencies{}, target) {
		t.Fatal("zero-value file-existence dependency reported an unprobed path as present")
	}
}

func TestRecordedFileExistenceRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileExistenceFilesystem{}

	if _, err := recording.assertedStatPaths(); err == nil {
		t.Fatal("empty recorded file-existence stat population passed")
	}
}

func TestOverwriteSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemOverwriteDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Overwrite selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Overwrite filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestOverwritePreservesCompleteLinesPathBytesModeAndDependencyError(t *testing.T) {
	path := "/complete file-overwrite/path with spaces/file.txt"
	writeError := errors.New("complete file-overwrite dependency error")
	tests := []struct {
		name      string
		lines     []string
		wantBytes []byte
		writeErr  error
	}{
		{name: "empty lines", lines: []string{}, wantBytes: []byte{}},
		{name: "single line", lines: []string{"single complete line"}, wantBytes: []byte("single complete line")},
		{
			name:      "multiple lines and dependency error",
			lines:     []string{"first complete line", "", "third complete line", ""},
			wantBytes: []byte("first complete line\n\nthird complete line\n"),
			writeErr:  writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileOverwriteFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("file-overwrite dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := overwrite(dependencies, test.lines, path)
			if !errors.Is(err, test.writeErr) {
				t.Fatalf("overwrite error was %v, want %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			want := []recordedFileOverwrite{{path: path, data: test.wantBytes, mode: 0644}}
			if !reflect.DeepEqual(writes, want) {
				t.Fatalf("recorded file-overwrite differs:\n got: %#v\nwant: %#v", writes, want)
			}
		})
	}
}

func TestOverwriteDependenciesDefaultToSafeNoMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist.txt")

	err := overwrite(overwriteDependencies{}, []string{"must not be written"}, target)

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe file-overwrite dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe file-overwrite dependency default mutated the target: %v", statErr)
	}
}

func TestRecordedFileOverwriteRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileOverwriteFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded file-overwrite population passed")
	}
}

func TestCreateFileSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemCreateFileDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("CreateFile selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("CreateFile filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestCreateFilePreservesCompletePathContentBytesModeAndDependencyError(t *testing.T) {
	path := "/complete file-create/path with spaces/file.txt"
	writeError := errors.New("complete file-create dependency error")
	tests := []struct {
		name      string
		content   string
		wantBytes []byte
		writeErr  error
	}{
		{name: "empty content", content: "", wantBytes: []byte{}},
		{
			name:      "multiline content and dependency error",
			content:   "first complete line\n\nthird complete line\n",
			wantBytes: []byte("first complete line\n\nthird complete line\n"),
			writeErr:  writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileCreateFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("file-create dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := createFile(dependencies, path, test.content)
			if !errors.Is(err, test.writeErr) {
				t.Fatalf("create-file error was %v, want %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			want := []recordedFileCreate{{path: path, data: test.wantBytes, mode: 0644}}
			if !reflect.DeepEqual(writes, want) {
				t.Fatalf("recorded file-create differs:\n got: %#v\nwant: %#v", writes, want)
			}
		})
	}
}

func TestCreateFileDependenciesDefaultToSafeNoMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist.txt")

	err := createFile(createFileDependencies{}, target, "must not be written")

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe file-create dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe file-create dependency default mutated the target: %v", statErr)
	}
}

func TestRecordedFileCreateRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileCreateFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded file-create population passed")
	}
}

func TestCreateDirectorySelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemCreateDirectoryDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("CreateDirectory selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("CreateDirectory filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestCreateDirectoryPreservesCompletePathMissingSelectionModeAndLegacyErrors(t *testing.T) {
	path := "/complete directory-create/path with spaces/nested"
	missingError := &os.PathError{Op: "stat", Path: path, Err: fs.ErrNotExist}
	mkdirError := errors.New("complete directory-create dependency error")
	tests := []struct {
		name      string
		mkdirErr  error
		wantError error
	}{
		{name: "missing path is created"},
		{
			name:      "creation failure returns original missing-path stat error",
			mkdirErr:  mkdirError,
			wantError: missingError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingDirectoryCreateFilesystem{statErr: missingError, mkdirErr: test.mkdirErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("directory-create dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := createDirectoryWithDependencies(dependencies, path)
			if err != test.wantError {
				t.Fatalf("directory-create error was %v, want original stat error %v", err, test.wantError)
			}
			statPaths, populationErr := recording.assertedStatPaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(statPaths, []string{path}) {
				t.Fatalf("directory-create dependency received stat paths %#v, want %#v", statPaths, []string{path})
			}
			creates, populationErr := recording.assertedCreates()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantCreates := []recordedDirectoryCreate{{path: path, mode: 0755}}
			if !reflect.DeepEqual(creates, wantCreates) {
				t.Fatalf("recorded directory creation differs:\n got: %#v\nwant: %#v", creates, wantCreates)
			}
		})
	}
}

func TestCreateDirectorySkipsCreationForExistingPathAndOtherStatErrors(t *testing.T) {
	path := "/complete directory-create/path with spaces/existing"
	otherError := errors.New("permission denied while probing complete directory path")
	tests := []struct {
		name string
		info fs.FileInfo
		err  error
	}{
		{name: "existing path", info: existingFileInfo{name: "existing"}},
		{name: "non-missing stat error", err: otherError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingDirectoryCreateFilesystem{statInfo: test.info, statErr: test.err}

			if err := createDirectoryWithDependencies(recording.dependencies(), path); err != nil {
				t.Fatalf("directory-create returned %v, want nil", err)
			}
			statPaths, populationErr := recording.assertedStatPaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(statPaths, []string{path}) {
				t.Fatalf("directory-create dependency received stat paths %#v, want %#v", statPaths, []string{path})
			}
			if len(recording.creates) != 0 {
				t.Fatalf("directory-create selected unexpected creations: %#v", recording.creates)
			}
		})
	}
}

func TestCreateDirectoryDependenciesDefaultToSafeNoMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist")

	if err := createDirectoryWithDependencies(createDirectoryDependencies{}, target); err != nil {
		t.Fatalf("safe directory-create dependency default returned %v, want nil", err)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe directory-create dependency default mutated the target: %v", statErr)
	}
}

func TestRecordedDirectoryCreateRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingDirectoryCreateFilesystem{}

	if _, err := recording.assertedStatPaths(); err == nil {
		t.Fatal("empty recorded directory-create stat population passed")
	}
	if _, err := recording.assertedCreates(); err == nil {
		t.Fatal("empty recorded directory-create creation population passed")
	}
}

func TestRelPath(t *testing.T) {
	relPath, err := RelPath(
		"/home/user/.ply/cloud-config/templates/flyway-demo",
		"/home/user/.ply/cloud-config/templates/flyway-demo/src/main/kotlin/no/ply/template/demo/flyway/Queue.kt")

	expected := "src/main/kotlin/no/ply/template/demo/flyway/Queue.kt"

	if err != nil {
		t.Errorf("%v\n", err)
	}

	if relPath != expected {
		t.Errorf("%s is not %s", relPath, expected)
	}
}
