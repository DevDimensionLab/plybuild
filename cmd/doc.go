package cmd

import (
	"github.com/devdimensionlab/co-pilot/pkg/webservice"
	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Documentation",
	Run: func(cmd *cobra.Command, args []string) {
		if err := webservice.OpenBrowser("https://devdimensionlab.github.io/"); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(docCmd)
}
