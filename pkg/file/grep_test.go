package file

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
)

type recordedGrepRecursiveCallback struct {
	path string
	info fs.FileInfo
	err  error
}

type recordingGrepRecursiveFilesystem struct {
	walkRoots       []string
	walkInputs      []recordedGrepRecursiveCallback
	callbackInputs  []recordedGrepRecursiveCallback
	callbackResults []error
	walkErr         error
}

func (*recordingGrepRecursiveFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected recursive-grep read")
}

func (*recordingGrepRecursiveFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected recursive-grep read directory")
}

func (*recordingGrepRecursiveFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected recursive-grep read directory entries")
}

func (*recordingGrepRecursiveFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected recursive-grep stat")
}

func (*recordingGrepRecursiveFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected recursive-grep mkdir")
}

func (*recordingGrepRecursiveFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected recursive-grep mkdir")
}

func (*recordingGrepRecursiveFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected recursive-grep write")
}

func (*recordingGrepRecursiveFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected recursive-grep open file")
}

func (*recordingGrepRecursiveFilesystem) Remove(string) error {
	return errors.New("unexpected recursive-grep remove")
}

func (*recordingGrepRecursiveFilesystem) RemoveAll(string) error {
	return errors.New("unexpected recursive-grep recursive remove")
}

func (*recordingGrepRecursiveFilesystem) Rename(string, string) error {
	return errors.New("unexpected recursive-grep rename")
}

func (*recordingGrepRecursiveFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected recursive-grep glob")
}

func (recording *recordingGrepRecursiveFilesystem) Walk(root string, callback filepath.WalkFunc) error {
	recording.walkRoots = append(recording.walkRoots, root)
	for _, input := range recording.walkInputs {
		recording.callbackInputs = append(recording.callbackInputs, input)
		callbackResult := callback(input.path, input.info, input.err)
		recording.callbackResults = append(recording.callbackResults, callbackResult)
		if callbackResult != nil {
			return callbackResult
		}
	}
	return recording.walkErr
}

func (*recordingGrepRecursiveFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected recursive-grep create")
}

func (*recordingGrepRecursiveFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected recursive-grep copy")
}

func (recording *recordingGrepRecursiveFilesystem) dependencies() grepRecursiveDependencies {
	return grepRecursiveDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingGrepRecursiveFilesystem) assertedWalkRoots() ([]string, error) {
	if len(recording.walkRoots) == 0 {
		return nil, errors.New("recorded recursive-grep walk population is empty")
	}
	return recording.walkRoots, nil
}

func (recording *recordingGrepRecursiveFilesystem) assertedCallbackInputs() ([]recordedGrepRecursiveCallback, error) {
	if len(recording.callbackInputs) == 0 {
		return nil, errors.New("recorded recursive-grep callback population is empty")
	}
	return recording.callbackInputs, nil
}

func TestGrepRecursiveSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemGrepRecursiveDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("GrepRecursive selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("GrepRecursive filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestGrepRecursivePreservesExactRootCallbackOrderPerPathGrepAndOrderedAllHits(t *testing.T) {
	root := "/complete recursive-grep/root path with spaces"
	callbackInputError := errors.New("ignored recursive-grep callback input error")
	firstMatch := filepath.Join("test", "Dockerfile.render")
	readFailure := "/developer/home/project/ordered-missing-path-must-not-exist"
	secondMatch := filepath.Join("test", "Dockerfile")
	recording := &recordingGrepRecursiveFilesystem{walkInputs: []recordedGrepRecursiveCallback{
		{path: firstMatch, info: existingFileInfo{name: "ignored directory metadata"}, err: callbackInputError},
		{path: readFailure, info: nil, err: callbackInputError},
		{path: secondMatch, info: existingFileInfo{name: "ignored file metadata"}, err: callbackInputError},
		{path: firstMatch, info: nil, err: callbackInputError},
	}}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("recursive-grep dependency lost its complete filesystem value: %#v", dependencies)
	}

	files, err := grepRecursive(dependencies, root, "LABEL")

	wantFiles := []string{firstMatch, secondMatch, firstMatch}
	if err != nil || !reflect.DeepEqual(files, wantFiles) {
		t.Fatalf("recursive-grep result was (%#v, %v), want ordered all-hit result (%#v, nil)", files, err, wantFiles)
	}
	roots, populationErr := recording.assertedWalkRoots()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(roots, []string{root}) {
		t.Fatalf("recursive-grep walk roots were %#v, want exact root %#v", roots, []string{root})
	}
	callbacks, populationErr := recording.assertedCallbackInputs()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(callbacks, recording.walkInputs) {
		t.Fatalf("recursive-grep callback order was %#v, want %#v", callbacks, recording.walkInputs)
	}
	wantCallbackResults := make([]error, len(recording.walkInputs))
	if !reflect.DeepEqual(recording.callbackResults, wantCallbackResults) {
		t.Fatalf("recursive-grep callback results were %#v, want all nil", recording.callbackResults)
	}
}

