package shell

import (
	"bytes"
	"os/exec"
)

func Wget(url, filepath string) (string, error) {
	cmd := exec.Command("wget", url, "-O", filepath)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return stdErr.String(), err
	}

	return stdOut.String(), nil
}

func Unzip(file string, outputDir string) (string, error) {
	cmd := exec.Command("unzip", file, "-d", outputDir)
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return stdErr.String(), err
	}

	return stdOut.String(), nil
}
