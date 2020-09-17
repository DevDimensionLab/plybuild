package spring

import (
	"co-pilot/pkg/config"
	"co-pilot/pkg/file"
	"co-pilot/pkg/http"
	"co-pilot/pkg/logger"
	"co-pilot/pkg/shell"
	"errors"
	"fmt"
	"os"
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

	springExec, err := file.FindFirst("bin/spring", targetDir)
	if err != nil {
		return "", err
	}

	return shell.Run(springExec, arg...)
}

func CheckCli() error {
	log.Infof("checking if Spring CLI is installed and for latest version")
	targetDir, err := binDir()
	if err != nil {
		return err
	}

	springExec, err := file.FindFirst("bin/spring", targetDir)
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

func InitFrom(config config.ProjectConfiguration, targetDir string) []string {
	var output []string

	output = append(output, "init")
	output = append(output, "-g="+config.GroupId)
	output = append(output, "-a="+config.ArtifactId)
	output = append(output, "--package="+config.Package)
	output = append(output, "--d="+strings.Join(config.Dependencies, ","))
	output = append(output, "-j=11")
	output = append(output, "--language=kotlin")
	output = append(output, "--description="+config.Description)
	output = append(output, "--name="+config.Name)
	output = append(output, targetDir) //output directory

	return output
}

func GetRoot() (IoRootResponse, error) {
	var deps IoRootResponse
	err := http.GetJson("http://start.spring.io", &deps)
	return deps, err
}

func GetInfo() (IoInfoResponse, error) {
	var deps IoInfoResponse
	err := http.GetJson("http://start.spring.io/actuator/info", &deps)
	return deps, err
}

func GetDependencies() (IoDependenciesResponse, error) {
	var deps IoDependenciesResponse
	err := http.GetJson("http://start.spring.io/dependencies", &deps)
	return deps, err
}

func Validate(config config.ProjectConfiguration) error {
	var invalidDependencies []string
	validDependencies, err := GetDependencies()
	if err != nil {
		return err
	}

	for _, userDefinedDependency := range config.Dependencies {
		valid := false
		for validDependency, _ := range validDependencies.Dependencies {
			if validDependency == userDefinedDependency {
				valid = true
			}
		}
		if !valid {
			invalidDependencies = append(invalidDependencies, userDefinedDependency)
		}
	}

	if len(invalidDependencies) > 0 {
		validKeys := make([]string, 0, len(validDependencies.Dependencies))
		for k, _ := range validDependencies.Dependencies {
			validKeys = append(validKeys, k)
		}
		return errors.New(fmt.Sprintf("%s not found in valid list of dependencies %s", invalidDependencies, validKeys))
	} else {
		return nil
	}
}
