package cmd

import (
	"co-pilot/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config settings for co-pilot",
	Long:  `config settings for co-pilot`,
}

var configShowCmd = &cobra.Command{
	Use:   "status",
	Short: "resets config for co-pilot",
	Long:  `resets config for co-pilot`,
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

var configDownloadCmd = &cobra.Command{
	Use:   "clone",
	Short: "clones global config for co-pilot",
	Long:  `clones global config for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Clone()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var configCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean current config for co-pilot",
	Long:  `clean current config for co-pilot`,
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
	configCmd.AddCommand(configDownloadCmd)
	configCmd.AddCommand(configCleanCmd)
}
