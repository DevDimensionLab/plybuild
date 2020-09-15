package file

import (
	"strings"
)

func OpenIgnoreFile(filePath string) (lines []string, err error) {
	ignoreLines, err := OpenLines(filePath)
	if err != nil {
		return lines, err
	}

	for _, line := range ignoreLines {
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "!") && line != "" {
			if strings.HasPrefix(line, "*") {
				lines = append(lines, strings.Replace(line, "*", "", 1))
			} else {
				lines = append(lines, line)
			}
		}
	}

	return lines, nil
}
