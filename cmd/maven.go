package cmd

import (
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/maven"
	"github.com/devdimensionlab/co-pilot/pkg/webservice"
	"github.com/spf13/cobra"
)

var mavenGraphExcludeTestScope bool
var mavenGraphExcludeFilters []string
var mavenGraphIncludeFilters []string

type AnalyzeOpts struct {
	Deps bool
}

func (analyzeOpts AnalyzeOpts) Any() bool {
	return analyzeOpts.Deps
}

var analyzeOpts AnalyzeOpts

var mavenCmd = &cobra.Command{
	Use:   "maven",
	Short: "Run maven (mvn) commands",
	Long:  `Run maven (mvn) commands`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return OkHelp(cmd, cleanOpts.Any)
	},
}

var mavenGraphCmd = &cobra.Command{
	Use:   "graph",
	Short: "creates a graph for dependencies in a project",
	Long:  `creates a graph for dependencies in a project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
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
		ctx.OnEachProject("creating graph for",
			maven.Graph(false, mavenGraphExcludeTestScope, mavenGraphIncludeFilters, mavenGraphExcludeFilters),
			maven.RunOn("dot",
				"-Tpng:cairo", "target/dependency-graph.dot", "-o", "target/dependency-graph.png"),
			openReportInBrowser("target/dependency-graph.png"),
		)
	},
}

var mavenGraph2PartyCmd = &cobra.Command{
	Use:   "2party",
	Short: "creates a graph only for 2party dependencies in a project",
	Long:  `creates a graph only for 2party dependencies in a project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachProject("creating 2party graph for",
			maven.Graph(true, mavenGraphExcludeTestScope, mavenGraphIncludeFilters, mavenGraphExcludeFilters),
			maven.RunOn("dot",
				"-Tpng:cairo", "target/dependency-graph.dot", "-o", "target/dependency-graph.png"),
			openReportInBrowser("target/dependency-graph.png"),
		)
	},
}

var mavenCheckstyleCmd = &cobra.Command{
	Use:   "checkstyle",
	Short: "runs checkstyle",
	Long:  `runs checkstyle`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachProject("running checkstyle analysis on",
			maven.RunOn("mvn", "checkstyle:checkstyle"),
		)
	},
}

var mavenOwaspCmd = &cobra.Command{
	Use:   "owasp",
	Short: "runs owasp",
	Long:  `runs owasp`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachProject("running wasp analysis on",
			maven.RunOn("mvn", "org.owasp:dependency-check-maven:check"),
			openReportInBrowser("target/dependency-check-report.html"),
		)
	},
}

var mavenSpringBootRunCmd = &cobra.Command{
	Use:   "boot-run",
	Short: "runs a spring boot application",
	Long:  `runs a spring boot application`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachProject("running spring-boot:run",
			maven.RunOn("mvn", "spring-boot:run"),
		)
	},
}

var mavenEnforcerCmd = &cobra.Command{
	Use:   "enforcer",
	Short: "runs enforcer",
	Long:  `runs enforcer`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachProject("running enforcer on",
			maven.RunOn(
				"mvn",
				"org.apache.maven.plugins:maven-enforcer-plugin:enforce",
				"-Drules=banDuplicatePomDependencyVersions,dependencyConvergence",
			),
		)
	},
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Perform an analyze on a projects",
	Long:  `Perform an analyze on a projects`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return OkHelp(cmd, analyzeOpts.Any)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachProject("Undeclared and unused dependencies", func(project config.Project) error {
			return maven.ListUnusedAndUndeclared(project.Type.FilePath())
		})
	},
}

func openReportInBrowser(reportPath string) func(project config.Project) error {
	return func(project config.Project) error {
		if ctx.OpenInBrowser {
			return webservice.OpenBrowser(fmt.Sprintf("%s/%s", project.Path, reportPath))
		}
		return nil
	}
}

func init() {
	RootCmd.AddCommand(mavenCmd)

	mavenCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	mavenCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")

	mavenCmd.AddCommand(mavenGraphCmd)
	mavenGraphCmd.AddCommand(mavenGraph2PartyCmd)
	mavenGraphCmd.PersistentFlags().BoolVar(&mavenGraphExcludeTestScope, "exclude-test-scope", false, "exclude test scope from graph")
	mavenGraphCmd.PersistentFlags().StringArrayVar(&mavenGraphExcludeFilters, "exclude-filters", []string{}, "exclude filter rules")
	mavenGraphCmd.PersistentFlags().StringArrayVar(&mavenGraphIncludeFilters, "include-filters", []string{}, "include filter rules")
	mavenGraphCmd.PersistentFlags().BoolVar(&ctx.OpenInBrowser, "open", false, "open report in browser")
	mavenCmd.AddCommand(mavenCheckstyleCmd)
	mavenCheckstyleCmd.PersistentFlags().BoolVar(&ctx.OpenInBrowser, "open", false, "open report in browser")
	mavenCmd.AddCommand(mavenOwaspCmd)
	mavenOwaspCmd.PersistentFlags().BoolVar(&ctx.OpenInBrowser, "open", false, "open report in browser")
	mavenCmd.AddCommand(mavenEnforcerCmd)
	mavenCmd.AddCommand(mavenSpringBootRunCmd)

	mavenCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().BoolVar(&analyzeOpts.Deps, "deps", false, "Show dependency usage")
}
