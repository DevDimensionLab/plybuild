package spring

import (
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/mvn_crud"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"io/ioutil"
	"strings"
)

func Upgrade(directory string) error {
	pomFile := directory + "/pom.xml"
	model, err := mvn_crud.GetPomModel(pomFile)
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

		err = writePomModel(model, pomFile+".new")
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("No update needed, project version is the newest of spring boot [%s]\n", newestVersion)
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
		for _, dep := range model.DependencyManagement.Dependencies.Dependency {
			if dep.ArtifactId == "spring-boot-dependencies" {
				if strings.HasPrefix(dep.Version, "${") {
					for _, a := range model.Properties.AnyElements {
						if a.XMLName.Local == strings.Trim(dep.Version, "${}") {
							return a.Value, nil
						}
					}
				} else {
					return dep.Version, nil
				}
			}
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
		for _, dep := range model.DependencyManagement.Dependencies.Dependency {
			if dep.ArtifactId == "spring-boot-dependencies" {
				if strings.HasPrefix(dep.Version, "${") {
					for i, a := range model.Properties.AnyElements {
						if a.XMLName.Local == strings.Trim(dep.Version, "${}") {
							model.Properties.AnyElements[i].Value = newestVersion
							return nil
						}
					}
				} else {
					dep.Version = newestVersion
					return nil
				}
			}
		}
	}

	return errors.New("could not update spring boot version to " + newestVersion)
}

func writePomModel(model *pom.Model, outputFile string) error {
	bytes, err := mvn_crud.Marshall(model)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFile, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
