package cmd

import (
	"github.com/spf13/cobra"
)

var mavenCmd = &cobra.Command{
	Use:   "maven [OPTIONS]",
	Short: "Maven options",
	Long:  `Various maven helper commands`,
}

var mavenRepositoriesCmd = &cobra.Command{
	Use:   "repositories",
	Short: "List repositories",
	Long:  `List repositories`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(mavenCmd)
	//lintCmd.AddCommand(mavenRepositoriesCmd)
}
