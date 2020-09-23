package cmd

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/maven"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [OPTIONS]",
	Short: "analyze options",
	Long:  `Perform analyze on existing projects`,
}

var analyzeDepsCmd = &cobra.Command{
	Use:   "deps",
	Short: "analyze dependencies of project",
	Long:  `analyze dependencies of project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}

		pomFile := file.Path("%s/pom.xml", targetDirectory)

		if err = maven.ListUnusedAndUndeclared(pomFile); err != nil {
			log.Warnln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	analyzeCmd.AddCommand(analyzeDepsCmd)
	analyzeCmd.PersistentFlags().String("target", ".", "Optional target directory")
}
