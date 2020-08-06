package spring

import (
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
)

func UpgradeSpringBoot(directory string) error {
	pomFile := directory + "/pom.alt2.xml"
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

func UpgradeSpringDependencies(directory string) error {
	pomFile := directory + "/pom.xml"
	project, err := pom.GetModelFrom(pomFile)
	if err != nil {
		return err
	}

	springBootDependencies, err := GetDependencies()
	if err != nil {
		return err
	}

	deps := getSpringBootDependenciesFromProject(project, springBootDependencies.Dependencies)

	for _, dep := range deps {
		fmt.Printf("Found spring-boot dependecy: %s:%s, on version: [%s] \n", dep.GroupId, dep.ArtifactId, dep.Version)
	}

	return errors.New("[NOT IMPLEMENTED] could not update any spring boot dependencies")
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
			return dep.GetVersion(model)
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
			return dep.SetVersion(model, newestVersion)
		}
	}

	return errors.New("could not update spring boot version to " + newestVersion)
}

func getSpringBootDependenciesFromProject(project *pom.Model, springBootDependencies map[string]Dependency) []pom.Dependency {

	var foundDependencies []pom.Dependency

	if project.Dependencies != nil {
		for _, projectDep := range project.Dependencies.Dependency {
			for _, bootDep := range springBootDependencies {
				if projectDep.ArtifactId == bootDep.ArtifactId {
					foundDependencies = append(foundDependencies, projectDep)
				}
			}
		}
	}

	return foundDependencies
}
