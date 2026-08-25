package filesystem

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/testutil"
)

type recordedArchiveFile struct {
	closed     bool
	closeCalls int
	closeErr   error
}

func (*recordedArchiveFile) Write(data []byte) (int, error) {
	return len(data), nil
}

func (file *recordedArchiveFile) Close() error {
	file.closeCalls++
	file.closed = true
	return file.closeErr
}

type recordingArchiveFilesystem struct {
	paths        []string
	zipPaths     []string
	destinations []File
	data         []byte
	file         *recordedArchiveFile
	zipReader    *zip.ReadCloser
	closedFiles  []File
	createErr    error
	openZipError error
	copyErr      error
	closeErr     error
}

func (*recordingArchiveFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected working directory")
}

func (*recordingArchiveFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected read")
}

func (*recordingArchiveFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected read directory")
}

func (*recordingArchiveFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected read directory entries")
}

func (*recordingArchiveFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected stat")
}

func (*recordingArchiveFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected mkdir")
}

func (*recordingArchiveFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected mkdir")
}

func (*recordingArchiveFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected write file")
}

func (*recordingArchiveFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected open file")
}

func (recording *recordingArchiveFilesystem) OpenZipReader(path string) (*zip.ReadCloser, error) {
	recording.zipPaths = append(recording.zipPaths, path)
	return recording.zipReader, recording.openZipError
}

func (*recordingArchiveFilesystem) Remove(string) error {
	return errors.New("unexpected remove")
}

func (*recordingArchiveFilesystem) RemoveAll(string) error {
	return errors.New("unexpected recursive remove")
}

func (*recordingArchiveFilesystem) Rename(string, string) error {
	return errors.New("unexpected rename")
}

func (*recordingArchiveFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected glob")
}

func (*recordingArchiveFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected walk")
}

func (recording *recordingArchiveFilesystem) Create(path string) (File, error) {
	recording.paths = append(recording.paths, path)
	if recording.createErr != nil {
		return nil, recording.createErr
	}
	recording.file = &recordedArchiveFile{}
	return recording.file, nil
}

func (recording *recordingArchiveFilesystem) Copy(destination File, source io.Reader) (int64, error) {
	recording.destinations = append(recording.destinations, destination)
	data, err := io.ReadAll(source)
	recording.data = append([]byte(nil), data...)
	if err != nil {
		return int64(len(data)), err
	}
	return int64(len(data)), recording.copyErr
}

func (recording *recordingArchiveFilesystem) Close(file File) error {
	recording.closedFiles = append(recording.closedFiles, file)
	return recording.closeErr
}

func (recording *recordingArchiveFilesystem) assertedZipPaths() ([]string, error) {
	if len(recording.zipPaths) == 0 {
		return nil, errors.New("recorded archive-open path population is empty")
	}
	return recording.zipPaths, nil
}

func (recording *recordingArchiveFilesystem) assertedClosedFiles() ([]File, error) {
	if len(recording.closedFiles) == 0 {
		return nil, errors.New("recorded archive file-close population is empty")
	}
	return recording.closedFiles, nil
}

func TestCloseDependenciesDefaultToSafeNoFileClose(t *testing.T) {
	file := &recordedArchiveFile{}

	err := Close(Dependencies{}, file)

	if err != ErrNoFilesystem {
		t.Fatalf("zero-value file close returned %v, want exact %v", err, ErrNoFilesystem)
	}
	if file.closeCalls != 0 || file.closed {
		t.Fatalf("zero-value file close invoked the supplied file %d times (closed=%t)", file.closeCalls, file.closed)
	}
}

func TestClosePassesExactFileIdentityOnceAndReturnsExactDependencyError(t *testing.T) {
	file := &recordedArchiveFile{}
	closeError := errors.New("complete file-close dependency error")
	recording := &recordingArchiveFilesystem{closeErr: closeError}

	err := Close(Dependencies{FileSystem: recording}, file)

	if err != closeError {
		t.Fatalf("file-close error was %v, want exact dependency error %v", err, closeError)
	}
	files, populationErr := recording.assertedClosedFiles()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(files) != 1 || files[0] != file {
		t.Fatalf("file-close dependency received %#v, want one exact file identity %p", files, file)
	}
	if file.closeCalls != 0 || file.closed {
		t.Fatalf("file-close forwarder bypassed its dependency and closed the file %d times", file.closeCalls)
	}
}

