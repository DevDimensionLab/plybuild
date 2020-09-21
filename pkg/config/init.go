package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
)

func NewLocalConfig(localConfigDir string) (cfg LocalConfig, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return
	}

	cfg.impl.Path = fmt.Sprintf("%s/%s", home, localConfigDir)
	return
}

func NewGitCloudConfig(localConfigDir string, cloudConfigDirName string) (cfg GitCloudConfig, err error) {
	home, err := homedir.Dir()
	if err != nil {
		return cfg, err
	}

	cfg.Impl.Path = fmt.Sprintf("%s/%s/%s", home, localConfigDir, cloudConfigDirName)
	return
}
