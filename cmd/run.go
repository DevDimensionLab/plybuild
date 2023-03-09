package cmd

import (
	"github.com/devdimensionlab/co-pilot/pkg/shell"
	"github.com/spf13/cobra"
	"strings"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run scrips in scripts directory",
	Long:  `Run scrips in scripts directory`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := InitGlobals(cmd); err != nil {
			log.Fatalln(err)
		}
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		scripts, err := ctx.CloudConfig.Scripts()
		if err != nil {
			log.Fatalln(err)
		}

		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatalln(err)
		}

		if list || len(args) == 0 {
			log.Infoln("Available scripts:")
			for _, script := range scripts {
				log.Infoln("- " + strings.TrimSuffix(script.Name, ".sh"))
			}
			return
		}

		requestedScript := args[0]
		script, err := ctx.CloudConfig.Script(requestedScript)
		if err != nil {
			log.Fatalln(err)
		}

		envMap := make(map[string]interface{})
		localConfigMap, err := ctx.LocalConfig.ConfigAsMap()
		if err != nil {
			log.Fatalln(err)
		}
		envMap["local_config"] = localConfigMap

		if len(ctx.Projects) > 0 {
			projectConfig, err := ctx.Projects[0].ConfigAsMap()
			if err != nil {
				log.Fatalln(err)
			}
			envMap["project_config"] = projectConfig
		}

		output, err := shell.RunWithEnvironment(envMap, script.Path, args[1:]...)
		if err != nil {
			log.Fatalln(err)
		}
		println(string(output))
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	runCmd.Flags().Bool("list", false, "list available scripts")
}
