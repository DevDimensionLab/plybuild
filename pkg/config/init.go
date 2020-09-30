package config

import (
	"co-pilot/pkg/file"
	"errors"
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
	if err != nil {
		return
	}
	err = config.Populate(strings.Replace(filePath, projectConfigFileName, "", 1))
	return
}

func InitProjectConfigurationFromDir(targetDir string) (config ProjectConfiguration, err error) {
	filePath := file.Path("%s/%s", targetDir, projectConfigFileName)
	err = file.ReadJson(filePath, &config)
	err = config.Populate(targetDir)
	return
}

func (project Project) InitProjectConfiguration() (err error) {
	if project.Type == nil || project.Type.Model() == nil {
		return errors.New("project type and model is nil")
	}

	if !project.Config.Empty() {
		return
	}

	model := project.Type.Model()
	project.Config.GroupId = model.GetGroupId()
	project.Config.ArtifactId = model.ArtifactId
	project.Config.Package = model.GetGroupId()
	project.Config.Name = model.Name
	project.Config.Description = model.Description
	err = project.Config.Populate(project.Path)
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
		project.Type = MavenProject{
			PomFile:  pomFile,
			PomModel: pomModel,
		}
	}

	project.ConfigFile = file.Path("%s/%s", targetDir, projectConfigFileName)
	project.Path = targetDir
	project.Config = config
	return
}
