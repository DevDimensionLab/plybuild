package upgrade

import (
	"co-pilot/pkg/springio"
	"errors"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"strings"
)

func Clean(model *pom.Model) error {
	springBootDependencies, err := springio.GetDependencies()
	if err != nil {
		return err
	}

	if model.Dependencies == nil {
		return errors.New("could not find any dependencies")
	}

	for _, dep := range model.Dependencies.Dependency {
		if dep.Version != "" {
			if inMap(dep, springBootDependencies.Dependencies) {
				log.Warnf("found hardcoded version on spring-boot dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
				err = model.SetDependencyVersion(dep, "")
				if err != nil {
					return err
				}
			} else if !strings.HasPrefix(dep.Version, "${") {
				log.Warnf("found hardcoded version on dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
				err = model.ReplaceVersionTagWithProperty(dep)
				if err != nil {
					return err
				}
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
