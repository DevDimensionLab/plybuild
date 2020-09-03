package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/deprecated"
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var deprecatedCmd = &cobra.Command{
	Use:   "deprecated",
	Short: "Deprecated functionalities for co-pilot",
	Long:  `Deprecated functionalities for co-pilot`,
}

var deprecatedShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows all deprecated dependencies for co-pilot",
	Long:  `Shows all deprecated dependencies for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		deprecated, err := config.GetDeprecated()
		if err != nil {
			log.Fatalln(err)
		}

		for _, dep := range deprecated.Data.Dependencies {
			log.Infof("deprecated dependency %s:%s", dep.GroupId, dep.ArtifactId)
		}
	},
}

var deprecatedStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows all deprecated dependencies for a project co-pilot",
	Long:  `Shows all deprecated dependencies for a project co-pilot`,
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

		d, err := config.GetDeprecated()
		if err != nil {
			log.Fatalln(err)
		}

		err = deprecated.FindDeprecated(model, d)
		if err != nil {
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
	RootCmd.AddCommand(deprecatedCmd)
	deprecatedCmd.AddCommand(deprecatedShowCmd)
	deprecatedCmd.AddCommand(deprecatedStatusCmd)
	deprecatedCmd.PersistentFlags().String("target", ".", "Optional target directory")
	deprecatedCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
}
