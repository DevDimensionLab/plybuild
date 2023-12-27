package cmd

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
)

var buildClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears a maven project with ply files and formatting",
	Long:  `Clears a maven project with ply files and formatting`,
	Run: func(cmd *cobra.Command, args []string) {

		// force defaults
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalln(err)
		}

		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to delete contents of: %s [yes/no]", ctx.TargetDirectory),
			Templates: templates,
			Validate: func(input string) error {
				if len(input) <= 0 || (input != "yes" && input != "no") {
					return errors.New("please enter 'yes' or 'no'")
				}
				return nil
			},
		}

		if !force {
			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				os.Exit(1)
			}
			if result == "no" {
				return
			}
		}

		log.Infof(fmt.Sprintf("Deleting all contents from: %s", ctx.TargetDirectory))
		if err := file.ClearDir(ctx.TargetDirectory, []string{".idea", "ply.json", ".iml"}); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	buildCmd.AddCommand(buildClearCmd)
}
