package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var structurizrCmd = &cobra.Command{
	Use:   "structurizr",
	Short: "Adding PNG-output support for structurizr with the help of graphviz",
	Long: `Adding PNG-output support for structurizr with the help of graphviz.

Support for structurizr requires binaries from structurizr-cli and graphviz installed:
- structurizr-cli -> https://structurizr.com/help/cli
- dot -> https://graphviz.org
`,
	Run: func(cmd *cobra.Command, args []string) {

		workspace, err := StructurizrGetMandatoryString(cmd, "workspace")
		StructurizrCheckIfError(err)

		tempDirectory := "target/"
		file.DeleteAll(tempDirectory)
		Run(exec.Command("structurizr-cli", "export", "-w", workspace, "-format", "dot", "-output", tempDirectory))

		files, err := file.FindAll("dot", []string{}, tempDirectory)
		StructurizrCheckIfError(err)

		for _, file := range files {
			outputPngFile := strings.Replace(strings.Replace(file, tempDirectory, "", 1), ".dot", "", 1) + ".png"
			println("Creating -> " + outputPngFile)
			err = RunWithOutputToFile(exec.Command("dot", file, "-Tpng"), outputPngFile)
			StructurizrCheckIfError(err)

			Run(exec.Command("open", outputPngFile))
		}
	},
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
	err = os.WriteFile(outputFile, out.Bytes(), 0644)
	return nil
}

func StructurizrCheckIfError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("\nerror: %s", err))
	os.Exit(1)
}

func StructurizrGetMandatoryString(cmd *cobra.Command, flag string) (string, error) {
	val := cmd.Flag(flag).Value.String()
	if "" == val {
		return "", errors.New(fmt.Sprintf("missing argument --%s", flag))
	}
	return val, nil
}

func init() {
	RootCmd.AddCommand(structurizrCmd)
	structurizrCmd.Flags().StringP("workspace", "w", "", "Path or URL to the workspace JSON file/DSL file(s)")
}
