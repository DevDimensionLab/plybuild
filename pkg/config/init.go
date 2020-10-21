package config

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/sirupsen/logrus"
	"strings"
)

var log = logger.Context()

func SetLogger(logger logrus.FieldLogger) {
	log = logger
}

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

func (project *Project) InitProjectConfiguration() (err error) {
	if project.Type == nil || project.Type.Model() == nil {
		return errors.New("project type and model is nil")
	}

	if !project.Config.Empty() {
		return
	}

	packageName, err := findProjectPackageName(project.Path)
	if err != nil {
		log.Warnf("Failed to get package name of root level source file %v", err)
	}

	model := project.Type.Model()
	project.Config.GroupId = model.GetGroupId()
	project.Config.ArtifactId = model.ArtifactId
	project.Config.Package = packageName
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

func findProjectPackageName(path string) (packageName string, err error) {
	packageName, err = findRootSourceFilePackageName(".kt", path)
	if err == nil {
		return packageName, err
	}
	packageName, err = findRootSourceFilePackageName(".java", path)
	if err == nil {
		return packageName, err
	}

	return
}

func findRootSourceFilePackageName(suffix string, path string) (packageName string, err error) {
	files, err := file.FindAll(suffix, []string{}, path)
	if err != nil {
		return
	}

	if len(files) == 0 {
		return packageName, errors.New(fmt.Sprintf("no files with suffix: %s in %s", suffix, path))
	}

	lines, err := file.OpenLines(files[0])
	if err != nil {
		return packageName, err
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "package") {
			packageParts := strings.Split(line, " ")
			return packageParts[1], nil
		}
	}

	return "",
		errors.New(fmt.Sprintf("failed to get any files with suffix %s and a 'package' line in %s", suffix, path))
}
