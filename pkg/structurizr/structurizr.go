package structurizr

import (
	"bytes"
	"os"
	"os/exec"
)

func Run(command *exec.Cmd) error {
	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}

func RunWithOutputToFile(command *exec.Cmd, outputFile string) error {
	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		return err
	}
	err = os.WriteFile(outputFile, out.Bytes(), 0644)
	return nil
}
