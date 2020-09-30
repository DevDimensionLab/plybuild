package cmd

import (
	"github.com/spf13/cobra"
)

type ConfigOpts struct {
	Sync  bool
	Reset bool
	Show  bool
}

func (configOpts ConfigOpts) Any() bool {
	return configOpts.Sync || configOpts.Reset
}

var configOpts ConfigOpts

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config settings for co-pilot",
	Long:  `Config settings for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		if !configOpts.Any() || configOpts.Show {
			if err := localCfg.Print(); err != nil {
				log.Fatalln(err)
			}
		}

		if configOpts.Sync {
			if err := cloudCfg.Refresh(localCfg); err != nil {
				log.Fatalln(err)
			}
		}

		if configOpts.Reset {
			if err := localCfg.TouchFile(); err != nil {
				log.Fatalln(err)
			} else {
				log.Infoln("new config file is generated")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVar(&configOpts.Sync, "cloud-sync", false, "sync with cloud config repo")
	configCmd.Flags().BoolVar(&configOpts.Show, "show", false, "show local config")
	configCmd.Flags().BoolVar(&configOpts.Reset, "reset", false, "reset local config")
}
