package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/upgrade"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format functionality for a project",
	Long:  `Format functionality for a project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cArgs.Recursive, cArgs.Err = cmd.Flags().GetBool("recursive"); cArgs.Err != nil {
			log.Fatalln(cArgs.Err)
		}
		if cArgs.TargetDirectory, cArgs.Err = cmd.Flags().GetString("target"); cArgs.Err != nil {
			log.Fatalln(cArgs)
		}
		if cArgs.Overwrite, cArgs.Err = cmd.Flags().GetBool("overwrite"); cArgs.Err != nil {
			log.Fatalln(cArgs.Err)
		}
		if cArgs.Recursive {
			if cArgs.PomFiles, cArgs.Err = file.FindAll("pom.xml", cArgs.TargetDirectory); cArgs.Err != nil {
				log.Fatalln(cArgs.Err)
			}
		} else {
			cArgs.PomFiles = append(cArgs.PomFiles, cArgs.TargetDirectory+"/pom.xml")
		}
	},
}

var formatInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Formats pom.xml and sorts dependencies",
	Long:  `Formats pom.xml and sorts dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("initializes pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Fatalln(err)
			}

			err = upgrade.Init(model, pomFile)
			if err != nil {
				log.Fatalln(err)
			}
			configFile := fmt.Sprintf("%s/co-pilot.json", pomFileToTargetDirectory(pomFile))
			initConfig, err := config.GenerateConfig(model)
			if err != nil {
				log.Fatalln(err)
			}

			log.Infof("writes co-pilot.json config file to %s", configFile)
			if err = config.WriteConfig(initConfig, configFile); err != nil {
				log.Fatalln(err)
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
				log.Fatalln(err)
			}

			if err = clean.VersionToPropertyTags(model); err != nil {
				log.Fatalln(err)
			}
			write(pomFile, model)
		}
	},
}

func init() {
	RootCmd.AddCommand(formatCmd)
	formatCmd.AddCommand(formatInitCmd)
	formatCmd.AddCommand(formatVersionCmd)

	formatCmd.PersistentFlags().Bool("recursive", false, "turn on recursive mode")
	formatCmd.PersistentFlags().String("target", ".", "Optional target directory")
	formatCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
}
