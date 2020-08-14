package springio

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/http"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CLI(springExec string, arg ...string) error {
	cmd := exec.Command(springExec, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InitFrom(config config.InitConfiguration) []string {
	var output []string

	output = append(output, "init")
	output = append(output, "-g="+config.GroupId)
	output = append(output, "-a="+config.ArtifactId)
	output = append(output, "--package="+config.Package)
	output = append(output, "--d="+strings.Join(config.Dependencies, ","))
	output = append(output, "-j=11")
	output = append(output, "--language=kotlin")
	output = append(output, "--description="+config.Description)
	output = append(output, "--name="+config.Name)
	output = append(output, "webservice") // outputdir

	return output
}

func GetRoot() (IoRootResponse, error) {
	var deps IoRootResponse
	err := http.GetJson("http://start.spring.io", &deps)
	return deps, err
}

func GetInfo() (IoInfoResponse, error) {
	var deps IoInfoResponse
	err := http.GetJson("http://start.spring.io/actuator/info", &deps)
	return deps, err
}

func GetDependencies() (IoDependenciesResponse, error) {
	var deps IoDependenciesResponse
	err := http.GetJson("http://start.spring.io/dependencies", &deps)
	return deps, err
}

func Validate(config config.InitConfiguration) error {
	var invalidDependencies []string
	validDependencies, err := GetDependencies()
	if err != nil {
		return err
	}

	for _, userDefinedDependency := range config.Dependencies {
		valid := false
		for validDependency, _ := range validDependencies.Dependencies {
			if validDependency == userDefinedDependency {
				valid = true
			}
		}
		if !valid {
			invalidDependencies = append(invalidDependencies, userDefinedDependency)
		}
	}

	if len(invalidDependencies) > 0 {
		validKeys := make([]string, 0, len(validDependencies.Dependencies))
		for k, _ := range validDependencies.Dependencies {
			validKeys = append(validKeys, k)
		}
		return errors.New(fmt.Sprintf("%s not found in valid list of dependencies %s", invalidDependencies, validKeys))
	} else {
		return nil
	}
}
