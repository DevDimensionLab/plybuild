package maven

import (
	"errors"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/mvn-pom-mutator/pkg/pom"
)

func UpgradeKotlin() func(project config.Project, args ...interface{}) error {
	return func(project config.Project, args ...interface{}) error {
		if project.Config.Settings.DisableKotlinUpgrade {
			return nil
		}
		return upgradeKotlinOnModel(project.Type.Model())
	}
}

func upgradeKotlinOnModel(model *pom.Model) error {
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
