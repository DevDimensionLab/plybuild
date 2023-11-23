package maven

import (
	"github.com/devdimensionlab/ply/pkg/config"
	"github.com/devdimensionlab/ply/pkg/file"
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

	repo, err := DefaultRepository()
	if err != nil {
		t.Error(err)
	}

	repo.upgradeDependencies(model, deps, project.Config.Settings, func(groupId string) bool {
		return true
	}, model.SetDependencyVersion)

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
