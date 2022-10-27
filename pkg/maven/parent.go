package maven

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/spring"
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
)

func UpgradeParent() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		if project.Config.Settings.DisableSpringBootUpgrade {
			return nil
		}
		return repository.upgradeParent(project.Type.Model())
	}
}

func (repository Repository) upgradeParent(model *pom.Model) error {
	parentGroupId := model.Parent.GroupId
	parentArtifactId := model.Parent.ArtifactId
	latestVersionMeta, err := repository.GetMetaData(parentGroupId, parentArtifactId)
	if err != nil {
		return err
	}

	latestVersion, err := latestVersionMeta.LatestRelease()
	if err != nil {
		return err
	}

	currentVersion, err := ParseVersion(model.Parent.Version)
	if err != nil {
		return err
	}

	if currentVersion.IsDifferentFrom(latestVersion) {
		msg := fmt.Sprintf("outdated parent version [%s => %s]", currentVersion.ToString(), latestVersion.ToString())
		if IsMajorUpgrade(currentVersion, latestVersion) {
			log.Warnf("major %s", msg)
		} else if !latestVersion.IsReleaseVersion() {
			log.Warnf("%s | not release", msg)
		} else {
			log.WithFields(logrus.Fields{
				"artifactId": parentArtifactId,
				"groupId":    parentGroupId,
				"oldVersion": currentVersion.ToString(),
				"newVersion": latestVersion.ToString(),
				"type":       "outdated parent",
			}).Info(msg)
		}
		model.Parent.Version = latestVersion.ToString()
	} else {
		log.Debugf("parent is the latest version [%s]", latestVersion.ToString())
	}

	return nil
}

func removeVersion(dependencies []pom.Dependency, springBootDependencies spring.IoDependenciesResponse, model *pom.Model) error {
	for _, dep := range dependencies {
		if dep.Version != "" && inMap(dep, springBootDependencies.Dependencies) {
			log.Warnf("found hardcoded version on spring-boot dependency %s:%s [%s]", dep.GroupId, dep.ArtifactId, dep.Version)
			err := model.SetDependencyVersion(dep, "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func inMap(dep pom.Dependency, springBootDeps map[string]spring.Dependency) bool {
	for _, springBootDep := range springBootDeps {
		if springBootDep.GroupId == dep.GroupId && springBootDep.ArtifactId == dep.ArtifactId {
			return true
		}
	}

	return false
}
