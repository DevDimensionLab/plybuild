package cmd

import (
	"co-pilot/pkg/analyze"
	"co-pilot/pkg/upgrade"
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

		model, err := analyze.GetModel(targetDirectory)
		if err != nil {
			log.Fatalln(err)
		}

		if localGroupId, err := analyze.GetLocalGroupId(model); err != nil {
			log.Fatalln(err)
		} else {
			log.Info("Local groupId domain is: " + localGroupId)
		}

		if err = upgrade.Dependencies(targetDirectory, true, true); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Dependencies(targetDirectory, false, true); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Kotlin(targetDirectory, true); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.SpringBoot(targetDirectory, true); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Plugin(targetDirectory, true); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().String("target", ".", "Optional target directory")
}
