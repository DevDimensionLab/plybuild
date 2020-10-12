package config

import (
	"co-pilot/pkg/file"
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

	project, _ := InitProjectFromPomFile(pomFile)
	err := project.SortAndWritePom()

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
