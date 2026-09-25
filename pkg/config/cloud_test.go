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

type deprecatedLoadResult struct {
	deprecated CloudDeprecated
	err        error
}

type recordingDeprecatedLoader struct {
	directories []Directory
	results     []deprecatedLoadResult
}

func (recording *recordingDeprecatedLoader) Load(directory Directory) (CloudDeprecated, error) {
	recording.directories = append(recording.directories, directory)
	index := len(recording.directories) - 1
	if index >= len(recording.results) {
		return CloudDeprecated{}, errors.New("recorded Deprecated result population is incomplete")
	}
	result := recording.results[index]
	return result.deprecated, result.err
}

func (recording *recordingDeprecatedLoader) dependencies() deprecatedDependencies {
	return deprecatedDependencies{Loader: recording}
}

func (recording *recordingDeprecatedLoader) assertedDirectories() ([]Directory, error) {
	if len(recording.directories) == 0 {
		return nil, errors.New("recorded Deprecated loader directory population is empty")
	}
	return recording.directories, nil
}

type recordingDeprecatedDirectory struct {
	dir       string
	path      string
	err       error
	fileNames []string
}

type globalConfigLoadResult struct {
	config GlobalCloudConfig
	err    error
}

type recordingGlobalConfigLoader struct {
	directories []Directory
	results     []globalConfigLoadResult
}

func (recording *recordingGlobalConfigLoader) Load(directory Directory) (GlobalCloudConfig, error) {
	recording.directories = append(recording.directories, directory)
	index := len(recording.directories) - 1
	if index >= len(recording.results) {
		return GlobalCloudConfig{}, errors.New("recorded GlobalCloudConfig result population is incomplete")
	}
	result := recording.results[index]
	return result.config, result.err
}

func (recording *recordingGlobalConfigLoader) dependencies() globalConfigDependencies {
	return globalConfigDependencies{Loader: recording}
}

func (recording *recordingGlobalConfigLoader) assertedDirectories() ([]Directory, error) {
	if len(recording.directories) == 0 {
		return nil, errors.New("recorded GlobalCloudConfig loader directory population is empty")
	}
	return recording.directories, nil
}

type recordingGlobalConfigDirectory struct {
	dir       string
	dirCalls  []string
	fileNames []string
}

func (directory *recordingGlobalConfigDirectory) Dir() string {
	directory.dirCalls = append(directory.dirCalls, directory.dir)
	return directory.dir
}

func (directory *recordingGlobalConfigDirectory) FilePath(fileName string) (string, error) {
	directory.fileNames = append(directory.fileNames, fileName)
	return "", errors.New("GlobalCloudConfig must not use Directory.FilePath")
}

func (directory *recordingGlobalConfigDirectory) assertedDirCalls() ([]string, error) {
	if len(directory.dirCalls) == 0 {
		return nil, errors.New("recorded GlobalCloudConfig Dir call population is empty")
	}
	return directory.dirCalls, nil
}

func (directory *recordingDeprecatedDirectory) Dir() string {
	return directory.dir
}

func (directory *recordingDeprecatedDirectory) FilePath(fileName string) (string, error) {
	directory.fileNames = append(directory.fileNames, fileName)
	return directory.path, directory.err
}

