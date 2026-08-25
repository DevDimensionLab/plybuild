package config

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/sirupsen/logrus"
)

func (*recordingTemplatesFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected templates archive open")
}

func (*recordingTemplatesFilesystem) Close(filesystem.File) error {
	return errors.New("unexpected Templates close")
}

func (*recordingTemplatesFilesystem) CloseReader(io.Closer) error {
	return errors.New("unexpected Templates reader close")
}

type recordedTemplatesCallback struct {
	path string
	info fs.FileInfo
	err  error
}

type recordingTemplatesFilesystem struct {
	walkRoots       []string
	walkInputs      []recordedTemplatesCallback
	callbackInputs  []recordedTemplatesCallback
	callbackResults []error
	walkErr         error
}

func (*recordingTemplatesFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected Templates working directory")
}

func (*recordingTemplatesFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected Templates read")
}

func (*recordingTemplatesFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, errors.New("unexpected Templates read directory")
}

func (*recordingTemplatesFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected Templates read directory entries")
}

func (*recordingTemplatesFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected Templates stat")
}

func (*recordingTemplatesFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected Templates mkdir")
}

func (*recordingTemplatesFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected Templates mkdir")
}

func (*recordingTemplatesFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected Templates write")
}

func (*recordingTemplatesFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected Templates open file")
}

func (*recordingTemplatesFilesystem) Remove(string) error {
	return errors.New("unexpected Templates remove")
}

func (*recordingTemplatesFilesystem) RemoveAll(string) error {
	return errors.New("unexpected Templates recursive remove")
}

func (*recordingTemplatesFilesystem) Rename(string, string) error {
	return errors.New("unexpected Templates rename")
}

func (*recordingTemplatesFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected Templates glob")
}

