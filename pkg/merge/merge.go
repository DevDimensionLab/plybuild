package merge

import (
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

func Directory(source string, target string) error {
	var files []string

	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		for _, nay := range DoNotCopyFilesContains() {
			if strings.Contains(info.Name(), nay) {
				return nil
			}
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}

	for _, f := range files {
		relPath, err := file.RelPath(source, f)
		if err != nil {
			return err
		}

		targetPath := fmt.Sprintf("%s/%s", target, relPath)
		err = file.Copy(f, targetPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func DoNotCopyFilesContains() []string {
	return []string{"Application"}
}
