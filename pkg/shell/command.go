package shell

import (
	"bytes"
	"os/exec"
)

func Run(name string, args ...string) (string, error) {
	return run(exec.Command(name, args...))
}

func run(cmd *exec.Cmd) (string, error) {
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return stdErr.String(), err
	}

	return stdOut.String(), nil
}

func Wget(url, filepath string) (string, error) {
	return run(exec.Command("wget", url, "-O", filepath))
}

func Unzip(file string, outputDir string) (string, error) {
	return run(exec.Command("unzip", file, "-d", outputDir))
}
