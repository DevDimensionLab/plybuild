package cmd

import (
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/maven"
	"github.com/devdimensionlab/plybuild/pkg/spring"
	"github.com/devdimensionlab/plybuild/pkg/template"
	"github.com/devdimensionlab/plybuild/pkg/webservice"
	"github.com/devdimensionlab/plybuild/pkg/webservice/api"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Short:   "Builds a maven project with ply files and formatting",
	Long:    `Builds a maven project with ply files and formatting`,
	Aliases: []string{"generate"},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := OpenDocumentationWebsite(cmd, "commands/build"); err != nil {
			log.Fatalln(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var orderConfig config.ProjectConfiguration
		var err error

		// fetch user input config
		var bootVersion, _ = cmd.Flags().GetString("boot-version")
		overrideGroupId, _ := cmd.Flags().GetString("group-id")
		overrideArtifactId, _ := cmd.Flags().GetString("artifact-id")
		overridePackage, _ := cmd.Flags().GetString("package")
		overrideName, _ := cmd.Flags().GetString("name")
		overrideApplicationName, _ := cmd.Flags().GetString("application-name")
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
		if overridePackage != "" {
			orderConfig.Package = overridePackage
		}
		if overrideName != "" {
			orderConfig.Name = overrideName
		}
		if overrideApplicationName != "" {
			orderConfig.ApplicationName = overrideApplicationName
		}

		err = spring.Validate(orderConfig)
		if err != nil {
			log.Fatalln(err)
		}

		// try to set spring-boot version manually
		if bootVersion == "" {
			bootVersion = orderConfig.Settings.MaxSpringBootVersion
		}
		// download from start.spring.io to targetDirectory
		err = spring.DownloadInitializer(ctx.TargetDirectory, spring.UrlValuesFrom(bootVersion, orderConfig))
		if err != nil {
			log.Fatalln(err)
		}

		// cleanup unwanted files from downloaded content
		spring.DeleteDemoFiles(ctx.TargetDirectory, orderConfig)

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
		ctx.OnEachMavenProject("removes version tags", maven.ChangeVersionToPropertyTags())

		// upgrade all ... maybe?
		disableUpgrade, _ := cmd.Flags().GetBool("disable-upgrading")
		if !disableUpgrade {
			ctx.OnEachMavenProject("upgrading everything",
				maven.UpgradeKotlin(),
				maven.UpgradeParent(),
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

var listTemplatesCmd = &cobra.Command{
	Use:   "list-templates",
	Short: "Lists all available templates",
	Long:  `Lists all available templates`,
	Run: func(cmd *cobra.Command, args []string) {

		markdownFormat, err := cmd.Flags().GetBool("markdown")
		if err != nil {
			log.Fatalln(err)
		}
		saveOutput, err := cmd.Flags().GetBool("save")
		if err != nil {
			log.Fatalln(err)
		}
		templates, err := ctx.CloudConfig.Templates()
		if err != nil {
			log.Fatalln(err)
		}
		terminalConfig, err := ctx.LocalConfig.GetTerminalConfig()
		if err != nil {
			log.Fatalln(err)
		}

		if markdownFormat || terminalConfig.Format == "markdown" {
			markdownDocument, err := template.ListAsMarkdown(ctx.CloudConfig, templates)
			if err != nil {
				log.Fatalln(err)
			}

			markdownForTerminal := markdown.Render(markdownDocument, terminalConfig.Width, 2)
			fmt.Println("\n" + string(markdownForTerminal))

			if saveOutput {
				fileRef, err := template.SaveTemplateListMarkdown(ctx.CloudConfig, markdownDocument)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println("File saved to " + fileRef)
			}

			gCloudCfg, err := ctx.CloudConfig.GlobalCloudConfig()
			if err != nil {
				log.Fatalln(err)
			}
			cloudSource := gCloudCfg.SourceFor(template.TemplatesDir, "README.md")
			if err != nil {
				log.Fatalln(err)
			}
			log.Infoln("Cloud source: " + cloudSource)
		} else {
			for _, folder := range templates {
				log.Infof("%s - %s (%s)", folder.Name, folder.Project.Config.Description, folder.Project.Config.Language)
			}
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
	RootCmd.AddCommand(buildCmd)

	buildCmd.AddCommand(listTemplatesCmd)

	buildCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	buildCmd.PersistentFlags().BoolVar(&ctx.ForceCloudSync, "cloud-sync", true, "Cloud sync")
	buildCmd.PersistentFlags().Bool("disable-upgrading", false, "dont upgrade dependencies")
	buildCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")
	buildCmd.Flags().String("config-file", "ply.json", "Optional config file")
	buildCmd.Flags().String("boot-version", "", "Defines spring-boot version to use")
	buildCmd.Flags().String("group-id", "", "Overrides groupId from config file")
	buildCmd.Flags().String("artifact-id", "", "Overrides artifactId from config file")
	buildCmd.Flags().String("package", "", "Overrides package from config file")
	buildCmd.Flags().String("name", "", "Overrides name from config file")
	buildCmd.Flags().String("application-name", "", "Overrides applicationName from config file")

	listTemplatesCmd.Flags().Bool("markdown", false, "Outputs templates as markdown in the terminal")
	listTemplatesCmd.Flags().Bool("save", false, "Saves the template markdown doc to cloud-config template-folder")

}
