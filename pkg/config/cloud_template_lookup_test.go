package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/file"
)

type templateListLoadResult struct {
	templates []CloudTemplate
	err       error
}

type recordingTemplateListLoader struct {
	receivers []GitCloudConfig
	results   []templateListLoadResult
}

func (recording *recordingTemplateListLoader) Load(gitCfg GitCloudConfig) ([]CloudTemplate, error) {
	recording.receivers = append(recording.receivers, gitCfg)
	index := len(recording.receivers) - 1
	if index >= len(recording.results) {
		return nil, errors.New("recorded Template list result population is incomplete")
	}
	result := recording.results[index]
	return result.templates, result.err
}

func (recording *recordingTemplateListLoader) dependencies() templateLookupDependencies {
	return templateLookupDependencies{Loader: recording}
}

func (recording *recordingTemplateListLoader) assertedReceivers() ([]GitCloudConfig, error) {
	if len(recording.receivers) == 0 {
		return nil, errors.New("recorded Template list-loader receiver population is empty")
	}
	return append([]GitCloudConfig(nil), recording.receivers...), nil
}

type templateLookupCloudConfigIdentity struct {
	GitCloudConfig
}

func TestTemplateSelectsGitTemplateListLoaderForProductionAndRetainsTrackedLookups(t *testing.T) {
	dependencies := systemTemplateLookupDependencies()
	if dependencies.Loader == nil {
		t.Fatal("Template selected an incomplete list-loader dependency")
	}
	if reflect.TypeOf(dependencies.Loader) != reflect.TypeOf(gitTemplateListLoader{}) {
		t.Fatalf("Template list-loader dependency is %T, want %T", dependencies.Loader, gitTemplateListLoader{})
	}

	isolateConfigTestHome(t)
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: file.Path("test/cloud-config")}}
	tests := []struct {
		name            string
		lookup          string
		wantConfigFile  string
		wantProjectName string
		wantMaven       bool
	}{
		{
			name:            "current project",
			lookup:          "test-template",
			wantConfigFile:  projectConfigFileName,
			wantProjectName: "test-template",
			wantMaven:       true,
		},
		{
			name:            "legacy project",
			lookup:          file.Path("legacy-category/legacy-template"),
			wantConfigFile:  legacyProjectConfigFileName,
			wantProjectName: "legacy-template",
		},
	}
	if len(tests) == 0 {
		t.Fatal("production Template lookup test-case population is empty")
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := gitCfg.Template(test.lookup)
			if err != nil {
				t.Fatalf("production Template lookup returned an error: %v", err)
			}
			wantDirectory := file.Path("test/cloud-config/templates/%s", test.lookup)
			if actual.Name != test.lookup || actual.Project.Path != wantDirectory ||
				actual.Project.ConfigFile != file.Path("%s/%s", wantDirectory, test.wantConfigFile) ||
				actual.Project.Config.Name != test.wantProjectName ||
				actual.Project.IsMavenProject() != test.wantMaven || actual.Project.CloudConfig == nil {
				t.Fatalf("production Template lookup returned %#v, want complete tracked %q project", actual, test.lookup)
			}
		})
	}
}

