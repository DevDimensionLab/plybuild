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

var log = logger.Context()

func MergeTemplate(cloudConfig config.CloudConfig, templateName string, targetDir string) error {
	template, err := cloudConfig.Template(templateName)
	if err != nil {
		return err
	}

	msg := logger.Info(fmt.Sprintf("merging template %s into %s", templateName, targetDir))
	log.Info(msg)
	if err := merge(template, targetDir); err != nil {
		return err
	}

	return nil
}

func merge(template config.CloudTemplate, targetDir string) error {
	sourceDir := template.Impl.Path
	files, err := FilesToCopy(sourceDir)
	if err != nil {
		return err
	}

	sourceConfig, err := config.InitProjectConfigurationFromDir(sourceDir)
	if err != nil {
		return err
	}

	targetConfig, err := config.InitProjectConfigurationFromDir(targetDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		sourceRelPath, err := file.RelPath(sourceDir, f)
		if err != nil {
			return err
		}

		sourceRelPath = ReplacePathForSource(sourceRelPath, sourceConfig, targetConfig)

		targetPath := file.Path("%s/%s", targetDir, sourceRelPath)
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

func FilesToCopy(sourceDir string) (files []string, err error) {
	ignores := GetIgnores(sourceDir)
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

func GetIgnores(sourceDir string) (ignores []string) {

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

	otherFilesToIgnore := []string{"pom.xml", "co-pilot.json", "Application", ".co-pilot.ignore", ".gitignore", ".mvn", "mvnw", "mvnw.cmd"}
	ignores = append(ignores, otherFilesToIgnore...)

	return
}

func Apply(cloudConfig config.CloudConfig, templates map[string]bool, targetDirectory string) {
	for k, _ := range templates {
		log.Infof("applying template %s", k)
		if err := MergeTemplate(cloudConfig, k, targetDirectory); err != nil {
			log.Warnf("%v", err)
		}
	}
}

func ReplacePathForSource(sourceRelPath string, sourceConfig config.ProjectConfiguration, targetConfig config.ProjectConfiguration) string {
	var output = sourceRelPath

	if strings.Contains(output, ".kt") || strings.Contains(output, ".java") {
		if strings.Contains(output, "src/main") {
			output = strings.Replace(sourceRelPath, sourceConfig.SourceMainPath(), targetConfig.SourceMainPath(), 1)
		} else if strings.Contains(output, "src/test") {
			output = strings.Replace(sourceRelPath, sourceConfig.SourceTestPath(), targetConfig.SourceTestPath(), 1)
		}
	}

	if output == sourceRelPath {
		log.Warnf("was not able to replace path for source, source and output path is the same %s", output)
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
