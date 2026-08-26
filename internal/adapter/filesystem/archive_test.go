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

type recordedArchiveReadCloser struct {
	reader     io.Reader
	closeCalls int
	closeErr   error
}

func (reader *recordedArchiveReadCloser) Read(data []byte) (int, error) {
	return reader.reader.Read(data)
}

func (reader *recordedArchiveReadCloser) Close() error {
	reader.closeCalls++
	return reader.closeErr
}

type recordingArchiveFilesystem struct {
	paths          []string
	zipPaths       []string
	zipEntries     []*zip.File
	destinations   []File
	data           []byte
	file           *recordedArchiveFile
	zipReader      *zip.ReadCloser
	zipEntryReader io.ReadCloser
	closedFiles    []File
	closedReaders  []io.Closer
	createErr      error
	openZipError   error
	openEntryError error
	copyErr        error
	closeErr       error
	closeReaderErr error
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

func (recording *recordingArchiveFilesystem) OpenZipEntry(entry *zip.File) (io.ReadCloser, error) {
	recording.zipEntries = append(recording.zipEntries, entry)
	return recording.zipEntryReader, recording.openEntryError
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

func (recording *recordingArchiveFilesystem) CloseReader(reader io.Closer) error {
	recording.closedReaders = append(recording.closedReaders, reader)
	return recording.closeReaderErr
}

func (recording *recordingArchiveFilesystem) assertedZipPaths() ([]string, error) {
	if len(recording.zipPaths) == 0 {
		return nil, errors.New("recorded archive-open path population is empty")
	}
	return recording.zipPaths, nil
}

func (recording *recordingArchiveFilesystem) assertedZipEntries() ([]*zip.File, error) {
	if len(recording.zipEntries) == 0 {
		return nil, errors.New("recorded archive-entry-open population is empty")
	}
	return recording.zipEntries, nil
}

func (recording *recordingArchiveFilesystem) assertedClosedFiles() ([]File, error) {
	if len(recording.closedFiles) == 0 {
		return nil, errors.New("recorded archive file-close population is empty")
	}
	return recording.closedFiles, nil
}

func (recording *recordingArchiveFilesystem) assertedClosedReaders() ([]io.Closer, error) {
	if len(recording.closedReaders) == 0 {
		return nil, errors.New("recorded archive reader-close population is empty")
	}
	return recording.closedReaders, nil
}

func TestCloseReaderDependenciesDefaultToSafeNoReaderClose(t *testing.T) {
	reader := &recordedArchiveReadCloser{reader: strings.NewReader("must not be read or closed")}

	err := CloseReader(Dependencies{}, reader)

	if err != ErrNoFilesystem {
		t.Fatalf("zero-value reader close returned %v, want exact %v", err, ErrNoFilesystem)
	}
	if reader.closeCalls != 0 {
		t.Fatalf("zero-value reader close invoked the supplied reader %d times", reader.closeCalls)
	}
}

func TestCloseReaderPassesExactIdentityOnceAndReturnsExactDependencyError(t *testing.T) {
	reader := &recordedArchiveReadCloser{reader: strings.NewReader("complete reader bytes")}
	closeError := errors.New("complete reader-close dependency error")
	recording := &recordingArchiveFilesystem{closeReaderErr: closeError}

	err := CloseReader(Dependencies{FileSystem: recording}, reader)

	if err != closeError {
		t.Fatalf("reader-close error was %v, want exact dependency error %v", err, closeError)
	}
	readers, populationErr := recording.assertedClosedReaders()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(readers) != 1 || readers[0] != reader {
		t.Fatalf("reader-close dependency received %#v, want one exact reader identity %p", readers, reader)
	}
	if reader.closeCalls != 0 {
		t.Fatalf("reader-close forwarder bypassed its dependency and closed the reader %d times", reader.closeCalls)
	}
}

func TestSystemCloseReaderInvokesExactReaderOnceAndReturnsExactError(t *testing.T) {
	closeError := errors.New("complete system reader-close error")
	reader := &recordedArchiveReadCloser{reader: strings.NewReader("system reader bytes"), closeErr: closeError}

	err := CloseReader(System(), reader)

	if err != closeError || reader.closeCalls != 1 {
		t.Fatalf("system reader close returned %v after %d attempts, want exact %v after one attempt",
			err, reader.closeCalls, closeError)
	}
}

func TestRecordedCloseReaderRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingArchiveFilesystem{}

	if _, err := recording.assertedClosedReaders(); err == nil {
		t.Fatal("empty recorded reader-close population passed")
	}
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

func TestOpenZipEntryDependenciesDefaultToSafeNoEntryOpen(t *testing.T) {
	entry := &zip.File{FileHeader: zip.FileHeader{Method: 99}}

	reader, err := OpenZipEntry(Dependencies{}, entry)

	if reader != nil || err != ErrNoFilesystem {
		t.Fatalf("zero-value archive-entry open returned (%#v, %v), want (nil, exact %v)", reader, err, ErrNoFilesystem)
	}
}

func TestOpenZipEntryPassesExactEntryAndReturnsExactReaderAndError(t *testing.T) {
	entry := &zip.File{FileHeader: zip.FileHeader{Name: "complete entry identity.bin", Method: 99}}
	wantReader := &recordedArchiveReadCloser{reader: strings.NewReader("complete injected entry bytes")}
	wantError := errors.New("complete archive-entry-open dependency error")
	recording := &recordingArchiveFilesystem{zipEntryReader: wantReader, openEntryError: wantError}

	reader, err := OpenZipEntry(Dependencies{FileSystem: recording}, entry)

	if reader != wantReader || err != wantError {
		t.Fatalf("archive-entry open returned (%p, %v), want exact (%p, %v)", reader, err, wantReader, wantError)
	}
	entries, populationErr := recording.assertedZipEntries()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(entries) != 1 || entries[0] != entry {
		t.Fatalf("archive-entry-open dependency received %#v, want one exact entry identity %p", entries, entry)
	}
}

func TestRecordedArchiveEntryOpenRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingArchiveFilesystem{}

	if _, err := recording.assertedZipEntries(); err == nil {
		t.Fatal("empty recorded archive-entry-open population passed")
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
		opened, err := OpenZipEntry(System(), reader.File[index])
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
