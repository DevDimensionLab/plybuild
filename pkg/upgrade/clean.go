package upgrade

import (
	"co-pilot/pkg/maven"
	"co-pilot/pkg/springio"
	"errors"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

func Clean(model *pom.Model) error {
	if model.Dependencies == nil {
		return errors.New("could not find any dependencies")
	}

	err := checkManualVersion(model)
	if err != nil {
		return err
	}

	err = checkBannedDependencies(model)
	if err != nil {
		return err
	}

	return nil
}

func checkManualVersion(model *pom.Model) error {
	springBootDependencies, err := springio.GetDependencies()
	if err != nil {
		return err
	}

	for _, dep := range model.Dependencies.Dependency {
		if dep.Version != "" {
			if inMap(dep, springBootDependencies.Dependencies) {
				log.Warnf("found hardcoded version on spring-boot dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
				err := model.SetDependencyVersion(dep, "")
				if err != nil {
					return err
				}
			} else if !strings.HasPrefix(dep.Version, "${") {
				log.Warnf("found hardcoded version on dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
				err := model.ReplaceVersionTagWithProperty(dep)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func checkBannedDependencies(model *pom.Model) error {
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
				log.Warnf("found banned dependency %s:%s", dep.GroupId, dep.ArtifactId)
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
