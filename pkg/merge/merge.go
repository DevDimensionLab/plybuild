package merge

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/maven"
	"co-pilot/pkg/upgrade"
	"errors"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"os"
	"path/filepath"
	"strings"
)

func TemplateName(templateName string, targetDirectory string) error {
	cloudConfigDir, err := config.GlobalConfigDir()
	if err != nil {
		return err
	}

	templatePath := fmt.Sprintf("%s/templates/%s", cloudConfigDir, templateName)
	if !file.Exists(templatePath) {
		return errors.New(fmt.Sprintf("no such template-directory: %s", templateName))
	}

	if err = Template(templatePath, targetDirectory); err != nil {
		return err
	}

	return nil
}

func Template(source string, target string) error {
	var files []string

	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, nayName := range DoNotCopyFileName() {
			if strings.Contains(info.Name(), nayName) {
				return nil
			}
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}

	_, err = config.FromProject(source)
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
	}

	return mergeAndWritePomFiles(source, target)
}

func DoNotCopyFileName() []string {
	return []string{"Application"}
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
		return err
	}

	toPomFile := target + "/pom.xml"
	projectModel, err := pom.GetModelFrom(toPomFile)
	if err != nil {
		return err
	}

	if err = maven.Merge(importModel, projectModel); err != nil {
		return err
	}

	return upgrade.SortAndWrite(projectModel, toPomFile)
}
