package upgrade

import (
	"co-pilot/pkg/maven"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"sort"
)

func Dependencies(directory string) error {
	pomFile := directory + "/pom.xml"
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	for _, dep := range model.Dependencies.Dependency {
		if dep.Version != "" {
			currentVersion, err := model.GetVersion(dep)
			metaData, err := maven.GetMetaData(dep.GroupId, dep.ArtifactId)
			if err == nil {
				if currentVersion != metaData.Versioning.Latest {
					fmt.Printf("[OUTDATED] %s:%s [%s => %s] \n", dep.GroupId, dep.ArtifactId, currentVersion, metaData.Versioning.Latest)
					_ = model.SetVersion(dep, metaData.Versioning.Latest)
				}
			} else {
				fmt.Printf("[ERROR] %v\n", err)
			}
		}
	}

	sort.Sort(DependencySort(model.Dependencies.Dependency))
	return model.WriteToFile(pomFile + ".new")
}
