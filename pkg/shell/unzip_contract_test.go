package shell

import (
	"archive/zip"
	"bytes"
	"errors"
	"hash/crc32"
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

type recordedUnzipCopy struct {
	destination filesystem.File
	data        []byte
}

type recordingUnzipEntryReadCloser struct {
	reader     io.Reader
	closeCalls int
	closeErr   error
}

func (reader *recordingUnzipEntryReadCloser) Read(data []byte) (int, error) {
	return reader.reader.Read(data)
}

func (reader *recordingUnzipEntryReadCloser) Close() error {
	reader.closeCalls++
	return reader.closeErr
}

type recordingUnzipFilesystem struct {
	archivePaths         []string
	openedArchives       map[int]*zip.ReadCloser
	archiveOpenErrors    map[int]error
	mkdirs               []recordedUnzipMkdirAll
	mkdirErrors          map[int]error
	openFiles            []recordedUnzipOpenFile
	openedFiles          map[int]*os.File
	openErrors           map[int]error
	copies               []recordedUnzipCopy
	copySources          []io.Reader
	copyCounts           map[int]int64
	copyErrors           map[int]error
	closedFiles          []filesystem.File
	closeErrors          map[int]error
	closedReaders        []io.ReadCloser
	closeReaderErrors    map[int]error
	closedArchives       []*zip.ReadCloser
	closeArchiveErrors   map[int]error
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
	return recording.openedFiles[len(recording.openFiles)], recording.openErrors[len(recording.openFiles)]
}

func (recording *recordingUnzipFilesystem) OpenZipReader(path string) (*zip.ReadCloser, error) {
	recording.archivePaths = append(recording.archivePaths, path)
	recording.operationOrder = append(recording.operationOrder, "open-zip-reader")
	return recording.openedArchives[len(recording.archivePaths)], recording.archiveOpenErrors[len(recording.archivePaths)]
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

func (recording *recordingUnzipFilesystem) Copy(destination filesystem.File, source io.Reader) (int64, error) {
	recording.copySources = append(recording.copySources, source)
	data, err := io.ReadAll(source)
	recording.copies = append(recording.copies, recordedUnzipCopy{
		destination: destination,
		data:        append([]byte(nil), data...),
	})
	recording.operationOrder = append(recording.operationOrder, "copy")
	if err != nil {
		return int64(len(data)), err
	}
	return recording.copyCounts[len(recording.copies)], recording.copyErrors[len(recording.copies)]
}

func (recording *recordingUnzipFilesystem) Close(file filesystem.File) error {
	recording.closedFiles = append(recording.closedFiles, file)
	recording.operationOrder = append(recording.operationOrder, "close")
	if err := recording.closeErrors[len(recording.closedFiles)]; err != nil {
		return err
	}
	return file.Close()
}

func (recording *recordingUnzipFilesystem) CloseReader(reader io.Closer) error {
	if archive, ok := reader.(*zip.ReadCloser); ok {
		recording.closedArchives = append(recording.closedArchives, archive)
		recording.operationOrder = append(recording.operationOrder, "close-archive")
		if err := recording.closeArchiveErrors[len(recording.closedArchives)]; err != nil {
			return err
		}
		return archive.Close()
	}
	entryReader, ok := reader.(io.ReadCloser)
	if !ok {
		return recording.unexpected("close reader")
	}
	recording.closedReaders = append(recording.closedReaders, entryReader)
	recording.operationOrder = append(recording.operationOrder, "close-reader")
	if err := recording.closeReaderErrors[len(recording.closedReaders)]; err != nil {
		return err
	}
	return entryReader.Close()
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

func (recording *recordingUnzipFilesystem) assertedArchivePaths() ([]string, error) {
	if len(recording.archivePaths) == 0 {
		return nil, errors.New("recorded unzip archive-open population is empty")
	}
	return recording.archivePaths, nil
}

func (recording *recordingUnzipFilesystem) assertedOpenFiles() ([]recordedUnzipOpenFile, error) {
	if len(recording.openFiles) == 0 {
		return nil, errors.New("recorded unzip file-open population is empty")
	}
	return recording.openFiles, nil
}

func (recording *recordingUnzipFilesystem) assertedCopies() ([]recordedUnzipCopy, error) {
	if len(recording.copies) == 0 {
		return nil, errors.New("recorded unzip copy population is empty")
	}
	return recording.copies, nil
}

func (recording *recordingUnzipFilesystem) assertedClosedFiles() ([]filesystem.File, error) {
	if len(recording.closedFiles) == 0 {
		return nil, errors.New("recorded unzip output-close population is empty")
	}
	return recording.closedFiles, nil
}

func (recording *recordingUnzipFilesystem) assertedClosedReaders() ([]io.ReadCloser, error) {
	if len(recording.closedReaders) == 0 {
		return nil, errors.New("recorded unzip entry-reader-close population is empty")
	}
	return recording.closedReaders, nil
}

func (recording *recordingUnzipFilesystem) assertedClosedArchives() ([]*zip.ReadCloser, error) {
	if len(recording.closedArchives) == 0 {
		return nil, errors.New("recorded unzip archive-close population is empty")
	}
	return recording.closedArchives, nil
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

func TestUnzipArchiveOpenPreservesExactSourceOneRequestReaderIdentityCloseAndTraversal(t *testing.T) {
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "z-first directory/", directory: true},
		{name: "a-second directory/nested/", directory: true},
	})
	injectedReader := openUnzipArchive(t, archive)
	recording := &recordingUnzipFilesystem{openedArchives: map[int]*zip.ReadCloser{1: injectedReader}}
	source := "/complete archive/source path with spaces/must-not-be-opened.zip"
	destination := filepath.Join(t.TempDir(), "injected archive destination")

	filenames, err := unzipWithDependencies(recording.dependencies(), source, destination)

	if err != nil {
		t.Fatalf("injected archive-reader unzip returned an error: %v", err)
	}
	wantFilenames := []string{
		filepath.Join(destination, "z-first directory"),
		filepath.Join(destination, "a-second directory/nested"),
	}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("injected archive-reader filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	paths, populationErr := recording.assertedArchivePaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(paths, []string{source}) {
		t.Fatalf("archive-open paths were %#v, want one exact source %#v", paths, []string{source})
	}
	assertRecordedUnzipArchiveClose(t, recording, injectedReader)
	if err := injectedReader.Close(); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("injected archive reader was not closed by Unzip: %v", err)
	}
	if !reflect.DeepEqual(recording.operationOrder, []string{"open-zip-reader", "mkdir-all", "mkdir-all", "close-archive"}) {
		t.Fatalf("injected archive-reader operation order differs: %#v", recording.operationOrder)
	}
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoRecordedUnzipCopies(t, recording)
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipArchiveOpenReturnsExactErrorNilPartialFilenamesAndNoLaterRequests(t *testing.T) {
	wantError := errors.New("complete unzip archive-open dependency error")
	recording := &recordingUnzipFilesystem{archiveOpenErrors: map[int]error{1: wantError}}
	source := "/complete failing archive/source path.zip"

	filenames, err := unzipWithDependencies(recording.dependencies(), source, filepath.Join(t.TempDir(), "must remain empty"))

	if filenames != nil || err != wantError {
		t.Fatalf("archive-open failure returned (%#v, %v), want (nil, exact %v)", filenames, err, wantError)
	}
	paths, populationErr := recording.assertedArchivePaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(paths, []string{source}) {
		t.Fatalf("failing archive-open paths were %#v, want %#v", paths, []string{source})
	}
	if !reflect.DeepEqual(recording.operationOrder, []string{"open-zip-reader"}) {
		t.Fatalf("archive-open failure operation order differs: %#v", recording.operationOrder)
	}
	if len(recording.closedArchives) != 0 {
		t.Fatalf("archive-open failure made archive-close requests: %#v", recording.closedArchives)
	}
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoRecordedUnzipCopies(t, recording)
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipDirectoryEntriesPreserveJoinedAppendOrderModeAndOneAttemptPerEntry(t *testing.T) {
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "z-first directory/", directory: true},
		{name: "a-second directory/nested/", directory: true},
	})
	destination := filepath.Join(t.TempDir(), "complete destination with spaces")
	recording := newRecordingUnzipFilesystem(t, archive)
	archiveCloseError := errors.New("complete ignored successful archive-close dependency error")
	recording.closeArchiveErrors = map[int]error{1: archiveCloseError}

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
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	wantOrder := []string{"open-zip-reader", "mkdir-all", "mkdir-all", "close-archive"}
	if !reflect.DeepEqual(recording.operationOrder, wantOrder) {
		t.Fatalf("directory-entry operation order differs:\n got: %#v\nwant: %#v", recording.operationOrder, wantOrder)
	}
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoRecordedUnzipCopies(t, recording)
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
	recording := newRecordingUnzipFilesystem(t, archive)
	recording.mkdirErrors = map[int]error{2: mkdirError}
	recording.closeArchiveErrors = map[int]error{1: errors.New("ignored archive-close error after directory error")}

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
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	wantOrder := []string{"open-zip-reader", "mkdir-all", "mkdir-all", "close-archive"}
	if !reflect.DeepEqual(recording.operationOrder, wantOrder) {
		t.Fatalf("directory-entry failure operation order differs:\n got: %#v\nwant: %#v", recording.operationOrder, wantOrder)
	}
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoRecordedUnzipCopies(t, recording)
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
	recording := newRecordingUnzipFilesystem(t, archive)
	recording.mkdirErrors = map[int]error{2: mkdirError}

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
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	wantOrder := []string{"open-zip-reader", "mkdir-all", "mkdir-all", "close-archive"}
	if !reflect.DeepEqual(recording.operationOrder, wantOrder) {
		t.Fatalf("file-parent failure operation order differs:\n got: %#v\nwant: %#v", recording.operationOrder, wantOrder)
	}
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoRecordedUnzipCopies(t, recording)
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
	recording := newRecordingUnzipFilesystem(t, archive)
	recording.openErrors = map[int]error{1: openError}

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
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	wantOrder := []string{"open-zip-reader", "mkdir-all", "mkdir-all", "open-file", "close-archive"}
	if !reflect.DeepEqual(recording.operationOrder, wantOrder) {
		t.Fatalf("unzip filesystem operation order differs:\n got: %#v\nwant: %#v", recording.operationOrder, wantOrder)
	}
	assertNoRecordedUnzipCopies(t, recording)
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipFileCopyPreservesOpenedDestinationsOrderedBytesAttemptsIgnoredResultsClosesAndTraversal(t *testing.T) {
	firstData := []byte{'f', 'i', 'r', 's', 't', 0x00, 0xff, '\n'}
	secondData := []byte("second arbitrary entry bytes\nwith a later line\n")
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "first directory/", directory: true},
		{name: "first parent/complete first.bin", data: firstData, mode: 0613, method: recordingUnzipCompressionMethod},
		{name: "second parent/complete second.txt", data: secondData, mode: 0642, method: recordingUnzipCompressionMethod},
	})
	destination := filepath.Join(t.TempDir(), "copy destination with spaces")
	firstReader, firstWriter := unzipPipe(t)
	secondReader, secondWriter := unzipPipe(t)
	copyError := errors.New("complete ignored unzip copy dependency error")
	recording := newRecordingUnzipFilesystem(t, archive)
	recording.openedFiles = map[int]*os.File{1: firstWriter, 2: secondWriter}
	recording.copyCounts = map[int]int64{1: 918273645, 2: -41}
	recording.copyErrors = map[int]error{1: copyError}
	recording.closeArchiveErrors = map[int]error{1: errors.New("complete ignored archive-close error after traversal")}
	var entryReaders []*recordingUnzipEntryReadCloser
	recordUnzipEntryCloses(recording.openedArchives[1], &entryReaders)
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("unzip dependency lost its complete filesystem value: %#v", dependencies)
	}

	filenames, err := unzipWithDependencies(dependencies, archive, destination)

	if err != nil {
		t.Fatalf("ignored copy count/error or later traversal changed unzip result: %v", err)
	}
	wantFilenames := []string{
		filepath.Join(destination, "first directory"),
		filepath.Join(destination, "first parent/complete first.bin"),
		filepath.Join(destination, "second parent/complete second.txt"),
	}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("file-copy filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	assertRecordedUnzipMkdirs(t, recording, []recordedUnzipMkdirAll{
		{path: wantFilenames[0], mode: os.ModePerm},
		{path: filepath.Dir(wantFilenames[1]), mode: os.ModePerm},
		{path: filepath.Dir(wantFilenames[2]), mode: os.ModePerm},
	})
	assertRecordedUnzipOpenFiles(t, recording, []recordedUnzipOpenFile{
		{path: wantFilenames[1], flags: os.O_WRONLY | os.O_CREATE | os.O_TRUNC, mode: 0613},
		{path: wantFilenames[2], flags: os.O_WRONLY | os.O_CREATE | os.O_TRUNC, mode: 0642},
	})
	assertRecordedUnzipCopies(t, recording, []recordedUnzipCopy{
		{destination: firstWriter, data: firstData},
		{destination: secondWriter, data: secondData},
	})
	assertRecordedUnzipCloses(t, recording, []filesystem.File{firstWriter, secondWriter})
	assertRecordedUnzipReaderCloses(t, recording, recordedUnzipCopyReaders(t, recording))
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	wantOrder := []string{
		"open-zip-reader",
		"mkdir-all",
		"mkdir-all", "open-file", "copy", "close", "close-reader",
		"mkdir-all", "open-file", "copy", "close", "close-reader",
		"close-archive",
	}
	if !reflect.DeepEqual(recording.operationOrder, wantOrder) {
		t.Fatalf("unzip file-copy operation order differs:\n got: %#v\nwant: %#v", recording.operationOrder, wantOrder)
	}
	assertUnzipOutputFileClosedWithoutBytes(t, firstReader, firstWriter)
	assertUnzipOutputFileClosedWithoutBytes(t, secondReader, secondWriter)
	if len(entryReaders) != 2 || entryReaders[0].closeCalls != 1 || entryReaders[1].closeCalls != 1 {
		t.Fatalf("successful output closes produced entry readers %#v with close counts, want two readers closed once", entryReaders)
	}
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipOutputCloseReturnsExactErrorAfterCopyAndStopsBeforeEntryCloseAndLaterTraversal(t *testing.T) {
	firstData := []byte("first bytes copied before output-close error\n")
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "first parent/first.bin", data: firstData, mode: 0613, method: recordingUnzipCompressionMethod},
		{name: "second parent/must-not-be-reached.txt", data: []byte("later bytes"), method: recordingUnzipCompressionMethod},
	})
	destination := filepath.Join(t.TempDir(), "output close error destination")
	_, firstWriter := unzipPipe(t)
	_, secondWriter := unzipPipe(t)
	closeError := errors.New("complete unzip output-close dependency error")
	recording := newRecordingUnzipFilesystem(t, archive)
	recording.openedFiles = map[int]*os.File{1: firstWriter, 2: secondWriter}
	recording.closeErrors = map[int]error{1: closeError}
	var entryReaders []*recordingUnzipEntryReadCloser
	recordUnzipEntryCloses(recording.openedArchives[1], &entryReaders)

	filenames, err := unzipWithDependencies(recording.dependencies(), archive, destination)

	if err != closeError {
		t.Fatalf("output-close error was %v, want exact dependency error %v", err, closeError)
	}
	wantFilenames := []string{filepath.Join(destination, "first parent/first.bin")}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("output-close partial filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	assertRecordedUnzipMkdirs(t, recording, []recordedUnzipMkdirAll{{
		path: filepath.Dir(wantFilenames[0]), mode: os.ModePerm,
	}})
	assertRecordedUnzipOpenFiles(t, recording, []recordedUnzipOpenFile{{
		path: wantFilenames[0], flags: os.O_WRONLY | os.O_CREATE | os.O_TRUNC, mode: 0613,
	}})
	assertRecordedUnzipCopies(t, recording, []recordedUnzipCopy{{destination: firstWriter, data: firstData}})
	assertRecordedUnzipCloses(t, recording, []filesystem.File{firstWriter})
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	wantOrder := []string{"open-zip-reader", "mkdir-all", "open-file", "copy", "close", "close-archive"}
	if !reflect.DeepEqual(recording.operationOrder, wantOrder) {
		t.Fatalf("output-close failure operation order differs:\n got: %#v\nwant: %#v", recording.operationOrder, wantOrder)
	}
	if len(entryReaders) != 1 || entryReaders[0].closeCalls != 0 {
		t.Fatalf("output-close failure produced entry readers %#v, want first reader left unclosed and no later reader", entryReaders)
	}
	if len(recording.closedReaders) != 0 {
		t.Fatalf("output-close failure made entry-reader-close requests: %#v", recording.closedReaders)
	}
	if _, statErr := firstWriter.Stat(); statErr != nil {
		t.Fatalf("recording output-close error unexpectedly closed the supplied file: %v", statErr)
	}
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipEntryReaderCloseReturnsExactErrorAfterCopyAndOutputCloseAndStopsLaterTraversal(t *testing.T) {
	firstData := []byte("first bytes copied before entry-reader-close error\n")
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "first parent/first.bin", data: firstData, mode: 0613, method: recordingUnzipCompressionMethod},
		{name: "second parent/must-not-be-reached.txt", data: []byte("later bytes"), method: recordingUnzipCompressionMethod},
	})
	destination := filepath.Join(t.TempDir(), "entry reader close error destination")
	firstReader, firstWriter := unzipPipe(t)
	_, secondWriter := unzipPipe(t)
	closeError := errors.New("complete unzip entry-reader-close dependency error")
	recording := newRecordingUnzipFilesystem(t, archive)
	recording.openedFiles = map[int]*os.File{1: firstWriter, 2: secondWriter}
	recording.closeReaderErrors = map[int]error{1: closeError}
	recording.closeArchiveErrors = map[int]error{1: errors.New("ignored archive-close error after entry-reader error")}
	var entryReaders []*recordingUnzipEntryReadCloser
	recordUnzipEntryCloses(recording.openedArchives[1], &entryReaders)

	filenames, err := unzipWithDependencies(recording.dependencies(), archive, destination)

	if err != closeError {
		t.Fatalf("entry-reader-close error was %v, want exact dependency error %v", err, closeError)
	}
	wantFilenames := []string{filepath.Join(destination, "first parent/first.bin")}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("entry-reader-close partial filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	assertRecordedUnzipMkdirs(t, recording, []recordedUnzipMkdirAll{{
		path: filepath.Dir(wantFilenames[0]), mode: os.ModePerm,
	}})
	assertRecordedUnzipOpenFiles(t, recording, []recordedUnzipOpenFile{{
		path: wantFilenames[0], flags: os.O_WRONLY | os.O_CREATE | os.O_TRUNC, mode: 0613,
	}})
	assertRecordedUnzipCopies(t, recording, []recordedUnzipCopy{{destination: firstWriter, data: firstData}})
	assertRecordedUnzipCloses(t, recording, []filesystem.File{firstWriter})
	if len(entryReaders) != 1 {
		t.Fatalf("entry-reader-close failure opened readers %#v, want exactly one", entryReaders)
	}
	assertRecordedUnzipReaderCloses(t, recording, recordedUnzipCopyReaders(t, recording))
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	wantOrder := []string{"open-zip-reader", "mkdir-all", "open-file", "copy", "close", "close-reader", "close-archive"}
	if !reflect.DeepEqual(recording.operationOrder, wantOrder) {
		t.Fatalf("entry-reader-close failure operation order differs:\n got: %#v\nwant: %#v", recording.operationOrder, wantOrder)
	}
	assertUnzipOutputFileClosedWithoutBytes(t, firstReader, firstWriter)
	if entryReaders[0].closeCalls != 0 {
		t.Fatalf("entry-reader-close dependency error directly closed the reader %d times", entryReaders[0].closeCalls)
	}
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipEntryOpenFailurePreventsCopyAfterParentAndFileOpen(t *testing.T) {
	archive := writeUnzipArchive(t, []unzipArchiveEntry{
		{name: "unsupported parent/unsupported entry.bin", data: []byte("raw unsupported bytes"), mode: 0607, method: 99},
		{name: "must-not-be-reached.txt", data: []byte("must not be opened or copied")},
	})
	destination := filepath.Join(t.TempDir(), "entry open error destination")
	_, writer := unzipPipe(t)
	recording := newRecordingUnzipFilesystem(t, archive)
	recording.openedFiles = map[int]*os.File{1: writer}

	filenames, err := unzipWithDependencies(recording.dependencies(), archive, destination)

	if err != zip.ErrAlgorithm {
		t.Fatalf("entry-open error was %v, want exact %v", err, zip.ErrAlgorithm)
	}
	wantFilenames := []string{filepath.Join(destination, "unsupported parent/unsupported entry.bin")}
	if !reflect.DeepEqual(filenames, wantFilenames) {
		t.Fatalf("entry-open partial filenames differ:\n got: %#v\nwant: %#v", filenames, wantFilenames)
	}
	assertRecordedUnzipMkdirs(t, recording, []recordedUnzipMkdirAll{{
		path: filepath.Dir(wantFilenames[0]), mode: os.ModePerm,
	}})
	assertRecordedUnzipOpenFiles(t, recording, []recordedUnzipOpenFile{{
		path: wantFilenames[0], flags: os.O_WRONLY | os.O_CREATE | os.O_TRUNC, mode: 0607,
	}})
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	assertNoRecordedUnzipCopies(t, recording)
	if !reflect.DeepEqual(recording.operationOrder, []string{"open-zip-reader", "mkdir-all", "open-file", "close-archive"}) {
		t.Fatalf("entry-open failure operation order differs: %#v", recording.operationOrder)
	}
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipZipSlipReturnsExactErrorBeforeTraversalAndStillDefersExactArchiveClose(t *testing.T) {
	archive := writeUnzipArchive(t, []unzipArchiveEntry{{
		name: "../must-not-escape.txt", data: []byte("must not be copied"),
	}})
	destination := filepath.Join(t.TempDir(), "zip slip destination")
	recording := newRecordingUnzipFilesystem(t, archive)
	recording.closeArchiveErrors = map[int]error{1: errors.New("ignored archive-close error after zip-slip")}

	filenames, err := unzipWithDependencies(recording.dependencies(), archive, destination)

	wantError := filepath.Join(destination, "../must-not-escape.txt") + ": illegal file path"
	if filenames != nil || err == nil || err.Error() != wantError {
		t.Fatalf("zip-slip returned (%#v, %v), want (nil, %q)", filenames, err, wantError)
	}
	assertRecordedUnzipArchiveClose(t, recording, recording.openedArchives[1])
	if !reflect.DeepEqual(recording.operationOrder, []string{"open-zip-reader", "close-archive"}) {
		t.Fatalf("zip-slip operation order differs: %#v", recording.operationOrder)
	}
	assertNoRecordedUnzipOpenFiles(t, recording)
	assertNoRecordedUnzipCopies(t, recording)
	if len(recording.closedReaders) != 0 {
		t.Fatalf("zip-slip made entry-reader-close requests: %#v", recording.closedReaders)
	}
	assertNoUnexpectedUnzipOperations(t, recording)
}

func TestUnzipDependenciesDefaultToSafeNoArchiveOpen(t *testing.T) {
	source := "/developer/archive/path/must-not-be-opened.zip"
	destination := filepath.Join(t.TempDir(), "safe default destination")

	filenames, err := unzipWithDependencies(unzipDependencies{}, source, destination)

	if err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe unzip default returned %v, want exact %v", err, filesystem.ErrNoFilesystem)
	}
	if filenames != nil {
		t.Fatalf("safe unzip default partial filenames were %#v, want nil", filenames)
	}
	if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
		t.Fatalf("safe unzip default created a destination: %v", statErr)
	}
}

