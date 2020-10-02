package spring

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/http"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/shell"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var log = logger.Context()

func UrlValuesFrom(config config.ProjectConfiguration) url.Values {
	// see https://github.com/spring-io/initializr#generating-a-project
	params := url.Values{}
	params.Add("groupId", config.GroupId)
	params.Add("artifactId", config.ArtifactId)
	params.Add("packageName", config.Package)
	params.Add("dependencies", strings.Join(config.Dependencies, ","))
	params.Add("javaVersion", "11")
	params.Add("language", "kotlin")
	params.Add("description", config.Description)
	params.Add("name", config.Name)
	//params.Add("baseDir", targetDir)

	return params
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

func Validate(config config.ProjectConfiguration) error {
	if config.Dependencies == nil || len(config.Dependencies) == 0 {
		return nil
	}

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

func DownloadInitializer(targetDir string, formData url.Values) error {
	targetFile := "/tmp/spring.zip"
	downloadUrl := "https://start.spring.io/starter.zip"
	log.Infof("Downloading from %s to %s", downloadUrl, targetFile)
	err := http.Wpost(downloadUrl, targetFile, formData)
	if err != nil {
		return err
	}

	_, err = shell.Unzip(targetFile, targetDir)
	return err
}
