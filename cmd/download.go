package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"spring-boot-co-pilot/pkg/shell"
)

var springBootDownloadUrl = "https://repo.spring.io/release/org/springframework/boot/spring-boot-cli/[RELEASE]/spring-boot-cli-[RELEASE]-bin.zip"

var downloadCmd = &cobra.Command{
	Use:   "download [OPTIONS]",
	Short: "Downloads ...",
	Long:  `Downloads ...`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

var downloadCliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Downloads spring-cli",
	Long:  `Downloads spring-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		download()
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	downloadCmd.AddCommand(downloadCliCmd)
}

func download() {
	targetDir := "./target"
	_ = os.RemoveAll(targetDir)
	_ = os.MkdirAll(targetDir, os.ModePerm)
	springBootCliZip := "./target/spring-boot-cli.zip"
	_ = shell.Wget(springBootDownloadUrl, springBootCliZip)
	_ = shell.Unzip(springBootCliZip, targetDir)
}
