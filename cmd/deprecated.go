package cmd

import (
	"co-pilot/pkg/deprecated"
	"co-pilot/pkg/maven"
	"co-pilot/pkg/template"
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

		ctx.OnEachPomProject("removes version tags", func(pair maven.PomPair, args ...interface{}) error {
			templates, err := deprecated.RemoveDeprecated(pair.Model, d)
			if err != nil {
				log.Warnln(err)
			} else {
				template.Apply(cloudCfg, templates, maven.PomFileToTargetDirectory(pair.PomFile))
			}
			return nil
		})
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
