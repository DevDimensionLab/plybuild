package maven

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"testing"
)

func TestSortAndWritePom_sort_enabled_by_default(t *testing.T) {
	original := file.Path("test/sorting/pom.xml.orig")
	pomFile := file.Path("test/sorting/pom.xml.sorted")
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	err = SortAndWritePom(PomWrapper{
		PomFile:       pomFile,
		Model:         model,
		ProjectConfig: config.ProjectConfiguration{},
	}, true)

	equal, err := file.Equal(original, pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if equal {
		t.Errorf("%s is equal to %s, therefore it does not appear to be sorted", pomFile, original)
	}
}

func TestSortAndWritePom_sort_enabled(t *testing.T) {
	original := file.Path("test/sorting/pom.xml.orig")
	pomFile := file.Path("test/sorting/pom.xml.sorted")
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	projectConfig := config.ProjectConfiguration{
		Settings: config.ProjectSettings{
			DisableDependencySort: false,
		},
	}

	err = SortAndWritePom(PomWrapper{
		PomFile:       pomFile,
		Model:         model,
		ProjectConfig: projectConfig,
	}, true)

	equal, err := file.Equal(original, pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if equal {
		t.Errorf("%s is equal to %s, therefore it does not appear to be sorted", pomFile, original)
	}
}

func TestSortAndWritePom_sort_disabled(t *testing.T) {
	original := file.Path("test/sorting/pom.xml.orig")
	pomFile := file.Path("test/sorting/pom.xml.unsorted")
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	projectConfig := config.ProjectConfiguration{
		Settings: config.ProjectSettings{
			DisableDependencySort: true,
		},
	}

	err = SortAndWritePom(PomWrapper{
		PomFile:       pomFile,
		Model:         model,
		ProjectConfig: projectConfig,
	}, true)

	equal, err := file.Equal(original, pomFile)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if !equal {
		t.Errorf("%s is not equal to %s, therefore it appears to be sorted", pomFile, original)
	}
}
