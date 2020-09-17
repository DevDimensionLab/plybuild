package cmd

import (
	"co-pilot/pkg/bitbucket"
	"co-pilot/pkg/config"
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
		cfg, err := config.GetLocalConfig()
		if err != nil {
			log.Fatalln(err)
		}

		bitbucketHost := cfg.SourceProvider.Host
		personalAccessToken := cfg.SourceProvider.AccessToken

		if ("" == bitbucketHost) || ("" == personalAccessToken) {
			log.Fatalln("Command requires host and access-token in config-file")
		}

		err = bitbucket.SynchronizeAllRepos(bitbucketHost, personalAccessToken)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(bitbucketCmd)
	bitbucketCmd.AddCommand(bitbucketSyncCmd)
}
