package cmd

import (
	"github.com/devdimensionlab/co-pilot/pkg/maven"
	"github.com/spf13/cobra"
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
}
