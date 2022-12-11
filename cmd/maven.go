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

var mavenGraph2PartyCmd = &cobra.Command{
	Use:   "2party",
	Short: "creates a graph only for 2party dependencies in a project",
	Long:  `creates a graph only for 2party dependencies in a project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachMavenProject("creating 2party graph for",
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
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachMavenProject("running checkstyle analysis on",
			maven.RunOn("mvn", "checkstyle:checkstyle"),
		)
	},
}

var mavenOwaspCmd = &cobra.Command{
	Use:   "owasp",
	Short: "runs owasp",
	Long:  `runs owasp`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachMavenProject("running wasp analysis on",
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
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachMavenProject("running spring-boot:run",
			maven.RunOn("mvn", "spring-boot:run"),
		)
	},
}

var mavenEnforcerCmd = &cobra.Command{
	Use:   "enforcer",
	Short: "runs enforcer",
	Long:  `runs enforcer`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachMavenProject("running enforcer on",
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
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachMavenProject("Undeclared and unused dependencies", func(repository maven.Repository, project config.Project) error {
			return maven.ListUnusedAndUndeclared(project.Type.FilePath())
		})
	},
}

func openReportInBrowser(reportPath string) func(repository maven.Repository, project config.Project) error {
	return func(repository maven.Repository, project config.Project) error {
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

	mavenCmd.AddCommand(mavenCheckstyleCmd)
	mavenCheckstyleCmd.PersistentFlags().BoolVar(&ctx.OpenInBrowser, "open", false, "open report in browser")

	mavenCmd.AddCommand(mavenOwaspCmd)
	mavenOwaspCmd.PersistentFlags().BoolVar(&ctx.OpenInBrowser, "open", false, "open report in browser")
	mavenCmd.AddCommand(mavenEnforcerCmd)
	mavenCmd.AddCommand(mavenSpringBootRunCmd)

	mavenCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().BoolVar(&analyzeOpts.Deps, "deps", false, "Show dependency usage")
}
