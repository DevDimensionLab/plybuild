package cmd

import (
	"github.com/devdimensionlab/plybuild/pkg/bitbucket"
	"github.com/devdimensionlab/plybuild/pkg/logger"
	"github.com/spf13/cobra"
)

var bitbucketCmd = &cobra.Command{
	Use:   "bitbucket",
	Short: "Bitbucket functionality",
	Long:  `Bitbucket functionality`,
}

var bitbucketSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronizes projects from bitbucket",
	Long:  `Synchronizes projects from bitbucket`,

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := ctx.LocalConfig.Config()
		if err != nil {
			log.Fatalln(err)
		}

		bitbucketHost := cfg.SourceProvider.Host
		personalAccessToken := cfg.SourceProvider.AccessToken

		if ("" == bitbucketHost) || ("" == personalAccessToken) {
			log.Fatalln("Command requires host and access-token in config-file")
		}

		err = bitbucket.With(logger.Context(), bitbucketHost, personalAccessToken).SynchronizeAllRepos(cfg.SourceProvider.ExcludeProjects)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	pluginCmd.AddCommand(bitbucketCmd)
	bitbucketCmd.AddCommand(bitbucketSyncCmd)
}
