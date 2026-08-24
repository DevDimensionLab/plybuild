package cmd

import (
	"errors"
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

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := ctx.LocalConfig.Config()
		if err != nil {
			return err
		}

		bitbucketHost := cfg.SourceProvider.Host
		personalAccessToken := cfg.SourceProvider.AccessToken

		if bitbucketHost == "" || personalAccessToken == "" {
			return errors.New("command requires host and access-token in config-file")
		}

		err = bitbucket.With(logger.Context(), bitbucketHost, personalAccessToken).SynchronizeAllRepos(cfg.SourceProvider.ExcludeProjects)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	pluginCmd.AddCommand(bitbucketCmd)
	bitbucketCmd.AddCommand(bitbucketSyncCmd)
}
