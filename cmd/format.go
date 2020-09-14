package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/logger"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format functionality for a project",
	Long:  `Format functionality for a project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd, args); err != nil {
			log.Fatalln(err)
		}
		populatePomFiles()
	},
}

var formatPomCmd = &cobra.Command{
	Use:   "pom",
	Short: "Formats pom.xml and sorts dependencies",
	Long:  `Formats pom.xml and sorts dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("formating pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err := write(pomFile, model); err != nil {
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
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("removes version tags for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}

			if err = clean.VersionToPropertyTags(model); err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err := write(pomFile, model); err != nil {
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

	formatCmd.PersistentFlags().BoolVarP(&cArgs.Recursive, "recursive", "r", false, "turn on recursive mode")
	formatCmd.PersistentFlags().StringVar(&cArgs.TargetDirectory, "target", ".", "Optional target directory")
	formatCmd.PersistentFlags().BoolVar(&cArgs.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	formatCmd.PersistentFlags().BoolVar(&cArgs.DryRun, "dry-run", false, "dry run does not write to pom.xml")
}
