package filesystem

import (
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
	closed bool
}

func (*recordedArchiveFile) Write(data []byte) (int, error) {
	return len(data), nil
}

func (file *recordedArchiveFile) Close() error {
	file.closed = true
	return nil
}

type recordingArchiveFilesystem struct {
	paths        []string
	destinations []File
	data         []byte
	file         *recordedArchiveFile
	createErr    error
	copyErr      error
}

func (*recordingArchiveFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected read")
}

func (*recordingArchiveFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected stat")
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
