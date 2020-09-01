package cmd

import (
	"co-pilot/pkg/analyze"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [OPTIONS]",
	Short: "analyze options",
	Long:  `Perform analyze on existing projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
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

		pomFile := targetDirectory + "/pom.xml"

		if err = analyze.Undeclared(pomFile); err != nil {
			log.Warnln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	analyzeCmd.AddCommand(analyzeDepsCmd)
	analyzeCmd.PersistentFlags().String("target", ".", "Optional target directory")
}
