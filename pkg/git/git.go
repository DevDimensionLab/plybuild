package git

import (
	"bytes"
	"os/exec"
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
