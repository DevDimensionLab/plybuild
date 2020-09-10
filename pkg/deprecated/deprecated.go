package deprecated

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/merge"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

var log = logger.Context()

func UpgradeDeprecated(model *pom.Model, deprecated config.CloudDeprecated, targetDirectory string, commitTemplates bool) error {
	var templates = make(map[string]bool)

	if model.Dependencies == nil {
		return nil
	}

	for _, modDep := range model.Dependencies.Dependency {
		for _, depRep := range deprecated.Data.Dependencies {
			if modDep.GroupId == depRep.GroupId && modDep.ArtifactId == depRep.ArtifactId {
				log.Infof("found deprecated dependency %s:%s", modDep.GroupId, modDep.ArtifactId)
				if err := model.RemoveDependency(modDep); err != nil {
					return err
				}
				if depRep.ReplacementTemplates != nil {
					for _, t := range depRep.ReplacementTemplates {
						templates[t] = true
					}
				}
			}
		}
	}

	for k, _ := range templates {
		if commitTemplates {
			log.Infof("applying template %s", k)
			if err := merge.TemplateName(k, targetDirectory); err != nil {
				log.Warnf("%v", err)
			}
		} else {
			log.Infof("template %s is ready for applying", k)
		}
	}

	return nil
}
