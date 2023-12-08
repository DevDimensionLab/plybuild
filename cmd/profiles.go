package cmd

import (
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

const DEFAULT_TERMINAL_WIDTH = 80

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
	Short: "Manage profiles settings for ply",
	Long:  `Manage profiles settings for ply`,
	Run: func(cmd *cobra.Command, args []string) {
		if configOpts.UseProfile != "" {
			log.Infof("switching to config: %s", configOpts.UseProfile)
			err := config.SwitchProfile(configOpts.UseProfile)
			if err != nil {
				log.Fatalln(err)
			}
			profilePath, err := config.GetProfilesPathFor(configOpts.UseProfile)
			if err != nil {
				log.Fatalln(err)
			}
			ctx.LoadProfile(profilePath)
			return
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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config display in the terminal",
	Long:  `Config display in the terminal`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := ctx.LocalConfig.Config()
		if err != nil {
			log.Fatalln(err)
		}
		width, err := cmd.Flags().GetInt("width")
		if err != nil {
			log.Fatalln(err)
		}
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Fatalln(err)
		}

		if width != DEFAULT_TERMINAL_WIDTH {
			cfg.TerminalConfig.Width = width
		}

		if format != "" {
			cfg.TerminalConfig.Format = format
		}

		err = ctx.LocalConfig.UpdateLocalConfig(cfg)
		if err != nil {
			log.Fatalln(err)
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

	profilesCmd.AddCommand(configCmd)

	configCmd.Flags().IntP("width", "w", DEFAULT_TERMINAL_WIDTH, "Configure width of rendering in the terminal")
	configCmd.Flags().StringP("format", "f", "", "Configure format of rendering in the terminal: markdown")

}
