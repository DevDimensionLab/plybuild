package template

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
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/sirupsen/logrus"
)

func (*recordingFilteredWalkFilesystem) OpenZipEntry(*zip.File) (io.ReadCloser, error) {
	return nil, errors.New("unexpected filtered-walk archive entry open")
}

func (*recordingFilteredWalkFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected filtered-template archive open")
}

func (*recordingFilteredWalkFilesystem) Close(filesystem.File) error {
	return errors.New("unexpected filtered-template close")
}

func (*recordingFilteredWalkFilesystem) CloseReader(io.Closer) error {
	return errors.New("unexpected filtered-template reader close")
}

type recordedFilteredWalkCallback struct {
	path string
	info fs.FileInfo
	err  error
}

type recordingFilteredWalkFilesystem struct {
	walkRoots       []string
	walkInputs      []recordedFilteredWalkCallback
	callbackInputs  []recordedFilteredWalkCallback
	callbackResults []error
	walkErr         error
}

func (*recordingFilteredWalkFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected filtered-walk working directory")
}

func (*recordingFilteredWalkFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected filtered-template read")
}

func (*recordingFilteredWalkFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected filtered-template read directory")
}

func (*recordingFilteredWalkFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected filtered-template read directory entries")
}

func (*recordingFilteredWalkFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected filtered-template stat")
}

func (*recordingFilteredWalkFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected filtered-template mkdir")
}

func (*recordingFilteredWalkFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected filtered-template mkdir")
}

func (*recordingFilteredWalkFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected filtered-template write")
}

func (*recordingFilteredWalkFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected filtered-template open file")
}

func (*recordingFilteredWalkFilesystem) Remove(string) error {
	return errors.New("unexpected filtered-template remove")
}

func (*recordingFilteredWalkFilesystem) RemoveAll(string) error {
	return errors.New("unexpected filtered-template recursive remove")
}

func (*recordingFilteredWalkFilesystem) Rename(string, string) error {
	return errors.New("unexpected filtered-template rename")
}

func (*recordingFilteredWalkFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected filtered-template glob")
}

