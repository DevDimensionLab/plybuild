package maven

import (
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/file"
	"testing"
)

func TestUpgradeDependency(t *testing.T) {
	project, err := config.InitProjectFromDirectory("test/dependency")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	err = project.InitProjectConfiguration()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	model := project.Type.Model()
	deps := model.Dependencies.Dependency
	upgradeDependencies(model, deps, project.Config.Settings, func(groupId string) bool {
		return true
	})

	if err := project.SortAndWritePom(); err != nil {
		t.Errorf("%v\n", err)
	}

	originPomFile := "test/dependency/origin.pom.xml"
	pomFile := "test/dependency/pom.xml"
	equal, err := file.Equal(originPomFile, pomFile)
	if !equal {
		t.Errorf("%s is not equal to %s", originPomFile, pomFile)
	}
}