func TestGrepRecursiveNoMatchUsesCaseSensitiveKeywordAndSuppressesPerPathReadFailure(t *testing.T) {
	root := "/complete recursive-grep/no match"
	recording := &recordingGrepRecursiveFilesystem{walkInputs: []recordedGrepRecursiveCallback{
		{path: filepath.Join("test", "Dockerfile.render"), info: existingFileInfo{name: "Dockerfile.render"}},
		{path: "/developer/home/project/missing-path-must-not-exist", err: errors.New("ignored missing-path callback error")},
		{path: filepath.Join("test", "Dockerfile"), info: existingFileInfo{name: "Dockerfile"}},
	}}

	files, err := grepRecursive(recording.dependencies(), root, "label")

	if files != nil || err != nil {
		t.Fatalf("no-match recursive-grep result was (%#v, %v), want (nil, nil)", files, err)
	}
	callbacks, populationErr := recording.assertedCallbackInputs()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(callbacks, recording.walkInputs) {
		t.Fatalf("no-match recursive-grep callback order was %#v, want %#v", callbacks, recording.walkInputs)
	}
	wantCallbackResults := make([]error, len(recording.walkInputs))
	if !reflect.DeepEqual(recording.callbackResults, wantCallbackResults) {
		t.Fatalf("suppressed-read-failure callback results were %#v, want all nil", recording.callbackResults)
	}
}

func TestGrepRecursiveReturnsExactNonEOFWalkErrorWithCurrentPartialResult(t *testing.T) {
	root := "/complete recursive-grep/walk error"
	matchingPath := filepath.Join("test", "Dockerfile")
	walkErrors := []struct {
		name string
		err  error
	}{
		{name: "ordinary error", err: errors.New("complete recursive-grep non-EOF walk error")},
		{name: "wrapped EOF is not equal to EOF", err: &os.PathError{Op: "walk", Path: root, Err: io.EOF}},
	}

	for _, walkError := range walkErrors {
		t.Run(walkError.name, func(t *testing.T) {
			recording := &recordingGrepRecursiveFilesystem{
				walkInputs: []recordedGrepRecursiveCallback{
					{path: matchingPath, err: errors.New("ignored matching callback error")},
					{path: "/developer/home/project/missing-after-partial-result"},
				},
				walkErr: walkError.err,
			}

			files, err := grepRecursive(recording.dependencies(), root, "LABEL")

			wantFiles := []string{matchingPath}
			if err != walkError.err || !reflect.DeepEqual(files, wantFiles) {
				t.Fatalf("walk-error recursive-grep result was (%#v, %v), want current partial result and exact error (%#v, %v)",
					files, err, wantFiles, walkError.err)
			}
			if !reflect.DeepEqual(recording.callbackResults, []error{nil, nil}) {
				t.Fatalf("walk-error callback results were %#v, want all nil", recording.callbackResults)
			}
		})
	}
}

func TestGrepRecursiveNormalizesOnlyFinalEOFWithCurrentPartialResult(t *testing.T) {
	root := "/complete recursive-grep/final EOF"
	matchingPath := filepath.Join("test", "Dockerfile.render")
	recording := &recordingGrepRecursiveFilesystem{
		walkInputs: []recordedGrepRecursiveCallback{{
			path: matchingPath,
			info: existingFileInfo{name: "ignored metadata"},
			err:  errors.New("ignored callback error"),
		}},
		walkErr: io.EOF,
	}

	files, err := grepRecursive(recording.dependencies(), root, "LABEL")

	if err != nil || !reflect.DeepEqual(files, []string{matchingPath}) {
		t.Fatalf("EOF recursive-grep result was (%#v, %v), want current partial result and nil", files, err)
	}
	if !reflect.DeepEqual(recording.callbackResults, []error{nil}) {
		t.Fatalf("EOF callback results were %#v, want nil", recording.callbackResults)
	}
}

func TestGrepRecursiveDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	files, err := grepRecursive(
		grepRecursiveDependencies{},
		"/developer/home/project/must-not-be-accessed",
		"LABEL",
	)

	if files != nil || !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe recursive-grep dependency default returned (%#v, %v), want (nil, %v)", files, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedGrepRecursiveRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingGrepRecursiveFilesystem{}

	if _, err := recording.assertedWalkRoots(); err == nil {
		t.Fatal("empty recorded recursive-grep walk population passed")
	}
	if _, err := recording.assertedCallbackInputs(); err == nil {
		t.Fatal("empty recorded recursive-grep callback population passed")
	}
}
