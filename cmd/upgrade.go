package cmd

import (
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/maven"
	"github.com/devdimensionlab/co-pilot/pkg/template"
	"github.com/devdimensionlab/co-pilot/pkg/webservice"
	"github.com/devdimensionlab/co-pilot/pkg/webservice/api"
	"github.com/spf13/cobra"
)

var groupId string
var artifactId string

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [OPTIONS]",
	Short: "Upgrade options",
	Long:  `Perform upgrade on existing projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := SyncActiveProfileCloudConfig(); err != nil {
			log.Warnln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
}

var upgradeSpringBootCmd = &cobra.Command{
	Use:   "spring-boot",
	Short: "Upgrade spring-boot to the latest version",
	Long:  `Upgrade spring-boot to the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachMavenProject("upgrading spring-boot", maven.UpgradeParent())
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
		ctx.OnEachMavenProject(description, maven.UpgradeDependency(groupId, artifactId))
	},
}

var upgrade2partyDependenciesCmd = &cobra.Command{
	Use:   "2party",
	Short: "Upgrade all 2party dependencies to project",
	Long:  `Upgrade all 2party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachMavenProject("upgrading 2party", maven.Upgrade2PartyDependencies())
	},
}

var upgrade3partyDependenciesCmd = &cobra.Command{
	Use:   "3party",
	Short: "Upgrade all 3party dependencies to project",
	Long:  `Upgrade all 3party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachMavenProject("upgrading 3party", maven.Upgrade3PartyDependencies())
	},
}

var upgradeKotlinCmd = &cobra.Command{
	Use:   "kotlin",
	Short: "Upgrade kotlin version in project",
	Long:  `Upgrade kotlin version in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachMavenProject("upgrading kotlin", maven.UpgradeKotlin())
	},
}

var upgradePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Upgrade all plugins found in project",
	Long:  `Upgrade all plugins found in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachMavenProject("upgrading plugins", maven.UpgradePlugins())
	},
}

var upgradeDeprecatedCmd = &cobra.Command{
	Use:   "deprecated",
	Short: "Remove and replace deprecated dependencies in a project",
	Long:  `Remove and replace deprecated dependencies in a project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachMavenProject("removes and replaces deprecated dependencies", func(repository maven.Repository, project config.Project) error {
			return upgradeDeprecated(project)
		})
	},
}

var upgradeAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Upgrade everything in project",
	Long:  `Upgrade everything in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachMavenProject("upgrading everything",
			maven.UpgradeKotlin(),
			maven.UpgradeParent(),
			maven.Upgrade2PartyDependencies(),
			maven.Upgrade3PartyDependencies(),
			maven.UpgradePlugins(),
		)
	},
}

var upgradeWithVersionsCmd = &cobra.Command{
	Use:   "with-versions",
	Short: "Upgrade using mvn versions in a project",
	Long:  `Upgrade using mvn versions in a project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.DryRun = true
		ctx.OnEachMavenProject("upgrading properties using maven versions",
			maven.UpgradeKotlinWithVersions(),
			maven.UpgradeDependenciesWithVersions(),
		)
	},
}

var upgradeInteractiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactively upgrade the project",
	Long:  `Interactively upgrade the project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnRootProject("starting interactive upgrade",
			webservice.InitAndBlockProject(webservice.Upgrade, api.CallbackChannel))
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
	upgradeCmd.AddCommand(upgradeDeprecatedCmd)
	upgradeCmd.AddCommand(upgradeAllCmd)
	upgradeCmd.AddCommand(upgradeInteractiveCmd)
	upgradeCmd.AddCommand(upgradeWithVersionsCmd)

	upgradeDependencyCmd.PersistentFlags().StringVarP(&groupId, "groupId", "g", "", "GroupId for upgrade")
	upgradeDependencyCmd.PersistentFlags().StringVarP(&artifactId, "artifactId", "a", "", "ArtifactId for upgrade")

	upgradeCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	upgradeCmd.PersistentFlags().BoolVar(&ctx.ForceCloudSync, "cloud-sync", false, "force cloud sync")
	upgradeCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "dry run does not write to pom.xml")
	upgradeCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	upgradeCmd.PersistentFlags().BoolVar(&ctx.StealthMode, "stealth", false, "use alternative pom.xml writer")
}

func upgradeDeprecated(project config.Project) error {
	templates, err := maven.RemoveDeprecated(project.CloudConfig, project.Type.Model())
	if err != nil {
		return err
	}
	template.MergeTemplates(templates, project)
	return nil
}
