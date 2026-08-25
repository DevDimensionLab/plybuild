package tips

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
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
)

func (*recordingTipsListFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected tips-list archive open")
}

func (recording *recordingTipsListFilesystem) Close(filesystem.File) error {
	return recording.unexpected("close")
}

func (recording *recordingTipsListFilesystem) CloseReader(io.Closer) error {
	return recording.unexpected("close reader")
}

type recordingTipsListFilesystem struct {
	readPaths            []string
	readEntries          []fs.DirEntry
	deliveredEntries     []fs.DirEntry
	readErr              error
	unexpectedOperations []string
}

var _ filesystem.FileSystem = (*recordingTipsListFilesystem)(nil)

func (recording *recordingTipsListFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingTipsListFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingTipsListFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory metadata")
}

func (recording *recordingTipsListFilesystem) ReadDirEntries(path string) ([]fs.DirEntry, error) {
	recording.readPaths = append(recording.readPaths, path)
	recording.deliveredEntries = append(recording.deliveredEntries, recording.readEntries...)
	return recording.readEntries, recording.readErr
}

func (recording *recordingTipsListFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingTipsListFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingTipsListFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("recursive mkdir")
}

func (recording *recordingTipsListFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingTipsListFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingTipsListFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingTipsListFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingTipsListFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingTipsListFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingTipsListFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingTipsListFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingTipsListFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingTipsListFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected tips-list " + operation)
}

func (recording *recordingTipsListFilesystem) dependencies() listDependencies {
	return listDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingTipsListFilesystem) assertedReadPaths() ([]string, error) {
	if len(recording.readPaths) == 0 {
		return nil, errors.New("recorded tips-list directory-entry-read population is empty")
	}
	return recording.readPaths, nil
}

func (recording *recordingTipsListFilesystem) assertedDeliveredEntries() ([]fs.DirEntry, error) {
	if len(recording.deliveredEntries) == 0 {
		return nil, errors.New("recorded tips-list delivered-entry population is empty")
	}
	return recording.deliveredEntries, nil
}

type recordingTipsDirEntry struct {
	name       string
	directory  bool
	nameCalls  int
	isDirCalls int
	typeCalls  int
	infoCalls  int
}

func (entry *recordingTipsDirEntry) Name() string {
	entry.nameCalls++
	return entry.name
}

func (entry *recordingTipsDirEntry) IsDir() bool {
	entry.isDirCalls++
	return entry.directory
}

func (entry *recordingTipsDirEntry) Type() fs.FileMode {
	entry.typeCalls++
	return 0
}

func (entry *recordingTipsDirEntry) Info() (fs.FileInfo, error) {
	entry.infoCalls++
	return nil, errors.New("tips-list must not call DirEntry.Info")
}

func (entry *recordingTipsDirEntry) observationCounts() (int, int, int, int) {
	return entry.isDirCalls, entry.nameCalls, entry.typeCalls, entry.infoCalls
}

type recordingTipsListDirectory struct {
	path          string
	dirCalls      int
	filePathCalls int
}

func (directory *recordingTipsListDirectory) Dir() string {
	directory.dirCalls++
	return directory.path
}

func (directory *recordingTipsListDirectory) FilePath(string) (string, error) {
	directory.filePathCalls++
	return "", errors.New("unexpected tips-list Directory.FilePath call")
}

type recordingTipsListCloudConfig struct {
	config.GitCloudConfig
	directory           config.Directory
	implementationCalls int
}

func (cloudConfig *recordingTipsListCloudConfig) Implementation() config.Directory {
	cloudConfig.implementationCalls++
	return cloudConfig.directory
}

func TestListSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemListDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("tips-list selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("tips-list filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestListPreservesExactLocalDirPathOneReadErrorIdentityNilResultAndNoEntryInspection(t *testing.T) {
	directory := &recordingTipsListDirectory{path: `/complete tips root//with spaces/../backslash\segment-ø`}
	cloudConfig := &recordingTipsListCloudConfig{directory: directory}
	partial := &recordingTipsDirEntry{name: "must-not-be-inspected.md"}
	readError := errors.New("complete tips-list directory-entry-read dependency error")
	recording := &recordingTipsListFilesystem{
		readEntries: []fs.DirEntry{partial},
		readErr:     readError,
	}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("tips-list dependency lost its complete filesystem value: %#v", dependencies)
	}

	actual, err := list(dependencies, cloudConfig)

	if actual != nil || err != readError {
		t.Fatalf("tips-list read-error result was (%#v, %v), want (nil, exact error %v)", actual, err, readError)
	}
	wantPath := file.Path("%s/%s", directory.path, TipsDir)
	paths, populationErr := recording.assertedReadPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(paths, []string{wantPath}) {
		t.Fatalf("tips-list read paths were %#v, want one exact LocalDir path %#v", paths, []string{wantPath})
	}
	if _, populationErr := recording.assertedDeliveredEntries(); populationErr != nil {
		t.Fatal(populationErr)
	}
	if cloudConfig.implementationCalls != 1 || directory.dirCalls != 1 || directory.filePathCalls != 0 {
		t.Fatalf("LocalDir evaluation calls were Implementation=%d Dir=%d FilePath=%d, want 1, 1, 0",
			cloudConfig.implementationCalls, directory.dirCalls, directory.filePathCalls)
	}
	if isDir, name, entryType, info := partial.observationCounts(); isDir != 0 || name != 0 || entryType != 0 || info != 0 {
		t.Fatalf("tips-list inspected a partial entry after read error: IsDir=%d Name=%d Type=%d Info=%d", isDir, name, entryType, info)
	}
	if len(recording.unexpectedOperations) != 0 {
		t.Fatalf("tips-list invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
	}
}

func TestListPreservesDeliveredOrderEntryIdentityAndExactShortCircuitFiltering(t *testing.T) {
	directory := &recordingTipsDirEntry{name: "directory-name-must-not-be-read.md", directory: true}
	first := &recordingTipsDirEntry{name: "z-first.md"}
	uppercase := &recordingTipsDirEntry{name: "case-mismatch.MD"}
	extraSuffix := &recordingTipsDirEntry{name: "not-a-suffix.md.more"}
	second := &recordingTipsDirEntry{name: "a-second.md"}
	entries := []fs.DirEntry{directory, first, uppercase, extraSuffix, second}
	recording := &recordingTipsListFilesystem{readEntries: entries}
	cloudConfig := config.GitCloudConfig{Impl: config.DirConfig{Path: "/complete tips filtering root"}}

	actual, err := list(recording.dependencies(), cloudConfig)

	if err != nil {
		t.Fatalf("tips-list filtering returned an error: %v", err)
	}
	want := []os.DirEntry{first, second}
	if len(actual) != len(want) {
		t.Fatalf("tips-list result was %#v, want two exact delivered entries %#v", actual, want)
	}
	for index := range want {
		if actual[index] != want[index] {
			t.Fatalf("tips-list result %d was %#v, want exact delivered identity %#v", index, actual[index], want[index])
		}
	}
	delivered, populationErr := recording.assertedDeliveredEntries()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(delivered, entries) {
		t.Fatalf("tips-list delivered entries were %#v, want exact order and identities %#v", delivered, entries)
	}
	for _, entry := range []*recordingTipsDirEntry{directory, first, uppercase, extraSuffix, second} {
		isDir, name, entryType, info := entry.observationCounts()
		wantNameCalls := 1
		if entry.directory {
			wantNameCalls = 0
		}
		if isDir != 1 || name != wantNameCalls || entryType != 0 || info != 0 {
			t.Fatalf("tips-list entry %q observations were IsDir=%d Name=%d Type=%d Info=%d, want 1, %d, 0, 0",
				entry.name, isDir, name, entryType, info, wantNameCalls)
		}
	}
	if len(recording.unexpectedOperations) != 0 {
		t.Fatalf("tips-list filtering invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
	}
}

func TestListReturnsNilForEmptyAndAllFilteredPopulations(t *testing.T) {
	tests := []struct {
		name    string
		entries []fs.DirEntry
	}{
		{name: "non-nil empty", entries: []fs.DirEntry{}},
		{name: "all filtered", entries: []fs.DirEntry{
			&recordingTipsDirEntry{name: "directory.md", directory: true},
			&recordingTipsDirEntry{name: "ordinary.txt"},
			&recordingTipsDirEntry{name: "case-mismatch.MD"},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingTipsListFilesystem{readEntries: test.entries}

			actual, err := list(recording.dependencies(), config.GitCloudConfig{Impl: config.DirConfig{Path: "/complete nil tips root"}})

			if actual != nil || err != nil {
				t.Fatalf("tips-list nil result was (%#v, %v), want (nil, nil)", actual, err)
			}
			if _, populationErr := recording.assertedReadPaths(); populationErr != nil {
				t.Fatal(populationErr)
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("tips-list nil-result case invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
		})
	}
}

func TestListDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	actual, err := list(
		listDependencies{},
		config.GitCloudConfig{Impl: config.DirConfig{Path: "/developer/home/project/must-not-be-accessed"}},
	)

	if actual != nil || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe tips-list dependency default returned (%#v, %v), want (nil, %v)", actual, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedListRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingTipsListFilesystem{}

	if _, err := recording.assertedReadPaths(); err == nil {
		t.Fatal("empty recorded tips-list directory-entry-read population passed")
	}
	if _, err := recording.assertedDeliveredEntries(); err == nil {
		t.Fatal("empty recorded tips-list delivered-entry population passed")
	}
}
