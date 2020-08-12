package cmd

import (
	"co-pilot/pkg/upgrade"
	"github.com/spf13/cobra"
	"log"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [OPTIONS]",
	Short: "Upgrade options",
	Long:  `Perform upgrade on existing projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var upgradeSpringBootCmd = &cobra.Command{
	Use:   "spring-boot",
	Short: "upgrade spring-boot to the latest version",
	Long:  `upgrade spring-boot to the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Println(err)
		}
		err = upgrade.SpringBoot(targetDirectory)
		if err != nil {
			log.Println(err)
		}
	},
}

var upgradeDependenciesCmd = &cobra.Command{
	Use:   "deps",
	Short: "upgrade dependencies to project",
	Long:  `upgrade dependencies to project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var upgradeAllDependenciesCmd = &cobra.Command{
	Use:   "all",
	Short: "upgrade all dependencies to project",
	Long:  `upgrade all dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Println(err)
		}
		err = upgrade.Dependencies(targetDirectory, true)
		err = upgrade.Dependencies(targetDirectory, false)
		if err != nil {
			log.Println(err)
		}
	},
}

var upgrade2partyDependenciesCmd = &cobra.Command{
	Use:   "2party",
	Short: "upgrade 2party dependencies to project",
	Long:  `upgrade 2party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Println(err)
		}
		err = upgrade.Dependencies(targetDirectory, true)
		if err != nil {
			log.Println(err)
		}
	},
}

var upgrade3partyDependenciesCmd = &cobra.Command{
	Use:   "3party",
	Short: "upgrade 3party dependencies to project",
	Long:  `upgrade 3party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Println(err)
		}
		err = upgrade.Dependencies(targetDirectory, false)
		if err != nil {
			log.Println(err)
		}
	},
}

var upgradeKotlinCmd = &cobra.Command{
	Use:   "kotlin",
	Short: "upgrade kotlin version in project",
	Long:  `upgrade kotlin version in project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Println(err)
		}
		err = upgrade.Kotlin(targetDirectory)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)

	upgradeCmd.AddCommand(upgradeSpringBootCmd)
	upgradeSpringBootCmd.Flags().String("target", ".", "Optional target directory")

	upgradeCmd.AddCommand(upgradeDependenciesCmd)
	upgradeDependenciesCmd.AddCommand(upgradeAllDependenciesCmd)
	upgradeAllDependenciesCmd.Flags().String("target", ".", "Optional target directory")
	upgradeDependenciesCmd.AddCommand(upgrade2partyDependenciesCmd)
	upgrade2partyDependenciesCmd.Flags().String("target", ".", "Optional target directory")
	upgradeDependenciesCmd.AddCommand(upgrade3partyDependenciesCmd)
	upgrade3partyDependenciesCmd.Flags().String("target", ".", "Optional target directory")

	upgradeCmd.AddCommand(upgradeKotlinCmd)
	upgradeKotlinCmd.Flags().String("target", ".", "Optional target directory")
}
