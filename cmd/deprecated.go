package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/deprecated"
	"co-pilot/pkg/file"
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
		if cArgs.Recursive, cArgs.Err = cmd.Flags().GetBool("recursive"); cArgs.Err != nil {
			log.Fatalln(cArgs.Err)
		}
		if cArgs.TargetDirectory, cArgs.Err = cmd.Flags().GetString("target"); cArgs.Err != nil {
			log.Fatalln(cArgs)
		}
		if cArgs.Overwrite, cArgs.Err = cmd.Flags().GetBool("overwrite"); cArgs.Err != nil {
			log.Fatalln(cArgs.Err)
		}
		if cArgs.Recursive {
			if cArgs.PomFiles, cArgs.Err = file.FindAll("pom.xml", cArgs.TargetDirectory); cArgs.Err != nil {
				log.Fatalln(cArgs.Err)
			}
		} else {
			cArgs.PomFiles = append(cArgs.PomFiles, cArgs.TargetDirectory+"/pom.xml")
		}
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

var deprecatedStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows all deprecated dependencies for a project co-pilot",
	Long:  `Shows all deprecated dependencies for a project co-pilot`,
	Run: func(cmd *cobra.Command, args []string) {
		d, err := config.GetDeprecated()
		if err != nil {
			log.Fatalln(err)
		}

		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("working on pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Fatalln(err)
			}

			err = deprecated.UpgradeDeprecated(model, d, pomFileToTargetDirectory(pomFile), false)
			if err != nil {
				log.Fatalln(err)
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
				log.Fatalln(err)
			}

			err = deprecated.UpgradeDeprecated(model, d, pomFileToTargetDirectory(pomFile), true)
			if err != nil {
				log.Fatalln(err)
			}

			if err := write(pomFile, model); err != nil {
				log.Fatalln(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(deprecatedCmd)
	deprecatedCmd.AddCommand(deprecatedShowCmd)
	deprecatedCmd.AddCommand(deprecatedStatusCmd)
	deprecatedCmd.AddCommand(deprecatedUpgradeCmd)

	deprecatedCmd.PersistentFlags().Bool("recursive", false, "turn on recursive mode")
	deprecatedCmd.PersistentFlags().String("target", ".", "Optional target directory")
	deprecatedCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
}
