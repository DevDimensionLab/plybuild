package deprecated

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/logger"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

var log = logger.Context()

func RemoveDeprecated(model *pom.Model, deprecated config.CloudDeprecated) (templates map[string]bool, err error) {
	templates = make(map[string]bool)

	if model.Dependencies == nil {
		return templates, err
	}

	for _, modDep := range model.Dependencies.Dependency {
		for _, depRep := range deprecated.Data.Dependencies {
			if modDep.GroupId == depRep.GroupId && modDep.ArtifactId == depRep.ArtifactId {
				log.Infof("found deprecated dependency %s:%s", modDep.GroupId, modDep.ArtifactId)
				if err := model.RemoveDependency(modDep); err != nil {
					return templates, err
				}
				if depRep.ReplacementTemplates != nil {
					for _, t := range depRep.ReplacementTemplates {
						templates[t] = true
					}
				}
			}
		}
	}

	return templates, err
}
