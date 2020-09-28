package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/maven"
	"co-pilot/pkg/spring"
	"co-pilot/pkg/template"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
	"os"
)

var springCmd = &cobra.Command{
	Use:   "spring [OPTIONS]",
	Short: "Spring boot tools",
	Long:  `Spring boot tools`,
}

var springInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Downloads and installs spring boot, and co-pilot templates, with default or provided settings",
	Long:  `Downloads and installs spring boot, and co-pilot templates, with default or provided settings`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.RemoveAll(ctx.TargetDirectory)

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

		// download cli
		if err := spring.CheckCli(localCfg); err != nil {
			log.Fatalln(err)
		}

		// execute cli with config
		msg, err := spring.RunCli(localCfg, spring.InitFrom(orderConfig, ctx.TargetDirectory)...)
		if err != nil {
			log.Fatalln(logger.ExternalError(err, msg))
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
			log.Fatalln(logger.ExternalError(err, msg))
		}

		// merge templates into the newly created project
		if orderConfig.Templates != nil {
			templates, err := cloudCfg.ValidTemplatesFrom(orderConfig.Templates)
			if err != nil {
				log.Fatalln(err)
			}
			for _, t := range templates {
				if err := template.MergeTemplate(t, project); err != nil {
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
		model, err := pom.GetModelFrom(project.Type.FilePath())
		if err != nil {
			log.Fatalln(err)
		}

		log.Info(logger.Info(fmt.Sprintf("formatting %s", project.Type.FilePath())))
		if err = maven.ChangeVersionToPropertyTagsOnModel(model); err != nil {
			log.Fatalln(err)
		}

		// upgrade all
		log.Info(logger.Info(fmt.Sprintf("upgrading all on %s", project.Type.FilePath())))
		if err = upgradeAll(model); err != nil {
			log.Fatalln(err)
		}

		// sorting and writing
		if err = project.SortAndWritePom(true); err != nil {
			log.Fatalln(err)
		}

		// git commit
		err = project.GitCommit("Clean up and upgrades")
		if err != nil {
			log.Fatalln(logger.ExternalError(err, msg))
		}
	},
}

var springInheritVersion = &cobra.Command{
	Use:   "inherit",
	Short: "Removes manual versions from spring dependencies",
	Long:  `Removes manual versions from spring dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		project, err := config.InitProjectFromDirectory(ctx.TargetDirectory)
		if err != nil {
			log.Fatalln(err)
		}

		if err = spring.CleanManualVersions(project.Type.Model()); err != nil {
			log.Fatalln(err)
		}

		if err = project.SortAndWritePom(ctx.Overwrite); err != nil {
			log.Fatalln(err)
		}
	},
}

var springDownloadCli = &cobra.Command{
	Use:   "download-cli",
	Short: "Downloads spring-cli",
	Long:  `Downloads spring-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := spring.CheckCli(localCfg); err != nil {
			log.Fatalln(err)
		}
	},
}

var springInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints info on spring boot and available dependencies",
	Long:  `Prints info on spring boot and available dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		root, err := spring.GetRoot()
		if err != nil {
			log.Fatalln(err)
		}
		log.Infof("Latest version of spring boot are: %s\n", root.BootVersion.Default)

		log.Infof(logger.Info(fmt.Sprintf("Valid dependencies: ")))
		for _, category := range root.Dependencies.Values {
			fmt.Println(logger.Info(fmt.Sprintf("%s", category.Name)))
			fmt.Printf("================================\n")
			for _, dep := range category.Values {
				fmt.Printf("[%s]\n    %s, (%s)\n", logger.Magenta(dep.Id), dep.Name, dep.Description)
			}
			fmt.Printf("\n")
		}

	},
}

var springManagedCmd = &cobra.Command{
	Use:   "managed",
	Short: "Prints info on spring-boot managed dependencies",
	Long:  `Prints info on spring-boot managed dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		deps, err := spring.GetDependencies()
		if err != nil {
			log.Fatalln(err)
		}

		log.Infof(logger.Info(fmt.Sprintf("Spring Boot managed dependencies:")))
		var organized = make(map[string][]pom.Dependency)
		for _, dep := range deps.Dependencies {
			mvnDep := pom.Dependency{
				GroupId:    dep.GroupId,
				ArtifactId: dep.ArtifactId,
			}
			organized[dep.GroupId] = append(organized[dep.GroupId], mvnDep)
		}

		for k, v := range organized {
			fmt.Println(logger.Info(fmt.Sprintf("GroupId: %s", k)))
			fmt.Printf("================================\n")
			for _, mvnDep := range v {
				fmt.Printf("  ArtifactId: %s\n", mvnDep.ArtifactId)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(springCmd)
	springCmd.AddCommand(springInitCmd)
	springCmd.AddCommand(springInfoCmd)
	springCmd.AddCommand(springManagedCmd)
	springCmd.AddCommand(springInheritVersion)
	springCmd.AddCommand(springDownloadCli)

	springCmd.PersistentFlags().StringVar(&ctx.TargetDirectory, "target", ".", "Optional target directory")
	springCmd.PersistentFlags().BoolVar(&ctx.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	springCmd.PersistentFlags().BoolVar(&ctx.Overwrite, "no-git", false, "Disables git for every step")
	springInitCmd.Flags().String("config-file", "", "Optional config file")
}
