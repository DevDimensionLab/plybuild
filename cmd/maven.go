package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"co-pilot/pkg/maven"
)

var mavenCmd = &cobra.Command{
	Use:   "maven [OPTIONS]",
	Short: "maven options",
	Long:  `Various maven helper commands`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var mavenRepositoriesCmd = &cobra.Command{
	Use:   "repositories",
	Short: "list repositories",
	Long:  `list repositories`,
	Run: func(cmd *cobra.Command, args []string) {
		repos, err := maven.GetRepositories()
		if err != nil {
			println(fmt.Sprintf("%v\n", err))
		}

		for _, repo := range repos {
			println(repo)
		}
	},
}

func init() {
	RootCmd.AddCommand(mavenCmd)
	mavenCmd.AddCommand(mavenRepositoriesCmd)
}
