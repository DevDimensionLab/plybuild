package upgrade

import (
	"co-pilot/pkg/analyze"
	"co-pilot/pkg/maven"
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func Dependencies(model *pom.Model, secondParty bool) error {
	secondPartyGroupId, err := analyze.GetSecondPartyGroupId(model)
	if err != nil {
		return err
	}

	if model.Dependencies == nil {
		return errors.New("could not find any dependencies")
	}

	for _, dep := range model.Dependencies.Dependency {
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

	return nil
}

func DependencyUpgrade(model *pom.Model, dep pom.Dependency) error {
	depVersion, err := model.GetDependencyVersion(dep)
	if err != nil {
		return err
	}

	currentVersion, err := ParseVersion(depVersion)
	if err != nil {
		return err
	}

	metaData, err := maven.GetMetaData(dep.GroupId, dep.ArtifactId)
	if err == nil {
		latestVersion, err := ParseVersion(metaData.Versioning.Release)
		if err != nil {
			return nil
		}

		if currentVersion.IsDifferentFrom(latestVersion) {
			msg := fmt.Sprintf("outdated dependency %s:%s [%s => %s] \n", dep.GroupId, dep.ArtifactId, currentVersion.ToString(), latestVersion.ToString())
			if IsMajorUpgrade(currentVersion, latestVersion) {
				log.Warnf("major %s", msg)
			} else if !latestVersion.IsReleaseVersion() {
				log.Warnf("%s | not release", msg)
			} else {
				log.Info(msg)
			}

			_ = model.SetDependencyVersion(dep, metaData.Versioning.Release)
		}
		return nil
	} else {
		return err
	}
}
