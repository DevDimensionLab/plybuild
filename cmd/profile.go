package cmd

import (
	"io"
	"os"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/spf13/cobra"
)

const DefaultTerminalWidth = 80

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

type profileEditorDependencies struct {
	Process process.Dependencies
	Stdin   io.Reader
	Stdout  io.Writer
}

func systemProfileEditorDependencies() profileEditorDependencies {
	dependencies := profileEditorDependencies{Process: process.System()}
	dependencies.Stdin = os.Stdin
	dependencies.Stdout = os.Stdout
	return dependencies
}

func selectProfileEditor(editor string) string {
	if editor == "" {
		return "vim"
	}
	return editor
}

func runProfileEditor(dependencies profileEditorDependencies, editor, configPath string) error {
	return process.Execute(dependencies.Process, process.Command{
		Name:   editor,
		Args:   []string{configPath},
		Dir:    "",
		Stdin:  dependencies.Stdin,
		Stdout: dependencies.Stdout,
		Stderr: nil,
	})
}

var profileCmd = &cobra.Command{
	Use:     "profile",
	Short:   "Manage profiles settings for ply",
	Long:    `Manage profiles settings for ply`,
	Aliases: []string{"profiles"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if configOpts.UseProfile != "" {
			log.Infof("switching to config: %s", configOpts.UseProfile)
			err := config.SwitchProfile(configOpts.UseProfile)
			if err != nil {
				return err
			}
			profilePath, err := config.GetProfilesPathFor(configOpts.UseProfile)
			if err != nil {
				return err
			}
			ctx.LoadProfile(profilePath)
			return nil
		}

		if configOpts.Edit {
			editor := os.Getenv("EDITOR")
			editor = selectProfileEditor(editor)
			configPath := ctx.LocalConfig.FilePath()
			err := runProfileEditor(systemProfileEditorDependencies(), editor, configPath)
			if err != nil {
				return err
			}
		}

		if configOpts.Sync {
			if err := ctx.CloudConfig.Refresh(ctx.LocalConfig); err != nil {
				return err
			}
		}

		if configOpts.Reset {
			if err := ctx.LocalConfig.TouchFile(); err != nil {
				return err
			}
		}

		if !configOpts.Reset || !configOpts.Sync || !configOpts.Edit {
			if err := ctx.LocalConfig.Print(); err != nil {
				return err
			}
		}
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config display in the terminal",
	Long:  `Config display in the terminal`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return InitGlobals(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := ctx.LocalConfig.Config()
		if err != nil {
			return err
		}
		width, err := cmd.Flags().GetInt("width")
		if err != nil {
			return err
		}
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return err
		}

		if width != DefaultTerminalWidth {
			cfg.TerminalConfig.Width = width
		}

		if format != "" {
			cfg.TerminalConfig.Format = format
		}

		err = ctx.LocalConfig.UpdateLocalConfig(cfg)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(profileCmd)

	profileCmd.Flags().BoolVar(&configOpts.Sync, "cloud-sync", false, "sync with cloud config repo")
	profileCmd.Flags().StringVar(&configOpts.UseProfile, "use", "", "switch to profile")
	profileCmd.Flags().BoolVar(&configOpts.Show, "show", false, "show local config")
	profileCmd.Flags().BoolVar(&configOpts.Edit, "edit", false, "edit active profile local config")
	profileCmd.Flags().BoolVar(&configOpts.Reset, "reset", false, "reset local config")

	profileCmd.AddCommand(configCmd)

	configCmd.Flags().IntP("width", "w", DefaultTerminalWidth, "Configure width of rendering in the terminal")
	configCmd.Flags().StringP("format", "f", "", "Configure format of rendering in the terminal: markdown")

}
