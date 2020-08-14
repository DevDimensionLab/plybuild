package cmd

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/springio"
	log "github.com/sirupsen/logrus"
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

var springInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "downloads and installs spring boot with default or provided settings",
	Long:  `downloads and installs spring boot with default or provided settings`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonConfigFile, _ := cmd.Flags().GetString("config-file")
		var initConfig = config.InitConfiguration{}

		_ = os.RemoveAll("webservice")

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
		err = springio.CLI(springExec, springio.InitFrom(initConfig)...)

		if err != nil {
			log.Println(err)
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
	springCmd.AddCommand(springInstallCmd)
	springCmd.AddCommand(springStatusCmd)
	springInstallCmd.Flags().String("config-file", "", "Optional config file")
}
