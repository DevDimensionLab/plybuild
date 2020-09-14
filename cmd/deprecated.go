package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/deprecated"
	"co-pilot/pkg/logger"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var deprecatedCmd = &cobra.Command{
	Use:   "deprecated",
	Short: "Deprecated detection and patching functionalities for projects",
	Long:  `Deprecated detection and patching functionalities for projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd, args); err != nil {
			log.Fatalln(err)
		}
		populatePomFiles()
	},
}

var deprecatedShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows all deprecated dependencies for co-pilot",
	Long:  `Shows all deprecated dependencies for co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		deprecated, err := config.GetDeprecated()
		if err != nil {
			log.Fatalln(err)
		}

		for _, dep := range deprecated.Data.Dependencies {
			log.Infof("== deprecated dependency %s:%s ==", dep.GroupId, dep.ArtifactId)
			if dep.Associated.Dependencies != nil {
				for _, assoc := range dep.Associated.Dependencies {
					log.Infof("\t <= associated deprecated dependency %s:%s", assoc.GroupId, assoc.ArtifactId)
				}
			}
			if dep.ReplacementTemplates != nil {
				for _, repTemp := range dep.ReplacementTemplates {
					log.Infof("\t <= replacement template %s", repTemp)
				}
			}
		}
	},
}

var deprecatedUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades deprecated dependencies for a project co-pilot",
	Long:  `Upgrades deprecated dependencies for a project co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := config.GetDeprecated()
		if err != nil {
			log.Fatalln(err)
		}

		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading deprecated dependencies for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}

			err = deprecated.UpgradeDeprecated(model, d, pomFileToTargetDirectory(pomFile), true)
			if err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err := write(pomFile, model); err != nil {
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

	deprecatedCmd.PersistentFlags().BoolVarP(&cArgs.Recursive, "recursive", "r", false, "turn on recursive mode")
	deprecatedCmd.PersistentFlags().StringVar(&cArgs.TargetDirectory, "target", ".", "Optional target directory")
	deprecatedCmd.PersistentFlags().BoolVar(&cArgs.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	deprecatedCmd.PersistentFlags().BoolVar(&cArgs.DryRun, "dry-run", false, "dry run does not write to pom.xml")
}
