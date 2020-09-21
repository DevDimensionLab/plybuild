package maven

import (
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

func UpgradeKotlin() func(pair PomWrapper, args ...interface{}) error {
	return func(pair PomWrapper, args ...interface{}) error {
		return UpgradeKotlinOnModel(pair.Model)
	}
}

func UpgradeKotlinOnModel(model *pom.Model) error {
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

	latestKotlinJdk8, err := GetMetaData("org.jetbrains.kotlin", "kotlin-maven-plugin")
	if err != nil {
		return err
	}

	latestVersion, err := latestKotlinJdk8.LatestRelease()
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
