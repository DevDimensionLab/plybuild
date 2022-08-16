package maven

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
)

func UpgradeKotlin() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		if project.Config.Settings.DisableKotlinUpgrade {
			return nil
		}
		model := project.Type.Model()
		return repository.upgradeKotlinOnModel(model, model.Properties.SetKey)
	}
}

func UpgradeKotlinWithVersions() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		if project.Config.Settings.DisableKotlinUpgrade {
			return nil
		}
		model := project.Type.Model()
		return repository.upgradeKotlinOnModel(model, func(propKey, version string) error {
			return RunOn("mvn", UpdateProperty(propKey, version)...)(repository, project)
		})
	}
}

func (repository Repository) upgradeKotlinOnModel(model *pom.Model, action func(propKey, version string) error) error {
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

	latestKotlinJdk8, err := repository.GetMetaData("org.jetbrains.kotlin", "kotlin-maven-plugin")
	if err != nil {
		return err
	}

	latestVersion, err := latestKotlinJdk8.LatestRelease()
	if err != nil {
		return err
	}

	if currentVersion.IsLessThan(latestVersion) {
		msg := fmt.Sprintf("outdated kotlin version [%s => %s]", currentVersion.ToString(), latestVersion.ToString())
		if IsMajorUpgrade(currentVersion, latestVersion) {
			log.Warnf("major %s", msg)
		} else if !latestVersion.IsReleaseVersion() {
			log.Warnf("%s | not release", msg)
		} else {
			log.Infof(msg)
		}

		//err = model.Properties.SetKey("kotlin.version", latestVersion.ToString())
		err = action("kotlin.version", latestVersion.ToString())
		if err != nil {
			return err
		}
	} else {
		log.Infof("No update needed, kotlin is the the latest version [%s]", currentVersion.ToString())
	}
	return nil
}
