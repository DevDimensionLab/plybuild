package spring

import (
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"spring-boot-co-pilot/pkg/maven"
)

func UpgradeSpringBoot(directory string) error {
	pomFile := directory + "/pom.xml"
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	springRootInfo, err := GetRoot()
	if err != nil {
		return err
	}

	modelVersion, err := getSpringBootVersion(model)
	if err != nil {
		return err
	}

	newestVersion := springRootInfo.BootVersion.Default

	if modelVersion != newestVersion {
		fmt.Printf("=> Upgrade needed !!! \n")
		fmt.Printf("Project uses spring boot version: %s\n", modelVersion)
		fmt.Printf("Latest version of spring boot: %s\n", newestVersion)
		err = updateSpringBootVersion(model, newestVersion)
		if err != nil {
			return err
		}

		err = model.WriteToFile(pomFile + ".new")
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("No update needed, model version is the newest of spring boot [%s]\n", newestVersion)
	}

	return nil
}

func UpgradeDependencies(directory string) error {
	pomFile := directory + "/pom.xml"
	model, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	deps := getDependenciesFromProject(model)

	for _, dep := range deps {
		if dep.Version != "" {
			currentVersion, err := model.GetVersion(dep)
			metaData, err := maven.GetMetaData(dep.GroupId, dep.ArtifactId)
			if err == nil {
				fmt.Printf("[OUTDATED] %s:%s [%s => %s] \n", dep.GroupId, dep.ArtifactId, currentVersion, metaData.Versioning.Latest)
				_ = model.SetVersion(dep, metaData.Versioning.Latest)
			} else {
				fmt.Printf("[ERROR] %v\n", err)
			}
		}
	}

	return model.WriteToFile(pomFile + ".new")
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
			return model.GetVersion(dep)
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
			return model.SetVersion(dep, newestVersion)
		}
	}

	return errors.New("could not update spring boot version to " + newestVersion)
}

func getDependenciesFromProject(model *pom.Model) []pom.Dependency {

	var foundDependencies []pom.Dependency

	if model.Dependencies != nil {
		for _, modelDep := range model.Dependencies.Dependency {
			foundDependencies = append(foundDependencies, modelDep)
		}
	}

	return foundDependencies
}
