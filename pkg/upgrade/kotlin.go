package upgrade

import (
	"co-pilot/pkg/maven"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"sort"
)

func Kotlin(directory string, dryRun bool) error {
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

	if currentVersion != latestKotlinJdk8.Versioning.Release {
		err = model.Properties.SetKey("kotlin.version", latestKotlinJdk8.Versioning.Release)
		if err != nil {
			return err
		}
		log.Warnf("outdated kotlin version [%s => %s]", currentVersion, latestKotlinJdk8.Versioning.Release)

		if !dryRun {
			sort.Sort(DependencySort(model.Dependencies.Dependency))
			return model.WriteToFile(pomFile)
		} else {
			return nil
		}
	} else {
		log.Infof("No update needed, kotlin is the the latest version [%s]", currentVersion)
		return nil
	}
}
