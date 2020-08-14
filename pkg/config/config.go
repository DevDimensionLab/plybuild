package config

import (
	"encoding/json"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"io/ioutil"
)

func DefaultConfiguration() InitConfiguration {
	return InitConfiguration{
		GroupId:     "com.example.demo",
		ArtifactId:  "demo-webservice",
		Package:     "com.example.demo",
		Name:        "webservice",
		Description: "demo webservice",
	}
}

func GenerateConfig(model *pom.Model) (InitConfiguration, error) {
	// needs to be implemented correctly...
	return InitConfiguration{
		GroupId:     model.GroupId,
		ArtifactId:  model.ArtifactId,
		Package:     model.GroupId,
		Name:        model.Name,
		Description: model.Description,
	}, nil
}

func WriteConfig(configuration InitConfiguration, targetFile string) error {
	data, err := json.MarshalIndent(configuration, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(targetFile, data, 0644)
}
