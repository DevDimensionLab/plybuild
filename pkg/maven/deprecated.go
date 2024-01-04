package maven

import (
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/pkg/config"
)

func RemoveDeprecated(cloudConfig config.CloudConfig, model *pom.Model) (templates []config.CloudTemplate, err error) {
	if model.Dependencies == nil {
		return
	}

	deprecated, err := cloudConfig.Deprecated()
	if err != nil {
		return
	}

	for _, modDep := range model.Dependencies.Dependency {
		for _, depRep := range deprecated.Data.Dependencies {
			if modDep.GroupId == depRep.GroupId && modDep.ArtifactId == depRep.ArtifactId {
				log.Infof("found deprecated dependency %s:%s", modDep.GroupId, modDep.ArtifactId)
				if err := model.RemoveDependency(modDep); err != nil {
					return templates, err
				}
				if depRep.ReplacementTemplates != nil {
					for _, replacementTemplate := range depRep.ReplacementTemplates {
						template, err := cloudConfig.Template(replacementTemplate)
						if err != nil {
							log.Warnln(err)
							continue
						}
						// TODO fix that it might add duplicates
						templates = append(templates, template)
					}
				}
			}
		}
	}

	return
}

func StatusDeprecated() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		cloudConfig := project.CloudConfig
		model := project.Type.Model()
		if model == nil || model.Dependencies == nil {
			return nil
		}

		deprecated, err := cloudConfig.Deprecated()
		if err != nil {
			return err
		}

		for _, modDep := range model.Dependencies.Dependency {
			for _, depRep := range deprecated.Data.Dependencies {
				if modDep.GroupId == depRep.GroupId && modDep.ArtifactId == depRep.ArtifactId {
					log.Infof("found deprecated dependency %s:%s", modDep.GroupId, modDep.ArtifactId)
				}
			}
		}
		return nil
	}
}
