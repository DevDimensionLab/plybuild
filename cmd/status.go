package cmd

import (
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/spf13/cobra"
)

type StatusOpts struct {
	Show bool
}

var statusOpts StatusOpts

func (sOpts StatusOpts) Any() bool {
	return sOpts.Show
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status functionality for a project",
	Long:  `Status functionality for a project`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := OpenDocumentationWebsite(cmd, "commands/status"); err != nil {
			log.Fatalln(err)
		}
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
	Run: func(cmd *cobra.Command, args []string) {
		if !statusOpts.Any() || statusOpts.Show {
			ctx.DryRun = true
			ctx.OnEachMavenProject("project status",
				maven.UpgradeKotlin(),
				maven.UpgradeParent(),
				maven.Upgrade2PartyDependencies(),
				maven.Upgrade3PartyDependencies(),
				maven.UpgradePlugins(),
				maven.ChangeVersionToPropertyTags(),
				maven.CleanManualVersions(),
				maven.StatusDeprecated(),
			)
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
	statusCmd.PersistentFlags().BoolVar(&statusOpts.Show, "show", false, "show project status")

	statusCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	statusCmd.PersistentFlags().BoolVar(&ctx.ForceCloudSync, "cloud-sync", false, "force cloud sync")
	statusCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")
}
