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
	Name        string
	Path        string
	Destination string
	Data        []byte
	Flags       int
	Mode        fs.FileMode
}

type recordedWalkInput struct {
	Path string
	Info fs.FileInfo
	Err  error
}

type recordingFilesystem struct {
	operations          []recordedFilesystemOperation
	globMatches         []string
	readData            []byte
	readDirEntries      []fs.FileInfo
	fileInfo            fs.FileInfo
	walkInputs          []recordedWalkInput
	walkCallbackResults []error
	globErr             error
	readErr             error
	readDirErr          error
	statErr             error
	mkdirErr            error
	writeErr            error
	opened              *os.File
	openErr             error
	removeErr           error
	removeAllErr        error
	renameErr           error
	walkErr             error
}

func (recording *recordingFilesystem) ReadFile(path string) ([]byte, error) {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "read", Path: path})
	return append([]byte(nil), recording.readData...), recording.readErr
}

func (recording *recordingFilesystem) ReadDir(path string) ([]fs.FileInfo, error) {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "read-dir", Path: path})
	return append([]fs.FileInfo{}, recording.readDirEntries...), recording.readDirErr
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

func (recording *recordingFilesystem) RemoveAll(path string) error {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "remove-all", Path: path})
	return recording.removeAllErr
}

func (recording *recordingFilesystem) Rename(source, destination string) error {
	recording.operations = append(recording.operations, recordedFilesystemOperation{
		Name: "rename", Path: source, Destination: destination,
	})
	return recording.renameErr
}

func (recording *recordingFilesystem) Glob(pattern string) ([]string, error) {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "glob", Path: pattern})
	return append([]string{}, recording.globMatches...), recording.globErr
}

