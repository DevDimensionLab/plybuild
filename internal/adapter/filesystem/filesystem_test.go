package filesystem

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/testutil"
)

type recordedFilesystemOperation struct {
	Name  string
	Path  string
	Data  []byte
	Flags int
	Mode  fs.FileMode
}

type recordingFilesystem struct {
	operations []recordedFilesystemOperation
	readData   []byte
	fileInfo   fs.FileInfo
	readErr    error
	statErr    error
	mkdirErr   error
	writeErr   error
	opened     *os.File
	openErr    error
	removeErr  error
}

func (recording *recordingFilesystem) ReadFile(path string) ([]byte, error) {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "read", Path: path})
	return append([]byte(nil), recording.readData...), recording.readErr
}

func (recording *recordingFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "stat", Path: path})
	return recording.fileInfo, recording.statErr
}

func (recording *recordingFilesystem) MkdirAll(path string, mode fs.FileMode) error {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "mkdir-all", Path: path, Mode: mode})
	return recording.mkdirErr
}

func (recording *recordingFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.operations = append(recording.operations, recordedFilesystemOperation{
		Name: "write", Path: path, Data: append([]byte(nil), data...), Mode: mode,
	})
	return recording.writeErr
}

func (recording *recordingFilesystem) OpenFile(path string, flags int, mode fs.FileMode) (*os.File, error) {
	recording.operations = append(recording.operations, recordedFilesystemOperation{
		Name: "open-file", Path: path, Flags: flags, Mode: mode,
	})
	return recording.opened, recording.openErr
}

func (recording *recordingFilesystem) Remove(path string) error {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "remove", Path: path})
	return recording.removeErr
}

func (*recordingFilesystem) Create(string) (File, error) {
	return nil, errors.New("unexpected create")
}

func (*recordingFilesystem) Copy(File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected copy")
}

func TestDependenciesDefaultToSafeNoFilesystemMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist.txt")

	if err := MkdirAll(Dependencies{}, filepath.Dir(target), 0755); !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value directory dependency returned %v, want %v", err, ErrNoFilesystem)
	}
	if err := WriteFile(Dependencies{}, target, []byte("must not be written"), 0644); !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value write dependency returned %v, want %v", err, ErrNoFilesystem)
	}
	if file, err := OpenFile(Dependencies{}, target, os.O_APPEND|os.O_WRONLY, 0644); file != nil || !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value append-open dependency returned (%p, %v), want (nil, %v)", file, err, ErrNoFilesystem)
	}
	if err := Remove(Dependencies{}, target); !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value remove dependency returned %v, want %v", err, ErrNoFilesystem)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("zero-value filesystem dependencies mutated the target: %v", err)
	}
}

func TestOperationsPassCompleteValuesToDependency(t *testing.T) {
	data := []byte("complete source bytes\nwith a second line\n")
	mode := fs.FileMode(0751)
	info := staticFileInfo{name: "complete-source.txt", size: int64(len(data)), mode: mode}
	recording := &recordingFilesystem{readData: data, fileInfo: info}
	dependencies := Dependencies{FileSystem: recording}
	source := "/complete source/path with spaces/source.txt"
	directory := "/complete destination/path with spaces"
	destination := directory + "/destination.txt"

	gotData, err := ReadFile(dependencies, source)
	if err != nil {
		t.Fatalf("read dependency returned an error: %v", err)
	}
	gotInfo, err := Stat(dependencies, source)
	if err != nil {
		t.Fatalf("stat dependency returned an error: %v", err)
	}
	if err := MkdirAll(dependencies, directory, 0755); err != nil {
		t.Fatalf("mkdir dependency returned an error: %v", err)
	}
	if err := WriteFile(dependencies, destination, data, mode); err != nil {
		t.Fatalf("write dependency returned an error: %v", err)
	}

	if !reflect.DeepEqual(gotData, data) || gotInfo != info {
		t.Fatalf("dependency results were incomplete: data=%q info=%#v", gotData, gotInfo)
	}
	want := []recordedFilesystemOperation{
		{Name: "read", Path: source},
		{Name: "stat", Path: source},
		{Name: "mkdir-all", Path: directory, Mode: 0755},
		{Name: "write", Path: destination, Data: data, Mode: mode},
	}
	if !reflect.DeepEqual(recording.operations, want) {
		t.Fatalf("filesystem dependency received incomplete operations:\n got: %#v\nwant: %#v", recording.operations, want)
	}
}

