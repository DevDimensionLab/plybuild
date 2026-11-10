package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/file"
)

type perNameTemplateLoadResult struct {
	template CloudTemplate
	err      error
}

type recordedPerNameTemplateLoad struct {
	receiver GitCloudConfig
	name     string
}

type recordingPerNameTemplateLoader struct {
	loads   []recordedPerNameTemplateLoad
	results []perNameTemplateLoadResult
}

func (recording *recordingPerNameTemplateLoader) Load(gitCfg GitCloudConfig, name string) (CloudTemplate, error) {
	recording.loads = append(recording.loads, recordedPerNameTemplateLoad{receiver: gitCfg, name: name})
	index := len(recording.loads) - 1
	if index >= len(recording.results) {
		return CloudTemplate{}, errors.New("recorded per-name Template result population is incomplete")
	}
	result := recording.results[index]
	return result.template, result.err
}

func (recording *recordingPerNameTemplateLoader) dependencies() validTemplatesDependencies {
	return validTemplatesDependencies{Loader: recording}
}

func (recording *recordingPerNameTemplateLoader) assertedLoads() ([]recordedPerNameTemplateLoad, error) {
	if len(recording.loads) == 0 {
		return nil, errors.New("recorded per-name Template load population is empty")
	}
	return append([]recordedPerNameTemplateLoad(nil), recording.loads...), nil
}

type validTemplatesCloudConfigIdentity struct {
	GitCloudConfig
}

func TestValidTemplatesFromSelectsGitPerNameTemplateLoaderForProduction(t *testing.T) {
	dependencies := systemValidTemplatesDependencies()
	if dependencies.Loader == nil {
		t.Fatal("ValidTemplatesFrom selected an incomplete per-name Template loader dependency")
	}
	if reflect.TypeOf(dependencies.Loader) != reflect.TypeOf(gitPerNameTemplateLoader{}) {
		t.Fatalf("ValidTemplatesFrom per-name Template loader dependency is %T, want %T", dependencies.Loader, gitPerNameTemplateLoader{})
	}

	isolateConfigTestHome(t)
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: file.Path("test/cloud-config")}}
	actual, err := dependencies.Load(gitCfg, "test-template")
	if err != nil {
		t.Fatalf("production per-name Template loader returned an error: %v", err)
	}
	wantDirectory := file.Path("test/cloud-config/templates/test-template")
	if actual.Name != "test-template" || actual.Project.Path != wantDirectory ||
		actual.Project.ConfigFile != file.Path("%s/%s", wantDirectory, projectConfigFileName) ||
		actual.Project.Config.Name != "test-template" || !actual.Project.IsMavenProject() ||
		actual.Project.Type == nil || actual.Project.CloudConfig == nil {
		t.Fatalf("production per-name Template loader returned %#v, want complete tracked template", actual)
	}
}

func TestPerNameTemplateLoaderDeliversCompleteReceiverAndNameOnceWithExactValueAndError(t *testing.T) {
	loadError := errors.New("complete per-name Template loader dependency error")
	cloud := &validTemplatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/complete template cloud identity"}}}
	projectType := &MavenProject{PomFile: "/complete template pom identity"}
	want := CloudTemplate{
		Name: "complete template value",
		Project: Project{
			Path:        "/complete template project",
			GitInfo:     GitInfo{IsRepo: true, IsDirty: true, EnableCommit: true},
			ConfigFile:  "/complete template project/ply.json",
			Config:      ProjectConfiguration{Name: "complete template project", Profile: "complete profile"},
			Type:        projectType,
			CloudConfig: cloud,
		},
	}
	recording := &recordingPerNameTemplateLoader{results: []perNameTemplateLoadResult{{template: want, err: loadError}}}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete arbitrary receiver/root with spaces"}}
	name := "complete arbitrary template/name with spaces"
	dependencies := recording.dependencies()
	if dependencies.Loader != recording {
		t.Fatalf("ValidTemplatesFrom dependency lost its complete loader value: %#v", dependencies)
	}

	actual, err := dependencies.Load(gitCfg, name)

	if !reflect.DeepEqual(actual, want) || err != loadError {
		t.Fatalf("per-name Template loader returned (%#v, %v), want exact (%#v, %v)", actual, err, want, loadError)
	}
	if actual.Project.CloudConfig != cloud || actual.Project.Type != projectType {
		t.Fatalf("per-name Template loader lost embedded identities: CloudConfig=%T %#v Type=%T %#v",
			actual.Project.CloudConfig, actual.Project.CloudConfig, actual.Project.Type, actual.Project.Type)
	}
	loads, populationErr := recording.assertedLoads()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(loads, []recordedPerNameTemplateLoad{{receiver: gitCfg, name: name}}) {
		t.Fatalf("per-name Template loads were %#v, want one complete receiver/name delivery", loads)
	}
}

