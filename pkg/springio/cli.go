package springio

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/shell"
	"fmt"
	"os"
	"os/exec"
)

var springBootDownloadUrl = "https://repo.spring.io/release/org/springframework/boot/spring-boot-cli/[RELEASE]/spring-boot-cli-[RELEASE]-bin.zip"

func BinDir() (string, error) {
	binDir := "spring-cli"
	configDir, err := config.LocalConfigDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", configDir, binDir), nil
}

func DownloadCli() error {
	targetDir, err := BinDir()
	if err != nil {
		return err
	}
	springBootCliZip := fmt.Sprintf("%s/spring-boot-cli.zip", targetDir)

	if err := os.RemoveAll(targetDir); err != nil {
		return err
	}
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return err
	}
	if err := shell.Wget(springBootDownloadUrl, springBootCliZip); err != nil {
		return err
	}
	if err := shell.Unzip(springBootCliZip, targetDir); err != nil {
		return err
	}
	return nil
}

func RunCli(arg ...string) error {
	targetDir, err := BinDir()
	if err != nil {
		return err
	}

	springExec, err := file.Find("bin/spring", targetDir)
	if err != nil {
		return err
	}

	cmd := exec.Command(springExec, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
