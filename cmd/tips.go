package cmd

import (
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var tipsCmd = &cobra.Command{
	Use:   "tips",
	Short: "Use a tips cheat-sheet to learn information faster",
	Long: `A concentrated version of everything you need to know for a topic, 
typically internal know-how that you can't find on the internet or via chatGPT`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			tipsListCmd.Run(cmd, args)
			return
		}
		tipsShowCmd.Run(cmd, args)
	},
}

var tipsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all tips (cheat-sheets) for current profile",
	Long:  `Lists all tips (cheat-sheets) for current profile`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Available tips (cheat-sheets):")
		cheatSheets, err := ctx.CloudConfig.CheatSheets()
		if err != nil {
			log.Fatalln(err)
		}
		for _, entry := range cheatSheets {
			log.Infof("- %s", strings.Replace(entry.Name(), ".md", "", 1))
		}
	},
}

var tipsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show tips (cheat-sheets)",
	Long:  `Show tips (cheat-sheets)`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := ctx.LocalConfig.Config()
		if err != nil {
			log.Fatalln(err)
		}

		if len(args) == 0 {
			log.Warnln("Missing tips (cheat sheet argument")
			tipsListCmd.Run(cmd, args)
			return
		}

		name := args[0]

		path := file.Path("%s/cheat-sheets/%s.md", ctx.CloudConfig.Implementation().Dir(), name)
		source, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to find any cheat-sheet file for [%s]: %s", name, path)
		}

		width := cfg.CheatSheetConfig.Width
		if 0 == width {
			width = 80
		}
		result := markdown.Render(string(source), width, 2)

		fmt.Println("\n" + string(result))
	},
}

var tipsShowConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Config display of tips (cheat-sheets)",
	Long:  `Config display of tips (cheat-sheets)`,
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

		width, err := getMandatoryString(cmd, "width")
		checkIfError(err)

		widthInt, err := strconv.Atoi(width)
		checkIfError(err)

		cfg.CheatSheetConfig.Width = widthInt
		ctx.LocalConfig.UpdateLocalConfig(cfg)
	},
}

func init() {
	RootCmd.AddCommand(tipsCmd)

	tipsCmd.AddCommand(tipsListCmd)
	tipsCmd.AddCommand(tipsShowCmd)
	tipsCmd.AddCommand(tipsShowConfigCmd)

	tipsShowConfigCmd.Flags().StringP("width", "w", "", "Configure width of tips (cheat-sheet) when displayed")
}
