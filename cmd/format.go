package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/upgrade"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format functionality for a project",
	Long:  `Format functionality for a project`,
}

var formatInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Formats pom.xml and sorts dependencies",
	Long:  `Formats pom.xml and sorts dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		err = upgrade.Init(model, pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		configFile := fmt.Sprintf("%s/co-pilot.json", targetDirectory)
		if !file.Exists(configFile) {
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
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = clean.VersionToPropertyTags(model); err != nil {
			log.Fatalln(err)
		}

		var writeToFile = pomFile
		if !overwrite {
			writeToFile = targetDirectory + "/pom.xml.new"
		}
		if err = upgrade.SortAndWrite(model, writeToFile); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(formatCmd)
	formatCmd.AddCommand(formatInitCmd)
	formatCmd.AddCommand(formatVersionCmd)
	formatCmd.PersistentFlags().String("target", ".", "Optional target directory")
	formatCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
}
