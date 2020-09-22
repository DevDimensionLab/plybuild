package cmd

import (
	"co-pilot/pkg/file"
	"github.com/spf13/cobra"
)

var ideCmd = &cobra.Command{
	Use:   "ide",
	Short: "IDE (editor) functionality",
	Long:  `IDE (editor) functionality`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd); err != nil {
			log.Fatalln(err)
		}
	},
}

var ideCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes IDE files and folders, e.g. .idea and *.iml",
	Long:  `Removes IDE files and folders, e.g. .idea and *.iml`,
	Run: func(cmd *cobra.Command, args []string) {
		report, err := file.CleanIntellijFiles(ctx.TargetDirectory, ctx.Recursive, ctx.DryRun)
		if err != nil {
			log.Fatalln(err)
		}
		log.Infof(report)
	},
}

func init() {
	RootCmd.AddCommand(ideCmd)
	ideCmd.AddCommand(ideCleanCmd)

	ideCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	ideCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	ideCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "disables delete of files and folders")
}
