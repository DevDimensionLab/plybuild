package cmd

import (
	"co-pilot/pkg/springio"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [OPTIONS]",
	Short: "Downloads options",
	Long:  `Downloads options`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var downloadCliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Downloads spring-cli",
	Long:  `Downloads spring-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := springio.DownloadCli(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	downloadCmd.AddCommand(downloadCliCmd)
}
