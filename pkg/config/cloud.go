package config

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/git"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
)

var globalConfigDir = "cloud-config"

func GlobalConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", home, localConfigDir, globalConfigDir), nil
}

func Clone() error {
	c, err := GetLocalConfig()
	if err != nil {
		log.Fatalln(err)
	}

	target, err := GlobalConfigDir()
	if err != nil {
		log.Fatalln(err)
	}

	if file.Exists(fmt.Sprintf("%s/.git", target)) {
		log.Infof("pulling cloud config on %s", target)
		err = git.Pull(target)
		if err != nil {
			return err
		}
	} else {
		log.Infof("cloning %s to %s", c.CloudConfig.Git.Url, target)
		err = git.Clone(c.CloudConfig.Git.Url, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDeprecated() (CloudDeprecated, error) {
	var deprecated CloudDeprecated

	path, err := GetCloudConfigFilePath("deprecated.json")
	if err != nil {
		return deprecated, err
	}

	err = file.ReadJson(path, &deprecated)
	if err != nil {
		return deprecated, err
	}

	return deprecated, nil
}

func GetCloudConfigFilePath(fileName string) (string, error) {
	err := Clone()
	if err != nil {
		return "", err
	}

	cloudConfigDir, err := GlobalConfigDir()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", cloudConfigDir, fileName)
	if !file.Exists(path) {
		return "", errors.New(fmt.Sprintf("could not find %s in cloud config", fileName))
	}

	return path, nil
}
