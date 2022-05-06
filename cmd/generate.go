package cmd

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/co-pilot/pkg/config"
	"github.com/devdimensionlab/co-pilot/pkg/file"
	"github.com/devdimensionlab/co-pilot/pkg/maven"
	"github.com/devdimensionlab/co-pilot/pkg/spring"
	"github.com/devdimensionlab/co-pilot/pkg/template"
	"github.com/devdimensionlab/co-pilot/pkg/webservice"
	"github.com/devdimensionlab/co-pilot/pkg/webservice/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Initializes a maven project with co-pilot files and formatting",
	Long:  `Initializes a maven project with co-pilot files and formatting`,
	Run: func(cmd *cobra.Command, args []string) {
		var orderConfig config.ProjectConfiguration
		var err error

		// fetch user input config
		overrideGroupId, _ := cmd.Flags().GetString("group-id")
		overrideArtifactId, _ := cmd.Flags().GetString("artifact-id")
		interactive, _ := cmd.Flags().GetBool("interactive")
		jsonConfigFile, _ := cmd.Flags().GetString("config-file")

		if jsonConfigFile != "" {
			orderConfig, err = config.InitProjectConfigurationFromFile(jsonConfigFile)
		}
		if orderConfig.Profile != "" {
			profilesPath, err := config.GetProfilesPathFor(orderConfig.Profile)
			if err != nil {
				log.Fatalln(err)
			}
			ctx.LoadProfile(profilesPath)
		}

		// sync cloud config
		if ctx.ForceCloudSync {
			if err := ctx.CloudConfig.Refresh(ctx.LocalConfig); err != nil {
				log.Fatalln(err)
			}
		}

		if interactive {
			interactiveWebService(&orderConfig)
		}
		if err != nil {
			log.Fatalln(err)
		}
		if err = orderConfig.Validate(); err != nil {
			log.Fatalln(err)
		}

		// validate templates
		var cloudTemplates []config.CloudTemplate
		if orderConfig.Templates != nil {
			cloudTemplates, err = ctx.CloudConfig.ValidTemplatesFrom(orderConfig.Templates)
			if err != nil {
				log.Fatalln(err)
			}
		}

		// check for override of groupId and artifactId
		if overrideArtifactId != "" {
			orderConfig.ArtifactId = overrideArtifactId
			//orderConfig.Package = fmt.Sprintf("%s.%s", orderConfig.GroupId, orderConfig.ArtifactId)
		}
		if overrideGroupId != "" {
			orderConfig.GroupId = overrideGroupId
			orderConfig.Package = orderConfig.GroupId
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

		// cleanup unwanted files from downloaded content
		spring.DeleteDemoFiles(ctx.TargetDirectory)

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
		err = project.GitInit(fmt.Sprintf("Adds project %s", project.Config.Name))
		if err != nil {
			log.Fatalln(err)
		}

		// merge templates into the newly created project
		if cloudTemplates != nil {
			for _, cloudTemplate := range cloudTemplates {
				if err := template.MergeTemplate(cloudTemplate, project, true); err != nil {
					log.Fatalln(err)
				}
			}
			// git commit
			err = project.GitCommit(fmt.Sprintf("Adds templates to %s", project.Config.Name))
			if err != nil {
				log.Fatalln(err)
			}
		}

		// load project into context
		if err := ctx.FindAndPopulateMavenProjects(); err != nil {
			log.Fatalln(err)
		}

		// format version
		ctx.OnEachProject("removes version tags", maven.ChangeVersionToPropertyTags())

		// upgrade all ... maybe?
		disableUpgrade, _ := cmd.Flags().GetBool("disable-upgrading")
		if !disableUpgrade {
			ctx.OnEachProject("upgrading everything",
				maven.UpgradeKotlin(),
				spring.UpgradeSpringBoot(),
				maven.Upgrade2PartyDependencies(),
				maven.Upgrade3PartyDependencies(),
				maven.UpgradePlugins(),
			)
		} else {
			// only apply sorting and writing
			if err = project.SortAndWritePom(); err != nil {
				log.Fatalln(err)
			}
		}

		// git commit
		err = project.GitCommit(fmt.Sprintf("Cleans up and upgrades for project %s", project.Config.Name))
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var generateCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans a maven project with co-pilot files and formatting",
	Long:  `Cleans a maven project with co-pilot files and formatting`,
	Run: func(cmd *cobra.Command, args []string) {

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
		result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
		if result == "no" {
			return
		}
		log.Infof(fmt.Sprintf("Deleting all contents from: %s", ctx.TargetDirectory))
		if err := file.ClearDir(ctx.TargetDirectory, []string{".idea", "co-pilot.json", ".iml"}); err != nil {
			log.Fatalln(err)
		}
	},
}

func interactiveWebService(orderConfig *config.ProjectConfiguration) {
	ioResp, err := spring.GetRoot()
	if err != nil {
		log.Fatalln(err)
	}
	api.GOptions = api.GenerateOptions{
		ProjectConfig: orderConfig,
		CloudConfig:   ctx.CloudConfig,
		IoResponse:    ioResp,
	}
	webservice.InitAndBlockStandalone(webservice.Generate, api.CallbackChannel)
}

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.AddCommand(generateCleanCmd)

	generateCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	generateCmd.PersistentFlags().BoolVar(&ctx.ForceCloudSync, "cloud-sync", true, "Cloud sync")
	generateCmd.PersistentFlags().Bool("disable-upgrading", false, "dont upgrade dependencies")
	generateCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")
	generateCmd.Flags().String("config-file", "co-pilot.json", "Optional config file")
	generateCmd.Flags().String("group-id", "", "Overrides groupId from config file")
	generateCmd.Flags().String("artifact-id", "", "Overrides artifactId from config file")
}