func TestTemplateListLoaderDeliversCompleteReceiverOnceAndPreservesListAndExactError(t *testing.T) {
	listError := errors.New("complete Template list-loader dependency error")
	cloud := &templateLookupCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/complete list cloud identity"}}}
	projectType := &MavenProject{PomFile: "/complete list pom identity"}
	want := []CloudTemplate{{
		Name: "complete list template",
		Project: Project{
			Path:        "/complete list project",
			ConfigFile:  "/complete list project/ply.json",
			Config:      ProjectConfiguration{Name: "complete list project", Profile: "complete profile"},
			Type:        projectType,
			CloudConfig: cloud,
		},
	}}
	recording := &recordingTemplateListLoader{results: []templateListLoadResult{{templates: want, err: listError}}}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete arbitrary receiver/root with spaces"}}
	dependencies := recording.dependencies()
	if dependencies.Loader != recording {
		t.Fatalf("Template list dependency lost its complete loader value: %#v", dependencies)
	}

	actual, err := dependencies.Load(gitCfg)

	if !reflect.DeepEqual(actual, want) || err != listError {
		t.Fatalf("Template list loader returned (%#v, %v), want exact (%#v, %v)", actual, err, want, listError)
	}
	if actual[0].Project.CloudConfig != cloud || actual[0].Project.Type != projectType {
		t.Fatalf("Template list loader lost embedded identities: CloudConfig=%T %#v Type=%T %#v",
			actual[0].Project.CloudConfig, actual[0].Project.CloudConfig, actual[0].Project.Type, actual[0].Project.Type)
	}
	receivers, populationErr := recording.assertedReceivers()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(receivers, []GitCloudConfig{gitCfg}) {
		t.Fatalf("Template list-loader receivers were %#v, want one complete receiver %#v", receivers, []GitCloudConfig{gitCfg})
	}
}

func TestTemplateLookupReturnsListErrorBeforeMatchingPartialResult(t *testing.T) {
	listError := errors.New("complete eager Template list error")
	matchingPartial := []CloudTemplate{{
		Name:    "matching-template",
		Project: Project{Path: "/matching partial project", Config: ProjectConfiguration{Name: "matching partial"}},
	}}
	recording := &recordingTemplateListLoader{results: []templateListLoadResult{{templates: matchingPartial, err: listError}}}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete list-error receiver"}}

	actual, err := gitCfg.template(recording.dependencies(), "matching-template")

	if !reflect.DeepEqual(actual, CloudTemplate{}) || err != listError {
		t.Fatalf("Template eager list-error result was (%#v, %v), want (zero, exact %v)", actual, err, listError)
	}
	receivers, populationErr := recording.assertedReceivers()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(receivers, []GitCloudConfig{gitCfg}) {
		t.Fatalf("Template eager list-error receivers were %#v, want one complete receiver", receivers)
	}
}

func TestTemplateLookupReturnsFirstExactCaseSensitiveMatchWithCompleteIdentities(t *testing.T) {
	firstCloud := &templateLookupCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/first exact cloud identity"}}}
	firstType := &MavenProject{PomFile: "/first exact pom identity"}
	first := CloudTemplate{
		Name: "Exact-Template",
		Project: Project{
			Path:        "/first exact project",
			GitInfo:     GitInfo{IsRepo: true, IsDirty: true, EnableCommit: true},
			ConfigFile:  "/first exact project/ply.json",
			Config:      ProjectConfiguration{Name: "first exact project", Profile: "first exact profile"},
			Type:        firstType,
			CloudConfig: firstCloud,
		},
	}
	lowerCase := CloudTemplate{Name: "exact-template", Project: Project{Path: "/case-sensitive lower project"}}
	laterExact := CloudTemplate{Name: "Exact-Template", Project: Project{Path: "/later exact project"}}
	recording := &recordingTemplateListLoader{results: []templateListLoadResult{{templates: []CloudTemplate{
		first, lowerCase, laterExact,
	}}}}

	actual, err := (GitCloudConfig{Impl: DirConfig{Path: "/complete exact-match receiver"}}).
		template(recording.dependencies(), "Exact-Template")

	if err != nil || !reflect.DeepEqual(actual, first) {
		t.Fatalf("Template exact match returned (%#v, %v), want first complete match %#v", actual, err, first)
	}
	if actual.Project.CloudConfig != firstCloud || actual.Project.Type != firstType {
		t.Fatalf("Template exact match lost embedded identities: CloudConfig=%T %#v Type=%T %#v",
			actual.Project.CloudConfig, actual.Project.CloudConfig, actual.Project.Type, actual.Project.Type)
	}
}

