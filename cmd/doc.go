package cmd

import (
	"github.com/devdimensionlab/ply/pkg/webservice"
	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Opens documentation in default browser",
	Run: func(cmd *cobra.Command, args []string) {
		if err := webservice.OpenBrowser("https://devdimensionlab.github.io/"); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(docCmd)
}
