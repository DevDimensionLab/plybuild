package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/sirupsen/logrus"
)

type recordedCopyOperation struct {
	Name string
	Path string
	Data []byte
	Mode fs.FileMode
}

type recordedStatResult struct {
	Path string
	Info fs.FileInfo
	Err  error
}

type recordingCopyFilesystem struct {
	operations  []recordedCopyOperation
	statResults []recordedStatResult
	statIndex   int
	readData    []byte
	readErr     error
	mkdirErr    error
	writeErr    error
}

func (recording *recordingCopyFilesystem) dependencies() filesystem.Dependencies {
	return filesystem.Dependencies{FileSystem: recording}
}

func (recording *recordingCopyFilesystem) ReadFile(path string) ([]byte, error) {
	recording.operations = append(recording.operations, recordedCopyOperation{Name: "read", Path: path})
	return append([]byte(nil), recording.readData...), recording.readErr
}

func (recording *recordingCopyFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.operations = append(recording.operations, recordedCopyOperation{Name: "stat", Path: path})
	if recording.statIndex >= len(recording.statResults) {
		return nil, fmt.Errorf("unexpected stat operation for %s", path)
	}
	result := recording.statResults[recording.statIndex]
	recording.statIndex++
	if result.Path != path {
		return nil, fmt.Errorf("stat path was %s, want %s", path, result.Path)
	}
	return result.Info, result.Err
}

func (recording *recordingCopyFilesystem) MkdirAll(path string, mode fs.FileMode) error {
	recording.operations = append(recording.operations, recordedCopyOperation{Name: "mkdir-all", Path: path, Mode: mode})
	return recording.mkdirErr
}

func (recording *recordingCopyFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.operations = append(recording.operations, recordedCopyOperation{
		Name: "write", Path: path, Data: append([]byte(nil), data...), Mode: mode,
	})
	return recording.writeErr
}

func (*recordingCopyFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected open file")
}

func (*recordingCopyFilesystem) Remove(string) error {
	return errors.New("unexpected remove")
}

func (*recordingCopyFilesystem) RemoveAll(string) error {
	return errors.New("unexpected recursive remove")
}

func (*recordingCopyFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected create")
}

func (*recordingCopyFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected copy")
}

func (recording *recordingCopyFilesystem) assertedOperations() ([]recordedCopyOperation, error) {
	if len(recording.operations) == 0 {
		return nil, errors.New("recorded template file copy population is empty")
	}
	return recording.operations, nil
}

func TestCopyOrMergeMissingTargetSelectsCopyWithCompleteDependency(t *testing.T) {
	source := "/complete template/source path/Complete.kt"
	destination := "/complete project/resolved target path/Complete.kt"
	destinationDirectory := filepath.Dir(destination)
	data := []byte("complete source bytes\nwith a second line\n")
	mode := fs.FileMode(0751)
	recording := &recordingCopyFilesystem{
		readData: data,
		statResults: []recordedStatResult{
			{Path: destination, Err: fs.ErrNotExist},
			{Path: destinationDirectory, Err: fs.ErrNotExist},
			{Path: destinationDirectory, Err: fs.ErrNotExist},
			{Path: source, Info: copyFileInfo{name: "Complete.kt", size: int64(len(data)), mode: mode}},
		},
	}
	logOutput := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(logOutput)
	previousLogger := log
	log = logger
	t.Cleanup(func() { log = previousLogger })

	err := copyOrMerge(recording.dependencies(), source, destination)

	if err != nil {
		t.Fatalf("copy-or-merge returned an error: %v", err)
	}
	want := []recordedCopyOperation{
		{Name: "stat", Path: destination},
		{Name: "read", Path: source},
		{Name: "stat", Path: destinationDirectory},
		{Name: "stat", Path: destinationDirectory},
		{Name: "mkdir-all", Path: destinationDirectory, Mode: 0755},
		{Name: "stat", Path: source},
		{Name: "write", Path: destination, Data: data, Mode: mode},
	}
	assertRecordedCopyOperations(t, recording, want)
	fromIndex := strings.Index(logOutput.String(), source)
	toIndex := strings.Index(logOutput.String(), destination)
	if fromIndex < 0 || toIndex < 0 || fromIndex >= toIndex {
		t.Fatalf("copy logging lost source-before-destination order:\n%s", logOutput.String())
	}
}

func TestCopyOrMergeExistingTargetRetainsMergeSelection(t *testing.T) {
	tests := []struct {
		name    string
		statErr error
	}{
		{name: "existing"},
		{name: "non-missing stat error remains existing", statErr: errors.New("destination stat failed")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := "/complete template/Existing.kt"
			destination := "/complete target/Existing.kt"
			recording := &recordingCopyFilesystem{statResults: []recordedStatResult{{
				Path: destination,
				Info: copyFileInfo{name: "Existing.kt", mode: 0644},
				Err:  test.statErr,
			}}}

			if err := copyOrMerge(recording.dependencies(), source, destination); err != nil {
				t.Fatalf("existing-target merge selection returned an error: %v", err)
			}
			assertRecordedCopyOperations(t, recording, []recordedCopyOperation{{Name: "stat", Path: destination}})
		})
	}
}

