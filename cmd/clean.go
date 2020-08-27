package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [OPTIONS]",
	Short: "Clean options",
	Long:  `Perform clean on existing projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var cleanManualVersion = &cobra.Command{
	Use:   "manual-version",
	Short: "removes manual versions from dependencies",
	Long:  `removes manual versions from dependencies`,
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

		if err = clean.ManualVersion(model); err != nil {
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

var cleanBlacklist = &cobra.Command{
	Use:   "blacklist",
	Short: "removes deps that are blacklisted",
	Long:  `removes deps that are blacklisted`,
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

		if err = clean.BlacklistedDependencies(model); err != nil {
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
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.AddCommand(cleanManualVersion)
	cleanCmd.AddCommand(cleanBlacklist)
	cleanCmd.PersistentFlags().String("target", ".", "Optional target directory")
	cleanCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
}
