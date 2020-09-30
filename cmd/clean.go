package cmd

import (
	"co-pilot/pkg/file"
	"github.com/spf13/cobra"
)

type CleanOpts struct {
	Ide bool
}

func (opts CleanOpts) Any() bool {
	return opts.Ide
}

var cleanOpts CleanOpts

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean files and folder in a project",
	Long:  `Clean files and folder in a project`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return OkHelp(cmd, cleanOpts.Any)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cleanOpts.Ide {
			report, err := file.CleanIntellijFiles(ctx.TargetDirectory, ctx.Recursive, ctx.DryRun)
			if err != nil {
				log.Fatalln(err)
			}
			log.Infof(report)
		}
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)

	cleanCmd.Flags().BoolVar(&cleanOpts.Ide, "ide", false, "removes .idea folders and *.iml files")
}
