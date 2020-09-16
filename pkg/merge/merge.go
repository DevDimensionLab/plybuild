package merge

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/maven"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var log = logger.Context()

func TemplateName(templateName string, targetDir string) error {
	cloudConfigDir, err := config.GlobalConfigDir()
	if err != nil {
		return err
	}

	templatePath := fmt.Sprintf("%s/templates/%s", cloudConfigDir, templateName)
	if !file.Exists(templatePath) {
		return errors.New(fmt.Sprintf("no such template-directory: %s", templateName))
	}

	msg := logger.Info(fmt.Sprintf("merging template %s into %s", templateName, targetDir))
	log.Info(msg)
	if err = Template(templatePath, targetDir); err != nil {
		return err
	}

	return nil
}

func Template(sourceDir string, targetDir string) error {
	var files []string

	ignores := GetIgnores(sourceDir)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
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
	if err != nil {
		return err
	}

	sourceConfig, err := config.FromProject(sourceDir)
	if err != nil {
		return err
	}

	targetConfig, err := config.FromProject(targetDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		sourceRelPath, err := file.RelPath(sourceDir, f)
		if err != nil {
			return err
		}

		sourceRelPath = replacePathForSource(sourceRelPath, sourceConfig, targetConfig)

		targetPath := fmt.Sprintf("%s/%s", targetDir, sourceRelPath)
		if err = file.CopyOrMerge(f, targetPath); err != nil {
			return err
		}

		if err = file.SearchReplace(targetPath, sourceConfig.Package, targetConfig.Package); err != nil {
			return err
		}

		if strings.HasSuffix(targetPath, ".render") {
			if err := renderAndDelete(targetPath, targetConfig); err != nil {
				return err
			}
		}
	}

	return maven.MergeAndWritePomFiles(sourceDir, targetDir)
}

func GetIgnores(sourceDir string) (ignores []string) {

	gitIgnores, err := file.OpenIgnoreFile(fmt.Sprintf("%s/.gitignore", sourceDir))
	if err != nil {
		log.Error(err)
	}
	ignores = append(ignores, gitIgnores...)

	coPilotIgnores, err := file.OpenIgnoreFile(fmt.Sprintf("%s/.co-pilot.ignore", sourceDir))
	if err != nil {
		log.Error(err)
	}
	ignores = append(ignores, coPilotIgnores...)

	otherFilesToIgnore := []string{"pom.xml", "co-pilot.json", "Application", ".co-pilot.ignore", ".gitignore"}
	ignores = append(ignores, otherFilesToIgnore...)

	return
}

func replacePathForSource(sourceRelPath string, sourceConfig config.ProjectConfiguration, targetConfig config.ProjectConfiguration) string {
	var output = sourceRelPath

	if strings.Contains(output, ".kt") || strings.Contains(output, ".java") {
		output = strings.Replace(sourceRelPath, sourceConfig.ProjectMainRoot(), targetConfig.ProjectMainRoot(), 1)
		output = strings.Replace(sourceRelPath, sourceConfig.ProjectTestRoot(), targetConfig.ProjectTestRoot(), 1)
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
	return file.Delete(targetPath)
}
