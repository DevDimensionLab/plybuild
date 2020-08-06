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

	projectVersion, err := getSpringBootVersion(model)
	if err != nil {
		return err
	}

	newestVersion := springRootInfo.BootVersion.Default

	if projectVersion != newestVersion {
		fmt.Printf("=> Upgrade needed !!! \n")
		fmt.Printf("Project uses spring boot version: %s\n", projectVersion)
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
		fmt.Printf("No update needed, project version is the newest of spring boot [%s]\n", newestVersion)
	}

	return nil
}

func UpgradeDependencies(directory string) error {
	pomFile := directory + "/pom.xml"
	project, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	deps := getDependenciesFromProject(project)

	for _, dep := range deps {
		if dep.Version != "" {
			currentVersion, err := project.GetVersion(dep)
			metaData, err := maven.GetMetaData(dep.GroupId, dep.ArtifactId)
			if err == nil {
				fmt.Printf("[INFO] Found: %s:%s with version: [%s], latest version is: [%s] \n", dep.GroupId, dep.ArtifactId, currentVersion, metaData.Versioning.Latest)
			} else {
				fmt.Printf("[ERROR] %v\n", err)
			}
		}
	}

	return errors.New("[NOT IMPLEMENTED] could not update any dependencies")
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
		for _, projectDep := range model.Dependencies.Dependency {
			foundDependencies = append(foundDependencies, projectDep)
		}
	}

	return foundDependencies
}
