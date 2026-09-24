package config

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
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
