/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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

		tempDirectory := "./target"
		Run(exec.Command("structurizr-cli", "export", "-w", workspace, "-format", "dot", "-output", tempDirectory))

		outputPng := workspace + ".png"
		RunWithOutput(exec.Command("dot", tempDirectory+"/structurizr-SystemContext.dot", "-Tpng"), outputPng)

		fmt.Println("Created " + workspace + ".png")
		Run(exec.Command("open", outputPng))
	},
}

func Run(command *exec.Cmd) error {
	err := command.Run()
	if err != nil {
		return err
	}
	return nil
}

func RunWithOutput(command *exec.Cmd, outputFile string) error {
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

func init() {
	RootCmd.AddCommand(structurizrCmd)
	structurizrCmd.Flags().StringP("workspace", "w", "", "Path or URL to the workspace JSON file/DSL file(s)")
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
