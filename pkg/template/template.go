package template

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/maven"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var defaultIgnores = []string{
	"pom.xml",
	"co-pilot.json",
	"Application",
	".co-pilot.ignore",
	".gitignore",
	".mvn",
	"mvnw",
	"mvnw.cmd",
	".idea",
	".iml",
}

func MergeTemplates(templates []config.CloudTemplate, target config.Project) {
	for _, template := range templates {
		log.Infof("applying Template %s", template.Name)
		if err := MergeTemplate(template, target); err != nil {
			log.Warnf("%v", err)
		}
	}
}

func MergeTemplate(cloudTemplate config.CloudTemplate, target config.Project) error {
	if target.IsDirtyGitRepo() {
		log.Warn(logger.White(fmt.Sprintf("merging Template %s into a dirty git repository %s", cloudTemplate.Name, target.Path)))
	} else {
		log.Info(logger.White(fmt.Sprintf("merging Template %s into %s", cloudTemplate.Name, target.Path)))
	}
	if err := merge(cloudTemplate.Project, target); err != nil {
		return err
	}

	return nil
}

func merge(sourceProject config.Project, targetProject config.Project) error {
	sourceDir := sourceProject.Path
	files, err := filesToCopy(sourceDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		sourceRelPath, err := file.RelPath(sourceDir, f)
		if err != nil {
			return err
		}

		sourceRelPath = replacePathForSource(sourceRelPath, sourceProject.Config, targetProject.Config)

		targetPath := file.Path("%s/%s", targetProject.Path, sourceRelPath)
		if err = file.CopyOrMerge(f, targetPath); err != nil {
			return err
		}

		if err = file.SearchReplace(targetPath, sourceProject.Config.Package, targetProject.Config.Package); err != nil {
			return err
		}

		if strings.HasSuffix(targetPath, ".render") {
			if err := renderAndDelete(targetPath, targetProject.Config); err != nil {
				return err
			}
		}
	}

	if sourceProject.IsMavenProject() && targetProject.IsMavenProject() {
		return maven.MergeAndWritePomFiles(sourceProject, targetProject)
	}
	return nil
}

func filesToCopy(sourceDir string) (files []string, err error) {
	ignores := getIgnores(sourceDir)
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, nayName := range ignores {
			if strings.Contains(path, nayName) {
				log.Debugf("ignoring %s", info.Name())
				return nil
			}
		}
		files = append(files, path)
		return nil
	})

	return
}

func getIgnores(sourceDir string) (ignores []string) {
	gitIgnores, err := file.OpenIgnoreFile(file.Path("%s/.gitignore", sourceDir))
	if err != nil {
		log.Error(err)
	}
	ignores = append(ignores, gitIgnores...)

	coPilotIgnores, err := file.OpenIgnoreFile(file.Path("%s/.co-pilot.ignore", sourceDir))
	if err != nil {
		log.Error(err)
	}
	ignores = append(ignores, coPilotIgnores...)

	ignores = append(ignores, defaultIgnores...)

	return
}

func replacePathForSource(sourceRelPath string, sourceConfig config.ProjectConfiguration, targetConfig config.ProjectConfiguration) string {
	var output = sourceRelPath

	if strings.Contains(output, ".kt") || strings.Contains(output, ".java") {
		if strings.Contains(output, "src/main") {
			output = strings.Replace(sourceRelPath, sourceConfig.SourceMainPath(), targetConfig.SourceMainPath(), 1)
		} else if strings.Contains(output, "src/test") {
			output = strings.Replace(sourceRelPath, sourceConfig.SourceTestPath(), targetConfig.SourceTestPath(), 1)
		}

		if output == sourceRelPath {
			log.Warnf("was not able to replace path for source file (.kt, .java), input and output path is the same %s", output)
		}
	}

	return output
}

func renderAndDelete(targetPath string, targetConfig interface{}) error {
	newTarget := strings.Replace(targetPath, ".render", "", 1)
	log.Infof("rendering %s into %s", targetPath, newTarget)
	if err := file.Render(targetPath, newTarget, targetConfig); err != nil {
		return err
	}

	log.Infof("deleting old render file %s", targetPath)
	return file.DeleteSingleFile(targetPath)
}
