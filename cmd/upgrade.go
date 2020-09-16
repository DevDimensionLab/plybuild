package cmd

import (
	"co-pilot/pkg/upgrade"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var groupId string
var artifactId string

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [OPTIONS]",
	Short: "Upgrade options",
	Long:  `Perform upgrade on existing projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
		ctx.FindAndPopulatePomModels()
	},
}

var upgradeSpringBootCmd = &cobra.Command{
	Use:   "spring-boot",
	Short: "Upgrade spring-boot to the latest version",
	Long:  `Upgrade spring-boot to the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachPomProject("upgrading spring-boot", func(model *pom.Model, args ...interface{}) error {
			return upgrade.SpringBoot(model)
		})
	},
}

var upgradeDependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Upgrade a specific dependency on a project",
	Long:  `Upgrade a specific dependency on a project`,
	Run: func(cmd *cobra.Command, args []string) {
		if groupId == "" || artifactId == "" {
			log.Fatal("--groupId (-g) and --artifactId (-a) must be set")
		}
		description := fmt.Sprintf("upgrading dependency %s:%s", groupId, artifactId)
		ctx.OnEachPomProject(description, func(model *pom.Model, args ...interface{}) error {
			return upgrade.Dependency(model, groupId, artifactId)
		})
	},
}

var upgrade2partyDependenciesCmd = &cobra.Command{
	Use:   "2party",
	Short: "Upgrade all 2party dependencies to project",
	Long:  `Upgrade all 2party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachPomProject("upgrading 2party", func(model *pom.Model, args ...interface{}) error {
			return upgrade.Dependencies(model, true)
		})
	},
}

var upgrade3partyDependenciesCmd = &cobra.Command{
	Use:   "3party",
	Short: "Upgrade all 3party dependencies to project",
	Long:  `Upgrade all 3party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachPomProject("upgrading 3party", func(model *pom.Model, args ...interface{}) error {
			return upgrade.Dependencies(model, false)
		})
	},
}

var upgradeKotlinCmd = &cobra.Command{
	Use:   "kotlin",
	Short: "Upgrade kotlin version in project",
	Long:  `Upgrade kotlin version in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachPomProject("upgrading kotlin", func(model *pom.Model, args ...interface{}) error {
			return upgrade.Kotlin(model)
		})
	},
}

var upgradePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Upgrade all plugins found in project",
	Long:  `Upgrade all plugins found in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachPomProject("upgrading plugins", func(model *pom.Model, args ...interface{}) error {
			return upgrade.Plugin(model)
		})
	},
}

var upgradeAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Upgrade everything in project",
	Long:  `Upgrade everything in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachPomProject("upgrading everything", func(model *pom.Model, args ...interface{}) error {
			return upgrade.All(model)
		})
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
	upgradeCmd.AddCommand(upgradeDependencyCmd)
	upgradeCmd.AddCommand(upgrade2partyDependenciesCmd)
	upgradeCmd.AddCommand(upgrade3partyDependenciesCmd)
	upgradeCmd.AddCommand(upgradeSpringBootCmd)
	upgradeCmd.AddCommand(upgradeKotlinCmd)
	upgradeCmd.AddCommand(upgradePluginsCmd)
	upgradeCmd.AddCommand(upgradeAllCmd)

	upgradeDependencyCmd.PersistentFlags().StringVarP(&groupId, "groupId", "g", "", "GroupId for upgrade")
	upgradeDependencyCmd.PersistentFlags().StringVarP(&artifactId, "artifactId", "a", "", "ArtifactId for upgrade")

	upgradeCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	upgradeCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	upgradeCmd.PersistentFlags().BoolVar(&ctx.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	upgradeCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "dry run does not write to pom.xml")
}
