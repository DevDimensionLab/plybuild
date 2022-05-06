package config

import (
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"testing"
)

func TestProjectConfiguration_SourceMainPath(t *testing.T) {
	sourceConfig, err := InitProjectConfigurationFromDir("test/cloud-config/templates/test-template")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	targetConfig, err := InitProjectConfigurationFromDir("test/target-app")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	expectedSourceRoot := "src/main/kotlin/no/copilot/template/test"
	if sourceConfig.SourceMainPath() != expectedSourceRoot {
		t.Errorf("expected %s, but got instead %s", expectedSourceRoot, sourceConfig.SourceMainPath())
	}

	expectedTargetRoot := "src/main/java/no/copilot/template/target"
	if targetConfig.SourceMainPath() != expectedTargetRoot {
		t.Errorf("expected %s, but got instead %s", expectedTargetRoot, targetConfig.SourceMainPath())
	}
}

func TestSortAndWritePom_sort_enabled_by_default(t *testing.T) {
	original := file.Path("test/sorting/pom.xml")
	pomFile := file.Path("test/sorting/sorted/pom.xml")

	project, err := InitProjectFromPomFile(pomFile)
	if err != nil {
		t.Fatal(err)
	}
	err = project.SortAndWritePom()
	if err != nil {
		t.Fatal(err)
	}

	equal, err := file.Equal(original, pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if equal {
		t.Errorf("%s is equal to %s, therefore it does not appear to be sorted", pomFile, original)
	}
}

func TestSortAndWritePom_sort_enabled(t *testing.T) {
	original := file.Path("test/sorting/pom.xml")
	pomFile := file.Path("test/sorting/sorted/pom.xml")

	projectConfig := ProjectConfiguration{
		Settings: ProjectSettings{
			DisableDependencySort: false,
		},
	}

	project, _ := InitProjectFromPomFile(pomFile)
	project.Config = projectConfig
	err := project.SortAndWritePom()

	equal, err := file.Equal(original, pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if equal {
		t.Errorf("%s is equal to %s, therefore it does not appear to be sorted", pomFile, original)
	}
}

func TestSortAndWritePom_sort_disabled(t *testing.T) {
	original := file.Path("test/sorting/pom.xml")
	pomFile := file.Path("test/sorting/unsorted/pom.xml")

	projectConfig := ProjectConfiguration{
		Settings: ProjectSettings{
			DisableDependencySort: true,
		},
	}

	project, _ := InitProjectFromPomFile(pomFile)
	project.Config = projectConfig
	err := project.SortAndWritePom()

	equal, err := file.Equal(original, pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if !equal {
		t.Errorf("%s is not equal to %s, therefore it appears to be sorted", pomFile, original)
	}
}

func TestCreateProjectConfig(t *testing.T) {
	project, err := InitProjectFromDirectory("test/project-config")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	err = project.InitProjectConfiguration()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if err := project.Config.WriteTo(project.ConfigFile); err != nil {
		t.Errorf("%v\n", err)
	}

	originConfig := "test/project-config/origin.co-pilot.json"
	newConfig := "test/project-config/co-pilot.json"
	equal, err := file.Equal(originConfig, newConfig)
	if !equal {
		t.Errorf("%s is not equal to %s", originConfig, newConfig)
	}
}
