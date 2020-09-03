package cmd

import (
	"co-pilot/pkg/bitbucket"
	config2 "co-pilot/pkg/config"
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
		config, err := config2.GetLocalConfig()
		if err != nil {
			log.Fatalln(err)
		}

		bitbucketHost := config.SourceProvider.Host
		personalAccessToken := config.SourceProvider.AccessToken

		if ("" == bitbucketHost) || ("" == personalAccessToken) {
			log.Fatalln("Command requires host and access-token in config-file")
		}

		err = bitbucket.Synchronize(bitbucketHost, personalAccessToken)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(bitbucketCmd)
	bitbucketCmd.AddCommand(bitbucketSyncCmd)
}