func TestRecordedUnzipRejectsEmptyMkdirAllPopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedMkdirs(); err == nil {
		t.Fatal("empty recorded unzip recursive-directory population passed")
	}
}

func TestRecordedUnzipRejectsEmptyArchiveOpenPopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedArchivePaths(); err == nil {
		t.Fatal("empty recorded unzip archive-open population passed")
	}
}

func TestRecordedUnzipRejectsEmptyOpenFilePopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedOpenFiles(); err == nil {
		t.Fatal("empty recorded unzip file-open population passed")
	}
}

func TestRecordedUnzipRejectsEmptyCopyPopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedCopies(); err == nil {
		t.Fatal("empty recorded unzip copy population passed")
	}
}

func TestRecordedUnzipRejectsEmptyOutputClosePopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedClosedFiles(); err == nil {
		t.Fatal("empty recorded unzip output-close population passed")
	}
}

func TestRecordedUnzipRejectsEmptyEntryReaderClosePopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedClosedReaders(); err == nil {
		t.Fatal("empty recorded unzip entry-reader-close population passed")
	}
}

func TestRecordedUnzipRejectsEmptyArchiveClosePopulation(t *testing.T) {
	recording := &recordingUnzipFilesystem{}

	if _, err := recording.assertedClosedArchives(); err == nil {
		t.Fatal("empty recorded unzip archive-close population passed")
	}
}

