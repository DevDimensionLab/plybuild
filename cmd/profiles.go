package cmd

import (
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/spf13/cobra"
)

type ConfigOpts struct {
	Sync       bool
	Reset      bool
	Show       bool
	UseProfile string
}

func (configOpts ConfigOpts) Any() bool {
	return configOpts.Sync
}

var configOpts ConfigOpts

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage profiles settings for co-pilot",
	Long:  `Manage profiles settings for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		if configOpts.UseProfile != "" {
			log.Infof("switching to config: %s", configOpts.UseProfile)
			err := config.SwitchProfile(configOpts.UseProfile)
			if err != nil {
				log.Fatalln(err)
			}
			loadConfig()
		}

		if configOpts.Sync {
			if err := cloudCfg.Refresh(localCfg); err != nil {
				log.Fatalln(err)
			}
		}

		if configOpts.Show {
			if err := localCfg.Print(); err != nil {
				log.Fatalln(err)
			}
		}

		if configOpts.Reset {
			if err := localCfg.TouchFile(); err != nil {
				log.Fatalln(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(profilesCmd)

	profilesCmd.Flags().BoolVar(&configOpts.Sync, "cloud-sync", false, "sync with cloud config repo")
	profilesCmd.Flags().StringVar(&configOpts.UseProfile, "use", "", "switch to profile")
	profilesCmd.Flags().BoolVar(&configOpts.Show, "show", false, "show local config")
	profilesCmd.Flags().BoolVar(&configOpts.Reset, "reset", false, "reset local config")
}
