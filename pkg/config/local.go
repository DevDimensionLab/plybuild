package config

import (
	"co-pilot/pkg/file"
	"co-pilot/pkg/logger"
	"fmt"
	"github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var log = logger.Context()

var localConfigFileName = "local-config.yaml"
var localConfigDir = ".co-pilot"

type LocalConfig struct {
}

type LocalConfigFile interface {
	DirPath() (string, error)
	FilePath() (string, error)
	CheckOrCreateConfigDir() error
	TouchFile() error
	Config() (LocalConfiguration, error)
	Print() error
	Exists() bool
}

func InitLocalConfig() (localConfig LocalConfig) {
	return
}

func (localConfig LocalConfig) DirPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", home, localConfigDir), nil
}

func (localConfig LocalConfig) FilePath() (string, error) {
	configDir, err := localConfig.DirPath()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", configDir, localConfigFileName), nil
}

func (localConfig LocalConfig) CheckOrCreateConfigDir() error {
	dir, err := localConfig.DirPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (localConfig LocalConfig) TouchFile() error {
	err := localConfig.CheckOrCreateConfigDir()
	if err != nil {
		return err
	}

	configFile, err := localConfig.FilePath()
	if err != nil {
		return err
	}

	config := LocalConfiguration{}
	d, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	log.Infof("creating new config file %s", configFile)

	f, err := os.Create(configFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFile, d, 0644)
	if err != nil {
		return err
	}

	return f.Close()
}

func (localConfig LocalConfig) Config() (LocalConfiguration, error) {
	config := LocalConfiguration{}
	localConfigFile, err := localConfig.FilePath()
	if err != nil {
		return config, err
	}

	b, err := file.Open(localConfigFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (localConfig LocalConfig) Print() error {
	c, err := localConfig.Config()
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	log.Infof("\n%s\n", b)

	return nil
}

func (localConfig LocalConfig) Exists() bool {
	localConfigFile, err := localConfig.FilePath()
	if err != nil {
		return false
	}

	return file.Exists(localConfigFile)
}
