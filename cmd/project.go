package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/service"
	"fmt"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project [OPTIONS]",
	Short: "Project options",
	Long:  `Various project helper commands`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		ctx.FindAndPopulatePomModels()
	},
}

var projectInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a maven project with co-pilot files and formatting",
	Long:  `Initializes a maven project with co-pilot files and formatting`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pair := range ctx.PomPairs {
			log.Info(logger.White(fmt.Sprintf("formating pom file %s", pair.PomFile)))

			if !ctx.DryRun {
				projectConfigFile := fmt.Sprintf("%sco-pilot.json", service.PomFileToTargetDirectory(pair.PomFile))
				projectCfg := config.InitProjectConfigurationFromModel(pair.Model)

				log.Infof("writes co-pilot.json config file to %s", projectConfigFile)
				if err := projectCfg.Write(projectConfigFile); err != nil {
					log.Warnln(err)
					continue
				}

				if err := service.Write(ctx.Overwrite, pair); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectInitCmd)

	projectCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	projectCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	projectCmd.PersistentFlags().BoolVar(&ctx.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	projectCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "dry run does not write to pom.xml")
}
