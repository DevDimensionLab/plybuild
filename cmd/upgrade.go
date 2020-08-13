package cmd

import (
	"co-pilot/pkg/upgrade"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.SpringBoot(model); err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.SortAndWrite(model, pomFile); err != nil {
			log.Fatalln(err)
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
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		err = upgrade.Dependencies(model, true)
		if err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.SortAndWrite(model, pomFile); err != nil {
			log.Fatalln(err)
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
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.Dependencies(model, false); err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.SortAndWrite(model, pomFile); err != nil {
			log.Fatalln(err)
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
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.Kotlin(model); err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.SortAndWrite(model, pomFile); err != nil {
			log.Fatalln(err)
		}
	},
}

var upgradePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "upgrade plugins found in project",
	Long:  `upgrade plugins found in project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.Plugin(model); err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.SortAndWrite(model, pomFile); err != nil {
			log.Fatalln(err)
		}
	},
}

var upgradeCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "cleans versions set in project that are redundant or unnecessary",
	Long:  `cleans versions set in project that are redundant or unnecessary`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.Clean(model); err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.SortAndWrite(model, pomFile); err != nil {
			log.Fatalln(err)
		}
	},
}

var upgradeAllDependenciesCmd = &cobra.Command{
	Use:   "all",
	Short: "upgrade everything in project",
	Long:  `upgrade everything in project`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}
		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.Dependencies(model, true); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Dependencies(model, false); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Kotlin(model); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.SpringBoot(model); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Plugin(model); err != nil {
			log.Fatalln(err)
		}
		if err = upgrade.Clean(model); err != nil {
			log.Fatalln(err)
		}

		if err = upgrade.SortAndWrite(model, pomFile); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
	upgradeCmd.AddCommand(upgradeAllDependenciesCmd)
	upgradeCmd.AddCommand(upgrade2partyDependenciesCmd)
	upgradeCmd.AddCommand(upgrade3partyDependenciesCmd)
	upgradeCmd.AddCommand(upgradeSpringBootCmd)
	upgradeCmd.AddCommand(upgradeKotlinCmd)
	upgradeCmd.AddCommand(upgradePluginsCmd)
	upgradeCmd.AddCommand(upgradeCleanCmd)

	upgradeCmd.PersistentFlags().String("target", ".", "Optional target directory")
}
