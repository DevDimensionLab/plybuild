package cmd

import (
	"bytes"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/spf13/cobra"
	"strings"
)

var diagramsCmd = &cobra.Command{
	Use:   "diagrams",
	Short: "Various tools for generating diagrams",
	Long:  `Various tools for generating diagrams`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitGlobals(cmd); err != nil {
			return err
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			return err
		}
		return nil
	},
}

var structurizrCmd = &cobra.Command{
	Use:   "structurizr",
	Short: "Adding PNG-output support for structurizr with the help of graphviz",
	Long: `Adding PNG-output support for structurizr with the help of graphviz.

Support for structurizr requires binaries from structurizr-cli and graphviz installed:
- structurizr-cli -> https://structurizr.com/help/cli
- dot -> https://graphviz.org
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		workspace, err := getMandatoryString(cmd, "workspace")
		if err != nil {
			return err
		}
		return runStructurizrDiagrams(
			systemPluginDiagramsExportDependencies(),
			systemPluginDiagramsGraphvizDependencies(),
			systemPluginDiagramsOpenDependencies(),
			workspace,
		)
	},
}

type pluginDiagramsExportDependencies struct {
	Process process.Dependencies
}

func systemPluginDiagramsExportDependencies() pluginDiagramsExportDependencies {
	return pluginDiagramsExportDependencies{Process: process.SystemRunner()}
}

type pluginDiagramsGraphvizDependencies struct {
	Process process.Dependencies
	Files   filesystem.Dependencies
}

func systemPluginDiagramsGraphvizDependencies() pluginDiagramsGraphvizDependencies {
	return pluginDiagramsGraphvizDependencies{
		Process: process.SystemRunner(),
		Files:   filesystem.System(),
	}
}

type pluginDiagramsOpenDependencies struct {
	Process process.Dependencies
}

func systemPluginDiagramsOpenDependencies() pluginDiagramsOpenDependencies {
	return pluginDiagramsOpenDependencies{Process: process.SystemRunner()}
}

func openStructurizrDiagram(dependencies pluginDiagramsOpenDependencies, outputPngFile string) {
	_ = process.Execute(dependencies.Process, process.Command{
		Name: "open",
		Args: []string{outputPngFile},
	})
}

func convertStructurizrDiagram(
	dependencies pluginDiagramsGraphvizDependencies,
	inputDotFile string,
	outputPngFile string,
) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := process.Execute(dependencies.Process, process.Command{
		Name:   "dot",
		Args:   []string{inputDotFile, "-Tpng"},
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return err
	}
	_ = filesystem.WriteFile(dependencies.Files, outputPngFile, stdout.Bytes(), 0644)
	return nil
}

func runStructurizrDiagrams(
	exportDependencies pluginDiagramsExportDependencies,
	graphvizDependencies pluginDiagramsGraphvizDependencies,
	openDependencies pluginDiagramsOpenDependencies,
	workspace string,
) error {
	tempDirectory := ".structurizr/"
	_ = file.DeleteAll(tempDirectory)
	_ = process.Execute(exportDependencies.Process, process.Command{
		Name: "structurizr-cli",
		Args: []string{"export", "-w", workspace, "-format", "dot", "-output", tempDirectory},
	})

	files, err := file.FindAll("dot", []string{}, tempDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		outputPngFile := strings.Replace(strings.Replace(file, tempDirectory, "", 1), ".dot", "", 1) + ".png"
		println("Creating -> " + outputPngFile)
		err = convertStructurizrDiagram(graphvizDependencies, file, outputPngFile)
		if err != nil {
			return err
		}

		openStructurizrDiagram(openDependencies, outputPngFile)
	}
	return nil
}

func init() {
	pluginCmd.AddCommand(diagramsCmd)

	diagramsCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	diagramsCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")

	diagramsCmd.AddCommand(structurizrCmd)
	structurizrCmd.Flags().StringP("workspace", "w", "", "Path or URL to the workspace JSON file/DSL file(s)")

}

func getMandatoryString(cmd *cobra.Command, flag string) (string, error) {
	val := cmd.Flag(flag).Value.String()
	if val == "" {
		return "", fmt.Errorf("missing argument --%s", flag)
	}
	return val, nil
}
