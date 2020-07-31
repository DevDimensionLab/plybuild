package cmd

import (
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


func init() {
	RootCmd.AddCommand(springCmd)
	springCmd.AddCommand(springInitCmd)
	springInitCmd.Flags().String("config-file", "", "Optional config file")
}
