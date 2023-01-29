package cmd

import (
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"github.com/devdimensionlab/co-pilot/pkg/tips"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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

		cfg, err := ctx.LocalConfig.Config()
		if err != nil {
			log.Fatalln(err)
		}

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

		width := cfg.TipsConfig.Width
		if 0 == width {
			width = 80
		}
		result := markdown.Render(string(source), width, 2)

		fmt.Println("\n" + string(result))

		log.Infoln("Local source: " + tipsPath)
		cloudSource, err := tips.CloudSource(name, ctx.CloudConfig)
		if err != nil {
			log.Fatalln(err)
		}
		log.Infoln("Cloud source: " + cloudSource)
	},
}

var tipsShowConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Config display of tips",
	Long:  `Config display of tips`,
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

		cfg.TipsConfig.Width = widthInt
		ctx.LocalConfig.UpdateLocalConfig(cfg)
	},
}

func init() {
	RootCmd.AddCommand(tipsCmd)

	tipsCmd.AddCommand(tipsListCmd)
	tipsCmd.AddCommand(tipsShowCmd)
	tipsCmd.AddCommand(tipsShowConfigCmd)

	tipsShowConfigCmd.Flags().StringP("width", "w", "", "Configure width of tips when rendered")
}