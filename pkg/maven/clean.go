package maven

import (
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"strings"
)

func ChangeVersionToPropertyTags() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		return ChangeVersionToPropertyTagsOnModel(project.Type.Model())
	}
}

func ChangeVersionToPropertyTagsOnModel(model *pom.Model) error {
	if model.Dependencies != nil {
		for _, dep := range model.Dependencies.Dependency {
			if dep.Version != "" && !strings.HasPrefix(dep.Version, "${") {
				log.Warnf("found hardcoded version on dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
				err := model.ReplaceVersionTagForDependency(dep)
				if err != nil {
					return err
				}
			}
		}
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		for _, managementDep := range model.DependencyManagement.Dependencies.Dependency {
			if managementDep.Version != "" && !strings.HasPrefix(managementDep.Version, "${") {
				log.Warnf("found hardcoded version on dependencyManagement dependency %s:%s [%s]",
					managementDep.GroupId, managementDep.ArtifactId, managementDep.Version)
				err := model.ReplaceVersionTagForDependencyManagement(managementDep)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
