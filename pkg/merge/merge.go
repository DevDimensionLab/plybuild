package merge

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
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
