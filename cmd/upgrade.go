package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"spring-boot-co-pilot/pkg/spring"
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
		err = spring.UpgradeSpringBoot(targetDirectory)
		if err != nil {
			log.Println(err)
		}
	},
}

var upgradeDependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "upgrade dependencies to project",
	Long:  `upgrade dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Println(err)
		}
		err = spring.UpgradeSpringDependencies(targetDirectory)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
	upgradeCmd.AddCommand(upgradeSpringBootCmd)
	upgradeCmd.AddCommand(upgradeDependenciesCmd)
	upgradeSpringBootCmd.Flags().String("target", ".", "Optional target directory")
	upgradeDependenciesCmd.Flags().String("target", ".", "Optional target directory")
}
