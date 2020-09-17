package cmd

import (
	"co-pilot/pkg/maven"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format functionality for a project",
	Long:  `Format functionality for a project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		ctx.FindAndPopulatePomModels()
	},
}

var formatPomCmd = &cobra.Command{
	Use:   "pom",
	Short: "Formats pom.xml and sorts dependencies",
	Long:  `Formats pom.xml and sorts dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachPomProject("formatting", func(model *pom.Model, args ...interface{}) error {
			return nil
		})
	},
}

var formatVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Removes version tags and replaces them with property tags",
	Long:  `Removes version tags and replaces them with property tags`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachPomProject("removes version tags", func(model *pom.Model, args ...interface{}) error {
			return maven.ChangeVersionToPropertyTags(model)
		})
	},
}

func init() {
	RootCmd.AddCommand(formatCmd)
	formatCmd.AddCommand(formatPomCmd)
	formatCmd.AddCommand(formatVersionCmd)

	formatCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	formatCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	formatCmd.PersistentFlags().BoolVar(&ctx.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	formatCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "dry run does not write to pom.xml")
}
