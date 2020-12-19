package spring

import (
	"errors"
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/file"
	"github.com/co-pilot-cli/co-pilot/pkg/http"
	"github.com/co-pilot-cli/co-pilot/pkg/shell"
	"net/url"
	"os"
	"strings"
	"time"
)

func UrlValuesFrom(config config.ProjectConfiguration) url.Values {
	// see https://github.com/spring-io/initializr#generating-a-project
	params := url.Values{}
	params.Add("groupId", config.GroupId)
	params.Add("artifactId", config.ArtifactId)
	params.Add("packageName", config.Package)
	params.Add("dependencies", strings.Join(config.Dependencies, ","))
	params.Add("javaVersion", "11")
	params.Add("language", config.Language)
	params.Add("description", config.Description)
	params.Add("name", config.Name)
	//params.Add("baseDir", targetDir)

	return params
}

func GetRoot() (IoRootResponse, error) {
	var deps IoRootResponse
	err := http.GetJson("https://start.spring.io", &deps)
	return deps, err
}

func GetDependencies() (IoDependenciesResponse, error) {
	var deps IoDependenciesResponse
	err := http.GetJson("https://start.spring.io/dependencies", &deps)
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
	targetArchiveFile, err := archivePath()
	if err != nil {
		return err
	}

	downloadUrl := "https://start.spring.io/starter.zip"
	log.Infof("Downloading from %s to %s", downloadUrl, targetArchiveFile)
	err = http.Wpost(downloadUrl, targetArchiveFile, formData)
	if err != nil {
		return err
	}

	_, err = shell.Unzip(targetArchiveFile, targetDir)
	if err != nil {
		return err
	}

	log.Debugf("Deleting archive file: %s", targetArchiveFile)
	err = file.DeleteSingleFile(targetArchiveFile)
	return err
}

func archivePath() (path string, err error) {
	curDir, err := os.Getwd()
	if err != nil {
		return
	}

	now := time.Now().Unix()
	path = file.Path("%s/spring-%d.zip", curDir, now)
	return
}
