package config

import (
	"co-pilot/pkg/file"
	"encoding/json"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"io/ioutil"
	"strings"
)

var projectFileName = "co-pilot.json"

type ProjectConfiguration struct {
	Language        string `json:"language"`
	GroupId         string `json:"groupId"`
	ArtifactId      string `json:"artifactId"`
	Package         string `json:"package"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ApplicationName string `json:"applicationName"`
	Team            struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"team"`
	Dependencies      []string          `json:"dependencies"`
	LocalDependencies []string          `json:"co-pilot-dependencies"`
	Render            map[string]string `json:"render"`
}

type ProjectConfig interface {
	Write(targetFile string) error
	SourceMainPath() string
	SourceTestPath() string
	FindApplicationName(targetDir string) (err error)
}

func InitProjectConfigurationFromFile(filePath string) (config ProjectConfiguration, err error) {
	err = file.ReadJson(filePath, &config)
	if config.ApplicationName == "" {
		err := config.FindApplicationName(strings.Replace(filePath, projectFileName, "", 1))
		if err != nil {
			return config, err
		}
	}
	return
}

func InitProjectConfigurationFromDir(targetDir string) (config ProjectConfiguration, err error) {
	filePath := fmt.Sprintf("%s/%s", targetDir, projectFileName)
	err = file.ReadJson(filePath, &config)

	if config.ApplicationName == "" {
		err := config.FindApplicationName(targetDir)
		if err != nil {
			return config, err
		}
	}
	return
}

func InitProjectConfigurationFromModel(model *pom.Model) (config ProjectConfiguration) {
	config.Language = "kotlin"
	config.GroupId = model.GetGroupId()
	config.ArtifactId = model.ArtifactId
	config.Package = model.GetGroupId()
	config.Name = model.Name
	config.Description = model.Description

	return
}

func (config ProjectConfiguration) Write(targetFile string) error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(targetFile, data, 0644)
}

func (config ProjectConfiguration) SourceMainPath() string {
	return fmt.Sprintf("%s", strings.Join(strings.Split(config.Package, "."), "/"))
}

func (config ProjectConfiguration) SourceTestPath() string {
	return fmt.Sprintf("%s", strings.Join(strings.Split(config.Package, "."), "/"))
}

func (config *ProjectConfiguration) FindApplicationName(targetDir string) (err error) {
	files, err := file.GrepRecursive(targetDir, "@SpringBootApplication")
	if err != nil {
		return
	}

	if len(files) == 1 {
		fileNamePath := strings.Split(files[0], "/")
		fileName := fileNamePath[len(fileNamePath)-1]
		fileNameParts := strings.Split(fileName, ".")
		config.ApplicationName = fileNameParts[0]
	}

	return
}
