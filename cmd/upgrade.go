package cmd

import (
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/maven"
	"github.com/co-pilot-cli/co-pilot/pkg/spring"
	"github.com/co-pilot-cli/co-pilot/pkg/template"
	"github.com/co-pilot-cli/co-pilot/pkg/webservice"
	"github.com/co-pilot-cli/co-pilot/pkg/webservice/api"
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
		if err := SyncCloudConfig(); err != nil {
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
		ctx.OnEachProject("upgrading spring-boot", spring.UpgradeSpringBoot())
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
		ctx.OnEachProject(description, maven.UpgradeDependency(groupId, artifactId))
	},
}

var upgrade2partyDependenciesCmd = &cobra.Command{
	Use:   "2party",
	Short: "Upgrade all 2party dependencies to project",
	Long:  `Upgrade all 2party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachProject("upgrading 2party", maven.Upgrade2PartyDependencies())
	},
}

var upgrade3partyDependenciesCmd = &cobra.Command{
	Use:   "3party",
	Short: "Upgrade all 3party dependencies to project",
	Long:  `Upgrade all 3party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachProject("upgrading 3party", maven.Upgrade3PartyDependencies())
	},
}

var upgradeKotlinCmd = &cobra.Command{
	Use:   "kotlin",
	Short: "Upgrade kotlin version in project",
	Long:  `Upgrade kotlin version in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachProject("upgrading kotlin", maven.UpgradeKotlin())
	},
}

var upgradePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Upgrade all plugins found in project",
	Long:  `Upgrade all plugins found in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachProject("upgrading plugins", maven.UpgradePlugins())
	},
}

var upgradeDeprecatedCmd = &cobra.Command{
	Use:   "deprecated",
	Short: "Remove and replace deprecated dependencies in a project",
	Long:  `Remove and replace deprecated dependencies in a project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachProject("removes and replaces deprecated dependencies", func(project config.Project) error {
			return upgradeDeprecated(project)
		})
	},
}

var upgradeAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Upgrade everything in project",
	Long:  `Upgrade everything in project`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx.OnEachProject("upgrading everything",
			maven.UpgradeKotlin(),
			spring.UpgradeSpringBoot(),
			maven.Upgrade2PartyDependencies(),
			maven.Upgrade3PartyDependencies(),
			maven.UpgradePlugins(),
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

	upgradeDependencyCmd.PersistentFlags().StringVarP(&groupId, "groupId", "g", "", "GroupId for upgrade")
	upgradeDependencyCmd.PersistentFlags().StringVarP(&artifactId, "artifactId", "a", "", "ArtifactId for upgrade")

	upgradeCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	upgradeCmd.PersistentFlags().BoolVar(&ctx.ForceCloudSync, "cloud-sync", false, "force cloud sync")
	upgradeCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "dry run does not write to pom.xml")
	upgradeCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
}

func upgradeDeprecated(project config.Project) error {
	templates, err := maven.RemoveDeprecated(cloudCfg, project.Type.Model())
	if err != nil {
		return err
	} else {
		template.MergeTemplates(templates, project)
	}
	return nil
}
