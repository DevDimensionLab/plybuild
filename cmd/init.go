package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/upgrade"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes a project",
	Long:  `initializes a project`,
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

		initConfig, err := config.GenerateConfig(model)
		if err != nil {
			log.Fatalln(err)
		}

		configFile := fmt.Sprintf("%s/co-pilot.json", targetDirectory)
		log.Infof("writes co-pilot.json config file to %s", configFile)
		if err = config.WriteConfig(initConfig, configFile); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().String("target", ".", "Optional target directory")
}