const recordingUnzipCompressionMethod uint16 = 93

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
			header.CRC32 = crc32.ChecksumIEEE(entry.data)
			header.CompressedSize64 = uint64(len(entry.data))
			header.UncompressedSize64 = uint64(len(entry.data))
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

func openUnzipArchive(t *testing.T, path string) *zip.ReadCloser {
	t.Helper()
	reader, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("open zip fixture: %v", err)
	}
	t.Cleanup(func() { _ = reader.Close() })
	return reader
}

func newRecordingUnzipFilesystem(t *testing.T, archive string) *recordingUnzipFilesystem {
	t.Helper()
	return &recordingUnzipFilesystem{openedArchives: map[int]*zip.ReadCloser{1: openUnzipArchive(t, archive)}}
}

func recordUnzipEntryCloses(reader *zip.ReadCloser, readers *[]*recordingUnzipEntryReadCloser) {
	reader.RegisterDecompressor(recordingUnzipCompressionMethod, func(source io.Reader) io.ReadCloser {
		opened := &recordingUnzipEntryReadCloser{reader: source}
		*readers = append(*readers, opened)
		return opened
	})
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

func assertRecordedUnzipCopies(t *testing.T, recording *recordingUnzipFilesystem, want []recordedUnzipCopy) {
	t.Helper()
	copies, err := recording.assertedCopies()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(copies, want) {
		t.Fatalf("recorded unzip copy calls differ:\n got: %#v\nwant: %#v", copies, want)
	}
}

func assertRecordedUnzipCloses(t *testing.T, recording *recordingUnzipFilesystem, want []filesystem.File) {
	t.Helper()
	files, err := recording.assertedClosedFiles()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("recorded unzip output-close calls differ:\n got: %#v\nwant: %#v", files, want)
	}
}

