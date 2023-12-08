package cmd

import (
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Linting commands",
	Long:  `Linting commands`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return OkHelp(cmd, cleanOpts.Any)
	},
}

var lintKotlinCmd = &cobra.Command{
	Use:   "kotlin",
	Short: "uses ktlint for linting kotlin code",
	Long:  `uses ktlint for linting kotlin code`,
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
		ctx.OnEachMavenProject("running ktlint",
			maven.RunOn("ktlint", "\"src/**/*.kt\""),
		)
	},
}

func init() {
	RootCmd.AddCommand(lintCmd)

	lintCmd.AddCommand(lintKotlinCmd)
	lintCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	lintCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")
}