func TestSystemCloseInvokesExactFileOnceAndReturnsExactError(t *testing.T) {
	closeError := errors.New("complete system file-close error")
	file := &recordedArchiveFile{closeErr: closeError}

	err := Close(System(), file)

	if err != closeError || file.closeCalls != 1 || !file.closed {
		t.Fatalf("system file close returned %v after %d attempts (closed=%t), want exact %v after one attempt",
			err, file.closeCalls, file.closed, closeError)
	}
}

func TestRecordedCloseRejectsEmptyFilePopulation(t *testing.T) {
	recording := &recordingArchiveFilesystem{}

	if _, err := recording.assertedClosedFiles(); err == nil {
		t.Fatal("empty recorded file-close population passed")
	}
}

func TestOpenZipReaderDependenciesDefaultToSafeNoArchiveOpen(t *testing.T) {
	reader, err := OpenZipReader(Dependencies{}, "/developer/archive/path/must-not-be-opened.zip")

	if reader != nil || !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value archive open returned (%#v, %v), want (nil, %v)", reader, err, ErrNoFilesystem)
	}
}

func TestOpenZipReaderPassesExactPathReaderIdentityAndError(t *testing.T) {
	wantReader := &zip.ReadCloser{}
	wantError := errors.New("complete archive-open dependency error")
	recording := &recordingArchiveFilesystem{zipReader: wantReader, openZipError: wantError}
	path := "/complete archive/source path with spaces.zip"

	reader, err := OpenZipReader(Dependencies{FileSystem: recording}, path)

	if reader != wantReader || err != wantError {
		t.Fatalf("archive open returned (%p, %v), want exact (%p, %v)", reader, err, wantReader, wantError)
	}
	paths, populationErr := recording.assertedZipPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(paths, []string{path}) {
		t.Fatalf("archive-open paths were %#v, want exact supplied path %#v", paths, []string{path})
	}
}

func TestRecordedArchiveOpenRejectsEmptyPathPopulation(t *testing.T) {
	recording := &recordingArchiveFilesystem{}

	if _, err := recording.assertedZipPaths(); err == nil {
		t.Fatal("empty recorded archive-open path population passed")
	}
}

func TestSystemOpenZipReaderPreservesEntryOrderNamesAndBytes(t *testing.T) {
	var contents bytes.Buffer
	writer := zip.NewWriter(&contents)
	entries := []struct {
		name string
		data string
	}{
		{name: "z-first entry.txt", data: "first archive bytes\n"},
		{name: "a-second entry.bin", data: "second\x00archive bytes"},
	}
	for _, entry := range entries {
		archiveEntry, err := writer.Create(entry.name)
		if err != nil {
			t.Fatalf("create archive entry %q: %v", entry.name, err)
		}
		if _, err := archiveEntry.Write([]byte(entry.data)); err != nil {
			t.Fatalf("write archive entry %q: %v", entry.name, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close archive fixture: %v", err)
	}
	path := filepath.Join(t.TempDir(), "system archive open.zip")
	if err := testutil.WriteFileOutsideWorkingTree(path, contents.Bytes(), 0600); err != nil {
		t.Fatalf("write archive fixture: %v", err)
	}

	reader, err := OpenZipReader(System(), path)
	if err != nil {
		t.Fatalf("system archive open returned an error: %v", err)
	}
	defer func() { _ = reader.Close() }()
	if len(reader.File) != len(entries) {
		t.Fatalf("system archive entry count was %d, want %d", len(reader.File), len(entries))
	}
	for index, entry := range entries {
		if reader.File[index].Name != entry.name {
			t.Fatalf("system archive entry %d name was %q, want %q", index, reader.File[index].Name, entry.name)
		}
		opened, err := reader.File[index].Open()
		if err != nil {
			t.Fatalf("open system archive entry %q: %v", entry.name, err)
		}
		data, readErr := io.ReadAll(opened)
		closeErr := opened.Close()
		if readErr != nil || closeErr != nil || string(data) != entry.data {
			t.Fatalf("system archive entry %q returned (%q, %v, %v), want (%q, nil, nil)", entry.name, data, readErr, closeErr, entry.data)
		}
	}
}

func TestArchiveDependenciesDefaultToSafeNoCreateOrCopy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "safe-default", "archive.zip")

	file, createErr := Create(Dependencies{}, path)
	copied, copyErr := Copy(Dependencies{}, &recordedArchiveFile{}, strings.NewReader("must not be copied"))

	if file != nil || !errors.Is(createErr, ErrNoFilesystem) {
		t.Fatalf("zero-value create returned (%#v, %v), want (nil, %v)", file, createErr, ErrNoFilesystem)
	}
	if copied != 0 || !errors.Is(copyErr, ErrNoFilesystem) {
		t.Fatalf("zero-value copy returned (%d, %v), want (0, %v)", copied, copyErr, ErrNoFilesystem)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("zero-value archive dependencies mutated the target: %v", err)
	}
}

