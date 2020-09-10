package cmd

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/upgrade"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [OPTIONS]",
	Short: "Upgrade options",
	Long:  `Perform upgrade on existing projects`,
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

var upgradeSpringBootCmd = &cobra.Command{
	Use:   "spring-boot",
	Short: "Upgrade spring-boot to the latest version",
	Long:  `Upgrade spring-boot to the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading spring-boot for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Fatalln(err)
			}

			if err = upgrade.SpringBoot(model); err != nil {
				log.Fatalln(err)
			}
			write(pomFile, model)
		}
	},
}

var upgrade2partyDependenciesCmd = &cobra.Command{
	Use:   "2party",
	Short: "Upgrade 2party dependencies to project",
	Long:  `Upgrade 2party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading 2party for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Fatalln(err)
			}
			if err = upgrade.Dependencies(model, true); err != nil {
				log.Fatalln(err)
			}
			write(pomFile, model)
		}
	},
}

var upgrade3partyDependenciesCmd = &cobra.Command{
	Use:   "3party",
	Short: "Upgrade 3party dependencies to project",
	Long:  `Upgrade 3party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading 3party for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Fatalln(err)
			}

			if err = upgrade.Dependencies(model, false); err != nil {
				log.Fatalln(err)
			}
			write(pomFile, model)
		}
	},
}

var upgradeKotlinCmd = &cobra.Command{
	Use:   "kotlin",
	Short: "Upgrade kotlin version in project",
	Long:  `Upgrade kotlin version in project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading kotlin for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Fatalln(err)
			}

			if err = upgrade.Kotlin(model); err != nil {
				log.Fatalln(err)
			}
			write(pomFile, model)
		}
	},
}

var upgradePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Upgrade plugins found in project",
	Long:  `Upgrade plugins found in project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading plugins for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Fatalln(err)
			}
			if err = upgrade.Plugin(model); err != nil {
				log.Fatalln(err)
			}
			write(pomFile, model)
		}
	},
}

var upgradeAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Upgrade everything in project",
	Long:  `Upgrade everything in project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading all for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Fatalln(err)
			}
			upgrade.All(model)
			write(pomFile, model)
		}
	},
}

var upgradeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "upgrade status for prosject",
	Long:  `upgrade status for prosject`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("working on pom file %s", f)))
			model, err := pom.GetModelFrom(f)
			if err != nil {
				log.Fatalln(err)
			}
			upgrade.All(model)
		}
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
	upgradeCmd.AddCommand(upgrade2partyDependenciesCmd)
	upgradeCmd.AddCommand(upgrade3partyDependenciesCmd)
	upgradeCmd.AddCommand(upgradeSpringBootCmd)
	upgradeCmd.AddCommand(upgradeKotlinCmd)
	upgradeCmd.AddCommand(upgradePluginsCmd)
	upgradeCmd.AddCommand(upgradeAllCmd)
	upgradeCmd.AddCommand(upgradeStatusCmd)

	upgradeCmd.PersistentFlags().Bool("recursive", false, "turn on recursive mode")
	upgradeCmd.PersistentFlags().String("target", ".", "Optional target directory")
	upgradeCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
}
