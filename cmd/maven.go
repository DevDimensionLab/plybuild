package cmd

import (
	"github.com/co-pilot-cli/co-pilot/pkg/maven"
	"github.com/spf13/cobra"
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
			maven.RunOn("mvn",
				"com.github.ferstl:depgraph-maven-plugin:graph",
				"-DshowVersions",
				"-DshowGroupIds",
				"-DshowConflicts",
				"-DshowDuplicates"),
			maven.RunOn("dot",
				"-Tpng:cairo", "target/dependency-graph.dot", "-o", "target/dependency-graph.png"),
		)
	},
}

func init() {
	RootCmd.AddCommand(mavenCmd)

	mavenCmd.AddCommand(mavenGraphCmd)
	mavenCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	mavenCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")
}