func TestCopyFileReturnsDependencyErrorsInLegacyOrder(t *testing.T) {
	source := "/complete template/error-source.txt"
	destination := "/complete target/missing/error-destination.txt"
	directory := filepath.Dir(destination)
	data := []byte("complete bytes")
	mode := fs.FileMode(0740)

	tests := []struct {
		name      string
		recording *recordingCopyFilesystem
		want      []recordedCopyOperation
		wantError string
	}{
		{
			name: "source read first",
			recording: &recordingCopyFilesystem{
				readErr:     errors.New("source read failed"),
				statResults: []recordedStatResult{{Path: destination, Err: fs.ErrNotExist}},
			},
			want:      []recordedCopyOperation{{Name: "stat", Path: destination}, {Name: "read", Path: source}},
			wantError: "source read failed",
		},
		{
			name: "directory creation before source mode",
			recording: &recordingCopyFilesystem{
				readData: data,
				mkdirErr: errors.New("directory creation failed"),
				statResults: []recordedStatResult{
					{Path: destination, Err: fs.ErrNotExist},
					{Path: directory, Err: fs.ErrNotExist},
					{Path: directory, Err: fs.ErrNotExist},
				},
			},
			want: []recordedCopyOperation{
				{Name: "stat", Path: destination}, {Name: "read", Path: source},
				{Name: "stat", Path: directory}, {Name: "stat", Path: directory},
				{Name: "mkdir-all", Path: directory, Mode: 0755},
			},
			// CopyFile historically returns the preceding missing-directory stat
			// error when MkdirAll fails; preserve that observable error behavior.
			wantError: fs.ErrNotExist.Error(),
		},
		{
			name: "source mode before destination write",
			recording: &recordingCopyFilesystem{
				readData: data,
				statResults: []recordedStatResult{
					{Path: destination, Err: fs.ErrNotExist},
					{Path: directory, Info: copyFileInfo{name: filepath.Base(directory), mode: fs.ModeDir | 0755}},
					{Path: source, Err: errors.New("source mode failed")},
				},
			},
			want: []recordedCopyOperation{
				{Name: "stat", Path: destination}, {Name: "read", Path: source},
				{Name: "stat", Path: directory}, {Name: "stat", Path: source},
			},
			wantError: "source mode failed",
		},
		{
			name: "destination write last",
			recording: &recordingCopyFilesystem{
				readData: data,
				writeErr: errors.New("destination write failed"),
				statResults: []recordedStatResult{
					{Path: destination, Err: fs.ErrNotExist},
					{Path: directory, Info: copyFileInfo{name: filepath.Base(directory), mode: fs.ModeDir | 0755}},
					{Path: source, Info: copyFileInfo{name: filepath.Base(source), size: int64(len(data)), mode: mode}},
				},
			},
			want: []recordedCopyOperation{
				{Name: "stat", Path: destination}, {Name: "read", Path: source},
				{Name: "stat", Path: directory}, {Name: "stat", Path: source},
				{Name: "write", Path: destination, Data: data, Mode: mode},
			},
			wantError: "destination write failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := copyOrMerge(test.recording.dependencies(), source, destination)
			if err == nil || err.Error() != test.wantError {
				t.Fatalf("dependency error was %v, want %q", err, test.wantError)
			}
			assertRecordedCopyOperations(t, test.recording, test.want)
		})
	}
}

func TestCopyDependenciesDefaultToNoMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe-default", "must-not-exist.txt")

	err := copyOrMerge(filesystem.Dependencies{}, "/must-not-read/source.txt", target)

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe copy dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe copy dependency default mutated the target: %v", statErr)
	}
}

func TestRecordedCopyOperationsRejectEmptyPopulation(t *testing.T) {
	recording := &recordingCopyFilesystem{}

	if _, err := recording.assertedOperations(); err == nil {
		t.Fatal("empty recorded template file copy population passed")
	}
}

func TestCopyFilePreservesSourceBytesModeAndCreatesDestinationDirectory(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source file.txt")
	destination := filepath.Join(root, "missing", "nested", "destination file.txt")
	data := []byte("source bytes preserved exactly\n")
	if err := testutil.WriteFileOutsideWorkingTree(source, data, 0751); err != nil {
		t.Fatalf("create copy source fixture: %v", err)
	}

	if err := CopyFile(source, destination); err != nil {
		t.Fatalf("copy file returned an error: %v", err)
	}

	written, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read copied destination: %v", err)
	}
	info, err := os.Stat(destination)
	if err != nil {
		t.Fatalf("stat copied destination: %v", err)
	}
	if !reflect.DeepEqual(written, data) || info.Mode().Perm() != 0751 {
		t.Fatalf("copied destination data/mode = (%q, %s), want (%q, %s)", written, info.Mode().Perm(), data, fs.FileMode(0751))
	}
}

func assertRecordedCopyOperations(t *testing.T, recording *recordingCopyFilesystem, want []recordedCopyOperation) {
	t.Helper()
	operations, err := recording.assertedOperations()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(operations, want) {
		t.Fatalf("recorded copy operations differ:\n got: %#v\nwant: %#v", operations, want)
	}
}

type copyFileInfo struct {
	name string
	size int64
	mode fs.FileMode
}

func (info copyFileInfo) Name() string       { return info.name }
func (info copyFileInfo) Size() int64        { return info.size }
func (info copyFileInfo) Mode() fs.FileMode  { return info.mode }
func (info copyFileInfo) ModTime() time.Time { return time.Time{} }
func (info copyFileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info copyFileInfo) Sys() interface{}   { return nil }