func (recording *recordingFilteredWalkFilesystem) Walk(root string, callback filepath.WalkFunc) error {
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

func (*recordingFilteredWalkFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected filtered-template create")
}

func (*recordingFilteredWalkFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected filtered-template copy")
}

func (recording *recordingFilteredWalkFilesystem) dependencies() filteredWalkDependencies {
	return filteredWalkDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFilteredWalkFilesystem) assertedWalkRoots() ([]string, error) {
	if len(recording.walkRoots) == 0 {
		return nil, errors.New("recorded filtered-template walk-root population is empty")
	}
	return recording.walkRoots, nil
}

func (recording *recordingFilteredWalkFilesystem) assertedCallbackInputs() ([]recordedFilteredWalkCallback, error) {
	if len(recording.callbackInputs) == 0 {
		return nil, errors.New("recorded filtered-template callback population is empty")
	}
	return recording.callbackInputs, nil
}

type filteredWalkFileInfo struct {
	name         string
	directory    bool
	size         int64
	mode         fs.FileMode
	modTime      time.Time
	system       interface{}
	nameCalls    int
	sizeCalls    int
	modeCalls    int
	modTimeCalls int
	isDirCalls   int
	systemCalls  int
}

func (info *filteredWalkFileInfo) Name() string {
	info.nameCalls++
	return info.name
}

func (info *filteredWalkFileInfo) Size() int64 {
	info.sizeCalls++
	return info.size
}

func (info *filteredWalkFileInfo) Mode() fs.FileMode {
	info.modeCalls++
	return info.mode
}

func (info *filteredWalkFileInfo) ModTime() time.Time {
	info.modTimeCalls++
	return info.modTime
}

func (info *filteredWalkFileInfo) IsDir() bool {
	info.isDirCalls++
	return info.directory
}

func (info *filteredWalkFileInfo) Sys() interface{} {
	info.systemCalls++
	return info.system
}

func (info *filteredWalkFileInfo) otherMetadataCalls() int {
	return info.sizeCalls + info.modeCalls + info.modTimeCalls + info.systemCalls
}

type filteredWalkMessageOnlyFormatter struct{}

func (filteredWalkMessageOnlyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func captureFilteredWalkDebugLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	output := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(output)
	logger.SetFormatter(filteredWalkMessageOnlyFormatter{})
	previousLogger := log
	log = logger
	t.Cleanup(func() { log = previousLogger })
	return output
}

func TestFilteredFilesFromTemplateSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemFilteredWalkDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("filtered template walk selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("filtered template walk filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestFilteredFilesFromTemplatePreservesExactRootCallbackOrderMetadataErrorsFilteringAndLogging(t *testing.T) {
	output := captureFilteredWalkDebugLog(t)
	sourceDir := "/complete filtered-template/root path with spaces"
	filter := []string{"pom.xml", "ignored-first", "ignored-second"}
	incomingError := errors.New("ignored filtered-template callback input error")
	directoryInfo := &filteredWalkFileInfo{name: "root path with spaces", directory: true, size: 41, mode: fs.ModeDir | 0751, modTime: time.Unix(1700000301, 0), system: "directory metadata"}
	rootPomInfo := &filteredWalkFileInfo{name: "pom.xml", size: 42, mode: 0601, modTime: time.Unix(1700000302, 0), system: "root pom metadata"}
	nestedPomInfo := &filteredWalkFileInfo{name: "pom.xml", size: 43, mode: 0602, modTime: time.Unix(1700000303, 0), system: "nested pom metadata"}
	ignoredInfo := &filteredWalkFileInfo{name: "prefix-ignored-first-suffix.txt", size: 44, mode: 0603, modTime: time.Unix(1700000304, 0), system: "ignored metadata"}
	renderInfo := &filteredWalkFileInfo{name: "pom.xml.render", size: 45, mode: 0604, modTime: time.Unix(1700000305, 0), system: "render metadata"}
	backslashPomInfo := &filteredWalkFileInfo{name: "pom.xml", size: 46, mode: 0605, modTime: time.Unix(1700000306, 0), system: "backslash root metadata"}
	laterFilterInfo := &filteredWalkFileInfo{name: "pom.xml", size: 47, mode: 0606, modTime: time.Unix(1700000307, 0), system: "later filter metadata"}
	includedInfo := &filteredWalkFileInfo{name: "included.txt", size: 48, mode: 0607, modTime: time.Unix(1700000308, 0), system: "included metadata"}
	rootPom := sourceDir + "/pom.xml"
	nestedPom := sourceDir + "/module/pom.xml"
	ignoredPath := sourceDir + "/nested/prefix-ignored-first-suffix.txt"
	renderPath := sourceDir + "/nested/ignored-first/pom.xml.render"
	backslashRootPom := sourceDir + `\pom.xml`
	laterFilterPath := sourceDir + "/ignored-second/pom.xml"
	includedPath := sourceDir + "/nested/included.txt"
	recording := &recordingFilteredWalkFilesystem{walkInputs: []recordedFilteredWalkCallback{
		{path: sourceDir, info: directoryInfo, err: incomingError},
		{path: rootPom, info: rootPomInfo, err: incomingError},
		{path: nestedPom, info: nestedPomInfo, err: incomingError},
		{path: ignoredPath, info: ignoredInfo, err: incomingError},
		{path: renderPath, info: renderInfo, err: incomingError},
		{path: backslashRootPom, info: backslashPomInfo, err: incomingError},
		{path: laterFilterPath, info: laterFilterInfo, err: incomingError},
		{path: includedPath, info: includedInfo, err: incomingError},
	}}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("filtered template walk dependency lost its complete filesystem value: %#v", dependencies)
	}

	files, err := filteredFilesFromTemplateWithDependencies(dependencies, sourceDir, filter)

	wantFiles := []string{nestedPom, renderPath, backslashRootPom, includedPath}
	if err != nil || !reflect.DeepEqual(files, wantFiles) {
		t.Fatalf("filtered template walk result was (%#v, %v), want ordered files (%#v, nil)", files, err, wantFiles)
	}
	roots, populationErr := recording.assertedWalkRoots()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(roots, []string{sourceDir}) {
		t.Fatalf("filtered template walk roots were %#v, want exact supplied root %#v", roots, []string{sourceDir})
	}
	callbacks, populationErr := recording.assertedCallbackInputs()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(callbacks, recording.walkInputs) {
		t.Fatalf("filtered template callback inputs were %#v, want ordered paths, metadata, and errors %#v", callbacks, recording.walkInputs)
	}
	wantCallbackResults := make([]error, len(recording.walkInputs))
	if !reflect.DeepEqual(recording.callbackResults, wantCallbackResults) {
		t.Fatalf("filtered template callback results were %#v, want all nil", recording.callbackResults)
	}
	if directoryInfo.isDirCalls != 1 || directoryInfo.nameCalls != 0 || directoryInfo.otherMetadataCalls() != 0 {
		t.Fatalf("directory metadata observations were IsDir=%d Name=%d other=%d, want 1, 0, 0",
			directoryInfo.isDirCalls, directoryInfo.nameCalls, directoryInfo.otherMetadataCalls())
	}
	metadata := []*filteredWalkFileInfo{rootPomInfo, nestedPomInfo, ignoredInfo, renderInfo, backslashPomInfo, laterFilterInfo, includedInfo}
	if len(metadata) == 0 {
		t.Fatal("filtered-template metadata expectation population is empty")
	}
	for _, info := range metadata {
		if info.isDirCalls != 1 || info.nameCalls != 1 || info.otherMetadataCalls() != 0 {
			t.Fatalf("non-directory %q metadata observations were IsDir=%d Name=%d other=%d, want 1, 1, 0",
				info.name, info.isDirCalls, info.nameCalls, info.otherMetadataCalls())
		}
	}
	wantLog := "ignoring " + rootPom + " in [pom.xml ignored-first ignored-second]\n" +
		"fileToCopy: " + nestedPom + "\n" +
		"ignoring " + ignoredPath + " in [pom.xml ignored-first ignored-second]\n" +
		"fileToCopy: " + renderPath + "\n" +
		"fileToCopy: " + backslashRootPom + "\n" +
		"ignoring " + laterFilterPath + " in [pom.xml ignored-first ignored-second]\n" +
		"fileToCopy: " + includedPath + "\n"
	if output.String() != wantLog {
		t.Fatalf("filtered template debug log was:\n%q\nwant:\n%q", output.String(), wantLog)
	}
}

func TestFilteredFilesFromTemplateRetainsNilInfoIncomingErrorPanic(t *testing.T) {
	sourceDir := "/complete filtered-template/nil-info root"
	incomingError := &os.PathError{Op: "lstat", Path: sourceDir, Err: errors.New("ordinary system walk error shape")}
	recording := &recordingFilteredWalkFilesystem{walkInputs: []recordedFilteredWalkCallback{{
		path: sourceDir,
		info: nil,
		err:  incomingError,
	}}}
	var recovered interface{}

	func() {
		defer func() { recovered = recover() }()
		_, _ = filteredFilesFromTemplateWithDependencies(recording.dependencies(), sourceDir, nil)
	}()

	if recovered == nil {
		t.Fatal("nil-info incoming walk error did not retain the existing panic")
	}
	roots, populationErr := recording.assertedWalkRoots()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(roots, []string{sourceDir}) {
		t.Fatalf("panic walk roots were %#v, want exact root %#v", roots, []string{sourceDir})
	}
	callbacks, populationErr := recording.assertedCallbackInputs()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(callbacks, recording.walkInputs) {
		t.Fatalf("panic callback inputs were %#v, want exact nil-info incoming-error callback %#v", callbacks, recording.walkInputs)
	}
	if recording.callbackResults != nil {
		t.Fatalf("panic callback unexpectedly returned results %#v", recording.callbackResults)
	}
}

func TestFilteredFilesFromTemplateReturnsEveryExactFinalWalkErrorWithCurrentOrderedPartialResult(t *testing.T) {
	walkErrors := []struct {
		name string
		err  error
	}{
		{name: "ordinary error", err: errors.New("complete filtered-template walk dependency error")},
		{name: "EOF is not normalized", err: io.EOF},
		{name: "wrapped EOF", err: &os.PathError{Op: "walk", Path: "/complete filtered-template root", Err: io.EOF}},
	}

	if len(walkErrors) == 0 {
		t.Fatal("filtered-template final walk-error test-case population is empty")
	}

	for _, walkError := range walkErrors {
		t.Run(walkError.name, func(t *testing.T) {
			sourceDir := "/complete filtered-template/walk-error root"
			includedPath := sourceDir + "/ordered/partial-result.txt"
			recording := &recordingFilteredWalkFilesystem{
				walkInputs: []recordedFilteredWalkCallback{{path: includedPath, info: &filteredWalkFileInfo{name: "partial-result.txt"}}},
				walkErr:    walkError.err,
			}

			files, err := filteredFilesFromTemplateWithDependencies(recording.dependencies(), sourceDir, []string{"ignored"})

			if err != walkError.err || !reflect.DeepEqual(files, []string{includedPath}) {
				t.Fatalf("walk-error filtered template result was (%#v, %v), want current partial result and exact error (%#v, %v)",
					files, err, []string{includedPath}, walkError.err)
			}
			if !reflect.DeepEqual(recording.callbackResults, []error{nil}) {
				t.Fatalf("walk-error callback results were %#v, want nil", recording.callbackResults)
			}
		})
	}
}

func TestFilteredFilesFromTemplateReturnsNilForDirectoriesAndOrderedIgnoredEntries(t *testing.T) {
	output := captureFilteredWalkDebugLog(t)
	sourceDir := "/complete filtered-template/nil-result root"
	rootPom := sourceDir + "/pom.xml"
	ignoredPath := sourceDir + "/nested/ignored-substring.txt"
	recording := &recordingFilteredWalkFilesystem{walkInputs: []recordedFilteredWalkCallback{
		{path: sourceDir, info: &filteredWalkFileInfo{name: "nil-result root", directory: true}},
		{path: rootPom, info: &filteredWalkFileInfo{name: "pom.xml"}},
		{path: ignoredPath, info: &filteredWalkFileInfo{name: "ignored-substring.txt"}},
	}}

	files, err := filteredFilesFromTemplateWithDependencies(recording.dependencies(), sourceDir, []string{"pom.xml", "ignored-substring"})

	if files != nil || err != nil {
		t.Fatalf("no-match filtered template result was (%#v, %v), want (nil, nil)", files, err)
	}
	wantLog := "ignoring " + rootPom + " in [pom.xml ignored-substring]\n" +
		"ignoring " + ignoredPath + " in [pom.xml ignored-substring]\n"
	if output.String() != wantLog {
		t.Fatalf("no-match filtered template debug log was %q, want %q", output.String(), wantLog)
	}
}

func TestFilteredFilesFromTemplateDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	files, err := filteredFilesFromTemplateWithDependencies(
		filteredWalkDependencies{},
		"/developer/home/project/must-not-be-accessed",
		[]string{"must-not-match"},
	)

	if files != nil || !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe filtered template dependency default returned (%#v, %v), want (nil, %v)", files, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedFilteredFilesFromTemplateRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingFilteredWalkFilesystem{}

	if _, err := recording.assertedWalkRoots(); err == nil {
		t.Fatal("empty recorded filtered-template walk-root population passed")
	}
	if _, err := recording.assertedCallbackInputs(); err == nil {
		t.Fatal("empty recorded filtered-template callback population passed")
	}
}
