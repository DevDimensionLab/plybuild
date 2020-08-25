package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func TouchConfigFile(configFile string) error {
	config := defaultConfig()

	d, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

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

func defaultConfig() GlobalConfiguration {
	return GlobalConfiguration{}
}
