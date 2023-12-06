package file

import (
	"strings"
)

func OpenIgnoreFile(sourceDir string) (lines []string, err error) {
	ignoreFile := Path("%s/.ply.ignore", sourceDir)
	if !Exists(ignoreFile) {
		// support legacy ignore file for now
		ignoreFile = Path("%s/.co-pilot.ignore", sourceDir)
	}
	ignoreLines, err := OpenLines(ignoreFile)

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
