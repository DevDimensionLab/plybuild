package shell

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/testutil"
)

type recordedUnzipMkdirAll struct {
	path string
	mode fs.FileMode
}

type recordedUnzipOpenFile struct {
	path  string
	flags int
	mode  fs.FileMode
}

type recordingUnzipFilesystem struct {
	mkdirs               []recordedUnzipMkdirAll
	mkdirErrors          map[int]error
	openFiles            []recordedUnzipOpenFile
	openErrors           map[int]error
	operationOrder       []string
	unexpectedOperations []string
}

var _ filesystem.FileSystem = (*recordingUnzipFilesystem)(nil)

func (recording *recordingUnzipFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingUnzipFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingUnzipFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingUnzipFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingUnzipFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingUnzipFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingUnzipFilesystem) MkdirAll(path string, mode fs.FileMode) error {
	recording.mkdirs = append(recording.mkdirs, recordedUnzipMkdirAll{path: path, mode: mode})
	recording.operationOrder = append(recording.operationOrder, "mkdir-all")
	return recording.mkdirErrors[len(recording.mkdirs)]
}

func (recording *recordingUnzipFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingUnzipFilesystem) OpenFile(path string, flags int, mode fs.FileMode) (*os.File, error) {
	recording.openFiles = append(recording.openFiles, recordedUnzipOpenFile{path: path, flags: flags, mode: mode})
	recording.operationOrder = append(recording.operationOrder, "open-file")
	return nil, recording.openErrors[len(recording.openFiles)]
}

func (recording *recordingUnzipFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingUnzipFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingUnzipFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingUnzipFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingUnzipFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingUnzipFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingUnzipFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingUnzipFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected unzip " + operation)
}

func (recording *recordingUnzipFilesystem) dependencies() unzipDependencies {
	return unzipDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingUnzipFilesystem) assertedMkdirs() ([]recordedUnzipMkdirAll, error) {
	if len(recording.mkdirs) == 0 {
		return nil, errors.New("recorded unzip recursive-directory population is empty")
	}
	return recording.mkdirs, nil
}

func (recording *recordingUnzipFilesystem) assertedOpenFiles() ([]recordedUnzipOpenFile, error) {
	if len(recording.openFiles) == 0 {
		return nil, errors.New("recorded unzip file-open population is empty")
	}
	return recording.openFiles, nil
}

func TestUnzipSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemUnzipDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Unzip selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Unzip filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestUnzipDirectoryEntriesPreserveJoinedAppendOrderModeAndOneAttemptPerEntry(t *testing.T) {
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "z-first directory/", directory: true},
		{name: "a-second directory/nested/", directory: true},
	})
	destination := filepath.Join(t.TempDir(), "complete destination with spaces")
	recording := &recordingUnzipFilesystem{}

	filenames, err := unzipWithDependencies(recording.dependencies(), archive, destination)

	if err != nil {
		t.Fatalf("directory-entry unzip returned an error: %v", err)
	}
	wantFilenames := []string{
		filepath.Join(destination, "z-first directory"),
		filepath.Join(destination, "a-second directory/nested"),
	}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("directory-entry filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	wantMkdirs := []recordedUnzipMkdirAll{
		{path: wantFilenames[0], mode: os.ModePerm},
		{path: wantFilenames[1], mode: os.ModePerm},
	}
	assertRecordedUnzipMkdirs(t, recording, wantMkdirs)
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipDirectoryEntryMkdirFailureReturnsExactErrorAndPartialFilenames(t *testing.T) {
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "first successful directory/", directory: true},
		{name: "second failing directory/", directory: true},
		{name: "must-not-be-reached/", directory: true},
	})
	destination := filepath.Join(t.TempDir(), "directory error destination")
	mkdirError := errors.New("complete unzip directory-entry dependency error")
	recording := &recordingUnzipFilesystem{mkdirErrors: map[int]error{2: mkdirError}}

	filenames, err := unzipWithDependencies(recording.dependencies(), archive, destination)

	if err != mkdirError {
		t.Fatalf("directory-entry mkdir error was %v, want exact dependency error %v", err, mkdirError)
	}
	wantFilenames := []string{
		filepath.Join(destination, "first successful directory"),
		filepath.Join(destination, "second failing directory"),
	}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("directory-entry partial filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	assertRecordedUnzipMkdirs(t, recording, []recordedUnzipMkdirAll{
		{path: wantFilenames[0], mode: os.ModePerm},
		{path: wantFilenames[1], mode: os.ModePerm},
	})
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipFileParentMkdirFailureReturnsExactErrorAndPartialFilenamesBeforeFileOpen(t *testing.T) {
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "first directory/", directory: true},
		{name: "nested parent/complete file.txt", data: []byte("must not be opened or copied")},
		{name: "must-not-be-reached/", directory: true},
	})
	destination := filepath.Join(t.TempDir(), "file error destination")
	mkdirError := errors.New("complete unzip file-parent dependency error")
	recording := &recordingUnzipFilesystem{mkdirErrors: map[int]error{2: mkdirError}}

	filenames, err := unzipWithDependencies(recording.dependencies(), archive, destination)

	if err != mkdirError {
		t.Fatalf("file-parent mkdir error was %v, want exact dependency error %v", err, mkdirError)
	}
	wantFilenames := []string{
		filepath.Join(destination, "first directory"),
		filepath.Join(destination, "nested parent/complete file.txt"),
	}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("file-parent partial filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	assertRecordedUnzipMkdirs(t, recording, []recordedUnzipMkdirAll{
		{path: wantFilenames[0], mode: os.ModePerm},
		{path: filepath.Dir(wantFilenames[1]), mode: os.ModePerm},
	})
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipFileOpenPreservesJoinedAppendEntryOrderParentFlagsModeErrorAndPrecedence(t *testing.T) {
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "first directory/", directory: true},
		{name: "nested parent/complete file.bin", mode: 0613, method: 99},
		{name: "must-not-be-reached/", directory: true},
	})
	destination := filepath.Join(t.TempDir(), "file open error destination")
	openError := errors.New("complete unzip file-open dependency error")
	recording := &recordingUnzipFilesystem{openErrors: map[int]error{1: openError}}

	filenames, err := unzipWithDependencies(recording.dependencies(), archive, destination)

	if err != openError {
		t.Fatalf("file-open error was %v, want exact dependency error %v before %v", err, openError, zip.ErrAlgorithm)
	}
	wantFilenames := []string{
		filepath.Join(destination, "first directory"),
		filepath.Join(destination, "nested parent/complete file.bin"),
	}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("file-open partial filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	assertRecordedUnzipMkdirs(t, recording, []recordedUnzipMkdirAll{
		{path: wantFilenames[0], mode: os.ModePerm},
		{path: filepath.Dir(wantFilenames[1]), mode: os.ModePerm},
	})
	assertRecordedUnzipOpenFiles(t, recording, []recordedUnzipOpenFile{{
		path:  wantFilenames[1],
		flags: os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
		mode:  0613,
	}})
	wantOrder := []string{"mkdir-all", "mkdir-all", "open-file"}
	if !reflect.DeepEqual(recording.operationOrder, wantOrder) {
		t.Fatalf("unzip filesystem operation order differs:\n got: %#v\nwant: %#v", recording.operationOrder, wantOrder)
	}
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipDependenciesDefaultToSafeNoFilesystemAtReachedDirectoryCreation(t *testing.T) {
	tests := []struct {
		name  string
		entry unzipArchiveEntry
	}{
		{name: "directory entry", entry: unzipArchiveEntry{name: "safe directory/", directory: true}},
		{name: "file parent", entry: unzipArchiveEntry{name: "safe parent/file.txt", data: []byte("must not be written")}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			archive := writeUnzipArchive(t, []unzipArchiveEntry{test.entry})
			destination := filepath.Join(t.TempDir(), "safe default destination")

			filenames, err := unzipWithDependencies(unzipDependencies{}, archive, destination)

			wantFilenames := []string{filepath.Join(destination, test.entry.name)}
			if err != filesystem.ErrNoFilesystem {
				t.Fatalf("safe unzip default returned %v, want exact %v", err, filesystem.ErrNoFilesystem)
			}
			if !reflect.DeepEqual(filenames, wantFilenames) {
				t.Fatalf("safe unzip partial filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
			}
			if _, statErr := os.Stat(filepath.Join(destination, "safe directory")); !os.IsNotExist(statErr) {
				t.Fatalf("safe unzip default created a directory: %v", statErr)
			}
			if _, statErr := os.Stat(filepath.Join(destination, "safe parent")); !os.IsNotExist(statErr) {
				t.Fatalf("safe unzip default created a file parent: %v", statErr)
			}
		})
	}
}