func TestValidTemplatesFromPreservesFirstOccurrenceOrderAndCaseSensitiveDuplicateSuppression(t *testing.T) {
	firstCloud := &validTemplatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/first cloud identity"}}}
	firstType := &MavenProject{PomFile: "/first pom identity"}
	want := []CloudTemplate{
		{Name: "first", Project: Project{Path: "/first project", Type: firstType, CloudConfig: firstCloud}},
		{Name: "", Project: Project{Path: "/empty-name project"}},
		{Name: "First", Project: Project{Path: "/case-distinct project"}},
		{Name: "second", Project: Project{Path: "/second project"}},
	}
	recording := &recordingPerNameTemplateLoader{results: []perNameTemplateLoadResult{
		{template: want[0]},
		{template: want[1]},
		{template: want[2]},
		{template: want[3]},
	}}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete ordered receiver"}}

	actual, err := gitCfg.validTemplatesFrom(recording.dependencies(), []string{
		"first", "first", "", "", "First", "second", "first", "second",
	})

	if err != nil || !reflect.DeepEqual(actual, want) {
		t.Fatalf("ValidTemplatesFrom ordered result was (%#v, %v), want exact (%#v, nil)", actual, err, want)
	}
	if actual[0].Project.CloudConfig != firstCloud || actual[0].Project.Type != firstType {
		t.Fatalf("ValidTemplatesFrom lost embedded first-template identities: CloudConfig=%T %#v Type=%T %#v",
			actual[0].Project.CloudConfig, actual[0].Project.CloudConfig, actual[0].Project.Type, actual[0].Project.Type)
	}
	loads, populationErr := recording.assertedLoads()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantLoads := []recordedPerNameTemplateLoad{
		{receiver: gitCfg, name: "first"},
		{receiver: gitCfg, name: ""},
		{receiver: gitCfg, name: "First"},
		{receiver: gitCfg, name: "second"},
	}
	if !reflect.DeepEqual(loads, wantLoads) {
		t.Fatalf("ValidTemplatesFrom per-name loads were %#v, want first-occurrence order %#v", loads, wantLoads)
	}
}

func TestValidTemplatesFromStopsAtFirstErrorWithExactPriorPartialResultAndDiscardsFailingValue(t *testing.T) {
	loadError := errors.New("complete per-name Template failure")
	firstCloud := &validTemplatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/partial cloud identity"}}}
	first := CloudTemplate{Name: "first-template", Project: Project{Path: "/first partial project", CloudConfig: firstCloud}}
	failing := CloudTemplate{Name: "failing value must be discarded", Project: Project{Path: "/failing project must be discarded"}}
	recording := &recordingPerNameTemplateLoader{results: []perNameTemplateLoadResult{
		{template: first},
		{template: failing, err: loadError},
		{template: CloudTemplate{Name: "later template must not load"}},
	}}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete partial-result receiver"}}

	actual, err := gitCfg.validTemplatesFrom(recording.dependencies(), []string{
		"first-template", "first-template", "failing-template", "later-template",
	})

	if err != loadError || !reflect.DeepEqual(actual, []CloudTemplate{first}) {
		t.Fatalf("ValidTemplatesFrom failure result was (%#v, %v), want exact prior partial result and error (%#v, %v)",
			actual, err, []CloudTemplate{first}, loadError)
	}
	if actual[0].Project.CloudConfig != firstCloud {
		t.Fatalf("ValidTemplatesFrom partial result lost embedded cloud identity: %T %#v", actual[0].Project.CloudConfig, actual[0].Project.CloudConfig)
	}
	loads, populationErr := recording.assertedLoads()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantLoads := []recordedPerNameTemplateLoad{
		{receiver: gitCfg, name: "first-template"},
		{receiver: gitCfg, name: "failing-template"},
	}
	if !reflect.DeepEqual(loads, wantLoads) {
		t.Fatalf("ValidTemplatesFrom failure loads were %#v, want stop before later name %#v", loads, wantLoads)
	}
}

