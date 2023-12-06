package cmd

import (
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/devdimensionlab/ply/pkg/file"
	"github.com/devdimensionlab/ply/pkg/tips"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var tipsCmd = &cobra.Command{
	Use:   "tips",
	Short: "Use tips to learn information faster",
	Long: `A concentrated version of things you need to know for a topic, 
typically internal know-how that you can't find on the internet`,
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
	Short: "Lists all tips for current profile",
	Long:  `Lists all tip for current profile`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Available tips:")
		tips, err := tips.List(ctx.CloudConfig)
		if err != nil {
			log.Fatalln(err)
		}
		for _, entry := range tips {
			log.Infof("- %s", strings.Replace(entry.Name(), ".md", "", 1))
		}
	},
}

var tipsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show tips",
	Long:  `Show tips`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			log.Warnln("Missing tips name argument")
			tipsListCmd.Run(cmd, args)
			return
		}

		name := args[0]

		tipsPath := file.Path("%s/%s.md", tips.LocalDir(ctx.CloudConfig), name)
		source, err := os.ReadFile(tipsPath)
		if err != nil {
			log.Fatalf("Failed to find any tips file for [%s]: %s", name, tipsPath)
		}

		terminalConfig, err := ctx.LocalConfig.GetTerminalConfig()
		if err != nil {
			log.Fatalln(err)
		}
		result := markdown.Render(string(source), terminalConfig.Width, 2)

		fmt.Println("\n" + string(result))

		log.Infoln("Local source: " + tipsPath)
		gCloudCfg, err := ctx.CloudConfig.GlobalCloudConfig()
		if err != nil {
			log.Fatalln(err)
		}
		log.Infof("Cloud source: %s\n", gCloudCfg.SourceFor(tips.TipsDir, name+".md"))
	},
}

func init() {
	RootCmd.AddCommand(tipsCmd)

	tipsCmd.AddCommand(tipsListCmd)
	tipsCmd.AddCommand(tipsShowCmd)
}
