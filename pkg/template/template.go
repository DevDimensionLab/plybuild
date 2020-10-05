package template

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/maven"
	"fmt"
	"github.com/sirupsen/logrus"
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

type Template struct {
	log logrus.FieldLogger
}

func With(log logrus.FieldLogger) Template {
	return Template{
		log: log,
	}
}

func (tmpl Template) MergeTemplates(templates []config.CloudTemplate, target config.Project) {
	for _, template := range templates {
		tmpl.log.Infof("applying Template %s", template.Name)
		if err := tmpl.MergeTemplate(template, target); err != nil {
			tmpl.log.Warnf("%v", err)
		}
	}
}

func (tmpl Template) MergeTemplate(cloudTemplate config.CloudTemplate, target config.Project) error {
	if target.IsDirtyGitRepo() {
		tmpl.log.Warn(logger.White(fmt.Sprintf("merging Template %s into a dirty git repository %s", cloudTemplate.Name, target.Path)))
	} else {
		tmpl.log.Info(logger.White(fmt.Sprintf("merging Template %s into %s", cloudTemplate.Name, target.Path)))
	}
	if err := tmpl.merge(cloudTemplate.Project, target); err != nil {
		return err
	}

	return nil
}

func (tmpl Template) merge(sourceProject config.Project, targetProject config.Project) error {
	sourceDir := sourceProject.Path
	files, err := tmpl.filesToCopy(sourceDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		sourceRelPath, err := file.RelPath(sourceDir, f)
		if err != nil {
			return err
		}

		sourceRelPath = tmpl.replacePathForSource(sourceRelPath, sourceProject.Config, targetProject.Config)

		targetPath := file.Path("%s/%s", targetProject.Path, sourceRelPath)
		if err = file.CopyOrMerge(f, targetPath); err != nil {
			return err
		}

		if err = file.SearchReplace(targetPath, sourceProject.Config.Package, targetProject.Config.Package); err != nil {
			return err
		}

		if strings.HasSuffix(targetPath, ".render") {
			if err := tmpl.renderAndDelete(targetPath, targetProject.Config); err != nil {
				return err
			}
		}
	}

	if sourceProject.IsMavenProject() && targetProject.IsMavenProject() {
		return maven.MergeAndWritePomFiles(sourceProject, targetProject)
	}
	return nil
}

func (tmpl Template) filesToCopy(sourceDir string) (files []string, err error) {
	ignores := tmpl.getIgnores(sourceDir)
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, nayName := range ignores {
			if strings.Contains(path, nayName) {
				tmpl.log.Debugf("ignoring %s", info.Name())
				return nil
			}
		}
		files = append(files, path)
		return nil
	})

	return
}

func (tmpl Template) getIgnores(sourceDir string) (ignores []string) {

	gitIgnores, err := file.OpenIgnoreFile(file.Path("%s/.gitignore", sourceDir))
	if err != nil {
		tmpl.log.Error(err)
	}
	ignores = append(ignores, gitIgnores...)

	coPilotIgnores, err := file.OpenIgnoreFile(file.Path("%s/.co-pilot.ignore", sourceDir))
	if err != nil {
		tmpl.log.Error(err)
	}
	ignores = append(ignores, coPilotIgnores...)

	ignores = append(ignores, defaultIgnores...)

	return
}

func (tmpl Template) replacePathForSource(sourceRelPath string, sourceConfig config.ProjectConfiguration, targetConfig config.ProjectConfiguration) string {
	var output = sourceRelPath

	if strings.Contains(output, ".kt") || strings.Contains(output, ".java") {
		if strings.Contains(output, "src/main") {
			output = strings.Replace(sourceRelPath, sourceConfig.SourceMainPath(), targetConfig.SourceMainPath(), 1)
		} else if strings.Contains(output, "src/test") {
			output = strings.Replace(sourceRelPath, sourceConfig.SourceTestPath(), targetConfig.SourceTestPath(), 1)
		}

		if output == sourceRelPath {
			tmpl.log.Warnf("was not able to replace path for source file (.kt, .java), input and output path is the same %s", output)
		}
	}

	return output
}

func (tmpl Template) renderAndDelete(targetPath string, targetConfig interface{}) error {
	newTarget := strings.Replace(targetPath, ".render", "", 1)
	tmpl.log.Infof("rendering %s into %s", targetPath, newTarget)
	if err := file.Render(targetPath, newTarget, targetConfig); err != nil {
		return err
	}

	tmpl.log.Infof("deleting old render file %s", targetPath)
	return file.DeleteSingleFile(targetPath)
}
