package merge

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var log = logger.Context()

func TextFiles(fromFile string, toFile string) error {
	fromBytes, err := file.Open(fromFile)
	if err != nil {
		return err
	}
	fromLines := strings.Split(string(fromBytes), "\n")

	toBytes, err := file.Open(toFile)
	if err != nil {
		return err
	}
	toLines := strings.Split(string(toBytes), "\n")

	var newLines []string
	for _, fromLine := range fromLines {
		var hasLine = false
		for _, toLine := range toLines {
			if fromLine == toLine {
				hasLine = true
			}
		}
		if !hasLine {
			newLines = append(newLines, fromLine)
			log.Infof("appending line: %s", fromLine)
		}
	}

	toLines = append(toLines, newLines...)

	return file.Overwrite(toLines, toFile)
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
		if err = file.Copy(f, targetPath); err != nil {
			return err
		}

		if err = file.SearchReplace(targetPath, sourceConfig.Package, targetConfig.Package); err != nil {
			return err
		}
	}

	return nil
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
