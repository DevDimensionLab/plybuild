package git

import (
	"bytes"
	"co-pilot/pkg/file"
	"fmt"
	"github.com/prometheus/common/log"
	"os/exec"
	"strings"
)

func Clone(url string, target string) (string, error) {
	cmd := exec.Command("git", "clone", url, target)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return stdErr.String(), err
	}

	return stdOut.String(), nil
}

func Pull(target string) (string, error) {
	cmd := exec.Command("git", "-C", target, "pull", "origin")
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return stdErr.String(), err
	}

	return stdOut.String(), nil
}

func OpenIgnore(targetDir string) (lines []string) {
	ignoreLines, err := file.OpenLines(fmt.Sprintf("%s/.gitignore", targetDir))
	if err != nil {
		log.Warnf("%v")
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

	return lines
}
