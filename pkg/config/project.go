package config

import (
	"co-pilot/pkg/file"
	"encoding/json"
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
	Dependencies []string          `json:"dependencies"`
	Templates    []string          `json:"templates"`
	Render       map[string]string `json:"render"`
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

	if config.Language == "" {
		kotlinFile, err := file.FindFirst(".kt", targetDir)
		if err == nil && kotlinFile != "" {
			config.Language = "kotlin"
			return nil
		}
		javaFile, err := file.FindFirst(".java", targetDir)
		if err == nil && javaFile != "" {
			config.Language = "java"
			return nil
		}

		// if all fails, fallback to kotlin
		config.Language = "kotlin"
	}

	return nil
}
