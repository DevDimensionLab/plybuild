package upgrade

import (
	"co-pilot/pkg/springio"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func Clean(model *pom.Model) error {

	springBootDependencies, err := springio.GetDependencies()
	if err != nil {
		return err
	}

	for _, dep := range model.Dependencies.Dependency {
		if dep.Version != "" && inMap(dep, springBootDependencies.Dependencies) {
			log.Warnf("found version on spring-boot dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
			err = model.SetDependencyVersion(dep, "")
			if err != nil {
				return err
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
