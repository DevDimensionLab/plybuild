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

func FromProject(target string) (ProjectConfiguration, error) {
	var config ProjectConfiguration
	err := file.ReadJson(target+"/co-pilot.json", &config)
	if err != nil {
		return config, err
	}

	return config, nil
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

func WriteConfig(configuration ProjectConfiguration, targetFile string) error {
	data, err := json.MarshalIndent(configuration, "", " ")
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