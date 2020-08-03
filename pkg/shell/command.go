package shell

import (
	"os"
	"os/exec"
)

func Wget(url, filepath string) error {
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Unzip(file string, outputDir string) error {
	cmd := exec.Command("unzip", file, "-d", outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
