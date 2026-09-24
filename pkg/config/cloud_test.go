package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/testutil"
)

type recordingServicesLoader struct {
	directories []Directory
	services    CloudServices
	err         error
}

func (recording *recordingServicesLoader) Load(directory Directory) (CloudServices, error) {
	recording.directories = append(recording.directories, directory)
	return recording.services, recording.err
}

func (recording *recordingServicesLoader) dependencies() servicesDependencies {
	return servicesDependencies{Loader: recording}
}

func (recording *recordingServicesLoader) assertedDirectories() ([]Directory, error) {
	if len(recording.directories) == 0 {
		return nil, errors.New("recorded Services loader directory population is empty")
	}
	return recording.directories, nil
}

type recordingServicesDirectory struct {
	dir       string
	path      string
	err       error
	fileNames []string
}

type projectDefaultsLoadResult struct {
	defaults CloudProjectDefaults
	err      error
}

type recordingProjectDefaultsLoader struct {
	directories []Directory
	results     []projectDefaultsLoadResult
}

func (recording *recordingProjectDefaultsLoader) Load(directory Directory) (CloudProjectDefaults, error) {
	recording.directories = append(recording.directories, directory)
	index := len(recording.directories) - 1
	if index >= len(recording.results) {
		return CloudProjectDefaults{}, errors.New("recorded ProjectDefaults result population is incomplete")
	}
	result := recording.results[index]
	return result.defaults, result.err
}

func (recording *recordingProjectDefaultsLoader) dependencies() projectDefaultsDependencies {
	return projectDefaultsDependencies{Loader: recording}
}

func (recording *recordingProjectDefaultsLoader) assertedDirectories() ([]Directory, error) {
	if len(recording.directories) == 0 {
		return nil, errors.New("recorded ProjectDefaults loader directory population is empty")
	}
	return recording.directories, nil
}

type recordingProjectDefaultsDirectory struct {
	dir       string
	path      string
	err       error
	fileNames []string
}

func (directory *recordingProjectDefaultsDirectory) Dir() string {
	return directory.dir
}

func (directory *recordingProjectDefaultsDirectory) FilePath(fileName string) (string, error) {
	directory.fileNames = append(directory.fileNames, fileName)
	return directory.path, directory.err
}

func (directory *recordingProjectDefaultsDirectory) assertedFileNames() ([]string, error) {
	if len(directory.fileNames) == 0 {
		return nil, errors.New("recorded ProjectDefaults filename population is empty")
	}
	return directory.fileNames, nil
}

func (directory *recordingServicesDirectory) Dir() string {
	return directory.dir
}

func (directory *recordingServicesDirectory) FilePath(fileName string) (string, error) {
	directory.fileNames = append(directory.fileNames, fileName)
	return directory.path, directory.err
}

func (directory *recordingServicesDirectory) assertedFileNames() ([]string, error) {
	if len(directory.fileNames) == 0 {
		return nil, errors.New("recorded Services filename population is empty")
	}
	return directory.fileNames, nil
}

func newMockCloudConfig() (cfg GitCloudConfig) {
	cfg.Impl.Path = "test/cloud-config"
	return
}

func TestServicesSelectsFileLoaderForProduction(t *testing.T) {
	dependencies := systemServicesDependencies()

	if dependencies.Loader == nil {
		t.Fatal("Services selected an incomplete loader dependency")
	}
	if reflect.TypeOf(dependencies.Loader) != reflect.TypeOf(fileServicesLoader{}) {
		t.Fatalf("Services loader dependency is %T, want %T", dependencies.Loader, fileServicesLoader{})
	}

	directory := &recordingServicesDirectory{path: filepath.Join("test", "cloud-config", "services.json")}
	services, err := dependencies.Load(directory)
	if err != nil {
		t.Fatalf("production Services loader returned an error: %v", err)
	}
	fileNames, populationErr := directory.assertedFileNames()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(fileNames, []string{"services.json"}) {
		t.Fatalf("production Services loader requested filenames %#v, want exact services.json", fileNames)
	}
	if services.Type != "services" || len(services.Data) == 0 {
		t.Fatalf("production Services loader decoded %#v, want the populated tracked fixture", services)
	}
}