func (recording *recordingFilesystem) Walk(root string, callback filepath.WalkFunc) error {
	recording.operations = append(recording.operations, recordedFilesystemOperation{Name: "walk", Path: root})
	for _, input := range recording.walkInputs {
		recording.walkCallbackResults = append(recording.walkCallbackResults, callback(input.Path, input.Info, input.Err))
	}
	return recording.walkErr
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
	if err := RemoveAll(Dependencies{}, target); !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value recursive-remove dependency returned %v, want %v", err, ErrNoFilesystem)
	}
	if err := Rename(Dependencies{}, target, target+".moved"); !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value rename dependency returned %v, want %v", err, ErrNoFilesystem)
	}
	if matches, err := Glob(Dependencies{}, filepath.Join(target, "*")); matches != nil || !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value glob dependency returned (%#v, %v), want (nil, %v)", matches, err, ErrNoFilesystem)
	}
	if entries, err := ReadDir(Dependencies{}, target); entries != nil || !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value directory-read dependency returned (%#v, %v), want (nil, %v)", entries, err, ErrNoFilesystem)
	}
	callbackCalls := 0
	if err := Walk(Dependencies{}, target, func(string, fs.FileInfo, error) error {
		callbackCalls++
		return nil
	}); !errors.Is(err, ErrNoFilesystem) {
		t.Fatalf("zero-value walk dependency returned %v, want %v", err, ErrNoFilesystem)
	}
	if callbackCalls != 0 {
		t.Fatalf("zero-value walk dependency invoked its callback %d times", callbackCalls)
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

func TestRemoveAllPassesCompletePathAndReturnsExactDependencyError(t *testing.T) {
	path := "/complete recursive-delete/path with spaces"
	removeError := errors.New("complete recursive-remove dependency error")
	recording := &recordingFilesystem{removeAllErr: removeError}
	dependencies := Dependencies{FileSystem: recording}

	err := RemoveAll(dependencies, path)

	if err != removeError {
		t.Fatalf("recursive-remove error was %v, want exact dependency error %v", err, removeError)
	}
	want := []recordedFilesystemOperation{{Name: "remove-all", Path: path}}
	if !reflect.DeepEqual(recording.operations, want) {
		t.Fatalf("recursive-remove dependency received incomplete values:\n got: %#v\nwant: %#v", recording.operations, want)
	}
}

func TestRenamePassesCompletePathsAndReturnsExactDependencyError(t *testing.T) {
	source := "/complete file-move/source path with spaces/source.txt"
	destination := "/complete file-move/destination path with spaces/destination.txt"
	renameError := errors.New("complete rename dependency error")
	recording := &recordingFilesystem{renameErr: renameError}
	dependencies := Dependencies{FileSystem: recording}

	err := Rename(dependencies, source, destination)

	if err != renameError {
		t.Fatalf("rename error was %v, want exact dependency error %v", err, renameError)
	}
	want := []recordedFilesystemOperation{{Name: "rename", Path: source, Destination: destination}}
	if !reflect.DeepEqual(recording.operations, want) {
		t.Fatalf("rename dependency received incomplete values:\n got: %#v\nwant: %#v", recording.operations, want)
	}
}

func TestGlobPassesExactPatternOrderedMatchesAndExactDependencyError(t *testing.T) {
	pattern := "/complete glob/path with spaces/*"
	matches := []string{
		"/complete glob/path with spaces/z-result",
		"/complete glob/path with spaces/a-result",
	}
	globError := errors.New("complete glob dependency error")
	recording := &recordingFilesystem{globMatches: matches, globErr: globError}
	dependencies := Dependencies{FileSystem: recording}

	actual, err := Glob(dependencies, pattern)

	if err != globError {
		t.Fatalf("glob error was %v, want exact dependency error %v", err, globError)
	}
	if !reflect.DeepEqual(actual, matches) {
		t.Fatalf("glob matches were %#v, want ordered %#v", actual, matches)
	}
	want := []recordedFilesystemOperation{{Name: "glob", Path: pattern}}
	if !reflect.DeepEqual(recording.operations, want) {
		t.Fatalf("glob dependency received incomplete values:\n got: %#v\nwant: %#v", recording.operations, want)
	}
}

func TestReadDirPassesExactPathOrderedEntriesAndExactDependencyError(t *testing.T) {
	path := "/complete directory-read/path with spaces"
	entries := []fs.FileInfo{
		staticFileInfo{name: "z-first.iml", size: 41, mode: 0751},
		staticFileInfo{name: "a-second.idea", size: 83, mode: fs.ModeDir | 0705},
	}
	readError := errors.New("complete directory-read dependency error")
	recording := &recordingFilesystem{readDirEntries: entries, readDirErr: readError}
	dependencies := Dependencies{FileSystem: recording}

	actual, err := ReadDir(dependencies, path)

	if err != readError {
		t.Fatalf("directory-read error was %v, want exact dependency error %v", err, readError)
	}
	if !reflect.DeepEqual(actual, entries) {
		t.Fatalf("directory-read entries were %#v, want ordered metadata %#v", actual, entries)
	}
	want := []recordedFilesystemOperation{{Name: "read-dir", Path: path}}
	if !reflect.DeepEqual(recording.operations, want) {
		t.Fatalf("directory-read dependency received incomplete values:\n got: %#v\nwant: %#v", recording.operations, want)
	}
}

func TestSystemReadDirPreservesIoutilFilenameSortingAndMetadata(t *testing.T) {
	root := t.TempDir()
	last := filepath.Join(root, "z-last.iml")
	first := filepath.Join(root, "a-first.txt")
	directory := filepath.Join(root, "m-directory.idea")
	if err := testutil.WriteFileOutsideWorkingTree(last, []byte("last"), 0641); err != nil {
		t.Fatalf("create final directory-read fixture: %v", err)
	}
	if err := testutil.WriteFileOutsideWorkingTree(first, []byte("first bytes"), 0751); err != nil {
		t.Fatalf("create first directory-read fixture: %v", err)
	}
	if err := MkdirAll(System(), directory, 0705); err != nil {
		t.Fatalf("create directory-read fixture directory: %v", err)
	}

	entries, err := ReadDir(System(), root)

	if err != nil {
		t.Fatalf("system directory-read returned an error: %v", err)
	}
	wantNames := []string{"a-first.txt", "m-directory.idea", "z-last.iml"}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	if !reflect.DeepEqual(names, wantNames) {
		t.Fatalf("system directory-read names were %#v, want filename-sorted %#v", names, wantNames)
	}
	if entries[0].Size() != int64(len("first bytes")) || entries[0].IsDir() || !entries[1].IsDir() || entries[2].Size() != int64(len("last")) {
		t.Fatalf("system directory-read metadata was %#v", entries)
	}
}

func TestWalkPassesExactRootAndCallbackAndReturnsExactDependencyError(t *testing.T) {
	root := "/complete walk/root path with spaces"
	callbackPath := root + "/callback target.kt"
	callbackInfo := staticFileInfo{name: "callback target.kt", size: 41, mode: 0751}
	callbackInputError := errors.New("complete walk callback input error")
	callbackReturnError := errors.New("complete walk callback return error")
	walkError := errors.New("complete walk dependency error")
	recording := &recordingFilesystem{
		walkInputs: []recordedWalkInput{{Path: callbackPath, Info: callbackInfo, Err: callbackInputError}},
		walkErr:    walkError,
	}
	dependencies := Dependencies{FileSystem: recording}
	var callbackInputs []recordedWalkInput

	err := Walk(dependencies, root, func(path string, info fs.FileInfo, err error) error {
		callbackInputs = append(callbackInputs, recordedWalkInput{Path: path, Info: info, Err: err})
		return callbackReturnError
	})

	if err != walkError {
		t.Fatalf("walk error was %v, want exact dependency error %v", err, walkError)
	}
	wantOperations := []recordedFilesystemOperation{{Name: "walk", Path: root}}
	if !reflect.DeepEqual(recording.operations, wantOperations) {
		t.Fatalf("walk dependency received incomplete values:\n got: %#v\nwant: %#v", recording.operations, wantOperations)
	}
	if !reflect.DeepEqual(callbackInputs, recording.walkInputs) {
		t.Fatalf("walk callback received incomplete values:\n got: %#v\nwant: %#v", callbackInputs, recording.walkInputs)
	}
	if !reflect.DeepEqual(recording.walkCallbackResults, []error{callbackReturnError}) {
		t.Fatalf("walk dependency observed callback results %#v, want exact error %#v", recording.walkCallbackResults, []error{callbackReturnError})
	}
}

func TestSystemWalkDelegatesRootAndPreservesFilepathTraversalOrder(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "a-directory")
	nested := filepath.Join(directory, "nested.txt")
	last := filepath.Join(root, "z-last.txt")
	if err := MkdirAll(System(), directory, 0755); err != nil {
		t.Fatalf("create system walk fixture directory: %v", err)
	}
	if err := testutil.WriteFileOutsideWorkingTree(nested, []byte("nested"), 0644); err != nil {
		t.Fatalf("create nested system walk fixture: %v", err)
	}
	if err := testutil.WriteFileOutsideWorkingTree(last, []byte("last"), 0644); err != nil {
		t.Fatalf("create final system walk fixture: %v", err)
	}
	var paths []string

	err := Walk(System(), root, func(path string, info fs.FileInfo, err error) error {
		if info == nil || err != nil {
			t.Fatalf("system walk callback received (%#v, %v) for %s", info, err, path)
		}
		paths = append(paths, path)
		return nil
	})

	if err != nil {
		t.Fatalf("system walk returned an error: %v", err)
	}
	want := []string{root, directory, nested, last}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("system walk paths were %#v, want filepath order %#v", paths, want)
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

func TestSystemRemoveAllRemovesTemporaryDirectoryTree(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "recursive remove", "nested")
	file := filepath.Join(target, "file.txt")
	if err := MkdirAll(System(), target, 0755); err != nil {
		t.Fatalf("create recursive-remove fixture directory: %v", err)
	}
	if err := WriteFile(System(), file, []byte("temporary fixture"), 0644); err != nil {
		t.Fatalf("create recursive-remove fixture file: %v", err)
	}

	if err := RemoveAll(System(), filepath.Dir(target)); err != nil {
		t.Fatalf("system recursive remove returned an error: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(target)); !os.IsNotExist(err) {
		t.Fatalf("system recursive remove left the directory tree present: %v", err)
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
