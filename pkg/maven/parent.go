package maven

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/logger"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/spring"
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
)

func UpgradeParent() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		if project.Config.Settings.DisableSpringBootUpgrade {
			return nil
		}
		return repository.upgradeParent(project.Type.Model(), project)
	}
}

func (repository Repository) upgradeParent(model *pom.Model, project config.Project) error {
	if model.Parent == nil {
		return errors.New("no parent found")
	}

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
		msg := fmt.Sprintf("outdated parent version [%s vs %s]", currentVersion.ToString(), latestVersion.ToString())
		var metaDataLogger = log
		if logger.IsFieldLogger() {
			metaDataLogger = log.WithFields(logrus.Fields{
				"artifactId": parentArtifactId,
				"groupId":    parentGroupId,
				"oldVersion": currentVersion.ToString(),
				"newVersion": latestVersion.ToString(),
				"type":       "outdated parent",
			})
		}
		if strings.Contains(parentGroupId, "org.springframework") && project.Config.Settings.MaxSpringBootVersion != "" {
			maxBootVersion, err := ParseVersion(project.Config.Settings.MaxSpringBootVersion)
			if err != nil {
				log.Errorln("Found non-parsable version in settings.maxSpringBootVersion")
				return nil
			}
			if maxBootVersion.IsLessThan(latestVersion) {
				log.Infof("Keeping spring-boot version at %s - newest version is %s", maxBootVersion.ToString(), latestVersion.ToString())
				return nil
			}
		}
		if IsMajorUpgrade(currentVersion, latestVersion) {
			metaDataLogger.Warnf("major %s", msg)
		} else if !latestVersion.IsReleaseVersion() {
			log.Warnf("%s | not release", msg)
		} else {
			metaDataLogger.Info(msg)
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