func assertRecordedUnzipReaderCloses(t *testing.T, recording *recordingUnzipFilesystem, want []io.ReadCloser) {
	t.Helper()
	readers, err := recording.assertedClosedReaders()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(readers, want) {
		t.Fatalf("recorded unzip entry-reader-close calls differ:\n got: %#v\nwant: %#v", readers, want)
	}
}

func assertRecordedUnzipArchiveClose(t *testing.T, recording *recordingUnzipFilesystem, want *zip.ReadCloser) {
	t.Helper()
	archives, err := recording.assertedClosedArchives()
	if err != nil {
		t.Fatal(err)
	}
	if len(archives) != 1 || archives[0] != want {
		t.Fatalf("recorded unzip archive-close calls were %#v, want one exact archive identity %p", archives, want)
	}
	for _, reader := range recording.closedReaders {
		if any(reader) == any(want) {
			t.Fatalf("archive-close identity was mixed into entry-reader-close requests: %#v", recording.closedReaders)
		}
	}
}

func recordedUnzipCopyReaders(t *testing.T, recording *recordingUnzipFilesystem) []io.ReadCloser {
	t.Helper()
	if len(recording.copySources) == 0 {
		t.Fatal("recorded unzip copy-reader population is empty")
	}
	readers := make([]io.ReadCloser, len(recording.copySources))
	for index, source := range recording.copySources {
		reader, ok := source.(io.ReadCloser)
		if !ok {
			t.Fatalf("recorded unzip copy source %d is %T, want exact opened io.ReadCloser", index, source)
		}
		readers[index] = reader
	}
	return readers
}

