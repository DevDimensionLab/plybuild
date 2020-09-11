package file

import (
	"os"
	"text/template"
)

func Render(inputFilePath string, outputFilePath string, r interface{}) error {
	inputBytes, err := Open(inputFilePath)
	if err != nil {
		return err
	}

	t := template.Must(template.New(inputFilePath).Parse(string(inputBytes)))

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	err = t.Execute(outputFile, r)
	if err != nil {
		return err
	}

	return nil
}
