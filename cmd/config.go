package cmd

import (
	"co-pilot/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config settings for co-pilot",
	Long:  `Config settings for co-pilot`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows local-config for co-pilot",
	Long:  `Shows local-config for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.GetLocalConfig()
		if err != nil {
			log.Fatalln(err)
		}

		err = config.PrintLocalConfig(c)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var configSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronizes cloud config for co-pilot",
	Long:  `Synchronizes cloud config for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Clone()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets local config for co-pilot",
	Long:  `Resets local config for co-pilot with empty values`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.TouchLocalConfigFile()
		if err != nil {
			log.Fatalln(err)
		}

		log.Infoln("new config file is generated")
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSyncCmd)
	configCmd.AddCommand(configResetCmd)
}
