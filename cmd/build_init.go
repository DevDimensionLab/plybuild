package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes build with ply files and formatting",
	Long:  `Initializes build project with ply files and formatting`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, project := range ctx.Projects {
			if project.Type == nil {
				log.Warnf("no project type defined for %s:", project.Path)
				continue
			}

			log.Info(fmt.Sprintf("formating pom file %s", project.Type.FilePath()))
			if !ctx.DryRun {
				err := project.InitProjectConfiguration()
				if err != nil {
					log.Warnln(err)
					continue
				}

				if err := project.Config.WriteTo(project.ConfigFile); err != nil {
					log.Warnln(err)
					continue
				}

				if err := project.SortAndWritePom(); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

func init() {
	buildCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().BoolVarP(&ctx.Recursive, "recursive", "r", false, "turn on recursive mode")
	initCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	initCmd.PersistentFlags().BoolVar(&ctx.DryRun, "dry-run", false, "dry run does not write to pom.xml")

	initCmd.Flags().String("config-file", "", "Optional config file")
}
