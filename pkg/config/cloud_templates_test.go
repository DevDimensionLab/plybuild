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

func (*recordingTemplatesFilesystem) OpenZipEntry(*zip.File) (io.ReadCloser, error) {
	return nil, errors.New("unexpected templates archive entry open")
}

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

type templateProjectLoadResult struct {
	project Project
	err     error
}

type recordingTemplateProjectLoader struct {
	directories []string
	results     []templateProjectLoadResult
}

func (recording *recordingTemplateProjectLoader) Load(directory string) (Project, error) {
	recording.directories = append(recording.directories, directory)
	index := len(recording.directories) - 1
	if index >= len(recording.results) {
		return Project{}, errors.New("recorded Templates project result population is incomplete")
	}
	result := recording.results[index]
	return result.project, result.err
}

func (recording *recordingTemplateProjectLoader) assertedDirectories() ([]string, error) {
	if len(recording.directories) == 0 {
		return nil, errors.New("recorded Templates project-loader directory population is empty")
	}
	return recording.directories, nil
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

func (recording *recordingTemplatesFilesystem) dependencies(loader templateProjectLoader) templatesDependencies {
	return templatesDependencies{
		Files:  filesystem.Dependencies{FileSystem: recording},
		Loader: loader,
	}
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

type templatesCloudConfigIdentity struct {
	GitCloudConfig
}

func TestTemplatesSelectsCompleteSystemFilesystemAndProjectLoaderDependencies(t *testing.T) {
	dependencies := systemTemplatesDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Templates selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Templates filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
	if dependencies.Loader == nil {
		t.Fatal("Templates selected an incomplete project-loader dependency")
	}
	if reflect.TypeOf(dependencies.Loader) != reflect.TypeOf(initTemplateProjectLoader{}) {
		t.Fatalf("Templates project-loader dependency is %T, want %T", dependencies.Loader, initTemplateProjectLoader{})
	}
}

func TestInitTemplateProjectLoaderPreservesCompleteProjectAndExactError(t *testing.T) {
	tests := []struct {
		name      string
		directory string
		assert    func(*testing.T, Project, error)
	}{
		{
			name:      "complete tracked project",
			directory: file.Path("test/cloud-config/templates/test-template"),
			assert: func(t *testing.T, project Project, err error) {
				if err != nil {
					t.Fatalf("production Templates project loader returned an error: %v", err)
				}
				wantDirectory := file.Path("test/cloud-config/templates/test-template")
				wantConfig := file.Path("%s/%s", wantDirectory, projectConfigFileName)
				if project.Path != wantDirectory || project.ConfigFile != wantConfig || project.Config.Name != "test-template" {
					t.Fatalf("production Templates project loader returned %#v, want complete tracked project", project)
				}
				if !project.IsMavenProject() || project.Type == nil || project.CloudConfig == nil {
					t.Fatalf("production Templates project loader lost embedded interfaces: Type=%T CloudConfig=%T", project.Type, project.CloudConfig)
				}
			},
		},
		{
			name:      "exact project error",
			directory: file.Path("test/templates-load-error/templates/broken-template"),
			assert: func(t *testing.T, project Project, err error) {
				wantError := fmt.Sprintf("%s directory detected, but language was not set in %s",
					file.Path("test/templates-load-error/templates/broken-template/src"),
					"config file (ply.json, or co-pilot.json")
				if !reflect.DeepEqual(project, Project{}) || err == nil || err.Error() != wantError {
					t.Fatalf("production Templates project-loader error result was (%#v, %v), want (zero, exact %q)", project, err, wantError)
				}
			},
		},
	}
	if len(tests) == 0 {
		t.Fatal("production Templates project-loader test-case population is empty")
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isolateConfigTestHome(t)
			project, err := (initTemplateProjectLoader{}).Load(test.directory)
			test.assert(t, project, err)
		})
	}
}

func TestTemplatesProductionRetainsTrackedCurrentAndLegacyProjects(t *testing.T) {
	isolateConfigTestHome(t)
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: file.Path("test/cloud-config")}}

	actual, err := gitCfg.Templates()

	if err != nil {
		t.Fatalf("production Templates returned an error: %v", err)
	}
	wantNames := []string{file.Path("legacy-category/legacy-template"), "test-template"}
	if len(actual) != len(wantNames) {
		t.Fatalf("production Templates returned %#v, want %d tracked projects", actual, len(wantNames))
	}
	for index, wantName := range wantNames {
		wantDirectory := file.Path("test/cloud-config/templates/%s", wantName)
		wantConfigName := projectConfigFileName
		wantProjectName := "test-template"
		wantMaven := true
		if index == 0 {
			wantConfigName = legacyProjectConfigFileName
			wantProjectName = "legacy-template"
			wantMaven = false
		}
		if actual[index].Name != wantName || actual[index].Project.Path != wantDirectory ||
			actual[index].Project.ConfigFile != file.Path("%s/%s", wantDirectory, wantConfigName) ||
			actual[index].Project.Config.Name != wantProjectName || actual[index].Project.IsMavenProject() != wantMaven ||
			actual[index].Project.CloudConfig == nil {
			t.Fatalf("production Templates result %d was %#v, want exact tracked %q project", index, actual[index], wantName)
		}
	}
}

