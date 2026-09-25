package structurizr

import (
	"bytes"
	"os/exec"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
)

type structurizrOutputWriteDependencies struct {
	Files filesystem.Dependencies
}

func systemStructurizrOutputWriteDependencies() structurizrOutputWriteDependencies {
	return structurizrOutputWriteDependencies{Files: filesystem.System()}
}

func writeStructurizrOutput(dependencies structurizrOutputWriteDependencies, outputFile string, data []byte) error {
	return filesystem.WriteFile(dependencies.Files, outputFile, data, 0644)
}

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
	_ = writeStructurizrOutput(systemStructurizrOutputWriteDependencies(), outputFile, out.Bytes())
	return nil
}
