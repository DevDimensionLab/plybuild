package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"spring-boot-co-pilot/pkg/file"
	"spring-boot-co-pilot/pkg/spring"
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
	Short: "Spring init",
	Long:  `Spring init`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonConfigFile, _ := cmd.Flags().GetString("config-file")
		var config = spring.InitConfiguration{}

		_ = os.RemoveAll("webservice")

		if jsonConfigFile != "" {
			err := file.ReadJson(jsonConfigFile, &config)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			err = spring.Validate(config)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		} else {
			config = spring.DefaultConfiguration()
		}

		springExec, err := file.Find("bin/spring", "./target")
		err = spring.SpringBootCLI(springExec, spring.InitFrom(config)...)

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
		root, err := spring.GetRoot()
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Latest version of spring boot are: %s\n", root.BootVersion.Default)
	},
}

func init() {
	RootCmd.AddCommand(springCmd)
	springCmd.AddCommand(springInitCmd)
	springCmd.AddCommand(springStatusCmd)
	springInitCmd.Flags().String("config-file", "", "Optional config file")
}
