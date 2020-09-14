package merge

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/git"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/maven"
	"co-pilot/pkg/upgrade"
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"os"
	"path/filepath"
	"strings"
)

var log = logger.Context()

func TemplateName(templateName string, targetDirectory string) error {
	cloudConfigDir, err := config.GlobalConfigDir()
	if err != nil {
		return err
	}

	templatePath := fmt.Sprintf("%s/templates/%s", cloudConfigDir, templateName)
	if !file.Exists(templatePath) {
		return errors.New(fmt.Sprintf("no such template-directory: %s", templateName))
	}

	msg := logger.Info(fmt.Sprintf("merging template %s into %s", templateName, targetDirectory))
	log.Info(msg)
	if err = Template(templatePath, targetDirectory); err != nil {
		return err
	}

	return nil
}

func Template(source string, target string) error {
	var files []string

	gitIgnores, err := git.OpenIgnore(source)
	if err != nil {
		log.Errorln(err)
	}
	gitIgnores = append(gitIgnores, "Application")

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, nayName := range gitIgnores {
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

	sourceConfig, err := config.FromProject(source)
	if err != nil {
		return err
	}

	targetConfig, err := config.FromProject(target)
	if err != nil {
		return err
	}

	for _, f := range files {
		sourceRelPath, err := file.RelPath(source, f)
		if err != nil {
			return err
		}

		sourceRelPath = replacePathForSource(sourceRelPath, sourceConfig, targetConfig)

		targetPath := fmt.Sprintf("%s/%s", target, sourceRelPath)
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

	return mergeAndWritePomFiles(source, target)
}

func replacePathForSource(sourceRelPath string, sourceConfig config.ProjectConfiguration, targetConfig config.ProjectConfiguration) string {
	var output = sourceRelPath

	if strings.Contains(output, ".kt") || strings.Contains(output, ".java") {
		output = strings.Replace(sourceRelPath, sourceConfig.ProjectMainRoot(), targetConfig.ProjectMainRoot(), 1)
		output = strings.Replace(sourceRelPath, sourceConfig.ProjectTestRoot(), targetConfig.ProjectTestRoot(), 1)
	}

	return output
}

func mergeAndWritePomFiles(source string, target string) error {
	fromPomFile := source + "/pom.xml"
	importModel, err := pom.GetModelFrom(fromPomFile)
	if err != nil {
		log.Warnln(err)
		return nil
	}

	toPomFile := target + "/pom.xml"
	projectModel, err := pom.GetModelFrom(toPomFile)
	if err != nil {
		log.Warnln(err)
		return nil
	}

	if err = maven.Merge(importModel, projectModel); err != nil {
		return err
	}

	return upgrade.SortAndWrite(projectModel, toPomFile)
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
