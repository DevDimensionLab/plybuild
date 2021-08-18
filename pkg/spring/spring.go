package spring

import (
	"errors"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/maven"
	"github.com/co-pilot-cli/mvn-pom-mutator/pkg/pom"
)

func CleanManualVersions() func(project config.Project, args ...interface{}) error {
	return func(project config.Project, args ...interface{}) error {
		return cleanManualVersions(project.Type.Model())
	}
}

func cleanManualVersions(model *pom.Model) error {
	springBootDependencies, err := GetDependencies()
	if err != nil {
		return err
	}

	if model.Dependencies != nil {
		err = removeVersion(model.Dependencies.Dependency, springBootDependencies, model)
		if err != nil {
			return err
		}
	}

	if model.DependencyManagement != nil && model.DependencyManagement.Dependencies != nil {
		err = removeVersion(model.DependencyManagement.Dependencies.Dependency, springBootDependencies, model)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpgradeSpringBoot() func(project config.Project, args ...interface{}) error {
	return func(project config.Project, args ...interface{}) error {
		if project.Config.Settings.DisableSpringBootUpgrade {
			return nil
		}
		return upgradeSpringBootOnModel(project.Type.Model())
	}
}

func upgradeSpringBootOnModel(model *pom.Model) error {
	latestVersionMeta, err := maven.GetMetaData("org.springframework.boot", "spring-boot")
	if err != nil {
		return err
	}

	latestVersion, err := latestVersionMeta.LatestRelease()
	if err != nil {
		return err
	}

	currentVersion, err := getSpringBootVersion(model)
	if err != nil {
		return err
	}

	if currentVersion.IsDifferentFrom(latestVersion) {
		msg := fmt.Sprintf("outdated spring-boot version [%s => %s]", currentVersion.ToString(), latestVersion.ToString())
		if maven.IsMajorUpgrade(currentVersion, latestVersion) {
			log.Warnf("major %s", msg)
		} else if !latestVersion.IsReleaseVersion() {
			log.Warnf("%s | not release", msg)
		} else {
			log.Info(msg)
		}

		return updateSpringBootVersion(model, latestVersion)
	} else {
		log.Infof("spring boot is the latest version [%s]", latestVersion.ToString())
	}

	return nil
}

func getSpringBootVersion(model *pom.Model) (maven.JavaVersion, error) {
	// check parent
	if model.Parent != nil && model.Parent.ArtifactId == "spring-boot-starter-parent" {
		return maven.ParseVersion(model.Parent.Version)
	}

	// check dependencyManagement
	if model.DependencyManagement != nil {
		dep, err := model.DependencyManagement.Dependencies.FindArtifact("spring-boot-dependencies")
		if err != nil {
			return maven.JavaVersion{}, err
		}
		version, err := model.GetDependencyVersion(dep)
		if err != nil {
			return maven.JavaVersion{}, err
		}
		return maven.ParseVersion(version)
	}

	return maven.JavaVersion{}, errors.New("could not extract spring boot version information")
}

func updateSpringBootVersion(model *pom.Model, newVersion maven.JavaVersion) error {
	// check parent
	if model.Parent != nil && model.Parent.ArtifactId == "spring-boot-starter-parent" {
		model.Parent.Version = newVersion.ToString()
		return nil
	}

	// check dependencyManagement
	if model.DependencyManagement != nil {
		dep, err := model.DependencyManagement.Dependencies.FindArtifact("spring-boot-dependencies")
		if err != nil {
			return err
		} else {
			return model.SetDependencyVersion(dep, newVersion.ToString())
		}
	}

	return errors.New("could not update spring boot version to " + newVersion.ToString())
}

func removeVersion(dependencies []pom.Dependency, springBootDependencies IoDependenciesResponse, model *pom.Model) error {
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

func inMap(dep pom.Dependency, springBootDeps map[string]Dependency) bool {
	for _, springBootDep := range springBootDeps {
		if springBootDep.GroupId == dep.GroupId && springBootDep.ArtifactId == dep.ArtifactId {
			return true
		}
	}

	return false
}
