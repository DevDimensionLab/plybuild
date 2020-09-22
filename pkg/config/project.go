package config

import (
	"co-pilot/pkg/file"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var projectConfigFileName = "co-pilot.json"

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
	Dependencies []string          `json:"dependencies"`
	Templates    []string          `json:"templates"`
	Settings     ProjectSettings   `json:"settings"`
	Render       map[string]string `json:"render"`
}

type ProjectSettings struct {
	DisableDependencySort bool `json:"disableDependencySort"`
}

type ProjectConfig interface {
	Write(targetFile string) error
	SourceMainPath() string
	SourceTestPath() string
	FindApplicationName(targetDir string) (err error)
	GetLanguage() string
	Populate(targetDir string) error
}

func (config ProjectConfiguration) Write(targetFile string) error {
	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(targetFile, data, 0644)
}

func (config ProjectConfiguration) SourceMainPath() string {
	return file.Path("src/main/%s/%s", config.GetLanguage(), strings.Join(strings.Split(config.Package, "."), "/"))
}

func (config ProjectConfiguration) SourceTestPath() string {
	return file.Path("src/test/%s/%s", config.GetLanguage(), strings.Join(strings.Split(config.Package, "."), "/"))
}

func (config *ProjectConfiguration) FindApplicationName(targetDir string) (err error) {
	files, err := file.GrepRecursive(targetDir, "@SpringBootApplication")
	if err != nil {
		log.Warnf("was not able to find application name in: %s", targetDir)
	}

	if len(files) == 1 {
		fileNamePath := strings.Split(files[0], "/")
		fileName := fileNamePath[len(fileNamePath)-1]
		fileNameParts := strings.Split(fileName, ".")
		config.ApplicationName = fileNameParts[0]
	}

	return
}

func (config *ProjectConfiguration) GetLanguage() string {
	if config.Language == "" || (config.Language != "kotlin" && config.Language != "java") {
		log.Warnf("language not set in config for package: %s, assuming kotlin...", config.Package)
		return "kotlin"
	}
	return config.Language
}

func (config *ProjectConfiguration) Populate(targetDir string) error {
	if config.ApplicationName == "" {
		err := config.FindApplicationName(targetDir)
		if err != nil {
			return err
		}
	}

	sourceTargetDir := file.Path("%s/src", targetDir)
	if config.Language == "" && file.Exists(sourceTargetDir) {
		kotlinFile, err := file.FindFirst(".kt", sourceTargetDir)
		if err == nil && kotlinFile != "" {
			log.Warnf("Language not set in %s, detected kotlin source files, setting language to kotlin",
				file.Path("%s/%s", targetDir, projectConfigFileName))
			config.Language = "kotlin"
			return nil
		}
		javaFile, err := file.FindFirst(".java", sourceTargetDir)
		if err == nil && javaFile != "" {
			log.Warnf("Language not set in %s, detected java source files, setting language to java",
				file.Path("%s/%s", targetDir, projectConfigFileName))
			config.Language = "java"
			return nil
		}

		return errors.New(fmt.Sprintf("%s directory detected, but language was not set in co-pilot.json",
			file.Path("%s/src", targetDir)))
	}

	return nil
}
