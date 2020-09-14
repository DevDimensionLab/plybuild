package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/logger"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project [OPTIONS]",
	Short: "Project options",
	Long:  `Various project helper commands`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd, args); err != nil {
			log.Fatalln(err)
		}
		populatePomFiles()
	},
}

var projectInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a maven project with co-pilot files and formatting",
	Long:  `Initializes a maven project with co-pilot files and formatting`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("formating pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				configFile := fmt.Sprintf("%sco-pilot.json", pomFileToTargetDirectory(pomFile))
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

				if err := write(pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectInitCmd)

	projectCmd.PersistentFlags().BoolVarP(&cArgs.Recursive, "recursive", "r", false, "turn on recursive mode")
	projectCmd.PersistentFlags().StringVar(&cArgs.TargetDirectory, "target", ".", "Optional target directory")
	projectCmd.PersistentFlags().BoolVar(&cArgs.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	projectCmd.PersistentFlags().BoolVar(&cArgs.DryRun, "dry-run", false, "dry run does not write to pom.xml")
}
