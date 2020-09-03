package cmd

import (
	"co-pilot/pkg/maven"
	"github.com/spf13/cobra"
)

var mavenCmd = &cobra.Command{
	Use:   "maven [OPTIONS]",
	Short: "Maven options",
	Long:  `Various maven helper commands`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var mavenRepositoriesCmd = &cobra.Command{
	Use:   "repositories",
	Short: "List repositories",
	Long:  `List repositories`,
	Run: func(cmd *cobra.Command, args []string) {
		repos, err := maven.GetRepositories()
		if err != nil {
			log.Fatalln(err)
		}

		for _, profileRepo := range repos.Profile {
			log.Infof("found profile repo: %s", profileRepo)
		}

		for _, mirrorRepo := range repos.Mirror {
			log.Infof("found mirror repo: %s", mirrorRepo)
		}

		log.Infof("fallback repo: %s", repos.Fallback)
	},
}

func init() {
	RootCmd.AddCommand(mavenCmd)
	mavenCmd.AddCommand(mavenRepositoriesCmd)
}