func assertNoRecordedUnzipOpenFiles(t *testing.T, recording *recordingUnzipFilesystem) {
	t.Helper()
	if len(recording.openFiles) != 0 {
		t.Fatalf("unzip unexpectedly opened files: %#v", recording.openFiles)
	}
}

func assertNoRecordedUnzipCopies(t *testing.T, recording *recordingUnzipFilesystem) {
	t.Helper()
	if len(recording.copies) != 0 {
		t.Fatalf("unzip unexpectedly copied files: %#v", recording.copies)
	}
}

func unzipPipe(t *testing.T) (*os.File, *os.File) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create unzip output pipe: %v", err)
	}
	t.Cleanup(func() {
		_ = reader.Close()
		_ = writer.Close()
	})
	return reader, writer
}

func assertUnzipOutputFileClosedWithoutBytes(t *testing.T, reader, writer *os.File) {
	t.Helper()
	if _, err := writer.Stat(); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("unzip output file remained open or returned an unexpected stat error: %v", err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read closed unzip output pipe: %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("recording unzip copy unexpectedly wrote destination bytes: %v", data)
	}
}

func assertNoUnexpectedUnzipOperations(t *testing.T, recording *recordingUnzipFilesystem) {
	t.Helper()
	if len(recording.unexpectedOperations) != 0 {
		t.Fatalf("unzip invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
	}
}
