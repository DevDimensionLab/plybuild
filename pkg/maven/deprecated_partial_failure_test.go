package maven

import (
	"errors"
	"reflect"
	"testing"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/sirupsen/logrus"
)

type partialFailureCloudConfig struct {
	config.GitCloudConfig
	deprecated     config.CloudDeprecated
	deprecatedErr  error
	templates      map[string]config.CloudTemplate
	templateErrors map[string]error
	templateCalls  []string
}

func (cloud *partialFailureCloudConfig) Deprecated() (config.CloudDeprecated, error) {
	return cloud.deprecated, cloud.deprecatedErr
}

func (cloud *partialFailureCloudConfig) Template(name string) (config.CloudTemplate, error) {
	cloud.templateCalls = append(cloud.templateCalls, name)
	if err := cloud.templateErrors[name]; err != nil {
		return config.CloudTemplate{}, err
	}
	return cloud.templates[name], nil
}

func TestRemoveDeprecatedReturnsExactOrderedPartialTemplatesWithRemovalError(t *testing.T) {
	dependencyA := pom.Dependency{GroupId: "a-group", ArtifactId: "a-artifact", Version: "1.0.0"}
	dependencyB := pom.Dependency{GroupId: "b-group", ArtifactId: "b-artifact", Version: "1.0.0"}
	dependencyC := pom.Dependency{GroupId: "c-group", ArtifactId: "c-artifact", Version: "1.0.0"}
	model := &pom.Model{Dependencies: &pom.Dependencies{Dependency: []pom.Dependency{dependencyA, dependencyB, dependencyC}}}
	deprecated := config.CloudDeprecated{}
	deprecated.Data.Dependencies = []config.CloudDeprecatedDependency{
		deprecatedDependency(dependencyA, "a-template"),
		deprecatedDependency(dependencyC, "c-template"),
	}
	wantTemplates := []config.CloudTemplate{
		{Name: "a-template", Project: config.Project{Path: "/templates/a"}},
		{Name: "c-template", Project: config.Project{Path: "/templates/c"}},
	}
	if len(wantTemplates) == 0 {
		t.Fatal("RemoveDeprecated partial-template expectation population is empty")
	}
	cloud := &partialFailureCloudConfig{
		deprecated: deprecated,
		templates: map[string]config.CloudTemplate{
			"a-template": wantTemplates[0],
			"c-template": wantTemplates[1],
		},
	}

	actual, err := RemoveDeprecated(cloud, model)

	wantError := errors.New("could not find dependency: c-group:c-artifact in model")
	if err == nil || err.Error() != wantError.Error() {
		t.Fatalf("RemoveDeprecated error was %v, want exact removal error %v", err, wantError)
	}
	if !reflect.DeepEqual(actual, wantTemplates) {
		t.Fatalf("RemoveDeprecated partial templates were %#v, want exact ordered partial result %#v", actual, wantTemplates)
	}
	wantCalls := []string{"a-template", "c-template"}
	if !reflect.DeepEqual(cloud.templateCalls, wantCalls) {
		t.Fatalf("RemoveDeprecated template calls were %#v, want exact successful prefix %#v", cloud.templateCalls, wantCalls)
	}
	if !reflect.DeepEqual(model.Dependencies.Dependency, []pom.Dependency{dependencyB}) {
		t.Fatalf("RemoveDeprecated model at failure was %#v, want exact populated mutation with only %#v", model.Dependencies.Dependency, dependencyB)
	}
}

func TestRemoveDeprecatedReportsReplacementFailureExactlyAndContinuesToLaterReplacement(t *testing.T) {
	dependency := pom.Dependency{GroupId: "deprecated-group", ArtifactId: "deprecated-artifact", Version: "1.0.0"}
	model := &pom.Model{Dependencies: &pom.Dependencies{Dependency: []pom.Dependency{dependency}}}
	deprecated := config.CloudDeprecated{}
	deprecated.Data.Dependencies = []config.CloudDeprecatedDependency{
		deprecatedDependency(dependency, "missing-template", "later-template"),
	}
	laterTemplate := config.CloudTemplate{Name: "later-template", Project: config.Project{Path: "/templates/later"}}
	sentinel := errors.New("complete replacement template failure")
	cloud := &partialFailureCloudConfig{
		deprecated:     deprecated,
		templates:      map[string]config.CloudTemplate{"later-template": laterTemplate},
		templateErrors: map[string]error{"missing-template": sentinel},
	}
	hook := captureMavenPartialFailureLogs(t)

	actual, err := RemoveDeprecated(cloud, model)

	if err != nil {
		t.Fatalf("RemoveDeprecated replacement continuation returned an error: %v", err)
	}
	if !reflect.DeepEqual(actual, []config.CloudTemplate{laterTemplate}) {
		t.Fatalf("RemoveDeprecated replacement result was %#v, want later template %#v", actual, laterTemplate)
	}
	if !reflect.DeepEqual(cloud.templateCalls, []string{"missing-template", "later-template"}) {
		t.Fatalf("RemoveDeprecated replacement calls were %#v, want exact failure-then-success order", cloud.templateCalls)
	}
	assertMavenPartialFailureLogs(t, hook, []mavenPartialFailureLog{{
		level: logrus.WarnLevel, message: sentinel.Error(),
	}})
}

func deprecatedDependency(dependency pom.Dependency, replacements ...string) config.CloudDeprecatedDependency {
	deprecated := config.CloudDeprecatedDependency{
		GroupId:              dependency.GroupId,
		ArtifactId:           dependency.ArtifactId,
		ReplacementTemplates: append([]string(nil), replacements...),
	}
	return deprecated
}
