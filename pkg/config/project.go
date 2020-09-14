package config

import (
	"co-pilot/pkg/file"
	"encoding/json"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"io/ioutil"
	"strings"
)

func DefaultConfiguration() ProjectConfiguration {
	return ProjectConfiguration{
		Language:          "kotlin",
		GroupId:           "com.example.demo",
		ArtifactId:        "demo-webservice",
		Package:           "com.example.demo",
		Name:              "webservice",
		Description:       "demo webservice",
		Dependencies:      []string{},
		LocalDependencies: []string{},
	}
}

func FromProject(targetDir string) (config ProjectConfiguration, err error) {
	err = file.ReadJson(targetDir+"/co-pilot.json", &config)
	if err != nil {
		return
	}

	if config.ApplicationName == "" {
		// populate applicationName field from targetDir
		appName, err := FindApplicationName(targetDir)
		if err != nil {
			return config, err
		} else {
			config.ApplicationName = appName
		}
	}
	return
}

func GenerateConfig(model *pom.Model) (ProjectConfiguration, error) {
	// needs to be implemented correctly...
	groupId := model.GroupId
	if groupId == "" {
		groupId = model.Parent.GroupId
	}
	return ProjectConfiguration{
		Language:          "kotlin",
		GroupId:           groupId,
		ArtifactId:        model.ArtifactId,
		Package:           groupId,
		Name:              model.Name,
		Description:       model.Description,
		Dependencies:      []string{},
		LocalDependencies: []string{},
	}, nil
}

func (config ProjectConfiguration) WriteConfig(targetFile string) error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(targetFile, data, 0644)
}

func (config ProjectConfiguration) ProjectMainRoot() string {
	return fmt.Sprintf("%s", strings.Join(strings.Split(config.Package, "."), "/"))
}

func (config ProjectConfiguration) ProjectTestRoot() string {
	return fmt.Sprintf("%s", strings.Join(strings.Split(config.Package, "."), "/"))
}

func FindApplicationName(targetDir string) (applicationName string, err error) {
	files, err := file.Recursive(targetDir, "@SpringBootApplication")
	if err != nil {
		return
	}

	if len(files) == 1 {
		fileNamePath := strings.Split(files[0], "/")
		fileName := fileNamePath[len(fileNamePath)-1]
		fileNameParts := strings.Split(fileName, ".")
		applicationName = fileNameParts[0]
	}

	return
}
