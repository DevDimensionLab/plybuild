package cmd

import (
	"co-pilot/pkg/config"
	"github.com/spf13/cobra"
)

var deprecatedCmd = &cobra.Command{
	Use:   "deprecated",
	Short: "deprecated settings for co-pilot",
	Long:  `deprecated settings for co-pilot`,
}

var deprecatedShowCmd = &cobra.Command{
	Use:   "status",
	Short: "show current deprecated for co-pilot",
	Long:  `show current deprecated for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		deprecated, err := config.GetDeprecated()
		if err != nil {
			log.Fatalln(err)
		}

		for _, dep := range deprecated.Data.Dependencies {
			log.Infof("deprecated dependency %s:%s", dep.GroupId, dep.ArtifactId)
		}
	},
}

func init() {
	RootCmd.AddCommand(deprecatedCmd)
	deprecatedCmd.AddCommand(deprecatedShowCmd)
}
