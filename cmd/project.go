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
		for pomFile, model := range ctx.PomModels {
			log.Info(logger.White(fmt.Sprintf("formating pom file %s", pomFile)))

			if !ctx.DryRun {
				configFile := fmt.Sprintf("%sco-pilot.json", service.PomFileToTargetDirectory(pomFile))
				initConfig, err := config.GenerateConfig(model)
				if err != nil {
					log.Warnln(err)
					continue
				}

				log.Infof("writes co-pilot.json config file to %s", configFile)
				if err = initConfig.WriteConfig(configFile); err != nil {
					log.Warnln(err)
					continue
				}

				if err := service.Write(ctx.Overwrite, pomFile, model); err != nil {
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
