package upgrade

import (
	"co-pilot/pkg/analyze"
	"co-pilot/pkg/maven"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"sort"
)

func Dependencies(directory string, local bool, dryRun bool) error {
	pomFile := directory + "/pom.xml"
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	localGroupId, err := analyze.GetLocalGroupId(model)
	if err != nil {
		return err
	}

	for _, dep := range model.Dependencies.Dependency {
		if dep.Version != "" {
			isLocal, err := analyze.IsLocalGroupId(dep.GroupId, localGroupId)
			if err == nil && isLocal == local {
				err = DependencyUpgrade(model, dep)
				if err != nil {
					log.Errorf("%v", err)
				}
			}
		}
	}

	if !dryRun {
		sort.Sort(DependencySort(model.Dependencies.Dependency))
		return model.WriteToFile(pomFile)
	} else {
		return nil
	}
}

func DependencyUpgrade(model *pom.Model, dep pom.Dependency) error {
	currentVersion, err := model.GetDependencyVersion(dep)
	metaData, err := maven.GetMetaData(dep.GroupId, dep.ArtifactId)
	if err == nil {
		if currentVersion != metaData.Versioning.Release {

			log.Warnf("outdated dependency %s:%s [%s => %s] \n", dep.GroupId, dep.ArtifactId, currentVersion, metaData.Versioning.Release)
			_ = model.SetDependencyVersion(dep, metaData.Versioning.Release)
		}
		return nil
	} else {
		return err
	}
}
