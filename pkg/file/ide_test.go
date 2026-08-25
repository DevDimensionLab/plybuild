package file

import (
	"archive/zip"
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
)

func (*recordingIntellijFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected IntelliJ archive open")
}

type recordingIntellijFilesystem struct {
	readDirPaths   []string
	readDirEntries []fs.FileInfo
	readDirErr     error
}

func (*recordingIntellijFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected non-recursive IDE working directory")
}

func (*recordingIntellijFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected non-recursive IDE read")
}

func (recording *recordingIntellijFilesystem) ReadDir(path string) ([]fs.FileInfo, error) {
	recording.readDirPaths = append(recording.readDirPaths, path)
	return append([]fs.FileInfo{}, recording.readDirEntries...), recording.readDirErr
}

func (*recordingIntellijFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected non-recursive IDE read directory entries")
}

func (*recordingIntellijFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected non-recursive IDE stat")
}

func (*recordingIntellijFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected non-recursive IDE mkdir")
}

func (*recordingIntellijFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected non-recursive IDE mkdir")
}

func (*recordingIntellijFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected non-recursive IDE write")
}

func (*recordingIntellijFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected non-recursive IDE open file")
}

func (*recordingIntellijFilesystem) Remove(string) error {
	return errors.New("unexpected non-recursive IDE remove")
}

func (*recordingIntellijFilesystem) RemoveAll(string) error {
	return errors.New("unexpected non-recursive IDE recursive remove")
}

func (*recordingIntellijFilesystem) Rename(string, string) error {
	return errors.New("unexpected non-recursive IDE rename")
}

func (*recordingIntellijFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected non-recursive IDE glob")
}

func (*recordingIntellijFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected non-recursive IDE walk")
}

func (*recordingIntellijFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected non-recursive IDE create")
}

func (*recordingIntellijFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected non-recursive IDE copy")
}

func (recording *recordingIntellijFilesystem) dependencies() removeIntellijFileDependencies {
	return removeIntellijFileDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingIntellijFilesystem) assertedReadDirPaths() ([]string, error) {
	if len(recording.readDirPaths) == 0 {
		return nil, errors.New("recorded non-recursive IDE directory-read population is empty")
	}
	return recording.readDirPaths, nil
}

func (recording *recordingIntellijFilesystem) assertedReadDirEntries() ([]fs.FileInfo, error) {
	if len(recording.readDirEntries) == 0 {
		return nil, errors.New("recorded non-recursive IDE entry population is empty")
	}
	return recording.readDirEntries, nil
}

type intellijFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	system  interface{}
}

func (info intellijFileInfo) Name() string       { return info.name }
func (info intellijFileInfo) Size() int64        { return info.size }
func (info intellijFileInfo) Mode() fs.FileMode  { return info.mode }
func (info intellijFileInfo) ModTime() time.Time { return info.modTime }
func (info intellijFileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info intellijFileInfo) Sys() interface{}   { return info.system }

func TestRemoveIntellijFileSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemRemoveIntellijFileDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("non-recursive IDE cleanup selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("non-recursive IDE cleanup filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestRemoveIntellijFilePreservesExactTargetDeliveredOrderClassificationPathsLogsCountsAndReport(t *testing.T) {
	target := "/complete non-recursive IDE/path with spaces"
	entries := []fs.FileInfo{
		intellijFileInfo{name: "z-first.iml.backup", size: 41, mode: 0641, system: "first metadata"},
		intellijFileInfo{name: "ignored.iml-directory", size: 82, mode: fs.ModeDir | 0751, system: "ignored directory metadata"},
		intellijFileInfo{name: "ignored.idea-file", size: 83, mode: 0601, system: "ignored file metadata"},
		intellijFileInfo{name: "a-second.idea.backup", size: 84, mode: fs.ModeDir | 0705, system: "second metadata"},
		intellijFileInfo{name: "case-mismatch.IML", size: 85, mode: 0644},
		intellijFileInfo{name: "neutral.txt", size: 86, mode: 0644},
	}
	recording := &recordingIntellijFilesystem{readDirEntries: entries}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("non-recursive IDE dependency lost its complete filesystem value: %#v", dependencies)
	}
	logOutput := captureFileDebugLog(t)

	report, err := removeIntellijFile(dependencies, target, true)

	if err != nil || report != "Iml files: 1, .idea dirs: 1" {
		t.Fatalf("non-recursive IDE result was (%q, %v), want exact report and nil", report, err)
	}
	paths, populationErr := recording.assertedReadDirPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(paths, []string{target}) {
		t.Fatalf("non-recursive IDE directory-read paths were %#v, want %#v", paths, []string{target})
	}
	readEntries, populationErr := recording.assertedReadDirEntries()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(readEntries, entries) {
		t.Fatalf("non-recursive IDE entry metadata was %#v, want %#v", readEntries, entries)
	}
	wantLogs := []string{
		"Found .iml file: " + Path("%s/%s", target, entries[0].Name()),
		"Found .iml file: " + Path("%s/%s", target, entries[3].Name()),
	}
	gotLogs := strings.Split(strings.TrimSuffix(logOutput.String(), "\n"), "\n")
	if !reflect.DeepEqual(gotLogs, wantLogs) {
		t.Fatalf("non-recursive IDE dry-run logs differ:\n got: %#v\nwant: %#v", gotLogs, wantLogs)
	}
}

func TestRemoveIntellijFileReturnsExactReadErrorWithEmptyReportBeforeLogs(t *testing.T) {
	target := "/complete non-recursive IDE/read error"
	readError := errors.New("complete non-recursive IDE directory-read dependency error")
	recording := &recordingIntellijFilesystem{
		readDirEntries: []fs.FileInfo{intellijFileInfo{name: "must-not-be-visited.iml", mode: 0644}},
		readDirErr:     readError,
	}
	logOutput := captureFileDebugLog(t)

	report, err := removeIntellijFile(recording.dependencies(), target, true)

	if report != "" || err != readError {
		t.Fatalf("non-recursive IDE read-error result was (%q, %v), want empty report and exact error %v", report, err, readError)
	}
	paths, populationErr := recording.assertedReadDirPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(paths, []string{target}) {
		t.Fatalf("non-recursive IDE read-error paths were %#v, want %#v", paths, []string{target})
	}
	if logOutput.Len() != 0 {
		t.Fatalf("non-recursive IDE cleanup logged after a read error: %q", logOutput.String())
	}
}

func TestRemoveIntellijFileDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	report, err := removeIntellijFile(
		removeIntellijFileDependencies{},
		"/developer/home/project/must-not-be-accessed",
		true,
	)

	if report != "" || !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe non-recursive IDE dependency default returned (%q, %v), want (empty, %v)", report, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedRemoveIntellijFileRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingIntellijFilesystem{}

	if _, err := recording.assertedReadDirPaths(); err == nil {
		t.Fatal("empty recorded non-recursive IDE directory-read population passed")
	}
	if _, err := recording.assertedReadDirEntries(); err == nil {
		t.Fatal("empty recorded non-recursive IDE entry population passed")
	}
}

func TestRemoveIntellijFilesPreservesRecursiveAndNonRecursiveSelection(t *testing.T) {
	target := t.TempDir()

	nonRecursive, err := RemoveIntellijFiles(target, false, true)
	if err != nil || nonRecursive != "Iml files: 0, .idea dirs: 0" {
		t.Fatalf("non-recursive selection returned (%q, %v)", nonRecursive, err)
	}
	recursive, err := RemoveIntellijFiles(target, true, true)
	if err != nil || recursive != "Found 0 files and 0 directories to delete" {
		t.Fatalf("recursive selection returned (%q, %v)", recursive, err)
	}
}
