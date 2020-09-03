package cmd

import (
	"co-pilot/pkg/analyze"
	"co-pilot/pkg/clean"
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints project status",
	Long:  `Prints project status`,
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

		if secondPartyGroupId, err := analyze.GetSecondPartyGroupId(model); err != nil {
			log.Fatalln(err)
		} else {
			log.Info("2party groupId = " + secondPartyGroupId)
		}
		if err = upgrade.Kotlin(model); err != nil {
			log.Warnln(err)
		}
		if err = upgrade.SpringBoot(model); err != nil {
			log.Warnln(err)
		}
		if err = upgrade.Dependencies(model, true); err != nil {
			log.Warnln(err)
		}
		if err = upgrade.Dependencies(model, false); err != nil {
			log.Warnln(err)
		}
		if err = upgrade.Plugin(model); err != nil {
			log.Warnln(err)
		}
		if err = clean.SpringManualVersion(model); err != nil {
			log.Warnln(err)
		}
		if err = clean.VersionToPropertyTags(model); err != nil {
			log.Warnln(err)
		}
		if err = analyze.Undeclared(pomFile); err != nil {
			log.Warnln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
	statusCmd.Flags().String("target", ".", "Optional target directory")
}
