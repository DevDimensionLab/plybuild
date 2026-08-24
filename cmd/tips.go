package cmd

import (
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/tips"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return tipsListCmd.RunE(cmd, args)
		}
		return tipsShowCmd.RunE(cmd, args)
	},
}

var tipsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all tips for current profile",
	Long:  `Lists all tip for current profile`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitGlobals(cmd); err != nil {
			return err
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("Available tips:")
		tips, err := tips.List(ctx.CloudConfig)
		if err != nil {
			return err
		}
		for _, entry := range tips {
			log.Infof("- %s", strings.Replace(entry.Name(), ".md", "", 1))
		}
		return nil
	},
}

var tipsShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show tips",
	Long:  `Show tips`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitGlobals(cmd); err != nil {
			return err
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			log.Warnln("Missing tips name argument")
			return tipsListCmd.RunE(cmd, args)
		}

		name := args[0]

		tipsPath := file.Path("%s/%s.md", tips.LocalDir(ctx.CloudConfig), name)
		source, err := os.ReadFile(tipsPath)
		if err != nil {
			return fmt.Errorf("failed to find any tips file for [%s]: %s: %w", name, tipsPath, err)
		}

		terminalConfig, err := ctx.LocalConfig.GetTerminalConfig()
		if err != nil {
			return err
		}
		result := markdown.Render(string(source), terminalConfig.Width, 2)

		fmt.Println("\n" + string(result))

		log.Infoln("Local source: " + tipsPath)
		gCloudCfg, err := ctx.CloudConfig.GlobalCloudConfig()
		if err != nil {
			return err
		}
		log.Infof("Cloud source: %s\n", gCloudCfg.SourceFor(tips.TipsDir, name+".md"))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(tipsCmd)

	tipsCmd.AddCommand(tipsListCmd)
	tipsCmd.AddCommand(tipsShowCmd)
}