func (recording *recordingTemplatesFilesystem) Walk(root string, callback filepath.WalkFunc) error {
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

func (*recordingTemplatesFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected Templates create")
}

func (*recordingTemplatesFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected Templates copy")
}

func (recording *recordingTemplatesFilesystem) dependencies() templatesDependencies {
	return templatesDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingTemplatesFilesystem) assertedWalkRoots() ([]string, error) {
	if len(recording.walkRoots) == 0 {
		return nil, errors.New("recorded Templates walk-root population is empty")
	}
	return recording.walkRoots, nil
}

func (recording *recordingTemplatesFilesystem) assertedCallbackInputs() ([]recordedTemplatesCallback, error) {
	if len(recording.callbackInputs) == 0 {
		return nil, errors.New("recorded Templates callback population is empty")
	}
	return recording.callbackInputs, nil
}

type templatesFileInfo struct {
	name         string
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

func (info *templatesFileInfo) Name() string {
	info.nameCalls++
	return info.name
}

func (info *templatesFileInfo) Size() int64 {
	info.sizeCalls++
	return info.size
}

func (info *templatesFileInfo) Mode() fs.FileMode {
	info.modeCalls++
	return info.mode
}

func (info *templatesFileInfo) ModTime() time.Time {
	info.modTimeCalls++
	return info.modTime
}

func (info *templatesFileInfo) IsDir() bool {
	info.isDirCalls++
	return info.mode.IsDir()
}

func (info *templatesFileInfo) Sys() interface{} {
	info.systemCalls++
	return info.system
}

func (info *templatesFileInfo) nonNameMetadataCalls() int {
	return info.sizeCalls + info.modeCalls + info.modTimeCalls + info.isDirCalls + info.systemCalls
}

type templatesMessageOnlyFormatter struct{}

func (templatesMessageOnlyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func TestTemplatesSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemTemplatesDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Templates selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Templates filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestTemplatesPreservesReceiverRootCallbackOrderPathsMetadataErrorsMatchingNamesAndProjects(t *testing.T) {
	isolateConfigTestHome(t)
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config receiver/root with spaces"}}
	fixtureRoot := file.Path("test/cloud-config/templates")
	currentDirectory := file.Path("%s/test-template", fixtureRoot)
	legacyDirectory := file.Path("%s/legacy-category/legacy-template", fixtureRoot)
	currentPath := file.Path("%s/%s", currentDirectory, projectConfigFileName)
	legacyPath := file.Path("%s/%s", legacyDirectory, legacyProjectConfigFileName)
	incomingError := errors.New("suppressed incoming Templates walk error")
	rootInfo := &templatesFileInfo{name: "templates", size: 41, mode: fs.ModeDir | 0751, modTime: time.Unix(1700000201, 0), system: "root metadata"}
	ignoredInfo := &templatesFileInfo{name: "notes.txt", size: 82, mode: 0641, modTime: time.Unix(1700000202, 0), system: "ignored metadata"}
	currentInfo := &templatesFileInfo{name: projectConfigFileName, size: 83, mode: 0601, modTime: time.Unix(1700000203, 0), system: "current metadata"}
	legacyInfo := &templatesFileInfo{name: legacyProjectConfigFileName, size: 84, mode: 0605, modTime: time.Unix(1700000204, 0), system: "legacy metadata"}
	recording := &recordingTemplatesFilesystem{walkInputs: []recordedTemplatesCallback{
		{path: fixtureRoot, info: rootInfo},
		{path: file.Path("%s/unreadable-entry", fixtureRoot), info: nil, err: incomingError},
		{path: file.Path("%s/notes.txt", fixtureRoot), info: ignoredInfo},
		{path: currentPath, info: currentInfo},
		{path: legacyPath, info: legacyInfo},
	}}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("Templates dependency lost its complete filesystem value: %#v", dependencies)
	}

	actual, err := gitCfg.templates(dependencies)

	if err != nil {
		t.Fatalf("Templates returned an error: %v", err)
	}
	wantRoot := file.Path("%s/templates", gitCfg.Implementation().Dir())
	roots, populationErr := recording.assertedWalkRoots()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(roots, []string{wantRoot}) {
		t.Fatalf("Templates walk roots were %#v, want exact receiver-derived root %#v", roots, []string{wantRoot})
	}
	callbacks, populationErr := recording.assertedCallbackInputs()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(callbacks, recording.walkInputs) {
		t.Fatalf("Templates callback inputs were %#v, want ordered paths, metadata, and errors %#v", callbacks, recording.walkInputs)
	}
	if !reflect.DeepEqual(recording.callbackResults, []error{nil, nil, nil, nil, nil}) {
		t.Fatalf("Templates callback results were %#v, want every incoming walk error suppressed and all callbacks nil", recording.callbackResults)
	}

	wantNames := []string{"test-template", file.Path("legacy-category/legacy-template")}
	if len(actual) != len(wantNames) {
		t.Fatalf("Templates returned %#v, want %d ordered matches", actual, len(wantNames))
	}
	for index, wantName := range wantNames {
		if actual[index].Name != wantName {
			t.Fatalf("Templates result %d name was %q, want exact relative name %q", index, actual[index].Name, wantName)
		}
	}
	wantDirectories := []string{currentDirectory, legacyDirectory}
	wantConfigFiles := []string{currentPath, legacyPath}
	for index := range actual {
		if actual[index].Project.Path != wantDirectories[index] || actual[index].Project.ConfigFile != wantConfigFiles[index] {
			t.Fatalf("Templates result %d loaded project path/config was (%q, %q), want (%q, %q)",
				index, actual[index].Project.Path, actual[index].Project.ConfigFile, wantDirectories[index], wantConfigFiles[index])
		}
	}
	if actual[0].Project.Config.Name != "test-template" || actual[1].Project.Config.Name != "legacy-template" {
		t.Fatalf("Templates loaded project configurations were (%q, %q), want exact fixture projects",
			actual[0].Project.Config.Name, actual[1].Project.Config.Name)
	}
	if !actual[0].Project.IsMavenProject() || actual[1].Project.IsMavenProject() {
		t.Fatalf("Templates loaded project types were (%T, %T), want Maven then non-Maven", actual[0].Project.Type, actual[1].Project.Type)
	}

	if rootInfo.nameCalls != 2 || ignoredInfo.nameCalls != 2 || currentInfo.nameCalls != 2 || legacyInfo.nameCalls != 3 {
		t.Fatalf("Templates info.Name calls were root=%d ignored=%d current=%d legacy=%d, want exact comparisons 2,2,2,3",
			rootInfo.nameCalls, ignoredInfo.nameCalls, currentInfo.nameCalls, legacyInfo.nameCalls)
	}
	for _, info := range []*templatesFileInfo{rootInfo, ignoredInfo, currentInfo, legacyInfo} {
		if calls := info.nonNameMetadataCalls(); calls != 0 {
			t.Fatalf("Templates inspected non-name metadata for %q %d times", info.name, calls)
		}
	}
}

func TestTemplatesReturnsAndLogsExactProjectLoadErrorWithOrderedPartialResult(t *testing.T) {
	isolateConfigTestHome(t)
	output := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetOutput(output)
	logger.SetFormatter(templatesMessageOnlyFormatter{})
	previousLogger := log
	log = logger
	t.Cleanup(func() { log = previousLogger })

	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete project-load-error receiver"}}
	fixtureRoot := file.Path("test/cloud-config/templates")
	currentDirectory := file.Path("%s/test-template", fixtureRoot)
	currentPath := file.Path("%s/%s", currentDirectory, projectConfigFileName)
	brokenDirectory := file.Path("test/templates-load-error/templates/broken-template")
	brokenPath := file.Path("%s/%s", brokenDirectory, projectConfigFileName)
	recording := &recordingTemplatesFilesystem{walkInputs: []recordedTemplatesCallback{
		{path: currentPath, info: &templatesFileInfo{name: projectConfigFileName}},
		{path: brokenPath, info: &templatesFileInfo{name: projectConfigFileName}},
		{path: file.Path("%s/legacy-category/legacy-template/%s", fixtureRoot, legacyProjectConfigFileName), info: &templatesFileInfo{name: legacyProjectConfigFileName}},
	}}

	actual, err := gitCfg.templates(recording.dependencies())

	wantError := fmt.Sprintf("%s directory detected, but language was not set in %s",
		file.Path("%s/src", brokenDirectory), "config file (ply.json, or co-pilot.json")
	if err == nil || err.Error() != wantError {
		t.Fatalf("Templates project-load error was %v, want exact error %q", err, wantError)
	}
	if len(recording.callbackResults) != 2 || recording.callbackResults[0] != nil || recording.callbackResults[1] != err {
		t.Fatalf("Templates callback results were %#v, want nil then the exact returned project-load error", recording.callbackResults)
	}
	if len(recording.callbackInputs) != 2 {
		t.Fatalf("Templates delivered %d callbacks, want stop at project-load error after 2", len(recording.callbackInputs))
	}
	if len(actual) != 1 || actual[0].Name != "test-template" || actual[0].Project.Path != currentDirectory {
		t.Fatalf("Templates project-load error partial result was %#v, want the first exact loaded template", actual)
	}
	if output.String() != err.Error()+"\n" {
		t.Fatalf("Templates project-load error log was %q, want exact returned error %q", output.String(), err.Error()+"\n")
	}
}

func TestTemplatesReturnsEveryExactFinalWalkErrorWithCurrentOrderedPartialResult(t *testing.T) {
	walkErrors := []struct {
		name string
		err  error
	}{
		{name: "ordinary error", err: errors.New("complete Templates walk dependency error")},
		{name: "EOF is not normalized", err: io.EOF},
		{name: "wrapped EOF", err: &os.PathError{Op: "walk", Path: "/complete Templates root", Err: io.EOF}},
	}

	for _, walkError := range walkErrors {
		t.Run(walkError.name, func(t *testing.T) {
			isolateConfigTestHome(t)
			fixtureRoot := file.Path("test/cloud-config/templates")
			currentDirectory := file.Path("%s/test-template", fixtureRoot)
			currentPath := file.Path("%s/%s", currentDirectory, projectConfigFileName)
			recording := &recordingTemplatesFilesystem{
				walkInputs: []recordedTemplatesCallback{{path: currentPath, info: &templatesFileInfo{name: projectConfigFileName}}},
				walkErr:    walkError.err,
			}

			actual, err := (GitCloudConfig{Impl: DirConfig{Path: "/complete walk-error receiver"}}).templates(recording.dependencies())

			if err != walkError.err {
				t.Fatalf("Templates walk error was %v, want exact dependency error %v", err, walkError.err)
			}
			if len(actual) != 1 || actual[0].Name != "test-template" || actual[0].Project.Path != currentDirectory {
				t.Fatalf("Templates walk-error partial result was %#v, want current exact ordered result", actual)
			}
			if !reflect.DeepEqual(recording.callbackResults, []error{nil}) {
				t.Fatalf("Templates walk-error callback results were %#v, want nil", recording.callbackResults)
			}
		})
	}
}

func TestTemplatesReturnsNilForOnlyIgnoredAndIncomingWalkErrorCallbacks(t *testing.T) {
	incomingError := errors.New("suppressed Templates no-match incoming error")
	ignoredInfo := &templatesFileInfo{name: "Ply.json", size: 41, mode: fs.ModeDir | 0751, system: "ignored metadata"}
	recording := &recordingTemplatesFilesystem{walkInputs: []recordedTemplatesCallback{
		{path: "/developer/home/cloud-config/must-not-be-accessed", info: nil, err: incomingError},
		{path: "/complete ignored/path/Ply.json", info: ignoredInfo},
	}}

	actual, err := (GitCloudConfig{Impl: DirConfig{Path: "/complete nil-result receiver"}}).templates(recording.dependencies())

	if actual != nil || err != nil {
		t.Fatalf("Templates no-match result was (%#v, %v), want (nil, nil)", actual, err)
	}
	if !reflect.DeepEqual(recording.callbackResults, []error{nil, nil}) {
		t.Fatalf("Templates no-match callback results were %#v, want suppressed incoming error and ignored entry", recording.callbackResults)
	}
	if ignoredInfo.nameCalls != 2 || ignoredInfo.nonNameMetadataCalls() != 0 {
		t.Fatalf("Templates ignored entry observations were Name=%d other=%d, want exact two name comparisons only",
			ignoredInfo.nameCalls, ignoredInfo.nonNameMetadataCalls())
	}
}

func TestTemplatesDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/developer/home/cloud-config/must-not-be-accessed"}}

	actual, err := gitCfg.templates(templatesDependencies{})

	if actual != nil || !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe Templates dependency default returned (%#v, %v), want (nil, %v)", actual, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedTemplatesRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingTemplatesFilesystem{}

	if _, err := recording.assertedWalkRoots(); err == nil {
		t.Fatal("empty recorded Templates walk-root population passed")
	}
	if _, err := recording.assertedCallbackInputs(); err == nil {
		t.Fatal("empty recorded Templates callback population passed")
	}
}
