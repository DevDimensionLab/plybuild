package clean

import (
	"co-pilot/pkg/maven"
	"co-pilot/pkg/springio"
	"errors"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

func SpringManualVersion(model *pom.Model) error {
	springBootDependencies, err := springio.GetDependencies()
	if err != nil {
		return err
	}

	if model.Dependencies != nil {
		err = springManualVersion(model.Dependencies.Dependency, springBootDependencies, model)
		if err != nil {
			return err
		}
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		err = springManualVersion(model.DependencyManagement.Dependencies.Dependency, springBootDependencies, model)
		if err != nil {
			return err
		}
	}

	return nil
}

func springManualVersion(dependencies []pom.Dependency, springBootDependencies springio.IoDependenciesResponse, model *pom.Model) error {
	for _, dep := range dependencies {
		if dep.Version != "" && inMap(dep, springBootDependencies.Dependencies) {
			log.Warnf("found hardcoded version on spring-boot dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
			err := model.SetDependencyVersion(dep, "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func VersionToPropertyTags(model *pom.Model) error {
	if model.Dependencies != nil {
		err := versionToPropertyTags(model.Dependencies.Dependency, model)
		if err != nil {
			return err
		}
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		err := versionToPropertyTags(model.DependencyManagement.Dependencies.Dependency, model)
		if err != nil {
			return err
		}
	}

	return nil
}

func versionToPropertyTags(dependencies []pom.Dependency, model *pom.Model) error {
	for _, dep := range dependencies {
		if dep.Version != "" && !strings.HasPrefix(dep.Version, "${") {
			log.Warnf("found hardcoded version on dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
			err := model.ReplaceVersionTagWithProperty(dep)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func BlacklistedDependencies(model *pom.Model) error {
	if model.Dependencies == nil {
		return errors.New("could not find any dependencies")
	}

	bannedPomUrl := viper.GetString("banned_pom_url")
	if bannedPomUrl == "" {
		return errors.New("could not extract key `banned_pom_url` from config file ~/.co-pilot.yaml")
	}
	bannedModel, err := maven.GetBannedModel(bannedPomUrl)
	if err != nil {
		return err
	}

	for _, dep := range model.Dependencies.Dependency {
		for _, bannedDep := range bannedModel.Dependencies.Dependency {
			if bannedDep.GroupId == dep.GroupId && bannedDep.ArtifactId == dep.ArtifactId {
				log.Warnf("found blacklisted dependency %s:%s", dep.GroupId, dep.ArtifactId)
				//err := model.RemoveDependency(dep)
				//if err != nil {
				//	return err
				//}
			}
		}
	}

	return nil
}

func inMap(dep pom.Dependency, springBootDeps map[string]springio.Dependency) bool {
	for _, v := range springBootDeps {
		if v.GroupId == dep.GroupId && v.ArtifactId == dep.ArtifactId {
			return true
		}
	}

	return false
}