func TestTemplatesPreservesReceiverRootCallbackOrderPathsMetadataErrorsMatchingNamesAndProjects(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config receiver/root with spaces"}}
	fixtureRoot := file.Path("%s/templates", gitCfg.Implementation().Dir())
	currentDirectory := file.Path("%s/current category/current template", fixtureRoot)
	legacyDirectory := file.Path("%s/legacy category/legacy template", fixtureRoot)
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
	currentCloud := &templatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/complete current cloud identity"}}}
	legacyCloud := &templatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/complete legacy cloud identity"}}}
	currentType := &MavenProject{PomFile: "/complete current pom identity"}
	legacyType := &MavenProject{PomFile: "/complete legacy pom identity"}
	wantProjects := []Project{
		{
			Path: currentDirectory,
			GitInfo: GitInfo{
				IsRepo: true, IsDirty: true, EnableCommit: true,
			},
			ConfigFile:  currentPath,
			Config:      ProjectConfiguration{Name: "complete current project", Profile: "current profile"},
			Type:        currentType,
			CloudConfig: currentCloud,
		},
		{
			Path: legacyDirectory,
			GitInfo: GitInfo{
				IsRepo: true, EnableCommit: true,
			},
			ConfigFile:  legacyPath,
			Config:      ProjectConfiguration{Name: "complete legacy project", Profile: "legacy profile"},
			Type:        legacyType,
			CloudConfig: legacyCloud,
		},
	}
	loader := &recordingTemplateProjectLoader{results: []templateProjectLoadResult{
		{project: wantProjects[0]},
		{project: wantProjects[1]},
	}}
	dependencies := recording.dependencies(loader)
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("Templates dependency lost its complete filesystem value: %#v", dependencies)
	}
	if dependencies.Loader != loader {
		t.Fatalf("Templates dependency lost its complete project-loader value: %#v", dependencies)
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

	directories, populationErr := loader.assertedDirectories()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantDirectories := []string{currentDirectory, legacyDirectory}
	if !reflect.DeepEqual(directories, wantDirectories) {
		t.Fatalf("Templates project-loader directories were %#v, want one exact load per current/legacy match %#v", directories, wantDirectories)
	}

	wantNames := []string{file.Path("current category/current template"), file.Path("legacy category/legacy template")}
	if len(wantNames) == 0 {
		t.Fatal("Templates name expectation population is empty")
	}
	if len(actual) != len(wantNames) {
		t.Fatalf("Templates returned %#v, want %d ordered matches", actual, len(wantNames))
	}
	for index, wantName := range wantNames {
		if actual[index].Name != wantName {
			t.Fatalf("Templates result %d name was %q, want exact relative name %q", index, actual[index].Name, wantName)
		}
		if !reflect.DeepEqual(actual[index].Project, wantProjects[index]) {
			t.Fatalf("Templates result %d project was %#v, want complete exact value %#v", index, actual[index].Project, wantProjects[index])
		}
		if actual[index].Project.CloudConfig != wantProjects[index].CloudConfig {
			t.Fatalf("Templates result %d cloud-config identity was %T %#v, want exact %T %#v",
				index, actual[index].Project.CloudConfig, actual[index].Project.CloudConfig,
				wantProjects[index].CloudConfig, wantProjects[index].CloudConfig)
		}
		if actual[index].Project.Type != wantProjects[index].Type {
			t.Fatalf("Templates result %d project-type identity was %T %#v, want exact %T %#v",
				index, actual[index].Project.Type, actual[index].Project.Type,
				wantProjects[index].Type, wantProjects[index].Type)
		}
	}

	if rootInfo.nameCalls != 2 || ignoredInfo.nameCalls != 2 || currentInfo.nameCalls != 2 || legacyInfo.nameCalls != 3 {
		t.Fatalf("Templates info.Name calls were root=%d ignored=%d current=%d legacy=%d, want exact comparisons 2,2,2,3",
			rootInfo.nameCalls, ignoredInfo.nameCalls, currentInfo.nameCalls, legacyInfo.nameCalls)
	}
	metadata := []*templatesFileInfo{rootInfo, ignoredInfo, currentInfo, legacyInfo}
	if len(metadata) == 0 {
		t.Fatal("Templates metadata expectation population is empty")
	}
	for _, info := range metadata {
		if calls := info.nonNameMetadataCalls(); calls != 0 {
			t.Fatalf("Templates inspected non-name metadata for %q %d times", info.name, calls)
		}
	}
}

