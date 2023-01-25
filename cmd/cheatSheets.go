/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"os"

	"github.com/devdimensionlab/co-pilot/pkg/file"
	"github.com/spf13/cobra"
	"strings"
)

// cheatSheetsCmd represents the cheatSheets command
var cheatSheetsCmd = &cobra.Command{
	Use:   "cheatSheets",
	Short: "Use a cheat sheet to learn information faster",
	Long: `A concentrated version of everything you need to know for a topic, 
typically internal knowhow that you can't find on the internet`,

	Run: func(cmd *cobra.Command, args []string) {
		cheatSheetsListCmd.Run(cmd, args)
	},
}

var cheatSheetsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all cheat-sheets for current profile",
	Long:  `Lists all cheat-sheets for current profile`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Available cheat-sheets:")
		cheatSheets, _ := ctx.CloudConfig.CheatSheets()
		for _, entry := range cheatSheets {
			log.Infof("- %s", strings.Replace(entry.Name(), ".md", "", 1))
		}
	},
}

var cheatSheetsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cheat-sheets",
	Long:  `Show cheat-sheets`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		name, err := getMandatoryString(cmd, "name")
		checkIfError(err)

		path := file.Path("%s/cheat-sheets/%s.md", ctx.CloudConfig.Implementation().Dir(), name)
		source, err := os.ReadFile(path)
		checkIfError(err)
		result := markdown.Render(string(source), 80, 6)

		fmt.Println()
		fmt.Println(string(result))
	},
}

func init() {
	RootCmd.AddCommand(cheatSheetsCmd)

	cheatSheetsCmd.AddCommand(cheatSheetsListCmd)
	cheatSheetsCmd.AddCommand(cheatSheetsShowCmd)
	cheatSheetsShowCmd.Flags().StringP("name", "n", "", "Name of cheat-sheet to show")
}
