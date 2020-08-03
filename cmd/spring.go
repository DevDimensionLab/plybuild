package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"os"
	"spring-boot-co-pilot/pkg/spring"
	"spring-boot-co-pilot/pkg/util"
)

var springCmd = &cobra.Command{
	Use:   "spring [OPTIONS]",
	Short: "Spring ...",
	Long:  `Spring ...`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var springInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Spring init",
	Long:  `Spring init`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonConfigFile, _ := cmd.Flags().GetString("config-file")
		var config = spring.InitConfiguration{}

		_ = os.RemoveAll("webservice")

		if jsonConfigFile != "" {
			err := util.ReadJson(jsonConfigFile, &config)
			if err != nil {
				log.Println(err)
				config = spring.DefaultConfiguration()
			}
		} else {
			config = spring.DefaultConfiguration()
		}

		springExec, err := util.FindFile("bin/spring", "./target")
		err = spring.SpringBootCLI(springExec, spring.InitFrom(config)...)

		if err != nil {
			log.Println(err)
		}
	},
}

var springRootCmd = &cobra.Command{
	Use:   "root",
	Short: "Spring root",
	Long:  `Spring root`,
	Run: func(cmd *cobra.Command, args []string) {
		root, err := spring.GetRoot()
		if err != nil {
			log.Println(err)
		}
		json, _ := json.MarshalIndent(root, "", "    ")
		println(string(json))
	},
}

var springInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Spring metadata info",
	Long:  `Spring metadata info`,
	Run: func(cmd *cobra.Command, args []string) {
		root, err := spring.GetInfo()
		if err != nil {
			log.Println(err)
		}
		json, _ := json.MarshalIndent(root, "", "    ")
		println(string(json))
	},
}


var springDependenciesCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Spring dependencies",
	Long:  `Spring dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		deps, err := spring.GetDependencies()
		if err != nil {
			log.Println(err)
		}
		json, _ := json.MarshalIndent(deps, "", "    ")
		println(string(json))
	},
}

func init() {
	RootCmd.AddCommand(springCmd)
	springCmd.AddCommand(springInitCmd)
	springCmd.AddCommand(springRootCmd)
	springCmd.AddCommand(springInfoCmd)
	springCmd.AddCommand(springDependenciesCmd)
	springInitCmd.Flags().String("config-file", "", "Optional config file")
}
