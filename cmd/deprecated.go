package cmd

import (
	"co-pilot/pkg/deprecated"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/service"
	"fmt"
	"github.com/spf13/cobra"
)

var deprecatedCmd = &cobra.Command{
	Use:   "deprecated",
	Short: "Deprecated detection and patching functionalities for projects",
	Long:  `Deprecated detection and patching functionalities for projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		ctx.FindAndPopulatePomModels()
	},
}

var deprecatedShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows all deprecated dependencies for co-pilot",
	Long:  `Shows all deprecated dependencies for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cloudCfg.ListDeprecated(); err != nil {
			log.Fatalln(err)
		}
	},
}

var deprecatedUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades deprecated dependencies for a project co-pilot",
	Long:  `Upgrades deprecated dependencies for a project co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := cloudCfg.Deprecated()
		if err != nil {
			log.Fatalln(err)
		}

		for pomFile, model := range ctx.PomModels {
			log.Info(logger.White(fmt.Sprintf("upgrading deprecated dependencies for pom file %s", pomFile)))

			templates, err := deprecated.RemoveDeprecated(model, d)
			if err != nil {
				log.Warnln(err)
				continue
			}

			deprecated.ApplyTemplates(cloudCfg, templates, service.PomFileToTargetDirectory(pomFile))

			if !ctx.DryRun {
				if err := service.Write(ctx.Overwrite, pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(deprecatedCmd)
	deprecatedCmd.AddCommand(deprecatedShowCmd)
	deprecatedCmd.AddCommand(deprecatedUpgradeCmd)

	deprecatedCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	deprecatedCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	deprecatedCmd.PersistentFlags().BoolVar(&ctx.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	deprecatedCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "dry run does not write to pom.xml")
}
