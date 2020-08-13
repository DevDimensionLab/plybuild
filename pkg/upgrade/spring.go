package upgrade

import (
	"co-pilot/pkg/springio"
	"errors"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
)

func SpringBoot(model *pom.Model) error {
	springRootInfo, err := springio.GetRoot()
	if err != nil {
		return err
	}

	modelVersion, err := getSpringBootVersion(model)
	if err != nil {
		return err
	}

	latestVersion := springRootInfo.BootVersion.Default

	if modelVersion != latestVersion {
		log.Warnf("outdated spring-boot version [%s => %s]", modelVersion, latestVersion)
		err = updateSpringBootVersion(model, latestVersion)
		if err != nil {
			return err
		}
	} else {
		log.Infof("Spring boot is the latest version [%s]", latestVersion)
	}

	return nil
}

func getSpringBootVersion(model *pom.Model) (string, error) {
	// check parent
	if model.Parent.ArtifactId == "spring-boot-starter-parent" {
		return model.Parent.Version, nil
	}

	// check dependencyManagement
	if model.DependencyManagement != nil {
		dep, err := model.DependencyManagement.Dependencies.FindArtifact("spring-boot-dependencies")
		if err != nil {
			return "", nil
		} else {
			return model.GetDependencyVersion(dep)
		}
	}

	return "", errors.New("could not extract spring boot version information")
}

func updateSpringBootVersion(model *pom.Model, newestVersion string) error {
	// check parent
	if model.Parent.ArtifactId == "spring-boot-starter-parent" {
		model.Parent.Version = newestVersion
		return nil
	}

	// check dependencyManagement
	if model.DependencyManagement != nil {
		dep, err := model.DependencyManagement.Dependencies.FindArtifact("spring-boot-dependencies")
		if err != nil {
			return err
		} else {
			return model.SetDependencyVersion(dep, newestVersion)
		}
	}

	return errors.New("could not update spring boot version to " + newestVersion)
}
