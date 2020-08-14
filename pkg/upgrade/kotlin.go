package upgrade

import (
	"co-pilot/pkg/maven"
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func Kotlin(model *pom.Model) error {
	if model.Properties == nil {
		return errors.New("could not kotlin version because pom does not contain any properties")
	}

	version, err := model.Properties.FindKey("kotlin.version")
	if err != nil {
		return err
	}

	currentVersion, err := ParseVersion(version)
	if err != nil {
		return err
	}

	latestKotlinJdk8, err := maven.GetMetaData("org.jetbrains.kotlin", "kotlin-maven-plugin")
	if err != nil {
		return err
	}

	latestVersion, err := ParseVersion(latestKotlinJdk8.Versioning.Release)
	if err != nil {
		return err
	}

	if currentVersion.IsDifferentFrom(latestVersion) {
		msg := fmt.Sprintf("outdated kotlin version [%s => %s]", currentVersion.ToString(), latestVersion.ToString())
		if IsMajorUpgrade(currentVersion, latestVersion) {
			log.Warnf("major %s", msg)
		} else if !latestVersion.IsReleaseVersion() {
			log.Warnf("%s | not release", msg)
		} else {
			log.Infof(msg)
		}

		err = model.Properties.SetKey("kotlin.version", latestVersion.ToString())
		if err != nil {
			return err
		}
	} else {
		log.Infof("No update needed, kotlin is the the latest version [%s]", currentVersion.ToString())
	}
	return nil
}
