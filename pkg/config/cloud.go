package config

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/git"
	"co-pilot/pkg/logger"
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
		msg := logger.Info(fmt.Sprintf("pulling cloud config on %s", target))
		log.Info(msg)
		out, err := git.Pull(target)
		if err != nil {
			return errors.New(fmt.Sprintf("pulling cloud config failed:\n%s, %v", out, err))
		}
	} else {
		msg := logger.Info(fmt.Sprintf("cloning %s to %s", c.CloudConfig.Git.Url, target))
		log.Info(msg)
		out, err := git.Clone(c.CloudConfig.Git.Url, target)
		if err != nil {
			return errors.New(fmt.Sprintf("ploning cloud config failed:\n%s, %v", out, err))
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
