package maven

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/sirupsen/logrus"
	"strings"
)

func Upgrade2PartyDependencies() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		return repository.upgradeDependenciesForProject(&project, true)
	}
}

func Upgrade3PartyDependencies() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		return repository.upgradeDependenciesForProject(&project, false)
	}
}

func UpgradeDependency(groupId string, artifactId string) func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		return repository.upgradeDependencyOnModel(project.Type.Model(), groupId, artifactId)
	}
}

func (repository Repository) upgradeDependencyOnModel(model *pom.Model, groupId string, artifactId string) (err error) {
	if model.Dependencies != nil {
		err = repository.specificDependencyUpgrade(model, model.Dependencies.Dependency, groupId, artifactId)
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		err = repository.specificDependencyUpgrade(model, model.DependencyManagement.Dependencies.Dependency, groupId, artifactId)
	}

	return err
}

func (repository Repository) upgradeDependenciesForProject(project *config.Project, enabledSecondParty bool) error {
	model := project.Type.Model()
	if model.Dependencies != nil {
		deps := model.Dependencies.Dependency
		repository.upgradeDependencies(model, deps, project.Config.Settings, isSecondParty(model, enabledSecondParty), model.SetDependencyVersion)
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		deps := model.DependencyManagement.Dependencies.Dependency
		repository.upgradeDependencies(model, deps, project.Config.Settings, isSecondParty(model, enabledSecondParty), model.SetDependencyVersion)
	}

	return nil
}

func UpgradeDependenciesWithVersions() func(repository Repository, project config.Project) error {
	return func(repository Repository, project config.Project) error {
		updateProp := func(dep pom.Dependency, version string) error {
			if strings.HasPrefix(dep.Version, "${") {
				versionKey := strings.Trim(dep.Version, "${}")
				args := UpdateProperty(versionKey, version)
				return RunOn("mvn", args...)(repository, project)
			} else {
				log.Debugf("Upgrading dependencies with `use-latest-version`")
				args := UseLatestVersion(dep.GroupId, dep.ArtifactId)
				return RunOn("mvn", args...)(repository, project)
			}
		}
		allDeps := func(groupId string) bool { return true }
		model := project.Type.Model()
		if model.Dependencies != nil {
			deps := model.Dependencies.Dependency
			repository.upgradeDependencies(model, deps, project.Config.Settings, allDeps, updateProp)
		}

		if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
			deps := model.DependencyManagement.Dependencies.Dependency
			repository.upgradeDependencies(model, deps, project.Config.Settings, allDeps, updateProp)
		}

		return nil
	}
}

func (repository Repository) specificDependencyUpgrade(model *pom.Model, availableDependencies []pom.Dependency, groupId string, artifactId string) error {
	for _, dep := range availableDependencies {
		if dep.Version != "" && dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return repository.upgradeDependency(model, dep, nil, model.SetDependencyVersion)
		}
	}

	return errors.New(fmt.Sprintf("could not find %s:%s in project", groupId, artifactId))
}

func isSecondParty(model *pom.Model, enabled bool) func(groupId string) bool {
	secondPartyGroupId, err := model.GetSecondPartyGroupId()
	if err != nil {
		log.Warnln(err)
	}
	log.Debugf("secondParty groupId is: %s", secondPartyGroupId)

	return func(groupId string) bool {
		isSecondParty, err := isSecondPartyGroupId(groupId, secondPartyGroupId)
		log.Debugf("upgrade secondParty %t, groupId:%s isSecondParty %t", enabled, groupId, isSecondParty)
		if err != nil {
			log.Debug(err)
			return false
		}
		return isSecondParty == enabled
	}
}

func (repository Repository) upgradeDependencies(
	model *pom.Model,
	dependencies []pom.Dependency,
	settings config.ProjectSettings,
	condition func(groupId string) bool,
	action func(dep pom.Dependency, version string) error) {

	for _, dep := range dependencies {
		if settings.DependencyIsIgnored(dep) {
			log.Debugf("ignoring dependency %s:%s", dep.GroupId, dep.ArtifactId)
			continue
		}
		depVersion, err := model.GetDependencyVersion(dep)
		if err != nil {
			log.Infof("failed to get version for %s:%s", dep.GroupId, dep.ArtifactId)
			continue
		}

		maxVersion := func() *JavaVersion {
			if maxVersionForDependency := settings.MaxVersionFor(dep); maxVersionForDependency != "" {
				maxVersion, err := ParseVersion(maxVersionForDependency)
				if err != nil {
					log.Debugf("failed to get max version for %s:%s", dep.GroupId, dep.ArtifactId)
					return nil
				}
				return &maxVersion
			}
			return nil
		}

		if depVersion != "" && condition(dep.GroupId) {
			err := repository.upgradeDependency(model, dep, maxVersion(), action)
			if err != nil {
				log.Warnf("%v", err)
			}
		}
	}
}

func (repository Repository) upgradeDependency(model *pom.Model, dep pom.Dependency, maxVersion *JavaVersion, action func(dep pom.Dependency, version string) error) error {
	if dep.Version == "${project.version}" || dep.Version == "${revision}" {
		return nil
	}

	depVersion, err := model.GetDependencyVersion(dep)
	if err != nil {
		return err
	}

	currentVersion, err := ParseVersion(depVersion)
	if err != nil {
		return err
	}

	//serviceUrl, err := config.LinkFromService(config.Services(), dep.GroupId, dep.ArtifactId, "info")
	//if err != nil {
	//	log.Infoln(err)
	//}
	//log.Debugf("%v", serviceUrl)

	repo := repository

	metaData, err := repo.GetMetaData(dep.GroupId, dep.ArtifactId)
	if err != nil {
		return err
	}

	latestVersion, err := metaData.LatestRelease()
	if err != nil {
		return nil
	}
	if maxVersion != nil {
		log.Warnf("dependency %s:%s is held back at version [%s]", dep.GroupId, dep.ArtifactId, maxVersion.ToString())
		latestVersion = *maxVersion
	}

	log.Debugf("comparing current version %s with latest version %s", currentVersion.ToString(), latestVersion.ToString())

	if currentVersion.IsLessThan(latestVersion) {
		msg := fmt.Sprintf("outdated dependency %s:%s [%s vs %s]", dep.GroupId, dep.ArtifactId, currentVersion.ToString(), latestVersion.ToString())
		var metaDataLogger = log
		if logger.IsFieldLogger() {
			metaDataLogger = log.WithFields(logrus.Fields{
				"artifactId":        dep.ArtifactId,
				"groupId":           dep.GroupId,
				"oldVersion":        currentVersion.ToString(),
				"newVersion":        latestVersion.ToString(),
				"versionIsProperty": strings.HasPrefix(dep.Version, "${") && strings.HasSuffix(dep.Version, "}"),
				"versionValue":      dep.Version,
				"type":              "outdated dependency",
			})
		}
		if IsMajorUpgrade(currentVersion, latestVersion) {
			metaDataLogger.Warnf("major %s", msg)
		}
		if !latestVersion.IsReleaseVersion() {
			log.Warnf("%s | not release", msg)
		} else {
			metaDataLogger.Info(msg)
		}

		//err = model.SetDependencyVersion(dep, latestVersion.ToString())
		err = action(dep, latestVersion.ToString())
	}

	return nil
}
