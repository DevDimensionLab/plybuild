package upgrade

import (
	"co-pilot/pkg/maven"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"sort"
)

func Kotlin(directory string) error {
	pomFile := directory + "/pom.xml"
	model, err := pom.GetModelFrom(pomFile)

	if err != nil {
		return err
	}

	currentVersion, err := model.Properties.FindKey("kotlin.version")
	latestKotlinJdk8, err := maven.GetMetaData("org.jetbrains.kotlin", "kotlin-maven-plugin")
	if err != nil {
		return err
	}

	if currentVersion != latestKotlinJdk8.Versioning.Latest {
		err = model.Properties.SetKey("kotlin.version", latestKotlinJdk8.Versioning.Latest)
		if err != nil {
			return err
		}
		fmt.Printf("[OUTDATED] kotlin version [%s => %s] \n", currentVersion, latestKotlinJdk8.Versioning.Latest)

		sort.Sort(DependencySort(model.Dependencies.Dependency))
		return model.WriteToFile(pomFile)
	} else {
		fmt.Printf("[INFO] No update needed, kotlin is the the latest version [%s]\n", currentVersion)
		return nil
	}
}
