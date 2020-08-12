package upgrade

import (
	"co-pilot/pkg/analyze"
	"co-pilot/pkg/maven"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"sort"
)

func Dependencies(directory string, local bool) error {
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
				err = Upgrade(model, dep)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			}
		}
	}

	sort.Sort(DependencySort(model.Dependencies.Dependency))
	return model.WriteToFile(pomFile)
}

func Upgrade(model *pom.Model, dep pom.Dependency) error {
	currentVersion, err := model.GetVersion(dep)
	metaData, err := maven.GetMetaData(dep.GroupId, dep.ArtifactId)
	if err == nil {
		if currentVersion != metaData.Versioning.Latest {
			fmt.Printf("[OUTDATED] %s:%s [%s => %s] \n", dep.GroupId, dep.ArtifactId, currentVersion, metaData.Versioning.Latest)
			_ = model.SetVersion(dep, metaData.Versioning.Latest)
		}
		return nil
	} else {
		return err
	}
}
