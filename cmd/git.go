package cmd

import (
	"fmt"
	"github.com/co-pilot-cli/co-pilot/pkg/config"
	"github.com/co-pilot-cli/co-pilot/pkg/shell"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git commands",
	Long:  `Git commands`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return OkHelp(cmd, cleanOpts.Any)
	},
}

var gitInstallHooksCmd = &cobra.Command{
	Use:   "install-hooks",
	Short: "install git hooks from cloud config",
	Long:  `install git hooks from cloud config`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
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
		gitHooksFolderName := "git-hooks"
		ctx.DryRun = true
		ctx.OnEachProject("installing git hooks",
			func(project config.Project) error {
				hooks, err := project.CloudConfig.GitHookFiles(gitHooksFolderName)
				if err != nil {
					return err
				}
				gitHookPath := fmt.Sprintf("%s/%s", project.CloudConfig.Implementation().Dir(), gitHooksFolderName)
				return shell.InstallGitHooks(gitHookPath, hooks, project.Path)
			},
		)
	},
}

func init() {
	RootCmd.AddCommand(gitCmd)

	gitCmd.AddCommand(gitInstallHooksCmd)
	gitCmd.PersistentFlags().BoolVar(&ctx.ForceCloudSync, "cloud-sync", false, "force cloud sync")
	gitCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	gitCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "optional target directory")
}
