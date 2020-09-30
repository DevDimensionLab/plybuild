package config

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/shell"
	"co-pilot/pkg/sorting"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"io/ioutil"
	"sort"
	"strings"
)

var projectConfigFileName = "co-pilot.json"

type Project struct {
	Path       string
	GitInfo    GitInfo
	ConfigFile string
	Config     ProjectConfiguration
	Type       ProjectType
}

type ValidProjectType string

const (
	Maven ValidProjectType = "Maven"
)

// only maven for now...
type ProjectType interface {
	Type() ValidProjectType
	FilePath() string
	Model() *pom.Model
}

type MavenProject struct {
	PomFile  string
	PomModel *pom.Model
}

func (mvnProject MavenProject) FilePath() string {
	return mvnProject.PomFile
}

func (mvnProject MavenProject) Model() *pom.Model {
	return mvnProject.PomModel
}

func (mvnProject MavenProject) Type() ValidProjectType {
	return Maven
}

type ProjectConfiguration struct {
	MavenProjectConfiguration
	Name        string `json:"name"`
	Description string `json:"description"`
	Team        struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"team"`
	Dependencies []string          `json:"dependencies"`
	Templates    []string          `json:"templates"`
	Settings     ProjectSettings   `json:"settings"`
	Render       map[string]string `json:"render"`
}

type MavenProjectConfiguration struct {
	Language        string `json:"language"`
	GroupId         string `json:"groupId"`
	ArtifactId      string `json:"artifactId"`
	Package         string `json:"package"`
	ApplicationName string `json:"applicationName"`
}

type ProjectSettings struct {
	DisableDependencySort bool `json:"disableDependencySort"`
}

type ProjectConfig interface {
	WriteTo(targetFile string) error
	SourceMainPath() string
	SourceTestPath() string
	FindApplicationName(targetDir string) (err error)
	GetLanguage() string
	Populate(targetDir string) error
}

func (config ProjectConfiguration) WriteTo(targetFile string) error {
	log.Infof("writes project config file to %s", targetFile)
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

func (config ProjectConfiguration) Empty() bool {
	return config.Name == "" ||
		config.Language == "" ||
		config.Package == "" ||
		config.GroupId == "" ||
		config.ArtifactId == ""
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
				ProjectConfigPath(targetDir))
			config.Language = "kotlin"
			return nil
		}
		javaFile, err := file.FindFirst(".java", sourceTargetDir)
		if err == nil && javaFile != "" {
			log.Warnf("Language not set in %s, detected java source files, setting language to java",
				ProjectConfigPath(targetDir))
			config.Language = "java"
			return nil
		}

		return errors.New(fmt.Sprintf("%s directory detected, but language was not set in %s",
			file.Path("%s/src", targetDir), projectConfigFileName))
	}

	return nil
}

func (project Project) IsMavenProject() bool {
	return project.Type != nil && project.Type.Type() == Maven
}

func (project Project) IsGitRepo() bool {
	return project.GitInfo.IsRepo
}

func (project Project) IsDirtyGitRepo() bool {
	return project.GitInfo.IsRepo && project.GitInfo.IsDirty
}

func (project Project) GitInit() error {
	if project.GitInfo.DisableCommit {
		return nil
	}
	if !project.GitInfo.IsRepo {
		init := shell.GitInit(project.Path)
		if init.Err != nil {
			return init.FormatError()
		}
	}

	return project.GitCommit("Initial commit")
}

func (project Project) GitCommit(message string) error {
	if project.GitInfo.DisableCommit {
		return nil
	}
	cmd := shell.GitAddAndCommit(project.Path, message)

	if cmd.Err != nil {
		return cmd.FormatError()
	}

	return nil
}

func ProjectConfigPath(targetDir string) string {
	return file.Path("%s/%s", targetDir, projectConfigFileName)
}

func GetGitInfoFromPath(targetDir string) (gitInfo GitInfo, err error) {
	isRepo, err := shell.GitIsRepo(targetDir)
	if err != nil {
		return
	}
	gitInfo.IsRepo = isRepo

	dirty, err := shell.GitDirty(targetDir)
	if err != nil {
		return
	}
	gitInfo.IsDirty = dirty
	return
}

func (project Project) SortAndWritePom() error {
	var disableDepSort = project.Config.Settings.DisableDependencySort

	if project.Type.Model().Dependencies != nil && !disableDepSort {
		sortKey, err := project.Type.Model().GetSecondPartyGroupId()
		if err != nil {
			log.Warn(err)
		} else {
			log.Infof("sorting pom file with default dependencySort")
			sort.Sort(sorting.DependencySort{
				Deps:    project.Type.Model().Dependencies.Dependency,
				SortKey: sortKey})
		}
	}

	var writeToFile = project.Type.FilePath()
	log.Infof("writing model to pom file: %s", writeToFile)
	return project.Type.Model().WriteToFile(writeToFile)
}
