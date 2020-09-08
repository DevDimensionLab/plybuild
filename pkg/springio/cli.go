package springio

import (
	"bytes"
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/shell"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var springBootDownloadUrl = "https://repo.spring.io/release/org/springframework/boot/spring-boot-cli/[RELEASE]/spring-boot-cli-[RELEASE]-bin.zip"
var log = logger.Context()

func binDir() (string, error) {
	binDir := "spring-cli"
	configDir, err := config.LocalConfigDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", configDir, binDir), nil
}

func RunCli(arg ...string) (string, error) {
	targetDir, err := binDir()
	if err != nil {
		return "", err
	}

	springExec, err := file.Find("bin/spring", targetDir)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(springExec, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

func CheckCli() error {
	log.Infof("checking if Spring CLI is installed and for latest version")
	targetDir, err := binDir()
	if err != nil {
		return err
	}

	springExec, err := file.Find("bin/spring", targetDir)
	if err != nil {
		return err
	}

	if springExec != "" {
		return upgrade()
	} else {
		return install()
	}
}

func install() error {
	log.Infof("installing Spring CLI into ~/.co-pilot directory")
	targetDir, err := binDir()
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

	_, err = shell.Wget(springBootDownloadUrl, springBootCliZip)
	if err != nil {
		return err
	}

	_, err = shell.Unzip(springBootCliZip, targetDir)
	if err != nil {
		return err
	}
	return nil
}

func upgrade() error {
	versionStr, err := RunCli("version")
	if err != nil {
		return err
	}
	versionStr = strings.TrimSpace(versionStr)

	response, err := GetRoot()
	if err != nil {
		return err
	}

	if strings.HasSuffix(versionStr, response.BootVersion.Default) {
		log.Infof("spring CLI is the latest version %s", response.BootVersion.Default)
		return nil
	} else {
		log.Infof("upgrading Spring CLI from %s, to the latest version %s", versionStr, response.BootVersion.Default)
		return install()
	}
}
