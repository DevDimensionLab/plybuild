package upgrade

import (
	"co-pilot/pkg/analyze"
	"co-pilot/pkg/maven"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func Dependencies(model *pom.Model, secondParty bool) error {
	secondPartyGroupId, err := analyze.GetSecondPartyGroupId(model)
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
			isSecondParty, err := analyze.IsSecondPartyGroupId(dep.GroupId, secondPartyGroupId)
			if err == nil && isSecondParty == secondParty {
				err = DependencyUpgrade(model, dep)
				if err != nil {
					log.Errorf("%v", err)
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
	if err == nil {
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

			_ = model.SetDependencyVersion(dep, latestVersion.ToString())
		}
		return nil
	} else {
		return err
	}
}
