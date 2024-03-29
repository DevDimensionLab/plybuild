package maven

import (
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/spring"
)

func CleanManualVersions() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		return cleanManualVersions(project.Type.Model())
	}
}

func cleanManualVersions(model *pom.Model) error {
	springBootDependencies, err := spring.GetDependencies()
	if err != nil {
		return err
	}

	if model.Dependencies != nil {
		err = removeVersion(model.Dependencies.Dependency, springBootDependencies, model)
		if err != nil {
			return err
		}
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		err = removeVersion(model.DependencyManagement.Dependencies.Dependency, springBootDependencies, model)
		if err != nil {
			return err
		}
	}

	return nil
}