func (directory *recordingDeprecatedDirectory) assertedFileNames() ([]string, error) {
	if len(directory.fileNames) == 0 {
		return nil, errors.New("recorded Deprecated filename population is empty")
	}
	return directory.fileNames, nil
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

func TestDeprecatedSelectsFileLoaderForProductionAndRequestsExactFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "deprecated.json")
	writeDeprecatedFixture(t, path, representativeDeprecatedJSON)
	dependencies := systemDeprecatedDependencies()

	if dependencies.Loader == nil {
		t.Fatal("Deprecated selected an incomplete loader dependency")
	}
	if reflect.TypeOf(dependencies.Loader) != reflect.TypeOf(fileDeprecatedLoader{}) {
		t.Fatalf("Deprecated loader dependency is %T, want %T", dependencies.Loader, fileDeprecatedLoader{})
	}

	directory := &recordingDeprecatedDirectory{dir: root, path: path}
	deprecated, err := dependencies.Load(directory)
	if err != nil {
		t.Fatalf("production Deprecated loader returned an error: %v", err)
	}
	fileNames, populationErr := directory.assertedFileNames()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(fileNames, []string{"deprecated.json"}) {
		t.Fatalf("production Deprecated loader requested filenames %#v, want exact deprecated.json", fileNames)
	}
	if !reflect.DeepEqual(deprecated, representativeDeprecated()) {
		t.Fatalf("production Deprecated loader decoded %#v, want representative complete recursive value %#v", deprecated, representativeDeprecated())
	}

	publicDeprecated, err := (GitCloudConfig{Impl: DirConfig{Path: root}}).Deprecated()
	if err != nil || !reflect.DeepEqual(publicDeprecated, representativeDeprecated()) {
		t.Fatalf("public Deprecated production wrapper returned (%#v, %v), want (%#v, nil)", publicDeprecated, err, representativeDeprecated())
	}
}

func TestDeprecatedLoadsExactlyOncePerIndependentInvocationWithCompleteDirectory(t *testing.T) {
	loadError := errors.New("complete Deprecated loader dependency error")
	first := CloudDeprecated{Type: "first complete deprecated value"}
	second := CloudDeprecated{Type: "second complete partial deprecated value"}
	recording := &recordingDeprecatedLoader{results: []deprecatedLoadResult{
		{deprecated: first},
		{deprecated: second, err: loadError},
	}}
	if len(recording.results) == 0 {
		t.Fatal("Deprecated loader result population is empty")
	}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config receiver/root with spaces"}}

	actual, err := gitCfg.deprecated(recording.dependencies())
	if !reflect.DeepEqual(actual, first) || err != nil {
		t.Fatalf("first Deprecated invocation returned (%#v, %v), want exact (%#v, nil)", actual, err, first)
	}
	if len(recording.directories) != 1 {
		t.Fatalf("first Deprecated invocation loaded %d times, want exactly one", len(recording.directories))
	}

	actual, err = gitCfg.deprecated(recording.dependencies())
	if !reflect.DeepEqual(actual, second) || err != loadError {
		t.Fatalf("second Deprecated invocation returned (%#v, %v), want exact (%#v, %v)", actual, err, second, loadError)
	}
	directories, populationErr := recording.assertedDirectories()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(directories) != 2 {
		t.Fatalf("two Deprecated invocations loaded %d times, want one independent load each", len(directories))
	}
	for index, directory := range directories {
		if !reflect.DeepEqual(directory, gitCfg.Implementation()) {
			t.Fatalf("Deprecated invocation %d received directory %#v, want complete implementation %#v", index+1, directory, gitCfg.Implementation())
		}
	}
}

func TestFileDeprecatedLoaderReturnsExactFilePathReadAndUnmarshalResults(t *testing.T) {
	pathError := errors.New("complete Deprecated FilePath dependency error")
	directory := &recordingDeprecatedDirectory{
		dir:  "/developer/home/cloud-config/must-not-be-accessed",
		err:  pathError,
		path: "/must/not/be/read/deprecated.json",
	}

	deprecated, err := (fileDeprecatedLoader{}).Load(directory)
	if !reflect.DeepEqual(deprecated, CloudDeprecated{}) || err != pathError {
		t.Fatalf("file Deprecated loader returned (%#v, %v), want (zero, exact error %v)", deprecated, err, pathError)
	}
	fileNames, populationErr := directory.assertedFileNames()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(fileNames, []string{"deprecated.json"}) {
		t.Fatalf("file Deprecated loader requested filenames %#v, want exact deprecated.json", fileNames)
	}

	root := t.TempDir()
	missing := filepath.Join(root, "missing-deprecated.json")
	deprecated, err = (fileDeprecatedLoader{}).Load(&recordingDeprecatedDirectory{path: missing})
	wantReadError := "Unable to read " + missing + ", " + (&os.PathError{Op: "open", Path: missing, Err: syscall.ENOENT}).Error()
	if !reflect.DeepEqual(deprecated, CloudDeprecated{}) || err == nil || err.Error() != wantReadError {
		t.Fatalf("file Deprecated read result was (%#v, %v), want (zero, exact %q)", deprecated, err, wantReadError)
	}

	malformed := filepath.Join(root, "malformed-deprecated.json")
	writeDeprecatedFixture(t, malformed, `{"type":"partial deprecated","data":{"dependencies":[{"groupId":"partial.group","artifactId":123}]}}`)
	deprecated, err = (fileDeprecatedLoader{}).Load(&recordingDeprecatedDirectory{path: malformed})
	wantPartial := CloudDeprecated{Type: "partial deprecated"}
	wantPartial.Data.Dependencies = []CloudDeprecatedDependency{{GroupId: "partial.group"}}
	wantUnmarshalError := "Unable to unmarshal " + malformed + ", json: cannot unmarshal number into Go struct field CloudDeprecatedDependency.data.dependencies.artifactId of type string"
	if !reflect.DeepEqual(deprecated, wantPartial) || err == nil || err.Error() != wantUnmarshalError {
		t.Fatalf("file Deprecated unmarshal result was (%#v, %v), want exact partial (%#v, %q)", deprecated, err, wantPartial, wantUnmarshalError)
	}
}

func TestDeprecatedDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	directory := &recordingDeprecatedDirectory{dir: "/developer/home/cloud-config/must-not-be-accessed"}

	deprecated, err := (deprecatedDependencies{}).Load(directory)
	if !reflect.DeepEqual(deprecated, CloudDeprecated{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe Deprecated dependency default returned (%#v, %v), want (zero, exact %v)", deprecated, err, filesystem.ErrNoFilesystem)
	}
	if len(directory.fileNames) != 0 {
		t.Fatalf("safe Deprecated dependency accessed developer filenames: %#v", directory.fileNames)
	}

	deprecated, err = (GitCloudConfig{Impl: DirConfig{Path: directory.dir}}).deprecated(deprecatedDependencies{})
	if !reflect.DeepEqual(deprecated, CloudDeprecated{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe Deprecated helper returned (%#v, %v), want (zero, exact %v)", deprecated, err, filesystem.ErrNoFilesystem)
	}
	if len(directory.fileNames) != 0 {
		t.Fatalf("safe Deprecated helper accessed developer filenames: %#v", directory.fileNames)
	}
}

func TestRecordedDeprecatedRejectsEmptyPopulations(t *testing.T) {
	if _, err := (&recordingDeprecatedLoader{}).assertedDirectories(); err == nil {
		t.Fatal("empty recorded Deprecated loader directory population passed")
	}
	if _, err := (&recordingDeprecatedDirectory{}).assertedFileNames(); err == nil {
		t.Fatal("empty recorded Deprecated filename population passed")
	}
}

func TestGlobalCloudConfigSelectsFileLoaderForProductionAndUsesDirectEnvironmentExpandedYAMLPath(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "global-config.yaml")
	t.Setenv("PLY_TEST_GLOBAL_CONFIG_DOCUMENT_KEY", "cloudConfigSource")
	t.Setenv("PLY_TEST_GLOBAL_CONFIG_ROOT_KEY", "rootUrl")
	t.Setenv("PLY_TEST_GLOBAL_CONFIG_ROOT", "https://expanded.example/cloud")
	t.Setenv("PLY_TEST_GLOBAL_CONFIG_RELATIVE", "/expanded/raw")
	writeGlobalConfigFixture(t, path, `${PLY_TEST_GLOBAL_CONFIG_DOCUMENT_KEY}:
  ${PLY_TEST_GLOBAL_CONFIG_ROOT_KEY}: "${PLY_TEST_GLOBAL_CONFIG_ROOT}"
  relativFileUrl: "${PLY_TEST_GLOBAL_CONFIG_RELATIVE}"
`)
	dependencies := systemGlobalConfigDependencies()

	if dependencies.Loader == nil {
		t.Fatal("GlobalCloudConfig selected an incomplete loader dependency")
	}
	if reflect.TypeOf(dependencies.Loader) != reflect.TypeOf(fileGlobalConfigLoader{}) {
		t.Fatalf("GlobalCloudConfig loader dependency is %T, want %T", dependencies.Loader, fileGlobalConfigLoader{})
	}

	directory := &recordingGlobalConfigDirectory{dir: root}
	config, err := dependencies.Load(directory)
	if err != nil {
		t.Fatalf("production GlobalCloudConfig loader returned an error: %v", err)
	}
	dirCalls, populationErr := directory.assertedDirCalls()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(dirCalls, []string{root}) {
		t.Fatalf("production GlobalCloudConfig loader requested directories %#v, want exact root %q once", dirCalls, root)
	}
	if len(directory.fileNames) != 0 {
		t.Fatalf("production GlobalCloudConfig loader called FilePath with %#v, want direct Dir path composition", directory.fileNames)
	}
	want := GlobalCloudConfig{CloudConfigSource: CloudConfigSource{
		RootUrl:        "https://expanded.example/cloud",
		RelativFileUrl: "/expanded/raw",
	}}
	if !reflect.DeepEqual(config, want) {
		t.Fatalf("production GlobalCloudConfig loader decoded %#v, want exact expanded value %#v", config, want)
	}

	publicConfig, err := (GitCloudConfig{Impl: DirConfig{Path: root}}).GlobalCloudConfig()
	if err != nil || !reflect.DeepEqual(publicConfig, want) {
		t.Fatalf("public GlobalCloudConfig production wrapper returned (%#v, %v), want (%#v, nil)", publicConfig, err, want)
	}
}

func TestGlobalCloudConfigLoadsExactlyOncePerIndependentInvocationWithCompleteDirectory(t *testing.T) {
	loadError := errors.New("complete GlobalCloudConfig loader dependency error")
	first := GlobalCloudConfig{CloudConfigSource: CloudConfigSource{RootUrl: "first complete config value"}}
	second := GlobalCloudConfig{CloudConfigSource: CloudConfigSource{RootUrl: "second complete partial config value"}}
	recording := &recordingGlobalConfigLoader{results: []globalConfigLoadResult{
		{config: first},
		{config: second, err: loadError},
	}}
	if len(recording.results) == 0 {
		t.Fatal("GlobalCloudConfig loader result population is empty")
	}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config receiver/root with spaces"}}

	actual, err := gitCfg.globalCloudConfig(recording.dependencies())
	if !reflect.DeepEqual(actual, first) || err != nil {
		t.Fatalf("first GlobalCloudConfig invocation returned (%#v, %v), want exact (%#v, nil)", actual, err, first)
	}
	if len(recording.directories) != 1 {
		t.Fatalf("first GlobalCloudConfig invocation loaded %d times, want exactly one", len(recording.directories))
	}

	actual, err = gitCfg.globalCloudConfig(recording.dependencies())
	if !reflect.DeepEqual(actual, second) || err != loadError {
		t.Fatalf("second GlobalCloudConfig invocation returned (%#v, %v), want exact (%#v, %v)", actual, err, second, loadError)
	}
	directories, populationErr := recording.assertedDirectories()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if len(directories) != 2 {
		t.Fatalf("two GlobalCloudConfig invocations loaded %d times, want one independent load each", len(directories))
	}
	for index, directory := range directories {
		if !reflect.DeepEqual(directory, gitCfg.Implementation()) {
			t.Fatalf("GlobalCloudConfig invocation %d received directory %#v, want complete implementation %#v", index+1, directory, gitCfg.Implementation())
		}
	}
}

func TestFileGlobalConfigLoaderReturnsExactRawReadAndPartialYAMLResults(t *testing.T) {
	root := t.TempDir()
	directory := &recordingGlobalConfigDirectory{dir: root}
	config, err := (fileGlobalConfigLoader{}).Load(directory)
	missing := filepath.Join(root, "global-config.yaml")
	wantReadError := &os.PathError{Op: "open", Path: missing, Err: syscall.ENOENT}
	if !reflect.DeepEqual(config, GlobalCloudConfig{}) || !reflect.DeepEqual(err, wantReadError) {
		t.Fatalf("file GlobalCloudConfig read result was (%#v, %#v), want (zero, exact raw %#v)", config, err, wantReadError)
	}
	dirCalls, populationErr := directory.assertedDirCalls()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(dirCalls, []string{root}) {
		t.Fatalf("file GlobalCloudConfig loader requested directories %#v, want exact root %q once", dirCalls, root)
	}
	if len(directory.fileNames) != 0 {
		t.Fatalf("file GlobalCloudConfig loader called FilePath with %#v, want no calls", directory.fileNames)
	}

	writeGlobalConfigFixture(t, missing, `cloudConfigSource:
  rootUrl: "https://partial.example/root"
  relativFileUrl:
    nested: invalid
`)
	config, err = (fileGlobalConfigLoader{}).Load(&recordingGlobalConfigDirectory{dir: root})
	wantPartial := GlobalCloudConfig{CloudConfigSource: CloudConfigSource{RootUrl: "https://partial.example/root"}}
	wantUnmarshalError := "yaml: unmarshal errors:\n  line 4: cannot unmarshal !!map into string"
	if !reflect.DeepEqual(config, wantPartial) || err == nil || err.Error() != wantUnmarshalError {
		t.Fatalf("file GlobalCloudConfig YAML result was (%#v, %v), want exact partial (%#v, %q)", config, err, wantPartial, wantUnmarshalError)
	}
}

func TestGlobalCloudConfigRereadsAndReexpandsEachPublicInvocation(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "global-config.yaml")
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: root}}
	t.Setenv("PLY_TEST_GLOBAL_CONFIG_DYNAMIC", "first-environment")
	writeGlobalConfigFixture(t, path, `cloudConfigSource:
  rootUrl: "${PLY_TEST_GLOBAL_CONFIG_DYNAMIC}/first-file"
  relativFileUrl: "/first-relative"
`)

	first, err := gitCfg.GlobalCloudConfig()
	wantFirst := GlobalCloudConfig{CloudConfigSource: CloudConfigSource{
		RootUrl:        "first-environment/first-file",
		RelativFileUrl: "/first-relative",
	}}
	if err != nil || !reflect.DeepEqual(first, wantFirst) {
		t.Fatalf("first public GlobalCloudConfig invocation returned (%#v, %v), want (%#v, nil)", first, err, wantFirst)
	}

	t.Setenv("PLY_TEST_GLOBAL_CONFIG_DYNAMIC", "second-environment")
	writeGlobalConfigFixture(t, path, `cloudConfigSource:
  rootUrl: "${PLY_TEST_GLOBAL_CONFIG_DYNAMIC}/second-file"
  relativFileUrl: "/second-relative"
`)
	second, err := gitCfg.GlobalCloudConfig()
	wantSecond := GlobalCloudConfig{CloudConfigSource: CloudConfigSource{
		RootUrl:        "second-environment/second-file",
		RelativFileUrl: "/second-relative",
	}}
	if err != nil || !reflect.DeepEqual(second, wantSecond) {
		t.Fatalf("second public GlobalCloudConfig invocation returned (%#v, %v), want fresh (%#v, nil)", second, err, wantSecond)
	}
}

func TestGlobalConfigDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	directory := &recordingGlobalConfigDirectory{dir: "/developer/home/cloud-config/must-not-be-accessed"}

	config, err := (globalConfigDependencies{}).Load(directory)
	if !reflect.DeepEqual(config, GlobalCloudConfig{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe GlobalCloudConfig dependency default returned (%#v, %v), want (zero, exact %v)", config, err, filesystem.ErrNoFilesystem)
	}
	if len(directory.dirCalls) != 0 || len(directory.fileNames) != 0 {
		t.Fatalf("safe GlobalCloudConfig dependency accessed developer directory: Dir=%#v FilePath=%#v", directory.dirCalls, directory.fileNames)
	}

	config, err = (GitCloudConfig{Impl: DirConfig{Path: directory.dir}}).globalCloudConfig(globalConfigDependencies{})
	if !reflect.DeepEqual(config, GlobalCloudConfig{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe GlobalCloudConfig helper returned (%#v, %v), want (zero, exact %v)", config, err, filesystem.ErrNoFilesystem)
	}
	if len(directory.dirCalls) != 0 || len(directory.fileNames) != 0 {
		t.Fatalf("safe GlobalCloudConfig helper accessed developer directory: Dir=%#v FilePath=%#v", directory.dirCalls, directory.fileNames)
	}
}

func TestGlobalCloudConfigSourceForPreservesExactSlashFormatting(t *testing.T) {
	config := GlobalCloudConfig{CloudConfigSource: CloudConfigSource{
		RootUrl:        "https://source.example/base/",
		RelativFileUrl: "/raw/content",
	}}

	actual := config.SourceFor("tips/advanced", "README.md")
	want := "https://source.example/base//raw/content/tips/advanced/README.md"
	if actual != want {
		t.Fatalf("GlobalCloudConfig SourceFor returned %q, want exact unchanged formatting %q", actual, want)
	}
}

func TestRecordedGlobalCloudConfigRejectsEmptyPopulations(t *testing.T) {
	if _, err := (&recordingGlobalConfigLoader{}).assertedDirectories(); err == nil {
		t.Fatal("empty recorded GlobalCloudConfig loader directory population passed")
	}
	if _, err := (&recordingGlobalConfigDirectory{}).assertedDirCalls(); err == nil {
		t.Fatal("empty recorded GlobalCloudConfig Dir call population passed")
	}
}

func writeGlobalConfigFixture(t *testing.T, path string, contents string) {
	t.Helper()
	if err := testutil.WriteFileOutsideWorkingTree(path, []byte(contents), 0600); err != nil {
		t.Fatalf("write disposable GlobalCloudConfig fixture: %v", err)
	}
}

const representativeDeprecatedJSON = `{
  "type": "complete deprecated",
  "data": {
    "dependencies": [{
      "groupId": "root.group",
      "artifactId": "root-artifact",
      "files": ["root-one", "root-two"],
      "associated": {
        "files": ["root-associated-file"],
        "dependencies": [{
          "groupId": "associated.group",
          "artifactId": "associated-artifact",
          "files": ["associated-file"],
          "associated": {
            "files": ["nested-associated-file"],
            "dependencies": [{
              "groupId": "nested.group",
              "artifactId": "nested-artifact",
              "files": ["nested-file"],
              "associated": {"files": ["leaf-file"], "dependencies": []},
              "replacement_templates": ["nested-template"]
            }]
          },
          "replacement_templates": ["associated-template"]
        }]
      },
      "replacement_templates": ["first-template", "second-template"]
    }, {
      "groupId": "nil.group",
      "artifactId": "nil-artifact"
    }]
  }
}`

func representativeDeprecated() CloudDeprecated {
	nested := CloudDeprecatedDependency{
		GroupId:              "nested.group",
		ArtifactId:           "nested-artifact",
		Files:                []string{"nested-file"},
		ReplacementTemplates: []string{"nested-template"},
	}
	nested.Associated.Files = []string{"leaf-file"}
	nested.Associated.Dependencies = []CloudDeprecatedDependency{}
	associated := CloudDeprecatedDependency{
		GroupId:              "associated.group",
		ArtifactId:           "associated-artifact",
		Files:                []string{"associated-file"},
		ReplacementTemplates: []string{"associated-template"},
	}
	associated.Associated.Files = []string{"nested-associated-file"}
	associated.Associated.Dependencies = []CloudDeprecatedDependency{nested}
	root := CloudDeprecatedDependency{
		GroupId:              "root.group",
		ArtifactId:           "root-artifact",
		Files:                []string{"root-one", "root-two"},
		ReplacementTemplates: []string{"first-template", "second-template"},
	}
	root.Associated.Files = []string{"root-associated-file"}
	root.Associated.Dependencies = []CloudDeprecatedDependency{associated}
	deprecated := CloudDeprecated{Type: "complete deprecated"}
	deprecated.Data.Dependencies = []CloudDeprecatedDependency{
		root,
		{GroupId: "nil.group", ArtifactId: "nil-artifact"},
	}
	return deprecated
}

func writeDeprecatedFixture(t *testing.T, path string, contents string) {
	t.Helper()
	if err := testutil.WriteFileOutsideWorkingTree(path, []byte(contents), 0600); err != nil {
		t.Fatalf("write disposable Deprecated fixture: %v", err)
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