func TestValidTemplatesFromReturnsNilWithoutLoadingForNilAndEmptyInputs(t *testing.T) {
	tests := []struct {
		name string
		list []string
	}{
		{name: "nil", list: nil},
		{name: "empty", list: []string{}},
	}
	if len(tests) == 0 {
		t.Fatal("ValidTemplatesFrom nil/empty test-case population is empty")
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingPerNameTemplateLoader{}
			actual, err := (GitCloudConfig{Impl: DirConfig{Path: "/must not be accessed"}}).
				validTemplatesFrom(recording.dependencies(), test.list)

			if actual != nil || err != nil {
				t.Fatalf("ValidTemplatesFrom %s result was (%#v, %v), want (nil, nil)", test.name, actual, err)
			}
			if len(recording.loads) != 0 {
				t.Fatalf("ValidTemplatesFrom %s input loaded templates: %#v", test.name, recording.loads)
			}
		})
	}
}

func TestValidTemplatesFromLoadsIndependentlyAcrossRepeatedInvocations(t *testing.T) {
	firstCloud := &validTemplatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/first repeated cloud identity"}}}
	secondCloud := &validTemplatesCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/second repeated cloud identity"}}}
	first := CloudTemplate{Name: "repeated-template", Project: Project{Path: "/first repeated project", CloudConfig: firstCloud}}
	second := CloudTemplate{Name: "repeated-template", Project: Project{Path: "/second repeated project", CloudConfig: secondCloud}}
	recording := &recordingPerNameTemplateLoader{results: []perNameTemplateLoadResult{
		{template: first},
		{template: second},
	}}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete repeated receiver"}}
	dependencies := recording.dependencies()

	firstActual, err := gitCfg.validTemplatesFrom(dependencies, []string{"repeated-template", "repeated-template"})
	if err != nil || !reflect.DeepEqual(firstActual, []CloudTemplate{first}) || firstActual[0].Project.CloudConfig != firstCloud {
		t.Fatalf("first ValidTemplatesFrom invocation returned (%#v, %v), want fresh exact first result", firstActual, err)
	}
	secondActual, err := gitCfg.validTemplatesFrom(dependencies, []string{"repeated-template", "repeated-template"})
	if err != nil || !reflect.DeepEqual(secondActual, []CloudTemplate{second}) || secondActual[0].Project.CloudConfig != secondCloud {
		t.Fatalf("second ValidTemplatesFrom invocation returned (%#v, %v), want independent exact second result", secondActual, err)
	}
	loads, populationErr := recording.assertedLoads()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantLoads := []recordedPerNameTemplateLoad{
		{receiver: gitCfg, name: "repeated-template"},
		{receiver: gitCfg, name: "repeated-template"},
	}
	if !reflect.DeepEqual(loads, wantLoads) {
		t.Fatalf("repeated ValidTemplatesFrom loads were %#v, want one fresh load per invocation %#v", loads, wantLoads)
	}
}

func TestValidTemplatesDependenciesDefaultToSafeExactErrorWithoutPathAccess(t *testing.T) {
	path := filepath.Join(t.TempDir(), "developer-home-cloud-config-must-not-be-accessed")
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: path}}

	template, err := (validTemplatesDependencies{}).Load(gitCfg, "must-not-load")
	if !reflect.DeepEqual(template, CloudTemplate{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe per-name Template dependency default returned (%#v, %v), want (zero, exact %v)",
			template, err, filesystem.ErrNoFilesystem)
	}
	actual, err := gitCfg.validTemplatesFrom(validTemplatesDependencies{}, []string{"must-not-load"})
	if actual != nil || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe ValidTemplatesFrom helper returned (%#v, %v), want (nil, exact %v)", actual, err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Lstat(path); !os.IsNotExist(statErr) {
		t.Fatalf("safe ValidTemplatesFrom dependency touched receiver path %q: %v", path, statErr)
	}
}

func TestRecordedPerNameTemplateLoaderRejectsEmptyPopulations(t *testing.T) {
	if _, err := (&recordingPerNameTemplateLoader{}).assertedLoads(); err == nil {
		t.Fatal("empty recorded per-name Template load population passed")
	}
}
