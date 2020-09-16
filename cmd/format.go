package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/service"
	"fmt"
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
		for pomFile, model := range ctx.PomModels {
			log.Info(logger.White(fmt.Sprintf("formating pom file %s", pomFile)))

			if !ctx.DryRun {
				if err := service.Write(ctx.Overwrite, pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

var formatVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Removes version tags and replaces them with property tags",
	Long:  `Removes version tags and replaces them with property tags`,
	Run: func(cmd *cobra.Command, args []string) {
		for pomFile, model := range ctx.PomModels {
			log.Info(logger.White(fmt.Sprintf("removes version tags for pom file %s", pomFile)))

			if err := clean.VersionToPropertyTags(model); err != nil {
				log.Warnln(err)
				continue
			}

			if !ctx.DryRun {
				if err := service.Write(ctx.Overwrite, pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
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
