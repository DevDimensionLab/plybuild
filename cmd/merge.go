package cmd

import (
	"co-pilot/pkg/maven"
	"co-pilot/pkg/merge"
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

		fromPomFile, err := cmd.Flags().GetString("from")
		if err != nil {
			log.Fatalln(err)
		}
		if fromPomFile == "" {
			log.Errorln("missing valid --from flag for pom.xml to merge from")
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

var mergeTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Merges two text files",
	Long:  `Merges two text files`,
	Run: func(cmd *cobra.Command, args []string) {
		fromFile, err := cmd.Flags().GetString("from")
		if err != nil {
			log.Fatalln(err)
		}
		if fromFile == "" {
			log.Errorln("missing valid --from file flag")
			os.Exit(-1)
		}

		toFile, err := cmd.Flags().GetString("to")
		if err != nil {
			log.Fatalln(err)
		}
		if toFile == "" {
			log.Errorln("missing valid --to file flag")
			os.Exit(-1)
		}

		if err := merge.TextFiles(fromFile, toFile); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(mergeCmd)
	mergeCmd.AddCommand(mergePomCmd)
	mergeCmd.AddCommand(mergeTextCmd)
	mergeCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
	mergeCmd.PersistentFlags().String("from", "", "file to merge")
	mergePomCmd.PersistentFlags().String("target", ".", "Optional target directory")
	mergeTextCmd.PersistentFlags().String("to", "", "target file to merge to")
}
