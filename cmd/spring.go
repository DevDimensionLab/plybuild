package cmd

import (
	"co-pilot/pkg/clean"
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
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

		targetDir := "webservice"
		_ = os.RemoveAll(targetDir)

		if jsonConfigFile != "" {
			err := file.ReadJson(jsonConfigFile, &initConfig)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			err = springio.Validate(initConfig)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		} else {
			initConfig = config.DefaultConfiguration()
		}

		springExec, err := file.Find("bin/spring", "./target")
		err = springio.CLI(springExec, springio.InitFrom(initConfig, targetDir)...)
		if err != nil {
			log.Println(err)
		}

		configFile := fmt.Sprintf("%s/co-pilot.json", targetDir)
		log.Infof("writes co-pilot.json config file to %s", configFile)
		if err = config.WriteConfig(initConfig, configFile); err != nil {
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
