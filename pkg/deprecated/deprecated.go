package deprecated

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/logger"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

var log = logger.Context()

func FindDeprecated(model *pom.Model, deprecated config.CloudDeprecated) error {

	for _, dep := range model.Dependencies.Dependency {
		for _, depr := range deprecated.Data.Dependencies {
			if dep.GroupId == depr.GroupId && dep.ArtifactId == depr.ArtifactId {
				log.Infof("found deprecated dependency %s:%s", dep.GroupId, dep.ArtifactId)
			}
		}
	}

	return nil
}
