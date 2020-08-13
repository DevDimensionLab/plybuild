package cmd

import (
	"co-pilot/pkg/analyze"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "analyzes a project",
	Long:  `analyzes a project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}

		model, err := analyze.GetModel(targetDirectory)
		if err != nil {
			log.Fatalln(err)
		}

		localGroupId, err := analyze.GetLocalGroupId(model)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Local groupId domain is: %s\n", localGroupId)

	},
}

func init() {
	RootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().String("target", ".", "Optional target directory")
}
