package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/kibana"
	"github.com/devdimensionlab/plybuild/pkg/structurizr"
	"github.com/spf13/cobra"
	"os/exec"
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

var kibanaCmd = &cobra.Command{
	Use:   "kibana",
	Short: "Specialized (experimental) command for executing a kibana-query based on a fetch-request [arg: fetch-file] and exporting the result to a json-file [arg: output-file]",
	Long:  `Specify the query in Kibana, then use developer tools to copy request as fetch (https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fetchFile, err := getMandatoryString(cmd, "fetch-file")
		if err != nil {
			return err
		}

		extractFieldsInput, err := getMandatoryString(cmd, "extract-fields")
		if err != nil {
			return err
		}

		outputFile := cmd.Flag("output-file").Value.String()

		fieldReMapInput := cmd.Flag("field-remap").Value.String()
		fieldFilterAndReMapping := kibana.CreateFilter(extractFieldsInput, fieldReMapInput)

		kibanaRequest, err := kibana.LoadFromFetchRequest(fetchFile)
		if err != nil {
			return err
		}

		timeInterval, err := kibana.ExtractTimeIntervalFrom(kibanaRequest)
		if err != nil {
			return err
		}

		resultExists := make(map[string]bool)
		err, _, result, _ := kibana.ExecuteKibanaQuery(kibanaRequest, timeInterval, fieldFilterAndReMapping, resultExists, "")
		if err != nil {
			return err
		}

		if outputFile == "" {
			for _, hit := range result {
				println(hit)
			}
		} else {
			content := strings.Join(kibana.RemoveDuplicateStr(result), "\n")
			_ = file.CreateFile(outputFile, content)
			println("Output written to [" + outputFile + "]")
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
		return runStructurizrDiagrams(systemPluginDiagramsExportDependencies(), workspace)
	},
}

type pluginDiagramsExportDependencies struct {
	Process process.Dependencies
}

func systemPluginDiagramsExportDependencies() pluginDiagramsExportDependencies {
	return pluginDiagramsExportDependencies{Process: process.System()}
}

func runStructurizrDiagrams(dependencies pluginDiagramsExportDependencies, workspace string) error {
	tempDirectory := ".structurizr/"
	_ = file.DeleteAll(tempDirectory)
	_ = process.Execute(dependencies.Process, process.Command{
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
		err = structurizr.RunWithOutputToFile(exec.Command("dot", file, "-Tpng"), outputPngFile)
		if err != nil {
			return err
		}

		_ = structurizr.Run(exec.Command("open", outputPngFile))
	}
	return nil
}

func init() {
	pluginCmd.AddCommand(diagramsCmd)

	diagramsCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	diagramsCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")

	diagramsCmd.AddCommand(kibanaCmd)
	kibanaCmd.Flags().StringP("fetch-file", "f", "", "Path to kibana.fetch-request")
	kibanaCmd.Flags().StringP("extract-fields", "e", "", "List of fields to extract pr hit")
	kibanaCmd.Flags().StringP("field-remap", "m", "", "List of fields matching list extract-fields with new names in output")
	kibanaCmd.Flags().StringP("output-file", "o", "", "Name of output-file to write results")

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
