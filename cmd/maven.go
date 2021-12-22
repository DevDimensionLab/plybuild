package cmd

import (
	"github.com/co-pilot-cli/co-pilot/pkg/maven"
	"github.com/spf13/cobra"
	"os"
)

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
		ctx.DryRun = true
		ctx.OnEachProject("creating graph for",
			maven.RunOn(os.Stdout, "mvn",
				"com.github.ferstl:depgraph-maven-plugin:graph",
				"-DshowVersions",
				"-DshowGroupIds",
				"-DshowConflicts",
				"-DshowDuplicates"),
			maven.RunOn(os.Stdout, "dot",
				"-Tpng:cairo", "target/dependency-graph.dot", "-o", "target/dependency-graph.png"),
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
			maven.RunOn(os.Stdout, "mvn", "checkstyle:checkstyle"),
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
			maven.RunOn(os.Stdout, "mvn", "org.owasp:dependency-check-maven:check"),
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
				os.Stdout,
				"mvn",
				"org.apache.maven.plugins:maven-enforcer-plugin:enforce",
				"-Drules=banDuplicatePomDependencyVersions,dependencyConvergence",
			),
		)
	},
}

func init() {
	RootCmd.AddCommand(mavenCmd)

	mavenCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	mavenCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")

	mavenCmd.AddCommand(mavenGraphCmd)
	mavenCmd.AddCommand(mavenCheckstyleCmd)
	mavenCmd.AddCommand(mavenOwaspCmd)
	mavenCmd.AddCommand(mavenEnforcerCmd)
}
