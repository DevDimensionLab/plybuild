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

func LocalConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", home, localConfigDir), nil
}

func LocalConfigFilePath() (string, error) {
	configDir, err := LocalConfigDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", configDir, localConfigFileName), nil
}

func CheckOrCreateConfigDir() error {
	dir, err := LocalConfigDir()
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

func TouchLocalConfigFile() error {
	err := CheckOrCreateConfigDir()
	if err != nil {
		return err
	}

	configFile, err := LocalConfigFilePath()
	if err != nil {
		return err
	}

	config := GlobalConfiguration{}
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

func GetLocalConfig() (GlobalConfiguration, error) {
	config := GlobalConfiguration{}
	localConfigFile, err := LocalConfigFilePath()
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

func PrintLocalConfig(config GlobalConfiguration) error {
	b, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	log.Infof("\n%s\n", b)

	return nil
}

func LocalConfigExists() bool {
	localConfigFile, err := LocalConfigFilePath()
	if err != nil {
		return false
	}

	return file.Exists(localConfigFile)
}
