package cmd

import (
	"co-pilot/pkg/maven"
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
	"os"
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge functionalities for files to a project",
	Long:  `Merge functionalities for files to a project`,
}

var mergePomCmd = &cobra.Command{
	Use:   "pom",
	Short: "Merges a pom-file into a project",
	Long:  `Merges a pom-file into a project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			log.Fatalln(err)
		}

		fromPomFile, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatalln(err)
		}
		if fromPomFile == "" {
			log.Errorln("missing valid --file flag for pom.xml to merge from")
			os.Exit(-1)
		}

		importModel, err := pom.GetModelFrom(fromPomFile)
		if err != nil {
			log.Fatalln(err)
		}

		toPomFile := targetDirectory + "/pom.xml"
		projectModel, err := pom.GetModelFrom(toPomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = maven.Merge(importModel, projectModel); err != nil {
			log.Fatalln(err)
		}

		var writeToFile = toPomFile
		if !overwrite {
			writeToFile = targetDirectory + "/pom.xml.new"
		}
		if err = upgrade.SortAndWrite(projectModel, writeToFile); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(mergeCmd)
	mergeCmd.AddCommand(mergePomCmd)
	mergeCmd.PersistentFlags().String("target", ".", "Optional target directory")
	mergeCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
	mergePomCmd.Flags().String("file", "", "pom file to merge into project")
}
