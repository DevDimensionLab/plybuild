package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/merge"
	"co-pilot/pkg/springio"
	"co-pilot/pkg/upgrade"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
	"os"
)

var springCmd = &cobra.Command{
	Use:   "spring [OPTIONS]",
	Short: "Spring boot tools",
	Long:  `Spring boot tools`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var springInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Downloads and installs spring boot with default or provided settings",
	Long:  `Downloads and installs spring boot with default or provided settings`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonConfigFile, _ := cmd.Flags().GetString("config-file")
		var initConfig = config.ProjectConfiguration{}

		targetDir, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}
		if targetDir == "" {
			targetDir = "webservice"
		}

		_ = os.RemoveAll(targetDir)

		// sync cloud config
		if err := config.Clone(); err != nil {
			log.Fatalln(err)
		}

		// fetch user input config
		if jsonConfigFile != "" {
			err := file.ReadJson(jsonConfigFile, &initConfig)
			if err != nil {
				log.Fatalln(err)
			}
			err = springio.Validate(initConfig)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			initConfig = config.DefaultConfiguration()
		}

		// download cli
		if err := springio.CheckCli(); err != nil {
			log.Fatalln(err)
		}

		// execute cli with config
		_, err = springio.RunCli(springio.InitFrom(initConfig, targetDir)...)
		if err != nil {
			log.Fatalln(err)
		}

		// write co-pilot.json to target directory
		configFile := fmt.Sprintf("%s/co-pilot.json", targetDir)
		msg := logger.Info(fmt.Sprintf("writes co-pilot.json config file to %s", configFile))
		log.Info(msg)
		if err := config.WriteConfig(initConfig, configFile); err != nil {
			log.Fatalln(err)
		}

		// merge templates
		if initConfig.LocalDependencies != nil {
			for _, d := range initConfig.LocalDependencies {
				if err := merge.TemplateName(d, targetDir); err != nil {
					log.Errorln(err)
				}
			}
		}

		// format version
		pomFile := targetDir + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		log.Info(logger.Info(fmt.Sprintf("formatting %s", pomFile)))
		if err = clean.VersionToPropertyTags(model); err != nil {
			log.Fatalln(err)
		}

		// upgrade all
		log.Info(logger.Info(fmt.Sprintf("upgrading %s", pomFile)))
		upgrade.All(model)

		// sorting and writing
		log.Info(logger.Info(fmt.Sprintf("Sorting and rewriting %s", pomFile)))
		if err = upgrade.SortAndWrite(model, pomFile); err != nil {
			log.Fatalln(err)
		}
	},
}

var springManagedCmd = &cobra.Command{
	Use:   "managed",
	Short: "Prints spring-boot managed dependencies",
	Long:  `Prints spring-boot managed dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		deps, err := springio.GetDependencies()
		if err != nil {
			log.Fatalln(err)
		}

		log.Infof("Spring Boot Managed Dependencies:")
		for _, dep := range deps.Dependencies {
			fmt.Printf("\t%s:%s [%s]\n", dep.GroupId, dep.ArtifactId, dep.Version)
		}
	},
}

var springInheritVersion = &cobra.Command{
	Use:   "inherit",
	Short: "Removes manual versions from spring dependencies",
	Long:  `Removes manual versions from spring dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		targetDirectory, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatalln(err)
		}
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			log.Fatalln(err)
		}

		pomFile := targetDirectory + "/pom.xml"
		model, err := pom.GetModelFrom(pomFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err = clean.SpringManualVersion(model); err != nil {
			log.Fatalln(err)
		}

		var writeToFile = pomFile
		if !overwrite {
			writeToFile = targetDirectory + "/pom.xml.new"
		}
		if err = upgrade.SortAndWrite(model, writeToFile); err != nil {
			log.Fatalln(err)
		}
	},
}

var springStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Spring status",
	Long:  `Spring status`,
	Run: func(cmd *cobra.Command, args []string) {
		root, err := springio.GetRoot()
		if err != nil {
			log.Fatalln(err)
		}
		log.Infof("Latest version of spring boot are: %s\n", root.BootVersion.Default)
	},
}

func init() {
	RootCmd.AddCommand(springCmd)
	springCmd.AddCommand(springInitCmd)
	springCmd.AddCommand(springStatusCmd)
	springCmd.AddCommand(springManagedCmd)
	springCmd.AddCommand(springInheritVersion)
	springCmd.PersistentFlags().String("target", ".", "Optional target directory")
	springCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
	springInitCmd.Flags().String("config-file", "", "Optional config file")
}
