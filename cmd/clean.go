package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [OPTIONS]",
	Short: "Clean options",
	Long:  `Perform clean on existing projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var cleanSpringManualVersion = &cobra.Command{
	Use:   "spring-manual-version",
	Short: "removes manual versions from spring dependencies",
	Long:  `removes manual versions from spring dependencies`,
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

		if err = clean.SpringManualVersion(model); err != nil {
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

var cleanVersionProps = &cobra.Command{
	Use:   "version-properties",
	Short: "removes version tags and replaces them with property tags",
	Long:  `removes version tags and replaces them with property tags`,
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

var cleanStatus = &cobra.Command{
	Use:   "status",
	Short: "outputs status",
	Long:  `outputs status, but does not write anything`,
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

		if err = clean.VersionToPropertyTags(model); err != nil {
			log.Warnln(err)
		}

		if err = clean.SpringManualVersion(model); err != nil {
			log.Warnln(err)
		}

		if err = clean.Undeclared(pomFile, model); err != nil {
			log.Warnln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.AddCommand(cleanSpringManualVersion)
	cleanCmd.AddCommand(cleanVersionProps)
	cleanCmd.AddCommand(cleanStatus)
	//cleanCmd.AddCommand(cleanBlacklist)
	cleanCmd.PersistentFlags().String("target", ".", "Optional target directory")
	cleanCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
}
