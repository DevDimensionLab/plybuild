package cmd

import (
	"co-pilot/pkg/analyze"
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "status",
	Short: "prints project status",
	Long:  `prints project status`,
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

		if localGroupId, err := analyze.GetLocalGroupId(model); err != nil {
			log.Fatalln(err)
		} else {
			log.Info("Local groupId domain is: " + localGroupId)
		}

		if err = upgrade.Dependencies(model, true); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Dependencies(model, false); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Kotlin(model); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.SpringBoot(model); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Plugin(model); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Clean(model); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().String("target", ".", "Optional target directory")
}
