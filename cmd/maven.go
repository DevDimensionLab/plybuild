package cmd

import (
	"co-pilot/pkg/maven"
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
		if err := maven.ListRepositories(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(mavenCmd)
	mavenCmd.AddCommand(mavenRepositoriesCmd)
}
