package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/maven"
	"github.com/spf13/cobra"
)

type AnalyzeOpts struct {
	Deps bool
}

func (analyzeOpts AnalyzeOpts) Any() bool {
	return analyzeOpts.Deps
}

var analyzeOpts AnalyzeOpts

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
		ctx.FindAndPopulateMavenProjects()
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachProject("Undeclared and unused dependencies", func(project config.Project, args ...interface{}) error {
			return maven.ListUnusedAndUndeclared(project.Type.FilePath())
		})
	},
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().BoolVar(&analyzeOpts.Deps, "deps", false, "Show dependency usage")

	analyzeCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "recursive mode")
	analyzeCmd.PersistentFlags().StringVarP(&ctx.TargetDirectory, "target", "t", ".", "optional target directory")
}