func TestTemplatesReturnsAndLogsExactProjectLoadErrorWithOrderedPartialResult(t *testing.T) {
	output := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetOutput(output)
	logger.SetFormatter(templatesMessageOnlyFormatter{})
	previousLogger := log
	log = logger
	t.Cleanup(func() { log = previousLogger })

	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete project-load-error receiver"}}
	fixtureRoot := file.Path("%s/templates", gitCfg.Implementation().Dir())
	currentDirectory := file.Path("%s/first-template", fixtureRoot)
	currentPath := file.Path("%s/%s", currentDirectory, projectConfigFileName)
	brokenDirectory := file.Path("%s/broken-template", fixtureRoot)
	brokenPath := file.Path("%s/%s", brokenDirectory, projectConfigFileName)
	thirdDirectory := file.Path("%s/third-template", fixtureRoot)
	projectLoadError := errors.New("complete Templates project-loader dependency error")
	firstProject := Project{
		Path:        currentDirectory,
		ConfigFile:  currentPath,
		Config:      ProjectConfiguration{Name: "complete first project"},
		Type:        &MavenProject{PomFile: "/complete first pom identity"},
		CloudConfig: &templatesCloudConfigIdentity{},
	}
	recording := &recordingTemplatesFilesystem{walkInputs: []recordedTemplatesCallback{
		{path: currentPath, info: &templatesFileInfo{name: projectConfigFileName}},
		{path: brokenPath, info: &templatesFileInfo{name: projectConfigFileName}},
		{path: file.Path("%s/%s", thirdDirectory, legacyProjectConfigFileName), info: &templatesFileInfo{name: legacyProjectConfigFileName}},
	}}
	loader := &recordingTemplateProjectLoader{results: []templateProjectLoadResult{
		{project: firstProject},
		{project: Project{Path: brokenDirectory, Config: ProjectConfiguration{Name: "complete failed partial project"}}, err: projectLoadError},
		{project: Project{Path: thirdDirectory}},
	}}

	actual, err := gitCfg.templates(recording.dependencies(loader))

	if err != projectLoadError {
		t.Fatalf("Templates project-load error was %v, want exact dependency error %v", err, projectLoadError)
	}
	if len(recording.callbackResults) != 2 || recording.callbackResults[0] != nil || recording.callbackResults[1] != err {
		t.Fatalf("Templates callback results were %#v, want nil then the exact returned project-load error", recording.callbackResults)
	}
	if len(recording.callbackInputs) != 2 {
		t.Fatalf("Templates delivered %d callbacks, want stop at project-load error after 2", len(recording.callbackInputs))
	}
	directories, populationErr := loader.assertedDirectories()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(directories, []string{currentDirectory, brokenDirectory}) {
		t.Fatalf("Templates project-load-error directories were %#v, want stop after exact failed load", directories)
	}
	if len(actual) != 1 || actual[0].Name != "first-template" || !reflect.DeepEqual(actual[0].Project, firstProject) {
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

	if len(walkErrors) == 0 {
		t.Fatal("Templates final walk-error test-case population is empty")
	}

	for _, walkError := range walkErrors {
		t.Run(walkError.name, func(t *testing.T) {
			fixtureRoot := file.Path("/complete walk-error receiver/templates")
			currentDirectory := file.Path("%s/current-template", fixtureRoot)
			currentPath := file.Path("%s/%s", currentDirectory, projectConfigFileName)
			currentProject := Project{
				Path:        currentDirectory,
				ConfigFile:  currentPath,
				Config:      ProjectConfiguration{Name: "complete walk-error project"},
				CloudConfig: &templatesCloudConfigIdentity{},
			}
			recording := &recordingTemplatesFilesystem{
				walkInputs: []recordedTemplatesCallback{{path: currentPath, info: &templatesFileInfo{name: projectConfigFileName}}},
				walkErr:    walkError.err,
			}
			loader := &recordingTemplateProjectLoader{results: []templateProjectLoadResult{{project: currentProject}}}

			actual, err := (GitCloudConfig{Impl: DirConfig{Path: "/complete walk-error receiver"}}).templates(recording.dependencies(loader))

			if err != walkError.err {
				t.Fatalf("Templates walk error was %v, want exact dependency error %v", err, walkError.err)
			}
			if len(actual) != 1 || actual[0].Name != "current-template" || !reflect.DeepEqual(actual[0].Project, currentProject) {
				t.Fatalf("Templates walk-error partial result was %#v, want current exact ordered result", actual)
			}
			if !reflect.DeepEqual(recording.callbackResults, []error{nil}) {
				t.Fatalf("Templates walk-error callback results were %#v, want nil", recording.callbackResults)
			}
			directories, populationErr := loader.assertedDirectories()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(directories, []string{currentDirectory}) {
				t.Fatalf("Templates walk-error project-loader directories were %#v, want one exact load", directories)
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
	loader := &recordingTemplateProjectLoader{}

	actual, err := (GitCloudConfig{Impl: DirConfig{Path: "/complete nil-result receiver"}}).templates(recording.dependencies(loader))

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
	if len(loader.directories) != 0 {
		t.Fatalf("Templates loaded projects for ignored or incoming-error callbacks: %#v", loader.directories)
	}
	if _, populationErr := loader.assertedDirectories(); populationErr == nil {
		t.Fatal("empty Templates project-loader recording unexpectedly passed")
	}
}

func TestTemplatesFilesystemDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/developer/home/cloud-config/must-not-be-accessed"}}

	actual, err := gitCfg.templates(templatesDependencies{})

	if actual != nil || !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe Templates dependency default returned (%#v, %v), want (nil, %v)", actual, err, filesystem.ErrNoFilesystem)
	}
}

func TestTemplatesProjectLoaderDefaultsToSafeExactErrorWithInjectedWalk(t *testing.T) {
	projectDirectory := file.Path("/developer/home/cloud-config/must-not-be-accessed/templates/complete-template")
	projectPath := file.Path("%s/%s", projectDirectory, projectConfigFileName)
	recording := &recordingTemplatesFilesystem{walkInputs: []recordedTemplatesCallback{{
		path: projectPath,
		info: &templatesFileInfo{name: projectConfigFileName},
	}}}

	actual, err := (GitCloudConfig{Impl: DirConfig{Path: "/complete missing-loader receiver"}}).templates(recording.dependencies(nil))

	if actual != nil || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe Templates project-loader default returned (%#v, %v), want (nil, exact %v)", actual, err, filesystem.ErrNoFilesystem)
	}
	if !reflect.DeepEqual(recording.callbackResults, []error{filesystem.ErrNoFilesystem}) {
		t.Fatalf("safe Templates project-loader callback results were %#v, want exact filesystem error", recording.callbackResults)
	}
	if len(recording.callbackInputs) != 1 {
		t.Fatalf("safe Templates project-loader received %d callbacks, want one injected match", len(recording.callbackInputs))
	}
}

func TestTemplatesLoadsProjectsIndependentlyAcrossRepeatedInvocations(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete repeated Templates receiver"}}
	projectDirectory := file.Path("%s/templates/repeated-template", gitCfg.Implementation().Dir())
	projectPath := file.Path("%s/%s", projectDirectory, projectConfigFileName)
	recording := &recordingTemplatesFilesystem{walkInputs: []recordedTemplatesCallback{{
		path: projectPath,
		info: &templatesFileInfo{name: projectConfigFileName},
	}}}
	firstCloud := &templatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/first repeated cloud identity"}}}
	secondCloud := &templatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/second repeated cloud identity"}}}
	first := Project{Path: "first repeated project", Config: ProjectConfiguration{Name: "first"}, CloudConfig: firstCloud}
	second := Project{Path: "second repeated project", Config: ProjectConfiguration{Name: "second"}, CloudConfig: secondCloud}
	loader := &recordingTemplateProjectLoader{results: []templateProjectLoadResult{{project: first}, {project: second}}}
	dependencies := recording.dependencies(loader)

	firstResult, err := gitCfg.templates(dependencies)
	if err != nil || len(firstResult) != 1 || !reflect.DeepEqual(firstResult[0].Project, first) || firstResult[0].Project.CloudConfig != firstCloud {
		t.Fatalf("first repeated Templates invocation returned (%#v, %v), want exact first project", firstResult, err)
	}
	secondResult, err := gitCfg.templates(dependencies)
	if err != nil || len(secondResult) != 1 || !reflect.DeepEqual(secondResult[0].Project, second) || secondResult[0].Project.CloudConfig != secondCloud {
		t.Fatalf("second repeated Templates invocation returned (%#v, %v), want independent exact second project", secondResult, err)
	}
	directories, populationErr := loader.assertedDirectories()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(directories, []string{projectDirectory, projectDirectory}) {
		t.Fatalf("repeated Templates project-loader directories were %#v, want one complete load per invocation", directories)
	}
	if !reflect.DeepEqual(recording.walkRoots, []string{
		file.Path("%s/templates", gitCfg.Implementation().Dir()),
		file.Path("%s/templates", gitCfg.Implementation().Dir()),
	}) {
		t.Fatalf("repeated Templates walk roots were %#v, want independent exact roots", recording.walkRoots)
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
	if _, err := (&recordingTemplateProjectLoader{}).assertedDirectories(); err == nil {
		t.Fatal("empty recorded Templates project-loader directory population passed")
	}
}
