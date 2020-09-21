package maven

import (
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

func UpgradeDependency(groupId string, artifactId string) func(pair PomWrapper, args ...interface{}) error {
	return func(pair PomWrapper, args ...interface{}) error {
		return UpgradeDependencyOnModel(pair.Model, groupId, artifactId)
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

func specificDependencyUpgrade(model *pom.Model, availableDependencies []pom.Dependency, groupId string, artifactId string) error {
	for _, dep := range availableDependencies {
		if dep.Version != "" && dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return upgradeDependency(model, dep)
		}
	}

	return errors.New(fmt.Sprintf("could not find %s:%s in project", groupId, artifactId))
}

func SecondParty(model *pom.Model, true bool) func(groupId string) bool {
	secondPartyGroupId, err := GetSecondPartyGroupId(model)
	if err != nil {
		log.Warnln(err)
	}
	log.Debugf("secondParty groupId is: %s", secondPartyGroupId)

	return func(groupId string) bool {
		isSecondParty, err := IsSecondPartyGroupId(groupId, secondPartyGroupId)
		log.Debugf("upgrade secondParty %t, groupId:%s isSecondParty %t", true, groupId, isSecondParty)
		if err != nil {
			log.Warnln(err)
			return false
		}
		return isSecondParty == true
	}
}

func Upgrade3PartyDependencies() func(pair PomWrapper, args ...interface{}) error {
	return func(pair PomWrapper, args ...interface{}) error {
		return Upgrade3PartyDependenciesOnModel(pair.Model)
	}
}

func Upgrade3PartyDependenciesOnModel(model *pom.Model) error {
	if model.Dependencies != nil {
		upgradeDependencies(model, model.Dependencies.Dependency, SecondParty(model, false))
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		upgradeDependencies(model, model.DependencyManagement.Dependencies.Dependency, SecondParty(model, false))
	}

	return nil
}

func Upgrade2PartyDependencies() func(pair PomWrapper, args ...interface{}) error {
	return func(pair PomWrapper, args ...interface{}) error {
		return Upgrade2PartyDependenciesOnModel(pair.Model)
	}
}

func Upgrade2PartyDependenciesOnModel(model *pom.Model) error {
	if model.Dependencies != nil {
		upgradeDependencies(model, model.Dependencies.Dependency, SecondParty(model, true))
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		upgradeDependencies(model, model.DependencyManagement.Dependencies.Dependency, SecondParty(model, true))
	}

	return nil
}

func upgradeDependencies(model *pom.Model, dependencies []pom.Dependency, condition func(groupId string) bool) {
	for _, dep := range dependencies {
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
