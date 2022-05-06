package cmd

import (
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

type ConfigOpts struct {
	Sync       bool
	Reset      bool
	Show       bool
	Edit       bool
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
			ctx.LoadProfile(configOpts.UseProfile)
		}

		if configOpts.Edit {
			var editor = os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim"
			}
			cmd := exec.Command(editor, ctx.LocalConfig.FilePath())
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				log.Fatalln(err)
			}
		}

		if configOpts.Sync {
			if err := ctx.CloudConfig.Refresh(ctx.LocalConfig); err != nil {
				log.Fatalln(err)
			}
		}

		if configOpts.Reset {
			if err := ctx.LocalConfig.TouchFile(); err != nil {
				log.Fatalln(err)
			}
		}

		if !configOpts.Reset || !configOpts.Sync || !configOpts.Edit {
			if err := ctx.LocalConfig.Print(); err != nil {
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
	profilesCmd.Flags().BoolVar(&configOpts.Edit, "edit", false, "edit active profile local config")
	profilesCmd.Flags().BoolVar(&configOpts.Reset, "reset", false, "reset local config")
}
