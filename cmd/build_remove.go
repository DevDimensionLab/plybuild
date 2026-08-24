package cmd

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/spf13/cobra"
)

type RemoveOpts struct {
	Generated bool
	Intellij  bool
}

func (opts RemoveOpts) Any() bool {
	return opts.Generated || opts.Intellij
}

var removeOpts RemoveOpts

var buildRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes files and folders in a ply build project",
	Long:  `Removes files and folders in a ply build project`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := InitGlobals(cmd); err != nil {
			return err
		}
		return OkHelp(cmd, removeOpts.Any)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		// force defaults
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		if removeOpts.Generated {
			if err := removeGeneratedFiles(force); err != nil {
				return err
			}
		}
		if removeOpts.Intellij {
			if err := removeIntellijFiles(force); err != nil {
				return err
			}
		}
		return nil
	},
}

func removeGeneratedFiles(force bool) error {
	answer, err := promptForContinue(fmt.Sprintf("Are you sure you want to delete contents of: %s", ctx.TargetDirectory), force)
	if err != nil {
		return err
	}
	if !answer {
		return nil
	}

	log.Infof(fmt.Sprintf("Deleting all contents from: %s", ctx.TargetDirectory))
	if err := file.ClearDir(ctx.TargetDirectory, []string{".idea", "ply.json", ".iml", ".git"}); err != nil {
		return err
	}
	return nil
}

func removeIntellijFiles(force bool) error {
	answer, err := promptForContinue(
		fmt.Sprintf("Are you sure you want to delete the intellij files in: %s", ctx.TargetDirectory), force)
	if err != nil {
		return err
	}
	if !answer {
		return nil
	}

	report, err := file.RemoveIntellijFiles(ctx.TargetDirectory, ctx.Recursive, ctx.DryRun)
	if err != nil {
		return err
	}
	log.Infof(report)
	return nil
}

func init() {
	buildCmd.AddCommand(buildRemoveCmd)

	buildRemoveCmd.Flags().BoolVar(&removeOpts.Generated, "generated", false, "removes plybuild generated files and folders")
	buildRemoveCmd.Flags().BoolVar(&removeOpts.Intellij, "intellij", false, "removes .idea folders and *.iml files")
	buildRemoveCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
}
