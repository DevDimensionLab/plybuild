package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/maven"
	"co-pilot/pkg/spring"
	"co-pilot/pkg/template"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Initializes a maven project with co-pilot files and formatting",
	Long:  `Initializes a maven project with co-pilot files and formatting`,
	Run: func(cmd *cobra.Command, args []string) {
		// remove targetDirectory
		if err := os.RemoveAll(ctx.TargetDirectory); err != nil {
			log.Fatalln(err)
		}

		// sync cloud config
		if err := cloudCfg.Refresh(localCfg); err != nil {
			log.Fatalln(err)
		}

		// fetch user input config
		jsonConfigFile, _ := cmd.Flags().GetString("config-file")
		if jsonConfigFile == "" {
			log.Fatalln("--config-file flag is required")
		}
		orderConfig, err := config.InitProjectConfigurationFromFile(jsonConfigFile)
		if err != nil {
			log.Fatalln(err)
		}
		err = spring.Validate(orderConfig)
		if err != nil {
			log.Fatalln(err)
		}

		// download from start.spring.io to targetDirectory
		err = spring.DownloadInitializer(ctx.TargetDirectory, spring.UrlValuesFrom(orderConfig))
		if err != nil {
			log.Fatalln(err)
		}

		// populate applicationName field in config
		if err := orderConfig.FindApplicationName(ctx.TargetDirectory); err != nil {
			log.Errorln(err)
		}

		// write project config to targetDir
		projectConfigFile := config.ProjectConfigPath(ctx.TargetDirectory)
		if err := orderConfig.WriteTo(projectConfigFile); err != nil {
			log.Fatalln(err)
		}

		// load the newly created project
		project, err := config.InitProjectFromDirectory(ctx.TargetDirectory)
		if err != nil {
			log.Fatalln(err)
		}

		// git init project
		err = project.GitInit()
		if err != nil {
			log.Fatalln(err)
		}

		// merge templates into the newly created project
		if orderConfig.Templates != nil {
			templates, err := cloudCfg.ValidTemplatesFrom(orderConfig.Templates)
			if err != nil {
				log.Fatalln(err)
			}
			for _, t := range templates {
				if err := template.With(logger.Context()).MergeTemplate(t, project); err != nil {
					log.Fatalln(err)
				}
			}
			// git commit
			err = project.GitCommit("Adds templates")
			if err != nil {
				log.Fatalln(err)
			}
		}

		// format version
		log.Info(logger.Info(fmt.Sprintf("formatting %s", project.Type.FilePath())))
		if err = maven.ChangeVersionToPropertyTagsOnModel(project.Type.Model()); err != nil {
			log.Fatalln(err)
		}

		// upgrade all
		log.Info(logger.Info(fmt.Sprintf("upgrading all on %s", project.Type.FilePath())))
		if err = UpgradeAll(project); err != nil {
			log.Fatalln(err)
		}

		// sorting and writing
		if err = project.SortAndWritePom(); err != nil {
			log.Fatalln(err)
		}

		// git commit
		err = project.GitCommit("Clean up and upgrades")
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	generateCmd.Flags().String("config-file", "", "Optional config file")
}
