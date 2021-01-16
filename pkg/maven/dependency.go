package maven

import (
	"errors"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

func UpgradeDependency(groupId string, artifactId string) func(project config.Project, args ...interface{}) error {
	return func(project config.Project, args ...interface{}) error {
		return UpgradeDependencyOnModel(project.Type.Model(), groupId, artifactId)
	}
}

func UpgradeDependencyOnModel(model *pom.Model, groupId string, artifactId string) (err error) {
	if model.Dependencies != nil {
		err = specificDependencyUpgrade(model, model.Dependencies.Dependency, groupId, artifactId)
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		err = specificDependencyUpgrade(model, model.DependencyManagement.Dependencies.Dependency, groupId, artifactId)
	}

	return err
}

func Upgrade3PartyDependencies() func(project config.Project, args ...interface{}) error {
	return func(project config.Project, args ...interface{}) error {
		return upgrade3PartyDependenciesOnModel(&project)
	}
}

func upgrade3PartyDependenciesOnModel(project *config.Project) error {
	model := project.Type.Model()
	if model.Dependencies != nil {
		deps := model.Dependencies.Dependency
		upgradeDependencies(model, deps, project.Config.Settings, isSecondParty(model, false))
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		deps := model.DependencyManagement.Dependencies.Dependency
		upgradeDependencies(model, deps, project.Config.Settings, isSecondParty(model, false))
	}

	return nil
}

func Upgrade2PartyDependencies() func(project config.Project, args ...interface{}) error {
	return func(project config.Project, args ...interface{}) error {
		return upgrade2PartyDependenciesOnModel(&project)
	}
}

func upgrade2PartyDependenciesOnModel(project *config.Project) error {
	model := project.Type.Model()
	if model.Dependencies != nil {
		deps := model.Dependencies.Dependency
		upgradeDependencies(model, deps, project.Config.Settings, isSecondParty(model, true))
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		deps := model.DependencyManagement.Dependencies.Dependency
		upgradeDependencies(model, deps, project.Config.Settings, isSecondParty(model, true))
	}

	return nil
}

func specificDependencyUpgrade(model *pom.Model, availableDependencies []pom.Dependency, groupId string, artifactId string) error {
	for _, dep := range availableDependencies {
		if dep.Version != "" && dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return upgradeDependency(model, dep)
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

func upgradeDependencies(
	model *pom.Model,
	dependencies []pom.Dependency,
	settings config.ProjectSettings,
	condition func(groupId string) bool) {

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
		if depVersion != "" && condition(dep.GroupId) {
			err := upgradeDependency(model, dep)
			if err != nil {
				log.Warnf("%v", err)
			}
		}
	}
}

func upgradeDependency(model *pom.Model, dep pom.Dependency) error {
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

	metaData, err := GetMetaData(dep.GroupId, dep.ArtifactId)
	if err != nil {
		return err
	}

	latestVersion, err := metaData.LatestRelease()
	if err != nil {
		return nil
	}

	if currentVersion.IsDifferentFrom(latestVersion) {
		msg := fmt.Sprintf("outdated dependency %s:%s [%s => %s]", dep.GroupId, dep.ArtifactId, currentVersion.ToString(), latestVersion.ToString())
		if IsMajorUpgrade(currentVersion, latestVersion) {
			log.Warnf("major %s", msg)
		}
		if !latestVersion.IsReleaseVersion() {
			log.Warnf("%s | not release", msg)
		} else {
			log.Info(msg)
		}

		err = model.SetDependencyVersion(dep, latestVersion.ToString())
	}

	return nil
}
