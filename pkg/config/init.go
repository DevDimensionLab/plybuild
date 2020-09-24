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
	log.Debugf("loading projectConfig: %s", filePath)
	err = file.ReadJson(filePath, &config)
	err = config.Populate(strings.Replace(filePath, projectConfigFileName, "", 1))
	return
}

func InitProjectConfigurationFromDir(targetDir string) (config ProjectConfiguration, err error) {
	filePath := file.Path("%s/%s", targetDir, projectConfigFileName)
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

func InitProjectFromPomFile(pomFile string) (project Project, err error) {
	targetDir := file.Path(strings.Replace(pomFile, "pom.xml", "", 1))
	return InitProjectFromDirectory(targetDir)
}

func InitProjectFromDirectory(targetDir string) (project Project, err error) {
	gitInfo, err := GetGitInfoFromPath(targetDir)
	if err != nil {
		log.Debugln(err)
	} else {
		project.GitInfo = gitInfo
	}

	config, err := InitProjectConfigurationFromDir(targetDir)
	if err != nil {
		return
	}

	pomFile := file.Path("%s/pom.xml", targetDir)
	if file.Exists(pomFile) {
		pomModel, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Warnln(err)
		}
		project.PomFile = pomFile
		project.PomModel = pomModel
	}

	project.ConfigFile = file.Path("%s/%s", targetDir, projectConfigFileName)
	project.Path = targetDir
	project.Config = config
	return
}