func TestOpenFilePassesCompleteValuesAndReturnsExactDependencyResult(t *testing.T) {
	path := "/complete append-open/path with spaces/file.txt"
	flags := os.O_APPEND | os.O_WRONLY
	mode := fs.FileMode(0644)
	opened := &os.File{}
	openError := errors.New("complete append-open dependency error")
	recording := &recordingFilesystem{opened: opened, openErr: openError}
	dependencies := Dependencies{FileSystem: recording}

	actual, err := OpenFile(dependencies, path, flags, mode)

	if actual != opened || err != openError {
		t.Fatalf("append-open result was (%p, %v), want exact dependency result (%p, %v)", actual, err, opened, openError)
	}
	want := []recordedFilesystemOperation{{
		Name: "open-file", Path: path, Flags: flags, Mode: mode,
	}}
	if !reflect.DeepEqual(recording.operations, want) {
		t.Fatalf("append-open dependency received incomplete values:\n got: %#v\nwant: %#v", recording.operations, want)
	}
}

func TestRemovePassesCompletePathAndReturnsExactDependencyError(t *testing.T) {
	path := "/complete single-file-delete/path with spaces/file.txt"
	removeError := errors.New("complete remove dependency error")
	recording := &recordingFilesystem{removeErr: removeError}
	dependencies := Dependencies{FileSystem: recording}

	err := Remove(dependencies, path)

	if err != removeError {
		t.Fatalf("remove error was %v, want exact dependency error %v", err, removeError)
	}
	want := []recordedFilesystemOperation{{Name: "remove", Path: path}}
	if !reflect.DeepEqual(recording.operations, want) {
		t.Fatalf("remove dependency received incomplete values:\n got: %#v\nwant: %#v", recording.operations, want)
	}
}

func TestOperationsReturnDependencyErrors(t *testing.T) {
	sentinel := errors.New("complete filesystem dependency error")
	tests := []struct {
		name      string
		recording *recordingFilesystem
		invoke    func(Dependencies) error
	}{
		{
			name:      "read",
			recording: &recordingFilesystem{readErr: sentinel},
			invoke: func(dependencies Dependencies) error {
				_, err := ReadFile(dependencies, "/source")
				return err
			},
		},
		{
			name:      "stat",
			recording: &recordingFilesystem{statErr: sentinel},
			invoke: func(dependencies Dependencies) error {
				_, err := Stat(dependencies, "/source")
				return err
			},
		},
		{
			name:      "mkdir",
			recording: &recordingFilesystem{mkdirErr: sentinel},
			invoke: func(dependencies Dependencies) error {
				return MkdirAll(dependencies, "/destination", 0755)
			},
		},
		{
			name:      "write",
			recording: &recordingFilesystem{writeErr: sentinel},
			invoke: func(dependencies Dependencies) error {
				return WriteFile(dependencies, "/destination/file", []byte("bytes"), 0640)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.invoke(Dependencies{FileSystem: test.recording})
			if !errors.Is(err, sentinel) {
				t.Fatalf("dependency error was %v, want %v", err, sentinel)
			}
		})
	}
}

func TestSystemPerformsCompleteFilesystemOperations(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source.txt")
	destinationDirectory := filepath.Join(root, "missing", "nested")
	destination := filepath.Join(destinationDirectory, "destination.txt")
	data := []byte("system source bytes\n")
	if err := testutil.WriteFileOutsideWorkingTree(source, data, 0751); err != nil {
		t.Fatalf("create source fixture: %v", err)
	}

	readData, err := ReadFile(System(), source)
	if err != nil {
		t.Fatalf("system read returned an error: %v", err)
	}
	info, err := Stat(System(), source)
	if err != nil {
		t.Fatalf("system stat returned an error: %v", err)
	}
	if err := MkdirAll(System(), destinationDirectory, 0755); err != nil {
		t.Fatalf("system mkdir returned an error: %v", err)
	}
	if err := WriteFile(System(), destination, readData, info.Mode()); err != nil {
		t.Fatalf("system write returned an error: %v", err)
	}

	written, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read system destination: %v", err)
	}
	writtenInfo, err := os.Stat(destination)
	if err != nil {
		t.Fatalf("stat system destination: %v", err)
	}
	if !reflect.DeepEqual(written, data) || writtenInfo.Mode().Perm() != info.Mode().Perm() {
		t.Fatalf("system destination data/mode = (%q, %s), want (%q, %s)", written, writtenInfo.Mode().Perm(), data, info.Mode().Perm())
	}
	if err := Remove(System(), destination); err != nil {
		t.Fatalf("system remove returned an error: %v", err)
	}
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatalf("system remove left the destination present: %v", err)
	}
}

type staticFileInfo struct {
	name string
	size int64
	mode fs.FileMode
}

func (info staticFileInfo) Name() string       { return info.name }
func (info staticFileInfo) Size() int64        { return info.size }
func (info staticFileInfo) Mode() fs.FileMode  { return info.mode }
func (info staticFileInfo) ModTime() time.Time { return time.Time{} }
func (info staticFileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info staticFileInfo) Sys() interface{}   { return nil }