func TestArchiveOperationsPassCompletePathFileAndBody(t *testing.T) {
	recording := &recordingArchiveFilesystem{}
	dependencies := Dependencies{FileSystem: recording}
	path := "/complete archive/path with spaces/spring-1700000000.zip"
	body := "complete archive bytes\nwith a second line\n"

	file, err := Create(dependencies, path)
	if err != nil {
		t.Fatalf("create dependency returned an error: %v", err)
	}
	copied, err := Copy(dependencies, file, strings.NewReader(body))
	if err != nil {
		t.Fatalf("copy dependency returned an error: %v", err)
	}

	if copied != int64(len(body)) || !reflect.DeepEqual(recording.paths, []string{path}) ||
		!reflect.DeepEqual(recording.destinations, []File{file}) || string(recording.data) != body {
		t.Fatalf("archive dependency received incomplete values: copied=%d paths=%#v destinations=%#v data=%q",
			copied, recording.paths, recording.destinations, recording.data)
	}
}

func TestArchiveOperationsReturnDependencyErrors(t *testing.T) {
	createError := errors.New("complete create dependency error")
	createRecording := &recordingArchiveFilesystem{createErr: createError}
	file, err := Create(Dependencies{FileSystem: createRecording}, "/complete/create/error.zip")
	if file != nil || !errors.Is(err, createError) {
		t.Fatalf("create error result was (%#v, %v), want (nil, %v)", file, err, createError)
	}

	copyError := errors.New("complete copy dependency error")
	copyRecording := &recordingArchiveFilesystem{copyErr: copyError}
	destination := &recordedArchiveFile{}
	copied, err := Copy(Dependencies{FileSystem: copyRecording}, destination, strings.NewReader("copied before error"))
	if copied != int64(len("copied before error")) || !errors.Is(err, copyError) {
		t.Fatalf("copy error result was (%d, %v), want (%d, %v)", copied, err, len("copied before error"), copyError)
	}
}

func TestSystemCreateTruncatesAndCopiesArchiveBytes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "system archive.zip")
	if err := testutil.WriteFileOutsideWorkingTree(path, []byte("long existing contents that must be truncated"), 0600); err != nil {
		t.Fatalf("create existing archive fixture: %v", err)
	}
	body := "short system archive"

	file, err := Create(System(), path)
	if err != nil {
		t.Fatalf("system create returned an error: %v", err)
	}
	copied, err := Copy(System(), file, strings.NewReader(body))
	if err != nil {
		t.Fatalf("system copy returned an error: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("system archive close returned an error: %v", err)
	}
	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read system archive: %v", err)
	}
	if copied != int64(len(body)) || string(written) != body {
		t.Fatalf("system archive copy was (%d, %q), want (%d, %q)", copied, written, len(body), body)
	}
}