func TestRecordedUnzipRejectsEmptyMkdirAllPopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedMkdirs(); err == nil {
		t.Fatal("empty recorded unzip recursive-directory population passed")
	}
}

func TestRecordedUnzipRejectsEmptyOpenFilePopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedOpenFiles(); err == nil {
		t.Fatal("empty recorded unzip file-open population passed")
	}
}

type unzipArchiveEntry struct {
	name      string
	directory bool
	data      []byte
	mode      fs.FileMode
	method    uint16
}

func writeUnzipArchive(t *testing.T, entries []unzipArchiveEntry) string {
	t.Helper()
	var contents bytes.Buffer
	writer := zip.NewWriter(&contents)
	for _, entry := range entries {
		header := &zip.FileHeader{Name: entry.name, Method: entry.method}
		if entry.directory {
			header.SetMode(fs.ModeDir | 0755)
		} else {
			mode := entry.mode
			if mode == 0 {
				mode = 0644
			}
			header.SetMode(mode)
		}
		var archiveEntry io.Writer
		var err error
		if entry.method == zip.Store {
			archiveEntry, err = writer.CreateHeader(header)
		} else {
			archiveEntry, err = writer.CreateRaw(header)
		}
		if err != nil {
			t.Fatalf("create zip entry %q: %v", entry.name, err)
		}
		if _, err := archiveEntry.Write(entry.data); err != nil {
			t.Fatalf("write zip entry %q: %v", entry.name, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close zip fixture: %v", err)
	}
	path := filepath.Join(t.TempDir(), "unzip-contract.zip")
	if err := testutil.WriteFileOutsideWorkingTree(path, contents.Bytes(), 0600); err != nil {
		t.Fatalf("write zip fixture: %v", err)
	}
	return path
}

func assertRecordedUnzipMkdirs(t *testing.T, recording *recordingUnzipFilesystem, want []recordedUnzipMkdirAll) {
	t.Helper()
	mkdirs, err := recording.assertedMkdirs()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(mkdirs, want) {
		t.Fatalf("recorded unzip recursive-directory calls differ:\n got: %#v\nwant: %#v", mkdirs, want)
	}
}

func assertRecordedUnzipOpenFiles(t *testing.T, recording *recordingUnzipFilesystem, want []recordedUnzipOpenFile) {
	t.Helper()
	openFiles, err := recording.assertedOpenFiles()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(openFiles, want) {
		t.Fatalf("recorded unzip file-open calls differ:\n got: %#v\nwant: %#v", openFiles, want)
	}
}

func assertNoRecordedUnzipOpenFiles(t *testing.T, recording *recordingUnzipFilesystem) {
	t.Helper()
	if len(recording.openFiles) != 0 {
		t.Fatalf("unzip unexpectedly opened files: %#v", recording.openFiles)
	}
}

func assertNoUnexpectedUnzipOperations(t *testing.T, recording *recordingUnzipFilesystem) {
	t.Helper()
	if len(recording.unexpectedOperations) != 0 {
		t.Fatalf("unzip invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
	}
}