func TestTemplateLookupReturnsExactNotFoundForNilEmptyCaseMismatchAndEmptyName(t *testing.T) {
	tests := []struct {
		name      string
		templates []CloudTemplate
		lookup    string
	}{
		{name: "nil list", templates: nil, lookup: "missing-template"},
		{name: "empty list", templates: []CloudTemplate{}, lookup: "missing-template"},
		{name: "case mismatch", templates: []CloudTemplate{{Name: "Exact-Template"}}, lookup: "exact-template"},
		{name: "empty lookup name", templates: []CloudTemplate{{Name: "non-empty-template"}}, lookup: ""},
	}
	if len(tests) == 0 {
		t.Fatal("Template not-found test-case population is empty")
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingTemplateListLoader{results: []templateListLoadResult{{templates: test.templates}}}

			actual, err := (GitCloudConfig{Impl: DirConfig{Path: "/complete not-found receiver"}}).
				template(recording.dependencies(), test.lookup)

			wantError := fmt.Sprintf("could not find any valid templates with name: %s", test.lookup)
			if !reflect.DeepEqual(actual, CloudTemplate{}) || err == nil || err.Error() != wantError {
				t.Fatalf("Template not-found result was (%#v, %v), want (zero, exact %q)", actual, err, wantError)
			}
			if receivers, populationErr := recording.assertedReceivers(); populationErr != nil || len(receivers) != 1 {
				t.Fatalf("Template not-found receiver recording was (%#v, %v), want one non-empty receiver", receivers, populationErr)
			}
		})
	}
}

func TestTemplateLookupLoadsIndependentlyAcrossRepeatedInvocations(t *testing.T) {
	firstCloud := &templateLookupCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/first repeated cloud identity"}}}
	secondCloud := &templateLookupCloudConfigIdentity{GitCloudConfig: GitCloudConfig{Impl: DirConfig{Path: "/second repeated cloud identity"}}}
	first := CloudTemplate{Name: "repeated-template", Project: Project{Path: "/first repeated project", CloudConfig: firstCloud}}
	second := CloudTemplate{Name: "repeated-template", Project: Project{Path: "/second repeated project", CloudConfig: secondCloud}}
	recording := &recordingTemplateListLoader{results: []templateListLoadResult{
		{templates: []CloudTemplate{first}},
		{templates: []CloudTemplate{second}},
	}}
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete repeated lookup receiver"}}
	dependencies := recording.dependencies()

	firstActual, err := gitCfg.template(dependencies, "repeated-template")
	if err != nil || !reflect.DeepEqual(firstActual, first) || firstActual.Project.CloudConfig != firstCloud {
		t.Fatalf("first repeated Template lookup returned (%#v, %v), want exact first result", firstActual, err)
	}
	secondActual, err := gitCfg.template(dependencies, "repeated-template")
	if err != nil || !reflect.DeepEqual(secondActual, second) || secondActual.Project.CloudConfig != secondCloud {
		t.Fatalf("second repeated Template lookup returned (%#v, %v), want fresh exact second result", secondActual, err)
	}
	receivers, populationErr := recording.assertedReceivers()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(receivers, []GitCloudConfig{gitCfg, gitCfg}) {
		t.Fatalf("repeated Template list-loader receivers were %#v, want one complete receiver per invocation", receivers)
	}
}

func TestTemplateLookupDependenciesDefaultToSafeExactErrorWithoutPathAccess(t *testing.T) {
	path := filepath.Join(t.TempDir(), "developer-home-cloud-config-must-not-be-accessed")
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: path}}

	actual, err := gitCfg.template(templateLookupDependencies{}, "must-not-load")

	if !reflect.DeepEqual(actual, CloudTemplate{}) || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe Template lookup dependency default returned (%#v, %v), want (zero, exact %v)", actual, err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Lstat(path); !os.IsNotExist(statErr) {
		t.Fatalf("safe Template lookup dependency touched receiver path %q: %v", path, statErr)
	}
}

func TestRecordedTemplateListLoaderRejectsEmptyPopulations(t *testing.T) {
	if _, err := (&recordingTemplateListLoader{}).assertedReceivers(); err == nil {
		t.Fatal("empty recorded Template list-loader receiver population passed")
	}
}
