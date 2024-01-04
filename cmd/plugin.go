package cmd

import (
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Plugin functionality for plybuild",
	Long:  `Plugin functionality for plybuild`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return OkHelp(cmd, removeOpts.Any)
	},
}

func init() {
	RootCmd.AddCommand(pluginCmd)
}
