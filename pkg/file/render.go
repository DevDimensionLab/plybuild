package file

import (
	"text/template"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
)

type renderDependencies struct {
	Files filesystem.Dependencies
}

func systemRenderDependencies() renderDependencies {
	return renderDependencies{Files: filesystem.System()}
}

func Render(inputFilePath string, outputFilePath string, r interface{}) error {
	return render(systemRenderDependencies(), inputFilePath, outputFilePath, r)
}

func render(dependencies renderDependencies, inputFilePath string, outputFilePath string, r interface{}) error {
	inputBytes, err := filesystem.ReadFile(dependencies.Files, inputFilePath)
	if err != nil {
		return err
	}

	t := template.Must(template.New(inputFilePath).Parse(string(inputBytes)))

	outputFile, err := filesystem.Create(dependencies.Files, outputFilePath)
	if err != nil {
		return err
	}

	err = t.Execute(outputFile, r)
	if err != nil {
		return err
	}

	return nil
}
