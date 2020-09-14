package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/merge"
	"co-pilot/pkg/springio"
	"co-pilot/pkg/upgrade"
	"encoding/xml"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
	"os"
	"strings"
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

		// populate applicationName field in config
		if err := initConfig.FindApplicationName(targetDir); err != nil {
			log.Errorln(err)
		}

		// write co-pilot.json to target directory
		configFile := fmt.Sprintf("%s/co-pilot.json", targetDir)
		msg := logger.Info(fmt.Sprintf("writes co-pilot.json config file to %s", configFile))
		log.Info(msg)
		if err := initConfig.WriteConfig(configFile); err != nil {
			log.Fatalln(err)
		}

		// merge templates
		if initConfig.LocalDependencies != nil {
			for _, d := range initConfig.LocalDependencies {
				if err := merge.TemplateName(d, targetDir); err != nil {
					log.Fatalln(err)
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

var springDownloadCli = &cobra.Command{
	Use:   "download-cli",
	Short: "Downloads spring-cli",
	Long:  `Downloads spring-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := springio.CheckCli(); err != nil {
			log.Fatalln(err)
		}
	},
}

var springInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints info on spring boot and available dependencies",
	Long:  `Prints info on spring boot and available dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		root, err := springio.GetRoot()
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
		deps, err := springio.GetDependencies()
		if err != nil {
			log.Fatalln(err)
		}

		log.Infof(logger.Info(fmt.Sprintf("Spring Boot Managed Dependencies:")))
		var organized = make(map[string][]pom.Dependency)
		for _, dep := range deps.Dependencies {
			mvnDep := pom.Dependency{
				GroupId:    dep.GroupId,
				ArtifactId: dep.ArtifactId,
			}
			organized[dep.GroupId] = append(organized[dep.GroupId], mvnDep)
		}

		for k, v := range organized {
			fmt.Println(logger.Info(fmt.Sprintf("%s", k)))
			fmt.Printf("================================\n")
			for _, mvnDep := range v {
				b, _ := xml.MarshalIndent(mvnDep, "", "    ")
				for _, line := range strings.Split(string(b), "\n") {
					fmt.Println(line)
				}
				println("")
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
	springCmd.PersistentFlags().String("target", ".", "Optional target directory")
	springCmd.PersistentFlags().Bool("overwrite", true, "Overwrite pom.xml file")
	springInitCmd.Flags().String("config-file", "", "Optional config file")
}