func TestServicesLoadsEagerlyOnceWithCompleteDirectoryAndReplaysStableResults(t *testing.T) {
	loadError := errors.New("complete Services loader dependency error")
	tests := []struct {
		name     string
		services CloudServices
		err      error
	}{
		{name: "success", services: CloudServices{Type: "complete services value"}},
		{name: "error", services: CloudServices{Type: "complete partial services value"}, err: loadError},
	}

	if len(tests) == 0 {
		t.Fatal("Services stable-result test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config receiver/root with spaces"}}
			recording := &recordingServicesLoader{services: test.services, err: test.err}

			loaded := gitCfg.services(recording.dependencies())

			directories, populationErr := recording.assertedDirectories()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if len(directories) != 1 || !reflect.DeepEqual(directories[0], gitCfg.Implementation()) {
				t.Fatalf("Services loader received directories %#v, want the complete implementation %#v", directories, gitCfg.Implementation())
			}

			for call := 1; call <= 2; call++ {
				actual, err := loaded()
				if !reflect.DeepEqual(actual, test.services) || err != test.err {
					t.Fatalf("Services closure call %d returned (%#v, %v), want stable exact (%#v, %v)", call, actual, err, test.services, test.err)
				}
			}
			if len(recording.directories) != 1 {
				t.Fatalf("Services loader ran %d times, want one eager load", len(recording.directories))
			}
		})
	}
}

func TestFileServicesLoaderReturnsExactFilePathError(t *testing.T) {
	pathError := errors.New("complete Services FilePath dependency error")
	directory := &recordingServicesDirectory{
		dir:  "/developer/home/cloud-config/must-not-be-accessed",
		err:  pathError,
		path: "/must/not/be/read/services.json",
	}

	services, err := (fileServicesLoader{}).Load(directory)

	if !reflect.DeepEqual(services, CloudServices{}) || err != pathError {
		t.Fatalf("file Services loader returned (%#v, %v), want (zero, exact error %v)", services, err, pathError)
	}
	fileNames, populationErr := directory.assertedFileNames()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(fileNames, []string{"services.json"}) {
		t.Fatalf("file Services loader requested filenames %#v, want exact services.json", fileNames)
	}
}

func TestServicesDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	directory := &recordingServicesDirectory{dir: "/developer/home/cloud-config/must-not-be-accessed"}

	services, err := (servicesDependencies{}).Load(directory)

	if !reflect.DeepEqual(services, CloudServices{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe Services dependency default returned (%#v, %v), want (zero, exact %v)", services, err, filesystem.ErrNoFilesystem)
	}
	if len(directory.fileNames) != 0 {
		t.Fatalf("safe Services dependency accessed developer filenames: %#v", directory.fileNames)
	}

	loaded := (GitCloudConfig{Impl: DirConfig{Path: directory.dir}}).services(servicesDependencies{})
	services, err = loaded()
	if !reflect.DeepEqual(services, CloudServices{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe Services helper returned (%#v, %v), want (zero, exact %v)", services, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedServicesRejectsEmptyPopulations(t *testing.T) {
	if _, err := (&recordingServicesLoader{}).assertedDirectories(); err == nil {
		t.Fatal("empty recorded Services loader directory population passed")
	}
	if _, err := (&recordingServicesDirectory{}).assertedFileNames(); err == nil {
		t.Fatal("empty recorded Services filename population passed")
	}
}

func TestProjectDefaultsSelectsFileLoaderForProductionAndRequestsExactFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "project-defaults.json")
	writeProjectDefaultsFixture(t, path, representativeProjectDefaultsJSON)
	dependencies := systemProjectDefaultsDependencies()

	if dependencies.Loader == nil {
		t.Fatal("ProjectDefaults selected an incomplete loader dependency")
	}
	if reflect.TypeOf(dependencies.Loader) != reflect.TypeOf(fileProjectDefaultsLoader{}) {
		t.Fatalf("ProjectDefaults loader dependency is %T, want %T", dependencies.Loader, fileProjectDefaultsLoader{})
	}

	directory := &recordingProjectDefaultsDirectory{dir: root, path: path}
	defaults, err := dependencies.Load(directory)
	if err != nil {
		t.Fatalf("production ProjectDefaults loader returned an error: %v", err)
	}
	fileNames, populationErr := directory.assertedFileNames()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(fileNames, []string{"project-defaults.json"}) {
		t.Fatalf("production ProjectDefaults loader requested filenames %#v, want exact project-defaults.json", fileNames)
	}
	if !reflect.DeepEqual(defaults, representativeProjectDefaults()) {
		t.Fatalf("production ProjectDefaults loader decoded %#v, want representative complete value %#v", defaults, representativeProjectDefaults())
	}

	publicDefaults, err := (GitCloudConfig{Impl: DirConfig{Path: root}}).ProjectDefaults()
	if err != nil || !reflect.DeepEqual(publicDefaults, representativeProjectDefaults()) {
		t.Fatalf("public ProjectDefaults production wrapper returned (%#v, %v), want (%#v, nil)", publicDefaults, err, representativeProjectDefaults())
	}
}

func TestProjectDefaultsLoadsExactlyOncePerIndependentInvocationWithCompleteDirectory(t *testing.T) {
	loadError := errors.New("complete ProjectDefaults loader dependency error")
	first := CloudProjectDefaults{Type: "first complete defaults value"}
	second := CloudProjectDefaults{Type: "second complete partial defaults value"}
	recording := &recordingProjectDefaultsLoader{results: []projectDefaultsLoadResult{
		{defaults: first},
		{defaults: second, err: loadError},
	}}
	if len(recording.results) == 0 {
		t.Fatal("ProjectDefaults loader result population is empty")
	}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config receiver/root with spaces"}}

	actual, err := gitCfg.projectDefaults(recording.dependencies())
	if !reflect.DeepEqual(actual, first) || err != nil {
		t.Fatalf("first ProjectDefaults invocation returned (%#v, %v), want exact (%#v, nil)", actual, err, first)
	}
	if len(recording.directories) != 1 {
		t.Fatalf("first ProjectDefaults invocation loaded %d times, want exactly one", len(recording.directories))
	}

	actual, err = gitCfg.projectDefaults(recording.dependencies())
	if !reflect.DeepEqual(actual, second) || err != loadError {
		t.Fatalf("second ProjectDefaults invocation returned (%#v, %v), want exact (%#v, %v)", actual, err, second, loadError)
	}
	directories, populationErr := recording.assertedDirectories()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(directories) != 2 {
		t.Fatalf("two ProjectDefaults invocations loaded %d times, want one independent load each", len(directories))
	}
	for index, directory := range directories {
		if !reflect.DeepEqual(directory, gitCfg.Implementation()) {
			t.Fatalf("ProjectDefaults invocation %d received directory %#v, want complete implementation %#v", index+1, directory, gitCfg.Implementation())
		}
	}
}

func TestFileProjectDefaultsLoaderReturnsExactFilePathReadAndUnmarshalResults(t *testing.T) {
	pathError := errors.New("complete ProjectDefaults FilePath dependency error")
	directory := &recordingProjectDefaultsDirectory{
		dir:  "/developer/home/cloud-config/must-not-be-accessed",
		err:  pathError,
		path: "/must/not/be/read/project-defaults.json",
	}

	defaults, err := (fileProjectDefaultsLoader{}).Load(directory)
	if !reflect.DeepEqual(defaults, CloudProjectDefaults{}) || err != pathError {
		t.Fatalf("file ProjectDefaults loader returned (%#v, %v), want (zero, exact error %v)", defaults, err, pathError)
	}
	fileNames, populationErr := directory.assertedFileNames()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(fileNames, []string{"project-defaults.json"}) {
		t.Fatalf("file ProjectDefaults loader requested filenames %#v, want exact project-defaults.json", fileNames)
	}

	root := t.TempDir()
	missing := filepath.Join(root, "missing-project-defaults.json")
	defaults, err = (fileProjectDefaultsLoader{}).Load(&recordingProjectDefaultsDirectory{path: missing})
	wantReadError := "Unable to read " + missing + ", " + (&os.PathError{Op: "open", Path: missing, Err: syscall.ENOENT}).Error()
	if !reflect.DeepEqual(defaults, CloudProjectDefaults{}) || err == nil || err.Error() != wantReadError {
		t.Fatalf("file ProjectDefaults read result was (%#v, %v), want (zero, exact %q)", defaults, err, wantReadError)
	}

	malformed := filepath.Join(root, "malformed-project-defaults.json")
	writeProjectDefaultsFixture(t, malformed, `{"type":"partial defaults","settings":{"disableDependencySort":true,"maxSpringBootVersion":123}}`)
	defaults, err = (fileProjectDefaultsLoader{}).Load(&recordingProjectDefaultsDirectory{path: malformed})
	wantPartial := CloudProjectDefaults{Type: "partial defaults", Settings: ProjectSettings{DisableDependencySort: true}}
	wantUnmarshalError := "Unable to unmarshal " + malformed + ", json: cannot unmarshal number into Go struct field ProjectSettings.settings.maxSpringBootVersion of type string"
	if !reflect.DeepEqual(defaults, wantPartial) || err == nil || err.Error() != wantUnmarshalError {
		t.Fatalf("file ProjectDefaults unmarshal result was (%#v, %v), want exact partial (%#v, %q)", defaults, err, wantPartial, wantUnmarshalError)
	}
}

func TestProjectDefaultsDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	directory := &recordingProjectDefaultsDirectory{dir: "/developer/home/cloud-config/must-not-be-accessed"}

	defaults, err := (projectDefaultsDependencies{}).Load(directory)
	if !reflect.DeepEqual(defaults, CloudProjectDefaults{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe ProjectDefaults dependency default returned (%#v, %v), want (zero, exact %v)", defaults, err, filesystem.ErrNoFilesystem)
	}
	if len(directory.fileNames) != 0 {
		t.Fatalf("safe ProjectDefaults dependency accessed developer filenames: %#v", directory.fileNames)
	}

	defaults, err = (GitCloudConfig{Impl: DirConfig{Path: directory.dir}}).projectDefaults(projectDefaultsDependencies{})
	if !reflect.DeepEqual(defaults, CloudProjectDefaults{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe ProjectDefaults helper returned (%#v, %v), want (zero, exact %v)", defaults, err, filesystem.ErrNoFilesystem)
	}
	if len(directory.fileNames) != 0 {
		t.Fatalf("safe ProjectDefaults helper accessed developer filenames: %#v", directory.fileNames)
	}
}

func TestRecordedProjectDefaultsRejectsEmptyPopulations(t *testing.T) {
	if _, err := (&recordingProjectDefaultsLoader{}).assertedDirectories(); err == nil {
		t.Fatal("empty recorded ProjectDefaults loader directory population passed")
	}
	if _, err := (&recordingProjectDefaultsDirectory{}).assertedFileNames(); err == nil {
		t.Fatal("empty recorded ProjectDefaults filename population passed")
	}
}

const representativeProjectDefaultsJSON = `{
  "type": "complete project defaults",
  "settings": {
    "disableDependencySort": true,
    "disableSpringBootUpgrade": true,
    "disableKotlinUpgrade": true,
    "searchReplacer": true,
    "pomFileIndentation": "  ",
    "disableUpgradesFor": [{"groupId": "ignored.group", "artifactId": "ignored-artifact"}],
    "maxSpringBootVersion": "3.2.7",
    "maxVersionForDependencies": [{"groupId": "limited.group", "artifactId": "limited-artifact", "maxVersion": "4.5.6"}]
  }
}`

func representativeProjectDefaults() CloudProjectDefaults {
	return CloudProjectDefaults{
		Type: "complete project defaults",
		Settings: ProjectSettings{
			DisableDependencySort:    true,
			DisableSpringBootUpgrade: true,
			DisableKotlinUpgrade:     true,
			UseStealthMode:           true,
			PomFileIndentation:       "  ",
			DisableUpgradesFor: []Artifact{{
				GroupId: "ignored.group", ArtifactId: "ignored-artifact",
			}},
			MaxSpringBootVersion: "3.2.7",
			MaxVersionForDependencies: []MaxArtifact{{
				Artifact:   Artifact{GroupId: "limited.group", ArtifactId: "limited-artifact"},
				MaxVersion: "4.5.6",
			}},
		},
	}
}

func writeProjectDefaultsFixture(t *testing.T, path string, contents string) {
	t.Helper()
	if err := testutil.WriteFileOutsideWorkingTree(path, []byte(contents), 0600); err != nil {
		t.Fatalf("write disposable ProjectDefaults fixture: %v", err)
	}
}

func TestGitCloudConfig_Services(t *testing.T) {
	cfg := newMockCloudConfig()

	services, err := cfg.Services()()
	if err != nil {
		t.Fatalf("load services fixture: %v", err)
	}

	expected := "services"
	if services.Type != expected {
		t.Errorf("expected services type %s, got %s\n", expected, services.Type)
	}
}

func TestGitCloudConfig_LinkFromService(t *testing.T) {
	cfg := newMockCloudConfig()

	link, err := cfg.LinkFromService(cfg.Services(), "com.example", "flyway-demo", "info")
	if err != nil {
		t.Fatalf("resolve service link: %v", err)
	}

	expected := "http://localhost:8080/actuator/info"
	if link != expected {
		t.Errorf("expected link %s, got %s\n", expected, link)
	}
}

func TestGitCloudConfig_DefaultServiceEnvironmentUrl(t *testing.T) {
	cfg := newMockCloudConfig()

	services, err := cfg.Services()()
	if err != nil {
		t.Fatalf("load services fixture: %v", err)
	}
	key := "info"
	defaultUrl, err := cfg.DefaultServiceEnvironmentUrl(services.Data[0], key)
	if err != nil {
		t.Fatalf("resolve default environment URL: %v", err)
	}

	expected := "http://localhost:8080/actuator/info"
	if defaultUrl != expected {
		t.Errorf("expected default-url %s, got %s\n", expected, defaultUrl)
	}
}

func TestGitCloudConfigValidTemplatesFromReturnsExactOrderedPartialResultWithError(t *testing.T) {
	isolateConfigTestHome(t)
	cfg := newMockCloudConfig()

	templates, err := cfg.ValidTemplatesFrom([]string{"test-template", "test-template", "missing-template"})

	wantError := "could not find any valid templates with name: missing-template"
	if err == nil || err.Error() != wantError {
		t.Fatalf("ValidTemplatesFrom error was %v, want exact error %q", err, wantError)
	}
	if len(templates) != 1 {
		t.Fatalf("ValidTemplatesFrom returned %#v, want one populated partial template", templates)
	}
	if templates[0].Name != "test-template" ||
		templates[0].Project.Path != filepath.Join("test", "cloud-config", "templates", "test-template") ||
		templates[0].Project.Config.Name != "test-template" {
		t.Fatalf("ValidTemplatesFrom partial result was %#v, want exact first unique template", templates[0])
	}
}
