package maven

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"testing"
)

func TestSortAndWritePom_sort_enabled_by_default(t *testing.T) {
	original := file.Path("test/sorting/pom.xml")
	pomFile := file.Path("test/sorting/sorted/pom.xml")

	project, _ := config.InitProjectFromPomFile(pomFile)
	err := SortAndWritePom(project, true)

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

	projectConfig := config.ProjectConfiguration{
		Settings: config.ProjectSettings{
			DisableDependencySort: false,
		},
	}

	project, _ := config.InitProjectFromPomFile(pomFile)
	project.Config = projectConfig
	err := SortAndWritePom(project, true)

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

	projectConfig := config.ProjectConfiguration{
		Settings: config.ProjectSettings{
			DisableDependencySort: true,
		},
	}

	project, _ := config.InitProjectFromPomFile(pomFile)
	project.Config = projectConfig
	err := SortAndWritePom(project, true)

	equal, err := file.Equal(original, pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if !equal {
		t.Errorf("%s is not equal to %s, therefore it appears to be sorted", pomFile, original)
	}
}
