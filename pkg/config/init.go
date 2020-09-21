package config

import (
	"co-pilot/pkg/file"
	"github.com/mitchellh/go-homedir"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"strings"
)

func NewLocalConfig(localConfigDir string) (cfg LocalConfig, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return
	}

	cfg.impl.Path = file.Path("%s/%s", home, localConfigDir)
	return
}

func NewGitCloudConfig(localConfigDir string, cloudConfigDirName string) (cfg GitCloudConfig, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return cfg, err
	}

	cfg.Impl.Path = file.Path("%s/%s/%s", home, localConfigDir, cloudConfigDirName)
	return
}

func InitProjectConfigurationFromFile(filePath string) (config ProjectConfiguration, err error) {
	err = file.ReadJson(filePath, &config)
	err = config.Populate(strings.Replace(filePath, projectFileName, "", 1))
	return
}

func InitProjectConfigurationFromDir(targetDir string) (config ProjectConfiguration, err error) {
	filePath := file.Path("%s/%s", targetDir, projectFileName)
	err = file.ReadJson(filePath, &config)
	err = config.Populate(targetDir)
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
