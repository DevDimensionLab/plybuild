package upgrade

import (
	"co-pilot/pkg/maven"
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

func Dependency(model *pom.Model, groupId string, artifactId string) (err error) {
	if model.Dependencies != nil {
		err = SpecificDependencyUpgrade(model, model.Dependencies.Dependency, groupId, artifactId)
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		err = SpecificDependencyUpgrade(model, model.DependencyManagement.Dependencies.Dependency, groupId, artifactId)
	}

	return err
}

func SpecificDependencyUpgrade(model *pom.Model, availableDependencies []pom.Dependency, groupId string, artifactId string) error {
	for _, dep := range availableDependencies {
		if dep.Version != "" && dep.GroupId == groupId && dep.ArtifactId == artifactId {
			return DependencyUpgrade(model, dep)
		}
	}

	return errors.New(fmt.Sprintf("could not find %s:%s in project", groupId, artifactId))
}

func Dependencies(model *pom.Model, secondParty bool) error {
	secondPartyGroupId, err := maven.GetSecondPartyGroupId(model)
	if err != nil {
		return err
	}

	if model.Dependencies != nil {
		DependenciesUpgrade(model.Dependencies.Dependency, secondPartyGroupId, secondParty, model)
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		DependenciesUpgrade(model.DependencyManagement.Dependencies.Dependency, secondPartyGroupId, secondParty, model)
	}

	return nil
}

func DependenciesUpgrade(dependencies []pom.Dependency, secondPartyGroupId string, secondParty bool, model *pom.Model) {
	for _, dep := range dependencies {
		if dep.Version != "" {
			isSecondParty, err := maven.IsSecondPartyGroupId(dep.GroupId, secondPartyGroupId)
			if err == nil && isSecondParty == secondParty {
				err = DependencyUpgrade(model, dep)
				if err != nil {
					log.Warnf("%v", err)
				}
			}
		}
	}
}

func DependencyUpgrade(model *pom.Model, dep pom.Dependency) error {
	if dep.Version == "${project.version}" || dep.Version == "${revision}" {
		return nil
	}

	depVersion, err := model.GetDependencyVersion(dep)
	if err != nil {
		return err
	}

	currentVersion, err := maven.ParseVersion(depVersion)
	if err != nil {
		return err
	}

	metaData, err := maven.GetMetaData(dep.GroupId, dep.ArtifactId)
	if err != nil {
		return err
	}

	latestVersion, err := metaData.LatestRelease()
	if err != nil {
		return nil
	}

	if currentVersion.IsDifferentFrom(latestVersion) {
		msg := fmt.Sprintf("outdated dependency %s:%s [%s => %s]", dep.GroupId, dep.ArtifactId, currentVersion.ToString(), latestVersion.ToString())
		if maven.IsMajorUpgrade(currentVersion, latestVersion) {
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
