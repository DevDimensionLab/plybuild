package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	_ = Wget(springBootDownloadUrl, springBootCliZip)
	_ = Unzip(springBootCliZip, targetDir)

	springExec, err := FindFile("bin/spring", targetDir)
	_ = SpringBootCLI(springExec)

	if err != nil {
		log.Println(err)
	}

}

func SpringBootCLI(springExec string) error {
	cmd := exec.Command(springExec)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func FindFile(fileSuffix string, dir string) (result string, err error) {
	err = filepath.Walk(dir,
		func(path string, fi os.FileInfo, errIn error) error {
			if strings.HasSuffix( path, fileSuffix) {
				result = path
				return io.EOF
			}
			return nil
		})

	if err == io.EOF {
		err = nil
	}
	return
}

func Wget(url, filepath string) error {
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Unzip(file string, outputDir string) error {
	cmd := exec.Command("unzip", file,"-d", outputDir )
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}