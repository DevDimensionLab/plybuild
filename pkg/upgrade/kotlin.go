package upgrade

import (
	"co-pilot/pkg/maven"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func Kotlin(model *pom.Model) error {
	currentVersion, err := model.Properties.FindKey("kotlin.version")
	latestKotlinJdk8, err := maven.GetMetaData("org.jetbrains.kotlin", "kotlin-maven-plugin")
	if err != nil {
		return err
	}

	latestVersion := latestKotlinJdk8.Versioning.Release
	if currentVersion != latestVersion {
		log.Warnf("outdated kotlin version [%s => %s]", currentVersion, latestVersion)
		err = model.Properties.SetKey("kotlin.version", latestVersion)
		if err != nil {
			return err
		}
	} else {
		log.Infof("No update needed, kotlin is the the latest version [%s]", currentVersion)
	}
	return nil
}
