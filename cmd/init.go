package cmd

import (
	"co-pilot/pkg/upgrade"
	"github.com/spf13/cobra"
	"log"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initializes a project",
	Long:  `initializes a project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}
		err = upgrade.Init(targetDirectory)

		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().String("target", ".", "Optional target directory")
}
