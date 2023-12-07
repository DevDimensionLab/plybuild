package cmd

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/kibana"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/devdimensionlab/plybuild/pkg/structurizr"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var diagramsCmd = &cobra.Command{
	Use:   "diagrams",
	Short: "Various tools for generating diagrams",
	Long:  `Various tools for generating diagrams`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
}

var mavenGraphCmd = &cobra.Command{
	Use:   "maven-graph",
	Short: "creates a graph using maven for dependencies in a project",
	Long:  `creates a graph using maven for dependencies in a project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, inc := range mavenGraphIncludeFilters {
			println(inc)
		}
		for _, ex := range mavenGraphExcludeFilters {
			println(ex)
		}
		ctx.DryRun = true
		ctx.OnEachMavenProject("creating graph for",
			maven.Graph(false, mavenGraphExcludeTestScope, mavenGraphIncludeFilters, mavenGraphExcludeFilters),
			maven.RunOn("dot",
				"-Tpng:cairo", "target/dependency-graph.dot", "-o", "target/dependency-graph.png"),
			openReportInBrowser("target/dependency-graph.png"),
		)
	},
}

var kibanaCmd = &cobra.Command{
	Use:   "kibana",
	Short: "Specialized (experimental) command for executing a kibana-query based on a fetch-request [arg: fetch-file] and exporting the result to a json-file [arg: output-file]",
	Long:  `Specify the query in Kibana, then use developer tools to copy request as fetch (https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)`,
	Run: func(cmd *cobra.Command, args []string) {
		fetchFile, err := getMandatoryString(cmd, "fetch-file")
		checkIfError(err)

		extractFieldsInput, err := getMandatoryString(cmd, "extract-fields")
		checkIfError(err)

		outputFile := cmd.Flag("output-file").Value.String()

		fieldReMapInput := cmd.Flag("field-remap").Value.String()
		fieldFilterAndReMapping := kibana.CreateFilter(extractFieldsInput, fieldReMapInput)

		kibanaRequest, err := kibana.LoadFromFetchRequest(fetchFile)
		checkIfError(err)

		timeInterval, err := kibana.ExtractTimeIntervalFrom(kibanaRequest)
		checkIfError(err)

		resultExists := make(map[string]bool)
		err, _, result, _ := kibana.ExecuteKibanaQuery(kibanaRequest, timeInterval, fieldFilterAndReMapping, resultExists, "")
		checkIfError(err)

		if outputFile == "" {
			for _, hit := range result {
				println(hit)
			}
		} else {
			content := strings.Join(kibana.RemoveDuplicateStr(result), "\n")
			file.CreateFile(outputFile, content)
			println("Output written to [" + outputFile + "]")
		}
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
	Run: func(cmd *cobra.Command, args []string) {

		workspace, err := getMandatoryString(cmd, "workspace")
		checkIfError(err)

		tempDirectory := ".structurizr/"
		file.DeleteAll(tempDirectory)
		structurizr.Run(exec.Command("structurizr-cli", "export", "-w", workspace, "-format", "dot", "-output", tempDirectory))

		files, err := file.FindAll("dot", []string{}, tempDirectory)
		checkIfError(err)

		for _, file := range files {
			outputPngFile := strings.Replace(strings.Replace(file, tempDirectory, "", 1), ".dot", "", 1) + ".png"
			println("Creating -> " + outputPngFile)
			err = structurizr.RunWithOutputToFile(exec.Command("dot", file, "-Tpng"), outputPngFile)
			checkIfError(err)

			structurizr.Run(exec.Command("open", outputPngFile))
		}
	},
}

func init() {
	RootCmd.AddCommand(diagramsCmd)

	diagramsCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	diagramsCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")

	diagramsCmd.AddCommand(mavenGraphCmd)
	mavenGraphCmd.AddCommand(mavenGraph2PartyCmd)
	mavenGraphCmd.PersistentFlags().BoolVar(&mavenGraphExcludeTestScope, "exclude-test-scope", false, "exclude test scope from graph")
	mavenGraphCmd.PersistentFlags().StringArrayVar(&mavenGraphExcludeFilters, "exclude-filters", []string{}, "exclude filter rules")
	mavenGraphCmd.PersistentFlags().StringArrayVar(&mavenGraphIncludeFilters, "include-filters", []string{}, "include filter rules")
	mavenGraphCmd.PersistentFlags().BoolVar(&ctx.OpenInBrowser, "open", false, "open report in browser")

	diagramsCmd.AddCommand(kibanaCmd)
	kibanaCmd.Flags().StringP("fetch-file", "f", "", "Path to kibana.fetch-request")
	kibanaCmd.Flags().StringP("extract-fields", "e", "", "List of fields to extract pr hit")
	kibanaCmd.Flags().StringP("field-remap", "m", "", "List of fields matching list extract-fields with new names in output")
	kibanaCmd.Flags().StringP("output-file", "o", "", "Name of output-file to write results")

	diagramsCmd.AddCommand(structurizrCmd)
	structurizrCmd.Flags().StringP("workspace", "w", "", "Path or URL to the workspace JSON file/DSL file(s)")

}

func checkIfError(err error) {
	if err == nil {
		return
	}
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("\nerror: %s", err))
	os.Exit(1)
}

func getMandatoryString(cmd *cobra.Command, flag string) (string, error) {
	val := cmd.Flag(flag).Value.String()
	if "" == val {
		return "", errors.New(fmt.Sprintf("missing argument --%s", flag))
	}
	return val, nil
}
